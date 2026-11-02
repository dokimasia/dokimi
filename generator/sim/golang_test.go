// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi/generator/sim"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// shopModule is the module of the Go tree.
const shopModule = "example.com/shop"

// The part of Go declares the target golang with the templates of the two
// families and the helpers of the templates.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the templates of the families for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, sim.New(), golang.Target, simTemplate, testTemplate)
	})

	t.Run("writes the simulations of the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t, registerPackage, registerSimPackage, plainPackage, plainSimPackage,
			genericPackage, genericSimPackage, remotePackage, remoteSimPackage, storePackage), golden.ShouldUpdate())
	})

	t.Run("writes the calls of a dependency that pass go vet then go test", func(t *testing.T) {
		t.Parallel()

		a := gotesting.New()
		if !toolchain.Require(t, a) {
			return
		}
		project := gentest.Project(rendered(t, remotePackage, remoteSimPackage, storePackage), goTree, shopModule)
		gotesting.AssertVets(t.Context(), t, project)
		toolchain.AssertTestsPass(t.Context(), t, a, project)
	})

	t.Run("writes the cases of a generic reply that pass go vet", func(t *testing.T) {
		t.Parallel()

		if !toolchain.Require(t, gotesting.New()) {
			return
		}
		gotesting.AssertVets(t.Context(), t,
			gentest.Project(rendered(t, genericPackage, genericSimPackage), goTree, shopModule))
	})

	t.Run("writes Go that passes go vet then go test", func(t *testing.T) {
		t.Parallel()

		a := gotesting.New()
		if !toolchain.Require(t, a) {
			return
		}
		project := gentest.Project(rendered(t, registerPackage, registerSimPackage), goTree, shopModule)
		gotesting.AssertVets(t.Context(), t, project)
		toolchain.AssertTestsPass(t.Context(), t, a, project)
	})

	t.Run("refuses the field of a record that the record lacks", func(t *testing.T) {
		t.Parallel()

		templates, provides := sim.New().(plugin.TemplateProvider)
		assert.True(t, provides, "sim declares templates")
		field, declared := templates.TemplateFuncs(golang.Target)["simField"].(func(*emit.Struct, string) (
			*emit.TypeRef, error))
		assert.True(t, declared, "the target golang declares the helper simField")
		_, err := field(&emit.Struct{Name: "MachineSimulation"}, "simulation")
		assert.HasError(t, err, "a template that reads a field that the record lacks fails")
		expect.Equal(t, err.Error(), "sim: the record MachineSimulation has no field simulation",
			"the error names the record and the field")
	})
}

// rendered runs the annotator fault and both halves of sim over the packages
// pkgs, and returns the files that the Go backend renders from them.
func rendered(t *testing.T, pkgs ...string) fstest.MapFS {
	t.Helper()

	r := generateWithFaults(t, pkgs...)
	return gentest.Rendered(t, r.f, r.res, r.p, gobackend.New())
}
