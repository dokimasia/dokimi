// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"strconv"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/fault"
)

// Func is a function that a call generator writes for a parameter of a
// function type, or for a method of the implementation of an interface.
// Signature is the signature of the function, whose parameters have no names
// and whose results have the names r0, r1 and so on. Field is the field of
// the implementation that contains the function of a method, as OnNext, and
// empty for a parameter. Results are the lists of the results that the case
// draws, one for each result, all of one length. Index is the name of the
// counter of the calls of a function with results, and empty for a function
// without results.
//
// The function returns the results of the lists in turn. After the last of
// them it returns the zero results, with the end error of each error result,
// so a caller that reads until an end stops.
type Func struct {
	Signature *emit.Method
	Field     string
	Index     string
	Results   []Result
}

// Result is the list of the values of one result of a generated function.
// Name is the name of the list, and Label the label of its draw. Return is
// the name of the result in the signature of the function, as r0. Type is
// the type of a value. Error reports an error result, whose values are nil
// and the errors of Errors: the declared errors of its method. The end error
// of the result is the first of Errors, and conform.ErrGenerated for an error
// result without declared errors.
type Result struct {
	Type   *emit.TypeRef
	Name   string
	Label  string
	Return string
	Errors []emit.Value
	Error  bool
}

// Input is the implementation of an interface that the suite declares for
// the parameters of the interface. Name is the name of the implementation,
// as StoreSourceInput, and the Go plan writes it as storeSourceInput. Type is
// the interface at the type arguments of the parameters. Methods are the
// exported methods of the interface, in the order of its member walk. The
// implementation embeds the interface, so it has each other method as well.
// Embedded is the alias of the interface that the implementation embeds, as
// StoreSourceInputEmbedded. The embedded field has the name of the alias, so a
// method with the name of the interface is a valid method of the
// implementation. Origin is the identity of the first callable that takes the
// implementation, the origin of the test of the implementations.
type Input struct {
	Type     *emit.TypeRef
	Name     string
	Embedded string
	Origin   symbol.Identity
	Methods  []InputMethod
}

// InputMethod is one method of an [Input]. Method is the signature of the
// method, whose parameters have the names a0, a1 and so on. The method
// passes them to the function of the field Field, as OnNext, and returns its
// results. Probe is the same signature with the results r0, r1 and so on,
// which the function of a test of the method declares.
type InputMethod struct {
	Method *emit.Method
	Probe  *emit.Method
	Field  string
}

// The words that the names of the parts of a generated function and of an
// implementation of an interface join: the counter of the calls of the
// function, the list of each of its results, the implementation, and the
// field of the function of each method.
const (
	atWord      = "At"
	resultsWord = "Results"
	inputWord   = "Input"
	onWord      = "On"
)

// sentinelParam is the key of the param of a mixin that names the error value
// of the callable, such as the sentinel of the mixin not-found.
const sentinelParam = "sentinel"

// funcOf returns the generated function of a value of the function type ref.
// name and label are the name and the label of the value. taken contains the
// names of the call generator, and funcOf adds the names of the lists and of
// the counter of the function. The results of the function type are roots of
// the overrides of the suite.
func (k *caller) funcOf(ref *node.TypeRef, name, label string, taken map[string]bool) Func {
	f := Func{Signature: &emit.Method{}}
	for _, p := range ref.Elems[:ref.Split] {
		param := &emit.Param{Type: rules.EmitRef(p)}
		if strings.HasPrefix(p.Spelling, variadicSpelling) {
			param.Type, param.Variadic = rules.EmitRef(p.Elems[0]), symbol.VariadicPositional
		}
		f.Signature.Params.Append(param)
	}
	returns := ref.Elems[ref.Split:]
	for i, r := range returns {
		res := k.result(r, name, label, i, len(returns), nil, taken)
		f.Signature.Returns.Append(&emit.Return{Name: res.Return, Type: res.Type})
		f.Results = append(f.Results, res)
	}
	f.Index = k.index(&f, name, taken)
	return f
}

// variadicSpelling opens the spelling of the last parameter of a variadic
// function type, and the spelling of no other parameter.
const variadicSpelling = "..."

// result returns the list of the result at the position i of n results of
// the generated function name with the label label. A function of one result
// names its list after itself, and a function of more results adds the
// position to the name and to the label of each list. errs are the declared
// errors of an error result. A result of another type is a root of the
// overrides of the suite.
func (k *caller) result(
	ref *node.TypeRef, name, label string, i, n int, errs []emit.Value, taken map[string]bool,
) Result {
	r := Result{
		Type: rules.EmitRef(ref), Label: label, Return: "r" + strconv.Itoa(i),
		Error: k.b.TypeOf(ref).Ref == rules.WellKnownError,
	}
	base := name + resultsWord
	if n > 1 {
		base += strconv.Itoa(i)
		r.Label += "-" + strconv.Itoa(i)
	}
	r.Name = naming.Unique(taken, base)
	taken[r.Name] = true
	if r.Error {
		r.Errors = errs
	} else {
		k.roots = append(k.roots, ref)
	}
	return r
}

// index returns the name of the counter of the calls of the function f
// named name, which taken does not contain, and adds it to taken. A function
// without results needs no counter, so its name is empty.
func (*caller) index(f *Func, name string, taken map[string]bool) string {
	if len(f.Results) == 0 {
		return ""
	}
	index := naming.Unique(taken, name+atWord)
	taken[index] = true
	return index
}

