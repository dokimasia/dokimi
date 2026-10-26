// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"go/token"
	"slices"
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
// value of any other struct: its authored sample, its value with the declared
// default of each field, or its zero value where the shape reader of the
// property engine reads none of its fields. Fields are the fields that a case
// keeps or changes.
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

// Generated is the generator of the values of a function type, or of an
// interface outside the standard library, that the calls of a suite take at
// some depth of their inputs. Each property run and each benchmark of the
// suite takes it for every value of the type.
//
// Type is the type. Funcs are the generated functions of a value: the one
// function of a function type, or the function of each method of the
// interface. Input is the name of the implementation of the interface that
// the suite declares, and empty for a function type. The case of a value
// draws the results of its functions, as the case of a call draws them for a
// parameter of the type.
//
// Leaves is the most base values of one value of a recursive type, and 0 for
// a type that is not recursive. A type is recursive where the results of its
// generated functions contain the type, directly or through another type at
// some depth, such as an interface with the method Children() []Node. The
// suite generates such a type through prop.Recursive. Its base is the value
// whose functions return their zero results at once, with the end error of
// each error result.
type Generated struct {
	Type   *emit.TypeRef
	Input  string
	Funcs  []Func
	Leaves int
}

// recursiveLeaves is the most base values of one value of a recursive type.
// A value of a tree of nodes with 10 bases costs 119 µs and 494 allocations to
// draw, a tenth of the cost under the default of 100 bases of the property
// engine. 8% of such values have a depth of 3 nodes or more.
const recursiveLeaves = 10

// valuesWord is the word that the name of the generator of a struct with a
// starting value joins to the struct and the declaration of the suite.
const valuesWord = "Values"

// The key and the value of the tag of prop that leaves a field out of the
// shape of its struct.
const (
	propTag  = "prop"
	propSkip = "-"
)

// overrides records into the suite s each enumeration, each struct with a
// starting value, and each function type and interface outside the standard
// library that the roots of the calls take at some depth, in the order in
// which a walk of the roots meets them. It records whether the roots take a
// value of the dynamic type, Go's any, at some depth, and the implementations
// of the interfaces that the calls and the generators take.
//
// The walk enters the elements and the type arguments of each type, the
// exported fields and the embedded types of a struct, the target of an alias,
// and the results of each generated function. It leaves the types of the
// standard library to the shape reader of the property engine, because the
// value sets of the standard library include sets of flags, such as
// fs.FileMode, and units, such as time.Duration. It does not override a type
// that the file of the suite cannot import, so a registration of the test
// process states the generator of such a type. It also leaves out each value
// that refers to a declaration that the file cannot refer to, such as a
// member of an enumeration that its package does not export.
func (k *caller) overrides(s *Suite) {
	w := walk{k: k, seen: map[symbol.Identity]bool{}, met: map[string]bool{}, names: map[string]bool{}}
	for i := 0; i < len(k.roots); i++ {
		w.visit(k.roots[i])
	}
	s.Using, s.Starts, s.Generated, s.Dynamic, s.Inputs = w.enums, w.starts, w.generated, w.dynamic, k.inputs
}

// walk is a walk of the types of the roots of the calls of a suite. seen
// contains each declaration that the walk visited, met the spelling of each
// function type and each interface outside the standard library that it met,
// and names the names of the generators of the starting values. enums,
// starts and generated are the enumerations, the structs with a starting
// value, and the function types and the interfaces that the walk overrides.
// dynamic reports that it met the dynamic type.
type walk struct {
	k         *caller
	seen      map[symbol.Identity]bool
	met       map[string]bool
	names     map[string]bool
	enums     []Enumeration
	starts    []Start
	generated []Generated
	dynamic   bool
}

