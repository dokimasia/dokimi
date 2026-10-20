/**
 * Represents a set of the workspace.
 *
 * @remarks
 * The class shadows the global class Set.
 */
export class Set<T> {
  /** Lists the members of the set. */
  items: T[] = [];
}

/** Represents a record with a property of the class Set of the workspace. */
//dokimi:builder
export class Shadowed {
  /**
   * Contains the tags in a Set of the workspace.
   *
   * @remarks
   * The Set of the workspace is no set of TypeScript.
   */
  tags: Set<string> = new Set<string>();
}
