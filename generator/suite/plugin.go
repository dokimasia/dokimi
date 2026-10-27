// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"errors"
	"fmt"
	"slices"
	"strings"

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
// subjects. Its param CasesParam contains the number of valid cases of each
// property run of the application, and its param OptionsParam the name of the
// function that returns the options of the application. CheckDirective binds
// a function of a consumer to a callable of a suite as a check. Its
// positional param CallableParam contains the callable, and its param
// DefectParam contains the name of the function of the planted defect. Its
// params ReferenceParam, SecondParam, RestartParam and CrashParam contain the
// names of the functions of the doors of the planted defect. Its param
// CountedParam marks a check that counts the allocations of the whole process,
// and its param AroundParam contains the name of the function that runs each
// run of the check inside an environment, such as a testing/synctest bubble.
// ExcuseDirective excuses a check of the suite from the conformance test of a
// subject. Its positional param IdentityParam contains the identity of the
// check, and its param ReasonParam contains the reason.
const (
	SuiteDirective   directive.Name     = "suite"
	SubjectDirective directive.Name     = "subject"
	CheckDirective   directive.Name     = "check"
	ExcuseDirective  directive.Name     = "excuse"
	DeclarationParam directive.ParamKey = "declaration"
	NewParam         directive.ParamKey = "new"
	CasesParam       directive.ParamKey = "cases"
	OptionsParam     directive.ParamKey = "options"
	CallableParam    directive.ParamKey = "callable"
	DefectParam      directive.ParamKey = "defect"
	ReferenceParam   directive.ParamKey = "reference"
	SecondParam      directive.ParamKey = "second"
	RestartParam     directive.ParamKey = "restart"
	CrashParam       directive.ParamKey = "crash"
	CountedParam     directive.ParamKey = "counted"
	AroundParam      directive.ParamKey = "around"
	IdentityParam    directive.ParamKey = "identity"
	ReasonParam      directive.ParamKey = "reason"
)

// KeyDeclaration, KeyCallable, KeySegment, KeyStatement, KeyDefect,
// KeyProperty and KeyCounted are the keys of the facts that the annotate half
// stamps on the function of each check that it binds. The facts contain the
// identity of the declaration under [SuiteDirective], the name of the callable
// such as Store.Get, the segment and the statement of the check, the identity
// of the function of its planted defect, whether the check is a property
// check, and whether it counts allocations. KeyOptions is the key of the fact
// that contains the identity of the function of the options of the property
// run of the check, and KeyAround the key of the fact that contains the
// identity of the function of its environment. KeyReference, KeySecond,
// KeyRestart and KeyCrash are the keys of the facts that contain the identity
// of the function of each door of the planted defect.
const (
	KeyDeclaration meta.KeyName = "suite.declaration"
	KeyCallable    meta.KeyName = "suite.callable"
	KeySegment     meta.KeyName = "suite.segment"
	KeyStatement   meta.KeyName = "suite.statement"
	KeyDefect      meta.KeyName = "suite.defect"
	KeyProperty    meta.KeyName = "suite.property"
	KeyCounted     meta.KeyName = "suite.counted"
	KeyOptions     meta.KeyName = "suite.options"
	KeyAround      meta.KeyName = "suite.around"
	KeyReference   meta.KeyName = "suite.reference"
	KeySecond      meta.KeyName = "suite.second"
	KeyRestart     meta.KeyName = "suite.restart"
	KeyCrash       meta.KeyName = "suite.crash"
)

// door is one door of the planted defect of a check. Its param of the
// directive check contains the name of the function of the door, and its key
// is the key of the fact of that function. field is the field of
// conform.Doors that the check of the suite sets. takes is the number of the
// parameters of the function, and needs contains the words for the
// parameters after the handle in the finding of another signature. A crash
// takes the call that it crashes as its third parameter. doc is the doc of
// the param.
type door struct {
	param directive.ParamKey
	key   meta.KeyName
	field string
	needs string
	doc   string
	takes int
}

// doors are the doors of a planted defect in the order of the fields of
// conform.Doors.
var doors = [...]door{
	{
		param: ReferenceParam, key: KeyReference, field: "Reference", takes: 1, needs: " alone",
		doc: "the function that returns a reference implementation from the handle of a check on the planted defect",
	},
	{
		param: SecondParam, key: KeySecond, field: "Second", takes: 1, needs: " alone",
		doc: "the function that returns a second subject from the handle of a check on the planted defect",
	},
	{
		param: RestartParam, key: KeyRestart, field: "Restart", takes: 2, needs: " and the planted defect",
		doc: "the function that returns the planted defect after a restart, from the handle of a check and the defect",
	},
	{
		param: CrashParam, key: KeyCrash, field: "Crash", takes: 3,
		needs: ", the planted defect and the call that it crashes",
		doc: "the function that runs a call against the planted defect, crashes the medium of the defect during the " +
			"call, and returns the defect as the crash left it",
	},
}

// checksThat contains the words after the name of a check that open its doc
// comment, as in "KeepsTheItem checks that Get returns the item".
const checksThat = " checks that"

// propPackage is the import path of the property engine of the assert
// library. A property check takes a pointer to its type caseType after the
// handle of its check, and the function of its options returns a list of its
// type optionType.
const (
	propPackage = "go.dokimi.dev/assert/prop"
	caseType    = "Case"
	optionType  = "Option"
)

// testingPackage is the import path of the testing package. The function of
// the environment of a check takes a pointer to its type testType and a
// function of such a pointer, as synctest.Test does.
const (
	testingPackage = "testing"
	testType       = "T"
)

