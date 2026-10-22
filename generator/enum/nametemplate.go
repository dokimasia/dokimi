// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
	"unicode"
	"unicode/utf8"
)

// enumField is the one field that a name template reads: the name of the
// enumeration, as in {{.Enum}}.
const enumField = "Enum"

// templateName is the name of a name template in the error of a template that
// does not parse.
const templateName = "name"

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
// after its actions {{.Enum}}. A name is these texts joined by the name of an
// enumeration. segments returns an error for a template that does not parse
// or that is empty, and for an action other than {{.Enum}}. It also returns
// an error for a character other than a letter, a digit, a space, a hyphen or
// an underscore, and for a first text that does not start with a letter.
func segments(text string) ([]string, error) {
	t, err := template.New(templateName).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("enum: the name template %q does not parse: %w", text, err)
	}
	if len(t.Root.Nodes) == 0 {
		return nil, errors.New("enum: the name template is empty")
	}
	parts := []string{""}
	for _, n := range t.Root.Nodes {
		switch n := n.(type) {
		case *parse.TextNode:
			if at := bytes.IndexFunc(n.Text, notInName); at >= 0 {
				r, _ := utf8.DecodeRune(n.Text[at:])
				return nil, fmt.Errorf("enum: the name template %q contains the character %q", text, r)
			}
			parts[len(parts)-1] += string(n.Text)
			continue
		case *parse.ActionNode:
			if readsEnum(n.Pipe) {
				parts = append(parts, "")
				continue
			}
		}
		return nil, fmt.Errorf("enum: the name template %q has an action other than {{.Enum}}", text)
	}
	if first, _ := utf8.DecodeRuneInString(parts[0]); parts[0] != "" && !unicode.IsLetter(first) {
		return nil, fmt.Errorf("enum: the name template %q starts with %q instead of a letter", text, first)
	}
	return parts, nil
}

// readsEnum reports whether the pipeline p of an action is the field Enum
// alone, as in {{.Enum}}.
func readsEnum(p *parse.PipeNode) bool {
	if len(p.Decl) > 0 || len(p.Cmds) != 1 || len(p.Cmds[0].Args) != 1 {
		return false
	}
	field, isField := p.Cmds[0].Args[0].(*parse.FieldNode)
	return isField && len(field.Ident) == 1 && field.Ident[0] == enumField
}

// notInName reports whether the text of a name template cannot contain r: a
// character other than a letter, a digit, a space, a hyphen or an underscore.
func notInName(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '-' && r != '_'
}
