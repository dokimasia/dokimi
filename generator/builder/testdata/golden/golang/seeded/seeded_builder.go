package seeded

import (
	"runtime"
)

// AccountBuilder constructs [Account] values.
//
// Start with [NewAccountBuilder] or [NewAccountBuilderFrom], set what you
// care about, and call [AccountBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Account breaks every
// literal at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [AccountBuilder.Clone] is how one becomes two.
type AccountBuilder struct {
	v Account
}

// NewAccountBuilder returns a builder seeded from the package's companion,
// with any declared per-member default applied over it.
//
// Two layers, in that order: whatever the companion supplies, then each
// member that declares its own default. A member the companion sets and
// the declaration does not keeps the companion's value, so the two sources
// compose rather than one replacing the other.
func NewAccountBuilder() *AccountBuilder {
	v := AccountDefaults()
	v.Owner = "owner"
	return &AccountBuilder{v: v}
}

// NewAccountBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewAccountBuilderFrom(v Account) *AccountBuilder {
	return &AccountBuilder{v: v}
}

// WithOwner replaces Owner.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *AccountBuilder) WithOwner(v string) *AccountBuilder {
	b.v.Owner = v
	return b
}

// WithPlan replaces Plan.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *AccountBuilder) WithPlan(v string) *AccountBuilder {
	b.v.Plan = v
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
func (b *AccountBuilder) Mutate(fn func(*Account)) *AccountBuilder {
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
// Use [AccountBuilder.Mutate] to replace such a member outright rather
// than editing it in place.
func (b *AccountBuilder) Clone() *AccountBuilder {
	out := &AccountBuilder{v: b.v}
	return out
}

// Build returns the constructed [Account].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [AccountBuilder.Clone] is what separates them.
func (b *AccountBuilder) Build() Account {
	return b.v
}

// PolicyBuilder constructs [Policy] values.
//
// Start with [NewPolicyBuilder] or [NewPolicyBuilderFrom], set what you
// care about, and call [PolicyBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Policy breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [PolicyBuilder.Clone] is how one becomes two.
type PolicyBuilder struct {
	v Policy
}

// NewPolicyBuilder returns a builder seeded from the package's companion,
// with any declared per-member default applied over it.
//
// Two layers, in that order: whatever the companion supplies, then each
// member that declares its own default. A member the companion sets and
// the declaration does not keeps the companion's value, so the two sources
// compose rather than one replacing the other.
func NewPolicyBuilder() *PolicyBuilder {
	v := defaultPolicy()
	return &PolicyBuilder{v: v}
}

// NewPolicyBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewPolicyBuilderFrom(v Policy) *PolicyBuilder {
	return &PolicyBuilder{v: v}
}

// WithAttempts replaces Attempts.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *PolicyBuilder) WithAttempts(v int) *PolicyBuilder {
	b.v.Attempts = v
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
func (b *PolicyBuilder) Mutate(fn func(*Policy)) *PolicyBuilder {
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
// Use [PolicyBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *PolicyBuilder) Clone() *PolicyBuilder {
	out := &PolicyBuilder{v: b.v}
	return out
}

// Build returns the constructed [Policy].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [PolicyBuilder.Clone] is what separates them.
func (b *PolicyBuilder) Build() Policy {
	return b.v
}

// TierBuilder constructs [Tier] values.
//
// Start with [NewTierBuilder] or [NewTierBuilderFrom], set what you care
// about, and call [TierBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Tier breaks every literal
// at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [TierBuilder.Clone] is how one becomes two.
type TierBuilder struct {
	v Tier
}

// NewTierBuilder returns a builder holding the zero value.
//
// No member declares a default, so there is nothing to seed. A caller that
// builds immediately gets the same value a composite literal with no
// fields would, and the setters are what make it differ.
func NewTierBuilder() *TierBuilder {
	return &TierBuilder{}
}

// NewTierBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewTierBuilderFrom(v Tier) *TierBuilder {
	return &TierBuilder{v: v}
}

