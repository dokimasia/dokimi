import { prop } from "@dokimi/assert";
import { compatV1 } from "@dokimi/runtime/conform/compat";
import type { Handle } from "@dokimi/runtime/conform/handle";
import type { Options } from "@dokimi/runtime/conform/options";
import { run } from "@dokimi/runtime/conform/run";
import { type Callable, type Check, type Id, type Suite, generated } from "@dokimi/runtime/conform/suite";
import { survives } from "@dokimi/runtime/probe/survives";
import type { OtherItem } from "./other";
import { type Box, type Flags, type Key, type Level, type Mode, type Names, type Pair, Point, type Signed, type Sink, type Status, type Tag, type Triple, newTag, sinkPShape } from "./shapes";
import type { Item } from "./store";
import type { ThirdItem } from "./third";

/**
 * The handle of a check of the suite of Sink.
 */
export type SinkHandle = Handle<Sink>;

/**
 * A check of the suite of Sink.
 */
export type SinkCheck = Check<Sink>;

/**
 * The options of an application of the suite of Sink.
 */
export type SinkOptions = Options<Sink>;

/**
 * The name of each callable of Sink, with the identity of each of its
 * checks.
 */
export const sinkChecks = {
  flag: {
    callable: "Sink.flag" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.flag/smoke",
  },
  bytes: {
    callable: "Sink.bytes" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.bytes/smoke",
  },
  texts: {
    callable: "Sink.texts" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.texts/smoke",
  },
  numbers: {
    callable: "Sink.numbers" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.numbers/smoke",
  },
  frozen: {
    callable: "Sink.frozen" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.frozen/smoke",
  },
  count: {
    callable: "Sink.count" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.count/smoke",
  },
  level: {
    callable: "Sink.level" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.level/smoke",
  },
  mode: {
    callable: "Sink.mode" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.mode/smoke",
  },
  status: {
    callable: "Sink.status" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.status/smoke",
  },
  flags: {
    callable: "Sink.flags" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.flags/smoke",
  },
  limit: {
    callable: "Sink.limit" as Callable<Sink>,
  },
  odd: {
    callable: "Sink.odd" as Callable<Sink>,
  },
  key: {
    callable: "Sink.key" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.key/smoke",
  },
  tree: {
    callable: "Sink.tree" as Callable<Sink>,
  },
  source: {
    callable: "Sink.source" as Callable<Sink>,
  },
  choice: {
    callable: "Sink.choice" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.choice/smoke",
  },
  pair: {
    callable: "Sink.pair" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.pair/smoke",
  },
  lookup: {
    callable: "Sink.lookup" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.lookup/smoke",
  },
  table: {
    callable: "Sink.table" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.table/smoke",
  },
  box: {
    callable: "Sink.box" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.box/smoke",
  },
  point: {
    callable: "Sink.point" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.point/smoke",
  },
  any: {
    callable: "Sink.any" as Callable<Sink>,
  },
  signed: {
    callable: "Sink.signed" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.signed/smoke",
  },
  many: {
    callable: "Sink.many" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.many/smoke",
  },
  maybe: {
    callable: "Sink.maybe" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.maybe/smoke",
  },
  triple: {
    callable: "Sink.triple" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.triple/smoke",
  },
  named: {
    callable: "Sink.named" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.named/smoke",
  },
  tag: {
    callable: "Sink.tag" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.tag/smoke",
  },
  step: {
    callable: "Sink.step" as Callable<Sink>,
  },
  node: {
    callable: "Sink.node" as Callable<Sink>,
  },
  names: {
    callable: "Sink.names" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.names/smoke",
  },
  unset: {
    callable: "Sink.unset" as Callable<Sink>,
    smoke: (): Id<Sink> => "Sink.unset/smoke",
  },
};

/**
 * The generator of the values of boolean that the generated bindings draw.
 */
const sinkOnShape = prop.ofShape<boolean>('{"shape":"bool"}');

