// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"fmt"
	"strings"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// Call is the call generator of one callable of a suite, which the template
// of the family of the suites writes as a variable.
//
// Var is the name of the variable, as StoreGetCalls, and the Go plan writes
// it as storeGetCalls. Field and Type are the field and the type of the
// callable, as its [Callable] contains them. Field is empty for a function,
// whose generated call calls the subject itself. Params are the parameters of
// the callable in order. Assign is the left side of the assignment of the
// results of the callable, as "_, err :=", and empty for a callable whose
// results the generated call does not assign. Returns is the result of the
// generated call after the assignment, err or nil. It is empty where the
// callable returns an error alone, because the generated call then returns
// the result of the callable.
//
// Budget is the budget of the callable, and Concurrent reports that the
// callable has the mixin concurrent or concurrent-readers. Over is the
// signature of the callable in the planted defect of a ceiling, whose
// parameters have no names and whose results have the names r0, r1 and so
// on.
type Call struct {
	Over       *emit.Method
	Var        string
	Field      string
	Type       string
	Assign     string
	Returns    string
	Params     []Param
	Budget     Budget
	Concurrent bool
}

// Param is one parameter of a call. Name is the name of the value that the
// call generator draws for the parameter, and Label is the label of the draw:
// the name of the parameter in kebab case, as user-id. Type is the type of
// the drawn value, or the type of the elements of a channel.
//
// Context reports a parameter that takes the context of the call, which has
// no value to draw. Variadic reports a parameter that takes a drawn list of
// its type, and Address a pointer that takes the address of the drawn value
// of its element type, because the callable states no mixin nil-safe.
// Receive reports a channel that the callable receives from. Its value is a
// drawn list of its elements, as inValues for the parameter in, and each call
// passes a new closed channel of the list.
//
// Funcs are the generated functions of a parameter of a function type or of
// an interface: the one function of a function type, or the function of each
// method of the interface. Input is the name of the implementation of the
// interface that the suite declares, as StoreSourceInput, and empty for any
// other parameter. Omitted reports a variadic parameter of functions or of
// such interfaces, for which the call passes no values.
type Param struct {
	Type     *emit.TypeRef
	Name     string
	Label    string
	Input    string
	Funcs    []Func
	Context  bool
	Variadic bool
	Address  bool
	Receive  bool
	Omitted  bool
}

// The names that the generated call of a callable declares for itself: the
// case of the call generator, the context and the subject of the call, the
// error of the callable, and the packages of the call generator. A parameter
// with one of these names takes another name.
var reserved = []string{"c", "ctx", "subject", "err", "prop", "context", "atomic", "conform"}

// predeclared are the types that Go predeclares. A parameter with the name of
// one of them takes another name, because a later draw of the call generator
// can spell the type.
var predeclared = []string{
	"any", "bool", "byte", "comparable", "complex64", "complex128", "error", "float32", "float64", "int", "int8",
	"int16", "int32", "int64", "rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
}

// caller writes the calls of the callables of the declaration origin, named
// decl. b contains the rules of the declaration. typeParams are the type
// parameters of a generic interface, and witnesses the types that its suite
// takes in their place. inputs are the implementations of interfaces that the
// calls take, and roots the types of the drawn values of the calls, which
// the overrides of the suite walk.
type caller struct {
	m          sdk.SubjectMatch
	g          *generator
	e          *sdk.Emitter
	b          rules.Bound
	origin     symbol.Identity
	decl       string
	typeParams []*node.TypeParam
	witnesses  []*node.TypeRef
	inputs     []Input
	roots      []*node.TypeRef
}

// call returns the call of the callable c. sym is the method or the function
// of the callable, and id its identity. A method of a declaration calls the
// method of the subject, and a function calls the subject itself.
func (k *caller) call(c *Callable, sym symbol.Symbol, id *symbol.Identity) Call {
	view, _ := k.b.CallableOf(sym)
	name := k.e.JoinName(callsWord, c.Type)
	if c.Field == "" {
		name = k.e.JoinName(callsWord, view.Name)
	}
	assign, returns := results(&view)
	return Call{
		Var: name, Field: c.Field, Type: c.Type, Params: k.params(&view, mixed(k.m, id, shape.NilSafe)),
		Assign: assign, Returns: returns, Budget: k.g.budgetFact(k.m, id),
		Concurrent: mixed(k.m, id, shape.Concurrent) || mixed(k.m, id, shape.ConcurrentReaders),
		Over:       k.over(&view, id),
	}
}

