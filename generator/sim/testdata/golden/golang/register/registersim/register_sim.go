package registersim

import (
	"errors"
	"time"

	"example.com/shop/register"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// machineAnatomy is how a world reads the events and the effects of
// Machine.
var machineAnatomy = sim.Anatomy[register.Event, register.Effect]{
	Started: func(reply uint64) register.Event {
		return register.Started{Reply: register.ReplyToken(reply)}
	},
	Stopped: func() register.Event { return register.Stopped{} },
	Cancelled: func(reply uint64) register.Event {
		return register.Cancelled{Reply: register.ReplyToken(reply)}
	},
	Call: func(ev register.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case register.Started:
			return uint64(ev.Reply), true
		case register.Watch:
			return uint64(ev.Reply), true
		case register.Write:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply: func(f register.Effect) (uint64, any, error, bool) {
		switch f := f.(type) {
		case register.Reply:
			return uint64(f.To), f.Value, f.Err, true
		default:
			return 0, nil, nil, false
		}
	},
	Unknown: func(err error) bool {
		return errors.Is(err, register.ErrStopped)
	},
	Token: func(f register.Effect) (uint64, bool) {
		switch f := f.(type) {
		case register.Commit:
			return uint64(f.Token), true
		default:
			return 0, false
		}
	},
	Dependency: func(register.Effect) (string, bool) { return "", false },
	Fail:       func(register.Effect, error) (register.Event, bool) { return nil, false },
	Errors:     func(register.Effect) []error { return nil },
	Timeout: func(f register.Effect) time.Duration {
		switch f := f.(type) {
		case register.Commit:
			return f.Timeout
		default:
			return 0
		}
	},
	Perform: performMachine,
}

// performMachine performs the effect f of the process p of Machine through
// the function or the call of its directive perform.
func performMachine(
	w *sim.World[register.Event, register.Effect],
	p *sim.Process[register.Event, register.Effect],
	f register.Effect,
) sim.Outcome[register.Event, register.Effect] {
	switch f := f.(type) {
	case register.Commit:
		return commit(w, p, f)
	default:
		return sim.Outcome[register.Event, register.Effect]{}
	}
}

// MachineSimulation returns the simulation of the system of the machines
// of Machine, with its checks.
func MachineSimulation() sim.Simulation[register.Event, register.Effect, string] {
	return sim.Simulation[register.Event, register.Effect, string]{
		Anatomy:     machineAnatomy,
		System:      func(*prop.Case) func(w *sim.World[register.Event, register.Effect]) { return system },
		Spec:        spec(),
		Fingerprint: fingerprint,
		Tolerate: map[sim.Fault]int{
			sim.FaultCrash:     1,
			sim.FaultPartition: 1,
		},
		Require: map[sim.Fault]float64{
			sim.FaultCrash: 0.1,
		},
		Actions: []sim.Action[register.Event, register.Effect]{
			{
				Name:  "pause",
				Input: func(c *prop.Case) any { return pause(c) },
				Run: func(w *sim.World[register.Event, register.Effect], _ int, in any) {
					run, _ := in.(func(*sim.World[register.Event, register.Effect]))
					run(w)
				},
			},
			{
				Name:  "rest",
				Drain: true,
				Input: func(c *prop.Case) any { return rest(c) },
				Run: func(w *sim.World[register.Event, register.Effect], _ int, in any) {
					run, _ := in.(func(*sim.World[register.Event, register.Effect]))
					run(w)
				},
			},
			{
				Name:  "settle",
				Call:  true,
				Drain: true,
				Input: func(c *prop.Case) any { return settle(c) },
				Run: func(w *sim.World[register.Event, register.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[register.Event, register.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "watch",
				Call:  true,
				Input: func(c *prop.Case) any { return watch(c) },
				Run: func(w *sim.World[register.Event, register.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[register.Event, register.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "write",
				Call:  true,
				Input: func(c *prop.Case) any { return write(c) },
				Run: func(w *sim.World[register.Event, register.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[register.Event, register.Effect], int))
					run(w, client)
				},
			},
		},
		Rules: []sim.Rule{
			{
				ID:        "Machine/replies-once",
				Statement: "each machine of Machine returns at most one reply to each inbound call",
				Claim:     "replies-once",
				Probe:     "replies-at-most-once",
				Reason:    "second-reply",
				Form:      "replies-twice",
			},
			{
				ID:        "Machine/replies-before-stop",
				Statement: "on the stop event, each machine of Machine replies to each inbound call that it has not replied to",
				Claim:     "replies-before-stop",
				Probe:     "replies-on-stop",
				Reason:    "call-without-reply",
				Form:      "forgets-a-call-on-stop",
			},
			{
				ID:        "Machine/quiet-after-stop",
				Statement: "from the stop event on, each machine of Machine returns replies alone",
				Claim:     "quiet-after-stop",
				Probe:     "only-replies-after-stop",
				Reason:    "effect-after-stop",
				Form:      "writes-after-stop",
			},
			{
				ID:        "Machine/quiet-after-error",
				Statement: "each machine of Machine returns no effect after the step in which its error method returns an error",
				Claim:     "quiet-after-error",
				Probe:     "no-effects-after-error",
				Reason:    "effect-after-error",
				Form:      "steps-after-error",
			},
			{
				ID:        "Machine/fresh-tokens",
				Statement: "each effect of a machine of Machine that starts a call or a wait has a token that no earlier effect of the machine in its epoch has",
				Claim:     "fresh-tokens",
				Probe:     "tokens-differ",
				Reason:    "reused-token",
				Form:      "reuses-a-token",
			},
		},
		Invariants: []sim.Invariant[register.Event, register.Effect]{
			{
				Statement: "The system has the node b",
				Check:     spawned,
				Defect:    func(*prop.Case) func(w *sim.World[register.Event, register.Effect]) { return lonely },
			},
		},
		Converges: []sim.Convergence[register.Event, register.Effect]{
			{
				Statement: "The system has the node b once the faults have ended",
				Check:     recovered,
				Defect:    func(*prop.Case) func(w *sim.World[register.Event, register.Effect]) { return lonely },
				Within:    bound(),
			},
		},
		Scenarios: []sim.Scenario[register.Event, register.Effect]{
			{Name: "writeOnA", Run: writeOnA},
		},
		Strategy: sim.PCT(2),
	}
}
