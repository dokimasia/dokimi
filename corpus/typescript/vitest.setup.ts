// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Registers the assertions of `@dokimi/assert` on the `expect` of vitest, so
 * the generated tests call them as `expect(got).isEqualTo(want, message)`.
 * It also registers the generator of each name that the shape definitions of
 * the generated suites refer to and no shape states.
 */

import { prop } from "@dokimi/assert";
import { register } from "@dokimi/assert/vitest";
import type { Node, Step } from "./iface/lang/recursive/iface";

register();

/**
 * The generator of the values of the type unknown, the witness of a type
 * parameter without a bound: a text, a number or a boolean.
 */
prop.register("unknown", prop.oneOf<unknown>(prop.string(), prop.float(-1e6, 1e6), prop.boolean()));

/** Returns a node with the name and the children. */
function nodeOf(name: string, children: Node[]): Node {
  return { name: () => name, children: () => children };
}

/** The generator of the trees of the fixture recursive: a leaf, or a node of up to three trees. */
prop.register(
  "Node",
  prop.recursive<Node>(
    prop.string().map((name) => nodeOf(name, [])),
    (tree) =>
      prop.composite((c) =>
        nodeOf(c.draw(prop.string(), "name"), c.draw(prop.list(tree, { maxSize: 3 }), "children")),
      ),
    { maxLeaves: 8 },
  ),
);

/** Returns a machine that takes n steps and then ends. */
function machineOf(n: number): Step {
  return () => (n > 1 ? machineOf(n - 1) : undefined);
}

/** The generator of the state machines of the fixture recursive: a machine of one to five steps. */
prop.register(
  "Step",
  prop.integer(1, 5).map((n) => machineOf(n)),
);
