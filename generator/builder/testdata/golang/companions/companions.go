// Package companions declares structs with functions that have the name of
// a companion and break its contract.
package companions

// Taken has a function with the name of its companion that takes a
// parameter.
//
//dokimi:builder
type Taken struct {
	// Name is set by the function.
	Name string
}

// TakenDefaults takes a parameter, so it is no companion.
func TakenDefaults(name string) Taken { return Taken{Name: name} }

// Paired has a function with the name of its companion that has two
// results.
//
//dokimi:builder
type Paired struct {
	// Name is set by the function.
	Name string
}

// PairedDefaults has two results, so it is no companion.
func PairedDefaults() (Paired, error) { return Paired{}, nil }

// Bare has a function with the name of its companion that has no result.
//
//dokimi:builder
type Bare struct {
	// Name is not set by the function.
	Name string
}

// BareDefaults has no result, so it is no companion.
func BareDefaults() {}

// Other has a function with the name of its companion that returns another
// type.
//
//dokimi:builder
type Other struct {
	// Name is not set by the function.
	Name string
}

// OtherDefaults returns a Taken, so it is no companion.
func OtherDefaults() Taken { return Taken{} }
