/** Represents the postal address of a customer. */
export class Address {
  /** Contains the street and the house number. */
  //dokimi:default "Dam 1"
  street: string = "";
  /** Contains the city. */
  city: string = "";
}

/**
 * Represents the floor level of an address.
 *
 * @remarks
 * A check that draws a level takes its default.
 */
//dokimi:default 3
export type Level = number;
