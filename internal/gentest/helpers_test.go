// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"cmp"
	"slices"
	"testing/fstest"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/gentest"
)

// The Go tree of the cases and its package catalog, the struct Address of
// the package and its field Street, and the struct Warehouse.
const (
	goTree         = "testdata/golang"
	catalogPackage = "example.com/shop/catalog"
	addressName    = "Address"
	streetName     = "Street"
	warehouseName  = "Warehouse"
	alikeName      = "Alike"
)

// The version of the plugins of the cases, the name of the recorder, and the
// templates of the two families of the recorder.
const (
	version                = "0.1.0"
	recorderID   plugin.ID = "recorder"
	recordFamily           = "recorder.tmpl"
	checkFamily            = "recorder.test.tmpl"
)

// recorded is the code of the findings that the recorder reports.
var recorded = diag.MustRegister("GENTESTTEST", diag.CodeSpec{
	Number:  2,
	Meaning: "a declaration of a fixture that the recorder of the tests of gentest reports",
})

// subject is the data of a record of the recorder in its primary family, and
// check the data of a record in the family of its tests.
type (
	subject struct{ Name string }
	check   struct{ Name string }
)

// recording is what the recorder does at each struct of a fixture: the
// structs that it records in both of its families, the structs whose record
// in the primary family also declares a struct and a function after its
// first struct, the number of findings that it reports at each struct, and
// the structs at whose first field it reports one finding. again are the
// structs that get a second record of their origin in the primary family,
// with data of its own, and alike the structs whose record in the primary
// family has the data of the name Alike. body is the template of the primary
// family, and an empty body writes nothing for a record. external puts the Go
// files of the family of the tests into the external test package of their
// directory.
type recording struct {
	records  []string
	beside   []string
	again    []string
	alike    []string
	findings map[string]int
	fields   []string
	body     string
	external bool
}

// generate runs the recorder of r over the package catalog of the Go tree
// under the target golang, and returns the recorder, the fixture and the
// result.
func generate(tb assert.TB, r recording) (plugin.Generator, *gentest.Fixture, plugintest.Result) {
	tb.Helper()

	p := recorder(r)
	f := gentest.New(tb, &gentest.Spec{
		Dir: goTree, Packages: []string{catalogPackage}, Target: golang.Target, Plugins: []plugin.Plugin{p},
	})
	res := f.Generate(tb, p)
	assert.NoError(tb, res.Err, "the recorder runs whole")
	return p, f, res
}

// recorder returns a generator of the target golang that does at each struct
// what r states.
func recorder(r recording) plugin.Generator {
	return sdk.NewPlugin(recorderID).
		Version(version).
		Output(plugin.Output{Per: plugin.PerSource, Word: string(recorderID), Template: recordFamily}).
		Output(plugin.Output{
			Tag:      sdk.TagTest,
			Per:      plugin.PerSource,
			Word:     string(recorderID),
			Template: checkFamily,
			External: r.external,
		}).
		For(golang.Target, sdk.Templates(fstest.MapFS{
			recordFamily: {Data: []byte(cmp.Or(r.body, "{{range .Records}}{{end}}"))},
			checkFamily:  {Data: []byte("{{range .Records}}{{end}}")},
		})).
		Handle(sdk.OnStruct(func(m *sdk.StructMatch, e *sdk.Emitter) error {
			s := m.Struct
			for range r.findings[s.Name] {
				m.Errorf(recorded, "%s is reported", s.Name)
			}
			if slices.Contains(r.fields, s.Name) {
				m.ErrorfAt(recorded, s.Fields[0].Pos, "the field of %s is reported", s.Name)
			}
			if slices.Contains(r.records, s.Name) {
				record := []symbol.Symbol{&emit.Struct{Origin: s.ID, Name: s.Name + "Record"}}
				if slices.Contains(r.beside, s.Name) {
					record = append(record, &emit.Struct{Origin: s.ID, Name: s.Name + "Option"},
						&emit.Function{Origin: s.ID, Name: "New" + s.Name + "Record"})
				}
				data := subject{Name: s.Name}
				if slices.Contains(r.alike, s.Name) {
					data.Name = alikeName
				}
				e.File().Record(data, record...)
				if slices.Contains(r.again, s.Name) {
					e.File().Record(subject{Name: s.Name + "Again"}, &emit.Struct{Origin: s.ID, Name: s.Name + "Again"})
				}
				e.File(sdk.TagTest).Record(check{Name: s.Name}, &emit.Struct{Origin: s.ID, Name: s.Name + "Check"})
			}
			return nil
		})).
		Generator()
}
