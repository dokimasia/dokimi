package commented

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Poor-0]
	_ = x[Fair-1]
	_ = x[Good-2]
}

const _Grade_name = "poorfair enoughGood"

var _Grade_index = [...]uint8{0, 4, 15, 19}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Grade) String() string {
	idx := int(v) - 0
	if v < 0 || idx >= len(_Grade_index)-1 {
		return "Grade(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Grade_name[_Grade_index[idx]:_Grade_index[idx+1]]
}

// ErrUnknownGrade is the error of a text that does not match a declared
// variant. The error of [ParseGrade] wraps it with the text.
var ErrUnknownGrade = errors.New("commented: unknown Grade value")

// ParseGrade returns the variant of Grade whose text is s. It is the
// inverse of [Grade.String] on the declared variants. It returns an error
// that wraps [ErrUnknownGrade] for a text that does not match a declared
// variant.
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
	return zero, fmt.Errorf("%w %q", ErrUnknownGrade, s)
}

// GradeValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func GradeValues() []Grade {
	return []Grade{
		Poor,
		Fair,
		Good,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Grade) IsValid() bool {
	switch v {
	case Poor, Fair, Good:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownGrade] for a value that [Grade.IsValid] refuses,
// whose text [ParseGrade] does not parse.
func (v Grade) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownGrade, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseGrade] for a text that does not match a declared variant.
func (v *Grade) UnmarshalText(text []byte) error {
	parsed, err := ParseGrade(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Calm-0]
	_ = x[Angry-1]
}

const _Mood_name = "CalmAngry"

var _Mood_index = [...]uint8{0, 4, 9}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Mood) String() string {
	idx := int(v) - 0
	if v < 0 || idx >= len(_Mood_index)-1 {
		return "Mood(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Mood_name[_Mood_index[idx]:_Mood_index[idx+1]]
}

// ErrUnknownMood is the error of a text that does not match a declared
// variant. The error of [ParseMood] wraps it with the text.
var ErrUnknownMood = errors.New("commented: unknown Mood value")

// ParseMood returns the variant of Mood whose text is s. It is the inverse
// of [Mood.String] on the declared variants. It returns an error that
// wraps [ErrUnknownMood] for a text that does not match a declared
// variant.
func ParseMood(s string) (Mood, error) {
	switch s {
	case "Calm":
		return Calm, nil
	case "Angry":
		return Angry, nil
	}
	var zero Mood
	return zero, fmt.Errorf("%w %q", ErrUnknownMood, s)
}

// MoodValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func MoodValues() []Mood {
	return []Mood{
		Calm,
		Angry,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Mood) IsValid() bool {
	switch v {
	case Calm, Angry:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownMood] for a value that [Mood.IsValid] refuses,
// whose text [ParseMood] does not parse.
func (v Mood) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownMood, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseMood] for a text that does not match a declared variant.
func (v *Mood) UnmarshalText(text []byte) error {
	parsed, err := ParseMood(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
