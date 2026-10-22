// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sharded is a Batched of the fixture batched-mixins that spreads its
// keys over shards.
package sharded

import (
	"context"
	"errors"
	"maps"
	"slices"
	"sync"
)

// ErrNotFound reports a key without a value.
var ErrNotFound = errors.New("sharded: no write stored the key")

// ListLimit is the most keys that List returns. It is the limit of the mixin
// bounded of List.
const ListLimit = 50

// Batched keeps each key in one of its shards, which a hash of the key
// selects.
type Batched struct {
	shards []shard
}

// shard is one part of the keys of a Batched.
type shard struct {
	values map[string]string
	mu     sync.Mutex
}

// New returns an empty Batched of n shards.
func New(n int) *Batched {
	b := &Batched{shards: make([]shard, n)}
	for i := range b.shards {
		b.shards[i].values = map[string]string{}
	}
	return b
}

// Put stores the value under the key in the shard of the key.
func (b *Batched) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s := b.shardOf(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
	return nil
}

// Read returns the value under the key from the shard of the key. It returns
// ErrNotFound for a key without a value.
func (b *Batched) Read(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	s := b.shardOf(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.values[key]
	if !found {
		return "", ErrNotFound
	}
	return value, nil
}

// List returns the keys of every shard in ascending order, at most ListLimit
// of them.
func (b *Batched) List(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var keys []string
	for i := range b.shards {
		s := &b.shards[i]
		s.mu.Lock()
		keys = slices.AppendSeq(keys, maps.Keys(s.values))
		s.mu.Unlock()
	}
	slices.Sort(keys)
	return keys[:min(len(keys), ListLimit)], nil
}

// shardOf returns the shard of the key, which a polynomial hash of the bytes
// of the key selects.
func (b *Batched) shardOf(key string) *shard {
	i := 0
	for j := range len(key) {
		i = (i*31 + int(key[j])) % len(b.shards)
	}
	return &b.shards[i]
}
