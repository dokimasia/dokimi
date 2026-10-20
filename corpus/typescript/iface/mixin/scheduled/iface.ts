// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the scheduled mixin, which states that
 * a task registered for a future instant has run once the clock passes it.
 *
 * @remarks
 * The directive names both halves, and the fired-count accessor is the
 * load-bearing one: a run that cannot count firings reports every scheduler
 * as correct, including one that fires nothing at all.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Represents the subject of the fixture.
 *
 * @remarks
 * The model tier runs clocked: both methods are the scheduled mixin's own
 * partners, and the consumer supplies a clocked model so the property owns a
 * test clock the scheduled-fires law can advance. Unclocked, no sequence
 * would ever move time and every firing claim would hold vacuously.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=scheduledtest/
export interface Mixed {
  /**
   * Registers a task for the given offset from now.
   *
   * @param after - The offset from now, in milliseconds.
   */
  //dokimi:mixin scheduled schedule=at runs=fired
  at(after: number): Promise<void>;

  /** Returns the count of the tasks that have run, which makes the claim checkable. */
  fired(): Promise<number>;
}
