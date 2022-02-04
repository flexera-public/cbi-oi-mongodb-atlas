package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/flexera/cbi-oi-mongodb-atlas/services/mongodbatlas"
)

type Download struct {
	Groups   bool        `help:"Download a JSON list of all projects (groups) in the MongoDB Atlas Org." required:"" xor:"cmd" short:"g"`
	Invoices bool        `help:"Download a JSON list of all invoices in the MongoDB Atlas Org." required:"" xor:"cmd" short:"i"`
	Last     int         `help:"Download the last INT JSON invoice(s) from the MongoDB Atlas Org." required:"" xor:"cmd" short:"l"`
	Month    []time.Time `help:"Download the JSON invoice for each MONTH from the MongoDB Atlas Org." required:"" xor:"cmd" format:"2006-01" short:"m"`

	atlasClient mongodbatlas.Client
}

func (download *Download) Run(cli *CLI) error {
	download.atlasClient = cli.AtlasClient()
	if download.Groups {
		err := download.groups(cli.AtlasOrgID)
		if err != nil {
			return err
		}
	} else if download.Invoices {
		err := download.invoices(cli.AtlasOrgID)
		if err != nil {
			return err
		}
	} else {
		invoices, err := download.atlasClient.IndexInvoices(cli.AtlasOrgID)
		if err != nil {
			return err
		}
		if len(download.Month) > 0 {
			invoicesByStartDate := make(map[time.Time]*mongodbatlas.Invoice)
			for _, invoice := range invoices {
				invoicesByStartDate[invoice.StartDate] = invoice
			}
			invoices = make([]*mongodbatlas.Invoice, 0, len(download.Month))
			for _, month := range download.Month {
				invoice, ok := invoicesByStartDate[month]
				if !ok {
					return fmt.Errorf("missing invoice for month: %v", month.Format(MonthLayout))
				}
				invoices = append(invoices, invoice)
			}
		} else if download.Last <= 0 {
			return fmt.Errorf("--last: expected an int greater than 0 but got %v", download.Last)
		} else {
			sort.SliceStable(invoices, func(i, j int) bool {
				return invoices[i].StartDate.Before(invoices[j].StartDate)
			})
			if len(invoices) > int(download.Last) {
				invoices = invoices[len(invoices)-download.Last:]
			}
		}
		for _, invoice := range invoices {
			err = download.invoice(cli.AtlasOrgID, invoice)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (download *Download) groups(orgID string) error {
	fmt.Printf("Downloading %v\n", GroupsJSON)
	b, err := download.atlasClient.RawIndexOrgGroups(orgID)
	if err != nil {
		return err
	}
	f, err := os.Create(GroupsJSON)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	return e.Encode(b)
}

func (download *Download) invoices(orgID string) error {
	fmt.Printf("Downloading %v\n", InvoicesJSON)
	b, err := download.atlasClient.RawIndexInvoices(orgID)
	if err != nil {
		return err
	}
	f, err := os.Create(InvoicesJSON)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	return e.Encode(b)
}

func (download *Download) invoice(orgID string, invoice *mongodbatlas.Invoice) error {
	n := fmt.Sprintf(InvoiceJSONFmt, invoice.StartDate.Format(MonthLayout))
	fmt.Printf("Downloading %v\n", n)
	b, err := download.atlasClient.RawShowInvoice(orgID, invoice.ID)
	if err != nil {
		return err
	}
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	return e.Encode(b)
}
