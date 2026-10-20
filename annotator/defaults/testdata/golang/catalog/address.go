// Package catalog declares a struct whose fields declare their defaults
// through the directive, through the tag and through both.
package catalog

// Address is the postal address of a customer.
type Address struct {
	// Street is the street and the house number.
	//
	//dokimi:default "Dam 1"
	Street string `default:"Main 1"`
	// City is the city.
	//
	//dokimi:default Amsterdam
	City string
	// Country is the ISO 3166-1 code of the country.
	Country string `default:"NL"`
	// Floor is the floor of the address, and zero for the ground floor.
	Floor int `default:"0"`
	// Lines are the extra lines of the address.
	Lines []string
}

// Key identifies an address, and a check that draws a key takes its default.
//
//dokimi:default "test-key"
type Key string

// Level is the floor level of an address, and declares no default.
type Level int
