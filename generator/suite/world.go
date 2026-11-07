// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"slices"

	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
)

// The words of the names of the world of a declaration: the variable of the
// world, as StoreWorld, and the function of its operations, as
// StoreOperations.
const (
	worldWord      = "World"
	operationsWord = "Operations"
)

// worldProbes are the probes of the claims of the world of a subject. Each
// takes the world of the declaration as its one input.
var worldProbes = []spec.ProbeID{spec.FaultyWorldHistory, spec.WorldRecovery}

// worldClients is the number of the clients of the world of a declaration
// whose each operation has the mixin concurrent. The world of any other
// declaration has one client, its zero Clients.
const worldClients = 3

// World is the world of the subject of a declaration that the checks of the
// world probes run. Name is the variable of the world, and Operations the
// function that returns its operations over the pools of the keys of a case.
// The operations are the actions of the generated sequences of the
// declaration, and each takes its keys from a probe.Pool. Corrupt is the
// declared error of a corrupted record that the mixin durable of a callable
// names, and nil for a declaration that names none. Clients is the number of
// the clients of the world, and 0 for one client.
type World struct {
	Corrupt    *Expr
	Name       string
	Operations string
	Clients    int
}

// worldOf returns the world of the suite s. The world reads the error of a
// corrupted record from the first callable whose mixin durable names one.
// Each client of the world calls each operation, so the world takes three
// clients only where each operation has the mixin concurrent. The operations
// are the steps of the sequences of s with an action.
func (k *caller) worldOf(s *Suite) *World {
	w := &World{
		Name: k.e.JoinName(worldWord, s.Declaration), Operations: k.e.JoinName(operationsWord, s.Declaration),
	}
	if at := slices.IndexFunc(s.Facts, func(f derive.Facts) bool { return paramOf(&f, corruptParam) != nil }); at >= 0 {
		w.Corrupt = k.valueOfParam(paramOf(&s.Facts[at], corruptParam))
	}
	concurrent := true
	for i := range s.Facts {
		if s.Sequences.Steps[i].Action && !mixed(k.m, &s.Facts[i].ID, shape.Concurrent) {
			concurrent = false
		}
	}
	if concurrent {
		w.Clients = worldClients
	}
	return w
}

// usedWorld returns w where a derived check of derived runs a world probe,
// and nil otherwise.
func usedWorld(w *World, derived []Derived) *World {
	if !slices.ContainsFunc(derived, func(d Derived) bool { return d.World }) {
		return nil
	}
	return w
}
