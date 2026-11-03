// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive

import (
	"fmt"
	"reflect"

	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/rules"

	"go.dokimi.dev/dokimi/spec"
)

// Reference is the reference form that covers every callable of one
// declaration, with the effect of each callable in the order of the
// callables.
type Reference struct {
	Form      spec.ReferenceForm
	Callables []Effected
}

// Effected is one callable of a declaration with a reference form, and its
// effect on the state of the form.
//
// Key is the key of a store or of a load of a form of the state map, by the
// rule of [Pairing]. Its Name is the parameter of the key of a load. Its Kind
// locates the key in a store: the parameter Param of the store, or the field
// Field of the value that the store stores. Every load selects the same key of
// one store, and a load has the key that it selects of the first store. An
// append has the zero Key.
//
// Value is the index in the inputs of the value that a store stores or that an
// append appends, and -1 for a load.
type Effected struct {
	Facts  *Facts
	Effect spec.Effect
	Key    Argument
	Value  int
}

// Unreferenced is the refusal of the gate reference-role on a callable that no
// reference form covers. Detail is what of the callable no form covers, such
// as "has the mixin ttl", with the callable as its subject.
type Unreferenced struct {
	Refusal
	Detail string
}

// ReferenceOf returns the reference form that covers every callable of the
// declaration decl, which has at least one callable, with the effect of each
// callable. It returns the refusal of the gate reference-role on a callable
// that no form covers instead.
//
// A callable with a role of a contract with a form has the effect of its role,
// and the callables with a role have their roles in one instance of one
// contract. The shape that the catalog detects on such a callable describes
// its signature, such as the answering-writer of the writer of a persister, so
// the effect of the role covers it. A callable without a role takes the effect
// of the writers of the form where the catalog classifies it as a writer or a
// composite-writer. Any other classification has no form: a mixin, a role of
// a contract without a form, roles in two contracts, or another shape without
// a role.
//
// In a form of the state map, each load pairs with each store by the rule of
// [Pairing], takes one key and returns one value beside an error. Each store
// returns nothing or the value that it stores. In a form of the state list,
// each append appends its one input, every append an input of one type, and
// returns an integer position.
func ReferenceOf(decl []Facts) (Reference, *Unreferenced) {
	form, effects, refused := formOf(decl)
	if refused != nil {
		return Reference{}, refused
	}
	ref := Reference{Form: form, Callables: make([]Effected, len(decl))}
	for i := range decl {
		ref.Callables[i] = Effected{Facts: &decl[i], Effect: effects[i], Value: -1}
	}
	if form.State == spec.StateList {
		refused = appends(ref.Callables)
	} else {
		refused = keys(ref.Callables)
	}
	if refused != nil {
		return Reference{}, refused
	}
	return ref, nil
}

// formOf returns the reference form that covers the classifications of the
// callables of decl and the effect of each callable. It judges each callable
// with a role first, because the form of the roles decides whether a callable
// without a role is a writer of the form.
func formOf(decl []Facts) (spec.ReferenceForm, []spec.Effect, *Unreferenced) {
	var form spec.ReferenceForm
	var instance *shape.Membership
	var first *Facts
	effects := make([]spec.Effect, len(decl))
	for i := range decl {
		f := &decl[i]
		c := &f.Class
		switch {
		case len(c.Mixins) > 0:
			return form, nil, unreferenced(f, "has the mixin "+string(c.Mixins[0]))
		case len(c.Contracts) > 1:
			return form, nil, unreferenced(f, fmt.Sprintf("has roles in the contracts %s and %s",
				c.Contracts[0].Contract, c.Contracts[1].Contract))
		case len(c.Contracts) == 0:
			continue
		}
		m := &c.Contracts[0]
		found, ok := referenceOf(m.Contract)
		switch {
		case !ok:
			return form, nil, unreferenced(f, fmt.Sprintf("has the role %s of the contract %s, which has no "+
				"reference form", m.Role, m.Contract))
		case instance != nil && (m.Contract != instance.Contract || m.ID != instance.ID):
			return form, nil, unreferenced(f, "has a role in another contract instance than "+first.Callable)
		case instance == nil:
			form, instance, first = found, m, f
		}
		effects[i] = form.Roles[string(m.Role)]
	}
	for i := range decl {
		f := &decl[i]
		s := f.Class.Shape
		switch {
		case effects[i] != "":
			continue
		case s == "":
			return form, nil, unreferenced(f, "has no classification")
		case s != shape.Writer && s != shape.CompositeWriter:
			return form, nil, unreferenced(f, "has the shape "+string(s))
		case instance == nil:
			return form, nil, unreferenced(f, fmt.Sprintf("has the shape %s, and no callable of %s has a role of a "+
				"contract with a reference form", s, f.Declaration))
		case form.Writers == "":
			return form, nil, unreferenced(f, fmt.Sprintf("has the shape %s, and the reference form %s takes no "+
				"writer", s, form.ID))
		}
		effects[i] = form.Writers
	}
	return form, effects, nil
}

