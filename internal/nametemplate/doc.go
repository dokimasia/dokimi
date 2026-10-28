// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package nametemplate parses the name templates of the options and the
// directive params of the generators of dokimi.
//
// A name template is a template in the syntax of text/template whose only
// action is one field of the generator, such as {{.Enum}} or {{.Method}}. The
// text around the actions contains letters, digits, spaces, hyphens and
// underscores, and the text before the first action starts with a letter.
// [Segments] returns the texts around the actions, so a generator renders a
// name as these texts joined by the value of the field. The part of a
// language then writes the name in its own case convention.
//
// # Dependency position
//
// The package imports the Go stdlib alone.
package nametemplate
