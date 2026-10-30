// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package htmlescape contains a control of the suite of Mixed: a renderer that
// escapes through html.EscapeString of the standard library.
package htmlescape

import (
	"context"
	"html"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/xsssafe"
)

// Renderer implements [xsssafe.Mixed] with html.EscapeString, which writes a
// quotation mark as &#34; where the subject writes &quot;. The corpus applies
// the suite of Mixed to it as a control of the stateless declarations.
//
//dokimi:subject xsssafe.Mixed new=New
type Renderer struct{}

// The compiler checks that Renderer implements [xsssafe.Mixed].
var _ xsssafe.Mixed = Renderer{}

// New returns a Renderer.
func New() Renderer { return Renderer{} }

// Render implements [xsssafe.Mixed].
func (Renderer) Render(ctx context.Context, in string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return html.EscapeString(in), nil
}
