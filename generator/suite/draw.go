// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"reflect"
	"slices"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
)

// Draw is one value that a generated check draws for the inputs of its calls
// with one key and one type. Param is the parameter of the call generator
// whose rules draw the value, under the name of the draw. Pair reports a draw
// of two distinct values, for a check that passes a second value of the
// input, such as the value of a second write.
//
// A draw of a generated function or implementation draws the lists of its
// results once. It declares a function that builds a new generated function
// or implementation over the lists, so two calls of a probe take two of them.
//
// Generator is the name of the shape definition of the suite whose generator
// draws the value, for a draw whose parameter has a shape file, and empty for
// any other draw.
type Draw struct {
	Generator string
	Param
	Pair bool
}

// Shape is the shape definition of one type whose values the generated
// bindings of a suite draw through a shape file, for a language whose
// property engine reads the type of no value. Name is the name of the
// variable of the definition in the spelling of the language, Type the type
// of its values and File its shape file. Start is the starting value of the
// values, which the generator keeps or changes, and nil for a type without
// one.
type Shape struct {
	Type  *emit.TypeRef
	Start *emit.Value
	Name  string
	File  string
}

// shapeWord is the word that the name of a shape definition joins to the name
// of its declaration and of the draw that first reads it, as StoreIdShape for
// the draw drawnId of the suite of Store.
const shapeWord = "Shape"

// shaped is one draw of a suite whose type has a shape file: its type, the
// starting value and the shape file of its values, the word that names its
// shape definition, and the name of the generator that it reads.
type shaped struct {
	typ       *emit.TypeRef
	start     *emit.Value
	generator *string
	file      string
	word      string
}

// shapesOf records into s a shape definition for each type with a shape file
// that the suite draws, and the name of that definition into each draw. The
// suite draws the inputs of its derived checks, the keys of each pool of its
// generated sequences, and each input of an action. The shape file of a draw
// follows from its type. A definition takes the name of the declaration of the
// suite and of the first draw of its type in the spelling of the part p, as
// sinkPShape for the draw p of the suite of Sink. The word of a draw is the
// name of the draw of a check, the type of the keys of a pool, or the name of
// the parameter of an action. The suites of one module so declare distinct
// names. A definition takes a number where another definition of the suite
// has that name, or where the file imports a type of that name for a draw.
func shapesOf(e *sdk.Emitter, p part, s *Suite) {
	var draws []shaped
	for i := range s.Derived {
		for _, d := range s.Derived[i].Draws {
			if d.Shape != "" {
				draws = append(draws, shaped{
					typ: d.Type, start: d.Start, generator: &d.Generator, file: d.Shape,
					word: strings.TrimPrefix(d.Name, drawnWord),
				})
			}
		}
	}
	if q := s.Sequences; q != nil {
		for i := range q.Pools {
			pool := &q.Pools[i]
			if pool.Shape != "" {
				draws = append(draws, shaped{
					typ: pool.Type, start: pool.Start, generator: &pool.Generator, file: pool.Shape,
					word: pool.Type.Spelling,
				})
			}
		}
		for i := range q.Steps {
			for j := range q.Steps[i].Params {
				sp := &q.Steps[i].Params[j]
				if sp.Shape != "" {
					draws = append(draws, shaped{
						typ: sp.Type, start: sp.Start, generator: &sp.Generator, file: sp.Shape, word: sp.Name,
					})
				}
			}
		}
	}
	taken := map[string]bool{}
	for _, d := range draws {
		imported(d.typ, taken)
	}
	for _, d := range draws {
		at := slices.IndexFunc(s.Shapes, func(sh Shape) bool { return reflect.DeepEqual(sh.Type, d.typ) })
		if at < 0 {
			word := e.JoinName(shapeWord, e.JoinName(d.word, s.Declaration))
			name := naming.Unique(taken, p.spelled(symbol.KindVariable, word))
			taken[name] = true
			s.Shapes = append(s.Shapes, Shape{Type: d.typ, Start: d.start, Name: name, File: d.file})
			at = len(s.Shapes) - 1
		}
		*d.generator = s.Shapes[at].Name
	}
}

