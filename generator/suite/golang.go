// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"embed"
	"fmt"
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
const (
	pascalHelper   = "suitePascal"
	camelHelper    = "suiteCamel"
	ceilingsHelper = "suiteCeilings"
	durationHelper = "suiteDuration"
)

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

// golangExternal ends the import path of the external test package of a
// directory, as go/packages reports it, such as example.com/shop/store_test.
const golangExternal = "_test"

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
			durationHelper: golangDuration,
		}),
	}
}

// external reports whether the package of factory is the external test
// package of the package of subject, whose import path ends with _test.
func (golangPart) external(subject, factory *symbol.Identity) bool {
	return factory.Package == subject.Package+golangExternal
}
