package declared

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Debug-1]
	_ = x[Info-2]
	_ = x[Warn-3]
}

const _Level_name = "DebugInfoWarn"

var _Level_index = [...]uint8{0, 5, 9, 13}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Level) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Level_index)-1 {
		return "Level(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Level_name[_Level_index[idx]:_Level_index[idx+1]]
}

// ParseLevel returns the variant of Level whose text is s. It is the
// inverse of [Level.String] on the declared variants. It returns an error
// that wraps [ErrUnknownLevel] for a text that does not match a declared
// variant.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "Debug":
		return Debug, nil
	case "Info":
		return Info, nil
	case "Warn":
		return Warn, nil
	}
	var zero Level
	return zero, fmt.Errorf("%w %q", ErrUnknownLevel, s)
}

// LevelValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func LevelValues() []Level {
	return []Level{
		Debug,
		Info,
		Warn,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Level) IsValid() bool {
	switch v {
	case Debug, Info, Warn:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Level) MarshalText() ([]byte, error) {
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

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Red-1]
	_ = x[Green-2]
}

const _Color_name = "RedGreen"

var _Color_index = [...]uint8{0, 3, 8}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Color) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Color_index)-1 {
		return "Color(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Color_name[_Color_index[idx]:_Color_index[idx+1]]
}

// ErrUnknownColor is the error of a text that does not match a declared
// variant. The error of [ParseColor] wraps it with the text.
var ErrUnknownColor = errors.New("declared: unknown Color value")

// ParseColor returns the variant of Color whose text is s. It is the
// inverse of [Color.String] on the declared variants. It returns an error
// that wraps [ErrUnknownColor] for a text that does not match a declared
// variant.
func ParseColor(s string) (Color, error) {
	switch s {
	case "Red":
		return Red, nil
	case "Green":
		return Green, nil
	}
	var zero Color
	return zero, fmt.Errorf("%w %q", ErrUnknownColor, s)
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Color) IsValid() bool {
	switch v {
	case Red, Green:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Color) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseColor] for a text that does not match a declared variant.
func (v *Color) UnmarshalText(text []byte) error {
	parsed, err := ParseColor(string(text))
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
	_ = x[Left-1]
	_ = x[Right-2]
}

const _Lane_name = "LeftRight"

var _Lane_index = [...]uint8{0, 4, 9}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Lane) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Lane_index)-1 {
		return "Lane(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Lane_name[_Lane_index[idx]:_Lane_index[idx+1]]
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
// encoding/json uses it for a value and for a map key.
func (v Lane) MarshalText() ([]byte, error) {
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

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Snug-1]
	_ = x[Slack-2]
}

const _Coil_name = "SnugSlack"

var _Coil_index = [...]uint8{0, 4, 9}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Coil) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Coil_index)-1 {
		return "Coil(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Coil_name[_Coil_index[idx]:_Coil_index[idx+1]]
}

// ErrUnknownCoil is the error of a text that does not match a declared
// variant. The error of [ParseCoil] wraps it with the text.
var ErrUnknownCoil = errors.New("declared: unknown Coil value")

// ParseCoil returns the variant of Coil whose text is s. It is the inverse
// of [Coil.String] on the declared variants. It returns an error that
// wraps [ErrUnknownCoil] for a text that does not match a declared
// variant.
func ParseCoil(s string) (Coil, error) {
	switch s {
	case "Snug":
		return Snug, nil
	case "Slack":
		return Slack, nil
	}
	var zero Coil
	return zero, fmt.Errorf("%w %q", ErrUnknownCoil, s)
}

// CoilValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func CoilValues() []Coil {
	return []Coil{
		Snug,
		Slack,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Coil) IsValid() bool {
	switch v {
	case Snug, Slack:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Coil) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseCoil] for a text that does not match a declared variant.
func (v *Coil) UnmarshalText(text []byte) error {
	parsed, err := ParseCoil(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
