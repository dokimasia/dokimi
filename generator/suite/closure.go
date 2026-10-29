// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"slices"
	"strconv"

	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
)

// The names that a closure declares: the value that a write writes, the
// value of a round trip and the delay of a registration, the parameter of a
// closure that takes a value, the handle of a call, the hook of a registrar,
// the body of a scope, the error of a call and the context of a call.
const (
	writtenName = "w"
	inName      = "in"
	delayName   = "d"
	takenName   = "v"
	handleName  = "handle"
	hookName    = "hook"
	bodyName    = "body"
	errName     = "err"
	ctxName     = "ctx"
)

// nilWord is the word of Go for the nil value.
const nilWord = "nil"

// input is one input of a probe that a closure takes. f is the callable that
// the closure calls, and nil where it calls the method of a handle that
// method names.
// pair is the pair of the check in which f reads or writes, and nil where it
// does neither. with is the value that the first input of the call takes in
// place of its sample, and nil for the sample. stamps is the field of the
// stamp of an entry, and handle the type of the handle of a call. changed is
// the field of the value of a guarded write that the predicate judges.
type input struct {
	sl      *slot
	f       *derive.Facts
	pair    *derive.Pair
	with    *Expr
	handle  *node.TypeRef
	method  string
	stamps  string
	changed string
}

// closure returns the closure of the input in of a probe.
func (k *caller) closure(in *input) *Closure {
	sl, f, pair, with := in.sl, in.f, in.pair, in.with
	c := &Closure{}
	if sl.form == formHandle {
		c.Param = &ClosureParam{Name: handleName, Type: k.typeOf(in.handle)}
		if f == nil {
			k.method(c, in.handle, in.method)
			return c
		}
	}
	if f.Declaration != "" {
		c.Method = f.Callable
	}
	c.Args = k.argsOf(f, pair)
	// The forms that write take the write of the pair as f.
	written := pair != nil && pair.Written >= 0
	switch sl.form {
	case formWritten, formRewritten, formActs:
		index := 0
		if written {
			index = pair.Written
		}
		at := inputAt(f, index)
		c.Lets = []Let{{Name: writtenName, Value: k.sampleOf(f, index, 0)}}
		switch sl.form {
		case formRewritten:
			if pair.Changed != "" {
				first, _ := k.fieldSamples(f, index, pair.Changed)
				c.Lets = append(c.Lets, Let{Name: writtenName, Select: pair.Changed, Value: first})
			}
			c.Later = []Let{k.rewrite(f, pair, index)}
			if pair.Changed == "" {
				k.repartition(c, f, pair)
			}
		case formActs:
			c.Later = []Let{{Name: writtenName, Value: k.sampleOf(f, index, 1)}}
		default:
		}
		c.Args[at] = Expr{Name: writtenName}
		if sl.form == formActs {
			k.act(c, f)
			return c
		}
		c.Returns = []*emit.TypeRef{k.typeOf(k.inputRef(f, index))}
		k.returnsBeside(c, f, writtenName)
		return c
	case formAct:
		if written && pair.Changed != "" {
			at := inputAt(f, pair.Written)
			c.Lets = []Let{{Name: writtenName, Value: k.sampleOf(f, pair.Written, 0)}, k.rewrite(f, pair, pair.Written)}
			c.Args[at] = Expr{Name: writtenName}
		}
		k.act(c, f)
		return c
	case formTaking, formInverse, formHook, formScope:
		param := &ClosureParam{Name: takenName}
		switch sl.form {
		case formHook:
			param = &ClosureParam{Name: hookName, Hook: true}
		case formScope:
			param = &ClosureParam{Name: bodyName, Body: true}
		default:
			param.Type = k.typeOf(k.inputRef(f, 0))
		}
		c.Param = param
		c.Args[inputAt(f, 0)] = Expr{Name: param.Name}
		if sl.form == formHook {
			// A registrar takes a function of its own type, such as
			// func(event string), which calls the hook of the probe.
			c.Args[inputAt(f, 0)].Hook = k.typeOf(k.inputRef(f, 0))
		}
		if sl.form == formInverse {
			k.values(c, f)
			return c
		}
		k.act(c, f)
		return c
	case formRoundTrip:
		c.Lets = []Let{{Name: inName, Value: k.sampleOf(f, 0, 0)}}
		c.Args[inputAt(f, 0)] = Expr{Name: inName}
		in := k.typeOf(k.inputRef(f, 0))
		out := k.typeOf(k.resultRef(f, 0))
		c.Returns, c.Pair, c.Error = []*emit.TypeRef{in, out}, true, true
		names := k.assign(c, f, true)
		c.Results = append([]Expr{{Pair: &PairExpr{
			Types:  [2]*emit.TypeRef{in, out},
			Values: [2]Expr{{Name: inName}, {Name: names[0]}},
		}}}, errorOf(c)...)
		return c
	case formDelay:
		c.Lets = []Let{{Name: delayName, Value: k.sampleOf(f, 0, 0)}}
		c.Args[inputAt(f, 0)] = Expr{Name: delayName}
		c.Returns = []*emit.TypeRef{k.typeOf(k.inputRef(f, 0))}
		k.returnsBeside(c, f, delayName)
		return c
	case formStamps:
		c.Stamps = in.stamps
		c.Returns = []*emit.TypeRef{k.fieldRef(k.resultRef(f, 0).Elems[0], in.stamps)}
		c.Error = true
		k.assign(c, f, true)
		c.Results = append([]Expr{{Name: stampsName}}, errorOf(c)...)
		return c
	case formHandle:
		c.Args[inputAt(f, 0)] = Expr{Name: handleName}
		k.act(c, f)
		return c
	case formMatched, formMatch:
		k.matched(c, in)
		return c
	case formLookup:
		c.Flag = true
	default:
	}
	if with != nil {
		c.Args[inputAt(f, 0)] = *with
	}
	k.values(c, f)
	return c
}

