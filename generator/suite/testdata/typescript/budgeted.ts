/** Represents a meter whose count has a budget of allocations. */
//dokimi:suite
export interface Meter {
  /** Returns the count of the meter. */
  //dokimi:budget allocs=0
  count(): number;
  /** Returns the size of the meter. */
  //dokimi:budget bytes=0
  size(): number;
  /** Returns the total of the meter. */
  //dokimi:budget p99=1ms
  total(): number;
}
