package catalog

import (
	"time"
)

// AddressBuilder constructs [Address] values.
//
// Start with [NewAddressBuilder] or [NewAddressBuilderFrom], set what you
// care about, and call [AddressBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Address breaks every
// literal at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [AddressBuilder.Clone] is how one becomes two.
type AddressBuilder struct {
	v Address
}

// NewAddressBuilder returns a builder seeded with the declared defaults.
//
// Only the members that declare one. Everything else holds its zero value,
// so a caller that sets nothing gets a mixture rather than either extreme
// — which is what makes a check against a built value worth writing.
func NewAddressBuilder() *AddressBuilder {
	var v Address
	v.Street = "Dam 1"
	v.City = "Amsterdam"
	v.Timeout = time.Second
	return &AddressBuilder{v: v}
}

// NewAddressBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewAddressBuilderFrom(v Address) *AddressBuilder {
	return &AddressBuilder{v: v}
}

// WithStreet replaces Street.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *AddressBuilder) WithStreet(v string) *AddressBuilder {
	b.v.Street = v
	return b
}

// WithCity replaces City.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *AddressBuilder) WithCity(v string) *AddressBuilder {
	b.v.City = v
	return b
}

// WithFloor replaces Floor, taking the value and addressing it, so a
// caller need not hold a variable to point at.
//
// Use [AddressBuilder.Mutate] to clear the member, or to point it at an
// address that already exists.
func (b *AddressBuilder) WithFloor(v int) *AddressBuilder {
	b.v.Floor = &v
	return b
}

// WithLines replaces Lines with the elements given.
//
// Variadic, and replacing: whatever the collection held is discarded,
// including anything a constructor seeded. AppendLines is the counterpart
// that keeps it. Called with no arguments this empties the member rather
// than leaving it alone, which is how a caller clears a seeded default.
func (b *AddressBuilder) WithLines(v ...string) *AddressBuilder {
	b.v.Lines = v
	return b
}

// AppendLines adds to Lines, keeping what is already there.
//
// The counterpart to WithLines, which discards it. Both are offered
// because a caller extending a seeded default and a caller replacing it
// want opposite things, and one method cannot mean both.
func (b *AddressBuilder) AppendLines(v ...string) *AddressBuilder {
	b.v.Lines = append(b.v.Lines, v...)
	return b
}

// WithLabels replaces Labels wholesale.
//
// The whole mapping, not one entry: what the member held is discarded
// rather than merged into, and the map is kept rather than copied — so a
// caller that writes to it afterwards changes what was already built.
// WithLabelsEntry is what sets a single key.
func (b *AddressBuilder) WithLabels(v map[string]string) *AddressBuilder {
	b.v.Labels = v
	return b
}

// WithLabelsEntry sets one key of Labels, allocating on first use.
//
// A key already present takes the new value. The allocation is why this is
// offered at all: a caller writing to the member directly has to check for
// nil first, and the builder is the thing that knows whether a constructor
// seeded one.
func (b *AddressBuilder) WithLabelsEntry(k string, v string) *AddressBuilder {
	if b.v.Labels == nil {
		b.v.Labels = make(map[string]string)
	}
	b.v.Labels[k] = v
	return b
}

// WithLabelsEntries merges entries into Labels.
//
// Merges rather than replaces: a key already present takes the new value,
// and keys absent from entries are left alone. Use WithLabels to discard
// what is there.
func (b *AddressBuilder) WithLabelsEntries(entries map[string]string) *AddressBuilder {
	if b.v.Labels == nil {
		b.v.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.v.Labels[k] = v
	}
	return b
}

// WithTags replaces Tags wholesale.
//
// The whole collection, not one key: what the member held is discarded
// rather than merged into, and the map is kept rather than copied.
// WithTagsEntry is what adds to what is already there.
func (b *AddressBuilder) WithTags(v map[string]struct{}) *AddressBuilder {
	b.v.Tags = v
	return b
}

// WithTagsEntry adds one key to Tags, allocating on first use.
//
// A set entry is a key. Every value in the collection is the same one, so
// there is nothing for a caller to supply and no parameter for it — a
// setter asking would be asking for the one thing they cannot vary.
func (b *AddressBuilder) WithTagsEntry(k string) *AddressBuilder {
	if b.v.Tags == nil {
		b.v.Tags = make(map[string]struct{})
	}
	b.v.Tags[k] = struct{}{}
	return b
}

