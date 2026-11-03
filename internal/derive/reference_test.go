// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// effected is the effect of one callable of a reference, with the binding of
// its key and of its value.
type effected struct {
	callable string
	effect   spec.Effect
	key      derive.Argument
	value    int
}

// effectsOf returns the effect of each callable of r, in order.
func effectsOf(r derive.Reference) []effected {
	out := make([]effected, 0, len(r.Callables))
	for _, c := range r.Callables {
		out = append(out, effected{callable: c.Facts.Callable, effect: c.Effect, key: c.Key, value: c.Value})
	}
	return out
}

// reading returns the facts of Store.Get(ctx, id string) (Item, error) with
// the role reader of the contract persister.
func reading() derive.Facts { return withRole(getter(), "persister", "reader") }

// putting returns the facts of Store.Put(ctx, it Item) error with the role
// writer of the contract persister.
func putting() derive.Facts { return withRole(putter(), "persister", "writer") }

// looking returns the facts of Store.Lookup(ctx, id string) (Item, error) with
// the role cache of the contract cache.
func looking() derive.Facts { return withRole(renamed(getter(), "Lookup"), "cache", "cache") }

// fetching returns the facts of Store.Fetch(ctx, id string) (Item, error) with
// the role backing of the contract cache.
func fetching() derive.Facts { return withRole(renamed(getter(), "Fetch"), "cache", "backing") }

// setting returns the facts of Store.Set(ctx, id string, it Item) error with
// the shape composite-writer.
func setting() derive.Facts { return withShape(keyedPutter(), shape.CompositeWriter) }

// appending returns the facts of Store.Append(ctx, it Item) (int64, error)
// with the role fn of the contract appender.
func appending() derive.Facts {
	inputs := []derive.Value{item("it")}
	return withRole(derive.Facts{
		Declaration: "Store", Callable: "Append", ID: idOf("Append"),
		Signature: signature(
			"Append", true, inputs, []rules.ReturnRole{rules.ReturnValue, rules.ReturnError}, rules.ErrorsLastReturn,
		),
		Inputs:  inputs,
		Results: []derive.Value{position(rules.ScalarInt)},
	}, "appender", "fn")
}

// position returns a value of an integer type of the class.
func position(class rules.ScalarClass) derive.Value {
	return derive.Value{Shape: rules.TypeShape{Form: symbol.FormScalar, Spelling: "int64", Class: class, Bits: 64}}
}

// withResults returns f with the results.
func withResults(f derive.Facts, results ...derive.Value) derive.Facts {
	f.Results = results
	return f
}

// withKeys returns f with the inputs in place of its inputs.
func withKeys(f derive.Facts, inputs ...derive.Value) derive.Facts {
	f.Inputs = inputs
	f.Signature.Params = f.Signature.Params[:1]
	for _, in := range inputs {
		f.Signature.Params = append(f.Signature.Params, rules.ParamView{Name: in.Name, Role: rules.ParamInput})
	}
	return f
}

// withInstance returns f with the id of the instance of its first contract.
func withInstance(f derive.Facts, id string) derive.Facts {
	f.Class.Contracts = append([]shape.Membership(nil), f.Class.Contracts...)
	f.Class.Contracts[0].ID = id
	return f
}

