package region

import (
	"errors"
	"fmt"
)

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Region) String() string {
	return string(v)
}

// ErrUnknownRegion is the error of a text that does not match a declared
// variant. The error of [ParseRegion] wraps it with the text.
var ErrUnknownRegion = errors.New("region: unknown Region value")

// ParseRegion returns the variant of Region whose text is s. It is the
// inverse of [Region.String] on the declared variants. It returns an error
// that wraps [ErrUnknownRegion] for a text that does not match a declared
// variant.
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
	return zero, fmt.Errorf("%w %q", ErrUnknownRegion, s)
}

// RegionValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func RegionValues() []Region {
	return []Region{
		US,
		EU,
		AP,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Region) IsValid() bool {
	switch v {
	case US, EU, AP:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownRegion] for a value that [Region.IsValid] refuses,
// whose text [ParseRegion] does not parse.
func (v Region) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownRegion, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseRegion] for a text that does not match a declared
// variant.
func (v *Region) UnmarshalText(text []byte) error {
	parsed, err := ParseRegion(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Zone) String() string {
	return string(v)
}

// ErrUnknownZone is the error of a text that does not match a declared
// variant. The error of [ParseZone] wraps it with the text.
var ErrUnknownZone = errors.New("region: unknown Zone value")

// ParseZone returns the variant of Zone whose text is s. It is the inverse
// of [Zone.String] on the declared variants. It returns an error that
// wraps [ErrUnknownZone] for a text that does not match a declared
// variant.
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
	return zero, fmt.Errorf("%w %q", ErrUnknownZone, s)
}

// ZoneValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func ZoneValues() []Zone {
	return []Zone{
		North,
		South,
		Unset,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Zone) IsValid() bool {
	switch v {
	case North, South, Unset:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownZone] for a value that [Zone.IsValid] refuses,
// whose text [ParseZone] does not parse.
func (v Zone) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownZone, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseZone] for a text that does not match a declared variant.
func (v *Zone) UnmarshalText(text []byte) error {
	parsed, err := ParseZone(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Lane) String() string {
	switch v {
	case Left:
		return "left"
	case Right:
		return "right"
	}
	return string(v)
}

// ErrUnknownLane is the error of a text that does not match a declared
// variant. The error of [ParseLane] wraps it with the text.
var ErrUnknownLane = errors.New("region: unknown Lane value")

// ParseLane returns the variant of Lane whose text is s. It is the inverse
// of [Lane.String] on the declared variants. It returns an error that
// wraps [ErrUnknownLane] for a text that does not match a declared
// variant.
func ParseLane(s string) (Lane, error) {
	switch s {
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	}
	var zero Lane
	return zero, fmt.Errorf("%w %q", ErrUnknownLane, s)
}

// LaneValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func LaneValues() []Lane {
	return []Lane{
		Left,
		Right,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Lane) IsValid() bool {
	switch v {
	case Left, Right:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownLane] for a value that [Lane.IsValid] refuses,
// whose text [ParseLane] does not parse.
func (v Lane) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownLane, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseLane] for a text that does not match a declared variant.
func (v *Lane) UnmarshalText(text []byte) error {
	parsed, err := ParseLane(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
