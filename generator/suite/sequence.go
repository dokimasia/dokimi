// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"slices"
	"strconv"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
)

// The words of the names of the parts of the generated sequences: the function
// of the actions of a declaration, as StoreActions, the function of its
// reference, as StoreReference, and the function that calls a callable with
// the arguments of a list, as StoreGetCall.
const (
	actionsWord   = "Actions"
	referenceWord = "Reference"
	callWord      = "Call"
)

// keyParam is the binding of a shape that names the key input of a callable,
// and keysWord the word that opens the name of a pool of keys.
const (
	keyParam = "key"
	keysWord = "keys"
)

// Sequences are the parts of the suite of a declaration that its checks of
// generated sequences call. Actions is the function of the actions of the
// declaration, and Reference the function that applies an operation to a
// reference by the name of its callable, which is empty where no check calls
// one. Normalize is the function of the param normalize of the directive
// suite, which maps the values of each call, and nil without the param. Pools
// are the pools of the keys that the actions of a case share, and Steps the
// callables of the declaration, in their order. Txn is the action of a
// transaction, and nil for a declaration without a check of transactions.
type Sequences struct {
	Normalize *emit.TypeRef
	Txn       *Txn
	Actions   string
	Reference string
	Pools     []Pool
	Steps     []Step
}

// Pool is the pool of the keys of the type Type that the actions of a case
// draw once, under the name Name and the label Label.
type Pool struct {
	Type  *emit.TypeRef
	Name  string
	Label string
}

// Step is one callable of the declaration in the generated sequences. Name is
// the name of the operation that the history records, which is the name of the
// callable. Call is the function that calls the callable with the arguments of
// a list, and Field the method of the subject that it calls, which is empty
// for a function. Params are the parameters of the call, and Keys the keys
// that the history records for a call. Assign are the names of the results of
// the call, Values the names of its values in order, and Error the name of its
// error, which is empty for a callable without one. Drain is the drain of a
// stream that the callable returns, Stream the name of the result that is the
// stream, and Elems the types of the elements of the stream. Action reports a
// callable with an action, which returns on its own. Builder is the builder of
// the mixin sample of the callable, whose values the call takes in place of
// the arguments of the list, and nil for a callable without one. An action of
// a callable with a builder draws no input, and its calls have no keys.
// Receiver is the type of the handle whose method a step of a transaction
// calls, and nil for a step that calls the subject. Produced is the key that a
// call of a callable with the mixin creates returns to the pool of its type,
// and nil for a callable without one.
type Step struct {
	Builder  *Builder
	Receiver *emit.TypeRef
	Produced *Produced
	Name     string
	Call     string
	Field    string
	Error    string
	Drain    string
	Stream   string
	Params   []StepParam
	Keys     []StepKey
	Assign   []string
	Values   []string
	Elems    []*emit.TypeRef
	Action   bool
}

// StepParam is one parameter of the call of a step: the parameter of the call
// generator of the callable, its position Arg in the list of the arguments,
// which is -1 for the context and for an omitted parameter, and the name
// Local of its value in the function of the call. Label is the label of its
// draw, Pool the pool of a key, and Fields the fields of its value that take
// keys of the pools.
type StepParam struct {
	Local  string
	Label  string
	Pool   string
	Fields []PoolField
	Param
	Arg int
}

// PoolField is a field of the value of a parameter that takes a key from the
// pool Pool under the label Label.
type PoolField struct {
	Name  string
	Pool  string
	Label string
}

// Produced is the key that a call of a callable with the mixin creates
// returns, which the world adds to the pool Pool once the call returns without
// an error. The key is the first value of the call, as a value of the type
// Type, or the field Field of that value. Pointer reports a first value that is
// a pointer, which the world skips where it is nil.
type Produced struct {
	Type    *emit.TypeRef
	Pool    string
	Field   string
	Pointer bool
}

// StepKey is one key of a call: the argument at the position Arg, or its field
// Field, which reads the argument as a value of the type Type.
type StepKey struct {
	Type  *emit.TypeRef
	Field string
	Arg   int
}

// Apply is a function that a probe of generated sequences calls with the
// arguments of a call of another callable. It calls the function Call with
// Args, and declares Lets for the arguments that Args read a field of.
type Apply struct {
	Call string
	Lets []ApplyLet
	Args []Expr
}

