// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// The helpers of a plan share one namespace, so the name of each helper of
// suite opens with suite. The templates of Go write each name that a record
// contains in the case that the settle of the Go plan gives its declaration:
// pascalHelper writes an exported name, and camelHelper a name that is not
// exported. ceilingsHelper returns the segment of the check of each ceiling
// of a [Budget], and durationHelper writes a duration as a Go expression.
//
// defectHelper and bodyHelper return the views of a planted defect that its
// templates read, joinHelper joins words with commas, namesHelper joins the
// names of the results of a planted defect, and identityHelper writes the
// identity of a callable of a suite. lockHelper returns the lines of the lock of
// a plan, sequencesHelper the view of the generated sequences of a
// declaration, worldHelper the view of the world of its subject,
// durableHelper the view of the calls of a durable probe, stepHelper the view
// of the step of a callable, and txnCallHelper the view of one call of a
// transaction.
const (
	pascalHelper    = "suitePascal"
	camelHelper     = "suiteCamel"
	ceilingsHelper  = "suiteCeilings"
	durationHelper  = "suiteDuration"
	defectHelper    = "suiteDefect"
	bodyHelper      = "suiteBody"
	joinHelper      = "suiteJoin"
	namesHelper     = "suiteNames"
	identityHelper  = "suiteIdentity"
	lockHelper      = "suiteLock"
	sequencesHelper = "suiteSequences"
	worldHelper     = "suiteWorld"
	durableHelper   = "suiteDurable"
	stepHelper      = "suiteStep"
	txnCallHelper   = "suiteTxnCall"
)

// durableView is the calls of a durable probe with the subject type of the
// suite, which the template writes as the type argument of probe.Durable and
// as the type of the subject that each call takes.
type durableView struct {
	Durable *Durable
	Subject string
}

// viewOfDurable returns the view of the calls d of a durable probe.
func viewOfDurable(d *Durable, subject string) durableView {
	return durableView{Durable: d, Subject: subject}
}

// sequencesView is the parts of the generated sequences of a declaration with
// the names of the suite that their template writes: the subject type and the
// name of the declaration.
type sequencesView struct {
	Sequences   *Sequences
	Subject     string
	Declaration string
}

// viewOfSequences returns the view of the parts q of the generated sequences.
func viewOfSequences(q *Sequences, subject, declaration string) sequencesView {
	return sequencesView{Sequences: q, Subject: subject, Declaration: declaration}
}

// worldView is the world of the subject of a declaration with the generated
// sequences whose steps are its operations, and with the names of the suite
// that its template writes: the subject type and the name of the
// declaration.
type worldView struct {
	World       *World
	Sequences   *Sequences
	Subject     string
	Declaration string
}

// viewOfWorld returns the view of the world w over the parts q of the
// generated sequences.
func viewOfWorld(w *World, q *Sequences, subject, declaration string) worldView {
	return worldView{World: w, Sequences: q, Subject: subject, Declaration: declaration}
}

// identityIn returns the identity of the callable of the declaration named
// declaration, such as Store.Get. methods reports a declaration with methods.
// The callable of a function is the function, so its identity is the name of
// the function, such as Slug.
func identityIn(declaration string, methods bool, callable string) string {
	if methods {
		return declaration + "." + callable
	}
	return declaration
}

// defectView is a planted defect with the names of the suite that its
// template writes: the alias of the subject type that its type embeds, the
// alias of the handle and the subject type.
type defectView struct {
	Defect   *Defect
	Embedded string
	Handle   string
	Subject  string
}

// viewOfDefect returns the view of the planted defect d.
func viewOfDefect(d *Defect, embedded, handle, subject string) defectView {
	return defectView{Defect: d, Embedded: embedded, Handle: handle, Subject: subject}
}

// bodyView is the body of the override of the planted defect Defect: Call is
// the callable of the fresh subject, and State the prefix of the counter, the
// clock, the handle and the dependency of the defect. Subject is the fresh
// subject that the defect wraps, on which the commit of a planted defect of
// transactions begins a late transaction.
type bodyView struct {
	Defect  *Defect
	Call    string
	State   string
	Subject string
}

// viewOfBody returns the view of the body of the planted defect d.
func viewOfBody(d *Defect, call, state, subject string) bodyView {
	return bodyView{Defect: d, Call: call, State: state, Subject: subject}
}

// stepView is the step of a callable with the names of the suite that its
// template writes: the subject type, and the function of the param normalize
// of the directive suite, which is nil without the param.
type stepView struct {
	Step      *Step
	Normalize *emit.TypeRef
	Subject   string
}

// viewOfStep returns the view of the step st.
func viewOfStep(st *Step, subject string, normalize *emit.TypeRef) stepView {
	return stepView{Step: st, Subject: subject, Normalize: normalize}
}

