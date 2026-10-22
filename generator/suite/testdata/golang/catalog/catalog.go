// Package catalog lists the products of the shop.
package catalog

// Catalog lists products by their identity, and is the subject of its own
// suite.
//
//dokimi:suite
//dokimi:subject Catalog new=New
type Catalog struct {
	products map[string]string
}

// New returns an empty catalog.
func New() *Catalog {
	return &Catalog{products: map[string]string{}}
}

// Find returns the name of the product with the identity id, and whether the
// catalog lists one.
//
//dokimi:budget allocs=0
func (c *Catalog) Find(id string) (string, bool) {
	name, found := c.products[id]
	return name, found
}

// Add lists a product under its identity.
func (c *Catalog) Add(id, name string) {
	c.products[id] = name
}

// reset removes every product.
func (c *Catalog) reset() {
	clear(c.products)
}
