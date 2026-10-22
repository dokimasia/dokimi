// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/internal/lang"
	"go.dokimi.dev/dokimi/internal/weave"
)

// ID is the name of the plugin.
const ID plugin.ID = "stub"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.2.0"

// Capability is the capability that stub provides. The suite generator
// requires it, because it builds its planted defects from doubles.
const Capability plugin.Capability = "stub"

// StubDirective puts an interface, a struct or a class under stub.
const StubDirective directive.Name = "stub"

// Word is the word of the family of the doubles and of the family of their
// tests, and the default of the option suffix. TestTag is the tag of the
// family of the tests.
const (
	Word    = "stub"
	TestTag = sdk.TagTest
)

// The tree of each target contains the file template of the doubles and the
// file template of their tests.
const (
	stubTemplate = "stub.tmpl"
	testTemplate = "stub.test.tmpl"
)

// The name of the type of the calls of a method joins callWord to the
// subject and the method, and the name of the type of its fixed result joins
// returnWord. The name of a method in a failure joins the subject and the
// method with separator. The name of the record of the tests of a double
// joins checksWord to the name of the double, so a reference to the double
// names the double alone.
const (
	callWord   = "call"
	returnWord = "return"
	checksWord = "checks"
	separator  = "."
)

// The receiver of a method of a double takes receiverName, the call in its body
// recordLocal, the answer answerLocal and the function that the test sets
// funcLocal. An unnamed result takes localPrefix and its position as its
// local. A method of the record of the tests takes argsPrefix, pinsPrefix or
// propsPrefix and its position in the record. Its body assigns the sample of
// the argument at position i to argLocal followed by i, and the sample of the
// result at position i to pinLocal followed by i. It assigns the sample of a
// property to the property of the double named doubleLocal.
const (
	receiverName = "s"
	recordLocal  = "call"
	answerLocal  = "r"
	funcLocal    = "fn"
	localPrefix  = "r"
	argsPrefix   = "args"
	pinsPrefix   = "pins"
	propsPrefix  = "props"
	argLocal     = "a"
	pinLocal     = "want"
	doubleLocal  = "double"
)

// Options are the options of stub. Suffix is the word that the name of each
// double joins to the name of its declaration, and takes [Word] where it is
// empty.
type Options struct {
	Suffix string `opt:"suffix" doc:"the word that the name of each double joins to the name of its declaration"`
}

// Model is the data of the record of one double, and the templates of both
// families read it.
//
// Double is the name of the double. Subject is the name of the declaration
// of the double. Producer is the
// declaration under the directive whose method returns the subject, and is
// empty for a subject under the directive. Value is the subject at the type
// parameters of the double, and nil where no double is a value of the type
// of its subject, as for a Go struct. Delegate is the type of the
// implementation that the option DelegateTo wraps, and nil where the
// template writes the interface of the methods of the double, as for a Go
// struct. Witnesses are the
// witnesses of the type parameters of a generic subject, and nil for a
// subject without type parameters or without witnesses. Guard is the type
// that the double is a value of at the witnesses, which a compile-time check
// asserts, and nil where no such check exists. Generic reports that the
// subject has type parameters.
//
// Methods are the methods of the double, in the order of the member walk of
// the subject. Properties are the fields that the double declares for the
// properties of its subject. Props is the position of the method of the
// record of the tests whose body assigns the sample of each property to a
// double, and -1 where no property has a sample. Fallible reports that a
// method of the double can fail, and Ordered that a method may follow another
// method alone.
//
// The record of a double names the struct of the double. In a language whose
// checks refer to the declarations beside a double through its package, the
// record also declares each type and function that the template writes
// beside the struct, and the template writes them at the struct. The record
// of the checks then has a field of the double and of each such declaration.
// The field has the name of the declaration, and its type refers to it, so
// the layout of eidos qualifies the reference where the checks are in
// another package than the double.
type Model struct {
	Value      *emit.TypeRef
	Delegate   *emit.TypeRef
	Guard      *emit.TypeRef
	Double     string
	Subject    string
	Producer   string
	Witnesses  []*emit.TypeRef
	Methods    []Method
	Properties []Value
	Props      int
	Generic    bool
	Fallible   bool
	Ordered    bool
}

// Checks is the data of the record of the tests of one double, and the file
// template of the tests reads it. It is the model of the double, as a type of
// its own, so each family records the data of its own type.
type Checks Model

