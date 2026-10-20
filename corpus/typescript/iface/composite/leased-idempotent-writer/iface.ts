// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture that stacks the lease contract with the idempotent
 * mixin on the acquire role.
 *
 * @remarks
 * It earns a composite fixture because the two make opposite demands of the
 * same call. Lease says a second acquire for a held key is turned away with
 * the declared sentinel. Idempotent says repeating the call changes nothing
 * observable. Both are true of a real lease, and they are about different
 * callers — the holder repeats and gets a no-op, a rival is refused.
 *
 * One caller cannot be both, so the lease claim is REFUSED here rather than
 * bound, and the header says which two directives disagree. That is the whole
 * point of the fixture: a generator binding the claim anyway emits a check
 * that reds an implementation which honoured everything the author wrote,
 * and the refusal is what a reader needs in order to move the role to a
 * method that is not idempotent.
 *
 * The directive names taken= exactly so the refusal is the generator's
 * judgement rather than a missing key. Leaving it out would produce the same
 * green suite for the wrong reason.
 *
 * @packageDocumentation
 */

/** Reports an acquire that loses to a holder. */
export const held = new Error("leasedidempotentwriter: already held");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=leasedidempotentwritertest/
export interface LeasedWriter {
  /**
   * Acquires the lease of the key, in the acquire role of the lease contract.
   *
   * @remarks
   * acquire has the mixin idempotent. Re-acquiring a lease this caller already
   * has is a no-op rather than a conflict, which is what makes the two
   * compatible at all.
   */
  //dokimi:contract lease role=acquire taken=held
  //dokimi:mixin idempotent
  acquire(key: string): Promise<void>;

  /** Releases the lease of the key, in the release role of the lease contract. */
  //dokimi:contract lease role=release
  release(key: string): Promise<void>;
}
