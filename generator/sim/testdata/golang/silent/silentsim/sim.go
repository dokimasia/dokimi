// Package silentsim contains a simulation with a call of a client.
package silentsim

import (
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/silent"
)

// system spawns no process.
//
//dokimi:system
func system(*sim.World[silent.Event, silent.Effect]) {}

// ask is a call of a client that the machine cannot answer.
//
//dokimi:action call
func ask(*prop.Case) func(*sim.World[silent.Event, silent.Effect], int) {
	return func(*sim.World[silent.Event, silent.Effect], int) {}
}
