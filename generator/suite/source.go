// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"slices"
	"strings"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// source is the generator of the values of one type in the inputs of a
// suite.
type source uint8

const (
	// sourceEngine is the shape reader of the property engine.
	sourceEngine source = iota
	// sourceMembers draws one member of an enumeration.
	sourceMembers
	// sourceStart starts from a sample, the value of a builder, declared
	// defaults or the zero value, and keeps or changes each readable field.
	sourceStart
	// sourceGenerated is a generated function or a generated implementation.
	sourceGenerated
	// sourceVariants draws one implementation of an interface with an
	// unexported method.
	sourceVariants
	// sourceStream is conform.Stream.
	sourceStream
	// sourceSink is the sink of the test.
	sourceSink
	// sourceBytes is a *conform.Bytes over drawn bytes, alone or as the
	// reader of a constructor of the standard library.
	sourceBytes
	// sourceConstant is one fixed value: nil, context.Background() or the
	// value of a constructor of the standard library.
	sourceConstant
	// sourceError draws nil or conform.ErrGenerated.
	sourceError
	// sourceDynamic draws a string, an int or a bool.
	sourceDynamic
	// sourceNone is no generator, and the suite reports Ungenerated.
	sourceNone
)

// The packages and the names of the declarations of the standard library
// that the suite generates by their identity.
const (
	contextPackage = "context"
	contextType    = "Context"
	unsafePackage  = "unsafe"
	pointerType    = "Pointer"
	ioPackage      = "io"
	bufioPackage   = "bufio"
	readerType     = "Reader"
)

// bytesInterfaces are the interfaces of the package io whose methods
// *conform.Bytes has. A value of each is a *conform.Bytes over drawn bytes,
// because a generated Read breaks the contract of io.Reader.
var bytesInterfaces = map[string]bool{
	"Reader": true, "ReaderAt": true, "ByteReader": true, "ByteScanner": true, "RuneReader": true,
	"RuneScanner": true, "Seeker": true, "ReadSeeker": true, "ReadCloser": true, "ReadSeekCloser": true,
	"WriterTo": true, "Writer": true, "StringWriter": true, "ByteWriter": true, "WriteCloser": true,
	"WriteSeeker": true, "ReadWriter": true, "ReadWriteCloser": true, "ReadWriteSeeker": true, "Closer": true,
}

// constructed are the values of the pointers of the standard library whose
// zero value and nil both fail at their first use, by the package and the
// name of the type that a pointer points to. Each constructor builds a value
// that uses no file, no socket and no process.
var constructed = map[string]emit.Value{
	"log/slog.Logger": emit.Call(stdIdentity("log/slog", "New", symbol.KindFunction),
		emit.Name(stdIdentity("log/slog", "DiscardHandler", symbol.KindVariable))),
	"log.Logger": emit.Call(stdIdentity("log", "New", symbol.KindFunction),
		emit.Name(stdIdentity(ioPackage, "Discard", symbol.KindVariable)), emit.Literal(emit.LiteralString, ""),
		emit.Literal(emit.LiteralInt, "0")),
	"bufio.Writer": emit.Call(stdIdentity(bufioPackage, "NewWriter", symbol.KindFunction),
		emit.Name(stdIdentity(ioPackage, "Discard", symbol.KindVariable))),
}

// shapes are the types that the shape reader of the property engine reads as
// a shape of their own, by their package and name: a time, an address, a
// wall time and the pointers of pointerShapes. The walk of the overrides does
// not enter them, so no override replaces a part of such a shape.
var shapes = map[string]bool{
	"time.Time": true, "net/netip.Addr": true, "go.dokimi.dev/assert/prop.WallTime": true,
}

// shaped reports whether the shape reader of the property engine reads the
// type ref as a shape of its own. A defined type over such a type is a type
// of its own, which the reader reads by its kind.
func (k *caller) shaped(ref *node.TypeRef) bool {
	ref = k.aliased(ref)
	if ref.Form == symbol.FormOptional {
		return pointerShapes[qualified(&ref.Elems[0].Target)]
	}
	return shapes[qualified(&ref.Target)]
}