// WithTagsEntries adds each key to Tags, keeping those it holds.
//
// Additive like WithTagsEntry, so two calls union rather than the second
// replacing the first. A key already present is left as it stands: a set
// holds a key or it does not, and adding one twice is the same as adding
// it once.
func (b *AddressBuilder) WithTagsEntries(keys ...string) *AddressBuilder {
	if b.v.Tags == nil {
		b.v.Tags = make(map[string]struct{}, len(keys))
	}
	for _, k := range keys {
		b.v.Tags[k] = struct{}{}
	}
	return b
}

// WithRaw replaces Raw with the bytes given.
//
// The bytes are kept, not copied: the builder and every value built from
// it share the array the caller passed, so a caller that mutates it
// afterwards changes what was already built. WithRawString is the form
// that takes text and has no such tie.
func (b *AddressBuilder) WithRaw(v []byte) *AddressBuilder {
	b.v.Raw = v
	return b
}

// WithRawString replaces Raw with the bytes of s.
//
// The conversion is available for this element type and no other, and a
// caller holding text should not have to write it at every call site.
func (b *AddressBuilder) WithRawString(s string) *AddressBuilder {
	b.v.Raw = []byte(s)
	return b
}

// WithGrid replaces Grid.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *AddressBuilder) WithGrid(v [2]int) *AddressBuilder {
	b.v.Grid = v
	return b
}

// WithTimeout replaces Timeout.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *AddressBuilder) WithTimeout(v time.Duration) *AddressBuilder {
	b.v.Timeout = v
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
func (b *AddressBuilder) Mutate(fn func(*Address)) *AddressBuilder {
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
// Use [AddressBuilder.Mutate] to replace such a member outright rather
// than editing it in place.
func (b *AddressBuilder) Clone() *AddressBuilder {
	out := &AddressBuilder{v: b.v}
	if b.v.Lines != nil {
		out.v.Lines = append([]string(nil), b.v.Lines...)
	}
	if b.v.Labels != nil {
		out.v.Labels = make(map[string]string, len(b.v.Labels))
		for k, v := range b.v.Labels {
			out.v.Labels[k] = v
		}
	}
	if b.v.Tags != nil {
		out.v.Tags = make(map[string]struct{}, len(b.v.Tags))
		for k, v := range b.v.Tags {
			out.v.Tags[k] = v
		}
	}
	if b.v.Raw != nil {
		out.v.Raw = append([]byte(nil), b.v.Raw...)
	}
	return out
}

// Build returns the constructed [Address].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [AddressBuilder.Clone] is what separates them.
func (b *AddressBuilder) Build() Address {
	return b.v
}

// ProfileBuilder constructs [Profile] values.
//
// Start with [NewProfileBuilder] or [NewProfileBuilderFrom], set what you
// care about, and call [ProfileBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Profile breaks every
// literal at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [ProfileBuilder.Clone] is how one becomes two.
type ProfileBuilder struct {
	v Profile
}

// NewProfileBuilder returns a builder seeded with the declared defaults.
//
// Only the members that declare one. Everything else holds its zero value,
// so a caller that sets nothing gets a mixture rather than either extreme
// — which is what makes a check against a built value worth writing.
func NewProfileBuilder() *ProfileBuilder {
	var v Profile
	v.Active = false
	return &ProfileBuilder{v: v}
}

// NewProfileBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewProfileBuilderFrom(v Profile) *ProfileBuilder {
	return &ProfileBuilder{v: v}
}

// WithRole replaces Role.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *ProfileBuilder) WithRole(v Role) *ProfileBuilder {
	b.v.Role = v
	return b
}

// WithAudit replaces Audit.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *ProfileBuilder) WithAudit(v *Audit) *ProfileBuilder {
	b.v.Audit = v
	return b
}