// Word is the word of the family of the suites and of the family of their
// tests, so the test of a suite is the spec of its file. TestTag is the tag of
// the family of the tests.
const (
	Word    = "suite"
	TestTag = sdk.TagTest
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
// suite, as in TestStoreSuite. The name of a test of allocation ceilings
// joins allocsWord to the name of a test, as in TestMemoryAllocs, and the
// name of a benchmark joins benchmarkWord to a subject, as in BenchmarkMemory.
const (
	checksWord      = "Checks"
	suiteWord       = "Suite"
	runWord         = "Run"
	fuzzWord        = "Fuzz"
	fromHandleWord  = "FromHandle"
	handleWord      = "Handle"
	checkWord       = "Check"
	optionWord      = "Option"
	methodsWord     = "Methods"
	funcWord        = "Func"
	callableWord    = "Callable"
	conformanceWord = "Conformance"
	proofsWord      = "Proofs"
	testWord        = "Test"
	conformsWord    = "ConformsTo"
	newWord         = "New"
	subjectWord     = "Subject"
	subjectsWord    = "Subjects"
	budgetsWord     = "Budgets"
	benchWord       = "Bench"
	parallelWord    = "Parallel"
	callsWord       = "Calls"
	overWord        = "Over"
	embeddedWord    = "Embedded"
	allocsWord      = "Allocs"
	benchmarkWord   = "Benchmark"
	declaredWord    = "Declared"
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
// Calls are the call generators of the callables, in the order of the
// callables. Using are the enumerations that the calls take at some depth,
// whose members the property runs and the benchmarks of the suite draw, and
// Starts are the generators of the structs with a starting value that the
// calls take. Generated are the generators of the function types and of the
// interfaces that the calls take at some depth. Variants are the generators
// of the interfaces with an unexported method, Streams the generators of the
// channels that a subject receives from, Readers the generators of the types
// that read drawn bytes, and Constants the generators of the types that take
// one fixed value. Sinks are the send channel types that the calls take,
// which the runtime gives the sink of each test. Dynamic reports that the
// calls take a value of the dynamic type, Go's any, at some depth, which the
// property runs and the benchmarks draw from a string, an int and a bool.
// Errors reports that the calls take an error outside the results of a
// generated function, which they draw from nil and conform.ErrGenerated.
// Inputs are the implementations of the interfaces that the callables and the
// generators take.
type Suite struct {
	Subject     *emit.TypeRef
	Methods     *emit.Interface
	Func        *emit.Method
	Declaration string
	Names       Names
	Callables   []Callable
	Checks      []Check
	Calls       []Call
	Using       []Enumeration
	Starts      []Start
	Generated   []Generated
	Variants    []Variants
	Streams     []Stream
	Readers     []Reader
	Constants   []Constant
	Sinks       []*emit.TypeRef
	Inputs      []Input
	Dynamic     bool
	Errors      bool
}

// Names are the names of the declarations of the suite of a declaration such
// as Store: the variable of the names of its callables StoreChecks, the
// function of its checks StoreSuite, the two forms RunStoreSuite and
// RunStoreSuiteFromHandle of its application, its fuzz form FuzzStoreSuite,
// the application of its allocation ceilings RunStoreBudgets, its benchmarks
// BenchStoreSuite and BenchStoreSuiteParallel, and the aliases StoreHandle,
// StoreCheck and StoreOption. Subjects is the function that turns the factory
// of an application into the factory of the subjects of the suite, as
// StoreSubjects. Each form of the application calls it, so the run of one
// form covers the factory that each form passes. Methods is the interface of
// the methods of a
// struct, as CatalogMethods, and Func the function type of a function, as
// SlugFunc. Each is empty for another declaration. Over is the type of the
// planted defects of the ceilings of an interface or a struct, as StoreOver,
// and NewOver the function that returns such a defect, as NewStoreOver.
// Embedded is the alias of the subject type that Over embeds, as
// StoreOverEmbedded. The embedded field has the name of the alias, so a
// callable with the name of the subject type is a valid method of Over. A
// function has no type of its defects, and a declaration without callables
// has none of the three names.
type Names struct {
	Checks        string
	Suite         string
	Run           string
	RunFromHandle string
	Fuzz          string
	Budgets       string
	Bench         string
	BenchParallel string
	Handle        string
	Check         string
	Option        string
	Subjects      string
	Methods       string
	Func          string
	Over          string
	NewOver       string
	Embedded      string
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
// "checks that". Property reports that the function takes a case of the
// property engine after the handle, so the suite passes it as the property of
// the check. Counted reports that the check counts the allocations of the
// whole process, so the runtime runs it with the allocation ceilings. Doors
// are the doors of the planted defect in the order of the fields of
// conform.Doors, and empty for a check without a door. Options refers to the
// function of the options of the property run of the check, and is nil for a
// check without options. Around refers to the function of the environment of
// the runs of the check, and is nil for a check without one.
type Check struct {
	Apply     *emit.TypeRef
	Defect    *emit.TypeRef
	Options   *emit.TypeRef
	Around    *emit.TypeRef
	Field     string
	Type      string
	Method    string
	Segment   string
	Statement string
	Doors     []Door
	Property  bool
	Counted   bool
}

// Door is one door of the planted defect of a check. Field is the field of
// conform.Doors that the check of the suite sets, as Restart, and Func refers
// to the function of the door.
type Door struct {
	Func  *emit.TypeRef
	Field string
}

// Conformance is the data of the record of one test of the family of the
// tests, and the template of the family reads it. Test is the name of the
// test function, and Declaration the name of the declaration whose suite the
// test applies. Subject is the name of the subject of a conformance test, and
// empty for the test of a suite. Excusals are the excusals of the subject in
// the order of their directives. The test of a suite applies a check of its
// own to the forms of the application, so its excusals are the checks of the
// suite. Function reports that the declaration is a function, whose form Run
// takes the subject itself. Builds and Receives are the segments of the
// checks that the test of a suite adds. The check of the form Run of the
// suite of a function has the segment Receives, and each other check has
// the segment Builds. Both are empty for a conformance test. Targets are the
// fuzz targets of the test: one for each callable in the conformance test of
// a subject and in the test of the suite of an interface or a struct, and one
// for the function in the test of the suite of a function. FromHandle reports
// that the factory of a subject takes the handle of a check.
//
// SuiteAllocs is the name of the test that falsifies the allocation ceilings
// and the checks that count allocations of the suite, as
// TestStoreSuiteAllocs. It is empty for the test of a subject and for a suite
// without either. Allocs is the name of the test that applies the allocation
// ceilings and the checks that count allocations to a subject or to a
// function, as TestMemoryAllocs, and empty where neither applies. Budgets is
// the name of the test that applies the allocation ceilings of the suite of an
// interface or a struct to its planted defect, as TestStoreSuiteBudgets, and
// empty for any other test. Benchmark is the name of the benchmark of a
// subject, a function or the suite of an interface or a struct, as
// BenchmarkMemory or BenchmarkStoreSuite, and Parallel the name of its
// parallel benchmark, as BenchmarkMemoryParallel. Parallel is empty where no
// callable of a level admits concurrent calls. Levels are the levels of the
// benchmarks of a subject or of a suite, one for each callable without the
// mixin waits. A test without a level declares no benchmark, and the test of
// the suite of a function that waits declares no benchmark and no allocation
// test. Ceilings are the budgets
// that the subject states for its callables, in the order of their
// directives. Cases is the number of valid cases of each property run that
// the subject states, and 0 for a subject without the param cases. The test
// passes it to the applications whose property runs it sets: the checks and
// the allocation ceilings. Identities are the identities of the checks of a
// suite, which the test of the suite compares with the names of the
// callables, and empty for the test of a subject. Options refers to the
// function of the param options of the subject, whose list each test of the
// subject passes after the options of its directives. It is nil for a subject
// without the param and for the test of a suite.
//
// The record of a conformance test is a struct with ten fields. Their types
// refer, in this order, to the form of the application, as RunStoreSuite, to
// the factory of the subjects, to the names of the callables of the suite, as
// StoreChecks, to its fuzz form, as FuzzStoreSuite, to its alias of the
// handle, to the type of the subjects that the factory returns, to the
// application of its allocation ceilings, as RunStoreBudgets, to its
// benchmark, as BenchStoreSuite, to its alias of the option, and to the type
// of the subjects of the suite. The record of a subject of a declaration with
// a callable that admits concurrent calls has an eleventh field, which refers
// to its parallel benchmark, as BenchStoreSuiteParallel. The record of the
// test of a suite ends with that field as well.
//
// The record of the test of a suite is a struct whose fields refer, in this
// order, to the function of the suite as StoreSuite, to the two forms of the
// application, to the aliases of the handle and the check, to the type of the
// subjects and to the names of the callables. The record of the suite of a
// function has four more fields, which refer to the function, to its fuzz
// form, to the application of its allocation ceilings and to its benchmark.
// The record of the suite of an interface or a struct with callables has three
// more fields, which refer to its fuzz form, to the application of its
// allocation ceilings and to its benchmark. The layout of eidos qualifies each
// reference to the suite by the package of the file of the suite.
type Conformance struct {
	Options     *emit.TypeRef
	Test        string
	Subject     string
	Declaration string
	Builds      string
	Receives    string
	SuiteAllocs string
	Allocs      string
	Budgets     string
	Benchmark   string
	Parallel    string
	Excusals    []Excusal
	Ceilings    []Ceiling
	Targets     []Target
	Levels      []Level
	Identities  []Identity
	Cases       int64
	Function    bool
	FromHandle  bool
}

// Identity is the identity of one check of a suite, which the test of the
// suite compares with the identity that the names of the callables return.
// Field is the field of its callable in the names of the callables, as Get,
// and empty for a function. Method is the method that returns the identity,
// as KeepsTheItem or AllocationBudget, and ID is the identity, as
// Store.Get/keeps-the-item.
type Identity struct {
	Field  string
	Method string
	ID     string
}

// Level is one level of the benchmarks of a subject or of a suite. Field is the
// field of its callable in the names of the callables, as Get. Concurrent
// reports that the callable admits concurrent calls, so the parallel benchmark
// has the level as well. Ceiled reports that the test of a suite adds an
// allocation ceiling of 0 to the callable: the declaration states no budget of
// it, and its calls build no implementation of an interface and no channel of
// their own, so a call of the planted defect allocates nothing.
type Level struct {
	Field      string
	Concurrent bool
	Ceiled     bool
}

// Target is one fuzz target of a test. Name is the name of the fuzz function,
// as FuzzStoreGet, and Field the field of its callable in the names of the
// callables, as Get. The target of the suite of a function has no field,
// because the function is the one callable of its suite.
type Target struct {
	Name  string
	Field string
}

// Excusal is one excusal of a check in the conformance test of a subject.
// Field is the field of the callable of the check in the names of the
// callables, as Get, and Method the method that returns the identity of the
// check, as KeepsTheItem or LatencyBudget. Reason is the reason why the
// subject does not pass the check.
type Excusal struct {
	Field  string
	Method string
	Reason string
}

// part is the part of suite for one language. It declares the target with its
// template tree and its helpers.
type part interface {
	lang.Part
}

// parts are the parts of suite, one for each language that it serves.
var parts = []part{golangPart{}}

// suiteSchema is the directive of a suite.
var suiteSchema = directive.Schema{
	Name: SuiteDirective,
	Doc: "generates the suite of the declaration: the names of its callables, the identities and the checks of " +
		"its callables, the functions that apply the checks to a subject, and the test of its planted defects",
}

// subjectSchema is the directive of a subject. eidos resolves the declaration
// to a type that the file of the subject can see, and the factory and the
// function of the options to functions in the scope of the subject. A struct
// takes one directive for each declaration whose suite it passes, so the
// directive repeats.
var subjectSchema = directive.Schema{
	Name:       SubjectDirective,
	Repeatable: true,
	Positional: []directive.ParamSpec{{
		Key: DeclarationParam, Type: directive.TypeReference, Resolution: directive.ResolveTypeInScope,
		Required: true, Doc: "the declaration under the directive suite whose checks the subject passes",
	}},
	Params: []directive.ParamSpec{
		{
			Key: NewParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
			Required: true,
			Doc: "the function that returns a fresh subject, without a parameter or from the handle of a check: " +
				"an exported function of the package of the struct, or a function of the package of its test",
		},
		{
			Key: CasesParam, Type: directive.TypeInt,
			Doc: "the number of valid cases of each property run of the application, 100 without the param",
		},
		{
			Key: OptionsParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
			Doc: "the function without parameters that returns the list of the options of the application, which " +
				"each test of the subject passes after the options of the directives",
		},
	},
	Doc: "generates the test that applies the checks of the suite of the declaration to the subject",
}

// checkSchema is the directive of a check of a consumer. suite resolves the
// callable among the declarations under the directive suite, and eidos
// resolves the defect and each door of the defect to a function in the scope
// of the check.
var checkSchema = directive.Schema{
	Name: CheckDirective,
	Positional: []directive.ParamSpec{{
		Key: CallableParam, Type: directive.TypeString, Required: true,
		Doc: "the callable of a declaration under the directive suite, such as Store.Get, or store.Store.Get with " +
			"the name of the package of the declaration",
	}},
	Params: append([]directive.ParamSpec{
		{
			Key: DefectParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
			Required: true,
			Doc:      "the function that returns the planted defect of the check from the handle of the check",
		},
		{
			Key: CountedParam, Type: directive.TypeBool,
			Doc: "true for a check that counts the allocations of the whole process, such as with assert.MaxAllocs. " +
				"The check runs with the allocation ceilings, in a test that does not run in parallel",
		},
		{
			Key: OptionsParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
			Doc: "the function that returns the options of the property run of the check in an application, " +
				"from a handle of the check",
		},
		{
			Key: AroundParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
			Doc: "the function that runs each run of the check inside an environment, such as a testing/synctest " +
				"bubble, from a *testing.T and the run, as synctest.Test does",
		},
	}, doorParams()...),
	Doc: "binds the function to a callable of a suite as a check of the consumer, with its planted defect",
}

// doorParams returns the param of the directive check for each door of the
// planted defect.
func doorParams() []directive.ParamSpec {
	out := make([]directive.ParamSpec, 0, len(doors))
	for _, d := range doors {
		out = append(out, directive.ParamSpec{
			Key: d.param, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope, Doc: d.doc,
		})
	}
	return out
}

// excuseSchema is the directive of an excusal. A subject excuses each check
// with a directive of its own, so the directive repeats, and it requires the
// directive subject on its struct.
var excuseSchema = directive.Schema{
	Name:       ExcuseDirective,
	Repeatable: true,
	Requires:   []directive.Name{SubjectDirective},
	Kinds:      []symbol.Kind{symbol.KindStruct},
	Positional: []directive.ParamSpec{{
		Key: IdentityParam, Type: directive.TypeString, Required: true,
		Doc: "the identity of the check of the suite that the subject does not pass, such as Store.Get/keeps-the-item",
	}},
	Params: []directive.ParamSpec{{
		Key: ReasonParam, Type: directive.TypeString, Required: true,
		Doc: "the reason why the subject does not pass the check",
	}},
	Doc: "excuses a check of the suite of the declaration from the conformance test of the subject, with the reason",
}

// generator is one instance of suite: the handles of the keys of the facts of
// a check and of a budget. The composition registers the keys before the
// first phase call, and the handlers read them afterwards.
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
	property    meta.Key[bool]
	counted     meta.Key[bool]
	allocs      meta.Key[int64]
	bytes       meta.Key[int64]
	latency     meta.Key[int64]
	mean        meta.Key[int64]
	warmup      meta.Key[int64]
	options     meta.Key[symbol.Identity]
	around      meta.Key[symbol.Identity]
	doors       [len(doors)]meta.Key[symbol.Identity]
}

// New returns suite, a plugin with an annotate half and a generate half. The
// annotate half stamps the facts of a check on each function under
// [CheckDirective] whose directive names a callable of a declaration under
// [SuiteDirective]. The generate half records a [Suite] into the family of the
// suites for each interface, struct or function under SuiteDirective, with the
// checks that the annotate half binds to its callables. The template of the
// family writes the suite in the conventions of the language of the plan.
//
// suite records a [Conformance] into the family of the tag [TestTag] for each
// declaration under SuiteDirective, and the template of the family writes the
// test that falsifies each check of the suite on its planted defect. Both
// families have the word [Word], so a Go suite and its test are
// <file>_suite.go and <file>_suite_test.go. suite also records a Conformance
// for each directive [SubjectDirective] of a struct, and the template writes
// a test that applies the suite of the declaration to a subject of the
// factory. The tests of the family test their package from outside, so a Go
// file of the family is in the external test package of its directory, or of
// the directory of its directive out. The factory of a subject is therefore an
// exported function of the package of the subject or a function of the package
// of its test. The test of a subject passes each excusal of the directive
// [ExcuseDirective] of its declaration on the struct to the application.
//
// suite serves Go, the language of its part, and leaves a declaration of
// another language alone. It reports [Uninstantiable] at a declaration whose
// type parameters its suite cannot instantiate, [Unimportable] at a
// declaration with a callable whose signature or parameter samples the file
// of its suite cannot spell, [NoSuite] at a subject whose declaration has no
// directive suite, [BadFactory] at a subject whose factory an application
// cannot call, and [BadOptions] at a subject whose function of the options a
// test cannot call.
// It reports [UnknownCallable] and [BadCheck] at the directive of a check that
// it cannot bind, and [UnknownExcusal] at an excusal of a check that the suite
// does not have.
// suite records and stamps nothing for a declaration that it reports. Each
// call returns a new plugin.
func New() sdk.Dual {
	g := &generator{}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Keys(g.register).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: suiteTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, External: true, Word: Word, Template: conformTemplate})
	lang.Register(b, parts)
	return b.Handle(
		sdk.Directive(suiteSchema,
			sdk.OnInterface(lang.Served(parts, g.suiteOfInterface)),
			sdk.OnStruct(lang.Served(parts, g.suiteOfStruct)),
			sdk.OnFunction(lang.Served(parts, g.suiteOfFunction))),
		sdk.Directive(subjectSchema, sdk.OnStruct(lang.Served(parts, g.conformanceOf))),
		sdk.Directive(checkSchema, sdk.OnFunction(lang.Served(parts, g.concludeCheck))),
		sdk.Directive(excuseSchema),
		sdk.Directive(budgetSchema,
			sdk.OnMethod(lang.Served(parts, g.budgetOfMethod)),
			sdk.OnFunction(lang.Served(parts, g.budgetOfFunction)),
			sdk.OnStruct(lang.Served(parts, g.budgetOfStruct))),
	).Dual()
}

