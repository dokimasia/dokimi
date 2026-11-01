// Package scoped declares at package scope the names of the packages that the
// builder of Basket imports.
package scoped

// The variables that take the names of the packages of the builder.
var (
	maps   = "the copy of the entries of a map"
	slices = "the copy of the elements of a list"
)

// Basket is a basket of items.
//
//dokimi:builder
type Basket struct {
	// Items are the names of the items in the basket.
	Items []string
	// Counts are the numbers of the items by their names.
	Counts map[string]int
}
