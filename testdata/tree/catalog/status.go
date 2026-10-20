package catalog

// Status is the state of an account.
//
//dokimi:enum
type Status int

// The states of an account.
const (
	Active Status = iota + 1
	Suspended
	Closed
)
