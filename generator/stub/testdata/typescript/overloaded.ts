/** Represents a parser with a method of two signatures, which one callable cannot implement. */
//dokimi:stub
export interface Overloaded {
  /** Returns the number that the text spells in base 10. */
  parse(text: string): number;
  /** Returns the number that the text spells in the radix. */
  parse(text: string, radix: number): number;
}
