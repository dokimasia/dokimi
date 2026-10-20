import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import { FeedStub } from './stream.stub';

/** Calls names of double with the samples of its arguments. */
function feedStubNames(double: FeedStub): Iterable<string> {
  return double.names();
}

/** Calls rows of double with the samples of its arguments. */
function feedStubRows(double: FeedStub): AsyncIterable<number> {
  return double.rows();
}

describe("FeedStub.names", () => {
  const checks = behaviour("Feed.names", (seat) => {
    const double = new FeedStub(seat);
    return {
      stub: double.onNames,
      call: () => feedStubNames(double),
      override: (mark) => {
        double.onNames.func(() => {
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
    const double = new FeedStub(seat);
    feedStubNames(double);
    double.onNames.assertCalledOnce(seat, "names records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new FeedStub(seat);
    let seen = 0;
    double.onNames.onRecord(() => {
      seen += 1;
    });
    feedStubNames(double);
    feedStubNames(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of names");
  });

  test("sets the function of withFeedNames at construction", async ({ seat }) => {
    let called = false;
    const option = FeedStub.withFeedNames(() => {
      called = true;
      return undefined as never;
    });
    feedStubNames(new FeedStub(seat, option));
    check.isTrue(seat, called, "withFeedNames sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("names-failing");
    const double = FeedStub.newFailing(seat, want);
    const outcome = await settle(() => feedStubNames(double));
    const message = "names fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps recording across a reset", async ({ seat }) => {
    const double = new FeedStub(seat);
    double.onNames.returns(undefined as never);
    feedStubNames(double);
    double.resetCalls();
    feedStubNames(double);
    double.onNames.assertCalledOnce(seat, "names records the call after the reset");
  });
});

describe("FeedStub.rows", () => {
  const checks = behaviour("Feed.rows", (seat) => {
    const double = new FeedStub(seat);
    return {
      stub: double.onRows,
      call: () => feedStubRows(double),
      override: (mark) => {
        double.onRows.func(() => {
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
    const double = new FeedStub(seat);
    feedStubRows(double);
    double.onRows.assertCalledOnce(seat, "rows records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new FeedStub(seat);
    let seen = 0;
    double.onRows.onRecord(() => {
      seen += 1;
    });
    feedStubRows(double);
    feedStubRows(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of rows");
  });

  test("sets the function of withFeedRows at construction", async ({ seat }) => {
    let called = false;
    const option = FeedStub.withFeedRows(() => {
      called = true;
      return undefined as never;
    });
    feedStubRows(new FeedStub(seat, option));
    check.isTrue(seat, called, "withFeedRows sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("rows-failing");
    const double = FeedStub.newFailing(seat, want);
    const outcome = await settle(() => feedStubRows(double));
    const message = "rows fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps recording across a reset", async ({ seat }) => {
    const double = new FeedStub(seat);
    double.onRows.returns(undefined as never);
    feedStubRows(double);
    double.resetCalls();
    feedStubRows(double);
    double.onRows.assertCalledOnce(seat, "rows records the call after the reset");
  });

  test("yields the items of yields", async ({ seat }) => {
    const double = new FeedStub(seat);
    double.onRows.yields(undefined as never, undefined as never);
    let seen = 0;
    for await (const _ of feedStubRows(double)) {
      seen += 1;
    }
    check.equal(seat, seen, 2, "rows yields each item");
  });

  test("stops yielding when the caller breaks", async ({ seat }) => {
    const double = new FeedStub(seat);
    double.onRows.yields(undefined as never, undefined as never, undefined as never);
    let seen = 0;
    for await (const _ of feedStubRows(double)) {
      seen += 1;
      break;
    }
    check.equal(seat, seen, 1, "rows stops when the caller breaks");
  });
});

describe("FeedStub", () => {
  const instance = (double: FeedStub) => ({
    stub: double.onNames,
    call: () => feedStubNames(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new FeedStub(seat)),
    withClock: (seat, clock) => instance(new FeedStub(seat, FeedStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(FeedStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new FeedStub(seat, FeedStub.withRandSource(random))),
    benchMode: (seat) => instance(new FeedStub(seat, FeedStub.benchMode())),
    strict: (seat) => instance(new FeedStub(seat, FeedStub.strict())),
    async: false,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes names to the implementation of delegateTo", async ({ seat }) => {
    const inner = new FeedStub(seat);
    await settle(() => feedStubNames(new FeedStub(seat, FeedStub.delegateTo(inner))));
    inner.onNames.assertCalledOnce(seat, "names reaches the implementation");
  });

  test("fails names with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("names-delegate");
    const inner = new FeedStub(seat);
    inner.onNames.faults(want, 1);
    const double = new FeedStub(seat, FeedStub.delegateTo(inner));
    const outcome = await settle(() => feedStubNames(double));
    const message = "names fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes rows to the implementation of delegateTo", async ({ seat }) => {
    const inner = new FeedStub(seat);
    await settle(() => feedStubRows(new FeedStub(seat, FeedStub.delegateTo(inner))));
    inner.onRows.assertCalledOnce(seat, "rows reaches the implementation");
  });

  test("fails rows with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("rows-delegate");
    const inner = new FeedStub(seat);
    inner.onRows.faults(want, 1);
    const double = new FeedStub(seat, FeedStub.delegateTo(inner));
    const outcome = await settle(() => feedStubRows(double));
    const message = "rows fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
