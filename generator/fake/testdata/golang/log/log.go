// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package log appends entries, the fixture of a fake of the contract appender.
package log

import "context"

// Entry is an entry of the log.
type Entry struct {
	Body string
}

// Log appends entries.
//
//dokimi:suite
//dokimi:fake
//dokimi:out path=logtest/
type Log interface {
	// Append appends e and returns its position.
	//
	//dokimi:contract appender role=fn
	Append(ctx context.Context, e Entry) (int64, error)
}
