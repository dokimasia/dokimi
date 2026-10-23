package handwritten

import (
	"fmt"
)

// WeekdayValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func WeekdayValues() []Weekday {
	return []Weekday{
		Monday,
		Tuesday,
		Wednesday,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Weekday) IsValid() bool {
	switch v {
	case Monday, Tuesday, Wednesday:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownWeekday] for a value that [Weekday.IsValid]
// refuses, whose text [ParseWeekday] does not parse.
func (v Weekday) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownWeekday, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseWeekday] for a text that does not match a declared
// variant.
func (v *Weekday) UnmarshalText(text []byte) error {
	parsed, err := ParseWeekday(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
