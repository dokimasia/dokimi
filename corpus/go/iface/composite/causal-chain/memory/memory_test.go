// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package memory_test

import (
	"context"
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi-go/conform"
	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/causalchaintest"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/memory"
)

// The log in memory passes the check of a consumer that Replay returns the
// body that Append admitted. The application falsifies the check on its
// planted defect first.
func TestLog(t *testing.T) {
	t.Parallel()

	causalchaintest.RunLog(t, memory.New, conform.With(keepsTheBody()))
}

// keepsTheBody returns the check that Replay returns the body that Append
// admitted. Its planted defect is the double of Log over a log in memory,
// whose Append drops the body of the entry.
func keepsTheBody() causalchaintest.LogCheck {
	return causalchaintest.LogCheck{
		ID:        causalchaintest.LogChecks.Replay.Check("keeps-the-body"),
		Statement: "Replay returns the body that Append admitted",
		Apply: func(h *causalchaintest.LogHandle) {
			s := h.Subject()
			assert.NoError(h, s.Append(h.Context(), causalchain.RootEntry()), "Append admits the root entry")
			got, err := s.Replay(h.Context())
			assert.NoError(h, err, "Replay reads the log")
			assert.Length(h, got, 1, "Replay returns the root entry")
			assert.Equal(h, got[0].Body, causalchain.RootEntry().Body, h.Statement())
		},
		Defect: func(h *causalchaintest.LogHandle) causalchain.Log {
			fake := memory.New()
			return causalchaintest.NewLogStub(h, causalchaintest.LogStubDelegateTo(fake),
				causalchaintest.WithLogAppend(func(ctx context.Context, e causalchain.Entry) error {
					e.Body = ""
					return fake.Append(ctx, e)
				}))
		},
	}
}