func TestReference(t *testing.T) {
	t.Parallel()

	t.Run("ReferenceOf", func(t *testing.T) {
		t.Parallel()

		byField := derive.Argument{Name: "id", Kind: derive.ArgField, Field: "ID"}
		byParam := derive.Argument{Name: "id", Kind: derive.ArgShared, Param: "id"}
		{
			tests := []struct {
				name     string
				give     []derive.Facts
				wantForm spec.ReferenceFormID
				want     []effected
			}{
				{
					name: "covers a persister by a store and a load of the field of the written value",
					give: []derive.Facts{putting(), reading()}, wantForm: "persister",
					want: []effected{
						{callable: "Put", effect: spec.EffectStore, key: byField, value: 0},
						{callable: "Get", effect: spec.EffectLoad, key: byField, value: -1},
					},
				},
				{
					name: "covers a store that returns the value that it stores",
					give: []derive.Facts{withResults(putting(), item("")), reading()}, wantForm: "persister",
					want: []effected{
						{callable: "Put", effect: spec.EffectStore, key: byField, value: 0},
						{callable: "Get", effect: spec.EffectLoad, key: byField, value: -1},
					},
				},
				{
					name:     "covers a role with the shape that the catalog detects on its signature",
					give:     []derive.Facts{withShape(putting(), "answering-writer"), withShape(reading(), "reader")},
					wantForm: "persister",
					want: []effected{
						{callable: "Put", effect: spec.EffectStore, key: byField, value: 0},
						{callable: "Get", effect: spec.EffectLoad, key: byField, value: -1},
					},
				},
				{
					name: "covers a cache by two loads and a writer of a shared parameter",
					give: []derive.Facts{setting(), looking(), fetching()}, wantForm: "cache",
					want: []effected{
						{callable: "Set", effect: spec.EffectStore, key: byParam, value: 1},
						{callable: "Lookup", effect: spec.EffectLoad, key: byParam, value: -1},
						{callable: "Fetch", effect: spec.EffectLoad, key: byParam, value: -1},
					},
				},
				{
					name: "binds the key of each store and the key of a load by the first store",
					give: []derive.Facts{setting(), withShape(putter(), shape.Writer), looking()}, wantForm: "cache",
					want: []effected{
						{callable: "Set", effect: spec.EffectStore, key: byParam, value: 1},
						{callable: "Put", effect: spec.EffectStore, key: byField, value: 0},
						{callable: "Lookup", effect: spec.EffectLoad, key: byParam, value: -1},
					},
				},
				{
					name:     "covers two loads that select one field of the written value by two names",
					give:     []derive.Facts{putting(), reading(), withKeys(renamed(reading(), "Find"), text("ID"))},
					wantForm: "persister",
					want: []effected{
						{callable: "Put", effect: spec.EffectStore, key: byField, value: 0},
						{callable: "Get", effect: spec.EffectLoad, key: byField, value: -1},
						{
							callable: "Find", effect: spec.EffectLoad, value: -1,
							key: derive.Argument{Name: "ID", Kind: derive.ArgField, Field: "ID"},
						},
					},
				},
				{
					name: "covers an appender by an append of its one input", give: []derive.Facts{appending()},
					wantForm: "appender",
					want:     []effected{{callable: "Append", effect: spec.EffectAppend, value: 0}},
				},
				{
					name:     "covers an append whose position is unsigned",
					give:     []derive.Facts{withResults(appending(), position(rules.ScalarUint))},
					wantForm: "appender",
					want:     []effected{{callable: "Append", effect: spec.EffectAppend, value: 0}},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, refused := derive.ReferenceOf(tt.give)

					assert.Nil(t, refused, "a reference form covers every callable")
					assert.Equal(t, got.Form.ID, tt.wantForm, "the form of the contract covers the callables")
					assert.Equal(t, effectsOf(got), tt.want, "each callable has the effect of its role or of a writer")
				})
			}
		}

		{
			byName := withKeys(renamed(reading(), "GetByName"), text("name"))
			tests := []struct {
				name         string
				give         []derive.Facts
				wantCallable string
				wantDetail   string
			}{
				{
					name:         "refuses a callable with a mixin",
					give:         []derive.Facts{putting(), withMixins(reading(), "ttl")},
					wantCallable: "Get", wantDetail: "has the mixin ttl",
				},
				{
					name:         "refuses a callable without a classification",
					give:         []derive.Facts{putting(), reading(), pure()},
					wantCallable: "Hash", wantDetail: "has no classification",
				},
				{
					name:         "refuses a callable with another shape and no role",
					give:         []derive.Facts{putting(), reading(), withShape(pure(), "deleter")},
					wantCallable: "Hash", wantDetail: "has the shape deleter",
				},
				{
					name: "refuses a role of a contract without a reference form",
					give: []derive.Facts{withRole(getter(), "tx", "begin")}, wantCallable: "Get",
					wantDetail: "has the role begin of the contract tx, which has no reference form",
				},
				{
					name:         "refuses roles in two contracts",
					give:         []derive.Facts{putting(), withRole(reading(), "cache", "cache")},
					wantCallable: "Get", wantDetail: "has roles in the contracts persister and cache",
				},
				{
					name:         "refuses a role in the instance of another contract",
					give:         []derive.Facts{putting(), looking()},
					wantCallable: "Lookup", wantDetail: "has a role in another contract instance than Put",
				},
				{
					name:         "refuses a role in another instance of one contract",
					give:         []derive.Facts{withInstance(putting(), "a"), withInstance(reading(), "b")},
					wantCallable: "Get", wantDetail: "has a role in another contract instance than Put",
				},
				{
					name:         "names the first callable with a role in the refusal of another instance",
					give:         []derive.Facts{putting(), reading(), looking()},
					wantCallable: "Lookup", wantDetail: "has a role in another contract instance than Put",
				},
				{
					name:         "refuses a writer of a form without writers",
					give:         []derive.Facts{putting(), reading(), setting()},
					wantCallable: "Set",
					wantDetail:   "has the shape composite-writer, and the reference form persister takes no writer",
				},
				{
					name:         "refuses a writer of a declaration without a role",
					give:         []derive.Facts{setting()},
					wantCallable: "Set",
					wantDetail: "has the shape composite-writer, and no callable of Store has a role of a contract " +
						"with a reference form",
				},
				{
					name: "refuses a map without a store", give: []derive.Facts{looking(), fetching()},
					wantCallable: "Lookup", wantDetail: "loads from a map that no callable of Store stores into",
				},
				{
					name: "refuses a map without a load", give: []derive.Facts{putting()},
					wantCallable: "Put", wantDetail: "stores into a map that no callable of Store loads from",
				},
				{
					name:         "refuses a load whose key selects nothing of the written value",
					give:         []derive.Facts{putting(), withRole(reader(), "persister", "reader")},
					wantCallable: "Read", wantDetail: "takes a key that selects nothing of the value that Put stores",
				},
				{
					name:         "refuses a load of another type than the value of the store",
					give:         []derive.Facts{withShape(marker(), shape.CompositeWriter), looking()},
					wantCallable: "Lookup", wantDetail: "returns another type than the value that Mark stores",
				},
				{
					name: "refuses a load without a key", give: []derive.Facts{putting(), withKeys(reading())},
					wantCallable: "Get", wantDetail: "takes no key",
				},
				{
					name:         "refuses a load of two keys",
					give:         []derive.Facts{putting(), withKeys(reading(), text("id"), text("name"))},
					wantCallable: "Get", wantDetail: "takes 2 keys",
				},
				{
					name:         "refuses a load whose key is a function",
					give:         []derive.Facts{putting(), withInput(reading(), symbol.FormFunc)},
					wantCallable: "Get", wantDetail: "takes a function as its key",
				},
				{
					name:         "refuses two loads that select two keys of one store",
					give:         []derive.Facts{putting(), reading(), byName},
					wantCallable: "GetByName", wantDetail: "selects another key of the value that Put stores than Get",
				},
				{
					name:         "refuses a load without an error",
					give:         []derive.Facts{putting(), withErrors(reading(), rules.ErrorsNone)},
					wantCallable: "Get", wantDetail: "returns no error",
				},
				{
					name:         "refuses a load of two values",
					give:         []derive.Facts{putting(), withRole(lookup(), "persister", "reader")},
					wantCallable: "Lookup", wantDetail: "returns 2 values",
				},
				{
					name:         "refuses a store that returns another value than the value that it stores",
					give:         []derive.Facts{withResults(putting(), text("")), reading()},
					wantCallable: "Put", wantDetail: "returns another value than the value that it stores",
				},
				{
					name:         "refuses a store that returns two values",
					give:         []derive.Facts{withResults(putting(), item(""), item("")), reading()},
					wantCallable: "Put", wantDetail: "returns another value than the value that it stores",
				},
				{
					name:         "refuses an append of two inputs",
					give:         []derive.Facts{withKeys(appending(), item("it"), text("note"))},
					wantCallable: "Append", wantDetail: "appends 2 inputs",
				},
				{
					name:         "refuses two appends of two types",
					give:         []derive.Facts{appending(), withKeys(renamed(appending(), "Note"), text("note"))},
					wantCallable: "Note", wantDetail: "appends another type than Append",
				},
				{
					name:         "refuses an append without a position",
					give:         []derive.Facts{withResults(appending())},
					wantCallable: "Append", wantDetail: "returns no integer position",
				},
				{
					name:         "refuses an append whose position is no integer",
					give:         []derive.Facts{withResults(appending(), position(rules.ScalarFloat))},
					wantCallable: "Append", wantDetail: "returns no integer position",
				},
				{
					name:         "refuses an append whose position is no number",
					give:         []derive.Facts{withResults(appending(), text(""))},
					wantCallable: "Append", wantDetail: "returns no integer position",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, refused := derive.ReferenceOf(tt.give)

					assert.NotNil(t, refused, "no reference form covers a callable")
					assert.Equal(t, []string{refused.Callable, refused.Detail},
						[]string{tt.wantCallable, tt.wantDetail}, "the refusal names the callable and what of it no "+
							"form covers")
					assert.Equal(t, got, derive.Reference{}, "a refused declaration has no reference")
				})
			}
		}

		t.Run("refuses with the reason and the remedy of the gate reference-role", func(t *testing.T) {
			t.Parallel()

			_, got := derive.ReferenceOf([]derive.Facts{putting(), withMixins(reading(), "ttl")})

			assert.NotNil(t, got, "no reference form covers Get")
			assert.Equal(t, got.Refusal, derive.Refusal{
				Callable: "Get", Key: "store.get", Gate: spec.ReferenceRole,
				Why: "Get has no role of a contract with a reference form",
				Remedy: "give Get a role of a contract with a reference form, or leave Store without the " +
					"directive fake",
			}, "the reason and the remedy name the callable and its declaration")
		})
	})
}
