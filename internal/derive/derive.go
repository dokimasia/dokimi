// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/spec"
)

// ErrKey is the error of [Derive] for two callables of a declaration whose
// names have one key, such as ReadAll and Read_All.
var ErrKey = errors.New("derive: two callables of the declaration have one key")

// Facts are the facts that eidos derives for one callable of a declaration.
type Facts struct {
	// Declaration and Callable are the names that the source declares.
	// Declaration is empty for a function, which is its own declaration.
	Declaration, Callable string

	// ID is the identity of the callable, by which a param of a
	// classification names it.
	ID symbol.Identity

	// Class is the classification of the callable.
	Class shape.Classification

	// Inputs pairs each parameter in the role input with what eidos derives
	// for its type. Results does the same for each return in the role value.
	Inputs, Results []Value

	// Signature is the normalized view of the callable.
	Signature rules.Callable
}

// Value is one parameter or return with what eidos derives for its type.
type Value struct {
	// Zero is the zero value of the type, where HasZero is true.
	Zero emit.Value

	Name string

	// Fields are the fields of the type where it is a struct, with the
	// shape of each.
	Fields []Field

	// Samples are two distinct samples of the type, or the reasons why
	// eidos refuses them.
	Samples [2]rules.Sample

	Shape rules.TypeShape

	HasZero bool
}

// Field is one field of a struct type.
type Field struct {
	Name  string
	Shape rules.TypeShape
}

// Key returns the key of the callable of f, which joins it across languages:
// store.read-all for the method ReadAll of Store, and slug for a function.
func (f *Facts) Key() string {
	return spec.Key(f.qualified())
}

// qualified returns the name of the callable of f in an identity: Store.Read
// for a method, and Slug for a function.
func (f *Facts) qualified() string {
	if f.Declaration == "" {
		return f.Callable
	}
	return f.Declaration + "." + f.Callable
}

// sampled reports whether eidos derives two distinct samples of v.
func (v *Value) sampled() bool {
	return v.Samples[0].OK() && v.Samples[1].OK() && !reflect.DeepEqual(v.Samples[0].Value, v.Samples[1].Value)
}

// made reports whether eidos derives both samples of the input v as a
// function literal. Each literal is a new value wherever a call writes it,
// so its two samples differ also where they have one spelling. The two makes
// of a channel differ in the same way, but nothing sends to them or closes
// them, so a callable that receives from one never returns.
func (v *Value) made() bool {
	return literal(&v.Samples[0].Value) && literal(&v.Samples[1].Value)
}

// literal reports whether v is a function literal, also under the
// conversions to a defined function type. The value of a refused sample is
// the zero value, which is no literal.
func literal(v *emit.Value) bool {
	for v.Kind == emit.ValueConversion && v.Inner != nil {
		v = v.Inner
	}
	return v.Kind == emit.ValueFunc
}

// Suite is the derivation of one declaration.
type Suite struct {
	Checks    []Check
	Refusals  []Refusal
	Unclaimed []Unclaimed
	Edges     []Edge
}

// Check is one claim bound to one callable.
type Check struct {
	// Claim is the claim of the check, and Facts the facts of its callable.
	Claim *spec.Claim
	Facts *Facts

	// Pairs contains the partner of each name of a partner of the claim.
	Pairs map[string]*Pair

	// Identity is the identity of the check, such as Store.Read/zero-on-error,
	// and Key the key that joins it across languages, such as
	// store.read/zero-on-error.
	Identity, Key string

	// Statement and Clause are the statement of the claim and the clause of
	// its planted defect, with the names of the declaration, the callable and
	// each partner in place.
	Statement, Clause string

	// Shallow is true when the machineries of the check do not include the
	// depth of its law.
	Shallow bool

	// Generated is true for a check whose inputs the property engine
	// generates. Such a check is the generated binding of a claim of samples,
	// or the check of a claim of the machinery generated-inputs.
	Generated bool
}

// Machineries returns the machineries of the check, the deepest first. They
// are the machineries of its claim. A generated binding of a claim of samples
// adds generated-inputs before nothing-beyond-a-caller, or at the end of a
// list without it.
func (c *Check) Machineries() []spec.MachineryID {
	m := c.Claim.Machineries
	if !c.Generated || slices.Contains(m, spec.GeneratedInputs) {
		return m
	}
	at := slices.Index(m, spec.NothingBeyondACaller)
	if at < 0 {
		at = len(m)
	}
	return slices.Insert(slices.Clone(m), at, spec.GeneratedInputs)
}

