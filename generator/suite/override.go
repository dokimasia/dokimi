// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"go/token"
	"slices"
	"strings"

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

// Variants is the generator of the values of an interface with an unexported
// method, which no code outside its package can implement. Type is the
// interface. Members are the implementations that the package of the
// interface exports and the file of the suite can import, in the order of
// their identities. Each property run and each benchmark of the suite draws
// one of them for each value of the interface.
type Variants struct {
	Type    *emit.TypeRef
	Members []Variant
}

// Variant is one implementation of the interface of a [Variants]. Type is the
// implementation, and Label the label of its draw: the name of the
// implementation in kebab case. Address reports an implementation whose
// pointer alone has the methods of the interface, so the generator returns
// the address of the drawn value.
type Variant struct {
	Type    *emit.TypeRef
	Label   string
	Address bool
}

// Stream is the generator of the values of a channel that a subject receives
// from: a receive channel or a channel of both directions. Type is the
// channel type, and Elem the type of its elements. A value is a closed
// channel with a drawn list of elements in its buffer.
type Stream struct {
	Type *emit.TypeRef
	Elem *emit.TypeRef
}

// Reader is the generator of the values of a type that reads drawn bytes: an
// interface of the package io whose methods *conform.Bytes has, and
// *bufio.Reader. Type is the type, and Label the label of the draw of the
// bytes. Wrap is the function that takes the *conform.Bytes, as
// bufio.NewReader, and the zero value for an interface of io.
type Reader struct {
	Type  *emit.TypeRef
	Wrap  emit.Value
	Label string
}