/**
 * The generator of the values of Uint8Array that the generated bindings
 * draw.
 */
const sinkDataShape = prop.ofShape<Uint8Array>('{"shape":"bytes"}');

/**
 * The generator of the values of string[] that the generated bindings
 * draw.
 */
const sinkItemsShape = prop.ofShape<string[]>('{"of":{"shape":"string"},"shape":"list"}');

/**
 * The generator of the values of Array<number> that the generated bindings
 * draw.
 */
const sinkItemsShape0 = prop.ofShape<Array<number>>('{"of":{"shape":"float","width":64},"shape":"list"}');

/**
 * The generator of the values of ReadonlyArray<string> that the generated
 * bindings draw.
 */
const sinkItemsShape1 = prop.ofShape<ReadonlyArray<string>>('{"of":{"shape":"string"},"shape":"list"}');

/**
 * The generator of the values of bigint that the generated bindings draw.
 */
const sinkNShape = prop.ofShape<bigint>('{"shape":"int","signed":true,"width":64}');

/**
 * The generator of the values of Level that the generated bindings draw.
 */
const sinkLShape = prop.ofShape<Level>('{"shape":"literal","values":[{"type":"float","value":0},{"type":"float","value":1}]}');

/**
 * The generator of the values of Mode that the generated bindings draw.
 */
const sinkMShape = prop.ofShape<Mode>('{"shape":"literal","values":[{"type":"string","value":"read"},{"type":"string","value":"write"}]}');

/**
 * The generator of the values of Status that the generated bindings draw.
 */
const sinkSShape = prop.ofShape<Status>('{"shape":"literal","values":[{"type":"float","value":1},{"type":"float","value":2}]}');

/**
 * The generator of the values of Flags that the generated bindings draw.
 */
const sinkFShape = prop.ofShape<Flags>('{"name":"Flags","shape":"ref"}');

/**
 * The generator of the values of Key that the generated bindings draw.
 */
const sinkKShape = prop.ofShape<Key>('{"shape":"string"}');

/**
 * The generator of the values of string | number that the generated
 * bindings draw.
 */
const sinkC0Shape = prop.ofShape<string | number>('{"name":"string | number","shape":"ref"}');

/**
 * The generator of the values of Pair that the generated bindings draw.
 */
const sinkPShape0 = prop.ofShape<Pair>('{"definitions":{"Item":{"fields":[["id",{"shape":"string"}],["count",{"shape":"float","width":64}],["label",{"of":{"shape":"string"},"shape":"optional"}]],"shape":"record"},"Item2":{"fields":[["name",{"shape":"string"}]],"shape":"record"},"Pair":{"fields":[["left",{"name":"Item","shape":"ref"}],["right",{"name":"Item2","shape":"ref"}]],"shape":"record"}},"name":"Pair","shape":"ref"}');

/**
 * The generator of the values of Map<string, number> that the generated
 * bindings draw.
 */
const sinkEntriesShape = prop.ofShape<Map<string, number>>('{"key":{"shape":"string"},"of":{"shape":"float","width":64},"shape":"map"}');

/**
 * The generator of the values of Record<string, boolean> that the
 * generated bindings draw.
 */
const sinkEntriesShape0 = prop.ofShape<Record<string, boolean>>('{"name":"Record<string, boolean>","shape":"ref"}');

/**
 * The generator of the values of Box<number> that the generated bindings
 * draw.
 */
const sinkBShape = prop.ofShape<Box<number>>('{"definitions":{"Box<number>":{"fields":[["value",{"shape":"float","width":64}]],"shape":"record"}},"name":"Box<number>","shape":"ref"}');

/**
 * The generator of the values of Point that the generated bindings draw.
 */
const sinkPShape1 = prop.ofShape<Point>('{"definitions":{"Point":{"fields":[["x",{"shape":"float","width":64}],["y",{"shape":"float","width":64}]],"shape":"record"}},"name":"Point","shape":"ref"}');

