// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package dokimi_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/cli"
	"go.dokimi.dev/eidos/core/ledger"
	"go.dokimi.dev/eidos/core/workspace"
	"go.dokimi.dev/eidos/core/workspace/workspacetest"
	golang "go.dokimi.dev/eidos/lang/go"
	gofrontend "go.dokimi.dev/eidos/lang/go/frontend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/lang/typescript"
	tsfrontend "go.dokimi.dev/eidos/lang/typescript/frontend"
	"go.dokimi.dev/eidos/plugin/shape/catalog"
	"go.dokimi.dev/eidos/sdk/output"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi"
	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/generator/suite"
)

// The tree of the composition and the directory of the files that the plans
// generate from it.
const (
	tree = "testdata/tree"
	want = "testdata/want"
)

// The file that the edit of the tree changes, the field of Box before and
// after the edit, and the mode of the file.
const (
	boxFile    = "catalog/account.go"
	sizeField  = "Size int"
	countField = "Count int"
	sourceMode = 0o644
)

// The extensions of the Go and the TypeScript files.
const (
	goExtension = ".go"
	tsExtension = ".ts"
)

// The composition of dokimi generates the files of testdata/want from the
// tree and passes the workspace suites. Its Go vets, and its tests pass and
// run in parallel. Every public declaration of its output has a doc comment.
func TestDokimi(t *testing.T) {
	t.Parallel()

	f := fixture(t)

	t.Run("Plans", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give plugin.ID
		}{
			{name: "returns two plans that share one builder", give: builder.ID},
			{name: "returns two plans that share one stub", give: stub.ID},
			{name: "returns two plans that share one enum", give: enum.ID},
			{name: "returns two plans that share one sentinel", give: sentinel.ID},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				plans := dokimi.Plans()
				assert.Length(t, plans, 2, "the composition has the plan go and the plan typescript")
				expect.Equal(t, generator(t, plans[0], tt.give), generator(t, plans[1], tt.give),
					"both plans run one instance of the generator", assert.ByIdentity())
			})
		}

		t.Run("returns a plan go that runs suite", func(t *testing.T) {
			t.Parallel()

			generator(t, dokimi.Plans()[0], suite.ID)
		})

		t.Run("returns a plan typescript that runs no suite", func(t *testing.T) {
			t.Parallel()

			suites := func(g plugin.Generator) bool { return g.Name() == suite.ID }
			expect.False(t, slices.ContainsFunc(dokimi.Plans()[1].Generators, suites), "the plan writes no suite")
		})

		sources := []struct {
			name string
			give string
			want workspace.Sources
		}{
			{
				name: "returns a plan go that reads the Go packages that the workspace imports", give: dokimi.GoPlan,
				want: workspace.Sources{Lang: golang.Lang},
			},
			{
				name: "returns a plan typescript that reads the modules that the workspace imports",
				give: dokimi.TypeScriptPlan, want: workspace.Sources{Lang: typescript.Lang},
			},
		}
		for _, tt := range sources {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				at := slices.IndexFunc(dokimi.Plans(), func(p workspace.Plan) bool { return p.Name == tt.give })
				assert.NotEqual(t, at, -1, "the composition has the plan "+tt.give)
				expect.Equal(t, dokimi.Plans()[at].Sources, tt.want,
					"a plan that its language alone scopes reads the packages of the language in the stores")
			})
		}
	})

	t.Run("Compose", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a composition that builds", func(t *testing.T) {
			t.Parallel()

			_, err := dokimi.Compose().Build()
			assert.NoError(t, err, "the composition of the binary builds")
		})

		tests := []struct {
			name string
			give plugin.ID
		}{
			{name: "returns a composition that runs the annotator of the Go satellite", give: golang.AnnotatorName},
			{name: "returns a composition that runs the annotator of the faults", give: fault.ID},
			{name: "returns a composition that runs the classifier of the shape catalog", give: catalog.ShapeID},
			{name: "returns a composition that runs the check of the shape catalog", give: catalog.CheckID},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				w, err := dokimi.Compose().Build()
				assert.NoError(t, err, "the composition of the binary builds")
				names := make([]plugin.ID, 0, len(w.Describe().Annotate))
				for _, c := range w.Describe().Annotate {
					names = append(names, c.Name)
				}
				assert.Contains(t, names, tt.give, "the composition runs the annotator")
			})
		}

		t.Run("writes the files of testdata/want from the tree", func(t *testing.T) {
			t.Parallel()

			workspacetest.AssertGolden(t, f, t.TempDir(), want)
		})

		t.Run("passes the workspace suite over the tree", func(t *testing.T) {
			t.Parallel()

			workspacetest.RunWorkspaceSuite(t, f)
		})

		t.Run("passes the warm-cold suite of the workspace over the tree", func(t *testing.T) {
			t.Parallel()

			workspacetest.RunWarmColdSuite(t, f)
		})

		t.Run("writes Go that vets", func(t *testing.T) {
			t.Parallel()

			if !toolchain.Require(t, gotesting.New()) {
				return
			}
			gotesting.AssertVets(t.Context(), t, generated(f, goExtension))
		})

		t.Run("writes Go whose public declarations have doc comments", func(t *testing.T) {
			t.Parallel()

			toolchain.AssertDocumented(t, gofrontend.New(nil), generated(f, goExtension))
		})

		t.Run("writes TypeScript whose public declarations have doc comments", func(t *testing.T) {
			t.Parallel()

			toolchain.AssertDocumented(t, tsfrontend.New(), generated(f, tsExtension))
		})

		t.Run("writes Go specs that pass", func(t *testing.T) {
			t.Parallel()

			a := gotesting.New()
			if !toolchain.Require(t, a) {
				return
			}
			toolchain.AssertTestsPass(t.Context(), t, a, generated(f, goExtension))
		})

		t.Run("writes Go specs that run in parallel", func(t *testing.T) {
			t.Parallel()

			gotesting.AssertParallel(t, generated(f, goExtension))
		})
	})
}

