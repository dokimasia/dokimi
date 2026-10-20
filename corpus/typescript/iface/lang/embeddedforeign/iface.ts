// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for an interface that extends one from
 * outside the tree.
 *
 * @remarks
 * Extending an interface of the TypeScript library is ordinary code rather
 * than an exotic case. A double missing a member does not satisfy the
 * interface it doubles, and a harness over part of a contract passes an
 * implementation that fails the rest. So both generators decline an
 * interface whose member set they cannot complete.
 *
 * What this fixture proves is that the completion reaches the output: name,
 * message, stack and cause arrive from Error with no declaration in this
 * module, and the double, the harness and the harness's own checks all carry
 * them beside retryable.
 *
 * The embedded fixture covers the same shape with every base declared
 * locally. The two must agree, which is the point of having both.
 *
 * @packageDocumentation
 */

/**
 * Represents a failure that extends Error of the TypeScript library beside a
 * method of its own.
 *
 * @remarks
 * name, message, stack and cause come from Error. retryable is declared here.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=embeddedforeigntest/
export interface Failure extends Error {
  /** Reports whether a retry of the failed call can succeed. */
  retryable(): boolean;
}