// aliased returns the type that ref refers to through its transparent
// aliases, at the type arguments of ref.
func (k *caller) aliased(ref *node.TypeRef) *node.TypeRef {
	for range aliasDepth {
		alias, is := k.declOf(ref).(*node.Alias)
		if !is || alias.Defined {
			return ref
		}
		ref = rules.Substitute(alias.Target, alias.TypeParams, ref.Args)
	}
	return ref
}

// bufferedReader is the constructor of a *bufio.Reader, which takes the
// *conform.Bytes of the drawn bytes.
var bufferedReader = emit.Name(stdIdentity(bufioPackage, "NewReader", symbol.KindFunction))

// stdIdentity returns the identity of the declaration name in the package
// pkg of the standard library. k is the kind of the declaration.
func stdIdentity(pkg, name string, k symbol.Kind) symbol.Identity {
	return symbol.Identity{Lang: golang.Lang, Package: pkg, Name: name, Kind: k}
}

// interfaceWord and structWord are the words of the names of the generators
// of an interface and a struct without a name, such as the type of the field
// Named interface{ Name() string }.
const (
	interfaceWord = "Interface"
	structWord    = "Struct"
)

// jsonTag is the key of the tag that names a field in the shape of its
// struct.
const jsonTag = "json"

// sourceOf returns the generator of the values of the type ref at its type
// arguments. For sourceNone it also returns the reason why no generator of
// the suite covers the type, as a sentence. A list, an array and a map take
// the shape reader, and the walk of the overrides enters their parts.
func (k *caller) sourceOf(ref *node.TypeRef) (source, string) {
	switch shape := k.b.TypeOf(ref); {
	case shape.Form == symbol.FormDynamic:
		return sourceDynamic, ""
	case shape.Ref == rules.WellKnownError:
		return sourceError, ""
	case k.shaped(ref):
		return sourceEngine, ""
	}
	if d, aliased := k.declOf(ref).(*node.Alias); aliased && d.Defined {
		if _, started := k.wholeStartOf(d, ref, map[string]bool{}); started {
			return sourceStart, ""
		}
	}
	under := k.underlying(ref)
	switch {
	case under.Form == symbol.FormFunc && k.importable(ref) && !k.importable(under):
		return sourceNone, fmt.Sprintf("The signature of %s refers to a declaration that the file of the suite "+
			"cannot import", ref.Spelling)
	case under.Form == symbol.FormFunc:
		return k.spelled(ref, sourceGenerated)
	case under.Form == symbol.FormStream:
		return k.channel(ref, under)
	case under.Form == symbol.FormInline:
		return k.inline(ref, under)
	case under.Form == symbol.FormOptional:
		return k.pointer(ref, under)
	case k.b.TypeOf(under).Form == symbol.FormDynamic:
		return k.spelled(ref, sourceGenerated)
	case under.Form == symbol.FormNamed:
		return k.named(ref, under)
	}
	return sourceEngine, ""
}

// channel returns the generator of the values of the channel type ref, whose
// underlying type is the channel under. A send channel takes the sink of the
// test, and any other channel a stream of drawn elements.
func (k *caller) channel(ref, under *node.TypeRef) (source, string) {
	if sends(under.Spelling) {
		return k.spelled(ref, sourceSink)
	}
	if !k.importable(ref, under.Elems[0]) {
		return sourceNone, fmt.Sprintf("The channel %s refers to a declaration that the file of the suite cannot "+
			"import", ref.Spelling)
	}
	return sourceStream, ""
}

// inline returns the generator of the values of the type ref, whose
// underlying type under is an inline struct or interface. An interface takes
// a generated implementation. A struct takes the shape reader where the
// reader reads one of its fields, and its zero value otherwise.
func (k *caller) inline(ref, under *node.TypeRef) (source, string) {
	b, _ := parsedBody(under)
	switch {
	case b.isInterface:
		it := &node.Interface{Methods: under.Methods, Embeds: b.embeds}
		if !k.implementable(ref, it) {
			return sourceNone, fmt.Sprintf("A method of %s refers to a declaration that the file of the suite "+
				"cannot import", ref.Spelling)
		}
		return sourceGenerated, ""
	case b.exported || slices.ContainsFunc(under.Fields, k.readable):
		return sourceEngine, ""
	}
	return k.spelled(ref, sourceStart)
}

