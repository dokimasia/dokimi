// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sampledparamtest contains the suite of Hasher. SHA256 is the subject
// of the suite.
package sampledparamtest

import (
	"crypto/sha256"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/sampledparam"
)

// interior is the lowest role of an interior node of a tree.
const interior sampledparam.Role = 0x80

// SHA256 implements [sampledparam.Hasher] with SHA-256. The corpus applies the
// suite of Hasher to it as a subject.
//
//dokimi:subject sampledparam.Hasher new=NewSHA256
type SHA256 struct{}

// The compiler checks that SHA256 implements [sampledparam.Hasher].
var _ sampledparam.Hasher = SHA256{}

// NewSHA256 returns a new SHA256.
func NewSHA256() SHA256 { return SHA256{} }

// HashTagged implements [sampledparam.Hasher]. It panics for the role of an
// interior node.
func (SHA256) HashTagged(r sampledparam.Role, data []byte) [32]byte {
	if r >= interior {
		//nolint:forbidigo // the contract of HashTagged refuses the role of an interior node with a panic
		panic("sampledparamtest: HashTagged takes a role below 0x80")
	}
	digest := sha256.Sum256(data)
	digest[0] ^= byte(r)
	return digest
}

// CombineTagged implements [sampledparam.Hasher]. It panics for the role of a
// leaf.
func (SHA256) CombineTagged(r sampledparam.Role, left, right [32]byte) [32]byte {
	if r < interior {
		//nolint:forbidigo // the contract of CombineTagged refuses the role of a leaf with a panic
		panic("sampledparamtest: CombineTagged takes a role of 0x80 or above")
	}
	var children [64]byte
	copy(children[:32], left[:])
	copy(children[32:], right[:])
	digest := sha256.Sum256(children[:])
	digest[0] ^= byte(r)
	return digest
}
