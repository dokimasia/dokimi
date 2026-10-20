// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the lease contract: an acquire
 * balanced by exactly one release.
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
 * Reports an acquire of a lease that someone already has.
 *
 * @remarks
 * The error is the contract's own sentinel, declared beside the shape so the
 * directive can name it.
 */
export const held = new Error("lease: the lease is already held");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=leasetest/
export interface Contract {
  /** Acquires the lease of the key, in the acquire role of the lease contract. */
  //dokimi:contract lease role=acquire timeout=100ms taken=held
  acquire(key: string): Promise<void>;

  /** Releases the lease of the key, in the release role of the lease contract. */
  //dokimi:contract lease role=release
  release(key: string): Promise<void>;
}
