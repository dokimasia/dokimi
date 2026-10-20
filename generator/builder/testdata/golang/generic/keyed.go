package generic

// Keyable bounds a type parameter with a type set that int is not in, so
// the rules derive no witness for it.
type Keyable interface {
	~string | ~float64
}

// Keyed has a type parameter without a witness, so its builder has no
// checks. It is the one builder of its source, so the file of its checks
// contains the note alone.
//
//dokimi:builder
type Keyed[K Keyable] struct {
	// Key refers to the type parameter.
	Key K
}
