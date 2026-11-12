// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Clock}.
 *
 * @packageDocumentation
 */

import type { Clock as WallClock } from "@dokimi/assert";
import { type Clock, type Instant, offset } from "../iface";
import type { ClockHandle } from "./iface.suite";

/** Contains the max offset of the wall time of a peer beyond the wall time of the clock, in nanoseconds. */
const maxOffset = 500_000_000n;

/** Returns the greatest of the wall times. */
function latest(...walls: bigint[]): bigint {
  return walls.reduce((a, b) => (a > b ? a : b));
}

/**
 * Implements {@link Clock} with the update rule of a hybrid logical clock, on
 * a clock of assert. The corpus applies the suite of Clock to it as a
 * subject.
 */
//dokimi:subject Clock new=newInMemoryFromHandle
export class InMemory implements Clock {
  readonly #wall: WallClock;
  #last: Instant = { wall: 0n, logical: 0 };

  /** Returns a hybrid logical clock that reads the wall time of wall. */
  constructor(wall: WallClock) {
    this.#wall = wall;
  }

  /**
   * Merges the instant of a peer into the clock, and returns the instant of
   * the clock after the merge.
   *
   * @throws {@link offset} for an instant beyond the max offset of the clock.
   */
  update(observed: Instant): Instant {
    const now = BigInt(this.#wall.now()) * 1_000_000n;
    if (observed.wall > now + maxOffset) {
      throw new Error(`refusedtest: wall time ${observed.wall}`, { cause: offset });
    }
    const wall = latest(now, this.#last.wall, observed.wall);
    let logical = 0;
    if (wall === this.#last.wall) {
      logical = this.#last.logical + 1;
      if (wall === observed.wall) {
        logical = Math.max(this.#last.logical, observed.logical) + 1;
      }
    } else if (wall === observed.wall) {
      logical = observed.logical + 1;
    }
    this.#last = { wall, logical };
    return this.#last;
  }
}

/** Returns a new {@link InMemory} on the clock of the handle. */
export function newInMemoryFromHandle(h: ClockHandle): InMemory {
  return new InMemory(h.clock());
}
