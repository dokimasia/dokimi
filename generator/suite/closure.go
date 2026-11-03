// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"reflect"
	"slices"
	"strconv"

	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
)

// The names that a closure declares: the value that a write writes, the
// request of a call, the value of a round trip and the delay of a
// registration, the parameter of a closure that takes a value, the two
// parameters of a closure that combines two values, the handle of a call, the
// hook of a registrar, the body of a scope, the position of an indexed call,
// the error of a call, the context of a call, the subject of a durable call
// and the presence flag of a read.
const (
	writtenName  = "w"
	requestName  = "req"
	inName       = "in"
	delayName    = "d"
	takenName    = "v"
	firstName    = "x"
	secondName   = "y"
	handleName   = "handle"
	hookName     = "hook"
	bodyName     = "body"
	positionName = "position"
	errName      = "err"
	ctxName      = "ctx"
	subjectName  = "s"
	flagName     = "ok"
)

// nilWord is the word of Go for the nil value, and newSubject the expression
// of a new subject of the factory of the handle.
const (
	nilWord    = "nil"
	newSubject = "h.New()"
)

// intSpelling is the spelling of the type of a position.
const intSpelling = "int"

// The drains of a stream into a list: each element of an iter.Seq, the first
// value of each pair of an iter.Seq2 until the error of its second value, and
// each pair of any other iter.Seq2.
const (
	drainElements   = "elements"
	drainUntilError = "until-error"
	drainPairs      = "pairs"
)

// input is one input of a probe that a closure takes. f is the callable that
// the closure calls, and nil where it calls the method of a handle that
// method names.
// pair is the pair of the check in which f reads or writes, and nil where it
// does neither. with is the value that the first input of the call takes in
// place of its sample, and nil for the sample. stamps is the field of the
// stamp of an entry, and handle the type of the handle of a call. changed is
// the field of the value of a guarded write that the predicate judges, and
// field the field of the request that the sample binding of a call of
// formFielded sets. drawn are the draws of a generated binding, whose values
// the calls take in place of their samples, and nil for a sample binding.
type input struct {
	sl      *slot
	f       *derive.Facts
	pair    *derive.Pair
	with    *Expr
	handle  *node.TypeRef
	drawn   *draws
	method  string
	stamps  string
	changed string
	field   string
}

// closure returns the closure of the input in of a probe. It calls the
// builder of each value of a builder that it takes once, before its call.
func (k *caller) closure(in *input) *Closure {
	c := k.closureOf(in)
	declareBuilt(c)
	return c
}