// visit visits the type ref and each type that it contains. A structural form
// has no target, and [golang.Standard] reports the empty package of such a
// target as a package of the standard library, so the walk declares nothing
// for it. The dynamic type contains no type.
func (w *walk) visit(ref *node.TypeRef) {
	if w.k.b.TypeOf(ref).Form == symbol.FormDynamic {
		w.dynamic = true
		return
	}
	if w.generate(ref) {
		return
	}
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
		if !w.k.importable(ref) {
			return
		}
		members := make([]emit.Value, 0, len(d.Variants))
		for _, v := range d.Variants {
			if member := emit.Name(v.ID); w.k.referable(&member) {
				members = append(members, member)
			}
		}
		if len(members) > 0 {
			w.enums = append(w.enums, Enumeration{Members: members})
		}
	case *node.Struct:
		if start, started := w.k.startOf(d, ref, w.names); started {
			w.starts = append(w.starts, start)
		}
	case *node.Alias:
		if start, started := w.k.wholeStartOf(d, ref, w.names); started {
			w.starts = append(w.starts, start)
		}
	}
	for _, part := range partsOf(decl, ref) {
		w.visit(part)
	}
}

// partsOf returns the types that a value of the declaration decl at the type
// ref contains, at the type arguments of ref: the exported fields and the
// embedded types of a struct, and the target of an alias.
func partsOf(decl symbol.Symbol, ref *node.TypeRef) []*node.TypeRef {
	switch d := decl.(type) {
	case *node.Struct:
		out := make([]*node.TypeRef, 0, len(d.Fields)+len(d.Embeds))
		for _, f := range d.Fields {
			if f.Visibility == symbol.VisibilityPublic {
				out = append(out, rules.Substitute(f.Type, d.TypeParams, ref.Args))
			}
		}
		for _, embed := range d.Embeds {
			out = append(out, rules.Substitute(embed.Ref, d.TypeParams, ref.Args))
		}
		return out
	case *node.Alias:
		return []*node.TypeRef{rules.Substitute(d.Target, d.TypeParams, ref.Args)}
	}
	return nil
}

// generate reports whether ref is a function type or an interface outside the
// standard library, and records the generator of the values of such a type
// the first time that the walk meets it. The walk does not enter such a type.
// A value of a function type is a generated function, and a value of an
// interface is the implementation of the suite with a generated function of
// each method. The results of these functions are roots of the walk. The
// suite does not record a generator of a type whose generated code refers to
// a type that the file of the suite cannot import.
//
// The name of the generated functions is the name of the type in camel case,
// and their label the name in kebab case. A function type without a name
// takes the name and the label func.
func (w *walk) generate(ref *node.TypeRef) bool {
	k := w.k
	under := k.underlying(ref)
	iface, implemented := k.implemented(ref)
	if under.Form != symbol.FormFunc && !implemented {
		return false
	}
	g := Generated{Type: rules.EmitRef(ref)}
	key := keyOf(g.Type)
	if w.met[key] {
		return true
	}
	w.met[key] = true
	word := ref.Target.Name
	if word == "" {
		word = funcWord
	}
	name, label, taken := naming.Camel(word), naming.Kebab(word), k.taken(ref)
	switch {
	case implemented && k.implementable(ref, iface):
		g.Input, g.Funcs = k.inputOf(ref, iface, name, label, taken)
	case !implemented && k.importable(ref, under):
		g.Funcs = []Func{k.funcOf(under, name, label, taken)}
	default:
		return true
	}
	if k.recursive(ref, key) {
		g.Leaves = recursiveLeaves
	}
	w.generated = append(w.generated, g)
	return true
}

// keyOf returns the key of the function type or the interface t in the walk:
// its spelling with the spellings of its type arguments.
func keyOf(t *emit.TypeRef) string {
	return t.Spelling + "[" + spelledArgs(t) + "]"
}

// recursive reports whether the results of the generated functions of a value
// of the type ref contain the type with the key target, at some depth. The
// search enters the elements of each type, the exported fields and the
// embedded types of a struct and the target of an alias at the type arguments
// of the type, and the results of each other function type and interface
// outside the standard library that it meets.
func (k *caller) recursive(ref *node.TypeRef, target string) bool {
	s := search{k: k, target: target, met: map[string]bool{}, seen: map[symbol.Identity]bool{}}
	return slices.ContainsFunc(k.resultsOf(ref), s.contains)
}