// ApplyLet is the argument at the position Index of the arguments of an
// [Apply], as a value of the type Type under the name Name.
type ApplyLet struct {
	Type  *emit.TypeRef
	Name  string
	Index int
}

// Written is a function that returns the argument at the position Index of a
// call of a write, as a value of the type Type.
type Written struct {
	Type  *emit.TypeRef
	Index int
}

// Countable is a call of an observer that returns the use of the subject: the
// function Call with Args, whose first value has the type Type. Count is len
// for a value whose length is the use, and int for an integer.
type Countable struct {
	Type  *emit.TypeRef
	Call  string
	Count string
	Args  []Expr
}

// The counts of a countable use.
const (
	countLen = "len"
	countInt = "int"
)

// sequenced reports whether the check c runs generated sequences.
func sequenced(c *derive.Check) bool {
	return slices.Contains(c.Claim.Machineries, spec.GeneratedSequences)
}

// sequencesOf returns the parts of the generated sequences of the suite s. A
// key input of a callable takes its key from the pool of its type, and so does
// a field of a written value that a pair of a check of the declaration binds
// to a key input of its read. The checks are visited in their order, and the
// partners of each in the order of its claim, so each run writes the keys in
// one order. [Sequences.used] drops the parts that no derived check calls.
// The action of a transaction joins the actions where a check of the
// derivation has the scope of a transaction. Once every pool exists, each
// step of a callable with the mixin creates gives its key to the pool of the
// type of the key.
func (k *caller) sequencesOf(s *Suite) *Sequences {
	q := &Sequences{
		Actions:   k.e.JoinName(actionsWord, s.Declaration),
		Reference: k.e.JoinName(referenceWord, s.Declaration),
		Normalize: k.normalizeOf(),
	}
	for i := range s.Facts {
		q.Steps = append(q.Steps, k.step(s, q, i))
	}
	for i := range s.Derivation.Checks {
		c := &s.Derivation.Checks[i]
		for _, p := range c.Claim.Partners {
			poolFields(s, q, c.Pairs[p.Name])
		}
		if c.Transaction != nil {
			q.Txn = k.txnOf(s, q, c.Transaction)
		}
	}
	for i := range s.Facts {
		q.Steps[i].Produced = k.producedOf(&s.Facts[i], q)
	}
	return q
}

// producedOf returns the key that a call of the callable of f produces, where
// the callable has the mixin creates: its first value, or the field of that
// value that the param key of the mixin names. It returns nil for a callable
// without the mixin or without a value, and for a key of a type that no pool
// of q has, because no call of a case takes such a key.
func (k *caller) producedOf(f *derive.Facts, q *Sequences) *Produced {
	if !mixed(k.m, &f.ID, shape.Creates) || len(f.Results) == 0 {
		return nil
	}
	ref := k.resultRef(f, 0)
	p := &Produced{Type: k.typeOf(ref), Pointer: ref.Form == symbol.FormOptional}
	key := p.Type
	if field, named := f.Class.Params[string(shape.Creates)][shape.CreatesKey].(symbol.Identity); named {
		p.Field, key = field.Name, k.fieldRef(ref, field.Name)
	}
	pool := q.poolOf(key)
	if pool == nil {
		return nil
	}
	p.Pool = pool.Name
	return p
}

// used returns q with the parts that the derived checks call: nil where no
// derived check runs generated sequences or a world, without the actions
// where no check runs generated sequences, and without the reference where no
// check applies an operation by its name. A world calls the steps of its
// operations, and applies the operations of its history to the reference.
func (q *Sequences) used(derived []Derived) *Sequences {
	sequences, world, reference := false, false, false
	for i := range derived {
		d := &derived[i]
		sequences, world = sequences || d.Sequence, world || d.World
		reference = reference || d.World ||
			slices.ContainsFunc(d.Args, func(a Arg) bool { return a.Ref == q.Reference })
	}
	if !sequences && !world {
		return nil
	}
	if !sequences {
		q.Actions = ""
	}
	if !reference {
		q.Reference = ""
	}
	return q
}

