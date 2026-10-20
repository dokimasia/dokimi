// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the zero-value half of the builder's defaults handling: no
 * defaults companion, no field directives, nothing to seed from.
 *
 * @remarks
 * It is the floor the other two fixtures are measured against. A generator
 * that always emits a defaults-aware constructor produces identical output
 * here and for the directive-bearing case, and only a fixture with no
 * defaults at all can tell those apart.
 *
 * Routing is declared once for the module rather than repeated on each
 * class: every builder in a module lands in the same companion directory, so
 * a per-class directive is the same statement written N times, and the Nth
 * copy is the one that gets forgotten.
 *
 * @packageDocumentation
 */
//dokimi:out path=plaintest/

/**
 * Represents an item without defaults of any kind.
 *
 * @remarks
 * newItemBuilder() must return a builder over the value of `new Item()`.
 */
//dokimi:builder
export class Item {
  /** Contains the identifier of the item. */
  id: string = "";
  /** Counts the units of the item. */
  count: number = 0;
  /** Lists the tags of the item. */
  tags: string[] = [];
}
