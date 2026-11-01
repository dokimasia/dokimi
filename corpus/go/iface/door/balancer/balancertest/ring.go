// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package balancertest contains the suite of Balancer and the subject of the
// suite.
package balancertest

import (
	"context"
	"fmt"
	"hash/fnv"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/balancer"
)

// instance is one instance behind a ring, which takes the key of a call.
type instance func(ctx context.Context, key string) error

// Ring is the Balancer over two instances. The corpus applies the suite of
// Balancer to it as a subject, over instances whose calls fail on demand.
//
//dokimi:subject balancer.Balancer new=NewRing options=RingOptions
type Ring struct {
	instances [2]instance
}

// The compiler checks that Ring implements [balancer.Balancer].
var _ balancer.Balancer = (*Ring)(nil)

// NewRing returns a ring over two instances that accept every call.
func NewRing() *Ring {
	accept := func(context.Context, string) error { return nil }
	return &Ring{instances: [2]instance{accept, accept}}
}

// Send implements [balancer.Balancer].
func (r *Ring) Send(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return r.instances[h.Sum32()%uint32(len(r.instances))](ctx, key)
}

// RingOptions returns the door of a ring over two instances whose calls fail
// on demand.
func RingOptions() []BalancerOption {
	return []BalancerOption{conform.Faultable(faulty)}
}

// faulty returns a ring over two instances, each of which asks a dependency of
// the injector of its own.
func faulty(_ *BalancerHandle, inj *conform.Injector) balancer.Balancer {
	r := &Ring{}
	for i := range r.instances {
		dependency := inj.Dependency()
		r.instances[i] = func(context.Context, string) error {
			if err := dependency.Call("Send"); err != nil {
				return fmt.Errorf("balancertest: send: %w", err)
			}
			return nil
		}
	}
	return r
}
