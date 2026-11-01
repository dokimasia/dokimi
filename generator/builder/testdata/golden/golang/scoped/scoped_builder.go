package scoped

import (
	maps2 "maps"
	slices2 "slices"
)

// BasketBuilder constructs [Basket] values. [NewBasketBuilder] and
// [NewBasketBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [BasketBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type BasketBuilder struct {
	v Basket
}

// NewBasketBuilder returns a builder of the zero value.
func NewBasketBuilder() *BasketBuilder {
	return &BasketBuilder{}
}

// NewBasketBuilderFrom returns a builder of v.
func NewBasketBuilderFrom(v Basket) *BasketBuilder {
	return &BasketBuilder{v: v}
}

// WithItems replaces the elements of Items with v. A call without
// arguments empties Items.
func (b *BasketBuilder) WithItems(v ...string) *BasketBuilder {
	b.v.Items = v
	return b
}

// AppendItems appends v to the elements of Items.
func (b *BasketBuilder) AppendItems(v ...string) *BasketBuilder {
	b.v.Items = append(b.v.Items, v...)
	return b
}

// WithCounts replaces Counts with v. The builder keeps v without a copy.
func (b *BasketBuilder) WithCounts(v map[string]int) *BasketBuilder {
	b.v.Counts = v
	return b
}

// WithCountsEntry sets the key k of Counts to v. It allocates the map
// where Counts is nil.
func (b *BasketBuilder) WithCountsEntry(k string, v int) *BasketBuilder {
	if b.v.Counts == nil {
		b.v.Counts = make(map[string]int)
	}
	b.v.Counts[k] = v
	return b
}

// WithCountsEntries copies each entry of entries into Counts. A key of
// entries replaces the value of that key in Counts. It allocates the map
// where Counts is nil.
func (b *BasketBuilder) WithCountsEntries(entries map[string]int) *BasketBuilder {
	if b.v.Counts == nil {
		b.v.Counts = make(map[string]int, len(entries))
	}
	maps2.Copy(b.v.Counts, entries)
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *BasketBuilder) Mutate(fn func(*Basket)) *BasketBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *BasketBuilder) Clone() *BasketBuilder {
	out := &BasketBuilder{v: b.v}
	out.v.Items = slices2.Clone(b.v.Items)
	out.v.Counts = maps2.Clone(b.v.Counts)
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [BasketBuilder.Clone] returns a builder that shares none of them.
func (b *BasketBuilder) Build() Basket {
	return b.v
}
