// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the authored check of {@link Slugger} with its planted defect.
 *
 * @remarks
 * The planted defect has the name survives, which the suite also imports
 * from the probe of smoke, so the suite imports one of the two under a name
 * of its own.
 *
 * @packageDocumentation
 */

import { check } from "@dokimi/assert";
import type { Slugger } from "../iface";
import type { SluggerHandle } from "./iface.suite";

/** Contains the titles that the check of slug reads. */
const titles = ["Hello World", "  Go  Dokimi  ", "one"];

/** Returns the slug of title that the sequential specification of Slugger states. */
function specOf(title: string): string {
  return title
    .split(/\s+/)
    .filter((word) => word !== "")
    .map((word) => word.toLowerCase())
    .join("-");
}

/** Checks that slug returns the slug that the specification states for each title. */
//dokimi:check Slugger.slug defect=survives
export async function followsTheSpecification(h: SluggerHandle): Promise<void> {
  for (const title of titles) {
    check.isEqualTo(h, h.subject().slug(title), specOf(title), h.statement());
  }
}

/** Returns a slugger whose slug keeps the case of the words. */
export function survives(_h: SluggerHandle): Slugger {
  return {
    slug: (title) =>
      title
        .split(/\s+/)
        .filter((word) => word !== "")
        .join("-"),
  };
}
