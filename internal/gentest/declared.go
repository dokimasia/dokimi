// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"io/fs"
	"slices"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/plugin"
)

// Code is one code of the findings of a plugin, with its spelling and its
// meaning.
type Code struct {
	// Name is the name of the variable of the code.
	Name string
	// WantCode is the spelling of the code: its prefix, a hyphen and its
	// number.
	WantCode string
	// WantMeaning is the meaning that the registry of the codes contains.
	WantMeaning string
	// Give is the code.
	Give diag.Code
}

// Templates fails the test where the plugin p declares no template tree for
// target, or where the tree lacks a template of names.
func Templates(tb assert.TB, p plugin.Plugin, target plugin.Target, names ...string) {
	tb.Helper()

	templates, provides := p.(plugin.TemplateProvider)
	assert.True(tb, provides, "the plugin declares templates")
	expect.Contains(tb, templates.TemplateTargets(), target, "the plugin serves the target "+string(target))
	tree, declared := templates.Templates(target)
	assert.True(tb, declared, "the target "+string(target)+" has a template tree")
	for _, name := range names {
		_, err := fs.Stat(tree, name)
		expect.NoError(tb, err, "the tree of "+string(target)+" contains "+name)
	}
}

// External fails the test where the plugin p does not declare the family of
// the tag, or where the Go files of that family are not in the external test
// package of their directory.
func External(tb assert.TB, p plugin.Plugin, tag string) {
	tb.Helper()

	families, declares := p.(plugin.OutputProvider)
	assert.True(tb, declares, "the plugin declares its families")
	at := slices.IndexFunc(families.Outputs(), func(o plugin.Output) bool { return o.Tag == tag })
	assert.NotEqual(tb, at, -1, "the plugin declares the family of the tag "+tag)
	expect.True(tb, families.Outputs()[at].External, "the files of the family test their package from outside")
}

// Codes fails the test for each code of codes that does not have its
// spelling, or whose meaning differs from the meaning in the registry of the
// codes of eidos. A package registers its codes when it initializes.
func Codes(tb assert.TB, codes ...Code) {
	tb.Helper()

	for _, c := range codes {
		expect.Equal(tb, c.Give.String(), c.WantCode, c.Name+" has the prefix of its package and its number")
		meaning, registered := diag.Kernel().Meaning(c.Give)
		expect.True(tb, registered, "the package registers "+c.Name)
		expect.Equal(tb, meaning, c.WantMeaning, "the meaning of "+c.Name+" anchors the code in the index")
	}
}
