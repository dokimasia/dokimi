// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the string-valued half of the enum corpus.
 *
 * @remarks
 * A string enum's declared value *is* its textual form, and it is already
 * written down. A generator deriving the identifier instead — `US` rather
 * than `us-east` — still round-trips the render against the parse, so its own
 * checks pass while every value arriving from JSON, a database column or a
 * query parameter fails to parse. Only a fixture whose values differ from
 * its identifiers can tell the two apart, which is why none of these matches
 * its name.
 *
 * The fallback differs too: a value outside the set renders as itself rather
 * than as a number, which is the more useful diagnostic anyway.
 *
 * @packageDocumentation
 */

/**
 * Represents a region of a cloud.
 *
 * @remarks
 * Each value differs from its identifier, and none is a prefix of another — a
 * prefix would make the fallback indistinguishable from a declared value
 * under a naive match.
 */
//dokimi:enum
export enum Region {
  US = "us-east",
  EU = "eu-west",
  AP = "ap-south",
}
