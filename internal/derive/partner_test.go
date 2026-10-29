// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// idOf returns the identity of the method name of Store.
func idOf(name string) symbol.Identity {
	return symbol.Identity{Lang: golang, Package: "example.com/store", Owner: "Store", Name: name}
}

// fieldID returns the identity of the field name of Item.
func fieldID(name string) symbol.Identity {
	return symbol.Identity{
		Lang:    golang,
		Package: "example.com/store",
		Owner:   "Item",
		Name:    name,
		Kind:    symbol.KindField,
	}
}

// getter returns the facts of Store.Get(ctx, id string) (Item, error).
func getter() derive.Facts {
	inputs := []derive.Value{text("id")}
	return derive.Facts{
		Declaration: "Store",
		Callable:    "Get",
		ID:          idOf("Get"),
		Signature: signature(
			"Get",
			true,
			inputs,
			[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError},
			rules.ErrorsLastReturn,
		),
		Inputs:  inputs,
		Results: []derive.Value{item("")},
	}
}

// putter returns the facts of Store.Put(ctx, it Item) error with its
// identity.
func putter() derive.Facts {
	f := writer()
	f.ID = idOf("Put")
	return f
}

// keyedPutter returns the facts of Store.Set(ctx, id string, it Item) error.
func keyedPutter() derive.Facts {
	inputs := []derive.Value{text("id"), item("it")}
	return derive.Facts{
		Declaration: "Store", Callable: "Set", ID: idOf("Set"),
		Signature: signature("Set", true, inputs, []rules.ReturnRole{rules.ReturnError}, rules.ErrorsLastReturn),
		Inputs:    inputs,
	}
}

// withParam returns f with the param key of the classification name set to
// value.
func withParam(f derive.Facts, name, key string, value any) derive.Facts {
	f.Class.Params = map[string]shape.Params{name: {directiveKey(key): value}}
	return f
}

// withShape returns f with the shape s.
func withShape(f derive.Facts, s shape.Shape) derive.Facts {
	f.Class.Shape = s
	return f
}

// withRole returns f with the role of an instance of the contract.
func withRole(f derive.Facts, contract shape.Contract, role shape.Role) derive.Facts {
	f.Class.Contracts = append(f.Class.Contracts, shape.Membership{Contract: contract, Role: role})
	return f
}

// argumentsOf returns the arguments of a pair.
func argumentsOf(p *derive.Pair) []derive.Argument {
	if p == nil {
		return nil
	}
	return p.Arguments
}