// Method is one method of a double.
//
// Name is the name of the method, and Qualified the name of the method in
// each failure that the runtime reports: the subject and the method, joined
// by a period. Field is the member of the double that configures the method.
// Call, Stub and Return are the names of the type of its calls, of the type
// that configures it and of the type of its fixed result. Return is empty
// where the language returns the one value of a method as its result.
// Receiver is the name of the receiver of the methods of the double, which no
// parameter of the method has. Record is the local of the call in the body of
// the method, Answer the local of its answer and Func the local of the
// function that the test sets for the method. No parameter, result or
// receiver has any of the three names. From is the declaration that
// contributes the method through an embed or a base, and empty for a method
// that the subject declares. After is the method that the method may follow
// alone, and empty for a method without the mixin order-after.
//
// Params are the parameters, and Returns the results. Err is the position of
// the error result, and -1 for a method without one. Failure is the field of
// the call that records the error that failed the call, and empty where the
// language returns the error as a result. Named reports that the methods of
// the double declare the names of the results. Async reports that the method
// returns its results asynchronously, and Fails that a fault can fail the
// method. Deprecated reports that the method has the mixin deprecated.
// Sequence is the sequence that the method returns, and nil for a method that
// returns no sequence. Defaulted reports that a result of the method has a
// default result.
//
// Args is the position of the method of the record of the tests whose body
// assigns the samples of the arguments, and Pins the position of the method
// whose body assigns the samples of the results. Each is -1 where no sample
// exists. Pinnable reports that a result has a sample.
//
// TypeParams are the type parameters of the double that a parameter or a
// result of the method mentions, in the order of the double. ArgTypeParams
// are the type parameters that a parameter mentions. ArgWitnesses are the
// witnesses of ArgTypeParams, and nil where the double has no witnesses.
// TypeScript refuses a type parameter that a declaration does not use, so
// each type of a method declares the type parameters that it uses.
//
// Errors are the error values that the helpers of the stub of the method
// arm, one helper each, in the order of the directive fault on the method.
// Retry is the attempt that the retry schedule of the method lets succeed,
// and 0 for a method without one. Partition is the position of the parameter
// whose argument selects the calls of a partition fault, and -1 for a method
// without one. The directive fault on the method declares all three.
type Method struct {
	Sequence      *Sequence
	Errors        []ErrorValue
	TypeParams    []*emit.TypeParam
	ArgTypeParams []*emit.TypeParam
	ArgWitnesses  []*emit.TypeRef
	Name          string
	Qualified     string
	Field         string
	Call          string
	Stub          string
	Return        string
	Receiver      string
	Record        string
	Answer        string
	Func          string
	From          string
	After         string
	Failure       string
	Params        []Param
	Returns       []Result
	Err           int
	Args          int
	Pins          int
	Retry         int
	Partition     int
	Named         bool
	Async         bool
	Fails         bool
	Deprecated    bool
	Pinnable      bool
	Defaulted     bool
}

// Checked returns the count of the results of the method that the test of
// its fixed results compares: each result with a sample, and the error.
func (m *Method) Checked() int {
	n := 0
	for _, r := range m.Returns {
		if r.Pinned || r.Error {
			n++
		}
	}
	return n
}

// Param is one parameter of a method of a double: its type, its name in the
// methods of the double, the field of the call that records it, and whether
// it is variadic or may be absent. Context reports that the rules of the
// language give the parameter the role of a context. A test passes its own
// context as the argument.
type Param struct {
	Type     *emit.TypeRef
	Name     string
	Field    string
	Variadic bool
	Optional bool
	Context  bool
}

// Result is one result of a method of a double: its type, the name that the
// methods of the double declare for it, the field of the call and of the
// fixed result that records it, and the local that the function of the
// method returns it into. The local is the name of a named result, and no
// parameter or receiver has it. Param is the parameter of the setter Returns
// that takes the result, which the receiver of the setter does not have.
// Error reports that the result is the error of the method, and Pinned that
// the tests have a sample of it. Default is the default result, and nil where
// an unconfigured call returns the zero value of the result.
type Result struct {
	Type    *emit.TypeRef
	Default *Default
	Name    string
	Field   string
	Local   string
	Param   string
	Error   bool
	Pinned  bool
}

// Default is the default result of a result whose type has a builder: a new
// value from the builder. Ref is the name of the declaration that the record
// of the double and the record of its tests refer to. In Go it is the
// constructor from the seed, such as NewAccountBuilder, and in TypeScript the
// builder class, such as AccountBuilder. Target is the type that the builder
// builds, whose record declares Ref. Seed is the name of the static method of
// the builder class that constructs a builder from the seed, such as create,
// and empty in Go. Args are the type arguments of the type that the builder
// builds. Optional reports that the result is an optional of that type, such
// as a Go pointer or Account | undefined in TypeScript. A Go double returns
// the address of the value.
type Default struct {
	Target   symbol.Identity
	Ref      string
	Seed     string
	Args     []*emit.TypeRef
	Optional bool
}

// ErrorValue is one error value that a helper of the stub of a method arms
// for the next call of the method: the name of the value in the source, the
// name of the helper, and the value as the name of its declaration, which the
// builtin value of a template writes with the import of its package.
type ErrorValue struct {
	Value  emit.Value
	Name   string
	Helper string
}

// Value is one property of a double: its type, its field, and whether the
// value may be absent.
type Value struct {
	Type     *emit.TypeRef
	Field    string
	Optional bool
}

// Sequence is the sequence that a method returns. Elem is the type of each
// element. Second is the type of the second value of a pair, and nil for a
// sequence of single values. Errors reports that the second value is an
// error.
type Sequence struct {
	Elem   *emit.TypeRef
	Second *emit.TypeRef
	Errors bool
}

// Site is one method of a double with the record of the tests, which a
// partial template of the tests reads to place the bodies that assign the
// samples of the method.
type Site struct {
	Test   *emit.Struct
	Method *Method
}

// The helpers of a plan share one namespace, so the name of each helper of
// stub opens with stub. siteHelper returns the [Site] of a method,
// camelHelper writes a name in camel case, and refHelper returns the
// reference of a record to a declaration that the record refers to, such as
// a declaration of the double or a builder.
const (
	siteHelper  = "stubSite"
	camelHelper = "stubCamel"
	refHelper   = "stubRef"
)

// siteOf returns the site of the method at the position i of methods, with
// the record test of the tests.
func siteOf(test *emit.Struct, methods []Method, i int) Site {
	return Site{Test: test, Method: &methods[i]}
}

// referentOf returns the reference of the record r to the declaration named
// name: the type of the field of r that refers to the declaration. The target
// of a plan writes the name of the field in its own case, such as
// accountBuilder in TypeScript, and the reference keeps the name of the
// declaration. referentOf returns an error where no field of r refers to the
// declaration, so a template that names a declaration that the record does
// not refer to fails.
func referentOf(r *emit.Struct, name string) (*emit.TypeRef, error) {
	for _, f := range r.Fields.Items() {
		if f.Type.Spelling == name {
			return f.Type, nil
		}
	}
	return nil, fmt.Errorf("stub: the record %s refers to no declaration %s", r.Name, name)
}

