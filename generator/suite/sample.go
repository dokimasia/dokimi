// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"reflect"
	"slices"
	"strconv"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/position"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
)

// Builder is the builder of the mixin sample of a callable: the function Func
// that returns the inputs of each derived call of the callable Callable, one
// value of each input in order. Subject reports a builder that takes the
// subject of the call, and a builder without it takes nothing.
type Builder struct {
	Func     *emit.TypeRef
	Callable string
	Subject  bool
}

// Built is the value at the position Index of the Inputs values that Builder
// returns, of the type Type. Name is the name under which a closure declares
// the value. A value without a name is an expression that calls the builder
// on the subject of the handle and returns the value.
type Built struct {
	Builder *Builder
	Type    *emit.TypeRef
	Name    string
	Index   int
	Inputs  int
}

// Values returns the left side of the assignment of the values of the builder
// of b in an expression without a name. The value of b takes the name v, and
// each other value takes the blank name, as in _, v.
func (b *Built) Values() string {
	names := make([]string, b.Inputs)
	for i := range names {
		names[i] = "_"
	}
	names[b.Index] = "v"
	return strings.Join(names, ", ")
}

// Build is one call of a builder that a closure makes before its call. Names
// are the names of the values that the closure takes, in the order of the
// values, with the blank name for a value that it does not take.
type Build struct {
	Builder *Builder
	Names   []string
}

// builtWord opens the name of each value of a builder that a closure
// declares, as builtSeal0.
const builtWord = "built"

// builderFact is the key of the fact of the param builder of the mixin
// sample, which refers to the builder of a callable.
var builderFact = func() meta.Key[symbol.Identity] {
	s, _ := shape.SpecOf(shape.Sample)
	//dokimi:mutate-skip ror-true: the mixin sample has the param builder alone
	at := slices.IndexFunc(s.Params, func(p shape.Param) bool { return p.Key == shape.SampleBuilder })
	return meta.Named[symbol.Identity](s.Params[at].Fact)
}()

// builderOf returns the builder of the mixin sample of the callable id named
// callable, and nil for a callable without one.
func (k *caller) builderOf(id *symbol.Identity, callable string) *Builder {
	fn, stated := sdk.FactOf(k.m, *id, builderFact)
	if !stated {
		return nil
	}
	found, _ := k.m.Reader().Lookup(fn)
	view, _ := k.b.CallableOf(found)
	return &Builder{
		Func: &emit.TypeRef{Spelling: fn.Name, Target: fn}, Callable: callable, Subject: len(view.Params) == 1,
	}
}

// builtOf returns the value of the i-th input of the callable of f that its
// builder returns, and false for a callable without a builder. The derivation
// takes the value before [caller.builders] checks the builder, so the count
// of the values is the count of the inputs of the callable.
func (k *caller) builtOf(f *derive.Facts, i int) (Expr, bool) {
	b := k.builderOf(&f.ID, f.Callable)
	if b == nil {
		return Expr{}, false
	}
	return Expr{Built: &Built{Builder: b, Type: k.typeOf(k.inputRef(f, i)), Index: i, Inputs: len(f.Inputs)}}, true
}

// declareBuilt makes c call each builder whose values its arguments and its
// assignments take, once before its call, and names each value that it
// takes after the callable and the position of the value, as builtSeal0.
func declareBuilt(c *Closure) {
	var visit func(e *Expr)
	visit = func(e *Expr) {
		switch {
		case e.Built != nil:
			// The closure gives a name to a copy, so an expression that shares
			// the value outside the closure keeps it without a name.
			b := new(*e.Built)
			e.Built = b
			at := slices.IndexFunc(c.Builds, func(o Build) bool { return o.Builder.Callable == b.Builder.Callable })
			if at < 0 {
				names := make([]string, b.Inputs)
				for i := range names {
					names[i] = "_"
				}
				c.Builds = append(c.Builds, Build{Builder: b.Builder, Names: names})
				at = len(c.Builds) - 1
			}
			b.Name = builtWord + naming.Pascal(b.Builder.Callable) + strconv.Itoa(b.Index)
			c.Builds[at].Names[b.Index] = b.Name
		case e.Of != nil:
			visit(e.Of)
		default:
		}
	}
	for i := range c.Args {
		visit(&c.Args[i])
	}
	for _, lets := range [][]Let{c.Lets, c.Later} {
		for i := range lets {
			visit(&lets[i].Value)
		}
	}
}

