// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for a relational mixin whose partner
 * spells the shared parameter differently from the method that names it.
 *
 * @remarks
 * Every other relational fixture in the corpus spells both halves alike, so
 * the derivation that matches them had one arm exercised and two dead. That
 * is the shape this fixture exists to vary: the classification, the types and
 * the arity are all ordinary, and only the identifiers differ.
 *
 * The methods are the four answers the derivation can give, one each:
 *
 * - touch/seen — one candidate of the type, so the correspondence is
 *   derivable and the check is generated.
 *
 * - move/at — two candidates of the type and nothing saying which, which is
 *   the ambiguity `partition` settles with `axis=`. Declined, and the run
 *   says so.
 *
 * - emit/count — no candidate at all, so a check would have to invent a
 *   value the method never takes. Declined, and for a different reason.
 *
 * - version/put/replace — two writes take the record whose revision version
 *   reads, so the partner write of monotonic-reads has two candidates.
 *   Declined with one-partner.
 *
 * This axis varies naming rather than the classification, so it breaks
 * generators independently of any directive — which is why it sits here
 * beside receivercollision rather than under `mixin/sideeffect`. That fixture
 * is about what the mixin claims; this one is about whether the generator
 * can write the claim down.
 *
 * @packageDocumentation
 */

/** Represents a record of {@link Store} with the revision that its last write stamped. */
export interface Rec {
  /** Contains the key of the record. */
  key: string;
  /** Contains the revision that the last write of the record stamped. */
  rev: number;
}

/** Represents a store with one relational pair per outcome of the partner-argument derivation. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=partnernamingtest/
export interface Store {
  /**
   * Marks the key with the weight.
   *
   * @remarks
   * touch mutates state the return value does not carry, and seen is what
   * makes the effect visible. The observer spells the key `k` where this
   * spells it `key`, and takes nothing of the other parameter's type — so
   * exactly one of touch's parameters can be it and no guess is involved.
   */
  //dokimi:mixin side-effect observe=seen
  touch(key: string, weight: number, signal: AbortSignal): Promise<void>;

  /** Returns what touch left behind under the key. */
  seen(k: string, signal: AbortSignal): Promise<number>;

  /**
   * Moves the state from one identifier to another of the same type.
   *
   * @remarks
   * at asks about one of them without saying which. Nothing in the source
   * decides it: `from` and `to` are both strings, and a generator picking one
   * would write a check about the wrong end of the move half the time. The
   * corresponding fix is to spell the observer's parameter as one of them.
   */
  //dokimi:mixin side-effect observe=at
  move(from: string, to: string, signal: AbortSignal): Promise<void>;

  /** Returns what is at one end of a move. */
  at(where: string, signal: AbortSignal): Promise<number>;

  /**
   * Records an event under the id.
   *
   * @remarks
   * count asks about a bucket the method never takes. Not an ambiguity but an
   * absence: a check receives emit's arguments and nothing else, so there is
   * no number to hand count. The fix is a fixture field, not a rename.
   */
  //dokimi:mixin side-effect observe=count
  emit(id: string, signal: AbortSignal): Promise<void>;

  /** Returns how many records the bucket contains. */
  count(bucket: number, signal: AbortSignal): Promise<number>;

  /**
   * Returns the record under the key, whose revision never goes back.
   *
   * @remarks
   * put and replace both take a Rec, so the derivation finds two writes that
   * the reads could follow, and refuses the claim with one-partner.
   */
  //dokimi:mixin monotonic-reads version=rev
  version(key: string, signal: AbortSignal): Promise<Rec>;

  /** Stores a record and stamps its revision. */
  put(r: Rec, signal: AbortSignal): Promise<void>;

  /** Stores a record in place of the record under its key and stamps its revision. */
  replace(r: Rec, signal: AbortSignal): Promise<void>;
}