// register claims the namespace of suite and registers the keys of the facts
// of a check and of a budget.
//
// # Errors
//
// register returns the error of the claim, such as a namespace that another
// registrant claimed, joined with the errors of the registrations, such as a
// key that the registrant registered with another type.
func (g *generator) register(r *meta.Registry) error {
	spec := func(name meta.KeyName, doc string) meta.KeySpec {
		return meta.KeySpec{Name: name, Kinds: []symbol.Kind{symbol.KindFunction}, Doc: doc}
	}
	budgetKey := func(name meta.KeyName, doc string) meta.KeySpec {
		return meta.KeySpec{Name: name, Kinds: []symbol.Kind{symbol.KindMethod, symbol.KindFunction}, Doc: doc}
	}
	var errs [15 + len(doors)]error
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
	g.property, errs[6] = meta.Register[bool](r, spec(KeyProperty,
		"reports that the check takes a case of the property engine after the handle of its check"))
	g.allocs, errs[7] = meta.Register[int64](r, budgetKey(KeyAllocs,
		"contains the most heap allocations of one call of the callable"))
	g.bytes, errs[8] = meta.Register[int64](r, budgetKey(KeyBytes,
		"contains the most heap bytes that one call of the callable allocates"))
	g.latency, errs[9] = meta.Register[int64](r, budgetKey(KeyLatency,
		"contains the highest p99 latency of one call of the callable, in nanoseconds"))
	g.mean, errs[10] = meta.Register[int64](r, budgetKey(KeyMean,
		"contains the highest mean latency of one call of the callable, in nanoseconds"))
	g.warmup, errs[11] = meta.Register[int64](r, budgetKey(KeyWarmup,
		"contains the number of iterations that each benchmark of the callable runs before it measures"))
	g.counted, errs[12] = meta.Register[bool](r, spec(KeyCounted,
		"reports that the check counts the allocations of the whole process"))
	g.options, errs[13] = meta.Register[symbol.Identity](r, spec(KeyOptions,
		"contains the identity of the function of the options of the property run of the check"))
	g.around, errs[14] = meta.Register[symbol.Identity](r, spec(KeyAround,
		"contains the identity of the function that runs each run of the check inside an environment"))
	for i, d := range doors {
		g.doors[i], errs[15+i] = meta.Register[symbol.Identity](r, spec(d.key,
			"contains the identity of the function of the door "+string(d.param)+" of the planted defect of the check"))
	}
	return errors.Join(errs[:]...)
}

