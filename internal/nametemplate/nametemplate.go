// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package nametemplate

import (
	"bytes"
	"fmt"
	"text/template"
	"text/template/parse"
	"unicode"
	"unicode/utf8"
)

// templateName is the name of a name template in the error of a template that
// does not parse.
const templateName = "name"

// Segments returns the texts of the name template text before, between and
// after its actions {{.F}}, where F is field. A name is these texts joined by
// the value of the field, so Segments returns one text for a template without
// an action. Segments returns an error for a template that does not parse or
// that is empty. It returns an error for an action
// other than {{.F}}, for a character other than a letter, a digit, a space, a
// hyphen or an underscore, and for a first text that does not start with a
// letter. The message of each error contains the template and the rule that
// the template breaks.
func Segments(text, field string) ([]string, error) {
	t, err := template.New(templateName).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("the name template %q does not parse: %w", text, err)
	}
	if len(t.Root.Nodes) == 0 {
		return nil, fmt.Errorf("the name template %q is empty", text)
	}
	parts := []string{""}
	for _, n := range t.Root.Nodes {
		switch n := n.(type) {
		case *parse.TextNode:
			if at := bytes.IndexFunc(n.Text, notInName); at >= 0 {
				r, _ := utf8.DecodeRune(n.Text[at:])
				return nil, fmt.Errorf("the name template %q contains the character %q", text, r)
			}
			parts[len(parts)-1] += string(n.Text)
			continue
		case *parse.ActionNode:
			if reads(n.Pipe, field) {
				parts = append(parts, "")
				continue
			}
		}
		return nil, fmt.Errorf("the name template %q has an action other than {{.%s}}", text, field)
	}
	if first, _ := utf8.DecodeRuneInString(parts[0]); parts[0] != "" && !unicode.IsLetter(first) {
		return nil, fmt.Errorf("the name template %q starts with %q instead of a letter", text, first)
	}
	return parts, nil
}

// reads reports whether the pipeline p of an action is the field named field
// alone, as in {{.Enum}}.
func reads(p *parse.PipeNode, field string) bool {
	if len(p.Decl) > 0 || len(p.Cmds) != 1 || len(p.Cmds[0].Args) != 1 {
		return false
	}
	f, isField := p.Cmds[0].Args[0].(*parse.FieldNode)
	return isField && len(f.Ident) == 1 && f.Ident[0] == field
}

// notInName reports whether the text of a name template cannot contain r: a
// character other than a letter, a digit, a space, a hyphen or an underscore.
func notInName(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '-' && r != '_'
}