// Refusal is a licensed claim whose gate failed on a callable.
type Refusal struct {
	Claim    spec.ClaimID
	Callable string

	// Key is the key of the check that the refusal replaces, such as
	// store.scan/zero-on-error.
	Key string

	// Gate is the gate that fails for the callable. Why and Remedy are the
	// reason and the remedy of the gate with the names of the callable.
	Gate spec.GateID

	Why, Remedy string
}

// Unclaimed is a classification of a callable without a bound or a refused
// claim, with its stance.
type Unclaimed struct {
	// Classification is shape:<name>, mixin:<name> or
	// contract:<name>/<role>.
	Classification, Callable string

	// Key is the key of the callable, such as store.put.
	Key string

	Stance spec.StanceID
	Note   string
}

// Edge is a pair of checks on one callable where the law of the deeper check
// implies the law of the shallower one.
type Edge struct {
	Deeper, Shallower string
}

// licensing is one way in which a callable has the licence of a claim: the
// classification that licenses it, and the id of the contract instance where
// the claim binds once for each instance.
type licensing struct {
	licence        *spec.Licence
	classification string
	instance       string
}

// Derive returns the checks of one declaration, with its refusals, its
// unclaimed classifications and its entailment edges. It visits the callables
// of decl in their order and the claims of the suite in table in their order.
// suite and runtime are the machineries that the suite generator and the
// runtime supply, and the machineries that the declaration states belong to
// suite.
//
// draws reports whether a check passes a sample to an input of one of its
// calls, which the language of the suite decides. Derive adds the generated
// binding of each such check after it. The generated binding has the case
// generated, as Store.Get/zero-on-error#generated, or the suffix -generated
// after the id of an instance, as Store.Get/cache#a-generated. A nil draws
// does not add a generated binding. A claim of the machinery generated-inputs
// binds once, as a check that the property engine generates. A claim of the
// machinery generated-sequences binds once too, because the actions of its
// sequences generate the inputs of its calls.
//
// # Errors
//
// Derive returns an error that wraps [ErrKey] for two callables of decl
// whose names have one key.
func Derive(
	decl []Facts, table []spec.Claim, suite, runtime []spec.MachineryID, draws func(c *Check) bool,
) (Suite, error) {
	d := derivation{supplied: slices.Concat(suite, runtime), draws: draws}
	keys := map[string]string{}
	for i := range decl {
		f := &decl[i]
		key := f.Key()
		if other, ok := keys[key]; ok {
			return Suite{}, fmt.Errorf("%w: %s and %s have the key %s", ErrKey, other, f.Callable, key)
		}
		keys[key] = f.Callable
		d.callable(f, decl, table)
	}
	return d.out, nil
}

// derivation is the derivation of one declaration. supplied are the
// machineries that the suite generator and the runtime supply, and draws
// reports whether a check draws an input in the language of the suite. out is
// the suite that the derivation adds to.
type derivation struct {
	draws    func(c *Check) bool
	supplied []spec.MachineryID
	out      Suite
}

// callable adds the checks, the refusals, the unclaimed classifications and
// the edges of the callable of f to the suite of d.
func (d *derivation) callable(f *Facts, decl []Facts, table []spec.Claim) {
	out, supplied := &d.out, d.supplied
	covered := map[string]bool{}
	needs := map[string]string{}
	first := len(out.Checks)
	for ci := range table {
		c := &table[ci]
		if c.Generator != "suite" {
			continue
		}
		licensings := licensed(c, f, decl)
		if c.Unless != nil && slices.Contains(f.Class.Mixins, c.Unless.Mixin) {
			cover(covered, licensings)
			continue
		}
		if missing := firstMissing(c.Machineries, supplied); missing != "" {
			for _, l := range licensings {
				if _, ok := needs[l.classification]; !ok {
					needs[l.classification] = needsNote(c.ID, missing)
				}
			}
			continue
		}
		for _, l := range licensings {
			check, refusal, silent := bind(c, l, f, decl)
			switch {
			case silent:
				continue
			case refusal != nil:
				out.Refusals = append(out.Refusals, *refusal)
			default:
				out.Checks = append(out.Checks, check)
				sequences := slices.Contains(c.Machineries, spec.GeneratedSequences)
				if !check.Generated && !sequences && d.draws != nil && d.draws(&check) {
					out.Checks = append(out.Checks, generatedOf(&check))
				}
			}
			covered[l.classification] = true
		}
	}
	out.Edges = append(out.Edges, edgesOf(out.Checks[first:])...)
	for _, cls := range classificationsOf(f) {
		if covered[cls] {
			continue
		}
		key := f.Key()
		if note, ok := needs[cls]; ok {
			out.Unclaimed = append(
				out.Unclaimed,
				Unclaimed{Classification: cls, Callable: f.Callable, Key: key, Stance: spec.Needs, Note: note},
			)
			continue
		}
		if u, ok := spec.UnclaimedOf(cls); ok {
			out.Unclaimed = append(
				out.Unclaimed,
				Unclaimed{Classification: cls, Callable: f.Callable, Key: key, Stance: u.Stance, Note: u.Note},
			)
		}
	}
}