func TestPartner(t *testing.T) {
	t.Parallel()

	t.Run("Partner", func(t *testing.T) {
		t.Parallel()

		t.Run("binds the named write and selects the field of the written value", func(t *testing.T) {
			t.Parallel()
			get := withParam(getter(), "read-after-write", "write", idOf("Put"))
			decl := []derive.Facts{get, putter()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Partner.Callable, "Put", "the partner is the named write")
			assert.Equal(t, argumentsOf(got), []derive.Argument{
				{Name: "ctx", Kind: derive.ArgContext},
				{Name: "id", Kind: derive.ArgField, Field: "ID"},
			}, "the read takes the context and the field of the written value")
		})

		t.Run("leaves a function parameter of the read to the probe", func(t *testing.T) {
			t.Parallel()
			body := derive.Value{Name: "body", Shape: rules.TypeShape{Form: symbol.FormFunc}}
			run := derive.Facts{
				Declaration: "Store", Callable: "Run", ID: idOf("Run"),
				Signature: signature("Run", true, []derive.Value{body}, []rules.ReturnRole{rules.ReturnError},
					rules.ErrorsLastReturn),
				Inputs: []derive.Value{body},
			}
			decl := []derive.Facts{withParam(run, "transaction", "write", idOf("Put")), putter()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, argumentsOf(got), []derive.Argument{
				{Name: "ctx", Kind: derive.ArgContext},
				{Name: "body", Kind: derive.ArgOwn},
			}, "the read takes the context and the function that the probe supplies")
		})

		t.Run("leaves each parameter of the read to its own sample for a write without inputs", func(t *testing.T) {
			t.Parallel()
			closer := derive.Facts{
				Declaration: "Store", Callable: "Close", ID: idOf("Close"),
				Signature: signature("Close", true, nil, []rules.ReturnRole{rules.ReturnError}, rules.ErrorsLastReturn),
			}
			decl := []derive.Facts{withParam(getter(), "lifecycle-after-close", "close", idOf("Close")), closer}

			got, gate := derive.Partner(
				&spec.Partner{Name: "close", Mode: spec.ModeNamed, Param: "close"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, argumentsOf(got), []derive.Argument{
				{Name: "ctx", Kind: derive.ArgContext},
				{Name: "id", Kind: derive.ArgOwn},
			}, "the read takes the context and its own sample of the id")
		})

		t.Run("shares one sample for a parameter of the same name in both calls", func(t *testing.T) {
			t.Parallel()
			get := withParam(getter(), "read-after-write", "write", idOf("Set"))
			decl := []derive.Facts{get, keyedPutter()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, argumentsOf(got)[1], derive.Argument{Name: "id", Kind: derive.ArgShared, Param: "id"},
				"the read shares the sample of the parameter id")
		})

		t.Run("writes the input that the read does not share as the written value", func(t *testing.T) {
			t.Parallel()
			readInputs := []derive.Value{text("key")}
			read := derive.Facts{
				Declaration: "Store", Callable: "Read", ID: idOf("Read"),
				Signature: signature("Read", true, readInputs,
					[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError}, rules.ErrorsLastReturn),
				Inputs: readInputs, Results: []derive.Value{text("")},
			}
			writeInputs := []derive.Value{text("key"), text("value")}
			write := derive.Facts{
				Declaration: "Store", Callable: "Write", ID: idOf("Write"),
				Signature: signature("Write", true, writeInputs, []rules.ReturnRole{rules.ReturnError},
					rules.ErrorsLastReturn),
				Inputs: writeInputs,
			}
			decl := []derive.Facts{withParam(read, "read-after-write", "write", idOf("Write")), write}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Written, 1, "the write writes its value, and the read shares its key")
		})

		t.Run("writes the input of the type of the read as the written value", func(t *testing.T) {
			t.Parallel()
			read := derive.Facts{
				Declaration: "Store", Callable: "Latest", ID: idOf("Latest"),
				Signature: signature("Latest", true, nil,
					[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError}, rules.ErrorsLastReturn),
				Results: []derive.Value{text("")},
			}
			writeInputs := []derive.Value{item("it"), text("value")}
			write := derive.Facts{
				Declaration: "Store", Callable: "Note", ID: idOf("Note"),
				Signature: signature("Note", true, writeInputs, []rules.ReturnRole{rules.ReturnError},
					rules.ErrorsLastReturn),
				Inputs: writeInputs,
			}
			decl := []derive.Facts{withParam(read, "read-after-write", "write", idOf("Note")), write}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Written, 1, "the write writes the text that the read returns, not the item before it")
		})

		t.Run("binds the named read of a write with the halves swapped", func(t *testing.T) {
			t.Parallel()
			put := withParam(putter(), "injection-safe", "read", idOf("Get"))
			decl := []derive.Facts{put, getter()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "read", Mode: spec.ModeReader, Param: "read"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Read.Callable, "Get", "the partner reads")
			assert.Equal(t, got.Write.Callable, "Put", "the licensed callable writes")
		})

		t.Run("derives the one other callable that takes the value whose field the param names", func(t *testing.T) {
			t.Parallel()
			get := withParam(getter(), "causal", "version", fieldID("Name"))
			decl := []derive.Facts{get, putter(), marker()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeDerived, Param: "version"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Partner.Callable, "Put", "the partner takes the Item that has the field Name")
			assert.Equal(t, got.Read.Callable, "Get", "the licensed callable reads")
		})

		t.Run("binds the role of the same contract", func(t *testing.T) {
			t.Parallel()
			put := withRole(putter(), "persister", "writer")
			get := withRole(getter(), "persister", "reader")
			decl := []derive.Facts{put, get}

			got, gate := derive.Partner(
				&spec.Partner{Name: "read", Mode: spec.ModeRole, Role: "reader"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Read.Callable, "Get", "the role that returns a value reads")
			assert.Equal(t, got.Write.Callable, "Put", "the role without a value writes")
		})

		{
			// A write that returns what it stored, as Put(ctx, it Item)
			// (Item, error), has a value of its own beside the read.
			stored := withRole(putter(), "persister", "writer")
			stored.Results = []derive.Value{item("")}
			stored.Signature.Returns = []rules.ReturnView{{Role: rules.ReturnValue}, {Role: rules.ReturnError}}
			tests := []struct {
				name      string
				giveFirst int
				giveRole  string
			}{
				{name: "binds the role that takes the value of the other as the write", giveRole: "reader"},
				{
					name:      "binds the role that returns the value of the other as the read",
					giveFirst: 1, giveRole: "writer",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					decl := []derive.Facts{stored, withRole(getter(), "persister", "reader")}

					got, gate := derive.Partner(
						&spec.Partner{Name: "partner", Mode: spec.ModeRole, Role: shape.Role(tt.giveRole)},
						&decl[tt.giveFirst],
						decl,
					)

					assert.Nil(t, gate, "the partner binds")
					assert.NotNil(t, got.Read, "one of the two roles reads")
					assert.Equal(t, got.Read.Callable, "Get", "the role that returns the type of the other reads")
					assert.Equal(t, got.Write.Callable, "Put", "the role that takes the type of the other writes")
				})
			}
		}

		{
			// Mark(ctx, id string) error returns nothing and takes no item, as
			// a write that marks the item under id.
			markInputs := []derive.Value{text("id")}
			mark := withRole(derive.Facts{
				Declaration: "Store",
				Callable:    "Mark",
				ID:          idOf("Mark"),
				Signature: signature(
					"Mark",
					true,
					markInputs,
					[]rules.ReturnRole{rules.ReturnError},
					rules.ErrorsLastReturn,
				),
				Inputs: markInputs,
			}, "persister", "writer")
			tests := []struct {
				name      string
				giveFirst int
				giveRole  string
			}{
				{name: "binds the role without a value as the write of the role with one", giveRole: "reader"},
				{
					name:      "binds the role with a value as the read of the role without one",
					giveFirst: 1,
					giveRole:  "writer",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					decl := []derive.Facts{mark, withRole(getter(), "persister", "reader")}

					got, gate := derive.Partner(
						&spec.Partner{Name: "partner", Mode: spec.ModeRole, Role: shape.Role(tt.giveRole)},
						&decl[tt.giveFirst],
						decl,
					)

					assert.Nil(t, gate, "the partner binds")
					assert.NotNil(t, got.Read, "one of the two roles reads")
					assert.Equal(t, got.Read.Callable, "Get", "the role with a value reads")
					assert.Equal(t, got.Write.Callable, "Mark", "the role without a value writes")
				})
			}
		}

		t.Run(
			"leaves the arguments of two roles that each take the value of the other to the probe",
			func(t *testing.T) {
				t.Parallel()
				coded := func(name, role string) derive.Facts {
					inputs := []derive.Value{text("in")}
					return withRole(derive.Facts{
						Declaration: "Store", Callable: name, ID: idOf(name),
						Signature: signature(name, true, inputs,
							[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError}, rules.ErrorsLastReturn),
						Inputs: inputs, Results: []derive.Value{text("")},
					}, "codec", shape.Role(role))
				}
				decl := []derive.Facts{coded("Encode", "forward"), coded("Decode", "inverse")}

				got, gate := derive.Partner(
					&spec.Partner{Name: "inverse", Mode: spec.ModeRole, Role: "inverse"},
					&decl[0],
					decl,
				)

				assert.Nil(t, gate, "the partner binds")
				assert.Nil(t, got.Read, "neither role is the write of the other")
				assert.Equal(t, got.Partner.Callable, "Decode", "the partner is the role inverse")
			},
		)

		t.Run("leaves the arguments of two roles that both read to the probe", func(t *testing.T) {
			t.Parallel()
			cache := withRole(getter(), "cache", "cache")
			backing := withRole(reader(), "cache", "backing")
			decl := []derive.Facts{cache, backing}

			got, gate := derive.Partner(
				&spec.Partner{
					Name: "backing",
					Mode: spec.ModeRole,
					Role: "backing",
				},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Nil(t, got.Read, "neither role writes")
			assert.Equal(t, got.Partner.Callable, "Read", "the partner is the role backing")
		})

		t.Run("binds the one establishing write of the classification", func(t *testing.T) {
			t.Parallel()
			decl := []derive.Facts{getter(), withShape(putter(), "writer")}

			got, gate := derive.Partner(
				&spec.Partner{
					Name:  "put",
					Mode:  spec.ModeEstablishing,
					Shape: "writer",
				},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Write.Callable, "Put", "the writer establishes the state")
		})

		t.Run("binds an interposed write whose value has a field that the read does not select by", func(t *testing.T) {
			t.Parallel()
			get := withParam(getter(), "point-in-time", "write", idOf("Put"))
			decl := []derive.Facts{get, putter()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeInterposed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Changed, "Name", "the second write changes the field Name")
		})

		t.Run("binds a rewriting read by its role", func(t *testing.T) {
			t.Parallel()
			put := withRole(putter(), "updater", "writer")
			get := withRole(getter(), "updater", "reader")
			decl := []derive.Facts{put, get}

			got, gate := derive.Partner(
				&spec.Partner{
					Name: "read",
					Mode: spec.ModeRewriter,
					Role: "reader",
				},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Read.Callable, "Get", "the role reader reads")
			assert.Equal(t, got.Changed, "Name", "the second write changes the field Name")
		})

		t.Run("binds the role that writes for a licensed callable that reads", func(t *testing.T) {
			t.Parallel()
			get := withRole(getter(), "persister", "reader")
			put := withRole(putter(), "persister", "writer")
			decl := []derive.Facts{get, put}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeRole, Role: "writer"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Write.Callable, "Put", "the role without a value writes")
		})

		t.Run("binds the one establishing write of a mixin", func(t *testing.T) {
			t.Parallel()
			put := putter()
			put.Class.Mixins = []shape.Mixin{"idempotent"}
			decl := []derive.Facts{getter(), put}

			got, gate := derive.Partner(
				&spec.Partner{
					Name:  "put",
					Mode:  spec.ModeEstablishing,
					Mixin: "idempotent",
				},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Write.Callable, "Put", "the write of the mixin establishes the state")
		})

		t.Run("binds the role of the same instance of the same contract", func(t *testing.T) {
			t.Parallel()
			put := putter()
			put.Class.Contracts = []shape.Membership{{Contract: "persister", Role: "writer", ID: "a"}}
			otherContract := getter()
			otherContract.Callable = "Find"
			otherContract.Class.Contracts = []shape.Membership{{Contract: "updater", Role: "reader", ID: "a"}}
			otherInstance := getter()
			otherInstance.Callable = "Load"
			otherInstance.Class.Contracts = []shape.Membership{{Contract: "persister", Role: "reader", ID: "b"}}
			otherRole := getter()
			otherRole.Callable = "Peek"
			otherRole.Class.Contracts = []shape.Membership{{Contract: "persister", Role: "writer", ID: "a"}}
			get := getter()
			get.Class.Contracts = []shape.Membership{{Contract: "persister", Role: "reader", ID: "a"}}
			decl := []derive.Facts{put, otherContract, otherInstance, otherRole, get}

			got, gate := derive.Partner(
				&spec.Partner{Name: "read", Mode: spec.ModeRole, Role: "reader"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Partner.Callable, "Get", "the reader of the instance a of the persister is the partner")
		})

		t.Run("leaves the arguments of two roles that both write to the probe", func(t *testing.T) {
			t.Parallel()
			put := withRole(putter(), "saga", "step")
			mark := withRole(marker(), "saga", "compensate")
			decl := []derive.Facts{put, mark}

			got, gate := derive.Partner(
				&spec.Partner{Name: "compensate", Mode: spec.ModeRole, Role: "compensate"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Nil(t, got.Read, "neither role reads")
		})

		t.Run("binds the establishing write whose arguments the read binds among two writers", func(t *testing.T) {
			t.Parallel()
			get := withShape(getter(), "writer")
			decl := []derive.Facts{get, withShape(putter(), "writer"), withShape(hasher(), "writer")}

			got, gate := derive.Partner(
				&spec.Partner{Name: "put", Mode: spec.ModeEstablishing, Shape: "writer"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, got.Write.Callable, "Put", "the write whose value the read selects establishes the state")
		})

		t.Run("binds each input of a read with two inputs", func(t *testing.T) {
			t.Parallel()
			get := getter()
			get.Inputs = []derive.Value{text("id"), text("name")}
			get.Signature = signature(
				"Get",
				true,
				get.Inputs,
				[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError},
				rules.ErrorsLastReturn,
			)
			get = withParam(get, "read-after-write", "write", idOf("Put"))
			decl := []derive.Facts{get, putter()}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(t, argumentsOf(got)[2], derive.Argument{Name: "name", Kind: derive.ArgField, Field: "Name"},
				"the second input takes the field Name")
		})

		t.Run("binds a named write whose value the read selects by every field", func(t *testing.T) {
			t.Parallel()
			get := withParam(getter(), "read-after-write", "write", idOf("Put"))
			decl := []derive.Facts{get, onlyID(putter())}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "a write that is not repeated needs no field to change")
			assert.Equal(t, got.Changed, "", "the pair changes no field")
		})

		t.Run("shares no sample between parameters of one name and two types", func(t *testing.T) {
			t.Parallel()
			set := keyedPutter()
			set.Inputs = append([]derive.Value(nil), set.Inputs...)
			set.Inputs[0].Shape = rules.TypeShape{Form: symbol.FormScalar, Spelling: "int"}
			get := withParam(getter(), "read-after-write", "write", idOf("Set"))
			decl := []derive.Facts{get, set}

			got, gate := derive.Partner(
				&spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				&decl[0],
				decl,
			)

			assert.Nil(t, gate, "the partner binds")
			assert.Equal(
				t,
				argumentsOf(got)[1].Kind,
				derive.ArgField,
				"the read takes the field of the value in place of the int",
			)
		})

		tests := []struct {
			name     string
			give     spec.Partner
			giveDecl []derive.Facts
			want     spec.GateID
		}{
			{
				name:     "refuses with partner-found for an establishing mixin that no callable has",
				give:     spec.Partner{Name: "put", Mode: spec.ModeEstablishing, Mixin: "idempotent"},
				giveDecl: []derive.Facts{getter(), withShape(putter(), "writer"), keyedPutter()},
				want:     spec.PartnerFound,
			},
			{
				name: "refuses with partner-found for a derived write of a callable without a value",
				give: spec.Partner{Name: "write", Mode: spec.ModeDerived, Param: "version"},
				giveDecl: []derive.Facts{
					withParam(putter(), "causal", "version", fieldID("ID")), keyedPutter(),
				},
				want: spec.PartnerFound,
			},
			{
				name: "refuses with partner-found for a field that the value of the read lacks",
				give: spec.Partner{Name: "write", Mode: spec.ModeDerived, Param: "version"},
				giveDecl: []derive.Facts{
					withParam(getter(), "causal", "version", fieldID("Gone")), keyedPutter(),
				},
				want: spec.PartnerFound,
			},
			{
				name:     "refuses with partner-found for a derived write without the param",
				give:     spec.Partner{Name: "write", Mode: spec.ModeDerived, Param: "version"},
				giveDecl: []derive.Facts{getter(), keyedPutter()},
				want:     spec.PartnerFound,
			},
			{
				name: "refuses with partner-arguments for a named read without a value",
				give: spec.Partner{Name: "read", Mode: spec.ModeReader, Param: "read"},
				giveDecl: []derive.Facts{
					withParam(putter(), "injection-safe", "read", idOf("Mark")), marker(),
				},
				want: spec.PartnerArguments,
			},
			{
				name: "refuses with partner-found for a param whose value is no callable",
				give: spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "read-after-write", "write", "Put"), putter(), reader(),
				},
				want: spec.PartnerFound,
			},
			{
				name: "refuses with partner-arguments for a field of another type",
				give: spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "read-after-write", "write", idOf("Put")), intID(putter()),
				},
				want: spec.PartnerArguments,
			},
			{
				name:     "refuses with partner-found for a named read that the classification does not name",
				give:     spec.Partner{Name: "read", Mode: spec.ModeReader, Param: "read"},
				giveDecl: []derive.Facts{putter(), getter()},
				want:     spec.PartnerFound,
			},
			{
				name: "refuses with partner-found for a classification that names another param",
				give: spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "read-after-write", "read", idOf("Put")), putter(),
				},
				want: spec.PartnerFound,
			},
			{
				name: "refuses with partner-arguments for an establishing write that the read cannot select",
				give: spec.Partner{
					Name:  "put",
					Mode:  spec.ModeEstablishing,
					Shape: "writer",
				},
				giveDecl: []derive.Facts{
					reader(), withShape(putter(), "writer"),
				},
				want: spec.PartnerArguments,
			},
			{
				name: "refuses with partner-arguments for a write that takes no value of the read",
				give: spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "read-after-write", "write", idOf("Hash")), hasher(),
				},
				want: spec.PartnerArguments,
			},
			{
				name: "refuses with rewritable for an interposed write that takes no value of the read",
				give: spec.Partner{Name: "write", Mode: spec.ModeInterposed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "point-in-time", "write", idOf("Mark")), marker(),
				},
				want: spec.Rewritable,
			},
			{
				name: "refuses with partner-found for a param that names no callable of the declaration",
				give: spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "read-after-write", "write", idOf("Gone")), putter(),
				},
				want: spec.PartnerFound,
			},
			{
				name:     "refuses with partner-found for a classification without the param",
				give:     spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{getter(), putter()},
				want:     spec.PartnerFound,
			},
			{
				name: "refuses with partner-arguments for a parameter that selects nothing of the write",
				give: spec.Partner{Name: "write", Mode: spec.ModeNamed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(reader(), "read-after-write", "write", idOf("Put")), putter(),
				},
				want: spec.PartnerArguments,
			},
			{
				name: "refuses with one-partner for two callables that take the value of the field",
				give: spec.Partner{Name: "write", Mode: spec.ModeDerived, Param: "version"},
				giveDecl: []derive.Facts{
					withParam(getter(), "causal", "version", fieldID("ID")), keyedPutter(), putter(),
				},
				want: spec.OnePartner,
			},
			{
				name: "refuses with partner-found for a value of the field that no other callable takes",
				give: spec.Partner{Name: "write", Mode: spec.ModeDerived, Param: "version"},
				giveDecl: []derive.Facts{
					withParam(getter(), "causal", "version", fieldID("ID")), marker(),
				},
				want: spec.PartnerFound,
			},
			{
				name: "refuses with partner-found for a role that no callable plays",
				give: spec.Partner{
					Name: "read",
					Mode: spec.ModeRole,
					Role: "reader",
				},
				giveDecl: []derive.Facts{withRole(putter(), "persister", "writer"), getter()},
				want:     spec.PartnerFound,
			},
			{
				name: "refuses with one-partner for two establishing writes",
				give: spec.Partner{
					Name:  "put",
					Mode:  spec.ModeEstablishing,
					Shape: "writer",
				},
				giveDecl: []derive.Facts{
					getter(), withShape(putter(), "writer"), withShape(keyedPutter(), "writer"),
				},
				want: spec.OnePartner,
			},
			{
				name: "refuses with partner-found for an establishing write that the declaration lacks",
				give: spec.Partner{
					Name:  "put",
					Mode:  spec.ModeEstablishing,
					Shape: "writer",
				},
				giveDecl: []derive.Facts{getter(), putter()},
				want:     spec.PartnerFound,
			},
			{
				name: "refuses with rewritable for a written value whose every field the read selects by",
				give: spec.Partner{Name: "write", Mode: spec.ModeInterposed, Param: "write"},
				giveDecl: []derive.Facts{
					withParam(getter(), "point-in-time", "write", idOf("Put")), onlyID(putter()),
				},
				want: spec.Rewritable,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				_, gate := derive.Partner(&tt.give, &tt.giveDecl[0], tt.giveDecl)

				assert.NotNil(t, gate, "the claim refuses")
				assert.Equal(t, gate.ID, tt.want, "the claim refuses with the gate")
			})
		}
	})
}

