// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package retrying_test

import (
	"errors"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/fault"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/reader"
	"go.dokimi.dev/dokimi/corpus/go/iface/detector/reader/readertest"
	"go.dokimi.dev/dokimi/corpus/go/iface/detector/reader/retrying"
)

// errUnavailable is the error of a source that cannot serve a read.
var errUnavailable = errors.New("the source is unavailable")

// outage is a fault of the consumer. It fails each read of a key that opens
// with the prefix of a region.
type outage struct {
	region string
}

// ShouldTrigger fails each call whose key opens with the region of the
// outage.
func (o outage) ShouldTrigger(c readertest.ReaderGetCall, _ fault.Clock) (bool, error) {
	if !strings.HasPrefix(c.Key, o.region) {
		return false, nil
	}
	return true, errUnavailable
}

// Reset keeps the region, because an outage has no counter.
func (outage) Reset() {}

// The retrying reader is tested against the double of Reader, whose faults
// make its source fail on demand.
func TestReader(t *testing.T) {
	t.Parallel()

	want := reader.Value{Key: "eu/1", Body: "body"}

	t.Run("Get", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value of the attempt at which the source recovers", func(t *testing.T) {
			t.Parallel()

			source := readertest.NewReaderStub(t)
			source.OnGet.Returns(want, nil).RetrySchedule(errUnavailable)
			got, err := retrying.New(source, 3).Get(t.Context(), want.Key)
			assert.NoError(t, err, "the third attempt succeeds")
			assert.Equal(t, got, want, "Get returns the value of the source")
			assert.Equal(t, source.OnGet.CallCount(), 3, "Get calls the source once for each attempt")
		})

		t.Run("returns the error of the last attempt of a source that does not recover", func(t *testing.T) {
			t.Parallel()

			source := readertest.NewReaderStub(t)
			source.OnGet.Returns(want, nil).RetrySchedule(errUnavailable)
			_, err := retrying.New(source, 2).Get(t.Context(), want.Key)
			assert.ErrorIs(t, err, errUnavailable, "both attempts fail before the source recovers")
			assert.Equal(t, source.OnGet.CallCount(), 2, "Get stops after its last attempt")
		})

		t.Run("returns ErrNotFound without a retry", func(t *testing.T) {
			t.Parallel()

			source := readertest.NewReaderStub(t)
			source.OnGet.Returns(want, nil).FaultNotFound()
			_, err := retrying.New(source, 3).Get(t.Context(), want.Key)
			assert.ErrorIs(t, err, reader.ErrNotFound, "Get reports the miss of the source")
			assert.Equal(t, source.OnGet.CallCount(), 1, "a retry does not find a missing key")
		})

		t.Run("fails the reads of the faulted key alone", func(t *testing.T) {
			t.Parallel()

			source := readertest.NewReaderStub(t)
			source.OnGet.Returns(want, nil).FaultForPartition("eu/1", errUnavailable, 1)
			r := retrying.New(source, 2)
			_, err := r.Get(t.Context(), "eu/1")
			assert.ErrorIs(t, err, errUnavailable, "each read of the faulted key fails")
			_, err = r.Get(t.Context(), "us/1")
			assert.NoError(t, err, "the read of another key succeeds")
		})

		t.Run("absorbs every second failed read of a region in an outage", func(t *testing.T) {
			t.Parallel()

			source := readertest.NewReaderStub(t)
			source.OnGet.Returns(want, nil).SetFault(fault.And[readertest.ReaderGetCall](
				outage{region: "eu/"},
				fault.NewCountedFault[readertest.ReaderGetCall](errUnavailable, 2),
			))
			r := retrying.New(source, 2)
			for _, key := range []string{"eu/1", "eu/2", "us/1"} {
				_, err := r.Get(t.Context(), key)
				assert.NoError(t, err, "the second attempt succeeds where the first fails")
			}
			assert.Equal(t, source.OnGet.CallCount(), 4, "the read of eu/2 takes a second attempt")
		})
	})
}
