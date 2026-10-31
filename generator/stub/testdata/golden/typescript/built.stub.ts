import type { Clock, Seat } from "@dokimi/assert";
import type { RandomSource } from "@dokimi/runtime/fault";
import { type Arms, answer, answerAsync } from "@dokimi/runtime/stub/answer";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import { MethodStub } from "@dokimi/runtime/stub/method-stub";
import type { Invoice } from "./billing";
import type { Account, Note, Store } from "./built";
import { AccountBuilder } from "./built.builder";

/** StoreGetCall records one call of Store.get. */
export interface StoreGetCall {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: Account | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StoreFindCall records one call of Store.find. */
export interface StoreFindCall {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: Account | undefined | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StoreNoteCall records one call of Store.note. */
export interface StoreNoteCall {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: Note | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StoreBillCall records one call of Store.bill. */
export interface StoreBillCall {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: Invoice | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * StoreGetStub configures get and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreGetStub extends MethodStub<StoreGetCall> {
  #fn: ((id: string) => Account | PromiseLike<Account>) | undefined;
  #fallback: { readonly value: Account } | undefined;

  /**
   * Sets the result of every call of get. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: Account): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of get from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (id: string) => Account | PromiseLike<Account>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of get with its arguments.
   */
  arms(id: string): Arms<StoreGetCall, Account, Account | PromiseLike<Account>> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(id)),
      fallback: this.#fallback,
      default: () => AccountBuilder.create().build(),
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
 * StoreFindStub configures find and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreFindStub extends MethodStub<StoreFindCall> {
  #fn: ((id: string) => Account | undefined) | undefined;
  #fallback: { readonly value: Account | undefined } | undefined;

  /**
   * Sets the result of every call of find. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: Account | undefined): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of find from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (id: string) => Account | undefined): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of find with its arguments.
   */
  arms(id: string): Arms<StoreFindCall, Account | undefined, Account | undefined> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(id)),
      fallback: this.#fallback,
      default: () => AccountBuilder.create().build(),
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
 * StoreNoteStub configures note and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreNoteStub extends MethodStub<StoreNoteCall> {
  #fn: ((id: string) => Note) | undefined;
  #fallback: { readonly value: Note } | undefined;

  /**
   * Sets the result of every call of note. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: Note): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of note from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (id: string) => Note): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of note with its arguments.
   */
  arms(id: string): Arms<StoreNoteCall, Note, Note> {
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

/**
 * StoreBillStub configures bill and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreBillStub extends MethodStub<StoreBillCall> {
  #fn: ((id: string) => Invoice) | undefined;
  #fallback: { readonly value: Invoice } | undefined;

  /**
   * Sets the result of every call of bill. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: Invoice): this {
    this.#fallback = { value: result };
    return this;
  }

  /**
   * Sets the function that computes the result of bill from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (id: string) => Invoice): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of bill with its arguments.
   */
  arms(id: string): Arms<StoreBillCall, Invoice, Invoice> {
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

/** Configures a StoreStub at its construction. */
export type StoreStubOption = (double: StoreStub) => void;

/**
 * StoreStub is a recording test double of Store. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns a new value from the builder of its result type. It returns
 * undefined where the type has no builder.
 */
export class StoreStub implements Store {
  /** Configures get. */
  readonly onGet: StoreGetStub;
  /** Configures find. */
  readonly onFind: StoreFindStub;
  /** Configures note. */
  readonly onNote: StoreNoteStub;
  /** Configures bill. */
  readonly onBill: StoreBillStub;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Store that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: StoreStubOption[]) {
    this.onGet = new StoreGetStub(seat, "Store.get");
    this.onFind = new StoreFindStub(seat, "Store.find");
    this.onNote = new StoreNoteStub(seat, "Store.note");
    this.onBill = new StoreBillStub(seat, "Store.bill");
    this.#group = new Group(this.onGet, this.onFind, this.onNote, this.onBill);
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
  static strict(): StoreStubOption {
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
  static delegateTo(impl: Store): StoreStubOption {
    return (double) => {
      double.onGet.func(impl.get.bind(impl));
      double.onFind.func(impl.find.bind(impl));
      double.onNote.func(impl.note.bind(impl));
      double.onBill.func(impl.bill.bind(impl));
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
  static withClock(clock: Clock): StoreStubOption {
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
  static withRandSource(random: RandomSource): StoreStubOption {
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
  static benchMode(): StoreStubOption {
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
  static withLatency(duration: number): StoreStubOption {
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
  static failing(error: unknown): StoreStubOption {
    return (double) => {
      double.onGet.faults(error, 1);
      double.onFind.faults(error, 1);
      double.onNote.faults(error, 1);
      double.onBill.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of get.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreGet(fn: (id: string) => Account | PromiseLike<Account>): StoreStubOption {
    return (double) => {
      double.onGet.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of find.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreFind(fn: (id: string) => Account | undefined): StoreStubOption {
    return (double) => {
      double.onFind.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of note.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreNote(fn: (id: string) => Note): StoreStubOption {
    return (double) => {
      double.onNote.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of bill.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreBill(fn: (id: string) => Invoice): StoreStubOption {
    return (double) => {
      double.onBill.func(fn);
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
    ...options: StoreStubOption[]
  ): StoreStub {
    return new StoreStub(
      seat,
      StoreStub.withClock(clock),
      StoreStub.withLatency(duration),
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
    ...options: StoreStubOption[]
  ): StoreStub {
    return new StoreStub(seat, StoreStub.failing(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }

  /**
   * Records the call of get and returns its result.
   */
  get(id: string): Promise<Account> {
    const call: StoreGetCall = { id };
    return answerAsync(this.onGet, call, this.onGet.arms(id));
  }

  /**
   * Records the call of find and returns its result.
   */
  find(id: string): Account | undefined {
    const call: StoreFindCall = { id };
    return answer(this.onFind, call, this.onFind.arms(id));
  }

  /**
   * Records the call of note and returns its result.
   */
  note(id: string): Note {
    const call: StoreNoteCall = { id };
    return answer(this.onNote, call, this.onNote.arms(id));
  }

  /**
   * Records the call of bill and returns its result.
   */
  bill(id: string): Invoice {
    const call: StoreBillCall = { id };
    return answer(this.onBill, call, this.onBill.arms(id));
  }
}

void compatV1;