// WithLevel replaces Level.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *TierBuilder) WithLevel(v int) *TierBuilder {
	b.v.Level = v
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
func (b *TierBuilder) Mutate(fn func(*Tier)) *TierBuilder {
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
// Use [TierBuilder.Mutate] to replace such a member outright rather than
// editing it in place.
func (b *TierBuilder) Clone() *TierBuilder {
	out := &TierBuilder{v: b.v}
	return out
}

// Build returns the constructed [Tier].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [TierBuilder.Clone] is what separates them.
func (b *TierBuilder) Build() Tier {
	return b.v
}

// DeploymentBuilder constructs [Deployment] values.
//
// Start with [NewDeploymentBuilder] or [NewDeploymentBuilderFrom], set
// what you care about, and call [DeploymentBuilder.Build].
//
// The point is what you do not have to write. A composite literal restates
// every member each time, so a member added to Deployment breaks every
// literal at once and a reader cannot tell which members a given call site
// actually cares about. Here the constructor supplies the rest, and each
// call states only what it varies.
//
// Every setter returns the builder, so calls chain. A builder is a value
// under construction and is not safe for concurrent use;
// [DeploymentBuilder.Clone] is how one becomes two.
type DeploymentBuilder struct {
	v Deployment
}

// NewDeploymentBuilder returns a builder seeded with the declared
// defaults.
//
// Only the members that declare one. Everything else holds its zero value,
// so a caller that sets nothing gets a mixture rather than either extreme
// — which is what makes a check against a built value worth writing.
func NewDeploymentBuilder() *DeploymentBuilder {
	var v Deployment
	v.Mode = "true"
	v.Region = Region
	v.Zone = zone()
	v.System = runtime.GOOS
	v.Name = "web"
	v.Label = "a b"
	v.Flag = "true"
	return &DeploymentBuilder{v: v}
}

// NewDeploymentBuilderFrom returns a builder seeded with v.
//
// For the case where a caller already holds a value and wants one member
// different — a round trip, a near-copy, a regression fixture built from
// something a previous step produced.
func NewDeploymentBuilderFrom(v Deployment) *DeploymentBuilder {
	return &DeploymentBuilder{v: v}
}

// WithMode replaces Mode.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithMode(v Code) *DeploymentBuilder {
	b.v.Mode = v
	return b
}

// WithRegion replaces Region.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithRegion(v string) *DeploymentBuilder {
	b.v.Region = v
	return b
}

// WithZone replaces Zone.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithZone(v string) *DeploymentBuilder {
	b.v.Zone = v
	return b
}

// WithSystem replaces System.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithSystem(v string) *DeploymentBuilder {
	b.v.System = v
	return b
}

// WithName replaces Name.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithName(v string) *DeploymentBuilder {
	b.v.Name = v
	return b
}

// WithLabel replaces Label.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithLabel(v string) *DeploymentBuilder {
	b.v.Label = v
	return b
}

// WithFlag replaces Flag.
//
// The later call wins: a setter overwrites rather than merges, so a caller
// building from a constructor that seeded this member states only what it
// varies.
func (b *DeploymentBuilder) WithFlag(v string) *DeploymentBuilder {
	b.v.Flag = v
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
func (b *DeploymentBuilder) Mutate(fn func(*Deployment)) *DeploymentBuilder {
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
// Use [DeploymentBuilder.Mutate] to replace such a member outright rather
// than editing it in place.
func (b *DeploymentBuilder) Clone() *DeploymentBuilder {
	out := &DeploymentBuilder{v: b.v}
	return out
}

// Build returns the constructed [Deployment].
//
// The builder stays usable afterwards, and a second call returns the value
// as it stands then — so a caller may build, set one more member, and
// build again. What it returns is a copy of the value under construction,
// but only as far as Go copies: a member holding a slice or a map is
// shared with the builder and with every value built from it, and
// [DeploymentBuilder.Clone] is what separates them.
func (b *DeploymentBuilder) Build() Deployment {
	return b.v
}
