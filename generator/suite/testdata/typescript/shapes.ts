import type { Item as OtherItem } from "./other";
import type { Item } from "./store";
import type { Item as ThirdItem } from "./third";

/** Represents a level, an enumeration without values. */
export enum Level {
  Low,
  High,
}

/** Represents a mode, an enumeration of strings. */
export enum Mode {
  Read = "read",
  Write = "write",
}

/** Represents a status, an enumeration of numbers from 1. */
export enum Status {
  Open = 1,
  Closed,
}

/** Represents flags, an enumeration of computed numbers. */
export enum Flags {
  Read = 1 << 0,
  Write = 1 << 1,
}

/** Represents a limit, an enumeration of an infinite number. */
export enum Limit {
  Max = Infinity,
}

/** Represents an odd value, an enumeration of a number that is not a number. */
export enum Odd {
  Value = NaN,
}

/** Represents a signed value, an enumeration of a negative number and the number after it. */
export enum Signed {
  Below = -1,
  Zero,
}

/** Represents a key, an alias of a string. */
export type Key = string;

/** Represents a step of a state machine, an alias of a function type, which no shape states. */
export type Step = (input: string) => Step | undefined;

/** Represents a node, an alias of an interface. */
export type Node = Tree;

/** Represents names, an alias of an array. */
export type Names = string[];

/** Represents a tag, an alias of a string with an authored sample. */
//dokimi:sample value=newTag()
export type Tag = string;

/** Returns the tag that the checks write. */
export function newTag(): Tag {
  return "release";
}

/** Represents a node of a tree, whose children are trees. */
export interface Tree {
  /** Contains the value of the node. */
  value: string;
  /** Contains the children of the node. */
  children: Tree[];
}

/** Represents a source of data, whose method no shape states. */
export interface Source {
  /** Returns the data. */
  read(): Uint8Array;
}

/** Represents a pair of two items with one name. */
export interface Pair {
  /** Contains the item of the store. */
  left: Item;
  /** Contains the item of the other module. */
  right: OtherItem;
}

/** Represents three items with one name. */
export interface Triple {
  /** Contains the item of the store. */
  a: Item;
  /** Contains the item of the other module. */
  b: OtherItem;
  /** Contains the item of the third module. */
  c: ThirdItem;
}

/** Represents a value of a class with the name of the shape definition of its draw. */
export class sinkPShape {
  /** Contains the value. */
  value: string = "";
}

/** Represents a box of a value of the type T. */
export interface Box<T> {
  /** Contains the value. */
  value: T;
}

/** Represents a point, a class of two coordinates. */
export class Point {
  /** Contains the horizontal coordinate. */
  x: number = 0;
  /** Contains the vertical coordinate. */
  y: number = 0;
}

/** Represents a sink of the values of each form of a shape. */
//dokimi:suite
export interface Sink {
  /** Writes a flag. */
  flag(on: boolean): void;
  /** Writes bytes. */
  bytes(data: Uint8Array): void;
  /** Writes texts. */
  texts(items: string[]): void;
  /** Writes numbers. */
  numbers(items: Array<number>): void;
  /** Writes texts that the sink does not change. */
  frozen(items: ReadonlyArray<string>): void;
  /** Writes a count. */
  count(n: bigint): void;
  /** Writes a level. */
  level(l: Level): void;
  /** Writes a mode. */
  mode(m: Mode): void;
  /** Writes a status. */
  status(s: Status): void;
  /** Writes flags. */
  flags(f: Flags): void;
  /** Writes a limit. */
  limit(l: Limit): void;
  /** Writes an odd value. */
  odd(o: Odd): void;
  /** Writes a key. */
  key(k: Key): void;
  /** Writes a tree. */
  tree(t: Tree): void;
  /** Writes a source. */
  source(s: Source): void;
  /** Writes a choice. */
  choice(c: string | number): void;
  /** Writes a pair. */
  pair(p: Pair): void;
  /** Writes a lookup. */
  lookup(entries: Map<string, number>): void;
  /** Writes a table. */
  table(entries: Record<string, boolean>): void;
  /** Writes a box. */
  box(b: Box<number>): void;
  /** Writes a point. */
  point(p: Point): void;
  /** Writes an item of any type. */
  any<T>(item: T): void;
  /** Writes a signed value. */
  signed(s: Signed): void;
  /** Writes texts as separate arguments. */
  many(...items: string[]): void;
  /** Writes a label that the caller can leave out. */
  maybe(label?: string): void;
  /** Writes three items. */
  triple(t: Triple): void;
  /** Writes a value of a type with the name of a shape definition. */
  named(p: Box<sinkPShape>): void;
  /** Writes a tag. */
  tag(t: Tag): void;
  /** Writes a step. */
  step(s: Step): void;
  /** Writes a node. */
  node(n: Node): void;
  /** Writes names. */
  names(n: Names): void;
  /** Writes a label that the caller can pass as undefined. */
  unset(label: string | undefined): void;
}
