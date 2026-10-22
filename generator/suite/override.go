// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"go/token"
	"strings"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/generator/builder"
)

// Enumeration is one enumeration that the calls of a suite take at some depth
// of their inputs. Members are the members of the enumeration in the order of
// its declaration. Each property run and each benchmark of the suite draws one
// of them for each value of the enumeration.
type Enumeration struct {
	Members []emit.Value
}

// Start is the generator of the values of a struct with a starting value,
// which each property run and each benchmark of the suite takes for every
// value of the struct. A value of the generator keeps or changes each field of
// the starting value, by a choice of the case, so a failure shrinks back to
// the starting value.
//
// Name is the name of the variable of the generator, as StoreItemValues, and
// the Go plan writes it as storeItemValues. Type is the struct, and Label the
// label of the draw of the changed values. Seed is the constructor of the
// builder of the struct, whose value is the starting value, and nil for a
// struct without a builder or with an authored sample. Value is the starting
// value of any other struct: its authored sample, or its value with the
// declared default of each field. Fields are the fields that a case keeps or
// changes.
//
// Whole reports a defined type that is no struct, whose value a case keeps or
// changes whole. Its starting value is its authored sample or its declared
// default.
type Start struct {
	Type   *emit.TypeRef
	Seed   *emit.TypeRef
	Value  emit.Value
	Name   string
	Label  string
	Fields []StartField
	Whole  bool
}

// StartField is one field of a [Start] that each case keeps or changes. Name
// is the name of the field, and Label the label of the choice of the case.
type StartField struct {
	Name  string
	Label string
}

// valuesWord is the word that the name of the generator of a struct with a
// starting value joins to the struct and the declaration of the suite.
const valuesWord = "Values"

// The key and the value of the tag of prop that leaves a field out of the
// shape of its struct.
const (
	propTag  = "prop"
	propSkip = "-"
)

// overrides returns each enumeration and each struct with a starting value
// that the roots of the calls take at some depth, in the order in which a
// walk of the roots meets them. The walk enters the elements and the type
// arguments of each type, the exported fields and the embedded types of a
// struct, and the target of an alias. It leaves the types of the standard
// library to the shape reader of the property engine, because the value sets
// of the standard library include sets of flags, such as fs.FileMode, and
// units, such as time.Duration.
func (k *caller) overrides() ([]Enumeration, []Start) {
	w := walk{k: k, seen: map[symbol.Identity]bool{}, names: map[string]bool{}}
	for _, root := range k.roots {
		w.visit(root)
	}
	return w.enums, w.starts
}

// walk is a walk of the types of the roots of the calls of a suite. seen
// contains each declaration that the walk visited, and names the names of
// the generators of the starting values. enums and starts are the
// enumerations and the structs with a starting value that the walk met.
type walk struct {
	k      *caller
	seen   map[symbol.Identity]bool
	names  map[string]bool
	enums  []Enumeration
	starts []Start
}

// visit visits the type ref and each type that it contains. A structural form
// has no target, and [golang.Standard] reports the empty package of such a
// target as a package of the standard library, so the walk declares nothing
// for it.
func (w *walk) visit(ref *node.TypeRef) {
	for _, elem := range ref.Elems {
		w.visit(elem)
	}
	for _, arg := range ref.Args {
		w.visit(arg)
	}
	if w.seen[ref.Target] || golang.Standard(ref.Target.Package) {
		return
	}
	w.seen[ref.Target] = true
	decl, _ := w.k.m.Reader().Lookup(ref.Target)
	switch d := decl.(type) {
	case *node.Enum:
		members := make([]emit.Value, 0, len(d.Variants))
		for _, v := range d.Variants {
			members = append(members, emit.Name(v.ID))
		}
		w.enums = append(w.enums, Enumeration{Members: members})
	case *node.Struct:
		if start, started := w.k.startOf(d, ref, w.names); started {
			w.starts = append(w.starts, start)
		}
		for _, f := range d.Fields {
			if f.Visibility == symbol.VisibilityPublic {
				w.visit(rules.Substitute(f.Type, d.TypeParams, ref.Args))
			}
		}
		for _, embed := range d.Embeds {
			w.visit(rules.Substitute(embed.Ref, d.TypeParams, ref.Args))
		}
	case *node.Alias:
		if start, started := w.k.wholeStartOf(d, ref, w.names); started {
			w.starts = append(w.starts, start)
		}
		w.visit(rules.Substitute(d.Target, d.TypeParams, ref.Args))
	}
}

