package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mongodbatlascbi "github.com/flexera/cbi-oi-mongodb-atlas"
	"github.com/flexera/cbi-oi-mongodb-atlas/services/mongodbatlas"
)

type Convert struct {
	Month []time.Time `help:"Convert the JSON invoice for each MONTH to the CBI CSV format." required:"" format:"2006-01" short:"m"`

	atlasClient mongodbatlas.Client
}

func (convert *Convert) Run(cli *CLI) error {
	convert.atlasClient = cli.AtlasClient()
	groups, err := convert.groups()
	if err != nil {
		return err
	}
	for _, month := range convert.Month {
		invoice, err := convert.invoice(month)
		if err != nil {
			return err
		}
		err = convert.convert(groups, invoice)
		if err != nil {
			return err
		}
	}
	return nil
}

func (convert *Convert) groups() (groups []*mongodbatlas.Group, err error) {
	fmt.Printf("Loading %v\n", GroupsJSON)
	f, err := os.Open(GroupsJSON)
	if err != nil {
		return
	}
	defer f.Close()
	d := json.NewDecoder(f)
	err = d.Decode(&groups)
	return
}

func (convert *Convert) invoice(month time.Time) (invoice *mongodbatlas.Invoice, err error) {
	n := fmt.Sprintf(InvoiceJSONFmt, month.Format(MonthLayout))
	fmt.Printf("Loading %v\n", n)
	f, err := os.Open(n)
	if err != nil {
		return
	}
	defer f.Close()
	d := json.NewDecoder(f)
	invoice = &mongodbatlas.Invoice{}
	err = d.Decode(invoice)
	return
}

func (convert *Convert) convert(groups []*mongodbatlas.Group, invoice *mongodbatlas.Invoice) error {
	n := fmt.Sprintf(BillCSVFmt, invoice.StartDate.Format(MonthLayout))
	fmt.Printf("Writing %v\n", n)
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	return mongodbatlascbi.Convert(groups, invoice, w)
}
