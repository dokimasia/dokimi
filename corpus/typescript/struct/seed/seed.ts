// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares values and a companion that the builders of another module draw
 * from, so the cross-module half of default resolution has somewhere to
 * point.
 *
 * @remarks
 * It declares no builder of its own. Its whole purpose is to be referenced
 * from outside, which is the case a single-module fixture cannot show.
 *
 * @packageDocumentation
 */

/** Contains a default that a field in another module names. */
export const region = "eu-west";

/**
 * Contains a numeric default.
 *
 * @remarks
 * The constant is present so the cross-module path is exercised for more
 * than one literal kind.
 */
export const retries = 3;

/**
 * Represents the class that {@link configDefaults} seeds.
 *
 * @remarks
 * It is declared here rather than beside the builder so the companion is
 * genuinely cross-module.
 */
export class Config {
  /** Contains the region. */
  region: string = "";
  /** Counts the retries. */
  retries: number = 0;
  /** Contains the label. */
  label: string = "";
}

/** Returns the seed that a builder in another module seeds from. */
export function configDefaults(): Config {
  const config = new Config();
  config.region = region;
  config.retries = retries;
  config.label = "seeded";
  return config;
}
