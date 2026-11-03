// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive

import (
	"reflect"
	"slices"
	"strings"

	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/spec"
)

// ArgumentKind is the source of one argument of the read of a pair.
type ArgumentKind uint8

// The sources of an argument of the read.
const (
	// ArgContext is the context of the check.
	ArgContext ArgumentKind = iota + 1

	// ArgShared is the sample of the parameter of the write with the same
	// name and type, which both calls take.
	ArgShared

	// ArgField is the field of the written value with the name and the type
	// of the parameter.
	ArgField
	// ArgOwn is the sample of the parameter of the read: a function, which
	// the probe replaces with the function that it supplies, such as the
	// body of a scope, or a parameter of a read whose write takes no input.
	ArgOwn
)

// Argument is one argument of the read of a pair.
type Argument struct {
	// Name is the name of the parameter of the read.
	Name string

	// Param is the parameter of the write of an ArgShared argument.
	Param string

	// Field is the field of the written value of an ArgField argument.
	Field string

	Kind ArgumentKind
}

// Pair is the partner of a claim on a callable, with the binding of the
// arguments where one of the two callables reads what the other writes.
type Pair struct {
	// Partner is the facts of the partner.
	Partner *Facts

	// Read and Write are the callable that reads and the callable that
	// writes, and nil where neither of the two writes.
	Read, Write *Facts

	// Changed is the field of the written value that the second write of
	// the modes interposed and rewriter changes.
	Changed string

	// Arguments are the arguments of the read, in the order of its
	// parameters.
	Arguments []Argument

	// Written is the index in the inputs of the write of the written value,
	// and -1 where the write takes no value of the type of the read.
	Written int
}

// Pairing binds the arguments of read to write by the rule of [Partner], with
// write as the partner of the pair. It returns the gate partner-arguments
// instead when a parameter of read selects nothing of write, beside the pair
// of the partner without a read and a write.
func Pairing(read, write *Facts) (*Pair, *spec.Gate) {
	return paired(spec.ModeNamed, write, read, write)
}

// Partner resolves the partner p of a claim on the callable of f, among the
// callables of decl, by the mode of p, and binds the arguments of the read.
// It returns the gate that refuses the claim instead: the gate partner-found
// when no callable matches, one-partner when two do, partner-arguments when
// a parameter of the read selects nothing of the write, and rewritable when
// a second write could not change what the read returns.
//
// Every mode but establishing and paired resolves the partner before it binds
// the arguments. For those modes, Partner returns the pair of the partner
// without a read and a write beside the gate partner-arguments. A claim whose
// probe calls the partner by its name takes that pair. The mode paired takes
// the one other callable whose read pairs with the callable of f by the rule
// of [Pairing], and whose result is the value that f writes.
func Partner(p *spec.Partner, f *Facts, decl []Facts) (*Pair, *spec.Gate) {
	switch p.Mode {
	case spec.ModeNamed, spec.ModeInterposed:
		partner, gate := named(p.Param, f, decl)
		if gate != nil {
			return nil, gate
		}
		pair, gate := paired(p.Mode, partner, f, partner)
		if gate != nil && p.Mode == spec.ModeNamed && len(f.Results) == 0 && len(partner.Results) > 0 {
			// A callable without a value whose parameters select nothing of
			// its named partner writes what the partner reads, as Add writes
			// the total that Total observes.
			return paired(p.Mode, partner, partner, f)
		}
		return pair, gate
	case spec.ModeReader, spec.ModeRewriter:
		partner, gate := named(p.Param, f, decl)
		if p.Param == "" {
			partner, gate = role(p, f, decl)
		}
		if gate != nil {
			return nil, gate
		}
		return paired(p.Mode, partner, partner, f)
	case spec.ModeDerived:
		if fieldOf(p.Param, f) == nil {
			return nil, gateOf(spec.PartnerFound)
		}
		// The write takes the value whose field the param names, as Put
		// takes the Item whose field Rev the reads compare.
		value := f.Results[0].Shape
		partner, gate := one(f, decl, func(c *Facts) bool {
			return slices.ContainsFunc(c.Inputs, func(v Value) bool { return reflect.DeepEqual(v.Shape, value) })
		})
		if gate != nil {
			return nil, gate
		}
		return paired(p.Mode, partner, f, partner)
	case spec.ModePaired:
		// A read that does not pair with f has no written value.
		partner, gate := one(f, decl, func(c *Facts) bool {
			pair, _ := Pairing(c, f)
			return pair.Written >= 0
		})
		if gate != nil {
			return nil, gate
		}
		return paired(p.Mode, partner, partner, f)
	case spec.ModeRole:
		partner, gate := role(p, f, decl)
		if gate != nil {
			return nil, gate
		}
		switch {
		case len(f.Results) > 0 && len(partner.Results) == 0:
			return paired(p.Mode, partner, f, partner)
		case len(f.Results) == 0 && len(partner.Results) > 0:
			return paired(p.Mode, partner, partner, f)
		case takes(partner, f) && !takes(f, partner):
			return paired(p.Mode, partner, f, partner)
		case takes(f, partner) && !takes(partner, f):
			return paired(p.Mode, partner, partner, f)
		default:
			return &Pair{Partner: partner, Written: -1}, nil
		}
	default:
		return establishing(p, f, decl)
	}
}

