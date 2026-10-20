import { Set } from "immutable";

/** Represents a record with a property of the class Set of another package. */
//dokimi:builder
export class Frozen {
  /**
   * Contains the tags in a Set of the package immutable.
   *
   * @remarks
   * The Set of immutable is no set of TypeScript.
   */
  tags: Set<string> = Set<string>();
}
