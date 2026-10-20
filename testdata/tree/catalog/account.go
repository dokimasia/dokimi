package catalog

// Account is the account of a customer.
//
//dokimi:builder
type Account struct {
	// Owner is the owner of the account.
	//
	//dokimi:default "owner"
	Owner string
	// Plan is the plan of the account.
	Plan string
}

// AccountDefaults returns the account of a new customer.
func AccountDefaults() Account {
	return Account{Owner: "seed", Plan: "basic"}
}

// Policy is the retry policy of a delivery.
//
//dokimi:builder defaults=defaultPolicy
type Policy struct {
	// Attempts is the number of attempts of a delivery.
	Attempts int
}

// defaultPolicy returns the policy of a new delivery.
func defaultPolicy() Policy {
	return Policy{Attempts: 3}
}

// Tag is a tag whose key is a kind of string. The rules derive no witness
// of its type parameter, so the directive witness states one.
//
//dokimi:builder
//dokimi:witness K=string
type Tag[K ~string] struct {
	// Key is the key of the tag.
	Key K
}

// Box holds one value of any type.
//
//dokimi:builder
type Box[T any] struct {
	// Value is the value in the box.
	Value T
	// Size is the size of the box.
	//
	//dokimi:default 1
	Size int
}