// normalizeOf returns the function of the param normalize of the directive
// suite of the declaration, and nil for a directive without the param. The
// function takes the name of a callable and the values of its call, and
// returns the values that the history records. normalizeOf reports
// [BadNormalize] at the directive and returns nil for a function of another
// type than func(string, []any) []any.
func (k *caller) normalizeOf() *emit.TypeRef {
	d := k.m.Directive()
	param, stated := d.Param(NormalizeParam)
	if !stated {
		return nil
	}
	fn := param.Target
	found, _ := k.m.Reader().Lookup(fn)
	view, _ := k.b.CallableOf(found)
	if got := funcType(&view); got != normalizeType {
		k.m.ErrorfAt(BadNormalize, d.Pos, "the function %s of the param normalize of %s has the type %s, and the "+
			"function of the param has the type %s, so the suite records the values of each call as they are",
			fn.Name, k.decl, got, normalizeType)
		return nil
	}
	return &emit.TypeRef{Spelling: fn.Name, Target: fn}
}

// normalizeType is the type of the function of the param normalize.
const normalizeType = "func(string, []any) []any"

// funcType returns the type of the callable of view as Go spells a function
// type, as func(string, []any) []any.
func funcType(view *rules.Callable) string {
	params := make([]string, 0, len(view.Params))
	for _, p := range view.Params {
		params = append(params, p.Ref.Spelling)
	}
	returns := make([]string, 0, len(view.Returns))
	for _, r := range view.Returns {
		returns = append(returns, r.Ref.Spelling)
	}
	out := "func(" + strings.Join(params, ", ") + ")"
	switch len(returns) {
	case 0:
		return out
	case 1:
		return out + " " + returns[0]
	default:
		return out + " (" + strings.Join(returns, ", ") + ")"
	}
}

// step returns the step of the i-th callable of s. A key input of the
// callable takes the pool of its type.
func (k *caller) step(s *Suite, q *Sequences, i int) Step {
	f, call := &s.Facts[i], &s.Calls[i]
	base := s.Declaration
	if f.Declaration != "" {
		base = s.Callables[i].Type
	}
	label := naming.Kebab(f.Callable)
	st := Step{
		Name: f.Callable, Call: k.e.JoinName(callWord, base), Field: call.Field, Action: call.repeatable(),
		Builder: call.Builder,
	}
	keys := keyNames(f)
	arg := 0
	for j := range call.Params {
		p := StepParam{Param: call.Params[j], Label: label + "." + call.Params[j].Label, Arg: -1}
		if !p.Context && !p.Omitted {
			p.Arg, p.Local = arg, "a"+strconv.Itoa(arg)
			arg++
			if call.Builder == nil && slices.Contains(keys, f.Signature.Params[j].Name) {
				p.Pool = q.pool(p.Type).Name
				st.Keys = append(st.Keys, StepKey{Arg: p.Arg})
			}
		}
		st.Params = append(st.Params, p)
	}
	k.results(&st, f)
	return st
}

// keyNames returns the names of the parameters of the callable of f that a
// binding key of its classifications names.
func keyNames(f *derive.Facts) []string {
	var out []string
	for _, params := range f.Class.Params {
		if id, ok := params[directive.ParamKey(keyParam)].(symbol.Identity); ok {
			out = append(out, id.Name)
		}
	}
	return out
}

// pool returns the pool of the keys of the type t, which it adds to q where q
// has none.
func (q *Sequences) pool(t *emit.TypeRef) *Pool {
	if p := q.poolOf(t); p != nil {
		return p
	}
	name := keysWord + naming.Pascal(t.Spelling)
	q.Pools = append(q.Pools, Pool{Type: t, Name: name, Label: keysWord + "." + naming.Kebab(t.Spelling)})
	return &q.Pools[len(q.Pools)-1]
}

// poolOf returns the pool of the keys of the type t, and nil where q has none.
func (q *Sequences) poolOf(t *emit.TypeRef) *Pool {
	at := slices.IndexFunc(q.Pools, func(p Pool) bool { return p.Type.Spelling == t.Spelling })
	if at < 0 {
		return nil
	}
	return &q.Pools[at]
}

