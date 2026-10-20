// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for a lifetime carried on the value rather
 * than fixed by the directive.
 *
 * @remarks
 * The sibling fixture `ttl` stamps `duration=1m`: one lifetime for every
 * write, declared once. This one stamps no duration at all, and each write
 * says how long its own entry lasts. Both are `ttl`, and the claim is the
 * same sentence — an entry stops being readable once its lifetime has run
 * out — but only this shape has a field a defect can reach for, which is what
 * lets the claim be proven rather than argued.
 *
 * lifetime= names the member, so the claim reads the number off the write's
 * answer instead of the directive.
 *
 * @packageDocumentation
 */

/**
 * Reports a read past the lifetime.
 *
 * @remarks
 * The not-found= param names it, so the law compares the identity the
 * declaration gave rather than merely noticing that something failed.
 */
export const expired = new Error("ttlperwrite: entry expired");

/**
 * Represents the payload, with its own lifetime.
 *
 * @remarks
 * Nothing declares what a drawn Entry looks like. A drawn value carries every
 * member, the lifetime included.
 */
export interface Entry {
  /** Identifies the entry. */
  key: string;

  /** Contains what a read returns. */
  body: string;

  /**
   * Contains how long the entry stays readable, in milliseconds.
   *
   * @remarks
   * Zero means forever, which is what makes stripping it a defect rather than
   * a refusal: the write still lands and the read still answers, and the only
   * thing that changed is that nothing ever expires.
   */
  lifetime: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=ttlperwritetest/
export interface Mixed {
  /** Stores an entry, which starts the lifetime that it carries. */
  put(e: Entry): Promise<void>;

  /**
   * Returns the entry under the key while its lifetime lasts, and rejects
   * with the sentinel after.
   *
   * @remarks
   * lifetime= rather than duration=: the entry says how long it lasts, so the
   * number does not exist until the write has answered.
   */
  //dokimi:mixin ttl put=put read=read lifetime=lifetime not-found=expired
  read(key: string): Promise<Entry>;
}
