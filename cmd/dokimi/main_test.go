// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/cli/acceptancetest"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"

	"go.dokimi.dev/dokimi"
)

// childVar is the variable of the environment that makes the test binary
// run main in place of its tests. The case of version runs the test binary
// again with it, so go test merges the coverage of main into the profile.
const childVar = "DOKIMI_TEST_MAIN"

// mainPackage is the import path of the main package, which the acceptance
// checks build into the binary that they run.
const mainPackage = "go.dokimi.dev/dokimi/cmd/dokimi"

// tree is the workspace tree of the acceptance checks, with the config file
// of the brand at its root.
const tree = "testdata/tree"

// The file that the Fail of the fixture writes: a struct under the directive
// builder with a default that is no value of its type, which builder reports
// as an Error.
const (
	shelfFile   = "catalog/shelf.go"
	shelfSource = "package catalog\n\n// Shelf has a default that is no value of its type.\n//\n" +
		"//dokimi:builder\ntype Shelf struct {\n\t// Slots is the number of slots.\n" +
		"\tSlots int `default:\"many\"`\n}\n"
	fileMode = 0o644
)

func TestMain(m *testing.M) {
	if os.Getenv(childVar) != "" {
		main()
	}
	os.Exit(m.Run())
}

// The binary writes its brand and its module for version, and meets the
// contract of the command line of eidos. Its run over the tree writes Go
// that compiles.
func TestDokimi(t *testing.T) {
	t.Parallel()

	t.Run("main", func(t *testing.T) {
		t.Parallel()

		t.Run("writes the brand and the module of the binary for version", func(t *testing.T) {
			t.Parallel()

			binary, err := os.Executable()
			assert.NoError(t, err, "the path of the test binary")
			cmd := exec.CommandContext(t.Context(), binary, "version")
			cmd.Env = append(os.Environ(), childVar+"=1")
			out, err := cmd.Output()
			assert.NoError(t, err, "dokimi version exits 0")
			assert.Contains(t, string(out), "brand dokimi\nmodule go.dokimi.dev/dokimi ",
				"dokimi version writes the brand and the module of the binary")
		})

		f := fixture()
		bin := acceptancetest.Build(t, mainPackage)
		checks := []struct {
			name  string
			check func(assert.TB, acceptancetest.Fixture, string, string)
		}{
			{name: "exits with the statuses of the contract", check: acceptancetest.AssertStatuses},
			{name: "runs each kernel command under its name", check: acceptancetest.AssertNames},
			{name: "finds the config file of the root", check: acceptancetest.AssertDiscovery},
			{name: "writes JSON lines alone under --format=json", check: acceptancetest.AssertJSON},
			{name: "rewrites no file on a second run", check: acceptancetest.AssertIdempotent},
			{name: "compiles the generated output", check: acceptancetest.AssertCompiles},
			{name: "runs each member of a list under a state directory of its own", check: acceptancetest.AssertLists},
			{name: "exits with status 1 while another process has the lock", check: acceptancetest.AssertLocked},
		}
		for _, c := range checks {
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				c.check(t, f, bin, t.TempDir())
			})
		}
	})
}

// fixture returns the fixture of the acceptance checks over the tree. Its
// Fail writes a struct with a default that is no value of its type, and its
// Compile compiles each package of the output and its tests.
func fixture() acceptancetest.Fixture {
	return acceptancetest.Fixture{
		Main:  mainPackage,
		Brand: dokimi.Brand,
		Tree:  os.DirFS(filepath.FromSlash(tree)),
		Fail: func(root string) error {
			return os.WriteFile(filepath.Join(root, filepath.FromSlash(shelfFile)), []byte(shelfSource), fileMode)
		},
		Compile: func(ctx context.Context, root string) error {
			return gotesting.New().TypeCheck(ctx, root)
		},
	}
}
