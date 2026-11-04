// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"regexp"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/rules"
)

// classified returns the spec of the shape catalog for a classification,
// and false for a classification that the catalog lacks or that does not
// have the form of its field.
func classified(c rules.Classification) (shape.Spec, bool) {
	var name string
	var form shape.Form
	switch {
	case c.Shape != "":
		name, form = string(c.Shape), shape.FormShape
	case c.Mixin != "":
		name, form = string(c.Mixin), shape.FormMixin
	default:
		name, form = string(c.Contract), shape.FormContract
	}
	s, ok := shape.SpecOf(name)
	if !ok || s.Form != form {
		return shape.Spec{}, false
	}
	hasRole := slices.ContainsFunc(s.Roles, func(r shape.RoleSpec) bool { return r.Name == c.Role })
	if c.Contract != "" && c.Role != "" && !hasRole {
		return shape.Spec{}, false
	}
	return s, true
}

// hasParam reports whether the classification of s declares the param key.
func hasParam(s shape.Spec, key string) bool {
	return slices.ContainsFunc(s.Params, func(p shape.Param) bool { return string(p.Key) == key })
}

// classificationsOf returns the specs of the shape catalog for the
// classifications of the licences of c and of their classifications With.
func classificationsOf(c rules.Claim) []shape.Spec {
	var specs []shape.Spec
	for _, l := range c.Licences {
		if l.Callable || l.Context {
			continue
		}
		if s, ok := classified(l.Classification); ok {
			specs = append(specs, s)
		}
		if l.With != nil {
			if s, ok := classified(*l.With); ok {
				specs = append(specs, s)
			}
		}
	}
	return specs
}

// directiveParams are the params of the directives of the generators enum
// and sentinel that a claim of either reads.
var directiveParams = map[string][]string{
	"enum":     nil,
	"sentinel": {"prefix", "no-overlap-with"},
}

// sources are the kinds of source that a claim of each generator binds.
var sources = map[string][]rules.SourceKind{
	"suite": {
		rules.FromCallable, rules.FromPartner, rules.FromInput, rules.FromSample, rules.FromZero, rules.FromParam,
		rules.FromClock, rules.FromDeclaration,
	},
	"enum": {
		rules.FromMembers,
		rules.FromName,
		rules.FromParse,
		rules.FromList,
		rules.FromValid,
		rules.FromOutside,
		rules.FromParam,
	},
	"sentinel": {
		rules.FromErrorValues, rules.FromErrorTypes, rules.FromCauses, rules.FromParam,
	},
}

// placeholder matches a name in braces in a statement or a clause.
var placeholder = regexp.MustCompile(`\{([a-z-]+)\}`)

