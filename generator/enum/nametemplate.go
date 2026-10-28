// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"fmt"
	"strings"

	"go.dokimi.dev/dokimi/internal/nametemplate"
)

// enumField is the one field that a name template reads: the name of the
// enumeration, as in {{.Enum}}.
const enumField = "Enum"

// NameTemplate is the value of an option that sets the name of a function, a
// variable or a class of the package of an enumeration. It is a template in
// the syntax of text/template whose only action is {{.Enum}}, the name of the
// enumeration, so Parse{{.Enum}} renders ParseStatus for Status. The text
// around the actions contains letters, digits, spaces, hyphens and
// underscores, and the text before the first action starts with a letter. A
// name template uses {{.Enum}} at least once, because a package with two
// enumerations would otherwise declare one name twice.
//
// Each language writes the rendered name in its own case convention. Go
// writes a function in PascalCase, and an error as a variable with the prefix
// Err. TypeScript writes a function in camelCase, and an error as a class with
// the suffix Error. So Unknown{{.Enum}} renders ErrUnknownStatus in Go and
// UnknownStatusError in TypeScript.
type NameTemplate string

// UnmarshalText sets t to text. It returns an error for a text that is not a
// name template, and leaves t unchanged. The decoder of a config file calls
// it, so the build of a composition fails on such a text.
func (t *NameTemplate) UnmarshalText(text []byte) error {
	parts, err := segments(string(text))
	if err != nil {
		return err
	}
	if len(parts) == 1 {
		return fmt.Errorf("enum: the name template %q does not use {{.Enum}}", text)
	}
	*t = NameTemplate(text)
	return nil
}

// MethodTemplate is the value of an option that sets the name of a method of
// an enumeration. It is a name template that need not use {{.Enum}}, because
// each enumeration declares its own methods. IsValid and IsA{{.Enum}} are
// both method templates. Go writes the rendered name in PascalCase.
type MethodTemplate string

// UnmarshalText sets t to text. It returns an error for a text that is not a
// method template, and leaves t unchanged. The decoder of a config file calls
// it, so the build of a composition fails on such a text.
func (t *MethodTemplate) UnmarshalText(text []byte) error {
	if _, err := segments(string(text)); err != nil {
		return err
	}
	*t = MethodTemplate(text)
	return nil
}

// render returns the name that the name template text renders for the
// enumeration named enum. It returns the error of segments.
func render(text, enum string) (string, error) {
	parts, err := segments(text)
	if err != nil {
		return "", err
	}
	return strings.Join(parts, enum), nil
}

// segments returns the texts of the name template text before, between and
// after its actions {{.Enum}}, as [nametemplate.Segments] returns them. The
// message of its error opens with enum.
func segments(text string) ([]string, error) {
	parts, err := nametemplate.Segments(text, enumField)
	if err != nil {
		return nil, fmt.Errorf("enum: %w", err)
	}
	return parts, nil
}
