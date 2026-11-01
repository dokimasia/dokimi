// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package packagescope is the language-axis fixture for a kit that declares
// at package scope the names of the packages that its suite imports. Go
// refuses an import whose name a declaration of the package takes, so the
// suite of the kit imports go.dokimi.dev/spec and go.dokimi.dev/conform/probe
// under names of their own.
package packagescope

// Slugger turns titles into slugs.
//
//dokimi:suite
//dokimi:out path=packagescopetest/
type Slugger interface {
	// Slug returns the words of title in lower case, joined by hyphens.
	Slug(title string) string
}
