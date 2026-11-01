// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus_test

import (
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/internal/corpus"
	"go.dokimi.dev/dokimi/spec"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("Update", func(t *testing.T) {
		t.Parallel()

		t.Run("rewrites the declarations of a fixture from its lock, its records and the spec", func(t *testing.T) {
			t.Parallel()
			stale := expectation()
			stale.Declarations = map[string]corpus.Declaration{"store": {}}

			got := corpus.Update([]corpus.Expectation{stale}, []corpus.Outcome{outcome()})

			assert.Equal(t, got, []corpus.Expectation{expectation()},
				"the declarations are those of the outcome, and the subjects and the controls stay")
		})

		t.Run("keeps the declarations of a fixture without an outcome", func(t *testing.T) {
			t.Parallel()
			other := expectation()
			other.Fixture = "iface/mixin/ttlperwrite"

			got := corpus.Update([]corpus.Expectation{other}, []corpus.Outcome{outcome()})

			assert.Equal(t, got, []corpus.Expectation{other}, "a fixture without an outcome keeps its declarations")
		})

		t.Run("keys the declaration of a function by the key of the function", func(t *testing.T) {
			t.Parallel()
			o := outcome()
			o.Lock = append(o.Lock, ttl+"\tslug-of\tshape:computation\tsilent\t-")

			got := corpus.Update([]corpus.Expectation{expectation()}, []corpus.Outcome{o})

			assert.Equal(t, got[0].Declarations["slug-of"], corpus.Declaration{
				Checks:   map[string]corpus.Check{},
				Refusals: map[string]spec.GateID{},
				Unclaimed: []corpus.Unclaimed{
					{Classification: "shape:computation", Callable: "slug-of", Stance: spec.Silent},
				},
			}, "a function is its own declaration")
		})

		t.Run("writes the partners of a check in the order of their names", func(t *testing.T) {
			t.Parallel()
			o := outcome()
			o.Lock[lineExpires] = ttl + "\tmixed.read/expires\texpires\tcheck\tread=mixed.get,write=mixed.put"

			got := corpus.Update([]corpus.Expectation{expectation()}, []corpus.Outcome{o})

			assert.Equal(t, got[0].Declarations["mixed"].Checks["mixed.read/expires"].Partners,
				map[string]string{"read": "mixed.get", "write": "mixed.put"}, "each partner has its callable")
			assert.Equal(t, got[0].Declarations["mixed"].Checks["mixed.read/expires"].Reason, spec.NotExpired,
				"the reason code is that of the claim")
		})
	})
}
