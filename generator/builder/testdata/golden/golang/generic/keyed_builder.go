package generic

// KeyedBuilder constructs [Keyed] values. [NewKeyedBuilder] and
// [NewKeyedBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [KeyedBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type KeyedBuilder[K Keyable] struct {
	v Keyed[K]
}

// NewKeyedBuilder returns a builder of the zero value.
func NewKeyedBuilder[K Keyable]() *KeyedBuilder[K] {
	return &KeyedBuilder[K]{}
}

// NewKeyedBuilderFrom returns a builder of v.
func NewKeyedBuilderFrom[K Keyable](v Keyed[K]) *KeyedBuilder[K] {
	return &KeyedBuilder[K]{v: v}
}

// WithKey replaces Key with v.
func (b *KeyedBuilder[K]) WithKey(v K) *KeyedBuilder[K] {
	b.v.Key = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *KeyedBuilder[K]) Mutate(fn func(*Keyed[K])) *KeyedBuilder[K] {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *KeyedBuilder[K]) Clone() *KeyedBuilder[K] {
	out := &KeyedBuilder[K]{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [KeyedBuilder.Clone] returns a builder that shares none of them.
func (b *KeyedBuilder[K]) Build() Keyed[K] {
	return b.v
}
