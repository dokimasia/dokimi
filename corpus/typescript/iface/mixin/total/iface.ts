// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the total mixin, which states that a
 * callable answers for every input in a named domain rather than failing on
 * some of them.
 *
 * @remarks
 * Two methods, differing in whether they return a promise. The claim is about
 * what a single call returns for a hostile input rather than about how two
 * calls relate, so one would state it — but the law's closure has to handle
 * an asynchronous and a synchronous call alike, and a fixture with one of
 * them leaves the other untested.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Contains the input that the directive names: the empty string, which is
 * where the domain of strings ends.
 *
 * @remarks
 * Declared rather than derived. A domain is prose for a reader — "strings"
 * says nothing a generator can build — and a drawn sample sits comfortably
 * inside it, so a subject that gives way at the boundary passes on one.
 */
export const edge = "";

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=totaltest/
export interface Mixed {
  /** Returns the class of any input in the declared domain. */
  //dokimi:mixin total domain=strings edge=edge
  classify(input: string): Promise<string>;

  /** Returns the normal form of any input in the declared domain. */
  //dokimi:mixin total domain=strings edge=edge
  normalize(input: string): string;
}
