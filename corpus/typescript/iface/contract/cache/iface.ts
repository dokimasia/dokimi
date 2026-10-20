// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the cache contract: a cache in front
 * of a backing store.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * The model tier rides the twin floor, and what a store would need is the
 * loader fetch calls on a miss: the cache's state is whatever that returned,
 * and no stamp names it.
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
//dokimi:out path=cachetest/
export interface Contract {
  /**
   * Stores the value under the key in the backing store.
   *
   * @remarks
   * store is on the interface for the reason tx's put is: a contract stating
   * that a cached read agrees with its backing states something about a store
   * with something IN it, and a state no conformance check can establish
   * through the interface is one no check can reach. Held off it, the
   * agreement bound against two misses — which agree by both finding nothing.
   */
  store(key: string, v: Value): Promise<void>;

  /** Returns the value under the key, in the cache role of the cache contract. */
  //dokimi:contract cache role=cache
  lookup(key: string): Promise<Value>;

  /** Returns the value under the key, in the backing role of the cache contract. */
  //dokimi:contract cache role=backing
  fetch(key: string): Promise<Value>;
}
