import { Recorder, Standard, byIdentity } from "@dokimi/assert";
import { behaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, resetBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { Invoice } from "./billing";
import { Account, Note } from "./built";
import { StoreStub } from "./built.stub";

/** Returns the samples of the arguments of get. An argument without a sample is undefined. */
function storeStubGetArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "get-a";
  return [a0];
}

/** Calls get of double with the samples of its arguments. */
function storeStubGet(double: StoreStub): Promise<Account> {
  const args = storeStubGetArgs();
  return double.get(args[0]);
}

/** Returns the sample of the result of get. */
function storeStubGetPin(): Account {
  let want0: Account;
  want0 = Object.assign(new Account(), { id: "id-a" });
  return want0;
}

/** Returns the samples of the arguments of find. An argument without a sample is undefined. */
function storeStubFindArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "find-a";
  return [a0];
}

/** Calls find of double with the samples of its arguments. */
function storeStubFind(double: StoreStub): Account | undefined {
  const args = storeStubFindArgs();
  return double.find(args[0]);
}

/** Returns the sample of the result of find. */
function storeStubFindPin(): Account | undefined {
  let want0: Account | undefined;
  want0 = Object.assign(new Account(), { id: "id-a" });
  return want0;
}

/** Returns the samples of the arguments of note. An argument without a sample is undefined. */
function storeStubNoteArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "note-a";
  return [a0];
}

/** Calls note of double with the samples of its arguments. */
function storeStubNote(double: StoreStub): Note {
  const args = storeStubNoteArgs();
  return double.note(args[0]);
}

/** Returns the sample of the result of note. */
function storeStubNotePin(): Note {
  let want0: Note;
  want0 = Object.assign(new Note(), { text: "text-a" });
  return want0;
}

/** Returns the samples of the arguments of bill. An argument without a sample is undefined. */
function storeStubBillArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "bill-a";
  return [a0];
}

/** Calls bill of double with the samples of its arguments. */
function storeStubBill(double: StoreStub): Invoice {
  const args = storeStubBillArgs();
  return double.bill(args[0]);
}

