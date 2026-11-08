// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim

import (
	"embed"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of sim for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Tested reports whether the file at the path is a test file of Go: its name
// ends in _test.go, as the go command reads it.
func (golangPart) Tested(path string) bool {
	return strings.HasSuffix(path, golang.TestSuffix+golang.Extension)
}

// The helpers of a plan share one namespace, so the name of each helper of
// sim opens with sim. quoteHelper writes a text as a Go string literal, and
// fieldHelper returns the type of the field of a record by its name.
const (
	quoteHelper = "simQuote"
	fieldHelper = "simField"
)

// Options returns the template tree and the helpers of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(golangTemplates, golang.Target),
		sdk.Funcs(template.FuncMap{quoteHelper: strconv.Quote, fieldHelper: fieldOf}),
	}
}

// fieldOf returns the type of the field named name of the record r, and an
// error where r has no such field. The layout of eidos respells the name of a
// field in the case of its language, such as Simulation for simulation, so
// fieldOf compares names without their case.
func fieldOf(r *emit.Struct, name string) (*emit.TypeRef, error) {
	for _, f := range r.Fields.Items() {
		if strings.EqualFold(f.Name, name) {
			return f.Type, nil
		}
	}
	return nil, fmt.Errorf("sim: the record %s has no field %s", r.Name, name)
}
