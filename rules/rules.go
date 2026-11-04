// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"cmp"
	"embed"
	"encoding/json"
	"maps"
	"slices"
	"sync"
)

// files contains the rendered JSON files of the rules. The test of the
// command of the renderer renders the YAML files again and compares the
// result with these files, so a file that does not decode fails the test,
// and no function reports a decoding error.
//
//go:embed claims.json laws.json defects.json gates.json unclaimed.json references.json
var files embed.FS

// table contains the entries of one section of a file in their order, and
// the position of each entry by its identifier.
type table[K comparable, E any] struct {
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

// add appends the entry e with the identifier id.
func (t *table[K, E]) add(id K, e E) {
	t.at[id] = len(t.entries)
	t.entries = append(t.entries, e)
}

// decode returns the rules of the section name of the file by their
// identifiers.
func decode[K comparable, E any](file, name string) map[K]E {
	data, _ := files.ReadFile(file)
	var doc map[string]map[K]E
	_ = json.Unmarshal(data, &doc)
	return doc[name]
}

// joiner is a pointer to a rule E that takes the runtime entry S of the spec
// with the identifier of the rule.
type joiner[S, E any] interface {
	*E
	join(s *S)
}

// joined returns a function that decodes the rules of the section name of
// the file on its first call, and returns the same table on every call. The
// table joins each runtime entry of the spec with the rule of the same
// identifier, in the order of the spec. The test of each kind fails for an
// entry of the spec without a rule, so no entry has the zero rule.
func joined[K comparable, S, E any, P joiner[S, E]](file, name string, runtime func() []S,
	id func(*S) K,
) func() table[K, E] {
	return sync.OnceValue(func() table[K, E] {
		byID := decode[K, E](file, name)
		t := table[K, E]{at: make(map[K]int, len(byID))}
		entries := runtime()
		for i := range entries {
			s := &entries[i]
			e := byID[id(s)]
			P(&e).join(s)
			t.add(id(s), e)
		}
		return t
	})
}

// namer is a pointer to a rule E without a runtime entry, which takes its
// identifier K from the key of its entry.
type namer[K, E any] interface {
	*E
	name(id K)
}

// listed returns a function that decodes the rules of the section name of
// the file on its first call, and returns the same table on every call. The
// table contains each rule with its identifier, in the order of the
// identifiers.
func listed[K cmp.Ordered, E any, P namer[K, E]](file, name string) func() table[K, E] {
	return sync.OnceValue(func() table[K, E] {
		byID := decode[K, E](file, name)
		t := table[K, E]{at: make(map[K]int, len(byID))}
		for _, id := range slices.Sorted(maps.Keys(byID)) {
			e := byID[id]
			P(&e).name(id)
			t.add(id, e)
		}
		return t
	})
}
