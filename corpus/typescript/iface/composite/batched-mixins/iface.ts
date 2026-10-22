// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture for several mixins on one line.
 *
 * @remarks
 * A comment line can contain several directives. Each directive after the
 * first opens with the mark dokimi: after a blank. A param belongs to the
 * directive that it follows. put and list declare three mixins on one line
 * each.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=batchedmixinstest/
export interface Batched {
  /**
   * Stores the value under the key.
   *
   * @remarks
   * put has three mixins. side-effect observes the write through list, which
   * returns a value before the write and after it. read refuses a key that no
   * write stored, so it cannot observe the change.
   */
  //dokimi:mixin idempotent dokimi:mixin concurrent dokimi:mixin side-effect observe=list
  put(key: string, value: string): Promise<void>;

  /** Returns the value under the key, with one mixin with a param. */
  //dokimi:mixin read-after-write write=put
  read(key: string): Promise<string>;

  /** Returns the values, with two mixins without params and one mixin with a param. */
  //dokimi:mixin cacheable dokimi:mixin pure dokimi:mixin bounded limit=50
  list(): Promise<string[]>;
}
