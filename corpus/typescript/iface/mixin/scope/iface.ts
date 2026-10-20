// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the scope mixin, which states that the
 * effect is confined to a named scope.
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
//dokimi:out path=scopetest/
export interface Mixed {
  /**
   * Stores the value under the key, confined to the named scope.
   *
   * @remarks
   * The law reads back from another scope to prove containment, so a reader
   * is required.
   *
   * The naming form beside the checkable one, which is the composition the
   * ruling describes: scope says what the axis MEANS and licenses no check,
   * partition names the observer as well and derives the isolation. Two
   * stamps rather than one mixin doing both, the way idempotent and
   * accumulates sit on the effect axis.
   */
  //dokimi:mixin scope name=tenant axis=scope
  //dokimi:mixin partition axis=scope read=get
  set(scope: string, key: string, value: string): Promise<void>;

  /** Returns the value under the key within a scope, which makes confinement observable. */
  get(scope: string, key: string): Promise<string>;
}
