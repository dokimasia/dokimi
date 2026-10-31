// Package hidden declares the builders whose samples, members and witnesses
// only this package can refer to.
package hidden

// Host is a host whose samples are the results of unexported functions.
//
//dokimi:sample value=newHost() alternate=otherHost()
type Host struct {
	// Name is the name of the host.
	Name string
}

// newHost returns the sample of Host.
func newHost() Host { return Host{Name: "eu"} }

// otherHost returns the alternate sample of Host.
func otherHost() Host { return Host{Name: "us"} }

// Zone is a zone whose sample is the result of an exported function and whose
// alternate is the result of an unexported function.
//
//dokimi:sample value=NewZone() alternate=otherZone()
type Zone struct {
	// Name is the name of the zone.
	Name string
}

// NewZone returns the sample of Zone.
func NewZone() Zone { return Zone{Name: "eu-west-1"} }

// otherZone returns the alternate sample of Zone.
func otherZone() Zone { return Zone{Name: "us-east-1"} }

// region is a region that only this package can refer to.
type region struct {
	// Name is the name of the region.
	Name string
}

// Site has a member whose samples only this package can refer to, a member
// whose alternate only this package can refer to, and a member of a type that
// only this package can refer to.
//
//dokimi:builder
type Site struct {
	// Host has the samples of Host.
	Host Host
	// Zone has the samples of Zone.
	Zone Zone
	// Region has a type that only this package can refer to.
	Region region
	// Hosts are the hosts of the site.
	Hosts []Host
}

// slot is a slot that only this package can refer to.
type slot struct {
	// N is the number of the slot.
	N int
}

// Pool is a generic struct whose witness only this package can refer to.
//
//dokimi:builder
//dokimi:witness T=slot
type Pool[T any] struct {
	// Items are the items of the pool.
	Items []T
}

// config is a configuration that only this package can refer to.
//
//dokimi:builder
type config struct {
	// Name is the name of the configuration.
	Name string `default:"web"`
}