// concludeCheck stamps the facts of a check on a function under the directive
// check. It stamps them where the callable of the directive belongs to one
// declaration under the directive suite. The function and its defect must also
// have the signatures of a check. The identity of a check of a function names
// no callable, so a check of a function must not have the segment of a check
// that the test of the suite adds. Otherwise concludeCheck reports
// [UnknownCallable] or [BadCheck] at the directive.
func (g *generator) concludeCheck(m *sdk.FunctionMatch, _ part, st *sdk.Stamper) error {
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
	property, ok := signed(m, d, fn)
	if !ok || !doored(m, d, fn) || !optioned(m, d, fn, property) || !arounded(m, d, fn) {
		return nil
	}
	segment := naming.Kebab(fn.Name)
	if !strings.Contains(callable, ".") && (segment == buildsSegment || segment == receivesSegment) {
		m.ErrorfAt(BadCheck, d.Pos, "the test of the suite of %s adds the check %s/%s to the forms of its "+
			"application, so %s is no check. Rename the function", callable, callable, segment, fn.Name)
		return nil
	}
	if slices.Contains(ceilingSegments, segment) {
		m.ErrorfAt(BadCheck, d.Pos, "a budget of %s adds the check %s/%s, so %s is no check. Rename the function",
			callable, callable, segment, fn.Name)
		return nil
	}
	sdk.Stamp(st, g.declaration, decl)
	sdk.Stamp(st, g.callable, callable)
	sdk.Stamp(st, g.segment, segment)
	sdk.Stamp(st, g.statement, statement)
	sdk.Stamp(st, g.defect, d.Params[DefectParam].Target)
	for i, door := range doors {
		if opened, stated := d.Param(door.param); stated {
			sdk.Stamp(st, g.doors[i], opened.Target)
		}
	}
	if options, stated := d.Param(OptionsParam); stated {
		sdk.Stamp(st, g.options, options.Target)
	}
	if around, stated := d.Param(AroundParam); stated {
		sdk.Stamp(st, g.around, around.Target)
	}
	if property {
		sdk.Stamp(st, g.property, true)
	}
	if counted, _ := d.Param(CountedParam); counted.Bool {
		sdk.Stamp(st, g.counted, true)
	}
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
	m.ErrorfAt(UnknownCallable, d.Pos, "%s has no callable %s, so %s is no check%s", head, method, check,
		callablesOf(head, callables))
	return symbol.Identity{}, "", false
}

// callablesOf returns the sentence that lists the callables of the
// declaration decl after a sentence that names a callable that decl lacks,
// such as ". The callables of Store are Get and Put". It returns the empty
// string for a declaration without callables.
func callablesOf(decl string, callables []string) string {
	switch {
	case len(callables) > 1:
		return fmt.Sprintf(". The callables of %s are %s", decl, listed(callables, "and"))
	case len(callables) == 1:
		return fmt.Sprintf(". The callable of %s is %s", decl, callables[0])
	}
	return ""
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
// the signatures of a check, and whether fn is a property check. A check takes
// the handle of its check and returns nothing. A property check also takes a
// *prop.Case after the handle. The defect returns the planted defect from the
// handle of its check. signed reports [BadCheck] at d where a signature
// differs.
func signed(m *sdk.FunctionMatch, d *directive.Directive, fn *node.Function) (property, ok bool) {
	defect := d.Params[DefectParam].Target
	b := m.Rules()
	check, _ := b.CallableOf(fn)
	found, _ := m.Reader().Lookup(defect)
	planted, _ := b.CallableOf(found)
	switch {
	case len(check.Params) == 2 && !takesCase(check.Params[1].Ref):
		m.ErrorfAt(BadCheck, d.Pos, "the second parameter of the check %s has the type %s, and a property check "+
			"takes a *prop.Case after the handle of its check", fn.Name, check.Params[1].Ref.Spelling)
	case len(check.Params) != 1 && len(check.Params) != 2:
		m.ErrorfAt(BadCheck, d.Pos, "the check %s takes %d parameters, and a check takes the handle of its check, "+
			"with a *prop.Case after it for a property check", fn.Name, len(check.Params))
	case len(check.Returns) > 0:
		m.ErrorfAt(BadCheck, d.Pos, "the check %s returns a result, and a check returns nothing", fn.Name)
	case len(planted.Params) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the defect %s of %s takes %d parameters, and a defect takes the handle of its "+
			"check alone", defect.Name, fn.Name, len(planted.Params))
	case len(planted.Returns) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the defect %s of %s returns %d results, and a defect returns the planted "+
			"defect alone", defect.Name, fn.Name, len(planted.Returns))
	default:
		return len(check.Params) == 2, true
	}
	return false, false
}

// doored reports whether each door of the planted defect of the check fn has
// the signature of its door. A door returns its subject alone. A reference
// and a second take the handle of the check, a restart the handle and the
// planted defect, and a crash the handle, the planted defect and the call that
// it crashes, a function. doored reports [BadCheck] at the directive d where a
// signature differs.
func doored(m *sdk.FunctionMatch, d *directive.Directive, fn *node.Function) bool {
	b := m.Rules()
	for _, door := range doors {
		param, stated := d.Param(door.param)
		if !stated {
			continue
		}
		found, _ := m.Reader().Lookup(param.Target)
		view, _ := b.CallableOf(found)
		name := param.Target.Name
		switch {
		case len(view.Params) != door.takes:
			m.ErrorfAt(BadCheck, d.Pos, "the door %s %s of %s takes %s, and the door %s takes the handle of its "+
				"check%s", door.param, name, fn.Name, parameters(len(view.Params)), door.param, door.needs)
		case len(view.Returns) != 1:
			m.ErrorfAt(BadCheck, d.Pos, "the door %s %s of %s returns %d results, and a door returns its subject "+
				"alone", door.param, name, fn.Name, len(view.Returns))
		case door.param == CrashParam && view.Params[2].Ref.Form != symbol.FormFunc:
			m.ErrorfAt(BadCheck, d.Pos, "the third parameter of the door %s %s of %s has the type %s, and the door "+
				"%s takes the call that it crashes as a func()", door.param, name, fn.Name,
				view.Params[2].Ref.Spelling, door.param)
		default:
			continue
		}
		return false
	}
	return true
}

// optioned reports whether the options of the check fn have the signature of
// the options of a property run. The options apply to the property run of a
// check, so fn takes a case of the property engine, which property reports.
// The function of the options takes the handle of the check and returns a
// list of prop.Option. optioned reports [BadCheck] at the directive d where
// the check or the signature differs.
func optioned(m *sdk.FunctionMatch, d *directive.Directive, fn *node.Function, property bool) bool {
	param, stated := d.Param(OptionsParam)
	if !stated {
		return true
	}
	found, _ := m.Reader().Lookup(param.Target)
	b := m.Rules()
	view, _ := b.CallableOf(found)
	name := param.Target.Name
	switch {
	case !property:
		m.ErrorfAt(BadCheck, d.Pos, "the options %s of %s apply to the property run of a check, and %s takes no "+
			"*prop.Case", name, fn.Name, fn.Name)
	case len(view.Params) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the options %s of %s take %s, and the options of a check take the handle of "+
			"its check alone", name, fn.Name, parameters(len(view.Params)))
	case len(view.Returns) != 1:
		m.ErrorfAt(BadCheck, d.Pos, "the options %s of %s return %d results, and the options of a check return "+
			"[]prop.Option", name, fn.Name, len(view.Returns))
	case !listsOptions(view.Returns[0].Ref):
		m.ErrorfAt(BadCheck, d.Pos, "the options %s of %s return the type %s, and the options of a check return "+
			"[]prop.Option", name, fn.Name, view.Returns[0].Ref.Spelling)
	default:
		return true
	}
	return false
}

// arounded reports whether the environment of the check fn has the signature
// of synctest.Test. The function of the environment takes a *testing.T and a
// func(*testing.T), and returns nothing. arounded reports [BadCheck] at the
// directive d where the signature differs.
func arounded(m *sdk.FunctionMatch, d *directive.Directive, fn *node.Function) bool {
	param, stated := d.Param(AroundParam)
	if !stated {
		return true
	}
	found, _ := m.Reader().Lookup(param.Target)
	b := m.Rules()
	view, _ := b.CallableOf(found)
	name := param.Target.Name
	const takes = "and the environment of a check takes a *testing.T and a func(*testing.T)"
	switch {
	case len(view.Params) != 2:
		m.ErrorfAt(BadCheck, d.Pos, "the environment %s of %s takes %s, %s", name, fn.Name,
			parameters(len(view.Params)), takes)
	case !takesTest(view.Params[0].Ref):
		m.ErrorfAt(BadCheck, d.Pos, "the first parameter of the environment %s of %s has the type %s, %s", name,
			fn.Name, view.Params[0].Ref.Spelling, takes)
	case !runsTest(view.Params[1].Ref):
		m.ErrorfAt(BadCheck, d.Pos, "the second parameter of the environment %s of %s has the type %s, %s", name,
			fn.Name, view.Params[1].Ref.Spelling, takes)
	case len(view.Returns) > 0:
		m.ErrorfAt(BadCheck, d.Pos, "the environment %s of %s returns a result, and the environment of a check "+
			"returns nothing", name, fn.Name)
	default:
		return true
	}
	return false
}