// WithActive replaces Active.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *ProfileBuilder) WithActive(v bool) *ProfileBuilder {
	b.v.Active = v
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
func (b *ProfileBuilder) Mutate(fn func(*Profile)) *ProfileBuilder {
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
// Use [ProfileBuilder.Mutate] to replace such a member outright rather
// than editing it in place.
func (b *ProfileBuilder) Clone() *ProfileBuilder {
	out := &ProfileBuilder{v: b.v}
	return out
}

// Build returns the constructed [Profile].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [ProfileBuilder.Clone] is what separates them.
func (b *ProfileBuilder) Build() Profile {
	return b.v
}

// EnvelopeBuilder constructs [Envelope] values.
//
// Start with [NewEnvelopeBuilder] or [NewEnvelopeBuilderFrom], set what
// you care about, and call [EnvelopeBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Envelope breaks every
// literal at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [EnvelopeBuilder.Clone] is how one becomes two.
type EnvelopeBuilder struct {
	v Envelope
}

// NewEnvelopeBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewEnvelopeBuilder() *EnvelopeBuilder {
	return &EnvelopeBuilder{}
}

// NewEnvelopeBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewEnvelopeBuilderFrom(v Envelope) *EnvelopeBuilder {
	return &EnvelopeBuilder{v: v}
}

// WithNotifier replaces Notifier.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *EnvelopeBuilder) WithNotifier(v Notifier) *EnvelopeBuilder {
	b.v.Notifier = v
	return b
}

// WithCodes replaces Codes.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *EnvelopeBuilder) WithCodes(v Codes) *EnvelopeBuilder {
	b.v.Codes = v
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
func (b *EnvelopeBuilder) Mutate(fn func(*Envelope)) *EnvelopeBuilder {
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
// Use [EnvelopeBuilder.Mutate] to replace such a member outright rather
// than editing it in place.
func (b *EnvelopeBuilder) Clone() *EnvelopeBuilder {
	out := &EnvelopeBuilder{v: b.v}
	return out
}

// Build returns the constructed [Envelope].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [EnvelopeBuilder.Clone] is what separates them.
func (b *EnvelopeBuilder) Build() Envelope {
	return b.v
}

// HiddenBuilder constructs [Hidden] values.
//
// Start with [NewHiddenBuilder] or [NewHiddenBuilderFrom], set what you
// care about, and call [HiddenBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Hidden breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [HiddenBuilder.Clone] is how one becomes two.
type HiddenBuilder struct {
	v Hidden
}

// NewHiddenBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewHiddenBuilder() *HiddenBuilder {
	return &HiddenBuilder{}
}

// NewHiddenBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewHiddenBuilderFrom(v Hidden) *HiddenBuilder {
	return &HiddenBuilder{v: v}
}

// WithName replaces Name.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *HiddenBuilder) WithName(v string) *HiddenBuilder {
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
func (b *HiddenBuilder) Mutate(fn func(*Hidden)) *HiddenBuilder {
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
// Use [HiddenBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *HiddenBuilder) Clone() *HiddenBuilder {
	out := &HiddenBuilder{v: b.v}
	return out
}

// Build returns the constructed [Hidden].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [HiddenBuilder.Clone] is what separates them.
func (b *HiddenBuilder) Build() Hidden {
	return b.v
}

// OrderBuilder constructs [Order] values.
//
// Start with [NewOrderBuilder] or [NewOrderBuilderFrom], set what you care
// about, and call [OrderBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Order breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [OrderBuilder.Clone] is how one becomes two.
type OrderBuilder struct {
	v Order
}

// NewOrderBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{}
}

// NewOrderBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewOrderBuilderFrom(v Order) *OrderBuilder {
	return &OrderBuilder{v: v}
}

// WithRemarks replaces Remarks with the elements given.
//
// Variadic, and replacing: whatever the collection held is discarded,
// including anything a constructor seeded. AppendRemarks is the
// counterpart that keeps it. Called with no arguments this empties the
// member rather than leaving it alone, which is how a caller clears a
// seeded default.
func (b *OrderBuilder) WithRemarks(v ...string) *OrderBuilder {
	b.v.Remarks = v
	return b
}

// AppendRemarks adds to Remarks, keeping what is already there.
//
// The counterpart to WithRemarks, which discards it. Both are offered
// because a caller extending a seeded default and a caller replacing it
// want opposite things, and one method cannot mean both.
func (b *OrderBuilder) AppendRemarks(v ...string) *OrderBuilder {
	b.v.Remarks = append(b.v.Remarks, v...)
	return b
}

// WithCodes replaces Codes.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *OrderBuilder) WithCodes(v Codes) *OrderBuilder {
	b.v.Codes = v
	return b
}

