// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"regexp"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/plugin/shape"

	"go.dokimi.dev/dokimi/spec"
)

// classified returns the spec of the shape catalog for a classification,
// and false for a classification that the catalog lacks or that does not
// have the form of its field.
func classified(c spec.Classification) (shape.Spec, bool) {
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
func classificationsOf(c spec.Claim) []shape.Spec {
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
var sources = map[string][]spec.SourceKind{
	"suite": {spec.FromCallable, spec.FromPartner, spec.FromSample, spec.FromZero, spec.FromParam, spec.FromClock},
	"enum": {
		spec.FromMembers,
		spec.FromName,
		spec.FromParse,
		spec.FromList,
		spec.FromValid,
		spec.FromOutside,
		spec.FromParam,
	},
	"sentinel": {
		spec.FromErrorValues, spec.FromErrorTypes, spec.FromCauses, spec.FromParam,
	},
}

// placeholder matches a name in braces in a statement or a clause.
var placeholder = regexp.MustCompile(`\{([a-z-]+)\}`)

func TestClaim(t *testing.T) {
	t.Parallel()

	t.Run("Claims", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the claims of claims.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(spec.Claims(), func(c spec.Claim) spec.ClaimID { return c.ID }),
				idsOf(t, "claims", "claims"),
				"the claims keep the order of claims.json",
			)
		})

		t.Run("returns claims that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "claims")
		})

		t.Run("returns claims whose law, probe, gates, machineries and defect the spec defines", func(t *testing.T) {
			t.Parallel()

			for _, c := range spec.Claims() {
				id := string(c.ID)
				_, ok := spec.LawOf(c.Law)
				expect.True(t, ok, "the spec has the law of claim "+id)
				p, ok := spec.ProbeOf(c.Probe)
				expect.True(t, ok, "the spec has the probe of claim "+id)
				for _, g := range c.Gates {
					_, found := spec.GateOf(g)
					expect.True(t, found, "the spec has the gate "+string(g)+" of claim "+id)
				}
				expect.True(t, len(c.Machineries) > 0, "claim "+id+" needs a machinery")
				for _, m := range c.Machineries {
					_, found := spec.MachineryOf(m)
					expect.True(t, found, "the spec has the machinery "+string(m)+" of claim "+id)
				}
				_, ok = spec.DefectFormOf(c.Defect.Form)
				expect.True(t, ok, "the spec has the defect form of claim "+id)
				expect.True(
					t,
					slices.Contains(p.Reasons, c.Defect.Reason),
					"the probe of claim "+id+" reports the reason code of its defect",
				)
			}
		})

		t.Run("returns claims that bind each input of their probe once", func(t *testing.T) {
			t.Parallel()

			for _, c := range spec.Claims() {
				id := string(c.ID)
				p, _ := spec.ProbeOf(c.Probe)
				want := keysOf(p.Inputs, func(in spec.Input) string { return in.Name })
				got := keysOf(c.Inputs, func(b spec.Binding) string { return b.Input })
				slices.Sort(want)
				slices.Sort(got)
				expect.Equal(t, got, want, "claim "+id+" binds the inputs of its probe")
			}
		})

		t.Run("returns claims whose sources their generator supplies", func(t *testing.T) {
			t.Parallel()

			for _, c := range spec.Claims() {
				id := string(c.ID)
				partners := keysOf(c.Partners, func(p spec.Partner) string { return p.Name })
				for _, b := range c.Inputs {
					expect.True(
						t,
						slices.Contains(sources[c.Generator], b.Kind),
						"the generator of claim "+id+" supplies the source of input "+b.Input,
					)
					if b.Kind == spec.FromPartner {
						expect.True(
							t,
							slices.Contains(partners, b.Partner),
							"input "+b.Input+" of claim "+id+" names a partner of the claim",
						)
					}
					if b.Kind == spec.FromParam {
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
						slices.ContainsFunc(c.Inputs, func(b spec.Binding) bool { return b.Partner == name }),
						"an input of claim "+id+" calls its partner "+name,
					)
				}
			}
		})

		t.Run("returns claims whose licences the shape catalog has", func(t *testing.T) {
			t.Parallel()

			for _, c := range spec.Claims() {
				id := string(c.ID)
				expect.Equal(
					t,
					len(c.Licences) > 0,
					c.Generator == "suite",
					"claim "+id+" has licences when its generator is suite",
				)
				for _, l := range c.Licences {
					if l.Callable || l.Context {
						expect.Equal(
							t,
							l.Classification,
							spec.Classification{},
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
					_, ok := classified(spec.Classification{Mixin: c.Unless.Mixin})
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

			for _, c := range spec.Claims() {
				id := string(c.ID)
				specs := classificationsOf(c)
				for _, p := range c.Partners {
					switch p.Mode {
					case spec.ModeNamed, spec.ModeReader, spec.ModeDerived, spec.ModeInterposed:
						expect.True(
							t,
							slices.ContainsFunc(specs, func(s shape.Spec) bool { return hasParam(s, p.Param) }),
							"partner "+p.Name+" of claim "+id+" reads a param of a licence",
						)
					case spec.ModeRole:
						expect.True(t, slices.ContainsFunc(specs, func(s shape.Spec) bool {
							return slices.ContainsFunc(s.Roles, func(r shape.RoleSpec) bool { return r.Name == p.Role })
						}), "partner "+p.Name+" of claim "+id+" is a role of a contract of a licence")
					case spec.ModeRewriter:
						expect.True(t, slices.ContainsFunc(specs, func(s shape.Spec) bool {
							return hasParam(s, p.Param) ||
								slices.ContainsFunc(s.Roles, func(r shape.RoleSpec) bool { return r.Name == p.Role })
						}), "partner "+p.Name+" of claim "+id+" reads a param or a role of a licence")
					case spec.ModeEstablishing:
						_, ok := classified(p.Classification)
						expect.True(
							t,
							ok,
							"the shape catalog has the classification of partner "+p.Name+" of claim "+id,
						)
					default:
						expect.True(t, false, "partner "+p.Name+" of claim "+id+" has a mode of the spec")
					}
				}
			}
		})

		t.Run("returns claims whose sentences name only the subject and the partners", func(t *testing.T) {
			t.Parallel()

			for _, c := range spec.Claims() {
				names := append(
					keysOf(c.Partners, func(p spec.Partner) string { return p.Name }),
					"callable",
					"declaration",
				)
				for _, m := range placeholder.FindAllStringSubmatch(c.Statement+" "+c.Defect.Clause, -1) {
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

			for _, c := range spec.Claims() {
				for _, e := range c.Entails {
					shallower, ok := spec.ClaimOf(e)
					expect.True(t, ok, "the spec has the claim "+string(e)+" that claim "+string(c.ID)+" entails")
					expect.True(
						t,
						shareLicence(c, shallower),
						"claim "+string(c.ID)+" and claim "+string(e)+" have a licence in common",
					)
				}
				expect.False(
					t,
					reaches(c.ID, c.Entails, map[spec.ClaimID]bool{}),
					"the entailments of claim "+string(c.ID)+" do not lead back to it",
				)
			}
		})
	})

	t.Run("ClaimOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the claim with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.ClaimOf(spec.ReadAfterWrite)

			assert.True(t, ok, "the spec has the claim")
			assert.Equal(t, got.Partners, []spec.Partner{{Name: "write", Mode: spec.ModeNamed, Param: "write"}},
				"the claim has its partner")
			assert.Equal(t, got.Inputs, []spec.Binding{
				{Input: "read", Kind: spec.FromCallable},
				{Input: "write", Kind: spec.FromPartner, Partner: "write"},
			}, "the claim has its inputs in order")
		})

		t.Run("returns a claim of enum or sentinel that a param of the directive withdraws", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.ClaimOf(spec.ErrorMessages)

			assert.True(t, ok, "the spec has the claim")
			assert.Equal(
				t,
				*got.Unless,
				spec.Unless{Param: "prefix", Value: "off"},
				"the param prefix withdraws the claim",
			)
		})

		t.Run("returns false for an identifier without a claim", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.ClaimOf("no-such-claim")

			assert.False(t, ok, "the spec has no such claim")
		})
	})
}

// paramOf reports whether the claim c has the param key: a param of the
// classification of a licence for a claim of suite, and a param of the
// directive for a claim of enum or sentinel.
func paramOf(c spec.Claim, key string) bool {
	if c.Generator != "suite" {
		return slices.Contains(directiveParams[c.Generator], key)
	}
	return slices.ContainsFunc(classificationsOf(c), func(s shape.Spec) bool { return hasParam(s, key) })
}

// shareLicence reports whether two claims have a licence in common. A
// licence of every callable is in common with each licence.
func shareLicence(a, b spec.Claim) bool {
	for _, x := range a.Licences {
		for _, y := range b.Licences {
			if x.Callable || y.Callable || x.Classification == y.Classification {
				return true
			}
		}
	}
	return false
}

// reaches reports whether the entailments from claims lead to the claim id.
func reaches(id spec.ClaimID, claims []spec.ClaimID, seen map[spec.ClaimID]bool) bool {
	for _, e := range claims {
		if e == id {
			return true
		}
		if seen[e] {
			continue
		}
		seen[e] = true
		next, _ := spec.ClaimOf(e)
		if reaches(id, next.Entails, seen) {
			return true
		}
	}
	return false
}
