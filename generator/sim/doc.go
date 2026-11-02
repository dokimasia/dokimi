// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sim generates the deterministic simulation of a system of sans-IO
// machines from the directives on their declarations, and the tests that run
// it through the module go.dokimi.dev/sim.
//
// The directive sim puts a machine interface under the generator:
//
//	//dokimi:sim start=Started stop=Stopped cancel=Cancelled tolerate=crash:1,partition:1 strategy=pct:3
//	type Machine interface {
//		Step(ev Event, out []Effect) []Effect
//		Err() error
//	}
//
// The events of the machine are the types that implement its event
// interface, and its effects the types that implement its effect interface.
// A type under the directive token is the type of the field that links an
// effect to the event of its result, and a type under the directive token
// reply the type of the field that links an inbound call to its reply. The
// directive reply marks the effect of a reply, and the directive perform
// names the function that performs an effect. The directives system, action,
// spec, invariant, converges, fingerprint and scenario mark the functions of
// the consumer that build the system and check it.
//
// The plugin is dual. Its annotate half stamps the facts of each directive on
// a type or a function. Its generate half records the model of the
// simulation of each machine interface: the anatomy of the machine, the
// performer of its effects, the simulation with the checks of the claims of
// the spec whose generator is sim, and the tests TestSimulation,
// TestSimulationFalsify and FuzzSimulation.
//
// # Binding
//
// A function under a directive binds to the machine interface whose event and
// effect interfaces are the type arguments of the world in its signature. The
// simulation is in the package of the directive out of the machine interface,
// and each bound function is in that package, so the simulation can call a
// function that the package does not export.
//
// # Diagnostics
//
// sim reports [BadMachine] at an interface without the step method and the
// error method of a sans-IO machine, [BadLifecycle] at a param start, stop or
// cancel that names no event of its form, and [BadParam] at a param tolerate,
// require or strategy that names no fault or no strategy. It reports
// [NoPerform] at an effect type without the directive perform, [NoReply] at a
// machine with an action call and without a reply, [BadReply] at a reply that
// lacks a field of its directive, [BadToken] at a token type that is not an
// unsigned integer type, and [BadFunction] at a function that does not have
// the signature of its directive. sim records nothing for a machine interface
// that it reports.
//
// # Dependency position
//
// The package imports the packages internal/lang of the module and spec of
// the module, the root package of the Go satellite in golang.go, the package
// naming of eidos-lang, the root, diag, directive, emit, meta, node, plugin,
// rules and symbol packages of the SDK facade of eidos, and the Go stdlib.
package sim
