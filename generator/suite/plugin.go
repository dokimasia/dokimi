// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
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
// declaration. Its positional param DeclarationParam contains the name of the
// declaration, and its param NewParam contains the name of the factory of the
// subjects. CheckDirective binds a function of a consumer to a callable of a
// suite as a check. Its positional param CallableParam contains the callable,
// and its param DefectParam contains the name of the function of the planted
// defect.
const (
	SuiteDirective   directive.Name     = "suite"
	SubjectDirective directive.Name     = "subject"
	CheckDirective   directive.Name     = "check"
	DeclarationParam directive.ParamKey = "declaration"
	NewParam         directive.ParamKey = "new"
	CallableParam    directive.ParamKey = "callable"
	DefectParam      directive.ParamKey = "defect"
)

// KeyDeclaration, KeyCallable, KeySegment, KeyStatement and KeyDefect are the
// keys of the facts that the annotate half stamps on the function of each
// check that it binds. The facts contain the identity of the declaration under
// [SuiteDirective], the name of the callable such as Store.Get, the segment
// and the statement of the check, and the identity of the function of its
// planted defect.
const (
	KeyDeclaration meta.KeyName = "suite.declaration"
	KeyCallable    meta.KeyName = "suite.callable"
	KeySegment     meta.KeyName = "suite.segment"
	KeyStatement   meta.KeyName = "suite.statement"
	KeyDefect      meta.KeyName = "suite.defect"
)

