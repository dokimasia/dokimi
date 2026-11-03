// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/internal/lang"
	"go.dokimi.dev/dokimi/internal/nametemplate"
	"go.dokimi.dev/dokimi/spec"
)

// ID is the name of the plugin.
const ID plugin.ID = "fake"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.1.0"

// FakeDirective puts an interface under fake.
const FakeDirective directive.Name = "fake"

// AbsentParam is the param of [FakeDirective] whose name template names the
// sentinel that a load of an absent key returns, such as
// absent=ErrNo{{.Name}}. The only action of the template is {{.Name}}, the
// name of the declaration. DefaultAbsent is its default, which gives the
// declaration Store the sentinel ErrStoreFakeNotFound.
const (
	AbsentParam   directive.ParamKey = "absent"
	DefaultAbsent string             = "Err{{.Name}}FakeNotFound"
)

// nameField is the one field that the template of [AbsentParam] reads.
const nameField = "Name"

// Word is the word of the family of the fakes and of the family of their
// tests. TestTag is the tag of the family of the tests.
const (
	Word    = "fake"
	TestTag = sdk.TagTest
)

// The tree of each target contains the file template of the fakes and the
// file template of their tests.
const (
	fakeTemplate = "fake.tmpl"
	testTemplate = "fake.test.tmpl"
)

// The name of a fake joins fakeWord to its declaration, and the name of its
// constructor joins newWord to the fake. The name of the test of a fake joins
// testWord to the fake, and the name of the record of the test joins
// conformanceWord to the fake.
const (
	fakeWord        = "Fake"
	newWord         = "New"
	testWord        = "Test"
	conformanceWord = "Conformance"
)

// The receiver of a method of a fake takes receiverName, the zero value of
// the result at position i takes resultLocal followed by i, and a load names
// the report of a present key presentLocal. Each name avoids the parameters
// of its method.
const (
	receiverName = "f"
	resultLocal  = "r"
	presentLocal = "present"
)

// Model is the data of the record of one fake, and the templates of both
// families read it.
//
// Fake is the name of the fake, New the name of its constructor and Absent the
// name of the sentinel of an absent key. Subject is the name of the
// declaration. Form is the reference form of the declaration and State the
// state of the form. Key is the type of the keys of a fake of the state map,
// and nil for the state list. Value is the type of the values that the fake
// stores or appends. Guard is the declaration at the witnesses of its type
// parameters, which a compile-time check asserts. Witnesses are the witnesses
// of the type parameters, and nil for a declaration without type parameters.
// Methods are the methods of the fake, in the order of the member walk of the
// declaration.
type Model struct {
	Key       *emit.TypeRef
	Value     *emit.TypeRef
	Guard     *emit.TypeRef
	Fake      string
	New       string
	Absent    string
	Subject   string
	Form      spec.ReferenceFormID
	State     spec.State
	Witnesses []*emit.TypeRef
	Methods   []Method
}

// Method is one method of a fake.
//
// Name is the name of the method and Receiver the name of its receiver. Effect
// is the effect of the method on the state of the fake. Context is the
// parameter of the context of a method that returns an error, and empty
// otherwise, so the method returns the error of a done context first. Key is
// the expression of the key of a store or of a load, such as v.Key, and empty
// for an append. Value is the parameter of the value that a store stores or
// that an append appends, and empty for a load. Present is the local of a load
// that reports a present key. Results are the results of the method, in
// order.
type Method struct {
	Name     string
	Receiver string
	Effect   spec.Effect
	Context  string
	Key      string
	Value    string
	Present  string
	Results  []Result
}

// Result is one result of a method of a fake. Type is its type and Error
// reports the error. Local is the local of the zero value of a result that is
// no error.
type Result struct {
	Type  *emit.TypeRef
	Local string
	Error bool
}

// parts are the parts of fake, one for each language that it serves.
var parts = []golangPart{{}}

// schema is the directive of fake.
var schema = directive.Schema{
	Name: FakeDirective,
	Params: []directive.ParamSpec{{
		Key: AbsentParam, Type: directive.TypeString, Default: DefaultAbsent,
		Doc: "the name template of the sentinel that a load of an absent key returns, such as ErrNo{{.Name}}",
	}},
	Doc: "generates an in-memory implementation of the declaration from the reference forms of the conformance " +
		"spec, and the test that applies the suite of the declaration to it",
}