// referenceOf returns the reference form of the contract, and false where the
// spec has none.
func referenceOf(contract shape.Contract) (spec.ReferenceForm, bool) {
	for _, form := range spec.References() {
		if form.Contract == string(contract) {
			return form, true
		}
	}
	return spec.ReferenceForm{}, false
}

// keys binds the key and the value of each store and each load of a form of
// the state map. It returns the refusal of the first callable whose key or
// results the form does not cover.
func keys(cs []Effected) *Unreferenced {
	stores, loads := effected(cs, spec.EffectStore), effected(cs, spec.EffectLoad)
	switch {
	case len(stores) == 0:
		f := loads[0].Facts
		return unreferenced(f, fmt.Sprintf("loads from a map that no callable of %s stores into", f.Declaration))
	case len(loads) == 0:
		f := stores[0].Facts
		return unreferenced(f, fmt.Sprintf("stores into a map that no callable of %s loads from", f.Declaration))
	}
	for _, l := range loads {
		switch f := l.Facts; {
		case len(f.Results) != 1:
			return unreferenced(f, fmt.Sprintf("returns %d values", len(f.Results)))
		case f.Signature.Errors == rules.ErrorsNone:
			return unreferenced(f, "returns no error")
		}
	}
	for si, s := range stores {
		for li, l := range loads {
			key, refused := keyOf(l, s)
			switch {
			case refused != nil:
				return refused
			case li == 0:
				s.Key = key
			case located(key) != located(s.Key):
				return unreferenced(l.Facts, fmt.Sprintf("selects another key of the value that %s stores than %s",
					s.Facts.Callable, loads[0].Facts.Callable))
			}
			if si == 0 {
				l.Key = key
			}
		}
		rs := s.Facts.Results
		if len(rs) > 1 || len(rs) == 1 && !reflect.DeepEqual(rs[0].Shape, s.Facts.Inputs[s.Value].Shape) {
			return unreferenced(s.Facts, "returns another value than the value that it stores")
		}
	}
	return nil
}

// keyOf returns the argument of the key of the load l in the store s, and
// sets the value of s. It returns the refusal of l where its one key selects
// nothing of s.
func keyOf(l, s *Effected) (Argument, *Unreferenced) {
	pair, gate := Pairing(l.Facts, s.Facts)
	if gate != nil {
		return Argument{}, unreferenced(l.Facts, "takes a key that selects nothing of the value that "+
			s.Facts.Callable+" stores")
	}
	if pair.Written < 0 {
		return Argument{}, unreferenced(l.Facts, "returns another type than the value that "+s.Facts.Callable+
			" stores")
	}
	var keys []Argument
	for _, a := range pair.Arguments {
		if a.Kind != ArgContext {
			keys = append(keys, a)
		}
	}
	switch {
	case len(keys) == 0:
		return Argument{}, unreferenced(l.Facts, "takes no key")
	case len(keys) > 1:
		return Argument{}, unreferenced(l.Facts, fmt.Sprintf("takes %d keys", len(keys)))
	case keys[0].Kind == ArgOwn:
		return Argument{}, unreferenced(l.Facts, "takes a function as its key")
	}
	s.Value = pair.Written
	return keys[0], nil
}

// located returns the key a without the name of the parameter of its load, so
// two loads that select one key of a store have one located key.
func located(a Argument) Argument {
	a.Name = ""
	return a
}

// appends binds the value of each append of a form of the state list. It
// returns the refusal of the first append whose inputs or results the form
// does not cover.
func appends(cs []Effected) *Unreferenced {
	for i := range cs {
		f := cs[i].Facts
		switch {
		case len(f.Inputs) != 1:
			return unreferenced(f, fmt.Sprintf("appends %d inputs", len(f.Inputs)))
		case !reflect.DeepEqual(f.Inputs[0].Shape, cs[0].Facts.Inputs[0].Shape):
			return unreferenced(f, "appends another type than "+cs[0].Facts.Callable)
		case len(f.Results) != 1 || !integer(&f.Results[0].Shape):
			return unreferenced(f, "returns no integer position")
		}
		cs[i].Value = 0
	}
	return nil
}

// integer reports whether s is the shape of an integer type. Only the shape
// of a number has a class.
func integer(s *rules.TypeShape) bool {
	return s.Class == rules.ScalarInt || s.Class == rules.ScalarUint
}

// effected returns the callables of cs with the effect, in order.
func effected(cs []Effected, effect spec.Effect) []*Effected {
	var out []*Effected
	for i := range cs {
		if cs[i].Effect == effect {
			out = append(out, &cs[i])
		}
	}
	return out
}

// unreferenced returns the refusal of the gate reference-role on the callable
// of f, with the reason and the remedy of the gate and the detail of what of
// the callable no form covers.
func unreferenced(f *Facts, detail string) *Unreferenced {
	return &Unreferenced{Refusal: refusal("", f, f.Key(), gateOf(spec.ReferenceRole), namesOf(f)), Detail: detail}
}
