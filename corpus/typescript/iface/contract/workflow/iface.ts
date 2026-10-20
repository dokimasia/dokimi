// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the workflow contract: a callable
 * driving a state machine.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * The states of the machine are plain strings: an enum here would be an enum
 * node, and neither `default` nor `sample` applies to one — which would leave
 * the zero-on-error claim of state with no value it can tell the zero apart
 * from. The empty state is where a key sits before anything ran. It is a real
 * state the graph names, not a stand-in for absence.
 *
 * @packageDocumentation
 */

/** Represents the payload that the roles of the contract exchange. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/** Contains the state of a key before anything ran. */
export const stateNew = "";

/** Contains the state of a draft. */
export const stateDraft = "Draft";

/** Contains the state of a live key. */
export const stateLive = "Live";

/**
 * Maps each state to the states that run may move a key to.
 *
 * @remarks
 * A declaration rather than a string the directive encodes: the graph is a
 * set of pairs, and a parameter names one thing. What can be named as one
 * thing is the constant that contains it, which leaves the shape here — a map
 * from a state to the states reachable from it — the fixture's own rather
 * than a notation every consumer would have to write in.
 */
export const transitions = new Map<string, string[]>([
  [stateNew, [stateDraft]],
  [stateDraft, [stateLive]],
]);

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=workflowtest/
export interface Contract {
  /**
   * Moves the key to its next state.
   *
   * @remarks
   * run is the workflow contract's fn role, and hosts the directive that
   * names the graph and the observation.
   */
  //dokimi:contract workflow role=fn transitions=transitions observe=state
  run(key: string): Promise<void>;

  /**
   * Returns the position of the key in the machine.
   *
   * @remarks
   * The position is the observation that observe= names, which a check
   * compares across a run.
   */
  state(key: string): Promise<string>;
}
