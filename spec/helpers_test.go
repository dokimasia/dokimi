// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"go.dokimi.dev/assert"
)

// validates checks that the schema of the rendered file name accepts the
// file.
func validates(t *testing.T, name string) {
	t.Helper()
	schemaData, err := os.ReadFile("schemas/" + name + ".schema.json")
	assert.NoError(t, err, "the schema is readable")
	var schema jsonschema.Schema
	assert.NoError(t, json.Unmarshal(schemaData, &schema), "the schema decodes")
	resolved, err := schema.Resolve(nil)
	assert.NoError(t, err, "the schema resolves")
	data, err := os.ReadFile(name + ".json")
	assert.NoError(t, err, "the rendered file is readable")
	var doc any
	assert.NoError(t, json.Unmarshal(data, &doc), "the rendered file decodes")

	assert.NoError(t, resolved.Validate(doc), "the schema accepts the rendered file")
}

// idsOf returns the identifiers of the entries of a section of a rendered
// file, in the order of the file, as encoding/json reads them.
func idsOf(t *testing.T, name, section string) []string {
	t.Helper()
	data, err := os.ReadFile(name + ".json")
	assert.NoError(t, err, "the rendered file is readable")
	var doc map[string][]struct {
		ID string `json:"id"`
	}
	assert.NoError(t, json.Unmarshal(data, &doc), "the rendered file decodes")
	ids := make([]string, 0, len(doc[section]))
	for _, e := range doc[section] {
		ids = append(ids, e.ID)
	}
	return ids
}

// keysOf returns the identifiers of entries.
func keysOf[E any, K ~string](entries []E, key func(E) K) []string {
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		ids = append(ids, string(key(e)))
	}
	return ids
}
