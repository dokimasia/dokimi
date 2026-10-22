import type { Clock as Clock2, Seat } from "@dokimi/assert";
import type { RandomSource } from "@dokimi/runtime/fault";
import { type Arms, answer } from "@dokimi/runtime/stub/answer";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import { MethodStub } from "@dokimi/runtime/stub/method-stub";
import type { Clock } from "./clock";

/** ClockNowCall records one call of Clock.now. */
export interface ClockNowCall {
  /** The result of the call, where it returned one. */
  result?: number | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** ClockEveryCall records one call of Clock.every. */
export interface ClockEveryCall {
  /** The argument fn. */
  fn: () => void;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** ClockStopperCall records one call of Clock.stopper. */
export interface ClockStopperCall {
  /** The result of the call, where it returned one. */
  result?: (() => void) | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * ClockNowStub configures now and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class ClockNowStub extends MethodStub<ClockNowCall> {
  #fn: (() => number) | undefined;
  #fallback: { readonly value: number } | undefined;

  /**
   * Sets the result of every call of now. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: number): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of now from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: () => number): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of now with its arguments.
   */
  arms(): Arms<ClockNowCall, number, number> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn()),
      fallback: this.#fallback,
      stamp: (call, result) => {
        call.result = result;
      },
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/**
 * ClockEveryStub configures every and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class ClockEveryStub extends MethodStub<ClockEveryCall> {
  #fn: ((fn: () => void) => void) | undefined;

  /**
   * Sets the function that computes the result of every from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (fn: () => void) => void): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of every with its arguments.
   */
  arms(fn: () => void): Arms<ClockEveryCall, void, void> {
    const fn0 = this.#fn;
    return {
      invoke: fn0 && (() => fn0(fn)),
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/**
 * ClockStopperStub configures stopper and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class ClockStopperStub extends MethodStub<ClockStopperCall> {
  #fn: (() => () => void) | undefined;
  #fallback: { readonly value: () => void } | undefined;

  /**
   * Sets the result of every call of stopper. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: () => void): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of stopper from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: () => () => void): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of stopper with its arguments.
   */
  arms(): Arms<ClockStopperCall, () => void, () => void> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn()),
      fallback: this.#fallback,
      stamp: (call, result) => {
        call.result = result;
      },
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/** Configures a ClockStub at its construction. */
export type ClockStubOption = (double: ClockStub) => void;

/**
 * ClockStub is a recording test double of Clock. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class ClockStub implements Clock {
  /** The property zone, which a test assigns. */
  zone!: string;
  /** Configures now. */
  readonly onNow: ClockNowStub;
  /** Configures every. */
  readonly onEvery: ClockEveryStub;
  /** Configures stopper. */
  readonly onStopper: ClockStopperStub;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Clock that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: ClockStubOption[]) {
    this.onNow = new ClockNowStub(seat, "Clock.now");
    this.onEvery = new ClockEveryStub(seat, "Clock.every");
    this.onStopper = new ClockStopperStub(seat, "Clock.stopper");
    this.#group = new Group(this.onNow, this.onEvery, this.onStopper);
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
  static strict(): ClockStubOption {
    return (double) => {
      double.#strict = true;
    };
  }

  /**
   * Returns the option that passes each call to impl. The double still records
   * the call.
   *
   * @param impl - The implementation.
   * @returns The option.
   */
  static delegateTo(impl: Clock): ClockStubOption {
    return (double) => {
      double.onNow.func(impl.now.bind(impl));
      double.onEvery.func(impl.every.bind(impl));
      double.onStopper.func(impl.stopper.bind(impl));
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
  static withClock(clock: Clock2): ClockStubOption {
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
  static withRandSource(random: RandomSource): ClockStubOption {
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
  static benchMode(): ClockStubOption {
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
  static withLatency(duration: number): ClockStubOption {
    return (double) => {
      double.#group.setLatency(duration);
    };
  }

  /**
   * Returns the option that arms error on each method that can fail. Each call
   * of such a method then fails with error.
   *
   * @param error - The error.
   * @returns The option.
   */
  static failing(error: unknown): ClockStubOption {
    return (double) => {
      double.onNow.faults(error, 1);
      double.onEvery.faults(error, 1);
      double.onStopper.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of now.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withClockNow(fn: () => number): ClockStubOption {
    return (double) => {
      double.onNow.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of every.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withClockEvery(fn: (fn: () => void) => void): ClockStubOption {
    return (double) => {
      double.onEvery.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of stopper.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withClockStopper(fn: () => () => void): ClockStubOption {
    return (double) => {
      double.onStopper.func(fn);
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
    clock: Clock2,
    ...options: ClockStubOption[]
  ): ClockStub {
    return new ClockStub(
      seat,
      ClockStub.withClock(clock),
      ClockStub.withLatency(duration),
      ...options,
    );
  }

  /**
   * Returns a double on which each call of a method that can fail fails with
   * error. It applies options after the failure.
   *
   * @param seat - The seat of the test.
   * @param error - The error.
   * @param options - The options.
   * @returns The double.
   */
  static newFailing(
    seat: Seat | undefined,
    error: unknown,
    ...options: ClockStubOption[]
  ): ClockStub {
    return new ClockStub(seat, ClockStub.failing(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }

  /**
   * Records the call of now and returns its result.
   */
  now(): number {
    const call: ClockNowCall = {};
    return answer(this.onNow, call, this.onNow.arms());
  }

  /**
   * Records the call of every and returns its result.
   */
  every(fn: () => void): void {
    const call: ClockEveryCall = { fn };
    return answer(this.onEvery, call, this.onEvery.arms(fn));
  }

  /**
   * Records the call of stopper and returns its result.
   */
  stopper(): () => void {
    const call: ClockStopperCall = {};
    return answer(this.onStopper, call, this.onStopper.arms());
  }
}

void compatV1;
