// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The part of Go declares the target golang with the template of each family
// and its helpers. It sets a pointer receiver on each method of a double, and
// returns the names of the member that configures a method, of the fields of
// a call and of the fields of the results in the case of Go. It declares the
// types and the functions that its templates write beside a double, and its
// templates write the doubles and their tests in the conventions of Go.
func TestGolang(t *testing.T) {
	t.Parallel()

	templates, provides := stub.New().(plugin.TemplateProvider)
	assert.True(t, provides, "stub declares templates")

	t.Run("declares the template of each family for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, stub.New(), golang.Target, stubTemplate, testTemplate)
	})

	t.Run("writes the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t), golden.ShouldUpdate())
	})

	t.Run("sets a pointer receiver on each method under the name of the model", func(t *testing.T) {
		t.Parallel()

		store := generate(t, golang.Target, goTree, catalogPackage).doubles[storeName]
		get := declared(t, store, "Get")
		assert.NotNil(t, get.Receiver, "the method of a Go double has a receiver")
		expect.Equal(t, get.Receiver.Name, "s", "the receiver is s")
		expect.Equal(t, get.Receiver.Type.Spelling, "*StoreStub", "the receiver is a pointer to the double")
		expect.Equal(t, declared(t, store, "Split").Receiver.Name, "s0", "the receiver of Split avoids its parameter s")
		repo := declared(t, generate(t, golang.Target, goTree, genericPackage).doubles[repoName], "Get")
		expect.Equal(t, repo.Receiver.Type.Spelling, "*RepoStub[T]", "the receiver restates the type parameter T")
	})

	t.Run("names the member that configures each method after the method with On", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, catalogPackage).models[storeName]
		fields := make([]string, 0, len(got.Methods))
		for i := range got.Methods {
			fields = append(fields, got.Methods[i].Field)
		}
		assert.Equal(t, fields, []string{"OnGet", "OnPut", "OnLen", "OnLookup", "OnLog", "OnSplit", "OnClose"},
			"the member joins On and the name of its method")
	})

	{
		tests := []struct {
			name string
			give string
			want string
		}{
			{
				name: "names the fields of a call after the exported names of its parameters", give: "Get",
				want: "Ctx ID",
			},
			{
				name: "names the fields of a call after the names of its unnamed parameters", give: "Put",
				want: "A0 A1",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], tt.give)
				assert.Equal(t, paramFields(got.Params), tt.want,
					"a field of a call has the exported name of a parameter")
			})
		}
	}

	{
		tests := []struct {
			name        string
			givePkg     string
			giveSubject string
			give        string
			want        string
		}{
			{
				name: "names the one value of the results Result beside the error Err", givePkg: catalogPackage,
				giveSubject: storeName, give: "Get", want: "Result Err",
			},
			{
				name: "names an error that is the only result Err", givePkg: catalogPackage, giveSubject: storeName,
				give: "Put", want: "Err",
			},
			{
				name: "names the one result of a method that cannot fail Result", givePkg: catalogPackage,
				giveSubject: storeName, give: "Len", want: "Result",
			},
			{
				name: "names the results after their names", givePkg: catalogPackage, giveSubject: storeName,
				give: "Lookup", want: "Item Ok",
			},
			{
				name: "names no field of a method without a result", givePkg: catalogPackage, giveSubject: storeName,
				give: "Log",
			},
			{
				name: "numbers two values of the results", givePkg: catalogPackage, giveSubject: storeName,
				give: "Split", want: "Result0 Result1 Err",
			},
			{
				name: "names the error result _ Err", givePkg: edgePackage, giveSubject: edgeName, give: "Read",
				want: "R Err",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := method(t, generate(t, golang.Target, goTree, tt.givePkg).models[tt.giveSubject], tt.give)
				assert.Equal(t, resultFields(got.Returns), tt.want, "the results have a field of each result")
				expect.NotEqual(t, got.Return, "", "a Go method has a type of its fixed result")
			})
		}
	}

	t.Run("declares the helpers of the templates for the target golang", func(t *testing.T) {
		t.Parallel()

		funcs := templates.TemplateFuncs(golang.Target)
		camel, declared := funcs["stubCamel"].(func(string) string)
		assert.True(t, declared, "the target golang declares the helper stubCamel")
		expect.Equal(t, camel("StoreStub"), "storeStub", "the helper writes a name in camel case")
		expect.Contains(t, funcs, "stubSite", "the target golang declares the helper stubSite")
	})

	{
		tests := []struct {
			name        string
			giveSubject string
			want        []string
		}{
			{
				name:        "declares the Failing option with its constructor beside a double that can fail",
				giveSubject: storeName,
				want: []string{
					"StoreGetCall", "StoreGetStub", "StoreGetReturn", "WithStoreGet",
					"StorePutCall", "StorePutStub", "StorePutReturn", "WithStorePut",
					"StoreLenCall", "StoreLenStub", "StoreLenReturn", "WithStoreLen",
					"StoreLookupCall", "StoreLookupStub", "StoreLookupReturn", "WithStoreLookup",
					"StoreLogCall", "StoreLogStub", "WithStoreLog",
					"StoreSplitCall", "StoreSplitStub", "StoreSplitReturn", "WithStoreSplit",
					"StoreCloseCall", "StoreCloseStub", "StoreCloseReturn", "WithStoreClose",
					"StoreStubOption", "StoreStubStrict", "StoreStubDelegateTo", "StoreStubWithClock",
					"StoreStubWithRandSource", "StoreStubBenchMode", "StoreStubWithLatency",
					"NewStoreStub", "NewStoreStubDelayed", "StoreStubFailing", "NewStoreStubFailing",
				},
			},
			{
				name: "declares no Failing option beside a double that cannot fail", giveSubject: clockName,
				want: []string{
					"ClockNowCall", "ClockNowStub", "ClockNowReturn", "WithClockNow",
					"ClockZoneCall", "ClockZoneStub", "ClockZoneReturn", "WithClockZone",
					"ClockEveryCall", "ClockEveryStub", "ClockEveryReturn", "WithClockEvery",
					"ClockStubOption", "ClockStubStrict", "ClockStubDelegateTo", "ClockStubWithClock",
					"ClockStubWithRandSource", "ClockStubBenchMode", "ClockStubWithLatency",
					"NewClockStub", "NewClockStubDelayed",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, golang.Target, goTree, catalogPackage)
				assert.Equal(t, got.beside[tt.giveSubject], tt.want,
					"the record of the double declares each type and function that the template writes beside it")
				assert.Equal(t, gentest.Referents(got.tests[tt.giveSubject]), append([]string{tt.giveSubject + "Stub"},
					tt.want...), "the record of the checks refers to the double and to each declaration beside it")
			})
		}
	}

	t.Run("returns the reference of the record of the checks to a declaration of the double", func(t *testing.T) {
		t.Parallel()

		ref, declared := templates.TemplateFuncs(golang.Target)["stubRef"].(func(*emit.Struct, string) (*emit.TypeRef, error))
		assert.True(t, declared, "the target golang declares the helper stubRef")
		got, err := ref(generate(t, golang.Target, goTree, catalogPackage).tests[storeName], "NewStoreStub")
		assert.NoError(t, err, "the record of the checks of Store refers to NewStoreStub")
		expect.Equal(t, got.Spelling, "NewStoreStub", "the reference names the constructor")
	})

	t.Run("refuses a declaration that the record of the checks does not refer to", func(t *testing.T) {
		t.Parallel()

		ref, declared := templates.TemplateFuncs(golang.Target)["stubRef"].(func(*emit.Struct, string) (*emit.TypeRef, error))
		assert.True(t, declared, "the target golang declares the helper stubRef")
		_, err := ref(generate(t, golang.Target, goTree, catalogPackage).tests[clockName], "NewClockStubFailing")
		expect.HasError(t, err, "a template that names a declaration of no record fails")
	})
}

// rendered runs stub over the packages of the Go tree whose doubles render,
// and returns the files that the Go backend renders from them.
func rendered(t *testing.T) fstest.MapFS {
	t.Helper()

	p := stub.New()
	res := fixture(t, golang.Target, goTree, p, catalogPackage, parityPackage, genericPackage, edgePackage,
		handlesPackage, mixedPackage, sequencePackage).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	return gentest.Rendered(t, res, p, gobackend.New())
}