// WithFlags replaces Flags wholesale.
//
// The whole collection, not one key: what the member held is discarded
// rather than merged into, and the map is kept rather than copied.
// WithFlagsEntry is what adds to what is already there.
func (b *OrderBuilder) WithFlags(v map[string]struct{}) *OrderBuilder {
	b.v.Flags = v
	return b
}

// WithFlagsEntry adds one key to Flags, allocating on first use.
//
// A set entry is a key. Every value in the collection is the same one, so
// there is nothing for a caller to supply and no parameter for it — a
// setter asking would be asking for the one thing they cannot vary.
func (b *OrderBuilder) WithFlagsEntry(k string) *OrderBuilder {
	if b.v.Flags == nil {
		b.v.Flags = make(map[string]struct{})
	}
	b.v.Flags[k] = struct{}{}
	return b
}

// WithFlagsEntries adds each key to Flags, keeping those it holds.
//
// Additive like WithFlagsEntry, so two calls union rather than the second
// replacing the first. A key already present is left as it stands: a set
// holds a key or it does not, and adding one twice is the same as adding
// it once.
func (b *OrderBuilder) WithFlagsEntries(keys ...string) *OrderBuilder {
	if b.v.Flags == nil {
		b.v.Flags = make(map[string]struct{}, len(keys))
	}
	for _, k := range keys {
		b.v.Flags[k] = struct{}{}
	}
	return b
}

// WithMarker replaces Marker.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *OrderBuilder) WithMarker(v struct{}) *OrderBuilder {
	b.v.Marker = v
	return b
}

// WithMarks replaces Marks with the elements given.
//
// Variadic, and replacing: whatever the collection held is discarded,
// including anything a constructor seeded. AppendMarks is the counterpart
// that keeps it. Called with no arguments this empties the member rather
// than leaving it alone, which is how a caller clears a seeded default.
func (b *OrderBuilder) WithMarks(v ...struct{}) *OrderBuilder {
	b.v.Marks = v
	return b
}

// AppendMarks adds to Marks, keeping what is already there.
//
// The counterpart to WithMarks, which discards it. Both are offered
// because a caller extending a seeded default and a caller replacing it
// want opposite things, and one method cannot mean both.
func (b *OrderBuilder) AppendMarks(v ...struct{}) *OrderBuilder {
	b.v.Marks = append(b.v.Marks, v...)
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
func (b *OrderBuilder) Mutate(fn func(*Order)) *OrderBuilder {
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
// Use [OrderBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *OrderBuilder) Clone() *OrderBuilder {
	out := &OrderBuilder{v: b.v}
	if b.v.Remarks != nil {
		out.v.Remarks = append(Remarks(nil), b.v.Remarks...)
	}
	if b.v.Flags != nil {
		out.v.Flags = make(map[string]struct{}, len(b.v.Flags))
		for k, v := range b.v.Flags {
			out.v.Flags[k] = v
		}
	}
	if b.v.Marks != nil {
		out.v.Marks = append([]struct{}(nil), b.v.Marks...)
	}
	return out
}

// Build returns the constructed [Order].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [OrderBuilder.Clone] is what separates them.
func (b *OrderBuilder) Build() Order {
	return b.v
}

// LooseBuilder constructs [Loose] values.
//
// Start with [NewLooseBuilder] or [NewLooseBuilderFrom], set what you care
// about, and call [LooseBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Loose breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [LooseBuilder.Clone] is how one becomes two.
type LooseBuilder struct {
	v Loose
}

// NewLooseBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewLooseBuilder() *LooseBuilder {
	return &LooseBuilder{}
}

// NewLooseBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewLooseBuilderFrom(v Loose) *LooseBuilder {
	return &LooseBuilder{v: v}
}

// WithItems replaces Items with the elements given.
//
// Variadic, and replacing: whatever the collection held is discarded,
// including anything a constructor seeded. AppendItems is the counterpart
// that keeps it. Called with no arguments this empties the member rather
// than leaving it alone, which is how a caller clears a seeded default.
func (b *LooseBuilder) WithItems(v ...Notifier) *LooseBuilder {
	b.v.Items = v
	return b
}

// AppendItems adds to Items, keeping what is already there.
//
// The counterpart to WithItems, which discards it. Both are offered
// because a caller extending a seeded default and a caller replacing it
// want opposite things, and one method cannot mean both.
func (b *LooseBuilder) AppendItems(v ...Notifier) *LooseBuilder {
	b.v.Items = append(b.v.Items, v...)
	return b
}

