// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// Part is the part of a plugin for one language. The interface of the parts
// of a plugin embeds it.
type Part interface {
	// Lang returns the language of the subjects of the part.
	Lang() symbol.Lang
	// Target returns the target of the part.
	Target() plugin.Target
	// Options returns the template tree and the helpers of the target.
	Options() []sdk.TargetOption
}

// Register declares the target of each of parts with its options on the
// builder b of a plugin.
func Register[P Part](b *sdk.Builder, parts []P) {
	for _, p := range parts {
		b.For(p.Target(), p.Options()...)
	}
}

// Served returns a handler that runs h with the part of the language of its
// match among parts. The handler takes the effect of h, so it generates where
// h generates and stamps where h stamps. For a match of a language without a
// part, the handler returns nil without running h.
func Served[M interface{ Lang() symbol.Lang }, P Part, E sdk.Effect](
	parts []P, h func(M, P, E) error,
) func(M, E) error {
	return func(m M, e E) error {
		for _, p := range parts {
			if p.Lang() == m.Lang() {
				return h(m, p, e)
			}
		}
		return nil
	}
}
