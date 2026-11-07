// Package routedsim contains a system of components whose simulation the
// directive out routes to another package.
package routedsim

import "go.dokimi.dev/sim"

// routed is a system of components whose simulation is in another package.
//
//dokimi:system
//dokimi:out path=gen/
func routed(*sim.World[sim.None, sim.None]) {}
