// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"reflect"
	"slices"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk/node"

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
type Draw struct {
	Param
	Pair bool
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
// implementation of the input itself is built anew for each call instead.
func (d *draws) param(f *derive.Facts, i int) *Param {
	j := slices.IndexFunc(d.facts, func(o derive.Facts) bool { return o.ID == f.ID })
	p := &d.calls[j].Params[inputAt(f, i)]
	if p.Sample != nil || p.Funcs == nil && d.k.consumes(d.k.inputRef(f, i), map[string]bool{}) {
		return nil
	}
	return p
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
	draw := d.draw(&Param{Name: p.Name + naming.Pascal(name), Label: p.Label + "." + naming.Kebab(name), Type: t})
	draw.Pair = true
	return Expr{Drawn: &Drawn{Draw: draw}}, Expr{Drawn: &Drawn{Draw: draw, Second: true}}
}

// draw returns the one draw of the check for the label and the type of the
// parameter p. A new draw takes the name of p after drawnWord, numbered where
// another draw has that name.
func (d *draws) draw(p *Param) *Draw {
	for _, o := range d.list {
		if o.Label == p.Label && reflect.DeepEqual(o.Type, p.Type) {
			return o
		}
	}
	o := &Draw{Param: *p}
	o.Name = naming.Unique(d.taken, drawnWord+naming.Pascal(p.Name))
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
		for _, l := range slices.Concat(c.Lets, c.Later) {
			visit(&l.Value)
		}
	}
	return out
}

// argOf returns the value which of the i-th input of the callable of f: its
// sample in a sample binding, where d is nil, and the value that d draws in a
// generated binding.
func (k *caller) argOf(d *draws, f *derive.Facts, i, which int) Expr {
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
