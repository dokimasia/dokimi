package region

import (
	"errors"
)

// String returns Region's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Region) String() string {
	return string(v)
}

// ErrUnknownRegion reports text that names no declared variant.
var ErrUnknownRegion = errors.New("region: unknown Region value")

// ParseRegion maps text back to a Region.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseRegion(s string) (Region, error) {
	switch s {
	case "us-east":
		return US, nil
	case "eu-west":
		return EU, nil
	case "ap-south":
		return AP, nil
	}
	var zero Region
	return zero, ErrUnknownRegion
}

// RegionValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func RegionValues() []Region {
	return []Region{
		US,
		EU,
		AP,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Region) IsValid() bool {
	switch v {
	case US:
		return true
	case EU:
		return true
	case AP:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Region legal as a map key —
// which a JSON marshaller alone does not.
func (v Region) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Region) UnmarshalText(text []byte) error {
	parsed, err := ParseRegion(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns Zone's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Zone) String() string {
	return string(v)
}

// ErrUnknownZone reports text that names no declared variant.
var ErrUnknownZone = errors.New("region: unknown Zone value")

// ParseZone maps text back to a Zone.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseZone(s string) (Zone, error) {
	switch s {
	case "north":
		return North, nil
	case "south":
		return South, nil
	case "":
		return Unset, nil
	}
	var zero Zone
	return zero, ErrUnknownZone
}

// ZoneValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func ZoneValues() []Zone {
	return []Zone{
		North,
		South,
		Unset,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Zone) IsValid() bool {
	switch v {
	case North:
		return true
	case South:
		return true
	case Unset:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Zone legal as a map key —
// which a JSON marshaller alone does not.
func (v Zone) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Zone) UnmarshalText(text []byte) error {
	parsed, err := ParseZone(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns Lane's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Lane) String() string {
	switch v {
	case Left:
		return "left"
	case Right:
		return "right"
	}
	return string(v)
}

// ErrUnknownLane reports text that names no declared variant.
var ErrUnknownLane = errors.New("region: unknown Lane value")

// ParseLane maps text back to a Lane.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseLane(s string) (Lane, error) {
	switch s {
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	}
	var zero Lane
	return zero, ErrUnknownLane
}

// LaneValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func LaneValues() []Lane {
	return []Lane{
		Left,
		Right,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Lane) IsValid() bool {
	switch v {
	case Left:
		return true
	case Right:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Lane legal as a map key —
// which a JSON marshaller alone does not.
func (v Lane) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Lane) UnmarshalText(text []byte) error {
	parsed, err := ParseLane(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
