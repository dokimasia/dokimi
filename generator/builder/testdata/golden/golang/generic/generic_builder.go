package generic

// PageBuilder constructs [Page] values.
//
// Start with [NewPageBuilder] or [NewPageBuilderFrom], set what you care
// about, and call [PageBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Page breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [PageBuilder.Clone] is how one becomes two.
type PageBuilder[T any] struct {
	v Page[T]
}

// NewPageBuilder returns a builder seeded with the declared defaults.
//
// Only the members that declare one. Everything else holds its zero value,
// so a caller that sets nothing gets a mixture rather than either extreme
// — which is what makes a check against a built value worth writing.
func NewPageBuilder[T any]() *PageBuilder[T] {
	var v Page[T]
	v.Next = 1
	return &PageBuilder[T]{v: v}
}

// NewPageBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewPageBuilderFrom[T any](v Page[T]) *PageBuilder[T] {
	return &PageBuilder[T]{v: v}
}

// WithItems replaces Items with the elements given.
//
// Variadic, and replacing: whatever the collection held is discarded,
// including anything a constructor seeded. AppendItems is the counterpart
// that keeps it. Called with no arguments this empties the member rather
// than leaving it alone, which is how a caller clears a seeded default.
func (b *PageBuilder[T]) WithItems(v ...T) *PageBuilder[T] {
	b.v.Items = v
	return b
}

// AppendItems adds to Items, keeping what is already there.
//
// The counterpart to WithItems, which discards it. Both are offered
// because a caller extending a seeded default and a caller replacing it
// want opposite things, and one method cannot mean both.
func (b *PageBuilder[T]) AppendItems(v ...T) *PageBuilder[T] {
	b.v.Items = append(b.v.Items, v...)
	return b
}

// WithNext replaces Next.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *PageBuilder[T]) WithNext(v int) *PageBuilder[T] {
	b.v.Next = v
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
func (b *PageBuilder[T]) Mutate(fn func(*Page[T])) *PageBuilder[T] {
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
// Use [PageBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *PageBuilder[T]) Clone() *PageBuilder[T] {
	out := &PageBuilder[T]{v: b.v}
	if b.v.Items != nil {
		out.v.Items = append([]T(nil), b.v.Items...)
	}
	return out
}

// Build returns the constructed [Page].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [PageBuilder.Clone] is what separates them.
func (b *PageBuilder[T]) Build() Page[T] {
	return b.v
}
