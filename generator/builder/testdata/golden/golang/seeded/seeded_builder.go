package seeded

import (
	"runtime"
)

// AccountBuilder constructs [Account] values. [NewAccountBuilder] and
// [NewAccountBuilderFrom] return a new builder. Each setter sets one
// member and returns the builder, so the calls chain.
// [AccountBuilder.Build] returns the value. A builder is not safe for
// concurrent use.
type AccountBuilder struct {
	v Account
}

// NewAccountBuilder returns a builder of the value that the companion
// function returns. The declared default of each member applies over that
// value.
func NewAccountBuilder() *AccountBuilder {
	v := AccountDefaults()
	v.Owner = "owner"
	return &AccountBuilder{v: v}
}

// NewAccountBuilderFrom returns a builder of v.
func NewAccountBuilderFrom(v Account) *AccountBuilder {
	return &AccountBuilder{v: v}
}

// WithOwner replaces Owner with v.
func (b *AccountBuilder) WithOwner(v string) *AccountBuilder {
	b.v.Owner = v
	return b
}

// WithPlan replaces Plan with v.
func (b *AccountBuilder) WithPlan(v string) *AccountBuilder {
	b.v.Plan = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *AccountBuilder) Mutate(fn func(*Account)) *AccountBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *AccountBuilder) Clone() *AccountBuilder {
	out := &AccountBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [AccountBuilder.Clone] returns a builder that shares none of them.
func (b *AccountBuilder) Build() Account {
	return b.v
}

// PolicyBuilder constructs [Policy] values. [NewPolicyBuilder] and
// [NewPolicyBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [PolicyBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type PolicyBuilder struct {
	v Policy
}

// NewPolicyBuilder returns a builder of the value that the companion
// function returns. The declared default of each member applies over that
// value.
func NewPolicyBuilder() *PolicyBuilder {
	v := defaultPolicy()
	return &PolicyBuilder{v: v}
}

// NewPolicyBuilderFrom returns a builder of v.
func NewPolicyBuilderFrom(v Policy) *PolicyBuilder {
	return &PolicyBuilder{v: v}
}

// WithAttempts replaces Attempts with v.
func (b *PolicyBuilder) WithAttempts(v int) *PolicyBuilder {
	b.v.Attempts = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *PolicyBuilder) Mutate(fn func(*Policy)) *PolicyBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *PolicyBuilder) Clone() *PolicyBuilder {
	out := &PolicyBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [PolicyBuilder.Clone] returns a builder that shares none of them.
func (b *PolicyBuilder) Build() Policy {
	return b.v
}

// TierBuilder constructs [Tier] values. [NewTierBuilder] and
// [NewTierBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [TierBuilder.Build] returns
// the value. A builder is not safe for concurrent use.
type TierBuilder struct {
	v Tier
}

// NewTierBuilder returns a builder of the zero value.
func NewTierBuilder() *TierBuilder {
	return &TierBuilder{}
}

// NewTierBuilderFrom returns a builder of v.
func NewTierBuilderFrom(v Tier) *TierBuilder {
	return &TierBuilder{v: v}
}

// WithLevel replaces Level with v.
func (b *TierBuilder) WithLevel(v int) *TierBuilder {
	b.v.Level = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *TierBuilder) Mutate(fn func(*Tier)) *TierBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *TierBuilder) Clone() *TierBuilder {
	out := &TierBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [TierBuilder.Clone] returns a builder that shares none of them.
func (b *TierBuilder) Build() Tier {
	return b.v
}

// DeploymentBuilder constructs [Deployment] values. [NewDeploymentBuilder]
// and [NewDeploymentBuilderFrom] return a new builder. Each setter sets
// one member and returns the builder, so the calls chain.
// [DeploymentBuilder.Build] returns the value. A builder is not safe for
// concurrent use.
type DeploymentBuilder struct {
	v Deployment
}

// NewDeploymentBuilder returns a builder of the declared defaults. A
// member without a declared default has its zero value.
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

// NewDeploymentBuilderFrom returns a builder of v.
func NewDeploymentBuilderFrom(v Deployment) *DeploymentBuilder {
	return &DeploymentBuilder{v: v}
}

// WithMode replaces Mode with v.
func (b *DeploymentBuilder) WithMode(v Code) *DeploymentBuilder {
	b.v.Mode = v
	return b
}

// WithRegion replaces Region with v.
func (b *DeploymentBuilder) WithRegion(v string) *DeploymentBuilder {
	b.v.Region = v
	return b
}

// WithZone replaces Zone with v.
func (b *DeploymentBuilder) WithZone(v string) *DeploymentBuilder {
	b.v.Zone = v
	return b
}

// WithSystem replaces System with v.
func (b *DeploymentBuilder) WithSystem(v string) *DeploymentBuilder {
	b.v.System = v
	return b
}

// WithName replaces Name with v.
func (b *DeploymentBuilder) WithName(v string) *DeploymentBuilder {
	b.v.Name = v
	return b
}

// WithLabel replaces Label with v.
func (b *DeploymentBuilder) WithLabel(v string) *DeploymentBuilder {
	b.v.Label = v
	return b
}

// WithFlag replaces Flag with v.
func (b *DeploymentBuilder) WithFlag(v string) *DeploymentBuilder {
	b.v.Flag = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *DeploymentBuilder) Mutate(fn func(*Deployment)) *DeploymentBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *DeploymentBuilder) Clone() *DeploymentBuilder {
	out := &DeploymentBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [DeploymentBuilder.Clone] returns a builder that shares none of them.
func (b *DeploymentBuilder) Build() Deployment {
	return b.v
}
