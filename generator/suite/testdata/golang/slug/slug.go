// Package slug writes the slugs of titles.
package slug

import "strings"

// Slug returns the slug of a title: its words in lower case, joined by
// hyphens.
//
//dokimi:suite
func Slug(title string, keep ...string) string {
	return strings.Join(append(strings.Fields(strings.ToLower(title)), keep...), "-")
}