// pointer returns the generator of the values of the pointer type ref, whose
// underlying type is the pointer under. A pointer of a constructor of the
// standard library takes the value of the constructor, and any other pointer
// the shape reader, which draws nil or the address of a drawn value.
func (k *caller) pointer(ref, under *node.TypeRef) (source, string) {
	target := &under.Elems[0].Target
	if _, built := constructed[qualified(target)]; built {
		return k.spelled(ref, sourceConstant)
	}
	if is(target, bufioPackage, readerType) {
		return k.spelled(ref, sourceBytes)
	}
	return sourceEngine, ""
}

// named returns the generator of the values of the type ref, whose
// underlying type under is a named type: a predeclared type, an enumeration,
// a struct, an interface or unsafe.Pointer. A type that the view does not
// contain takes the shape reader, or the registration of the test process.
func (k *caller) named(ref, under *node.TypeRef) (source, string) {
	if is(&under.Target, unsafePackage, pointerType) {
		return k.spelled(ref, sourceConstant)
	}
	switch d := k.declOf(under).(type) {
	case *node.Enum:
		if ref.Target == d.ID && !golang.Standard(d.ID.Package) && len(k.membersOf(d, ref)) > 0 {
			return sourceMembers, ""
		}
	case *node.Struct:
		return k.structure(ref, d)
	case *node.Interface:
		return k.iface(ref, under, d)
	}
	return sourceEngine, ""
}

// structure returns the generator of the values of the type ref, whose
// underlying type is the struct st. A struct with two fields of one JSON name
// has no generator, because the shape reader names a field by its JSON name.
func (k *caller) structure(ref *node.TypeRef, st *node.Struct) (source, string) {
	if first, second, name, twins := k.twins(st); twins {
		return sourceNone, fmt.Sprintf("The fields %s and %s of %s are both named %q",
			first, second, ref.Spelling, name)
	}
	if _, started := k.startOf(st, ref, map[string]bool{}); started {
		return sourceStart, ""
	}
	if k.unreadable(st) {
		return sourceNone, fmt.Sprintf("The struct %s has no field that the property engine reads, and the file of "+
			"the suite cannot import it", ref.Spelling)
	}
	return sourceEngine, ""
}

// unimplemented is the reason of the finding Ungenerated for an interface
// with an unexported method whose package does not export an implementation
// that the file of the suite can import.
const unimplemented = "The interface %s has an unexported method, and the package %s exports no " +
	"implementation that the file of the suite can import"

// unbuildable is the reason of the finding Ungenerated for an interface of the
// standard library with an unexported method.
const unbuildable = "The interface %s of the standard library has an unexported method, so a test cannot build " +
	"an implementation of it"

// iface returns the generator of the values of the type ref, whose
// underlying type under refers to the interface it. context.Context takes
// context.Background(), and an interface of io whose methods *conform.Bytes
// has takes a *conform.Bytes over drawn bytes. An interface of the standard
// library with an unexported method has no generator, because the values of
// its package depend on state that only the package sets, such as a
// *testing.T of a running test. Any other interface with an unexported method
// takes one of the implementations that its package exports, and any other
// interface a generated implementation.
func (k *caller) iface(ref, under *node.TypeRef, it *node.Interface) (source, string) {
	switch {
	case !k.importable(ref):
		return sourceNone, unimported(ref)
	case is(&it.ID, contextPackage, contextType):
		return k.spelled(ref, sourceConstant)
	case it.ID.Package == ioPackage && bytesInterfaces[it.ID.Name]:
		return k.spelled(ref, sourceBytes)
	case k.sealed(it) && golang.Standard(it.ID.Package):
		return sourceNone, fmt.Sprintf(unbuildable, ref.Spelling)
	case k.sealed(it):
		if len(k.variantsOf(under, it)) == 0 {
			return sourceNone, fmt.Sprintf(unimplemented, ref.Spelling, it.ID.Package)
		}
		return sourceVariants, ""
	case !k.implementable(ref, it):
		return sourceNone, fmt.Sprintf("A method of %s refers to a declaration that the file of the suite cannot "+
			"import", ref.Spelling)
	}
	return sourceGenerated, ""
}