/**
 * The generator of the values of Signed that the generated bindings draw.
 */
const sinkSShape0 = prop.ofShape<Signed>('{"shape":"literal","values":[{"type":"float","value":-1},{"type":"float","value":0}]}');

/**
 * The generator of the values of string that the generated bindings draw.
 */
const sinkItemsShape2 = prop.ofShape<string>('{"shape":"string"}');

/**
 * The generator of the values of Triple that the generated bindings draw.
 */
const sinkTShape = prop.ofShape<Triple>('{"definitions":{"Item":{"fields":[["id",{"shape":"string"}],["count",{"shape":"float","width":64}],["label",{"of":{"shape":"string"},"shape":"optional"}]],"shape":"record"},"Item2":{"fields":[["name",{"shape":"string"}]],"shape":"record"},"Item3":{"fields":[["flag",{"shape":"bool"}]],"shape":"record"},"Triple":{"fields":[["a",{"name":"Item","shape":"ref"}],["b",{"name":"Item2","shape":"ref"}],["c",{"name":"Item3","shape":"ref"}]],"shape":"record"}},"name":"Triple","shape":"ref"}');

/**
 * The generator of the values of Box<sinkPShape> that the generated
 * bindings draw.
 */
const sinkPShape2 = prop.ofShape<Box<sinkPShape>>('{"definitions":{"Box<sinkPShape>":{"fields":[["value",{"name":"sinkPShape","shape":"ref"}]],"shape":"record"},"sinkPShape":{"fields":[["value",{"shape":"string"}]],"shape":"record"}},"name":"Box<sinkPShape>","shape":"ref"}');

/**
 * The generator of the values of Tag that the generated bindings draw: its
 * starting value, or a value of its shape.
 */
const sinkTShape0 = prop.oneOf<Tag>(prop.just(newTag()), prop.ofShape<Tag>('{"shape":"string"}'));

/**
 * The generator of the values of Names that the generated bindings draw.
 */
const sinkNShape0 = prop.ofShape<Names>('{"of":{"shape":"string"},"shape":"list"}');

/**
 * Returns the checks of Sink, the claims that the generator refuses for
 * its callables, and the classifications of its callables without a check.
 */
