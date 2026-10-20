// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the partition mixin, which states that
 * state is isolated per partition key.
 *
 * @remarks
 * The interface carries whatever methods the law needs to be stateable. A
 * mixin whose law spans two calls cannot be hosted by a single method: there
 * would be nothing to compare against, and the generated subtest would pass
 * by having nothing to check.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=partitiontest/
export interface Mixed {
  /**
   * Stores the value under the key in the partition.
   *
   * @remarks
   * Two keys in different partitions never interfere. The law needs a
   * partition parameter distinct from the key, or there is nothing to isolate
   * by.
   */
  //dokimi:mixin partition read=read axis=partition
  put(partition: string, key: string, value: string): Promise<void>;

  /**
   * Returns the value under the key in the partition.
   *
   * @remarks
   * The law observes isolation by reading the same key from another
   * partition.
   */
  read(partition: string, key: string): Promise<string>;
}