// hasher returns the facts of Store.Hash(key string) string with its
// identity.
func hasher() derive.Facts {
	f := pure()
	f.ID = idOf("Hash")
	return f
}

// marker returns the facts of Store.Mark(ctx, id string) error, which takes
// no value of the read.
func marker() derive.Facts {
	inputs := []derive.Value{text("id")}
	return derive.Facts{
		Declaration: "Store", Callable: "Mark", ID: idOf("Mark"),
		Signature: signature("Mark", true, inputs, []rules.ReturnRole{rules.ReturnError}, rules.ErrorsLastReturn),
		Inputs:    inputs,
	}
}

// intID returns f with a written value whose field ID is an int.
func intID(f derive.Facts) derive.Facts {
	f.Inputs = append([]derive.Value(nil), f.Inputs...)
	f.Inputs[0].Fields = []derive.Field{{Name: "ID", Shape: rules.TypeShape{Form: symbol.FormScalar, Spelling: "int"}}}
	return f
}

// onlyID returns f with a written value whose one field is ID.
func onlyID(f derive.Facts) derive.Facts {
	f.Inputs = append([]derive.Value(nil), f.Inputs...)
	f.Inputs[0].Fields = f.Inputs[0].Fields[:1]
	return f
}

// directiveKey returns the key of a param of a classification.
func directiveKey(key string) directive.ParamKey { return directive.ParamKey(key) }
