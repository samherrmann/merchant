package memdb

import (
	"embed"
	"encoding/json"
	"io"
	"text/template"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

//go:embed summary.tpl
var embeddedFS embed.FS

type Operations struct {
	tmpl *template.Template
	// NewProducts is a list of new products.
	NewProducts []goshopify.Product `json:",omitempty"`
	// ProductUpdates is a list of product updates.
	ProductUpdates []goshopify.Product `json:",omitempty"`
	// NewVariants is a list of new variants.
	NewVariants []goshopify.Variant `json:",omitempty"`
	// VariantUpdates is a list of variant updates.
	VariantUpdates []goshopify.Variant `json:",omitempty"`
}

// CreateProduct appends p to the NewProducts slice.
func (s *Operations) CreateProduct(p goshopify.Product) {
	s.NewProducts = append(s.NewProducts, p)
}

// UpdateProduct appends p to the ProductUpdates slice.
func (s *Operations) UpdateProduct(p goshopify.Product) {
	// Remove variants because they are updated separately.
	p.Variants = nil
	s.ProductUpdates = append(s.ProductUpdates, p)
}

// CreateVariant appends v to the NewVariants slice.
func (s *Operations) CreateVariant(v goshopify.Variant) {
	s.NewVariants = append(s.NewVariants, v)
}

// UpdateVariant appends v to the VariantUpdates slice.
func (s *Operations) UpdateVariant(v goshopify.Variant) {
	s.VariantUpdates = append(s.VariantUpdates, v)
}

// PrintJSON prints the JSON encoding of Operations to w.
func (s *Operations) PrintJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	return encoder.Encode(s)
}

// PrintSummary prints a summary of Operations to w.
func (s *Operations) PrintSummary(w io.Writer) error {
	if s.tmpl == nil {
		tmpl, err := template.New("summary").ParseFS(embeddedFS, "*")
		if err != nil {
			return err
		}
		s.tmpl = tmpl
	}
	return s.tmpl.ExecuteTemplate(w, "summary.tpl", s)
}
