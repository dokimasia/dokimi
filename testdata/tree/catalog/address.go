// Package catalog declares the structs that the composition writes builders
// of, and the errors that it writes checks of.
//
//dokimi:sentinel
package catalog

import "time"

// Address is the postal address of a customer.
//
//dokimi:builder
type Address struct {
	// Street is the street and the house number.
	//
	//dokimi:default "Dam 1"
	Street string
	// City is the city of the address.
	//
	//dokimi:sample value=Utrecht alternate=Leiden
	City string `default:"Amsterdam"`
	// Floor is the floor of the address, where it has one.
	Floor *int
	// Lines are the lines of the address.
	Lines []string
	// Labels are the labels of the address.
	Labels map[string]string
	// Tags are the tags of the address.
	Tags map[string]struct{}
	// Raw is the address as the customer typed it.
	Raw []byte
	// Grid is the position of the address on a map.
	Grid [2]int
	// Timeout bounds a delivery to the address.
	Timeout time.Duration `default:"time.Second"`
	// note is a note of the courier.
	note string
}

// Note returns the note of the courier.
func (a Address) Note() string { return a.note }

// Role is the role of a profile.
type Role struct {
	// Name is the name of the role.
	Name string
}

// Stamp is the time of a change.
type Stamp struct {
	// At is the time of the change.
	At string
}

// Audit records the author of a profile.
type Audit struct {
	*Stamp
	// CreatedBy is the author of the profile.
	CreatedBy string
	// Notes are the notes of the audit.
	Notes []string
}

// Profile is the profile of a customer.
//
//dokimi:builder
type Profile struct {
	Role
	*Audit
	// Active reports whether the profile is in use.
	Active bool
}

// Mirror is the role of a profile under another name.
//
//dokimi:builder
type Mirror = Role
