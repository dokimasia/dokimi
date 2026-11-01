// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"

	"go.dokimi.dev/dokimi/internal/corpus"
	"go.dokimi.dev/dokimi/internal/report"
	"go.dokimi.dev/dokimi/spec"
)

// header is the comment lines that open a file of expected outcomes.
const header = "# Copyright Dokimasia B.V. 2026\n# SPDX-License-Identifier: Apache-2.0\n\n"

// priority is the file of the expected outcomes of the fixture enum/priority,
// as WriteExpectation writes it.
const priority = header + `fixture: enum/priority
level: 3
subjects: []
controls: []
declarations: {}
`

// ttlFile is the file of [expectation], as WriteExpectation writes it.
const ttlFile = header + `fixture: iface/mixin/ttl
level: 2
subjects:
  - type: iface/mixin/ttl/ttltest.InMemory
controls:
  - iface/mixin/ttl/deadline.Store
declarations:
  mixed:
    checks:
      mixed.read/cancel:
        claim: cancel
        law: cancelled-context-is-reported
        reason: ignored-cancel
      mixed.read/expires:
        claim: expires
        partners:
          write: mixed.put
        law: entry-lapses-on-the-clock
        reason: not-expired
    authored:
      - mixed.read/keeps-the-item
    ceilings:
      - mixed.read/latency-budget
    refusals:
      mixed.scan/zero-on-error: input-samples
    unclaimed:
      - classification: shape:writer
        callable: mixed.put
        stance: silent
`