// fixture returns the fixture of the workspace suites over the tree. Its
// composition is Base with a disk sink and a ledger at the root of a check
// and the config file of the root, as the binary applies them. Its stores
// are the stores that the frontends of the composition locate in the
// environment, as the command line passes them. The fixture runs the plans
// of Plans and matches the files of testdata/want. Its edit renames a field
// of Box.
func fixture(t *testing.T) workspacetest.Fixture {
	t.Helper()

	w, err := dokimi.Compose().Build()
	assert.NoError(t, err, "the composition builds")
	stores, err := w.Stores(os.Getenv)
	assert.NoError(t, err, "the frontends of the composition locate their stores")
	return workspacetest.Fixture{
		Tree:   os.DirFS(filepath.FromSlash(tree)),
		Stores: stores,
		Compose: func(root string) (*workspace.Builder, error) {
			b := dokimi.Base().
				Output(func() (output.Sink, error) { return output.NewDisk(root, dokimi.Brand) }).
				Ledger(func() (ledger.Ledger, error) { return ledger.OpenDir(root, dokimi.Brand) })
			return b, cli.Configure(b, root)
		},
		Plans: dokimi.Plans,
		Want:  workspacetest.Want(t, want),
		Edit:  edit,
	}
}

// edit renames the field Size of Box in the tree at root, the way a person
// edits a source file. It opens the file through an os.Root of root, and
// keeps the mode of the file. edit returns the error of a root or a file
// that does not open, read or write. It returns an error for a file without
// the field Size.
func edit(root string) error {
	r, err := os.OpenRoot(root)
	if err != nil {
		return fmt.Errorf("dokimi: %w", err)
	}
	defer r.Close()
	name := filepath.FromSlash(boxFile)
	b, err := r.ReadFile(name)
	if err != nil {
		return fmt.Errorf("dokimi: %w", err)
	}
	src := string(b)
	if !strings.Contains(src, sizeField) {
		return fmt.Errorf("dokimi: %s declares no field %s to rename", boxFile, sizeField)
	}
	if err := r.WriteFile(name, []byte(strings.Replace(src, sizeField, countField, 1)), sourceMode); err != nil {
		return fmt.Errorf("dokimi: %w", err)
	}
	return nil
}

// generated returns the files of the fixture f with the extension ext as
// the output of a toolchain run, with the tree as its sources. The tree has
// the go.mod of its module.
func generated(f workspacetest.Fixture, ext string) toolchain.Generated {
	files := map[string][]byte{}
	for p, b := range f.Want {
		if path.Ext(p) == ext {
			files[p] = b
		}
	}
	return toolchain.Generated{Files: files, Sources: f.Tree}
}

// generator returns the generator named id of the plan p. It stops the
// test where the plan has no such generator.
func generator(t *testing.T, p workspace.Plan, id plugin.ID) plugin.Generator {
	t.Helper()

	at := slices.IndexFunc(p.Generators, func(g plugin.Generator) bool { return g.Name() == id })
	assert.NotEqual(t, at, -1, "the plan "+p.Name+" runs "+string(id))
	return p.Generators[at]
}
