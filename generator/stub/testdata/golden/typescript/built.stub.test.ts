import { type Seat, Standard, byIdentity, check } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { Invoice } from "./billing";
import { Account, Note } from "./built";
import { AccountBuilder } from "./built.builder";
import { type StoreBillCall, type StoreFindCall, type StoreGetCall, type StoreNoteCall, StoreStub } from "./built.stub";

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

/** Returns the subject of behaviour and resetBehaviour for get, on a new double that reports to seat. */
function storeStubGetSubject(seat: Seat): Subject<StoreGetCall, Account | undefined> {
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
    default: () => AccountBuilder.create().build(),
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

/** Returns the subject of behaviour and resetBehaviour for find, on a new double that reports to seat. */
function storeStubFindSubject(seat: Seat): Subject<StoreFindCall, Account | undefined | undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onFind,
    call: () => storeStubFind(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = storeStubFindArgs();
      check.isEqualTo(s, c.id, args[0], "the call records id");
    },
    result: async () => storeStubFind(double),
    default: () => AccountBuilder.create().build(),
    returns: () => {
      const want = storeStubFindPin();
      double.onFind.returns(want);
      return want;
    },
    override: (seen) => {
      double.onFind.func((...received: unknown[]) => {
        seen({ id: received[0] } as StoreFindCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubFind(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
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

/** Returns the subject of behaviour and resetBehaviour for note, on a new double that reports to seat. */
function storeStubNoteSubject(seat: Seat): Subject<StoreNoteCall, Note | undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onNote,
    call: () => storeStubNote(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = storeStubNoteArgs();
      check.isEqualTo(s, c.id, args[0], "the call records id");
    },
    result: async () => storeStubNote(double),
    returns: () => {
      const want = storeStubNotePin();
      double.onNote.returns(want);
      return want;
    },
    override: (seen) => {
      double.onNote.func((...received: unknown[]) => {
        seen({ id: received[0] } as StoreNoteCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubNote(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
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

/** Returns the subject of behaviour and resetBehaviour for bill, on a new double that reports to seat. */
function storeStubBillSubject(seat: Seat): Subject<StoreBillCall, Invoice | undefined> {
  const double = new StoreStub(seat);
  return {
    stub: double.onBill,
    call: () => storeStubBill(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = storeStubBillArgs();
      check.isEqualTo(s, c.id, args[0], "the call records id");
    },
    result: async () => storeStubBill(double),
    returns: () => {
      const want = storeStubBillPin();
      double.onBill.returns(want);
      return want;
    },
    override: (seen) => {
      double.onBill.func((...received: unknown[]) => {
        seen({ id: received[0] } as StoreBillCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => storeStubBill(double));
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

  describe("find", () => {
    for (const c of behaviour("Store.find", storeStubFindSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("note", () => {
    for (const c of behaviour("Store.note", storeStubNoteSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("bill", () => {
    for (const c of behaviour("Store.bill", storeStubBillSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("get", () => {
      for (const c of resetBehaviour("Store.get", storeStubGetSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("find", () => {
      for (const c of resetBehaviour("Store.find", storeStubFindSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("note", () => {
      for (const c of resetBehaviour("Store.note", storeStubNoteSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("bill", () => {
      for (const c of resetBehaviour("Store.bill", storeStubBillSubject)) {
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