/** Returns the sample of the result of bill. */
function storeStubBillPin(): Invoice {
  let want0: Invoice;
  want0 = Object.assign(new Invoice(), { number: "number-a" });
  return want0;
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
        default: () => AccountBuilder.create().build(),
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

  describe("find", () => {
    const checks = behaviour("Store.find", (seat) => {
      const double = new StoreStub(seat);
      return {
        stub: double.onFind,
        call: () => storeStubFind(double),
        override: (mark) => {
          double.onFind.func(() => {
            mark();
            return undefined as never;
          });
        },
        default: () => AccountBuilder.create().build(),
      };
    });
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new StoreStub(undefined);
      storeStubFind(double);
      expect(double.onFind.callCount).isEqualTo(1, "find records the call");
      const got = double.onFind.calls[0];
      const args = storeStubFindArgs();
      expect(got?.id).isEqualTo(args[0], "the call records id");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onFind.onRecord(() => {
        seen += 1;
      });
      storeStubFind(double);
      storeStubFind(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of find");
    });
  });

  describe("note", () => {
    const checks = behaviour("Store.note", (seat) => {
      const double = new StoreStub(seat);
      return {
        stub: double.onNote,
        call: () => storeStubNote(double),
        override: (mark) => {
          double.onNote.func(() => {
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
      storeStubNote(double);
      expect(double.onNote.callCount).isEqualTo(1, "note records the call");
      const got = double.onNote.calls[0];
      const args = storeStubNoteArgs();
      expect(got?.id).isEqualTo(args[0], "the call records id");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onNote.onRecord(() => {
        seen += 1;
      });
      storeStubNote(double);
      storeStubNote(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of note");
    });
  });

  describe("bill", () => {
    const checks = behaviour("Store.bill", (seat) => {
      const double = new StoreStub(seat);
      return {
        stub: double.onBill,
        call: () => storeStubBill(double),
        override: (mark) => {
          double.onBill.func(() => {
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
      storeStubBill(double);
      expect(double.onBill.callCount).isEqualTo(1, "bill records the call");
      const got = double.onBill.calls[0];
      const args = storeStubBillArgs();
      expect(got?.id).isEqualTo(args[0], "the call records id");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new StoreStub(undefined);
      let seen = 0;
      double.onBill.onRecord(() => {
        seen += 1;
      });
      storeStubBill(double);
      storeStubBill(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of bill");
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

    it("clears the recorded calls of find", async () => {
      const double = new StoreStub(undefined);
      storeStubFind(double);
      double.resetCalls();
      expect(double.onFind.callCount).isEqualTo(0, "resetCalls clears the calls of find");
    });

    it("keeps the result of returns of find", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubFindPin();
      double.onFind.returns(want);
      double.resetCalls();
      expect(storeStubFind(double)).isEqualTo(want, "resetCalls keeps the result of find");
    });

    it("clears the recorded calls of note", async () => {
      const double = new StoreStub(undefined);
      storeStubNote(double);
      double.resetCalls();
      expect(double.onNote.callCount).isEqualTo(0, "resetCalls clears the calls of note");
    });

    it("keeps the result of returns of note", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubNotePin();
      double.onNote.returns(want);
      double.resetCalls();
      expect(storeStubNote(double)).isEqualTo(want, "resetCalls keeps the result of note");
    });

    it("clears the recorded calls of bill", async () => {
      const double = new StoreStub(undefined);
      storeStubBill(double);
      double.resetCalls();
      expect(double.onBill.callCount).isEqualTo(0, "resetCalls clears the calls of bill");
    });

    it("keeps the result of returns of bill", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubBillPin();
      double.onBill.returns(want);
      double.resetCalls();
      expect(storeStubBill(double)).isEqualTo(want, "resetCalls keeps the result of bill");
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

    it("fails find with its error", async () => {
      const want = new Error("find-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubFind(double));
      const message = "find fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails note with its error", async () => {
      const want = new Error("note-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubNote(double));
      const message = "note fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails bill with its error", async () => {
      const want = new Error("bill-failing");
      const double = StoreStub.newFailing(undefined, want);
      const outcome = await settle(() => storeStubBill(double));
      const message = "bill fails with the error of newFailing";
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

    it("passes find to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubFind(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onFind.callCount).isEqualTo(1, "find calls the implementation");
    });

    it("fails find with the error of the implementation", async () => {
      const want = new Error("find-delegate");
      const inner = new StoreStub(undefined);
      inner.onFind.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubFind(double));
      const message = "find fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes note to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubNote(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onNote.callCount).isEqualTo(1, "note calls the implementation");
    });

    it("fails note with the error of the implementation", async () => {
      const want = new Error("note-delegate");
      const inner = new StoreStub(undefined);
      inner.onNote.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubNote(double));
      const message = "note fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes bill to the implementation", async () => {
      const inner = new StoreStub(undefined);
      await settle(() => storeStubBill(new StoreStub(undefined, StoreStub.delegateTo(inner))));
      expect(inner.onBill.callCount).isEqualTo(1, "bill calls the implementation");
    });

    it("fails bill with the error of the implementation", async () => {
      const want = new Error("bill-delegate");
      const inner = new StoreStub(undefined);
      inner.onBill.faults(want, 1);
      const double = new StoreStub(undefined, StoreStub.delegateTo(inner));
      const outcome = await settle(() => storeStubBill(double));
      const message = "bill fails with the error of the implementation";
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

    it("fails find with its error", async () => {
      const want = new Error("find-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubFind(double));
      const message = "find fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails note with its error", async () => {
      const want = new Error("note-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubNote(double));
      const message = "note fails with the error of StoreStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails bill with its error", async () => {
      const want = new Error("bill-failing");
      const double = new StoreStub(undefined, StoreStub.failing(want));
      const outcome = await settle(() => storeStubBill(double));
      const message = "bill fails with the error of StoreStub.failing";
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

  describe("withStoreFind", () => {
    it("sets the function of find", async () => {
      let called = false;
      const option = StoreStub.withStoreFind(() => {
        called = true;
        return undefined as never;
      });
      storeStubFind(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreFind sets the function of find");
    });
  });

  describe("withStoreNote", () => {
    it("sets the function of note", async () => {
      let called = false;
      const option = StoreStub.withStoreNote(() => {
        called = true;
        return undefined as never;
      });
      storeStubNote(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreNote sets the function of note");
    });
  });

  describe("withStoreBill", () => {
    it("sets the function of bill", async () => {
      let called = false;
      const option = StoreStub.withStoreBill(() => {
        called = true;
        return undefined as never;
      });
      storeStubBill(new StoreStub(undefined, option));
      expect(called).isTrue("withStoreBill sets the function of bill");
    });
  });
});

describe("StoreGetStub", () => {
  describe("returns", () => {
    it("sets the result of each call of get", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubGetPin();
      const unset = AccountBuilder.create().build();
      expect(want).isNotEqualTo(unset, "the sample of the result differs from its default result");
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
});

describe("StoreFindStub", () => {
  describe("returns", () => {
    it("sets the result of each call of find", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubFindPin();
      const unset = AccountBuilder.create().build();
      expect(want).isNotEqualTo(unset, "the sample of the result differs from its default result");
      double.onFind.returns(want);
      expect(storeStubFind(double)).isEqualTo(want, "find returns the result that returns sets");
    });

    it("configures find for strict mode", async () => {
      const failures = new Recorder();
      const double = new StoreStub(failures);
      double.onFind.strict();
      double.onFind.returns(undefined as never);
      await settle(() => storeStubFind(double));
      expect(failures.failed).isFalse("strict mode accepts the call of find");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of find to the function", async () => {
      const double = new StoreStub(undefined);
      let got: unknown[] = [];
      double.onFind.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      storeStubFind(double);
      const args = storeStubFindArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});

describe("StoreNoteStub", () => {
  describe("returns", () => {
    it("sets the result of each call of note", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubNotePin();
      double.onNote.returns(want);
      expect(storeStubNote(double)).isEqualTo(want, "note returns the result that returns sets");
    });

    it("configures note for strict mode", async () => {
      const failures = new Recorder();
      const double = new StoreStub(failures);
      double.onNote.strict();
      double.onNote.returns(undefined as never);
      await settle(() => storeStubNote(double));
      expect(failures.failed).isFalse("strict mode accepts the call of note");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of note to the function", async () => {
      const double = new StoreStub(undefined);
      let got: unknown[] = [];
      double.onNote.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      storeStubNote(double);
      const args = storeStubNoteArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});

describe("StoreBillStub", () => {
  describe("returns", () => {
    it("sets the result of each call of bill", async () => {
      const double = new StoreStub(undefined);
      const want = storeStubBillPin();
      double.onBill.returns(want);
      expect(storeStubBill(double)).isEqualTo(want, "bill returns the result that returns sets");
    });

    it("configures bill for strict mode", async () => {
      const failures = new Recorder();
      const double = new StoreStub(failures);
      double.onBill.strict();
      double.onBill.returns(undefined as never);
      await settle(() => storeStubBill(double));
      expect(failures.failed).isFalse("strict mode accepts the call of bill");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of bill to the function", async () => {
      const double = new StoreStub(undefined);
      let got: unknown[] = [];
      double.onBill.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      storeStubBill(double);
      const args = storeStubBillArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});
