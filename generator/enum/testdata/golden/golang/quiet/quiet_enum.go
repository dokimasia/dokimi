package quiet

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Low-1]
	_ = x[High-2]
}

const _Tone_name = "LowHigh"

var _Tone_index = [...]uint8{0, 3, 7}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Tone) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Tone_index)-1 {
		return "Tone(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Tone_name[_Tone_index[idx]:_Tone_index[idx+1]]
}

// ErrUnknownTone is the error of text that names no declared variant. The
// error of [ParseTone] wraps it with the text.
var ErrUnknownTone = errors.New("quiet: unknown Tone value")

// ParseTone returns the variant of Tone whose text is s. It is the inverse
// of [Tone.String] on the declared variants. It returns an error that
// wraps [ErrUnknownTone] for text that names no declared variant.
func ParseTone(s string) (Tone, error) {
	switch s {
	case "Low":
		return Low, nil
	case "High":
		return High, nil
	}
	var zero Tone
	return zero, fmt.Errorf("%w %q", ErrUnknownTone, s)
}

// ToneValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func ToneValues() []Tone {
	return []Tone{
		Low,
		High,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Tone) IsValid() bool {
	switch v {
	case Low, High:
		return true
	default:
		return false
	}
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseTone] for text that names no declared variant.
func (v *Tone) UnmarshalText(text []byte) error {
	parsed, err := ParseTone(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Switch) String() string {
	switch v {
	case Up:
		return "Up"
	case Down:
		return "Down"
	}
	return fmt.Sprintf("Switch(%t)", bool(v))
}

// ErrUnknownSwitch is the error of text that names no declared variant.
// The error of [ParseSwitch] wraps it with the text.
var ErrUnknownSwitch = errors.New("quiet: unknown Switch value")

// ParseSwitch returns the variant of Switch whose text is s. It is the
// inverse of [Switch.String] on the declared variants. It returns an error
// that wraps [ErrUnknownSwitch] for text that names no declared variant.
func ParseSwitch(s string) (Switch, error) {
	switch s {
	case "Up":
		return Up, nil
	case "Down":
		return Down, nil
	}
	var zero Switch
	return zero, fmt.Errorf("%w %q", ErrUnknownSwitch, s)
}

// SwitchValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func SwitchValues() []Switch {
	return []Switch{
		Up,
		Down,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Switch) IsValid() bool {
	switch v {
	case Up, Down:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Switch) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseSwitch] for text that names no declared variant.
func (v *Switch) UnmarshalText(text []byte) error {
	parsed, err := ParseSwitch(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// PhaseValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func PhaseValues() []Phase {
	return []Phase{
		Warm,
		Cool,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Phase) IsValid() bool {
	switch v {
	case Warm, Cool:
		return true
	default:
		return false
	}
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Slow-1]
	_ = x[Fast-2]
}

const _Speed_name = "SlowFast"

var _Speed_index = [...]uint8{0, 4, 8}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Speed) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Speed_index)-1 {
		return "Speed(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Speed_name[_Speed_index[idx]:_Speed_index[idx+1]]
}

// ErrUnknownSpeed is the error of text that names no declared variant. The
// error of [ParseSpeed] wraps it with the text.
var ErrUnknownSpeed = errors.New("quiet: unknown Speed value")

// ParseSpeed returns the variant of Speed whose text is s. It is the
// inverse of [Speed.String] on the declared variants. It returns an error
// that wraps [ErrUnknownSpeed] for text that names no declared variant.
func ParseSpeed(s string) (Speed, error) {
	switch s {
	case "Slow":
		return Slow, nil
	case "Fast":
		return Fast, nil
	}
	var zero Speed
	return zero, fmt.Errorf("%w %q", ErrUnknownSpeed, s)
}

// SpeedValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func SpeedValues() []Speed {
	return []Speed{
		Slow,
		Fast,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Speed) IsValid() bool {
	switch v {
	case Slow, Fast:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Speed) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
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
	return fmt.Sprintf("Level(%g)", float64(v))
}

// ErrUnknownLevel is the error of text that names no declared variant. The
// error of [ParseLevel] wraps it with the text.
var ErrUnknownLevel = errors.New("quiet: unknown Level value")

// ParseLevel returns the variant of Level whose text is s. It is the
// inverse of [Level.String] on the declared variants. It returns an error
// that wraps [ErrUnknownLevel] for text that names no declared variant.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "Dim":
		return Dim, nil
	case "Bright":
		return Bright, nil
	}
	var zero Level
	return zero, fmt.Errorf("%w %q", ErrUnknownLevel, s)
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
// encoding/json uses it for a value and for a map key.
func (v Level) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseLevel] for text that names no declared variant.
func (v *Level) UnmarshalText(text []byte) error {
	parsed, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Flag) String() string {
	switch v {
	case Lit:
		return "Lit"
	}
	return fmt.Sprintf("Flag(%t)", bool(v))
}

// ErrUnknownFlag is the error of text that names no declared variant. The
// error of [ParseFlag] wraps it with the text.
var ErrUnknownFlag = errors.New("quiet: unknown Flag value")

// ParseFlag returns the variant of Flag whose text is s. It is the inverse
// of [Flag.String] on the declared variants. It returns an error that
// wraps [ErrUnknownFlag] for text that names no declared variant.
func ParseFlag(s string) (Flag, error) {
	switch s {
	case "Lit":
		return Lit, nil
	}
	var zero Flag
	return zero, fmt.Errorf("%w %q", ErrUnknownFlag, s)
}

// FlagValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func FlagValues() []Flag {
	return []Flag{
		Lit,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Flag) IsValid() bool {
	switch v {
	case Lit:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Flag) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseFlag] for text that names no declared variant.
func (v *Flag) UnmarshalText(text []byte) error {
	parsed, err := ParseFlag(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
