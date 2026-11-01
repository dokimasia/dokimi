// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package packagescopetest

import (
	"strings"
	"unicode"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/packagescope"
)

// probe contains the titles that the check of Slug reads.
var probe = []string{"Hello World", "  Go  Dokimi  ", "one"}

// spec returns the slug of title that the sequential specification of
// Slugger states.
func spec(title string) string {
	var b strings.Builder
	for i, word := range strings.Fields(title) {
		if i > 0 {
			b.WriteByte('-')
		}
		b.WriteString(strings.Map(unicode.ToLower, word))
	}
	return b.String()
}

// FollowsTheSpecification checks that Slug returns the slug that spec returns
// for each title of probe.
//
//dokimi:check packagescope.Slugger.Slug defect=keepsTheCase
func FollowsTheSpecification(h *SluggerHandle) {
	for _, title := range probe {
		assert.Equal(h, h.Subject().Slug(title), spec(title), h.Statement())
	}
}

// cased is a slugger whose Slug keeps the case of the words.
type cased struct{}

// Slug returns the words of title joined by hyphens, in their case.
func (cased) Slug(title string) string { return strings.Join(strings.Fields(title), "-") }

// keepsTheCase returns a slugger whose Slug keeps the case of the words.
func keepsTheCase(*SluggerHandle) packagescope.Slugger { return cased{} }
