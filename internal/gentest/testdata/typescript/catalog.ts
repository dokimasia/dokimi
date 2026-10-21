/** Represents the postal address of a customer. */
export class Address {
  /** Contains the street and the house number. */
  street: string = "";

  /** Contains the labels of the address, in a Map that the library of TypeScript declares. */
  labels: Map<string, string> = new Map();
}
