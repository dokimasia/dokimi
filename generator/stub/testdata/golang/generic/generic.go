// Package generic declares generic interfaces under the directive stub.
package generic

// Repo has a type parameter with a witness that the rules derive.
//
//dokimi:stub
type Repo[T any] interface {
	// Get returns a value of the type parameter.
	Get(id string) (T, error)
	// Put takes a value of the type parameter. Its authored sample is no
	// value of the type parameter in the generic checks.
	Put(
		//dokimi:sample value=7
		v T,
	)
	// List returns a slice of the type parameter.
	List() ([]T, error)
	// Boxed returns a generic struct at the type parameter.
	Boxed() Box[T]
}

// Pairs has two type parameters whose names contain a digit and an
// underscore, and an authored witness of each.
//
//dokimi:stub
//dokimi:witness K1=string V_=int
type Pairs[K1 comparable, V_ any] interface {
	// All returns a map of the type parameters.
	All() map[K1]V_
}

// Box contains a value of its type parameter.
type Box[T any] struct {
	// V is the value.
	V T
}

// Keyable has type-set elements, so it only constrains type parameters.
//
//dokimi:stub
type Keyable interface {
	~string | ~float64
}

// Bag has a type parameter without a witness, so its double has no test.
//
//dokimi:stub
type Bag[K Keyable] interface {
	// Add takes a value of the type parameter.
	Add(k K)
}
