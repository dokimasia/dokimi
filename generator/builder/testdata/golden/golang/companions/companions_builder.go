package companions

// TakenBuilder constructs [Taken] values. [NewTakenBuilder] and
// [NewTakenBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [TakenBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type TakenBuilder struct {
	v Taken
}

// NewTakenBuilder returns a builder of the zero value.
func NewTakenBuilder() *TakenBuilder {
	return &TakenBuilder{}
}

// NewTakenBuilderFrom returns a builder of v.
func NewTakenBuilderFrom(v Taken) *TakenBuilder {
	return &TakenBuilder{v: v}
}

// WithName replaces Name with v.
func (b *TakenBuilder) WithName(v string) *TakenBuilder {
	b.v.Name = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *TakenBuilder) Mutate(fn func(*Taken)) *TakenBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *TakenBuilder) Clone() *TakenBuilder {
	out := &TakenBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [TakenBuilder.Clone] returns a builder that shares none of them.
func (b *TakenBuilder) Build() Taken {
	return b.v
}

// PairedBuilder constructs [Paired] values. [NewPairedBuilder] and
// [NewPairedBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [PairedBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type PairedBuilder struct {
	v Paired
}

// NewPairedBuilder returns a builder of the zero value.
func NewPairedBuilder() *PairedBuilder {
	return &PairedBuilder{}
}

// NewPairedBuilderFrom returns a builder of v.
func NewPairedBuilderFrom(v Paired) *PairedBuilder {
	return &PairedBuilder{v: v}
}

// WithName replaces Name with v.
func (b *PairedBuilder) WithName(v string) *PairedBuilder {
	b.v.Name = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *PairedBuilder) Mutate(fn func(*Paired)) *PairedBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *PairedBuilder) Clone() *PairedBuilder {
	out := &PairedBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [PairedBuilder.Clone] returns a builder that shares none of them.
func (b *PairedBuilder) Build() Paired {
	return b.v
}

// BareBuilder constructs [Bare] values. [NewBareBuilder] and
// [NewBareBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [BareBuilder.Build] returns
// the value. A builder is not safe for concurrent use.
type BareBuilder struct {
	v Bare
}

// NewBareBuilder returns a builder of the zero value.
func NewBareBuilder() *BareBuilder {
	return &BareBuilder{}
}

// NewBareBuilderFrom returns a builder of v.
func NewBareBuilderFrom(v Bare) *BareBuilder {
	return &BareBuilder{v: v}
}

// WithName replaces Name with v.
func (b *BareBuilder) WithName(v string) *BareBuilder {
	b.v.Name = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *BareBuilder) Mutate(fn func(*Bare)) *BareBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *BareBuilder) Clone() *BareBuilder {
	out := &BareBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [BareBuilder.Clone] returns a builder that shares none of them.
func (b *BareBuilder) Build() Bare {
	return b.v
}

// OtherBuilder constructs [Other] values. [NewOtherBuilder] and
// [NewOtherBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [OtherBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type OtherBuilder struct {
	v Other
}

// NewOtherBuilder returns a builder of the zero value.
func NewOtherBuilder() *OtherBuilder {
	return &OtherBuilder{}
}

// NewOtherBuilderFrom returns a builder of v.
func NewOtherBuilderFrom(v Other) *OtherBuilder {
	return &OtherBuilder{v: v}
}

// WithName replaces Name with v.
func (b *OtherBuilder) WithName(v string) *OtherBuilder {
	b.v.Name = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *OtherBuilder) Mutate(fn func(*Other)) *OtherBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *OtherBuilder) Clone() *OtherBuilder {
	out := &OtherBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [OtherBuilder.Clone] returns a builder that shares none of them.
func (b *OtherBuilder) Build() Other {
	return b.v
}
