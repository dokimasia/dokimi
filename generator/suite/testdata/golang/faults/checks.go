package faults

import "go.dokimi.dev/conform"

// ServesTheStoredValue checks that Get returns the value that the backing
// store keeps, after a failed call of the store. Its callable has the name of
// the package, because other packages declare a Cache.
//
//dokimi:check faults.Cache.Get defect=servesStaleValues faults=faulty
func ServesTheStoredValue(h *CacheHandle) {}

// servesStaleValues returns a Cache whose Get serves the value that it cached
// before a failed call of the store.
func servesStaleValues(*CacheHandle) Cache { return nil }

// faulty returns a Cache over a backing store whose calls ask inj.
func faulty(*CacheHandle, *conform.Injector) Cache { return nil }
