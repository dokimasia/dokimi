// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"embed"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of suite for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Options returns the template tree of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{sdk.TemplatesOf(golangTemplates, golang.Target)}
}
