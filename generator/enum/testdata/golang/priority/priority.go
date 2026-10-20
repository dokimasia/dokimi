// Package priority declares an enumeration whose zero value is a variant.
package priority

// Priority is the urgency of a ticket.
//
//dokimi:enum
type Priority int

// The urgencies of a ticket, from the zero value.
const (
	Low Priority = iota
	Medium
	High
	Critical
)
