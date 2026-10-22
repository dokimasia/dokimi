// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Registers the assertions of `@dokimi/assert` on the `expect` of vitest, so
 * the generated tests call them as `expect(got).isEqualTo(want, message)`.
 */

import { register } from "@dokimi/assert/vitest";

register();
