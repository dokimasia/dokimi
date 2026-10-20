// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for a roled BARE PARAMETER: a method
 * taking a named type directly rather than a request record.
 *
 * @remarks
 * The shape decides where a role can be written. A request record carries
 * the stamp on the field holding the value; a bare parameter has no field, so
 * the only declaration left is the named type itself. Both arms exist in this
 * fixture, so the pool derivation is exercised by the gate rather than by the
 * validated packs alone.
 *
 * @packageDocumentation
 */

/**
 * Identifies one entry.
 *
 * @remarks
 * A named type rather than a bare string, and the role is ON the type: put
 * and get take Key directly, so there is no field to carry the stamp.
 */
//dokimi:default "test-key"
export type Key = string;

/** Represents what an entry contains — the request-record arm, whose role is on the field. */
export interface Payload {
  /** Contains the body of the entry. */
  //dokimi:default "test-body"
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=roledtypestest/
export interface Store {
  /** Writes a payload under a key, which draws both roles. */
  put(key: Key, payload: Payload): Promise<void>;

  /** Returns the payload under the key, which draws the key alone. */
  get(key: Key): Promise<Payload>;
}
