// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package archive is a Log of the fixture causal-chain that keeps its entries
// in a log in memory. Its spec builds each archive through a factory of the
// external test package, so the package declares no constructor for tests.
package archive

import (
	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/memory"
)

// Log keeps the entries of an archive in a log in memory.
//
//dokimi:subject causalchain.Log new=newLog
type Log struct {
	*memory.Log
}

// The compiler checks that Log implements [causalchain.Log].
var _ causalchain.Log = Log{}