// takes reports whether the callable of write takes an input of the type of
// the first result of the callable of read, as Put(ctx, v Value) takes what
// Get(ctx, key) returns.
func takes(write, read *Facts) bool {
	return len(read.Results) > 0 && slices.ContainsFunc(write.Inputs, func(v Value) bool {
		return reflect.DeepEqual(v.Shape, read.Results[0].Shape)
	})
}

// named returns the callable of decl that the param key of a classification
// of f names.
func named(key string, f *Facts, decl []Facts) (*Facts, *spec.Gate) {
	id, ok := paramID(key, f)
	if ok {
		for i := range decl {
			if decl[i].ID == id {
				return &decl[i], nil
			}
		}
	}
	return nil, gateOf(spec.PartnerFound)
}

// fieldOf returns the field of the value that f returns, which the param key
// of a classification of f names. It returns nil where f has no such param or
// no such field. Every field has a name, so the zero identity of a missing
// param names no field.
func fieldOf(key string, f *Facts) *Field {
	if len(f.Results) == 0 {
		return nil
	}
	id, _ := paramID(key, f)
	fields := f.Results[0].Fields
	for i := range fields {
		if fields[i].Name == id.Name {
			return &fields[i]
		}
	}
	return nil
}

// paramID returns the declaration that the param key of a classification of
// f resolves to.
func paramID(key string, f *Facts) (symbol.Identity, bool) {
	for _, params := range f.Class.Params {
		if id, ok := params[directive.ParamKey(key)].(symbol.Identity); ok {
			return id, true
		}
	}
	return symbol.Identity{}, false
}

// role returns the callable of decl with the role of p in a contract
// instance of f.
func role(p *spec.Partner, f *Facts, decl []Facts) (*Facts, *spec.Gate) {
	return one(f, decl, func(c *Facts) bool {
		for _, m := range f.Class.Contracts {
			for _, n := range c.Class.Contracts {
				if n.Contract == m.Contract && n.ID == m.ID && n.Role == p.Role {
					return true
				}
			}
		}
		return false
	})
}

// establishing returns the pair of the one other callable of decl with the
// classification of p whose arguments the read of f binds.
func establishing(p *spec.Partner, f *Facts, decl []Facts) (*Pair, *spec.Gate) {
	var found *Pair
	classified := false
	for i := range decl {
		c := &decl[i]
		if c == f || !hasClassification(c, p) {
			continue
		}
		classified = true
		pair, gate := paired(p.Mode, c, f, c)
		if gate != nil {
			continue
		}
		if found != nil {
			return nil, gateOf(spec.OnePartner)
		}
		found = pair
	}
	switch {
	case found != nil:
		return found, nil
	case classified:
		return nil, gateOf(spec.PartnerArguments)
	default:
		return nil, gateOf(spec.PartnerFound)
	}
}

// one returns the one other callable of decl that match accepts.
func one(f *Facts, decl []Facts, match func(c *Facts) bool) (*Facts, *spec.Gate) {
	var found *Facts
	for i := range decl {
		c := &decl[i]
		if c == f || !match(c) {
			continue
		}
		if found != nil {
			return nil, gateOf(spec.OnePartner)
		}
		found = c
	}
	if found == nil {
		return nil, gateOf(spec.PartnerFound)
	}
	return found, nil
}

// hasClassification reports whether c has the shape or the mixin of p.
func hasClassification(c *Facts, p *spec.Partner) bool {
	if p.Shape != "" {
		return c.Class.Shape == p.Shape
	}
	return slices.Contains(c.Class.Mixins, p.Mixin)
}

