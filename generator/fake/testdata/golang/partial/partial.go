// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package partial embeds an interface of the package store. A spec that
// loads partial without store leaves the members of Partial incomplete.
package partial

import "example.com/shop/store"

// Partial embeds the interface Store.
//
//dokimi:suite
//dokimi:fake
type Partial interface {
	store.Store
}
