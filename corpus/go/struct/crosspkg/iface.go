// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package crosspkg is the fixture for a default that names a symbol outside
// the declaring package.
//
// Two notations, both through the tag. A qualifier resolves against the
// file's own imports, which reads naturally and is what an author writes
// for a package the file already uses — Timeout, Retries and Region all
// declare theirs this way. A full import path needs no import at all, which
// is what an author reaches for when the package would otherwise be
// imported solely to feed a default — Zone declares its through that form,
// naming the same constant Region does. Both resolve for a string-typed
// field exactly as they do for a numeric one.
//
// [go.dokimi.dev/dokimi/corpus/go/struct/seed] holds the far side.
//
// Routing is declared once for the package rather than repeated on each
// type: every builder here lands in the same companion package, so a
// per-struct directive is the same statement written N times, and the Nth
// copy is the one that gets forgotten.
//
//dokimi:out path=crosspkgtest/
package crosspkg

import (
	"time"

	"go.dokimi.dev/dokimi/corpus/go/struct/seed"
)

// Settings draws every default from another package, through the tag.
//
//dokimi:builder
type Settings struct {
	// Timeout names a constant in a package imported for real code below,
	// through the tag.
	Timeout time.Duration `default:"time.Second"`

	// Retries names a constant in a package this file imports, through the
	// tag.
	Retries int `default:"seed.Retries"`

	// Region names a constant in a package this file imports, through the
	// tag.
	Region string `default:"seed.Region"`

	// Zone names the same constant by full import path, through the tag.
	Zone string `default:"seed.Region"`

	// Deadline exists so the file's `time` import is used by real code and not
	// only by a directive.
	Deadline time.Time
}

// Mirrored seeds from a companion in another package, named by full path.
//
//dokimi:builder defaults=seed.ConfigDefaults
type Mirrored = seed.Config

// Shorthand names the same companion by qualifier rather than by path, which
// is the notation an author reaches for first and the one every field in this
// file already uses.
//
// It exists because the two notations take different routes: a full path names
// itself, while a qualifier is only meaningful against the import block of the
// file that wrote it. The qualifier route was documented and unreachable — no
// fixture wrote one on a `defaults=` key, so nothing established that the
// generator could get from a struct to its file at all.
//
//dokimi:builder defaults=seed.ConfigDefaults
type Shorthand = seed.Config

var _ = seed.Region
