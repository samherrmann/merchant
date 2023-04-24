package csv

import (
	"fmt"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

// NewMetafield returns a new metafield from the given path. A metafield path is
// of the following form:
//
//	<product/variant>.metafields.<namespace>.<key>
func NewMetafield(path string) (*goshopify.Metafield, error) {
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