// closureOf returns the closure of the input in of a probe, by the form of
// the input.
func (k *caller) closureOf(in *input) *Closure {
	sl, f, pair, with, d := in.sl, in.f, in.pair, in.with, in.drawn
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
	if sl.fresh {
		c.Receiver = &Expr{Word: newSubject}
	}
	c.Args = k.argsOf(d, f, pair, 0)
	// The forms that write take the write of the pair as f.
	written := pair != nil && pair.Written >= 0
	switch sl.form {
	case formWritten, formRewritten, formActs:
		index := 0
		if written {
			index = pair.Written
		}
		at := inputAt(f, index)
		c.Lets = []Let{{Name: writtenName, Value: k.argOf(d, f, index, 0)}}
		switch sl.form {
		case formRewritten:
			if pair.Changed != "" {
				first, _ := k.fieldValues(d, f, index, pair.Changed)
				c.Lets = append(c.Lets, Let{Name: writtenName, Select: pair.Changed, Value: first})
			}
			c.Later = []Let{k.rewrite(d, f, pair, index)}
			if pair.Changed == "" {
				k.repartition(c, d, f, pair)
			}
		case formActs:
			c.Later = []Let{{Name: writtenName, Value: k.argOf(d, f, index, 1)}}
		default:
		}
		c.Args[at] = Expr{Name: writtenName}
		if sl.form == formActs {
			k.act(c, f)
			return c
		}
		if len(f.Results) == 1 && reflect.DeepEqual(f.Results[0].Shape, f.Inputs[index].Shape) {
			// A write that returns a value of the written type, such as a
			// write that stamps a revision, returns what it stored.
			k.values(c, f)
			return c
		}
		c.Returns = []*emit.TypeRef{k.typeOf(k.inputRef(f, index))}
		k.returnsBeside(c, f, writtenName)
		return c
	case formAct:
		if written && pair.Changed != "" {
			// The first write sets the changed field to its first value, and
			// every later write to its second value, so a read between two
			// writes sees the field change.
			at := inputAt(f, pair.Written)
			first, _ := k.fieldValues(d, f, pair.Written, pair.Changed)
			c.Lets = []Let{
				{Name: writtenName, Value: k.argOf(d, f, pair.Written, 0)},
				{Name: writtenName, Select: pair.Changed, Value: first},
			}
			c.Later = []Let{k.rewrite(d, f, pair, pair.Written)}
			c.Args[at] = Expr{Name: writtenName}
		}
		k.act(c, f)
		return c
	case formCombine:
		c.Param = &ClosureParam{Name: firstName, Other: secondName, Type: k.typeOf(k.inputRef(f, 0))}
		c.Args[inputAt(f, 0)] = Expr{Name: firstName}
		c.Args[inputAt(f, 1)] = Expr{Name: secondName}
		k.values(c, f)
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
		c.Lets = []Let{{Name: inName, Value: k.argOf(d, f, 0, 0)}}
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
		c.Lets = []Let{{Name: delayName, Value: k.argOf(d, f, 0, 0)}}
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
	case formBatch:
		// The variadic input of a batch read is its last input, and it takes
		// the key of each of the two writes: the two values of the field of
		// the written value that it selects, or the two written values.
		a := &pair.Arguments[len(pair.Arguments)-1]
		var first, second Expr
		if a.Kind == derive.ArgField {
			first, second = k.fieldValues(d, pair.Write, pair.Written, a.Field)
		} else {
			first, second = k.argOf(d, pair.Write, pair.Written, 0), k.argOf(d, pair.Write, pair.Written, 1)
		}
		c.Args = append(c.Args[:len(c.Args)-1], first, second)
		k.values(c, f)
		return c
	case formPosition:
		c.Param = &ClosureParam{Name: positionName, Type: &emit.TypeRef{Spelling: intSpelling}}
		for i := range f.Inputs {
			// A shape has a scalar class only where it is a number.
			shape := &f.Inputs[i].Shape
			if shape.Class != rules.ScalarInt && shape.Class != rules.ScalarUint {
				continue
			}
			position := Expr{Name: positionName}
			if shape.Spelling != intSpelling {
				position.Convert = k.typeOf(k.inputRef(f, i))
			}
			c.Args[inputAt(f, i)] = position
		}
		k.act(c, f)
		return c
	case formLookup:
		c.Flag = true
	case formPull:
		c.Pull = true
	case formComputing:
		k.compute(c, f)
	case formFielded:
		// A failed call returns the field of its request. The sample of a
		// request sets one field alone, so a sample binding sets the field to
		// a sample that is not its zero value. A drawn request draws the field.
		if at := requestAt(f, in.field); at >= 0 && d == nil {
			first, _, _ := k.fieldSamples(f, at, in.field)
			c.Lets = []Let{
				{Name: requestName, Value: k.argOf(nil, f, at, 0)},
				{Name: requestName, Select: in.field, Value: first},
			}
			c.Args[inputAt(f, at)] = Expr{Name: requestName}
		}
	case formChannel:
		// The probe reads the channel, so the closure returns it as a channel
		// that it receives from. Each value of a Go callable precedes its flag
		// and its error, so the channel is the first result.
		c.Channel, c.Error, c.Returns = true, true, []*emit.TypeRef{k.typeOf(k.channelElem(f))}
		k.assign(c, f, false)
		c.Assign[0] = "r0"
		c.Results = append([]Expr{{Name: c.Assign[0]}}, errorOf(c)...)
		return c
	default:
	}
	if with != nil {
		c.Args[inputAt(f, 0)] = *with
	}
	if at := streamAt(f); at >= 0 {
		k.drain(c, f, at)
		return c
	}
	k.values(c, f)
	return c
}

