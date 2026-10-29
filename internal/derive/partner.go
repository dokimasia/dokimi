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

// Partner resolves the partner p of a claim on the callable of f, among the
// callables of decl, by the mode of p, and binds the arguments of the read.
// It returns the gate that refuses the claim instead: the gate partner-found
// when no callable matches, one-partner when two do, partner-arguments when
// a parameter of the read selects nothing of the write, and rewritable when
// a second write could not change what the read returns.
func Partner(p *spec.Partner, f *Facts, decl []Facts) (*Pair, *spec.Gate) {
	switch p.Mode {
	case spec.ModeNamed, spec.ModeInterposed:
		partner, gate := named(p.Param, f, decl)
		if gate != nil {
			return nil, gate
		}
		return paired(p.Mode, partner, f, partner)
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
		field := fieldOf(p.Param, f)
		if field == nil {
			return nil, gateOf(spec.PartnerFound)
		}
		partner, gate := one(f, decl, func(c *Facts) bool {
			return slices.ContainsFunc(c.Inputs, func(v Value) bool { return reflect.DeepEqual(v.Shape, field.Shape) })
		})
		if gate != nil {
			return nil, gate
		}
		return paired(p.Mode, partner, f, partner)
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
		default:
			return &Pair{Partner: partner, Written: -1}, nil
		}
	default:
		return establishing(p, f, decl)
	}
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
func paired(mode spec.PartnerMode, partner, read, write *Facts) (*Pair, *spec.Gate) {
	pair := &Pair{Partner: partner, Read: read, Write: write, Written: -1}
	if len(read.Results) > 0 {
		pair.Written = slices.IndexFunc(write.Inputs, func(v Value) bool {
			return reflect.DeepEqual(v.Shape, read.Results[0].Shape)
		})
	}
	inputs := read.Inputs
	for _, param := range read.Signature.Params {
		if param.Role == rules.ParamContext {
			pair.Arguments = append(pair.Arguments, Argument{Name: param.Name, Kind: ArgContext})
			continue
		}
		arg, ok := argumentOf(&inputs[0], write, pair.Written)
		if !ok {
			return nil, gateOf(spec.PartnerArguments)
		}
		pair.Arguments = append(pair.Arguments, arg)
		inputs = inputs[1:]
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

// argumentOf binds the input value of the read. The inputs of a callable
// are its parameters in the role input, in the order of its parameters.
func argumentOf(value *Value, write *Facts, written int) (Argument, bool) {
	shared := slices.IndexFunc(write.Inputs, func(v Value) bool {
		return v.Name == value.Name && reflect.DeepEqual(v.Shape, value.Shape)
	})
	if shared >= 0 {
		return Argument{Name: value.Name, Kind: ArgShared, Param: write.Inputs[shared].Name}, true
	}
	if written < 0 {
		return Argument{}, false
	}
	fields := write.Inputs[written].Fields
	field := slices.IndexFunc(fields, func(fd Field) bool {
		return strings.EqualFold(fd.Name, value.Name) && reflect.DeepEqual(fd.Shape, value.Shape)
	})
	if field < 0 {
		return Argument{}, false
	}
	return Argument{Name: value.Name, Kind: ArgField, Field: fields[field].Name}, true
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
