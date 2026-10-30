// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package durabletest contains the suite of Store, the checks that a consumer
// adds to the suite with the doors of their planted defects, and the subject
// of the suite.
package durabletest

import (
	"context"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/prop"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/durable"
)

// KeepsEachPutAcrossARestart checks that Get returns the value that Put
// stored before a restart.
//
//dokimi:check durable.Store.Get defect=forgetsOnRestart restart=forget crash=crash
func KeepsEachPutAcrossARestart(h *StoreHandle) {
	s := h.Subject()
	assert.NoError(h, s.Put(h.Context(), "k", "v"), "Put stores the value")
	got, _ := h.Restart(s).Get(h.Context(), "k")
	assert.Equal(h, got, "v", h.Statement())
}

// keptOld and keptNew are the labels of the outcomes of a crash during a Put:
// the value before the Put, and the whole value of the Put.
const (
	keptOld = "the crash keeps the value before the Put"
	keptNew = "the crash keeps the value of the Put"
)

// LeavesTheWholeValue checks that Get returns the value before a Put or the
// whole value of the Put after a crash during the Put.
//
//dokimi:check durable.Store.Put defect=tearsWrites restart=reopenTorn crash=tear options=crashShares
func LeavesTheWholeValue(h *StoreHandle, c *prop.Case) {
	before := h.Subject()
	assert.NoError(h, before.Put(h.Context(), "k", "old"), "Put stores the first value")
	after := h.Crash(before, func() { _ = before.Put(h.Context(), "k", "new value") })
	got, _ := after.Get(h.Context(), "k")
	switch got {
	case "old":
		c.Classify(keptOld)
	case "new value":
		c.Classify(keptNew)
	}
	assert.True(h, got == "old" || got == "new value", h.Statement())
}

// crashShares returns the shares of the outcomes of a crash: each outcome is
// the outcome of a fifth of the cases at least.
func crashShares(*StoreHandle) []prop.Option {
	return []prop.Option{prop.Require(keptOld, 0.2), prop.Require(keptNew, 0.2)}
}

// forgetsOnRestart returns a journal, whose restart forget replaces with an
// empty journal. Its crash is the crash of a journal.
func forgetsOnRestart(*StoreHandle) durable.Store { return NewJournal() }

// forget returns a journal over an empty medium in place of the restarted
// store, as a store that loses its values on a restart.
func forget(*StoreHandle, durable.Store) durable.Store { return NewJournal() }

// torn is a journal whose Put writes the first half of each value and then
// the whole value. A cut of its medium between the two writes leaves half of
// the value.
type torn struct {
	*Journal
}

// Put writes the first half of the value, then the whole value.
func (t torn) Put(ctx context.Context, key, value string) error {
	if err := t.Journal.Put(ctx, key, value[:len(value)/2]); err != nil {
		return err
	}
	return t.Journal.Put(ctx, key, value)
}

// tearsWrites returns a torn journal over an empty medium.
func tearsWrites(*StoreHandle) durable.Store { return torn{Journal: NewJournal()} }

// reopenTorn returns a new torn journal over the medium of the torn journal
// s, as reopen does for a journal.
func reopenTorn(_ *StoreHandle, s durable.Store) durable.Store {
	return torn{Journal: &Journal{medium: tornMedium(s)}}
}

// tear cuts the medium of the torn journal s after the first write of call,
// runs call, and returns a torn journal over the medium.
func tear(_ *StoreHandle, s durable.Store, call func()) durable.Store {
	medium := tornMedium(s)
	medium.Cut(1)
	call()
	medium.Heal()
	return torn{Journal: &Journal{medium: medium}}
}

// tornMedium returns the medium of the torn journal s, and a new medium for a
// store that is no torn journal.
func tornMedium(s durable.Store) *Medium {
	if t, isTorn := s.(torn); isTorn {
		return t.medium
	}
	return NewMedium()
}
