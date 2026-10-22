package catalog

import (
	"maps"
	"slices"
	"time"
)

// AddressBuilder constructs [Address] values. [NewAddressBuilder] and
// [NewAddressBuilderFrom] return a new builder. Each setter sets one
// member and returns the builder, so the calls chain.
// [AddressBuilder.Build] returns the value. A builder is not safe for
// concurrent use.
type AddressBuilder struct {
	v Address
}

// NewAddressBuilder returns a builder of the declared defaults. A member
// without a declared default has its zero value.
func NewAddressBuilder() *AddressBuilder {
	var v Address
	v.Street = "Dam 1"
	v.City = "Amsterdam"
	v.Timeout = time.Second
	return &AddressBuilder{v: v}
}

// NewAddressBuilderFrom returns a builder of v.
func NewAddressBuilderFrom(v Address) *AddressBuilder {
	return &AddressBuilder{v: v}
}

// WithStreet replaces Street with v.
func (b *AddressBuilder) WithStreet(v string) *AddressBuilder {
	b.v.Street = v
	return b
}

// WithCity replaces City with v.
func (b *AddressBuilder) WithCity(v string) *AddressBuilder {
	b.v.City = v
	return b
}

// WithFloor points Floor at a copy of v. [AddressBuilder.Mutate] can set
// Floor to nil or to another address.
func (b *AddressBuilder) WithFloor(v int) *AddressBuilder {
	b.v.Floor = &v
	return b
}

// WithLines replaces the elements of Lines with v. A call without
// arguments empties Lines.
func (b *AddressBuilder) WithLines(v ...string) *AddressBuilder {
	b.v.Lines = v
	return b
}

// AppendLines appends v to the elements of Lines.
func (b *AddressBuilder) AppendLines(v ...string) *AddressBuilder {
	b.v.Lines = append(b.v.Lines, v...)
	return b
}

// WithLabels replaces Labels with v. The builder keeps v without a copy.
func (b *AddressBuilder) WithLabels(v map[string]string) *AddressBuilder {
	b.v.Labels = v
	return b
}

// WithLabelsEntry sets the key k of Labels to v. It allocates the map
// where Labels is nil.
func (b *AddressBuilder) WithLabelsEntry(k, v string) *AddressBuilder {
	if b.v.Labels == nil {
		b.v.Labels = make(map[string]string)
	}
	b.v.Labels[k] = v
	return b
}

// WithLabelsEntries copies each entry of entries into Labels. A key of
// entries replaces the value of that key in Labels. It allocates the map
// where Labels is nil.
func (b *AddressBuilder) WithLabelsEntries(entries map[string]string) *AddressBuilder {
	if b.v.Labels == nil {
		b.v.Labels = make(map[string]string, len(entries))
	}
	maps.Copy(b.v.Labels, entries)
	return b
}

// WithTags replaces Tags with v. The builder keeps v without a copy.
func (b *AddressBuilder) WithTags(v map[string]struct{}) *AddressBuilder {
	b.v.Tags = v
	return b
}

// WithTagsEntry adds the key k to Tags. It allocates the set where Tags is
// nil.
func (b *AddressBuilder) WithTagsEntry(k string) *AddressBuilder {
	if b.v.Tags == nil {
		b.v.Tags = make(map[string]struct{})
	}
	b.v.Tags[k] = struct{}{}
	return b
}

// WithTagsEntries adds each key of keys to Tags. It allocates the set
// where Tags is nil.
func (b *AddressBuilder) WithTagsEntries(keys ...string) *AddressBuilder {
	if b.v.Tags == nil {
		b.v.Tags = make(map[string]struct{}, len(keys))
	}
	for _, k := range keys {
		b.v.Tags[k] = struct{}{}
	}
	return b
}

// WithRaw replaces Raw with v. The builder keeps v without a copy.
func (b *AddressBuilder) WithRaw(v []byte) *AddressBuilder {
	b.v.Raw = v
	return b
}

