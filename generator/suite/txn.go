// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"slices"
	"strconv"
	"strings"

	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
)

// Txn is the action of a transaction, which the function of the actions of a
// declaration adds after the action of each callable. Each probe of an
// isolation level runs the action alone, and every other probe of generated
// sequences leaves it out, because its history is a history of transactions.
//
// The input of the action draws one to four micro-operations over the keys of
// the pool Pool, of the type Key, and whether the transaction commits. Each
// micro-operation reads a key or writes it. A write takes the next value of a
// counter of the case, of the type Value, which a string type takes as its
// decimal text where Text is set. Function is the function of a write in the
// history: append where Appends reports a read that returns a list of the
// written values, and write for a register.
//
// The action begins a transaction on the subject, calls the write or the read
// of each micro-operation with the handle, and then commits or rolls back. It
// records one call of the operation txn whose arguments are the
// micro-operations. Each call of the transaction is the call of the function
// of a step: a step of the declaration, or a step of Steps for a method of the
// handle of the type Handle.
type Txn struct {
	Key, Value, Handle                   *emit.TypeRef
	Pool                                 string
	Function                             string
	Steps                                []Step
	Begin, Write, Read, Commit, Rollback TxnCall
	Text                                 bool
	Appends                              bool
}

// TxnCall is one call of a transaction: the function Call of the step of its
// callable with the arguments Args, which takes the subject, or the handle of
// the transaction where Handled is set. An argument is the handle tx, the key,
// the value v, or the first sample of another input.
type TxnCall struct {
	Call    string
	Args    []Expr
	Handled bool
}

// The names of the handle, the key and the value of a call of a transaction,
// the operation of a transaction in the history, and the functions of its
// micro-operations.
const (
	txName         = "tx"
	keyName        = "key"
	valueName      = "v"
	txnOperation   = "txn"
	appendFunction = "append"
	writeFunction  = "write"
)

// txnOf returns the action of the transactions of the scope tx in the
// sequences q of the suite s. The keys of its writes and reads take the pool
// of their type.
func (k *caller) txnOf(s *Suite, q *Sequences, tx *derive.Transaction) *Txn {
	write := tx.Write.Facts
	key := k.typeOf(k.inputRef(write, tx.Key))
	t := &Txn{
		Key: key, Value: k.typeOf(k.inputRef(write, tx.Value)), Handle: k.typeOf(k.resultRef(tx.Begin.Facts, 0)),
		Pool: q.pool(key).Name, Function: writeFunction, Appends: tx.Appends,
		Text: write.Inputs[tx.Value].Shape.Form == symbol.FormText,
	}
	if tx.Appends {
		t.Function = appendFunction
	}
	t.Begin = k.txnCall(s, q, t, tx.Begin, nil)
	t.Write = k.txnCall(s, q, t, tx.Write, map[int]string{tx.Key: keyName, tx.Value: valueName})
	t.Read = k.txnCall(s, q, t, tx.Read, map[int]string{tx.ReadKey: keyName})
	t.Commit = k.txnCall(s, q, t, tx.Commit, nil)
	t.Rollback = k.txnCall(s, q, t, tx.Rollback, nil)
	return t
}

// txnCall returns the call of the part p of a transaction: the call of the
// step of its callable, or of a step of t for a method of the handle, with
// the handle at its position, the input of each position of names under its
// name, and the first sample of each other input.
func (k *caller) txnCall(s *Suite, q *Sequences, t *Txn, p derive.Part, names map[int]string) TxnCall {
	f := p.Facts
	c := TxnCall{}
	if at := slices.IndexFunc(s.Facts, func(o derive.Facts) bool { return o.ID == f.ID }); at >= 0 {
		c.Call = q.Steps[at].Call
	} else {
		c.Handled = true
		c.Call = k.handleStep(t, f).Call
	}
	for i := range f.Inputs {
		name, named := names[i]
		switch {
		case i == p.Handle:
			c.Args = append(c.Args, Expr{Name: txName})
		case named:
			c.Args = append(c.Args, Expr{Name: name})
		default:
			c.Args = append(c.Args, k.sampleOf(f, i, 0))
		}
	}
	return c
}

// handleStep returns the step of the method of f of the handle of t, which
// it adds to the steps of t. Its function takes the handle in place of the
// subject, and the arguments of a list.
func (k *caller) handleStep(t *Txn, f *derive.Facts) Step {
	found, _ := k.m.Reader().Lookup(f.ID)
	view, _ := k.b.CallableOf(found)
	st := Step{
		Name: f.Callable, Call: k.e.JoinName(callWord, k.e.JoinName(f.Callable, f.Declaration)), Field: f.Callable,
		Receiver: t.Handle,
	}
	arg := 0
	for _, p := range k.built(&view) {
		sp := StepParam{Param: p, Arg: -1}
		if !p.Context {
			sp.Arg, sp.Local = arg, "a"+strconv.Itoa(arg)
			arg++
		}
		st.Params = append(st.Params, sp)
	}
	k.results(&st, f)
	t.Steps = append(t.Steps, st)
	return st
}