// spelled returns the generator src of the type ref, whose override spells
// ref. A type that the file of the suite cannot import has no generator.
func (k *caller) spelled(ref *node.TypeRef, src source) (source, string) {
	if !k.importable(ref) {
		return sourceNone, unimported(ref)
	}
	return src, ""
}

// unimported returns the reason of the finding Ungenerated for the type ref,
// which the file of the suite cannot import.
func unimported(ref *node.TypeRef) string {
	return fmt.Sprintf("The type %s refers to a declaration that the file of the suite cannot import", ref.Spelling)
}

// declOf returns the declaration that the target of ref refers to, and nil
// for a reference without a target or with a target outside the view.
func (k *caller) declOf(ref *node.TypeRef) symbol.Symbol {
	decl, _ := k.m.Reader().Lookup(ref.Target)
	return decl
}

// is reports whether id is the declaration name of the package pkg.
func is(id *symbol.Identity, pkg, name string) bool {
	return id.Package == pkg && id.Name == name
}

// qualified returns the package and the name of the declaration id, joined
// by a dot, as log/slog.Logger.
func qualified(id *symbol.Identity) string {
	return id.Package + "." + id.Name
}

// sends reports whether spelling spells a send channel. Go associates the
// arrow with the leftmost chan, so the spelling opens with the direction of
// the outer channel, whatever white space it has.
func sends(spelling string) bool {
	return strings.HasPrefix(strings.Join(strings.Fields(spelling), ""), "chan<-")
}

// receives reports whether spelling spells a receive channel, as sends reads
// the spelling.
func receives(spelling string) bool {
	return strings.HasPrefix(strings.Join(strings.Fields(spelling), ""), "<-chan")
}

// sealed reports whether the interface it has a method that its package does
// not export, so no code outside its package implements it.
func (k *caller) sealed(it *node.Interface) bool {
	set, _ := k.b.MembersOf(it)
	return slices.ContainsFunc(set.Members, func(m rules.Member) bool {
		method, isMethod := m.Symbol.(*node.Method)
		return isMethod && method.Visibility != symbol.VisibilityPublic
	})
}

// variantsOf returns the implementations of the interface it at the type
// under that its package exports and the file of the suite can import, in
// the order of their identities. Each is a struct or an enumeration without
// type parameters, with a generator of its own. An implementation whose
// value lacks a method of the interface takes its address.
func (k *caller) variantsOf(under *node.TypeRef, it *node.Interface) []Variant {
	var out []Variant
	for _, id := range k.b.Implementers(under) {
		if id.Package != it.ID.Package || !token.IsExported(id.Name) {
			continue
		}
		value := &node.TypeRef{Spelling: id.Name, Target: id}
		switch d := k.declOf(value).(type) {
		case *node.Struct:
			if len(d.TypeParams) > 0 {
				continue
			}
		case *node.Enum:
		default:
			continue
		}
		if src, _ := k.sourceOf(value); src == sourceNone || !k.importable(value) {
			continue
		}
		satisfied, _ := k.b.Satisfies(value, under)
		out = append(out, Variant{
			Type: rules.EmitRef(value), Label: naming.Kebab(id.Name), Address: !satisfied.Satisfied(),
		})
	}
	return out
}

// membersOf returns the members of the enumeration d at the type ref that the
// file of the suite can refer to, and nil where the file cannot import ref.
func (k *caller) membersOf(d *node.Enum, ref *node.TypeRef) []emit.Value {
	if !k.importable(ref) {
		return nil
	}
	members := make([]emit.Value, 0, len(d.Variants))
	for _, v := range d.Variants {
		if member := emit.Name(v.ID); k.referable(&member) {
			members = append(members, member)
		}
	}
	return members
}