// cover marks the classification of each licensing as covered.
func cover(covered map[string]bool, licensings []licensing) {
	for _, l := range licensings {
		covered[l.classification] = true
	}
}

// bind binds the claim c through the licensing l on the callable of f. It
// returns the check, or the refusal of the claim, or silent for a claim of
// every callable whose gate fails. The identity of a claim of a declaration
// names the declaration, as Store/agrees-with-the-reference.
func bind(c *spec.Claim, l licensing, f *Facts, decl []Facts) (Check, *Refusal, bool) {
	everyCallable := l.licence.Callable
	names := namesOf(f)
	identity := f.qualified() + "/" + string(c.ID)
	if l.licence.Declaration || c.Scope == spec.ScopeDeclaration {
		identity = names["declaration"] + "/" + string(c.ID)
	}
	if l.instance != "" {
		identity += "#" + l.instance
	}
	refuse := func(gate *spec.Gate, names map[string]string) (Check, *Refusal, bool) {
		if everyCallable {
			return Check{}, nil, true
		}
		r := refusal(c.ID, f, spec.Key(identity), gate, names)
		return Check{}, &r, false
	}
	if g, failed := failing(c.Gates, f); failed {
		return refuse(gateOf(g), names)
	}
	for i := range f.Inputs {
		in := &f.Inputs[i]
		if !in.sampled() && !in.made() && !slices.Contains(c.Machineries, spec.GeneratedSequences) {
			return refuse(gateOf(spec.InputSamples), mapWith(names, "param", in.Name))
		}
	}
	for _, b := range c.Inputs {
		if b.Kind == spec.FromParam && b.Default == "" && !hasClassParam(f, b.Param) {
			return refuse(gateOf(spec.ClassificationParams), mapWith(names, "param", b.Param))
		}
	}
	pairs := map[string]*Pair{}
	for pi := range c.Partners {
		p := &c.Partners[pi]
		pair, gate := Partner(p, f, decl)
		if gate != nil && (pair == nil || !unbound(c, p.Name)) {
			return refuse(gate, mapWith(names, "partner", p.Name))
		}
		if slices.Contains(c.Gates, spec.WriteReturnsTheValue) && !returnsTheValue(pair.Partner, f) {
			return refuse(gateOf(spec.WriteReturnsTheValue), mapWith(names, "partner", p.Name))
		}
		pairs[p.Name] = pair
		names[p.Name] = pair.Partner.Callable
	}
	law, _ := spec.LawOf(c.Law)
	return Check{
		Claim: c, Facts: f, Pairs: pairs,
		Identity: identity, Key: spec.Key(identity),
		Statement: fill(c.Statement, names), Clause: fill(c.Defect.Clause, names),
		Shallow:   !slices.Contains(c.Machineries, law.Depth),
		Generated: slices.Contains(c.Machineries, spec.GeneratedInputs),
	}, nil, false
}

// unbound reports whether the claim c of generated sequences takes the partner
// name only as inputs of the kinds operation and call. Its probe calls such a
// partner by its name or with its own samples, so the arguments of the
// partner need no binding.
func unbound(c *spec.Claim, name string) bool {
	if !slices.Contains(c.Machineries, spec.GeneratedSequences) {
		return false
	}
	probe, _ := spec.ProbeOf(c.Probe)
	return !slices.ContainsFunc(c.Inputs, func(b spec.Binding) bool {
		if b.Partner != name {
			return false
		}
		at := slices.IndexFunc(probe.Inputs, func(in spec.Input) bool { return in.Name == b.Input })
		kind := probe.Inputs[at].Kind
		return kind != spec.KindOperation && kind != spec.KindCall
	})
}

