package metafields

import (
	"fmt"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

// NewMetafield returns a new metafield from the given path. A metafield path is
// of the following form:
//
//	<product/variant>.metafields.<namespace>.<key>
func New(path string) (*goshopify.Metafield, error) {
	parts := strings.Split(path, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("path %q does not have four parts", path)
	}
	if parts[0] != "product" && parts[0] != "variant" {
		return nil, fmt.Errorf(`first part must be "product" or "variant", got %q`, parts[0])
	}
	if parts[1] != "metafields" {
		return nil, fmt.Errorf(`second part must be "metafields", got %q`, parts[1])
	}

	m := &goshopify.Metafield{
		OwnerResource: parts[0],
		Namespace:     parts[2],
		Key:           parts[3],
	}

	return m, nil
}

// IndexOf returns the first index at which a matching metafield can be found in
// the slice, or -1 if it is not present.
func IndexOf(slice []goshopify.Metafield, m *goshopify.Metafield) int {
	for i, v := range slice {
		if v.Namespace == m.Namespace &&
			v.Key == m.Key &&
			v.OwnerResource == m.OwnerResource {
			return i
		}
	}
	return -1
}

// Attach appends m to slice if slice does not already contain a matching
// metafield. If slice already has a matching metafield, then that metafield is
// overwritten with m.
func Attach(slice []goshopify.Metafield, m *goshopify.Metafield) []goshopify.Metafield {
	i := IndexOf(slice, m)
	if i == -1 {
		slice = append(slice, *m)
	} else {
		slice[i] = *m
	}
	return slice
}