// inputOf returns the name of the implementation of the interface iface at
// the type ref, and the generated function of each method of the interface
// for a parameter of the callable id. name is the name of the parameter, and
// label is its label. The suite declares one implementation of each interface
// at its type arguments.
// The function of a method takes the declared errors of the method for its
// error result.
func (k *caller) inputOf(
	ref *node.TypeRef, iface *node.Interface, id *symbol.Identity, name, label string, taken map[string]bool,
) (string, []Func) {
	b := k.b
	methods := methodsOf(&b, iface)
	in := k.input(ref, iface, id, methods)
	funcs := make([]Func, 0, len(methods))
	for i, method := range methods {
		view, _ := b.CallableOf(method)
		fn := name + naming.Pascal(method.Name)
		f := Func{Signature: &emit.Method{}, Field: in.Methods[i].Field}
		for _, p := range view.Params {
			param := &emit.Param{Type: rules.EmitRef(rules.Substitute(p.Ref, iface.TypeParams, ref.Args))}
			if p.Variadic {
				param.Variadic = symbol.VariadicPositional
			}
			f.Signature.Params.Append(param)
		}
		errs := k.errorsOf(&method.ID)
		for j, r := range view.Returns {
			res := k.result(rules.Substitute(r.Ref, iface.TypeParams, ref.Args), fn,
				label+"-"+naming.Kebab(method.Name), j, len(view.Returns), errs, taken)
			f.Signature.Returns.Append(&emit.Return{Name: res.Return, Type: res.Type})
			f.Results = append(f.Results, res)
		}
		f.Index = k.index(&f, fn, taken)
		funcs = append(funcs, f)
	}
	return in.Name, funcs
}

// input returns the implementation of the interface iface at the type ref,
// whose methods are the methods of the interface. It declares the
// implementation where the suite has none of the interface at the type
// arguments of ref, with the callable id as its origin. The name of a second
// implementation of one interface ends with a number.
func (k *caller) input(
	ref *node.TypeRef, iface *node.Interface, id *symbol.Identity, methods []*node.Method,
) *Input {
	typ := rules.EmitRef(ref)
	for i := range k.inputs {
		if k.inputs[i].Type.Spelling == typ.Spelling && spelledArgs(k.inputs[i].Type) == spelledArgs(typ) {
			return &k.inputs[i]
		}
	}
	taken := make(map[string]bool, len(k.inputs))
	for i := range k.inputs {
		taken[k.inputs[i].Name] = true
	}
	in := Input{Type: typ, Name: naming.Unique(taken, k.e.JoinName(inputWord, k.e.JoinName(iface.Name, k.decl)))}
	in.Embedded, in.Origin = k.e.JoinName(embeddedWord, in.Name), *id
	for _, method := range methods {
		view, _ := k.b.CallableOf(method)
		m := &emit.Method{Origin: method.ID, Name: method.Name, Visibility: symbol.VisibilityPublic}
		for i, p := range view.Params {
			param := &emit.Param{
				Name: "a" + strconv.Itoa(i), Type: rules.EmitRef(rules.Substitute(p.Ref, iface.TypeParams, ref.Args)),
			}
			if p.Variadic {
				param.Variadic = symbol.VariadicPositional
			}
			m.Params.Append(param)
		}
		probe := &emit.Method{
			Origin: method.ID, Name: method.Name, Visibility: symbol.VisibilityPublic, Params: m.Params,
		}
		for i, r := range view.Returns {
			typ := rules.EmitRef(rules.Substitute(r.Ref, iface.TypeParams, ref.Args))
			m.Returns.Append(&emit.Return{Type: typ})
			probe.Returns.Append(&emit.Return{Name: "r" + strconv.Itoa(i), Type: typ})
		}
		in.Methods = append(in.Methods, InputMethod{Method: m, Probe: probe, Field: k.e.JoinName(method.Name, onWord)})
	}
	k.inputs = append(k.inputs, in)
	return &k.inputs[len(k.inputs)-1]
}

// spelledArgs returns the spellings of the type arguments of ref, joined by
// commas.
func spelledArgs(ref *emit.TypeRef) string {
	args := make([]string, 0, len(ref.Args))
	for _, a := range ref.Args {
		args = append(args, a.Spelling)
	}
	return strings.Join(args, ",")
}

// errorsOf returns the declared errors of the method id: the error values
// that its directive fault lists, and otherwise the error value of each
// mixin of the method that names one, such as the sentinel of the mixin
// not-found. It returns nil for a method without declared errors.
func (k *caller) errorsOf(id *symbol.Identity) []emit.Value {
	var out []emit.Value
	listed, _ := sdk.FactOf(k.m, *id, meta.Named[[]symbol.Identity](fault.KeyErrors))
	for _, e := range listed {
		out = append(out, emit.Name(e))
	}
	if len(out) > 0 {
		return out
	}
	specs := shape.Specs()
	for i := range specs {
		spec := &specs[i]
		if spec.Form != shape.FormMixin {
			continue
		}
		for j := range spec.Params {
			p := &spec.Params[j]
			if p.Key != sentinelParam {
				continue
			}
			if sentinel, stated := sdk.FactOf(k.m, *id, meta.Named[symbol.Identity](p.Fact)); stated {
				out = append(out, emit.Name(sentinel))
			}
		}
	}
	return out
}
