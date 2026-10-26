// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package resolved is the fixture of an application that states the values of
// a type that the suite cannot build. A resolution has no exported field, and
// only a registry builds one with a limit, so the suite starts each resolution
// from its zero value. The application of the suite in the package memory
// passes the resolutions of a registry through conform.Using.
package resolved

import "context"

// Resolution contains the settings that a registry resolves. Its zero value
// has the limit 0.
type Resolution struct {
	limit int
}

// Limit returns the limit of the resolution.
func (r Resolution) Limit() int { return r.limit }

// Registry resolves the settings of an application.
type Registry struct {
	limit int
}

// NewRegistry returns a registry whose resolutions have the limit limit.
func NewRegistry(limit int) *Registry { return &Registry{limit: limit} }

// Resolve returns the resolution of the registry.
func (r *Registry) Resolve() Resolution { return Resolution{limit: r.limit} }

// ApplyRequest is the request of Apply.
type ApplyRequest struct {
	Key        string
	Resolution Resolution
}

// Applier applies requests under the limit of their resolution.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=resolvedtest/
type Applier interface {
	// Apply applies the request under its key.
	Apply(ctx context.Context, req ApplyRequest) error
	// Limit returns the limit of the resolution of the request that Apply
	// applied under key, and 0 for a key without a request.
	Limit(ctx context.Context, key string) (int, error)
}
