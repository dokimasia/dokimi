// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/internal/gentest"
)

// Project returns rendered files as the output of a run over a tree, with
// the paths of the files in the tree.
func TestProject(t *testing.T) {
	t.Parallel()

	t.Run("Project", func(t *testing.T) {
		t.Parallel()

		files := fstest.MapFS{
			"example.com/shop/catalog/address_recorder.go": {Data: []byte("package catalog\n")},
			"other/file.go": {Data: []byte("package other\n")},
		}
		got := gentest.Project(files, goTree, "example.com/shop")

		t.Run("removes the module path from the path of each file of the module", func(t *testing.T) {
			t.Parallel()

			assert.Contains(t, got.Files, "catalog/address_recorder.go", "the file is in the directory of its package")
			expect.Equal(t, string(got.Files["catalog/address_recorder.go"]), "package catalog\n",
				"the file keeps its content")
		})

		t.Run("keeps the path of a file outside the module", func(t *testing.T) {
			t.Parallel()

			expect.Contains(t, got.Files, "other/file.go", "the file keeps its path")
		})

		t.Run("returns the tree as the sources", func(t *testing.T) {
			t.Parallel()

			_, err := fs.Stat(got.Sources, "go.mod")
			expect.NoError(t, err, "the sources are the tree with its module file")
		})
	})
}
