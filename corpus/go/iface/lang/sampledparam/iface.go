// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sampledparam is the language-axis fixture for the sample of a
// parameter. A //dokimi:sample on a parameter gives each generated call of
// its callable the value of the sample.
//
// HashTagged takes a role below 0x80, and CombineTagged takes a role of 0x80
// or above, so a leaf of a tree and an interior node never share a role. Each
// method panics for a role of the other half. The benchmarks and the budgets
// of the suite take the calls of the call generator. They run only because
// the role of each method states a sample of its own half.
package sampledparam

// Role is the role of a hash.
type Role uint8

// Hasher hashes the leaves and the interior nodes of a tree under tagged
// roles.
//
//dokimi:suite
//dokimi:out path=sampledparamtest/
type Hasher interface {
	// HashTagged hashes the data of a leaf under the role r. It panics for a
	// role of 0x80 or above.
	//dokimi:budget allocs=0
	HashTagged(
		//dokimi:sample value=0x01
		r Role,
		data []byte,
	) [32]byte

	// CombineTagged hashes the digests of two children under the role r. It
	// panics for a role below 0x80.
	//dokimi:budget allocs=0
	CombineTagged(
		//dokimi:sample value=0x81
		r Role,
		left, right [32]byte,
	) [32]byte
}
