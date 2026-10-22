import { describe, expect, it } from "vitest";
import { LabelError } from "./label";

function labelErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new LabelError("cause-a", "op-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("contains each text member", () => {
    const got = new LabelError("cause-a", "op-a");
    expect(got.message).doesContain(got.op, "the message contains op");
  });
}

describe("LabelError", () => {
  describe("message", labelErrorMessageCases);
});