// twins returns the names of the first two fields of the struct st that the
// shape reader of the property engine reads under one name, and that name.
// The reader reads each readable field and each exported embedded type. It
// names one by the name of its json tag, and by its Go name where the tag
// states none.
func (k *caller) twins(st *node.Struct) (first, second, name string, found bool) {
	seen := map[string]string{}
	add := func(field string, tag reflect.StructTag) bool {
		text := tag.Get(jsonTag)
		read, _, _ := strings.Cut(text, ",")
		if text == propSkip || read == "" {
			read = field
		}
		if earlier, taken := seen[read]; taken {
			first, second, name = earlier, field, read
			return true
		}
		seen[read] = field
		return false
	}
	for _, f := range st.Fields {
		if k.readable(f) && add(f.Name, reflect.StructTag(f.Tag)) {
			return first, second, name, true
		}
	}
	for _, e := range st.Embeds {
		tag := reflect.StructTag(e.Tag)
		if field := embeddedName(e.Ref); token.IsExported(field) && tag.Get(propTag) != propSkip && add(field, tag) {
			return first, second, name, true
		}
	}
	return "", "", "", false
}

// interfaceOf returns the interface that the type ref refers to, through its
// aliases and its defined types, with the reference under to it. An inline
// interface returns an interface of the methods and the embedded types of
// its body, and the dynamic type an interface without members. interfaceOf
// reports false for a type that is no interface.
func (k *caller) interfaceOf(ref *node.TypeRef) (*node.Interface, *node.TypeRef, bool) {
	under := k.underlying(ref)
	if it, is := k.declOf(under).(*node.Interface); is {
		return it, under, true
	}
	if k.b.TypeOf(under).Form == symbol.FormDynamic {
		return &node.Interface{}, under, true
	}
	if b, parsed := parsedBody(under); parsed && b.isInterface {
		return &node.Interface{Methods: under.Methods, Embeds: b.embeds}, under, true
	}
	return nil, nil, false
}

// body is the body of an inline struct or interface. isInterface reports an
// interface, and embeds are the types that the interface embeds. exported
// reports a struct that embeds an exported type. The lowering of eidos keeps
// the embedded types of an inline body among the children of its reference
// alone.
type body struct {
	embeds      []*node.Embed
	isInterface bool
	exported    bool
}

// parsedBody parses the body of the inline type ref from its spelling, and
// reports false for a spelling that is not a struct or an interface. The
// children of a struct are the types of its entries, and the children of an
// interface are the types of its embeds and its constraint elements and the
// types of the parameters and the results of its methods, in the order of
// the body.
func parsedBody(ref *node.TypeRef) (body, bool) {
	expr, err := parser.ParseExpr(ref.Spelling)
	if err != nil {
		return body{}, false
	}
	switch t := expr.(type) {
	case *ast.StructType:
		var b body
		for i, f := range t.Fields.List {
			if len(f.Names) == 0 && token.IsExported(embeddedName(ref.Elems[i])) {
				b.exported = true
			}
		}
		return b, true
	case *ast.InterfaceType:
		b := body{isInterface: true}
		at := 0
		for _, f := range t.Methods.List {
			sig, isMethod := f.Type.(*ast.FuncType)
			if isMethod {
				at += fieldCount(sig.Params) + fieldCount(sig.Results)
				continue
			}
			b.embeds = append(b.embeds, &node.Embed{Ref: ref.Elems[at]})
			at++
		}
		return b, true
	}
	return body{}, false
}

// fieldCount returns the number of the types that the lowering of eidos
// takes from the list fields: one for each name, and one for an entry
// without a name.
func fieldCount(fields *ast.FieldList) int {
	if fields == nil {
		return 0
	}
	n := 0
	for _, f := range fields.List {
		n += max(len(f.Names), 1)
	}
	return n
}
