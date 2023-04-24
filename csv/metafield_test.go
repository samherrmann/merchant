package csv

import (
	"reflect"
	"testing"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

func TestNewMetafield(t *testing.T) {
	tests := []struct {
		path    string
		want    *goshopify.Metafield
		wantErr bool
	}{
		{
			path:    "",
			wantErr: true,
		},
		{
			path:    "foo.bar",
			wantErr: true,
		},
		{
			path:    "foo.bar.baz.qux",
			wantErr: true,
		},
		{
			path:    "foo.bar.baz.qux.quux",
			wantErr: true,
		},
		{
			path:    "product.bar.baz.qux",
			wantErr: true,
		},
		{
			path:    "foo.metafields.baz.qux",
			wantErr: true,
		},
		{
			path: "product.metafields.baz.qux",
			want: &goshopify.Metafield{
				Key:           "qux",
				Namespace:     "baz",
				OwnerResource: "product",
			},
		},
		{
			path: "variant.metafields.baz.qux",
			want: &goshopify.Metafield{
				Key:           "qux",
				Namespace:     "baz",
				OwnerResource: "variant",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := NewMetafield(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
