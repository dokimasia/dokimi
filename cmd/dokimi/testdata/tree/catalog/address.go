// Package catalog declares the struct that the binary writes a builder of,
// and the enumeration that it writes a surface of.
package catalog

// Kind is the kind of an address.
//
//dokimi:enum
type Kind int

// The kinds of an address.
const (
	Home Kind = iota + 1
	Work
)

// Address is the postal address of a customer.
//
//dokimi:builder
type Address struct {
	// Street is the street and the house number.
	//
	//dokimi:default "Dam 1"
	Street string
	// Lines are the lines of the address.
	Lines []string
}
