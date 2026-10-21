// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"embed"
	"text/template"

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
// exported.
const (
	pascalHelper = "suitePascal"
	camelHelper  = "suiteCamel"
)

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
		sdk.Funcs(template.FuncMap{pascalHelper: naming.Pascal, camelHelper: naming.Camel}),
	}
}
