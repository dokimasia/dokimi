package remotesim

import (
	"time"

	"example.com/shop/remote"
	"example.com/shop/remote/store"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// machineAnatomy is how a world reads the events and the effects of
// Machine.
var machineAnatomy = sim.Anatomy[remote.Event, remote.Effect]{
	Started: func(reply uint64) remote.Event {
		return remote.Started{Reply: remote.ReplyToken(reply)}
	},
	Stopped: func() remote.Event { return remote.Stopped{} },
	Call: func(ev remote.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case remote.Started:
			return uint64(ev.Reply), true
		case remote.Watch:
			return uint64(ev.Reply), true
		case remote.Write:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply: func(f remote.Effect) (uint64, any, error, bool) {
		switch f := f.(type) {
		case remote.Reply:
			return uint64(f.To), f.Value, f.Err, true
		default:
			return 0, nil, nil, false
		}
	},
	Token: func(f remote.Effect) (uint64, bool) {
		switch f := f.(type) {
		case remote.Get:
			return uint64(f.Token), true
		case remote.Load:
			return uint64(f.Token), true
		case remote.Put:
			return uint64(f.Token), true
		default:
			return 0, false
		}
	},
	Dependency: func(f remote.Effect) (string, bool) {
		switch f.(type) {
		case remote.Get:
			return "remote.Writer", true
		case remote.Load:
			return "store.Reader", true
		case remote.Put:
			return "remote.Writer", true
		default:
			return "", false
		}
	},
	Fail: func(f remote.Effect, err error) (remote.Event, bool) {
		switch f := f.(type) {
		case remote.Get:
			return remote.Fetched{Token: f.Token, Err: err}, true
		case remote.Load:
			return remote.Loaded{Token: f.Token, Err: err}, true
		case remote.Put:
			return remote.Stored{Token: f.Token, Err: err}, true
		default:
			return nil, false
		}
	},
	Errors: func(f remote.Effect) []error {
		switch f.(type) {
		case remote.Put:
			return []error{remote.ErrFull}
		default:
			return nil
		}
	},
	Timeout: func(f remote.Effect) time.Duration {
		switch f := f.(type) {
		case remote.Put:
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
	w *sim.World[remote.Event, remote.Effect],
	p *sim.Process[remote.Event, remote.Effect],
	f remote.Effect,
) sim.Outcome[remote.Event, remote.Effect] {
	switch f := f.(type) {
	case remote.Get:
		ev := remote.Fetched{Token: f.Token}
		ev.Value, ev.Err = sim.Dependency[remote.Writer](w, p).Get(w.Context(), f.Key)
		return sim.Outcome[remote.Event, remote.Effect]{Event: ev}
	case remote.Load:
		ev := remote.Loaded{Token: f.Token}
		ev.Value, ev.Err = sim.Dependency[store.Reader](w, p).Read(w.Context(), f.Key)
		return sim.Outcome[remote.Event, remote.Effect]{Event: ev}
	case remote.Put:
		ev := remote.Stored{Token: f.Token}
		ev.Err = sim.Dependency[remote.Writer](w, p).Put(w.Context(), f.Key, f.Value)
		return sim.Outcome[remote.Event, remote.Effect]{Event: ev}
	default:
		return sim.Outcome[remote.Event, remote.Effect]{}
	}
}

// MachineSimulation returns the simulation of the system of the machines
// of Machine, with its checks.
func MachineSimulation() sim.Simulation[remote.Event, remote.Effect, struct{}] {
	return sim.Simulation[remote.Event, remote.Effect, struct{}]{
		Anatomy: machineAnatomy,
		System:  func(*prop.Case) func(w *sim.World[remote.Event, remote.Effect]) { return system },
		Actions: []sim.Action[remote.Event, remote.Effect]{
			{
				Name:  "watch",
				Call:  true,
				Input: func(c *prop.Case) any { return watch(c) },
				Run: func(w *sim.World[remote.Event, remote.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[remote.Event, remote.Effect], int))
					run(w, client)
				},
			},
			{
				Name:  "write",
				Call:  true,
				Input: func(c *prop.Case) any { return write(c) },
				Run: func(w *sim.World[remote.Event, remote.Effect], client int, in any) {
					run, _ := in.(func(*sim.World[remote.Event, remote.Effect], int))
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
	}
}
