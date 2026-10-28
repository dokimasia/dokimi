// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

import (
	"embed"
	"encoding/json"
	"strings"
	"sync"
)

// files contains the version and the rendered JSON files of the spec. The
// test of the package compares the entries that each function returns with
// the entries that encoding/json reads from the same file, so a file that
// does not decode fails the test, and no function reports a decoding error.
//
//go:embed VERSION *.json
var files embed.FS

// Version returns the version of the spec from the file VERSION, such as
// 1.0.0.
func Version() string {
	return version()
}

var version = sync.OnceValue(func() string {
	data, _ := files.ReadFile("VERSION")
	return strings.TrimSpace(string(data))
})

// keyed is a pointer to an entry E of the spec with the identifier K.
type keyed[K ~string, E any] interface {
	*E
	key() K
}

// table contains the entries of one section of a file in the order of the
// file, and the position of each entry by its identifier.
type table[K ~string, E any] struct {
	at      map[K]int
	entries []E
}

// of returns the entry with the identifier id, and false when the section
// has no such entry.
func (t table[K, E]) of(id K) (E, bool) {
	i, ok := t.at[id]
	if !ok {
		var zero E
		return zero, false
	}
	return t.entries[i], true
}

// section returns a function that decodes the section of a file on its first
// call and returns the same table on every call.
func section[K ~string, E any, P keyed[K, E]](file, name string) func() table[K, E] {
	return sync.OnceValue(func() table[K, E] {
		data, _ := files.ReadFile(file)
		var doc map[string][]E
		_ = json.Unmarshal(data, &doc)
		entries := doc[name]
		at := make(map[K]int, len(entries))
		for i := range entries {
			at[P(&entries[i]).key()] = i
		}
		return table[K, E]{entries: entries, at: at}
	})
}
