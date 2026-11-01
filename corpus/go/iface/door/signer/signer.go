// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package signer is the fixture of the doors of a value. Its signer signs the
// request of a body, and the backend of the signer serves the body to a client
// without the secret of the signer. The suite opens the backend of a signer
// through a door of a value. The signer in memory of signertest is the
// subject of the door.
package signer

import (
	"context"
	"errors"
)

// ErrNotFound is the error of Sign for a key without a body.
var ErrNotFound = errors.New("signer: not found")

// Signer stores bodies and signs the requests that read them.
//
//dokimi:suite
//dokimi:out path=signertest/
type Signer interface {
	// Put stores the body under the key.
	Put(ctx context.Context, key, body string) error
	// Sign returns a signed request of the body under the key, and
	// ErrNotFound for a key without a body.
	Sign(ctx context.Context, key string) (string, error)
}
