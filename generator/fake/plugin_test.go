// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/fake"
	"go.dokimi.dev/dokimi/internal/gentest"
	"go.dokimi.dev/dokimi/spec"
)

// The packages of the Go tree that this spec alone reads.
const (
	refusedPackage = "example.com/shop/refused"
	partialPackage = "example.com/shop/partial"
	keysPackage    = "example.com/shop/lib/internal/keys"
	meterPackage   = "example.com/shop/lib/meter"
)

// method is a method of a fake with the type of each result spelled.
type method struct {
	name, receiver, context, key, value, present string
	effect                                       spec.Effect
	results                                      []string
}

// methodsOf returns each method of the model m with the local and the
// spelling of the type of each result, and error for an error.
func methodsOf(m *fake.Model) []method {
	out := make([]method, 0, len(m.Methods))
	for _, c := range m.Methods {
		results := make([]string, 0, len(c.Results))
		for _, r := range c.Results {
			if r.Error {
				results = append(results, "error")
				continue
			}
			results = append(results, r.Local+" "+r.Type.Spelling)
		}
		out = append(out, method{
			name: c.Name, receiver: c.Receiver, context: c.Context, key: c.Key, value: c.Value,
			present: c.Present, effect: c.Effect, results: results,
		})
	}
	return out
}

// spelled returns the spelling of t, and the empty string for nil.
func spelled(t *emit.TypeRef) string {
	if t == nil {
		return ""
	}
	return t.Spelling
}

// spellings returns the spelling of each type of ts.
func spellings(ts []*emit.TypeRef) []string {
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		out = append(out, t.Spelling)
	}
	return out
}

