package companions

// TakenBuilder constructs [Taken] values.
//
// Start with [NewTakenBuilder] or [NewTakenBuilderFrom], set what you care
// about, and call [TakenBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Taken breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [TakenBuilder.Clone] is how one becomes two.
type TakenBuilder struct {
	v Taken
}

// NewTakenBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewTakenBuilder() *TakenBuilder {
	return &TakenBuilder{}
}

// NewTakenBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewTakenBuilderFrom(v Taken) *TakenBuilder {
	return &TakenBuilder{v: v}
}

// WithName replaces Name.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *TakenBuilder) WithName(v string) *TakenBuilder {
	b.v.Name = v
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
func (b *TakenBuilder) Mutate(fn func(*Taken)) *TakenBuilder {
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
// Use [TakenBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *TakenBuilder) Clone() *TakenBuilder {
	out := &TakenBuilder{v: b.v}
	return out
}

// Build returns the constructed [Taken].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [TakenBuilder.Clone] is what separates them.
func (b *TakenBuilder) Build() Taken {
	return b.v
}

// PairedBuilder constructs [Paired] values.
//
// Start with [NewPairedBuilder] or [NewPairedBuilderFrom], set what you
// care about, and call [PairedBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Paired breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [PairedBuilder.Clone] is how one becomes two.
type PairedBuilder struct {
	v Paired
}

// NewPairedBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewPairedBuilder() *PairedBuilder {
	return &PairedBuilder{}
}

// NewPairedBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewPairedBuilderFrom(v Paired) *PairedBuilder {
	return &PairedBuilder{v: v}
}

// WithName replaces Name.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *PairedBuilder) WithName(v string) *PairedBuilder {
	b.v.Name = v
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
func (b *PairedBuilder) Mutate(fn func(*Paired)) *PairedBuilder {
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
// Use [PairedBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *PairedBuilder) Clone() *PairedBuilder {
	out := &PairedBuilder{v: b.v}
	return out
}

// Build returns the constructed [Paired].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [PairedBuilder.Clone] is what separates them.
func (b *PairedBuilder) Build() Paired {
	return b.v
}

// BareBuilder constructs [Bare] values.
//
// Start with [NewBareBuilder] or [NewBareBuilderFrom], set what you care
// about, and call [BareBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Bare breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [BareBuilder.Clone] is how one becomes two.
type BareBuilder struct {
	v Bare
}

// NewBareBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewBareBuilder() *BareBuilder {
	return &BareBuilder{}
}

// NewBareBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewBareBuilderFrom(v Bare) *BareBuilder {
	return &BareBuilder{v: v}
}

// WithName replaces Name.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *BareBuilder) WithName(v string) *BareBuilder {
	b.v.Name = v
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
func (b *BareBuilder) Mutate(fn func(*Bare)) *BareBuilder {
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
// Use [BareBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *BareBuilder) Clone() *BareBuilder {
	out := &BareBuilder{v: b.v}
	return out
}

// Build returns the constructed [Bare].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [BareBuilder.Clone] is what separates them.
func (b *BareBuilder) Build() Bare {
	return b.v
}

// OtherBuilder constructs [Other] values.
//
// Start with [NewOtherBuilder] or [NewOtherBuilderFrom], set what you care
// about, and call [OtherBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Other breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [OtherBuilder.Clone] is how one becomes two.
type OtherBuilder struct {
	v Other
}

// NewOtherBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewOtherBuilder() *OtherBuilder {
	return &OtherBuilder{}
}

// NewOtherBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewOtherBuilderFrom(v Other) *OtherBuilder {
	return &OtherBuilder{v: v}
}

// WithName replaces Name.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *OtherBuilder) WithName(v string) *OtherBuilder {
	b.v.Name = v
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
func (b *OtherBuilder) Mutate(fn func(*Other)) *OtherBuilder {
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
// Use [OtherBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *OtherBuilder) Clone() *OtherBuilder {
	out := &OtherBuilder{v: b.v}
	return out
}

// Build returns the constructed [Other].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [OtherBuilder.Clone] is what separates them.
func (b *OtherBuilder) Build() Other {
	return b.v
}
