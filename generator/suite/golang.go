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
// identity of a check of a suite. lockHelper returns the lines of the lock of
// a plan.
const (
	pascalHelper   = "suitePascal"
	camelHelper    = "suiteCamel"
	ceilingsHelper = "suiteCeilings"
	durationHelper = "suiteDuration"
	defectHelper   = "suiteDefect"
	bodyHelper     = "suiteBody"
	joinHelper     = "suiteJoin"
	namesHelper    = "suiteNames"
	identityHelper = "suiteIdentity"
	lockHelper     = "suiteLock"
)

// identityIn returns the identity of the check segment of the callable of the
// declaration named declaration, such as Store.Get/smoke, and the name of the
// callable alone for an empty segment. methods reports a declaration with
// methods. The callable of a function is the function, so its identity has no
// declaration, such as Slug/smoke.
func identityIn(declaration string, methods bool, callable, segment string) string {
	name := declaration
	if methods {
		name += "." + callable
	}
	if segment == "" {
		return name
	}
	return name + "/" + segment
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
// the callable of the fresh subject, and State the prefix of the counter and
// the clock of the defect.
type bodyView struct {
	Defect *Defect
	Call   string
	State  string
}

// viewOfBody returns the view of the body of the planted defect d.
func viewOfBody(d *Defect, call, state string) bodyView {
	return bodyView{Defect: d, Call: call, State: state}
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
// unit that divides it, such as time.Millisecond or 2 * time.Microsecond. A
// template that writes the expression imports the package time.
func golangDuration(d time.Duration) string {
	unit, name := time.Nanosecond, "Nanosecond"
	for _, u := range units {
		if d%u.unit == 0 {
			unit, name = u.unit, u.name
			break
		}
	}
	if n := d / unit; n != 1 {
		return fmt.Sprintf("%d * time.%s", n, name)
	}
	return "time." + name
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
		}),
	}
}
