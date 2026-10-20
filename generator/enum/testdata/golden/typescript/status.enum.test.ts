import { check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { Grade, Header, Mode, Region, Status, modeToString } from './status';
import { UnknownGradeError, UnknownHeaderError, UnknownRegionError, UnknownStatusError, gradeToString, gradeValues, headerToString, headerValues, isGrade, isHeader, isMode, isRegion, isStatus, parseGrade, parseHeader, parseRegion, parseStatus, regionToString, regionValues, statusToString, statusValues } from './status.enum';

describe("Status variants", () => {
  const all = statusValues();

  test("declares exactly 3 variants", ({ seat }) => {
    check.length(seat, all, 3, "the declared set must not change unnoticed");
  });

  test("no two variants share a value", ({ seat }) => {
    const seen = new Set<Status>();
    for (const value of all) {
      check.isFalse(seat, seen.has(value), "each variant must hold a distinct value");
      seen.add(value);
    }
  });

  test("no two variants share a textual form", ({ seat }) => {
    const seen = new Set<string>();
    for (const value of all) {
      check.isFalse(seat, seen.has(statusToString(value)), "each variant must render distinctly");
      seen.add(statusToString(value));
    }
  });

  test("a value outside the set is not valid", ({ seat }) => {
    check.isFalse(seat, isStatus(4), "an undeclared value must not be valid");
  });

  test("every declared variant is valid", ({ seat }) => {
    for (const value of all) {
      check.isTrue(seat, isStatus(value), "a declared variant must be valid");
    }
  });
});

describe("Status text", () => {
  const all = statusValues();

  test("every variant survives a text round trip", ({ seat }) => {
    for (const value of all) {
      check.equal(seat, parseStatus(statusToString(value)), value, "the round trip must be lossless");
    }
  });

  test("text naming no variant is refused", ({ seat }) => {
    const err = check.throws(seat, () => parseStatus("??????????"), "unknown text must be refused");
    check.errorAs(seat, err, UnknownStatusError, "unknown text must be refused with the declared refusal");
  });

  test("a value outside the set does not render as one inside it", ({ seat }) => {
    const out = 4 as unknown as Status;
    for (const value of all) {
      check.notEqual(seat, statusToString(out), statusToString(value), "an undeclared value must render distinctly");
    }
  });
});

describe("Region variants", () => {
  const all = regionValues();

  test("declares exactly 3 variants", ({ seat }) => {
    check.length(seat, all, 3, "the declared set must not change unnoticed");
  });

  test("no two variants share a value", ({ seat }) => {
    const seen = new Set<Region>();
    for (const value of all) {
      check.isFalse(seat, seen.has(value), "each variant must hold a distinct value");
      seen.add(value);
    }
  });

  test("every variant renders the value it declares", ({ seat }) => {
    check.equal(seat, regionToString(Region.US), "us-east", "the declared value is the textual form");
    check.equal(seat, regionToString(Region.EU), "eu-west", "the declared value is the textual form");
    check.equal(seat, regionToString(Region.AP), "ap-south", "the declared value is the textual form");
  });

  test("no two variants share a textual form", ({ seat }) => {
    const seen = new Set<string>();
    for (const value of all) {
      check.isFalse(seat, seen.has(regionToString(value)), "each variant must render distinctly");
      seen.add(regionToString(value));
    }
  });

  test("a value outside the set is not valid", ({ seat }) => {
    check.isFalse(seat, isRegion(""), "an undeclared value must not be valid");
  });

  test("every declared variant is valid", ({ seat }) => {
    for (const value of all) {
      check.isTrue(seat, isRegion(value), "a declared variant must be valid");
    }
  });
});

describe("Region text", () => {
  const all = regionValues();

  test("every variant survives a text round trip", ({ seat }) => {
    for (const value of all) {
      check.equal(seat, parseRegion(regionToString(value)), value, "the round trip must be lossless");
    }
  });

  test("text naming no variant is refused", ({ seat }) => {
    const err = check.throws(seat, () => parseRegion("?????????"), "unknown text must be refused");
    check.errorAs(seat, err, UnknownRegionError, "unknown text must be refused with the declared refusal");
  });

  test("a value outside the set does not render as one inside it", ({ seat }) => {
    const out = "" as unknown as Region;
    for (const value of all) {
      check.notEqual(seat, regionToString(out), regionToString(value), "an undeclared value must render distinctly");
    }
  });
});

describe("Header variants", () => {
  const all = headerValues();

  test("declares exactly 4 variants", ({ seat }) => {
    check.length(seat, all, 4, "the declared set must not change unnoticed");
  });

  test("no two variants share a value", ({ seat }) => {
    const seen = new Set<Header>();
    for (const value of all) {
      check.isFalse(seat, seen.has(value), "each variant must hold a distinct value");
      seen.add(value);
    }
  });

  test("every variant renders the value it declares", ({ seat }) => {
    check.equal(seat, headerToString(Header["content-type"]), "ct", "the declared value is the textual form");
    check.equal(seat, headerToString(Header.Accept), "a", "the declared value is the textual form");
    check.equal(seat, headerToString(Header.Html), "<html>", "the declared value is the textual form");
    check.equal(seat, headerToString(Header["1st"]), "first", "the declared value is the textual form");
  });

  test("no two variants share a textual form", ({ seat }) => {
    const seen = new Set<string>();
    for (const value of all) {
      check.isFalse(seat, seen.has(headerToString(value)), "each variant must render distinctly");
      seen.add(headerToString(value));
    }
  });

  test("a value outside the set is not valid", ({ seat }) => {
    check.isFalse(seat, isHeader(""), "an undeclared value must not be valid");
  });

  test("every declared variant is valid", ({ seat }) => {
    for (const value of all) {
      check.isTrue(seat, isHeader(value), "a declared variant must be valid");
    }
  });
});

describe("Header text", () => {
  const all = headerValues();

  test("every variant survives a text round trip", ({ seat }) => {
    for (const value of all) {
      check.equal(seat, parseHeader(headerToString(value)), value, "the round trip must be lossless");
    }
  });

  test("text naming no variant is refused", ({ seat }) => {
    const err = check.throws(seat, () => parseHeader("???????"), "unknown text must be refused");
    check.errorAs(seat, err, UnknownHeaderError, "unknown text must be refused with the declared refusal");
  });

  test("a value outside the set does not render as one inside it", ({ seat }) => {
    const out = "" as unknown as Header;
    for (const value of all) {
      check.notEqual(seat, headerToString(out), headerToString(value), "an undeclared value must render distinctly");
    }
  });
});

describe("Grade variants", () => {
  const all = gradeValues();

  test("declares exactly 3 variants", ({ seat }) => {
    check.length(seat, all, 3, "the declared set must not change unnoticed");
  });

  test("no two variants share a value", ({ seat }) => {
    const seen = new Set<Grade>();
    for (const value of all) {
      check.isFalse(seat, seen.has(value), "each variant must hold a distinct value");
      seen.add(value);
    }
  });

  test("every variant renders the value it declares", ({ seat }) => {
    check.equal(seat, gradeToString(Grade.Poor), "poor", "the declared value is the textual form");
    check.equal(seat, gradeToString(Grade.Fair), "fair enough", "the declared value is the textual form");
  });

  test("no two variants share a textual form", ({ seat }) => {
    const seen = new Set<string>();
    for (const value of all) {
      check.isFalse(seat, seen.has(gradeToString(value)), "each variant must render distinctly");
      seen.add(gradeToString(value));
    }
  });

  test("a value outside the set is not valid", ({ seat }) => {
    check.isFalse(seat, isGrade(3), "an undeclared value must not be valid");
  });

  test("every declared variant is valid", ({ seat }) => {
    for (const value of all) {
      check.isTrue(seat, isGrade(value), "a declared variant must be valid");
    }
  });
});

describe("Grade text", () => {
  const all = gradeValues();

  test("every variant survives a text round trip", ({ seat }) => {
    for (const value of all) {
      check.equal(seat, parseGrade(gradeToString(value)), value, "the round trip must be lossless");
    }
  });

  test("text naming no variant is refused", ({ seat }) => {
    const err = check.throws(seat, () => parseGrade("????????????"), "unknown text must be refused");
    check.errorAs(seat, err, UnknownGradeError, "unknown text must be refused with the declared refusal");
  });

  test("a value outside the set does not render as one inside it", ({ seat }) => {
    const out = 3 as unknown as Grade;
    for (const value of all) {
      check.notEqual(seat, gradeToString(out), gradeToString(value), "an undeclared value must render distinctly");
    }
  });
});

describe("Mode variants", () => {
  const all = [
    Mode.Off,
    Mode.On,
  ];

  test("declares exactly 2 variants", ({ seat }) => {
    check.length(seat, all, 2, "the declared set must not change unnoticed");
  });

  test("no two variants share a value", ({ seat }) => {
    const seen = new Set<Mode>();
    for (const value of all) {
      check.isFalse(seat, seen.has(value), "each variant must hold a distinct value");
      seen.add(value);
    }
  });

  test("no two variants share a textual form", ({ seat }) => {
    const seen = new Set<string>();
    for (const value of all) {
      check.isFalse(seat, seen.has(modeToString(value)), "each variant must render distinctly");
      seen.add(modeToString(value));
    }
  });

  test("a value outside the set is not valid", ({ seat }) => {
    check.isFalse(seat, isMode(2), "an undeclared value must not be valid");
  });

  test("every declared variant is valid", ({ seat }) => {
    for (const value of all) {
      check.isTrue(seat, isMode(value), "a declared variant must be valid");
    }
  });
});