// referent returns a field of a record whose type refers to the declaration
// named name of the record of origin: a declaration of the double, or the
// constructor of a builder. The reference is an emitted reference, so the
// layout of eidos qualifies it where the file of the record is in another
// package than the declaration.
func referent(origin *symbol.Identity, name string) *emit.Field {
	return &emit.Field{Name: name, Type: &emit.TypeRef{Spelling: name, Target: *origin, Emitted: true}}
}

// part is the part of stub for one language. It declares the target with its
// template tree and its helpers. It also reports each fact of the language
// that the facts of eidos leave to it, and returns the names of the members
// that the templates write.
type part interface {
	lang.Part
	// structural reports whether the language types a value by its members
	// alone. A double of a struct of such a language is a value of the
	// struct, and a double declares the fields of its subject.
	structural() bool
	// referable reports whether the file of a double can refer to a subject
	// of the visibility v.
	referable(v symbol.Visibility) bool
	// receiver returns the method m of a double with the receiver named name
	// that the language writes. tps are the type parameters of the double.
	receiver(m *emit.Method, tps []*emit.TypeParam, name string) *emit.Method
	// constraint reports whether the interface id only constrains type
	// parameters, so that no value has its type.
	constraint(m sdk.Matcher, id *symbol.Identity) bool
	// delegate returns the type of the implementation that the option
	// DelegateTo of a double of the subject value wraps. class reports that
	// the subject is a struct or a class. delegate returns nil where the
	// template writes the interface of the methods of the double.
	delegate(value *emit.TypeRef, class bool) *emit.TypeRef
	// names returns the names of the method of a double named method, with
	// the parameters named params and the returns rs. taken are the names
	// of the parameters and the receiver of the methods of the double.
	names(method string, params []string, rs []rules.ReturnView, taken map[string]bool) Names
	// declarations returns the declarations that the templates of the
	// language write beside the struct of the double with the model m,
	// where the checks of the double refer to them through the package of
	// the double. It returns nil where the checks import each of them from
	// the module of the double.
	declarations(double *emit.Struct, m *Model) []symbol.Symbol
	// fault returns the name of the helper of the stub of a method that
	// arms the error value with the base name base for the next call, such
	// as NotFound for the error value ErrNotFound.
	fault(base string) string
	// construction returns the declaration that a default result refers to,
	// and the method of that declaration that constructs a builder from the
	// seed, for the builder named name whose constructor from the seed is
	// named seed. The method is empty where the declaration is the
	// constructor itself.
	construction(name, seed string) (ref, method string)
}

// Names are the names that the part of a language gives one method of a
// double and its call.
//
// Field is the member of the double that configures the method. Args are the
// fields of the call that record the arguments, one for each parameter.
// Results are the fields of the call that record the results, at the
// positions of the returns, and Grouped reports that the method returns its
// results as a value of a type of its own. Setters are the parameters of the
// setter Returns, at the positions of the returns. Named reports that the
// methods of the double declare the names of the results. Failure is the
// field of the call that records the error that failed the call, and empty
// where the language returns the error as a result.
type Names struct {
	Field   string
	Failure string
	Args    []string
	Results []string
	Setters []string
	Grouped bool
	Named   bool
}

// parts are the parts of stub, one for each language that it serves.
var parts = []part{golangPart{}, typescriptPart{}}

// schema is the directive of stub.
var schema = directive.Schema{
	Name: StubDirective,
	Doc: "generates a recording double of the declaration, which records each call, returns the results that the " +
		"test configures and injects faults, and the test of the double",
}

// generator is the state of one instance of stub: its options.
//
// # Concurrency
//
// The composition sets the options before the first phase call, and the
// handlers read them afterwards, so the workers of a phase call read them
// concurrently.
type generator struct {
	options *Options
}

// subject is a declaration that stub writes a double of. value reports that
// a double is a value of the type of the declaration, and class that the
// declaration is a struct, whose constructors a double leaves out. producer
// is the name of the declaration whose method returns a handle.
type subject struct {
	origin     symbol.Identity
	decl       symbol.Symbol
	name       string
	producer   string
	params     []*node.TypeParam
	visibility symbol.Visibility
	value      bool
	class      bool
}

// members are the methods and the fields of the member walk of a subject
// that its double implements, in the order of the walk, and the declaration
// that contributes each method through an embed or a base.
type members struct {
	methods []*node.Method
	from    []string
	fields  []*node.Field
}

// invocation is the state of one double. It contains the match of the
// subject, the part of the language of the subject, the emitter, the subject,
// the rules of the language and the suffix of the names of the doubles.
type invocation struct {
	m      sdk.SubjectMatch
	p      part
	e      *sdk.Emitter
	sub    *subject
	b      rules.Bound
	suffix string
}

