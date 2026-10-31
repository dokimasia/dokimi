// Package moved declares the structs whose builders the directive out moves
// into the package movedtest. A builder there cannot refer to an unexported
// declaration of this package.
//
//dokimi:out path=movedtest/
package moved

// config is a configuration that only this package can refer to.
//
//dokimi:builder
type config struct {
	// Name is the name of the configuration.
	Name string
}

// region is a region that only this package can refer to.
type region struct {
	// Name is the name of the region.
	Name string
}

// Site has a member of a type that only this package can refer to.
//
//dokimi:builder
type Site struct {
	// Region has a type that only this package can refer to.
	Region region
}

// zone returns the zone of a new deployment.
func zone() string { return "eu-west-1" }

// Deployment has a default that calls a function that only this package can
// call.
//
//dokimi:builder
type Deployment struct {
	// Zone calls the function zone.
	Zone string `default:"zone()"`
}

// defaultPolicy seeds the builder of Policy.
func defaultPolicy() Policy { return Policy{Attempts: 3} }

// Policy has a seed function that only this package can call.
//
//dokimi:builder defaults=defaultPolicy
type Policy struct {
	// Attempts is the number of attempts.
	Attempts int
}

// Plain has a member that every package can refer to.
//
//dokimi:builder
type Plain struct {
	// Name is the name of the value.
	Name string
}
