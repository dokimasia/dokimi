// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the xsssafe mixin, which states that a
 * value rendered into markup cannot close a tag or open a script.
 *
 * @remarks
 * One method is enough here: the claim is about what a single call returns
 * for a hostile or unusual input, not about how two calls relate.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Contains the value that the directive names: markup that closes a tag and
 * opens a script, which a renderer must not return.
 *
 * @remarks
 * Declared rather than derived. A drawn string is "test-in" and nothing in it
 * is dangerous, so a check over one passes a renderer that escapes nothing at
 * all.
 */
export const unsafe = "<script>alert(1)</script>";

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=xsssafetest/
export interface Mixed {
  /** Returns the raw input escaped for embedding in markup. */
  //dokimi:mixin xss-safe unsafe=unsafe
  render(input: string): Promise<string>;
}