// New returns fake. fake records a [Model] into the family [Word] for each
// interface under [FakeDirective] whose callables a reference form of the
// conformance spec covers, by the rule of [derive.ReferenceOf]. The template
// of the family writes the fake, its constructor and the sentinel of an
// absent key in the conventions of the language of the plan. fake records the
// test of each fake into the family [Word] under the tag [TestTag]. The
// template of that family writes a test that applies the suite of the
// declaration to a new fake. The test uses the fake through the exported API
// of its package, so a Go file of the family is in the external test package
// of its directory. fake serves Go, the language of its part, and leaves a
// declaration of another language alone.
//
// fake reports [NoSuite] at an interface without the directive suite,
// [BadAbsent] at an interface whose param [AbsentParam] is no name template of
// a sentinel, [Unimplementable] at an interface that a fake cannot implement
// or that the test of a fake cannot refer to, and [Refused] at the callable of
// an interface that no reference form covers. fake records nothing for an
// interface that it reports. Each call returns a new plugin.
func New() plugin.Generator {
	b := sdk.NewPlugin(ID).
		Version(Version).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: fakeTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: Word, Template: testTemplate, External: true})
	lang.Register(b, parts)
	return b.Handle(sdk.Directive(schema, sdk.OnInterface(lang.Served(parts, generate)))).Generator()
}

// invocation is the state of one fake: the match of its interface, the part
// of its language, the emitter and the rules of the language.
type invocation struct {
	m *sdk.InterfaceMatch
	p golangPart
	e *sdk.Emitter
	b rules.Bound
}

// generate records the fake of the interface of m and its test.
func generate(m *sdk.InterfaceMatch, p golangPart, e *sdk.Emitter) error {
	in := &invocation{m: m, p: p, e: e, b: m.Rules()}
	i := m.Interface
	template := m.Directive().Params[AbsentParam].Str
	absent, err := nametemplate.Segments(template, nameField)
	switch {
	case !suite.Suited(m, &i.ID):
		m.Errorf(NoSuite, "%s has no directive suite, so it has no fake", i.Name)
		return nil
	case err != nil:
		m.Errorf(BadAbsent, "the param absent of %s is not the name template of a sentinel: %v, so it has no fake",
			i.Name, err)
		return nil
	}
	methods, witnesses, why := in.implementable()
	if why != "" {
		m.Errorf(Unimplementable, "%s %s, so it has no fake", i.Name, why)
		return nil
	}
	facts := make([]derive.Facts, len(methods))
	for k, method := range methods {
		facts[k] = in.factsOf(method)
	}
	ref, refused := derive.ReferenceOf(facts)
	if refused != nil {
		at := slices.IndexFunc(methods, func(method *node.Method) bool { return method.Name == refused.Callable })
		m.WarnfAt(Refused, methods[at].Pos, "%s.%s %s, so %s has no fake: %s", i.Name, refused.Callable,
			refused.Detail, i.Name, refused.Remedy)
		return nil
	}
	if why := in.keyed(&ref, methods); why != "" {
		m.Errorf(Unimplementable, "%s %s, so it has no fake", i.Name, why)
		return nil
	}
	in.record(&ref, methods, witnesses, strings.Join(absent, i.Name))
	return nil
}

// implementable returns the methods of the interface that its fake
// implements, in the order of the member walk, and the witnesses of its type
// parameters. It returns the reason why a fake cannot implement the interface
// or why its test cannot refer to it instead: an interface that its package
// does not export, that the file of the fake cannot import, whose members are
// incomplete, that has a method that its package does not export or no method,
// or whose type parameters have no witnesses that the files of the fake can
// refer to.
func (in *invocation) implementable() ([]*node.Method, []*node.TypeRef, string) {
	i := in.m.Interface
	if i.Visibility != symbol.VisibilityPublic {
		return nil, nil, "is not exported"
	}
	if pkg, placed := in.e.File().Package(); placed == nil {
		if err := in.b.Importable(pkg, &node.TypeRef{Spelling: i.Name, Target: i.ID}); err != nil {
			return nil, nil, fmt.Sprintf("cannot be imported by the file of its fake in the package %s: %v",
				pkg.Package, err)
		}
	}
	set, _ := in.b.MembersOf(i)
	if !set.Complete() {
		return nil, nil, "has incomplete members"
	}
	var methods []*node.Method
	for k := range set.Members {
		method, _ := set.Members[k].Symbol.(*node.Method)
		if method.Visibility != symbol.VisibilityPublic {
			return nil, nil, "has the method " + method.Name + " that its package does not export"
		}
		methods = append(methods, method)
	}
	if len(methods) == 0 {
		return nil, nil, "has no method"
	}
	witnesses := in.b.Witnesses(i.TypeParams)
	for _, out := range []*sdk.Out{in.e.File(), in.e.File(TestTag)} {
		pkg, placed := out.Package()
		unimportable := func(w *node.TypeRef) bool { return in.b.Importable(pkg, w) != nil }
		if placed == nil && slices.ContainsFunc(witnesses, unimportable) {
			witnesses = nil
		}
	}
	if len(i.TypeParams) > 0 && witnesses == nil {
		return nil, nil, "has type parameters without witnesses that the files of its fake can refer to"
	}
	return methods, witnesses, ""
}

