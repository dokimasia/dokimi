// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/plugin/shape"
	sdkrules "go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/rules"
)

// The shapes of the values of a transaction: the handle Tx, a key, an int
// value and the list of the values of a key. A float value, an unsigned value
// and the address of a value fit no history, or another one.
var (
	txShape    = sdkrules.TypeShape{Form: symbol.FormReference, Spelling: "Tx"}
	intShape   = sdkrules.Scalar("int", sdkrules.ScalarInt, 0)
	listShape  = sdkrules.TypeShape{Form: symbol.FormList, Spelling: "[]int", Elems: []sdkrules.TypeShape{intShape}}
	floatShape = sdkrules.Scalar("float64", sdkrules.ScalarFloat, 0)
	uintShape  = sdkrules.Scalar("uint", sdkrules.ScalarUint, 0)
	addrShape  = sdkrules.TypeShape{Form: symbol.FormOptional, Spelling: "*int", Elems: []sdkrules.TypeShape{intShape}}
)

// listOf returns the shape of a list of the values of the shape s.
func listOf(s sdkrules.TypeShape) sdkrules.TypeShape {
	return sdkrules.TypeShape{Form: symbol.FormList, Spelling: "[]" + s.Spelling, Elems: []sdkrules.TypeShape{s}}
}

// valueOf returns a value named name of the shape s.
func valueOf(name string, s sdkrules.TypeShape) derive.Value {
	return derive.Value{Name: name, Shape: s}
}

// part returns the facts of the callable name of the declaration decl, with
// a context, the inputs and the results, an error, and the role of the
// contract tx. A part without a role has the empty role.
func part(decl, name string, role shape.Role, inputs, results []derive.Value) derive.Facts {
	returns := make([]sdkrules.ReturnRole, 0, len(results)+1)
	for range results {
		returns = append(returns, sdkrules.ReturnValue)
	}
	f := derive.Facts{
		Declaration: decl, Callable: name,
		ID:        symbol.Identity{Lang: golang, Package: "example.com/store", Owner: decl, Name: name},
		Signature: signature(name, true, inputs, append(returns, sdkrules.ReturnError), sdkrules.ErrorsLastReturn),
		Inputs:    inputs, Results: results,
	}
	if role != "" {
		f = withRole(f, "tx", role)
	}
	return f
}

// writes returns the facts of Store.Append(ctx, tx Tx, key string, v int)
// error with the mixin, whose param read names Get.
func writes(mixin shape.Mixin) derive.Facts {
	f := part("Store", "Append", "", []derive.Value{valueOf("tx", txShape), text("key"), valueOf("v", intShape)}, nil)
	return withParam(withMixins(f, mixin), string(mixin), "read", idOf("Get"))
}

// scoped returns the callables of a store whose transactions thread the
// handle Tx through each call: Begin, Commit, Rollback, Append with the mixin
// serializable and Get, whose read returns the list of the values of a key.
func scoped() []derive.Facts {
	tx := valueOf("tx", txShape)
	return []derive.Facts{
		part("Store", "Begin", "begin", nil, []derive.Value{valueOf("", txShape)}),
		part("Store", "Commit", "commit", []derive.Value{tx}, nil),
		part("Store", "Rollback", "rollback", []derive.Value{tx}, nil),
		writes(shape.Serializable),
		part("Store", "Get", "", []derive.Value{tx, text("key")}, []derive.Value{valueOf("", listShape)}),
	}
}

// handled returns the callables of a store whose Begin returns the handle Tx,
// and the methods of Tx that commit, roll back, write and read: Append with
// the mixin serializable and Get, whose param read names Get of Tx.
func handled() []derive.Facts {
	begin := part("Store", "Begin", "begin", nil, []derive.Value{valueOf("", txShape)})
	write := part("Tx", "Append", "", []derive.Value{text("key"), valueOf("v", intShape)}, nil)
	read := part("Tx", "Get", "", []derive.Value{text("key")}, []derive.Value{valueOf("", listShape)})
	begin.Handle = []derive.Facts{
		part("Tx", "Commit", "commit", nil, nil),
		part("Tx", "Rollback", "rollback", nil, nil),
		withParam(withMixins(write, shape.Serializable), "serializable", "read", read.ID),
		read,
	}
	return []derive.Facts{begin}
}

