// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the enum-kind fixture whose zero value is not a declared variant.
 *
 * @remarks
 * The production surface — statusToString, parseStatus, statusValues,
 * isStatus and the parse error — is generated rather than written here, so
 * the fixture declares only the enum and its members. Writing a
 * statusToString by hand would suppress the generated one, which is why none
 * appears below.
 *
 * Draft is deliberately non-zero. That makes an unset Status invalid rather
 * than silently meaning the first variant, and it is what lets the generated
 * checks assert something about the zero value at all — the priority fixture
 * makes the opposite choice, and the two assertions read differently.
 *
 * The sequence has no gap: a gap would be indistinguishable from an
 * out-of-range value in the fallback check, and the boundary the check probes
 * is the one past the last declared variant.
 *
 * No routing directive. The generated surface imports the enum from the
 * module beside it, and the generated checks import the surface from the
 * module beside them, so an `out` sending either elsewhere breaks the import.
 *
 * @packageDocumentation
 */

/** Represents the status of a document. */
//dokimi:enum
export enum Status {
  Draft = 1,
  Published,
  Archived,
}