func TestCorpus(t *testing.T) {
	t.Parallel()

	t.Run("ReadExpectations", func(t *testing.T) {
		t.Parallel()

		t.Run("reads the expectation of each file in the order of the files", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				"iface/mixin/ttl.yaml":        files.Text(ttlFile),
				"enum/priority.yaml":          files.Text(priority),
				"enum/archive.yaml/README.md": files.Text("# A directory whose name ends in .yaml\n"),
				"README.md":                   files.Text("# The expected outcomes\n"),
			})

			got, err := corpus.ReadExpectations(os.DirFS(dir))

			assert.NoError(t, err, "every file decodes")
			assert.Equal(t, got, []corpus.Expectation{
				{
					Fixture:      "enum/priority",
					Declarations: map[string]corpus.Declaration{},
					Subjects:     []corpus.Subject{},
					Controls:     []string{},
					Level:        3,
				},
				expectation(),
			}, "the expectations are in the order of their files")
		})

		tests := []struct {
			name     string
			giveFile string
			want     string
		}{
			{
				name:     "refuses a file that does not decode",
				giveFile: "fixture: [enum/priority\n",
				want:     "enum/priority.yaml",
			},
			{
				name:     "refuses a field that an expectation does not have",
				giveFile: "fixture: enum/priority\nlevels: 3\n",
				want:     "field levels not found",
			},
			{
				name:     "refuses a fixture that is not the path of its file",
				giveFile: "fixture: enum/status\nlevel: 3\n",
				want:     `enum/priority.yaml has the fixture "enum/status", and its path is enum/priority`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				dir := files.Workspace(t, files.Tree{"enum/priority.yaml": files.Text(tt.giveFile)})

				_, err := corpus.ReadExpectations(os.DirFS(dir))

				assert.ErrorIs(t, err, corpus.ErrFile, "the error is the error of an invalid file")
				assert.Contains(t, err.Error(), tt.want, "the error names the file and its fault")
			})
		}

		t.Run("returns the error of a directory that does not exist", func(t *testing.T) {
			t.Parallel()

			_, err := corpus.ReadExpectations(os.DirFS(filepath.Join(t.TempDir(), "expect")))

			assert.ErrorIs(t, err, fs.ErrNotExist, "the directory of the expected outcomes does not exist")
		})
	})

	t.Run("ReadOverlay", func(t *testing.T) {
		t.Parallel()

		t.Run("reads the overlay of a language", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{"overlay.yaml": files.Text(header + `language: typescript
level: 1
diverge:
  - fixture: iface/mixin/ttl
    key: mixed.read/zero-on-error
    expect: check
    actual: refused
    gate: returns-beside-error
    why: A call that fails throws and returns no value
`)})

			got, err := corpus.ReadOverlay(os.DirFS(dir), "overlay.yaml")

			assert.NoError(t, err, "the overlay decodes")
			assert.Equal(t, got, corpus.Overlay{Language: "typescript", Level: 1, Diverge: []corpus.Entry{{
				Fixture: ttl, Key: "mixed.read/zero-on-error", Expect: "check", Actual: "refused",
				Gate: spec.ReturnsBesideError, Why: "A call that fails throws and returns no value",
			}}}, "the overlay has its level and its entries")
		})

		t.Run("refuses a field that an overlay does not have", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{"overlay.yaml": files.Text("language: go\nlevels: 4\n")})

			_, err := corpus.ReadOverlay(os.DirFS(dir), "overlay.yaml")

			assert.ErrorIs(t, err, corpus.ErrFile, "the error is the error of an invalid file")
			assert.Contains(t, err.Error(), "overlay.yaml", "the error names the file")
		})

		t.Run("returns the error of a file that does not exist", func(t *testing.T) {
			t.Parallel()

			_, err := corpus.ReadOverlay(os.DirFS(t.TempDir()), "overlay.yaml")

			assert.ErrorIs(t, err, fs.ErrNotExist, "the overlay does not exist")
		})
	})

	t.Run("WriteExpectation", func(t *testing.T) {
		t.Parallel()

		t.Run("writes the expectation after the comment lines of the file that it replaces", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				"iface/mixin/ttl.yaml": files.Text(header + "fixture: iface/mixin/ttl\n"),
			})
			e := expectation()

			err := corpus.WriteExpectation(dir, &e)

			assert.NoError(t, err, "the file is written")
			files.HasContent(t, filepath.Join(dir, "iface", "mixin", "ttl.yaml"), ttlFile,
				"the file has the comment lines and then the expectation")
		})

		t.Run("returns the error of a fixture without a file", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{"iface/mixin": files.Dir()})
			e := expectation()

			err := corpus.WriteExpectation(dir, &e)

			assert.ErrorIs(t, err, fs.ErrNotExist, "the fixture has no file to replace")
			files.Absent(t, filepath.Join(dir, "iface", "mixin", "ttl.yaml"), "the write creates no file")
		})

		t.Run("returns the error of a file that it cannot write", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				"iface/mixin/ttl.yaml": files.Text(header).WithMode(0o444),
			})
			e := expectation()

			err := corpus.WriteExpectation(dir, &e)

			assert.ErrorIs(t, err, fs.ErrPermission, "the file is read-only")
		})
	})

	t.Run("Fixtures", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each directory of the groups of the expected fixtures", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				"enum/status/status.go":     files.Text("package status\n"),
				"enum/priority/priority.go": files.Text("package priority\n"),
				"enum/README.md":            files.Text("# The enumerations\n"),
				"iface/mixin/ttl/iface.go":  files.Text("package ttl\n"),
				"struct/plain/plain.go":     files.Text("package plain\n"),
			})
			expected := []corpus.Expectation{{Fixture: "enum/priority"}, {Fixture: ttl}, {Fixture: "errors/store"}}

			got, err := corpus.Fixtures(os.DirFS(dir), expected)

			assert.NoError(t, err, "the groups list")
			assert.Equal(t, got, []string{"enum/priority", "enum/status", ttl},
				"the fixtures are the directories of the groups enum and iface/mixin")
		})

		t.Run("returns the error of a group that is not a directory", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{"enum": files.Text("a file\n")})

			_, err := corpus.Fixtures(os.DirFS(dir), []corpus.Expectation{{Fixture: "enum/priority"}})

			failed := assert.ErrorAs[*fs.PathError](t, err, "the group enum is a file")
			assert.Equal(t, *failed, fs.PathError{Op: "readdir", Path: "enum", Err: syscall.ENOTDIR},
				"the error is the readdir error of the group on every system")
		})

		t.Run("returns the error of a group that is not a valid path", func(t *testing.T) {
			t.Parallel()

			_, err := corpus.Fixtures(os.DirFS(t.TempDir()), []corpus.Expectation{{Fixture: "../priority"}})

			assert.ErrorIs(t, err, fs.ErrInvalid, "the group .. is not a valid path of the corpus")
		})
	})

	t.Run("Outcomes", func(t *testing.T) {
		t.Parallel()

		t.Run("groups the lines and the records by fixture without the module path", func(t *testing.T) {
			t.Parallel()
			const module = "example.com/corpus"
			lock := []string{
				"# dokimi lock 1",
				"",
				module + "/iface/mixin/ttl\tmixed.read/cancel\tcancel\tcheck\t-",
				module + "/iface/mixin/ttl/clock\tclock.now/smoke\tsmoke\tcheck\t-",
				module + "/iface/mixin/new\tnew.get/smoke\tsmoke\tcheck\t-",
				module + "/iface/mixin/new\tnew.put/smoke\tsmoke\tcheck\t-",
			}
			r := record(spec.OnSubject, module+"/"+subject)
			r.Suite.Package = module + "/" + ttl

			got := corpus.Outcomes(module, []string{"enum/priority", ttl}, lock, []report.Record{r})

			want := record(spec.OnSubject, subject)
			assert.Equal(t, got, []corpus.Outcome{
				{Fixture: "enum/priority"},
				{
					Fixture: ttl,
					Lock: []string{
						"iface/mixin/ttl\tmixed.read/cancel\tcancel\tcheck\t-",
						"iface/mixin/ttl/clock\tclock.now/smoke\tsmoke\tcheck\t-",
					},
					Records: []report.Record{want},
				},
				{
					Fixture: "iface/mixin/new",
					Lock: []string{
						"iface/mixin/new\tnew.get/smoke\tsmoke\tcheck\t-",
						"iface/mixin/new\tnew.put/smoke\tsmoke\tcheck\t-",
					},
				},
			}, "each fixture has its lines and its records, and a package outside the fixtures has its own")
		})
	})
}
