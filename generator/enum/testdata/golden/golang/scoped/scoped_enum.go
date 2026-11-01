package scoped

import (
	errors2 "errors"
	fmt2 "fmt"
	strconv2 "strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Draft-1]
	_ = x[Published-2]
}

const _Status_name = "DraftPublished"

var _Status_index = [...]uint8{0, 5, 14}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Status) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Status_index)-1 {
		return "Status(" + strconv2.FormatInt(int64(v), 10) + ")"
	}
	return _Status_name[_Status_index[idx]:_Status_index[idx+1]]
}

// ErrUnknownStatus is the error of a text that does not match a declared
// variant. The error of [ParseStatus] wraps it with the text.
var ErrUnknownStatus = errors2.New("scoped: unknown Status value")

// ParseStatus returns the variant of Status whose text is s. It is the
// inverse of [Status.String] on the declared variants. It returns an error
// that wraps [ErrUnknownStatus] for a text that does not match a declared
// variant.
func ParseStatus(s string) (Status, error) {
	switch s {
	case "Draft":
		return Draft, nil
	case "Published":
		return Published, nil
	}
	var zero Status
	return zero, fmt2.Errorf("%w %q", ErrUnknownStatus, s)
}

// StatusValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func StatusValues() []Status {
	return []Status{
		Draft,
		Published,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Status) IsValid() bool {
	switch v {
	case Draft, Published:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownStatus] for a value that [Status.IsValid] refuses,
// whose text [ParseStatus] does not parse.
func (v Status) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt2.Errorf("%w %q", ErrUnknownStatus, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseStatus] for a text that does not match a declared
// variant.
func (v *Status) UnmarshalText(text []byte) error {
	parsed, err := ParseStatus(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Level) String() string {
	switch v {
	case Dim:
		return "Dim"
	case Bright:
		return "Bright"
	}
	return fmt2.Sprintf("Level(%g)", float64(v))
}

// ErrUnknownLevel is the error of a text that does not match a declared
// variant. The error of [ParseLevel] wraps it with the text.
var ErrUnknownLevel = errors2.New("scoped: unknown Level value")

// ParseLevel returns the variant of Level whose text is s. It is the
// inverse of [Level.String] on the declared variants. It returns an error
// that wraps [ErrUnknownLevel] for a text that does not match a declared
// variant.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "Dim":
		return Dim, nil
	case "Bright":
		return Bright, nil
	}
	var zero Level
	return zero, fmt2.Errorf("%w %q", ErrUnknownLevel, s)
}

// LevelValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func LevelValues() []Level {
	return []Level{
		Dim,
		Bright,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Level) IsValid() bool {
	switch v {
	case Dim, Bright:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownLevel] for a value that [Level.IsValid] refuses,
// whose text [ParseLevel] does not parse.
func (v Level) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt2.Errorf("%w %q", ErrUnknownLevel, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseLevel] for a text that does not match a declared variant.
func (v *Level) UnmarshalText(text []byte) error {
	parsed, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