// takesTest reports whether ref is a pointer to the test of the testing
// package, whatever name the file gives the package.
func takesTest(ref *node.TypeRef) bool {
	return pointerTo(ref, testingPackage, testType)
}

// runsTest reports whether ref is a function type that takes a pointer to the
// test of the testing package alone and returns nothing. Split is the index of
// the first result of a function type and 0 for every other form, so any other
// form has no parameter.
func runsTest(ref *node.TypeRef) bool {
	params, results := ref.Elems[:ref.Split], ref.Elems[ref.Split:]
	return len(params) == 1 && len(results) == 0 && takesTest(params[0])
}

// listsOptions reports whether ref is a list of the options of the property
// engine, whatever name the file gives the package of the engine.
func listsOptions(ref *node.TypeRef) bool {
	if ref.Form != symbol.FormList {
		return false
	}
	elem := ref.Elems[0]
	return elem.Package == propPackage && elem.Spelling[strings.LastIndex(elem.Spelling, ".")+1:] == optionType
}

// parameters returns the count of n parameters, as 1 parameter or 2
// parameters.
func parameters(n int) string {
	if n == 1 {
		return "1 parameter"
	}
	return fmt.Sprintf("%d parameters", n)
}

// takesCase reports whether ref is a pointer to the case of the property
// engine, whatever name the file gives the package of the engine.
func takesCase(ref *node.TypeRef) bool {
	return pointerTo(ref, propPackage, caseType)
}

// pointerTo reports whether ref is a pointer to the type named name of the
// package pkg, whatever name the file gives the package. Go spells a pointer
// in the optional form, which has one child.
func pointerTo(ref *node.TypeRef, pkg, name string) bool {
	if ref.Form != symbol.FormOptional {
		return false
	}
	elem := ref.Elems[0]
	return elem.Package == pkg && elem.Spelling[strings.LastIndex(elem.Spelling, ".")+1:] == name
}

// bound is a check that the annotate half binds to a declaration: the
// function of the check and the facts that the annotate half stamps on it.
// doors are the functions of the doors of the planted defect in the order of
// [doors], and the zero identity for a door that the directive does not state.
type bound struct {
	fn        symbol.Identity
	callable  string
	segment   string
	statement string
	defect    symbol.Identity
	options   symbol.Identity
	around    symbol.Identity
	doors     [len(doors)]symbol.Identity
	property  bool
	counted   bool
}

// boundTo returns the checks that the annotate half binds to the declaration
// origin, in the order of the identities of their functions.
func (g *generator) boundTo(m sdk.Matcher, origin *symbol.Identity) []bound {
	var out []bound
	for id := range sdk.Subjects(m, g.declaration) {
		if decl, _ := sdk.FactOf(m, id, g.declaration); decl != *origin {
			continue
		}
		b := bound{fn: id}
		b.callable, _ = sdk.FactOf(m, id, g.callable)
		b.segment, _ = sdk.FactOf(m, id, g.segment)
		b.statement, _ = sdk.FactOf(m, id, g.statement)
		b.defect, _ = sdk.FactOf(m, id, g.defect)
		b.options, _ = sdk.FactOf(m, id, g.options)
		b.around, _ = sdk.FactOf(m, id, g.around)
		for i := range doors {
			b.doors[i], _ = sdk.FactOf(m, id, g.doors[i])
		}
		b.property, _ = sdk.FactOf(m, id, g.property)
		b.counted, _ = sdk.FactOf(m, id, g.counted)
		out = append(out, b)
	}
	return out
}

// checksOf returns the checks that the annotate half binds to the callables
// of the declaration origin, in the order of the callables. The checks of one
// callable are in the order of the identities of their functions.
func (g *generator) checksOf(m sdk.Matcher, origin *symbol.Identity, callables []Callable) []Check {
	byCallable := map[string][]Check{}
	checks := g.boundTo(m, origin)
	for i := range checks {
		b := &checks[i]
		var opened []Door
		for i, door := range doors {
			if fn := b.doors[i]; fn != (symbol.Identity{}) {
				opened = append(opened, Door{Field: door.field, Func: &emit.TypeRef{Spelling: fn.Name, Target: fn}})
			}
		}
		var options, around *emit.TypeRef
		if fn := b.options; fn != (symbol.Identity{}) {
			options = &emit.TypeRef{Spelling: fn.Name, Target: fn}
		}
		if fn := b.around; fn != (symbol.Identity{}) {
			around = &emit.TypeRef{Spelling: fn.Name, Target: fn}
		}
		byCallable[b.callable] = append(byCallable[b.callable], Check{
			Apply:  &emit.TypeRef{Spelling: b.fn.Name, Target: b.fn},
			Defect: &emit.TypeRef{Spelling: b.defect.Name, Target: b.defect},
			Method: b.fn.Name, Segment: b.segment, Statement: b.statement, Doors: opened, Options: options,
			Around: around, Property: b.property, Counted: b.counted,
		})
	}
	var out []Check
	for _, c := range callables {
		bound := byCallable[c.Name]
		for i := range bound {
			bound[i].Field, bound[i].Type = c.Field, c.Type
			out = append(out, bound[i])
		}
	}
	return out
}

// excusalsOf returns the excusals of the subject of m in the order of their
// directives. Each excusal excuses a check that the annotate half binds to the
// declaration declared, or the check of a ceiling of the budget of one of its
// callables. excusalsOf reports [UnknownExcusal] at each excusal whose
// identity does not match such a check, with the identities of the checks of
// its callable: the checks in the order of their functions, then the
// ceilings. It returns false where it reports an excusal. A subject of a
// struct with several directives subject reads the excusals of the callables
// of its own declaration, as [skips] states.
func (g *generator) excusalsOf(
	m *sdk.StructMatch, declared *symbol.Identity, callables []*node.Method,
) ([]Excusal, bool) {
	var checks []excusable
	bindings := g.boundTo(m, declared)
	for i := range bindings {
		b := &bindings[i]
		checks = append(checks, excusable{callable: b.callable, segment: b.segment, method: b.fn.Name})
	}
	for _, c := range callables {
		for _, segment := range g.budgetFact(m, &c.ID).segments() {
			checks = append(checks, excusable{
				callable: declared.Name + "." + c.Name, segment: segment, method: naming.Pascal(segment),
			})
		}
	}
	var out []Excusal
	excused := true
	for x := range m.Directives(ExcuseDirective) {
		identity := x.Args[0].Str
		callable, _, _ := strings.Cut(identity, "/")
		if skipped, refused := skips(m, callable); skipped {
			excused = excused && !refused
			continue
		}
		if at := slices.IndexFunc(checks, func(c excusable) bool { return c.identity() == identity }); at >= 0 {
			_, field, _ := strings.Cut(callable, ".")
			out = append(out, Excusal{Field: field, Method: checks[at].method, Reason: x.Params[ReasonParam].Str})
			continue
		}
		var ids []string
		for _, c := range checks {
			if c.callable == callable {
				ids = append(ids, c.identity())
			}
		}
		listing := fmt.Sprintf("The suite of %s has no check of %s", declared.Name, callable)
		switch {
		case len(ids) > 1:
			listing = fmt.Sprintf("The checks of %s are %s", callable, listed(ids, "and"))
		case len(ids) == 1:
			listing = fmt.Sprintf("The check of %s is %s", callable, ids[0])
		}
		m.ErrorfAt(UnknownExcusal, x.Pos, "%s is no check of the suite of %s, so %s has no conformance test. %s",
			identity, declared.Name, m.Struct.Name, listing)
		excused = false
	}
	return out, excused
}

// excusable is a check that a subject can excuse: a check of a consumer, or
// the check of a ceiling of a budget. callable is its callable, as Store.Get,
// segment the last part of its identity, and method the method of the names
// of the callables that returns its identity.
type excusable struct {
	callable, segment, method string
}

// identity returns the identity of the check, as Store.Get/keeps-the-item.
func (c excusable) identity() string {
	return c.callable + "/" + c.segment
}

// listed returns words, two or more, as an English list whose last two words
// the conjunction joins, such as "Get and Put" or "a, b or c".
func listed(words []string, conjunction string) string {
	last := len(words) - 1
	return strings.Join(words[:last], ", ") + " " + conjunction + " " + words[last]
}