// imported adds to taken the name under which the file of a suite refers to
// typ and to each type in it, such as Box for Box<number>. The spelling of a
// structural type, such as string[], is no name of a shape definition, so
// taken may contain it.
func imported(typ *emit.TypeRef, taken map[string]bool) {
	name, _, _ := strings.Cut(typ.Spelling, "<")
	taken[name] = true
	for _, part := range slices.Concat(typ.Elems, typ.Args) {
		imported(part, taken)
	}
}

// choosing reports whether the draw takes a choice of the case: a draw of a
// value, or a generated function or implementation with a result.
func (d *Draw) choosing() bool {
	return d.Funcs == nil || slices.ContainsFunc(d.Funcs, func(f Func) bool { return len(f.Results) > 0 })
}

// Drawn is the value of a draw that an expression refers to: the first value,
// or the second value where Second is set. Address reports an input that
// takes the address of the value, because its call generator draws a pointer
// as the address of a value of its element type.
type Drawn struct {
	Draw    *Draw
	Second  bool
	Address bool
}

// drawnWord is the word that opens the name of each draw. A name with an
// upper-case letter after it is no name that a closure, a package or a
// generated function of the check declares.
const drawnWord = "drawn"

// draws are the draws of one generated check. facts are the facts of the
// callables of the suite and calls their call generators, in one order. taken
// are the names of the draws.
type draws struct {
	k     *caller
	taken map[string]bool
	facts []derive.Facts
	calls []Call
	list  []*Draw
}

// newDraws returns the draws of a generated check of a suite. facts are the
// facts of the callables of the suite, and calls their call generators.
func (k *caller) newDraws(facts []derive.Facts, calls []Call) *draws {
	return &draws{k: k, taken: map[string]bool{}, facts: facts, calls: calls}
}

// param returns the parameter of the call generator of the i-th input of the
// callable of f, a callable of the suite. It returns nil for an input with a
// sample of its own or with a part that a call consumes, such as a variadic
// input of generated functions, because such an input keeps its sample. A
// probe that calls twice passes one drawn value to both calls, so the second
// call would see what the first consumed. A generated function or
// implementation of the input itself is built anew for each call instead,
// and keeps its sample only where its results have such a part.
func (d *draws) param(f *derive.Facts, i int) *Param {
	j := slices.IndexFunc(d.facts, func(o derive.Facts) bool { return o.ID == f.ID })
	p := &d.calls[j].Params[inputAt(f, i)]
	ref := d.k.inputRef(f, i)
	switch {
	case p.Sample != nil:
		return nil
	case p.Funcs != nil:
		if d.k.yieldsConsumable(d.k.instance(ref)) {
			return nil
		}
		return p
	case d.k.consumes(ref, map[string]bool{}):
		return nil
	default:
		return p
	}
}

// yieldsConsumable reports whether a generated function or implementation of
// the type ref returns a value with a part that a call consumes: a result of
// the function type, or of a method of the interface. Each call of a check
// takes a new function or implementation, which returns the values of the
// lists that the case drew once, so two calls share such a value.
func (k *caller) yieldsConsumable(ref *node.TypeRef) bool {
	var results []*node.TypeRef
	if iface, under, isInterface := k.interfaceOf(ref); isInterface {
		b := k.b
		for _, m := range methodsOf(&b, iface) {
			for _, r := range m.Returns {
				results = append(results, rules.Substitute(r.Type, iface.TypeParams, under.Args))
			}
		}
	} else {
		under := k.underlying(ref)
		results = under.Elems[under.Split:]
	}
	met := map[string]bool{}
	return slices.ContainsFunc(results, func(r *node.TypeRef) bool { return k.consumes(r, met) })
}

// consumes reports whether a value of the type ref has a part that a call can
// consume at any depth. Such a part is a generated function or
// implementation, a variant of an interface, a stream, a sink or a reader of
// drawn bytes. met contains the key of each type that the walk visited.
func (k *caller) consumes(ref *node.TypeRef, met map[string]bool) bool {
	key := k.typeKey(ref)
	if met[key] {
		return false
	}
	met[key] = true
	switch src, _ := k.sourceOf(ref); src {
	case sourceGenerated, sourceVariants, sourceStream, sourceSink, sourceBytes:
		return true
	default:
	}
	parts := slices.Concat(ref.Elems, ref.Args)
	for _, p := range partsOf(k.declOf(ref), ref) {
		parts = append(parts, p.ref)
	}
	return slices.ContainsFunc(parts, func(p *node.TypeRef) bool { return k.consumes(p, met) })
}

