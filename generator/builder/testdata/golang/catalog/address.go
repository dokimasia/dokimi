// Package catalog declares the structs that builder reads in the specs.
package catalog

import "time"

// Address has a field of each form of a setter.
//
//dokimi:builder
type Address struct {
	// Street has a default through the directive.
	//
	//dokimi:default "Dam 1"
	Street string
	// City has a default through the tag.
	City string `default:"Amsterdam"`
	// Floor is optional.
	Floor *int
	// Lines is a sequence.
	Lines []string
	// Labels is a mapping.
	Labels map[string]string
	// Tags is a set.
	Tags map[string]struct{}
	// Raw is bytes.
	Raw []byte
	// Grid has a fixed length.
	Grid [2]int
	// Timeout has a default that names a constant of another package.
	Timeout time.Duration `default:"time.Second"`
	// note is unexported, so it has no setter.
	note string
}

// Role is embedded by value in Profile.
type Role struct {
	// Name is promoted to Profile, and has a default that Profile does not
	// assign.
	//
	//dokimi:default "admin"
	Name string
}

// Stamp is embedded by pointer in Audit.
type Stamp struct {
	// At is promoted to Profile through two pointers.
	At string
}

// Audit is embedded by pointer in Profile.
type Audit struct {
	*Stamp
	// CreatedBy is promoted to Profile through the pointer.
	CreatedBy string
	// Notes is a sequence that the pointer promotes to Profile.
	Notes []string
}

// tag is embedded in Profile under an unexported name.
type tag struct {
	// Label is promoted to Profile through the unexported embed.
	Label string
}

// Profile embeds a struct by value, a struct by pointer and a struct with an
// unexported name.
//
//dokimi:builder
type Profile struct {
	Role
	*Audit
	tag
	// Active is declared on Profile, and its default is the zero value.
	//
	//dokimi:default false
	Active bool
}

// Notifier is an interface that Envelope embeds.
type Notifier interface {
	// Notify sends a notice.
	Notify() error
}

// listener is an interface with an unexported name that Envelope embeds.
type listener interface {
	// Listen waits for a notice.
	Listen()
}

// marks is a defined type with an unexported name that Envelope embeds.
type marks []string

// Envelope embeds an interface, a defined type and an interface of the
// language, and an interface and a defined type with unexported names.
//
//dokimi:builder
type Envelope struct {
	Notifier
	Codes
	error
	listener
	marks
}

// Hidden has a field that the tag builder leaves out.
//
//dokimi:builder
type Hidden struct {
	// Name has a setter.
	Name string
	// Secret has the tag that leaves it out.
	Secret string `builder:"-"`
}

// Remarks is a transparent alias of a sequence.
type Remarks = []string

// Codes is a defined type of a sequence.
type Codes []string

// Order has fields whose types an alias, a defined type and a spelling with
// a space give.
//
//dokimi:builder
type Order struct {
	// Remarks is a sequence through an alias.
	Remarks Remarks
	// Codes is a reference to a defined type.
	Codes Codes
	// Flags is a set whose empty struct has a space.
	Flags map[string]struct {}
	// Marker is a scalar of one value, so it has a sample and no alternate.
	Marker struct{}
	// Marks is a sequence whose elements have one value, so it has a sample
	// and no alternate.
	Marks []struct{}
}

// Loose has a field of each form whose type has no sample. eidos derives no
// sample of an interface with a method, such as Notifier.
//
//dokimi:builder
type Loose struct {
	// Items is a sequence of notifiers.
	Items []Notifier
	// Index is a mapping to notifiers.
	Index map[string]Notifier
	// Keys is a set of notifiers.
	Keys map[Notifier]struct{}
	// Extra is optional.
	Extra *Notifier
	// Hook is a notifier.
	Hook Notifier
}