// factsOf returns the facts of the method that the reference forms read: its
// classification, its signature, and the type and the fields of each input
// and of each result in the role value.
func (in *invocation) factsOf(method *node.Method) derive.Facts {
	view, _ := in.b.CallableOf(method)
	f := derive.Facts{
		Declaration: in.m.Interface.Name, Callable: view.Name, ID: method.ID,
		Class: shape.OfDeclaration(in.m, method.ID), Signature: view,
	}
	for _, p := range view.Params {
		if p.Role == rules.ParamInput {
			f.Inputs = append(f.Inputs, in.valueOf(p.Name, p.Ref))
		}
	}
	for _, r := range view.Returns {
		if r.Role == rules.ReturnValue {
			f.Results = append(f.Results, in.valueOf(r.Name, r.Ref))
		}
	}
	return f
}

// valueOf returns the value named name of the type ref with its shape and the
// exported fields of the struct that ref refers to. A pointer has no fields,
// because a fake that keys a nil pointer by its field would panic.
func (in *invocation) valueOf(name string, ref *node.TypeRef) derive.Value {
	v := derive.Value{Name: name, Shape: in.b.TypeOf(ref)}
	sym, _ := in.m.Reader().Lookup(ref.Target)
	if st, isStruct := sym.(*node.Struct); isStruct {
		for _, f := range st.Fields {
			if f.Visibility == symbol.VisibilityPublic {
				v.Fields = append(v.Fields, derive.Field{Name: f.Name, Shape: in.b.FieldTypeOf(f)})
			}
		}
	}
	return v
}

// keyed returns the reason why the language cannot key the values of a fake
// of the state map, and the empty string where it can or where the state is a
// list.
func (in *invocation) keyed(ref *derive.Reference, methods []*node.Method) string {
	for k := range ref.Callables {
		c := &ref.Callables[k]
		if c.Effect != spec.EffectLoad {
			continue
		}
		key := keyParam(methods[k], c.Key.Name)
		if why := in.p.keyable(key, in.m.Interface.TypeParams); why != "" {
			return fmt.Sprintf("keys its values by the parameter %s of the type %s, %s", key.Name, key.Type.Spelling,
				why)
		}
	}
	return ""
}

// keyParam returns the parameter named name of the method.
func keyParam(method *node.Method, name string) *node.Param {
	at := slices.IndexFunc(method.Params, func(p *node.Param) bool { return p.Name == name })
	return method.Params[at]
}

