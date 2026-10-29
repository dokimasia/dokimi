// Package cached caches integers in memory.
package cached

import "example.com/shop/generic"

// Cache is the cache of integers by integer in memory. The directive suite of
// generic.Cache has no prefix. It excuses a derived check of Get, whose
// derivation reads the witnesses of generic.Cache.
//
//dokimi:subject generic.Cache new=New
//dokimi:excuse Cache.Get/zero-on-miss reason="the cache returns the zero value of int for each miss already"
type Cache struct {
	values map[int]int
}

// New returns an empty cache.
func New() *Cache {
	return &Cache{values: map[int]int{}}
}

// Get returns the value under the key, and whether the cache has one.
func (c *Cache) Get(key int) (int, bool) {
	value, found := c.values[key]
	return value, found
}

// Set stores the value under the key.
func (c *Cache) Set(key, value int) {
	c.values[key] = value
}

// Compile-time proof that Cache satisfies generic.Cache at its witnesses.
var _ generic.Cache[int, int] = (*Cache)(nil)
