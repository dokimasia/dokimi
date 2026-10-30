// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package xsssafetest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package xsssafetest

import (
	"context"
	"errors"
	"strings"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/xsssafe"
)

// InMemory implements [xsssafe.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject xsssafe.Mixed new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [xsssafe.Mixed].
var _ xsssafe.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Render implements [xsssafe.Mixed].
func (*InMemory) Render(ctx context.Context, in string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	escaped := strings.NewReplacer(
		"&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;",
	).Replace(in)
	return escaped, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("xsssafetest: nil context")
	}
	return ctx.Err()
}
