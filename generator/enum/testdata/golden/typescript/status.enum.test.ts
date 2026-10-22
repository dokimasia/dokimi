import { byIdentity } from "@dokimi/assert";
import { describe, expect, it } from "vitest";
import { Grade, Header, Mode, Region, Status, modeToString } from "./status";
import { UnknownGradeError, UnknownHeaderError, UnknownRegionError, UnknownStatusError, gradeToString, gradeValues, headerToString, headerValues, isGrade, isHeader, isMode, isRegion, isStatus, parseGrade, parseHeader, parseRegion, parseStatus, regionToString, regionValues, statusToString, statusValues } from "./status.enum";

describe("Status", () => {
  it("gives each variant a value of its own", () => {
    const values = [Status.Draft, Status.Published, Status.Archived];
    expect(new Set(values)).hasLength(values.length, "each variant has a value of its own");
  });

  describe("statusValues", () => {
    it("returns the declared variants in the order of their declaration", () => {
      const want = [Status.Draft, Status.Published, Status.Archived];
      expect(statusValues()).isEqualTo(want, "statusValues returns the declared variants");
    });

    it("returns a new array for each call", () => {
      expect(statusValues()).isNotEqualTo(statusValues(), "the two calls return two arrays", byIdentity());
    });
  });

  describe("statusToString", () => {
    it("returns a text of its own for each variant", () => {
      const texts = [Status.Draft, Status.Published, Status.Archived].map((value) => statusToString(value));
      expect(new Set(texts)).hasLength(texts.length, "each variant has a text of its own");
    });

    it("renders a value outside the declared set with a text of its own", () => {
      const out = 4 as unknown as Status;
      const texts = [Status.Draft, Status.Published, Status.Archived].map((value) => statusToString(value));
      expect(texts).doesNotContain(statusToString(out), "an undeclared value has a text of its own");
    });
  });

  describe("isStatus", () => {
    const tests = [
      { name: "accepts Draft", give: Status.Draft },
      { name: "accepts Published", give: Status.Published },
      { name: "accepts Archived", give: Status.Archived },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(isStatus(tt.give)).isTrue("isStatus accepts each declared variant");
      });
    }

    it("refuses a value outside the declared set", () => {
      expect(isStatus(4)).isFalse("isStatus refuses an undeclared value");
    });
  });

  describe("parseStatus", () => {
    const tests = [
      { name: "returns Draft for its text", give: Status.Draft },
      { name: "returns Published for its text", give: Status.Published },
      { name: "returns Archived for its text", give: Status.Archived },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(parseStatus(statusToString(tt.give))).isEqualTo(tt.give, "parseStatus returns the variant of the text");
      });
    }

    it("refuses a text that no variant has", () => {
      let thrown: unknown;
      try {
        parseStatus("??????????");
      } catch (error) {
        thrown = error;
      }
      expect(thrown).errorAs(UnknownStatusError, "parseStatus refuses the text with UnknownStatusError");
    });
  });
});

