// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package balancer is the fixture of a balancer over two instances whose calls
// fail on demand. The balancer sends each call with one key to the instance
// that a hash of the key selects.
package balancer

import "context"

// Balancer sends calls to its instances.
//
//dokimi:suite machineries=[faults]
//dokimi:out path=balancertest/
type Balancer interface {
	// Send sends a call with the key to the instance of the key.
	//
	//dokimi:mixin sticky key=key
	Send(ctx context.Context, key string) error
}
