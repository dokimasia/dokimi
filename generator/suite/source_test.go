// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/suite"
)

// The generators of the values of a suite, by the list of the overrides of
// the suite that contains the type, and fromEngine for a type that the shape
// reader of the property engine reads.
const (
	fromEngine    = "engine"
	fromStarts    = "starts"
	fromGenerated = "generated"
	fromStreams   = "streams"
	fromConstants = "constants"
)

// The suite gives each type in the inputs of a callable a generator: the
// shape reader of the property engine, or an override of the suite.
func TestSource(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		p := suite.New()
		f := fixture(t, goTree, p, kindsPackage)
		got := recorded(t, f, p)
		kinds := map[reflect.Kind]string{
			reflect.Bool: fromEngine, reflect.Int: fromEngine, reflect.Int8: fromEngine, reflect.Int16: fromEngine,
			reflect.Int32: fromEngine, reflect.Int64: fromEngine, reflect.Uint: fromEngine, reflect.Uint8: fromEngine,
			reflect.Uint16: fromEngine, reflect.Uint32: fromEngine, reflect.Uint64: fromEngine,
			reflect.Uintptr: fromEngine, reflect.Float32: fromEngine, reflect.Float64: fromEngine,
			reflect.Complex64: fromEngine, reflect.Complex128: fromEngine, reflect.Array: fromEngine,
			reflect.Chan: fromStreams, reflect.Func: fromGenerated, reflect.Interface: fromGenerated,
			reflect.Map: fromEngine, reflect.Pointer: fromEngine, reflect.Slice: fromEngine,
			reflect.String: fromEngine, reflect.Struct: fromStarts, reflect.UnsafePointer: fromConstants,
		}
		tests := []struct {
			name      string
			giveDecl  string
			giveOwner string
			giveKind  symbol.Kind
			giveType  func(node.Declaration) *node.TypeRef
		}{
			{
				name: "records a generator of each kind of reflect in a field", giveDecl: "Every", giveOwner: "Fields",
				giveKind: symbol.KindField,
				giveType: func(d node.Declaration) *node.TypeRef { return fieldType(t, d) },
			},
			{
				name:     "records a generator of each kind of reflect in an element of a list",
				giveDecl: "Every", giveOwner: "Elements", giveKind: symbol.KindField,
				giveType: func(d node.Declaration) *node.TypeRef { return fieldType(t, d).Elems[0] },
			},
			{
				name:     "records a generator of each kind of reflect in a result of a method of an input",
				giveDecl: "Every", giveOwner: "Results", giveKind: symbol.KindMethod,
				giveType: func(d node.Declaration) *node.TypeRef { return methodOf(t, d).Returns[0].Type },
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				s := got.suites[tt.giveDecl]
				for kind := reflect.Bool; kind <= reflect.UnsafePointer; kind++ {
					want, listed := kinds[kind]
					assert.True(t, listed, "the spec states the generator of the kind "+kind.String())
					typ := tt.giveType(f.Find(t, tt.giveKind, tt.giveOwner, kindName(kind)))
					expect.Equal(t, generatorOf(s, typ.Spelling), want,
						"the suite of "+tt.giveDecl+" has the generator of the type "+typ.Spelling)
				}
				expect.Empty(t, ungenerated(got, tt.giveDecl), "each type in the inputs of "+tt.giveDecl+
					" has a generator")
			})
		}

		t.Run("records a generator of each kind of reflect in a parameter", func(t *testing.T) {
			t.Parallel()

			s := got.suites["Params"]
			for kind := reflect.Bool; kind <= reflect.UnsafePointer; kind++ {
				want, listed := kinds[kind]
				assert.True(t, listed, "the spec states the generator of the kind "+kind.String())
				param := called(t, s, kindName(kind)).Params[0]
				generator := generatorOf(s, spelled(param.Type))
				if len(param.Funcs) > 0 {
					generator = fromGenerated
				}
				expect.Equal(t, generator, want, "the parameter of Params."+kindName(kind)+" has its generator")
			}
			expect.Empty(t, ungenerated(got, "Params"), "each type in the inputs of Params has a generator")
		})
	})
}

// kindName returns the name of the declarations of the package kinds that
// take a value of the kind k: the name of the kind in Pascal case, as Ptr and
// UnsafePointer.
func kindName(k reflect.Kind) string {
	return naming.Pascal(strings.ReplaceAll(k.String(), ".", " "))
}

// fieldType returns the type of the field d. It stops the test where d is no
// field.
func fieldType(t *testing.T, d node.Declaration) *node.TypeRef {
	t.Helper()

	f, isField := d.(*node.Field)
	assert.True(t, isField, d.Identity().Name+" is a field")
	return f.Type
}

// methodOf returns the method d. It stops the test where d is no method.
func methodOf(t *testing.T, d node.Declaration) *node.Method {
	t.Helper()

	m, isMethod := d.(*node.Method)
	assert.True(t, isMethod, d.Identity().Name+" is a method")
	return m
}

// generatorOf returns the generator that the suite s records for the type
// spelled spelling: the list of the overrides that contains the type, or
// fromEngine for a type that no override covers.
func generatorOf(s suite.Suite, spelling string) string {
	has := func(t *emit.TypeRef) bool { return spelled(t) == spelling }
	switch {
	case slices.ContainsFunc(s.Starts, func(st suite.Start) bool { return has(st.Type) }):
		return fromStarts
	case slices.ContainsFunc(s.Generated, func(g suite.Generated) bool { return has(g.Type) }):
		return fromGenerated
	case slices.ContainsFunc(s.Streams, func(st suite.Stream) bool { return has(st.Type) }):
		return fromStreams
	case slices.ContainsFunc(s.Constants, func(c suite.Constant) bool { return has(c.Type) }):
		return fromConstants
	}
	return fromEngine
}
