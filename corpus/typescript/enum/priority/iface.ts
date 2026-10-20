// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the enum-kind fixture whose zero value is a declared variant.
 *
 * @remarks
 * It differs from the status fixture in exactly one respect, and that is the
 * point: Low is both the zero value and the first declaration, where status
 * starts at 1. The generated zero-value check reads differently in each case
 * — one asserts the zero is a variant, the other that it is not — and a
 * fixture set where every enum makes the same choice cannot tell the two
 * renderings apart.
 *
 * Four variants rather than three, so the count in the exhaustiveness check
 * is not the same number as the neighbouring fixture's. A check that
 * hard-coded one enum's arity would pass against the other by coincidence.
 *
 * No routing directive. The generated surface imports the enum from the
 * module beside it, and the generated checks import the surface from the
 * module beside them, so an `out` sending either elsewhere breaks the import.
 *
 * @packageDocumentation
 */

/** Represents the priority of a task, starting at zero. */
//dokimi:enum
export enum Priority {
  Low,
  Medium,
  High,
  Critical,
}
