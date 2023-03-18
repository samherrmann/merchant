package memdb

import (
	"fmt"
	"testing"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

func Test_validateProductID(t *testing.T) {
	tests := []struct {
		name     string
		current  *goshopify.Variant
		incoming *goshopify.Variant
		wantErr  bool
	}{
		{
			name:     "variant ID mismatch",
			current:  &goshopify.Variant{ID: 1, ProductID: 1},
			incoming: &goshopify.Variant{ID: 2, ProductID: 1},
			wantErr:  true,
		},
		{
			name:     "product ID mismatch",
			current:  &goshopify.Variant{ID: 1, ProductID: 1},
			incoming: &goshopify.Variant{ID: 1, ProductID: 2},
			wantErr:  true,
		},
		{
			name:     "product ID match",
			current:  &goshopify.Variant{ID: 1, ProductID: 2},
			incoming: &goshopify.Variant{ID: 1, ProductID: 2},
			wantErr:  false,
		},
		{
			name:     "incoming product ID not set",
			current:  &goshopify.Variant{ID: 1, ProductID: 2},
			incoming: &goshopify.Variant{ID: 1, ProductID: 0},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProductID(tt.current, tt.incoming)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateStringProperty(t *testing.T) {

	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{
			name: "both strings are empty",
			s1:   "",
			s2:   "",
			want: true,
		},
		{
			name: "s1 is empty",
			s1:   "",
			s2:   "123",
			want: true,
		},
		{
			name: "s2 is empty",
			s1:   "123",
			s2:   "",
			want: true,
		},
		{
			name: "strings are equal",
			s1:   "123",
			s2:   "123",
			want: true,
		},
		{
			name: "strings are not equal",
			s1:   "123",
			s2:   "456",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := equalNonEmptyStrings(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionEqual(t *testing.T) {
	tests := []struct {
		o1   string
		o2   string
		want bool
	}{
		{"", "", true},
		{defaultOptionValue, "", true},
		{"", defaultOptionValue, true},
		{defaultOptionValue, defaultOptionValue, true},
		{"foo", defaultOptionValue, false},
		{defaultOptionValue, "foo", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("o1 = %v, o2 = %v", tt.o1, tt.o2), func(t *testing.T) {
			if got := optionEqual(tt.o1, tt.o2); got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyVariantIDs(t *testing.T) {
	tests := []struct {
		src     *goshopify.Variant
		dst     *goshopify.Variant
		wantErr bool
	}{
		{
			src:     &goshopify.Variant{},
			dst:     &goshopify.Variant{},
			wantErr: false,
		},
		{
			src:     &goshopify.Variant{ID: 1},
			dst:     &goshopify.Variant{ID: 1},
			wantErr: false,
		},
		{
			src:     &goshopify.Variant{ID: 1},
			dst:     &goshopify.Variant{ID: 2},
			wantErr: true,
		},
		{
			src:     &goshopify.Variant{ProductID: 1},
			dst:     &goshopify.Variant{ProductID: 1},
			wantErr: false,
		},
		{
			src:     &goshopify.Variant{ProductID: 1},
			dst:     &goshopify.Variant{ProductID: 2},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %v", i), func(t *testing.T) {
			if err := copyVariantIDs(tt.src, tt.dst); (err != nil) != tt.wantErr {
				t.Fatalf("got error %q, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				if tt.dst.ID != tt.src.ID {
					t.Fatalf("dst = %v, src = %v", tt.dst.ID, tt.src.ID)
				}
				if tt.dst.ProductID != tt.src.ProductID {
					t.Fatalf("dst = %v, src = %v", tt.dst.ProductID, tt.src.ProductID)
				}
			}
		})
	}
}
