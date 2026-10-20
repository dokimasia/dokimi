// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package fault stamps the faults that the double of a method offers.
//
// fault is an annotator. It reads the directive fault on a method, in Go and
// in TypeScript:
//
//	//dokimi:fault errors=[ErrNotFound] retry=3 partition=key
//	Get(ctx context.Context, key string) (Value, error)
//
// The param errors lists error values of the package, and eidos resolves each
// name to a variable of the package. The double fails the next call with one
// of them through a helper of its own. The param retry is the attempt that a
// retry schedule lets succeed. The double fails each call before it, so retry
// is 2 or more. The param partition names a parameter of the method, and eidos
// resolves the name. The double fails the calls whose argument of the
// parameter matches a key, or the calls whose argument differs from it.
//
// # Facts
//
// fault stamps the identities of the error values of errors under
// [KeyErrors], the attempt of retry under [KeyRetry], and the identity of the
// parameter of partition under [KeyPartition]. The stub generator reads them
// and writes the helpers on the stub of the method: FaultNotFound for the
// error value ErrNotFound, RetrySchedule, FaultForPartition and
// FaultForOtherPartitions. fault provides the capability [Capability], so a
// plugin that reads the facts requires it and runs after it.
//
// # Findings
//
// fault reports [Infallible] at a directive on a method that cannot fail,
// such as a Go method without an error result, and [BadRetry] at a retry
// below 2.
//
// # Dependency position
//
// The package imports the root, diag, directive, meta, plugin, rules and
// symbol packages of the SDK facade of eidos, and the Go stdlib.
package fault
