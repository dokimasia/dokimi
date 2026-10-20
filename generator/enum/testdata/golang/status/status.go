// Package status declares an enumeration whose zero value is no variant.
package status

// Status is the state of an order.
//
//dokimi:enum
type Status int

// The states of an order.
const (
	Draft Status = iota + 1
	Published
	Archived
)
