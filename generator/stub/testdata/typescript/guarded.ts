/** Represents a task with a protected method, so no other class implements it. */
//dokimi:stub
export class Guarded {
  /** Runs the task, and returns the result of check. */
  run(): boolean {
    return this.check();
  }

  /** Reports whether the task may run. */
  protected check(): boolean {
    return true;
  }
}
