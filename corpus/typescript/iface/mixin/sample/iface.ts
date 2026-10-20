// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the sample mixin, which states that a
 * declared value is used as generated input rather than a random draw.
 *
 * @remarks
 * The interface carries whatever methods the law needs to be stateable. A
 * mixin whose law spans two calls cannot be hosted by a single method: there
 * would be nothing to compare against, and the generated subtest would pass
 * by having nothing to check.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=sampletest/
export interface Mixed {
  /**
   * Processes an input that the mixin pins, and returns the result.
   *
   * @remarks
   * The parameter is the point: with no argument there is nothing for a
   * sample to replace.
   */
  //dokimi:mixin sample builder=newInput
  process(input: string): Promise<string>;

  /**
   * Returns a new input, as the builder that the mixin names.
   *
   * @remarks
   * It sits on the interface because that is where the sibling resolver
   * looks: a mixin on a method resolves its partner against the owner's
   * method list, so a module-level factory would not be found however well it
   * fitted.
   */
  newInput(): Promise<string>;
}
