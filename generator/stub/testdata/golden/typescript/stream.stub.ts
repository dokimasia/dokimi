import type { Clock, Seat } from "@dokimi/assert";
import type { RandomSource } from "@dokimi/runtime/fault";
import { type Arms, answer } from "@dokimi/runtime/stub/answer";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import { MethodStub } from "@dokimi/runtime/stub/method-stub";
import type { Feed } from "./stream";

/** FeedNamesCall records one call of Feed.names. */
export interface FeedNamesCall {
  /** The result of the call, where it returned one. */
  result?: Iterable<string> | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** FeedRowsCall records one call of Feed.rows. */
export interface FeedRowsCall {
  /** The result of the call, where it returned one. */
  result?: AsyncIterable<number> | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * FeedNamesStub configures names and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class FeedNamesStub extends MethodStub<FeedNamesCall> {
  #fn: (() => Iterable<string>) | undefined;
  #fallback: { readonly value: Iterable<string> } | undefined;

  /**
   * Sets the result of every call of names. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: Iterable<string>): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of names from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: () => Iterable<string>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of names with its arguments.
   */
  arms(): Arms<FeedNamesCall, Iterable<string>, Iterable<string>> {
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
 * FeedRowsStub configures rows and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class FeedRowsStub extends MethodStub<FeedRowsCall> {
  #fn: (() => AsyncIterable<number>) | undefined;
  #fallback: { readonly value: AsyncIterable<number> } | undefined;

  /**
   * Sets the result of every call of rows. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: AsyncIterable<number>): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of rows from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: () => AsyncIterable<number>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Makes rows return an asynchronous sequence of the items.
   *
   * @param items - The items.
   * @returns The stub.
   */
  yields(...items: number[]): this {
    return this.func(async function* () {
      yield* items;
    });
  }

  /**
   * Returns the arms of a call of rows with its arguments.
   */
  arms(): Arms<FeedRowsCall, AsyncIterable<number>, AsyncIterable<number>> {
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

/** Configures a FeedStub at its construction. */
export type FeedStubOption = (double: FeedStub) => void;

/**
 * FeedStub is a recording test double of Feed. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class FeedStub implements Feed {
  /** Configures names. */
  readonly onNames: FeedNamesStub;
  /** Configures rows. */
  readonly onRows: FeedRowsStub;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Feed that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: FeedStubOption[]) {
    this.onNames = new FeedNamesStub(seat, "Feed.names");
    this.onRows = new FeedRowsStub(seat, "Feed.rows");
    this.#group = new Group(this.onNames, this.onRows);
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
  static strict(): FeedStubOption {
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
  static delegateTo(impl: Feed): FeedStubOption {
    return (double) => {
      double.onNames.func(impl.names.bind(impl));
      double.onRows.func(impl.rows.bind(impl));
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
  static withClock(clock: Clock): FeedStubOption {
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
  static withRandSource(random: RandomSource): FeedStubOption {
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
  static benchMode(): FeedStubOption {
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
  static withLatency(duration: number): FeedStubOption {
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
  static failing(error: unknown): FeedStubOption {
    return (double) => {
      double.onNames.faults(error, 1);
      double.onRows.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of names.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withFeedNames(fn: () => Iterable<string>): FeedStubOption {
    return (double) => {
      double.onNames.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of rows.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withFeedRows(fn: () => AsyncIterable<number>): FeedStubOption {
    return (double) => {
      double.onRows.func(fn);
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
    ...options: FeedStubOption[]
  ): FeedStub {
    return new FeedStub(
      seat,
      FeedStub.withClock(clock),
      FeedStub.withLatency(duration),
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
    ...options: FeedStubOption[]
  ): FeedStub {
    return new FeedStub(seat, FeedStub.failing(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }

  /**
   * Records the call of names and returns its result.
   */
  names(): Iterable<string> {
    const call: FeedNamesCall = {};
    return answer(this.onNames, call, this.onNames.arms());
  }

  /**
   * Records the call of rows and returns its result.
   */
  rows(): AsyncIterable<number> {
    const call: FeedRowsCall = {};
    return answer(this.onRows, call, this.onRows.arms());
  }
}

void compatV1;
