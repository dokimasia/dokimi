/** Represents a clock with a method without a declared result type. */
//dokimi:stub
export class Undeclared {
  /** Returns a number without a type annotation of its result. */
  now() {
    return 0;
  }
}

/** Represents a counter with a property without a declared type. */
//dokimi:stub
export class Counted {
  /** Contains the count, without a declared type. */
  count = 0;
}

/** Represents a task with a parameter without a declared type. */
//dokimi:stub
export interface Loose {
  /** Runs the task with a parameter without a type annotation. */
  run(x): void;
}