// checksThat contains the words after the name of a check that open its doc
// comment, as in "KeepsTheItem checks that Get returns the item".
const checksThat = " checks that"

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
// its declaration, as StoreChecks, RunStoreSuite and RunStoreSuiteFromHandle
// join them to Store. The name of a conformance test joins testWord, the
// subject, conformsWord and the declaration, as in TestMemoryConformsToStore.
// The name of a falsification test joins testWord and the function of the
// suite, as in TestStoreSuite.
const (
	checksWord     = "Checks"
	suiteWord      = "Suite"
	runWord        = "Run"
	fromHandleWord = "FromHandle"
	handleWord     = "Handle"
	checkWord      = "Check"
	optionWord     = "Option"
	methodsWord    = "Methods"
	funcWord       = "Func"
	callableWord   = "Callable"
	conformanceWd  = "Conformance"
	proofsWord     = "Proofs"
	testWord       = "Test"
	conformsWord   = "ConformsTo"
	newWord        = "New"
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
// function itself for a function. Checks are the checks of the consumer that
// the annotate half binds to the callables, in the order of the callables.
type Suite struct {
	Subject     *emit.TypeRef
	Methods     *emit.Interface
	Func        *emit.Method
	Declaration string
	Names       Names
	Callables   []Callable
	Checks      []Check
}

// Names are the names of the declarations of the suite of a declaration such
// as Store: the variable of the names of its callables StoreChecks, the
// function of its checks StoreSuite, the two forms RunStoreSuite and
// RunStoreSuiteFromHandle of its application, and the aliases StoreHandle,
// StoreCheck and StoreOption. Methods is the interface of the methods of a
// struct, as CatalogMethods, and Func the function type of a function, as
// SlugFunc. Each is empty for another declaration.
type Names struct {
	Checks        string
	Suite         string
	Run           string
	RunFromHandle string
	Handle        string
	Check         string
	Option        string
	Methods       string
	Func          string
}

// Callable is one callable of a declaration. Name is the name of the callable
// in the identities of its checks, as Store.Get, Field the field of the
// variable of the names that contains it, as Get, and Type the type of that
// field, as StoreGet. The type is not exported, so the Go plan writes it as
// storeGet. A function is its own callable, so its Name is the name of the
// function, and its Field is empty.
type Callable struct {
	Name  string
	Field string
	Type  string
}

// Check is one check of a consumer that suite binds to a callable. Apply
// refers to the function of the check, and Defect to the function of its
// planted defect. Field and Type are the field and the type of the callable,
// as its [Callable] contains them. Method is the name of the function of the
// check, and the template writes the method that returns the identity of the
// check under that name. Segment is the name of the function in kebab case,
// and the last part of the identity of the check. Statement is the rest of
// the first sentence of the doc comment of the function after the words
// "checks that".
type Check struct {
	Apply     *emit.TypeRef
	Defect    *emit.TypeRef
	Field     string
	Type      string
	Method    string
	Segment   string
	Statement string
}

// Conformance is the data of the record of one test of the family of the
// tests, and the template of the family reads it. Test is the name of the
// test function, and Declaration the name of the declaration whose suite the
// test applies. Subject is the name of the subject of a conformance test, and
// empty for the falsification test of a suite.
//
// The record of a conformance test is a struct with two fields. The type of
// the first refers to the form of the application, as RunStoreSuite, and the
// type of the second to the factory of the subjects. The record of a
// falsification test is a struct with one field, whose type refers to the
// function of the suite, as StoreSuite. The layout of eidos qualifies each
// reference to the suite by the package of the file of the suite.
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
		"its callables, the functions that apply the checks to a subject, and the test of its planted defects",
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

// checkSchema is the directive of a check of a consumer. suite resolves the
// callable among the declarations under the directive suite, and eidos
// resolves the defect to a function in the scope of the check.
var checkSchema = directive.Schema{
	Name: CheckDirective,
	Positional: []directive.ParamSpec{{
		Key: CallableParam, Type: directive.TypeString, Required: true,
		Doc: "the callable of a declaration under the directive suite, such as Store.Get, or store.Store.Get with " +
			"the name of the package of the declaration",
	}},
	Params: []directive.ParamSpec{{
		Key: DefectParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope, Required: true,
		Doc: "the function that returns the planted defect of the check from the handle of the check",
	}},
	Doc: "binds the function to a callable of a suite as a check of the consumer, with its planted defect",
}

// generator is one instance of suite: the handles of the keys of the facts of
// a check. The composition registers the keys before the first phase call,
// and the handlers read them afterwards.
//
// # Concurrency
//
// The workers of a phase call read the handles concurrently, and no handler
// writes them.
type generator struct {
	declaration meta.Key[symbol.Identity]
	defect      meta.Key[symbol.Identity]
	callable    meta.Key[string]
	segment     meta.Key[string]
	statement   meta.Key[string]
}

// New returns suite, a plugin with an annotate half and a generate half. The
// annotate half stamps the facts of a check on each function under
// [CheckDirective] whose directive names a callable of a declaration under
// [SuiteDirective]. The generate half records a [Suite] into the family [Word]
// for each interface, struct or function under SuiteDirective, with the
// checks that the annotate half binds to its callables. The template of the
// family writes the suite in the conventions of the language of the plan.
//
// suite records a [Conformance] into the family [ConformWord] under the tag
// [TestTag] for each declaration under SuiteDirective, and the template of the
// family writes the test that falsifies each check of the suite on its
// planted defect. suite also records a Conformance for each struct under
// [SubjectDirective], and the template writes a test that applies the suite
// of the declaration to a subject of the factory. The test calls the
// factory, which need not be exported, so a Go file of the family is in the
// package of the subject.
//
// suite runs after stub, whose doubles the planted defects of a suite build
// on. It serves Go, the language of its part, and leaves a declaration of
// another language alone. It reports [Uninstantiable] at a declaration whose
// type parameters its suite cannot instantiate, [NoSuite] at a subject whose
// declaration has no directive suite, and [BadFactory] at a subject whose
// factory an application cannot call. It reports [UnknownCallable] and
// [BadCheck] at the directive of a check that it cannot bind. suite records
// and stamps nothing for a declaration that it reports. Each call returns a
// new plugin.
func New() sdk.Dual {
	g := &generator{}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Requires(stub.Capability).
		Keys(g.register).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: suiteTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: ConformWord, Template: conformTemplate})
	lang.Register(b, parts)
	return b.Handle(
		sdk.Directive(suiteSchema,
			sdk.OnInterface(lang.Served(parts, g.suiteOfInterface)),
			sdk.OnStruct(lang.Served(parts, g.suiteOfStruct)),
			sdk.OnFunction(lang.Served(parts, g.suiteOfFunction))),
		sdk.Directive(subjectSchema, sdk.OnStruct(lang.Served(parts, conformanceOf))),
		sdk.Directive(checkSchema, sdk.OnFunction(lang.Served(parts, g.concludeCheck))),
	).Dual()
}

