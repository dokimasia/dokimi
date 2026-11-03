// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
)

// golang is the language of the raw values of the samples.
const golang symbol.Lang = "go"

// idOf returns the identity of the method name of Store.
func idOf(name string) symbol.Identity {
	return symbol.Identity{Lang: golang, Package: "example.com/store", Owner: "Store", Name: name}
}

// text returns a value of a string type with two distinct samples.
func text(name string) derive.Value {
	return derive.Value{
		Name:  name,
		Shape: rules.TypeShape{Form: symbol.FormText, Spelling: "string"},
		Samples: [2]rules.Sample{
			rules.Of(emit.Literal(emit.LiteralString, "a")),
			rules.Of(emit.Literal(emit.LiteralString, "b")),
		},
		Zero:    emit.Literal(emit.LiteralString, ""),
		HasZero: true,
	}
}

// item returns a value of the struct Item, whose fields are ID and Name,
// with two distinct samples.
func item(name string) derive.Value {
	return derive.Value{
		Name:  name,
		Shape: rules.TypeShape{Form: symbol.FormReference, Spelling: "Item"},
		Samples: [2]rules.Sample{
			rules.Of(emit.Raw(golang, `Item{ID: "a"}`)),
			rules.Of(emit.Raw(golang, `Item{ID: "b"}`)),
		},
		Zero:    emit.Raw(golang, "Item{}"),
		HasZero: true,
		Fields: []derive.Field{
			{Name: "ID", Shape: rules.TypeShape{Form: symbol.FormText, Spelling: "string"}},
			{Name: "Name", Shape: rules.TypeShape{Form: symbol.FormText, Spelling: "string"}},
		},
	}
}

// signature returns the view of a callable with a context, the inputs and
// the returns of the roles, and the error model.
func signature(
	name string,
	ctx bool,
	inputs []derive.Value,
	returns []rules.ReturnRole,
	errs rules.ErrorModel,
) rules.Callable {
	c := rules.Callable{Name: name, Errors: errs}
	if ctx {
		c.Params = append(c.Params, rules.ParamView{Name: "ctx", Role: rules.ParamContext})
	}
	for _, in := range inputs {
		c.Params = append(c.Params, rules.ParamView{Name: in.Name, Role: rules.ParamInput})
	}
	for _, r := range returns {
		c.Returns = append(c.Returns, rules.ReturnView{Role: r})
	}
	return c
}

// reader returns the facts of Store.Read(ctx, key string) (Item, error).
func reader() derive.Facts {
	inputs := []derive.Value{text("key")}
	return derive.Facts{
		Declaration: "Store",
		Callable:    "Read",
		Signature: signature(
			"Read",
			true,
			inputs,
			[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError},
			rules.ErrorsLastReturn,
		),
		Inputs:  inputs,
		Results: []derive.Value{item("")},
	}
}

// writer returns the facts of Store.Put(ctx, it Item) error.
func writer() derive.Facts {
	inputs := []derive.Value{item("it")}
	return derive.Facts{
		Declaration: "Store", Callable: "Put",
		Signature: signature("Put", true, inputs, []rules.ReturnRole{rules.ReturnError}, rules.ErrorsLastReturn),
		Inputs:    inputs,
	}
}

// lookup returns the facts of Store.Lookup(ctx, key string) (Item, string,
// bool).
func lookup() derive.Facts {
	inputs := []derive.Value{text("key")}
	return derive.Facts{
		Declaration: "Store", Callable: "Lookup",
		Signature: signature("Lookup", true, inputs,
			[]rules.ReturnRole{rules.ReturnValue, rules.ReturnValue, rules.ReturnOkBool}, rules.ErrorsNone),
		Inputs:  inputs,
		Results: []derive.Value{item(""), text("")},
	}
}

// pure returns the facts of Store.Hash(key string) string.
func pure() derive.Facts {
	inputs := []derive.Value{text("key")}
	return derive.Facts{
		Declaration: "Store", Callable: "Hash",
		Signature: signature("Hash", false, inputs, []rules.ReturnRole{rules.ReturnValue}, rules.ErrorsNone),
		Inputs:    inputs,
		Results:   []derive.Value{text("")},
	}
}

// combiner returns the facts of Store.Merge(first, second Item) Item.
func combiner() derive.Facts {
	inputs := []derive.Value{item("first"), item("second")}
	return derive.Facts{
		Declaration: "Store", Callable: "Merge",
		Signature: signature("Merge", false, inputs, []rules.ReturnRole{rules.ReturnValue}, rules.ErrorsNone),
		Inputs:    inputs,
		Results:   []derive.Value{item("")},
	}
}

