// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/internal/lang"
)

// ID is the name of the plugin.
const ID plugin.ID = "suite"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.1.0"

// SuiteDirective puts an interface, a struct or a function under suite.
// SubjectDirective marks a struct of a consumer as a subject of the suite of a
// declaration. DeclarationParam is its positional param that names the
// declaration, and NewParam its param that names the factory of the subjects.
const (
	SuiteDirective   directive.Name     = "suite"
	SubjectDirective directive.Name     = "subject"
	DeclarationParam directive.ParamKey = "declaration"
	NewParam         directive.ParamKey = "new"
)

// Word is the word of the family of the suites, and ConformWord the word of
// the family of the conformance tests. TestTag is the tag of the family of the
// conformance tests.
const (
	Word        = "suite"
	ConformWord = "conform"
	TestTag     = sdk.TagTest
)

// The tree of each target contains the file template of the suites and the
// file template of the conformance tests.
const (
	suiteTemplate   = "suite.tmpl"
	conformTemplate = "conform.test.tmpl"
)

// The words that the names of the declarations of a suite join to the name of
// its declaration, as StoreChecks, RunStore and RunStoreOn join them to Store.
// The name of a conformance test joins testWord, the subject, conformsWord and
// the declaration, as in TestMemoryConformsToStore.
const (
	checksWord    = "Checks"
	suiteWord     = "Suite"
	runWord       = "Run"
	controlWord   = "Control"
	onWord        = "On"
	handleWord    = "Handle"
	checkWord     = "Check"
	optionWord    = "Option"
	methodsWord   = "Methods"
	funcWord      = "Func"
	callableWord  = "Callable"
	conformanceWd = "Conformance"
	testWord      = "Test"
	conformsWord  = "ConformsTo"
	newWord       = "New"
)

// Suite is the data of the record of the suite of one declaration, and the
// template of the family of the suites reads it.
//
// Subject is the type of the subjects of the suite: the interface at its
// witnesses, the interface of the methods of a struct, or the function type of
// a function. Methods is the interface of the exported methods of a struct, and
// nil for another declaration. Func is the signature of a function, and nil for
// another declaration. Declaration is the name of the declaration, and Names
// are the names of the declarations of the suite. Callables are the exported
// methods of an interface or a struct in the order of its member walk, and the
// function itself for a function.
type Suite struct {
	Subject     *emit.TypeRef
	Methods     *emit.Interface
	Func        *emit.Method
	Declaration string
	Names       Names
	Callables   []Callable
}

// Names are the names of the declarations of the suite of a declaration such
// as Store: the variable of the names of its callables StoreChecks, the
// function of its checks StoreSuite, the applications RunStore, RunStoreOn,
// ControlStore and ControlStoreOn, and the aliases StoreHandle, StoreCheck and
// StoreOption. Methods is the interface of the methods of a struct, as
// CatalogMethods, and Func the function type of a function, as SlugFunc. Each
// is empty for another declaration.
type Names struct {
	Checks    string
	Suite     string
	Run       string
	RunOn     string
	Control   string
	ControlOn string
	Handle    string
	Check     string
	Option    string
	Methods   string
	Func      string
}

// Callable is one callable of a declaration. Name is the name of the callable
// in the identities of its checks, as Store.Get, Field the field of the
// variable of the names that contains it, as Get, and Type the type of that
// field, as storeGet. A function is its own callable, so its Name is the name
// of the function, and its Field is empty.
type Callable struct {
	Name  string
	Field string
	Type  string
}

// Conformance is the data of the record of the conformance test of one
// subject, and the template of the family of the conformance tests reads it.
// Test is the name of the test function, Subject the name of the subject and
// Declaration the name of the declaration whose suite the test applies.
//
// The record of the test is a struct with two fields. The type of the first
// refers to the application of the suite, as RunStore, and the type of the
// second to the factory of the subjects. The layout of eidos qualifies the
// first by the package of the file of the suite.
type Conformance struct {
	Test        string
	Subject     string
	Declaration string
}

