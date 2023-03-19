module github.com/samherrmann/shopctl

go 1.20

replace github.com/bold-commerce/go-shopify/v3 => github.com/samherrmann/go-shopify/v3 v3.11.2

require (
	github.com/bold-commerce/go-shopify/v3 v3.11.0
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	github.com/spf13/cobra v1.2.1
)

require (
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