// builders reports [BadBuilder] at each callable of the suite s whose builder
// takes another parameter than the subject of the suite, or returns other
// values than one value of the type of each input of the callable in order.
// It reports [Unimportable] at a builder that the file of the suite cannot
// import. It returns false where it reports a builder. A builder takes the
// subject of the suite of an interface alone, because the subject of any
// other suite has a type that the suite declares.
func (k *caller) builders(s *Suite) bool {
	valid := true
	for i := range s.Facts {
		f := &s.Facts[i]
		b := k.builderOf(&f.ID, f.Callable)
		if b == nil {
			continue
		}
		found, _ := k.m.Reader().Lookup(b.Func.Target)
		view, _ := k.b.CallableOf(found)
		callable := identityIn(s.Declaration, s.Func == nil, f.Callable)
		at := k.positionOf(&f.ID)
		params := make([]*node.TypeRef, 0, len(view.Params))
		for _, p := range view.Params {
			params = append(params, p.Ref)
		}
		returns := make([]*node.TypeRef, 0, len(view.Returns))
		for _, r := range view.Returns {
			returns = append(returns, r.Ref)
		}
		inputs := make([]*node.TypeRef, 0, len(f.Inputs))
		for j := range f.Inputs {
			inputs = append(inputs, k.inputRef(f, j))
		}
		// The origin of the suite of a function is the function, which is no
		// type of a parameter.
		takesSubject := len(params) == 1 && s.Methods == nil && params[0].Target == k.origin
		// Where the plan does not preview the package of the file, the layout
		// checks each reference itself.
		pkg, placed := k.e.File().Package()
		var unimportable error
		if placed == nil {
			unimportable = k.b.Importable(pkg, &node.TypeRef{Spelling: b.Func.Spelling, Target: b.Func.Target})
		}
		switch {
		case len(params) > 0 && !takesSubject:
			k.m.ErrorfAt(BadBuilder, at, "the builder %s of %s takes %s, and a builder takes the subject of the "+
				"suite or nothing, so %s has no suite", b.Func.Spelling, callable, spellings(params), s.Declaration)
		case !k.returnsInputs(&view, f):
			k.m.ErrorfAt(BadBuilder, at, "the builder %s of %s returns %s, and a builder returns one value of the "+
				"type of each input of %s in order, %s, so %s has no suite", b.Func.Spelling, callable,
				spellings(returns), f.Callable, spellings(inputs), s.Declaration)
		case unimportable != nil:
			k.m.ErrorfAt(Unimportable, at, "the builder %s of %s is a function that the file of the suite in the "+
				"package %s cannot import, so %s has no suite: %v", b.Func.Spelling, callable, pkg.Package,
				s.Declaration, unimportable)
		default:
			continue
		}
		valid = false
	}
	return valid
}

// positionOf returns the position of the declaration of the callable id.
func (k *caller) positionOf(id *symbol.Identity) position.Pos {
	found, _ := k.m.Reader().Lookup(*id)
	decl, _ := found.(interface{ Position() position.Pos })
	return decl.Position()
}

// returnsInputs reports whether the builder of view returns one value of the
// type of each input of the callable of f, in order.
func (k *caller) returnsInputs(view *rules.Callable, f *derive.Facts) bool {
	if len(view.Returns) != len(f.Inputs) {
		return false
	}
	for i, r := range view.Returns {
		if !reflect.DeepEqual(k.b.TypeOf(k.instance(r.Ref)), f.Inputs[i].Shape) {
			return false
		}
	}
	return true
}

// spellings returns the spellings of refs as Go writes a list of types, as
// (string, int).
func spellings(refs []*node.TypeRef) string {
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		out = append(out, r.Spelling)
	}
	return "(" + strings.Join(out, ", ") + ")"
}
