// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package fake is the generator of verified fakes. For an interface under the
// directive fake whose callables the reference forms of the conformance spec
// cover, it writes an in-memory implementation of the interface that is safe
// for concurrent use, and a test that applies the suite of the interface to
// it. The forms cover the contracts persister, cache and appender. A fake
// keeps the values that its callables store in a map under their keys, or
// appends them to a list. A load of an absent key returns the sentinel that
// the file of the fake declares. A generic interface gets a generic fake, and
// its test applies the suite at the witnesses of the interface.
//
// A declaration with a callable outside every form does not get a fake. fake
// reports a warning at the callable with the remedy of the gate
// reference-role.
package fake
