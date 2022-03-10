package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/icholy/digest"
)

const (
	BaseURL = "https://cloud.mongodb.com/api/atlas/v1.0"
)

type (
	Client interface {
		IndexOrgGroups(orgID string) ([]*Group, error)
		RawIndexOrgGroups(orgID string) ([]json.RawMessage, error)
		IndexInvoices(orgID string) ([]*Invoice, error)
		RawIndexInvoices(orgID string) ([]json.RawMessage, error)
		ShowInvoice(orgID, invoiceID string) (*Invoice, error)
		RawShowInvoice(orgID, invoiceID string) (json.RawMessage, error)
	}

	client struct {
		baseURL string
		client  *http.Client
	}

	paginatedResponse struct {
		Results    []json.RawMessage `json:"results"`
		TotalCount uint              `json:"totalCount"`
	}

	Group struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Invoice struct {
		ID        string      `json:"id"`
		StartDate time.Time   `json:"startDate"`
		LineItems []*LineItem `json:"lineItems"`
	}

	LineItem struct {
		StartDate        time.Time `json:"startDate"`
		EndDate          time.Time `json:"endDate"`
		GroupID          string    `json:"groupId"`
		ClusterName      string    `json:"clusterName"`
		Note             string    `json:"note"`
		SKU              string    `json:"sku"`
		DiscountCents    float64   `json:"discountCents"`
		PrecentDiscount  float64   `json:"percentDiscount"`
		Quantity         float64   `json:"quantity"`
		TotalPriceCents  float64   `json:"totalPriceCents"`
		Unit             string    `json:"unit"`
		UnitPriceDollars float64   `json:"unitPriceDollars"`
	}
)

func NewClient(baseURL, publicKey, privateKey string) Client {
	return &client{
		baseURL: baseURL,
		client: &http.Client{
			Transport: &digest.Transport{
				Username: publicKey,
				Password: privateKey,
			},
		},
	}
}

func (c *client) IndexOrgGroups(orgID string) (groups []*Group, err error) {
	rawGroups, err := c.RawIndexOrgGroups(orgID)
	if err != nil {
		return
	}
	for _, r := range rawGroups {
		g := &Group{}
		err = json.Unmarshal(r, g)
		if err != nil {
			return
		}
		groups = append(groups, g)
	}
	return
}

func (c *client) RawIndexOrgGroups(orgID string) (groups []json.RawMessage, err error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/orgs/"+orgID+"/groups", nil)
	if err != nil {
		return
	}
	groups, err = c.paginate(req)
	return
}

func (c *client) IndexInvoices(orgID string) (invoices []*Invoice, err error) {
	rawInvoices, err := c.RawIndexInvoices(orgID)
	if err != nil {
		return
	}
	for _, r := range rawInvoices {
		i := &Invoice{}
		err = json.Unmarshal(r, i)
		if err != nil {
			return
		}
		invoices = append(invoices, i)
	}
	return
}

func (c *client) RawIndexInvoices(orgID string) (invoices []json.RawMessage, err error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/orgs/"+orgID+"/invoices/", nil)
	if err != nil {
		return
	}
	invoices, err = c.paginate(req)
	return
}

func (c *client) ShowInvoice(orgID, invoiceID string) (invoice *Invoice, err error) {
	rawInvoice, err := c.RawShowInvoice(orgID, invoiceID)
	if err != nil {
		return
	}
	invoice = &Invoice{}
	err = json.Unmarshal(rawInvoice, invoice)
	return
}

func (c *client) RawShowInvoice(orgID, invoiceID string) (invoice json.RawMessage, err error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/orgs/"+orgID+"/invoices/"+invoiceID, nil)
	if err != nil {
		return
	}
	invoice, err = c.do(req)
	return
}

func (c *client) do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		m := "empty body"
		if err != nil {
			m = fmt.Sprintf("failed to read body: %v", err)
		} else if len(b) > 0 {
			m = string(b)
		}
		return nil, fmt.Errorf("unexpected response: %v: %v", resp.Status, m)
	}
	return io.ReadAll(resp.Body)
}

func (c *client) paginate(req *http.Request) ([]json.RawMessage, error) {
	var (
		results []json.RawMessage
		q       = req.URL.Query()
	)
	q.Set("itemsPerPage", "500")
	req.URL.RawQuery = q.Encode()
	for page := 1; ; page++ {
		q = req.URL.Query()
		if page > 1 {
			q.Set("pageNum", fmt.Sprint(page))
		}
		req.URL.RawQuery = q.Encode()
		b, err := c.do(req)
		if err != nil {
			return nil, err
		}
		r := &paginatedResponse{}
		err = json.Unmarshal(b, r)
		if err != nil {
			return nil, err
		}
		results = append(results, r.Results...)
		if len(results) >= int(r.TotalCount) {
			return results, nil
		}
	}
}