// stampsName is the name of the stamps that a closure of formStamps returns.
const stampsName = "stamps"

// matched makes c call the writer or the predicate of a guarded write. The
// value of the first write has the first sample, with the first sample of its
// field in.changed. On every later call of the writer and on each call of the
// predicate, that field has its second sample. The predicate thus judges the
// value of the second write,
// which differs from the first write in that field alone. The predicate takes
// the value under the name of the input of the writer, because the
// derivation binds the two by name.
func (k *caller) matched(c *Closure, in *input) {
	write := in.pair.Write
	index := slices.IndexFunc(write.Inputs, func(v derive.Value) bool {
		return slices.ContainsFunc(v.Fields, func(fd derive.Field) bool { return fd.Name == in.changed })
	})
	first, second := k.fieldSamples(write, index, in.changed)
	changed := Let{Name: writtenName, Select: in.changed, Value: second}
	c.Lets = []Let{{Name: writtenName, Value: k.sampleOf(write, index, 0)}}
	if in.sl.form == formMatched {
		c.Lets = append(c.Lets, Let{Name: writtenName, Select: in.changed, Value: first})
		c.Later = []Let{changed}
		c.Args[inputAt(write, index)] = Expr{Name: writtenName}
		k.act(c, write)
		return
	}
	c.Lets = append(c.Lets, changed)
	at := slices.IndexFunc(in.f.Inputs, func(v derive.Value) bool { return v.Name == write.Inputs[index].Name })
	c.Args[inputAt(in.f, at)] = Expr{Name: writtenName}
	k.values(c, in.f)
}

// method makes c call the method name of the handle of the type ref for its
// error, with the context of the closure and the first sample of each input.
func (k *caller) method(c *Closure, ref *node.TypeRef, name string) {
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	decl, _ := k.m.Reader().Lookup(ref.Target)
	b := k.b
	methods := methodsOf(&b, decl)
	f := k.factsOf(methods[slices.IndexFunc(methods, func(m *node.Method) bool { return m.Name == name })])
	c.Receiver, c.Method = &Expr{Name: handleName}, name
	c.Args = k.argsOf(&f, nil)
	k.act(c, &f)
}

