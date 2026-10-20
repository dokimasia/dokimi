/** Represents the postal address of a customer. */
//dokimi:builder
export class Address {
  /** Contains the street and the house number. */
  //dokimi:default "Dam 1"
  street: string = "";
  /** Lists the lines of the address. */
  lines: string[] = [];
}
