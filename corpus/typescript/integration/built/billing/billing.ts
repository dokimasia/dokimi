// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares a class whose builder is in another module than the double of the
 * fixture built.
 *
 * @packageDocumentation
 */
//dokimi:out path=billingtest/

/** Represents an invoice, a class with a builder. */
//dokimi:builder
export class Invoice {
  /** Contains the number of the invoice. */
  number: string = "";
}