// New returns stub. stub records a [Model] into the family [Word] for each
// interface, struct or class under [StubDirective], and for each interface and
// struct on which the suite generator stamps the fact [weave.KeyDouble]. A
// declaration under both directives gets one double. stub runs after the
// capability [weave.Suite], so in a composition without the suite generator
// the gate on the fact does not select a declaration. The template of the
// family writes the double in the conventions of the language of the plan.
// stub records the [Checks] of each double into the family [Word] under the
// tag [TestTag]. The template of that family writes a test of a double with a
// method, and a test of the stub of each method. Each test has a level for
// each function that it checks. The levels run the behaviour suites of the
// runtime, and check the results that need a value of the signature of a
// method. The test of a double without a method builds the double under each
// option, and checks that a reset keeps the sample of each property. The test
// uses the double through the exported API of its package, so a Go file of
// the family is in the external test package of its directory.
//
// The methods of a double are the methods of the member walk of its subject,
// without the members at type level and the constructors of a class. A
// double declares the fields of its subject where the part of the language
// types a value by its members. A method of a double of an interface returns
// an interface of the package of the subject without the directive, called a
// handle. Such a handle gets a double in the file of the double of the first
// subject of the package whose method returns it. An out directive on that
// subject moves the double of the handle with the double of the subject.
// stub serves Go and TypeScript, the languages of its parts, and leaves a
// declaration of another language alone.
//
// The stub of a method writes the helpers of the faults that the directive
// fault on the method declares, which the annotator fault stamps. So stub
// requires the capability of fault.
//
// stub reports [Constraint] at an interface that only constrains type
// parameters, [IncompleteMembers] at a declaration with a gap in its member
// walk, [Unimplementable] at a declaration with a member that a double cannot
// implement or that the file of its double cannot refer to, [Collision] at a
// declaration whose double would have two members of one name, and [Empty] at
// a declaration without a member. stub records nothing for a declaration that
// it reports. stub also reports [Incomparable] at a declaration with a
// partition fault that cannot compare the arguments of its parameter, and
// records the double without the partition helpers. Each call returns a new
// plugin.
func New() plugin.Generator {
	g := &generator{options: &Options{Suffix: Word}}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Provides(Capability).
		Requires(fault.Capability).
		After(weave.Suite).
		Options(g.options).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: stubTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: Word, Template: testTemplate, External: true})
	lang.Register(b, parts)
	asked := sdk.HasKey(meta.Named[bool](weave.KeyDouble))
	return b.Handle(
		sdk.Directive(schema, sdk.Where(sdk.Not(asked), g.rules()...)),
		sdk.Where(asked, g.rules()...),
	).Generator()
}

// rules returns the rules of the doubles of an interface and of a struct.
func (g *generator) rules() []sdk.Rule {
	return []sdk.Rule{
		sdk.OnInterface(lang.Served(parts, g.generateInterface)), sdk.OnStruct(lang.Served(parts, g.generateStruct)),
	}
}

// generateInterface records the double of an interface under the directive,
// and the double of each handle that the interface produces. It reports
// [Constraint] at an interface that only constrains type parameters.
func (g *generator) generateInterface(m *sdk.InterfaceMatch, p part, e *sdk.Emitter) error {
	i := m.Interface
	if p.constraint(m, &i.ID) {
		m.Errorf(Constraint, "%s has type-set elements and only constrains type parameters, so it has no double",
			i.Name)
		return nil
	}
	in := g.invocation(m, p, e, &subject{
		origin: i.ID, decl: i, name: i.Name, params: i.TypeParams, visibility: i.Visibility, value: true,
	})
	for _, h := range in.handles(in.write().methods) {
		g.invocation(m, p, e, h).write()
	}
	return nil
}

// generateStruct records the double of a struct or a class under the
// directive. The double is a value of the struct where the language types a
// value by its members.
func (g *generator) generateStruct(m *sdk.StructMatch, p part, e *sdk.Emitter) error {
	s := m.Struct
	g.invocation(m, p, e, &subject{
		origin: s.ID, decl: s, name: s.Name, params: s.TypeParams, visibility: s.Visibility,
		value: p.structural(), class: true,
	}).write()
	return nil
}

// invocation returns the invocation of the double of sub.
func (g *generator) invocation(m sdk.SubjectMatch, p part, e *sdk.Emitter, sub *subject) *invocation {
	return &invocation{m: m, p: p, e: e, sub: sub, b: m.Rules(), suffix: cmp.Or(g.options.Suffix, Word)}
}

// write records the model of the double of the subject into the family of
// the doubles and into the family of the tests, and returns the members of
// the double. Where write reports a finding, it does not record the model
// and returns no members. A handle without a member has no double and no
// finding, because a method may return an interface without methods.
func (in *invocation) write() members {
	sub := in.sub
	if !in.p.referable(sub.visibility) {
		in.m.Errorf(Unimplementable, "%s is not exported, so the file of its double cannot refer to it", sub.name)
		return members{}
	}
	walked, complete := in.walk()
	if !complete {
		return members{}
	}
	if len(walked.methods) == 0 && len(walked.fields) == 0 {
		if sub.producer == "" {
			in.m.Errorf(Empty, "%s has no member that a double can implement, so it has no double", sub.name)
		}
		return members{}
	}
	host := in.e.JoinName(in.suffix, sub.name)
	double := &emit.Struct{
		Origin:     sub.origin,
		Name:       host,
		Visibility: sub.visibility,
		TypeParams: sdk.MirrorTypeParams(sub.params),
	}
	model := in.modelOf(double)
	test := &emit.Struct{
		Origin: sub.origin, Name: in.e.JoinName(checksWord, host), Visibility: sub.visibility,
		TypeParams: sdk.MirrorTypeParams(sub.params),
	}
	for i, method := range walked.methods {
		mirrored := naming.Params(sdk.Mirror(host, method))
		c := in.methodOf(method, mirrored, walked.from[i], test)
		typeParamsOf(&c, double.TypeParams, model.Witnesses)
		double.Methods.Append(in.p.receiver(mirrored, double.TypeParams, c.Receiver))
		model.Methods = append(model.Methods, c)
		model.Fallible = model.Fallible || c.Fails
		model.Ordered = model.Ordered || c.After != ""
	}
	props := make([]emit.Value, 0, len(walked.fields))
	targets := make([]string, 0, len(walked.fields))
	for _, f := range walked.fields {
		v := Value{Type: rules.EmitRef(f.Type), Field: f.Name, Optional: f.Optional}
		model.Properties = append(model.Properties, v)
		props = append(props, in.sample(&f.ID, f.Type, v.Type, f.Name))
		targets = append(targets, doubleLocal+"."+f.Name)
	}
	model.Props = samples(test, &sub.origin, propsPrefix, targets, props)
	if name, collides := collision(&model, double); collides {
		in.m.Errorf(Collision, "the double of %s would have two members named %s, so it has no double", sub.name, name)
		return members{}
	}
	declared := in.p.declarations(double, &model)
	record := append(make([]symbol.Symbol, 0, 1+len(declared)), double)
	record = append(record, declared...)
	calls := constructions(model.Methods)
	for _, c := range calls {
		double.Fields.Append(referent(&c.Target, c.Ref))
	}
	in.e.File().Record(model, record...)
	if declared != nil {
		for _, d := range record {
			test.Fields.Append(referent(&sub.origin, emit.DeclaredName(d)))
		}
	}
	for _, c := range calls {
		test.Fields.Append(referent(&c.Target, c.Ref))
	}
	in.e.File(TestTag).Record(Checks(model), test)
	return walked
}