// txnCallView is one call of a transaction with the expression of the subject
// on which a call that is no method of the handle calls its callable.
type txnCallView struct {
	Subject string
	Call    TxnCall
}

// viewOfTxnCall returns the view of the call c of a transaction.
func viewOfTxnCall(c TxnCall, subject string) txnCallView {
	return txnCallView{Call: c, Subject: subject}
}

// namesOfOutcomes returns the names of the results outcomes, joined by
// commas.
func namesOfOutcomes(outcomes []Outcome) string {
	names := make([]string, 0, len(outcomes))
	for _, o := range outcomes {
		names = append(names, o.Name)
	}
	return strings.Join(names, ", ")
}

// joined returns words joined by commas.
func joined(words []string) string { return strings.Join(words, ", ") }

// units are the units of the time package above a nanosecond from the
// largest down, with the names of their constants.
var units = []struct {
	name string
	unit time.Duration
}{
	{"Hour", time.Hour},
	{"Minute", time.Minute},
	{"Second", time.Second},
	{"Millisecond", time.Millisecond},
	{"Microsecond", time.Microsecond},
}

// golangDuration returns the Go expression of the duration d in its largest
// unit that divides it, such as time.Millisecond or 2 * time.Microsecond. pkg
// is the name under which the file imports the package time.
func golangDuration(d time.Duration, pkg string) string {
	unit, name := time.Nanosecond, "Nanosecond"
	for _, u := range units {
		if d%u.unit == 0 {
			unit, name = u.unit, u.name
			break
		}
	}
	if n := d / unit; n != 1 {
		return fmt.Sprintf("%d * %s.%s", n, pkg, name)
	}
	return pkg + "." + name
}

// golangPart is the part of suite for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Options returns the template tree and the helpers of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(golangTemplates, golang.Target),
		sdk.Funcs(template.FuncMap{
			pascalHelper: naming.Pascal, camelHelper: naming.Camel, ceilingsHelper: Budget.segments,
			durationHelper: golangDuration, defectHelper: viewOfDefect, bodyHelper: viewOfBody, joinHelper: joined,
			namesHelper: namesOfOutcomes, identityHelper: identityIn, lockHelper: lockLines,
			sequencesHelper: viewOfSequences, worldHelper: viewOfWorld, durableHelper: viewOfDurable,
			stepHelper: viewOfStep, txnCallHelper: viewOfTxnCall,
		}),
	}
}

// golangContext is the name of the context of the call of a closure.
const golangContext = "ctx"

// context returns ctx, the parameter of a Go closure that the probe passes
// the context of its call.
func (golangPart) context() string { return golangContext }

// typed returns v as a value of the type ref: the conversion of a literal or
// a name, whose constant can be untyped, and any other value as it is, which
// already has its type.
func (golangPart) typed(ref *emit.TypeRef, v *emit.Value) *emit.Value {
	if v.Kind == emit.ValueLiteral || v.Kind == emit.ValueName {
		return new(emit.Conversion(ref, *v))
	}
	return v
}

// converted returns the conversion of v to the type ref, because a Go call
// generator assigns the sample to a variable whose type the value states.
func (golangPart) converted(ref *emit.TypeRef, v *emit.Value) *emit.Value {
	return new(emit.Conversion(ref, *v))
}

// names returns the names that [NamesOf] returns.
func (golangPart) names(e *sdk.Emitter, decl string) Names { return NamesOf(e, decl) }

// shape returns the empty string, because the property engine of Go reads the
// type of each value that it draws.
func (golangPart) shape(*caller, *node.TypeRef) string { return "" }

// checkOpening returns the name of the check and the words checks that, as a
// doc comment of Go opens with the name of its declaration.
func (golangPart) checkOpening(name string) string { return name + checksThat }

// takesCase reports whether ref is a pointer to the case of the property
// engine, whatever name the file gives the package of the engine.
func (golangPart) takesCase(ref *node.TypeRef) bool { return pointerTo(ref, propPackage, caseType) }

// caseSpelling returns *prop.Case, the spelling of the case of the property
// engine in a message.
func (golangPart) caseSpelling() string { return "*prop.Case" }

// countsAllocations reports true, because the testing package of Go counts
// the allocations of a call.
func (golangPart) countsAllocations() bool { return true }

// start returns nil, because the suite of Go registers the generator of the
// starting value of each type with one.
func (golangPart) start(*caller, *node.TypeRef) *emit.Value { return nil }

// spelled returns name, because the settle of a Go plan keeps the name of
// each exported declaration of a suite, which is in Pascal case already.
func (golangPart) spelled(_ symbol.Kind, name string) string { return name }

// pascal returns word in Pascal case under the initialisms of Go, as ID for
// id.
func (golangPart) pascal(word string) string { return naming.Pascal(word) }
