// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package signertest contains the suite of Signer, a check that a consumer
// adds to the suite with the door of a value of its planted defect, and the
// subject of the suite.
package signertest

import (
	"context"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/signer"
)

// Backend is the backend of a signer, which serves the body of a signed
// request to a client without the secret of the signer.
type Backend interface {
	// Serve returns the body that the signed request reads, and false for a
	// request that the backend refuses.
	Serve(request string) (string, bool)
}

// ServesTheSignedBody checks that the backend serves the body that Put stored
// for a request that Sign signed.
//
//dokimi:check signer.Signer.Sign defect=signsWithAnotherSecret values=[defectBackend]
func ServesTheSignedBody(h *SignerHandle) {
	s := h.Subject()
	assert.NoError(h, s.Put(h.Context(), "k", "body"), "Put stores the body")
	request, err := s.Sign(h.Context(), "k")
	assert.NoError(h, err, "Sign signs the request of the body")
	body, served := conform.Open[Backend](h, s).Serve(request)
	assert.True(h, served && body == "body", h.Statement())
}

// anotherSecret is a signer in memory whose Sign signs each request with
// another secret, so its backend refuses the request.
type anotherSecret struct {
	*Memory
}

// Sign returns the key with the signature of another secret.
func (s anotherSecret) Sign(ctx context.Context, key string) (string, error) {
	return s.signed(ctx, key, "another secret")
}

// signsWithAnotherSecret returns a signer in memory whose Sign signs with
// another secret than its backend reads.
func signsWithAnotherSecret(*SignerHandle) signer.Signer { return anotherSecret{Memory: NewMemory()} }

// defectBackend returns the backend of the signer in memory that the planted
// defect s wraps.
func defectBackend(_ *SignerHandle, s signer.Signer) Backend {
	if d, isDefect := s.(anotherSecret); isDefect {
		return d.Backend()
	}
	return refusing{}
}
