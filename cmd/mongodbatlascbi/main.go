package main

import (
	"net/url"

	"github.com/flexera/cbi-oi-mongodb-atlas/services/mongodbatlas"

	"github.com/alecthomas/kong"
)

const (
	MonthLayout    = "2006-01"
	GroupsJSON     = "groups.json"
	InvoicesJSON   = "invoices.json"
	InvoiceJSONFmt = "invoice-%v.json"
	BillCSVFmt     = "bill-%v.csv"
)

type CLI struct {
	Debug           bool     `help:"Enable debug mode." short:"d"`
	AtlasURL        *url.URL `help:"Set the MongoDB Atlas API base URL." default:"${atlasURL}" short:"u"`
	AtlasOrgID      string   `help:"Set the MongoDB Atlas Org ID." required:"" env:"MONGODB_ATLAS_ORG_ID" short:"o"`
	AtlasPublicKey  string   `help:"Set the MongoDB Atlas API public key." required:"" env:"MONGODB_ATLAS_PUBLIC_KEY" short:"p"`
	AtlasPrivateKey string   `help:"Set the MongoDB Atlas API private key." required:"" env:"MONGODB_ATLAS_PRIVATE_KEY" short:"P"`

	Download Download `cmd:"" help:"Download JSON billing files from the MongoDB Atlas API."`
	Convert  Convert  `cmd:"" help:"Convert JSON billing files for MongoDB Atlas to the CBI CSV format."`

	atlasClient mongodbatlas.Client
}

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli, kong.Vars{
		"atlasURL":    mongodbatlas.BaseURL,
		"monthLayout": MonthLayout,
	})
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func (cli *CLI) AtlasClient() mongodbatlas.Client {
	if cli.atlasClient != nil {
		return cli.atlasClient
	}
	cli.atlasClient = mongodbatlas.NewClient(cli.AtlasURL.String(), cli.AtlasPublicKey, cli.AtlasPrivateKey)
	return cli.atlasClient
}