// register claims the namespace of suite and registers the keys of the facts
// of a check.
//
// Error modes: the error of the claim, such as a namespace that another
// registrant claimed, and the errors of the registrations, such as a key that
// the registrant registered with another type, joined.
func (g *generator) register(r *meta.Registry) error {
	spec := func(name meta.KeyName, doc string) meta.KeySpec {
		return meta.KeySpec{Name: name, Kinds: []symbol.Kind{symbol.KindFunction}, Doc: doc}
	}
	var errs [6]error
	errs[0] = r.ClaimNamespace(KeyDeclaration.Namespace())
	g.declaration, errs[1] = meta.Register[symbol.Identity](r, spec(KeyDeclaration,
		"contains the identity of the declaration under the directive suite whose suite the check joins"))
	g.callable, errs[2] = meta.Register[string](r, spec(KeyCallable,
		"contains the name of the callable of the check, such as Store.Get"))
	g.segment, errs[3] = meta.Register[string](r, spec(KeySegment,
		"contains the last part of the identity of the check, such as keeps-the-item"))
	g.statement, errs[4] = meta.Register[string](r, spec(KeyStatement,
		"contains the sentence that the check proves of a subject"))
	g.defect, errs[5] = meta.Register[symbol.Identity](r, spec(KeyDefect,
		"contains the identity of the function that returns the planted defect of the check"))
	return errors.Join(errs[:]...)
}

// concludeCheck stamps the facts of a check on a function under the directive
// check. It stamps them where the callable of the directive belongs to one
// declaration under the directive suite. The function and its defect must also
// have the signatures of a check. Otherwise concludeCheck reports
// [UnknownCallable] or [BadCheck] at the directive.
func (g *generator) concludeCheck(m *sdk.FunctionMatch, _ lang.Part, st *sdk.Stamper) error {
	fn := m.Function
	d := m.Directive()
	decl, callable, resolved := resolveCallable(m, d, fn.Name)
	if !resolved {
		return nil
	}
	statement, stated := statementOf(fn)
	if !stated {
		m.ErrorfAt(BadCheck, d.Pos, "the doc comment of %s does not open with %q, so the check has no statement",
			fn.Name, fn.Name+checksThat)
		return nil
	}
	if !signed(m, d, fn) {
		return nil
	}
	sdk.Stamp(st, g.declaration, decl)
	sdk.Stamp(st, g.callable, callable)
	sdk.Stamp(st, g.segment, naming.Kebab(fn.Name))
	sdk.Stamp(st, g.statement, statement)
	sdk.Stamp(st, g.defect, d.Params[DefectParam].Target)
	return nil
}

// named is a declaration under the directive suite that the directive of a
// check names. name is the part of the callable of the directive that names
// the declaration. qualified is the callable with the name of the package of
// the declaration.
type named struct {
	decl      node.Declaration
	name      string
	qualified string
}

