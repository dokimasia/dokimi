// Package opaque ranks keys under a constraint of its own.
package opaque

// Keyed is a constraint of the keys of a ranking, which no type of the
// language satisfies without a method of its own.
type Keyed interface {
	Key() string
}

// Ranked ranks keys whose type the constraint Keyed admits.
//
//dokimi:suite
type Ranked[K Keyed] interface {
	// Rank returns the rank of the key.
	Rank(key K) int
}
