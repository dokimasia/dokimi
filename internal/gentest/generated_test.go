// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/gentest"
)

// The openings of the contracts of the records that a failed helper writes.
const (
	templateContract = "the family of the records claims its files through its template"
	tagContract      = "the records are in the family of their tag"
	countContract    = "the generator reports one finding at "
	recordContract   = "the generator records nothing for "
)

// Records returns the records of one type of data from a phase call of a
// generator by their key, List returns their data in order, and Beside
// returns what a record declares beside its struct.
// Referents returns the references of the fields of a struct, and Finding
// returns the one finding at a declaration. Unrecorded checks that a subject
// has no record, and Assigned returns the names that statements assign.
func TestGenerated(t *testing.T) {
	t.Parallel()

	t.Run("Records", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the records of the type of the data by their key", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}})
			got := gentest.Records(t, res, recordFamily, "", subject.key)
			assert.Equal(t, keys(got), []string{addressName}, "the recorder records Address in its primary family")
			expect.Equal(t, got[addressName].Data, subject{Name: addressName}, "the record has its data")
			assert.NotNil(t, got[addressName].Decl, "the declaration of the record is a struct")
			expect.Equal(t, got[addressName].Decl.Name, "AddressRecord", "the record has its declaration")
		})

		t.Run("returns the first declaration of a record that declares more", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}, beside: []string{addressName}})
			got := gentest.Records(t, res, recordFamily, "", subject.key)
			assert.Equal(t, keys(got), []string{addressName}, "the later declarations of Address are no records")
			assert.NotNil(t, got[addressName].Decl, "the declaration of the record is a struct")
			expect.Equal(t, got[addressName].Decl.Name, "AddressRecord", "the record has its first struct")
		})

		t.Run("returns each record of one origin with the data of its own call", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{
				records: []string{addressName}, beside: []string{addressName}, again: []string{addressName},
			})
			got := gentest.Records(t, res, recordFamily, "", subject.key)
			assert.Equal(t, keys(got), []string{addressName, addressName + "Again"},
				"the recorder records Address twice in its primary family")
			expect.Equal(t, got[addressName+"Again"].Decl.Name, "AddressAgain", "the second record has its own struct")
		})

		t.Run("returns the records of the family of the tag", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}})
			got := gentest.Records(t, res, checkFamily, sdk.TagTest, check.key)
			assert.Equal(t, keys(got), []string{addressName}, "the recorder records Address in the family of its tests")
			expect.Equal(t, got[addressName].Decl.Name, "AddressCheck", "the record has its declaration")
		})

		t.Run("fails at a record whose unit has another template", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}})
			records := assert.Rejects(t, "a record of another template", func(tb assert.TB) {
				gentest.Records(tb, res, checkFamily, "", subject.key)
			})
			assert.Length(t, records, 1, "the helper fails once, at the record of Address")
			expect.Equal(t, records[0].Contract, templateContract, "the helper fails at the template")
		})

		t.Run("fails at a record in the family of another tag", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}})
			records := assert.Rejects(t, "a record of another tag", func(tb assert.TB) {
				gentest.Records(tb, res, recordFamily, sdk.TagTest, subject.key)
			})
			assert.Length(t, records, 1, "the helper fails once, at the record of Address")
			expect.Equal(t, records[0].Contract, tagContract, "the helper fails at the tag")
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give recording
			want []subject
		}{
			{
				name: "returns the data of a record that declares more once",
				give: recording{records: []string{addressName}, beside: []string{addressName}},
				want: []subject{{Name: addressName}},
			},
			{
				name: "returns the data of each record of one origin",
				give: recording{records: []string{addressName}, again: []string{addressName}},
				want: []subject{{Name: addressName}, {Name: addressName + "Again"}},
			},
			{
				name: "returns the data of each record of two origins with equal data",
				give: recording{
					records: []string{addressName, warehouseName}, alike: []string{addressName, warehouseName},
				},
				want: []subject{{Name: alikeName}, {Name: alikeName}},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				_, _, res := generate(t, tt.give)
				assert.Equal(t, gentest.List[subject](res), tt.want, "the list follows the order of the declarations")
			})
		}
	})

	t.Run("Beside", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the declarations after the first declaration of a record", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}, beside: []string{addressName}})
			got := gentest.Beside(res, subject.key)
			assert.Equal(t, got, map[string][]string{addressName: {"AddressOption", "NewAddressRecord"}},
				"the record of Address declares AddressOption and NewAddressRecord after AddressRecord")
		})

		t.Run("returns nothing for a record of a struct alone", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}})
			expect.Empty(t, gentest.Beside(res, subject.key), "the record of Address declares its struct alone")
		})

		t.Run("returns nothing of a record whose data has another type", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}, beside: []string{addressName}})
			expect.Empty(t, gentest.Beside(res, check.key), "the declarations after AddressRecord belong to a subject")
		})
	})

	t.Run("Referents", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the spelling of the type of each field in order", func(t *testing.T) {
			t.Parallel()

			s := &emit.Struct{Name: "AddressCheck"}
			s.Fields.Append(&emit.Field{Name: "New", Type: &emit.TypeRef{Spelling: "NewAddressRecord"}},
				&emit.Field{Name: "From", Type: &emit.TypeRef{Spelling: "NewAddressRecordFrom"}})
			assert.Equal(t, gentest.Referents(s), []string{"NewAddressRecord", "NewAddressRecordFrom"},
				"the spellings follow the order of the fields")
		})
	})

	t.Run("Finding", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the finding at the declaration", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{findings: map[string]int{warehouseName: 1}})
			got := gentest.Finding(t, res, f.Find(t, symbol.KindStruct, "", warehouseName))
			expect.Equal(t, got.Code, recorded, "the finding has its code")
			expect.Equal(t, got.Msg, "Warehouse is reported", "the finding is the finding at Warehouse")
		})

		t.Run("returns the finding at a member", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{fields: []string{warehouseName}})
			got := gentest.Finding(t, res, f.Find(t, symbol.KindField, warehouseName, streetName))
			expect.Equal(t, got.Msg, "the field of Warehouse is reported", "the finding is the finding at Street")
		})

		t.Run("fails where the generator reports no finding at the declaration", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{findings: map[string]int{warehouseName: 1}})
			records := assert.Rejects(t, "a declaration without a finding", func(tb assert.TB) {
				gentest.Finding(tb, res, f.Find(tb, symbol.KindStruct, "", addressName))
			})
			assert.Length(t, records, 1, "the helper stops at the count of the findings")
			expect.Equal(t, records[0].Contract, countContract+addressName, "the helper stops at the declaration")
		})

		t.Run("fails where the generator reports two findings at the declaration", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{findings: map[string]int{warehouseName: 2}})
			records := assert.Rejects(t, "a declaration with two findings", func(tb assert.TB) {
				gentest.Finding(tb, res, f.Find(tb, symbol.KindStruct, "", warehouseName))
			})
			assert.Length(t, records, 1, "the helper stops at the count of the findings")
			expect.Equal(t, records[0].Contract, countContract+warehouseName, "the helper stops at the declaration")
		})
	})

	t.Run("Unrecorded", func(t *testing.T) {
		t.Parallel()

		t.Run("passes where the generator records another declaration", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{records: []string{addressName}})
			gentest.Unrecorded(t, res, f.Find(t, symbol.KindStruct, "", warehouseName))
		})

		t.Run("fails at each record of the declaration", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{records: []string{addressName}})
			records := assert.Rejects(t, "a declaration with a record in each family", func(tb assert.TB) {
				gentest.Unrecorded(tb, res, f.Find(tb, symbol.KindStruct, "", addressName))
			})
			assert.Length(t, records, 2, "the helper fails at the record of each family")
			expect.Equal(t, records[0].Contract, recordContract+addressName, "the helper fails at the declaration")
		})
	})

	t.Run("Assigned", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give []emit.Stmt
			want string
		}{
			{name: "returns nothing for no statements"},
			{
				name: "writes each statement as its names joined by commas",
				give: []emit.Stmt{{Names: []string{"seed.Street"}}, {Names: []string{"a", "b"}}},
				want: "seed.Street a,b",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, gentest.Assigned(tt.give), tt.want, "the names follow the order of the statements")
			})
		}
	})
}

// key returns the key of the record of a subject: the name of the struct.
func (s subject) key() string { return s.Name }

// key returns the key of the record of a check: the name of the struct.
func (c check) key() string { return c.Name }

// keys returns the keys of the records, sorted.
func keys[D any](records map[string]gentest.Record[D]) []string {
	out := make([]string, 0, len(records))
	for k := range records {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
