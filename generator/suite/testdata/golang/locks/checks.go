package locks

import (
	"testing"
	"testing/synctest"
)

// ReleaseEndsTheWait checks that Release ends the wait of Acquire.
//
//dokimi:check Lock.Acquire defect=neverReleases around=bubble
func ReleaseEndsTheWait(h *LockHandle) {}

// neverReleases returns a Lock whose Release keeps the lock.
func neverReleases(*LockHandle) Lock { return nil }

// bubble runs each run of a check in a testing/synctest bubble of its own.
func bubble(t *testing.T, run func(*testing.T)) { synctest.Test(t, run) }
