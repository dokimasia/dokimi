// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/rules"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// runtime are the machineries that the runtime supplies to every check.
var runtime = []spec.MachineryID{spec.NothingBeyondACaller, spec.ClockTheRunAdvances}

// identities returns the identities of the checks of s.
func identities(s derive.Suite) []string {
	ids := make([]string, 0, len(s.Checks))
	for _, c := range s.Checks {
		ids = append(ids, c.Identity)
	}
	return ids
}

// checkOf returns the check of s with the identity id.
func checkOf(t *testing.T, s derive.Suite, id string) derive.Check {
	t.Helper()
	i := slices.IndexFunc(s.Checks, func(c derive.Check) bool { return c.Identity == id })
	assert.True(t, i >= 0, "the suite has the check "+id)
	return s.Checks[i]
}

// claimOn returns a claim of the suite with the licence l, the gates and the
// machinery nothing-beyond-a-caller.
func claimOn(id spec.ClaimID, l spec.Licence, gates ...spec.GateID) spec.Claim {
	return spec.Claim{
		ID: id, Generator: "suite", Licences: []spec.Licence{l}, Gates: gates,
		Law: spec.NoValueBesideAnError, Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
		Statement: "{callable} holds", Defect: spec.Defect{Clause: "a {declaration} that breaks"},
	}
}

