import type { Clock, Seat } from "@dokimi/assert";
import { type RandomSource, retry } from "@dokimi/runtime/fault";
import { type Arms, answer, answerAsync } from "@dokimi/runtime/stub/answer";
import { compatV1 } from "@dokimi/runtime/stub/compat";
import { Group } from "@dokimi/runtime/stub/group";
import { MethodStub } from "@dokimi/runtime/stub/method-stub";
import { type Item, type Repo, type Store, type Table, type Wrap, notFound } from "./store";

/** StoreGetCall records one call of Store.get. */
export interface StoreGetCall {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: Item | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StorePutCall records one call of Store.put. */
export interface StorePutCall {
  /** The argument item. */
  item: Item;
  /** The argument note. */
  note?: string | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StoreSizeCall records one call of Store.size. */
export interface StoreSizeCall {
  /** The result of the call, where it returned one. */
  result?: number | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StoreLogCall records one call of Store.log. */
export interface StoreLogCall {
  /** The argument format. */
  format: string;
  /** The argument args. */
  args: string[];
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** StoreCloseCall records one call of Store.close. */
export interface StoreCloseCall {
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * StoreGetStub configures get and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreGetStub extends MethodStub<StoreGetCall> {
  #fn: ((id: string) => Item | PromiseLike<Item>) | undefined;
  #fallback: { readonly value: Item } | undefined;

  /**
   * Sets the result of every call of get. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: Item): this {
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
  func(fn: (id: string) => Item | PromiseLike<Item>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Makes the next call of get fail with notFound. The calls after it
   * succeed.
   *
   * @returns The stub.
   */
  faultNotFound(): this {
    return this.faultsOnce(notFound);
  }

  /**
   * Makes get fail with error on every call before call number
   * 3. Call 3 and every call after it succeed.
   *
   * @param error - The error.
   * @returns The stub.
   */
  retrySchedule(error: unknown): this {
    return this.setFault(retry(error, 3));
  }

  /**
   * Makes get fail every nth call whose id matches key. The
   * fault leaves the calls of the other partitions alone. A test of the
   * isolation of a partition needs a fault that selects its calls.
   *
   * @param key - The id of the calls that fail.
   * @param error - The error.
   * @param n - The period of the failed calls among the calls of the partition.
   * @returns The stub.
   */
  faultForPartition(key: string, error: unknown, n: number): this {
    return this.faultsWhen((call) => call.id === key, error, n);
  }

  /**
   * Makes get fail every nth call whose id differs from key.
   * The fault leaves the calls of the partition of key alone.
   *
   * @param key - The id of the calls that pass.
   * @param error - The error.
   * @param n - The period of the failed calls among the other calls.
   * @returns The stub.
   */
  faultForOtherPartitions(key: string, error: unknown, n: number): this {
    return this.faultsWhen((call) => call.id !== key, error, n);
  }

  /**
   * Returns the arms of a call of get with its arguments.
   */
  arms(id: string): Arms<StoreGetCall, Item, Item | PromiseLike<Item>> {
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
 * StorePutStub configures put and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StorePutStub extends MethodStub<StorePutCall> {
  #fn: ((item: Item, note?: string) => void) | undefined;

  /**
   * Sets the function that computes the result of put from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (item: Item, note?: string) => void): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of put with its arguments.
   */
  arms(item: Item, note?: string): Arms<StorePutCall, void, void> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(item, note)),
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/**
 * StoreSizeStub configures size and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreSizeStub extends MethodStub<StoreSizeCall> {
  #fn: (() => number) | undefined;
  #fallback: { readonly value: number } | undefined;

  /**
   * Sets the result of every call of size. A function and an injected
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
   * Sets the function that computes the result of size from its
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
   * Returns the arms of a call of size with its arguments.
   */
  arms(): Arms<StoreSizeCall, number, number> {
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
 * StoreLogStub configures log and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreLogStub extends MethodStub<StoreLogCall> {
  #fn: ((format: string, ...args: string[]) => void) | undefined;

  /**
   * Sets the function that computes the result of log from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (format: string, ...args: string[]) => void): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of log with its arguments.
   */
  arms(format: string, ...args: string[]): Arms<StoreLogCall, void, void> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(format, ...args)),
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/**
 * StoreCloseStub configures close and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class StoreCloseStub extends MethodStub<StoreCloseCall> {
  #fn: (() => void | PromiseLike<void>) | undefined;

  /**
   * Sets the function that computes the result of close from its
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
   * Returns the arms of a call of close with its arguments.
   */
  arms(): Arms<StoreCloseCall, void, void | PromiseLike<void>> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn()),
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
 * and returns undefined.
 */
export class StoreStub implements Store {
  /** The property zone, which a test assigns. */
  zone!: string;
  /** The property label, which a test assigns. */
  label?: string;
  /** Configures get. */
  readonly onGet: StoreGetStub;
  /** Configures put. */
  readonly onPut: StorePutStub;
  /** Configures size. */
  readonly onSize: StoreSizeStub;
  /** Configures log. */
  readonly onLog: StoreLogStub;
  /** Configures close, which Closer contributes. */
  readonly onClose: StoreCloseStub;
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
    this.onPut = new StorePutStub(seat, "Store.put");
    this.onSize = new StoreSizeStub(seat, "Store.size");
    this.onLog = new StoreLogStub(seat, "Store.log");
    this.onClose = new StoreCloseStub(seat, "Store.close");
    this.#group = new Group(this.onGet, this.onPut, this.onSize, this.onLog, this.onClose);
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
      double.onPut.func(impl.put.bind(impl));
      double.onSize.func(impl.size.bind(impl));
      double.onLog.func(impl.log.bind(impl));
      double.onClose.func(impl.close.bind(impl));
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
      double.onPut.faults(error, 1);
      double.onSize.faults(error, 1);
      double.onLog.faults(error, 1);
      double.onClose.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of get.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreGet(fn: (id: string) => Item | PromiseLike<Item>): StoreStubOption {
    return (double) => {
      double.onGet.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of put.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStorePut(fn: (item: Item, note?: string) => void): StoreStubOption {
    return (double) => {
      double.onPut.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of size.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreSize(fn: () => number): StoreStubOption {
    return (double) => {
      double.onSize.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of log.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreLog(fn: (format: string, ...args: string[]) => void): StoreStubOption {
    return (double) => {
      double.onLog.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of close.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withStoreClose(fn: () => void | PromiseLike<void>): StoreStubOption {
    return (double) => {
      double.onClose.func(fn);
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
  get(id: string): Promise<Item> {
    const call: StoreGetCall = { id };
    return answerAsync(this.onGet, call, this.onGet.arms(id));
  }

  /**
   * Records the call of put and returns its result.
   */
  put(item: Item, note?: string): void {
    const call: StorePutCall = { item, note };
    return answer(this.onPut, call, this.onPut.arms(item, note));
  }

  /**
   * Records the call of size and returns its result.
   */
  size(): number {
    const call: StoreSizeCall = {};
    return answer(this.onSize, call, this.onSize.arms());
  }

  /**
   * Records the call of log and returns its result.
   */
  log(format: string, ...args: string[]): void {
    const call: StoreLogCall = { format, args };
    return answer(this.onLog, call, this.onLog.arms(format, ...args));
  }

  /**
   * Records the call of close and returns its result.
   */
  close(): Promise<void> {
    const call: StoreCloseCall = {};
    return answerAsync(this.onClose, call, this.onClose.arms());
  }
}

// Item is a handle that a method of Store returns.

/** Configures a ItemStub at its construction. */
export type ItemStubOption = (double: ItemStub) => void;

/**
 * ItemStub is a recording test double of Item. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class ItemStub implements Item {
  /** The property id, which a test assigns. */
  id!: string;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Item that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: ItemStubOption[]) {
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
  static strict(): ItemStubOption {
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
  static withClock(clock: Clock): ItemStubOption {
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
  static withRandSource(random: RandomSource): ItemStubOption {
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
  static benchMode(): ItemStubOption {
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
  static withLatency(duration: number): ItemStubOption {
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
    ...options: ItemStubOption[]
  ): ItemStub {
    return new ItemStub(
      seat,
      ItemStub.withClock(clock),
      ItemStub.withLatency(duration),
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

/** RepoFindCall records one call of Repo.find. */
export interface RepoFindCall<T> {
  /** The argument id. */
  id: string;
  /** The result of the call, where it returned one. */
  result?: T | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * RepoFindStub configures find and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class RepoFindStub<T> extends MethodStub<RepoFindCall<T>> {
  #fn: ((id: string) => T) | undefined;
  #fallback: { readonly value: T } | undefined;

  /**
   * Sets the result of every call of find. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: T): this {
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
  func(fn: (id: string) => T): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of find with its arguments.
   */
  arms(id: string): Arms<RepoFindCall<T>, T, T> {
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

/** Configures a RepoStub at its construction. */
export type RepoStubOption<T> = (double: RepoStub<T>) => void;

/**
 * RepoStub is a recording test double of Repo. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class RepoStub<T> implements Repo<T> {
  /** Configures find. */
  readonly onFind: RepoFindStub<T>;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Repo that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: RepoStubOption<T>[]) {
    this.onFind = new RepoFindStub<T>(seat, "Repo.find");
    this.#group = new Group(this.onFind);
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
  static strict<T>(): RepoStubOption<T> {
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
  static delegateTo<T>(impl: Repo<T>): RepoStubOption<T> {
    return (double) => {
      double.onFind.func(impl.find.bind(impl));
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
  static withClock<T>(clock: Clock): RepoStubOption<T> {
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
  static withRandSource<T>(random: RandomSource): RepoStubOption<T> {
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
  static benchMode<T>(): RepoStubOption<T> {
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
  static withLatency<T>(duration: number): RepoStubOption<T> {
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
  static failing<T>(error: unknown): RepoStubOption<T> {
    return (double) => {
      double.onFind.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of find.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withRepoFind<T>(fn: (id: string) => T): RepoStubOption<T> {
    return (double) => {
      double.onFind.func(fn);
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
    ...options: RepoStubOption<T>[]
  ): RepoStub<T> {
    return new RepoStub<T>(
      seat,
      RepoStub.withClock<T>(clock),
      RepoStub.withLatency<T>(duration),
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
  static newFailing<T>(
    seat: Seat | undefined,
    error: unknown,
    ...options: RepoStubOption<T>[]
  ): RepoStub<T> {
    return new RepoStub<T>(seat, RepoStub.failing<T>(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }

  /**
   * Records the call of find and returns its result.
   */
  find(id: string): T {
    const call: RepoFindCall<T> = { id };
    return answer(this.onFind, call, this.onFind.arms(id));
  }
}

/** WrapPutCall records one call of Wrap.put. */
export interface WrapPutCall<T$> {
  /** The argument items. */
  items: T$[];
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * WrapPutStub configures put and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class WrapPutStub<T$> extends MethodStub<WrapPutCall<T$>> {
  #fn: ((items: T$[]) => void) | undefined;

  /**
   * Sets the function that computes the result of put from its
   * arguments. An injected fault takes precedence over it.
   *
   * @param fn - The function.
   * @returns The stub.
   */
  func(fn: (items: T$[]) => void): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of put with its arguments.
   */
  arms(items: T$[]): Arms<WrapPutCall<T$>, void, void> {
    const fn = this.#fn;
    return {
      invoke: fn && (() => fn(items)),
      fail: (call, error) => {
        call.error = error;
      },
    };
  }
}

/** Configures a WrapStub at its construction. */
export type WrapStubOption<T$> = (double: WrapStub<T$>) => void;

/**
 * WrapStub is a recording test double of Wrap. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class WrapStub<T$> implements Wrap<T$> {
  /** Configures put. */
  readonly onPut: WrapPutStub<T$>;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Wrap that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: WrapStubOption<T$>[]) {
    this.onPut = new WrapPutStub<T$>(seat, "Wrap.put");
    this.#group = new Group(this.onPut);
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
  static strict<T$>(): WrapStubOption<T$> {
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
  static delegateTo<T$>(impl: Wrap<T$>): WrapStubOption<T$> {
    return (double) => {
      double.onPut.func(impl.put.bind(impl));
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
  static withClock<T$>(clock: Clock): WrapStubOption<T$> {
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
  static withRandSource<T$>(random: RandomSource): WrapStubOption<T$> {
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
  static benchMode<T$>(): WrapStubOption<T$> {
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
  static withLatency<T$>(duration: number): WrapStubOption<T$> {
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
  static failing<T$>(error: unknown): WrapStubOption<T$> {
    return (double) => {
      double.onPut.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of put.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withWrapPut<T$>(fn: (items: T$[]) => void): WrapStubOption<T$> {
    return (double) => {
      double.onPut.func(fn);
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
  static newDelayed<T$>(
    seat: Seat | undefined,
    duration: number,
    clock: Clock,
    ...options: WrapStubOption<T$>[]
  ): WrapStub<T$> {
    return new WrapStub<T$>(
      seat,
      WrapStub.withClock<T$>(clock),
      WrapStub.withLatency<T$>(duration),
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
  static newFailing<T$>(
    seat: Seat | undefined,
    error: unknown,
    ...options: WrapStubOption<T$>[]
  ): WrapStub<T$> {
    return new WrapStub<T$>(seat, WrapStub.failing<T$>(error), ...options);
  }

  /**
   * Clears the recorded calls, the counters of the faults and the expected
   * counts. It keeps the functions, the fixed results and the faults.
   */
  resetCalls(): void {
    this.#group.reset();
  }

  /**
   * Records the call of put and returns its result.
   */
  put(items: T$[]): void {
    const call: WrapPutCall<T$> = { items };
    return answer(this.onPut, call, this.onPut.arms(items));
  }
}

/** TableGetCall records one call of Table.get. */
export interface TableGetCall<K, V> {
  /** The argument key. */
  key: K;
  /** The result of the call, where it returned one. */
  result?: V | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/** TableSizeCall records one call of Table.size. */
export interface TableSizeCall {
  /** The result of the call, where it returned one. */
  result?: number | undefined;
  /** The error that failed the call, where one did. */
  error?: unknown;
}

/**
 * TableGetStub configures get and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class TableGetStub<K, V> extends MethodStub<TableGetCall<K, V>> {
  #fn: ((key: K) => V | PromiseLike<V>) | undefined;
  #fallback: { readonly value: V } | undefined;

  /**
   * Sets the result of every call of get. A function and an injected
   * fault take precedence over it.
   *
   * @param result - The result.
   * @returns The stub.
   */
  returns(result: V): this {
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
  func(fn: (key: K) => V | PromiseLike<V>): this {
    this.#fn = fn;
    return this;
  }

  /**
   * Returns the arms of a call of get with its arguments.
   */
  arms(key: K): Arms<TableGetCall<K, V>, V, V | PromiseLike<V>> {
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

/**
 * TableSizeStub configures size and records its calls. A fault armed
 * through it takes precedence over the function and the fixed result.
 */
export class TableSizeStub extends MethodStub<TableSizeCall> {
  #fn: (() => number) | undefined;
  #fallback: { readonly value: number } | undefined;

  /**
   * Sets the result of every call of size. A function and an injected
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
   * Sets the function that computes the result of size from its
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
   * Returns the arms of a call of size with its arguments.
   */
  arms(): Arms<TableSizeCall, number, number> {
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

/** Configures a TableStub at its construction. */
export type TableStubOption<K, V> = (double: TableStub<K, V>) => void;

/**
 * TableStub is a recording test double of Table. Each property
 * on<Method> configures its method. An unconfigured method records the call
 * and returns undefined.
 */
export class TableStub<K, V> implements Table<K, V> {
  /** Configures get. */
  readonly onGet: TableGetStub<K, V>;
  /** Configures size. */
  readonly onSize: TableSizeStub;
  readonly #group: Group;
  #strict = false;

  /**
   * Creates a double of Table that reports its failures to seat and
   * applies options. The double checks the expected counts of calls at the end
   * of the test. With an undefined seat, the double reports nothing.
   *
   * @param seat - The seat of the test.
   * @param options - The options.
   */
  constructor(seat: Seat | undefined, ...options: TableStubOption<K, V>[]) {
    this.onGet = new TableGetStub<K, V>(seat, "Table.get");
    this.onSize = new TableSizeStub(seat, "Table.size");
    this.#group = new Group(this.onGet, this.onSize);
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
  static strict<K, V>(): TableStubOption<K, V> {
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
  static delegateTo<K, V>(impl: Table<K, V>): TableStubOption<K, V> {
    return (double) => {
      double.onGet.func(impl.get.bind(impl));
      double.onSize.func(impl.size.bind(impl));
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
  static withClock<K, V>(clock: Clock): TableStubOption<K, V> {
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
  static withRandSource<K, V>(random: RandomSource): TableStubOption<K, V> {
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
  static benchMode<K, V>(): TableStubOption<K, V> {
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
  static withLatency<K, V>(duration: number): TableStubOption<K, V> {
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
  static failing<K, V>(error: unknown): TableStubOption<K, V> {
    return (double) => {
      double.onGet.faults(error, 1);
      double.onSize.faults(error, 1);
    };
  }

  /**
   * Returns the option that sets the function of get.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withTableGet<K, V>(fn: (key: K) => V | PromiseLike<V>): TableStubOption<K, V> {
    return (double) => {
      double.onGet.func(fn);
    };
  }

  /**
   * Returns the option that sets the function of size.
   *
   * @param fn - The function.
   * @returns The option.
   */
  static withTableSize<K, V>(fn: () => number): TableStubOption<K, V> {
    return (double) => {
      double.onSize.func(fn);
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
  static newDelayed<K, V>(
    seat: Seat | undefined,
    duration: number,
    clock: Clock,
    ...options: TableStubOption<K, V>[]
  ): TableStub<K, V> {
    return new TableStub<K, V>(
      seat,
      TableStub.withClock<K, V>(clock),
      TableStub.withLatency<K, V>(duration),
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
  static newFailing<K, V>(
    seat: Seat | undefined,
    error: unknown,
    ...options: TableStubOption<K, V>[]
  ): TableStub<K, V> {
    return new TableStub<K, V>(seat, TableStub.failing<K, V>(error), ...options);
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
  get(key: K): Promise<V> {
    const call: TableGetCall<K, V> = { key };
    return answerAsync(this.onGet, call, this.onGet.arms(key));
  }

  /**
   * Records the call of size and returns its result.
   */
  size(): number {
    const call: TableSizeCall = {};
    return answer(this.onSize, call, this.onSize.arms());
  }
}

void compatV1;
