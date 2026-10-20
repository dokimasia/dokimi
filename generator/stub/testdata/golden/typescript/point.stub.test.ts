import { Controlled, type Seat, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { BoxStub, PointStub } from './point.stub';

/** Assigns the sample of each property of PointStub to double. */
function pointStubProps(double: PointStub): void {
  double.x = 1.5;
  double.y = 1.5;
}

describe("PointStub", () => {
  const builds: [string, (seat: Seat) => PointStub][] = [
    ["new", (seat) => new PointStub(seat)],
    ["strict", (seat) => new PointStub(seat, PointStub.strict())],
    ["withClock", (seat) => new PointStub(seat, PointStub.withClock(new Controlled()))],
    ["withRandSource", (seat) => new PointStub(seat, PointStub.withRandSource(() => 0))],
    ["benchMode", (seat) => new PointStub(seat, PointStub.benchMode())],
    ["withLatency", (seat) => new PointStub(seat, PointStub.withLatency(1))],
    ["newDelayed", (seat) => PointStub.newDelayed(seat, 1, new Controlled())],
  ];
  for (const [name, build] of builds) {
    test(`keeps each property across a reset of a double built by ${name}`, ({ seat }) => {
      const double = build(seat);
      pointStubProps(double);
      double.resetCalls();
      const want = new PointStub(undefined);
      pointStubProps(want);
      check.equal(seat, double.x, want.x, "the reset keeps the property x");
      check.equal(seat, double.y, want.y, "the reset keeps the property y");
    });
  }
});

describe("BoxStub", () => {
  const builds: [string, (seat: Seat) => BoxStub<unknown>][] = [
    ["new", (seat) => new BoxStub<unknown>(seat)],
    ["strict", (seat) => new BoxStub<unknown>(seat, BoxStub.strict<unknown>())],
    ["withClock", (seat) => new BoxStub<unknown>(seat, BoxStub.withClock<unknown>(new Controlled()))],
    ["withRandSource", (seat) => new BoxStub<unknown>(seat, BoxStub.withRandSource<unknown>(() => 0))],
    ["benchMode", (seat) => new BoxStub<unknown>(seat, BoxStub.benchMode<unknown>())],
    ["withLatency", (seat) => new BoxStub<unknown>(seat, BoxStub.withLatency<unknown>(1))],
    ["newDelayed", (seat) => BoxStub.newDelayed<unknown>(seat, 1, new Controlled())],
  ];
  for (const [name, build] of builds) {
    test(`keeps each property across a reset of a double built by ${name}`, ({ seat }) => {
      const double = build(seat);
      double.resetCalls();
      const want = new BoxStub<unknown>(undefined);
      check.equal(seat, double.value, want.value, "the reset keeps the property value");
    });
  }
});