describe("Region", () => {
  it("gives each variant a value of its own", () => {
    const values = [Region.US, Region.EU, Region.AP];
    expect(new Set(values)).hasLength(values.length, "each variant has a value of its own");
  });

  describe("regionValues", () => {
    it("returns the declared variants in the order of their declaration", () => {
      const want = [Region.US, Region.EU, Region.AP];
      expect(regionValues()).isEqualTo(want, "regionValues returns the declared variants");
    });

    it("returns a new array for each call", () => {
      expect(regionValues()).isNotEqualTo(regionValues(), "the two calls return two arrays", byIdentity());
    });
  });

  describe("regionToString", () => {
    const tests = [
      { name: "returns the declared text of US", give: Region.US, want: "us-east" },
      { name: "returns the declared text of EU", give: Region.EU, want: "eu-west" },
      { name: "returns the declared text of AP", give: Region.AP, want: "ap-south" },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(regionToString(tt.give)).isEqualTo(tt.want, "regionToString returns the declared text of the variant");
      });
    }

    it("returns a text of its own for each variant", () => {
      const texts = [Region.US, Region.EU, Region.AP].map((value) => regionToString(value));
      expect(new Set(texts)).hasLength(texts.length, "each variant has a text of its own");
    });

    it("renders a value outside the declared set with a text of its own", () => {
      const out = "" as unknown as Region;
      const texts = [Region.US, Region.EU, Region.AP].map((value) => regionToString(value));
      expect(texts).doesNotContain(regionToString(out), "an undeclared value has a text of its own");
    });
  });

  describe("isRegion", () => {
    const tests = [
      { name: "accepts US", give: Region.US },
      { name: "accepts EU", give: Region.EU },
      { name: "accepts AP", give: Region.AP },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(isRegion(tt.give)).isTrue("isRegion accepts each declared variant");
      });
    }

    it("refuses a value outside the declared set", () => {
      expect(isRegion("")).isFalse("isRegion refuses an undeclared value");
    });
  });

  describe("parseRegion", () => {
    const tests = [
      { name: "returns US for its text", give: Region.US },
      { name: "returns EU for its text", give: Region.EU },
      { name: "returns AP for its text", give: Region.AP },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(parseRegion(regionToString(tt.give))).isEqualTo(tt.give, "parseRegion returns the variant of the text");
      });
    }

    it("refuses a text that no variant has", () => {
      let thrown: unknown;
      try {
        parseRegion("?????????");
      } catch (error) {
        thrown = error;
      }
      expect(thrown).errorAs(UnknownRegionError, "parseRegion refuses the text with UnknownRegionError");
    });
  });
});

describe("Header", () => {
  it("gives each variant a value of its own", () => {
    const values = [Header["content-type"], Header.Accept, Header.Html, Header["1st"]];
    expect(new Set(values)).hasLength(values.length, "each variant has a value of its own");
  });

  describe("headerValues", () => {
    it("returns the declared variants in the order of their declaration", () => {
      const want = [Header["content-type"], Header.Accept, Header.Html, Header["1st"]];
      expect(headerValues()).isEqualTo(want, "headerValues returns the declared variants");
    });

    it("returns a new array for each call", () => {
      expect(headerValues()).isNotEqualTo(headerValues(), "the two calls return two arrays", byIdentity());
    });
  });

  describe("headerToString", () => {
    const tests = [
      { name: "returns the declared text of content-type", give: Header["content-type"], want: "ct" },
      { name: "returns the declared text of Accept", give: Header.Accept, want: "a" },
      { name: "returns the declared text of Html", give: Header.Html, want: "<html>" },
      { name: "returns the declared text of 1st", give: Header["1st"], want: "first" },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(headerToString(tt.give)).isEqualTo(tt.want, "headerToString returns the declared text of the variant");
      });
    }

    it("returns a text of its own for each variant", () => {
      const texts = [Header["content-type"], Header.Accept, Header.Html, Header["1st"]].map((value) => headerToString(value));
      expect(new Set(texts)).hasLength(texts.length, "each variant has a text of its own");
    });

    it("renders a value outside the declared set with a text of its own", () => {
      const out = "" as unknown as Header;
      const texts = [Header["content-type"], Header.Accept, Header.Html, Header["1st"]].map((value) => headerToString(value));
      expect(texts).doesNotContain(headerToString(out), "an undeclared value has a text of its own");
    });
  });

  describe("isHeader", () => {
    const tests = [
      { name: "accepts content-type", give: Header["content-type"] },
      { name: "accepts Accept", give: Header.Accept },
      { name: "accepts Html", give: Header.Html },
      { name: "accepts 1st", give: Header["1st"] },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(isHeader(tt.give)).isTrue("isHeader accepts each declared variant");
      });
    }

    it("refuses a value outside the declared set", () => {
      expect(isHeader("")).isFalse("isHeader refuses an undeclared value");
    });
  });

  describe("parseHeader", () => {
    const tests = [
      { name: "returns content-type for its text", give: Header["content-type"] },
      { name: "returns Accept for its text", give: Header.Accept },
      { name: "returns Html for its text", give: Header.Html },
      { name: "returns 1st for its text", give: Header["1st"] },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(parseHeader(headerToString(tt.give))).isEqualTo(tt.give, "parseHeader returns the variant of the text");
      });
    }

    it("refuses a text that no variant has", () => {
      let thrown: unknown;
      try {
        parseHeader("???????");
      } catch (error) {
        thrown = error;
      }
      expect(thrown).errorAs(UnknownHeaderError, "parseHeader refuses the text with UnknownHeaderError");
    });
  });
});

