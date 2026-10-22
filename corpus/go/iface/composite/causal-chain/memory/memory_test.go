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

// FuzzLogReplay fuzzes the property checks of Log.Replay on the log in
// memory. go test runs the two seeds, and go test -fuzz generates more inputs.
func FuzzLogReplay(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("a body of the root entry"))
	causalchaintest.FuzzLogSuite(f, causalchaintest.LogChecks.Replay.Callable,
		func(*causalchaintest.LogHandle) *memory.Log { return memory.New() })
}

// refusesAnOrphan returns the check that Append refuses an entry whose
// dependency the log does not have. Only the log in memory reports that entry
// with ErrUnadmitted. Its planted defect is the log in memory without the
// dependencies of each entry, so it admits every entry and passes the other
// checks of the suite.
func refusesAnOrphan() causalchaintest.LogCheck {
	return causalchaintest.LogCheck{
		ID:        causalchaintest.LogChecks.Append.Check("refuses-an-orphan"),
		Statement: "Append refuses an entry that depends on an entry that the log does not have",
		Apply: func(h *causalchaintest.LogHandle) {
			orphan := causalchain.Entry{ID: "orphan", DependsOn: []string{"missing"}, Body: "lost"}
			assert.ErrorIs(h, h.Subject().Append(h.Context(), orphan), memory.ErrUnadmitted, h.Statement())
		},
		Defect: func(h *causalchaintest.LogHandle) causalchain.Log {
			log := memory.New()
			return causalchaintest.NewLogStub(h, causalchaintest.LogStubDelegateTo(log),
				causalchaintest.WithLogAppend(func(ctx context.Context, e causalchain.Entry) error {
					e.DependsOn = nil
					return log.Append(ctx, e)
				}))
		},
	}
}