// durableOf returns the calls of a durable probe over the write and the read
// of pair. For each sample, the write and the read call the subject that the
// closures take, with the values of the sample. The read of a callable with a
// presence flag returns fs.ErrNotExist for a key without a value, because the
// probe reads an error where no write made a value durable. d are the draws
// of a generated binding, and nil for a sample binding.
func (k *caller) durableOf(d *draws, pair *derive.Pair) *Durable {
	out := &Durable{Value: k.typeOf(k.resultRef(pair.Read, 0))}
	for which := range out.Writes {
		write := &Closure{
			Receiver: &Expr{Name: subjectName}, Method: pair.Write.Callable,
			Args: k.argsOf(d, pair.Write, pair, which),
		}
		k.act(write, pair.Write)
		read := &Closure{
			Receiver: &Expr{Name: subjectName}, Method: pair.Read.Callable,
			Args: k.argsOf(d, pair.Read, pair, which),
		}
		k.values(read, pair.Read)
		for i, r := range pair.Read.Signature.Returns {
			if r.Role == rules.ReturnOkBool {
				read.Absent, read.Assign[i] = true, flagName
			}
		}
		out.Wants[which] = k.argOf(d, pair.Write, pair.Written, which)
		out.Keys[which] = slices.DeleteFunc(slices.Clone(read.Args), func(e Expr) bool { return e.Word == ctxName })
		declareBuilt(write)
		declareBuilt(read)
		out.Writes[which], out.Reads[which] = write, read
	}
	return out
}

// compute makes c pass a computation to each input of a function type of the
// callable of f, in place of its value. The computation records a call of the
// callable on a new dependency of the injector of the handle and returns its
// zero results, so the injector counts each computation, and Hold keeps each
// computation waiting.
func (k *caller) compute(c *Closure, f *derive.Facts) {
	for i := range f.Inputs {
		under := k.underlying(k.instance(k.inputRef(f, i)))
		if under.Form != symbol.FormFunc {
			continue
		}
		m := &emit.Method{}
		appendParams(m, under)
		for j, r := range under.Elems[under.Split:] {
			m.Returns.Append(&emit.Return{Name: "r" + strconv.Itoa(j), Type: rules.EmitRef(r)})
		}
		c.Args[inputAt(f, i)] = Expr{Compute: m, Name: f.Callable}
	}
}

// streamAt returns the position of the stream among the results of the
// callable of f, and -1 for a callable without one.
func streamAt(f *derive.Facts) int {
	return slices.IndexFunc(f.Signature.Returns, func(r rules.ReturnView) bool { return r.Role == rules.ReturnStream })
}

// elemsOf returns the types of the elements of the stream ref: the value of an
// iter.Seq, or the two values of an iter.Seq2, through the aliases of the
// stream.
func (k *caller) elemsOf(ref *node.TypeRef) []*node.TypeRef {
	return rules.Unaliased(k.instance(ref), k.b.View()).Args
}