// Constant is the generator of one fixed value of a type, which each case
// builds again: nil for unsafe.Pointer, context.Background() for
// context.Context, and the value of a constructor of the standard library
// for a type whose zero value and nil fail at their first use. Type is the
// type, and Value the value.
type Constant struct {
	Type  *emit.TypeRef
	Value emit.Value
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

// root is a type of the drawn values of a call, which the overrides of the
// suite walk. callable is the name of the callable whose call draws the
// value, as Placer.Place. at is the path from the type to the parameter of
// the callable, as "a result of Timer.C" and then "its parameter clock".
type root struct {
	ref      *node.TypeRef
	callable string
	at       []string
}

// to returns the root of the type ref that a value of the type of r
// contains. hop is the step from the type of r to ref in the path of a
// finding.
func (r root) to(ref *node.TypeRef, hop string) root {
	return root{ref: ref, callable: r.callable, at: append([]string{hop}, r.at...)}
}

// overrides records into the suite s the generator of each type that the
// roots of the calls take at some depth and that the shape reader of the
// property engine does not read, in the order in which a walk of the roots
// meets them. [caller.sourceOf] returns the generator of each type. The suite
// records the enumerations, the structs with a starting value, the generated
// functions and implementations, the variants of the interfaces with an
// unexported method, the streams, the readers of drawn bytes, the constants
// and the sinks. It records whether the roots take a value of the dynamic
// type and an error, and the implementations of the interfaces that the
// calls and the generators take. It reports [Ungenerated] for each type
// without a generator, once, at the first path that the walk meets.
//
// The walk enters the elements and the type arguments of each type that the
// shape reader reads, the exported fields and the embedded types of a struct
// and the target of an alias. It enters the elements of a stream, the
// implementations of the variants of an interface and the results of each
// generated function.
func (k *caller) overrides(s *Suite) {
	w := walk{k: k, met: map[string]bool{}, names: map[string]bool{}}
	for i := 0; i < len(k.roots); i++ {
		w.visit(k.roots[i])
	}
	s.Using, s.Starts, s.Generated, s.Variants = w.enums, w.starts, w.generated, w.variants
	s.Streams, s.Readers, s.Constants, s.Sinks = w.streams, w.readers, w.constants, w.sinks
	s.Dynamic, s.Errors, s.Inputs = w.dynamic, w.errors, k.inputs
}

// walk is a walk of the types of the roots of the calls of a suite. met
// contains the key of each type that the walk visited, and names the names
// of the generators of the starting values. The other fields are the
// overrides that the walk records, as [Suite] states them.
type walk struct {
	k         *caller
	met       map[string]bool
	names     map[string]bool
	enums     []Enumeration
	starts    []Start
	generated []Generated
	variants  []Variants
	streams   []Stream
	readers   []Reader
	constants []Constant
	sinks     []*emit.TypeRef
	dynamic   bool
	errors    bool
}

// visit records the generator of the type of r, and visits each type that a
// value of the type contains, once for each type.
func (w *walk) visit(r root) {
	k := w.k
	ref := r.ref
	key := k.typeKey(ref)
	if w.met[key] {
		return
	}
	w.met[key] = true
	src, reason := k.sourceOf(ref)
	switch src {
	case sourceDynamic:
		w.dynamic = true
	case sourceError:
		w.errors = true
	case sourceNone:
		k.m.Warnf(Ungenerated, "%s takes %s in %s. %s. Each run of %s fails before its first case until the test "+
			"process registers a generator of %s", r.callable, ref.Spelling, strings.Join(r.at, " in "), reason,
			r.callable, ref.Spelling)
	case sourceGenerated:
		w.generate(r)
	case sourceVariants:
		w.vary(r)
	case sourceStream:
		elem := k.underlying(ref).Elems[0]
		w.streams = append(w.streams, Stream{Type: rules.EmitRef(ref), Elem: rules.EmitRef(elem)})
		w.visit(r.to(elem, "an element of "+ref.Spelling))
	case sourceSink:
		w.sinks = append(w.sinks, rules.EmitRef(ref))
	case sourceBytes:
		w.readers = append(w.readers, k.readerOf(ref))
	case sourceConstant:
		w.constants = append(w.constants, Constant{Type: rules.EmitRef(ref), Value: k.constantOf(ref)})
	case sourceEngine:
		if !k.shaped(ref) {
			w.enter(r, src)
		}
	default:
		w.enter(r, src)
	}
}

// enter visits the elements and the type arguments of the type of r, records
// its enumeration or its starting value, and visits the parts of its
// declaration. src is the generator of the type: the shape reader, the
// members of an enumeration or a starting value.
func (w *walk) enter(r root, src source) {
	k := w.k
	ref := r.ref
	for i, elem := range ref.Elems {
		w.visit(r.to(elem, hopOf(ref, i)))
	}
	for _, arg := range ref.Args {
		w.visit(r.to(arg, "a type argument of "+ref.Spelling))
	}
	decl := k.declOf(ref)
	switch src {
	case sourceMembers:
		enum, _ := decl.(*node.Enum)
		w.enums = append(w.enums, Enumeration{Members: k.membersOf(enum, ref)})
	case sourceStart:
		w.starts = append(w.starts, k.starting(ref, w.names))
	default:
		// The shape reader reads the type, and the walk enters its parts.
	}
	for _, p := range partsOf(decl, ref) {
		w.visit(r.to(p.ref, p.hop))
	}
}

// hopOf returns the hop from the structural type ref to its child at the
// position i: a key or a value of a map, the target of a pointer, a field of
// an inline struct, and an element of a list or an array.
func hopOf(ref *node.TypeRef, i int) string {
	switch {
	case ref.Form == symbol.FormMap && i == 0:
		return "a key of " + ref.Spelling
	case ref.Form == symbol.FormMap:
		return "a value of " + ref.Spelling
	case ref.Form == symbol.FormOptional:
		return "the target of " + ref.Spelling
	case ref.Form == symbol.FormInline:
		return "a field of " + ref.Spelling
	}
	return "an element of " + ref.Spelling
}

// typeKey returns the key of the type ref in the walk: the identity of its
// target with the spellings of its type arguments, and the spelling of a
// type without a target.
func (*caller) typeKey(ref *node.TypeRef) string {
	t := rules.EmitRef(ref)
	if ref.Target.IsZero() {
		return keyOf(t)
	}
	return ref.Target.String() + "[" + spelledArgs(t) + "]"
}

// inner is a type that a value of a declaration contains, with the hop from
// the declaration to the type in the path of a finding.
type inner struct {
	ref *node.TypeRef
	hop string
}

// partsOf returns the types that a value of the declaration decl at the type
// ref contains, at the type arguments of ref: the exported fields and the
// embedded types of a struct, and the target of an alias.
func partsOf(decl symbol.Symbol, ref *node.TypeRef) []inner {
	switch d := decl.(type) {
	case *node.Struct:
		out := make([]inner, 0, len(d.Fields)+len(d.Embeds))
		for _, f := range d.Fields {
			if f.Visibility == symbol.VisibilityPublic {
				out = append(out, inner{
					ref: rules.Substitute(f.Type, d.TypeParams, ref.Args),
					hop: "the field " + f.Name + " of " + ref.Spelling,
				})
			}
		}
		for _, embed := range d.Embeds {
			out = append(out, inner{
				ref: rules.Substitute(embed.Ref, d.TypeParams, ref.Args),
				hop: "the embedded field " + embeddedName(embed.Ref) + " of " + ref.Spelling,
			})
		}
		return out
	case *node.Alias:
		return []inner{{ref: rules.Substitute(d.Target, d.TypeParams, ref.Args), hop: "the target of " + ref.Spelling}}
	}
	return nil
}

// generate records the generator of the values of the function type or the
// interface of r. A value of a function type is a generated function, and a
// value of an interface is the address of the implementation of the suite,
// with a generated function of each method. The results of these functions
// are roots of the walk.
//
// The name of the generated functions is the name of the type in camel case,
// and their label the name in kebab case. A function type without a name
// takes the name and the label func, and an interface without a name the
// name and the label interface.
func (w *walk) generate(r root) {
	k := w.k
	ref := r.ref
	g := Generated{Type: rules.EmitRef(ref)}
	iface, under, isInterface := k.interfaceOf(ref)
	word := ref.Target.Name
	switch {
	case word != "":
	case isInterface:
		word = interfaceWord
	default:
		word = funcWord
	}
	name, label, taken := naming.Camel(word), naming.Kebab(word), k.taken(ref)
	if isInterface {
		g.Input, g.Funcs = k.inputOf(under, iface, inputWordOf(ref, iface), name, label, taken, r)
	} else {
		g.Funcs = []Func{k.funcOf(k.underlying(ref), name, label, taken, r, ref.Spelling)}
	}
	if k.recursive(ref, keyOf(g.Type)) {
		g.Leaves = recursiveLeaves
	}
	w.generated = append(w.generated, g)
}

// vary records the variants of the interface of r, and visits each of its
// implementations.
func (w *walk) vary(r root) {
	k := w.k
	iface, under, _ := k.interfaceOf(r.ref)
	v := Variants{Type: rules.EmitRef(r.ref), Members: k.variantsOf(under, iface)}
	w.variants = append(w.variants, v)
	for _, m := range v.Members {
		w.visit(r.to(&node.TypeRef{Spelling: m.Type.Spelling, Target: m.Type.Target},
			"the implementation "+m.Type.Spelling+" of "+r.ref.Spelling))
	}
}

// readerOf returns the generator of the values of the type ref, whose values
// read drawn bytes. A *bufio.Reader wraps the bytes in its constructor. The
// label of the draw of the bytes is the name of the type in kebab case.
func (k *caller) readerOf(ref *node.TypeRef) Reader {
	under := k.underlying(ref)
	r := Reader{Type: rules.EmitRef(ref), Label: naming.Kebab(under.Target.Name)}
	if under.Form == symbol.FormOptional {
		r.Label, r.Wrap = naming.Kebab(under.Elems[0].Target.Name), bufferedReader
	}
	return r
}

// constantOf returns the fixed value of the type ref: nil for unsafe.Pointer,
// context.Background() for context.Context, and the value of the constructor
// of a pointer of the standard library.
func (k *caller) constantOf(ref *node.TypeRef) emit.Value {
	under := k.underlying(ref)
	switch {
	case under.Form == symbol.FormOptional:
		return constructed[qualified(&under.Elems[0].Target)]
	case is(&under.Target, contextPackage, contextType):
		return emit.Call(stdIdentity(contextPackage, "Background", symbol.KindFunction))
	}
	return emit.Literal(emit.LiteralNil, "")
}

// keyOf returns the key of the function type or the interface t in the
// search of a recursion: its spelling with the spellings of its type
// arguments.
func keyOf(t *emit.TypeRef) string {
	return t.Spelling + "[" + spelledArgs(t) + "]"
}

// recursive reports whether the results of the generated functions of a value
// of the type ref contain the type with the key target, at some depth. The
// search enters the elements of each type, the exported fields and the
// embedded types of a struct and the target of an alias at the type arguments
// of the type, and the results of each other generated type that it meets.
func (k *caller) recursive(ref *node.TypeRef, target string) bool {
	s := search{k: k, target: target, met: map[string]bool{}, seen: map[string]bool{}}
	return slices.ContainsFunc(k.resultsOf(ref), s.contains)
}

// search is a search for the generated type with the key target. met
// contains the key of each generated type that the search entered, and seen
// the key of each declaration that it entered.
type search struct {
	k      *caller
	met    map[string]bool
	seen   map[string]bool
	target string
}

// contains reports whether a value of the type ref contains the type of the
// search.
func (s *search) contains(ref *node.TypeRef) bool {
	k := s.k
	if src, _ := k.sourceOf(ref); src == sourceGenerated {
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
	key := k.typeKey(ref)
	if s.seen[key] {
		return false
	}
	s.seen[key] = true
	return slices.ContainsFunc(partsOf(k.declOf(ref), ref), func(p inner) bool { return s.contains(p.ref) })
}

// resultsOf returns the types of the results of the generated functions of a
// value of the function type or the interface at the type ref.
func (k *caller) resultsOf(ref *node.TypeRef) []*node.TypeRef {
	iface, under, isInterface := k.interfaceOf(ref)
	if !isInterface {
		fn := k.underlying(ref)
		return fn.Elems[fn.Split:]
	}
	b := k.b
	var out []*node.TypeRef
	for _, method := range methodsOf(&b, iface) {
		view, _ := b.CallableOf(method)
		for _, r := range view.Returns {
			out = append(out, rules.Substitute(r.Ref, iface.TypeParams, under.Args))
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
// without a field that the shape reader of the property engine reads starts
// from its zero value, because code outside its package builds no other value
// without a constructor. A struct of the standard library and a generic
// struct at the type arguments of ref start so too, and so does a struct that
// ref refers to through an alias or a defined type. Such a struct with a
// readable field takes the shape of the property engine alone, and so does a
// struct that the file of the suite cannot import. taken contains the names
// of the generators of the suite, and startOf adds the name that it takes.
func (k *caller) startOf(st *node.Struct, ref *node.TypeRef, taken map[string]bool) (Start, bool) {
	if !k.importable(ref) {
		return Start{}, false
	}
	s := Start{Type: rules.EmitRef(ref)}
	if len(st.TypeParams) > 0 || ref.Target != st.ID {
		if !k.unreadable(st) {
			return Start{}, false
		}
		s.Value = emit.Composite(rules.EmitRef(ref))
	} else if sample, _ := k.b.View().Authored(k.b.Source(), st.ID, ref); sample.OK() && k.referable(&sample.Value) {
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

// starting returns the generator of the values of the type ref, whose source
// is a starting value: the whole value of a defined type, the starting value
// of a struct, or the zero value of an inline struct without a readable
// field. taken contains the names of the generators of the suite, and
// starting adds the name that it takes.
func (k *caller) starting(ref *node.TypeRef, taken map[string]bool) Start {
	if d, aliased := k.declOf(ref).(*node.Alias); aliased {
		if s, started := k.wholeStartOf(d, ref, taken); started {
			return s
		}
	}
	if st, isStruct := k.declOf(k.underlying(ref)).(*node.Struct); isStruct {
		s, _ := k.startOf(st, ref, taken)
		return s
	}
	s := Start{Type: rules.EmitRef(ref), Value: emit.Composite(rules.EmitRef(ref))}
	k.name(&s, structWord, taken)
	return s
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
