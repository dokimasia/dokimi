// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"encoding/json"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"
)

func TestRules(t *testing.T) {
	t.Parallel()

	t.Run("files", func(t *testing.T) {
		t.Parallel()

		for _, name := range []string{"claims", "laws", "defects", "gates", "unclaimed", "references"} {
			t.Run("has a schema that accepts "+name+".json", func(t *testing.T) {
				t.Parallel()
				var schema jsonschema.Schema
				assert.NoError(t, json.Unmarshal([]byte(files.Read(t, "schemas/"+name+".schema.json")), &schema),
					"the schema decodes")
				resolved, err := schema.Resolve(nil)
				assert.NoError(t, err, "the schema resolves")
				var doc any
				assert.NoError(
					t,
					json.Unmarshal([]byte(files.Read(t, name+".json")), &doc),
					"the rendered file decodes",
				)

				assert.NoError(t, resolved.Validate(doc), "the schema accepts the rendered file")
			})
		}
	})
}
