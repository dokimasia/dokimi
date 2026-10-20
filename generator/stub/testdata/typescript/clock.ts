/** Represents a clock, whose double implements its public instance members. */
//dokimi:stub
export class Clock {
  /** Contains the epoch, which is static, so the double leaves it out. */
  static epoch: number = 0;
  /** Contains the zone, which is a public property. */
  zone: string;

  /**
   * Creates a clock in the zone.
   *
   * @remarks
   * The double leaves the constructor out.
   */
  constructor(zone: string) {
    this.zone = zone;
  }

  /** Returns the instant of the clock. */
  now(): number {
    return Clock.epoch;
  }

  /** Calls fn, a parameter without a sample. */
  every(fn: () => void): void {
    fn();
  }

  /** Returns a function, a result without a sample. */
  stopper(): () => void {
    return () => {};
  }

  /** Returns zero, and is static, so the double leaves it out. */
  static zero(): number {
    return 0;
  }
}