// search is a search for the function type or the interface with the key
// target. met contains the key of each function type and interface that the
// search entered, and seen each declaration that it entered.
type search struct {
	k      *caller
	met    map[string]bool
	seen   map[symbol.Identity]bool
	target string
}

// contains reports whether a value of the type ref contains the type of the
// search.
func (s *search) contains(ref *node.TypeRef) bool {
	k := s.k
	under := k.underlying(ref)
	if _, implemented := k.implemented(ref); under.Form == symbol.FormFunc || implemented {
		key := keyOf(rules.EmitRef(ref))
		if key == s.target {
			return true
		}
		if s.met[key] {
			return false
		}
		s.met[key] = true
		return slices.ContainsFunc(k.resultsOf(ref), s.contains)
	}
	if slices.ContainsFunc(ref.Elems, s.contains) {
		return true
	}
	//dokimi:mutate-skip lcr-left: a declaration of the standard library has no part that the workspace declares, so its parts never contain the type of the search
	if s.seen[ref.Target] || golang.Standard(ref.Target.Package) {
		return false
	}
	s.seen[ref.Target] = true
	decl, _ := k.m.Reader().Lookup(ref.Target)
	return slices.ContainsFunc(partsOf(decl, ref), s.contains)
}

// resultsOf returns the types of the results of the generated functions of a
// value of the function type or the interface at the type ref.
func (k *caller) resultsOf(ref *node.TypeRef) []*node.TypeRef {
	iface, implemented := k.implemented(ref)
	if !implemented {
		under := k.underlying(ref)
		return under.Elems[under.Split:]
	}
	b := k.b
	var out []*node.TypeRef
	for _, method := range methodsOf(&b, iface) {
		view, _ := b.CallableOf(method)
		for _, r := range view.Returns {
			out = append(out, rules.Substitute(r.Ref, iface.TypeParams, ref.Args))
		}
	}
	return out
}

// startOf returns the generator of the values of the struct st at the type
// ref, and false for a struct without a starting value. The starting value is
// the authored sample of the struct, the value of its builder, or its value
// with the declared defaults of its fields, in that order. A sample that the
// file of the suite cannot refer to yields to the next source. The value of a
// builder applies to a struct of the
// package of the declaration alone, because the template finds the builder by
// the name of its constructor, which is unique in one package. A struct
// without a field that the shape reader of the
// property engine reads starts from its zero value, because code outside its
// package builds no other value without a constructor. A generic struct takes
// the shape of the property engine alone, and so does a struct that the file
// of the suite cannot import. taken contains the names of the generators of
// the suite, and startOf adds the name that it takes.
func (k *caller) startOf(st *node.Struct, ref *node.TypeRef, taken map[string]bool) (Start, bool) {
	if len(st.TypeParams) > 0 || !k.importable(ref) {
		return Start{}, false
	}
	s := Start{Type: rules.EmitRef(ref)}
	if sample, _ := k.b.View().Authored(k.b.Source(), st.ID, ref); sample.OK() && k.referable(&sample.Value) {
		s.Value = sample.Value
	} else if seed, built := k.seedOf(st); built {
		s.Seed = &emit.TypeRef{Spelling: seed, Target: st.ID, Emitted: true}
	} else if value, defaulted := k.defaultsOf(st, ref); defaulted {
		s.Value = value
	} else if k.unreadable(st) {
		s.Value = emit.Composite(rules.EmitRef(ref))
	} else {
		return Start{}, false
	}
	k.name(&s, st.Name, taken)
	for _, f := range st.Fields {
		if k.readable(f) {
			s.Fields = append(s.Fields, StartField{Name: f.Name, Label: s.Label + "." + naming.Kebab(f.Name)})
		}
	}
	for _, embed := range st.Embeds {
		if name := embeddedName(embed.Ref); token.IsExported(name) {
			s.Fields = append(s.Fields, StartField{Name: name, Label: s.Label + "." + naming.Kebab(name)})
		}
	}
	return s, true
}