export function sinkSuite(): Suite<Sink> {
  return {
    declaration: "Sink",
    package: "shapes",
    generator: { name: "dokimi", version: "0.1.0" },
    checks: [
      {
        id: sinkChecks.flag.smoke(),
        statement: "flag returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().flag(true));
        },
        defect: sinkFlagPanics,
      },
      {
        id: generated(sinkChecks.flag.smoke()),
        statement: "flag returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnOn = c.draw(sinkOnShape, "on");
          await survives(h, () => h.subject().flag(drawnOn));
        },
        defect: sinkFlagPanics,
      },
      {
        id: sinkChecks.bytes.smoke(),
        statement: "bytes returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().bytes(new Uint8Array([1])));
        },
        defect: sinkBytesPanics,
      },
      {
        id: generated(sinkChecks.bytes.smoke()),
        statement: "bytes returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnData = c.draw(sinkDataShape, "data");
          await survives(h, () => h.subject().bytes(drawnData));
        },
        defect: sinkBytesPanics,
      },
      {
        id: sinkChecks.texts.smoke(),
        statement: "texts returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().texts(["items-a"]));
        },
        defect: sinkTextsPanics,
      },
      {
        id: generated(sinkChecks.texts.smoke()),
        statement: "texts returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnItems = c.draw(sinkItemsShape, "items");
          await survives(h, () => h.subject().texts(drawnItems));
        },
        defect: sinkTextsPanics,
      },
      {
        id: sinkChecks.numbers.smoke(),
        statement: "numbers returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().numbers([1.5]));
        },
        defect: sinkNumbersPanics,
      },
      {
        id: generated(sinkChecks.numbers.smoke()),
        statement: "numbers returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnItems = c.draw(sinkItemsShape0, "items");
          await survives(h, () => h.subject().numbers(drawnItems));
        },
        defect: sinkNumbersPanics,
      },
      {
        id: sinkChecks.frozen.smoke(),
        statement: "frozen returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().frozen(["items-a"]));
        },
        defect: sinkFrozenPanics,
      },
      {
        id: generated(sinkChecks.frozen.smoke()),
        statement: "frozen returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnItems = c.draw(sinkItemsShape1, "items");
          await survives(h, () => h.subject().frozen(drawnItems));
        },
        defect: sinkFrozenPanics,
      },
      {
        id: sinkChecks.count.smoke(),
        statement: "count returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().count(42n));
        },
        defect: sinkCountPanics,
      },
      {
        id: generated(sinkChecks.count.smoke()),
        statement: "count returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnN = c.draw(sinkNShape, "n");
          await survives(h, () => h.subject().count(drawnN));
        },
        defect: sinkCountPanics,
      },
      {
        id: sinkChecks.level.smoke(),
        statement: "level returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().level(0 as Level));
        },
        defect: sinkLevelPanics,
      },
      {
        id: generated(sinkChecks.level.smoke()),
        statement: "level returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnL = c.draw(sinkLShape, "l");
          await survives(h, () => h.subject().level(drawnL));
        },
        defect: sinkLevelPanics,
      },
      {
        id: sinkChecks.mode.smoke(),
        statement: "mode returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().mode("read" as Mode));
        },
        defect: sinkModePanics,
      },
      {
        id: generated(sinkChecks.mode.smoke()),
        statement: "mode returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnM = c.draw(sinkMShape, "m");
          await survives(h, () => h.subject().mode(drawnM));
        },
        defect: sinkModePanics,
      },
      {
        id: sinkChecks.status.smoke(),
        statement: "status returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().status(1 as Status));
        },
        defect: sinkStatusPanics,
      },
      {
        id: generated(sinkChecks.status.smoke()),
        statement: "status returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnS = c.draw(sinkSShape, "s");
          await survives(h, () => h.subject().status(drawnS));
        },
        defect: sinkStatusPanics,
      },
      {
        id: sinkChecks.flags.smoke(),
        statement: "flags returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().flags(1 as Flags));
        },
        defect: sinkFlagsPanics,
      },
      {
        id: generated(sinkChecks.flags.smoke()),
        statement: "flags returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnF = c.draw(sinkFShape, "f");
          await survives(h, () => h.subject().flags(drawnF));
        },
        defect: sinkFlagsPanics,
      },
      {
        id: sinkChecks.key.smoke(),
        statement: "key returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().key("k-a"));
        },
        defect: sinkKeyPanics,
      },
      {
        id: generated(sinkChecks.key.smoke()),
        statement: "key returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnK = c.draw(sinkKShape, "k");
          await survives(h, () => h.subject().key(drawnK));
        },
        defect: sinkKeyPanics,
      },
      {
        id: sinkChecks.choice.smoke(),
        statement: "choice returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().choice("c-a"));
        },
        defect: sinkChoicePanics,
      },
      {
        id: generated(sinkChecks.choice.smoke()),
        statement: "choice returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnC0 = c.draw(sinkC0Shape, "c");
          await survives(h, () => h.subject().choice(drawnC0));
        },
        defect: sinkChoicePanics,
      },
      {
        id: sinkChecks.pair.smoke(),
        statement: "pair returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().pair({ left: { id: "id-a", count: 1.5 }, right: { name: "name-a" } }));
        },
        defect: sinkPairPanics,
      },
      {
        id: generated(sinkChecks.pair.smoke()),
        statement: "pair returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnP = c.draw(sinkPShape0, "p");
          await survives(h, () => h.subject().pair(drawnP));
        },
        defect: sinkPairPanics,
      },
      {
        id: sinkChecks.lookup.smoke(),
        statement: "lookup returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().lookup(new Map<string, number>([["entries-a", 1.5]])));
        },
        defect: sinkLookupPanics,
      },
      {
        id: generated(sinkChecks.lookup.smoke()),
        statement: "lookup returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnEntries = c.draw(sinkEntriesShape, "entries");
          await survives(h, () => h.subject().lookup(drawnEntries));
        },
        defect: sinkLookupPanics,
      },
      {
        id: sinkChecks.table.smoke(),
        statement: "table returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().table({ ["entries-a"]: true }));
        },
        defect: sinkTablePanics,
      },
      {
        id: generated(sinkChecks.table.smoke()),
        statement: "table returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnEntries = c.draw(sinkEntriesShape0, "entries");
          await survives(h, () => h.subject().table(drawnEntries));
        },
        defect: sinkTablePanics,
      },
      {
        id: sinkChecks.box.smoke(),
        statement: "box returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().box({ value: 1.5 }));
        },
        defect: sinkBoxPanics,
      },
      {
        id: generated(sinkChecks.box.smoke()),
        statement: "box returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnB = c.draw(sinkBShape, "b");
          await survives(h, () => h.subject().box(drawnB));
        },
        defect: sinkBoxPanics,
      },
      {
        id: sinkChecks.point.smoke(),
        statement: "point returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().point(Object.assign(new Point(), { x: 1.5 })));
        },
        defect: sinkPointPanics,
      },
      {
        id: generated(sinkChecks.point.smoke()),
        statement: "point returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnP = c.draw(sinkPShape1, "p");
          await survives(h, () => h.subject().point(drawnP));
        },
        defect: sinkPointPanics,
      },
      {
        id: sinkChecks.signed.smoke(),
        statement: "signed returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().signed(-1 as Signed));
        },
        defect: sinkSignedPanics,
      },
      {
        id: generated(sinkChecks.signed.smoke()),
        statement: "signed returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnS = c.draw(sinkSShape0, "s");
          await survives(h, () => h.subject().signed(drawnS));
        },
        defect: sinkSignedPanics,
      },
      {
        id: sinkChecks.many.smoke(),
        statement: "many returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().many("items-a"));
        },
        defect: sinkManyPanics,
      },
      {
        id: generated(sinkChecks.many.smoke()),
        statement: "many returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnItems = c.draw(sinkItemsShape2, "items");
          await survives(h, () => h.subject().many(drawnItems));
        },
        defect: sinkManyPanics,
      },
      {
        id: sinkChecks.maybe.smoke(),
        statement: "maybe returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().maybe("label-a"));
        },
        defect: sinkMaybePanics,
      },
      {
        id: generated(sinkChecks.maybe.smoke()),
        statement: "maybe returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnLabel = c.draw(sinkItemsShape2, "label");
          await survives(h, () => h.subject().maybe(drawnLabel));
        },
        defect: sinkMaybePanics,
      },
      {
        id: sinkChecks.triple.smoke(),
        statement: "triple returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().triple({ a: { id: "id-a", count: 1.5 }, b: { name: "name-a" }, c: { flag: true } }));
        },
        defect: sinkTriplePanics,
      },
      {
        id: generated(sinkChecks.triple.smoke()),
        statement: "triple returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnT = c.draw(sinkTShape, "t");
          await survives(h, () => h.subject().triple(drawnT));
        },
        defect: sinkTriplePanics,
      },
      {
        id: sinkChecks.named.smoke(),
        statement: "named returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().named({ value: Object.assign(new sinkPShape(), { value: "value-a" }) }));
        },
        defect: sinkNamedPanics,
      },
      {
        id: generated(sinkChecks.named.smoke()),
        statement: "named returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnP = c.draw(sinkPShape2, "p");
          await survives(h, () => h.subject().named(drawnP));
        },
        defect: sinkNamedPanics,
      },
      {
        id: sinkChecks.tag.smoke(),
        statement: "tag returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().tag(newTag()));
        },
        defect: sinkTagPanics,
      },
      {
        id: generated(sinkChecks.tag.smoke()),
        statement: "tag returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnT = c.draw(sinkTShape0, "t");
          await survives(h, () => h.subject().tag(drawnT));
        },
        defect: sinkTagPanics,
      },
      {
        id: sinkChecks.names.smoke(),
        statement: "names returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().names(["n-a"]));
        },
        defect: sinkNamesPanics,
      },
      {
        id: generated(sinkChecks.names.smoke()),
        statement: "names returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnN = c.draw(sinkNShape0, "n");
          await survives(h, () => h.subject().names(drawnN));
        },
        defect: sinkNamesPanics,
      },
      {
        id: sinkChecks.unset.smoke(),
        statement: "unset returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().unset("label-a"));
        },
        defect: sinkUnsetPanics,
      },
      {
        id: generated(sinkChecks.unset.smoke()),
        statement: "unset returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnLabel = c.draw(sinkItemsShape2, "label");
          await survives(h, () => h.subject().unset(drawnLabel));
        },
        defect: sinkUnsetPanics,
      },
    ],
    refusals: [],
    unclaimed: [
      {
        classification: "shape:writer",
        callable: "Sink.flag",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.bytes",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.texts",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.numbers",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.frozen",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.count",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.level",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.mode",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.status",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.flags",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.limit",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.odd",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.key",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.tree",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.source",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.choice",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.pair",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.lookup",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.table",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.box",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.point",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.any",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.signed",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.many",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.maybe",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.triple",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.named",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.tag",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.step",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.node",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.names",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Sink.unset",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
    ],
  };
}