// suiteOfInterface records the suite of an interface. The subject type is the
// interface at the witnesses of its type parameters.
func (g *generator) suiteOfInterface(m *sdk.InterfaceMatch, _ part, e *sdk.Emitter) error {
	i := m.Interface
	b := m.Rules()
	subject := &emit.TypeRef{Spelling: i.Name, Target: i.ID}
	k := caller{m: m, g: g, e: e, b: b, origin: i.ID, decl: i.Name, typeParams: i.TypeParams}
	if len(i.TypeParams) > 0 {
		k.witnesses = b.Witnesses(i.TypeParams)
		if k.witnesses == nil {
			m.Errorf(Uninstantiable, "the type parameters of %s have no witnesses, so %s has no suite", i.Name, i.Name)
			return nil
		}
		for _, w := range k.witnesses {
			subject.Args = append(subject.Args, rules.EmitRef(w))
		}
	}
	methods := methodsOf(&b, i)
	if unimportable(m, &b, e.File(), i.Name, methods) {
		return nil
	}
	s := Suite{Subject: subject, Declaration: i.Name, Names: namesOf(e, i.Name)}
	g.callsOf(&k, &s, methods)
	write(e, &i.ID, &s)
	return nil
}

// suiteOfStruct records the suite of a struct. The subject type is the
// interface of the exported methods of the struct, which the suite declares.
func (g *generator) suiteOfStruct(m *sdk.StructMatch, _ part, e *sdk.Emitter) error {
	st := m.Struct
	if len(st.TypeParams) > 0 {
		m.Errorf(Uninstantiable, "%s is a generic struct, so it has no suite", st.Name)
		return nil
	}
	b := m.Rules()
	callables := methodsOf(&b, st)
	if unimportable(m, &b, e.File(), st.Name, callables) {
		return nil
	}
	names := namesOf(e, st.Name)
	names.Methods = e.JoinName(methodsWord, st.Name)
	methods := &emit.Interface{Origin: st.ID, Name: names.Methods, Visibility: symbol.VisibilityPublic}
	s := Suite{
		Subject: &emit.TypeRef{Spelling: names.Methods}, Methods: methods, Declaration: st.Name, Names: names,
	}
	for _, method := range callables {
		methods.Methods.Append(sdk.Mirror(names.Methods, method))
	}
	g.callsOf(&caller{m: m, g: g, e: e, b: b, origin: st.ID, decl: st.Name}, &s, callables)
	write(e, &st.ID, &s)
	return nil
}

// callsOf records into the suite s the callable and the call of each method
// of methods, the checks that the annotate half binds to them, the overrides
// that the calls take, and the implementations of the interfaces that they
// take. k writes the calls. A suite with callables names the type of its
// planted defects and the function that returns one.
func (g *generator) callsOf(k *caller, s *Suite, methods []*node.Method) {
	for _, method := range methods {
		c := callableOf(k.e, s.Declaration, method.Name)
		s.Callables = append(s.Callables, c)
		s.Calls = append(s.Calls, k.call(&c, method, &method.ID))
	}
	k.overrides(s)
	s.Checks = g.checksOf(k.m, &k.origin, s.Callables)
	if len(methods) > 0 {
		s.Names.Over = k.e.JoinName(overWord, s.Declaration)
		s.Names.NewOver = k.e.JoinName(s.Names.Over, newWord)
		s.Names.Embedded = k.e.JoinName(embeddedWord, s.Names.Over)
	}
}

// suiteOfFunction records the suite of a function. The subject type is the
// function type of the function, which the suite declares.
func (g *generator) suiteOfFunction(m *sdk.FunctionMatch, _ part, e *sdk.Emitter) error {
	f := m.Function
	if len(f.TypeParams) > 0 {
		m.Errorf(Uninstantiable, "%s is a generic function, so it has no suite", f.Name)
		return nil
	}
	b := m.Rules()
	signature := &node.Method{ID: f.ID, Name: f.Name, Visibility: f.Visibility, Params: f.Params, Returns: f.Returns}
	if unreferable(m, &b, e.File(), f.Name, []*node.Method{signature}) {
		return nil
	}
	names := namesOf(e, f.Name)
	names.Func = e.JoinName(funcWord, f.Name)
	names.NewOver = e.JoinName(e.JoinName(overWord, f.Name), newWord)
	callables := []Callable{{Name: f.Name, Type: e.JoinName(callableWord, f.Name)}}
	k := caller{m: m, g: g, e: e, b: b, origin: f.ID, decl: f.Name}
	s := Suite{
		Subject: &emit.TypeRef{Spelling: names.Func}, Func: sdk.Mirror(names.Func, signature), Declaration: f.Name,
		Names: names, Callables: callables, Checks: g.checksOf(m, &f.ID, callables),
		Calls: []Call{k.call(&callables[0], f, &f.ID)},
	}
	k.overrides(&s)
	write(e, &f.ID, &s)
	return nil
}

// write records the suite s of the declaration origin into the family of the
// suites. The record declares each name that the template writes: the variable
// of the names of the callables first, then the interface of the methods or the
// function type, the type of each callable, the variable of each call and of
// each generator of a starting value, each implementation of an interface,
// the type of the planted defects and the function that returns one, the
// functions and the aliases. Each implementation of an interface and the type
// of the planted defects declare the alias of the type that they embed after
// their own name. A suite without a call that admits concurrent
// calls has no parallel benchmark, so its name of the benchmark is empty. write
// also records the falsification test of the suite into the family of the
// tests.
func write(e *sdk.Emitter, origin *symbol.Identity, s *Suite) {
	public := func(name string) symbol.Symbol {
		return &emit.Function{Origin: *origin, Name: name, Visibility: symbol.VisibilityPublic}
	}
	alias := func(name string) symbol.Symbol {
		return &emit.Alias{Origin: *origin, Name: name, Visibility: symbol.VisibilityPublic}
	}
	n := &s.Names
	if !slices.ContainsFunc(s.Calls, func(c Call) bool { return c.Concurrent }) {
		n.BenchParallel = ""
	}
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
	for i := range s.Calls {
		c := &s.Calls[i]
		decls = append(decls, &emit.Variable{Origin: *origin, Name: c.Var, Visibility: symbol.VisibilityPackage})
	}
	for i := range s.Starts {
		st := &s.Starts[i]
		decls = append(decls, &emit.Variable{Origin: *origin, Name: st.Name, Visibility: symbol.VisibilityPackage})
	}
	for i := range s.Inputs {
		in := &s.Inputs[i]
		decls = append(decls,
			&emit.Struct{Origin: *origin, Name: in.Name, Visibility: symbol.VisibilityPackage},
			&emit.Alias{Origin: *origin, Name: in.Embedded, Visibility: symbol.VisibilityPackage})
	}
	if n.Over != "" {
		decls = append(decls,
			&emit.Struct{Origin: *origin, Name: n.Over, Visibility: symbol.VisibilityPackage},
			&emit.Alias{Origin: *origin, Name: n.Embedded, Visibility: symbol.VisibilityPackage})
	}
	if n.NewOver != "" {
		decls = append(decls, &emit.Function{Origin: *origin, Name: n.NewOver, Visibility: symbol.VisibilityPackage})
	}
	decls = append(decls, public(n.Suite), public(n.Run), public(n.RunFromHandle), public(n.Fuzz),
		public(n.Budgets), public(n.Bench))
	if n.BenchParallel != "" {
		decls = append(decls, public(n.BenchParallel))
	}
	decls = append(decls, alias(n.Handle), alias(n.Check), alias(n.Option),
		&emit.Function{Origin: *origin, Name: n.Subjects, Visibility: symbol.VisibilityPackage})
	e.File().Record(*s, decls...)
	record := &emit.Struct{
		Origin: *origin, Name: e.JoinName(proofsWord, s.Declaration), Visibility: symbol.VisibilityPackage,
	}
	emitted := func(name string) *emit.TypeRef { return &emit.TypeRef{Spelling: name, Target: *origin, Emitted: true} }
	subject := s.Subject
	if s.Methods != nil || s.Func != nil {
		subject = emitted(s.Subject.Spelling)
	}
	record.Fields.Append(
		&emit.Field{Name: suiteWord, Type: emitted(n.Suite)},
		&emit.Field{Name: runWord, Type: emitted(n.Run)},
		&emit.Field{Name: fromHandleWord, Type: emitted(n.RunFromHandle)},
		&emit.Field{Name: handleWord, Type: emitted(n.Handle)},
		&emit.Field{Name: checkWord, Type: emitted(n.Check)},
		&emit.Field{Name: subjectWord, Type: subject},
		&emit.Field{Name: checksWord, Type: emitted(n.Checks)},
	)
	excused := make([]Excusal, 0, len(s.Checks))
	for i := range s.Checks {
		c := &s.Checks[i]
		excused = append(excused, Excusal{Field: c.Field, Method: c.Method, Reason: aloneReason})
	}
	test := Conformance{
		Test: e.JoinName(n.Suite, testWord), Declaration: s.Declaration, Builds: buildsSegment,
		Receives: receivesSegment, Excusals: excused, Function: s.Func != nil,
	}
	allocated := slices.ContainsFunc(s.Calls, func(c Call) bool { return c.Budget.Allocs != nil }) ||
		slices.ContainsFunc(s.Checks, func(c Check) bool { return c.Counted })
	if allocated {
		test.SuiteAllocs = e.JoinName(allocsWord, test.Test)
	}
	if s.Func != nil {
		record.Fields.Append(
			&emit.Field{Name: funcWord, Type: &emit.TypeRef{Spelling: s.Declaration, Target: *origin}},
			&emit.Field{Name: fuzzWord, Type: emitted(n.Fuzz)},
			&emit.Field{Name: budgetsWord, Type: emitted(n.Budgets)},
			&emit.Field{Name: benchWord, Type: emitted(n.Bench)},
		)
		test.Targets = []Target{{Name: e.JoinName(s.Declaration, fuzzWord)}}
		if !s.Calls[0].Waits {
			test.benchmarks(e, s.Declaration, allocated, s.Calls[0].Concurrent)
		}
	} else if len(s.Calls) > 0 {
		record.Fields.Append(
			&emit.Field{Name: fuzzWord, Type: emitted(n.Fuzz)},
			&emit.Field{Name: budgetsWord, Type: emitted(n.Budgets)},
			&emit.Field{Name: benchWord, Type: emitted(n.Bench)},
		)
		test.idle(e, n, s.Calls)
	}
	if n.BenchParallel != "" {
		record.Fields.Append(&emit.Field{Name: parallelWord, Type: emitted(n.BenchParallel)})
	}
	test.Identities = identitiesOf(s)
	e.File(TestTag).Record(test, record)
}

