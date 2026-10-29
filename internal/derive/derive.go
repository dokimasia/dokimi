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

	"go.dokimi.dev/eidos/lang/naming"
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

// qualified returns the name of the callable of f in an identity: Store.Read
// for a method, and Slug for a function.
func (f *Facts) qualified() string {
	if f.Declaration == "" {
		return f.Callable
	}
	return f.Declaration + "." + f.Callable
}

// key returns the name of the callable of f in a key: store.read for a
// method, and slug for a function.
func (f *Facts) key() string {
	if f.Declaration == "" {
		return keyOf(f.Callable)
	}
	return keyOf(f.Declaration) + "." + keyOf(f.Callable)
}

// sampled reports whether eidos derives two distinct samples of v.
func (v *Value) sampled() bool {
	return v.Samples[0].OK() && v.Samples[1].OK() && !reflect.DeepEqual(v.Samples[0].Value, v.Samples[1].Value)
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

	// Shallow is true when the machineries of the claim do not include the
	// depth of its law.
	Shallow bool
}

// Refusal is a licensed claim whose gate failed on a callable.
type Refusal struct {
	Claim       spec.ClaimID
	Callable    string
	Why, Remedy string
}

// Unclaimed is a classification of a callable without a bound or a refused
// claim, with its stance.
type Unclaimed struct {
	// Classification is shape:<name>, mixin:<name> or
	// contract:<name>/<role>.
	Classification, Callable string

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
// # Errors
//
// Derive returns an error that wraps [ErrKey] for two callables of decl
// whose names have one key.
func Derive(decl []Facts, table []spec.Claim, suite, runtime []spec.MachineryID) (Suite, error) {
	supplied := slices.Concat(suite, runtime)
	var out Suite
	keys := map[string]string{}
	for i := range decl {
		f := &decl[i]
		key := f.key()
		if other, ok := keys[key]; ok {
			return Suite{}, fmt.Errorf("%w: %s and %s have the key %s", ErrKey, other, f.Callable, key)
		}
		keys[key] = f.Callable
		deriveCallable(&out, f, decl, table, supplied)
	}
	return out, nil
}

// deriveCallable adds the checks, the refusals, the unclaimed classifications
// and the edges of the callable of f to out.
func deriveCallable(out *Suite, f *Facts, decl []Facts, table []spec.Claim, supplied []spec.MachineryID) {
	covered := map[string]bool{}
	needs := map[string]string{}
	first := len(out.Checks)
	for ci := range table {
		c := &table[ci]
		if c.Generator != "suite" {
			continue
		}
		licensings := licensed(c, f)
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
			}
			covered[l.classification] = true
		}
	}
	out.Edges = append(out.Edges, edgesOf(out.Checks[first:])...)
	for _, cls := range classificationsOf(f) {
		if covered[cls] {
			continue
		}
		if note, ok := needs[cls]; ok {
			out.Unclaimed = append(
				out.Unclaimed,
				Unclaimed{Classification: cls, Callable: f.Callable, Stance: spec.Needs, Note: note},
			)
			continue
		}
		if u, ok := spec.UnclaimedOf(cls); ok {
			out.Unclaimed = append(
				out.Unclaimed,
				Unclaimed{Classification: cls, Callable: f.Callable, Stance: u.Stance, Note: u.Note},
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
// every callable whose gate fails.
func bind(c *spec.Claim, l licensing, f *Facts, decl []Facts) (Check, *Refusal, bool) {
	everyCallable := l.licence.Callable
	refuse := func(why, remedy string) (Check, *Refusal, bool) {
		if everyCallable {
			return Check{}, nil, true
		}
		return Check{}, &Refusal{Claim: c.ID, Callable: f.Callable, Why: why, Remedy: remedy}, false
	}
	declaration := f.Declaration
	if declaration == "" {
		declaration = f.Callable
	}
	names := map[string]string{"callable": f.Callable, "declaration": declaration}
	for _, g := range c.Gates {
		if passes, defined := Passes(g, f); defined && !passes {
			gate, _ := spec.GateOf(g)
			return refuse(fill(gate.Why, names), fill(gate.Remedy, names))
		}
	}
	for i := range f.Inputs {
		if in := &f.Inputs[i]; !in.sampled() {
			return refuse(
				fmt.Sprintf("eidos derives no sample of the parameter %s of %s", in.Name, f.Callable),
				fmt.Sprintf("author a sample of the parameter %s with the directive sample", in.Name),
			)
		}
	}
	for _, b := range c.Inputs {
		if b.Kind == spec.FromParam && !hasClassParam(f, b.Param) {
			return refuse(
				fmt.Sprintf("the classification of %s sets no param %s", f.Callable, b.Param),
				fmt.Sprintf("set the param %s in the directive of %s", b.Param, f.Callable),
			)
		}
	}
	pairs := map[string]*Pair{}
	for pi := range c.Partners {
		p := &c.Partners[pi]
		pair, gate := Partner(p, f, decl)
		if gate != nil {
			named := mapWith(names, "partner", p.Name)
			return refuse(fill(gate.Why, named), fill(gate.Remedy, named))
		}
		pairs[p.Name] = pair
		names[p.Name] = pair.Partner.Callable
	}
	law, _ := spec.LawOf(c.Law)
	identity := f.qualified() + "/" + string(c.ID)
	key := f.key() + "/" + string(c.ID)
	if l.instance != "" {
		identity += "#" + l.instance
		key += "#" + l.instance
	}
	return Check{
		Claim: c, Facts: f, Pairs: pairs,
		Identity: identity, Key: key,
		Statement: fill(c.Statement, names), Clause: fill(c.Defect.Clause, names),
		Shallow: !slices.Contains(c.Machineries, law.Depth),
	}, nil, false
}

// licensed returns each way in which the callable of f has a licence of c.
// A callable with two instances of the contract of a licence has it once for
// each instance.
func licensed(c *spec.Claim, f *Facts) []licensing {
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
			"The claim %s needs the param machineries=%s of the directive suite.",
			claim,
			machinery.Param,
		)
	}
	return fmt.Sprintf("The claim %s needs %s.", claim, m)
}

// edgesOf returns the edges between the checks of one callable.
func edgesOf(checks []Check) []Edge {
	var out []Edge
	for _, c := range checks {
		for _, e := range c.Claim.Entails {
			i := slices.IndexFunc(checks, func(o Check) bool { return o.Claim.ID == e })
			if i >= 0 {
				out = append(out, Edge{Deeper: c.Identity, Shallower: checks[i].Identity})
			}
		}
	}
	return out
}

// keyOf returns a name in the key form: its words in lower case, joined by
// hyphens.
func keyOf(name string) string {
	return strings.ToLower(strings.Join(naming.Words(name), "-"))
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
