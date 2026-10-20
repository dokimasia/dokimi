// Package seeded declares the structs whose builders a function seeds.
package seeded

import "runtime"

// Account has the companion AccountDefaults.
//
//dokimi:builder
type Account struct {
	// Owner has a default through the directive, which New assigns over the
	// value of the companion.
	//
	//dokimi:default "owner"
	Owner string
	// Plan is set by the companion.
	Plan string
}

// AccountDefaults is the companion of Account.
func AccountDefaults() Account {
	return Account{Owner: "seed", Plan: "basic"}
}

// Policy names the function that seeds its builder.
//
//dokimi:builder defaults=defaultPolicy
type Policy struct {
	// Attempts is set by the function of the directive.
	Attempts int
}

// defaultPolicy seeds the builder of Policy.
func defaultPolicy() Policy {
	return Policy{Attempts: 3}
}

// Tier has the function TierSeed, which seeds its builder where the option
// companion-word is seed.
//
//dokimi:builder
type Tier struct {
	// Level is set by TierSeed.
	Level int
}

// TierSeed seeds the builder of Tier where the option companion-word is
// seed.
func TierSeed() Tier {
	return Tier{Level: 2}
}

// Region is the region of a new deployment.
const Region = "eu-west"

// zone returns the zone of a new deployment.
func zone() string { return "eu-west-1" }

// Code is a defined type of text.
type Code string

// Deployment has defaults of text that name a declaration, and defaults of
// text that are the text itself.
//
//dokimi:builder
type Deployment struct {
	// Mode is a defined type of text whose default is a truth value, which
	// is no string.
	Mode Code `default:"true"`
	// Region names the constant Region.
	Region string `default:"Region"`
	// Zone calls the function zone.
	Zone string `default:"zone()"`
	// System names the constant GOOS of the package runtime.
	System string `default:"runtime.GOOS"`
	// Name names no declaration, so its default is the text web.
	Name string `default:"web"`
	// Label is a quoted string.
	Label string `default:"\"a b\""`
	// Flag is a truth value, which is no string.
	Flag string `default:"true"`
}

// system is the system of a new deployment, the code that uses the import
// of runtime.
var system = runtime.GOOS
