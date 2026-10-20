// Package clash declares an interface whose call would have two fields of
// one name.
package clash

// Clash has a parameter named err beside its error result, so the fields of
// both are Err.
//
//dokimi:stub
type Clash interface {
	// Fail takes an error and returns one.
	Fail(err error) error
}
