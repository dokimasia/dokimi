/**
 * Represents a toggle with the property onRun.
 *
 * @remarks
 * onRun is also the name of the property of the callable of run.
 */
//dokimi:stub
export interface Toggle {
  /** Is true when the toggle runs. */
  onRun: boolean;
  /** Runs the toggle. */
  run(): void;
}