func TestDerive(t *testing.T) {
	t.Parallel()

	t.Run("Derive", func(t *testing.T) {
		t.Parallel()

		t.Run("binds the claims of a reader with their identities and keys", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withShape(reader(), "reader")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(
				t,
				identities(got),
				[]string{"Store.Read/smoke", "Store.Read/cancel", "Store.Read/zero-on-error"},
				"the reader gets its checks in the order of the claims",
			)
			check := checkOf(t, got, "Store.Read/zero-on-error")
			assert.Equal(t, check.Key, "store.read/zero-on-error", "the key splits the names into words")
			assert.Equal(
				t,
				check.Statement,
				"Read returns the zero value beside an error",
				"the statement names the callable",
			)
			assert.Equal(
				t,
				check.Clause,
				"a Store whose Read returns a value beside an error",
				"the clause names the declaration",
			)
		})

		t.Run("refuses a licensed claim whose gate fails, with the reason of the gate", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withShape(withRefusedSample(reader()), "reader")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Refusals, []derive.Refusal{{
				Claim: spec.ZeroOnError, Callable: "Read",
				Why:    "eidos derives no sample of a result of Read",
				Remedy: "author a sample of the result type with the directive sample",
			}}, "the claim refuses with the reason and the remedy of the gate")
		})

		t.Run("leaves a claim of every callable silent when its gate fails", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withShape(pure(), "computation")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.False(
				t,
				slices.ContainsFunc(got.Refusals, func(r derive.Refusal) bool { return r.Claim == spec.Cancel }),
				"a callable without a context neither binds nor refuses cancel",
			)
		})

		t.Run("withdraws a claim from a callable with its mixin unless", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withMixins(withShape(reader(), "reader"), "default-on-error")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.False(
				t,
				slices.Contains(identities(got), "Store.Read/zero-on-error"),
				"the default withdraws the claim",
			)
			assert.Equal(t, got.Unclaimed, []derive.Unclaimed{{
				Classification: "mixin:default-on-error", Callable: "Read", Stance: spec.Blocked,
				Note: "The mixin declares no default value for a check to compare.",
			}}, "the mixin has its stance")
		})

		t.Run(
			"lists the classification of a claim whose machinery the run does not supply under the stance needs",
			func(t *testing.T) {
				t.Parallel()
				get := withParam(withMixins(getter(), "ttl"), "ttl", "put", idOf("Put"))
				decl := []derive.Facts{get, putter()}

				got, err := derive.Derive(decl, spec.Claims(), nil, []spec.MachineryID{spec.NothingBeyondACaller})

				assert.NoError(t, err, "the declaration derives")
				assert.False(t, slices.Contains(identities(got), "Store.Get/expires"), "the claim needs a clock")
				assert.True(t, slices.Contains(got.Unclaimed, derive.Unclaimed{
					Classification: "mixin:ttl", Callable: "Get", Stance: spec.Needs,
					Note: "The claim expires needs clock-the-run-advances.",
				}), "the mixin needs the machinery")
			},
		)

		t.Run("binds a claim with a partner and names the partner in the statement", func(t *testing.T) {
			t.Parallel()
			get := withParam(withMixins(getter(), "read-after-write"), "read-after-write", "write", idOf("Put"))
			decl := []derive.Facts{get, putter()}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			check := checkOf(t, got, "Store.Get/read-after-write")
			assert.Equal(t, check.Statement, "Get returns what Put stored before it", "the statement names the partner")
			assert.Equal(t, check.Pairs["write"].Partner.Callable, "Put", "the check has its partner")
		})

		t.Run("refuses a claim whose partner the classification does not name", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withMixins(getter(), "read-after-write"), putter()}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.True(t, slices.ContainsFunc(got.Refusals, func(r derive.Refusal) bool {
				return r.Claim == spec.ReadAfterWrite &&
					r.Why == "no callable of the declaration is the partner write of Get"
			}), "the claim refuses with the gate partner-found")
		})

		t.Run("refuses a claim whose input reads a param that the classification does not set", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withMixins(getter(), "not-found")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.True(t, slices.Contains(got.Refusals, derive.Refusal{
				Claim: spec.Miss, Callable: "Get",
				Why:    "the classification of Get sets no param sentinel",
				Remedy: "set the param sentinel in the directive of Get",
			}), "the claim refuses with the name of the param")
		})

		t.Run("lists a classification without a claim with its stance", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{withShape(putter(), "writer")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Unclaimed, []derive.Unclaimed{{
				Classification: "shape:writer", Callable: "Put", Stance: spec.Silent,
				Note: "A write is an input of the other checks, and its result is an error alone.",
			}}, "the shape writer is silent")
		})

		t.Run("adds an edge where a claim entails another that binds on the callable", func(t *testing.T) {
			t.Parallel()
			closer := derive.Facts{
				Declaration: "Store", Callable: "Close",
				Signature: signature("Close", false, nil, nil, rules.ErrorsNone),
			}
			decl := []derive.Facts{withShape(closer, "void-lifecycle")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(
				t,
				got.Edges,
				[]derive.Edge{{Deeper: "Store.Close/repeat-survives", Shallower: "Store.Close/smoke"}},
				"the repeat entails the smoke",
			)
		})

		t.Run("adds the id of the instance to a claim that binds twice on one callable", func(t *testing.T) {
			t.Parallel()
			get := getter()
			get.Class.Contracts = []shape.Membership{
				{Contract: "cache", Role: "cache", ID: "a"},
				{Contract: "cache", Role: "cache", ID: "b"},
			}
			backing := reader()
			backing.Class.Contracts = []shape.Membership{
				{Contract: "cache", Role: "backing", ID: "a"},
				{Contract: "cache", Role: "backing", ID: "b"},
			}
			decl := []derive.Facts{get, backing, withShape(putter(), "writer")}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.True(t, slices.Contains(identities(got), "Store.Get/cache#a"), "the first instance has its id")
			assert.True(t, slices.Contains(identities(got), "Store.Get/cache#b"), "the second instance has its id")
		})

		t.Run("refuses a licensed claim and leaves smoke silent for an input without samples", func(t *testing.T) {
			t.Parallel()
			hash := withShape(pure(), "computation")
			hash.Inputs = append([]derive.Value(nil), hash.Inputs...)
			hash.Inputs[0].Samples[0] = rules.Refused(rules.RefusedNoLiteral)
			decl := []derive.Facts{hash}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Empty(t, got.Checks, "no claim binds without samples of the input")
			assert.Equal(t, got.Refusals, []derive.Refusal{{
				Claim: spec.Deterministic, Callable: "Hash",
				Why:    "eidos derives no sample of the parameter key of Hash",
				Remedy: "author a sample of the parameter key with the directive sample",
			}}, "the claim of the shape refuses, and smoke is silent")
		})

		t.Run("licenses a claim with its param and withdraws the licence of a conflict", func(t *testing.T) {
			t.Parallel()
			commit := withParam(withMixins(getter(), "retry-succeeds"), "retry-succeeds", "attempts", int64(3))
			commit.Class.Contracts = []shape.Membership{{Contract: "tx", Role: "commit"}}
			begin := withRole(putter(), "tx", "begin")
			decl := []derive.Facts{commit, begin}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.True(
				t,
				slices.Contains(identities(got), "Store.Get/retry-succeeds-tx"),
				"the commit of a tx binds the claim with",
			)
			assert.False(
				t,
				slices.Contains(identities(got), "Store.Get/retry-succeeds"),
				"the commit withdraws the claim of its conflict",
			)
		})

		t.Run("leaves a claim whose licence needs a param that the classification lacks", func(t *testing.T) {
			t.Parallel()
			get := withParam(withMixins(getter(), "retry-succeeds"), "retry-succeeds", "other", int64(3))
			decl := []derive.Facts{get}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.False(
				t,
				slices.Contains(identities(got), "Store.Get/retry-succeeds"),
				"the licence needs the param attempts",
			)
			assert.False(
				t,
				slices.ContainsFunc(got.Refusals, func(r derive.Refusal) bool { return r.Claim == spec.RetrySucceeds }),
				"a claim without its licence does not refuse",
			)
		})

		t.Run("names the param of a machinery of the consumer that the declaration does not state", func(t *testing.T) {
			t.Parallel()
			table := []spec.Claim{{
				ID: "agrees", Generator: "suite", Law: spec.ReadReturnsTheWriteBeforeIt,
				Licences:    []spec.Licence{{Shape: "reader"}},
				Machineries: []spec.MachineryID{spec.ReferenceImplementation},
			}}
			decl := []derive.Facts{withShape(reader(), "reader")}

			got, err := derive.Derive(decl, table, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Unclaimed, []derive.Unclaimed{{
				Classification: "shape:reader", Callable: "Read", Stance: spec.Needs,
				Note: "The claim agrees needs the param machineries=reference of the directive suite.",
			}}, "the note names the param that states the machinery")
		})

		t.Run("binds only the claims of the generator suite", func(t *testing.T) {
			t.Parallel()
			table := []spec.Claim{claimOn("own", spec.Licence{Callable: true})}
			table[0].Generator = "enum"

			got, err := derive.Derive([]derive.Facts{reader()}, table, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Empty(t, got.Checks, "a claim of enum binds to no callable")
		})

		t.Run("leaves out the stance of a classification that a claim binds or withdraws", func(t *testing.T) {
			t.Parallel()
			bound := claimOn("bound", spec.Licence{Shape: "writer"})
			withdrawn := claimOn("withdrawn", spec.Licence{Mixin: "deprecated"})
			withdrawn.Unless = &spec.Unless{Mixin: "deprecated"}
			decl := []derive.Facts{withMixins(withShape(putter(), "writer"), "deprecated")}

			got, err := derive.Derive(decl, []spec.Claim{bound, withdrawn}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, identities(got), []string{"Store.Put/bound"}, "the claim of the writer binds")
			assert.Empty(t, got.Unclaimed, "neither the bound writer nor the withdrawn mixin has a stance")
		})

		t.Run("lists the stance of each role of a contract once", func(t *testing.T) {
			t.Parallel()
			stats := getter()
			stats.Class.Contracts = []shape.Membership{
				{Contract: "pool", Role: "stats", ID: "a"},
				{Contract: "pool", Role: "stats", ID: "b"},
			}

			got, err := derive.Derive([]derive.Facts{stats}, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Unclaimed, []derive.Unclaimed{{
				Classification: "contract:pool/stats", Callable: "Get", Stance: spec.Gap,
				Note: "No claim states that the stats of a pool count each loan and each return.",
			}}, "the role has its stance once")
		})

		t.Run("binds a claim of one contract instance without the id of the instance", func(t *testing.T) {
			t.Parallel()
			c := claimOn("one", spec.Licence{Contract: "cache", Role: "cache"})
			get := getter()
			get.Class.Contracts = []shape.Membership{{Contract: "cache", Role: "cache", ID: "a"}}
			other := withRole(withRole(reader(), "cache", "backing"), "pool", "cache")

			got, err := derive.Derive([]derive.Facts{get, other}, []spec.Claim{c}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(
				t,
				identities(got),
				[]string{"Store.Get/one"},
				"only the role cache of the contract cache binds",
			)
			assert.Equal(t, got.Checks[0].Key, "store.get/one", "the key has no id")
		})

		t.Run("adds the id of the instance to the key of a claim that binds twice", func(t *testing.T) {
			t.Parallel()
			c := claimOn("twice", spec.Licence{Contract: "cache", Role: "cache"})
			get := getter()
			get.Class.Contracts = []shape.Membership{
				{Contract: "cache", Role: "cache", ID: "a"},
				{Contract: "cache", Role: "cache", ID: "b"},
			}

			got, err := derive.Derive([]derive.Facts{get}, []spec.Claim{c}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Checks[1].Key, "store.get/twice#b", "the key has the id of the instance")
		})

		t.Run("binds a claim with of a callable that has both classifications alone", func(t *testing.T) {
			t.Parallel()
			c := claimOn("both", spec.Licence{
				Shape: "reader", With: &spec.Classification{Mixin: "causal"},
			})
			decl := []derive.Facts{withShape(reader(), "reader"), withMixins(withShape(getter(), "reader"), "causal")}

			got, err := derive.Derive(decl, []spec.Claim{c}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, identities(got), []string{"Store.Get/both"}, "only the reader with the mixin causal binds")
		})

		t.Run("matches a contract without a role and a classification of another form", func(t *testing.T) {
			t.Parallel()
			anyRole := claimOn("any-role", spec.Licence{
				Shape: "reader", With: &spec.Classification{Contract: "cache"},
			})
			mixin := claimOn("mixin", spec.Licence{Mixin: "causal"})
			get := withRole(withShape(getter(), "reader"), "cache", "backing")
			other := withRole(withMixins(withShape(reader(), "lookup"), "pure"), "pool", "get")

			got, err := derive.Derive([]derive.Facts{get, other}, []spec.Claim{anyRole, mixin}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(
				t,
				identities(got),
				[]string{"Store.Get/any-role"},
				"a role of the contract cache licenses the claim",
			)
		})

		t.Run("binds a claim with a role of a contract only on a callable that plays the role", func(t *testing.T) {
			t.Parallel()
			c := claimOn("commit", spec.Licence{
				Shape: "reader", With: &spec.Classification{Contract: "tx", Role: "commit"},
			})
			begin := withRole(withShape(getter(), "reader"), "tx", "begin")
			otherContract := withRole(withShape(reader(), "reader"), "saga", "commit")
			commit := withRole(withShape(lookup(), "reader"), "tx", "commit")

			got, err := derive.Derive([]derive.Facts{begin, otherContract, commit}, []spec.Claim{c}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(
				t,
				identities(got),
				[]string{"Store.Lookup/commit"},
				"only the role commit of the contract tx licenses the claim",
			)
		})

		t.Run("judges a gate with a condition and leaves a gate of the partners to the partners", func(t *testing.T) {
			t.Parallel()
			c := claimOn("gated", spec.Licence{Callable: true}, spec.PartnerFound, spec.Failure)
			decl := []derive.Facts{reader(), pure()}

			got, err := derive.Derive(decl, []spec.Claim{c}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, identities(got), []string{"Store.Read/gated"}, "the callable that can fail binds")
		})

		t.Run("marks a check shallow when its machineries lack the depth of its law", func(t *testing.T) {
			t.Parallel()
			get := withParam(
				withMixins(withShape(getter(), "reader"), "read-after-write"),
				"read-after-write",
				"write",
				idOf("Put"),
			)
			decl := []derive.Facts{get, putter()}

			got, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.True(
				t,
				checkOf(t, got, "Store.Get/read-after-write").Shallow,
				"a read after a write is deeper over sequences",
			)
			assert.False(t, checkOf(t, got, "Store.Get/cancel").Shallow, "cancel witnesses its law at its depth")
		})

		t.Run("adds an edge only to an entailed claim that binds on the callable", func(t *testing.T) {
			t.Parallel()
			deeper := claimOn("deeper", spec.Licence{Callable: true})
			deeper.Entails = []spec.ClaimID{"absent", "shallower"}
			shallower := claimOn("shallower", spec.Licence{Callable: true})
			other := claimOn("other", spec.Licence{Callable: true})

			got, err := derive.Derive([]derive.Facts{reader()}, []spec.Claim{deeper, other, shallower}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Edges, []derive.Edge{{Deeper: "Store.Read/deeper", Shallower: "Store.Read/shallower"}},
				"the edge joins the bound shallower claim")
		})

		t.Run("names the check of a function by the function alone", func(t *testing.T) {
			t.Parallel()
			slug := pure()
			slug.Declaration, slug.Callable = "", "Slug"
			c := claimOn("own", spec.Licence{Callable: true})

			got, err := derive.Derive([]derive.Facts{slug}, []spec.Claim{c}, nil, runtime)

			assert.NoError(t, err, "the declaration derives")
			assert.Equal(t, got.Checks[0].Identity, "Slug/own", "the identity has no declaration")
			assert.Equal(t, got.Checks[0].Key, "slug/own", "the key has no declaration")
			assert.Equal(t, got.Checks[0].Clause, "a Slug that breaks", "the function is its own declaration")
		})

		t.Run("returns an error for two callables with one key", func(t *testing.T) {
			t.Parallel()
			a, b := pure(), pure()
			a.Callable, b.Callable = "ReadAll", "Read_All"
			decl := []derive.Facts{a, b}

			_, err := derive.Derive(decl, spec.Claims(), nil, runtime)

			assert.ErrorIs(t, err, derive.ErrKey, "two callables share the key store.read-all")
		})
	})
}
