/** Reports a commit before prepare. */
export const notReady = new Error("not ready");

/** Represents a session that commits after it prepares, and keeps an old method. */
//dokimi:stub
export interface Session {
  /** Prepares the commit. */
  prepare(): Promise<void>;
  /** Commits the session, which is valid after prepare alone. */
  //dokimi:mixin order-after fn=prepare unready=notReady
  commit(): Promise<void>;
  /** Returns the value of the key, and is scheduled for removal. */
  //dokimi:mixin deprecated
  old(key: string): string;
}
