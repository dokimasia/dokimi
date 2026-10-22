// Package custom declares the validity method of its enumeration as Valid.
package custom

// Gear is a gear of a bicycle, whose validity method is Valid.
//
//dokimi:enum
type Gear int

// The gears of a bicycle.
const (
	Low Gear = iota + 1
	High
)

// Valid reports whether g is a gear.
func (g Gear) Valid() bool { return g == Low || g == High }
