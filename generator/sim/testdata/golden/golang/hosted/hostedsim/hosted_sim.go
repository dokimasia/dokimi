package hostedsim

import (
	"time"

	"example.com/shop/hosted"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// machineAnatomy is how a world reads the events and the effects of
// Machine.
var machineAnatomy = sim.Anatomy[hosted.Event, hosted.Effect]{
	Started: func(reply uint64) hosted.Event {
		return hosted.Started{Reply: hosted.ReplyToken(reply)}
	},
	Stopped: func() hosted.Event { return hosted.Stopped{} },
	Call: func(ev hosted.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case hosted.Await:
			return uint64(ev.Reply), true
		case hosted.Find:
			return uint64(ev.Reply), true
		case hosted.Get:
			return uint64(ev.Reply), true
		case hosted.Put:
			return uint64(ev.Reply), true
		case hosted.Started:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply: func(f hosted.Effect) (uint64, any, error, bool) {
		switch f := f.(type) {
		case hosted.Reply:
			return uint64(f.To), f.Value, f.Err, true
		default:
			return 0, nil, nil, false
		}
	},
	Token: func(f hosted.Effect) (uint64, bool) {
		switch f := f.(type) {
		case hosted.Load:
			return uint64(f.Token), true
		case hosted.Lookup:
			return uint64(f.Token), true
		case hosted.Persist:
			return uint64(f.Token), true
		default:
			return 0, false
		}
	},
	Dependency: func(f hosted.Effect) (string, bool) {
		switch f.(type) {
		case hosted.Load:
			return "hosted.Storage", true
		case hosted.Lookup:
			return "hosted.Registry", true
		case hosted.Persist:
			return "hosted.Storage", true
		default:
			return "", false
		}
	},
	Scope: func(dependency string) sim.Scope {
		switch dependency {
		case "hosted.Storage":
			return sim.ScopeHost
		default:
			return sim.ScopeWorld
		}
	},
	Fail: func(f hosted.Effect, err error) (hosted.Event, bool) {
		switch f := f.(type) {
		case hosted.Load:
			return hosted.Loaded{Token: f.Token, Err: err}, true
		case hosted.Lookup:
			return hosted.Found{Token: f.Token, Err: err}, true
		case hosted.Persist:
			return hosted.Persisted{Token: f.Token, Err: err}, true
		default:
			return nil, false
		}
	},
	Errors: func(f hosted.Effect) []error {
		switch f.(type) {
		case hosted.Load:
			return []error{hosted.ErrUnavailable}
		case hosted.Persist:
			return []error{hosted.ErrUnavailable}
		default:
			return nil
		}
	},
	Timeout: func(hosted.Effect) time.Duration { return 0 },
	Perform: performMachine,
}

// performMachine performs the effect f of the process p of Machine through
// the function or the call of its directive perform.
func performMachine(
	w *sim.World[hosted.Event, hosted.Effect],
	p *sim.Process[hosted.Event, hosted.Effect],
	f hosted.Effect,
) sim.Outcome[hosted.Event, hosted.Effect] {
	switch f := f.(type) {
	case hosted.Load:
		ev := hosted.Loaded{Token: f.Token}
		ev.Value, ev.Err = sim.Dependency[hosted.Storage](w, p).Read(w.Context(), f.Key)
		return sim.Outcome[hosted.Event, hosted.Effect]{Event: ev}
	case hosted.Lookup:
		ev := hosted.Found{Token: f.Token}
		ev.Address, ev.Err = sim.Dependency[hosted.Registry](w, p).Lookup(w.Context(), f.Name)
		return sim.Outcome[hosted.Event, hosted.Effect]{Event: ev}
	case hosted.Persist:
		return persist(w, p, f)
	default:
		return sim.Outcome[hosted.Event, hosted.Effect]{}
	}
}

// MachineSimulation returns the simulation of the system of the machines
// of Machine, with its checks.
func MachineSimulation() sim.Simulation[hosted.Event, hosted.Effect, struct{}] {
	return sim.Simulation[hosted.Event, hosted.Effect, struct{}]{
		Anatomy: machineAnatomy,
		System:  func(*prop.Case) func(w *sim.World[hosted.Event, hosted.Effect]) { return system },
		Tolerate: map[sim.Fault]int{
			sim.FaultCrash:  1,
			sim.FaultOutage: 1,
		},
		Actions: []sim.Action[hosted.Event, hosted.Effect]{
			{
				Name:  "await",
				Call:  true,
				Input: func(c *prop.Case) any { return await(c) },
				Run: func(w *sim.World[hosted.Event, hosted.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[hosted.Event, hosted.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "find",
				Call:  true,
				Input: func(c *prop.Case) any { return find(c) },
				Run: func(w *sim.World[hosted.Event, hosted.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[hosted.Event, hosted.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "get",
				Call:  true,
				Input: func(c *prop.Case) any { return get(c) },
				Run: func(w *sim.World[hosted.Event, hosted.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[hosted.Event, hosted.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "put",
				Call:  true,
				Input: func(c *prop.Case) any { return put(c) },
				Run: func(w *sim.World[hosted.Event, hosted.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[hosted.Event, hosted.Effect], int))
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
		Invariants: []sim.Invariant[hosted.Event, hosted.Effect]{
			{
				Statement: "The world has the member b",
				Check:     spawned,
				Defect:    func(*prop.Case) func(w *sim.World[hosted.Event, hosted.Effect]) { return lonely },
			},
		},
	}
}