// rewrite returns the assignment of the second write of the written input i
// of the callable of f: the field that the second write of pair changes, from
// the second sample of the field, or the whole second sample where pair names
// no field. Each probe that writes twice has the pair of the write and its
// read.
func (k *caller) rewrite(f *derive.Facts, pair *derive.Pair, i int) Let {
	if pair.Changed == "" {
		return Let{Name: writtenName, Select: "", Value: k.sampleOf(f, i, 1)}
	}
	_, second := k.fieldSamples(f, i, pair.Changed)
	return Let{Name: writtenName, Select: pair.Changed, Value: second}
}

// fieldSamples returns the two samples of the field name of the struct that
// the i-th input of the callable of f takes, or points to, each converted to
// the type of the field. The sample of a struct sets one field alone, so a
// write that changes another field takes its values from the samples of the
// field itself.
func (k *caller) fieldSamples(f *derive.Facts, i int, name string) (first, second Expr) {
	ref := k.inputRef(f, i)
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	sym, _ := k.m.Reader().Lookup(ref.Target)
	st, _ := sym.(*node.Struct)
	// The derivation and the contract if-match name only a field of the
	// struct.
	field := st.Fields[slices.IndexFunc(st.Fields, func(fd *node.Field) bool { return fd.Name == name })]
	v := k.valueOf(&field.ID, field.Name, field.Type)
	t := k.typeOf(k.instance(field.Type))
	return Expr{Value: typed(t, &v.Samples[0].Value)}, Expr{Value: typed(t, &v.Samples[1].Value)}
}

// repartition makes the later calls of c write into another partition: each
// input of the write of pair that the read takes by its name, such as the
// scope and the key of Set(ctx, scope, key, value), takes a variable of its
// own, with the first sample on the first call and the second sample on every
// later call.
func (k *caller) repartition(c *Closure, f *derive.Facts, pair *derive.Pair) {
	for i := range f.Inputs {
		if !slices.ContainsFunc(pair.Read.Inputs, func(r derive.Value) bool { return r.Name == f.Inputs[i].Name }) {
			continue
		}
		name := selectorName + strconv.Itoa(i)
		c.Lets = append(c.Lets, Let{Name: name, Value: k.sampleOf(f, i, 0)})
		c.Later = append(c.Later, Let{Name: name, Value: k.sampleOf(f, i, 1)})
		c.Args[inputAt(f, i)] = Expr{Name: name}
	}
}

// selectorName is the prefix of the variable of an input that selects the
// partition of a write.
const selectorName = "s"

// argsOf returns the arguments of a call of the callable of f: the context of
// the closure for a parameter in the role context, and the first sample of
// each input. The read of pair takes the arguments that the pair binds, and
// the write of pair the first sample of the written input.
func (k *caller) argsOf(f *derive.Facts, pair *derive.Pair) []Expr {
	out := make([]Expr, 0, len(f.Signature.Params))
	if pair != nil && same(pair.Read, f) {
		for _, a := range pair.Arguments {
			switch a.Kind {
			case derive.ArgContext:
				out = append(out, Expr{Word: ctxName})
			case derive.ArgShared:
				i := slices.IndexFunc(pair.Write.Inputs, func(v derive.Value) bool { return v.Name == a.Param })
				out = append(out, k.sampleOf(pair.Write, i, 0))
			case derive.ArgField:
				written := k.sampleOf(pair.Write, pair.Written, 0)
				out = append(out, Expr{Of: &written, Select: a.Field})
			}
		}
		return out
	}
	i := 0
	for _, p := range f.Signature.Params {
		if p.Role == rules.ParamContext {
			out = append(out, Expr{Word: ctxName})
			continue
		}
		out = append(out, k.sampleOf(f, i, 0))
		i++
	}
	return out
}

// sampleOf returns the sample which of the i-th input of the callable of f,
// converted to the type of the input.
func (k *caller) sampleOf(f *derive.Facts, i, which int) Expr {
	return Expr{Value: typed(k.typeOf(k.inputRef(f, i)), &f.Inputs[i].Samples[which].Value)}
}

