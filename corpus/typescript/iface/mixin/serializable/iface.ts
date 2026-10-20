// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the serializable mixin, which states
 * that concurrent transactions are equivalent to some serial order — so no
 * anomaly is observable at all.
 *
 * @remarks
 * The sibling of snapshotisolation, and the fixture pair states why they are
 * siblings rather than levels. Snapshot isolation forbids dirty writes and
 * dirty, intermediate and circular reads, and *permits* write skew: two
 * transactions each reading what the other is about to invalidate, both
 * committing. Serializability forbids that too. So the anti-dependency-cycle
 * check is correct against this claim and wrong against a snapshot, and a
 * store that declares only the snapshot must not earn it.
 *
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
//dokimi:out path=serializabletest/
export interface Mixed {
  /** Appends one operation to the history. */
  record(e: Entry): Promise<void>;

  /** Returns what the transactions did, which is what the anomaly checks read. */
  //dokimi:mixin serializable
  history(): Promise<Entry[]>;

  /**
   * Returns the latest recorded entry of the key.
   *
   * @remarks
   * get is the read half a serializable store genuinely has, and the draw
   * that gives the anomaly law a key type to instantiate at.
   */
  get(key: string): Promise<Entry>;
}
