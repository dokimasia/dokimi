package hidden

import (
	"slices"
)

// SiteBuilder constructs [Site] values. [NewSiteBuilder] and
// [NewSiteBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [SiteBuilder.Build] returns
// the value. A builder is not safe for concurrent use.
type SiteBuilder struct {
	v Site
}

// NewSiteBuilder returns a builder of the zero value.
func NewSiteBuilder() *SiteBuilder {
	return &SiteBuilder{}
}

// NewSiteBuilderFrom returns a builder of v.
func NewSiteBuilderFrom(v Site) *SiteBuilder {
	return &SiteBuilder{v: v}
}

// WithHost replaces Host with v.
func (b *SiteBuilder) WithHost(v Host) *SiteBuilder {
	b.v.Host = v
	return b
}

// WithZone replaces Zone with v.
func (b *SiteBuilder) WithZone(v Zone) *SiteBuilder {
	b.v.Zone = v
	return b
}

// WithRegion replaces Region with v.
func (b *SiteBuilder) WithRegion(v region) *SiteBuilder {
	b.v.Region = v
	return b
}

// WithHosts replaces the elements of Hosts with v. A call without
// arguments empties Hosts.
func (b *SiteBuilder) WithHosts(v ...Host) *SiteBuilder {
	b.v.Hosts = v
	return b
}

// AppendHosts appends v to the elements of Hosts.
func (b *SiteBuilder) AppendHosts(v ...Host) *SiteBuilder {
	b.v.Hosts = append(b.v.Hosts, v...)
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *SiteBuilder) Mutate(fn func(*Site)) *SiteBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *SiteBuilder) Clone() *SiteBuilder {
	out := &SiteBuilder{v: b.v}
	out.v.Hosts = slices.Clone(b.v.Hosts)
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [SiteBuilder.Clone] returns a builder that shares none of them.
func (b *SiteBuilder) Build() Site {
	return b.v
}

// PoolBuilder constructs [Pool] values. [NewPoolBuilder] and
// [NewPoolBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [PoolBuilder.Build] returns
// the value. A builder is not safe for concurrent use.
type PoolBuilder[T any] struct {
	v Pool[T]
}

// NewPoolBuilder returns a builder of the zero value.
func NewPoolBuilder[T any]() *PoolBuilder[T] {
	return &PoolBuilder[T]{}
}

// NewPoolBuilderFrom returns a builder of v.
func NewPoolBuilderFrom[T any](v Pool[T]) *PoolBuilder[T] {
	return &PoolBuilder[T]{v: v}
}

// WithItems replaces the elements of Items with v. A call without
// arguments empties Items.
func (b *PoolBuilder[T]) WithItems(v ...T) *PoolBuilder[T] {
	b.v.Items = v
	return b
}

// AppendItems appends v to the elements of Items.
func (b *PoolBuilder[T]) AppendItems(v ...T) *PoolBuilder[T] {
	b.v.Items = append(b.v.Items, v...)
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *PoolBuilder[T]) Mutate(fn func(*Pool[T])) *PoolBuilder[T] {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *PoolBuilder[T]) Clone() *PoolBuilder[T] {
	out := &PoolBuilder[T]{v: b.v}
	out.v.Items = slices.Clone(b.v.Items)
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [PoolBuilder.Clone] returns a builder that shares none of them.
func (b *PoolBuilder[T]) Build() Pool[T] {
	return b.v
}

// configBuilder constructs [config] values. [newConfigBuilder] and
// [newConfigBuilderFrom] return a new builder. Each setter sets one member
// and returns the builder, so the calls chain. [configBuilder.Build]
// returns the value. A builder is not safe for concurrent use.
type configBuilder struct {
	v config
}

// newConfigBuilder returns a builder of the declared defaults. A member
// without a declared default has its zero value.
func newConfigBuilder() *configBuilder {
	var v config
	v.Name = "web"
	return &configBuilder{v: v}
}

// newConfigBuilderFrom returns a builder of v.
func newConfigBuilderFrom(v config) *configBuilder {
	return &configBuilder{v: v}
}

// WithName replaces Name with v.
func (b *configBuilder) WithName(v string) *configBuilder {
	b.v.Name = v
	return b
}

// Mutate applies fn to the value under construction. It changes what no
// setter changes, such as several members in one step.
func (b *configBuilder) Mutate(fn func(*config)) *configBuilder {
	fn(&b.v)
	return b
}

// Clone returns a builder of a copy of the value under construction. The copy
// has its own lists, maps, sets and bytes, so a change of one builder leaves
// the other unchanged. The two builders share the elements of those members
// and the value of each optional member.
func (b *configBuilder) Clone() *configBuilder {
	out := &configBuilder{v: b.v}
	return out
}

// Build returns the value under construction. The builder remains usable,
// and a later call returns the value with the changes since. The value
// shares its lists, maps, sets and bytes with the builder, and
// [configBuilder.Clone] returns a builder that shares none of them.
func (b *configBuilder) Build() config {
	return b.v
}
