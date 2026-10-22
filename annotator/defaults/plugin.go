// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package defaults

import (
	"errors"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// ID is the name of the plugin.
const ID plugin.ID = "defaults"

// Version is the version of the plugin. The fingerprint of a composition
// folds it in.
const Version = "0.1.0"

// Capability is the capability that defaults provides. A plugin that reads
// the declared defaults requires it, so it runs after defaults.
const Capability plugin.Capability = "defaults"

// DefaultDirective is the directive that declares the default value of a
// field. ValueParam is its one positional param, which is the value.
const (
	DefaultDirective directive.Name     = "default"
	ValueParam       directive.ParamKey = "value"
)

// KeyValue is the key of the fact that contains the text of the default
// value of a field. KeySource is the key of the fact that contains the
// form that declares the value: [SourceDirective] or [SourceTag].
const (
	KeyValue  meta.KeyName = "defaults.value"
	KeySource meta.KeyName = "defaults.source"
)

// SourceDirective and SourceTag are the values of [KeySource]. The
// directive declares the default of a field in every language, and a tag
// declares it in a language whose rules state tags.
const (
	SourceDirective = "directive"
	SourceTag       = "tag"
)

// DefaultTag is the default of the option tag-key: the key of the tag that
// contains the default value of a field.
const DefaultTag = "default"

// schema is the directive of defaults.
var schema = directive.Schema{
	Name: DefaultDirective,
	Positional: []directive.ParamSpec{{
		Key: ValueParam, Type: directive.TypeString, Required: true,
		Doc: "the default value of the field or the type, as text that the rules of its language read",
	}},
	Doc: "declares the default value of the field, or of each value of the type",
}

// Options are the options of defaults. TagKey is the key of the tag that the
// tag rule reads. NoTags turns the tag rule off, so the directive is the only
// form of a default.
type Options struct {
	TagKey string `opt:"tag-key" doc:"the key of the field tag that contains the default value of the field"`
	NoTags bool   `opt:"no-tags" doc:"reads no field tag, so the directive is the only form of a default"`
}

// conclusion is the state of one instance of defaults: its options, and
// the handles of its keys, which the registration sets and the rules stamp
// through.
//
// # Concurrency
//
// The composition sets the options and registers the keys before the first
// phase call, and the handlers read them afterwards, so the workers of a
// phase call read them concurrently.
type conclusion struct {
	options *Options
	value   meta.Key[string]
	source  meta.Key[string]
}

// New returns defaults, an annotator. It stamps the text of the default
// value of a field or a type under [KeyValue], and the form that declares it
// under [KeySource]. The directive rule reads [DefaultDirective] on a field
// and on a type alias or a defined type, in every language. A type has no
// tag, so the directive is the only form of its default. The tag rule reads
// the tag of the key that the option tag-key sets, in each language whose
// rules state tags, unless the option no-tags is set. Within one plugin, the
// first claim in rule order takes precedence, and the directive rule comes
// first. So on a field with both forms, the value of the directive becomes
// the fact. A reader turns the text into a value with [Value]. Each call
// returns a new plugin, because a plugin instance belongs to one workspace.
func New() plugin.Annotator {
	c := &conclusion{options: &Options{TagKey: DefaultTag}}
	return sdk.NewPlugin(ID).
		Version(Version).
		Provides(Capability).
		Options(c.options).
		Keys(c.register).
		Handle(
			sdk.Directive(schema, sdk.OnField(c.declaredField), sdk.OnAlias(c.declaredType)),
			sdk.OnField(c.tagged),
		).
		Annotator()
}

// register claims the namespace of defaults and registers [KeyValue] and
// [KeySource].
//
// # Errors
//
// register returns the error of the claim, such as a namespace that another
// registrant claimed, joined with the errors of the two registrations, such as
// a key that the registrant registered with another type.
func (c *conclusion) register(r *meta.Registry) error {
	errClaim := r.ClaimNamespace(KeyValue.Namespace())
	var errValue, errSource error
	c.value, errValue = meta.Register[string](r, meta.KeySpec{
		Name:  KeyValue,
		Kinds: []symbol.Kind{symbol.KindField, symbol.KindAlias},
		Doc:   "contains the text of the default value of a field or a type",
	})
	c.source, errSource = meta.Register[string](r, meta.KeySpec{
		Name:  KeySource,
		Kinds: []symbol.Kind{symbol.KindField, symbol.KindAlias},
		Doc:   "contains the form that declares the default value of a field or a type: directive or tag",
	})
	return errors.Join(errClaim, errValue, errSource)
}

// declaredField stamps the value of the directive on its field.
func (c *conclusion) declaredField(m *sdk.FieldMatch, st *sdk.Stamper) error {
	c.declared(m.Directive(), st)
	return nil
}

// declaredType stamps the value of the directive on its type alias or its
// defined type.
func (c *conclusion) declaredType(m *sdk.AliasMatch, st *sdk.Stamper) error {
	c.declared(m.Directive(), st)
	return nil
}

// declared stamps the value of the instance d with the source
// [SourceDirective]. The schema requires the one positional value, so
// validation passes no instance without it.
func (c *conclusion) declared(d *directive.Directive, st *sdk.Stamper) {
	sdk.Stamp(st, c.value, d.Args[0].Str)
	sdk.Stamp(st, c.source, SourceDirective)
}

// tagged stamps the value of the tag of the field under the key of the
// option tag-key, with the source [SourceTag]. It does not stamp a field
// without the tag, a field whose tag has an empty value, a field of a
// language whose rules state no tags, or any field under the option
// no-tags. An empty value declares no default, because a reader would assign
// the zero value over a value that a seed set.
func (c *conclusion) tagged(m *sdk.FieldMatch, st *sdk.Stamper) error {
	if c.options.NoTags {
		return nil
	}
	if text, _ := m.Rules().Tag(m.Field, c.options.TagKey); text != "" {
		sdk.Stamp(st, c.value, text)
		sdk.Stamp(st, c.source, SourceTag)
	}
	return nil
}
