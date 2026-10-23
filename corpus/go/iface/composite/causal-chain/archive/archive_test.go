// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package archive_test

import (
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/archive"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/memory"
)

// newLog returns an archive over an empty log in memory, the subject of the
// conformance test that dokimi generates for Log.
func newLog() archive.Log {
	return archive.Log{Log: memory.New()}
}
