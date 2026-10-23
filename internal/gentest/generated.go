// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"reflect"
	"strings"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// Record is one record of a family of a generator: the data that the
// generator attached to it and its first declaration. Decl is nil where the
// first declaration is no struct.
type Record[D any] struct {
	Data D
	Decl *emit.Struct
}

// Records returns each record of res whose data has the type D, by the key
// that key returns for the data. The declarations of a record are the
// declarations that one call of the generator places with the data, and the
// first of them is the declaration of the record. The declarations after it,
// such as a function that a record declares beside its struct, are no records
// of their own. Records fails the test for a record in a unit whose template
// is not template or whose tag is not tag, so a spec checks that each family
// claims its files through its own template. The primary family has the empty
// tag.
func Records[D any](
	tb assert.TB,
	res plugintest.Result,
	template, tag string,
	key func(D) string,
) map[string]Record[D] {
	tb.Helper()

	out := map[string]Record[D]{}
	for u := range res.Emit.Units() {
		opens := opened(&u)
		for i, d := range u.Decls {
			data, typed := u.Data[d].(D)
			if !typed || !opens[i] {
				continue
			}
			expect.Equal(tb, u.Template, template, "the family of the records claims its files through its template")
			expect.Equal(tb, u.Tag, tag, "the records are in the family of their tag")
			decl, _ := d.(*emit.Struct)
			out[key(data)] = Record[D]{Data: data, Decl: decl}
		}
	}
	return out
}

// List returns the data of each record of res whose data has the type D, in
// the order of the units and of their declarations. A record that a generator
// places twice is in the list twice.
func List[D any](res plugintest.Result) []D {
	var out []D
	for u := range res.Emit.Units() {
		opens := opened(&u)
		for i, d := range u.Decls {
			if data, typed := u.Data[d].(D); typed && opens[i] {
				out = append(out, data)
			}
		}
	}
	return out
}

// Beside returns the names of the declarations that each record of res
// whose data has the type D declares after its first declaration, by the key
// that key returns for the data, in the order of the record.
func Beside[D any](res plugintest.Result, key func(D) string) map[string][]string {
	out := map[string][]string{}
	for u := range res.Emit.Units() {
		opens := opened(&u)
		for i, d := range u.Decls {
			data, typed := u.Data[d].(D)
			if !typed || opens[i] {
				continue
			}
			out[key(data)] = append(out[key(data)], emit.DeclaredName(d))
		}
	}
	return out
}

// opened reports for each declaration of u whether it opens a record. The
// declarations that one call places share their origin and their data, so a
// declaration opens a record where the declaration before it has another
// origin or other data.
func opened(u *plugin.Unit) []bool {
	out := make([]bool, len(u.Decls))
	var origin symbol.Identity
	var data any
	for i, d := range u.Decls {
		own := u.Data[d]
		at, _ := emit.OriginOf(d)
		out[i] = at != origin || !reflect.DeepEqual(own, data)
		origin, data = at, own
	}
	return out
}

// Referents returns the spelling of the type of each field of the struct s,
// in order. A record of the checks refers to a declaration of another record
// through such a field.
func Referents(s *emit.Struct) []string {
	out := make([]string, 0, s.Fields.Len())
	for _, f := range s.Fields.Items() {
		out = append(out, f.Type.Spelling)
	}
	return out
}

// Finding returns the finding that res reports at the position of the
// declaration at, and fails the test where res reports no finding or two
// findings there.
func Finding(tb assert.TB, res plugintest.Result, at node.Declaration) diag.Diag {
	tb.Helper()

	var found []diag.Diag
	for d := range res.Sink.All() {
		if d.Pos == at.Position() {
			found = append(found, d)
		}
	}
	assert.Length(tb, found, 1, "the generator reports one finding at "+at.Identity().Name)
	return found[0]
}

// Unrecorded fails the test for each declaration of res whose origin is the
// declaration subject, so a spec checks that a generator records nothing for
// a subject that it reports.
func Unrecorded(tb assert.TB, res plugintest.Result, subject node.Declaration) {
	tb.Helper()

	id := subject.Identity()
	for u := range res.Emit.Units() {
		for _, d := range u.Decls {
			origin, _ := emit.OriginOf(d)
			expect.NotEqual(tb, origin, id, "the generator records nothing for "+id.Name)
		}
	}
}

// Assigned returns the names that each statement of stmts assigns. The names
// of one statement are joined by commas, and the statements by spaces.
func Assigned(stmts []emit.Stmt) string {
	out := make([]string, 0, len(stmts))
	for i := range stmts {
		out = append(out, strings.Join(stmts[i].Names, ","))
	}
	return strings.Join(out, " ")
}
