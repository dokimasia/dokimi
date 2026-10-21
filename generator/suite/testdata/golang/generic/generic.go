// Package generic caches values by key.
package generic

// Cache keeps values by key.
//
//dokimi:suite
type Cache[K comparable, V any] interface {
	// Get returns the value under the key, and whether the cache has one.
	Get(key K) (V, bool)
	// Set stores the value under the key.
	Set(key K, value V)
}
