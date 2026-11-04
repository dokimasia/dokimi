// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package aeadtest contains the double and the suite of Mixed. InMemory is the
// subject of the suite.
package aeadtest

import (
	"crypto/aes"
	"crypto/cipher"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/aead"
)

// InMemory implements [aead.Mixed] with AES-GCM under a key that the
// constructor sets. The corpus applies the suite of Mixed to it as a subject.
//
//dokimi:subject aead.Mixed new=NewInMemoryFromHandle
type InMemory struct {
	cipher.AEAD
}

// The compiler checks that InMemory implements [aead.Mixed].
var _ aead.Mixed = (*InMemory)(nil)

// NewInMemoryFromHandle returns AES-128-GCM under a key of 16 zero bytes. It
// fails the check of h where the standard library refuses the key or the
// block.
func NewInMemoryFromHandle(h *MixedHandle) *InMemory {
	block, err := aes.NewCipher(make([]byte, 16))
	assert.NoError(h, err, "AES takes a key of 16 bytes")
	gcm, err := cipher.NewGCM(block)
	assert.NoError(h, err, "GCM takes a block of 16 bytes")
	return &InMemory{AEAD: gcm}
}