// resolveCallable resolves the callable of the directive d among the
// declarations under the directive suite. It returns the identity of the
// declaration and the name of the callable, such as Store.Get. check is the
// name of the function of the check. The callable of a directive is a method
// of an interface or a struct, or a function. Store.Get is a method, and Slug
// is a function. A callable such as store.Store.Get opens with the name of a
// package, and it belongs to a declaration of that package. resolveCallable
// reports [UnknownCallable] at d where no declaration has the name, where two
// declarations have it, and where the one declaration with the name has no
// such callable.
func resolveCallable(m *sdk.FunctionMatch, d *directive.Directive, check string) (symbol.Identity, string, bool) {
	written := d.Args[0].Str
	dot := strings.LastIndex(written, ".")
	head, method := written[:max(dot, 0)], written[dot+1:]
	var found []named
	for _, decl := range suites(m) {
		id := decl.Identity()
		pkg, _ := m.Reader().PackageOf(id)
		qualified := pkg.Name + "." + id.Name
		_, isFunction := decl.(*node.Function)
		switch {
		case isFunction && (written == id.Name || written == qualified):
			found = append(found, named{decl: decl, name: written, qualified: qualified})
		case !isFunction && (head == id.Name || head == qualified):
			found = append(found, named{decl: decl, name: head, qualified: qualified + "." + method})
		}
	}
	if len(found) == 0 {
		m.ErrorfAt(UnknownCallable, d.Pos, "no declaration under the directive suite has the callable %s, so %s is "+
			"no check", written, check)
		return symbol.Identity{}, "", false
	}
	if len(found) > 1 {
		spellings := make([]string, 0, len(found))
		for _, n := range found {
			spellings = append(spellings, n.qualified)
		}
		m.ErrorfAt(UnknownCallable, d.Pos, "%d declarations under the directive suite have the name %s, so %s is no "+
			"check. Qualify the callable with the name of its package, as %s",
			len(found), found[0].name, check, listed(spellings, "or"))
		return symbol.Identity{}, "", false
	}
	decl := found[0].decl
	id := decl.Identity()
	if _, isFunction := decl.(*node.Function); isFunction {
		return id, id.Name, true
	}
	b := m.Rules()
	var callables []string
	for _, c := range methodsOf(&b, decl) {
		callables = append(callables, c.Name)
	}
	if slices.Contains(callables, method) {
		return id, id.Name + "." + method, true
	}
	reason := fmt.Sprintf("%s has no callable %s, so %s is no check", head, method, check)
	switch {
	case len(callables) > 1:
		reason += fmt.Sprintf(". The callables of %s are %s", head, listed(callables, "and"))
	case len(callables) == 1:
		reason += fmt.Sprintf(". The callable of %s is %s", head, callables[0])
	}
	m.ErrorfAt(UnknownCallable, d.Pos, "%s", reason)
	return symbol.Identity{}, "", false
}

// statementOf returns the statement of the check fn, and false where the doc
// comment of fn does not open with the name of fn and the words checks that.
// The statement is the rest of the first sentence of the doc comment.
func statementOf(fn *node.Function) (string, bool) {
	text := strings.Join(strings.Fields(strings.Join(fn.Doc, " ")), " ")
	rest, opens := strings.CutPrefix(text, fn.Name+checksThat+" ")
	sentence, _, _ := strings.Cut(rest, ". ")
	return strings.TrimSuffix(sentence, "."), opens
}

// signed reports whether the check fn and the defect of its directive d have
// the signatures of a check. A check takes the handle of its check and returns
// nothing, and its defect returns the planted defect from the same handle.
// signed reports [BadCheck] at d where a signature differs.
func signed(m *sdk.FunctionMatch, d *directive.Directive, fn *node.Function) bool {
	defect := d.Params[DefectParam].Target
	b := m.Rules()
	check, _ := b.CallableOf(fn)
	found, _ := m.Reader().Lookup(defect)
	planted, _ := b.CallableOf(found)
	switch {
	case len(check.Params) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the check %s takes %d parameters, and a check takes the handle of its check "+
			"alone", fn.Name, len(check.Params))
	case len(check.Returns) > 0:
		m.ErrorfAt(BadCheck, d.Pos, "the check %s returns a result, and a check returns nothing", fn.Name)
	case len(planted.Params) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the defect %s of %s takes %d parameters, and a defect takes the handle of its "+
			"check alone", defect.Name, fn.Name, len(planted.Params))
	case len(planted.Returns) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the defect %s of %s returns %d results, and a defect returns the planted "+
			"defect alone", defect.Name, fn.Name, len(planted.Returns))
	default:
		return true
	}
	return false
}

