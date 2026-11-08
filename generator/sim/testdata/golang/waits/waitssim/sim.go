// Package waitssim contains the simulation of the leader of waits.
package waitssim

import (
	"context"

	"go.dokimi.dev/sim"

	"example.com/shop/waits"
)

// world is the world of the leader of waits, and outcome the outcome of its
// effects.
type (
	world   = sim.World[waits.Event, waits.Effect]
	outcome = sim.Outcome[waits.Event, waits.Effect]
)

// store is the double of the store. It grants each campaign the term 1.
type store struct{}

// Campaign campaigns for the term of the scope.
func (store) Campaign(context.Context, string) error { return nil }

// system spawns the leader over a store.
//
//dokimi:system
func system(w *world) {
	w.Spawn("leader", func(*sim.Clock) sim.Machine[waits.Event, waits.Effect] {
		return waits.New()
	}, waits.Store(store{}))
}

// campaign campaigns for the scope of the effect f through the store of the
// process p. Its wait delivers the grant of the term 1, and the end of the
// term once the world cancels the campaign.
func campaign(w *world, p *sim.Process[waits.Event, waits.Effect], f waits.Campaign) outcome {
	if err := sim.Dependency[waits.Store](w, p).Campaign(w.Context(), f.Scope); err != nil {
		return outcome{Event: waits.Campaigned{Err: err, Token: f.Token}}
	}
	return stream(waits.Campaigned{Token: f.Token, Term: 1}, waits.TermEnded{Token: f.Token})
}

// resign resigns the term of the campaign of the effect f. The cancel of the
// campaign ends the term, so the resign has no outcome.
func resign(*world, *sim.Process[waits.Event, waits.Effect], waits.Resign) outcome { return outcome{} }

// watch watches the feed of the leader. Its wait delivers one item, and ends
// without an item once the world cancels the watch. It takes a pointer to the
// effect.
func watch(_ *world, _ *sim.Process[waits.Event, waits.Effect], f *waits.Watch) outcome {
	return stream(waits.Watched{Token: f.Token}, nil)
}

// stream returns the outcome of a wait that delivers first, and then last once
// the world cancels the wait. A nil last ends the wait without an event.
func stream(first, last waits.Event) outcome {
	delivered, cancelled, ended := false, false, false
	return outcome{
		Next: func() (waits.Event, bool, bool) {
			switch {
			case !delivered:
				delivered = true
				return first, true, false
			case cancelled && !ended:
				ended = true
				return last, last != nil, true
			default:
				return nil, false, ended
			}
		},
		Cancel: func() { cancelled = true },
	}
}
