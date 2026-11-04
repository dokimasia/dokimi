//go:build tools

// Package tools keeps the requirement of go.dokimi.dev/assert, which the tests that the specs
// generate in this module import, so that go mod tidy keeps it.
package tools

import _ "go.dokimi.dev/assert"
