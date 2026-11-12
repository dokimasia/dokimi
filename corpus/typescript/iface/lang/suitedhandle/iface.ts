// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for a declaration under the directive
 * suite alone that a method of a stubbed interface returns.
 *
 * @remarks
 * The declaration has one double, which the module of the producer contains.
 *
 * @packageDocumentation
 */

/** Reports a name that the file system does not have. */
export const notFound = new Error("suitedhandle: file not found");

/**
 * Represents an open file under the directive suite alone. The method open of
 * {@link FS} returns it, so it gets its double as a handle of FS.
 */
//dokimi:suite
//dokimi:out path=suitedhandletest/
export interface File {
  /** Returns the name of the file. */
  name(): string;
  /** Closes the file. */
  close(): void;
}

/** Represents a file system that opens its files. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=suitedhandletest/
export interface FS {
  /**
   * Returns the file with the name.
   *
   * @throws {@link notFound} for a name that the file system does not have.
   */
  open(name: string): File;
}
