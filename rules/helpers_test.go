// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"encoding/json"
	"maps"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"
)

// keysOf returns the identifiers of entries.
func keysOf[E any, K ~string](entries []E, key func(E) K) []string {
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		ids = append(ids, string(key(e)))
	}
	return ids
}

// ruleIDs returns the identifiers of the rules of a section of a rendered
// file, in sorted order, as encoding/json reads them.
func ruleIDs(t *testing.T, name, section string) []string {
	t.Helper()
	var doc map[string]map[string]json.RawMessage
	assert.NoError(t, json.Unmarshal([]byte(files.Read(t, name+".json")), &doc), name+".json decodes")
	return slices.Sorted(maps.Keys(doc[section]))
}
