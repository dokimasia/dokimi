// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/sim"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "BadMachine", Give: sim.BadMachine, WantCode: "DOKIMISIM-0001",
			WantMeaning: "an interface under the directive sim does not have the step method and the error method " +
				"of a sans-IO machine",
		},
		{
			Name: "BadLifecycle", Give: sim.BadLifecycle, WantCode: "DOKIMISIM-0002",
			WantMeaning: "a param start, stop or cancel of the directive sim names no event type of the form of its " +
				"event",
		},
		{
			Name: "BadParam", Give: sim.BadParam, WantCode: "DOKIMISIM-0003",
			WantMeaning: "a param tolerate, require, strategy, mean-steps or max-steps of the directive sim or " +
				"system has a value that is not a fault, a mode, a strategy or a count in its range, such a param " +
				"or the param unknown is on the directive system of a machine interface, or the directive action " +
				"has a word after final",
		},
		{
			Name: "NoPerform", Give: sim.NoPerform, WantCode: "DOKIMISIM-0004",
			WantMeaning: "an effect type of a machine interface under the directive sim has no directive perform",
		},
		{
			Name: "NoReply", Give: sim.NoReply, WantCode: "DOKIMISIM-0005",
			WantMeaning: "a machine interface under the directive sim has an action call and no effect type under " +
				"the directive reply",
		},
		{
			Name: "BadReply", Give: sim.BadReply, WantCode: "DOKIMISIM-0006",
			WantMeaning: "an effect type under the directive reply lacks a field that the directive names, or has a " +
				"type parameter whose constraint lists no types",
		},
		{
			Name: "BadToken", Give: sim.BadToken, WantCode: "DOKIMISIM-0007",
			WantMeaning: "a type under the directive token is not an unsigned integer type",
		},
		{
			Name: "BadFunction", Give: sim.BadFunction, WantCode: "DOKIMISIM-0008",
			WantMeaning: "a function under a directive of sim does not have the signature of its directive, is a " +
				"second function under a directive that a simulation takes once, or is not in the package of the " +
				"simulation of its machine interface",
		},
		{
			Name: "BadPerform", Give: sim.BadPerform, WantCode: "DOKIMISIM-0009",
			WantMeaning: "a directive perform names both a function and a call or neither, a call without a result, " +
				"a dependency without a function or a result, a function with a result and without a dependency, a " +
				"dependency that is no interface, the end of a wait beside a call or a result, an end event without " +
				"a dependency, or a call, a dependency or an end that the effect and its events do not bind",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