// checksOf returns the checks that the annotate half binds to the callables
// of the declaration origin, in the order of the callables. The checks of one
// callable are in the order of the identities of their functions.
func (g *generator) checksOf(m sdk.Matcher, origin *symbol.Identity, callables []Callable) []Check {
	bound := map[string][]Check{}
	for id := range sdk.Subjects(m, g.declaration) {
		if decl, _ := sdk.FactOf(m, id, g.declaration); decl != *origin {
			continue
		}
		callable, _ := sdk.FactOf(m, id, g.callable)
		segment, _ := sdk.FactOf(m, id, g.segment)
		statement, _ := sdk.FactOf(m, id, g.statement)
		defect, _ := sdk.FactOf(m, id, g.defect)
		bound[callable] = append(bound[callable], Check{
			Apply:  &emit.TypeRef{Spelling: id.Name, Target: id},
			Defect: &emit.TypeRef{Spelling: defect.Name, Target: defect},
			Method: id.Name, Segment: segment, Statement: statement,
		})
	}
	var out []Check
	for _, c := range callables {
		for _, check := range bound[c.Name] {
			check.Field, check.Type = c.Field, c.Type
			out = append(out, check)
		}
	}
	return out
}

// listed returns words, two or more, as an English list whose last two words
// the conjunction joins, such as "Get and Put" or "a, b or c".
func listed(words []string, conjunction string) string {
	last := len(words) - 1
	return strings.Join(words[:last], ", ") + " " + conjunction + " " + words[last]
}

// suiteOfInterface records the suite of an interface. The subject type is the
// interface at the witnesses of its type parameters.
func (g *generator) suiteOfInterface(m *sdk.InterfaceMatch, _ lang.Part, e *sdk.Emitter) error {
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
	s.Checks = g.checksOf(m, &i.ID, s.Callables)
	write(e, &i.ID, &s)
	return nil
}

// suiteOfStruct records the suite of a struct. The subject type is the
// interface of the exported methods of the struct, which the suite declares.
func (g *generator) suiteOfStruct(m *sdk.StructMatch, _ lang.Part, e *sdk.Emitter) error {
	st := m.Struct
	if len(st.TypeParams) > 0 {
		m.Errorf(Uninstantiable, "%s is a generic struct, so it has no suite", st.Name)
		return nil
	}
	b := m.Rules()
	names := namesOf(e, st.Name)
	names.Methods = e.JoinName(methodsWord, st.Name)
	methods := &emit.Interface{Origin: st.ID, Name: names.Methods, Visibility: symbol.VisibilityPublic}
	s := Suite{
		Subject: &emit.TypeRef{Spelling: names.Methods}, Methods: methods, Declaration: st.Name, Names: names,
	}
	for _, method := range methodsOf(&b, st) {
		methods.Methods.Append(sdk.Mirror(names.Methods, method))
		s.Callables = append(s.Callables, callableOf(e, st.Name, method.Name))
	}
	s.Checks = g.checksOf(m, &st.ID, s.Callables)
	write(e, &st.ID, &s)
	return nil
}

// suiteOfFunction records the suite of a function. The subject type is the
// function type of the function, which the suite declares.
func (g *generator) suiteOfFunction(m *sdk.FunctionMatch, _ lang.Part, e *sdk.Emitter) error {
	f := m.Function
	if len(f.TypeParams) > 0 {
		m.Errorf(Uninstantiable, "%s is a generic function, so it has no suite", f.Name)
		return nil
	}
	names := namesOf(e, f.Name)
	names.Func = e.JoinName(funcWord, f.Name)
	signature := &node.Method{ID: f.ID, Name: f.Name, Visibility: f.Visibility, Params: f.Params, Returns: f.Returns}
	callables := []Callable{{Name: f.Name, Type: e.JoinName(callableWord, f.Name)}}
	write(e, &f.ID, &Suite{
		Subject: &emit.TypeRef{Spelling: names.Func}, Func: sdk.Mirror(names.Func, signature), Declaration: f.Name,
		Names: names, Callables: callables, Checks: g.checksOf(m, &f.ID, callables),
	})
	return nil
}