// isolation derives the claim of the isolation level id over decl with the
// machineries of concurrent transactions.
func isolation(t *testing.T, id spec.ClaimID, decl []derive.Facts) derive.Suite {
	t.Helper()
	c, ok := rules.ClaimOf(id)
	assert.True(t, ok, "the spec has the claim "+string(id))
	supplied := append(slices.Clone(runtime), spec.ConcurrentCallers, spec.GeneratedSequences)
	got, err := derive.Derive(decl, []rules.Claim{c}, nil, supplied, nil)
	assert.NoError(t, err, "the declaration derives")
	return got
}

// transactionsRefusal is the refusal of the claim serializable of Store by
// the gate transactions, on the callable anchor with the write Append.
func transactionsRefusal(anchor string) derive.Refusal {
	return derive.Refusal{
		Claim: spec.Serializable, Callable: anchor, Key: "store/serializable", Gate: spec.Transactions,
		Why: "Append and the read of its mixin do not run in the scope of a transaction of Store",
		Remedy: "give Append and its read the handle that begin returns, or declare them on an interface of the " +
			"package that the handle implements, with a key and a value of an integer or a string type",
	}
}

func TestTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Derive", func(t *testing.T) {
		t.Parallel()

		t.Run("binds an isolation level once on the write of a store whose callables take the handle",
			func(t *testing.T) {
				t.Parallel()
				decl := scoped()

				got := isolation(t, spec.Serializable, decl)

				assert.Equal(t, identities(got), []string{"Store/serializable"}, "the declaration binds the claim once")
				check := got.Checks[0]
				assert.Equal(t, check.Facts.Callable, "Append", "the claim binds on the write")
				assert.Equal(t, *check.Transaction, derive.Transaction{
					Begin:    derive.Part{Facts: &decl[0], Handle: -1},
					Commit:   derive.Part{Facts: &decl[1], Handle: 0},
					Rollback: derive.Part{Facts: &decl[2], Handle: 0},
					Write:    derive.Part{Facts: &decl[3], Handle: 0},
					Read:     derive.Part{Facts: &decl[4], Handle: 0},
					Key:      1, Value: 2, ReadKey: 1, Appends: true,
				}, "the transaction threads the handle through each call after the begin")
			})

		t.Run("binds an isolation level once on the begin of a handle whose methods write, read, commit and roll "+
			"back", func(t *testing.T) {
			t.Parallel()
			decl := handled()

			got := isolation(t, spec.Serializable, decl)

			assert.Equal(t, identities(got), []string{"Store/serializable"}, "the declaration binds the claim once")
			check := got.Checks[0]
			assert.Equal(t, check.Facts.Callable, "Begin", "the claim binds on the begin of the handle")
			handle := decl[0].Handle
			assert.Equal(t, *check.Transaction, derive.Transaction{
				Begin:    derive.Part{Facts: &decl[0], Handle: -1},
				Commit:   derive.Part{Facts: &handle[0], Handle: -1},
				Rollback: derive.Part{Facts: &handle[1], Handle: -1},
				Write:    derive.Part{Facts: &handle[2], Handle: -1},
				Read:     derive.Part{Facts: &handle[3], Handle: -1},
				Key:      0, Value: 1, ReadKey: 0, Appends: true,
			}, "each call after the begin is a method of the handle")
		})

		t.Run("reads a history of registers where the read returns the written value", func(t *testing.T) {
			t.Parallel()
			decl := scoped()
			decl[4].Results = []derive.Value{valueOf("", intShape)}

			got := isolation(t, spec.Serializable, decl)

			assert.False(t, got.Checks[0].Transaction.Appends, "the write sets the value of a register")
		})

		t.Run("binds a write of a value of an unsigned type", func(t *testing.T) {
			t.Parallel()
			decl := scoped()
			decl[3].Inputs = []derive.Value{decl[3].Inputs[0], text("key"), valueOf("v", uintShape)}
			decl[4].Results = []derive.Value{valueOf("", listOf(uintShape))}

			got := isolation(t, spec.Serializable, decl)

			assert.Equal(t, identities(got), []string{"Store/serializable"}, "a counter gives each write an integer")
		})

		t.Run("binds the begin of the contract after a callable that returns the handle", func(t *testing.T) {
			t.Parallel()
			decl := append([]derive.Facts{part("Store", "Peek", "", nil, []derive.Value{valueOf("", txShape)})},
				scoped()...)

			got := isolation(t, spec.Serializable, decl)

			assert.Length(t, got.Checks, 1, "the declaration binds the claim")
			assert.Equal(t, got.Checks[0].Transaction.Begin.Facts.Callable, "Begin",
				"the transaction begins with the role begin")
		})

		t.Run("binds snapshot-isolation on the write of its mixin", func(t *testing.T) {
			t.Parallel()
			decl := scoped()
			decl[3] = writes(shape.SnapshotIsolation)

			got := isolation(t, spec.SnapshotIsolation, decl)

			assert.Equal(t, identities(got), []string{"Store/snapshot-isolation"}, "the declaration binds the claim")
		})

		t.Run("licenses no isolation level of a write without the param read", func(t *testing.T) {
			t.Parallel()
			decl := scoped()
			decl[3].Class.Params = nil

			got := isolation(t, spec.Serializable, decl)

			assert.Empty(t, got.Checks, "the licence needs the param read")
			assert.Empty(t, got.Refusals, "a claim without its licence does not refuse")
		})

		{
			tests := []struct {
				name string
				give func(decl []derive.Facts) []derive.Facts
			}{
				{
					name: "refuses a write that takes no handle",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[3].Inputs = decl[3].Inputs[1:]
						return decl
					},
				},
				{
					name: "refuses a write of a value of a struct type",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[3].Inputs = []derive.Value{decl[3].Inputs[0], text("key"), item("v")}
						return decl
					},
				},
				{
					name: "refuses a write of a key of a struct type",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[3].Inputs = []derive.Value{decl[3].Inputs[0], item("key"), valueOf("v", intShape)}
						return decl
					},
				},
				{
					name: "refuses a write with an input beside the key and the value",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[3].Inputs = append(slices.Clone(decl[3].Inputs), text("note"))
						return decl
					},
				},
				{
					name: "refuses a read that takes no handle",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Inputs = decl[4].Inputs[1:]
						return decl
					},
				},
				{
					name: "refuses a read that takes no key",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Inputs = decl[4].Inputs[:1]
						return decl
					},
				},
				{
					name: "refuses a read that returns a value of another type",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Results = []derive.Value{text("")}
						return decl
					},
				},
				{
					name: "refuses a read that returns no value",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Results = nil
						return decl
					},
				},
				{
					name: "refuses a write whose param read names no callable",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[3] = withParam(decl[3], "serializable", "read", idOf("Missing"))
						return decl
					},
				},
				{
					name: "refuses a store without a begin",
					give: func(decl []derive.Facts) []derive.Facts { return decl[1:] },
				},
				{
					name: "refuses a store without a rollback",
					give: func(decl []derive.Facts) []derive.Facts { return slices.Delete(decl, 2, 3) },
				},
				{
					name: "refuses a commit of another instance of the contract",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[1].Class.Contracts = []shape.Membership{{Contract: "tx", Role: "commit", ID: "other"}}
						return decl
					},
				},
				{
					name: "refuses a rollback that takes no handle",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[2].Inputs = nil
						return decl
					},
				},
				{
					name: "refuses a begin that returns no handle",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[0].Results = nil
						return decl
					},
				},
				{
					name: "refuses a write of a value of a float type",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[3].Inputs = []derive.Value{decl[3].Inputs[0], text("key"), valueOf("v", floatShape)}
						decl[4].Results = []derive.Value{valueOf("", listOf(floatShape))}
						return decl
					},
				},
				{
					name: "refuses a write of a key of a float type",
					give: func(decl []derive.Facts) []derive.Facts {
						key := valueOf("key", floatShape)
						decl[3].Inputs = []derive.Value{decl[3].Inputs[0], key, valueOf("v", intShape)}
						decl[4].Inputs = []derive.Value{decl[4].Inputs[0], key}
						return decl
					},
				},
				{
					name: "refuses a read whose key has another type than the key of the write",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Inputs = []derive.Value{decl[4].Inputs[0], valueOf("key", intShape)}
						return decl
					},
				},
				{
					name: "refuses a read that returns a list of another type",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Results = []derive.Value{valueOf("", listOf(text("").Shape))}
						return decl
					},
				},
				{
					name: "refuses a read that returns the address of the value",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[4].Results = []derive.Value{valueOf("", addrShape)}
						return decl
					},
				},
				{
					name: "refuses a begin of another contract",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[0].Class.Contracts = []shape.Membership{{Contract: "lease", Role: "begin"}}
						return decl
					},
				},
				{
					name: "refuses a commit of another contract",
					give: func(decl []derive.Facts) []derive.Facts {
						decl[1].Class.Contracts = []shape.Membership{{Contract: "lease", Role: "commit"}}
						return decl
					},
				},
				{
					name: "refuses a commit that returns the handle in a store without a begin",
					give: func(decl []derive.Facts) []derive.Facts {
						decl = decl[1:]
						decl[0].Results = []derive.Value{valueOf("", txShape)}
						return decl
					},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					decl := tt.give(scoped())

					got := isolation(t, spec.Serializable, decl)

					assert.Empty(t, got.Checks, "the claim does not bind")
					assert.Equal(t, got.Refusals, []derive.Refusal{transactionsRefusal("Append")},
						"the write refuses the claim with the gate transactions")
				})
			}
		}

		t.Run("names the write of a handle in the refusal of the begin", func(t *testing.T) {
			t.Parallel()
			decl := handled()
			decl[0].Handle[2].Inputs = []derive.Value{text("key"), item("v")}

			got := isolation(t, spec.Serializable, decl)

			assert.Equal(t, got.Refusals, []derive.Refusal{transactionsRefusal("Begin")},
				"the begin refuses the claim, and the reason names the write")
		})

		t.Run("leaves the methods of a handle out of a claim of the declaration without the gate transactions",
			func(t *testing.T) {
				t.Parallel()
				c := claimOn("whole", rules.Licence{Mixin: "concurrent"})
				c.Scope = rules.ScopeDeclaration
				decl := handled()
				decl[0].Handle[3] = withMixins(decl[0].Handle[3], "concurrent")

				got, err := derive.Derive(decl, []rules.Claim{c}, nil, runtime, nil)

				assert.NoError(t, err, "the declaration derives")
				assert.Empty(t, got.Checks, "a mixin of a method of the handle licenses no other claim")
			})

		t.Run("names the write in the refusal of a gate that the suite generator judges", func(t *testing.T) {
			t.Parallel()
			decl := handled()
			check := isolation(t, spec.Serializable, decl).Checks[0]

			got := derive.Refuse(&check, spec.ImportableSamples)

			assert.Equal(t, got.Callable, "Begin", "the refusal is on the begin")
			assert.Equal(t, got.Why, "a sample of the check refers to a declaration that the file of the suite cannot "+
				"import", "the gate has its reason")
			assert.Equal(t, got.Remedy, "author a sample of each input of Append with the directive sample",
				"the remedy names the write")
		})
	})

	t.Run("Begins", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give shape.Membership
			want bool
		}{
			{
				name: "reports the role begin of the contract tx",
				give: shape.Membership{Contract: "tx", Role: "begin"}, want: true,
			},
			{name: "reports no other role of the contract tx", give: shape.Membership{Contract: "tx", Role: "commit"}},
			{
				name: "reports no role begin of another contract",
				give: shape.Membership{Contract: "lease", Role: "begin"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				f := part("Store", "Begin", "", nil, nil)
				f.Class.Contracts = []shape.Membership{tt.give}

				assert.Equal(t, derive.Begins(&f), tt.want, "the role begin of the contract tx opens a transaction")
			})
		}
	})
}
