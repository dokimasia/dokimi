// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"fmt"
	"maps"
	"slices"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
)

// RefuseDirective withdraws a derived claim of a callable for every subject
// of its suite. Its positional param ClaimParam contains the claim, and its
// param [ReasonParam] the reason why the declaration does not make the claim.
const (
	RefuseDirective directive.Name     = "refuse"
	ClaimParam      directive.ParamKey = "claim"
)

// KeyRefused and KeyRefusedWhy are the keys of the facts that the annotate
// half stamps on a callable under [RefuseDirective]: the claim of each
// directive, and the reason of each, in the order of the directives.
const (
	KeyRefused    meta.KeyName = "suite.refused"
	KeyRefusedWhy meta.KeyName = "suite.refused-why"
)

// refuseSchema is the directive of a refusal. A callable refuses each claim
// with a directive of its own, so the directive repeats.
var refuseSchema = directive.Schema{
	Name:       RefuseDirective,
	Repeatable: true,
	Kinds:      []symbol.Kind{symbol.KindMethod, symbol.KindFunction},
	Positional: []directive.ParamSpec{{
		Key: ClaimParam, Type: directive.TypeString, Required: true,
		Doc: "the derived claim that the declaration does not make for the callable, such as zero-on-error",
	}},
	Params: []directive.ParamSpec{{
		Key: ReasonParam, Type: directive.TypeString, Required: true,
		Doc: "the reason why the declaration does not make the claim",
	}},
	Doc: "refuses a derived claim of the callable for every subject of its suite, with the reason",
}

// refusalsOfMethod stamps the refusals of a method of a declaration under the
// directive suite. It reports [UnknownRefusal] at the refusal of a method of
// another declaration, because no suite derives its claims.
func (g *generator) refusalsOfMethod(m *sdk.MethodMatch, _ part, st *sdk.Stamper) error {
	method := m.Method
	name := method.ID.Owner + "." + method.Name
	if !Suited(m, &method.Host) {
		m.ErrorfAt(UnknownRefusal, m.Directive().Pos, "%s has no directive suite, so the refusal of %s by %s refuses "+
			"no claim", method.ID.Owner, m.Directive().Args[0].Str, name)
		return nil
	}
	g.stampRefusals(m, name, st)
	return nil
}

// refusalsOfFunction stamps the refusals of a function under the directive
// suite. It reports [UnknownRefusal] at the refusal of another function.
func (g *generator) refusalsOfFunction(m *sdk.FunctionMatch, _ part, st *sdk.Stamper) error {
	fn := m.Function
	if !Suited(m, &fn.ID) {
		m.ErrorfAt(UnknownRefusal, m.Directive().Pos, "%s has no directive suite, so the refusal of %s refuses no "+
			"claim", fn.Name, m.Directive().Args[0].Str)
		return nil
	}
	g.stampRefusals(m, fn.Name, st)
	return nil
}

// stampRefusals stamps the claim and the reason of each directive refuse of
// the callable name on the subject of m, once, at its first directive. A
// callable refuses each claim once, so stampRefusals reports
// [UnknownRefusal] at each later directive of a claim, and stamps the first.
func (g *generator) stampRefusals(m sdk.SubjectMatch, name string, st *sdk.Stamper) {
	var claims, reasons []string
	for d := range m.Directives(RefuseDirective) {
		if len(claims) == 0 && d.Instance != m.Directive().Instance {
			return
		}
		claim := d.Args[0].Str
		if slices.Contains(claims, claim) {
			m.ErrorfAt(UnknownRefusal, d.Pos, "%s refuses the claim %s twice, and a callable refuses each claim once",
				name, claim)
			continue
		}
		claims, reasons = append(claims, claim), append(reasons, d.Params[ReasonParam].Str)
	}
	sdk.Stamp(st, g.refused, claims)
	sdk.Stamp(st, g.refusedWhy, reasons)
}

// refusedOf returns the claims that the directives refuse of the callable id
// name, with the reason of each, and nil for a callable without one.
func (g *generator) refusedOf(m sdk.Matcher, id *symbol.Identity) map[spec.ClaimID]string {
	claims, _ := sdk.FactOf(m, *id, g.refused)
	reasons, _ := sdk.FactOf(m, *id, g.refusedWhy)
	if len(claims) == 0 {
		return nil
	}
	out := make(map[spec.ClaimID]string, len(claims))
	for i, claim := range claims {
		out[spec.ClaimID(claim)] = reasons[i]
	}
	return out
}

// unknownRefusals reports [UnknownRefusal] at the declaration of the suite s
// for each claim that a directive refuse of a callable of s names, where the
// derivation of s does not refuse the claim of that callable. The derivation
// refuses such a claim by the gate unrefused before any other gate. It returns
// false where it reports a claim. The message lists the claims that the
// derivation binds or refuses for the callable.
func (k *caller) unknownRefusals(s *Suite) bool {
	known := true
	for i := range s.Facts {
		f := &s.Facts[i]
		for _, claim := range slices.Sorted(maps.Keys(f.Refused)) {
			refused := slices.ContainsFunc(s.Derivation.Refusals, func(r derive.Refusal) bool {
				return r.Claim == claim && r.Callable == f.Callable
			})
			if refused {
				continue
			}
			k.m.Errorf(UnknownRefusal, "%s refuses the claim %s, which the derivation does not license for %s, so "+
				"%s has no suite%s", identityIn(s.Declaration, s.Func == nil, f.Callable, ""), claim, f.Callable,
				s.Declaration, claimsOf(f.Callable, &s.Derivation))
			known = false
		}
	}
	return known
}

// claimsOf returns the sentence that lists the claims that the derivation d
// binds or refuses for the callable, after a sentence that names a claim
// that it does not license, such as ". The claims of Record are smoke and
// cancel". It returns the empty string for a callable without claims.
func claimsOf(callable string, d *derive.Suite) string {
	var claims []string
	add := func(id spec.ClaimID) {
		if !slices.Contains(claims, string(id)) {
			claims = append(claims, string(id))
		}
	}
	for i := range d.Checks {
		if c := &d.Checks[i]; c.Facts.Callable == callable {
			add(c.Claim.ID)
		}
	}
	for _, r := range d.Refusals {
		if r.Callable == callable {
			add(r.Claim)
		}
	}
	switch {
	case len(claims) > 1:
		return fmt.Sprintf(". The claims of %s are %s", callable, listed(claims, "and"))
	case len(claims) == 1:
		return fmt.Sprintf(". The claim of %s is %s", callable, claims[0])
	}
	return ""
}
