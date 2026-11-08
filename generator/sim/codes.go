// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that sim reports under.
const CodePrefix diag.Prefix = "DOKIMISIM"

// BadMachine reports an interface under the directive sim whose step method
// does not have the form Step(ev E, out []F) []F with interface types E and
// F, or whose error method is not Err() error. A world cannot run such a
// machine, so the finding is an Error at the interface, and sim writes nothing
// for it.
var BadMachine = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "an interface under the directive sim does not have the step method and the error method of a sans-IO machine",
})

// BadLifecycle reports a param start, stop or cancel of the directive sim
// that names no event type of its form: an event with one field of the reply
// token for start and cancel, and an event for stop. The world builds these
// events, so the finding is an Error at the interface, and sim writes nothing
// for it.
var BadLifecycle = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "a param start, stop or cancel of the directive sim names no event type of the form of its event",
})

// BadParam reports a param tolerate, require or strategy of the directive sim
// or of the directive system of a system of components with a value that is
// not a fault, a mode or a strategy, or with a count or a share out of its
// range. It also reports a param of a simulation on the directive system of a
// machine interface, because the directive sim of the interface has the
// params of its simulation. The finding is an Error at the interface or at
// the function, and sim writes nothing for its simulation.
var BadParam = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 3,
	Meaning: "a param tolerate, require or strategy of the directive sim or system has a value that is not a fault, " +
		"a mode or a strategy, or is on the directive system of a machine interface",
})

// NoPerform reports an effect type of a machine interface under the directive
// sim without the directive perform, other than the reply. The world cannot
// perform the effect, so the finding is an Error at the effect type, and sim
// writes nothing for the interface.
var NoPerform = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  4,
	Meaning: "an effect type of a machine interface under the directive sim has no directive perform",
})

// NoReply reports a machine interface under the directive sim with an action
// call and without an effect type under the directive reply. The run cannot
// complete a call in the history, so the finding is an Error at the
// interface, and sim writes nothing for it.
var NoReply = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  5,
	Meaning: "a machine interface under the directive sim has an action call and no effect type under the directive reply",
})

// BadReply reports an effect type under the directive reply that lacks a
// field that the directive names, or whose type parameter has a constraint
// that lists no types. The anatomy reads a reply through one case for each
// type of the constraint, so the finding is an Error at the effect type, and
// sim writes nothing for its interface.
var BadReply = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  6,
	Meaning: "an effect type under the directive reply lacks a field that the directive names, or has a type parameter whose constraint lists no types",
})

// BadToken reports a type under the directive token that is not an unsigned
// integer type. The world numbers the tokens of a run as uint64, so the
// finding is an Error at the type, and sim writes nothing for the interfaces
// whose events or effects have it.
var BadToken = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  7,
	Meaning: "a type under the directive token is not an unsigned integer type",
})

// BadFunction reports a function under a directive of sim that does not have
// the signature of its directive, that is a second function under a directive
// that a simulation takes once, or that is not in the package of the
// simulation of its machine interface. The simulation cannot call the
// function, so the finding is an Error at the function, and sim writes
// nothing for its interface or its system.
var BadFunction = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 8,
	Meaning: "a function under a directive of sim does not have the signature of its directive, is a second " +
		"function under a directive that a simulation takes once, or is not in the package of the simulation of " +
		"its machine interface",
})

// BadPerform reports a directive perform that names both a function and a
// call or neither, a call without the event of its result, a dependency
// without a function or without the event of the failure, or a function with
// a result and without a dependency. It also reports a dependency that is no
// interface, the param ends beside a call, a dependency, a result, a timeout
// or an end, and the param end without a dependency. It reports a call, a
// dependency, an end event or an ended wait whose parameters, results, events
// and tokens the fields of the effect and of the events do not bind. The world
// cannot perform the effect, so the finding is an Error at the effect type,
// and sim writes nothing for its interface.
var BadPerform = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 9,
	Meaning: "a directive perform names both a function and a call or neither, a call without a result, a " +
		"dependency without a function or a result, a function with a result and without a dependency, a " +
		"dependency that is no interface, the end of a wait beside a call or a result, an end event without a " +
		"dependency, or a call, a dependency or an end that the effect and its events do not bind",
})
