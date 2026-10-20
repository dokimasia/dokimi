import { Controlled, type Seat, byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import { type Item, notFound } from './store';
import { ItemStub, RepoStub, StoreStub, TableStub, WrapStub } from './store.stub';

/** Returns the samples of the arguments of get, or undefined. */
function storeStubGetArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'get-a';
  return [a0];
}

/** Calls get of double with the samples of its arguments. */
function storeStubGet(double: StoreStub): Promise<Item> {
  const args = storeStubGetArgs();
  return double.get(args[0]);
}

/** Returns the sample of the result of get. */
function storeStubGetPin(): Item {
  let want0: Item;
  want0 = {id: 'id-a'} as Item;
  return want0;
}

/** Returns the samples of the arguments of put, or undefined. */
function storeStubPutArgs(): [Item, string | undefined] {
  let a0 = undefined as unknown as Item;
  let a1 = undefined as unknown as string;
  a0 = {id: 'id-a'} as Item;
  a1 = 'put-a';
  return [a0, a1];
}

/** Calls put of double with the samples of its arguments. */
function storeStubPut(double: StoreStub): void {
  const args = storeStubPutArgs();
  return double.put(args[0], args[1]);
}

/** Calls size of double with the samples of its arguments. */
function storeStubSize(double: StoreStub): number {
  return double.size();
}

/** Returns the sample of the result of size. */
function storeStubSizePin(): number {
  let want0: number;
  want0 = 1.5;
  return want0;
}

/** Returns the samples of the arguments of log, or undefined. */
function storeStubLogArgs(): [string, string[]] {
  let a0 = undefined as unknown as string;
  const a1: string[] = [];
  a0 = 'log-a';
  return [a0, a1];
}

/** Calls log of double with the samples of its arguments. */
function storeStubLog(double: StoreStub): void {
  const args = storeStubLogArgs();
  return double.log(args[0], ...args[1]);
}

/** Calls close of double with the samples of its arguments. */
function storeStubClose(double: StoreStub): Promise<void> {
  return double.close();
}

