package generic

import (
	"slices"
)

// PageBuilder constructs [Page] values. [NewPageBuilder] and
// [NewPageBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [PageBuilder.Build] returns
// the value. A builder is not safe for concurrent use.
type PageBuilder[T any] struct {
	v Page[T]
}

// NewPageBuilder returns a builder of the declared defaults. A member
// without a declared default has its zero value.
func NewPageBuilder[T any]() *PageBuilder[T] {
	var v Page[T]
	v.Next = 1
	return &PageBuilder[T]{v: v}
}

// NewPageBuilderFrom returns a builder of v.
func NewPageBuilderFrom[T any](v Page[T]) *PageBuilder[T] {
	return &PageBuilder[T]{v: v}
}

// WithItems replaces the elements of Items with v. A call without
// arguments empties Items.
func (b *PageBuilder[T]) WithItems(v ...T) *PageBuilder[T] {
	b.v.Items = v
	return b
}

// AppendItems appends v to the elements of Items.
func (b *PageBuilder[T]) AppendItems(v ...T) *PageBuilder[T] {
	b.v.Items = append(b.v.Items, v...)
	return b
}

// WithNext replaces Next with v.
func (b *PageBuilder[T]) WithNext(v int) *PageBuilder[T] {
	b.v.Next = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *PageBuilder[T]) Mutate(fn func(*Page[T])) *PageBuilder[T] {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *PageBuilder[T]) Clone() *PageBuilder[T] {
	out := &PageBuilder[T]{v: b.v}
	out.v.Items = slices.Clone(b.v.Items)
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [PageBuilder.Clone] returns a builder that shares none of them.
func (b *PageBuilder[T]) Build() Page[T] {
	return b.v
}