// parts are the parts of suite, one for each language that it serves.
var parts = []lang.Part{golangPart{}}

// suiteSchema is the directive of a suite.
var suiteSchema = directive.Schema{
	Name: SuiteDirective,
	Doc: "generates the suite of the declaration: the names of its callables, the identities and the checks of " +
		"its callables, and the functions that apply the checks to a subject and to a control",
}

// subjectSchema is the directive of a subject. eidos resolves the declaration
// to a type that the file of the subject can see, and the factory to a
// function in the scope of the subject.
var subjectSchema = directive.Schema{
	Name: SubjectDirective,
	Positional: []directive.ParamSpec{{
		Key: DeclarationParam, Type: directive.TypeReference, Resolution: directive.ResolveTypeInScope,
		Required: true, Doc: "the declaration under the directive suite whose checks the subject passes",
	}},
	Params: []directive.ParamSpec{{
		Key: NewParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope, Required: true,
		Doc: "the function that returns a fresh subject, without a parameter or from the handle of a check",
	}},
	Doc: "generates the test that applies the checks of the suite of the declaration to the subject",
}

// New returns suite. suite records a [Suite] into the family [Word] for each
// interface, struct or function under [SuiteDirective], and the template of
// the family writes the suite in the conventions of the language of the plan.
// The suite of a declaration has no checks of its own. A consumer adds checks
// to an application of the suite.
//
// suite records a [Conformance] into the family [ConformWord] under the tag
// [TestTag] for each struct under [SubjectDirective]. The template of the
// family writes a test that applies the suite of the declaration to a subject
// of the factory. The test calls the factory, which need not be exported, so
// a Go file of the family is in the package of the subject.
//
// suite runs after stub, whose doubles the planted defects of a suite build
// on. It serves Go, the language of its part, and leaves a declaration of
// another language alone. It reports [Uninstantiable] at a declaration whose
// type parameters its suite cannot instantiate, [NoSuite] at a subject whose
// declaration has no directive suite, and [BadFactory] at a subject whose
// factory an application cannot call. suite records nothing for a declaration
// that it reports. Each call returns a new plugin.
func New() plugin.Generator {
	b := sdk.NewPlugin(ID).
		Version(Version).
		Requires(stub.Capability).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: suiteTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: ConformWord, Template: conformTemplate})
	lang.Register(b, parts)
	return b.Handle(
		sdk.Directive(suiteSchema,
			sdk.OnInterface(lang.Served(parts, suiteOfInterface)),
			sdk.OnStruct(lang.Served(parts, suiteOfStruct)),
			sdk.OnFunction(lang.Served(parts, suiteOfFunction))),
		sdk.Directive(subjectSchema, sdk.OnStruct(lang.Served(parts, conformanceOf))),
	).Generator()
}

// suiteOfInterface records the suite of an interface. The subject type is the
// interface at the witnesses of its type parameters.
func suiteOfInterface(m *sdk.InterfaceMatch, _ lang.Part, e *sdk.Emitter) error {
	i := m.Interface
	b := m.Rules()
	subject := &emit.TypeRef{Spelling: i.Name, Target: i.ID}
	if len(i.TypeParams) > 0 {
		witnesses := b.Witnesses(i.TypeParams)
		if witnesses == nil {
			m.Errorf(Uninstantiable, "the type parameters of %s have no witnesses, so %s has no suite", i.Name, i.Name)
			return nil
		}
		for _, w := range witnesses {
			subject.Args = append(subject.Args, rules.EmitRef(w))
		}
	}
	s := Suite{Subject: subject, Declaration: i.Name, Names: namesOf(e, i.Name)}
	for _, method := range methodsOf(&b, i) {
		s.Callables = append(s.Callables, callableOf(e, i.Name, method.Name))
	}
	write(e, &i.ID, &s)
	return nil
}