// poolFields gives the write of pair the keys of the pools that its read takes
// as keys: a shared parameter of the write takes the pool, and a field of the
// written value takes it in the action of the write. Each becomes a key of
// the calls of the write. A pair without a read has no write either. A read
// with a builder takes no pool, and a write with a builder takes the values of
// the builder in place of keys.
func poolFields(s *Suite, q *Sequences, pair *derive.Pair) {
	if pair.Read == nil {
		return
	}
	read, write := stepOf(s, q, pair.Read), stepOf(s, q, pair.Write)
	if write.Builder != nil {
		return
	}
	for _, a := range pair.Arguments {
		rp := paramNamed(read, pair.Read, a.Name)
		if rp.Pool == "" {
			continue
		}
		switch a.Kind {
		case derive.ArgShared:
			if wp := paramNamed(write, pair.Write, a.Param); wp.Pool == "" {
				wp.Pool = rp.Pool
				write.Keys = append(write.Keys, StepKey{Arg: wp.Arg})
			}
		case derive.ArgField:
			wp := &write.Params[inputAt(pair.Write, pair.Written)]
			if slices.ContainsFunc(wp.Fields, func(f PoolField) bool { return f.Name == a.Field }) {
				continue
			}
			wp.Fields = append(wp.Fields, PoolField{
				Name: a.Field, Pool: rp.Pool, Label: wp.Label + "." + naming.Kebab(a.Field),
			})
			write.Keys = append(write.Keys, StepKey{Arg: wp.Arg, Type: wp.Type, Field: a.Field})
		default:
		}
	}
}

// stepOf returns the step of the callable of f.
func stepOf(s *Suite, q *Sequences, f *derive.Facts) *Step {
	i := slices.IndexFunc(s.Facts, func(o derive.Facts) bool { return o.ID == f.ID })
	return &q.Steps[i]
}

// paramNamed returns the parameter of the step st of the callable of f whose
// parameter in the signature has the name name. A pair names each argument
// after a parameter of its read, and each shared argument after a parameter
// of its write.
func paramNamed(st *Step, f *derive.Facts, name string) *StepParam {
	j := slices.IndexFunc(f.Signature.Params, func(p rules.ParamView) bool { return p.Name == name })
	return &st.Params[j]
}

// results sets the results of the call of the step st of the callable of f:
// r and its position for a value, ok for the flag and err for the error. A
// stream takes the drain of its elements into a list.
func (k *caller) results(st *Step, f *derive.Facts) {
	for i, r := range f.Signature.Returns {
		switch r.Role {
		case rules.ReturnError:
			st.Error = errName
			st.Assign = append(st.Assign, errName)
		case rules.ReturnOkBool:
			st.Assign = append(st.Assign, okName)
			st.Values = append(st.Values, okName)
		default:
			name := "r" + strconv.Itoa(i)
			st.Assign = append(st.Assign, name)
			st.Values = append(st.Values, name)
		}
	}
	at := streamAt(f)
	if at < 0 {
		return
	}
	c := &Closure{}
	k.drain(c, f, at)
	st.Drain, st.Elems, st.Stream = c.Drain, c.Returns, "r"+strconv.Itoa(at)
}

// okName is the name of the flag of a call.
const okName = "ok"

// stepCall returns the name of the function that calls the callable of f in q.
func (q *Sequences) stepCall(s *Suite, f *derive.Facts) string { return stepOf(s, q, f).Call }

