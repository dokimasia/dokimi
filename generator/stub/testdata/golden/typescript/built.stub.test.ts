import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import type { Invoice } from './billing';
import type { Account, Note } from './built';
import { StoreStub } from './built.stub';

/** Returns the samples of the arguments of get, or undefined. */
function storeStubGetArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'get-a';
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
  want0 = {id: 'id-a'} as Account;
  return want0;
}

/** Returns the samples of the arguments of find, or undefined. */
function storeStubFindArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'find-a';
  return [a0];
}

/** Calls find of double with the samples of its arguments. */
function storeStubFind(double: StoreStub): Account|undefined {
  const args = storeStubFindArgs();
  return double.find(args[0]);
}

/** Returns the sample of the result of find. */
function storeStubFindPin(): Account|undefined {
  let want0: Account|undefined;
  want0 = {id: 'id-a'} as Account;
  return want0;
}

/** Returns the samples of the arguments of note, or undefined. */
function storeStubNoteArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'note-a';
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
  want0 = {text: 'text-a'} as Note;
  return want0;
}

/** Returns the samples of the arguments of bill, or undefined. */
function storeStubBillArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'bill-a';
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
  want0 = {number: 'number-a'} as Invoice;
  return want0;
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
      default: () => AccountBuilder.create().build(),
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubGetPin();
    const unset = AccountBuilder.create().build();
    check.notEqual(seat, want, unset, "the sample of the result differs from its default result");
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

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubGetPin();
    double.onGet.returns(want);
    await storeStubGet(double);
    double.resetCalls();
    check.equal(seat, await storeStubGet(double), want, "a reset keeps the result of returns");
  });
});

describe("StoreStub.find", () => {
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
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubFindPin();
    const unset = AccountBuilder.create().build();
    check.notEqual(seat, want, unset, "the sample of the result differs from its default result");
    double.onFind.returns(want);
    check.equal(seat, storeStubFind(double), want, "find returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    storeStubFind(double);
    const got = double.onFind.assertCalledOnce(seat, "find records the call");
    const args = storeStubFindArgs();
    check.equal(seat, got?.id, args[0], "the call records id");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onFind.onRecord(() => {
      seen += 1;
    });
    storeStubFind(double);
    storeStubFind(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of find");
  });

  test("sets the function of withStoreFind at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreFind(() => {
      called = true;
      return undefined as never;
    });
    storeStubFind(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreFind sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("find-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubFind(double));
    const message = "find fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubFindPin();
    double.onFind.returns(want);
    storeStubFind(double);
    double.resetCalls();
    check.equal(seat, storeStubFind(double), want, "a reset keeps the result of returns");
  });
});

describe("StoreStub.note", () => {
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
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubNotePin();
    double.onNote.returns(want);
    check.equal(seat, storeStubNote(double), want, "note returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    storeStubNote(double);
    const got = double.onNote.assertCalledOnce(seat, "note records the call");
    const args = storeStubNoteArgs();
    check.equal(seat, got?.id, args[0], "the call records id");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onNote.onRecord(() => {
      seen += 1;
    });
    storeStubNote(double);
    storeStubNote(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of note");
  });

  test("sets the function of withStoreNote at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreNote(() => {
      called = true;
      return undefined as never;
    });
    storeStubNote(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreNote sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("note-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubNote(double));
    const message = "note fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubNotePin();
    double.onNote.returns(want);
    storeStubNote(double);
    double.resetCalls();
    check.equal(seat, storeStubNote(double), want, "a reset keeps the result of returns");
  });
});

describe("StoreStub.bill", () => {
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
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubBillPin();
    double.onBill.returns(want);
    check.equal(seat, storeStubBill(double), want, "bill returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new StoreStub(seat);
    storeStubBill(double);
    const got = double.onBill.assertCalledOnce(seat, "bill records the call");
    const args = storeStubBillArgs();
    check.equal(seat, got?.id, args[0], "the call records id");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new StoreStub(seat);
    let seen = 0;
    double.onBill.onRecord(() => {
      seen += 1;
    });
    storeStubBill(double);
    storeStubBill(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of bill");
  });

  test("sets the function of withStoreBill at construction", async ({ seat }) => {
    let called = false;
    const option = StoreStub.withStoreBill(() => {
      called = true;
      return undefined as never;
    });
    storeStubBill(new StoreStub(seat, option));
    check.isTrue(seat, called, "withStoreBill sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("bill-failing");
    const double = StoreStub.newFailing(seat, want);
    const outcome = await settle(() => storeStubBill(double));
    const message = "bill fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new StoreStub(seat);
    const want = storeStubBillPin();
    double.onBill.returns(want);
    storeStubBill(double);
    double.resetCalls();
    check.equal(seat, storeStubBill(double), want, "a reset keeps the result of returns");
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

  test("passes find to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubFind(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onFind.assertCalledOnce(seat, "find reaches the implementation");
  });

  test("fails find with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("find-delegate");
    const inner = new StoreStub(seat);
    inner.onFind.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubFind(double));
    const message = "find fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes note to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubNote(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onNote.assertCalledOnce(seat, "note reaches the implementation");
  });

  test("fails note with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("note-delegate");
    const inner = new StoreStub(seat);
    inner.onNote.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubNote(double));
    const message = "note fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes bill to the implementation of delegateTo", async ({ seat }) => {
    const inner = new StoreStub(seat);
    await settle(() => storeStubBill(new StoreStub(seat, StoreStub.delegateTo(inner))));
    inner.onBill.assertCalledOnce(seat, "bill reaches the implementation");
  });

  test("fails bill with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("bill-delegate");
    const inner = new StoreStub(seat);
    inner.onBill.faults(want, 1);
    const double = new StoreStub(seat, StoreStub.delegateTo(inner));
    const outcome = await settle(() => storeStubBill(double));
    const message = "bill fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
