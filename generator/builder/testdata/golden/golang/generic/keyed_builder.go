package generic

// KeyedBuilder constructs [Keyed] values.
//
// Start with [NewKeyedBuilder] or [NewKeyedBuilderFrom], set what you care
// about, and call [KeyedBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Keyed breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [KeyedBuilder.Clone] is how one becomes two.
type KeyedBuilder[K Keyable] struct {
	v Keyed[K]
}

// NewKeyedBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewKeyedBuilder[K Keyable]() *KeyedBuilder[K] {
	return &KeyedBuilder[K]{}
}

// NewKeyedBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewKeyedBuilderFrom[K Keyable](v Keyed[K]) *KeyedBuilder[K] {
	return &KeyedBuilder[K]{v: v}
}

// WithKey replaces Key.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *KeyedBuilder[K]) WithKey(v K) *KeyedBuilder[K] {
	b.v.Key = v
	return b
}

// Mutate applies fn to the value under construction.
//
// The way out for what the setters do not reach: a member this
// builder cannot name, a shape it does not model, or a change that
// touches several members at once and only makes sense as one edit.
//
// Prefer a setter where one exists. A setter says which member it
// writes, and a reader of the call site learns that without following
// a closure.
func (b *KeyedBuilder[K]) Mutate(fn func(*Keyed[K])) *KeyedBuilder[K] {
	fn(&b.v)
	return b
}

// Clone returns a builder that can be configured without affecting this
// one.
//
// One level deep, and the boundary matters. The sequence, byte-sequence
// and keyed members are copied, so appending or setting an entry through
// one clone is not visible through another.
//
// What is shared: anything held *inside* those — a struct in a slice that
// owns a slice of its own — and every optional member, which is what a
// reference means. Mutating through one of those does reach the original.
// Use [KeyedBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *KeyedBuilder[K]) Clone() *KeyedBuilder[K] {
	out := &KeyedBuilder[K]{v: b.v}
	return out
}

// Build returns the constructed [Keyed].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [KeyedBuilder.Clone] is what separates them.
func (b *KeyedBuilder[K]) Build() Keyed[K] {
	return b.v
}