// applyOf returns the function that calls the callable target with the
// arguments of a call of the callable source. pair is the pair of the check
// whose read is target and whose write is source, and nil where the check has
// none. The read of pair takes the arguments that the pair binds from its
// write. Any other target takes for each input the argument of an input of
// the source with its name and its type, and otherwise its first sample.
func (k *caller) applyOf(s *Suite, q *Sequences, target, source *derive.Facts, pair *derive.Pair) *Apply {
	a := &Apply{Call: q.stepCall(s, target)}
	if pair != nil {
		for _, arg := range pair.Arguments {
			switch arg.Kind {
			case derive.ArgShared:
				i := slices.IndexFunc(source.Inputs, func(v derive.Value) bool { return v.Name == arg.Param })
				a.Args = append(a.Args, argAt(i))
			case derive.ArgField:
				// Each field is a field of the one written value.
				a.Lets = []ApplyLet{{
					Name: writtenName, Index: pair.Written, Type: k.typeOf(k.inputRef(source, pair.Written)),
				}}
				a.Args = append(a.Args, Expr{Of: &Expr{Name: writtenName}, Select: arg.Field})
			case derive.ArgOwn:
				i := slices.IndexFunc(target.Inputs, func(v derive.Value) bool { return v.Name == arg.Name })
				a.Args = append(a.Args, k.sampleOf(target, i, 0))
			default:
			}
		}
		return a
	}
	for i := range target.Inputs {
		in := &target.Inputs[i]
		j := slices.IndexFunc(source.Inputs, func(v derive.Value) bool {
			return v.Name == in.Name && v.Shape.Spelling == in.Shape.Spelling
		})
		if j < 0 {
			a.Args = append(a.Args, k.sampleOf(target, i, 0))
			continue
		}
		a.Args = append(a.Args, argAt(j))
	}
	return a
}

// replicable reports whether the callable of f takes one input that another
// subject of the suite s can be: a pointer to the declaration, or an interface
// that a pointer to the declaration satisfies, such as the declaration itself.
func (k *caller) replicable(s *Suite, f *derive.Facts) bool {
	if len(f.Inputs) != 1 {
		return false
	}
	if f.Inputs[0].Shape.Spelling == "*"+s.Declaration {
		return true
	}
	subject := &node.TypeRef{Form: symbol.FormOptional, Elems: []*node.TypeRef{{Target: k.origin}}}
	satisfaction, ok := k.b.Satisfies(subject, k.inputRef(f, 0))
	return ok && satisfaction.Satisfied()
}

// graphed reports whether the param v of a classification names a variable
// that can be a graph of states: a variable of a map type, or a variable
// whose type its initializer states. eidos resolves the param transitions of
// the contract workflow to a variable of the package.
func (k *caller) graphed(v any) bool {
	id, _ := v.(symbol.Identity)
	sym, _ := k.m.Reader().Lookup(id)
	variable, _ := sym.(*node.Variable)
	return variable.Type == nil || variable.Type.Form == symbol.FormMap
}

// argAt returns the expression of the argument at the position i of a list.
func argAt(i int) Expr { return Expr{Name: "args[" + strconv.Itoa(i) + "]"} }

// countableOf returns the call of the observer f that returns the use of the
// subject with the samples of its inputs, and false where its first value is
// neither an integer nor a value whose length is the use.
func (k *caller) countableOf(s *Suite, q *Sequences, f *derive.Facts) (*Countable, bool) {
	if len(f.Results) == 0 {
		return nil, false
	}
	shape := f.Results[0].Shape
	c := &Countable{Call: q.stepCall(s, f), Type: k.typeOf(k.resultRef(f, 0))}
	// Only a scalar has a class.
	switch {
	case shape.Class == rules.ScalarInt || shape.Class == rules.ScalarUint:
		c.Count = countInt
	case shape.Form == symbol.FormList || shape.Form == symbol.FormText:
		c.Count = countLen
	default:
		return nil, false
	}
	for i := range f.Inputs {
		c.Args = append(c.Args, k.sampleOf(f, i, 0))
	}
	return c, true
}

// elementOf returns the types of the elements of the list or the stream that
// the callable of f returns first: one type, or the two types of the
// probe.Pair of the elements of a stream of pairs.
func (k *caller) elementOf(f *derive.Facts) []*emit.TypeRef {
	refs := k.elementRefs(f)
	out := make([]*emit.TypeRef, 0, len(refs))
	for _, ref := range refs {
		out = append(out, k.typeOf(ref))
	}
	return out
}

// elementRefs returns the types of the elements of the list or the stream
// that the callable of f returns first. A stream of values and errors has the
// type of its values alone.
func (k *caller) elementRefs(f *derive.Facts) []*node.TypeRef {
	if at := streamAt(f); at >= 0 {
		elems := k.elemsOf(f.Signature.Returns[at].Ref)
		if len(elems) == 2 && k.b.TypeOf(elems[1]).Ref == rules.WellKnownError {
			return elems[:1]
		}
		return elems
	}
	return []*node.TypeRef{k.resultRef(f, 0).Elems[0]}
}
