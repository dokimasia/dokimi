// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the prefix=off half of the sentinel directive.
 *
 * @remarks
 * The directive has two behaviours behind one name: prefix=<value> overrides
 * the prefix derived from the module name, and prefix=off suppresses the
 * prefix check entirely. A corpus covering only the override branch is green
 * with the suppression branch never generated, because the directive name
 * appears either way.
 *
 * The errors here deliberately carry no module prefix. That is what makes
 * the fixture able to detect a failure rather than merely exercise a path:
 * with the check suppressed they are legal, and a generator that parsed
 * prefix=off and ignored it would emit the assertion and fail against them.
 *
 * @packageDocumentation
 */
//dokimi:sentinel prefix=off

/** Reports an operation that the caller should retry, without a prefix. */
export class RetryError extends Error {
  constructor() {
    super("retry the operation");
    this.name = "RetryError";
  }
}

/** Reports an operation without attempts left, without a prefix. */
export class GiveUpError extends Error {
  constructor() {
    super("no attempts remain");
    this.name = "GiveUpError";
  }
}