// WithRawString replaces Raw with the bytes of s.
func (b *AddressBuilder) WithRawString(s string) *AddressBuilder {
	b.v.Raw = []byte(s)
	return b
}

// WithGrid replaces Grid with v.
func (b *AddressBuilder) WithGrid(v [2]int) *AddressBuilder {
	b.v.Grid = v
	return b
}

// WithTimeout replaces Timeout with v.
func (b *AddressBuilder) WithTimeout(v time.Duration) *AddressBuilder {
	b.v.Timeout = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *AddressBuilder) Mutate(fn func(*Address)) *AddressBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *AddressBuilder) Clone() *AddressBuilder {
	out := &AddressBuilder{v: b.v}
	out.v.Lines = slices.Clone(b.v.Lines)
	out.v.Labels = maps.Clone(b.v.Labels)
	out.v.Tags = maps.Clone(b.v.Tags)
	out.v.Raw = slices.Clone(b.v.Raw)
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [AddressBuilder.Clone] returns a builder that shares none of them.
func (b *AddressBuilder) Build() Address {
	return b.v
}

// ProfileBuilder constructs [Profile] values. [NewProfileBuilder] and
// [NewProfileBuilderFrom] return a new builder. Each setter sets one
// member and returns the builder, so the calls chain.
// [ProfileBuilder.Build] returns the value. A builder is not safe for
// concurrent use.
type ProfileBuilder struct {
	v Profile
}

// NewProfileBuilder returns a builder of the declared defaults. A member
// without a declared default has its zero value.
func NewProfileBuilder() *ProfileBuilder {
	var v Profile
	v.Active = false
	return &ProfileBuilder{v: v}
}

// NewProfileBuilderFrom returns a builder of v.
func NewProfileBuilderFrom(v Profile) *ProfileBuilder {
	return &ProfileBuilder{v: v}
}

// WithRole replaces Role with v.
func (b *ProfileBuilder) WithRole(v Role) *ProfileBuilder {
	b.v.Role = v
	return b
}

// WithAudit replaces Audit with v.
func (b *ProfileBuilder) WithAudit(v *Audit) *ProfileBuilder {
	b.v.Audit = v
	return b
}

// WithActive replaces Active with v.
func (b *ProfileBuilder) WithActive(v bool) *ProfileBuilder {
	b.v.Active = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *ProfileBuilder) Mutate(fn func(*Profile)) *ProfileBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *ProfileBuilder) Clone() *ProfileBuilder {
	out := &ProfileBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [ProfileBuilder.Clone] returns a builder that shares none of them.
func (b *ProfileBuilder) Build() Profile {
	return b.v
}

// EnvelopeBuilder constructs [Envelope] values. [NewEnvelopeBuilder] and
// [NewEnvelopeBuilderFrom] return a new builder. Each setter sets one
// member and returns the builder, so the calls chain.
// [EnvelopeBuilder.Build] returns the value. A builder is not safe for
// concurrent use.
type EnvelopeBuilder struct {
	v Envelope
}

// NewEnvelopeBuilder returns a builder of the zero value.
func NewEnvelopeBuilder() *EnvelopeBuilder {
	return &EnvelopeBuilder{}
}

// NewEnvelopeBuilderFrom returns a builder of v.
func NewEnvelopeBuilderFrom(v Envelope) *EnvelopeBuilder {
	return &EnvelopeBuilder{v: v}
}

// WithNotifier replaces Notifier with v.
func (b *EnvelopeBuilder) WithNotifier(v Notifier) *EnvelopeBuilder {
	b.v.Notifier = v
	return b
}

// WithCodes replaces Codes with v.
func (b *EnvelopeBuilder) WithCodes(v Codes) *EnvelopeBuilder {
	b.v.Codes = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *EnvelopeBuilder) Mutate(fn func(*Envelope)) *EnvelopeBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *EnvelopeBuilder) Clone() *EnvelopeBuilder {
	out := &EnvelopeBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [EnvelopeBuilder.Clone] returns a builder that shares none of them.
func (b *EnvelopeBuilder) Build() Envelope {
	return b.v
}

// HiddenBuilder constructs [Hidden] values. [NewHiddenBuilder] and
// [NewHiddenBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [HiddenBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type HiddenBuilder struct {
	v Hidden
}