// params returns the parameters of the callable of view. Under the mixin
// nil-safe a pointer takes a pointer that can be nil, and otherwise the
// address of a value of its element type. A parameter of a function type
// takes a generated function, and a parameter of an interface outside the
// standard library an implementation whose methods are generated functions.
// A function type or an interface whose generated code refers to a type that
// the file of the suite cannot import takes the generator of its
// registration. A channel that the callable receives from takes a closed
// channel of a drawn list. A variadic parameter of functions or of interfaces
// outside the standard library does not take a value.
//
// Each parameter takes its own name where the name is free, and otherwise the
// first free name of [naming.Unique]. A name is free where [caller.taken]
// does not contain it. A parameter without a name takes the name and the
// label a and its position, as a0.
func (k *caller) params(view *rules.Callable, nilSafe bool) []Param {
	refs := make([]*node.TypeRef, 0, len(view.Params))
	for _, p := range view.Params {
		refs = append(refs, p.Ref)
	}
	taken := k.taken(refs...)
	out := make([]Param, 0, len(view.Params))
	for i, p := range view.Params {
		if p.Role == rules.ParamContext {
			out = append(out, Param{Context: true})
			continue
		}
		base, label := p.Name, naming.Kebab(p.Name)
		if base == "" || base == "_" {
			base = fmt.Sprintf("a%d", i)
			label = base
		}
		name := naming.Unique(taken, base)
		taken[name] = true
		ref := k.instance(p.Ref)
		param := Param{Type: rules.EmitRef(ref), Name: name, Label: label, Variadic: p.Variadic}
		under := k.underlying(ref)
		iface, implemented := k.implemented(ref)
		switch {
		case p.Variadic && (under.Form == symbol.FormFunc || implemented):
			param.Omitted = true
		case p.Variadic:
			k.roots = append(k.roots, ref)
		case under.Form == symbol.FormFunc && k.importable(under):
			param.Funcs = []Func{k.funcOf(under, name, label, taken)}
		case strings.HasPrefix(under.Spelling, receiveSpelling):
			param.Type, param.Receive = rules.EmitRef(under.Elems[0]), true
			param.Name = naming.Unique(taken, base+valuesWord)
			taken[param.Name] = true
			k.roots = append(k.roots, under.Elems[0])
		case ref.Form == symbol.FormOptional && !nilSafe && !pointerShapes[ref.Elems[0].Target.Package+"."+
			ref.Elems[0].Target.Name]:
			param.Type, param.Address = rules.EmitRef(ref.Elems[0]), true
			k.roots = append(k.roots, ref.Elems[0])
		case implemented && k.implementable(ref, iface):
			param.Input, param.Funcs = k.inputOf(ref, iface, name, label, taken)
		default:
			k.roots = append(k.roots, ref)
		}
		out = append(out, param)
	}
	return out
}

// taken returns the names that a generator of the suite does not declare for
// a value of its own: the names that the generated call declares, the
// predeclared types, and the name of each package that the generator can
// spell. These packages are the package of the declaration and the packages
// of the named types of refs. A local variable with such a name would hide
// the package from a later spelling of a type.
func (k *caller) taken(refs ...*node.TypeRef) map[string]bool {
	taken := map[string]bool{}
	for _, names := range [][]string{reserved, predeclared} {
		for _, name := range names {
			taken[name] = true
		}
	}
	if pkg, found := k.m.Reader().PackageOf(k.origin); found {
		taken[pkg.Name] = true
	}
	for _, ref := range refs {
		k.packages(ref, taken)
	}
	return taken
}

// importable reports whether the file of the suite can refer to each type of
// refs. It reports true where the plan does not preview the package of the
// file, because the layout then checks each reference itself.
func (k *caller) importable(refs ...*node.TypeRef) bool {
	pkg, placed := k.e.File().Package()
	if placed != nil {
		return true
	}
	for _, ref := range refs {
		if k.b.Importable(pkg, ref) != nil {
			return false
		}
	}
	return true
}

// implementable reports whether the file of the suite can refer to the
// interface iface at the type ref, and to each type of the signature of each
// of its exported methods. The implementation of the interface spells each of
// these types.
func (k *caller) implementable(ref *node.TypeRef, iface *node.Interface) bool {
	b := k.b
	refs := []*node.TypeRef{ref}
	for _, method := range methodsOf(&b, iface) {
		for _, p := range method.Params {
			refs = append(refs, p.Type)
		}
		for _, r := range method.Returns {
			refs = append(refs, r.Type)
		}
	}
	return k.importable(refs...)
}

// receiveSpelling opens the spelling of a channel that the callable receives
// from, and the spelling of no other type.
const receiveSpelling = "<-chan"

// pointerShapes are the types, by their package and name, whose pointer the
// shape reader of the property engine reads as a value of its own shape: a
// *big.Int as an integer, a *big.Rat as a decimal and a *time.Location as a
// zone. The reader reads no value of such a type, so a pointer to one of them
// takes the drawn pointer itself in place of the address of a drawn value.
var pointerShapes = map[string]bool{"math/big.Int": true, "math/big.Rat": true, "time.Location": true}