// stream returns the facts of Store.Watch(ctx) (<-chan Item, error).
func stream() derive.Facts {
	return derive.Facts{
		Declaration: "Store",
		Callable:    "Watch",
		Signature: signature(
			"Watch",
			true,
			nil,
			[]rules.ReturnRole{rules.ReturnStream, rules.ReturnError},
			rules.ErrorsLastReturn,
		),
	}
}

// getter returns the facts of Store.Get(ctx, id string) (Item, error).
func getter() derive.Facts {
	inputs := []derive.Value{text("id")}
	return derive.Facts{
		Declaration: "Store",
		Callable:    "Get",
		ID:          idOf("Get"),
		Signature: signature(
			"Get",
			true,
			inputs,
			[]rules.ReturnRole{rules.ReturnValue, rules.ReturnError},
			rules.ErrorsLastReturn,
		),
		Inputs:  inputs,
		Results: []derive.Value{item("")},
	}
}

// renamed returns f as the callable name of Store.
func renamed(f derive.Facts, name string) derive.Facts {
	f.Callable, f.ID, f.Signature.Name = name, idOf(name), name
	return f
}

// putter returns the facts of Store.Put(ctx, it Item) error with its
// identity.
func putter() derive.Facts {
	f := writer()
	f.ID = idOf("Put")
	return f
}

// keyedPutter returns the facts of Store.Set(ctx, id string, it Item) error.
func keyedPutter() derive.Facts {
	inputs := []derive.Value{text("id"), item("it")}
	return derive.Facts{
		Declaration: "Store", Callable: "Set", ID: idOf("Set"),
		Signature: signature("Set", true, inputs, []rules.ReturnRole{rules.ReturnError}, rules.ErrorsLastReturn),
		Inputs:    inputs,
	}
}

// marker returns the facts of Store.Mark(ctx, id string) error, which takes
// no value of the read.
func marker() derive.Facts {
	inputs := []derive.Value{text("id")}
	return derive.Facts{
		Declaration: "Store", Callable: "Mark", ID: idOf("Mark"),
		Signature: signature("Mark", true, inputs, []rules.ReturnRole{rules.ReturnError}, rules.ErrorsLastReturn),
		Inputs:    inputs,
	}
}

// withParam returns f with the param key of the classification name set to
// value.
func withParam(f derive.Facts, name, key string, value any) derive.Facts {
	f.Class.Params = map[string]shape.Params{name: {directiveKey(key): value}}
	return f
}

// directiveKey returns the key of a param of a classification.
func directiveKey(key string) directive.ParamKey { return directive.ParamKey(key) }

// withShape returns f with the shape s.
func withShape(f derive.Facts, s shape.Shape) derive.Facts {
	f.Class.Shape = s
	return f
}

// withRole returns f with the role of an instance of the contract.
func withRole(f derive.Facts, contract shape.Contract, role shape.Role) derive.Facts {
	f.Class.Contracts = append(f.Class.Contracts, shape.Membership{Contract: contract, Role: role})
	return f
}

// withMixins returns f with the mixins.
func withMixins(f derive.Facts, mixins ...shape.Mixin) derive.Facts {
	f.Class.Mixins = mixins
	return f
}

// withErrors returns f with the error model errs.
func withErrors(f derive.Facts, errs rules.ErrorModel) derive.Facts {
	f.Signature.Errors = errs
	return f
}

// withResult returns f with a first result of the form.
func withResult(f derive.Facts, form symbol.TypeForm) derive.Facts {
	f.Results = append([]derive.Value(nil), f.Results...)
	f.Results[0].Shape = rules.TypeShape{Form: form}
	return f
}

// withInput returns f with a first input of the form.
func withInput(f derive.Facts, form symbol.TypeForm) derive.Facts {
	f.Inputs = append([]derive.Value(nil), f.Inputs...)
	f.Inputs[0].Shape = rules.TypeShape{Form: form}
	return f
}

// withSameSamples returns f with a first result whose two samples are equal.
func withSameSamples(f derive.Facts) derive.Facts {
	f.Results = append([]derive.Value(nil), f.Results...)
	f.Results[0].Samples[1] = f.Results[0].Samples[0]
	return f
}

// withRefusedSample returns f with a first result whose second sample is
// refused.
func withRefusedSample(f derive.Facts) derive.Facts {
	f.Results = append([]derive.Value(nil), f.Results...)
	f.Results[0].Samples[1] = rules.Refused(rules.RefusedNoLiteral)
	return f
}