// NewHiddenBuilder returns a builder of the zero value.
func NewHiddenBuilder() *HiddenBuilder {
	return &HiddenBuilder{}
}

// NewHiddenBuilderFrom returns a builder of v.
func NewHiddenBuilderFrom(v Hidden) *HiddenBuilder {
	return &HiddenBuilder{v: v}
}

// WithName replaces Name with v.
func (b *HiddenBuilder) WithName(v string) *HiddenBuilder {
	b.v.Name = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *HiddenBuilder) Mutate(fn func(*Hidden)) *HiddenBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *HiddenBuilder) Clone() *HiddenBuilder {
	out := &HiddenBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [HiddenBuilder.Clone] returns a builder that shares none of them.
func (b *HiddenBuilder) Build() Hidden {
	return b.v
}

// OrderBuilder constructs [Order] values. [NewOrderBuilder] and
// [NewOrderBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [OrderBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type OrderBuilder struct {
	v Order
}

// NewOrderBuilder returns a builder of the zero value.
func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{}
}

// NewOrderBuilderFrom returns a builder of v.
func NewOrderBuilderFrom(v Order) *OrderBuilder {
	return &OrderBuilder{v: v}
}

// WithRemarks replaces the elements of Remarks with v. A call without
// arguments empties Remarks.
func (b *OrderBuilder) WithRemarks(v ...string) *OrderBuilder {
	b.v.Remarks = v
	return b
}

// AppendRemarks appends v to the elements of Remarks.
func (b *OrderBuilder) AppendRemarks(v ...string) *OrderBuilder {
	b.v.Remarks = append(b.v.Remarks, v...)
	return b
}

// WithCodes replaces Codes with v.
func (b *OrderBuilder) WithCodes(v Codes) *OrderBuilder {
	b.v.Codes = v
	return b
}

// WithFlags replaces Flags with v. The builder keeps v without a copy.
func (b *OrderBuilder) WithFlags(v map[string]struct{}) *OrderBuilder {
	b.v.Flags = v
	return b
}

// WithFlagsEntry adds the key k to Flags. It allocates the set where Flags
// is nil.
func (b *OrderBuilder) WithFlagsEntry(k string) *OrderBuilder {
	if b.v.Flags == nil {
		b.v.Flags = make(map[string]struct{})
	}
	b.v.Flags[k] = struct{}{}
	return b
}

// WithFlagsEntries adds each key of keys to Flags. It allocates the set
// where Flags is nil.
func (b *OrderBuilder) WithFlagsEntries(keys ...string) *OrderBuilder {
	if b.v.Flags == nil {
		b.v.Flags = make(map[string]struct{}, len(keys))
	}
	for _, k := range keys {
		b.v.Flags[k] = struct{}{}
	}
	return b
}

// WithMarker replaces Marker with v.
func (b *OrderBuilder) WithMarker(v struct{}) *OrderBuilder {
	b.v.Marker = v
	return b
}

// WithMarks replaces the elements of Marks with v. A call without
// arguments empties Marks.
func (b *OrderBuilder) WithMarks(v ...struct{}) *OrderBuilder {
	b.v.Marks = v
	return b
}

// AppendMarks appends v to the elements of Marks.
func (b *OrderBuilder) AppendMarks(v ...struct{}) *OrderBuilder {
	b.v.Marks = append(b.v.Marks, v...)
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *OrderBuilder) Mutate(fn func(*Order)) *OrderBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *OrderBuilder) Clone() *OrderBuilder {
	out := &OrderBuilder{v: b.v}
	out.v.Remarks = slices.Clone(b.v.Remarks)
	out.v.Flags = maps.Clone(b.v.Flags)
	out.v.Marks = slices.Clone(b.v.Marks)
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [OrderBuilder.Clone] returns a builder that shares none of them.
func (b *OrderBuilder) Build() Order {
	return b.v
}

// LooseBuilder constructs [Loose] values. [NewLooseBuilder] and
// [NewLooseBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [LooseBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type LooseBuilder struct {
	v Loose
}