// drain makes c drain the stream that the at-th result of the callable of f
// returns into a list, and return the list beside the error of the call. A
// drain that ends on an error of the stream returns that error.
func (k *caller) drain(c *Closure, f *derive.Facts, at int) {
	elems := k.elemsOf(f.Signature.Returns[at].Ref)
	c.Drain, c.Error = drainElements, true
	if len(elems) == 2 {
		c.Drain, c.Pair = drainPairs, true
		if k.b.TypeOf(elems[1]).Ref == rules.WellKnownError {
			c.Drain, c.Pair, elems = drainUntilError, false, elems[:1]
		}
	}
	for _, e := range elems {
		c.Returns = append(c.Returns, k.typeOf(e))
	}
	k.assign(c, f, false)
	c.Assign[at] = "r" + strconv.Itoa(at)
	c.Results = append([]Expr{{Name: c.Assign[at]}}, errorOf(c)...)
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
	first, second := k.fieldValues(in.drawn, write, index, in.changed)
	changed := Let{Name: writtenName, Select: in.changed, Value: second}
	c.Lets = []Let{{Name: writtenName, Value: k.argOf(in.drawn, write, index, 0)}}
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
// A method of a handle has no call generator in the suite, so a generated
// binding passes the samples of its inputs too.
func (k *caller) method(c *Closure, ref *node.TypeRef, name string) {
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	decl, _ := k.m.Reader().Lookup(ref.Target)
	b := k.b
	methods := methodsOf(&b, decl)
	f := k.factsOf(methods[slices.IndexFunc(methods, func(m *node.Method) bool { return m.Name == name })])
	c.Receiver, c.Method = &Expr{Name: handleName}, name
	c.Args = k.argsOf(nil, &f, nil, 0)
	k.act(c, &f)
}

// rewrite returns the assignment of the second write of the written input i
// of the callable of f: the field that the second write of pair changes, from
// the second value of the field, or the whole second value where pair names
// no field. Each probe that writes twice has the pair of the write and its
// read. d are the draws of a generated binding, and nil for a sample binding.
func (k *caller) rewrite(d *draws, f *derive.Facts, pair *derive.Pair, i int) Let {
	if pair.Changed == "" {
		return Let{Name: writtenName, Select: "", Value: k.argOf(d, f, i, 1)}
	}
	_, second := k.fieldValues(d, f, i, pair.Changed)
	return Let{Name: writtenName, Select: pair.Changed, Value: second}
}

// fieldSamples returns the two samples of the field name of the struct that
// the i-th input of the callable of f takes, or points to, each converted to
// the type t of the field. The sample of a struct sets one field alone, so a
// write that changes another field takes its values from the samples of the
// field itself.
func (k *caller) fieldSamples(f *derive.Facts, i int, name string) (first, second Expr, t *emit.TypeRef) {
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
	t = k.typeOf(k.instance(field.Type))
	return Expr{Value: typed(t, &v.Samples[0].Value)}, Expr{Value: typed(t, &v.Samples[1].Value)}, t
}

// repartition makes the later calls of c write into another partition: each
// input of the write of pair that the read takes by its name, such as the
// scope and the key of Set(ctx, scope, key, value), takes a variable of its
// own, with the first value on the first call and the second value on every
// later call. d are the draws of a generated binding, and nil for a sample
// binding.
func (k *caller) repartition(c *Closure, d *draws, f *derive.Facts, pair *derive.Pair) {
	for i := range f.Inputs {
		if !slices.ContainsFunc(pair.Read.Inputs, func(r derive.Value) bool { return r.Name == f.Inputs[i].Name }) {
			continue
		}
		name := selectorName + strconv.Itoa(i)
		c.Lets = append(c.Lets, Let{Name: name, Value: k.argOf(d, f, i, 0)})
		c.Later = append(c.Later, Let{Name: name, Value: k.argOf(d, f, i, 1)})
		c.Args[inputAt(f, i)] = Expr{Name: name}
	}
}

// selectorName is the prefix of the variable of an input that selects the
// partition of a write.
const selectorName = "s"

// argsOf returns the arguments of a call of the callable of f: the context of
// the closure for a parameter in the role context, and the value which of
// each input. The read of pair takes the arguments that the pair binds, and
// the write of pair the value which of the written input. A parameter that
// the pair leaves to its own value takes it, and a closure of the probe
// replaces the value of a function that the probe supplies. d are the draws
// of a generated binding, and nil for a sample binding.
func (k *caller) argsOf(d *draws, f *derive.Facts, pair *derive.Pair, which int) []Expr {
	out := make([]Expr, 0, len(f.Signature.Params))
	if pair != nil && same(pair.Read, f) {
		for _, a := range pair.Arguments {
			switch a.Kind {
			case derive.ArgContext:
				out = append(out, Expr{Word: ctxName})
			case derive.ArgShared:
				i := slices.IndexFunc(pair.Write.Inputs, func(v derive.Value) bool { return v.Name == a.Param })
				out = append(out, k.argOf(d, pair.Write, i, which))
			case derive.ArgField:
				written := k.argOf(d, pair.Write, pair.Written, which)
				out = append(out, Expr{Of: &written, Select: a.Field})
			case derive.ArgOwn:
				i := slices.IndexFunc(f.Inputs, func(v derive.Value) bool { return v.Name == a.Name })
				out = append(out, k.argOf(d, f, i, which))
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
		out = append(out, k.argOf(d, f, i, which))
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
