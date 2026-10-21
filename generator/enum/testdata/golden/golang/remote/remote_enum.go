package remote

import (
	"errors"
	"fmt"
)

// String returns Season's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Season) String() string {
	switch v {
	case Spring:
		return "Spring"
	case Autumn:
		return "Autumn"
	}
	return fmt.Sprintf("Season(%d)", int(v))
}

// ErrUnknownSeason reports text that names no declared variant.
var ErrUnknownSeason = errors.New("remote: unknown Season value")

// ParseSeason maps text back to a Season.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseSeason(s string) (Season, error) {
	switch s {
	case "Spring":
		return Spring, nil
	case "Autumn":
		return Autumn, nil
	}
	var zero Season
	return zero, ErrUnknownSeason
}

// SeasonValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func SeasonValues() []Season {
	return []Season{
		Spring,
		Autumn,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Season) IsValid() bool {
	switch v {
	case Spring:
		return true
	case Autumn:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Season legal as a map key —
// which a JSON marshaller alone does not.
func (v Season) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Season) UnmarshalText(text []byte) error {
	parsed, err := ParseSeason(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
