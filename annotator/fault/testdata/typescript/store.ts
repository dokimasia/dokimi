/** Reports an id without an item. */
export const notFound = new Error("not found");

/** Represents a store of items. */
export interface Store {
  /**
   * Returns the item of an id.
   *
   * @remarks
   * The double recovers on the third attempt. A partition fault selects its
   * calls by the id, and the double arms the error {@link notFound}.
   */
  //dokimi:fault errors=[notFound] retry=3 partition=id
  get(id: string): Promise<string>;
  /**
   * Returns the count of the items.
   *
   * @remarks
   * Every TypeScript method can fail, so a synchronous method takes a fault
   * too.
   */
  //dokimi:fault retry=2
  size(): number;
}