// value returns the value which of the i-th input of the callable of f, 0 for
// the first value and 1 for the second. An input takes the value that its
// call generator draws, a generated function or implementation among them,
// and a variadic input takes one value of its element type. The second value
// of an input is the second of two distinct draws. An input that keeps its
// sample, and the second value of a generated function or implementation,
// take the sample which.
func (d *draws) value(f *derive.Facts, i, which int) Expr {
	p := d.param(f, i)
	if p == nil || p.Funcs != nil && which == 1 {
		return d.k.sampleOf(f, i, which)
	}
	draw := d.draw(p)
	draw.Pair = draw.Pair || which == 1
	return Expr{Drawn: &Drawn{Draw: draw, Second: which == 1, Address: p.Address}}
}

// field returns the first and the second value of the field name of the
// struct that the i-th input of the callable of f takes or points to. An
// input that the check draws takes two distinct draws of the type of the
// field, under the key of the input and the field. Any other input takes the
// two samples of the field.
func (d *draws) field(f *derive.Facts, i int, name string) (first, second Expr) {
	first, second, t := d.k.fieldSamples(f, i, name)
	p := d.param(f, i)
	if p == nil {
		return first, second
	}
	draw := d.draw(&Param{Name: p.Name + d.k.p.pascal(name), Label: p.Label + "." + naming.Kebab(name), Type: t})
	draw.Pair = true
	return Expr{Drawn: &Drawn{Draw: draw}}, Expr{Drawn: &Drawn{Draw: draw, Second: true}}
}

// draw returns the one draw of the check for the label and the type of the
// parameter p. A new draw takes the name of p after drawnWord in the Pascal
// case of the part, as drawnID in Go and drawnId in TypeScript, numbered
// where another draw has that name.
func (d *draws) draw(p *Param) *Draw {
	for _, o := range d.list {
		if o.Label == p.Label && reflect.DeepEqual(o.Type, p.Type) {
			return o
		}
	}
	o := &Draw{Param: *p}
	o.Name = naming.Unique(d.taken, drawnWord+d.k.p.pascal(p.Name))
	d.taken[o.Name] = true
	d.list = append(d.list, o)
	return o
}

// used returns the draws that args refer to, in the order of their first
// reference. A form that fills an input from the probe, such as the value
// that formTaking takes, replaces the value that the check drew for the
// input, so the draws of a check are those that its arguments refer to.
func used(args []Arg) []*Draw {
	var out []*Draw
	var visit func(e *Expr)
	visit = func(e *Expr) {
		switch {
		case e.Drawn != nil:
			if !slices.Contains(out, e.Drawn.Draw) {
				out = append(out, e.Drawn.Draw)
			}
		case e.Of != nil:
			visit(e.Of)
		default:
		}
	}
	for i := range args {
		if e := args[i].Expr; e != nil {
			visit(e)
		}
		c := args[i].Closure
		if c == nil {
			continue
		}
		for j := range c.Args {
			visit(&c.Args[j])
		}
		lets := slices.Concat(c.Lets, c.Later)
		for j := range lets {
			visit(&lets[j].Value)
		}
	}
	return out
}

// argOf returns the value which of the i-th input of the callable of f: the
// value that its builder returns, its sample in a sample binding, where d is
// nil, and the value that d draws in a generated binding. A builder returns
// one value of each input, so both values of a callable with a builder are
// the value of its builder.
func (k *caller) argOf(d *draws, f *derive.Facts, i, which int) Expr {
	if built, ok := k.builtOf(f, i); ok {
		return built
	}
	if d == nil {
		return k.sampleOf(f, i, which)
	}
	return d.value(f, i, which)
}

// fieldValues returns the first and the second value of the field name of
// the struct that the i-th input of the callable of f takes or points to: its
// two samples in a sample binding, where d is nil, and the values that d
// draws in a generated binding.
func (k *caller) fieldValues(d *draws, f *derive.Facts, i int, name string) (first, second Expr) {
	if d == nil {
		first, second, _ = k.fieldSamples(f, i, name)
		return first, second
	}
	return d.field(f, i, name)
}