// seedOf returns the name of the constructor of the builder of the struct st
// from its seed, and false where st has no builder or is a struct of another
// package than the declaration.
func (k *caller) seedOf(st *node.Struct) (string, bool) {
	seed, built := sdk.FactOf(k.m, st.ID, meta.Named[string](builder.KeySeed))
	return seed, built && st.ID.Package == k.origin.Package
}

// readable reports whether the shape reader of the property engine reads the
// field f: a public field without the tag prop:"-".
func (k *caller) readable(f *node.Field) bool {
	tag, _ := k.b.Tag(f, propTag)
	return f.Visibility == symbol.VisibilityPublic && tag != propSkip
}

// unreadable reports whether the shape reader of the property engine reads no
// field of the struct st and no type that st embeds.
func (k *caller) unreadable(st *node.Struct) bool {
	exported := func(e *node.Embed) bool { return token.IsExported(embeddedName(e.Ref)) }
	return !slices.ContainsFunc(st.Fields, k.readable) && !slices.ContainsFunc(st.Embeds, exported)
}

// embeddedName returns the name of the field of the embedded type ref: the
// name of the type without its qualifier and without the star of a pointer.
func embeddedName(ref *node.TypeRef) string {
	return strings.TrimPrefix(ref.Spelling[strings.LastIndex(ref.Spelling, ".")+1:], "*")
}

// wholeStartOf returns the generator of the values of the defined type d at
// the type ref, whose value a case keeps or changes whole, and false for a
// type without a starting value. The starting value is the authored sample of
// the type, or its declared default where the file of the suite cannot refer
// to the sample. A default that the file cannot refer to leaves the type
// without a starting value. A transparent alias is the type of its target.
// Generic types take the shape of the property engine alone, and so does a
// type that the file of the suite cannot import. taken contains the names of
// the generators of the suite, and wholeStartOf adds the name that it takes.
func (k *caller) wholeStartOf(d *node.Alias, ref *node.TypeRef, taken map[string]bool) (Start, bool) {
	if !d.Defined || len(d.TypeParams) > 0 || !k.importable(ref) {
		return Start{}, false
	}
	s := Start{Type: rules.EmitRef(ref), Whole: true}
	if sample, _ := k.b.View().Authored(k.b.Source(), d.ID, ref); sample.OK() && k.referable(&sample.Value) {
		s.Value = sample.Value
	} else if text, declared := sdk.FactOf(k.m, d.ID, meta.Named[string](defaults.KeyValue)); declared {
		scope := k.b.View().ScopeOf(d.ID)
		value, read := defaults.Value(&k.b, &scope, ref, text)
		if !read || !k.referable(&value) {
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
// as a value of the type of its field leaves the field at its zero value. So
// does a default of a field that the file of the suite cannot set, or a
// default that refers to a declaration that the file cannot refer to.
func (k *caller) defaultsOf(st *node.Struct, ref *node.TypeRef) (emit.Value, bool) {
	var fields []emit.ValueField
	for _, f := range st.Fields {
		text, declared := sdk.FactOf(k.m, f.ID, meta.Named[string](defaults.KeyValue))
		field := emit.Name(f.ID)
		if !declared || !k.referable(&field) {
			continue
		}
		scope := k.b.View().ScopeOf(f.ID)
		if value, read := defaults.Value(&k.b, &scope, f.Type, text); read && k.referable(&value) {
			fields = append(fields, emit.NamedField(f.Name, value))
		}
	}
	if len(fields) == 0 {
		return emit.Value{}, false
	}
	return emit.Composite(rules.EmitRef(ref), fields...), true
}
