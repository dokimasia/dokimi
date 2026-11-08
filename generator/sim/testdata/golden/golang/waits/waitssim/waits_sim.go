package waitssim

import (
	"time"

	"example.com/shop/waits"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"
)

// leaderAnatomy is how a world reads the events and the effects of Leader.
var leaderAnatomy = sim.Anatomy[waits.Event, waits.Effect]{
	Started: func(reply uint64) waits.Event {
		return waits.Started{Reply: waits.ReplyToken(reply)}
	},
	Stopped: func() waits.Event { return waits.Stopped{} },
	Call: func(ev waits.Event) (uint64, bool) {
		switch ev := ev.(type) {
		case waits.Started:
			return uint64(ev.Reply), true
		default:
			return 0, false
		}
	},
	Reply: func(f waits.Effect) (uint64, any, error, bool) {
		switch f := f.(type) {
		case waits.Reply:
			return uint64(f.To), f.Value, f.Err, true
		default:
			return 0, nil, nil, false
		}
	},
	Token: func(f waits.Effect) (uint64, bool) {
		switch f := f.(type) {
		case waits.Campaign:
			return uint64(f.Token), true
		case waits.Watch:
			return uint64(f.Token), true
		default:
			return 0, false
		}
	},
	Ends: func(f waits.Effect) (uint64, bool) {
		switch f := f.(type) {
		case waits.Resign:
			return uint64(f.Token), true
		case waits.Unwatch:
			return uint64(f.Token), true
		default:
			return 0, false
		}
	},
	Dependency: func(f waits.Effect) (string, bool) {
		switch f.(type) {
		case waits.Campaign:
			return "waits.Store", true
		default:
			return "", false
		}
	},
	Fail: func(f waits.Effect, err error) (waits.Event, bool) {
		switch f := f.(type) {
		case waits.Campaign:
			return waits.Campaigned{Token: f.Token, Err: err}, true
		default:
			return nil, false
		}
	},
	End: func(f waits.Effect, err error) (waits.Event, bool) {
		switch f := f.(type) {
		case waits.Campaign:
			return waits.TermEnded{Token: f.Token}, true
		default:
			return nil, false
		}
	},
	Errors: func(f waits.Effect) []error {
		switch f.(type) {
		case waits.Campaign:
			return []error{waits.ErrUnavailable}
		default:
			return nil
		}
	},
	Timeout: func(waits.Effect) time.Duration { return 0 },
	Perform: performLeader,
}

// performLeader performs the effect f of the process p of Leader through
// the function or the call of its directive perform.
func performLeader(
	w *sim.World[waits.Event, waits.Effect],
	p *sim.Process[waits.Event, waits.Effect],
	f waits.Effect,
) sim.Outcome[waits.Event, waits.Effect] {
	switch f := f.(type) {
	case waits.Campaign:
		return campaign(w, p, f)
	case waits.Resign:
		return resign(w, p, f)
	case waits.Watch:
		return watch(w, p, &f)
	default:
		return sim.Outcome[waits.Event, waits.Effect]{}
	}
}

// LeaderSimulation returns the simulation of the system of the machines of
// Leader, with its checks.
func LeaderSimulation() sim.Simulation[waits.Event, waits.Effect, struct{}] {
	return sim.Simulation[waits.Event, waits.Effect, struct{}]{
		Anatomy: leaderAnatomy,
		System:  func(*prop.Case) func(w *sim.World[waits.Event, waits.Effect]) { return system },
		Tolerate: map[sim.Fault]int{
			sim.FaultOutage:    1,
			sim.FaultLoseAfter: 1,
		},
		Rules: []sim.Rule{
			{
				ID:        "Leader/replies-once",
				Statement: "each machine of Leader returns at most one reply to each inbound call",
				Claim:     "replies-once",
				Probe:     "replies-at-most-once",
				Reason:    "second-reply",
				Form:      "replies-twice",
			},
			{
				ID:        "Leader/replies-before-stop",
				Statement: "on the stop event, each machine of Leader replies to each inbound call that it has not replied to",
				Claim:     "replies-before-stop",
				Probe:     "replies-on-stop",
				Reason:    "call-without-reply",
				Form:      "forgets-a-call-on-stop",
			},
			{
				ID:        "Leader/quiet-after-stop",
				Statement: "from the stop event on, each machine of Leader returns replies alone",
				Claim:     "quiet-after-stop",
				Probe:     "only-replies-after-stop",
				Reason:    "effect-after-stop",
				Form:      "writes-after-stop",
			},
			{
				ID:        "Leader/quiet-after-error",
				Statement: "each machine of Leader returns no effect after the step in which its error method returns an error",
				Claim:     "quiet-after-error",
				Probe:     "no-effects-after-error",
				Reason:    "effect-after-error",
				Form:      "steps-after-error",
			},
			{
				ID:        "Leader/fresh-tokens",
				Statement: "each effect of a machine of Leader that starts a call or a wait has a token that no earlier effect of the machine in its epoch has",
				Claim:     "fresh-tokens",
				Probe:     "tokens-differ",
				Reason:    "reused-token",
				Form:      "reuses-a-token",
			},
		},
		MeanSteps: 1,
		MaxSteps:  200,
	}
}
