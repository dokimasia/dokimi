// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the tamperevident mixin, which states
 * that a stored value altered outside the interface is detected rather than
 * served as if intact.
 *
 * @remarks
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=tamperevidenttest/
export interface Mixed {
  /** Stores the body under its integrity check. */
  //dokimi:mixin tamper-evident tamper=corrupt verify=verify
  store(body: string): Promise<void>;

  /**
   * Alters the stored bytes behind the back of the interface.
   *
   * @remarks
   * corrupt is the only way to reach the state that verify exists to detect.
   */
  corrupt(): Promise<void>;

  /** Rejects when what is stored is not what was stored. */
  verify(): Promise<void>;
}
