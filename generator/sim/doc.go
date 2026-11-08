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
// A function under the directive system whose world is
// *sim.World[sim.None, sim.None] is a system of components, which has no
// machine interface:
//
//	//dokimi:system tolerate=[crash:1, partition:1, blocking] strategy=uniform
//	func cluster(w *sim.World[sim.None, sim.None]) {
//		sim.SpawnComponent(w, "cell", cells.Open)
//	}
//
// Its directive system states the params tolerate, require, strategy and
// isolation of its simulation. The values tasks and blocking of tolerate are
// its modes. tasks runs the hosts of its components under a task scheduler of
// each case, and blocking lets a world in native mode strike the blocking
// faults. The simulation runs each case in conform.Bubble, and has no rule of
// sim, because each rule of sim judges the steps of machines.
//
// The events of the machine are the types that implement its event
// interface, and its effects the types that implement its effect interface.
// A type under the directive token is the type of the field that links an
// effect to the event of its result, and a type under the directive token
// reply the type of the field that links an inbound call to its reply. The
// directive reply marks the effect of a reply. The directive perform names
// the function that performs an effect, or the method of a dependency that
// performs it. A function that calls a dependency names the interface of the
// dependency in the param dependency and the event of the failure of the
// effect in the param result, so the faults of the dependency strike the
// effect:
//
//	//dokimi:perform func=write dependency=Store result=Written timeout=Timeout
//	type Write struct{ ... }
//
// An outage of the dependency ends the wait of such an effect. A wait that
// has delivered an event ends with the event of the param end, and any other
// wait ends with the event of the param result:
//
//	//dokimi:perform func=campaign dependency=Store result=Campaigned end=TermEnded
//	type Campaign struct{ ... }
//
// The param ends marks an effect that ends the waits of the effects that it
// lists, such as an unwatch. The anatomy returns no token for such an effect,
// and the world cancels the wait of its token. Its function is optional:
//
//	//dokimi:perform ends=[Watch]
//	type Unwatch struct{ ... }
//
// The events and the effects of a machine leave out the types of the test
// files, because the anatomy refers to each event and each effect. A file
// outside the tests of a package cannot refer to a type of a test file.
//
// The directives system, action, spec, invariant, converges, fingerprint and
// scenario mark the functions of the consumer that build the system and
// check it. A system and the defect of an invariant that take inputs from the
// case have the form of an action, func(c *prop.Case) func(w *sim.World[E, F]).
//
// The plugin is dual. Its annotate half stamps the facts of each directive on
// a type or a function. Its generate half records the model of the
// simulation of each machine interface and of each system of components: the
// anatomy of the machine, the performer of its effects, the simulation with
// the checks of the claims of the spec whose generator is sim, and the tests
// TestSimulation, TestSimulationFalsify and FuzzSimulation.
//
// # Binding
//
// A function under a directive binds to the machine interface whose event and
// effect interfaces are the type arguments of the world in its signature. The
// simulation is in the package of the directive out of the machine interface,
// and each bound function is in that package, so the simulation can call a
// function that the package does not export. The simulation of a system of
// components is in the package of its function under the directive system,
// and binds the functions of that package whose world is
// *sim.World[sim.None, sim.None].
//
// # Diagnostics
//
// sim reports [BadMachine] at an interface without the step method and the
// error method of a sans-IO machine, [BadLifecycle] at a param start, stop or
// cancel that names no event of its form, and [BadParam] at a value of a
// param tolerate, require or strategy that is not a fault, a mode or a
// strategy. It also reports [BadParam] at a param of a simulation on the
// directive system of a machine interface. It reports [NoPerform] at an effect type
// without the directive perform, [BadPerform] at a directive perform whose
// params or fields do not bind, [NoReply] at a machine with an action call
// and without a reply, [BadReply] at a reply that lacks a field of its
// directive, [BadToken] at a token type that is not an unsigned integer
// type, and [BadFunction] at a function that does not have the signature of
// its directive. sim records nothing for a machine interface or a system of
// components that it reports.
//
// # Dependency position
//
// The package imports the packages internal/lang of the module and spec of
// the module, the root package of the Go satellite in golang.go, the package
// naming of eidos-lang, the root, diag, directive, emit, meta, node, plugin,
// rules and symbol packages of the SDK facade of eidos, and the Go stdlib.
package sim