describe("Grade", () => {
  it("gives each variant a value of its own", () => {
    const values = [Grade.Poor, Grade.Fair, Grade.Good];
    expect(new Set(values)).hasLength(values.length, "each variant has a value of its own");
  });

  describe("gradeValues", () => {
    it("returns the declared variants in the order of their declaration", () => {
      const want = [Grade.Poor, Grade.Fair, Grade.Good];
      expect(gradeValues()).isEqualTo(want, "gradeValues returns the declared variants");
    });

    it("returns a new array for each call", () => {
      expect(gradeValues()).isNotEqualTo(gradeValues(), "the two calls return two arrays", byIdentity());
    });
  });

  describe("gradeToString", () => {
    const tests = [
      { name: "returns the declared text of Poor", give: Grade.Poor, want: "poor" },
      { name: "returns the declared text of Fair", give: Grade.Fair, want: "fair enough" },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(gradeToString(tt.give)).isEqualTo(tt.want, "gradeToString returns the declared text of the variant");
      });
    }

    it("returns a text of its own for each variant", () => {
      const texts = [Grade.Poor, Grade.Fair, Grade.Good].map((value) => gradeToString(value));
      expect(new Set(texts)).hasLength(texts.length, "each variant has a text of its own");
    });

    it("renders a value outside the declared set with a text of its own", () => {
      const out = 3 as unknown as Grade;
      const texts = [Grade.Poor, Grade.Fair, Grade.Good].map((value) => gradeToString(value));
      expect(texts).doesNotContain(gradeToString(out), "an undeclared value has a text of its own");
    });
  });

  describe("isGrade", () => {
    const tests = [
      { name: "accepts Poor", give: Grade.Poor },
      { name: "accepts Fair", give: Grade.Fair },
      { name: "accepts Good", give: Grade.Good },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(isGrade(tt.give)).isTrue("isGrade accepts each declared variant");
      });
    }

    it("refuses a value outside the declared set", () => {
      expect(isGrade(3)).isFalse("isGrade refuses an undeclared value");
    });
  });

  describe("parseGrade", () => {
    const tests = [
      { name: "returns Poor for its text", give: Grade.Poor },
      { name: "returns Fair for its text", give: Grade.Fair },
      { name: "returns Good for its text", give: Grade.Good },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(parseGrade(gradeToString(tt.give))).isEqualTo(tt.give, "parseGrade returns the variant of the text");
      });
    }

    it("refuses a text that no variant has", () => {
      let thrown: unknown;
      try {
        parseGrade("????????????");
      } catch (error) {
        thrown = error;
      }
      expect(thrown).errorAs(UnknownGradeError, "parseGrade refuses the text with UnknownGradeError");
    });
  });
});

describe("Mode", () => {
  it("gives each variant a value of its own", () => {
    const values = [Mode.Off, Mode.On];
    expect(new Set(values)).hasLength(values.length, "each variant has a value of its own");
  });

  describe("modeToString", () => {
    it("returns a text of its own for each variant", () => {
      const texts = [Mode.Off, Mode.On].map((value) => modeToString(value));
      expect(new Set(texts)).hasLength(texts.length, "each variant has a text of its own");
    });
  });

  describe("isMode", () => {
    const tests = [
      { name: "accepts Off", give: Mode.Off },
      { name: "accepts On", give: Mode.On },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        expect(isMode(tt.give)).isTrue("isMode accepts each declared variant");
      });
    }

    it("refuses a value outside the declared set", () => {
      expect(isMode(2)).isFalse("isMode refuses an undeclared value");
    });
  });
});