// constructions returns the default results of the methods with each builder
// once, in the order of the methods and of their results. The name of the
// constructor identifies the builder, because a double takes the builders of
// the package of its subject alone.
func constructions(methods []Method) []*Default {
	var out []*Default
	for i := range methods {
		for _, r := range methods[i].Returns {
			if r.Default != nil && !slices.ContainsFunc(out, func(d *Default) bool { return d.Ref == r.Default.Ref }) {
				out = append(out, r.Default)
			}
		}
	}
	return out
}

// modelOf returns the model of the double of the subject without its
// methods and properties.
func (in *invocation) modelOf(double *emit.Struct) Model {
	sub := in.sub
	model := Model{Double: double.Name, Subject: sub.name, Producer: sub.producer, Generic: len(sub.params) > 0}
	value := &emit.TypeRef{Spelling: sub.name, Target: sub.origin, Args: sdk.TypeArgs(double.TypeParams)}
	model.Delegate = in.p.delegate(value, sub.class)
	for _, w := range in.b.Witnesses(sub.params) {
		model.Witnesses = append(model.Witnesses, rules.EmitRef(w))
	}
	if !sub.value {
		return model
	}
	model.Value = value
	if !model.Generic || model.Witnesses != nil {
		model.Guard = &emit.TypeRef{Spelling: sub.name, Target: sub.origin, Args: model.Witnesses}
	}
	return model
}

// walk returns the members of the member walk of the subject that its double
// implements. It leaves out a member at type level, a constructor of a
// struct, and a field where the language does not type a value by its
// members. It reports [IncompleteMembers] where the walk has a gap, and
// [Unimplementable] at a member that a double cannot implement, and returns
// false after either report.
func (in *invocation) walk() (members, bool) {
	sub := in.sub
	set, _ := in.b.MembersOf(sub.decl)
	if !set.Complete() {
		in.m.Errorf(IncompleteMembers, "the members of %s are incomplete, so it has no double", sub.name)
		return members{}, false
	}
	var out members
	for i := range set.Members {
		var reason string
		switch s := set.Members[i].Symbol.(type) {
		case *node.Method:
			if s.Level == symbol.LevelType || s.Constructs && sub.class {
				continue
			}
			reason = unimplemented(s)
			out.methods = append(out.methods, s)
			out.from = append(out.from, contributor(&set.Members[i]))
		case *node.Field:
			if !in.p.structural() || s.Level == symbol.LevelType {
				continue
			}
			switch {
			case hidden(s.Visibility):
				reason = "the member " + s.Name + " that is not public"
			case s.Computed:
				reason = "the property " + s.Name + " with a computed name"
			case s.Type == nil:
				reason = "the property " + s.Name + " without a declared type"
			}
			out.fields = append(out.fields, s)
		}
		if reason != "" {
			in.m.Errorf(Unimplementable, "%s has %s, which a double cannot implement, so it has no double",
				sub.name, reason)
			return members{}, false
		}
	}
	return out, true
}

// contributor returns the name of the declaration that contributes the
// member through an embed or a base, and the empty string for a member that
// the subject declares.
func contributor(member *rules.Member) string {
	if len(member.Through) == 0 {
		return ""
	}
	return member.Owner.Name
}

// unimplemented returns the description of the method s where a double
// cannot implement it, and an empty string where it can. The member walk
// reports a gap at an overloaded method, so walk never passes an overload to
// unimplemented. A double names the members that configure a method after
// the name of the method. A computed name such as [Symbol.iterator] is not
// an identifier, so a method with a computed name has no implementation. A
// double writes the type of each field of a call and of each result, so a
// parameter or a result without a declared type has no implementation
// either.
func unimplemented(s *node.Method) string {
	switch {
	case s.Indexer:
		return "an index signature"
	case s.Constructs:
		return "a construct signature"
	case s.Accessor != symbol.AccessorNone:
		return "the accessor " + s.Name
	case s.Computed:
		return "the method " + s.Name + " with a computed name"
	case len(s.TypeParams) > 0:
		return "the method " + s.Name + " with type parameters"
	case hidden(s.Visibility):
		return "the member " + s.Name + " that is not public"
	}
	for _, p := range s.Params {
		if p.Type == nil {
			return "the parameter " + p.Name + " of " + s.Name + " without a declared type"
		}
	}
	for _, r := range s.Returns {
		if r.Type == nil {
			return "the method " + s.Name + " without a declared result type"
		}
	}
	return ""
}

