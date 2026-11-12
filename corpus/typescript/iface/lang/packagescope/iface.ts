// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for a kit that declares the names that
 * its suite imports.
 *
 * @remarks
 * The suite imports the authored check and the planted defect of the kit
 * module beside the probes of the runtime. The kit declares its planted
 * defect under the name of a probe that the suite imports, so the suite
 * imports one of the two under a name of its own.
 *
 * @packageDocumentation
 */

/** Represents a slugger that turns titles into slugs. */
//dokimi:suite
//dokimi:out path=packagescopetest/
export interface Slugger {
  /** Returns the words of title in lower case, joined by hyphens. */
  slug(title: string): string;
}
