// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the hooks mixin, which states that the
 * method calls registered callbacks as part of its work.
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
//dokimi:out path=hookstest/
export interface Mixed {
  /**
   * Calls whatever onEvent registered with the event.
   *
   * @remarks
   * Without a registration method there is no way to observe that it did.
   */
  //dokimi:mixin hooks register=onEvent
  fire(event: string): Promise<void>;

  /** Registers the callback that fire calls. */
  onEvent(fn: (event: string) => void): void;
}