// write records the suite s of the declaration origin into the family of the
// suites. The record declares each name that the template writes: the variable
// of the names of the callables first, then the interface of the methods or the
// function type, the type of each callable, the functions and the aliases.
// write also records the falsification test of the suite into the family of
// the tests.
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
	decls = append(decls, public(n.Suite), public(n.Run), public(n.RunFromHandle),
		alias(n.Handle), alias(n.Check), alias(n.Option))
	e.File().Record(*s, decls...)
	record := &emit.Struct{
		Origin: *origin, Name: e.JoinName(proofsWord, s.Declaration), Visibility: symbol.VisibilityPackage,
	}
	record.Fields.Append(&emit.Field{
		Name: suiteWord, Type: &emit.TypeRef{Spelling: n.Suite, Target: *origin, Emitted: true},
	})
	e.File(TestTag).Record(Conformance{Test: e.JoinName(n.Suite, testWord), Declaration: s.Declaration}, record)
}

// conformanceOf records the conformance test of a struct under the directive
// subject. It reports [NoSuite] where the declaration has no directive suite,
// and [BadFactory] where an application cannot call the factory. A factory
// without a parameter goes to the form Run of the application, and a factory
// of one parameter, the handle of a check, to the form RunFromHandle.
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
	names := namesOf(e, declared.Name)
	run := names.Run
	if len(view.Params) == 1 {
		run = names.RunFromHandle
	}
	if len(view.Returns) != 1 {
		m.Errorf(BadFactory, "the factory %s returns %d results, so %s has no conformance test",
			factory.Name, len(view.Returns), st.Name)
		return nil
	}
	record := &emit.Struct{
		Origin:     st.ID,
		Name:       e.JoinName(conformanceWd, st.Name),
		Visibility: symbol.VisibilityPackage,
	}
	record.Fields.Append(
		&emit.Field{Name: runWord, Type: &emit.TypeRef{Spelling: run, Target: declared, Emitted: true}},
		&emit.Field{Name: newWord, Type: &emit.TypeRef{Spelling: factory.Name, Target: factory}},
	)
	e.File(TestTag).Record(Conformance{
		Test:    e.JoinName(declared.Name, e.JoinName(conformsWord, e.JoinName(st.Name, testWord))),
		Subject: st.Name, Declaration: declared.Name,
	}, record)
	return nil
}

// suited reports whether the declaration id has the directive suite.
func suited(m *sdk.StructMatch, id *symbol.Identity) bool {
	return slices.ContainsFunc(suites(m), func(decl node.Declaration) bool { return decl.Identity() == *id })
}

// suites returns the declarations under the directive suite, in either
// spelling. The store lists the declarations of a directive by the name as
// written.
func suites(m sdk.SubjectMatch) []node.Declaration {
	var out []node.Declaration
	for _, name := range []directive.Name{SuiteDirective, canonical} {
		for s := range m.Reader().ByDirective(name) {
			if decl, declares := s.(node.Declaration); declares {
				out = append(out, decl)
			}
		}
	}
	return out
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
// function type. Each name is in the neutral convention of an emitted name,
// and the backend of the plan writes it in the case of its target.
func namesOf(e *sdk.Emitter, decl string) Names {
	suite := e.JoinName(suiteWord, decl)
	run := e.JoinName(suite, runWord)
	return Names{
		Checks: e.JoinName(checksWord, decl), Suite: suite, Run: run, RunFromHandle: e.JoinName(fromHandleWord, run),
		Handle: e.JoinName(handleWord, decl), Check: e.JoinName(checkWord, decl), Option: e.JoinName(optionWord, decl),
	}
}

// callableOf returns the callable of the method named method of the
// declaration named decl.
func callableOf(e *sdk.Emitter, decl, method string) Callable {
	return Callable{Name: decl + "." + method, Field: method, Type: e.JoinName(method, decl)}
}
