// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"embed"
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// typescriptTemplates contains the file templates of the target typescript.
//
//go:embed templates/typescript
var typescriptTemplates embed.FS

// typescriptCaser cases the names that the templates of TypeScript write, as
// the settle of a TypeScript plan cases a declared name. It recognises no
// initialism, so each word keeps the case that its author wrote: the method
// of the check zero-on-error of findById is zeroOnError on the field
// findById.
var typescriptCaser = naming.New()

// typescriptContext is the parameter of a TypeScript closure that the probe
// passes the signal of its call.
const typescriptContext = "signal"

// optionsWord is the word that the alias of the options of an application
// joins to the name of its declaration, as StoreOptions joins it to Store.
const optionsWord = "Options"

// The templates of TypeScript spell a type of a module of the runtime through
// typeHelper, which imports the type for the type checker alone, and write a
// text as a string literal through stringHelper, in the quote of the plan.
// moduleHelper returns the module of a probe. identityOfHelper and
// defectOfHelper return the views of an identity and of a planted defect.
const (
	typeHelper       = "suiteType"
	stringHelper     = "suiteString"
	moduleHelper     = "suiteModule"
	identityOfHelper = "suiteIdentityOf"
	defectOfHelper   = "suiteDefectOf"
)

// probeModules is the directory of the modules of the probes of the runtime.
const probeModules = "@dokimi/runtime/probe/"

// typescriptPart is the part of suite for TypeScript.
type typescriptPart struct{}

// Lang returns the language TypeScript.
func (typescriptPart) Lang() symbol.Lang { return typescript.Lang }

// Target returns the target typescript.
func (typescriptPart) Target() plugin.Target { return typescript.Target }

// Options returns the template tree and the helpers of the target
// typescript.
func (typescriptPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(typescriptTemplates, typescript.Target),
		sdk.Funcs(template.FuncMap{
			pascalHelper: typescriptCaser.Pascal, camelHelper: typescriptCaser.Camel, typeHelper: typeOf,
			stringHelper: stringOf, moduleHelper: moduleOf, identityHelper: identityIn, lockHelper: lockLines,
			bodyHelper: viewOfBody, identityOfHelper: viewOfIdentity, defectOfHelper: viewOfTypescriptDefect,
			stepHelper: viewOfStep,
		}),
	}
}

// context returns signal, the parameter of a TypeScript closure that the
// probe passes the signal of its call.
func (typescriptPart) context() string { return typescriptContext }

// typed returns v as it is, because a TypeScript literal is a value of each
// alias of its type.
func (typescriptPart) typed(_ *emit.TypeRef, v *emit.Value) *emit.Value { return v }

// converted returns v as it is, as typed does.
func (typescriptPart) converted(_ *emit.TypeRef, v *emit.Value) *emit.Value { return v }

// names returns the names that [NamesOf] returns, with the alias of the
// options of an application after the type Options of the runtime, as
// StoreOptions.
func (typescriptPart) names(e *sdk.Emitter, decl string) Names {
	n := NamesOf(e, decl)
	n.Option = e.JoinName(optionsWord, decl)
	return n
}

// spelled returns the name that the settle of a TypeScript plan gives a
// declaration of the kind named name: Pascal case for a type and camel case
// for a function or a variable, as runStoreSuite.
func (typescriptPart) spelled(kind symbol.Kind, name string) string {
	switch kind {
	case symbol.KindStruct, symbol.KindInterface, symbol.KindAlias, symbol.KindEnum:
		return typescriptCaser.Pascal(name)
	default:
		return typescriptCaser.Camel(name)
	}
}

// pascal returns word in Pascal case without initialisms, as Id for id,
// because TypeScript writes an initialism as a word.
func (typescriptPart) pascal(word string) string { return typescriptCaser.Pascal(word) }

