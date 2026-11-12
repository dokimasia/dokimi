// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for inputs whose fields have function
 * types and interface types, at each depth of a request.
 *
 * @packageDocumentation
 */

/** Reports a dispatch whose request did not take its parcel, because done or notify of the request threw. */
export const undelivered = new Error("nested: the request did not take the parcel");

/** Represents the parcel of a dispatch. */
export interface Parcel {
  /** Returns the weight of the parcel in grams. */
  weight(): number;
}

/** Represents the notification of the sender of the parcel with the identity id. */
export type Notify = (id: string) => void;

/** Represents the request of a dispatch. parcel has an interface type, done a function type and notify an alias of a function type. */
export interface DispatchRequest {
  /** Contains the parcel of the dispatch. */
  parcel: Parcel;
  /** Contains the function that the dispatch calls once the parcel left. */
  done: () => void;
  /** Contains the function that notifies the sender. */
  notify: Notify;
  /** Contains the identity of the request. */
  id: string;
}

/** Represents a batch of dispatches. */
export interface Batch {
  /** Contains the requests of the batch. */
  requests: DispatchRequest[];
}

/** Represents a dispatcher of the parcels of requests. */
//dokimi:suite
//dokimi:out path=nestedtest/
export interface Dispatcher {
  /**
   * Dispatches the parcel of the request.
   *
   * @throws {@link undelivered} where the request does not take the parcel.
   */
  //dokimi:budget allocs=0
  dispatch(req: DispatchRequest, signal: AbortSignal): Promise<void>;

  /**
   * Dispatches the parcel of each request of the batch, and stops at the
   * first request that does not take its parcel.
   *
   * @throws {@link undelivered} where a request does not take its parcel.
   */
  //dokimi:budget allocs=0
  dispatchAll(b: Batch, signal: AbortSignal): Promise<void>;

  /** Returns the total weight of the dispatched parcels in grams. */
  weight(): number;
}
