import type { Clock, Seat } from "@dokimi/assert";
import type { RandomSource } from "@dokimi/runtime/fault";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import type { Box, Point } from "./point";

/** Configures a PointStub at its construction. */
export type PointStubOption = (double: PointStub) => void;

/**
 * PointStub is a recording test double of Point. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class PointStub implements Point {
  /** The property x, which a test assigns. */
  x!: number;
  /** The property y, which a test assigns. */
  y!: number;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Point that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: PointStubOption[]) {
    this.#group = new Group();
    for (const option of options) {
      option(this);
    }
    if (this.#strict) {
      this.#group.strictAll();
    }
    this.#group.bind(seat);
  }

  /**
   * Returns the option that makes an unconfigured method fail the test.
   *
   * @returns The option.
   */
  static strict(): PointStubOption {
    return (double) => {
      double.#strict = true;
    };
  }

  /**
   * Returns the option that sets the clock of the double to clock. The double
   * takes the latency, the windows of the faults and the instants of the calls
   * from clock.
   *
   * @param clock - The clock.
   * @returns The option.
   */
  static withClock(clock: Clock): PointStubOption {
    return (double) => {
      double.#group.setClock(clock);
    };
  }

  /**
   * Returns the option that sets the random source of the probabilistic faults
   * to random. A seeded source repeats the faults of a run.
   *
   * @param random - The source.
   * @returns The option.
   */
  static withRandSource(random: RandomSource): PointStubOption {
    return (double) => {
      double.#group.setRandSource(random);
    };
  }

  /**
   * Returns the option that turns the recording of calls off. The methods
   * still return their results.
   *
   * @returns The option.
   */
  static benchMode(): PointStubOption {
    return (double) => {
      double.#group.benchMode();
    };
  }

  /**
   * Returns the option that makes every asynchronous method wait duration on
   * the clock of the double before it returns.
   *
   * @param duration - The latency, in milliseconds.
   * @returns The option.
   */
  static withLatency(duration: number): PointStubOption {
    return (double) => {
      double.#group.setLatency(duration);
    };
  }

  /**
   * Returns a double whose asynchronous methods wait duration on clock before
   * they return. It applies options after the clock and the latency.
   *
   * @param seat - The seat of the test.
   * @param duration - The latency, in milliseconds.
   * @param clock - The clock.
   * @param options - The options.
   * @returns The double.
   */
  static newDelayed(
    seat: Seat | undefined,
    duration: number,
    clock: Clock,
    ...options: PointStubOption[]
  ): PointStub {
    return new PointStub(
      seat,
      PointStub.withClock(clock),
      PointStub.withLatency(duration),
      ...options,
    );
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }
}

/** Configures a BoxStub at its construction. */
export type BoxStubOption<T> = (double: BoxStub<T>) => void;

/**
 * BoxStub is a recording test double of Box. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class BoxStub<T> implements Box<T> {
  /** The property value, which a test assigns. */
  value!: T;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Box that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: BoxStubOption<T>[]) {
    this.#group = new Group();
    for (const option of options) {
      option(this);
    }
    if (this.#strict) {
      this.#group.strictAll();
    }
    this.#group.bind(seat);
  }

  /**
   * Returns the option that makes an unconfigured method fail the test.
   *
   * @returns The option.
   */
  static strict<T>(): BoxStubOption<T> {
    return (double) => {
      double.#strict = true;
    };
  }

  /**
   * Returns the option that sets the clock of the double to clock. The double
   * takes the latency, the windows of the faults and the instants of the calls
   * from clock.
   *
   * @param clock - The clock.
   * @returns The option.
   */
  static withClock<T>(clock: Clock): BoxStubOption<T> {
    return (double) => {
      double.#group.setClock(clock);
    };
  }

  /**
   * Returns the option that sets the random source of the probabilistic faults
   * to random. A seeded source repeats the faults of a run.
   *
   * @param random - The source.
   * @returns The option.
   */
  static withRandSource<T>(random: RandomSource): BoxStubOption<T> {
    return (double) => {
      double.#group.setRandSource(random);
    };
  }

  /**
   * Returns the option that turns the recording of calls off. The methods
   * still return their results.
   *
   * @returns The option.
   */
  static benchMode<T>(): BoxStubOption<T> {
    return (double) => {
      double.#group.benchMode();
    };
  }

  /**
   * Returns the option that makes every asynchronous method wait duration on
   * the clock of the double before it returns.
   *
   * @param duration - The latency, in milliseconds.
   * @returns The option.
   */
  static withLatency<T>(duration: number): BoxStubOption<T> {
    return (double) => {
      double.#group.setLatency(duration);
    };
  }

  /**
   * Returns a double whose asynchronous methods wait duration on clock before
   * they return. It applies options after the clock and the latency.
   *
   * @param seat - The seat of the test.
   * @param duration - The latency, in milliseconds.
   * @param clock - The clock.
   * @param options - The options.
   * @returns The double.
   */
  static newDelayed<T>(
    seat: Seat | undefined,
    duration: number,
    clock: Clock,
    ...options: BoxStubOption<T>[]
  ): BoxStub<T> {
    return new BoxStub<T>(
      seat,
      BoxStub.withClock<T>(clock),
      BoxStub.withLatency<T>(duration),
      ...options,
    );
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }
}

void compatV1;
