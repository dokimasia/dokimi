// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive

import (
	"reflect"
	"slices"

	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/spec"
)

// Transaction is the scope of the transactions of a claim with the gate
// transactions: the callables that begin, write, read, commit and roll back
// one transaction of an instance of the contract tx. The begin returns the
// handle of the transaction. Each other part takes the handle as an input, or
// is a method of the handle.
type Transaction struct {
	Begin, Write, Read, Commit, Rollback Part

	// Key and Value are the positions of the key and the value among the
	// inputs of the write, and ReadKey is the position of the key among the
	// inputs of the read.
	Key, Value, ReadKey int

	// Appends reports a read that returns a list of the type of the written
	// value, so the write appends its value to the list of its key. Any other
	// read returns the value of a register, which the write sets.
	Appends bool
}

// Part is one callable of a transaction. Handle is the position of the handle
// among the inputs of the callable, and -1 for the begin and for a method of
// the handle.
type Part struct {
	Facts  *Facts
	Handle int
}

// readParam is the param of the mixins serializable and snapshot-isolation
// that names the read of a transaction.
const readParam = "read"

// transactional reports whether the claim c has the gate transactions. The
// write of such a claim can be a method of the handle that a begin returns.
func transactional(c *spec.Claim) bool { return slices.Contains(c.Gates, spec.Transactions) }

// licensee returns the callable that has the licence l of the claim c on
// behalf of the callable g: g itself, or a method of the handle that g
// returns for a claim with the gate transactions. It returns nil where
// neither has the licence.
func licensee(c *spec.Claim, l *spec.Licence, g *Facts) *Facts {
	if holds(l, g) {
		return g
	}
	if !transactional(c) {
		return nil
	}
	for i := range g.Handle {
		if holds(l, &g.Handle[i]) {
			return &g.Handle[i]
		}
	}
	return nil
}

// transactionOf returns the scope of the transactions of the write: a begin
// of the contract tx among the callables of decl, with the read that the
// param read of the write names and the commit and the rollback of the
// instance of the begin. It returns false where no begin of decl has a scope
// that fits the gate transactions.
func transactionOf(write *Facts, decl []Facts) (*Transaction, bool) {
	for i := range decl {
		begin := &decl[i]
		instance, ok := roleIn(begin, shape.TxBegin)
		if !ok || len(begin.Results) == 0 {
			continue
		}
		if tx, fits := scope(write, begin, instance, decl); fits {
			return tx, true
		}
	}
	return nil, false
}

// scope returns the transaction of the write whose begin is begin, of the
// instance with the id instance. Each part after the begin is a callable of
// decl that takes the handle that begin returns, or a method of the handle.
// The write takes a key and a value of an integer or a string type besides
// the handle, and the read the key alone. scope returns false where a part is
// missing or does not fit.
func scope(write, begin *Facts, instance string, decl []Facts) (*Transaction, bool) {
	handle := begin.Results[0].Shape
	// The licence of the claim needs the param read, so the write names it.
	read, _ := paramID(readParam, write)
	tx := &Transaction{Begin: Part{Facts: begin, Handle: -1}}
	finds := []struct {
		part  *Part
		match func(f *Facts) bool
	}{
		{&tx.Write, func(f *Facts) bool { return f == write }},
		{&tx.Read, func(f *Facts) bool { return f.ID == read }},
		{&tx.Commit, func(f *Facts) bool { return playsIn(f, shape.TxCommit, instance) }},
		{&tx.Rollback, func(f *Facts) bool { return playsIn(f, shape.TxRollback, instance) }},
	}
	for _, find := range finds {
		p, found := locate(begin, &handle, decl, find.match)
		if !found {
			return nil, false
		}
		*find.part = p
	}
	values := others(tx.Write)
	keys := others(tx.Read)
	if len(values) != 2 || len(keys) != 1 || len(tx.Read.Facts.Results) != 1 {
		return nil, false
	}
	tx.Key, tx.Value, tx.ReadKey = values[0], values[1], keys[0]
	key, value := &tx.Write.Facts.Inputs[tx.Key].Shape, &tx.Write.Facts.Inputs[tx.Value].Shape
	result := &tx.Read.Facts.Results[0].Shape
	tx.Appends = result.Form == symbol.FormList && reflect.DeepEqual(result.Elems[0], *value)
	fits := counted(key) && counted(value) &&
		reflect.DeepEqual(tx.Read.Facts.Inputs[tx.ReadKey].Shape, *key) &&
		(tx.Appends || reflect.DeepEqual(*result, *value))
	return tx, fits
}

// locate returns the part that match accepts: a callable of decl that takes
// the handle as an input, or a method of the handle that begin returns. It
// returns false where neither has such a callable.
func locate(begin *Facts, handle *rules.TypeShape, decl []Facts, match func(f *Facts) bool) (Part, bool) {
	for i := range decl {
		f := &decl[i]
		if !match(f) {
			continue
		}
		at := slices.IndexFunc(f.Inputs, func(v Value) bool { return reflect.DeepEqual(v.Shape, *handle) })
		return Part{Facts: f, Handle: at}, at >= 0
	}
	for i := range begin.Handle {
		if f := &begin.Handle[i]; match(f) {
			return Part{Facts: f, Handle: -1}, true
		}
	}
	return Part{}, false
}

// others returns the positions of the inputs of the part p other than the
// handle.
func others(p Part) []int {
	var out []int
	for i := range p.Facts.Inputs {
		if i != p.Handle {
			out = append(out, i)
		}
	}
	return out
}

// Begins reports whether the callable of f has the role begin of the contract
// tx, so it returns the handle of a transaction.
func Begins(f *Facts) bool {
	_, begins := roleIn(f, shape.TxBegin)
	return begins
}

// roleIn returns the id of the instance of the contract tx in which the
// callable of f has the role, and false where it has no such role.
func roleIn(f *Facts, role shape.Role) (string, bool) {
	for _, m := range f.Class.Contracts {
		if m.Contract == shape.Tx && m.Role == role {
			return m.ID, true
		}
	}
	return "", false
}

// playsIn reports whether the callable of f has the role in the instance of
// the contract tx with the id instance.
func playsIn(f *Facts, role shape.Role, instance string) bool {
	return slices.ContainsFunc(f.Class.Contracts, func(m shape.Membership) bool {
		return m.Contract == shape.Tx && m.Role == role && m.ID == instance
	})
}

// counted reports whether a value of the shape s is an integer or a string,
// which a counter of the case can give each write of the case. Only a number
// has a scalar class.
func counted(s *rules.TypeShape) bool {
	return s.Form == symbol.FormText || s.Class == rules.ScalarInt || s.Class == rules.ScalarUint
}
