// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// runtime contains the machineries that the runtime supplies to every check.
var runtime = func() []spec.MachineryID {
	var out []spec.MachineryID
	for _, m := range spec.Machineries() {
		if m.By == spec.ByRuntime {
			out = append(out, m.ID)
		}
	}
	return out
}()

// derivation records into s the facts of each callable of the suite and the
// checks that the claims of the conformance spec derive from them. syms are
// the methods or the function of the callables, in the order of s.Callables.
// The suite of a function has no declaration beside the function.
func (k *caller) derivation(s *Suite, syms []symbol.Symbol) {
	for _, sym := range syms {
		f := k.factsOf(sym)
		if s.Func == nil {
			f.Declaration = s.Declaration
		}
		s.Facts = append(s.Facts, f)
	}
	// The callables of one declaration have distinct names, so their keys
	// differ and Derive returns no error.
	s.Derivation, _ = derive.Derive(s.Facts, spec.Claims(), nil, runtime)
}

// factsOf returns the facts of the method or the function sym. Each input and
// each result in the role value has the shape, the samples, the zero value
// and the fields of its type at the witnesses of the declaration. The facts
// have the classification that the shape catalog stamps on sym.
func (k *caller) factsOf(sym symbol.Symbol) derive.Facts {
	view, _ := k.b.CallableOf(sym)
	decl, _ := sym.(node.Declaration)
	id := decl.Identity()
	f := derive.Facts{Callable: view.Name, ID: id, Class: shape.OfDeclaration(k.m, id), Signature: view}
	var params []*node.Param
	switch d := sym.(type) {
	case *node.Method:
		params = d.Params
	case *node.Function:
		params = d.Params
	}
	for _, p := range view.Params {
		if p.Role != rules.ParamInput {
			continue
		}
		var subject symbol.Identity
		for _, n := range params {
			if n.Type == p.Ref {
				subject = n.ID
			}
		}
		f.Inputs = append(f.Inputs, k.valueOf(&subject, p.Name, p.Ref))
	}
	var none symbol.Identity
	for _, r := range view.Returns {
		if r.Role == rules.ReturnValue {
			f.Results = append(f.Results, k.valueOf(&none, r.Name, r.Ref))
		}
	}
	return f
}

// valueOf returns the value named name of the type ref, with what eidos
// derives for the type. subject is the declaration of a parameter, and zero
// for a return. The first sample is the sample that an author states on
// subject or on the type, else the declared default of the type, else the
// sample that eidos derives. The second sample is the authored alternate, else
// a derived sample that differs from the first.
func (k *caller) valueOf(subject *symbol.Identity, name string, ref *node.TypeRef) derive.Value {
	ref = k.instance(ref)
	v := derive.Value{Name: name, Shape: k.b.TypeOf(ref), Fields: k.fieldsOf(ref)}
	v.Zero, v.HasZero = k.b.ZeroValue(ref)
	sample, alternate := k.b.View().Authored(k.b.Source(), *subject, ref)
	if !sample.OK() && sample.Refusal == rules.RefusedNone {
		if value, declared := k.defaultOf(ref); declared {
			sample = rules.Of(value)
		}
	}
	derived, derivedAlternate := k.b.Source().DeriveSamples(ref, name, k.b.View())
	v.Samples[0], v.Samples[1] = rules.Complete(sample, alternate, derived, derivedAlternate)
	return v
}

// defaultOf returns the declared default of the type ref: the default of a
// defined type or an alias, or the value of a struct with the declared
// default of each field that declares one.
func (k *caller) defaultOf(ref *node.TypeRef) (emit.Value, bool) {
	sym, _ := k.m.Reader().Lookup(ref.Target)
	switch d := sym.(type) {
	case *node.Alias:
		text, declared := sdk.FactOf(k.m, d.ID, meta.Named[string](defaults.KeyValue))
		if !declared {
			return emit.Value{}, false
		}
		scope := k.b.View().ScopeOf(d.ID)
		return defaults.Value(&k.b, &scope, ref, text)
	case *node.Struct:
		return k.defaultsOf(d, ref)
	default:
		return emit.Value{}, false
	}
}

// fieldsOf returns the exported fields of the struct that ref refers to, or
// that the pointer ref points to, with the shape of each.
func (k *caller) fieldsOf(ref *node.TypeRef) []derive.Field {
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	sym, _ := k.m.Reader().Lookup(ref.Target)
	st, isStruct := sym.(*node.Struct)
	if !isStruct {
		return nil
	}
	var out []derive.Field
	for _, f := range st.Fields {
		if f.Visibility == symbol.VisibilityPublic {
			out = append(out, derive.Field{Name: f.Name, Shape: k.b.FieldTypeOf(f)})
		}
	}
	return out
}
