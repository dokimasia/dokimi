// Package bad declares the structs whose builders builder reports.
package bad

// Shelf has a default that no integer spells.
//
//dokimi:builder
type Shelf struct {
	// Slots has a default of the wrong type.
	Slots int `default:"many"`
}

// Cell is a struct that no default names.
type Cell struct {
	// X is the column of the cell.
	X int
}

// Board has a default that is no value of a struct.
//
//dokimi:builder
type Board struct {
	// Origin has a default that names nothing.
	Origin Cell `default:"origin"`
}

// Ledger has a field whose setter has the name of the entry setter of
// another field.
//
//dokimi:builder
type Ledger struct {
	// Labels is a mapping, whose entry setter is WithLabelsEntry.
	Labels map[string]string
	// LabelsEntry has the setter WithLabelsEntry too.
	LabelsEntry string
}

// Typo has a tag builder with another value than -.
//
//dokimi:builder
type Typo struct {
	// Name has a tag with a typo, so it keeps its setter.
	Name string `builder:"skip"`
}

// Blank has no field that a builder can set.
//
//dokimi:builder
type Blank struct {
	// hidden is unexported.
	hidden string
}
