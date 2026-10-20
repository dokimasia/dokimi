// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the field-directive half of the builder's defaults handling.
 *
 * @remarks
 * The generator seeds the builder from one of three sources: a defaults
 * companion, a per-field //dokimi:default directive, or the value of
 * `new Config()`. Each produces a different constructor, so they are three
 * fixtures rather than one — the domain fixture carries the companion, this
 * carries the per-field form, and the plain fixture carries neither.
 *
 * There is deliberately no defaults function here. Adding one would make the
 * fixture test precedence between the two sources rather than the directive
 * path itself, and precedence is a question for a composite.
 *
 * Routing is declared once for the module rather than repeated on each
 * class: every builder here lands in the same companion directory, so a
 * per-class directive is the same statement written N times, and the Nth
 * copy is the one that gets forgotten.
 *
 * @packageDocumentation
 */
//dokimi:out path=configtest/

/**
 * Represents a config with a default on every literal kind that the
 * directive accepts.
 *
 * @remarks
 * Each kind is parsed differently: quoting for strings, numbers, keywords for
 * booleans, and undefined, which has no literal of its own. Each directive is
 * the trailing comment of its field.
 */
//dokimi:builder
export class Config {
  /** Contains the host. */
  host: string = ""; //dokimi:default "localhost"

  /** Contains the port. */
  port: number = 0; //dokimi:default 8080

  /** Is true when the client logs every request. */
  verbose: boolean = false; //dokimi:default true

  /** Contains the share of the requests that the client samples. */
  ratio: number = 0; //dokimi:default 0.75

  /**
   * Counts the retries of a request.
   *
   * @remarks
   * The zero default is stated explicitly. A generator that treats "zero" as
   * "no directive" emits the same constructor either way, and this is the
   * field that catches it.
   */
  retries: number = 0; //dokimi:default 0

  /**
   * Contains the fallback host, where the config has one.
   *
   * @remarks
   * The default is undefined. Its literal is the same as the value of a new
   * config, so it tests that the directive is read rather than inferred.
   */
  fallback?: string; //dokimi:default undefined

  /**
   * Contains a field without a directive.
   *
   * @remarks
   * The field must keep the value of a new config even though its neighbours
   * do not.
   */
  undirected: string = "";
}
