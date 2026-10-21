// Package shaped declares enumerations whose values are integers of a type
// that a table of Go does not serve.
package shaped

// Ratio is a float whose values are integers. An index of an array takes no
// constant of a float type.
//
//dokimi:enum
type Ratio float64

// The ratios, each a whole number.
const (
	Single Ratio = 1
	Double Ratio = 2
)

// Base is the defined type under Tier.
type Base int

// Tier is a defined type over another defined type.
//
//dokimi:enum
type Tier Base

// The tiers of a plan.
const (
	Bronze Tier = iota + 1
	Silver
)