// NewLooseBuilder returns a builder of the zero value.
func NewLooseBuilder() *LooseBuilder {
	return &LooseBuilder{}
}

// NewLooseBuilderFrom returns a builder of v.
func NewLooseBuilderFrom(v Loose) *LooseBuilder {
	return &LooseBuilder{v: v}
}

// WithItems replaces the elements of Items with v. A call without
// arguments empties Items.
func (b *LooseBuilder) WithItems(v ...Notifier) *LooseBuilder {
	b.v.Items = v
	return b
}

// AppendItems appends v to the elements of Items.
func (b *LooseBuilder) AppendItems(v ...Notifier) *LooseBuilder {
	b.v.Items = append(b.v.Items, v...)
	return b
}

// WithIndex replaces Index with v. The builder keeps v without a copy.
func (b *LooseBuilder) WithIndex(v map[string]Notifier) *LooseBuilder {
	b.v.Index = v
	return b
}

// WithIndexEntry sets the key k of Index to v. It allocates the map where
// Index is nil.
func (b *LooseBuilder) WithIndexEntry(k string, v Notifier) *LooseBuilder {
	if b.v.Index == nil {
		b.v.Index = make(map[string]Notifier)
	}
	b.v.Index[k] = v
	return b
}

// WithIndexEntries copies each entry of entries into Index. A key of
// entries replaces the value of that key in Index. It allocates the map
// where Index is nil.
func (b *LooseBuilder) WithIndexEntries(entries map[string]Notifier) *LooseBuilder {
	if b.v.Index == nil {
		b.v.Index = make(map[string]Notifier, len(entries))
	}
	maps.Copy(b.v.Index, entries)
	return b
}

// WithKeys replaces Keys with v. The builder keeps v without a copy.
func (b *LooseBuilder) WithKeys(v map[Notifier]struct{}) *LooseBuilder {
	b.v.Keys = v
	return b
}

// WithKeysEntry adds the key k to Keys. It allocates the set where Keys is
// nil.
func (b *LooseBuilder) WithKeysEntry(k Notifier) *LooseBuilder {
	if b.v.Keys == nil {
		b.v.Keys = make(map[Notifier]struct{})
	}
	b.v.Keys[k] = struct{}{}
	return b
}

// WithKeysEntries adds each key of keys to Keys. It allocates the set
// where Keys is nil.
func (b *LooseBuilder) WithKeysEntries(keys ...Notifier) *LooseBuilder {
	if b.v.Keys == nil {
		b.v.Keys = make(map[Notifier]struct{}, len(keys))
	}
	for _, k := range keys {
		b.v.Keys[k] = struct{}{}
	}
	return b
}

// WithExtra points Extra at a copy of v. [LooseBuilder.Mutate] can set
// Extra to nil or to another address.
func (b *LooseBuilder) WithExtra(v Notifier) *LooseBuilder {
	b.v.Extra = &v
	return b
}

// WithHook replaces Hook with v.
func (b *LooseBuilder) WithHook(v Notifier) *LooseBuilder {
	b.v.Hook = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *LooseBuilder) Mutate(fn func(*Loose)) *LooseBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *LooseBuilder) Clone() *LooseBuilder {
	out := &LooseBuilder{v: b.v}
	out.v.Items = slices.Clone(b.v.Items)
	out.v.Index = maps.Clone(b.v.Index)
	out.v.Keys = maps.Clone(b.v.Keys)
	return out
}

// Build returns the value under construction. The builder stays usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [LooseBuilder.Clone] returns a builder that shares none of them.
func (b *LooseBuilder) Build() Loose {
	return b.v
}