// packages adds to taken the name of the package of each named type of ref:
// the qualifier of its spelling, such as context for context.Context, and
// the name that its package declares. A local variable with such a name
// would hide the package from a later spelling of the type.
func (k *caller) packages(ref *node.TypeRef, taken map[string]bool) {
	//dokimi:mutate-skip ror-true: the spelling of a structural form opens with a token that is no identifier, such as []time, so its qualifier is no name of a parameter
	if ref.Form == symbol.FormNamed {
		if qualifier, _, qualified := strings.Cut(ref.Spelling, "."); qualified {
			taken[qualifier] = true
		}
		if pkg, found := k.m.Reader().PackageOf(ref.Target); found {
			taken[pkg.Name] = true
		}
	}
	for _, elem := range ref.Elems {
		k.packages(elem, taken)
	}
	for _, arg := range ref.Args {
		k.packages(arg, taken)
	}
}

// results returns the left side of the assignment of the results of the
// callable of view, and the result of the generated call after it. Each
// result takes the blank name, and the error takes err.
func results(view *rules.Callable) (assign, returns string) {
	if len(view.Returns) == 0 {
		return "", "nil"
	}
	names := make([]string, 0, len(view.Returns))
	errs := false
	for _, r := range view.Returns {
		if r.Role == rules.ReturnError {
			names, errs = append(names, "err"), true
			continue
		}
		names = append(names, "_")
	}
	switch {
	case errs && len(names) == 1:
		return "", ""
	case errs:
		return strings.Join(names, ", ") + " :=", "err"
	}
	return strings.Join(names, ", ") + " =", "nil"
}

// over returns the signature of the callable of view in the planted defect of
// a ceiling: each parameter without a name, and each result with the name r
// and its position, so the defect returns its zero results with a bare
// return. id is the identity of the callable.
func (k *caller) over(view *rules.Callable, id *symbol.Identity) *emit.Method {
	m := &emit.Method{Origin: *id, Name: view.Name, Visibility: symbol.VisibilityPublic}
	for _, p := range view.Params {
		param := &emit.Param{Type: k.typeOf(p.Ref)}
		if p.Variadic {
			param.Variadic = symbol.VariadicPositional
		}
		m.Params.Append(param)
	}
	for i, r := range view.Returns {
		m.Returns.Append(&emit.Return{Name: fmt.Sprintf("r%d", i), Type: k.typeOf(r.Ref)})
	}
	return m
}

// typeOf returns the type ref at the witnesses of the type parameters of the
// declaration.
func (k *caller) typeOf(ref *node.TypeRef) *emit.TypeRef {
	return rules.EmitRef(k.instance(ref))
}

// instance returns ref with the witnesses in place of the type parameters of
// the declaration, and ref itself for a declaration without type parameters.
func (k *caller) instance(ref *node.TypeRef) *node.TypeRef {
	return rules.Substitute(ref, k.typeParams, k.witnesses)
}

// aliasDepth is the most aliases and defined types that [caller.underlying]
// follows, so a cycle of a broken source ends.
const aliasDepth = 8

// underlying returns the type that ref refers to through its aliases and its
// defined types, such as the function type of a defined function type, at the
// type arguments of ref. A structural form has no target, so the lookup
// finds no declaration of it.
func (k *caller) underlying(ref *node.TypeRef) *node.TypeRef {
	for range aliasDepth {
		decl, _ := k.m.Reader().Lookup(ref.Target)
		alias, is := decl.(*node.Alias)
		if !is {
			return ref
		}
		ref = rules.Substitute(alias.Target, alias.TypeParams, ref.Args)
	}
	return ref
}

// implemented returns the interface that the type ref refers to, where the
// interface is outside the standard library. A call generator takes an
// implementation of such an interface whose methods return generated
// results. An interface of the standard library, such as io.Reader, takes
// the generator of its registration. A structural form and a builtin have no
// target, and [golang.Standard] reports the empty package of such a target as
// a package of the standard library.
func (k *caller) implemented(ref *node.TypeRef) (*node.Interface, bool) {
	if golang.Standard(ref.Target.Package) {
		return nil, false
	}
	decl, _ := k.m.Reader().Lookup(ref.Target)
	iface, is := decl.(*node.Interface)
	return iface, is
}

// mixed reports whether the shape catalog stamps the mixin n on the callable
// id.
func mixed(m sdk.Matcher, id *symbol.Identity, n shape.Mixin) bool {
	spec, _ := shape.SpecOf(n)
	stamped, _ := sdk.FactOf(m, *id, meta.Named[bool](spec.Key))
	return stamped
}
