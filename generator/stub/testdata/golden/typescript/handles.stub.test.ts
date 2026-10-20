import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import type { Cursor } from './handles';
import { CursorStub, ProducerStub } from './handles.stub';

/** Calls open of double with the samples of its arguments. */
function producerStubOpen(double: ProducerStub): Promise<Cursor> {
  return double.open();
}

describe("ProducerStub.open", () => {
  const checks = behaviour("Producer.open", (seat) => {
    const double = new ProducerStub(seat);
    return {
      stub: double.onOpen,
      call: () => producerStubOpen(double),
      override: (mark) => {
        double.onOpen.func(() => {
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
    const double = new ProducerStub(seat);
    await producerStubOpen(double);
    double.onOpen.assertCalledOnce(seat, "open records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new ProducerStub(seat);
    let seen = 0;
    double.onOpen.onRecord(() => {
      seen += 1;
    });
    await producerStubOpen(double);
    await producerStubOpen(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of open");
  });

  test("sets the function of withProducerOpen at construction", async ({ seat }) => {
    let called = false;
    const option = ProducerStub.withProducerOpen(() => {
      called = true;
      return undefined as never;
    });
    await producerStubOpen(new ProducerStub(seat, option));
    check.isTrue(seat, called, "withProducerOpen sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("open-failing");
    const double = ProducerStub.newFailing(seat, want);
    const outcome = await settle(() => producerStubOpen(double));
    const message = "open fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps recording across a reset", async ({ seat }) => {
    const double = new ProducerStub(seat);
    double.onOpen.returns(undefined as never);
    await producerStubOpen(double);
    double.resetCalls();
    await producerStubOpen(double);
    double.onOpen.assertCalledOnce(seat, "open records the call after the reset");
  });
});

describe("ProducerStub", () => {
  const instance = (double: ProducerStub) => ({
    stub: double.onOpen,
    call: () => producerStubOpen(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new ProducerStub(seat)),
    withClock: (seat, clock) => instance(new ProducerStub(seat, ProducerStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(ProducerStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new ProducerStub(seat, ProducerStub.withRandSource(random))),
    benchMode: (seat) => instance(new ProducerStub(seat, ProducerStub.benchMode())),
    strict: (seat) => instance(new ProducerStub(seat, ProducerStub.strict())),
    async: true,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes open to the implementation of delegateTo", async ({ seat }) => {
    const inner = new ProducerStub(seat);
    await settle(() => producerStubOpen(new ProducerStub(seat, ProducerStub.delegateTo(inner))));
    inner.onOpen.assertCalledOnce(seat, "open reaches the implementation");
  });

  test("fails open with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("open-delegate");
    const inner = new ProducerStub(seat);
    inner.onOpen.faults(want, 1);
    const double = new ProducerStub(seat, ProducerStub.delegateTo(inner));
    const outcome = await settle(() => producerStubOpen(double));
    const message = "open fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

/** Calls next of double with the samples of its arguments. */
function cursorStubNext(double: CursorStub): Promise<string> {
  return double.next();
}

/** Returns the sample of the result of next. */
function cursorStubNextPin(): string {
  let want0: string;
  want0 = 'next-a';
  return want0;
}

/** Calls close of double with the samples of its arguments. */
function cursorStubClose(double: CursorStub): void {
  return double.close();
}

describe("CursorStub.next", () => {
  const checks = behaviour("Cursor.next", (seat) => {
    const double = new CursorStub(seat);
    return {
      stub: double.onNext,
      call: () => cursorStubNext(double),
      override: (mark) => {
        double.onNext.func(() => {
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
    const double = new CursorStub(seat);
    const want = cursorStubNextPin();
    double.onNext.returns(want);
    check.equal(seat, await cursorStubNext(double), want, "next returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new CursorStub(seat);
    await cursorStubNext(double);
    double.onNext.assertCalledOnce(seat, "next records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new CursorStub(seat);
    let seen = 0;
    double.onNext.onRecord(() => {
      seen += 1;
    });
    await cursorStubNext(double);
    await cursorStubNext(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of next");
  });

  test("sets the function of withCursorNext at construction", async ({ seat }) => {
    let called = false;
    const option = CursorStub.withCursorNext(() => {
      called = true;
      return undefined as never;
    });
    await cursorStubNext(new CursorStub(seat, option));
    check.isTrue(seat, called, "withCursorNext sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("next-failing");
    const double = CursorStub.newFailing(seat, want);
    const outcome = await settle(() => cursorStubNext(double));
    const message = "next fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new CursorStub(seat);
    const want = cursorStubNextPin();
    double.onNext.returns(want);
    await cursorStubNext(double);
    double.resetCalls();
    check.equal(seat, await cursorStubNext(double), want, "a reset keeps the result of returns");
  });
});

describe("CursorStub.close", () => {
  const checks = behaviour("Cursor.close", (seat) => {
    const double = new CursorStub(seat);
    return {
      stub: double.onClose,
      call: () => cursorStubClose(double),
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
    const double = new CursorStub(seat);
    cursorStubClose(double);
    double.onClose.assertCalledOnce(seat, "close records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new CursorStub(seat);
    let seen = 0;
    double.onClose.onRecord(() => {
      seen += 1;
    });
    cursorStubClose(double);
    cursorStubClose(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of close");
  });

  test("sets the function of withCursorClose at construction", async ({ seat }) => {
    let called = false;
    const option = CursorStub.withCursorClose(() => {
      called = true;
      return undefined as never;
    });
    cursorStubClose(new CursorStub(seat, option));
    check.isTrue(seat, called, "withCursorClose sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("close-failing");
    const double = CursorStub.newFailing(seat, want);
    const outcome = await settle(() => cursorStubClose(double));
    const message = "close fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("CursorStub", () => {
  const instance = (double: CursorStub) => ({
    stub: double.onNext,
    call: () => cursorStubNext(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new CursorStub(seat)),
    withClock: (seat, clock) => instance(new CursorStub(seat, CursorStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(CursorStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new CursorStub(seat, CursorStub.withRandSource(random))),
    benchMode: (seat) => instance(new CursorStub(seat, CursorStub.benchMode())),
    strict: (seat) => instance(new CursorStub(seat, CursorStub.strict())),
    async: true,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes next to the implementation of delegateTo", async ({ seat }) => {
    const inner = new CursorStub(seat);
    await settle(() => cursorStubNext(new CursorStub(seat, CursorStub.delegateTo(inner))));
    inner.onNext.assertCalledOnce(seat, "next reaches the implementation");
  });

  test("fails next with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("next-delegate");
    const inner = new CursorStub(seat);
    inner.onNext.faults(want, 1);
    const double = new CursorStub(seat, CursorStub.delegateTo(inner));
    const outcome = await settle(() => cursorStubNext(double));
    const message = "next fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes close to the implementation of delegateTo", async ({ seat }) => {
    const inner = new CursorStub(seat);
    await settle(() => cursorStubClose(new CursorStub(seat, CursorStub.delegateTo(inner))));
    inner.onClose.assertCalledOnce(seat, "close reaches the implementation");
  });

  test("fails close with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("close-delegate");
    const inner = new CursorStub(seat);
    inner.onClose.faults(want, 1);
    const double = new CursorStub(seat, CursorStub.delegateTo(inner));
    const outcome = await settle(() => cursorStubClose(double));
    const message = "close fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
