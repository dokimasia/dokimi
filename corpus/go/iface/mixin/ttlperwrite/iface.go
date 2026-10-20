// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ttlperwrite is the mixin-axis fixture for a lifetime carried on
// the value rather than fixed by the directive.
//
// The sibling fixture `ttl` stamps `duration=1m`: one lifetime for every
// write, declared once. This one stamps no duration at all, and each
// write says how long its own entry lasts. Both are `ttl`, and the claim
// is the same sentence — an entry stops being readable once its lifetime
// has run out — but only this shape has a field a defect can reach for,
// which is what lets the claim be proven rather than argued.
//
// Both halves state it now. lifetime= names the member, so the claim
// reads the number off the write's answer instead of the directive, and
// the comparability walk reaches the time.Duration the Entry holds. The
// fixture spent a while stating nothing while those were upstream asks.
package ttlperwrite

import (
	"context"
	"errors"
	"time"
)

// ErrExpired is what a read past the lifetime reports.
//
// Named by the not-found= param, so the law compares the identity the
// declaration gave rather than merely noticing that something failed.
var ErrExpired = errors.New("ttlperwrite: entry expired")

// Entry is the payload, carrying its own lifetime.
//
// Nothing declares what a drawn Entry looks like. It used to: a derived
// value set the one member the roles aim by and left the rest at their
// zero, and a zero Lifetime is this mixin's own spelling for "forever",
// so the lapse check bound and asserted nothing. A drawn value now
// carries every member, this one included.
type Entry struct {
	// Key identifies the entry.
	Key string

	// Body is what a read answers with.
	Body string

	// Lifetime is how long the entry stays readable. Zero means forever,
	// which is what makes stripping it a defect rather than a refusal:
	// the write still lands and the read still answers, and the only
	// thing that changed is that nothing ever expires.
	Lifetime time.Duration
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=ttlperwritetest/
type Mixed interface {
	// Put stores an entry, starting the lifetime it carries.
	Put(ctx context.Context, e Entry) error

	// Read returns it while that lifetime holds, and the sentinel after.
	//
	// lifetime= rather than duration=: the entry says how long it lasts,
	// so the number does not exist until the write has answered.
	//dokimi:mixin ttl put=Put read=Read lifetime=Lifetime not-found=ErrExpired
	Read(ctx context.Context, key string) (Entry, error)
}