// startOf returns the generator of the values of the struct st at the type
// ref, and false for a struct without a starting value. The starting value is
// the authored sample of the struct, the value of its builder, or its value
// with the declared defaults of its fields, in that order. A generic struct
// takes the shape of the property engine alone. taken contains the names of
// the generators of the suite, and startOf adds the name that it takes.
func (k *caller) startOf(st *node.Struct, ref *node.TypeRef, taken map[string]bool) (Start, bool) {
	if len(st.TypeParams) > 0 {
		return Start{}, false
	}
	s := Start{Type: rules.EmitRef(ref)}
	if sample, _ := k.b.View().Authored(k.b.Source(), st.ID, ref); sample.OK() {
		s.Value = sample.Value
	} else if seed, built := sdk.FactOf(k.m, st.ID, meta.Named[string](builder.KeySeed)); built {
		s.Seed = &emit.TypeRef{Spelling: seed, Target: st.ID, Emitted: true}
	} else if value, defaulted := k.defaultsOf(st, ref); defaulted {
		s.Value = value
	} else {
		return Start{}, false
	}
	k.name(&s, st.Name, taken)
	for _, f := range st.Fields {
		if tag, _ := k.b.Tag(f, propTag); f.Visibility != symbol.VisibilityPublic || tag == propSkip {
			continue
		}
		s.Fields = append(s.Fields, StartField{Name: f.Name, Label: s.Label + "." + naming.Kebab(f.Name)})
	}
	for _, embed := range st.Embeds {
		name := strings.TrimPrefix(embed.Ref.Spelling[strings.LastIndex(embed.Ref.Spelling, ".")+1:], "*")
		if token.IsExported(name) {
			s.Fields = append(s.Fields, StartField{Name: name, Label: s.Label + "." + naming.Kebab(name)})
		}
	}
	return s, true
}

// wholeStartOf returns the generator of the values of the defined type d at
// the type ref, whose value a case keeps or changes whole, and false for a
// type without a starting value. The starting value is the authored sample of
// the type, or its declared default. A transparent alias is the type of its
// target, and a generic type takes the shape of the property engine alone.
// taken contains the names of the generators of the suite, and wholeStartOf
// adds the name that it takes.
func (k *caller) wholeStartOf(d *node.Alias, ref *node.TypeRef, taken map[string]bool) (Start, bool) {
	if !d.Defined || len(d.TypeParams) > 0 {
		return Start{}, false
	}
	s := Start{Type: rules.EmitRef(ref), Whole: true}
	if sample, _ := k.b.View().Authored(k.b.Source(), d.ID, ref); sample.OK() {
		s.Value = sample.Value
	} else if text, declared := sdk.FactOf(k.m, d.ID, meta.Named[string](defaults.KeyValue)); declared {
		scope := k.b.View().ScopeOf(d.ID)
		value, read := defaults.Value(&k.b, &scope, ref, text)
		if !read {
			return Start{}, false
		}
		s.Value = value
	} else {
		return Start{}, false
	}
	k.name(&s, d.Name, taken)
	return s, true
}

// name gives the generator s of the values of the type named typ its name,
// as StoreItemValues, which taken does not contain, and its label, the name
// of the type in kebab case. It adds the name to taken.
func (k *caller) name(s *Start, typ string, taken map[string]bool) {
	s.Name = naming.Unique(taken, k.e.JoinName(valuesWord, k.e.JoinName(typ, k.decl)))
	taken[s.Name] = true
	s.Label = naming.Kebab(typ)
}

// defaultsOf returns the value of the struct st at the type ref with the
// declared default of each field that declares one, and false for a struct
// without a declared default. The rules of the language read the text of
// each default in the scope of its field, and a text that they do not read
// as a value of the type of its field leaves the field at its zero value.
func (k *caller) defaultsOf(st *node.Struct, ref *node.TypeRef) (emit.Value, bool) {
	var fields []emit.ValueField
	for _, f := range st.Fields {
		text, declared := sdk.FactOf(k.m, f.ID, meta.Named[string](defaults.KeyValue))
		if !declared {
			continue
		}
		scope := k.b.View().ScopeOf(f.ID)
		if value, read := defaults.Value(&k.b, &scope, f.Type, text); read {
			fields = append(fields, emit.NamedField(f.Name, value))
		}
	}
	if len(fields) == 0 {
		return emit.Value{}, false
	}
	return emit.Composite(rules.EmitRef(ref), fields...), true
}
