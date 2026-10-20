// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for a trailing rest parameter, which
 * changes both the forwarding of a call and the shape of the recorded call.
 *
 * @remarks
 * This axis varies the type system rather than the classification, so these
 * break generators independently of any directive.
 *
 * @packageDocumentation
 */

/**
 * Represents a finder that takes a rest list of keys.
 *
 * @remarks
 * Forwarding the list needs the spread at the call site, and recording it
 * needs an array field rather than a scalar.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=variadictest/
export interface Finder {
  /** Returns the keys that the finder knows among the keys. */
  find(...keys: string[]): Promise<string[]>;
  /** Returns at most limit of the keys that the finder knows among the keys. */
  findWithLimit(limit: number, ...keys: string[]): Promise<string[]>;
}
