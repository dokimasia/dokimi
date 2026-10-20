// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"os"
	"path/filepath"
	"strings"
	"testing/fstest"

	"go.dokimi.dev/eidos/sdk/toolchain"
)

// Project returns the files that [Rendered] returns as the output of a run
// over the tree dir, for a check of the toolchain of the language. The
// render routes each file of a Go package under the import path of its
// package, so Project removes the module path module and its slash from the
// start of each path. The project has the tree dir as its sources, with the
// module file of the tree.
func Project(files fstest.MapFS, dir, module string) toolchain.Generated {
	out := make(map[string][]byte, len(files))
	for path, f := range files {
		out[strings.TrimPrefix(path, module+"/")] = f.Data
	}
	return toolchain.Generated{Files: out, Sources: os.DirFS(filepath.FromSlash(dir))}
}
