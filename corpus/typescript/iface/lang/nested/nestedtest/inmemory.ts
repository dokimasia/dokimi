// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Dispatcher}.
 *
 * @packageDocumentation
 */

import { type Batch, type Dispatcher, type DispatchRequest, undelivered } from "../iface";

/**
 * Implements {@link Dispatcher} in memory. Each dispatch weighs the parcel of
 * its request and calls done and notify of the request, so it reads each
 * field of a function type and of an interface type. The corpus applies the
 * suite of Dispatcher to it as a subject.
 */
//dokimi:subject Dispatcher new=newInMemory
export class InMemory implements Dispatcher {
  #weight = 0;

  /**
   * Dispatches the parcel of the request.
   *
   * @throws {@link undelivered} where the request does not take the parcel.
   */
  async dispatch(req: DispatchRequest, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#dispatch(req);
  }

  /**
   * Dispatches the parcel of each request of the batch, and stops at the
   * first request that does not take its parcel.
   *
   * @throws {@link undelivered} where a request does not take its parcel.
   */
  async dispatchAll(b: Batch, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    for (const req of b.requests) {
      this.#dispatch(req);
    }
  }

  /** Returns the total weight of the dispatched parcels in grams. */
  weight(): number {
    return this.#weight;
  }

  /** Weighs the parcel of the request, and calls done and notify of the request. */
  #dispatch(req: DispatchRequest): void {
    this.#weight += req.parcel.weight();
    try {
      req.done();
      req.notify(req.id);
    } catch (cause) {
      throw new Error("nestedtest: the request did not take the parcel", { cause: [undelivered, cause] });
    }
  }
}

/** Returns an {@link InMemory} without dispatches. */
export function newInMemory(): InMemory {
  return new InMemory();
}
