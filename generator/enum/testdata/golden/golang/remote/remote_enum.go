package remote

import (
	"errors"
	"fmt"
)

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Season) String() string {
	switch v {
	case Spring:
		return "Spring"
	case Autumn:
		return "Autumn"
	}
	return fmt.Sprintf("Season(%d)", int(v))
}

// ErrUnknownSeason is the error of a text that does not match a declared
// variant. The error of [ParseSeason] wraps it with the text.
var ErrUnknownSeason = errors.New("remote: unknown Season value")

// ParseSeason returns the variant of Season whose text is s. It is the
// inverse of [Season.String] on the declared variants. It returns an error
// that wraps [ErrUnknownSeason] for a text that does not match a declared
// variant.
func ParseSeason(s string) (Season, error) {
	switch s {
	case "Spring":
		return Spring, nil
	case "Autumn":
		return Autumn, nil
	}
	var zero Season
	return zero, fmt.Errorf("%w %q", ErrUnknownSeason, s)
}

// SeasonValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func SeasonValues() []Season {
	return []Season{
		Spring,
		Autumn,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Season) IsValid() bool {
	switch v {
	case Spring, Autumn:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownSeason] for a value that [Season.IsValid] refuses,
// whose text [ParseSeason] does not parse.
func (v Season) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownSeason, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseSeason] for a text that does not match a declared
// variant.
func (v *Season) UnmarshalText(text []byte) error {
	parsed, err := ParseSeason(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
