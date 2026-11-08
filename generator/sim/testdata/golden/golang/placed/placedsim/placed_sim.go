package placedsim

import (
	"time"

	"example.com/shop/placed"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// machineAnatomy is how a world reads the events and the effects of
// Machine.
var machineAnatomy = sim.Anatomy[placed.Event, placed.Effect]{
	Started: func(reply uint64) placed.Event {
		return placed.Started{Reply: placed.ReplyToken(reply)}
	},
	Stopped: func() placed.Event { return placed.Stopped{} },
	Call: func(ev placed.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case placed.Assign:
			return uint64(ev.Reply), true
		case placed.Await:
			return uint64(ev.Reply), true
		case placed.Started:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply: func(f placed.Effect) (uint64, any, error, bool) {
		switch f := f.(type) {
		case placed.Reply:
			return uint64(f.To), f.Value, f.Err, true
		default:
			return 0, nil, nil, false
		}
	},
	Token: func(f placed.Effect) (uint64, bool) {
		switch f := f.(type) {
		case placed.Write:
			return uint64(f.Token), true
		default:
			return 0, false
		}
	},
	Dependency: func(f placed.Effect) (string, bool) {
		switch f.(type) {
		case placed.Write:
			return "placed.Store", true
		default:
			return "", false
		}
	},
	Fail: func(f placed.Effect, err error) (placed.Event, bool) {
		switch f := f.(type) {
		case placed.Write:
			return placed.Written{Token: f.Token, Err: err}, true
		default:
			return nil, false
		}
	},
	Errors: func(f placed.Effect) []error {
		switch f.(type) {
		case placed.Write:
			return []error{placed.ErrUnavailable, placed.ErrStale}
		default:
			return nil
		}
	},
	Timeout: func(f placed.Effect) time.Duration {
		switch f := f.(type) {
		case placed.Write:
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
	w *sim.World[placed.Event, placed.Effect],
	p *sim.Process[placed.Event, placed.Effect],
	f placed.Effect,
) sim.Outcome[placed.Event, placed.Effect] {
	switch f := f.(type) {
	case placed.Write:
		return write(w, p, f)
	default:
		return sim.Outcome[placed.Event, placed.Effect]{}
	}
}

// MachineSimulation returns the simulation of the system of the machines
// of Machine, with its checks.
func MachineSimulation() sim.Simulation[placed.Event, placed.Effect, struct{}] {
	return sim.Simulation[placed.Event, placed.Effect, struct{}]{
		Anatomy: machineAnatomy,
		System:  system,
		Tolerate: map[sim.Fault]int{
			sim.FaultCrash:     1,
			sim.FaultOutage:    1,
			sim.FaultLoseAfter: 1,
		},
		Actions: []sim.Action[placed.Event, placed.Effect]{
			{
				Name:  "assign",
				Call:  true,
				Input: func(c *prop.Case) any { return assign(c) },
				Run: func(w *sim.World[placed.Event, placed.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[placed.Event, placed.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "await",
				Call:  true,
				Input: func(c *prop.Case) any { return await(c) },
				Run: func(w *sim.World[placed.Event, placed.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[placed.Event, placed.Effect], int))
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
		Invariants: []sim.Invariant[placed.Event, placed.Effect]{
			{Statement: "The world has the agent b", Check: spawned, Defect: lonely},
		},
	}
}