describe("StoreStub.get", () => {
  const checks = behaviour("Store.get", (seat) => {
    const double = new StoreStub(seat);
    return {
      stub: double.onGet,
      call: () => storeStubGet(double),
      override: (mark) => {
        double.onGet.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubGetPin();
    double.onGet.returns(want);
    check.equal(seat, await storeStubGet(double), want, "get returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    await storeStubGet(double);
    const got = double.onGet.assertCalledOnce(seat, "get records the call");
    const args = storeStubGetArgs();
    check.equal(seat, got?.id, args[0], "the call records id");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onGet.onRecord(() => {
      seen += 1;
    });
    await storeStubGet(double);
    await storeStubGet(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of get");
  });

  test("sets the function of withStoreGet at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreGet(() => {
      called = true;
      return undefined as never;
    });
    await storeStubGet(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreGet sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("get-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubGet(double));
    const message = "get fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("fails the next call with notFound once faultNotFound arms it", async ({ seat }) => {
    const double = new StoreStub(seat);
    double.onGet.faultNotFound();
    const first = await settle(() => storeStubGet(double));
    const message = "faultNotFound fails the next call of get with notFound";
    check.equal(seat, first.threw && first.error, notFound, message, byIdentity());
    const second = await settle(() => storeStubGet(double));
    check.isFalse(seat, second.threw, "the call after it succeeds");
  });

  test("fails until attempt 3 then succeeds", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = new Error("get-retry");
    double.onGet.retrySchedule(want);
    for (let attempt = 1; attempt < 3; attempt += 1) {
      const outcome = await settle(() => storeStubGet(double));
      const message = "a call before attempt 3 fails";
      check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
    }
    const outcome = await settle(() => storeStubGet(double));
    check.isFalse(seat, outcome.threw, "attempt 3 succeeds");
  });

  test("faults only the targeted partition", async ({ seat }) => {
    const double = new StoreStub(seat);
    const args = storeStubGetArgs();
    const want = new Error("get-partition");
    double.onGet.faultForPartition(args[0], want, 1);
    const outcome = await settle(() => storeStubGet(double));
    const message = "a call in the targeted partition fails";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("spares the partition that faultForOtherPartitions names", async ({ seat }) => {
    const double = new StoreStub(seat);
    const args = storeStubGetArgs();
    double.onGet.faultForOtherPartitions(args[0], new Error("get-partition"), 1);
    const outcome = await settle(() => storeStubGet(double));
    check.isFalse(seat, outcome.threw, "a call in the named partition passes");
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubGetPin();
    double.onGet.returns(want);
    await storeStubGet(double);
    double.resetCalls();
    check.equal(seat, await storeStubGet(double), want, "a reset keeps the result of returns");
  });
});

describe("StoreStub.put", () => {
  const checks = behaviour("Store.put", (seat) => {
    const double = new StoreStub(seat);
    return {
      stub: double.onPut,
      call: () => storeStubPut(double),
      override: (mark) => {
        double.onPut.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    storeStubPut(double);
    const got = double.onPut.assertCalledOnce(seat, "put records the call");
    const args = storeStubPutArgs();
    check.equal(seat, got?.item, args[0], "the call records item");
    check.equal(seat, got?.note, args[1], "the call records note");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onPut.onRecord(() => {
      seen += 1;
    });
    storeStubPut(double);
    storeStubPut(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of put");
  });

  test("sets the function of withStorePut at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStorePut(() => {
      called = true;
      return undefined as never;
    });
    storeStubPut(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStorePut sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("put-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubPut(double));
    const message = "put fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("StoreStub.size", () => {
  const checks = behaviour("Store.size", (seat) => {
    const double = new StoreStub(seat);
    return {
      stub: double.onSize,
      call: () => storeStubSize(double),
      override: (mark) => {
        double.onSize.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubSizePin();
    double.onSize.returns(want);
    check.equal(seat, storeStubSize(double), want, "size returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    storeStubSize(double);
    double.onSize.assertCalledOnce(seat, "size records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onSize.onRecord(() => {
      seen += 1;
    });
    storeStubSize(double);
    storeStubSize(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of size");
  });

  test("sets the function of withStoreSize at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreSize(() => {
      called = true;
      return undefined as never;
    });
    storeStubSize(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreSize sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("size-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubSize(double));
    const message = "size fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubSizePin();
    double.onSize.returns(want);
    storeStubSize(double);
    double.resetCalls();
    check.equal(seat, storeStubSize(double), want, "a reset keeps the result of returns");
  });
});

describe("StoreStub.log", () => {
  const checks = behaviour("Store.log", (seat) => {
    const double = new StoreStub(seat);
    return {
      stub: double.onLog,
      call: () => storeStubLog(double),
      override: (mark) => {
        double.onLog.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    storeStubLog(double);
    const got = double.onLog.assertCalledOnce(seat, "log records the call");
    const args = storeStubLogArgs();
    check.equal(seat, got?.format, args[0], "the call records format");
    check.equal(seat, got?.args, args[1], "the call records args");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onLog.onRecord(() => {
      seen += 1;
    });
    storeStubLog(double);
    storeStubLog(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of log");
  });

  test("sets the function of withStoreLog at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreLog(() => {
      called = true;
      return undefined as never;
    });
    storeStubLog(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreLog sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("log-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubLog(double));
    const message = "log fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("StoreStub.close", () => {
  const checks = behaviour("Store.close", (seat) => {
    const double = new StoreStub(seat);
    return {
      stub: double.onClose,
      call: () => storeStubClose(double),
      override: (mark) => {
        double.onClose.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    await storeStubClose(double);
    double.onClose.assertCalledOnce(seat, "close records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onClose.onRecord(() => {
      seen += 1;
    });
    await storeStubClose(double);
    await storeStubClose(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of close");
  });

  test("sets the function of withStoreClose at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreClose(() => {
      called = true;
      return undefined as never;
    });
    await storeStubClose(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreClose sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("close-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubClose(double));
    const message = "close fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("StoreStub", () => {
  const instance = (double: StoreStub) => ({
    stub: double.onGet,
    call: () => storeStubGet(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new StoreStub(seat)),
    withClock: (seat, clock) => instance(new StoreStub(seat, StoreStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(StoreStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new StoreStub(seat, StoreStub.withRandSource(random))),
    benchMode: (seat) => instance(new StoreStub(seat, StoreStub.benchMode())),
    strict: (seat) => instance(new StoreStub(seat, StoreStub.strict())),
    async: true,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes get to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubGet(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onGet.assertCalledOnce(seat, "get reaches the implementation");
  });

  test("fails get with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("get-delegate");
    const inner = new StoreStub(seat);
    inner.onGet.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubGet(double));
    const message = "get fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes put to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubPut(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onPut.assertCalledOnce(seat, "put reaches the implementation");
  });

  test("fails put with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("put-delegate");
    const inner = new StoreStub(seat);
    inner.onPut.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubPut(double));
    const message = "put fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes size to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubSize(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onSize.assertCalledOnce(seat, "size reaches the implementation");
  });

  test("fails size with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("size-delegate");
    const inner = new StoreStub(seat);
    inner.onSize.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubSize(double));
    const message = "size fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes log to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubLog(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onLog.assertCalledOnce(seat, "log reaches the implementation");
  });

  test("fails log with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("log-delegate");
    const inner = new StoreStub(seat);
    inner.onLog.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubLog(double));
    const message = "log fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes close to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubClose(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onClose.assertCalledOnce(seat, "close reaches the implementation");
  });

  test("fails close with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("close-delegate");
    const inner = new StoreStub(seat);
    inner.onClose.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubClose(double));
    const message = "close fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

/** Assigns the sample of each property of ItemStub to double. */
function itemStubProps(double: ItemStub): void {
  double.id = 'id-a';
}

describe("ItemStub", () => {
  const builds: [string, (seat: Seat) => ItemStub][] = [
    ["new", (seat) => new ItemStub(seat)],
    ["strict", (seat) => new ItemStub(seat, ItemStub.strict())],
    ["withClock", (seat) => new ItemStub(seat, ItemStub.withClock(new Controlled()))],
    ["withRandSource", (seat) => new ItemStub(seat, ItemStub.withRandSource(() => 0))],
    ["benchMode", (seat) => new ItemStub(seat, ItemStub.benchMode())],
    ["withLatency", (seat) => new ItemStub(seat, ItemStub.withLatency(1))],
    ["newDelayed", (seat) => ItemStub.newDelayed(seat, 1, new Controlled())],
  ];
  for (const [name, build] of builds) {
    test(`keeps each property across a reset of a double built by ${name}`, ({ seat }) => {
      const double = build(seat);
      itemStubProps(double);
      double.resetCalls();
      const want = new ItemStub(undefined);
      itemStubProps(want);
      check.equal(seat, double.id, want.id, "the reset keeps the property id");
    });
  }
});

/** Returns the samples of the arguments of find, or undefined. */
function repoStubFindArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'find-a';
  return [a0];
}

/** Calls find of double with the samples of its arguments. */
function repoStubFind<T>(double: RepoStub<T>): T {
  const args = repoStubFindArgs();
  return double.find(args[0]);
}

describe("RepoStub.find", () => {
  const checks = behaviour("Repo.find", (seat) => {
    const double = new RepoStub<unknown>(seat);
    return {
      stub: double.onFind,
      call: () => repoStubFind(double),
      override: (mark) => {
        double.onFind.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new RepoStub<unknown>(seat);
    repoStubFind(double);
    const got = double.onFind.assertCalledOnce(seat, "find records the call");
    const args = repoStubFindArgs();
    check.equal(seat, got?.id, args[0], "the call records id");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new RepoStub<unknown>(seat);
    let seen = 0;
    double.onFind.onRecord(() => {
      seen += 1;
    });
    repoStubFind(double);
    repoStubFind(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of find");
  });

  test("sets the function of withRepoFind at construction", async ({ seat }) => {
    let called = false;
    const option = RepoStub.withRepoFind<unknown>(() => {
      called = true;
      return undefined as never;
    });
    repoStubFind(new RepoStub<unknown>(seat, option));
    check.isTrue(seat, called, "withRepoFind sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("find-failing");
    const double = RepoStub.newFailing<unknown>(seat, want);
    const outcome = await settle(() => repoStubFind(double));
    const message = "find fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps recording across a reset", async ({ seat }) => {
    const double = new RepoStub<unknown>(seat);
    double.onFind.returns(undefined as never);
    repoStubFind(double);
    double.resetCalls();
    repoStubFind(double);
    double.onFind.assertCalledOnce(seat, "find records the call after the reset");
  });
});

describe("RepoStub", () => {
  const instance = (double: RepoStub<unknown>) => ({
    stub: double.onFind,
    call: () => repoStubFind(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new RepoStub<unknown>(seat)),
    withClock: (seat, clock) => instance(new RepoStub<unknown>(seat, RepoStub.withClock<unknown>(clock))),
    withLatency: (seat, duration, clock) => instance(RepoStub.newDelayed<unknown>(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new RepoStub<unknown>(seat, RepoStub.withRandSource<unknown>(random))),
    benchMode: (seat) => instance(new RepoStub<unknown>(seat, RepoStub.benchMode<unknown>())),
    strict: (seat) => instance(new RepoStub<unknown>(seat, RepoStub.strict<unknown>())),
    async: false,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes find to the implementation of delegateTo", async ({ seat }) => {
    const inner = new RepoStub<unknown>(seat);
    await settle(() => repoStubFind(new RepoStub<unknown>(seat, RepoStub.delegateTo<unknown>(inner))));
    inner.onFind.assertCalledOnce(seat, "find reaches the implementation");
  });

  test("fails find with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("find-delegate");
    const inner = new RepoStub<unknown>(seat);
    inner.onFind.faults(want, 1);
    const double = new RepoStub<unknown>(seat, RepoStub.delegateTo<unknown>(inner));
    const outcome = await settle(() => repoStubFind(double));
    const message = "find fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

/** Returns the samples of the arguments of put, or undefined. */
function wrapStubPutArgs<T$>(): [T$[]] {
  let a0 = undefined as unknown as T$[];
  return [a0];
}

/** Calls put of double with the samples of its arguments. */
function wrapStubPut<T$>(double: WrapStub<T$>): void {
  const args = wrapStubPutArgs<T$>();
  return double.put(args[0]);
}

describe("WrapStub.put", () => {
  const checks = behaviour("Wrap.put", (seat) => {
    const double = new WrapStub<unknown>(seat);
    return {
      stub: double.onPut,
      call: () => wrapStubPut(double),
      override: (mark) => {
        double.onPut.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new WrapStub<unknown>(seat);
    wrapStubPut(double);
    const got = double.onPut.assertCalledOnce(seat, "put records the call");
    const args = wrapStubPutArgs<unknown>();
    check.equal(seat, got?.items, args[0], "the call records items");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new WrapStub<unknown>(seat);
    let seen = 0;
    double.onPut.onRecord(() => {
      seen += 1;
    });
    wrapStubPut(double);
    wrapStubPut(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of put");
  });

  test("sets the function of withWrapPut at construction", async ({ seat }) => {
    let called = false;
    const option = WrapStub.withWrapPut<unknown>(() => {
      called = true;
      return undefined as never;
    });
    wrapStubPut(new WrapStub<unknown>(seat, option));
    check.isTrue(seat, called, "withWrapPut sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("put-failing");
    const double = WrapStub.newFailing<unknown>(seat, want);
    const outcome = await settle(() => wrapStubPut(double));
    const message = "put fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("WrapStub", () => {
  const instance = (double: WrapStub<unknown>) => ({
    stub: double.onPut,
    call: () => wrapStubPut(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new WrapStub<unknown>(seat)),
    withClock: (seat, clock) => instance(new WrapStub<unknown>(seat, WrapStub.withClock<unknown>(clock))),
    withLatency: (seat, duration, clock) => instance(WrapStub.newDelayed<unknown>(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new WrapStub<unknown>(seat, WrapStub.withRandSource<unknown>(random))),
    benchMode: (seat) => instance(new WrapStub<unknown>(seat, WrapStub.benchMode<unknown>())),
    strict: (seat) => instance(new WrapStub<unknown>(seat, WrapStub.strict<unknown>())),
    async: false,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes put to the implementation of delegateTo", async ({ seat }) => {
    const inner = new WrapStub<unknown>(seat);
    await settle(() => wrapStubPut(new WrapStub<unknown>(seat, WrapStub.delegateTo<unknown>(inner))));
    inner.onPut.assertCalledOnce(seat, "put reaches the implementation");
  });

  test("fails put with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("put-delegate");
    const inner = new WrapStub<unknown>(seat);
    inner.onPut.faults(want, 1);
    const double = new WrapStub<unknown>(seat, WrapStub.delegateTo<unknown>(inner));
    const outcome = await settle(() => wrapStubPut(double));
    const message = "put fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

/** Returns the samples of the arguments of get, or undefined. */
function tableStubGetArgs<K>(): [K] {
  let a0 = undefined as unknown as K;
  return [a0];
}

/** Calls get of double with the samples of its arguments. */
function tableStubGet<K, V>(double: TableStub<K, V>): Promise<V> {
  const args = tableStubGetArgs<K>();
  return double.get(args[0]);
}

/** Calls size of double with the samples of its arguments. */
function tableStubSize<K, V>(double: TableStub<K, V>): number {
  return double.size();
}

/** Returns the sample of the result of size. */
function tableStubSizePin(): number {
  let want0: number;
  want0 = 1.5;
  return want0;
}

describe("TableStub.get", () => {
  const checks = behaviour("Table.get", (seat) => {
    const double = new TableStub<unknown, unknown>(seat);
    return {
      stub: double.onGet,
      call: () => tableStubGet(double),
      override: (mark) => {
        double.onGet.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    await tableStubGet(double);
    const got = double.onGet.assertCalledOnce(seat, "get records the call");
    const args = tableStubGetArgs<unknown>();
    check.equal(seat, got?.key, args[0], "the call records key");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    let seen = 0;
    double.onGet.onRecord(() => {
      seen += 1;
    });
    await tableStubGet(double);
    await tableStubGet(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of get");
  });

  test("sets the function of withTableGet at construction", async ({ seat }) => {
    let called = false;
    const option = TableStub.withTableGet<unknown, unknown>(() => {
      called = true;
      return undefined as never;
    });
    await tableStubGet(new TableStub<unknown, unknown>(seat, option));
    check.isTrue(seat, called, "withTableGet sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("get-failing");
    const double = TableStub.newFailing<unknown, unknown>(seat, want);
    const outcome = await settle(() => tableStubGet(double));
    const message = "get fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps recording across a reset", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    double.onGet.returns(undefined as never);
    await tableStubGet(double);
    double.resetCalls();
    await tableStubGet(double);
    double.onGet.assertCalledOnce(seat, "get records the call after the reset");
  });
});

describe("TableStub.size", () => {
  const checks = behaviour("Table.size", (seat) => {
    const double = new TableStub<unknown, unknown>(seat);
    return {
      stub: double.onSize,
      call: () => tableStubSize(double),
      override: (mark) => {
        double.onSize.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    const want = tableStubSizePin();
    double.onSize.returns(want);
    check.equal(seat, tableStubSize(double), want, "size returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    tableStubSize(double);
    double.onSize.assertCalledOnce(seat, "size records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    let seen = 0;
    double.onSize.onRecord(() => {
      seen += 1;
    });
    tableStubSize(double);
    tableStubSize(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of size");
  });

  test("sets the function of withTableSize at construction", async ({ seat }) => {
    let called = false;
    const option = TableStub.withTableSize<unknown, unknown>(() => {
      called = true;
      return undefined as never;
    });
    tableStubSize(new TableStub<unknown, unknown>(seat, option));
    check.isTrue(seat, called, "withTableSize sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("size-failing");
    const double = TableStub.newFailing<unknown, unknown>(seat, want);
    const outcome = await settle(() => tableStubSize(double));
    const message = "size fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new TableStub<unknown, unknown>(seat);
    const want = tableStubSizePin();
    double.onSize.returns(want);
    tableStubSize(double);
    double.resetCalls();
    check.equal(seat, tableStubSize(double), want, "a reset keeps the result of returns");
  });
});

describe("TableStub", () => {
  const instance = (double: TableStub<unknown, unknown>) => ({
    stub: double.onGet,
    call: () => tableStubGet(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new TableStub<unknown, unknown>(seat)),
    withClock: (seat, clock) => instance(new TableStub<unknown, unknown>(seat, TableStub.withClock<unknown, unknown>(clock))),
    withLatency: (seat, duration, clock) => instance(TableStub.newDelayed<unknown, unknown>(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new TableStub<unknown, unknown>(seat, TableStub.withRandSource<unknown, unknown>(random))),
    benchMode: (seat) => instance(new TableStub<unknown, unknown>(seat, TableStub.benchMode<unknown, unknown>())),
    strict: (seat) => instance(new TableStub<unknown, unknown>(seat, TableStub.strict<unknown, unknown>())),
    async: true,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes get to the implementation of delegateTo", async ({ seat }) => {
    const inner = new TableStub<unknown, unknown>(seat);
    await settle(() => tableStubGet(new TableStub<unknown, unknown>(seat, TableStub.delegateTo<unknown, unknown>(inner))));
    inner.onGet.assertCalledOnce(seat, "get reaches the implementation");
  });

  test("fails get with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("get-delegate");
    const inner = new TableStub<unknown, unknown>(seat);
    inner.onGet.faults(want, 1);
    const double = new TableStub<unknown, unknown>(seat, TableStub.delegateTo<unknown, unknown>(inner));
    const outcome = await settle(() => tableStubGet(double));
    const message = "get fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes size to the implementation of delegateTo", async ({ seat }) => {
    const inner = new TableStub<unknown, unknown>(seat);
    await settle(() => tableStubSize(new TableStub<unknown, unknown>(seat, TableStub.delegateTo<unknown, unknown>(inner))));
    inner.onSize.assertCalledOnce(seat, "size reaches the implementation");
  });

  test("fails size with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("size-delegate");
    const inner = new TableStub<unknown, unknown>(seat);
    inner.onSize.faults(want, 1);
    const double = new TableStub<unknown, unknown>(seat, TableStub.delegateTo<unknown, unknown>(inner));
    const outcome = await settle(() => tableStubSize(double));
    const message = "size fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
