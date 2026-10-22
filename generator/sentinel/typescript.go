// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel

import (
	"embed"
	"slices"
	"text/template"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// typescriptTemplates contains the file template of the target typescript.
//
//go:embed templates/typescript
var typescriptTemplates embed.FS

// errorFields are the fields of the global class Error of TypeScript: its
// name, its message, its stack and its cause.
var errorFields = []string{"name", "message", "stack", "cause"}

// casesHelper is the name of the helper that writes the name of a function of
// cases, and typescriptCases ends that name.
const (
	casesHelper     = "sentinelCases"
	typescriptCases = "Cases"
)

// typescriptCaser cases the names of the functions of cases. It recognises no
// initialism, so each word keeps the case that its author wrote.
var typescriptCaser = naming.New()

// typescriptPart is the part of sentinel for TypeScript.
type typescriptPart struct{}

// Lang returns the language TypeScript.
func (typescriptPart) Lang() symbol.Lang { return typescript.Lang }

// Target returns the target typescript.
func (typescriptPart) Target() plugin.Target { return typescript.Target }

// Options returns the template tree and the helper of the target typescript.
func (typescriptPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(typescriptTemplates, typescript.Target),
		sdk.Funcs(template.FuncMap{casesHelper: casesOf}),
	}
}

// casesOf returns the name of the function that declares the cases of the
// member named member of the error class named class, such as
// notFoundErrorMessageCases for message of NotFoundError.
func casesOf(class, member string) string {
	return typescriptCaser.Camel(class) + typescriptCaser.Pascal(member) + typescriptCases
}

// prefix returns the empty prefix, because a message of TypeScript has no
// convention for its opening.
func (typescriptPart) prefix(string) string { return "" }

// inherent reports whether name is a field of the global class Error.
func (typescriptPart) inherent(name string) bool { return slices.Contains(errorFields, name) }