// GeneratedOf returns the identity of the generated binding of the check of
// samples with the identity id: id with the case generated, as
// Store.Get/zero-on-error#generated, or with the suffix -generated after the
// id of an instance, as Store.Get/cache#a-generated. For the key of a check
// it returns the key of the generated binding.
func GeneratedOf(id string) string {
	if strings.Contains(id, "#") {
		return id + "-generated"
	}
	return id + "#generated"
}

// generatedOf returns the generated binding of the check c of samples: the
// check with the case generated in its identity and its key, and the depth of
// its machineries.
func generatedOf(c *Check) Check {
	g := *c
	g.Identity = GeneratedOf(g.Identity)
	g.Key = spec.Key(g.Identity)
	g.Generated = true
	law, _ := spec.LawOf(g.Claim.Law)
	g.Shallow = !slices.Contains(g.Machineries(), law.Depth)
	return g
}

// Refuse returns the refusal of the check c by the gate id, with the names of
// the callable of c and of each of its partners in the reason and the remedy.
// The suite generator of a language refuses a derived check by a gate on the
// code that it writes, such as importable-samples.
func Refuse(c *Check, id spec.GateID) Refusal {
	names := namesOf(c.Facts)
	for name, pair := range c.Pairs {
		names[name] = pair.Partner.Callable
	}
	return refusal(c.Claim.ID, c.Facts, c.Key, gateOf(id), names)
}

// refusal returns the refusal of the claim on the callable of f by gate. The
// reason and the remedy of the refusal are those of gate with names in place.
func refusal(claim spec.ClaimID, f *Facts, key string, gate *spec.Gate, names map[string]string) Refusal {
	return Refusal{
		Claim: claim, Callable: f.Callable, Key: key, Gate: gate.ID,
		Why: fill(gate.Why, names), Remedy: fill(gate.Remedy, names),
	}
}

// namesOf returns the names that the text of a claim or of a gate refers to:
// the callable of f and its declaration. A function is its own declaration.
func namesOf(f *Facts) map[string]string {
	declaration := f.Declaration
	if declaration == "" {
		declaration = f.Callable
	}
	return map[string]string{"callable": f.Callable, "declaration": declaration}
}

// licensed returns each way in which the callable of f has a licence of c.
// A callable with two instances of the contract of a licence has it once for
// each instance. One callable of the declaration decl has the licence of a
// claim of the declaration: the first callable that passes the gates of c, or
// the first callable where none does. For a claim of the scope declaration
// whose licence is a classification, that callable is among the callables
// with the classification.
func licensed(c *spec.Claim, f *Facts, decl []Facts) []licensing {
	var out []licensing
	for li := range c.Licences {
		l := &c.Licences[li]
		if slices.ContainsFunc(l.Conflicts, func(cls spec.Classification) bool { return has(f, cls) }) {
			continue
		}
		if l.With != nil && !has(f, *l.With) {
			continue
		}
		switch {
		case l.Callable:
			out = append(out, licensing{licence: l})
		case l.Context:
			if hasParamRole(f) {
				out = append(out, licensing{licence: l})
			}
		case l.Declaration:
			if f == anchorOf(c, decl, func(*Facts) bool { return true }) {
				out = append(out, licensing{licence: l})
			}
		case c.Scope == spec.ScopeDeclaration:
			cls := l.Classification
			if f == anchorOf(c, decl, func(g *Facts) bool { return has(g, cls) }) {
				out = append(out, licensing{licence: l, classification: labelOf(cls)})
			}
		case l.Contract != "":
			var memberships []shape.Membership
			for _, m := range f.Class.Contracts {
				if m.Contract == l.Contract && m.Role == l.Role {
					memberships = append(memberships, m)
				}
			}
			for _, m := range memberships {
				instance := ""
				if len(memberships) > 1 {
					instance = m.ID
				}
				out = append(out, licensing{licence: l, classification: labelOf(l.Classification), instance: instance})
			}
		case has(f, l.Classification) && (l.NeedsParam == "" || hasClassParam(f, l.NeedsParam)):
			out = append(out, licensing{licence: l, classification: labelOf(l.Classification)})
		}
	}
	return out
}

