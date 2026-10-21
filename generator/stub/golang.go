// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub

import (
	"embed"
	"strconv"
	"text/template"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// The member that configures a method joins golangOn and the name of the
// method. A result without a name is golangErr for the error and
// golangResult for the one other result. Where a method has two or more other
// results, each is golangResult with its position among them.
const (
	golangOn     = "On"
	golangErr    = "Err"
	golangResult = "Result"
)

// The templates of Go write these declarations beside a double such as
// StoreStub: the type of its options StoreStubOption, the options
// StoreStubStrict, StoreStubDelegateTo, StoreStubWithClock,
// StoreStubWithRandSource, StoreStubBenchMode and StoreStubWithLatency, and
// the constructors NewStoreStub and NewStoreStubDelayed. A double whose
// method can fail also has the option StoreStubFailing and the constructor
// NewStoreStubFailing. The option of the method Get of the subject Store is
// WithStoreGet.
const (
	golangOption     = "Option"
	golangStrict     = "Strict"
	golangDelegateTo = "DelegateTo"
	golangWithClock  = "WithClock"
	golangWithRand   = "WithRandSource"
	golangBenchMode  = "BenchMode"
	golangLatency    = "WithLatency"
	golangFailing    = "Failing"
	golangNew        = "New"
	golangDelayed    = "Delayed"
	golangWith       = "With"
)

// golangFault opens the name of the helper of the stub of a method that arms
// an error value for the next call, as in FaultNotFound for ErrNotFound.
const golangFault = "Fault"

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of stub for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Options returns the template tree and the helpers of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(golangTemplates, golang.Target),
		sdk.Funcs(template.FuncMap{camelHelper: naming.Camel, siteHelper: siteOf, refHelper: referentOf}),
	}
}

// structural reports false, because a value of a Go struct is a value of the
// struct alone. A double of a struct implements its methods, and an
// interface of Go has no fields.
func (golangPart) structural() bool { return false }

// referable reports true, because the file of a double is in the package of
// its subject.
func (golangPart) referable(symbol.Visibility) bool { return true }

// receiver sets a pointer receiver named name on the method m of a double,
// as Go declares a method that changes its value, and returns m. The
// receiver restates the type parameters tps of a generic double.
func (golangPart) receiver(m *emit.Method, tps []*emit.TypeParam, name string) *emit.Method {
	m.Receives.Args = sdk.TypeArgs(tps)
	golang.PointerReceiver(m).Receiver.Name = name
	return m
}

// constraint reports whether the satellite of Go marks the interface id as an
// interface with type-set elements.
func (golangPart) constraint(m sdk.Matcher, id *symbol.Identity) bool {
	constraint, _ := sdk.FactOf(m, *id, meta.Named[bool](golang.ConstraintInterfaceKey))
	return constraint
}

// names returns the names of a Go method of a double. The member that
// configures the method is On followed by the exported name of the method,
// such as OnGet for Get. A field of the call has the exported name of its
// parameter, such as ID for id, so a test in another package reads the field.
// The fields of the results follow [golangResults], and a method answers its
// results as a value of a type of its own. A parameter of the setter Returns
// is the field of its result in camel case, such as item for Item. The methods
// of the double declare the names of the results where [golangNamed] reports
// it. The error of a Go method is a result, so a call has no field of a
// failure.
func (golangPart) names(method string, params []string, rs []rules.ReturnView, taken map[string]bool) Names {
	n := Names{Field: golangOn + naming.Pascal(method), Grouped: true, Named: golangNamed(rs, taken)}
	for _, p := range params {
		n.Args = append(n.Args, naming.Pascal(p))
	}
	n.Results = golangResults(rs)
	for _, f := range n.Results {
		n.Setters = append(n.Setters, naming.Camel(f))
	}
	return n
}

// golangResults returns the exported name of each result with a name. A
// result without a name, or with the name _, is Err for the error and Result
// for the one other result. Where a method has two or more other results,
// each is Result followed by its position. The error of a Go method is its
// last result, so the position of each other result is its position among
// the other results.
func golangResults(rs []rules.ReturnView) []string {
	values := 0
	for _, r := range rs {
		if r.Role != rules.ReturnError {
			values++
		}
	}
	out := make([]string, 0, len(rs))
	for i, r := range rs {
		switch {
		case r.Name != "" && r.Name != "_":
			out = append(out, naming.Pascal(r.Name))
		case r.Role == rules.ReturnError:
			out = append(out, golangErr)
		case values == 1:
			out = append(out, golangResult)
		default:
			out = append(out, golangResult+strconv.Itoa(i))
		}
	}
	return out
}

// golangNamed reports whether the methods of a double declare the names of
// the returns rs. Go declares the names of every result or of none, so each
// return needs a name other than _. A name that a parameter or the receiver
// in taken has does not compile, so the double then declares no names. Go
// refuses two results of one name, so no two returns share a name.
func golangNamed(rs []rules.ReturnView, taken map[string]bool) bool {
	for _, r := range rs {
		if r.Name == "" || r.Name == "_" || taken[r.Name] {
			return false
		}
	}
	return true
}

// delegate returns value for an interface, and a pointer to value for a
// struct, whose method set includes the methods with a pointer receiver.
func (golangPart) delegate(value *emit.TypeRef, class bool) *emit.TypeRef {
	if !class {
		return value
	}
	return &emit.TypeRef{Spelling: "*" + value.Spelling, Form: symbol.FormOptional, Elems: []*emit.TypeRef{value}}
}

// fault returns Fault followed by base in Pascal case, such as FaultNotFound
// for the base NotFound of ErrNotFound.
func (golangPart) fault(base string) string { return golangFault + naming.Pascal(base) }

// construction returns seed, the function that constructs a builder from the
// seed, such as NewAccountBuilder. A Go double calls the function itself.
func (golangPart) construction(_, seed string) (ref, method string) { return seed, "" }

// declarations returns the declarations that the templates of Go write
// beside the double: the type of the calls, the stub and the fixed results of
// each method and the option of each method, then the type of the options,
// the options and the constructors of the double. A method without results
// has no type of fixed results. A check of Go names a declaration of another
// package through the name of that package, as in catalog.NewStoreStub.
func (golangPart) declarations(double *emit.Struct, m *Model) []symbol.Symbol {
	name := double.Name
	function := func(n string) symbol.Symbol {
		return &emit.Function{Origin: double.Origin, Name: n, Visibility: double.Visibility}
	}
	typ := func(n string) symbol.Symbol {
		return &emit.Struct{Origin: double.Origin, Name: n, Visibility: double.Visibility}
	}
	var out []symbol.Symbol
	for i := range m.Methods {
		c := &m.Methods[i]
		out = append(out, typ(c.Call), typ(c.Stub))
		if len(c.Returns) > 0 {
			out = append(out, typ(c.Return))
		}
		out = append(out, function(golangWith+m.Subject+c.Name))
	}
	out = append(out,
		&emit.Alias{Origin: double.Origin, Name: name + golangOption, Visibility: double.Visibility, Defined: true},
		function(name+golangStrict), function(name+golangDelegateTo), function(name+golangWithClock),
		function(name+golangWithRand), function(name+golangBenchMode), function(name+golangLatency),
		function(golangNew+name), function(golangNew+name+golangDelayed),
	)
	if m.Fallible {
		out = append(out, function(name+golangFailing), function(golangNew+name+golangFailing))
	}
	return out
}
