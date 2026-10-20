// Package catalog declares the structs that the tests of gentest find.
package catalog

// Address is the postal address of a customer.
type Address struct {
	// Street is the street and the house number.
	Street string
}

// Warehouse is a struct with a field of the same name as a field of
// Address.
type Warehouse struct {
	// Street is the street of the warehouse.
	Street string
}
