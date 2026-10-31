// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/files"
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

// The source of the package store of the shop, its file and the file of the
// lock of the plan go. Store has a check of zero-on-error that the samples of
// Item need, and Index has a check of read-after-write with the partner Add.
const (
	storeFile   = "store/store.go"
	lockFile    = "dokimi.go.lock"
	storeSource = `// Package store keeps the items of the shop under two suites.
package store

import (
	"context"
	"fmt"
)

// Item is an item of the shop. eidos derives no value of its field Owner, so
// the directive sample names the functions that a check writes it with.
//
//dokimi:sample value=NewItem() alternate=OtherItem()
type Item struct {
	// Owner is the owner of the item.
	Owner fmt.Stringer
}

// NewItem returns the item that the checks write.
func NewItem() Item { return Item{Owner: owner("shop")} }

// OtherItem returns the second item that the checks write.
func OtherItem() Item { return Item{Owner: owner("supplier")} }

// owner is a fmt.Stringer with its name.
type owner string

// String returns the name of the owner.
func (o owner) String() string { return string(o) }

// Store keeps the items by key.
//
//dokimi:suite
type Store interface {
	// Get returns the item under the key.
	Get(ctx context.Context, key string) (Item, error)
}

// Index keeps the values by key.
//
//dokimi:suite
type Index interface {
	// Find returns the value that Add stored under the key.
	//
	//dokimi:mixin read-after-write write=Add
	Find(ctx context.Context, key string) (string, error)

	// Add stores the value under the key.
	Add(ctx context.Context, key, value string) error
}
`
)

// shop returns the tree of a module with the package store of source.
func shop(source string) files.Tree {
	return files.Tree{
		"go.mod":       files.Text("module example.com/shop\n\ngo 1.27.2\n"),
		".dokimi.yaml": files.Text("version: 1\n"),
		storeFile:      files.Text(source),
	}
}

// dokimiRun runs the binary bin with args and --no-color in the directory dir,
// and returns its output and its exit status.
func dokimiRun(t *testing.T, bin, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), bin, append(args, "--no-color")...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		return string(out), exit.ExitCode()
	}
	assert.NoError(t, err, "the binary runs")
	return string(out), 0
}

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

		t.Run("writes the brand of the binary with its module for version", func(t *testing.T) {
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

		t.Run("writes the lock of the plan go at the root of the workspace", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, shop(storeSource))

			out, status := dokimiRun(t, bin, dir, "run")

			assert.Equal(t, status, 0, "the run passes: "+out)
			files.HasContent(t, filepath.Join(dir, lockFile), "# dokimi lock 1\n"+
				"# Each line is one check, ceiling, refusal or unclaimed classification. Commit this file.\n"+
				"example.com/shop/store\tindex.add\tshape:composite-writer\tsilent\t-\n"+
				"example.com/shop/store\tindex.add/cancel\tcancel\tcheck\t-\n"+
				"example.com/shop/store\tindex.add/smoke\tsmoke\tcheck\t-\n"+
				"example.com/shop/store\tindex.find/cancel\tcancel\tcheck\t-\n"+
				"example.com/shop/store\tindex.find/read-after-write\tread-after-write\tcheck\twrite=index.add\n"+
				"example.com/shop/store\tindex.find/smoke\tsmoke\tcheck\t-\n"+
				"example.com/shop/store\tindex.find/zero-on-error\tzero-on-error\tcheck\t-\n"+
				"example.com/shop/store\tstore.get/cancel\tcancel\tcheck\t-\n"+
				"example.com/shop/store\tstore.get/smoke\tsmoke\tcheck\t-\n"+
				"example.com/shop/store\tstore.get/zero-on-error\tzero-on-error\tcheck\t-\n",
				"the lock has a line for each check of both suites")
		})

		losses := []struct {
			name     string
			giveFrom string
			giveTo   string
			wantLost []string
		}{
			{
				name:     "fails --check with LostCheck for each check of a suite without its directive",
				giveFrom: "//\n//dokimi:suite\ntype Index", giveTo: "type Index",
				wantLost: []string{"index.find/read-after-write (read-after-write)", "index.add/smoke (smoke)"},
			},
			{
				name:     "fails --check with LostCheck for a check that becomes a refusal",
				giveFrom: "//dokimi:sample value=NewItem() alternate=OtherItem()\n", giveTo: "",
				wantLost: []string{"store.get/zero-on-error (zero-on-error)"},
			},
			{
				name:     "fails --check with LostCheck for a check whose partner changes",
				giveFrom: "Add", giveTo: "Insert",
				wantLost: []string{"index.find/read-after-write (read-after-write)"},
			},
		}
		for _, tt := range losses {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				dir := files.Workspace(t, shop(storeSource))
				out, status := dokimiRun(t, bin, dir, "run")
				assert.Equal(t, status, 0, "the first run writes the lock: "+out)
				edited := strings.ReplaceAll(storeSource, tt.giveFrom, tt.giveTo)
				files.Write(t, dir, files.Tree{storeFile: files.Text(edited)})

				out, status = dokimiRun(t, bin, dir, "run", "--check")

				assert.Equal(t, status, 1, "the check fails")
				for _, lost := range tt.wantLost {
					expect.Contains(t, out, lost+" of example.com/shop/store is no longer generated as the lock states",
						"the check names the lost check")
				}
			})
		}

		t.Run("fails --check with the lock out of date where no suite of the plan remains", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, shop(storeSource))
			out, status := dokimiRun(t, bin, dir, "run")
			assert.Equal(t, status, 0, "the first run writes the lock: "+out)
			unsuited := strings.ReplaceAll(storeSource, "//\n//dokimi:suite\n", "")
			files.Write(t, dir, files.Tree{storeFile: files.Text(unsuited)})

			out, status = dokimiRun(t, bin, dir, "run", "--check")

			assert.Equal(t, status, 1, "the check fails")
			expect.Contains(t, out, lockFile+" is out of date", "the check reports the lock that the run would delete")
		})

		t.Run("passes --check once the run writes the regenerated lock", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, shop(storeSource))
			out, status := dokimiRun(t, bin, dir, "run")
			assert.Equal(t, status, 0, "the first run writes the lock: "+out)
			files.Write(t, dir, files.Tree{storeFile: files.Text(strings.Replace(storeSource,
				"//\n//dokimi:suite\ntype Index", "type Index", 1))})
			out, status = dokimiRun(t, bin, dir, "run")
			assert.Equal(t, status, 0, "the run writes the regenerated lock: "+out)

			out, status = dokimiRun(t, bin, dir, "run", "--check")

			assert.Equal(t, status, 0, "the check passes with the regenerated lock: "+out)
		})

		t.Run("writes the same bytes of the lock in every run of one source", func(t *testing.T) {
			t.Parallel()
			first, second := files.Workspace(t, shop(storeSource)), files.Workspace(t, shop(storeSource))
			_, firstStatus := dokimiRun(t, bin, first, "run")
			_, secondStatus := dokimiRun(t, bin, second, "run")
			assert.Equal(t, []int{firstStatus, secondStatus}, []int{0, 0}, "both runs pass")

			files.HasContent(t, filepath.Join(second, lockFile), files.Read(t, filepath.Join(first, lockFile)),
				"two runs of one source write one lock")
		})
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