func TestNew(t *testing.T) {
	t.Parallel()

	{
		tests := []struct {
			name        string
			givePkg     string
			giveSubject string
			wantNames   []string
			wantTypes   []string
			wantForm    spec.ReferenceFormID
			wantState   spec.State
			want        []method
		}{
			{
				name:    "records a fake of a persister that keys its values by the field of the written value",
				givePkg: storePackage, giveSubject: "Store",
				wantNames: []string{"StoreFake", "NewStoreFake", "ErrStoreFakeNotFound"},
				wantTypes: []string{"string", "Value", "Store"}, wantForm: "persister", wantState: spec.StateMap,
				want: []method{
					{
						name: "Put", receiver: "f", context: "ctx", key: "v.Key", value: "v",
						effect: spec.EffectStore, results: []string{"r0 Value", "error"},
					},
					{
						name: "Get", receiver: "f", context: "ctx", key: "key", present: "present",
						effect: spec.EffectLoad, results: []string{"r0 Value", "error"},
					},
				},
			},
			{
				name:    "records a fake of a cache that keys its values by the parameter of the writer",
				givePkg: cachePackage, giveSubject: "Cache",
				wantNames: []string{"CacheFake", "NewCacheFake", "ErrNoValue"},
				wantTypes: []string{"string", "Value", "Cache"}, wantForm: "cache", wantState: spec.StateMap,
				want: []method{
					{
						name: "Store", receiver: "f", context: "ctx", key: "key", value: "v",
						effect: spec.EffectStore, results: []string{"error"},
					},
					{
						name: "Lookup", receiver: "f", context: "ctx", key: "key", present: "present",
						effect: spec.EffectLoad, results: []string{"r0 Value", "error"},
					},
					{
						name: "Fetch", receiver: "f", context: "ctx", key: "key", present: "present",
						effect: spec.EffectLoad, results: []string{"r0 Value", "error"},
					},
				},
			},
			{
				name: "records a fake of an appender without a key and a sentinel", givePkg: logPackage,
				giveSubject: "Log", wantNames: []string{"LogFake", "NewLogFake", ""},
				wantTypes: []string{"", "Entry", "Log"}, wantForm: "appender", wantState: spec.StateList,
				want: []method{{
					name: "Append", receiver: "f", context: "ctx", value: "e", effect: spec.EffectAppend,
					results: []string{"r0 int64", "error"},
				}},
			},
			{
				name: "records a fake of a store that takes a context and returns no error", givePkg: quietPackage,
				giveSubject: "Quiet", wantNames: []string{"QuietFake", "NewQuietFake", "ErrQuietFakeNotFound"},
				wantTypes: []string{"string", "Value", "Quiet"}, wantForm: "persister", wantState: spec.StateMap,
				want: []method{
					{
						name: "Put", receiver: "f", key: "v.Key", value: "v", effect: spec.EffectStore,
						results: []string{},
					},
					{
						name: "Get", receiver: "f", context: "ctx", key: "key", present: "present",
						effect: spec.EffectLoad, results: []string{"r0 Value", "error"},
					},
				},
			},
			{
				name: "records a fake of a generic cache over its type parameters", givePkg: genericPackage,
				giveSubject: "Cache", wantNames: []string{"CacheFake", "NewCacheFake", "ErrCacheFakeNotFound"},
				wantTypes: []string{"K", "V", "Cache"}, wantForm: "cache", wantState: spec.StateMap,
				want: []method{
					{
						name: "Store", receiver: "f", context: "ctx", key: "key", value: "v",
						effect: spec.EffectStore, results: []string{"error"},
					},
					{
						name: "Lookup", receiver: "f", context: "ctx", key: "key", present: "present",
						effect: spec.EffectLoad, results: []string{"r0 V", "error"},
					},
					{
						name: "Fetch", receiver: "f", context: "ctx", key: "key", present: "present",
						effect: spec.EffectLoad, results: []string{"r0 V", "error"},
					},
				},
			},
			{
				name: "names the receiver and the locals apart from the parameters", givePkg: clashPackage,
				giveSubject: "Shadow", wantNames: []string{"ShadowFake", "NewShadowFake", "ErrShadowFakeNotFound"},
				wantTypes: []string{"string", "Value", "Shadow"}, wantForm: "persister", wantState: spec.StateMap,
				want: []method{
					{
						name: "Put", receiver: "f0", context: "f", key: "r0.Present", value: "r0",
						effect: spec.EffectStore, results: []string{"r00 Value", "error"},
					},
					{
						name: "Get", receiver: "f0", context: "f", key: "present", present: "present0",
						effect: spec.EffectLoad, results: []string{"r0 Value", "error"},
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, recorded := generate(t, tt.givePkg).models(t)[tt.giveSubject]

				assert.True(t, recorded, "fake records the fake of "+tt.giveSubject)
				m := &got.Data
				assert.Equal(t, []string{m.Fake, m.New, m.Absent}, tt.wantNames,
					"the fake, its constructor and its sentinel join the words of fake to the declaration")
				assert.Equal(t, []string{spelled(m.Key), spelled(m.Value), spelled(m.Guard)}, tt.wantTypes,
					"the fake keeps the values of the form under their keys")
				assert.Equal(t, []any{m.Subject, m.Form, m.State}, []any{tt.giveSubject, tt.wantForm, tt.wantState},
					"the reference form of the contract of the declaration covers its callables")
				assert.Equal(t, methodsOf(m), tt.want, "each method has the effect of its callable")
				assert.Equal(t, got.Decl.Name, tt.wantNames[0], "the record of the model is the struct of the fake")
			})
		}
	}

	t.Run("records the witnesses of a generic declaration", func(t *testing.T) {
		t.Parallel()

		got := generate(t, genericPackage).models(t)["Cache"].Data

		assert.Equal(t, spellings(got.Witnesses), []string{"int", "string"},
			"the fake has the witnesses of the type parameters of Cache")
		assert.Equal(t, spellings(got.Guard.Args), []string{"int", "string"},
			"the compile-time check of the fake instantiates Cache at its witnesses")
	})

	t.Run("leaves the imports of a declaration and its witnesses to the layout where the plan previews no package",
		func(t *testing.T) {
			t.Parallel()

			got := phases(t, nil, keysPackage, meterPackage).models(t)

			expect.Contains(t, got, "Keys", "fake records the fake of Keys, and the layout refuses it")
			expect.Contains(t, got, "Meter", "fake records the fake of Meter, and the layout refuses its witness")
		})

	t.Run("records the struct of a generic fake with the type parameters of its declaration", func(t *testing.T) {
		t.Parallel()

		got := generate(t, genericPackage).models(t)["Cache"].Decl

		assert.NotNil(t, got, "fake records the struct of the fake")
		assert.Length(t, got.TypeParams, 2, "the fake restates the type parameters K and V")
		receiver := got.Methods.Items()[0].Receives
		assert.Equal(t, []string{receiver.Spelling, spellings(receiver.Args)[0], spellings(receiver.Args)[1]},
			[]string{"CacheFake", "V", "K"}, "the receiver of a method restates the type parameters")
	})

	{
		tests := []struct {
			name        string
			givePkg     string
			giveSubject string
			want        []string
			wantArgs    []string
		}{
			{
				name: "records the test of a fake that applies the suite of its declaration", givePkg: storePackage,
				giveSubject: "Store", want: []string{"TestStoreFake", "StoreFake", "RunStoreSuite", "NewStoreFake"},
				wantArgs: []string{},
			},
			{
				name: "records the test of a generic fake at the witnesses of its declaration", givePkg: genericPackage,
				giveSubject: "Cache", want: []string{"TestCacheFake", "CacheFake", "RunCacheSuite", "NewCacheFake"},
				wantArgs: []string{"int", "string"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				r := generate(t, tt.givePkg)
				got, recorded := r.tests(t)[tt.giveSubject]

				assert.True(t, recorded, "fake records the test of the fake of "+tt.giveSubject)
				refs := got.Decl.Fields.Items()
				assert.Equal(t, append([]string{got.Data.Test, got.Data.Fake}, gentest.Referents(got.Decl)...), tt.want,
					"the test refers to the function Run of the suite and to the constructor of the fake")
				id := r.f.Find(t, symbol.KindInterface, "", tt.giveSubject).Identity()
				for _, ref := range refs {
					expect.Equal(t, []any{ref.Type.Target, ref.Type.Emitted}, []any{id, true},
						"the test refers to a declaration that a record of "+tt.giveSubject+" declares")
				}
				expect.Equal(t, spellings(refs[1].Type.Args), tt.wantArgs,
					"the test instantiates the constructor at the witnesses")
				expect.Equal(t, got.Decl.Visibility, symbol.VisibilityPackage, "the record of the test is unexported")
			})
		}
	}

	t.Run("declares the constructor and the sentinel of an absent key beside the fake", func(t *testing.T) {
		t.Parallel()

		key := func(m fake.Model) string { return m.Subject }
		assert.Equal(t, gentest.Beside(generate(t, storePackage).res, key)["Store"],
			[]string{"NewStoreFake", "ErrStoreFakeNotFound"}, "the record of a map declares its sentinel")
		assert.Equal(t, gentest.Beside(generate(t, logPackage).res, key)["Log"], []string{"NewLogFake"},
			"the record of a list declares no sentinel")
	})

	{
		tests := []struct {
			name         string
			givePkgs     []string
			giveOwner    string
			giveName     string
			wantCode     diag.Code
			wantSeverity diag.Severity
			wantMsg      string
		}{
			{
				name: "reports NoSuite at an interface without the directive suite", givePkgs: []string{refusedPackage},
				giveName: "Unsuited", wantCode: fake.NoSuite, wantSeverity: diag.SeverityError,
				wantMsg: "Unsuited has no directive suite, so it has no fake",
			},
			{
				name:     "reports BadAbsent at an interface whose param absent reads another field",
				givePkgs: []string{refusedPackage}, giveName: "Misnamed", wantCode: fake.BadAbsent,
				wantSeverity: diag.SeverityError,
				wantMsg: `the param absent of Misnamed is not the name template of a sentinel: the name template ` +
					`"Err{{.Method}}" has an action other than {{.Name}}, so it has no fake`,
			},
			{
				name:     "reports Refused at a callable with a mixin",
				givePkgs: []string{refusedPackage}, giveOwner: "Expiring", giveName: "Get", wantCode: fake.Refused,
				wantSeverity: diag.SeverityWarning,
				wantMsg: "Expiring.Get has the mixin ttl, so Expiring has no fake: give Get a role of a contract " +
					"with a reference form, or leave Expiring without the directive fake",
			},
			{
				name:     "reports Refused at a callable with a shape and no role",
				givePkgs: []string{refusedPackage}, giveOwner: "Counted", giveName: "Since", wantCode: fake.Refused,
				wantSeverity: diag.SeverityWarning,
				wantMsg: "Counted.Since has the shape aggregator, so Counted has no fake: give Since a role of a " +
					"contract with a reference form, or leave Counted without the directive fake",
			},
			{
				name:     "reports Refused at a load whose key selects only an unexported field of the written value",
				givePkgs: []string{refusedPackage}, giveOwner: "Private", giveName: "Get", wantCode: fake.Refused,
				wantSeverity: diag.SeverityWarning,
				wantMsg: "Private.Get takes a key that selects nothing of the value that Put stores, so Private has " +
					"no fake: give Get a role of a contract with a reference form, or leave Private without the " +
					"directive fake",
			},
			{
				name:     "reports Refused at a load whose key selects nothing of the written value",
				givePkgs: []string{refusedPackage}, giveOwner: "Unkeyed", giveName: "Get", wantCode: fake.Refused,
				wantSeverity: diag.SeverityWarning,
				wantMsg: "Unkeyed.Get takes a key that selects nothing of the value that Put stores, so Unkeyed has " +
					"no fake: give Get a role of a contract with a reference form, or leave Unkeyed without the " +
					"directive fake",
			},
			{
				name:     "reports Unimplementable at an interface that its package does not export",
				givePkgs: []string{refusedPackage}, giveName: "hidden", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError, wantMsg: "hidden is not exported, so it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface that the file of its fake cannot import",
				givePkgs: []string{keysPackage}, giveName: "Keys", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg: "Keys cannot be imported by the file of its fake in the package example.com/shop/keystest: " +
					"golang: example.com/shop/lib/internal/keys is internal to example.com/shop/lib, and the " +
					"package example.com/shop/keystest is outside it, so it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface that embeds an interface of no loaded package",
				givePkgs: []string{partialPackage}, giveName: "Partial", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError, wantMsg: "Partial has incomplete members, so it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface with a method that its package does not export",
				givePkgs: []string{refusedPackage}, giveName: "Sealed", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg:      "Sealed has the method get that its package does not export, so it has no fake",
			},
			{
				name: "reports Unimplementable at an interface without a method", givePkgs: []string{refusedPackage},
				giveName: "Nothing", wantCode: fake.Unimplementable, wantSeverity: diag.SeverityError,
				wantMsg: "Nothing has no method, so it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface whose type parameter has no witness",
				givePkgs: []string{refusedPackage}, giveName: "Tagged", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg: "Tagged has type parameters without witnesses that the files of its fake can refer to, so " +
					"it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface whose witness the files of its fake cannot import",
				givePkgs: []string{meterPackage, keysPackage}, giveName: "Meter", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg: "Meter has type parameters without witnesses that the files of its fake can refer to, so " +
					"it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface with a key that a Go map cannot key",
				givePkgs: []string{refusedPackage}, giveName: "Bytes", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg: "Bytes keys its values by the parameter key of the type []byte, which a Go map cannot " +
					"key, so it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface with a key of a map type",
				givePkgs: []string{refusedPackage}, giveName: "Mapped", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg: "Mapped keys its values by the parameter key of the type map[string]string, which a Go map " +
					"cannot key, so it has no fake",
			},
			{
				name:     "reports Unimplementable at an interface with a key of a type parameter of any type",
				givePkgs: []string{refusedPackage}, giveName: "Loose", wantCode: fake.Unimplementable,
				wantSeverity: diag.SeverityError,
				wantMsg: "Loose keys its values by the parameter key of the type K, whose constraint permits types " +
					"that a Go map cannot key, so it has no fake",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				r := generate(t, tt.givePkgs...)
				kind := symbol.KindInterface
				if tt.giveOwner != "" {
					kind = symbol.KindMethod
				}
				at := r.f.Find(t, kind, tt.giveOwner, tt.giveName)
				got := gentest.Finding(t, r.res, at)

				assert.Equal(t, []any{got.Code, got.Severity, got.Msg}, []any{tt.wantCode, tt.wantSeverity, tt.wantMsg},
					"fake reports the finding with the cause and the remedy")
				subject := tt.giveName
				if tt.giveOwner != "" {
					subject = tt.giveOwner
				}
				_, recorded := r.models(t)[subject]
				expect.False(t, recorded, "fake records no fake of "+subject)
				_, tested := r.tests(t)[subject]
				expect.False(t, tested, "fake records no test of "+subject)
			})
		}
	}
}
