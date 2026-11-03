// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim_test

import (
	"path"
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/sim"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// otherRegistrant is the registrant that claims the namespace of sim before
// sim registers its keys.
const otherRegistrant = "other"

// declOf returns the declaration of the run r that the text qualified names
// by the last element of the path of its package and its name, such as
// crookedsim.again.
func declOf(t *testing.T, r run, qualified string) node.Declaration {
	t.Helper()

	pkg, name, _ := strings.Cut(qualified, ".")
	for p := range r.f.Graph.Packages() {
		if path.Base(p.ID.Package) != pkg {
			continue
		}
		for _, file := range p.Files {
			for _, d := range file.Decls {
				if decl, declares := d.(node.Declaration); declares && decl.Identity().Name == name {
					return decl
				}
			}
		}
	}
	t.Fatalf("the run has no declaration %s", qualified)
	return nil
}

// spelled returns the spelling of each type of fields with its field, as
// Write.Reply.
func spelled(fields []sim.Field) []string {
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		out = append(out, f.Type.Spelling+"."+f.Field)
	}
	return out
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("records the anatomy of a machine interface", func(t *testing.T) {
		t.Parallel()

		r := generate(t, registerPackage, registerSimPackage)
		expect.Empty(t, slices.Concat(slices.Collect(r.stamps.Sink.All()), slices.Collect(r.res.Sink.All())),
			"the register has no finding")
		m := r.models(t)["Machine"]
		expect.Equal(t, m.Names, sim.Names{
			Anatomy: "machineAnatomy", Perform: "performMachine", Simulation: "MachineSimulation",
		}, "the names of the declarations of the simulation join the name of the interface")
		expect.Equal(t, []string{m.Event.Spelling, m.Effect.Spelling}, []string{"Event", "Effect"},
			"the event and the effect interfaces are the types of the step method")
		expect.Equal(t, []string{m.Started.Type.Spelling + "." + m.Started.Field, m.Started.Token.Spelling},
			[]string{"Started.Reply", "ReplyToken"}, "the start is the event of the param start with its reply token")
		expect.Equal(t, m.Stopped.Spelling, "Stopped", "the stop is the event of the param stop")
		assert.NotNil(t, m.Cancelled, "the param cancel names the cancel event")
		expect.Equal(t, m.Cancelled.Type.Spelling+"."+m.Cancelled.Field, "Cancelled.Reply",
			"the cancel is the event of the param cancel with its reply token")
		expect.Equal(t, spelled(m.Calls), []string{"Started.Reply", "Watch.Reply", "Write.Reply"},
			"each event with a reply token but the cancel is an inbound call")
		expect.Equal(t, spelled(m.Tokens), []string{"Commit.Token"}, "each effect with a token starts a call")
		expect.Equal(t, m.Replies, []sim.Reply{{Type: m.Replies[0].Type, Value: "Value", Err: "Err", To: "To"}},
			"the reply has the fields of its directive")
		assert.Length(t, m.Unknown, 1, "the reply names one error of an unknown outcome")
		expect.Equal(t, m.Unknown[0].Spelling, "ErrStopped", "the error of an unknown outcome is the error of a stop")
		expect.Equal(t, m.Performs, []sim.Perform{{Type: m.Performs[0].Type, Func: "commit", Timeout: "Timeout"}},
			"each effect has the performer and the timeout of its directive")
		expect.True(t, m.Timeouts, "an effect has a timeout")
	})

	t.Run("records the functions of the consumer in the package of the simulation", func(t *testing.T) {
		t.Parallel()

		m := generate(t, registerPackage, registerSimPackage).models(t)["Machine"]
		expect.Equal(t, []string{m.System, m.Fingerprint}, []string{"system", "fingerprint"},
			"the simulation has the system and the fingerprint of the consumer")
		expect.Equal(t, *m.Spec, sim.Spec{Func: "spec"}, "the spec returns a history.Spec")
		expect.Equal(t, m.State.Spelling, "string", "the state of the spec is the type argument of its spec")
		expect.Equal(t, m.Actions, []sim.Action{
			{Func: "pause"},
			{Func: "settle", Call: true, Final: true},
			{Func: "watch", Call: true},
			{Func: "write", Call: true},
		}, "the actions are in the order of their identities")
		expect.Equal(t, m.Invariants, []sim.Invariant{{
			Statement: "The system has the node b", Func: "spawned", Defect: "lonely",
		}}, "an invariant has the first sentence of its doc comment and its defect")
		expect.Equal(t, m.Converges, []sim.Convergence{{
			Statement: "The log of the run has a step once the faults have ended", Func: "replied", Within: "bound",
		}}, "a convergence has the first sentence of its doc comment and its bound")
		expect.Equal(t, m.Scenarios, []sim.Scenario{{Func: "writeOnA"}}, "a scenario has the name of its function")
	})

	t.Run("records the params of the directive", func(t *testing.T) {
		t.Parallel()

		m := generate(t, registerPackage, registerSimPackage).models(t)["Machine"]
		expect.Equal(t, m.Tolerate, []sim.Bound{{Fault: "FaultCrash", N: 1}, {Fault: "FaultPartition", N: 1}},
			"the param tolerate bounds each fault that it lists")
		expect.Equal(t, m.Require, []sim.Share{{Fault: "FaultCrash", Share: "0.1"}},
			"the param require gives each fault that it lists its share")
		expect.Equal(t, []any{m.Strategy, m.Depth}, []any{"pct", 2}, "the param strategy is PCT of depth 2")
	})

	t.Run("records the params of a uniform simulation", func(t *testing.T) {
		t.Parallel()

		m := generate(t, plainPackage, plainSimPackage).models(t)["Plain"]
		expect.Equal(t, []string{m.Strategy, m.Isolation}, []string{"uniform", "serializable"},
			"the params strategy and isolation of the directive")
		expect.Equal(t, m.Tolerate, []sim.Bound{{Fault: "FaultPause", N: 1}}, "the param tolerate bounds the pause")
		expect.Equal(t, m.Require, []sim.Share{{Fault: "FaultPause", Share: "1"}},
			"a share of 1 requires the pause in every case")
	})

	t.Run("records the faults of the media of a process", func(t *testing.T) {
		t.Parallel()

		m := generate(t, disksPackage, disksSimPackage).models(t)["Disks"]
		expect.Equal(t, m.Tolerate, []sim.Bound{
			{Fault: "FaultTornWrite", N: 1}, {Fault: "FaultCorruption", N: 1}, {Fault: "FaultFullDisk", N: 1},
		}, "the param tolerate bounds each fault of the media")
		expect.Equal(t, m.Require, []sim.Share{{Fault: "FaultFullDisk", Share: "0.5"}},
			"the param require gives the full disk its share")
	})

	t.Run("records a spec of a subject", func(t *testing.T) {
		t.Parallel()

		m := generate(t, plainPackage, plainSimPackage).models(t)["Plain"]
		expect.Equal(t, *m.Spec, sim.Spec{Func: "subject", From: true}, "history.SpecFrom wraps the subject")
		expect.Nil(t, m.State, "the state of history.SpecFrom is the list of the applied operations")
	})

	t.Run("records a machine without a cancel event", func(t *testing.T) {
		t.Parallel()

		m := generate(t, plainPackage, plainSimPackage).models(t)["Plain"]
		expect.Nil(t, m.Cancelled, "the machine has no cancel event")
		expect.Equal(t, spelled(m.Calls), []string{"Started.Reply"}, "the start is the one inbound call")
	})

	t.Run("records an effect without a timeout", func(t *testing.T) {
		t.Parallel()

		m := generate(t, plainPackage, plainSimPackage).models(t)["Plain"]
		expect.Equal(t, m.Performs, []sim.Perform{{Type: m.Performs[0].Type, Func: "beta"}},
			"the effect has the performer of its directive")
		expect.False(t, m.Timeouts, "no effect has a timeout")
	})

	t.Run("records the call of a dependency that performs an effect", func(t *testing.T) {
		t.Parallel()

		r := generateWithFaults(t, remotePackage, remoteSimPackage, storePackage)
		expect.Empty(t, slices.Concat(slices.Collect(r.stamps.Sink.All()), slices.Collect(r.res.Sink.All())),
			"the remote register has no finding")
		m := r.models(t)["Machine"]
		assert.Length(t, m.Performs, 3, "each effect but the reply has a performer")
		for _, p := range m.Performs {
			assert.NotNil(t, p.Call, "a call of a dependency performs the effect "+p.Type.Spelling)
			expect.Empty(t, p.Func, "the effect "+p.Type.Spelling+" has no function that performs it")
			expect.Equal(t, []string{p.Call.Token, p.Call.EffectToken, p.Call.Err}, []string{"Token", "Token", "Err"},
				"the event takes the token of the effect and the error of the call")
			expect.True(t, p.Call.Fails, "each method of a dependency returns an error")
		}
		get, load, put := m.Performs[0].Call, m.Performs[1].Call, m.Performs[2].Call
		for _, call := range []*sim.Call{get, put} {
			expect.Equal(t, []string{call.Dependency.Spelling, call.Name}, []string{"Writer", "remote.Writer"},
				"the writer of the node performs the effect "+call.Method)
		}
		expect.Equal(t, []string{load.Dependency.Spelling, load.Name, load.Method, load.Event.Spelling},
			[]string{"Reader", "store.Reader", "Read", "Loaded"}, "the reader of another package performs a load")
		expect.Equal(t, load.Values, []string{"Value"},
			"the event of a load takes the value of the type that the reader declares")
		expect.Equal(t, []string{get.Method, get.Event.Spelling}, []string{"Get", "Fetched"},
			"a get performs the effect Get, and its event is Fetched")
		expect.Equal(t, get.Values, []string{"Value"}, "the event of a get takes the value of the call")
		expect.Equal(t, get.Args, []sim.Arg{{Context: true}, {Field: "Key"}},
			"the get takes the context of the world and the key of the effect")
		expect.Empty(t, get.Errors, "the get has no directive fault")
		expect.Equal(t, []string{put.Method, put.Event.Spelling, m.Performs[2].Timeout},
			[]string{"Put", "Stored", "Timeout"}, "a put performs the effect Put within its timeout")
		expect.Empty(t, put.Values, "the event of a put takes no value")
		expect.Equal(t, put.Args, []sim.Arg{{Context: true}, {Field: "Key"}, {Field: "Value"}},
			"the put takes the context of the world, and the key and the value of the effect")
		assert.Length(t, put.Errors, 1, "the directive fault of the put declares one error")
		expect.Equal(t, put.Errors[0].Spelling, "ErrFull", "the put fails with the error of its directive fault")
	})

	misfits := []struct {
		name     string
		giveDecl string
		wantMsg  string
	}{
		{
			name: "refuses a perform of a function and a call", giveDecl: "Twice",
			wantMsg: "the directive perform of Twice names both a function and a call. Name one of them",
		},
		{
			name: "refuses a perform of neither a function nor a call", giveDecl: "Neither",
			wantMsg: "the directive perform of Neither names neither a function nor a call. Write " +
				"func=<function>, or call=<callable> result=<event>",
		},
		{
			name: "refuses a perform of a call without a result", giveDecl: "Resultless",
			wantMsg: "the directive perform of Resultless names a call without a result. Write result=<event>",
		},
	}
	for _, tt := range misfits {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := generate(t, "example.com/shop/misfit")
			got := gentest.Finding(t, r.stamps, r.f.Find(t, symbol.KindStruct, "", tt.giveDecl))
			expect.Equal(t, got.Code, sim.BadPerform, "the finding has the code of a bad perform")
			expect.Equal(t, got.Msg, tt.wantMsg, "the finding names the effect and the remedy")
		})
	}

	t.Run("records a convergence without a bound", func(t *testing.T) {
		t.Parallel()

		m := generate(t, plainPackage, plainSimPackage).models(t)["Plain"]
		expect.Equal(t, m.Converges, []sim.Convergence{{
			Statement: "The world settles once the faults end", Func: "settles",
		}}, "the convergence has no bound")
	})

	t.Run("returns from Keys the error of a namespace that another registrant claimed", func(t *testing.T) {
		t.Parallel()

		r := meta.NewRegistry()
		assert.NoError(t, r.For(otherRegistrant).ClaimNamespace(sim.KeyDefect.Namespace()),
			"another registrant claims the namespace of sim")
		keys, provides := sim.New().(plugin.KeyProvider)
		assert.True(t, provides, "sim registers keys")
		assert.HasError(t, keys.Keys(r.For(string(sim.ID))), "sim cannot register its keys")
	})

	t.Run("passes the plugin suite over the register", func(t *testing.T) {
		t.Parallel()

		plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
			p := sim.New()
			return p, fixture(tb, p, registerPackage, registerSimPackage).Fixture
		})
	})

	t.Run("records the check of each claim of the spec whose generator is sim", func(t *testing.T) {
		t.Parallel()

		m := generate(t, registerPackage, registerSimPackage).models(t)["Machine"]
		ids := make([]string, 0, len(m.Rules))
		for _, rule := range m.Rules {
			ids = append(ids, rule.ID)
		}
		expect.Equal(t, ids, []string{
			"Machine/replies-once", "Machine/replies-before-stop", "Machine/quiet-after-stop",
			"Machine/quiet-after-error", "Machine/fresh-tokens",
		}, "each claim of sim binds to the interface")
		expect.Equal(t, m.Rules[0], sim.Rule{
			ID:        "Machine/replies-once",
			Statement: "each machine of Machine returns at most one reply to each inbound call",
			Claim:     "replies-once",
			Probe:     "replies-at-most-once",
			Reason:    "second-reply",
			Form:      "replies-twice",
		}, "a rule has the statement of its claim with the name of the interface")
	})

	t.Run("records the tests of the simulation", func(t *testing.T) {
		t.Parallel()

		got := generate(t, registerPackage, registerSimPackage).checks(t)["Machine"]
		expect.Equal(t, got.Data, sim.Checks{Machine: "Machine", Simulation: "MachineSimulation"},
			"the tests call the function of the simulation")
		expect.Equal(t, gentest.Referents(got.Decl), []string{"MachineSimulation"},
			"the record of the tests refers to the function of the simulation")
	})

	refusals := []struct {
		name     string
		giveDecl string
		wantCode diag.Code
		wantMsg  string
	}{
		{
			name: "refuses an effect without a performer", giveDecl: "faulty.Send", wantCode: sim.NoPerform,
			wantMsg: "the effect Send of Sender has no directive perform. Write //dokimi:perform " +
				"func=<function> on Send",
		},
		{
			name: "refuses a reply without a field of its directive", giveDecl: "partial.Reply", wantCode: sim.BadReply,
			wantMsg: "the reply Reply has no field Result that its directive names",
		},
		{
			name: "refuses an effect whose performer the simulation lacks", giveDecl: "partial.Flush",
			wantCode: sim.BadFunction,
			wantMsg: "the package of the simulation of Flusher has no function flush with the signature " +
				"func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E], with T the effect Flush",
		},
		{
			name: "refuses a call without a reply", giveDecl: "silent.Asker", wantCode: sim.NoReply,
			wantMsg: "Asker has an action call and no effect under the directive reply, so the run cannot complete " +
				"a call in the history",
		},
		{
			name: "refuses a simulation without a system", giveDecl: "silent.Mute", wantCode: sim.BadFunction,
			wantMsg: "the package of the simulation of Mute has no function under the directive system, so the " +
				"simulation has no processes",
		},
		{
			name: "refuses a second system in the order of the identities", giveDecl: "crookedsim.system",
			wantCode: sim.BadFunction,
			wantMsg:  "system is a second function under the directive system of Crooked, after again",
		},
		{
			name: "refuses a system that takes no world", giveDecl: "crookedsim.broken", wantCode: sim.BadFunction,
			wantMsg: "broken is under the directive system, so it needs the signature func(w *sim.World[E, F]), " +
				"with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a second spec", giveDecl: "crookedsim.second", wantCode: sim.BadFunction,
			wantMsg: "second is a second function under the directive spec of Crooked, after first",
		},
		{
			name: "refuses a spec that takes a parameter", giveDecl: "crookedsim.parameterized",
			wantCode: sim.BadFunction,
			wantMsg: "parameterized is under the directive spec, so it needs the signature func() history.Spec[S] " +
				"or func() history.Subject, with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a second fingerprint", giveDecl: "crookedsim.rehashed", wantCode: sim.BadFunction,
			wantMsg: "rehashed is a second function under the directive fingerprint of Crooked, after hashed",
		},
		{
			name: "refuses an action that takes no case", giveDecl: "crookedsim.caseless", wantCode: sim.BadFunction,
			wantMsg: "caseless is under the directive action, so it needs the signature func(c *prop.Case) " +
				"func(w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a step without the slice of the effects", giveDecl: "bad.ListlessStep",
			wantCode: sim.BadMachine,
			wantMsg: "ListlessStep has no method Step(ev E, out []F) []F with interface types E and F, so it is " +
				"no sans-IO machine",
		},
		{
			name: "refuses a step on a struct", giveDecl: "bad.ConcreteStep", wantCode: sim.BadMachine,
			wantMsg: "ConcreteStep has no method Step(ev E, out []F) []F with interface types E and F, so it is " +
				"no sans-IO machine",
		},
		{
			name: "refuses a start that is no event", giveDecl: "bad.StartOutside", wantCode: sim.BadLifecycle,
			wantMsg: "the param start of StartOutside names no event of Event with one field of a reply token",
		},
		{
			name: "refuses a fingerprint that returns a signed count", giveDecl: "crookedsim.counted",
			wantCode: sim.BadFunction,
			wantMsg: "counted is under the directive fingerprint, so it needs the signature " +
				"func(w *sim.World[E, F]) uint64, with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a spec that returns no spec", giveDecl: "crookedsim.numbered", wantCode: sim.BadFunction,
			wantMsg: "numbered is under the directive spec, so it needs the signature func() history.Spec[S] or " +
				"func() history.Subject, with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an action that runs without a world", giveDecl: "crookedsim.idle", wantCode: sim.BadFunction,
			wantMsg: "idle is under the directive action, so it needs the signature func(c *prop.Case) " +
				"func(w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an action call without a client", giveDecl: "crookedsim.anonymous",
			wantCode: sim.BadFunction,
			wantMsg: "anonymous is under the directive action, so it needs the signature func(c *prop.Case) " +
				"func(w *sim.World[E, F], client int), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an invariant without the handle of a check", giveDecl: "crookedsim.unchecked",
			wantCode: sim.BadFunction,
			wantMsg: "unchecked is under the directive invariant, so it needs the signature func(h *sim.Check, " +
				"w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an invariant whose defect is no system", giveDecl: "crookedsim.misplanted",
			wantCode: sim.BadFunction,
			wantMsg:  "the defect counted of the invariant misplanted needs the signature func(w *sim.World[E, F])",
		},
		{
			name: "refuses a convergence whose bound is no duration", giveDecl: "crookedsim.unbounded",
			wantCode: sim.BadFunction,
			wantMsg:  "the bound late of the convergence unbounded needs the signature func() time.Duration",
		},
		{
			name: "refuses a convergence without the handle of a check", giveDecl: "crookedsim.uncheckedConvergence",
			wantCode: sim.BadFunction,
			wantMsg: "uncheckedConvergence is under the directive converges, so it needs the signature " +
				"func(h *sim.Check, w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a scenario that takes two worlds", giveDecl: "crookedsim.twice", wantCode: sim.BadFunction,
			wantMsg: "twice is under the directive scenario, so it needs the signature func(w *sim.World[E, F]), " +
				"with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an interface without a step method", giveDecl: "bad.NoStep", wantCode: sim.BadMachine,
			wantMsg: "NoStep has no method Step(ev E, out []F) []F with interface types E and F, so it is no " +
				"sans-IO machine",
		},
		{
			name: "refuses a step method without the slice of the effects", giveDecl: "bad.BadStep",
			wantCode: sim.BadMachine,
			wantMsg: "BadStep has no method Step(ev E, out []F) []F with interface types E and F, so it is no " +
				"sans-IO machine",
		},
		{
			name: "refuses an interface without an error method", giveDecl: "bad.NoErr", wantCode: sim.BadMachine,
			wantMsg: "NoErr has no method Err() error, so it is no sans-IO machine",
		},
		{
			name: "refuses an error method that returns no error", giveDecl: "bad.BadErr", wantCode: sim.BadMachine,
			wantMsg: "BadErr has no method Err() error, so it is no sans-IO machine",
		},
		{
			name: "refuses a start without a reply token", giveDecl: "bad.StartWithoutToken",
			wantCode: sim.BadLifecycle,
			wantMsg: "the param start of StartWithoutToken names no event of Event with one field of a " +
				"reply token",
		},
		{
			name: "refuses a start with two reply tokens", giveDecl: "bad.StartTwice", wantCode: sim.BadLifecycle,
			wantMsg: "the param start of StartTwice names no event of Event with one field of a reply token",
		},
		{
			name: "refuses a stop that is no event", giveDecl: "bad.StopOutside", wantCode: sim.BadLifecycle,
			wantMsg: "the param stop of StopOutside names no event of Event",
		},
		{
			name: "refuses a cancel without a reply token", giveDecl: "bad.CancelWithoutToken",
			wantCode: sim.BadLifecycle,
			wantMsg: "the param cancel of CancelWithoutToken names no event of Event with one field of a " +
				"reply token",
		},
		{
			name: "refuses a tolerated fault without a count", giveDecl: "bad.CountlessFault", wantCode: sim.BadParam,
			wantMsg: `the param tolerate of CountlessFault has the entry "crash", which is no fault with a count ` +
				"of at least 1, such as crash:1",
		},
		{
			name: "refuses a tolerated fault that the world does not have", giveDecl: "bad.UnknownFault",
			wantCode: sim.BadParam,
			wantMsg: `the param tolerate of UnknownFault has the entry "meteor:1", which is no fault with a count ` +
				"of at least 1, such as crash:1",
		},
		{
			name: "refuses a required share above 1", giveDecl: "bad.WholeShare", wantCode: sim.BadParam,
			wantMsg: `the param require of WholeShare has the entry "crash:2", which is no fault with a share ` +
				"above 0 and at most 1, such as crash:0.2",
		},
		{
			name: "refuses PCT without depth", giveDecl: "bad.ShallowStrategy", wantCode: sim.BadParam,
			wantMsg: `the param strategy of ShallowStrategy is "pct:0", which is neither uniform nor pct with a ` +
				"depth of at least 1, such as pct:3",
		},
		{
			name: "refuses a strategy that the world does not have", giveDecl: "bad.UnknownStrategy",
			wantCode: sim.BadParam,
			wantMsg: `the param strategy of UnknownStrategy is "random:3", which is neither uniform nor pct with a ` +
				"depth of at least 1, such as pct:3",
		},
		{
			name: "refuses an error method that takes a parameter", giveDecl: "bad.ParamErr", wantCode: sim.BadMachine,
			wantMsg: "ParamErr has no method Err() error, so it is no sans-IO machine",
		},
		{
			name: "refuses a step method with a second result", giveDecl: "bad.TwoStepReturns",
			wantCode: sim.BadMachine,
			wantMsg: "TwoStepReturns has no method Step(ev E, out []F) []F with interface types E and F, so it is " +
				"no sans-IO machine",
		},
		{
			name: "refuses a step method that returns one effect", giveDecl: "bad.PointerReturn",
			wantCode: sim.BadMachine,
			wantMsg: "PointerReturn has no method Step(ev E, out []F) []F with interface types E and F, so it is " +
				"no sans-IO machine",
		},
		{
			name: "refuses a step method that returns events", giveDecl: "bad.MismatchedStep",
			wantCode: sim.BadMachine,
			wantMsg: "MismatchedStep has no method Step(ev E, out []F) []F with interface types E and F, so it is " +
				"no sans-IO machine",
		},
		{
			name: "refuses a step method on effects of a struct", giveDecl: "bad.ConcreteEffects",
			wantCode: sim.BadMachine,
			wantMsg: "ConcreteEffects has no method Step(ev E, out []F) []F with interface types E and F, so it is " +
				"no sans-IO machine",
		},
		{
			name: "refuses a tolerated count beyond the range of an int", giveDecl: "bad.HugeCount",
			wantCode: sim.BadParam,
			wantMsg: `the param tolerate of HugeCount has the entry "crash:99999999999999999999", which is no ` +
				"fault with a count of at least 1, such as crash:1",
		},
		{
			name: "refuses a tolerated count of 0", giveDecl: "bad.ZeroCount", wantCode: sim.BadParam,
			wantMsg: `the param tolerate of ZeroCount has the entry "crash:0", which is no fault with a count of ` +
				"at least 1, such as crash:1",
		},
		{
			name: "refuses a required share of 0", giveDecl: "bad.ZeroShare", wantCode: sim.BadParam,
			wantMsg: `the param require of ZeroShare has the entry "crash:0", which is no fault with a share ` +
				"above 0 and at most 1, such as crash:0.2",
		},
		{
			name: "refuses a required fault that the world does not have", giveDecl: "bad.UnknownShare",
			wantCode: sim.BadParam,
			wantMsg: `the param require of UnknownShare has the entry "meteor:0.5", which is no fault with a ` +
				"share above 0 and at most 1, such as crash:0.2",
		},
		{
			name: "refuses a depth of PCT beyond the range of an int", giveDecl: "bad.HugeDepth",
			wantCode: sim.BadParam,
			wantMsg: `the param strategy of HugeDepth is "pct:99999999999999999999", which is neither uniform ` +
				"nor pct with a depth of at least 1, such as pct:3",
		},
		{
			name: "refuses an invariant that returns an error", giveDecl: "crookedsim.fallible",
			wantCode: sim.BadFunction,
			wantMsg: "fallible is under the directive invariant, so it needs the signature func(h *sim.Check, " +
				"w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a convergence whose bound takes a parameter", giveDecl: "crookedsim.drifting",
			wantCode: sim.BadFunction,
			wantMsg:  "the bound lateBy of the convergence drifting needs the signature func() time.Duration",
		},
		{
			name: "refuses an action whose run returns the world", giveDecl: "crookedsim.inverted",
			wantCode: sim.BadFunction,
			wantMsg: "inverted is under the directive action, so it needs the signature func(c *prop.Case) " +
				"func(w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an action with a second result", giveDecl: "crookedsim.doubled", wantCode: sim.BadFunction,
			wantMsg: "doubled is under the directive action, so it needs the signature func(c *prop.Case) " +
				"func(w *sim.World[E, F]), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses an action call whose client is no int", giveDecl: "crookedsim.misclient",
			wantCode: sim.BadFunction,
			wantMsg: "misclient is under the directive action, so it needs the signature func(c *prop.Case) " +
				"func(w *sim.World[E, F], client int), with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a spec with a second result", giveDecl: "crookedsim.twofold", wantCode: sim.BadFunction,
			wantMsg: "twofold is under the directive spec, so it needs the signature func() history.Spec[S] or " +
				"func() history.Subject, with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a fingerprint that takes no world", giveDecl: "crookedsim.unseen",
			wantCode: sim.BadFunction,
			wantMsg: "unseen is under the directive fingerprint, so it needs the signature " +
				"func(w *sim.World[E, F]) uint64, with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a scenario on a type that is no world", giveDecl: "crookedsim.paired",
			wantCode: sim.BadFunction,
			wantMsg: "paired is under the directive scenario, so it needs the signature func(w *sim.World[E, F]), " +
				"with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a scenario on a world of other interfaces", giveDecl: "crookedsim.foreign",
			wantCode: sim.BadFunction,
			wantMsg: "foreign is under the directive scenario, so it needs the signature func(w *sim.World[E, F]), " +
				"with E and F the event and the effect interfaces of Crooked",
		},
		{
			name: "refuses a performer that takes a world in place of a process", giveDecl: "partial.Wrong",
			wantCode: sim.BadFunction,
			wantMsg: "the package of the simulation of Flusher has no function wrong with the signature " +
				"func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E], with T the effect Wrong",
		},
		{
			name: "refuses a performer of another effect", giveDecl: "partial.Pair", wantCode: sim.BadFunction,
			wantMsg: "the package of the simulation of Flusher has no function pair with the signature " +
				"func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E], with T the effect Pair",
		},
		{
			name: "refuses a performer that returns an outcome of another package", giveDecl: "partial.Odd",
			wantCode: sim.BadFunction,
			wantMsg: "the package of the simulation of Flusher has no function odd with the signature " +
				"func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E], with T the effect Odd",
		},
		{
			name: "refuses a performer that returns an outcome of an effect", giveDecl: "partial.Skew",
			wantCode: sim.BadFunction,
			wantMsg: "the package of the simulation of Flusher has no function skew with the signature " +
				"func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E], with T the effect Skew",
		},
		{
			name: "refuses a call of a function of the package", giveDecl: "misfit.Unhosted",
			wantCode: sim.BadPerform, wantMsg: "the call helper of Unhosted is no method of an interface",
		},
		{
			name: "refuses a call of an effect without a token", giveDecl: "misfit.Tokenless",
			wantCode: sim.BadPerform,
			wantMsg:  "the effect Tokenless has no field of a token, which the event of its result takes",
		},
		{
			name: "refuses a call whose parameter the effect has no field for", giveDecl: "misfit.Keyless",
			wantCode: sim.BadPerform,
			wantMsg:  "the effect Keyless has no field Key for the parameter key of Store.Put",
		},
		{
			name: "refuses a call whose event has no field of the token", giveDecl: "misfit.Lost",
			wantCode: sim.BadPerform,
			wantMsg:  "the event Untokened of the result of Lost has no field of the type Token for the token of Lost",
		},
		{
			name: "refuses a call whose event has no field of the error", giveDecl: "misfit.Silent",
			wantCode: sim.BadPerform,
			wantMsg: "the event Errorless of the result of Silent has no field of the type error for the error of " +
				"Store.Put",
		},
		{
			name: "refuses a call whose event has no field of the value", giveDecl: "misfit.Blind",
			wantCode: sim.BadPerform,
			wantMsg: "the event Valueless of the result of Blind has no field of the type string for the value of " +
				"Store.Get",
		},
	}
	fixtures := map[string][]string{
		"misfit":     {"example.com/shop/misfit"},
		"bad":        {"example.com/shop/bad"},
		"faulty":     {"example.com/shop/faulty"},
		"partial":    {"example.com/shop/partial", "example.com/shop/partial/partialsim"},
		"silent":     {"example.com/shop/silent", "example.com/shop/silent/silentsim"},
		"crookedsim": {"example.com/shop/crooked", "example.com/shop/crooked/crookedsim"},
	}
	for _, tt := range refusals {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pkg, _, _ := strings.Cut(tt.giveDecl, ".")
			r := generate(t, fixtures[pkg]...)
			got := gentest.Finding(t, r.res, declOf(t, r, tt.giveDecl))
			expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of the defect")
			expect.Equal(t, got.Msg, tt.wantMsg, "the finding names the declaration and the defect")
			expect.Empty(t, r.models(t), "sim records no simulation of a machine that it refuses")
		})
	}

	t.Run("binds a function whose world has a name of its own", func(t *testing.T) {
		t.Parallel()

		r := generate(t, fixtures["crookedsim"]...)
		for d := range r.res.Sink.All() {
			expect.NotEqual(t, d.Pos, declOf(t, r, "crookedsim.aliased").Position(),
				"the scenario aliased has the signature of a scenario")
		}
	})

	t.Run("refuses a token that is no unsigned integer type", func(t *testing.T) {
		t.Parallel()

		r := generate(t, "example.com/shop/bad")
		got := gentest.Finding(t, r.stamps, r.f.Find(t, symbol.KindAlias, "", "Name"))
		expect.Equal(t, got.Code, sim.BadToken, "the finding has the code of a bad token")
		expect.Equal(t, got.Msg, "Name is not an unsigned integer type, so it is no token",
			"the finding names the type")
	})

	t.Run("records a reply for each type of the union of a generic reply", func(t *testing.T) {
		t.Parallel()

		m := generate(t, genericPackage, genericSimPackage).models(t)["Replier"]
		instances := make([]string, 0, len(m.Replies))
		for _, reply := range m.Replies {
			instances = append(instances, reply.Type.Spelling+"["+reply.Type.Args[0].Spelling+"]")
		}
		expect.Equal(t, instances, []string{"Note[struct{}]", "Note[int]", "Reply[struct{}]", "Reply[Ack]"},
			"each reply has an instantiation for each type of its union, in the order of the source")
	})

	vague := []struct {
		name     string
		giveDecl string
	}{
		{name: "refuses a generic reply over every type", giveDecl: "Loose"},
		{name: "refuses a generic reply over an approximation", giveDecl: "Rough"},
		{name: "refuses a generic reply over two type parameters", giveDecl: "Twin"},
		{name: "refuses a generic reply over a constraint without a union", giveDecl: "Bare"},
	}
	for _, tt := range vague {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := generate(t, "example.com/shop/vague", "example.com/shop/vague/vaguesim")
			got := gentest.Finding(t, r.res, r.f.Find(t, symbol.KindStruct, "", tt.giveDecl))
			expect.Equal(t, got.Code, sim.BadReply, "the finding has the code of a bad reply")
			expect.Equal(t, got.Msg, "the reply "+tt.giveDecl+" is generic over another constraint than a union of "+
				"types without an approximation, so the anatomy cannot list its instantiations",
				"the finding names the reply and the constraint that its anatomy needs")
			gentest.Unrecorded(t, r.res, r.f.Find(t, symbol.KindInterface, "", "Vague"))
		})
	}
}
