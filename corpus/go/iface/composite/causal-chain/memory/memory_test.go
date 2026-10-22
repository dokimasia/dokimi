// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package memory_test

import (
	"context"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"

	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/causalchaintest"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/memory"
)

// The log in memory passes the suite of Log and a check of its own. The
// application falsifies that check on its planted defect first.
func TestLog(t *testing.T) {
	t.Parallel()

	causalchaintest.RunLogSuite(t, memory.New, conform.With(refusesAnOrphan()))
}

// refusesAnOrphan returns the check that Append refuses an entry whose
// dependency the log does not have. Only the log in memory reports that entry
// with ErrUnadmitted. Its planted defect admits every entry.
func refusesAnOrphan() causalchaintest.LogCheck {
	return causalchaintest.LogCheck{
		ID:        causalchaintest.LogChecks.Append.Check("refuses-an-orphan"),
		Statement: "Append refuses an entry that depends on an entry that the log does not have",
		Apply: func(h *causalchaintest.LogHandle) {
			orphan := causalchain.Entry{ID: "orphan", DependsOn: []string{"missing"}, Body: "lost"}
			assert.ErrorIs(h, h.Subject().Append(h.Context(), orphan), memory.ErrUnadmitted, h.Statement())
		},
		Defect: func(h *causalchaintest.LogHandle) causalchain.Log {
			return causalchaintest.NewLogStub(h,
				causalchaintest.WithLogAppend(func(context.Context, causalchain.Entry) error { return nil }))
		},
	}
}
