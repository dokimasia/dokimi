// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"go.dokimi.dev/eidos/cli"

	"go.dokimi.dev/dokimi"
)

// help is the help text of the app.
const help = "Generates the tests and the test helpers of the declarations of a workspace."

// main runs the app of the binary over the composition of dokimi, and exits
// with the status of the run.
func main() {
	cli.Main(cli.New(help, dokimi.Compose), dokimi.Compose)
}