// anchorOf returns the callable of decl on which a claim c of the declaration
// binds: the first callable that candidate accepts and that passes the gates
// of c, or the first callable that candidate accepts where none passes. It
// returns nil where candidate accepts no callable.
func anchorOf(c *spec.Claim, decl []Facts, candidate func(f *Facts) bool) *Facts {
	var first *Facts
	for i := range decl {
		f := &decl[i]
		if !candidate(f) {
			continue
		}
		if _, failed := failing(c.Gates, f); !failed {
			return f
		}
		if first == nil {
			first = f
		}
	}
	return first
}

// failing returns the first gate of gates with a condition that the facts f
// fail, and false where f pass the condition of each gate.
func failing(gates []spec.GateID, f *Facts) (spec.GateID, bool) {
	for _, g := range gates {
		if passes, defined := Passes(g, f); defined && !passes {
			return g, true
		}
	}
	return "", false
}

// has reports whether the callable of f has the classification cls.
func has(f *Facts, cls spec.Classification) bool {
	switch {
	case cls.Shape != "":
		return f.Class.Shape == cls.Shape
	case cls.Mixin != "":
		return slices.Contains(f.Class.Mixins, cls.Mixin)
	default:
		return slices.ContainsFunc(f.Class.Contracts, func(m shape.Membership) bool {
			return m.Contract == cls.Contract && (cls.Role == "" || m.Role == cls.Role)
		})
	}
}

// hasParamRole reports whether a parameter of the callable of f has the role
// context.
func hasParamRole(f *Facts) bool {
	passes, _ := Passes(spec.Cancellation, f)
	return passes
}

// hasClassParam reports whether a classification of the callable of f sets
// the param key.
func hasClassParam(f *Facts, key string) bool {
	for _, params := range f.Class.Params {
		if _, ok := params[directive.ParamKey(key)]; ok {
			return true
		}
	}
	return false
}

// labelOf returns the classification in the form of an entry of unclaimed.
func labelOf(cls spec.Classification) string {
	switch {
	case cls.Shape != "":
		return "shape:" + string(cls.Shape)
	case cls.Mixin != "":
		return "mixin:" + string(cls.Mixin)
	default:
		return "contract:" + string(cls.Contract) + "/" + string(cls.Role)
	}
}

// classificationsOf returns each classification of the callable of f in the
// form of an entry of unclaimed: its shape, its mixins and its roles. A
// callable without a shape has the label shape:, which no entry of unclaimed
// matches.
func classificationsOf(f *Facts) []string {
	out := []string{"shape:" + string(f.Class.Shape)}
	for _, m := range f.Class.Mixins {
		out = append(out, labelOf(spec.Classification{Mixin: m}))
	}
	for _, m := range f.Class.Contracts {
		label := labelOf(spec.Classification{Contract: m.Contract, Role: m.Role})
		if !slices.Contains(out, label) {
			out = append(out, label)
		}
	}
	return out
}

// firstMissing returns the first machinery of needed that supplied lacks,
// and the empty string where supplied has every one.
func firstMissing(needed, supplied []spec.MachineryID) spec.MachineryID {
	for _, m := range needed {
		if !slices.Contains(supplied, m) {
			return m
		}
	}
	return ""
}

// needsNote returns the note of a classification whose claim needs the
// machinery m. A machinery of the consumer names the param that states it.
func needsNote(claim spec.ClaimID, m spec.MachineryID) string {
	machinery, _ := spec.MachineryOf(m)
	if machinery.Param != "" {
		return fmt.Sprintf(
			"The claim %s needs the param machineries=[%s] of the directive suite.",
			claim,
			machinery.Param,
		)
	}
	return fmt.Sprintf("The claim %s needs %s.", claim, m)
}

// edgesOf returns the edges between the checks of one callable. An edge joins
// two sample bindings or two generated bindings.
func edgesOf(checks []Check) []Edge {
	var out []Edge
	for _, c := range checks {
		for _, e := range c.Claim.Entails {
			i := slices.IndexFunc(checks, func(o Check) bool { return o.Claim.ID == e && o.Generated == c.Generated })
			if i >= 0 {
				out = append(out, Edge{Deeper: c.Identity, Shallower: checks[i].Identity})
			}
		}
	}
	return out
}

// fill replaces each name in braces in text with its value in names.
func fill(text string, names map[string]string) string {
	for name, value := range names {
		text = strings.ReplaceAll(text, "{"+name+"}", value)
	}
	return text
}

// mapWith returns a copy of names with key set to value.
func mapWith(names map[string]string, key, value string) map[string]string {
	out := make(map[string]string, len(names)+1)
	maps.Copy(out, names)
	out[key] = value
	return out
}
