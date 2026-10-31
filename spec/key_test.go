// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/spec"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("Key", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name         string
			giveIdentity string
			want         string
		}{
			{
				name:         "writes the words of a name in lower case with hyphens",
				giveIdentity: "Store.ReadAll/zero-on-error",
				want:         "store.read-all/zero-on-error",
			},
			{
				name:         "splits a name in camel case",
				giveIdentity: "Store.readAll/zero-on-error",
				want:         "store.read-all/zero-on-error",
			},
			{
				name:         "splits a name in snake case",
				giveIdentity: "Store.read_all/zero-on-error",
				want:         "store.read-all/zero-on-error",
			},
			{
				name:         "keeps an initialism as one word",
				giveIdentity: "HTTPServer.Serve/cancel",
				want:         "http-server.serve/cancel",
			},
			{
				name:         "keeps the segment with its case",
				giveIdentity: "Store.Get/expires#0-generated",
				want:         "store.get/expires#0-generated",
			},
			{
				name:         "returns the key of a check without a callable",
				giveIdentity: "OrderStore/agrees-with-the-reference",
				want:         "order-store/agrees-with-the-reference",
			},
			{
				name:         "returns the key of a callable",
				giveIdentity: "Store.PutAll",
				want:         "store.put-all",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, spec.Key(tt.giveIdentity), tt.want, "the key joins the check across languages")
			})
		}
	})
}