// shape returns the shape file of the values of the type ref, which
// prop.ofShape of @dokimi/assert reads into a generator. Each named record,
// such as an interface of data, is a definition of the file that a ref
// names, so a recursive type has a shape. A type that no shape of the
// vocabulary states, such as a function type or an interface with a method,
// is a ref to the name of the type, whose generator the consumer registers
// with prop.register in the setup file of Vitest.
func (typescriptPart) shape(k *caller, ref *node.TypeRef) string {
	w := shapeWalk{k: k, defs: map[string]any{}, names: map[recordKey]string{}}
	root := w.of(ref)
	if len(w.defs) > 0 {
		root["definitions"] = w.defs
	}
	// The encoder keeps the angle brackets of a type argument and the
	// ampersand of an intersection as they are.
	var b strings.Builder
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	// A shape of maps, lists, numbers and strings always encodes.
	_ = enc.Encode(root)
	return strings.TrimSuffix(b.String(), "\n")
}

// checkOpening returns the words Checks that, as the summary of a TSDoc
// comment opens with its verb.
func (typescriptPart) checkOpening(string) string { return tsChecksThat }

// tsChecksThat contains the words that open the TSDoc comment of a check.
const tsChecksThat = "Checks that"

// assertModule is the module of the assertions and of the property engine of
// TypeScript.
const assertModule = "@dokimi/assert"

// countsAllocations reports false, because V8 counts no allocations of a
// call.
func (typescriptPart) countsAllocations() bool { return false }

// takesCase reports whether ref is the case of the property engine of
// @dokimi/assert, as prop.Case, whatever name the module gives prop. A ref
// with a package is a named type.
func (typescriptPart) takesCase(ref *node.TypeRef) bool {
	return ref.Package == assertModule && ref.Spelling[strings.LastIndex(ref.Spelling, ".")+1:] == caseType
}

// caseSpelling returns prop.Case, the spelling of the case of the property
// engine in a message.
func (typescriptPart) caseSpelling() string { return "prop.Case" }

// start returns the starting value of the values of the type ref: the
// authored sample or the declared default of the alias that ref refers to,
// which a generated binding keeps or changes. A suite of TypeScript registers
// no generator of a starting value, so the generator of the draws of the type
// contains it.
func (typescriptPart) start(k *caller, ref *node.TypeRef) *emit.Value {
	d, aliased := k.declOf(ref).(*node.Alias)
	if !aliased {
		return nil
	}
	value, started := k.startValueOf(d, ref)
	if !started {
		return nil
	}
	return &value
}

// The keys of a shape and the shapes of the vocabulary of the shape files
// that a TypeScript suite writes, and the width of a number of TypeScript.
const (
	shapeKey      = "shape"
	nameKey       = "name"
	shapeBool     = "bool"
	shapeInt      = "int"
	shapeFloat    = "float"
	shapeString   = "string"
	shapeBytes    = "bytes"
	shapeList     = "list"
	shapeMap      = "map"
	shapeOptional = "optional"
	shapeRecord   = "record"
	shapeLiteral  = "literal"
	shapeRef      = "ref"
	numberWidth   = 64
)

// shapeWalk walks a type into its shape. defs are the definitions of the
// named records that the walk entered, and names the name of the definition
// of each of their instances.
type shapeWalk struct {
	k     *caller
	defs  map[string]any
	names map[recordKey]string
}

// recordKey is an instance of a record: its declaration and the TypeScript
// spelling of its type arguments, such as <number> of Box<number>.
type recordKey struct {
	id   symbol.Identity
	args string
}

// of returns the shape of the type ref: the shape of the vocabulary that
// states it, and otherwise a ref to its registered name.
func (w *shapeWalk) of(ref *node.TypeRef) map[string]any {
	switch ref.Form {
	case symbol.FormOptional:
		return w.wrap(shapeOptional, ref.Elems[0])
	case symbol.FormList, symbol.FormArray:
		return w.wrap(shapeList, ref.Elems[0])
	case symbol.FormNamed:
		return w.named(ref)
	default:
		return referred(ref.Spelling)
	}
}

