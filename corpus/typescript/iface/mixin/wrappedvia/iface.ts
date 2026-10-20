// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the wrappedvia mixin, which states that
 * the annotated callable delegates to the named sibling, wrapping it with
 * cross-cutting concerns — logging, metrics, retries.
 *
 * @remarks
 * Delegation is a fact about the implementation's insides: both methods sit
 * on one subject, and nothing observable at the interface floor says whether
 * one ran through the other. The stamp documents the relationship and the
 * resolver validates the named sibling; what a delegating call preserves is
 * the model tier's question, asked with instrumentation the floor does not
 * have.
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
//dokimi:out path=wrappedviatest/
export interface Mixed {
  /**
   * Opens the resource of the name through cause, which it wraps.
   *
   * @remarks
   * The fn parameter resolves to a sibling, so the delegate has to be
   * reachable from this interface.
   */
  //dokimi:mixin wrapped-via fn=cause
  open(name: string): Promise<void>;

  /** Does the work that open delegates to, as the delegate that fn names. */
  cause(): Promise<void>;
}
