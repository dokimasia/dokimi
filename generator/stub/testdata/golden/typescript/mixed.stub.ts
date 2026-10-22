import type { Clock, Seat } from "@dokimi/assert";
import type { RandomSource } from "@dokimi/runtime/fault";
import { type Arms, answer, answerAsync } from "@dokimi/runtime/stub/answer";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import { MethodStub } from "@dokimi/runtime/stub/method-stub";
import { OrderTracker } from "@dokimi/runtime/stub/order-tracker";
import type { Session } from "./mixed";

/** SessionPrepareCall records one call of Session.prepare. */
export interface SessionPrepareCall {
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** SessionCommitCall records one call of Session.commit. */
export interface SessionCommitCall {
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** SessionOldCall records one call of Session.old. */
export interface SessionOldCall {
  /** The argument key. */
  key: string;
  /** The result of the call, where it returned one. */
  result?: string | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * SessionPrepareStub configures prepare and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class SessionPrepareStub extends MethodStub<SessionPrepareCall> {
  #fn: (() => void | PromiseLike<void>) | undefined;

  /**
   * Sets the function that computes the result of prepare from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: () => void | PromiseLike<void>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of prepare with its arguments.
   */
  arms(): Arms<SessionPrepareCall, void, void | PromiseLike<void>> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn()),
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/**
 * SessionCommitStub configures commit and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class SessionCommitStub extends MethodStub<SessionCommitCall> {
  #fn: (() => void | PromiseLike<void>) | undefined;

  /**
   * Sets the function that computes the result of commit from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: () => void | PromiseLike<void>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of commit with its arguments.
   */
  arms(): Arms<SessionCommitCall, void, void | PromiseLike<void>> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn()),
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/**
 * SessionOldStub configures old and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class SessionOldStub extends MethodStub<SessionOldCall> {
  #fn: ((key: string) => string) | undefined;
  #fallback: { readonly value: string } | undefined;

  /**
   * Sets the result of every call of old. A function and an injected
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
   * Sets the function that computes the result of old from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (key: string) => string): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of old with its arguments.
   */
  arms(key: string): Arms<SessionOldCall, string, string> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(key)),
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

/** Configures a SessionStub at its construction. */
export type SessionStubOption = (double: SessionStub) => void;

/**
 * SessionStub is a recording test double of Session. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class SessionStub implements Session {
  /** Configures prepare. */
  readonly onPrepare: SessionPrepareStub;
  /** Configures commit. */
  readonly onCommit: SessionCommitStub;
  /** Configures old. */
  readonly onOld: SessionOldStub;
  readonly #group: Group;
  readonly #order: OrderTracker;
  #strict = false;

  /**
   * Creates a double of Session that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: SessionStubOption[]) {
    this.onPrepare = new SessionPrepareStub(seat, "Session.prepare");
    this.onCommit = new SessionCommitStub(seat, "Session.commit");
    this.onOld = new SessionOldStub(seat, "Session.old");
    this.#group = new Group(this.onPrepare, this.onCommit, this.onOld);
    for (const option of options) {
      option(this);
    }
    if (this.#strict) {
      this.#group.strictAll();
    }
    this.#order = new OrderTracker(seat, this.#strict);
    this.#group.bind(seat);
  }

  /**
   * Returns the option that makes an unconfigured method fail the test.
   *
   * @returns The option.
   */
  static strict(): SessionStubOption {
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
  static delegateTo(impl: Session): SessionStubOption {
    return (double) => {
      double.onPrepare.func(impl.prepare.bind(impl));
      double.onCommit.func(impl.commit.bind(impl));
      double.onOld.func(impl.old.bind(impl));
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
  static withClock(clock: Clock): SessionStubOption {
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
  static withRandSource(random: RandomSource): SessionStubOption {
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
  static benchMode(): SessionStubOption {
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
  static withLatency(duration: number): SessionStubOption {
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
  static failing(error: unknown): SessionStubOption {
    return (double) => {
      double.onPrepare.faults(error, 1);
      double.onCommit.faults(error, 1);
      double.onOld.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of prepare.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withSessionPrepare(fn: () => void | PromiseLike<void>): SessionStubOption {
    return (double) => {
      double.onPrepare.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of commit.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withSessionCommit(fn: () => void | PromiseLike<void>): SessionStubOption {
    return (double) => {
      double.onCommit.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of old.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withSessionOld(fn: (key: string) => string): SessionStubOption {
    return (double) => {
      double.onOld.func(fn);
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
    ...options: SessionStubOption[]
  ): SessionStub {
    return new SessionStub(
      seat,
      SessionStub.withClock(clock),
      SessionStub.withLatency(duration),
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
    ...options: SessionStubOption[]
  ): SessionStub {
    return new SessionStub(seat, SessionStub.failing(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
    this.#order.reset();
  }

  /**
   * Records the call of prepare and returns its result.
   */
  prepare(): Promise<void> {
    this.#order.record("prepare");
    const call: SessionPrepareCall = {};
    return answerAsync(this.onPrepare, call, this.onPrepare.arms());
  }

  /**
   * Records the call of commit and returns its result.
   */
  commit(): Promise<void> {
    this.#order.assertAfter("commit", "prepare");
    this.#order.record("commit");
    const call: SessionCommitCall = {};
    return answerAsync(this.onCommit, call, this.onCommit.arms());
  }

  /**
   * Records the call of old and returns its result.
   */
  old(key: string): string {
    this.#order.record("old");
    const call: SessionOldCall = { key };
    return answer(this.onOld, call, this.onOld.arms(key));
  }
}

void compatV1;
