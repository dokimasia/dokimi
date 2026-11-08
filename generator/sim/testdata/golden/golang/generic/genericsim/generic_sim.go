package genericsim

import (
	"time"

	"example.com/shop/generic"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// replierAnatomy is how a world reads the events and the effects of
// Replier.
var replierAnatomy = sim.Anatomy[generic.Event, generic.Effect]{
	Started: func(reply uint64) generic.Event {
		return generic.Started{Reply: generic.ReplyToken(reply)}
	},
	Stopped: func() generic.Event { return generic.Stopped{} },
	Call: func(ev generic.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case generic.Started:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply: func(f generic.Effect) (uint64, any, error, bool) {
		switch f := f.(type) {
		case generic.Note[struct{}]:
			return uint64(f.To), f.Value, f.Err, true
		case generic.Note[int]:
			return uint64(f.To), f.Value, f.Err, true
		case generic.Reply[struct{}]:
			return uint64(f.To), f.Value, f.Err, true
		case generic.Reply[generic.Ack]:
			return uint64(f.To), f.Value, f.Err, true
		default:
			return 0, nil, nil, false
		}
	},
	Token:      func(generic.Effect) (uint64, bool) { return 0, false },
	Dependency: func(generic.Effect) (string, bool) { return "", false },
	Fail:       func(generic.Effect, error) (generic.Event, bool) { return nil, false },
	Errors:     func(generic.Effect) []error { return nil },
	Timeout:    func(generic.Effect) time.Duration { return 0 },
	Perform:    performReplier,
}

// performReplier performs the effect f of the process p of Replier through
// the function or the call of its directive perform.
func performReplier(
	w *sim.World[generic.Event, generic.Effect],
	p *sim.Process[generic.Event, generic.Effect],
	f generic.Effect,
) sim.Outcome[generic.Event, generic.Effect] {
	return sim.Outcome[generic.Event, generic.Effect]{}
}

// ReplierSimulation returns the simulation of the system of the machines
// of Replier, with its checks.
func ReplierSimulation() sim.Simulation[generic.Event, generic.Effect, struct{}] {
	return sim.Simulation[generic.Event, generic.Effect, struct{}]{
		Anatomy: replierAnatomy,
		System:  func(*prop.Case) func(w *sim.World[generic.Event, generic.Effect]) { return system },
		Rules: []sim.Rule{
			{
				ID:        "Replier/replies-once",
				Statement: "each machine of Replier returns at most one reply to each inbound call",
				Claim:     "replies-once",
				Probe:     "replies-at-most-once",
				Reason:    "second-reply",
				Form:      "replies-twice",
			},
			{
				ID:        "Replier/replies-before-stop",
				Statement: "on the stop event, each machine of Replier replies to each inbound call that it has not replied to",
				Claim:     "replies-before-stop",
				Probe:     "replies-on-stop",
				Reason:    "call-without-reply",
				Form:      "forgets-a-call-on-stop",
			},
			{
				ID:        "Replier/quiet-after-stop",
				Statement: "from the stop event on, each machine of Replier returns replies alone",
				Claim:     "quiet-after-stop",
				Probe:     "only-replies-after-stop",
				Reason:    "effect-after-stop",
				Form:      "writes-after-stop",
			},
			{
				ID:        "Replier/quiet-after-error",
				Statement: "each machine of Replier returns no effect after the step in which its error method returns an error",
				Claim:     "quiet-after-error",
				Probe:     "no-effects-after-error",
				Reason:    "effect-after-error",
				Form:      "steps-after-error",
			},
			{
				ID:        "Replier/fresh-tokens",
				Statement: "each effect of a machine of Replier that starts a call or a wait has a token that no earlier effect of the machine in its epoch has",
				Claim:     "fresh-tokens",
				Probe:     "tokens-differ",
				Reason:    "reused-token",
				Form:      "reuses-a-token",
			},
		},
	}
}