// named returns the shape of the named type ref. A type of a library of
// TypeScript, such as Uint8Array, takes the shape of its form, and every
// other declaration the shape of its members. An Array and a ReadonlyArray
// are a list of their type argument, and a Map and a ReadonlyMap a map of
// their two type arguments. The reader converts a map to a Map, so a Record,
// a plain object, is a ref to its spelling.
func (w *shapeWalk) named(ref *node.TypeRef) map[string]any {
	if !typescript.Library(ref.Target.Package) {
		if decl, found := w.k.m.Reader().Lookup(ref.Target); found {
			return w.declared(ref, decl)
		}
	}
	switch t := w.k.b.TypeOf(ref); {
	case t.Form == symbol.FormBool:
		return map[string]any{shapeKey: shapeBool}
	case t.Form == symbol.FormText:
		return map[string]any{shapeKey: shapeString}
	case t.Form == symbol.FormBytes:
		return map[string]any{shapeKey: shapeBytes}
	case t.Form == symbol.FormScalar:
		return map[string]any{shapeKey: shapeFloat, "width": numberWidth}
	case t.Form == symbol.FormList:
		// The fold refuses a list without one type argument.
		return w.wrap(shapeList, ref.Args[0])
	case t.Form == symbol.FormMap && ref.Spelling != typescript.TypeRecord:
		// The fold refuses a map without two type arguments.
		return map[string]any{shapeKey: shapeMap, "key": w.of(ref.Args[0]), "of": w.of(ref.Args[1])}
	case ref.Spelling == typescript.TypeBigInt:
		// The reader converts an int wider than 32 bits to a bigint.
		return map[string]any{shapeKey: shapeInt, "width": numberWidth, "signed": true}
	default:
		return referred(spelling(ref))
	}
}

// spelling returns the TypeScript spelling of the type ref with its type
// arguments, such as Record<string, boolean>.
func spelling(ref *node.TypeRef) string { return ref.Spelling + argsOf(ref) }

// argsOf returns the TypeScript spelling of the type arguments of ref, such
// as <string, boolean>, and an empty string for a type without them.
func argsOf(ref *node.TypeRef) string {
	if len(ref.Args) == 0 {
		return ""
	}
	args := make([]string, 0, len(ref.Args))
	for _, a := range ref.Args {
		args = append(args, spelling(a))
	}
	return "<" + strings.Join(args, ", ") + ">"
}

// referred returns a ref to the shape definition name: a definition of the
// shape file, or a name whose generator the consumer registers.
func referred(name string) map[string]any { return map[string]any{shapeKey: shapeRef, nameKey: name} }

// wrap returns the shape named shape of the values of elem, such as a list of
// them.
func (w *shapeWalk) wrap(shape string, elem *node.TypeRef) map[string]any {
	return map[string]any{shapeKey: shape, "of": w.of(elem)}
}

// declared returns the shape of the type ref of the declaration decl: the
// shape of the target of an alias, the literal of the values of an
// enumeration, the ref of the definition of a class or an interface, and a
// ref to the spelling of a type parameter of a generic method. An alias of a
// structural type that no shape states, such as a function type, is a ref to
// the name of the alias, which the consumer registers.
func (w *shapeWalk) declared(ref *node.TypeRef, decl symbol.Symbol) map[string]any {
	switch d := decl.(type) {
	case *node.Alias:
		target := rules.Substitute(d.Target, d.TypeParams, ref.Args)
		shape := w.of(target)
		if target.Form != symbol.FormNamed && shape[shapeKey] == shapeRef {
			return referred(d.Name + argsOf(ref))
		}
		return shape
	case *node.Enum:
		return enumShape(d)
	case *node.Struct:
		return w.record(ref, d, d.TypeParams)
	case *node.Interface:
		return w.record(ref, d, d.TypeParams)
	default:
		return referred(spelling(ref))
	}
}

// record returns the ref of the definition of the record decl at the type
// arguments of ref, which the walk defines on its first entry. A record has
// the fields of the member walk of decl in its order, with the type
// arguments in place of params. A declaration with a method is no record,
// because a shape states data alone, so its shape is a ref to its registered
// name.
func (w *shapeWalk) record(ref *node.TypeRef, decl node.Declaration, params []*node.TypeParam) map[string]any {
	id, args := decl.Identity(), argsOf(ref)
	key := recordKey{id: id, args: args}
	if name, entered := w.names[key]; entered {
		return referred(name)
	}
	set, _ := w.k.b.MembersOf(decl)
	if slices.ContainsFunc(set.Members, func(m rules.Member) bool { return isMethod(m.Symbol) }) {
		return referred(id.Name + args)
	}
	name := id.Name + args
	for n := 2; w.defined(name); n++ {
		name = id.Name + strconv.Itoa(n) + args
	}
	w.names[key] = name
	fields := make([]any, 0, len(set.Members))
	for i := range set.Members {
		// A record without a method has fields alone.
		m, _ := set.Members[i].Symbol.(*node.Field)
		typ := rules.Substitute(m.Type, params, ref.Args)
		if m.Optional {
			typ = &node.TypeRef{Form: symbol.FormOptional, Elems: []*node.TypeRef{typ}}
		}
		fields = append(fields, []any{m.Name, w.of(typ)})
	}
	w.defs[name] = map[string]any{shapeKey: shapeRecord, "fields": fields}
	return referred(name)
}

