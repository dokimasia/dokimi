// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the singleflight contract:
 * concurrent calls for one key sharing a single computation.
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

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=singleflighttest/
export interface Contract {
  /**
   * Returns the result of the computation for the key, which concurrent calls
   * for the key share.
   *
   * @remarks
   * run is the singleflight contract's fn role, and hosts the directive that
   * names its partners. It takes the computation it deduplicates — the shape
   * a coalescing claim needs, because a run that computes nothing observable
   * cannot show how often it computed.
   */
  //dokimi:contract single-flight role=fn
  run(key: string, compute: () => string): Promise<string>;

  /**
   * Returns the count of the computations that have ever run.
   *
   * @remarks
   * The count is the observable that a coalescing claim counts. It is also
   * what keeps the sequences driving: the run role takes a callable no pool
   * can draw, so this is the method the twin comparison exercises between law
   * rounds.
   */
  flights(): Promise<number>;
}
