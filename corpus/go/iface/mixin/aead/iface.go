// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package aead is the mixin-axis fixture of the builders of the mixin sample
// on the methods that an interface promotes from an embed, beside the fixture
// sample.
//
// Seal and Open panic for a nonce that is not NonceSize bytes, as
// [cipher.AEAD] documents, so a derived sample of a nonce fails each check.
// Mixed declares both methods again with the signatures of the embed, and
// each takes its inputs from the builder of its mixin sample.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package aead

import "crypto/cipher"

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=aeadtest/
type Mixed interface {
	cipher.AEAD

	// Seal encrypts and authenticates plaintext, as [cipher.AEAD] documents.
	//
	//dokimi:mixin sample builder=SealInputs
	Seal(dst, nonce, plaintext, additionalData []byte) []byte

	// Open authenticates and decrypts ciphertext, as [cipher.AEAD] documents.
	//
	//dokimi:mixin sample builder=OpenInputs
	Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error)
}

// SealInputs returns the inputs of a call of Seal on a, with a nonce of
// NonceSize bytes.
func SealInputs(a Mixed) (dst, nonce, plaintext, additionalData []byte) {
	return nil, make([]byte, a.NonceSize()), []byte("plaintext"), nil
}

// OpenInputs returns the inputs of a call of Open on a: the ciphertext that
// Seal returns for the inputs of SealInputs.
func OpenInputs(a Mixed) (dst, nonce, ciphertext, additionalData []byte) {
	_, nonce, plaintext, additionalData := SealInputs(a)
	return nil, nonce, a.Seal(nil, nonce, plaintext, additionalData), additionalData
}