// hidden reports whether a member of the visibility v is hidden from a type
// that does not extend its declaration, so that a double cannot implement
// it.
func hidden(v symbol.Visibility) bool {
	return v == symbol.VisibilityPrivate || v == symbol.VisibilityProtected
}

// methodOf returns the method of the double that mirrored mirrors in the
// double of the method of the source. from is the declaration that
// contributes the method. methodOf appends the methods whose bodies assign
// the samples of the arguments and of the results to the record test of the
// tests. The part of the language returns the names of the field of the
// method, of the fields of its call and of the fields of its results. A fault
// can fail the method where the error model of the language has a failure
// channel.
func (in *invocation) methodOf(method *node.Method, mirrored *emit.Method, from string, test *emit.Struct) Method {
	view, _ := in.b.CallableOf(method)
	joined := in.e.JoinName(method.Name, in.sub.name)
	c := Method{
		Name:       method.Name,
		Qualified:  in.sub.name + separator + method.Name,
		Call:       in.e.JoinName(callWord, joined),
		Stub:       in.e.JoinName(in.suffix, joined),
		From:       from,
		Err:        -1,
		Args:       -1,
		Pins:       -1,
		Async:      view.Async,
		Fails:      view.Errors != rules.ErrorsNone,
		Deprecated: in.mixed(&method.ID, shape.Deprecated),
		After:      in.after(&method.ID),
		Sequence:   in.sequenceOf(view.Returns),
	}
	params := mirrored.Params.Items()
	taken := make(map[string]bool, len(params)+1)
	names := make([]string, 0, len(params))
	for _, p := range params {
		taken[p.Name] = true
		names = append(names, p.Name)
	}
	c.Receiver = naming.Unique(taken, receiverName)
	taken[c.Receiver] = true
	named := in.p.names(method.Name, names, view.Returns, taken)
	c.Field, c.Failure, c.Named = named.Field, named.Failure, named.Named
	for i, p := range params {
		variadic := p.Variadic != symbol.VariadicNone
		c.Params = append(c.Params, Param{
			Name: p.Name, Field: named.Args[i], Type: p.Type,
			Variadic: variadic, Optional: view.Params[i].Optional && !variadic,
			Context: view.Params[i].Role == rules.ParamContext,
		})
	}
	if named.Grouped {
		c.Return = in.e.JoinName(returnWord, joined)
	}
	setter := map[string]bool{receiverName: true}
	for i, r := range mirrored.Returns.Items() {
		result := Result{Field: named.Results[i], Type: r.Type, Error: view.Returns[i].Role == rules.ReturnError}
		if c.Named {
			result.Name = r.Name
		}
		result.Local = naming.Unique(taken, cmp.Or(result.Name, localPrefix+strconv.Itoa(i)))
		taken[result.Local] = true
		result.Param = naming.Unique(setter, named.Setters[i])
		setter[result.Param] = true
		if result.Error {
			c.Err = i
		}
		result.Default = in.defaultOf(r.Type)
		c.Defaulted = c.Defaulted || result.Default != nil
		r.Name = result.Name
		c.Returns = append(c.Returns, result)
	}
	c.Record = naming.Unique(taken, recordLocal)
	c.Answer = naming.Unique(taken, answerLocal)
	c.Func = naming.Unique(taken, funcLocal)
	c.Args = samples(test, &method.ID, argsPrefix, locals(argLocal, len(c.Params)), in.argSamples(method, &c))
	c.Pins = samples(test, &method.ID, pinsPrefix, locals(pinLocal, len(c.Returns)), in.resultSamples(method, &c))
	values, _ := sdk.FactOf(in.m, method.ID, meta.Named[[]symbol.Identity](fault.KeyErrors))
	for _, id := range values {
		base, named := in.b.SentinelBase(id.Name)
		if !named {
			base = id.Name
		}
		c.Errors = append(c.Errors, ErrorValue{Value: emit.Name(id), Name: id.Name, Helper: in.p.fault(base)})
	}
	retry, _ := sdk.FactOf(in.m, method.ID, meta.Named[int64](fault.KeyRetry))
	c.Retry = int(retry)
	c.Partition = in.partition(method, &c)
	return c
}

// defaultOf returns the default result of a result of the type t, and nil
// where the type that a builder builds for t has no builder. A builder builds
// t itself, or the type of an optional t. The facts of the builder name the
// builder of a type. The templates find a builder by the name of its
// constructor, which is unique in one package. So a type of another package
// than the subject has no default result.
func (in *invocation) defaultOf(t *emit.TypeRef) *Default {
	built, optional := t, t.Form == symbol.FormOptional
	if optional {
		built = t.Elems[0]
	}
	if built.Target.Package != in.sub.origin.Package {
		return nil
	}
	name, has := sdk.FactOf(in.m, built.Target, meta.Named[string](builder.KeyName))
	if !has {
		return nil
	}
	seed, _ := sdk.FactOf(in.m, built.Target, meta.Named[string](builder.KeySeed))
	ref, method := in.p.construction(name, seed)
	return &Default{Ref: ref, Seed: method, Args: built.Args, Target: built.Target, Optional: optional}
}

