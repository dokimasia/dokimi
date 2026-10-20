// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package lang runs the handlers of a plugin of dokimi with the part of the
// language of each subject.
//
// A plugin has a part for each language that it serves. The part of a
// language declares the target of the language with its template tree and
// its helpers. It also takes each step that the facts of eidos leave to the
// language. The plugin declares its own interface of parts, which embeds
// [Part], and lists one part for each language.
//
// [Register] declares the target of each part on the builder of the plugin.
// [Served] wraps a handler, so the handler runs with the part of the language
// of its match and leaves a subject of another language alone.
//
// # Dependency position
//
// The package imports the root, plugin and symbol packages of the SDK facade
// of eidos.
package lang
