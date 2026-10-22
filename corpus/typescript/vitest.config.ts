// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/** Runs the setup file of the assertions before each test file. */

import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    setupFiles: ["./vitest.setup.ts"],
  },
});
