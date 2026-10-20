/** Represents a factory with a construct signature, which a double cannot implement. */
//dokimi:stub
export interface Constructed {
  /** Creates an object in the zone. */
  new (zone: string): object;
}
