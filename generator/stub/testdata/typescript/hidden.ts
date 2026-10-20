/** Represents a task that is not exported, so the module of its double cannot import it. */
//dokimi:stub
interface Hidden {
  /** Runs the task. */
  run(): void;
}

/** Contains a task of the type Hidden, so the module does not declare an unused type. */
export const hidden: Hidden = { run: () => {} };
