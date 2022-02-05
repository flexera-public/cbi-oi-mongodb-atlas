package mongodbatlascbi

import (
	"time"

	"github.com/jszwec/csvutil"

	"github.com/flexera/cbi-oi-mongodb-atlas/services/mongodbatlas"
)

const (
	YearMonthLayout = "200601"
)

type (
	Row struct {
		CloudVendorAccountID   string
		CloudVendorAccountName string
		Category               string
		InstanceType           string
		LineItemType           string
		Region                 string
		ResourceGroup          string
		ResourceType           string
		ResourceID             string
		Service                string
		UsageType              string
		Tags                   string
		UsageAmount            float64
		UsageUnit              string
		Cost                   float64
		CurrencyCode           string
		UsageStartTime         time.Time
		InvoiceYearMonth       YearMonth
		InvoiceID              string
	}

	YearMonth struct {
		time.Time
	}
)

func Convert(groups []*mongodbatlas.Group, invoice *mongodbatlas.Invoice, w csvutil.Writer) error {
	groupNamesByID := make(map[string]string, len(groups))
	for _, group := range groups {
		groupNamesByID[group.ID] = group.Name
	}
	e := csvutil.NewEncoder(w)
	err := e.EncodeHeader(&Row{})
	if err != nil {
		return err
	}
	for _, lineItem := range invoice.LineItems {
		err = e.Encode(&Row{
			CloudVendorAccountID:   lineItem.GroupID,
			CloudVendorAccountName: groupNamesByID[lineItem.GroupID],
			ResourceID:             lineItem.ClusterName,
			UsageType:              lineItem.SKU,
			UsageAmount:            lineItem.Quantity,
			UsageUnit:              lineItem.Unit,
			Cost:                   lineItem.TotalPriceCents / 100,
			CurrencyCode:           "USD",
			UsageStartTime:         lineItem.StartDate,
			InvoiceYearMonth:       YearMonth{invoice.StartDate},
			InvoiceID:              invoice.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t YearMonth) MarshalCSV() ([]byte, error) {
	b := make([]byte, 0, len(YearMonthLayout))
	return t.AppendFormat(b, YearMonthLayout), nil
}