func TestClaim(t *testing.T) {
	t.Parallel()

	t.Run("Claims", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each claim of the spec with its rule, in the order of the spec", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(rules.Claims(), func(c rules.Claim) spec.ClaimID { return c.ID }),
				keysOf(spec.Claims(), func(c spec.Claim) spec.ClaimID { return c.ID }),
				"each claim of the spec has a rule",
			)
			for _, c := range rules.Claims() {
				want, _ := spec.ClaimOf(c.ID)
				expect.Equal(t, c.Claim, want, "claim "+string(c.ID)+" has the runtime contract of the spec")
			}
		})

		t.Run("returns claims that each have a rule of claims.json", func(t *testing.T) {
			t.Parallel()
			ids := keysOf(spec.Claims(), func(c spec.Claim) spec.ClaimID { return c.ID })
			slices.Sort(ids)

			assert.Equal(t, ruleIDs(t, "claims", "claims"), ids,
				"claims.json has a rule for each claim of the spec and for no other claim")
		})

		t.Run("returns claims whose gates the spec defines", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				for _, g := range c.Gates {
					_, ok := spec.GateOf(g)
					expect.True(t, ok, "the spec has the gate "+string(g)+" of claim "+string(c.ID))
				}
			}
		})

		t.Run("returns claims that bind each input of their probe once", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				p, _ := spec.ProbeOf(c.Probe)
				want := keysOf(p.Inputs, func(in spec.Input) string { return in.Name })
				got := keysOf(c.Inputs, func(b rules.Binding) string { return b.Input })
				slices.Sort(want)
				slices.Sort(got)
				expect.Equal(t, got, want, "claim "+string(c.ID)+" binds the inputs of its probe")
			}
		})

		t.Run("returns claims whose sources their generator supplies", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				id := string(c.ID)
				partners := keysOf(c.Partners, func(p rules.Partner) string { return p.Name })
				for _, b := range c.Inputs {
					expect.True(
						t,
						slices.Contains(sources[c.Generator], b.Kind),
						"the generator of claim "+id+" supplies the source of input "+b.Input,
					)
					if b.Kind == rules.FromPartner {
						expect.True(
							t,
							slices.Contains(partners, b.Partner),
							"input "+b.Input+" of claim "+id+" names a partner of the claim",
						)
					}
					if b.Kind == rules.FromParam {
						expect.True(
							t,
							paramOf(c, b.Param),
							"input "+b.Input+" of claim "+id+" reads a param that the claim has",
						)
					}
				}
				for _, name := range partners {
					expect.True(
						t,
						slices.ContainsFunc(c.Inputs, func(b rules.Binding) bool { return b.Partner == name }),
						"an input of claim "+id+" calls its partner "+name,
					)
				}
			}
		})

		t.Run("returns claims whose licences the shape catalog has", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				id := string(c.ID)
				expect.Equal(
					t,
					len(c.Licences) > 0,
					c.Generator == "suite",
					"claim "+id+" has licences when its generator is suite",
				)
				for _, l := range c.Licences {
					if l.Callable || l.Context || l.Declaration {
						expect.Equal(
							t,
							l.Classification,
							rules.Classification{},
							"a licence of claim "+id+" sets one kind",
						)
						continue
					}
					s, ok := classified(l.Classification)
					expect.True(t, ok, "the shape catalog has a licence of claim "+id)
					if l.NeedsParam != "" {
						expect.True(
							t,
							ok && hasParam(s, l.NeedsParam),
							"a licence of claim "+id+" needs a param of its classification",
						)
					}
					if l.With != nil {
						_, ok := classified(*l.With)
						expect.True(t, ok, "the shape catalog has the classification With of a licence of claim "+id)
					}
					for _, conflict := range l.Conflicts {
						_, ok := classified(conflict)
						expect.True(t, ok, "the shape catalog has a conflict of a licence of claim "+id)
					}
				}
				if c.Unless != nil && c.Unless.Mixin != "" {
					_, ok := classified(rules.Classification{Mixin: c.Unless.Mixin})
					expect.True(t, ok, "the shape catalog has the mixin unless of claim "+id)
				}
				if c.Unless != nil && c.Unless.Param != "" {
					expect.True(
						t,
						slices.Contains(directiveParams[c.Generator], c.Unless.Param),
						"the directive of claim "+id+" has the param unless",
					)
				}
			}
		})

		t.Run("returns claims whose partners their mode can find", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				id := string(c.ID)
				specs := classificationsOf(c)
				for _, p := range c.Partners {
					switch p.Mode {
					case rules.ModeNamed, rules.ModeReader, rules.ModeDerived, rules.ModeInterposed:
						expect.True(
							t,
							slices.ContainsFunc(specs, func(s shape.Spec) bool { return hasParam(s, p.Param) }),
							"partner "+p.Name+" of claim "+id+" reads a param of a licence",
						)
					case rules.ModeRole:
						expect.True(t, slices.ContainsFunc(specs, func(s shape.Spec) bool {
							return slices.ContainsFunc(s.Roles, func(r shape.RoleSpec) bool { return r.Name == p.Role })
						}), "partner "+p.Name+" of claim "+id+" is a role of a contract of a licence")
					case rules.ModeRewriter:
						expect.True(t, slices.ContainsFunc(specs, func(s shape.Spec) bool {
							return hasParam(s, p.Param) ||
								slices.ContainsFunc(s.Roles, func(r shape.RoleSpec) bool { return r.Name == p.Role })
						}), "partner "+p.Name+" of claim "+id+" reads a param or a role of a licence")
					case rules.ModeEstablishing:
						_, ok := classified(p.Classification)
						expect.True(
							t,
							ok,
							"the shape catalog has the classification of partner "+p.Name+" of claim "+id,
						)
					case rules.ModePaired:
						expect.Equal(t, p, rules.Partner{Name: p.Name, Mode: rules.ModePaired},
							"partner "+p.Name+" of claim "+id+" reads no param and no classification")
					default:
						expect.True(t, false, "partner "+p.Name+" of claim "+id+" has a mode of the rules")
					}
				}
			}
		})

		t.Run("returns claims whose sentences name only the subject and the partners", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				names := append(
					keysOf(c.Partners, func(p rules.Partner) string { return p.Name }),
					"callable",
					"declaration",
				)
				for _, m := range placeholder.FindAllStringSubmatch(c.Statement+" "+c.Clause, -1) {
					expect.True(
						t,
						slices.Contains(names, m[1]),
						"claim "+string(c.ID)+" names {"+m[1]+"}, its subject or a partner",
					)
				}
			}
		})

		t.Run("returns entailments between claims with a licence in common and without a cycle", func(t *testing.T) {
			t.Parallel()

			for _, c := range rules.Claims() {
				for _, e := range c.Entails {
					shallower, ok := rules.ClaimOf(e)
					expect.True(t, ok, "the spec has the claim "+string(e)+" that claim "+string(c.ID)+" entails")
					expect.True(
						t,
						shareLicence(c, shallower),
						"claim "+string(c.ID)+" and claim "+string(e)+" have a licence in common",
					)
				}
				expect.False(
					t,
					leadsTo(c.ID, c.Entails, map[spec.ClaimID]bool{}),
					"the entailments of claim "+string(c.ID)+" do not lead back to it",
				)
			}
		})
	})

	t.Run("ClaimOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the claim with its runtime contract, its partners and its inputs", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.ClaimOf(spec.ReadAfterWrite)

			assert.True(t, ok, "the rules have the claim")
			assert.Equal(t, got.Statement, "{callable} returns what {write} stored before it",
				"the claim has the statement of the spec")
			assert.Equal(t, got.Partners, []rules.Partner{{Name: "write", Mode: rules.ModeNamed, Param: "write"}},
				"the claim has its partner")
			assert.Equal(t, got.Inputs, []rules.Binding{
				{Input: "read", Kind: rules.FromCallable},
				{Input: "write", Kind: rules.FromPartner, Partner: "write"},
			}, "the claim has its inputs in order")
		})

		t.Run("returns a claim of enum or sentinel that a param of the directive withdraws", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.ClaimOf(spec.ErrorMessages)

			assert.True(t, ok, "the rules have the claim")
			assert.Equal(
				t,
				*got.Unless,
				rules.Unless{Param: "prefix", Value: "off"},
				"the param prefix withdraws the claim",
			)
		})

		t.Run("returns false for an identifier without a claim", func(t *testing.T) {
			t.Parallel()

			_, ok := rules.ClaimOf("no-such-claim")

			assert.False(t, ok, "the rules have no such claim")
		})
	})
}

// paramOf reports whether the claim c has the param key: a param of the
// classification of a licence for a claim of suite, and a param of the
// directive for a claim of enum or sentinel.
func paramOf(c rules.Claim, key string) bool {
	if c.Generator != "suite" {
		return slices.Contains(directiveParams[c.Generator], key)
	}
	return slices.ContainsFunc(classificationsOf(c), func(s shape.Spec) bool { return hasParam(s, key) })
}

// shareLicence reports whether two claims have a licence in common. A
// licence of every callable is in common with each licence.
func shareLicence(a, b rules.Claim) bool {
	for _, x := range a.Licences {
		for _, y := range b.Licences {
			if x.Callable || y.Callable || x.Classification == y.Classification {
				return true
			}
		}
	}
	return false
}

// leadsTo reports whether the entailments from claims lead to the claim id.
func leadsTo(id spec.ClaimID, claims []spec.ClaimID, seen map[spec.ClaimID]bool) bool {
	for _, e := range claims {
		if e == id {
			return true
		}
		if seen[e] {
			continue
		}
		seen[e] = true
		next, _ := rules.ClaimOf(e)
		if leadsTo(id, next.Entails, seen) {
			return true
		}
	}
	return false
}
