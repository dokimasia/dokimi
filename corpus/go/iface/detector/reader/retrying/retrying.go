// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package retrying is a Reader of the fixture reader that retries the failed
// reads of another Reader.
package retrying

import (
	"context"
	"errors"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/reader"
)

// Reader retries each failed Get of its source, up to its number of attempts.
// A miss is no failure, so Reader returns ErrNotFound without a retry.
type Reader struct {
	source   reader.Reader
	attempts int
}

// New returns a Reader that calls source at most attempts times for each Get.
func New(source reader.Reader, attempts int) *Reader {
	return &Reader{source: source, attempts: attempts}
}

// Get returns the value of the first attempt that succeeds, and ErrNotFound
// at the first attempt that reports it. It returns the error of the last
// attempt where every attempt fails.
func (r *Reader) Get(ctx context.Context, key string) (reader.Value, error) {
	var err error
	for range r.attempts {
		var v reader.Value
		v, err = r.source.Get(ctx, key)
		if err == nil || errors.Is(err, reader.ErrNotFound) {
			//nolint:wrapcheck // a caller matches ErrNotFound of the source with errors.Is
			return v, err
		}
	}
	//nolint:wrapcheck // a caller matches the error of the last attempt with errors.Is
	return reader.Value{}, err
}