// partition returns the position of the parameter of the method c whose
// argument selects the calls of a partition fault, which the annotator fault
// stamps, and -1 for a method without one. A partition fault compares each
// argument of the parameter with a key. So partition reports [Incomparable]
// at a variadic parameter and at a parameter whose type mentions a type
// parameter of the double, and returns -1 for either.
func (in *invocation) partition(method *node.Method, c *Method) int {
	id, partitioned := sdk.FactOf(in.m, method.ID, meta.Named[symbol.Identity](fault.KeyPartition))
	if !partitioned {
		return -1
	}
	at := slices.IndexFunc(method.Params, func(p *node.Param) bool { return p.ID == id })
	p := &c.Params[at]
	var reason string
	switch {
	case p.Variadic:
		reason = "is variadic"
	case in.mentions(p.Type):
		reason = "has a type that mentions a type parameter"
	default:
		return at
	}
	in.m.Errorf(Incomparable, "the partition parameter %s of %s %s, so a partition fault cannot compare its "+
		"arguments with a key", p.Name, method.Name, reason)
	return -1
}

// mixed reports whether the shape catalog stamps the mixin n on the method
// id.
func (in *invocation) mixed(id *symbol.Identity, n shape.Mixin) bool {
	spec, _ := shape.SpecOf(n)
	mixed, _ := sdk.FactOf(in.m, *id, meta.Named[bool](spec.Key))
	return mixed
}

// after returns the name of the method that the method id may follow alone,
// which the param fn of its mixin order-after names, and the empty string for
// a method without the mixin. The spec of order-after lists fn as its first
// param.
func (in *invocation) after(id *symbol.Identity) string {
	spec, _ := shape.SpecOf(shape.OrderAfter)
	fn, _ := sdk.FactOf(in.m, *id, meta.Named[symbol.Identity](spec.Params[0].Fact))
	return fn.Name
}

// sequenceOf returns the sequence that a method with the returns rs returns,
// and nil where the method does not return one sequence alone. The elements
// of the sequence are the type arguments of the stream, also where an alias
// names the stream. A stream without type arguments does not compile, and
// has no sequence.
func (in *invocation) sequenceOf(rs []rules.ReturnView) *Sequence {
	if len(rs) != 1 || rs[0].Role != rules.ReturnStream {
		return nil
	}
	args := rules.Unaliased(rs[0].Ref, in.b.View()).Args
	if len(args) == 0 {
		return nil
	}
	seq := &Sequence{Elem: rules.EmitRef(args[0])}
	if len(args) > 1 {
		seq.Second = rules.EmitRef(args[1])
		seq.Errors = in.b.TypeOf(args[1]).Ref == rules.WellKnownError
	}
	return seq
}

// argSamples returns the sample of each parameter of the method c, and nil
// for a parameter without one. A variadic parameter takes no sample, because
// a sample of its element is no value of the list that the method receives,
// and neither does a parameter whose type mentions a type parameter of the
// subject.
func (in *invocation) argSamples(method *node.Method, c *Method) []emit.Value {
	out := make([]emit.Value, len(c.Params))
	for i, p := range method.Params {
		if c.Params[i].Variadic {
			continue
		}
		out[i] = in.sample(&p.ID, p.Type, c.Params[i].Type, method.Name)
	}
	return out
}

// resultSamples returns the sample of each result of the method c, and nil
// for a result without one. It marks each result with a sample as pinned,
// and the method as pinnable where a result has one.
func (in *invocation) resultSamples(method *node.Method, c *Method) []emit.Value {
	out := make([]emit.Value, len(c.Returns))
	for i, r := range method.Returns {
		out[i] = in.sample(&r.ID, r.Type, c.Returns[i].Type, method.Name)
		c.Returns[i].Pinned = !out[i].IsZero()
		c.Pinnable = c.Pinnable || c.Returns[i].Pinned
	}
	return out
}

// sample returns the sample of the type ref of the declaration id, whose
// emitted type is emitted, and the zero value for a type without a sample.
// hint is the text of a sample of text. A type that mentions a type
// parameter of the double has no sample, because the checks of a generic
// double run in a function that takes the type parameters. A refused sample
// has the zero value.
func (in *invocation) sample(id *symbol.Identity, ref *node.TypeRef, emitted *emit.TypeRef, hint string) emit.Value {
	if in.mentions(emitted) {
		return emit.Value{}
	}
	s, _ := in.b.SamplesOf(*id, ref, hint)
	return s.Value
}

// mentions reports whether the emitted type t mentions a type parameter of
// the double.
func (in *invocation) mentions(t *emit.TypeRef) bool {
	return len(rules.Mentioned(sdk.MirrorTypeParams(in.sub.params), t)) > 0
}

// typeParamsOf sets the type parameters of the method c from the type
// parameters tps of the double, as [Method] states them. witnesses are the
// witnesses of tps in their order, and nil where tps has none.
func typeParamsOf(c *Method, tps []*emit.TypeParam, witnesses []*emit.TypeRef) {
	args := make([]*emit.TypeRef, 0, len(c.Params))
	for _, p := range c.Params {
		args = append(args, p.Type)
	}
	signature := slices.Clone(args)
	for _, r := range c.Returns {
		signature = append(signature, r.Type)
	}
	c.TypeParams = mentioned(tps, signature)
	c.ArgTypeParams = mentioned(tps, args)
	if witnesses == nil {
		return
	}
	for _, tp := range c.ArgTypeParams {
		c.ArgWitnesses = append(c.ArgWitnesses, witnesses[slices.Index(tps, tp)])
	}
}

// mentioned returns the type parameters of tps that one of the types
// mentions, in the order of tps.
func mentioned(tps []*emit.TypeParam, types []*emit.TypeRef) []*emit.TypeParam {
	var out []*emit.TypeParam
	for i := range tps {
		one := tps[i : i+1]
		if slices.ContainsFunc(types, func(t *emit.TypeRef) bool { return len(rules.Mentioned(one, t)) > 0 }) {
			out = append(out, tps[i])
		}
	}
	return out
}

