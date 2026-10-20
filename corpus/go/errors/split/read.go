// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package split

import "errors"

// The read half of the contract. First in source order, so this file's name is
// what the generated checks are named after.
var (
	// ErrNotFound reports a key that was never written.
	ErrNotFound = errors.New("split: not found")

	// ErrDenied reports a read the caller may not perform.
	ErrDenied = errors.New("split: denied")
)