// WithIndex replaces Index wholesale.
//
// The whole mapping, not one entry: what the member held is discarded
// rather than merged into, and the map is kept rather than copied — so a
// caller that writes to it afterwards changes what was already built.
// WithIndexEntry is what sets a single key.
func (b *LooseBuilder) WithIndex(v map[string]Notifier) *LooseBuilder {
	b.v.Index = v
	return b
}

// WithIndexEntry sets one key of Index, allocating on first use.
//
// A key already present takes the new value. The allocation is why this is
// offered at all: a caller writing to the member directly has to check for
// nil first, and the builder is the thing that knows whether a constructor
// seeded one.
func (b *LooseBuilder) WithIndexEntry(k string, v Notifier) *LooseBuilder {
	if b.v.Index == nil {
		b.v.Index = make(map[string]Notifier)
	}
	b.v.Index[k] = v
	return b
}

// WithIndexEntries merges entries into Index.
//
// Merges rather than replaces: a key already present takes the new value,
// and keys absent from entries are left alone. Use WithIndex to discard
// what is there.
func (b *LooseBuilder) WithIndexEntries(entries map[string]Notifier) *LooseBuilder {
	if b.v.Index == nil {
		b.v.Index = make(map[string]Notifier, len(entries))
	}
	for k, v := range entries {
		b.v.Index[k] = v
	}
	return b
}

// WithKeys replaces Keys wholesale.
//
// The whole collection, not one key: what the member held is discarded
// rather than merged into, and the map is kept rather than copied.
// WithKeysEntry is what adds to what is already there.
func (b *LooseBuilder) WithKeys(v map[Notifier]struct{}) *LooseBuilder {
	b.v.Keys = v
	return b
}

// WithKeysEntry adds one key to Keys, allocating on first use.
//
// A set entry is a key. Every value in the collection is the same one, so
// there is nothing for a caller to supply and no parameter for it — a
// setter asking would be asking for the one thing they cannot vary.
func (b *LooseBuilder) WithKeysEntry(k Notifier) *LooseBuilder {
	if b.v.Keys == nil {
		b.v.Keys = make(map[Notifier]struct{})
	}
	b.v.Keys[k] = struct{}{}
	return b
}

// WithKeysEntries adds each key to Keys, keeping those it holds.
//
// Additive like WithKeysEntry, so two calls union rather than the second
// replacing the first. A key already present is left as it stands: a set
// holds a key or it does not, and adding one twice is the same as adding
// it once.
func (b *LooseBuilder) WithKeysEntries(keys ...Notifier) *LooseBuilder {
	if b.v.Keys == nil {
		b.v.Keys = make(map[Notifier]struct{}, len(keys))
	}
	for _, k := range keys {
		b.v.Keys[k] = struct{}{}
	}
	return b
}

// WithExtra replaces Extra, taking the value and addressing it, so a
// caller need not hold a variable to point at.
//
// Use [LooseBuilder.Mutate] to clear the member, or to point it at an
// address that already exists.
func (b *LooseBuilder) WithExtra(v Notifier) *LooseBuilder {
	b.v.Extra = &v
	return b
}

// WithHook replaces Hook.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *LooseBuilder) WithHook(v Notifier) *LooseBuilder {
	b.v.Hook = v
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
func (b *LooseBuilder) Mutate(fn func(*Loose)) *LooseBuilder {
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
// Use [LooseBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *LooseBuilder) Clone() *LooseBuilder {
	out := &LooseBuilder{v: b.v}
	if b.v.Items != nil {
		out.v.Items = append([]Notifier(nil), b.v.Items...)
	}
	if b.v.Index != nil {
		out.v.Index = make(map[string]Notifier, len(b.v.Index))
		for k, v := range b.v.Index {
			out.v.Index[k] = v
		}
	}
	if b.v.Keys != nil {
		out.v.Keys = make(map[Notifier]struct{}, len(b.v.Keys))
		for k, v := range b.v.Keys {
			out.v.Keys[k] = v
		}
	}
	return out
}

// Build returns the constructed [Loose].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [LooseBuilder.Clone] is what separates them.
func (b *LooseBuilder) Build() Loose {
	return b.v
}
