// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package weave contains the names through which the suite generator asks the
// stub generator for the double of a declaration under the directive suite.
//
// The annotate half of suite stamps a fact under [KeyDouble] on each interface
// and struct under the directive suite, and provides [Suite] in the annotator
// role. A rule of stub has a gate on the fact, and stub runs after [Suite]. In
// a composition without suite, the gate does not select a subject. suite
// requires the capability of stub, so the generator half of suite still runs
// after stub.
//
// # Dependency position
//
// The package imports the meta and plugin packages of the SDK facade of
// eidos.
package weave

import (
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
)

// Suite is the capability that the annotate half of suite provides, and that
// stub runs after.
const Suite plugin.Capability = "suite"

// KeyDouble is the key of the fact that suite stamps on each interface and
// struct under the directive suite. The fact asks stub for the double of the
// declaration. The key is in the namespace of the keys of suite.
const KeyDouble meta.KeyName = "suite.double"