// isMethod reports whether the member s is a method.
func isMethod(s symbol.Symbol) bool {
	_, method := s.(*node.Method)
	return method
}

// defined reports whether the walk defines a record under name.
func (w *shapeWalk) defined(name string) bool {
	for _, n := range w.names {
		if n == name {
			return true
		}
	}
	return false
}

// enumShape returns the literal of the values of the enumeration d: the text
// of each string member, the number of each numeric member, and for a member
// without an initializer the number after the member before it, from 0. The
// numbers are float literals, because the reader converts an int literal to
// a bigint. An enumeration with a member that no literal initializes, such
// as 1 << 2 or Infinity, is a ref to its registered name.
func enumShape(d *node.Enum) map[string]any {
	values := make([]any, 0, len(d.Variants))
	next := 0.0
	for _, v := range d.Variants {
		if text, quoted := typescript.Unquote(v.Value); quoted {
			values = append(values, map[string]any{"type": shapeString, "value": text})
			continue
		}
		number := next
		if v.Value != "" {
			// A sign is no part of a numeric literal.
			digits, negative := strings.CutPrefix(v.Value, "-")
			parsed, numeric := typescript.ParseNumber(digits)
			if !numeric {
				return referred(d.Name)
			}
			number = parsed
			if negative {
				number = -parsed
			}
		}
		values = append(values, map[string]any{"type": shapeFloat, "value": number})
		next = number + 1
	}
	return map[string]any{shapeKey: shapeLiteral, "values": values}
}

// typeOf returns the reference of the type name of the module specifier with
// the type arguments args, such as Handle<Store> of
// @dokimi/runtime/conform/handle. The reference has no target, so the speller
// imports it from its module as a type.
func typeOf(name, specifier string, args ...*emit.TypeRef) *emit.TypeRef {
	return &emit.TypeRef{Spelling: name, Package: specifier, Args: args}
}

// stringOf returns the string literal of text, which the value builtin
// writes in the quote of the plan.
func stringOf(text string) emit.Value { return emit.Literal(emit.LiteralString, text) }

// moduleOf returns the module of the probe id in the runtime, such as
// @dokimi/runtime/probe/under-cancelled.
func moduleOf(id string) string { return probeModules + id }

// identityView is the identity of a check that a TypeScript suite writes
// through the names of its callables. Checks is the variable of the names,
// Field the field of the callable, empty for a function, and Method the
// method that returns the identity. Base is the method of the identity of
// the sample binding of a generated binding, and empty for any other check.
type identityView struct {
	Checks string
	Field  string
	Method string
	Base   string
}

// viewOfIdentity returns the view of the identity of a check.
func viewOfIdentity(checks, field, method, base string) identityView {
	return identityView{Checks: checks, Field: field, Method: method, Base: base}
}

// typescriptDefectView is a planted defect of a TypeScript suite with the
// callables and the fields of its declaration, which the object of the
// defect forwards to the fresh subject, the alias of the handle, and the
// subject type.
type typescriptDefectView struct {
	Defect    *Defect
	Handle    string
	Subject   string
	Callables []Callable
	Fields    []string
}

// viewOfTypescriptDefect returns the view of the planted defect d.
func viewOfTypescriptDefect(
	d *Defect, callables []Callable, fields []string, handle, subject string,
) typescriptDefectView {
	return typescriptDefectView{Defect: d, Callables: callables, Fields: fields, Handle: handle, Subject: subject}
}
