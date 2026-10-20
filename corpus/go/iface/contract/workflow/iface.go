// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package workflow is the contract-axis fixture for the workflow contract:
// a callable driving a state machine.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package workflow

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// The states this fixture's machine moves through, as plain strings: a
// named type here would be an enum node, and neither `default` nor
// `sample` applies to one — which would leave State's zero-on-error claim
// with no value it can tell the zero apart from.
//
// The empty state is where a key sits before anything ran. It is a real
// state the graph names, not a stand-in for absence.
const (
	StateNew   = ""
	StateDraft = "Draft"
	StateLive  = "Live"
)

// Transitions is the graph the directive names, and the moves a check
// holds Run to.
//
// A declaration rather than a string the directive encodes: the graph is
// a set of pairs, and a parameter names one thing. What can be named as
// one thing is the var that holds it, which leaves the shape here — a
// map from a state to the states reachable from it — the fixture's own
// rather than a notation every consumer would have to write in.
var Transitions = map[string][]string{
	StateNew:   {StateDraft},
	StateDraft: {StateLive},
}

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=workflowtest/
type Contract interface {
	// Run is the workflow contract's fn role, and hosts the directive
	// that names the graph and the observation.
	//dokimi:contract workflow role=fn transitions=Transitions observe=State
	Run(ctx context.Context, key string) error

	// State reports the key's position in the machine — the observation
	// observe= names, which a check compares across a Run.
	State(ctx context.Context, key string) (string, error)
}
