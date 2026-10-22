import { Controlled, Recorder, Standard, byIdentity } from "@dokimi/assert";
import { behaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, resetBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { type Item, notFound } from "./store";
import { ItemStub, RepoStub, StoreStub, TableStub, WrapStub } from "./store.stub";

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

/** Calls close of double with the samples of its arguments. */
function storeStubClose(double: StoreStub): Promise<void> {
  return double.close();
}

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * get alone.
 */
function storeStubInstance(double: StoreStub) {
  return {
    stub: double.onGet,
    call: () => storeStubGet(double),
    reset: () => double.resetCalls(),
  };
}

describe("StoreStub", () => {
  describe("get", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new StoreStub(undefined);
      await storeStubGet(double);
      expect(double.onGet.callCount).isEqualTo(1, "get records the call");
      const got = double.onGet.calls[0];
      const args = storeStubGetArgs();
      expect(got?.id).isEqualTo(args[0], "the call records id");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onGet.onRecord(() => {
        seen += 1;
      });
      await storeStubGet(double);
      await storeStubGet(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of get");
    });
  });

  describe("put", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new StoreStub(undefined);
      storeStubPut(double);
      expect(double.onPut.callCount).isEqualTo(1, "put records the call");
      const got = double.onPut.calls[0];
      const args = storeStubPutArgs();
      expect.soft(got?.item).isEqualTo(args[0], "the call records item");
      expect.soft(got?.note).isEqualTo(args[1], "the call records note");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onPut.onRecord(() => {
        seen += 1;
      });
      storeStubPut(double);
      storeStubPut(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of put");
    });
  });

  describe("size", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onSize.onRecord(() => {
        seen += 1;
      });
      storeStubSize(double);
      storeStubSize(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of size");
    });
  });

  describe("log", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new StoreStub(undefined);
      storeStubLog(double);
      expect(double.onLog.callCount).isEqualTo(1, "log records the call");
      const got = double.onLog.calls[0];
      const args = storeStubLogArgs();
      expect.soft(got?.format).isEqualTo(args[0], "the call records format");
      expect.soft(got?.args).isEqualTo(args[1], "the call records args");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onLog.onRecord(() => {
        seen += 1;
      });
      storeStubLog(double);
      storeStubLog(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of log");
    });
  });

  describe("close", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onClose.onRecord(() => {
        seen += 1;
      });
      await storeStubClose(double);
      await storeStubClose(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of close");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => storeStubInstance(new StoreStub(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of get", async () => {
      const double = new StoreStub(undefined);
      await storeStubGet(double);
      double.resetCalls();
      expect(double.onGet.callCount).isEqualTo(0, "resetCalls clears the calls of get");
    });

    it("keeps the result of returns of get", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubGetPin();
      double.onGet.returns(want);
      double.resetCalls();
      expect(await storeStubGet(double)).isEqualTo(want, "resetCalls keeps the result of get");
    });

    it("clears the recorded calls of put", async () => {
      const double = new StoreStub(undefined);
      storeStubPut(double);
      double.resetCalls();
      expect(double.onPut.callCount).isEqualTo(0, "resetCalls clears the calls of put");
    });

    it("clears the recorded calls of size", async () => {
      const double = new StoreStub(undefined);
      storeStubSize(double);
      double.resetCalls();
      expect(double.onSize.callCount).isEqualTo(0, "resetCalls clears the calls of size");
    });

    it("keeps the result of returns of size", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubSizePin();
      double.onSize.returns(want);
      double.resetCalls();
      expect(storeStubSize(double)).isEqualTo(want, "resetCalls keeps the result of size");
    });

    it("clears the recorded calls of log", async () => {
      const double = new StoreStub(undefined);
      storeStubLog(double);
      double.resetCalls();
      expect(double.onLog.callCount).isEqualTo(0, "resetCalls clears the calls of log");
    });

    it("clears the recorded calls of close", async () => {
      const double = new StoreStub(undefined);
      await storeStubClose(double);
      double.resetCalls();
      expect(double.onClose.callCount).isEqualTo(0, "resetCalls clears the calls of close");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => storeStubInstance(StoreStub.newDelayed(seat, duration, clock)),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails get with its error", async () => {
      const want = new Error("get-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubGet(double));
      const message = "get fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails put with its error", async () => {
      const want = new Error("put-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubPut(double));
      const message = "put fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails size with its error", async () => {
      const want = new Error("size-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubSize(double));
      const message = "size fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails log with its error", async () => {
      const want = new Error("log-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubLog(double));
      const message = "log fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails close with its error", async () => {
      const want = new Error("close-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubClose(double));
      const message = "close fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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
  describe("returns", () => {
    it("sets the result of each call of get", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubGetPin();
      double.onGet.returns(want);
      expect(await storeStubGet(double)).isEqualTo(want, "get returns the result that returns sets");
    });

    it("configures get for strict mode", async () => {
      const failures = new Recorder();
      const double = new StoreStub(failures);
      double.onGet.strict();
      double.onGet.returns(undefined as never);
      await settle(() => storeStubGet(double));
      expect(failures.failed).isFalse("strict mode accepts the call of get");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of get to the function", async () => {
      const double = new StoreStub(undefined);
      let got: unknown[] = [];
      double.onGet.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      await storeStubGet(double);
      const args = storeStubGetArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });

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

describe("StorePutStub", () => {
  describe("func", () => {
    it("passes the arguments of each call of put to the function", async () => {
      const double = new StoreStub(undefined);
      let got: unknown[] = [];
      double.onPut.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      storeStubPut(double);
      const args = storeStubPutArgs();
      expect(got).isEqualTo([args[0], args[1]], "the function receives the arguments of the call");
    });
  });
});

describe("StoreSizeStub", () => {
  describe("returns", () => {
    it("sets the result of each call of size", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubSizePin();
      double.onSize.returns(want);
      expect(storeStubSize(double)).isEqualTo(want, "size returns the result that returns sets");
    });

    it("configures size for strict mode", async () => {
      const failures = new Recorder();
      const double = new StoreStub(failures);
      double.onSize.strict();
      double.onSize.returns(undefined as never);
      await settle(() => storeStubSize(double));
      expect(failures.failed).isFalse("strict mode accepts the call of size");
    });
  });
});

describe("StoreLogStub", () => {
  describe("func", () => {
    it("passes the arguments of each call of log to the function", async () => {
      const double = new StoreStub(undefined);
      let got: unknown[] = [];
      double.onLog.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      storeStubLog(double);
      const args = storeStubLogArgs();
      expect(got).isEqualTo([args[0], ...args[1]], "the function receives the arguments of the call");
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
      {
        name: "keeps each property of a delayed double",
        giveDouble: () => ItemStub.newDelayed(undefined, 1, new Controlled()),
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

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * find alone.
 */
function repoStubInstance(double: RepoStub<unknown>) {
  return {
    stub: double.onFind,
    call: () => repoStubFind(double),
    reset: () => double.resetCalls(),
  };
}

describe("RepoStub", () => {
  describe("find", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new RepoStub<unknown>(undefined);
      repoStubFind(double);
      expect(double.onFind.callCount).isEqualTo(1, "find records the call");
      const got = double.onFind.calls[0];
      const args = repoStubFindArgs();
      expect(got?.id).isEqualTo(args[0], "the call records id");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new RepoStub<unknown>(undefined);
      let seen = 0;
      double.onFind.onRecord(() => {
        seen += 1;
      });
      repoStubFind(double);
      repoStubFind(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of find");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => repoStubInstance(new RepoStub<unknown>(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of find", async () => {
      const double = new RepoStub<unknown>(undefined);
      repoStubFind(double);
      double.resetCalls();
      expect(double.onFind.callCount).isEqualTo(0, "resetCalls clears the calls of find");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => repoStubInstance(RepoStub.newDelayed<unknown>(seat, duration, clock)),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails find with its error", async () => {
      const want = new Error("find-failing");
      const double = RepoStub.newFailing<unknown>(undefined, want);
      const outcome = await settle(() => repoStubFind(double));
      const message = "find fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("RepoFindStub", () => {
  describe("returns", () => {
    it("configures find for strict mode", async () => {
      const failures = new Recorder();
      const double = new RepoStub<unknown>(failures);
      double.onFind.strict();
      double.onFind.returns(undefined as never);
      await settle(() => repoStubFind(double));
      expect(failures.failed).isFalse("strict mode accepts the call of find");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of find to the function", async () => {
      const double = new RepoStub<unknown>(undefined);
      let got: unknown[] = [];
      double.onFind.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      repoStubFind(double);
      const args = repoStubFindArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
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

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * put alone.
 */
function wrapStubInstance(double: WrapStub<unknown>) {
  return {
    stub: double.onPut,
    call: () => wrapStubPut(double),
    reset: () => double.resetCalls(),
  };
}

describe("WrapStub", () => {
  describe("put", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new WrapStub<unknown>(undefined);
      wrapStubPut(double);
      expect(double.onPut.callCount).isEqualTo(1, "put records the call");
      const got = double.onPut.calls[0];
      const args = wrapStubPutArgs<unknown>();
      expect(got?.items).isEqualTo(args[0], "the call records items");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new WrapStub<unknown>(undefined);
      let seen = 0;
      double.onPut.onRecord(() => {
        seen += 1;
      });
      wrapStubPut(double);
      wrapStubPut(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of put");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => wrapStubInstance(new WrapStub<unknown>(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of put", async () => {
      const double = new WrapStub<unknown>(undefined);
      wrapStubPut(double);
      double.resetCalls();
      expect(double.onPut.callCount).isEqualTo(0, "resetCalls clears the calls of put");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => wrapStubInstance(WrapStub.newDelayed<unknown>(seat, duration, clock)),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails put with its error", async () => {
      const want = new Error("put-failing");
      const double = WrapStub.newFailing<unknown>(undefined, want);
      const outcome = await settle(() => wrapStubPut(double));
      const message = "put fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("WrapPutStub", () => {
  describe("func", () => {
    it("passes the arguments of each call of put to the function", async () => {
      const double = new WrapStub<unknown>(undefined);
      let got: unknown[] = [];
      double.onPut.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      wrapStubPut(double);
      const args = wrapStubPutArgs<unknown>();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
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

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * get alone.
 */
function tableStubInstance(double: TableStub<unknown, unknown>) {
  return {
    stub: double.onGet,
    call: () => tableStubGet(double),
    reset: () => double.resetCalls(),
  };
}

describe("TableStub", () => {
  describe("get", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      await tableStubGet(double);
      expect(double.onGet.callCount).isEqualTo(1, "get records the call");
      const got = double.onGet.calls[0];
      const args = tableStubGetArgs<unknown>();
      expect(got?.key).isEqualTo(args[0], "the call records key");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      let seen = 0;
      double.onGet.onRecord(() => {
        seen += 1;
      });
      await tableStubGet(double);
      await tableStubGet(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of get");
    });
  });

  describe("size", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      let seen = 0;
      double.onSize.onRecord(() => {
        seen += 1;
      });
      tableStubSize(double);
      tableStubSize(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of size");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => tableStubInstance(new TableStub<unknown, unknown>(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of get", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      await tableStubGet(double);
      double.resetCalls();
      expect(double.onGet.callCount).isEqualTo(0, "resetCalls clears the calls of get");
    });

    it("clears the recorded calls of size", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      tableStubSize(double);
      double.resetCalls();
      expect(double.onSize.callCount).isEqualTo(0, "resetCalls clears the calls of size");
    });

    it("keeps the result of returns of size", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      const want = tableStubSizePin();
      double.onSize.returns(want);
      double.resetCalls();
      expect(tableStubSize(double)).isEqualTo(want, "resetCalls keeps the result of size");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => tableStubInstance(TableStub.newDelayed<unknown, unknown>(seat, duration, clock)),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails get with its error", async () => {
      const want = new Error("get-failing");
      const double = TableStub.newFailing<unknown, unknown>(undefined, want);
      const outcome = await settle(() => tableStubGet(double));
      const message = "get fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails size with its error", async () => {
      const want = new Error("size-failing");
      const double = TableStub.newFailing<unknown, unknown>(undefined, want);
      const outcome = await settle(() => tableStubSize(double));
      const message = "size fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("TableGetStub", () => {
  describe("returns", () => {
    it("configures get for strict mode", async () => {
      const failures = new Recorder();
      const double = new TableStub<unknown, unknown>(failures);
      double.onGet.strict();
      double.onGet.returns(undefined as never);
      await settle(() => tableStubGet(double));
      expect(failures.failed).isFalse("strict mode accepts the call of get");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of get to the function", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      let got: unknown[] = [];
      double.onGet.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      await tableStubGet(double);
      const args = tableStubGetArgs<unknown>();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});

describe("TableSizeStub", () => {
  describe("returns", () => {
    it("sets the result of each call of size", async () => {
      const double = new TableStub<unknown, unknown>(undefined);
      const want = tableStubSizePin();
      double.onSize.returns(want);
      expect(tableStubSize(double)).isEqualTo(want, "size returns the result that returns sets");
    });

    it("configures size for strict mode", async () => {
      const failures = new Recorder();
      const double = new TableStub<unknown, unknown>(failures);
      double.onSize.strict();
      double.onSize.returns(undefined as never);
      await settle(() => tableStubSize(double));
      expect(failures.failed).isFalse("strict mode accepts the call of size");
    });
  });
});
