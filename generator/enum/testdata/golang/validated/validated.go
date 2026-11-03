// Package validated declares enumerations that kanon encodes as a
// kanon.Validator through their underlying types, and ones that it does not.
package validated

//go:generate go tool kanon -validate=valid -type=Grade,Record

// Grade is a grade that the kanon directive of its file names.
//
//dokimi:enum
type Grade uint8

// The grades.
const (
	Poor Grade = iota + 1
	Good
)

// Record is a record of a grade.
//
//go:generate go tool kanon -type=Voice
type Record struct {
	Grade Grade `kanon:"1"`
}

// Mean returns the mean grade of the record.
//
//go:generate go tool kanon -type=Clef
func (r Record) Mean() Grade { return r.Grade }

// grades is the number of the grades.
//
//go:generate go tool kanon -type=Tempo
const grades = 2

// Passing is the grade that passes.
//
//go:generate go tool kanon -type=Meter
var Passing = Good

// Graded is a value with a grade.
//
//go:generate go tool kanon -type=Chord
type Graded interface{ Grade() Grade }

// Mark is another name of Grade.
//
//go:generate go tool kanon -type=Pitch
type Mark = Grade

// Passes reports whether g passes.
//
//go:generate go tool kanon -type=Scale
func Passes(g Grade) bool { return g == Passing && grades > 1 }

// Tier is a tier with the method ValidateKanon, as kanon writes it.
//
//dokimi:enum
type Tier uint8

// The tiers.
const (
	Free Tier = iota + 1
	Paid
)

// ValidateKanon accepts every tier.
func (Tier) ValidateKanon() error { return nil }

// Hue is a hue that a directive of another command names.
//
//go:generate go tool kanonical -type=Hue
//dokimi:enum
type Hue uint8

// The hues.
const (
	Warm Hue = iota + 1
	Cold
)
