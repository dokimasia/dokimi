// Package values declares a default of each kind of type.
package values

// Code is a declared type of text.
type Code string

// Street is the street of a default address.
const Street = "Dam 2"

// Count is a declared type of a number.
type Count int

// Spot is a struct.
type Spot struct {
	X int
}

// Loop1 and Loop2 are declared types over each other, which no compiler
// accepts.
type (
	Loop1 Loop2
	Loop2 Loop1
)

// Values declares a default on a field of each kind of type.
type Values struct {
	Number int    `default:"7"`
	Text   string `default:"Dam 1"`
	Named  string `default:"Street"`
	Code   Code   `default:"abc"`
	Bad    int    `default:"nope"`
	Count  Count  `default:"nope"`
	Spot   Spot   `default:"nope"`
	Loop   Loop1  `default:"nope"`
}
