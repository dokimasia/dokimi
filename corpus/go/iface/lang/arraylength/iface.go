// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package arraylength is the fixture of a suite whose inputs are arrays whose
// lengths are constants. The length of a salt is a constant of another
// package, the length of a pad is a constant of this package, and the length
// of a header is an expression of both. The length of a key is a constant
// that the package of the suite cannot refer to, so the suite refuses the
// checks of Key that pass a sample of a key request.
package arraylength

import (
	"context"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/arraylength/sizes"
)

// PadSize is the size of a pad in bytes.
const PadSize = 8

// keySize is the size of a key in bytes.
const keySize = 4

// Request is a request to stage a payload.
type Request struct {
	Salt [sizes.SaltSize]byte // the salt of the payload
	Pad  [PadSize]byte        // the pad of the payload
}

// KeyRequest is a request to stage a key.
type KeyRequest struct {
	Key [keySize]byte // the key
}

// Stager stages the bytes of payloads.
//
//dokimi:suite
//dokimi:out path=arraylengthtest/
type Stager interface {
	// Stage stages the salt and the pad of the request.
	Stage(ctx context.Context, req Request) error
	// Seal stages the salt.
	Seal(ctx context.Context, salt [sizes.SaltSize]byte) error
	// Pad stages the pad.
	Pad(ctx context.Context, pad [PadSize]byte) error
	// Header stages the header, which is a salt and a pad.
	Header(ctx context.Context, header [sizes.SaltSize + PadSize]byte) error
	// Key stages the key of the request.
	Key(ctx context.Context, req KeyRequest) error
}