/**
 * Registers the application of the checks of Sink, and of the checks that
 * options add, to a fresh subject of newSubject for each check.
 */
export function runSinkSuite(newSubject: () => Sink | Promise<Sink>, options?: SinkOptions): void {
  runSinkSuiteFromHandle(() => newSubject(), options);
}

/**
 * Registers the same application to a subject that newSubject builds from
 * the handle of each check.
 */
export function runSinkSuiteFromHandle(newSubject: (h: SinkHandle) => Sink | Promise<Sink>, options?: SinkOptions): void {
  run(sinkSuite(), newSubject, options);
}

/**
 * Returns the planted defect of a Sink whose flag panics over a fresh
 * subject.
 */
async function sinkFlagPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: () => {
      throw new TypeError("the planted defect of a Sink whose flag panics");
    },
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose bytes panics over a fresh
 * subject.
 */
async function sinkBytesPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: () => {
      throw new TypeError("the planted defect of a Sink whose bytes panics");
    },
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose texts panics over a fresh
 * subject.
 */
async function sinkTextsPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: () => {
      throw new TypeError("the planted defect of a Sink whose texts panics");
    },
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose numbers panics over a fresh
 * subject.
 */
async function sinkNumbersPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: () => {
      throw new TypeError("the planted defect of a Sink whose numbers panics");
    },
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose frozen panics over a fresh
 * subject.
 */