// inputRef returns the type of the i-th input of the callable of f.
func (*caller) inputRef(f *derive.Facts, i int) *node.TypeRef {
	return f.Signature.Params[inputAt(f, i)].Ref
}

// inputAt returns the position among the parameters of the callable of f of
// its i-th input.
func inputAt(f *derive.Facts, i int) int {
	var inputs []int
	for at, p := range f.Signature.Params {
		if p.Role != rules.ParamContext {
			inputs = append(inputs, at)
		}
	}
	return inputs[i]
}

// act makes c return the error of the call of the callable of f: the call
// itself for a callable that returns an error alone, the assigned error for a
// callable that also returns values, and nil after the call for a callable
// without an error.
func (k *caller) act(c *Closure, f *derive.Facts) {
	c.Error = true
	switch {
	case len(f.Signature.Returns) == 1 && f.Signature.Returns[0].Role == rules.ReturnError:
	case hasError(f):
		k.assign(c, f, false)
		c.Results = []Expr{{Name: errName}}
	default:
		c.Discard = true
		c.Results = []Expr{{Word: nilWord}}
	}
}

// values makes c return the values of the call of the callable of f beside
// its error, or beside its flag where c has Flag: the call itself where the
// results of the callable have the shape of c, one value or a pair of two. A
// presence flag of Go is the second of two results, so a lookup returns the
// results of its call as they are.
func (k *caller) values(c *Closure, f *derive.Facts) {
	n := len(f.Results)
	for i := range n {
		c.Returns = append(c.Returns, k.typeOf(k.resultRef(f, i)))
	}
	c.Pair = n == 2
	c.Error = !c.Flag
	exact := len(f.Signature.Returns) == 2 && f.Signature.Returns[1].Role == rules.ReturnError
	if c.Flag || exact {
		return
	}
	names := k.assign(c, f, true)
	var value Expr
	if c.Pair {
		value = Expr{Pair: &PairExpr{
			Types:  [2]*emit.TypeRef{c.Returns[0], c.Returns[1]},
			Values: [2]Expr{{Name: names[0]}, {Name: names[1]}},
		}}
	} else {
		value = Expr{Name: names[0]}
	}
	c.Results = append([]Expr{value}, errorOf(c)...)
}

// returnsBeside makes c return the variable name beside the error of the
// call of the callable of f, or beside nil for a callable without an error.
func (k *caller) returnsBeside(c *Closure, f *derive.Facts, name string) {
	c.Error = true
	if !hasError(f) {
		c.Discard = true
		c.Results = []Expr{{Name: name}, {Word: nilWord}}
		return
	}
	k.assign(c, f, false)
	c.Results = []Expr{{Name: name}, {Name: errName}}
}

// assign makes c assign the results of the call of the callable of f, and
// returns the names of its values in order. A value takes the name r and its
// position where keep is set, and the error err. A closure that does not
// return a result assigns it to the blank name, the flag among them.
func (*caller) assign(c *Closure, f *derive.Facts, keep bool) []string {
	var values []string
	for i, r := range f.Signature.Returns {
		name := "_"
		switch r.Role {
		case rules.ReturnError:
			if c.Error {
				name = errName
			}
		case rules.ReturnOkBool:
		default:
			if keep {
				name = "r" + strconv.Itoa(i)
				values = append(values, name)
			}
		}
		c.Assign = append(c.Assign, name)
	}
	return values
}

// errorOf returns the error that c returns: the assigned error, or nil for a
// callable without one.
func errorOf(c *Closure) []Expr {
	if slices.Contains(c.Assign, errName) {
		return []Expr{{Name: errName}}
	}
	return []Expr{{Word: nilWord}}
}

// hasError reports whether the callable of f returns an error.
func hasError(f *derive.Facts) bool {
	return slices.ContainsFunc(
		f.Signature.Returns,
		func(r rules.ReturnView) bool { return r.Role == rules.ReturnError },
	)
}
