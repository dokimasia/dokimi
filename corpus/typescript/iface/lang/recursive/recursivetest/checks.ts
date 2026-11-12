// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the checks that a consumer adds to the suite of {@link Counter}.
 *
 * @packageDocumentation
 */

import { check, prop } from "@dokimi/assert";
import type { Counter, Node, Step } from "../iface";
import type { CounterHandle } from "./iface.suite";

/** Checks that nodes counts the root of a tree that the case draws. */
//dokimi:check Counter.nodes defect=countsNothing
export async function countsTheRoot(h: CounterHandle, c: prop.Case): Promise<void> {
  const root = c.draw(prop.of<Node>("Node"), "root");
  check.isTrue(h, h.subject().nodes(root) >= 1, h.statement());
}

/** Checks that steps counts the first step of a machine that the case draws. */
//dokimi:check Counter.steps defect=countsNothing
export async function takesTheFirstStep(h: CounterHandle, c: prop.Case): Promise<void> {
  const start = c.draw(prop.of<Step>("Step"), "start");
  check.isTrue(h, h.subject().steps(start) >= 1, h.statement());
}

/** Returns the planted defect of both checks, a counter that counts no node and no step. */
export function countsNothing(_h: CounterHandle): Counter {
  return { nodes: () => 0, steps: () => 0 };
}