// samples appends to the record test a method whose body assigns each value
// of values to the name at the same position of names, and returns the
// position of the method. The name of the method is prefix followed by its
// position. origin is the declaration of the source. samples appends nothing
// and returns -1 where each value is zero.
func samples(test *emit.Struct, origin *symbol.Identity, prefix string, names []string, values []emit.Value) int {
	var body []emit.Stmt
	for i := range values {
		if values[i].IsZero() {
			continue
		}
		body = append(body, emit.Stmt{
			Kind: emit.StmtAssign, Names: []string{names[i]}, Value: emit.ValueExpr(values[i]),
		})
	}
	if len(body) == 0 {
		return -1
	}
	test.Methods.Append(&emit.Method{
		Origin: *origin, Name: prefix + strconv.Itoa(test.Methods.Len()), Visibility: symbol.VisibilityPackage,
		Body: emit.Body{Stmts: body},
	})
	return test.Methods.Len() - 1
}

// locals returns the n locals that join local and a position, such as a0 and
// a1 for the local a.
func locals(local string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = local + strconv.Itoa(i)
	}
	return out
}

// collision returns the first name that two members of one type of the
// double share, and reports false where each member has a name of its own.
// The members of the double are its methods, the fields that configure its
// methods and the fields of its properties. The members of a call are the
// fields of the parameters, of the results and of the failure. A call
// without a field of the failure adds the empty name, which no other field
// has. The helpers of the error values of a method are members of the stub of
// the method.
func collision(model *Model, double *emit.Struct) (string, bool) {
	names := make([]string, 0, double.Methods.Len()+len(model.Methods)+len(model.Properties))
	for _, m := range double.Methods.Items() {
		names = append(names, m.Name)
	}
	for i := range model.Methods {
		names = append(names, model.Methods[i].Field)
	}
	for _, v := range model.Properties {
		names = append(names, v.Field)
	}
	if name, twice := repeated(names); twice {
		return name, true
	}
	for i := range model.Methods {
		c := &model.Methods[i]
		fields := make([]string, 0, len(c.Params)+len(c.Returns)+1)
		for _, p := range c.Params {
			fields = append(fields, p.Field)
		}
		for _, r := range c.Returns {
			fields = append(fields, r.Field)
		}
		if name, twice := repeated(append(fields, c.Failure)); twice {
			return c.Call + separator + name, true
		}
		helpers := make([]string, 0, len(c.Errors))
		for j := range c.Errors {
			helpers = append(helpers, c.Errors[j].Helper)
		}
		if name, twice := repeated(helpers); twice {
			return c.Stub + separator + name, true
		}
	}
	return "", false
}

// repeated returns the first name of names that an earlier name repeats, and
// reports false where each name differs.
func repeated(names []string) (string, bool) {
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		if seen[name] {
			return name, true
		}
		seen[name] = true
	}
	return "", false
}

// handles returns the subjects of the handles of the methods of the subject:
// each interface of the package of the subject that a method returns. A
// handle has no directive of its own, and the subject is its first producer.
// So two subjects that return one handle write one double of it.
func (in *invocation) handles(methods []*node.Method) []*subject {
	var out []*subject
	seen := map[symbol.Identity]bool{}
	for _, m := range methods {
		for _, r := range m.Returns {
			decl, _ := in.b.View().Lookup(r.Type.Target)
			i, isInterface := decl.(*node.Interface)
			if !isInterface || seen[i.ID] {
				continue
			}
			seen[i.ID] = true
			if i.ID.Package != in.sub.origin.Package || in.annotated(&i.ID) || in.producer(&i.ID) != in.sub.origin {
				continue
			}
			out = append(out, &subject{
				origin: i.ID, decl: i, name: i.Name, producer: in.sub.name, params: i.TypeParams,
				visibility: i.Visibility, value: true,
			})
		}
	}
	return out
}

// annotated reports whether the declaration id has the directive of stub,
// in either spelling.
func (in *invocation) annotated(id *symbol.Identity) bool {
	return slices.ContainsFunc(in.carriers(), func(d node.Declaration) bool { return d.Identity() == *id })
}

// producer returns the first producer of the handle id in identity order. A
// producer is an interface under the directive in the package of the handle,
// whose member walk has a method that returns the handle. The subject that
// asks is one of them.
func (in *invocation) producer(id *symbol.Identity) symbol.Identity {
	var producers []symbol.Identity
	for _, d := range in.carriers() {
		if i, isInterface := d.(*node.Interface); isInterface && i.ID.Package == id.Package && in.returns(i, id) {
			producers = append(producers, i.ID)
		}
	}
	return slices.MinFunc(producers, symbol.Identity.Compare)
}

// carriers returns the declarations that have the directive of stub,
// written bare or with the prefix of the plugin. The store lists the
// declarations of a directive by the name as written.
func (in *invocation) carriers() []node.Declaration {
	var out []node.Declaration
	for _, name := range []directive.Name{StubDirective, canonical} {
		for s := range in.b.View().Reader().ByDirective(name) {
			if d, declares := s.(node.Declaration); declares {
				out = append(out, d)
			}
		}
	}
	return out
}

// returns reports whether a method of the member walk of the interface i
// returns the declaration id.
func (in *invocation) returns(i *node.Interface, id *symbol.Identity) bool {
	set, _ := in.b.MembersOf(i)
	return slices.ContainsFunc(set.Members, func(member rules.Member) bool {
		m, isMethod := member.Symbol.(*node.Method)
		return isMethod && slices.ContainsFunc(m.Returns, func(r *node.Return) bool {
			return r.Type != nil && r.Type.Target == *id
		})
	})
}

// canonical is the canonical name of the directive of stub, under which the
// reader of a view lists the declarations under it.
var canonical = directive.Schema{Plugin: string(ID), Name: StubDirective}.Canonical()