// TxnPlant is what an override of a planted defect of transactions writes.
// Txn is the action of the transactions, whose calls the commit of a late
// transaction makes. Handle is the expression of the handle of the
// transaction of the call: a parameter of the override, or the wrapper t of
// a handle. Key and Value are the parameters of the key and the value of a
// write, and the read takes the key alone. Context is the parameter of the
// context. Wrapper is the wrapper of the handle that the override of a begin
// returns, and nil for any other override.
type TxnPlant struct {
	Txn     *Txn
	Wrapper *Defect
	Handle  string
	Key     string
	Value   string
	Context string
}

// The bodies of the overrides of the planted defects of transactions. A
// write of reads-uncommitted keeps its value as the open value of its key, a
// read returns the open value of another transaction, and a settlement drops
// the open values of its transaction. A write of writes-without-a-snapshot
// stages its value, a read adds the staged values of its own transaction, a
// commit applies the staged values in a late transaction of its own, and a
// rollback drops them. The begin of a handle wraps the handle that it opens.
const (
	bodyKeepOpen   = "keep-open"
	bodyReadOpen   = "read-open"
	bodyDropOpen   = "drop-open"
	bodyStage      = "stage"
	bodyReadStaged = "read-staged"
	bodyApplyLate  = "apply-late"
	bodyDropStaged = "drop-staged"
	bodyWrapHandle = "wrap-handle"
	wrapperWord    = "Txn"
	wrapperName    = "t"
)

// txnDefect sets the overrides of the planted defect d of the form id of the
// transactions tx of the check c on the suite s. The callables of the
// declaration in the transaction override their callable on the type of d. A
// method of the handle overrides its method on the wrapper of the handle,
// which the begin of d returns in place of the handle that it opens.
func (k *caller) txnDefect(s *Suite, d *Defect, id spec.DefectFormID, tx *derive.Transaction) {
	bodies := map[*derive.Facts]string{tx.Write.Facts: bodyKeepOpen, tx.Read.Facts: bodyReadOpen}
	bodies[tx.Commit.Facts], bodies[tx.Rollback.Facts] = bodyDropOpen, bodyDropOpen
	if id == spec.WritesWithoutASnapshot {
		bodies = map[*derive.Facts]string{tx.Write.Facts: bodyStage, tx.Read.Facts: bodyReadStaged}
		bodies[tx.Commit.Facts], bodies[tx.Rollback.Facts] = bodyApplyLate, bodyDropStaged
	}
	t := s.Sequences.Txn
	d.Keeps, d.Txn, d.Body = true, &TxnPlant{Txn: t}, bodies[tx.Write.Facts]
	parts := []derive.Part{tx.Write, tx.Read, tx.Commit, tx.Rollback}
	for _, p := range parts {
		o := k.txnOverride(p, bodies[p.Facts], t)
		if p.Handle >= 0 {
			o.Txn.Handle = o.Args[inputAt(p.Facts, p.Handle)]
			d.Overrides = append(d.Overrides, o)
			continue
		}
		if d.Wrapper == nil {
			d.Wrapper = &Defect{Name: k.e.JoinName(wrapperWord, d.Name), Handled: t.Handle}
			begin := k.txnOverride(tx.Begin, bodyWrapHandle, t)
			begin.Txn.Wrapper = d.Wrapper
			d.Overrides = append(d.Overrides, begin)
		}
		o.Txn.Handle = wrapperName
		d.Wrapper.Overrides = append(d.Wrapper.Overrides, o)
	}
}

// EmbeddedField returns the name of the field of the handle that the wrapper
// d of a handle embeds: the name of the type of the handle without its
// package.
func (d *Defect) EmbeddedField() string {
	return d.Handled.Spelling[strings.LastIndex(d.Handled.Spelling, ".")+1:]
}

// txnOverride returns the override of the callable of the part p of a
// transaction with the body body, over the action t of the transactions.
func (k *caller) txnOverride(p derive.Part, body string, t *Txn) *Defect {
	f := p.Facts
	o := &Defect{Body: body, Delegates: true, Txn: &TxnPlant{Txn: t}}
	o.Method, o.Args = k.signatureOf(f)
	o.Results = k.outcomesOf(f)
	guard(o, f)
	o.Guards = false
	o.Txn.Context = o.Context
	inputs := make([]int, 0, len(f.Inputs))
	for i := range f.Inputs {
		inputs = append(inputs, inputAt(f, i))
	}
	keyed := slices.DeleteFunc(slices.Clone(inputs), func(at int) bool {
		return p.Handle >= 0 && at == inputAt(f, p.Handle)
	})
	if len(keyed) > 0 {
		o.Txn.Key = o.Args[keyed[0]]
	}
	if len(keyed) > 1 {
		o.Txn.Value = o.Args[keyed[1]]
	}
	return o
}