// identitiesOf returns the identity of each check of the suite s, in the order
// of the callables: the checks of the consumer, then the ceilings of the
// budget of the callable.
func identitiesOf(s *Suite) []Identity {
	var out []Identity
	for i := range s.Callables {
		c := &s.Callables[i]
		for j := range s.Checks {
			if check := &s.Checks[j]; check.Field == c.Field {
				out = append(out, Identity{Field: c.Field, Method: check.Method, ID: c.Name + "/" + check.Segment})
			}
		}
		for _, segment := range s.Calls[i].Budget.segments() {
			out = append(out, Identity{Field: c.Field, Method: naming.Pascal(segment), ID: c.Name + "/" + segment})
		}
	}
	return out
}

// idle names the tests of the suite of an interface or a struct that apply the
// fuzz form, the allocation ceilings and the benchmarks of the suite to its
// planted defect with a spend of nothing, which answers each call at once. The
// test has a fuzz target for each call of calls, and a level for each call
// that does not wait. It has a parallel benchmark where a callable of a level
// admits concurrent calls.
func (c *Conformance) idle(e *sdk.Emitter, n *Names, calls []Call) {
	c.Budgets = e.JoinName(budgetsWord, c.Test)
	c.Benchmark = e.JoinName(n.Suite, benchmarkWord)
	fuzz := e.JoinName(n.Suite, fuzzWord)
	concurrent := false
	for i := range calls {
		call := &calls[i]
		c.Targets = append(c.Targets, Target{Name: e.JoinName(call.Field, fuzz), Field: call.Field})
		if call.Waits {
			continue
		}
		builds := slices.ContainsFunc(call.Params, func(p Param) bool { return p.Input != "" || p.Receive })
		c.Levels = append(c.Levels, Level{
			Field: call.Field, Concurrent: call.Concurrent, Ceiled: call.Budget == Budget{} && !builds,
		})
		concurrent = concurrent || call.Concurrent
	}
	if concurrent {
		c.Parallel = e.JoinName(parallelWord, c.Benchmark)
	}
}

// benchmarks names the test of the allocation ceilings and the benchmarks of
// the test of the subject or the function name. allocated reports that a
// ceiling of allocations applies to a callable, and concurrent that a
// callable admits concurrent calls. The test has a test of the allocation
// ceilings where one applies, and a parallel benchmark where a callable
// admits concurrent calls.
func (c *Conformance) benchmarks(e *sdk.Emitter, name string, allocated, concurrent bool) {
	if allocated {
		c.Allocs = e.JoinName(allocsWord, e.JoinName(name, testWord))
	}
	c.Benchmark = e.JoinName(name, benchmarkWord)
	if concurrent {
		c.Parallel = e.JoinName(parallelWord, c.Benchmark)
	}
}

// aloneReason is the reason of each excusal of the test of a suite, which
// applies a check of its own alone to the forms of the application.
const aloneReason = "the test applies the check of its option alone"

// buildsSegment and receivesSegment are the segments of the checks that the
// test of a suite adds to the forms of its application, as the record of the
// test states them in [Conformance].
const (
	buildsSegment   = "builds-the-subject"
	receivesSegment = "receives-the-function"
)

// conformanceOf records the conformance test of a struct under the directive
// subject, with the excusals and the budgets of the struct. It reports
// [NoSuite] where the declaration has no directive suite, and [BadFactory]
// where an application cannot call the factory. The test calls the factory
// from the package of its routed file: the external test package of the
// directory of the struct, or of the directory of its directive out. A factory
// that the package of the struct does not export, or a function of another
// external test package, is out of its reach. A factory without a parameter
// goes to the form Run of the application, and a factory of one parameter,
// the handle of a check, to the form RunFromHandle. The test of a factory of
// its own package calls the factory without a qualifier. The test applies the
// allocation ceilings to the subject where the budget of a callable states
// one, and has a level of its benchmarks for each callable. A struct has a
// conformance test for each declaration of its directives subject, and
// conformanceOf reports [RepeatedSubject] at a second directive of one
// declaration.
func (g *generator) conformanceOf(m *sdk.StructMatch, _ part, e *sdk.Emitter) error {
	st := m.Struct
	d := m.Directive()
	declared := d.Args[0].Target
	if !suited(m, &declared) {
		m.Errorf(NoSuite, "%s has no directive suite, so %s has no conformance test", declared.Name, st.Name)
		return nil
	}
	if first := firstSubject(m, &declared); first.Instance != d.Instance {
		m.ErrorfAt(RepeatedSubject, d.Pos, "%s has a directive subject of %s at line %d already, so this one has "+
			"no conformance test", st.Name, declared.Name, first.Pos.Line)
		return nil
	}
	factory := d.Params[NewParam].Target
	found, _ := m.Reader().Lookup(factory)
	b := m.Rules()
	view, _ := b.CallableOf(found)
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
	// The test calls the factory from the package of its routed file. Where
	// the plan previews no package, the layout checks the reference itself.
	if pkg, placed := e.File(TestTag).Package(); placed == nil {
		if err := b.Importable(pkg, &node.TypeRef{Spelling: factory.Name, Target: factory}); err != nil {
			m.Errorf(BadFactory, "the conformance test of %s in the package %s cannot call the factory %s, so %s "+
				"has no conformance test: %v", st.Name, pkg.Package, factory.Name, st.Name, err)
			return nil
		}
	}
	options, called := optionsOf(m, e, &b, d)
	if !called {
		return nil
	}
	decl, _ := m.Reader().Lookup(declared)
	methods := methodsOf(&b, decl)
	excusals, excused := g.excusalsOf(m, &declared, methods)
	ceilings, stated := g.ceilingsOf(m, &declared, methods)
	if !excused || !stated {
		return nil
	}
	allocated := slices.ContainsFunc(ceilings, func(c Ceiling) bool { return c.Budget.Allocs != nil }) ||
		slices.ContainsFunc(g.boundTo(m, &declared), func(b bound) bool { return b.counted })
	levels := make([]Level, 0, len(methods))
	for _, method := range methods {
		allocated = allocated || g.budgetFact(m, &method.ID).Allocs != nil
		if mixed(m, &method.ID, shape.Waits) {
			continue
		}
		concurrent := mixed(m, &method.ID, shape.Concurrent) || mixed(m, &method.ID, shape.ConcurrentReaders)
		levels = append(levels, Level{Field: method.Name, Concurrent: concurrent})
	}
	// A struct can be the subject of several declarations, so each name of its
	// test joins the struct and the declaration, as StoreConformsToStore.
	conforms := e.JoinName(declared.Name, e.JoinName(conformsWord, st.Name))
	record := &emit.Struct{
		Origin:     st.ID,
		Name:       e.JoinName(conformanceWord, conforms),
		Visibility: symbol.VisibilityPackage,
	}
	emitted := func(name string) *emit.TypeRef { return &emit.TypeRef{Spelling: name, Target: declared, Emitted: true} }
	record.Fields.Append(
		&emit.Field{Name: runWord, Type: emitted(run)},
		&emit.Field{Name: newWord, Type: &emit.TypeRef{Spelling: factory.Name, Target: factory}},
		&emit.Field{Name: checksWord, Type: emitted(names.Checks)},
		&emit.Field{Name: fuzzWord, Type: emitted(names.Fuzz)},
		&emit.Field{Name: handleWord, Type: emitted(names.Handle)},
		&emit.Field{Name: subjectWord, Type: rules.EmitRef(view.Returns[0].Ref)},
		&emit.Field{Name: budgetsWord, Type: emitted(names.Budgets)},
		&emit.Field{Name: benchWord, Type: emitted(names.Bench)},
		&emit.Field{Name: optionWord, Type: emitted(names.Option)},
		&emit.Field{Name: declaredWord, Type: g.subjectOf(e, &b, decl, emitted)},
	)
	if slices.ContainsFunc(levels, func(l Level) bool { return l.Concurrent }) {
		record.Fields.Append(&emit.Field{Name: parallelWord, Type: emitted(names.BenchParallel)})
	}
	var targets []Target
	fuzz := e.JoinName(conforms, fuzzWord)
	for _, method := range methods {
		targets = append(targets, Target{Name: e.JoinName(method.Name, fuzz), Field: method.Name})
	}
	test := Conformance{
		Test:    e.JoinName(conforms, testWord),
		Subject: st.Name, Declaration: declared.Name, Excusals: excusals, Ceilings: ceilings, Targets: targets,
		Levels: levels, Options: options, FromHandle: len(view.Params) == 1,
	}
	if cases, stated := d.Param(CasesParam); stated {
		test.Cases = cases.Int
	}
	test.benchmarks(e, conforms, allocated, slices.ContainsFunc(levels, func(l Level) bool { return l.Concurrent }))
	e.File(TestTag).Record(test, record)
	return nil
}

