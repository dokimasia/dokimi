// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package packagescopetest contains the suite of Slugger, the sequential
// specification of Slug and the titles that a check of Slug reads, and the
// subject of the suite. The specification is spec and the titles are probe,
// so the kit declares the names of two packages that its suite imports.
package packagescopetest

import (
	"strings"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/packagescope"
)

// InMemory implements [packagescope.Slugger] with the functions of the
// package strings. The corpus applies the suite of Slugger to it as a
// subject.
//
//dokimi:subject packagescope.Slugger new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [packagescope.Slugger].
var _ packagescope.Slugger = (*InMemory)(nil)

// NewInMemory returns a slugger.
func NewInMemory() *InMemory { return &InMemory{} }

// Slug implements [packagescope.Slugger].
func (*InMemory) Slug(title string) string {
	return strings.Join(strings.Fields(strings.ToLower(title)), "-")
}
