/** Represents a counter with a private member, so no other class implements it. */
//dokimi:stub
export class Sealed {
  private count: number = 0;

  /** Returns the next count. */
  next(): number {
    this.count += 1;
    return this.count;
  }
}
