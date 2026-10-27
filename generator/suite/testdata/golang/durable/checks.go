package durable

import "go.dokimi.dev/assert/prop"

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