// record records the model of the fake of the reference ref, whose callables
// are the methods, into the family of the fakes, and the record of its test
// into the family of the tests. absent is the name of the sentinel of an
// absent key.
func (in *invocation) record(ref *derive.Reference, methods []*node.Method, witnesses []*node.TypeRef,
	absent string,
) {
	i := in.m.Interface
	name := in.e.JoinName(fakeWord, i.Name)
	fake := &emit.Struct{
		Origin: i.ID, Name: name, Visibility: symbol.VisibilityPublic, TypeParams: sdk.MirrorTypeParams(i.TypeParams),
	}
	model := Model{
		Fake: name, New: in.e.JoinName(name, newWord), Subject: i.Name, Form: ref.Form.ID, State: ref.Form.State,
	}
	for _, w := range witnesses {
		model.Witnesses = append(model.Witnesses, rules.EmitRef(w))
	}
	model.Guard = &emit.TypeRef{Spelling: i.Name, Target: i.ID, Args: model.Witnesses}
	for k, method := range methods {
		mirrored := naming.Params(sdk.Mirror(name, method))
		c := in.methodOf(&ref.Callables[k], method, mirrored)
		fake.Methods.Append(in.p.receiver(mirrored, fake.TypeParams, c.Receiver))
		model.Methods = append(model.Methods, c)
		in.typesOf(&model, &ref.Callables[k], method)
	}
	record := []symbol.Symbol{fake, &emit.Function{Origin: i.ID, Name: model.New, Visibility: symbol.VisibilityPublic}}
	if model.State == spec.StateMap {
		model.Absent = absent
		record = append(record, &emit.Variable{Origin: i.ID, Name: absent, Visibility: symbol.VisibilityPublic})
	}
	in.e.File().Record(model, record...)
	emitted := func(name string, args []*emit.TypeRef) *emit.TypeRef {
		return &emit.TypeRef{Spelling: name, Target: i.ID, Emitted: true, Args: args}
	}
	test := &emit.Struct{Origin: i.ID, Name: in.e.JoinName(conformanceWord, name), Visibility: symbol.VisibilityPackage}
	test.Fields.Append(
		&emit.Field{Name: "Run", Type: emitted(suite.NamesOf(in.e, i.Name).Run, nil)},
		&emit.Field{Name: newWord, Type: emitted(model.New, model.Witnesses)},
	)
	in.e.File(TestTag).Record(Test{Test: in.e.JoinName(name, testWord), Fake: name, Subject: i.Name}, test)
}

// Test is the data of the record of the test of one fake: the name of the
// test, of the fake and of its declaration. The record has the field Run,
// which refers to the function of the suite of the declaration that applies
// the suite to a subject, and the field New, which refers to the constructor
// of the fake at the witnesses of the declaration.
type Test struct {
	Test    string
	Fake    string
	Subject string
}

// methodOf returns the method of a fake for the callable c of the method,
// whose mirror on the fake is mirrored. The names of the receiver and of the
// locals avoid the parameters of the mirror. They start with distinct words,
// and a method has one result besides its error at most, so they avoid each
// other too.
func (in *invocation) methodOf(c *derive.Effected, method *node.Method, mirrored *emit.Method) Method {
	params := mirrored.Params.Items()
	taken := make(map[string]bool, len(params))
	var inputs []string
	out := Method{Name: method.Name, Effect: c.Effect}
	errs := c.Facts.Signature.Errors != rules.ErrorsNone
	for k, p := range params {
		taken[p.Name] = true
		switch c.Facts.Signature.Params[k].Role {
		case rules.ParamContext:
			if errs {
				out.Context = p.Name
			}
		case rules.ParamInput:
			inputs = append(inputs, p.Name)
		}
	}
	out.Receiver = naming.Unique(taken, receiverName)
	switch c.Effect {
	case spec.EffectLoad:
		out.Key = c.Key.Name
		out.Present = naming.Unique(taken, presentLocal)
	case spec.EffectStore:
		out.Value = inputs[c.Value]
		out.Key = c.Key.Param
		if c.Key.Kind == derive.ArgField {
			out.Key = in.p.field(out.Value, c.Key.Field)
		}
	case spec.EffectAppend:
		out.Value = inputs[c.Value]
	}
	for k, r := range mirrored.Returns.Items() {
		result := Result{Type: r.Type, Error: c.Facts.Signature.Returns[k].Role == rules.ReturnError}
		if !result.Error {
			result.Local = naming.Unique(taken, resultLocal+strconv.Itoa(k))
		}
		out.Results = append(out.Results, result)
	}
	return out
}

// typesOf sets the type of the keys and the type of the values of the model
// from the callable c of the method: the parameter of the key and the value
// result of a load, and the input of an append. A load returns one value
// before its error, as Go returns an error last. A store stores the type that
// each load returns, so the loads set the types of a map.
func (*invocation) typesOf(model *Model, c *derive.Effected, method *node.Method) {
	switch c.Effect {
	case spec.EffectStore:
	case spec.EffectLoad:
		model.Key = rules.EmitRef(keyParam(method, c.Key.Name).Type)
		model.Value = rules.EmitRef(method.Returns[0].Type)
	case spec.EffectAppend:
		at := slices.IndexFunc(c.Facts.Signature.Params, func(p rules.ParamView) bool {
			return p.Role == rules.ParamInput
		})
		model.Value = rules.EmitRef(method.Params[at].Type)
	}
}