// suiteOfStruct records the suite of a struct. The subject type is the
// interface of the exported methods of the struct, which the suite declares.
func suiteOfStruct(m *sdk.StructMatch, _ lang.Part, e *sdk.Emitter) error {
	st := m.Struct
	if len(st.TypeParams) > 0 {
		m.Errorf(Uninstantiable, "%s is a generic struct, so it has no suite", st.Name)
		return nil
	}
	b := m.Rules()
	names := namesOf(e, st.Name)
	names.Methods = pascal(e, methodsWord, st.Name)
	methods := &emit.Interface{Origin: st.ID, Name: names.Methods, Visibility: symbol.VisibilityPublic}
	s := Suite{
		Subject: &emit.TypeRef{Spelling: names.Methods}, Methods: methods, Declaration: st.Name, Names: names,
	}
	for _, method := range methodsOf(&b, st) {
		methods.Methods.Append(sdk.Mirror(names.Methods, method))
		s.Callables = append(s.Callables, callableOf(e, st.Name, method.Name))
	}
	write(e, &st.ID, &s)
	return nil
}

// suiteOfFunction records the suite of a function. The subject type is the
// function type of the function, which the suite declares.
func suiteOfFunction(m *sdk.FunctionMatch, _ lang.Part, e *sdk.Emitter) error {
	f := m.Function
	if len(f.TypeParams) > 0 {
		m.Errorf(Uninstantiable, "%s is a generic function, so it has no suite", f.Name)
		return nil
	}
	names := namesOf(e, f.Name)
	names.Func = pascal(e, funcWord, f.Name)
	signature := &node.Method{ID: f.ID, Name: f.Name, Visibility: f.Visibility, Params: f.Params, Returns: f.Returns}
	write(e, &f.ID, &Suite{
		Subject: &emit.TypeRef{Spelling: names.Func}, Func: sdk.Mirror(names.Func, signature), Declaration: f.Name,
		Names: names, Callables: []Callable{{Name: f.Name, Type: camel(e, callableWord, f.Name)}},
	})
	return nil
}

// write records the suite s of the declaration origin into the family of the
// suites. The record declares each name that the template writes: the variable
// of the names of the callables first, then the interface of the methods or the
// function type, the type of each callable, the functions and the aliases.
func write(e *sdk.Emitter, origin *symbol.Identity, s *Suite) {
	public := func(name string) symbol.Symbol {
		return &emit.Function{Origin: *origin, Name: name, Visibility: symbol.VisibilityPublic}
	}
	alias := func(name string) symbol.Symbol {
		return &emit.Alias{Origin: *origin, Name: name, Visibility: symbol.VisibilityPublic}
	}
	n := &s.Names
	decls := []symbol.Symbol{&emit.Variable{Origin: *origin, Name: n.Checks, Visibility: symbol.VisibilityPublic}}
	if s.Methods != nil {
		decls = append(decls, s.Methods)
	}
	if s.Func != nil {
		decls = append(decls, alias(n.Func))
	}
	for _, c := range s.Callables {
		decls = append(decls, &emit.Struct{Origin: *origin, Name: c.Type, Visibility: symbol.VisibilityPackage})
	}
	decls = append(decls, public(n.Suite), public(n.Run), public(n.RunOn), public(n.Control), public(n.ControlOn),
		alias(n.Handle), alias(n.Check), alias(n.Option))
	e.File().Record(*s, decls...)
}

