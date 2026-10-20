// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel

import (
	"embed"
	"text/template"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangTemplates contains the file template of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of sentinel for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// pascalHelper is the name of the helper that writes the name of a package
// in Pascal case, as the name of a test function contains it.
const pascalHelper = "pascal"

// Options returns the template tree and the helper of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(golangTemplates, golang.Target),
		sdk.Funcs(template.FuncMap{pascalHelper: naming.Pascal}),
	}
}

// prefix returns base and the separator, because a message of Go opens with
// the name of its package, as in "store: not found".
func (golangPart) prefix(base string) string { return base + PrefixSeparator }

// inherent reports false, because the error of Go is an interface without
// fields.
func (golangPart) inherent(string) bool { return false }
