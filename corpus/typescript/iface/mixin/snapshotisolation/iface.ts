// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the snapshotisolation mixin, which
 * states that concurrent transactions are free of the read and write
 * anomalies a serial order would prevent.
 *
 * @remarks
 * The interface reports its own transaction history, because that is the
 * only thing an isolation claim can be checked against: the anomalies are
 * patterns across transactions, not properties of any one call.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents one recorded operation, with the transaction that made it. */
export interface Entry {
  /** Contains the transaction that made the operation. */
  txn: number;
  /** Contains the key that the operation touched. */
  key: string;
  /** Contains the version that the operation read or wrote. */
  version: number;
  /** Is true when the operation is a write. */
  write: boolean;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=snapshotisolationtest/
export interface Mixed {
  /** Appends one operation to the history. */
  record(e: Entry): Promise<void>;

  /** Returns what the transactions did, which is what the anomaly checks read. */
  //dokimi:mixin snapshot-isolation
  history(): Promise<Entry[]>;

  /**
   * Returns the latest recorded entry of the key.
   *
   * @remarks
   * get is the read half a store under isolation genuinely has, and the draw
   * that gives the anomaly laws a key type to instantiate at.
   */
  get(key: string): Promise<Entry>;
}
