import type { Clock, Seat } from '@dokimi/assert';
import type { RandomSource } from '@dokimi/runtime/fault';
import { compatV1 } from '@dokimi/runtime/stub/compat';
import { Group } from '@dokimi/runtime/stub/group';
import type { Box, Point } from './point';

/** An option of a PointStub, which its constructor applies. */
export type PointStubOption = (double: PointStub) => void;

/**
 * PointStub is a recording test double of Point. Each property
 * on<Method> configures its method. Left alone, a method returns undefined
 * and records the call.
 */
export class PointStub implements Point {
  /** The property x, which a test assigns. */
  x!: number;
  /** The property y, which a test assigns. */
  y!: number;
  readonly #group: Group;
  #strict = false;

  /**
   * Returns a double of Point that reports to seat, with the options.
   * The double registers the check of the expected counts of calls for the
   * end of the test. An undefined seat reports nothing.
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
   * Returns the option that reads the latency, the windows of faults and the
   * instants of the calls from clock.
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
   * Returns the option that takes the numbers of probabilistic faults from
   * random, so a seeded source repeats the faults of a run.
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
   * Returns a double whose every asynchronous method waits duration on clock
   * before it returns, with the options after the pair.
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

/** An option of a BoxStub, which its constructor applies. */
export type BoxStubOption<T> = (double: BoxStub<T>) => void;

/**
 * BoxStub is a recording test double of Box. Each property
 * on<Method> configures its method. Left alone, a method returns undefined
 * and records the call.
 */
export class BoxStub<T> implements Box<T> {
  /** The property value, which a test assigns. */
  value!: T;
  readonly #group: Group;
  #strict = false;

  /**
   * Returns a double of Box that reports to seat, with the options.
   * The double registers the check of the expected counts of calls for the
   * end of the test. An undefined seat reports nothing.
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
   * Returns the option that reads the latency, the windows of faults and the
   * instants of the calls from clock.
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
   * Returns the option that takes the numbers of probabilistic faults from
   * random, so a seeded source repeats the faults of a run.
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
   * Returns a double whose every asynchronous method waits duration on clock
   * before it returns, with the options after the pair.
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