async function sinkFrozenPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: () => {
      throw new TypeError("the planted defect of a Sink whose frozen panics");
    },
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose count panics over a fresh
 * subject.
 */
async function sinkCountPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: () => {
      throw new TypeError("the planted defect of a Sink whose count panics");
    },
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose level panics over a fresh
 * subject.
 */
async function sinkLevelPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: () => {
      throw new TypeError("the planted defect of a Sink whose level panics");
    },
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose mode panics over a fresh
 * subject.
 */
async function sinkModePanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: () => {
      throw new TypeError("the planted defect of a Sink whose mode panics");
    },
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose status panics over a fresh
 * subject.
 */
async function sinkStatusPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: () => {
      throw new TypeError("the planted defect of a Sink whose status panics");
    },
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose flags panics over a fresh
 * subject.
 */
async function sinkFlagsPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: () => {
      throw new TypeError("the planted defect of a Sink whose flags panics");
    },
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose key panics over a fresh
 * subject.
 */
async function sinkKeyPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: () => {
      throw new TypeError("the planted defect of a Sink whose key panics");
    },
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose choice panics over a fresh
 * subject.
 */
async function sinkChoicePanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: () => {
      throw new TypeError("the planted defect of a Sink whose choice panics");
    },
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose pair panics over a fresh
 * subject.
 */
