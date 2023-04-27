package csv

import (
	"reflect"
	"testing"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/shopspring/decimal"
)

func Test_groupVariants(t *testing.T) {
	tests := []struct {
		name    string
		rows    [][]string
		want    []goshopify.Product
		wantErr bool
	}{
		{
			name: "no header",
			rows: [][]string{},
			want: []goshopify.Product{},
		},
		{
			name: "no variants",
			rows: [][]string{{}},
			want: []goshopify.Product{},
		},
		{
			name: "no title column",
			rows: [][]string{
				{"Foo"},
				{"abc"},
			},
			wantErr: true,
		},
		{
			name: "variant with empty title",
			rows: [][]string{
				{keyTitle},
				{""},
			},
			wantErr: true,
		},
		{
			name: "title column header but variant with no fields",
			rows: [][]string{
				{keyTitle},
				{},
			},
			wantErr: true,
		},
		{
			name: "variant with title",
			rows: [][]string{
				{keyTitle},
				{"foo"},
			},
			want: func() []goshopify.Product {
				p := goshopify.Product{Title: "foo"}
				p.Variants = append(p.Variants, goshopify.Variant{})
				return []goshopify.Product{p}
			}(),
		},
		{
			name: "multiple variants belonging to same product",
			rows: [][]string{
				{keyTitle},
				{"foo"},
				{"foo"},
			},
			want: func() []goshopify.Product {
				p := goshopify.Product{Title: "foo"}
				p.Variants = append(p.Variants, goshopify.Variant{}, goshopify.Variant{})
				return []goshopify.Product{p}
			}(),
		},
		{
			name: "multiple variants belonging to different products",
			rows: [][]string{
				{keyTitle},
				{"foo"},
				{"foo"},
				{"bar"},
				{"bar"},
			},
			want: func() []goshopify.Product {
				p1 := goshopify.Product{Title: "foo"}
				p2 := goshopify.Product{Title: "bar"}
				p1.Variants = append(p1.Variants, goshopify.Variant{}, goshopify.Variant{})
				p2.Variants = append(p2.Variants, goshopify.Variant{}, goshopify.Variant{})
				return []goshopify.Product{p1, p2}
			}(),
		},
		{
			name: "all fields",
			rows: [][]string{
				{
					keyProductID,
					keyVariantID,
					keySKU,
					keyBarcode,
					keyTitle,
					keyVendor,
					keyProductType,
					keyWeight,
					keyWeightUnit,
					keyPrice,
					keyOption1Name,
					keyOption1Value,
					keyOption2Name,
					keyOption2Value,
					keyOption3Name,
					keyOption3Value,
					"product.metafields.foo.bar",
					"variant.metafields.foo.bar",
				},
				{
					"123",
					"456",
					"mySku",
					"myBarcode",
					"myTitle",
					"myVendor",
					"myProductType",
					"123.456",
					"myWeightUnit",
					"7.89",
					"myOption1Name",
					"myOption1Value",
					"myOption2Name",
					"myOption2Value",
					"myOption3Name",
					"myOption3Value",
					"myProductMetafield",
					"myVariantMetafield",
				},
			},
			want: func() []goshopify.Product {

				var productID int64 = 123
				var variantID int64 = 456
				weight := decimal.NewFromFloat(123.456)
				price := decimal.NewFromFloat(7.89)

				p := goshopify.Product{
					ID:          productID,
					Title:       "myTitle",
					Vendor:      "myVendor",
					ProductType: "myProductType",
					Options: []goshopify.ProductOption{
						{Name: "myOption1Name"},
						{Name: "myOption2Name"},
						{Name: "myOption3Name"},
					},
					Metafields: []goshopify.Metafield{
						{
							Key:           "bar",
							Namespace:     "foo",
							OwnerResource: "product",
							Value:         "myProductMetafield",
						},
					},
					Variants: []goshopify.Variant{{
						ID:         variantID,
						ProductID:  productID,
						Price:      &price,
						Weight:     &weight,
						WeightUnit: "myWeightUnit",
						Sku:        "mySku",
						Barcode:    "myBarcode",
						Option1:    "myOption1Value",
						Option2:    "myOption2Value",
						Option3:    "myOption3Value",
						Metafields: []goshopify.Metafield{
							{
								Key:           "bar",
								Namespace:     "foo",
								OwnerResource: "variant",
								Value:         "myVariantMetafield",
							},
						},
					}},
				}
				return []goshopify.Product{p}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := groupVariants(tt.rows)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("\ngot: %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}