// optionsOf returns the reference to the function of the param options of the
// directive subject d, and nil for a directive without the param. The
// function takes no parameter and returns one list, and the test of the
// subject can call it. Otherwise optionsOf reports [BadOptions] at the subject
// and returns false.
func optionsOf(m *sdk.StructMatch, e *sdk.Emitter, b *rules.Bound, d *directive.Directive) (*emit.TypeRef, bool) {
	param, stated := d.Param(OptionsParam)
	if !stated {
		return nil, true
	}
	fn, subject := param.Target, m.Struct.Name
	found, _ := m.Reader().Lookup(fn)
	view, _ := b.CallableOf(found)
	switch {
	case len(view.Params) > 0:
		m.Errorf(BadOptions, "the function %s of the param options takes %s, so %s has no conformance test",
			fn.Name, parameters(len(view.Params)), subject)
		return nil, false
	case len(view.Returns) != 1:
		m.Errorf(BadOptions, "the function %s of the param options returns %d results, so %s has no conformance "+
			"test", fn.Name, len(view.Returns), subject)
		return nil, false
	case view.Returns[0].Ref.Form != symbol.FormList:
		m.Errorf(BadOptions, "the function %s of the param options returns the type %s, and the options of an "+
			"application are a list, so %s has no conformance test", fn.Name, view.Returns[0].Ref.Spelling, subject)
		return nil, false
	}
	// The test calls the function from the package of its routed file, as it
	// calls the factory.
	if pkg, placed := e.File(TestTag).Package(); placed == nil {
		if err := b.Importable(pkg, &node.TypeRef{Spelling: fn.Name, Target: fn}); err != nil {
			m.Errorf(BadOptions, "the conformance test of %s in the package %s cannot call the function %s of the "+
				"param options, so %s has no conformance test: %v", subject, pkg.Package, fn.Name, subject, err)
			return nil, false
		}
	}
	return &emit.TypeRef{Spelling: fn.Name, Target: fn}, true
}

// firstSubject returns the first directive subject of the struct of m whose
// declaration is declared, in source order.
func firstSubject(m *sdk.StructMatch, declared *symbol.Identity) *directive.Directive {
	first := m.Directive()
	for x := range m.Directives(SubjectDirective) {
		if x.Args[0].Target == *declared {
			first = x
			break
		}
	}
	return first
}

// ownerOf returns the first directive subject of the struct of m whose
// declaration has the name before the first dot of the callable, as Store of
// Store.Get. It returns nil where no directive subject of the struct names
// that declaration. Each directive subject reads the excusals and the budgets
// of the callables of its own declaration.
func ownerOf(m *sdk.StructMatch, callable string) *directive.Directive {
	owner, _, _ := strings.Cut(callable, ".")
	for x := range m.Directives(SubjectDirective) {
		if x.Args[0].Target.Name == owner {
			return x
		}
	}
	return nil
}

// skips reports whether the directive subject of m leaves the excusal or the
// budget of the callable to another directive. It leaves the callable of
// another declaration of the struct to the directive of that declaration.
// It leaves a callable of no declaration of the struct to the first
// directive subject, which reports it. refused reports a callable of no
// declaration, which no directive of the struct records a test for.
func skips(m *sdk.StructMatch, callable string) (skipped, refused bool) {
	gate := m.Directive()
	if owner := ownerOf(m, callable); owner != nil {
		return owner.Instance != gate.Instance, false
	}
	return gate.Instance != 0, true
}

// subjectOf returns the type of the subjects of the suite of the declaration
// decl: the interface at its witnesses, or the interface of the methods of a
// struct, which emitted refers to.
func (*generator) subjectOf(
	e *sdk.Emitter, b *rules.Bound, decl symbol.Symbol, emitted func(string) *emit.TypeRef,
) *emit.TypeRef {
	if st, isStruct := decl.(*node.Struct); isStruct {
		return emitted(e.JoinName(methodsWord, st.Name))
	}
	i, _ := decl.(*node.Interface)
	subject := &emit.TypeRef{Spelling: i.Name, Target: i.ID}
	for _, w := range b.Witnesses(i.TypeParams) {
		subject.Args = append(subject.Args, rules.EmitRef(w))
	}
	return subject
}

// suited reports whether the declaration id has the directive suite.
func suited(m sdk.SubjectMatch, id *symbol.Identity) bool {
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

// unimportable reports [Unimportable] at the declaration named decl where a
// parameter or a result of a callable of methods refers to a declaration that
// the code of the package of the file out cannot refer to, and returns true.
// It checks the package that the plan previews for out. Where the plan
// previews no package, the layout checks each reference itself, and
// unimportable returns false. It then checks the sample of each parameter, as
// [unreferable] does. A function is its own callable, and its package can
// refer to each type of its signature, so the suite of a function needs the
// check of the samples alone.
func unimportable(m sdk.SubjectMatch, b *rules.Bound, out *sdk.Out, decl string, methods []*node.Method) bool {
	pkg, placed := out.Package()
	if placed != nil {
		return false
	}
	for _, method := range methods {
		refs := make([]*node.TypeRef, 0, len(method.Params)+len(method.Returns))
		for _, p := range method.Params {
			refs = append(refs, p.Type)
		}
		for _, r := range method.Returns {
			refs = append(refs, r.Type)
		}
		for _, ref := range refs {
			if err := b.Importable(pkg, ref); err != nil {
				m.Errorf(Unimportable, "the callable %s.%s refers to %s, which the file of the suite in the "+
					"package %s cannot import, so %s has no suite: %v",
					decl, method.Name, ref.Spelling, pkg.Package, decl, err)
				return true
			}
		}
	}
	return unreferable(m, b, out, decl, methods)
}

// unreferable reports [Unimportable] at the declaration named decl where the
// sample that a parameter of a callable of methods states on itself refers to
// a declaration that the code of the package of the file out cannot refer to,
// and returns true. It checks the package that the plan previews for out.
// Where the plan previews no package, the layout checks each reference
// itself, and unreferable returns false.
func unreferable(m sdk.SubjectMatch, b *rules.Bound, out *sdk.Out, decl string, methods []*node.Method) bool {
	pkg, placed := out.Package()
	if placed != nil {
		return false
	}
	for _, method := range methods {
		for _, p := range method.Params {
			// A parameter without a sample has the zero value, which refers to
			// no declaration.
			if err := b.ImportableValue(pkg, sampleOf(m, b, p, p.Type)); err != nil {
				text, _ := sdk.FactOf(m, p.ID, meta.Named[string](meta.SampleKey))
				m.Errorf(Unimportable, "the parameter %s of %s has the sample %s, which the file of the suite in "+
					"the package %s cannot refer to, so %s has no suite: %v",
					p.Name, method.Name, text, pkg.Package, decl, err)
				return true
			}
		}
	}
	return false
}

// namesOf returns the names of the declarations of the suite of the
// declaration named decl, without the interface of the methods and the
// function type. Each name is in the neutral convention of an emitted name,
// and the backend of the plan writes it in the case of its target.
func namesOf(e *sdk.Emitter, decl string) Names {
	suite := e.JoinName(suiteWord, decl)
	run := e.JoinName(suite, runWord)
	bench := e.JoinName(suite, benchWord)
	return Names{
		Checks: e.JoinName(checksWord, decl), Suite: suite, Run: run, RunFromHandle: e.JoinName(fromHandleWord, run),
		Fuzz: e.JoinName(suite, fuzzWord), Budgets: e.JoinName(budgetsWord, e.JoinName(decl, runWord)),
		Bench: bench, BenchParallel: e.JoinName(parallelWord, bench), Handle: e.JoinName(handleWord, decl),
		Check: e.JoinName(checkWord, decl), Option: e.JoinName(optionWord, decl),
		Subjects: e.JoinName(subjectsWord, decl),
	}
}

// callableOf returns the callable of the method named method of the
// declaration named decl.
func callableOf(e *sdk.Emitter, decl, method string) Callable {
	return Callable{Name: decl + "." + method, Field: method, Type: e.JoinName(method, decl)}
}
