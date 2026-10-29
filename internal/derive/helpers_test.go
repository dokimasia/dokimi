// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
)

// golang is the language of the raw values of the samples.
const golang symbol.Lang = "go"

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
