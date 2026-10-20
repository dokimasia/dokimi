/** Represents a ticket, whose constructor takes a parameter that a caller must pass. */
//dokimi:builder
export class Ticket {
  /** Contains the code that the constructor sets. */
  code: string;

  constructor(code: string) {
    this.code = code;
  }
}

/**
 * Represents a voucher, whose constructor takes parameters that a caller can
 * leave out.
 *
 * @remarks
 * code has a default, and note has the optional mark.
 */
//dokimi:builder
export class Voucher {
  /** Contains the code that the constructor sets. */
  code: string;
  /** Contains the note that the constructor sets. */
  note: string;

  constructor(code = "", note?: string) {
    this.code = code;
    this.note = note ?? "";
  }
}

/**
 * Represents a coupon, whose constructor takes a parameter that a caller must
 * pass after a parameter with a default.
 */
//dokimi:builder
export class Coupon {
  /** Contains the code that the constructor sets. */
  code: string;
  /** Contains the owner that the constructor sets. */
  owner: string;

  constructor(code = "", owner: string) {
    this.code = code;
    this.owner = owner;
  }
}
