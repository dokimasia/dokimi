/** Represents a gauge with a default on a property without a type annotation. */
//dokimi:builder
export class Gauge {
  /** Contains the level, whose default has no type to check it against. */
  //dokimi:default 5
  level = 0;
}

/** Represents a record with two properties whose setters have one name. */
//dokimi:builder
export class Named {
  /** Contains a name, which has the setter withName. */
  name: string = "";
  /** Contains a second name, which has the setter withName too. */
  Name: string = "";
}

/** Represents a record without a property that a builder can set. */
//dokimi:builder
export class Sealed {
  /** Contains a private count. */
  private count: number = 0;

  /** Returns the private count. */
  read(): number {
    return this.count;
  }
}
