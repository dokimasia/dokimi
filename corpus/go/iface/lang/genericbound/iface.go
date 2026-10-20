// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package genericbound is the language-axis fixture for a type parameter whose
// constraint the generator cannot reason about.
//
// `any` and `comparable` are the two bounds whose type set is known without
// reading anything: the first admits every type, the second every basic one.
// A named constraint is a reference into a package the generator never loaded,
// so nothing tells it which types satisfy the bound — and inventing one
// produces a companion that fails to compile for a reason the author could not
// have predicted.
//
// The witness directive is the answer: the source names the type each
// parameter is instantiated at, the witness annotator resolves and stamps
// them, and every generator reads the stamp through its rules. This fixture
// holds that path, as [generic] holds the derived one.
package genericbound

import (
	"context"
)

// Ordered is a constraint declared here rather than being one of Go's
// predeclared bounds, which is exactly what makes it opaque to the generator.
type Ordered interface {
	~int | ~int64 | ~string
}

// Score is the value type, so a witness naming a type from this package
// exercises the qualification a bare identifier needs — the companion lives in
// an external test package and reaches nothing here unqualified.
type Score struct {
	Points int
}

// Ranked is generic over an opaque constraint. Without the witness its
// companion would be a note, because no derivation can tell that `int`
// satisfies [Ordered].
//
// V names [Score] by full import path rather than bare: the annotator keeps
// a bare name unqualified — the spelling for a builtin — and the companion
// lives in an external test package, where an unqualified Score names
// nothing.
//
//dokimi:stub
//dokimi:suite
//dokimi:witness K=int V=Score
//dokimi:out path=genericboundtest/
type Ranked[K Ordered, V any] interface {
	Rank(ctx context.Context, key K) (V, error)
	Reset(ctx context.Context) error
}
