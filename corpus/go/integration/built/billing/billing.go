// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package billing declares a struct whose builder is in another package than
// the double of the package built.
//
//dokimi:out path=billingtest/
package billing

// Invoice is a struct with a builder.
//
//dokimi:builder
type Invoice struct {
	Number string
}
