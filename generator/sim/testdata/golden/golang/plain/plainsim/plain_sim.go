package plainsim

import (
	"time"

	"example.com/shop/plain"
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// plainAnatomy is how a world reads the events and the effects of Plain.
var plainAnatomy = sim.Anatomy[plain.Event, plain.Effect]{
	Started: func(reply uint64) plain.Event {
		return plain.Started{Reply: plain.ReplyToken(reply)}
	},
	Stopped: func() plain.Event { return plain.Stopped{} },
	Call: func(ev plain.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case plain.Started:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply:      func(plain.Effect) (uint64, any, error, bool) { return 0, nil, nil, false },
	Token:      func(plain.Effect) (uint64, bool) { return 0, false },
	Dependency: func(plain.Effect) (string, bool) { return "", false },
	Fail:       func(plain.Effect, error) (plain.Event, bool) { return nil, false },
	Errors:     func(plain.Effect) []error { return nil },
	Timeout:    func(plain.Effect) time.Duration { return 0 },
	Perform:    performPlain,
}

// performPlain performs the effect f of the process p of Plain through the
// function or the call of its directive perform.
func performPlain(
	w *sim.World[plain.Event, plain.Effect],
	p *sim.Process[plain.Event, plain.Effect],
	f plain.Effect,
) sim.Outcome[plain.Event, plain.Effect] {
	switch f := f.(type) {
	case plain.Beta:
		return beta(w, p, f)
	default:
		return sim.Outcome[plain.Event, plain.Effect]{}
	}
}

// PlainSimulation returns the simulation of the system of the machines of
// Plain, with its checks.
func PlainSimulation() sim.Simulation[plain.Event, plain.Effect, []history.Operation] {
	return sim.Simulation[plain.Event, plain.Effect, []history.Operation]{
		Anatomy: plainAnatomy,
		System:  func(*prop.Case) func(w *sim.World[plain.Event, plain.Effect]) { return system },
		Spec:    history.SpecFrom(subject),
		Tolerate: map[sim.Fault]int{
			sim.FaultPause: 1,
		},
		Require: map[sim.Fault]float64{
			sim.FaultPause: 1,
		},
		Rules: []sim.Rule{
			{
				ID:        "Plain/replies-once",
				Statement: "each machine of Plain returns at most one reply to each inbound call",
				Claim:     "replies-once",
				Probe:     "replies-at-most-once",
				Reason:    "second-reply",
				Form:      "replies-twice",
			},
			{
				ID:        "Plain/replies-before-stop",
				Statement: "on the stop event, each machine of Plain replies to each inbound call that it has not replied to",
				Claim:     "replies-before-stop",
				Probe:     "replies-on-stop",
				Reason:    "call-without-reply",
				Form:      "forgets-a-call-on-stop",
			},
			{
				ID:        "Plain/quiet-after-stop",
				Statement: "from the stop event on, each machine of Plain returns replies alone",
				Claim:     "quiet-after-stop",
				Probe:     "only-replies-after-stop",
				Reason:    "effect-after-stop",
				Form:      "writes-after-stop",
			},
			{
				ID:        "Plain/quiet-after-error",
				Statement: "each machine of Plain returns no effect after the step in which its error method returns an error",
				Claim:     "quiet-after-error",
				Probe:     "no-effects-after-error",
				Reason:    "effect-after-error",
				Form:      "steps-after-error",
			},
			{
				ID:        "Plain/fresh-tokens",
				Statement: "each effect of a machine of Plain that starts a call or a wait has a token that no earlier effect of the machine in its epoch has",
				Claim:     "fresh-tokens",
				Probe:     "tokens-differ",
				Reason:    "reused-token",
				Form:      "reuses-a-token",
			},
		},
		Converges: []sim.Convergence[plain.Event, plain.Effect]{
			{Statement: "The world settles once the faults end", Check: settles},
		},
		Strategy:  sim.Uniform(),
		Isolation: sim.Serializable,
	}
}