async function sinkPairPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: () => {
      throw new TypeError("the planted defect of a Sink whose pair panics");
    },
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose lookup panics over a fresh
 * subject.
 */
async function sinkLookupPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: () => {
      throw new TypeError("the planted defect of a Sink whose lookup panics");
    },
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose table panics over a fresh
 * subject.
 */
async function sinkTablePanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: () => {
      throw new TypeError("the planted defect of a Sink whose table panics");
    },
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose box panics over a fresh
 * subject.
 */
async function sinkBoxPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: () => {
      throw new TypeError("the planted defect of a Sink whose box panics");
    },
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose point panics over a fresh
 * subject.
 */
async function sinkPointPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: () => {
      throw new TypeError("the planted defect of a Sink whose point panics");
    },
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose signed panics over a fresh
 * subject.
 */
async function sinkSignedPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: () => {
      throw new TypeError("the planted defect of a Sink whose signed panics");
    },
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose many panics over a fresh
 * subject.
 */
async function sinkManyPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: () => {
      throw new TypeError("the planted defect of a Sink whose many panics");
    },
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose maybe panics over a fresh
 * subject.
 */
async function sinkMaybePanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: () => {
      throw new TypeError("the planted defect of a Sink whose maybe panics");
    },
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose triple panics over a fresh
 * subject.
 */
async function sinkTriplePanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: () => {
      throw new TypeError("the planted defect of a Sink whose triple panics");
    },
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose named panics over a fresh
 * subject.
 */
async function sinkNamedPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: () => {
      throw new TypeError("the planted defect of a Sink whose named panics");
    },
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose tag panics over a fresh
 * subject.
 */
async function sinkTagPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: () => {
      throw new TypeError("the planted defect of a Sink whose tag panics");
    },
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose names panics over a fresh
 * subject.
 */
async function sinkNamesPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: () => {
      throw new TypeError("the planted defect of a Sink whose names panics");
    },
    unset: s.unset.bind(s),
  };
}

/**
 * Returns the planted defect of a Sink whose unset panics over a fresh
 * subject.
 */
async function sinkUnsetPanics(h: SinkHandle): Promise<Sink> {
  const s = await h.fresh();
  return {
    flag: s.flag.bind(s),
    bytes: s.bytes.bind(s),
    texts: s.texts.bind(s),
    numbers: s.numbers.bind(s),
    frozen: s.frozen.bind(s),
    count: s.count.bind(s),
    level: s.level.bind(s),
    mode: s.mode.bind(s),
    status: s.status.bind(s),
    flags: s.flags.bind(s),
    limit: s.limit.bind(s),
    odd: s.odd.bind(s),
    key: s.key.bind(s),
    tree: s.tree.bind(s),
    source: s.source.bind(s),
    choice: s.choice.bind(s),
    pair: s.pair.bind(s),
    lookup: s.lookup.bind(s),
    table: s.table.bind(s),
    box: s.box.bind(s),
    point: s.point.bind(s),
    any: s.any.bind(s),
    signed: s.signed.bind(s),
    many: s.many.bind(s),
    maybe: s.maybe.bind(s),
    triple: s.triple.bind(s),
    named: s.named.bind(s),
    tag: s.tag.bind(s),
    step: s.step.bind(s),
    node: s.node.bind(s),
    names: s.names.bind(s),
    unset: () => {
      throw new TypeError("the planted defect of a Sink whose unset panics");
    },
  };
}

void compatV1;
