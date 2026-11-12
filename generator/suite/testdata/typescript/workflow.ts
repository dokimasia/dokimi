/** Contains the state of a key before anything ran. */
export const stateNew = "";

/** Contains the state of a draft. */
export const stateDraft = "Draft";

/** Maps each state to the states that run may move a key to. */
export const transitions = new Map<string, string[]>([[stateNew, [stateDraft]]]);

/** Lists the states in the order that restart moves a key through, which is no graph. */
export const order: readonly string[] = [stateNew, stateDraft];

/** Represents a workflow over the graph of its transitions. */
//dokimi:suite
export interface Flow {
  /** Moves the key to its next state. */
  //dokimi:contract workflow role=fn transitions=transitions observe=state
  run(key: string, signal: AbortSignal): Promise<void>;
  /** Moves the key to its first state. */
  //dokimi:contract workflow role=fn transitions=order observe=state
  restart(key: string, signal: AbortSignal): Promise<void>;
  /** Returns the state of the key. */
  state(key: string, signal: AbortSignal): Promise<string>;
}