// conformanceOf records the conformance test of a struct under the directive
// subject. It reports [NoSuite] where the declaration has no directive suite,
// and [BadFactory] where an application cannot call the factory. A factory
// without a parameter goes to the application Run, and a factory of one
// parameter, the handle of a check, to RunOn.
func conformanceOf(m *sdk.StructMatch, _ lang.Part, e *sdk.Emitter) error {
	st := m.Struct
	d := m.Directive()
	declared := d.Args[0].Target
	if !suited(m, &declared) {
		m.Errorf(NoSuite, "%s has no directive suite, so %s has no conformance test", declared.Name, st.Name)
		return nil
	}
	factory := d.Params[NewParam].Target
	found, _ := m.Reader().Lookup(factory)
	view, _ := m.Rules().CallableOf(found)
	if len(view.Params) > 1 {
		m.Errorf(BadFactory, "the factory %s takes %d parameters, so %s has no conformance test",
			factory.Name, len(view.Params), st.Name)
		return nil
	}
	run := pascal(e, declared.Name, runWord)
	if len(view.Params) == 1 {
		run = pascal(e, onWord, run)
	}
	if len(view.Returns) != 1 {
		m.Errorf(BadFactory, "the factory %s returns %d results, so %s has no conformance test",
			factory.Name, len(view.Returns), st.Name)
		return nil
	}
	record := &emit.Struct{Origin: st.ID, Name: camel(e, conformanceWd, st.Name), Visibility: symbol.VisibilityPackage}
	record.Fields.Append(
		&emit.Field{Name: runWord, Type: &emit.TypeRef{Spelling: run, Target: declared, Emitted: true}},
		&emit.Field{Name: newWord, Type: &emit.TypeRef{Spelling: factory.Name, Target: factory}},
	)
	e.File(TestTag).Record(Conformance{
		Test:    testWord + naming.Pascal(st.Name) + conformsWord + naming.Pascal(declared.Name),
		Subject: st.Name, Declaration: declared.Name,
	}, record)
	return nil
}

// suited reports whether the declaration id has the directive suite, in either
// spelling. The store lists the declarations of a directive by the name as
// written.
func suited(m *sdk.StructMatch, id *symbol.Identity) bool {
	decl, _ := m.Reader().Lookup(*id)
	for _, name := range []directive.Name{SuiteDirective, canonical} {
		for s := range m.Reader().ByDirective(name) {
			if s == decl {
				return true
			}
		}
	}
	return false
}

// canonical is the canonical name of the directive suite, under which the
// reader lists the declarations under it.
var canonical = directive.Schema{Plugin: string(ID), Name: SuiteDirective}.Canonical()

// methodsOf returns the exported methods of the member walk of the
// declaration decl, in the order of the walk.
func methodsOf(b *rules.Bound, decl symbol.Symbol) []*node.Method {
	set, _ := b.MembersOf(decl)
	var out []*node.Method
	for i := range set.Members {
		method, isMethod := set.Members[i].Symbol.(*node.Method)
		if isMethod && method.Visibility == symbol.VisibilityPublic {
			out = append(out, method)
		}
	}
	return out
}

// namesOf returns the names of the declarations of the suite of the
// declaration named decl, without the interface of the methods and the
// function type. Each name is exported in Go.
func namesOf(e *sdk.Emitter, decl string) Names {
	run, control := pascal(e, decl, runWord), pascal(e, decl, controlWord)
	return Names{
		Checks: pascal(e, checksWord, decl), Suite: pascal(e, suiteWord, decl),
		Run: run, RunOn: pascal(e, onWord, run), Control: control, ControlOn: pascal(e, onWord, control),
		Handle: pascal(e, handleWord, decl), Check: pascal(e, checkWord, decl), Option: pascal(e, optionWord, decl),
	}
}

// callableOf returns the callable of the method named method of the
// declaration named decl.
func callableOf(e *sdk.Emitter, decl, method string) Callable {
	return Callable{Name: decl + "." + method, Field: method, Type: camel(e, method, decl)}
}

// pascal returns head joined to tail in the case of an exported name of Go,
// which the settle of the Go plan writes.
func pascal(e *sdk.Emitter, tail, head string) string { return naming.Pascal(e.JoinName(tail, head)) }

// camel returns head joined to tail in the case of a name of Go that is not
// exported, which the settle of the Go plan writes.
func camel(e *sdk.Emitter, tail, head string) string { return naming.Camel(e.JoinName(tail, head)) }
