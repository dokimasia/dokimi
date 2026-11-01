package durable

import (
	"time"

	"go.dokimi.dev/assert/prop"
)

// LeavesTheWholeValue checks that a crash during Put leaves the whole value or
// no value.
//
//dokimi:check Keeper.Put defect=tearsWrites reference=reference second=second crash=tear
func LeavesTheWholeValue(h *KeeperHandle) {}

// KeepsEachPutAcrossARestart checks that Get returns each value that Put
// stored before a restart.
//
//dokimi:check Keeper.Get defect=forgetsOnRestart restart=forget
func KeepsEachPutAcrossARestart(h *KeeperHandle) {}

// KeepsTheLastPut checks that Get returns the value of the last Put of a key.
//
//dokimi:check Keeper.Get defect=keepsTheFirstPut options=putShares
func KeepsTheLastPut(h *KeeperHandle, c *prop.Case) {}

// ServesTheStoredValue checks that the backend of a Keeper serves the value
// that Put stored.
//
//dokimi:check Keeper.Get defect=servesAnotherValue values=[backendOf, clockOf]
func ServesTheStoredValue(h *KeeperHandle) {}

// Backend is the backend of a Keeper, which serves its values to a client
// without the Keeper.
type Backend interface {
	// Serve returns the value under the key.
	Serve(key string) string
}

// Clock is the clock of the backend of a Keeper.
type Clock interface {
	// Now returns the time of the clock.
	Now() time.Time
}

// servesAnotherValue returns a Keeper whose backend serves another value.
func servesAnotherValue(*KeeperHandle) Keeper { return nil }

// backendOf returns the backend of the Keeper k.
func backendOf(*KeeperHandle, Keeper) Backend { return nil }

// clockOf returns the clock of the backend of the Keeper k.
func clockOf(*KeeperHandle, Keeper) Clock { return nil }

// tearsWrites returns a Keeper whose Put writes a value in parts.
func tearsWrites(*KeeperHandle) Keeper { return nil }

// reference returns the reference implementation of the checks on tearsWrites.
func reference(*KeeperHandle) Keeper { return nil }

// second returns the second subject of the checks on tearsWrites.
func second(*KeeperHandle) Keeper { return nil }

// tear crashes the medium of k in the middle of the write of call.
func tear(*KeeperHandle, Keeper, func()) Keeper { return nil }

// forgetsOnRestart returns a Keeper whose restart forget replaces.
func forgetsOnRestart(*KeeperHandle) Keeper { return nil }

// forget returns an empty Keeper in place of k.
func forget(*KeeperHandle, Keeper) Keeper { return nil }

// keepsTheFirstPut returns a Keeper whose Put keeps the first value of a key.
func keepsTheFirstPut(*KeeperHandle) Keeper { return nil }

// putShares returns the shares of the outcomes of the Puts of each case.
func putShares(*KeeperHandle) []prop.Option { return nil }