// paired returns the pair of partner where read reads what write writes,
// with the arguments of the read bound by the rule of the derivation: a
// parameter of the same name and type in both calls shares one sample, and
// any other parameter takes the field of the written value with its name
// and type.
//
// A batch read returns the list of the values that the write writes one at a
// time, and selects them by the keys of its variadic input. The keys select
// the field of the written value of their name, or are the written values
// themselves. The second write changes that field or the whole value, so the
// two writes store two keys.
func paired(mode spec.PartnerMode, partner, read, write *Facts) (*Pair, *spec.Gate) {
	pair := &Pair{Partner: partner, Read: read, Write: write, Written: -1}
	batch := batched(read)
	if len(read.Results) > 0 {
		value := read.Results[0].Shape
		if batch {
			value = value.Elems[0]
		}
		// An input that the read takes by its name is a selector, such as the
		// key of Write(key, value), so it is not the value. A read input of
		// that name and another type binds nothing, so its claim refuses
		// whichever input the write writes.
		pair.Written = slices.IndexFunc(write.Inputs, func(v Value) bool {
			shared := slices.ContainsFunc(read.Inputs, func(r Value) bool { return r.Name == v.Name })
			return !shared && reflect.DeepEqual(v.Shape, value)
		})
	}
	inputs := read.Inputs
	for _, param := range read.Signature.Params {
		if param.Role == rules.ParamContext {
			pair.Arguments = append(pair.Arguments, Argument{Name: param.Name, Kind: ArgContext})
			continue
		}
		arg, ok := argumentOf(&inputs[0], param.Variadic, write, pair.Written)
		if !ok {
			return unpaired(partner), gateOf(spec.PartnerArguments)
		}
		pair.Arguments = append(pair.Arguments, arg)
		inputs = inputs[1:]
	}
	if batch {
		keys := &pair.Arguments[len(pair.Arguments)-1]
		if keys.Kind != ArgField && (pair.Written < 0 || keys.Param != write.Inputs[pair.Written].Name) {
			return unpaired(partner), gateOf(spec.PartnerArguments)
		}
		pair.Changed = keys.Field
	}
	if mode != spec.ModeInterposed && mode != spec.ModeRewriter {
		return pair, nil
	}
	pair.Changed = changed(pair, write)
	if pair.Changed == "" {
		return nil, gateOf(spec.Rewritable)
	}
	return pair, nil
}

// unpaired returns the pair of partner without a read and a write, whose
// arguments do not bind.
func unpaired(partner *Facts) *Pair { return &Pair{Partner: partner, Written: -1} }

// argumentOf binds the input value of the read. The inputs of a callable
// are its parameters in the role input, in the order of its parameters. A
// value of a function type, and any value of a read whose write takes no
// input, takes its own sample, because nothing of the write selects it. A
// variadic input takes several values of the name without its last letter s,
// as keys takes several keys.
func argumentOf(value *Value, variadic bool, write *Facts, written int) (Argument, bool) {
	if value.Shape.Form == symbol.FormFunc || len(write.Inputs) == 0 {
		return Argument{Name: value.Name, Kind: ArgOwn}, true
	}
	name := value.Name
	if variadic {
		name = strings.TrimSuffix(name, "s")
	}
	shared := slices.IndexFunc(write.Inputs, func(v Value) bool {
		return v.Name == name && reflect.DeepEqual(v.Shape, value.Shape)
	})
	if shared >= 0 {
		return Argument{Name: value.Name, Kind: ArgShared, Param: write.Inputs[shared].Name}, true
	}
	if written < 0 {
		return Argument{}, false
	}
	fields := write.Inputs[written].Fields
	field := slices.IndexFunc(fields, func(fd Field) bool {
		return strings.EqualFold(fd.Name, name) && reflect.DeepEqual(fd.Shape, value.Shape)
	})
	if field < 0 {
		return Argument{}, false
	}
	return Argument{Name: value.Name, Kind: ArgField, Field: fields[field].Name}, true
}

// batched reports whether the callable of f reads in a batch: it takes a
// variadic input and returns a list as its first value.
func batched(f *Facts) bool {
	return len(f.Results) > 0 && f.Results[0].Shape.Form == symbol.FormList &&
		slices.ContainsFunc(f.Signature.Params, func(p rules.ParamView) bool { return p.Variadic })
}

// changed returns the first field of the written value that no argument of
// the read selects by, and the empty string where every field is selected.
func changed(pair *Pair, write *Facts) string {
	if pair.Written < 0 {
		return ""
	}
	fields := write.Inputs[pair.Written].Fields
	for i := range fields {
		name := fields[i].Name
		if !slices.ContainsFunc(pair.Arguments, func(a Argument) bool { return a.Field == name }) {
			return name
		}
	}
	return ""
}

// gateOf returns the gate id of the spec.
func gateOf(id spec.GateID) *spec.Gate {
	g, _ := spec.GateOf(id)
	return &g
}
