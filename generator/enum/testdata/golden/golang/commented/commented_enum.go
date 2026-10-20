package commented

import (
	"errors"
	"fmt"
)

// String returns Grade's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Grade) String() string {
	switch v {
	case Poor:
		return "poor"
	case Fair:
		return "fair enough"
	case Good:
		return "Good"
	}
	return fmt.Sprintf("Grade(%d)", int(v))
}

// ErrUnknownGrade reports text that names no declared variant.
var ErrUnknownGrade = errors.New("commented: unknown Grade value")

// ParseGrade maps text back to a Grade.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseGrade(s string) (Grade, error) {
	switch s {
	case "poor":
		return Poor, nil
	case "fair enough":
		return Fair, nil
	case "Good":
		return Good, nil
	}
	var zero Grade
	return zero, ErrUnknownGrade
}

// GradeValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func GradeValues() []Grade {
	return []Grade{
		Poor,
		Fair,
		Good,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Grade) IsValid() bool {
	switch v {
	case Poor:
		return true
	case Fair:
		return true
	case Good:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Grade legal as a map key —
// which a JSON marshaller alone does not.
func (v Grade) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Grade) UnmarshalText(text []byte) error {
	parsed, err := ParseGrade(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns Mood's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Mood) String() string {
	switch v {
	case Calm:
		return "Calm"
	case Angry:
		return "Angry"
	}
	return fmt.Sprintf("Mood(%d)", int(v))
}

// ErrUnknownMood reports text that names no declared variant.
var ErrUnknownMood = errors.New("commented: unknown Mood value")

// ParseMood maps text back to a Mood.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseMood(s string) (Mood, error) {
	switch s {
	case "Calm":
		return Calm, nil
	case "Angry":
		return Angry, nil
	}
	var zero Mood
	return zero, ErrUnknownMood
}

// MoodValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func MoodValues() []Mood {
	return []Mood{
		Calm,
		Angry,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Mood) IsValid() bool {
	switch v {
	case Calm:
		return true
	case Angry:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Mood legal as a map key —
// which a JSON marshaller alone does not.
func (v Mood) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Mood) UnmarshalText(text []byte) error {
	parsed, err := ParseMood(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
