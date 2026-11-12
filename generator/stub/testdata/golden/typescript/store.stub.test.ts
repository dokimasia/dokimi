import { Controlled, type Seat, Standard, byIdentity, check } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { type Item, notFound } from "./store";
import { ItemStub, type RepoFindCall, RepoStub, type StoreCloseCall, type StoreGetCall, type StoreLogCall, type StorePutCall, type StoreSizeCall, StoreStub, type TableGetCall, type TableSizeCall, TableStub, type WrapPutCall, WrapStub } from "./store.stub";

/** Returns the samples of the arguments of get. An argument without a sample is undefined. */
function storeStubGetArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "get-a";
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
  want0 = { id: "id-a" };
  return want0;
}

/** Returns the subject of behaviour and resetBehaviour for get, on a new double that reports to seat. */
function storeStubGetSubject(seat: Seat): Subject<StoreGetCall, Item | undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onGet,
    call: () => storeStubGet(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = storeStubGetArgs();
      check.isEqualTo(s, c.id, args[0], "the call records id");
    },
    result: async () => await storeStubGet(double),
    returns: () => {
      const want = storeStubGetPin();
      double.onGet.returns(want);
      return want;
    },
    override: (seen) => {
      double.onGet.func((...received: unknown[]) => {
        seen({ id: received[0] } as StoreGetCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubGet(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Returns the samples of the arguments of put. An argument without a sample is undefined. */
function storeStubPutArgs(): [Item, string | undefined] {
  let a0 = undefined as unknown as Item;
  let a1 = undefined as unknown as string;
  a0 = { id: "id-a" };
  a1 = "put-a";
  return [a0, a1];
}

/** Calls put of double with the samples of its arguments. */
function storeStubPut(double: StoreStub): void {
  const args = storeStubPutArgs();
  return double.put(args[0], args[1]);
}

/** Returns the subject of behaviour and resetBehaviour for put, on a new double that reports to seat. */
function storeStubPutSubject(seat: Seat): Subject<StorePutCall, undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onPut,
    call: () => storeStubPut(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = storeStubPutArgs();
      check.isEqualTo(s, c.item, args[0], "the call records item");
      check.isEqualTo(s, c.note, args[1], "the call records note");
    },
    override: (seen) => {
      double.onPut.func((...received: unknown[]) => {
        seen({ item: received[0], note: received[1] } as StorePutCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubPut(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
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

/** Returns the subject of behaviour and resetBehaviour for size, on a new double that reports to seat. */
function storeStubSizeSubject(seat: Seat): Subject<StoreSizeCall, number | undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onSize,
    call: () => storeStubSize(double),
    reset: () => double.resetCalls(),
    result: async () => storeStubSize(double),
    returns: () => {
      const want = storeStubSizePin();
      double.onSize.returns(want);
      return want;
    },
    override: (seen) => {
      double.onSize.func(() => {
        seen({} as StoreSizeCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubSize(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Returns the samples of the arguments of log. An argument without a sample is undefined. */
function storeStubLogArgs(): [string, string[]] {
  let a0 = undefined as unknown as string;
  const a1: string[] = [];
  a0 = "log-a";
  return [a0, a1];
}

/** Calls log of double with the samples of its arguments. */
function storeStubLog(double: StoreStub): void {
  const args = storeStubLogArgs();
  return double.log(args[0], ...args[1]);
}

/** Returns the subject of behaviour and resetBehaviour for log, on a new double that reports to seat. */
function storeStubLogSubject(seat: Seat): Subject<StoreLogCall, undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onLog,
    call: () => storeStubLog(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = storeStubLogArgs();
      check.isEqualTo(s, c.format, args[0], "the call records format");
      check.isEqualTo(s, c.args, args[1], "the call records args");
    },
    override: (seen) => {
      double.onLog.func((...received: unknown[]) => {
        seen({ format: received[0], args: received.slice(1) } as StoreLogCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubLog(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Calls close of double with the samples of its arguments. */
function storeStubClose(double: StoreStub): Promise<void> {
  return double.close();
}

/** Returns the subject of behaviour and resetBehaviour for close, on a new double that reports to seat. */
function storeStubCloseSubject(seat: Seat): Subject<StoreCloseCall, undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onClose,
    call: () => storeStubClose(double),
    reset: () => double.resetCalls(),
    override: (seen) => {
      double.onClose.func(() => {
        seen({} as StoreCloseCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubClose(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call get alone.
 */
function storeStubInstance(double: StoreStub) {
  return {
    stub: double.onGet,
    call: () => storeStubGet(double),
  };
}

describe("StoreStub", () => {
  describe("get", () => {
    for (const c of behaviour("Store.get", storeStubGetSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("put", () => {
    for (const c of behaviour("Store.put", storeStubPutSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("size", () => {
    for (const c of behaviour("Store.size", storeStubSizeSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("log", () => {
    for (const c of behaviour("Store.log", storeStubLogSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("close", () => {
    for (const c of behaviour("Store.close", storeStubCloseSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("get", () => {
      for (const c of resetBehaviour("Store.get", storeStubGetSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("put", () => {
      for (const c of resetBehaviour("Store.put", storeStubPutSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("size", () => {
      for (const c of resetBehaviour("Store.size", storeStubSizeSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("log", () => {
      for (const c of resetBehaviour("Store.log", storeStubLogSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("close", () => {
      for (const c of resetBehaviour("Store.close", storeStubCloseSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => storeStubInstance(new StoreStub(seat, StoreStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes get to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubGet(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onGet.callCount).isEqualTo(1, "get calls the implementation");
    });

    it("fails get with the error of the implementation", async () => {
      const want = new Error("get-delegate");
      const inner = new StoreStub(undefined);
      inner.onGet.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubGet(double));
      const message = "get fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes put to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubPut(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onPut.callCount).isEqualTo(1, "put calls the implementation");
    });

    it("fails put with the error of the implementation", async () => {
      const want = new Error("put-delegate");
      const inner = new StoreStub(undefined);
      inner.onPut.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubPut(double));
      const message = "put fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes size to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubSize(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onSize.callCount).isEqualTo(1, "size calls the implementation");
    });

    it("fails size with the error of the implementation", async () => {
      const want = new Error("size-delegate");
      const inner = new StoreStub(undefined);
      inner.onSize.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubSize(double));
      const message = "size fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes log to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubLog(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onLog.callCount).isEqualTo(1, "log calls the implementation");
    });

    it("fails log with the error of the implementation", async () => {
      const want = new Error("log-delegate");
      const inner = new StoreStub(undefined);
      inner.onLog.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubLog(double));
      const message = "log fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes close to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubClose(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onClose.callCount).isEqualTo(1, "close calls the implementation");
    });

    it("fails close with the error of the implementation", async () => {
      const want = new Error("close-delegate");
      const inner = new StoreStub(undefined);
      inner.onClose.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubClose(double));
      const message = "close fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => storeStubInstance(new StoreStub(seat, StoreStub.withClock(clock))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      storeStubInstance(new StoreStub(seat, StoreStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => storeStubInstance(new StoreStub(seat, StoreStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        storeStubInstance(new StoreStub(seat, StoreStub.withClock(clock), StoreStub.withLatency(duration))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails get with its error", async () => {
      const want = new Error("get-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubGet(double));
      const message = "get fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails put with its error", async () => {
      const want = new Error("put-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubPut(double));
      const message = "put fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails size with its error", async () => {
      const want = new Error("size-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubSize(double));
      const message = "size fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails log with its error", async () => {
      const want = new Error("log-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubLog(double));
      const message = "log fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails close with its error", async () => {
      const want = new Error("close-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubClose(double));
      const message = "close fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withStoreGet", () => {
    it("sets the function of get", async () => {
      let called = false;
      const option = StoreStub.withStoreGet(() => {
        called = true;
        return undefined as never;
      });
      await storeStubGet(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreGet sets the function of get");
    });
  });

  describe("withStorePut", () => {
    it("sets the function of put", async () => {
      let called = false;
      const option = StoreStub.withStorePut(() => {
        called = true;
        return undefined as never;
      });
      storeStubPut(new StoreStub(undefined, option));
      expect(called).isTrue("withStorePut sets the function of put");
    });
  });

  describe("withStoreSize", () => {
    it("sets the function of size", async () => {
      let called = false;
      const option = StoreStub.withStoreSize(() => {
        called = true;
        return undefined as never;
      });
      storeStubSize(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreSize sets the function of size");
    });
  });

  describe("withStoreLog", () => {
    it("sets the function of log", async () => {
      let called = false;
      const option = StoreStub.withStoreLog(() => {
        called = true;
        return undefined as never;
      });
      storeStubLog(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreLog sets the function of log");
    });
  });

  describe("withStoreClose", () => {
    it("sets the function of close", async () => {
      let called = false;
      const option = StoreStub.withStoreClose(() => {
        called = true;
        return undefined as never;
      });
      await storeStubClose(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreClose sets the function of close");
    });
  });
});

describe("StoreGetStub", () => {
  describe("faultNotFound", () => {
    it("fails the next call of get with notFound", async () => {
      const double = new StoreStub(undefined);
      double.onGet.faultNotFound();
      const outcome = await settle(() => storeStubGet(double));
      const message = "the next call of get fails with notFound";
      expect(outcome.threw && outcome.error).isEqualTo(notFound, message, byIdentity());
    });

    it("spares the call after the failed call", async () => {
      const double = new StoreStub(undefined);
      double.onGet.faultNotFound();
      await settle(() => storeStubGet(double));
      const outcome = await settle(() => storeStubGet(double));
      expect(outcome.threw).isFalse("the call of get after the failed call succeeds");
    });
  });

  describe("retrySchedule", () => {
    it("fails each call of get before call 3", async () => {
      const double = new StoreStub(undefined);
      const want = new Error("get-retry");
      double.onGet.retrySchedule(want);
      for (let attempt = 1; attempt < 3; attempt += 1) {
        const outcome = await settle(() => storeStubGet(double));
        const message = "a call of get before call 3 fails";
        expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
      }
    });

    it("spares call 3 of get", async () => {
      const double = new StoreStub(undefined);
      double.onGet.retrySchedule(new Error("get-retry"));
      for (let attempt = 1; attempt < 3; attempt += 1) {
        await settle(() => storeStubGet(double));
      }
      const outcome = await settle(() => storeStubGet(double));
      expect(outcome.threw).isFalse("call 3 of get succeeds");
    });
  });

  describe("faultForPartition", () => {
    it("fails a call of get in the partition of the key", async () => {
      const double = new StoreStub(undefined);
      const args = storeStubGetArgs();
      const want = new Error("get-partition");
      double.onGet.faultForPartition(args[0], want, 1);
      const outcome = await settle(() => storeStubGet(double));
      const message = "a call in the partition of the key fails";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("faultForOtherPartitions", () => {
    it("spares a call of get in the partition of the key", async () => {
      const double = new StoreStub(undefined);
      const args = storeStubGetArgs();
      double.onGet.faultForOtherPartitions(args[0], new Error("get-partition"), 1);
      const outcome = await settle(() => storeStubGet(double));
      expect(outcome.threw).isFalse("a call in the partition of the key succeeds");
    });
  });
});

/** Assigns the sample of each property of ItemStub to double. */
function itemStubProps(double: ItemStub): void {
  double.id = "id-a";
}

describe("ItemStub", () => {
  describe("resetCalls", () => {
    const tests: { name: string; giveDouble: () => ItemStub }[] = [
      {
        name: "keeps each property of a new double",
        giveDouble: () => new ItemStub(undefined),
      },
      {
        name: "keeps each property of a strict double",
        giveDouble: () => new ItemStub(undefined, ItemStub.strict()),
      },
      {
        name: "keeps each property of a double with a clock",
        giveDouble: () => new ItemStub(undefined, ItemStub.withClock(new Controlled())),
      },
      {
        name: "keeps each property of a double with a random source",
        giveDouble: () => new ItemStub(undefined, ItemStub.withRandSource(() => 0)),
      },
      {
        name: "keeps each property of a double in bench mode",
        giveDouble: () => new ItemStub(undefined, ItemStub.benchMode()),
      },
      {
        name: "keeps each property of a double with a latency",
        giveDouble: () => new ItemStub(undefined, ItemStub.withLatency(1)),
      },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        const double = tt.giveDouble();
        itemStubProps(double);
        double.resetCalls();
        const want = new ItemStub(undefined);
        itemStubProps(want);
        expect(double.id).isEqualTo(want.id, "resetCalls keeps the property id");
      });
    }
  });
});

/** Returns the samples of the arguments of find. An argument without a sample is undefined. */
function repoStubFindArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "find-a";
  return [a0];
}

/** Calls find of double with the samples of its arguments. */
function repoStubFind<T>(double: RepoStub<T>): T {
  const args = repoStubFindArgs();
  return double.find(args[0]);
}

/** Returns the subject of behaviour and resetBehaviour for find, on a new double that reports to seat. */
function repoStubFindSubject<T>(seat: Seat): Subject<RepoFindCall<T>, T | undefined> {
  const double = new RepoStub<T>(seat);
  return {
    stub: double.onFind,
    call: () => repoStubFind(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = repoStubFindArgs();
      check.isEqualTo(s, c.id, args[0], "the call records id");
    },
    result: async () => repoStubFind(double),
    override: (seen) => {
      double.onFind.func((...received: unknown[]) => {
        seen({ id: received[0] } as RepoFindCall<T>);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => repoStubFind(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call find alone.
 */
function repoStubInstance(double: RepoStub<unknown>) {
  return {
    stub: double.onFind,
    call: () => repoStubFind(double),
  };
}

describe("RepoStub", () => {
  describe("find", () => {
    for (const c of behaviour("Repo.find", repoStubFindSubject<unknown>)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("find", () => {
      for (const c of resetBehaviour("Repo.find", repoStubFindSubject<unknown>)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => repoStubInstance(new RepoStub<unknown>(seat, RepoStub.strict<unknown>())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes find to the implementation", async () => {
      const inner = new RepoStub<unknown>(undefined);
      await settle(() => repoStubFind(new RepoStub<unknown>(undefined, RepoStub.delegateTo<unknown>(inner))));
      expect(inner.onFind.callCount).isEqualTo(1, "find calls the implementation");
    });

    it("fails find with the error of the implementation", async () => {
      const want = new Error("find-delegate");
      const inner = new RepoStub<unknown>(undefined);
      inner.onFind.faults(want, 1);
      const double = new RepoStub<unknown>(undefined, RepoStub.delegateTo<unknown>(inner));
      const outcome = await settle(() => repoStubFind(double));
      const message = "find fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => repoStubInstance(new RepoStub<unknown>(seat, RepoStub.withClock<unknown>(clock))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      repoStubInstance(new RepoStub<unknown>(seat, RepoStub.withRandSource<unknown>(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => repoStubInstance(new RepoStub<unknown>(seat, RepoStub.benchMode<unknown>())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        repoStubInstance(new RepoStub<unknown>(seat, RepoStub.withClock<unknown>(clock), RepoStub.withLatency<unknown>(duration))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails find with its error", async () => {
      const want = new Error("find-failing");
      const double = new RepoStub<unknown>(undefined, RepoStub.failing<unknown>(want));
      const outcome = await settle(() => repoStubFind(double));
      const message = "find fails with the error of RepoStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withRepoFind", () => {
    it("sets the function of find", async () => {
      let called = false;
      const option = RepoStub.withRepoFind<unknown>(() => {
        called = true;
        return undefined as never;
      });
      repoStubFind(new RepoStub<unknown>(undefined, option));
      expect(called).isTrue("withRepoFind sets the function of find");
    });
  });
});

/** Returns the samples of the arguments of put. An argument without a sample is undefined. */
function wrapStubPutArgs<T$>(): [T$[]] {
  let a0 = undefined as unknown as T$[];
  return [a0];
}

/** Calls put of double with the samples of its arguments. */
function wrapStubPut<T$>(double: WrapStub<T$>): void {
  const args = wrapStubPutArgs<T$>();
  return double.put(args[0]);
}

/** Returns the subject of behaviour and resetBehaviour for put, on a new double that reports to seat. */
function wrapStubPutSubject<T$>(seat: Seat): Subject<WrapPutCall<T$>, undefined> {
  const double = new WrapStub<T$>(seat);
  return {
    stub: double.onPut,
    call: () => wrapStubPut(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = wrapStubPutArgs<T$>();
      check.isEqualTo(s, c.items, args[0], "the call records items");
    },
    override: (seen) => {
      double.onPut.func((...received: unknown[]) => {
        seen({ items: received[0] } as WrapPutCall<T$>);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => wrapStubPut(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call put alone.
 */
function wrapStubInstance(double: WrapStub<unknown>) {
  return {
    stub: double.onPut,
    call: () => wrapStubPut(double),
  };
}

describe("WrapStub", () => {
  describe("put", () => {
    for (const c of behaviour("Wrap.put", wrapStubPutSubject<unknown>)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("put", () => {
      for (const c of resetBehaviour("Wrap.put", wrapStubPutSubject<unknown>)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => wrapStubInstance(new WrapStub<unknown>(seat, WrapStub.strict<unknown>())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes put to the implementation", async () => {
      const inner = new WrapStub<unknown>(undefined);
      await settle(() => wrapStubPut(new WrapStub<unknown>(undefined, WrapStub.delegateTo<unknown>(inner))));
      expect(inner.onPut.callCount).isEqualTo(1, "put calls the implementation");
    });

    it("fails put with the error of the implementation", async () => {
      const want = new Error("put-delegate");
      const inner = new WrapStub<unknown>(undefined);
      inner.onPut.faults(want, 1);
      const double = new WrapStub<unknown>(undefined, WrapStub.delegateTo<unknown>(inner));
      const outcome = await settle(() => wrapStubPut(double));
      const message = "put fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => wrapStubInstance(new WrapStub<unknown>(seat, WrapStub.withClock<unknown>(clock))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      wrapStubInstance(new WrapStub<unknown>(seat, WrapStub.withRandSource<unknown>(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => wrapStubInstance(new WrapStub<unknown>(seat, WrapStub.benchMode<unknown>())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        wrapStubInstance(new WrapStub<unknown>(seat, WrapStub.withClock<unknown>(clock), WrapStub.withLatency<unknown>(duration))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails put with its error", async () => {
      const want = new Error("put-failing");
      const double = new WrapStub<unknown>(undefined, WrapStub.failing<unknown>(want));
      const outcome = await settle(() => wrapStubPut(double));
      const message = "put fails with the error of WrapStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withWrapPut", () => {
    it("sets the function of put", async () => {
      let called = false;
      const option = WrapStub.withWrapPut<unknown>(() => {
        called = true;
        return undefined as never;
      });
      wrapStubPut(new WrapStub<unknown>(undefined, option));
      expect(called).isTrue("withWrapPut sets the function of put");
    });
  });
});

/** Returns the samples of the arguments of get. An argument without a sample is undefined. */
function tableStubGetArgs<K>(): [K] {
  let a0 = undefined as unknown as K;
  return [a0];
}

/** Calls get of double with the samples of its arguments. */
function tableStubGet<K, V>(double: TableStub<K, V>): Promise<V> {
  const args = tableStubGetArgs<K>();
  return double.get(args[0]);
}

/** Returns the subject of behaviour and resetBehaviour for get, on a new double that reports to seat. */
function tableStubGetSubject<K, V>(seat: Seat): Subject<TableGetCall<K, V>, V | undefined> {
  const double = new TableStub<K, V>(seat);
  return {
    stub: double.onGet,
    call: () => tableStubGet(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = tableStubGetArgs<K>();
      check.isEqualTo(s, c.key, args[0], "the call records key");
    },
    result: async () => await tableStubGet(double),
    override: (seen) => {
      double.onGet.func((...received: unknown[]) => {
        seen({ key: received[0] } as TableGetCall<K, V>);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => tableStubGet(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
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

/** Returns the subject of behaviour and resetBehaviour for size, on a new double that reports to seat. */
function tableStubSizeSubject<K, V>(seat: Seat): Subject<TableSizeCall, number | undefined> {
  const double = new TableStub<K, V>(seat);
  return {
    stub: double.onSize,
    call: () => tableStubSize(double),
    reset: () => double.resetCalls(),
    result: async () => tableStubSize(double),
    returns: () => {
      const want = tableStubSizePin();
      double.onSize.returns(want);
      return want;
    },
    override: (seen) => {
      double.onSize.func(() => {
        seen({} as TableSizeCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => tableStubSize(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call get alone.
 */
function tableStubInstance(double: TableStub<unknown, unknown>) {
  return {
    stub: double.onGet,
    call: () => tableStubGet(double),
  };
}

describe("TableStub", () => {
  describe("get", () => {
    for (const c of behaviour("Table.get", tableStubGetSubject<unknown, unknown>)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("size", () => {
    for (const c of behaviour("Table.size", tableStubSizeSubject<unknown, unknown>)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("get", () => {
      for (const c of resetBehaviour("Table.get", tableStubGetSubject<unknown, unknown>)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("size", () => {
      for (const c of resetBehaviour("Table.size", tableStubSizeSubject<unknown, unknown>)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => tableStubInstance(new TableStub<unknown, unknown>(seat, TableStub.strict<unknown, unknown>())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes get to the implementation", async () => {
      const inner = new TableStub<unknown, unknown>(undefined);
      await settle(() => tableStubGet(new TableStub<unknown, unknown>(undefined, TableStub.delegateTo<unknown, unknown>(inner))));
      expect(inner.onGet.callCount).isEqualTo(1, "get calls the implementation");
    });

    it("fails get with the error of the implementation", async () => {
      const want = new Error("get-delegate");
      const inner = new TableStub<unknown, unknown>(undefined);
      inner.onGet.faults(want, 1);
      const double = new TableStub<unknown, unknown>(undefined, TableStub.delegateTo<unknown, unknown>(inner));
      const outcome = await settle(() => tableStubGet(double));
      const message = "get fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes size to the implementation", async () => {
      const inner = new TableStub<unknown, unknown>(undefined);
      await settle(() => tableStubSize(new TableStub<unknown, unknown>(undefined, TableStub.delegateTo<unknown, unknown>(inner))));
      expect(inner.onSize.callCount).isEqualTo(1, "size calls the implementation");
    });

    it("fails size with the error of the implementation", async () => {
      const want = new Error("size-delegate");
      const inner = new TableStub<unknown, unknown>(undefined);
      inner.onSize.faults(want, 1);
      const double = new TableStub<unknown, unknown>(undefined, TableStub.delegateTo<unknown, unknown>(inner));
      const outcome = await settle(() => tableStubSize(double));
      const message = "size fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => tableStubInstance(new TableStub<unknown, unknown>(seat, TableStub.withClock<unknown, unknown>(clock))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      tableStubInstance(new TableStub<unknown, unknown>(seat, TableStub.withRandSource<unknown, unknown>(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => tableStubInstance(new TableStub<unknown, unknown>(seat, TableStub.benchMode<unknown, unknown>())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        tableStubInstance(new TableStub<unknown, unknown>(seat, TableStub.withClock<unknown, unknown>(clock), TableStub.withLatency<unknown, unknown>(duration))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails get with its error", async () => {
      const want = new Error("get-failing");
      const double = new TableStub<unknown, unknown>(undefined, TableStub.failing<unknown, unknown>(want));
      const outcome = await settle(() => tableStubGet(double));
      const message = "get fails with the error of TableStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails size with its error", async () => {
      const want = new Error("size-failing");
      const double = new TableStub<unknown, unknown>(undefined, TableStub.failing<unknown, unknown>(want));
      const outcome = await settle(() => tableStubSize(double));
      const message = "size fails with the error of TableStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withTableGet", () => {
    it("sets the function of get", async () => {
      let called = false;
      const option = TableStub.withTableGet<unknown, unknown>(() => {
        called = true;
        return undefined as never;
      });
      await tableStubGet(new TableStub<unknown, unknown>(undefined, option));
      expect(called).isTrue("withTableGet sets the function of get");
    });
  });

  describe("withTableSize", () => {
    it("sets the function of size", async () => {
      let called = false;
      const option = TableStub.withTableSize<unknown, unknown>(() => {
        called = true;
        return undefined as never;
      });
      tableStubSize(new TableStub<unknown, unknown>(undefined, option));
      expect(called).isTrue("withTableSize sets the function of size");
    });
  });
});
