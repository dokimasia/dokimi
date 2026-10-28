// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Command cmd renders the spec in the working directory, which go generate
// sets to the directory of the package spec.
package main

import (
	"fmt"
	"os"

	"go.dokimi.dev/dokimi/spec/internal/render"
)

func main() {
	if err := render.Render("."); err != nil {
		fmt.Fprintln(os.Stderr, "render:", err)
		os.Exit(1)
	}
}
