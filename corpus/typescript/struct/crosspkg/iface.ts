// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture for a default that names a symbol outside the
 * declaring module.
 *
 * @remarks
 * A default resolves a name against the imports of its file, which reads
 * naturally and is what an author writes for a module the file already uses
 * — retries and region declare theirs this way, and zone names the same
 * constant through a namespace import. Both resolve for a string-typed field
 * exactly as they do for a numeric one.
 *
 * The seed fixture holds the far side.
 *
 * Routing is declared once for the module rather than repeated on each
 * class: every builder here lands in the same companion directory, so a
 * per-class directive is the same statement written N times, and the Nth
 * copy is the one that gets forgotten.
 *
 * @packageDocumentation
 */
//dokimi:out path=crosspkgtest/

import { Config, configDefaults, region, retries } from "../seed/seed";
import * as seed from "../seed/seed";

/** Represents settings that draw every default from another module. */
//dokimi:builder
export class Settings {
  /** Counts the retries, with a default that names a constant of the imported module. */
  //dokimi:default retries
  retries: number = 0;

  /** Contains the region, with a default that names a constant of the imported module. */
  //dokimi:default region
  region: string = "";

  /** Contains the zone, with a default that names the same constant through the namespace import. */
  //dokimi:default seed.region
  zone: string = "";

  /**
   * Contains the deadline.
   *
   * @remarks
   * The field exists so the class has a field of a library type beside the
   * defaulted ones.
   */
  deadline: Date = new Date(0);
}

/** Represents an alias that seeds from a companion of another module, named by an import. */
//dokimi:builder defaults=configDefaults
export type Mirrored = Config;

/**
 * Represents the same companion named through the namespace import.
 *
 * @remarks
 * It exists because the two notations take different routes: a named import
 * binds the function itself, while a namespace import is only meaningful
 * against the import declarations of the file that wrote it.
 */
//dokimi:builder defaults=seed.configDefaults
export type Shorthand = Config;

/** Lists the companion and the constants, so each import is used by code. */
export const used = [configDefaults, region, retries, seed.region] as const;
