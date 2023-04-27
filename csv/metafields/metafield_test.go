package metafields

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
			got, err := New(tt.path)
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

func TestAttach(t *testing.T) {
	tests := []struct {
		name  string
		slice []goshopify.Metafield
		m     *goshopify.Metafield
		want  []goshopify.Metafield
	}{
		{
			name: "should add new metafield",
			m: &goshopify.Metafield{
				Key:           "foo",
				Namespace:     "bar",
				OwnerResource: "product",
			},
			want: []goshopify.Metafield{
				{
					Key:           "foo",
					Namespace:     "bar",
					OwnerResource: "product",
				},
			},
		}, {
			name: "should update existing metafield",
			slice: []goshopify.Metafield{
				{
					Key:           "foo",
					Namespace:     "bar",
					OwnerResource: "product",
					Value:         1,
				},
			},
			m: &goshopify.Metafield{
				Key:           "foo",
				Namespace:     "bar",
				OwnerResource: "product",
				Value:         2,
			},
			want: []goshopify.Metafield{
				{
					Key:           "foo",
					Namespace:     "bar",
					OwnerResource: "product",
					Value:         2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Attach(tt.slice, tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
