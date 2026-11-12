// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture of inputs whose values contain their own type: a tree
 * of nodes and a state function.
 *
 * @remarks
 * No shape states an interface with methods or a function type, so the setup
 * file of the corpus registers the generators of Node and Step. Each value of
 * them ends at its base values.
 *
 * @packageDocumentation
 */

/** Represents a node of a tree. Its children are nodes. */
export interface Node {
  /** Returns the name of the node. */
  name(): string;
  /** Returns the children of the node. */
  children(): Node[];
}

/** Represents a step of a state machine. It returns the next step, and undefined at the end. */
export type Step = (input: string) => Step | undefined;

/** Represents a counter of the parts of recursive values. */
//dokimi:suite
//dokimi:out path=recursivetest/
export interface Counter {
  /** Counts the nodes of the tree under root. */
  nodes(root: Node): number;
  /** Counts the steps of a run of the machine from start. */
  steps(start: Step): number;
}
