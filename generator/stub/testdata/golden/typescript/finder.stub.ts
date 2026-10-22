import type { Clock, Seat } from "@dokimi/assert";
import type { RandomSource } from "@dokimi/runtime/fault";
import { type Arms, answer } from "@dokimi/runtime/stub/answer";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import { MethodStub } from "@dokimi/runtime/stub/method-stub";
import type { Finder } from "./finder";

/** FinderFindByIdCall records one call of Finder.findById. */
export interface FinderFindByIdCall {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: string | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * FinderFindByIdStub configures findById and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class FinderFindByIdStub extends MethodStub<FinderFindByIdCall> {
  #fn: ((id: string) => string) | undefined;
  #fallback: { readonly value: string } | undefined;

  /**
   * Sets the result of every call of findById. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: string): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of findById from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (id: string) => string): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of findById with its arguments.
   */
  arms(id: string): Arms<FinderFindByIdCall, string, string> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(id)),
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

/** Configures a FinderStub at its construction. */
export type FinderStubOption = (double: FinderStub) => void;

/**
 * FinderStub is a recording test double of Finder. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class FinderStub implements Finder {
  /** Configures findById. */
  readonly onFindById: FinderFindByIdStub;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Finder that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: FinderStubOption[]) {
    this.onFindById = new FinderFindByIdStub(seat, "Finder.findById");
    this.#group = new Group(this.onFindById);
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
  static strict(): FinderStubOption {
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
  static delegateTo(impl: Finder): FinderStubOption {
    return (double) => {
      double.onFindById.func(impl.findById.bind(impl));
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
  static withClock(clock: Clock): FinderStubOption {
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
  static withRandSource(random: RandomSource): FinderStubOption {
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
  static benchMode(): FinderStubOption {
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
  static withLatency(duration: number): FinderStubOption {
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
  static failing(error: unknown): FinderStubOption {
    return (double) => {
      double.onFindById.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of findById.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withFinderFindById(fn: (id: string) => string): FinderStubOption {
    return (double) => {
      double.onFindById.func(fn);
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
    ...options: FinderStubOption[]
  ): FinderStub {
    return new FinderStub(
      seat,
      FinderStub.withClock(clock),
      FinderStub.withLatency(duration),
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
    ...options: FinderStubOption[]
  ): FinderStub {
    return new FinderStub(seat, FinderStub.failing(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }

  /**
   * Records the call of findById and returns its result.
   */
  findById(id: string): string {
    const call: FinderFindByIdCall = { id };
    return answer(this.onFindById, call, this.onFindById.arms(id));
  }
}

void compatV1;
