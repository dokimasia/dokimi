// Package linkssim contains a system of components whose directive tolerates
// the faults of the links between its hosts.
package linkssim

import (
	"go.dokimi.dev/sim"

	"example.com/shop/cells"
)

// linked spawns two cells.
//
//dokimi:system tolerate=[reset:1, slow-link:1, disk-replacement:1]
func linked(w *sim.World[sim.None, sim.None]) {
	sim.SpawnComponent(w, "a", cells.Open)
	sim.SpawnComponent(w, "b", cells.Open)
}
