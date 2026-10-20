// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the pool contract: a pool whose
 * every get is returned by a put.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/** Represents the payload that the roles of the contract exchange. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/**
 * Represents the accounting that the pool laws balance.
 *
 * @remarks
 * The stats role of the contract returns it, so the generator knows where the
 * numbers are.
 */
export interface Stats {
  /** Counts the gets. */
  gets: number;
  /** Counts the puts. */
  puts: number;
  /** Counts the values that a get took and no put returned. */
  outstanding: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=pooltest/
export interface Contract {
  /** Takes a value from the pool, in the get role of the pool contract. */
  //dokimi:contract pool role=get
  get(): Promise<Value>;

  /** Returns a value to the pool, in the put role of the pool contract. */
  //dokimi:contract pool role=put
  put(v: Value): Promise<void>;

  /**
   * Returns the accounting of the pool, in the stats role of the pool
   * contract.
   *
   * @remarks
   * The balanced and leak-free laws verify the accounting against the cycle
   * count.
   */
  //dokimi:contract pool role=stats
  stats(): Promise<Stats>;
}
