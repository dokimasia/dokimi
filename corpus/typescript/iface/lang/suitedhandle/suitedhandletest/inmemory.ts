// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subjects of the suites of {@link File} and {@link FS}.
 *
 * @packageDocumentation
 */

import { type File, type FS, notFound } from "../iface";

/**
 * Implements {@link File} in memory. The corpus applies the suite of File to
 * it as a subject.
 */
//dokimi:subject File new=newInMemoryFile
export class InMemoryFile implements File {
  readonly #name: string;

  /** Returns a file with the name. */
  constructor(name: string) {
    this.#name = name;
  }

  /** Returns the name of the file. */
  name(): string {
    return this.#name;
  }

  /** Closes the file. A second call returns as well. */
  close(): void {}
}

/** Returns a new {@link InMemoryFile} without a name. */
export function newInMemoryFile(): InMemoryFile {
  return new InMemoryFile("");
}

/**
 * Implements {@link FS} in memory. The corpus applies the suite of FS to it as
 * a subject.
 */
//dokimi:subject FS new=newInMemoryFS
export class InMemoryFS implements FS {
  readonly #files = new Map<string, InMemoryFile>();

  /**
   * Returns the file with the name.
   *
   * @throws {@link notFound} for a name that the file system does not have.
   */
  open(name: string): File {
    const file = this.#files.get(name);
    if (file === undefined) {
      throw notFound;
    }
    return file;
  }
}

/** Returns a new {@link InMemoryFS} without files. */
export function newInMemoryFS(): InMemoryFS {
  return new InMemoryFS();
}
