// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture that stacks the tx contract with the retrysucceeds
 * mixin on the commit role.
 *
 * @remarks
 * It earns a composite fixture because retry changes what the contract's
 * terminal-state rule means. Tx says a transaction reaches exactly one
 * terminal state: once committed it cannot be rolled back, and once rolled
 * back it cannot be committed. Retry says a call failing transiently
 * succeeds on a later attempt.
 *
 * Together they ask whether a commit that failed reached a terminal state. If
 * it did, the retry is a second terminal operation and the contract rejects
 * it; if it did not, the retry is the same operation continuing. The answer
 * here is the second — a failed commit leaves the transaction open — and a
 * generator that assumes the first emits a suite failing on the retry it was
 * told to expect.
 *
 * The generator answers it by binding the pairing rather than the mixin: the
 * begin of the same tx contract runs first, and the free-standing retry claim
 * is withheld because that one would be false here.
 *
 * @packageDocumentation
 */

/** Reports a terminal operation on a settled transaction. */
export const closed = new Error("txwithretry: transaction already settled");

/**
 * Reports a commit that may succeed on retry.
 *
 * @remarks
 * It is deliberately distinct from {@link closed}: the difference between
 * them is the whole question this fixture poses.
 */
export const transient = new Error("txwithretry: transient commit failure");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=txwithretrytest/
export interface TxWithRetry {
  /** Opens the transaction, in the begin role of the tx contract. */
  //dokimi:contract tx role=begin
  begin(): Promise<void>;

  /**
   * Applies the transaction, in the commit role of the tx contract, and
   * retries on a transient failure.
   *
   * @remarks
   * A transient failure leaves the transaction open, so the retry is the same
   * terminal operation rather than a second one.
   *
   * attempts= is declared, and which claim it binds is the pairing's answer.
   * Read alone the mixin says any caller may retry until it converges, and a
   * check driving commit with nothing open meets a transaction still idle,
   * reports {@link closed} on every attempt, and reds a subject that honoured
   * both directives. Read beside the begin of the same tx contract, it says
   * the retry converges once the transaction is open — which is what the two
   * of them state and neither states alone.
   */
  //dokimi:contract tx role=commit
  //dokimi:mixin retry-succeeds attempts=2
  commit(): Promise<void>;

  /** Discards the transaction, in the rollback role of the tx contract. */
  //dokimi:contract tx role=rollback
  rollback(): Promise<void>;
}
