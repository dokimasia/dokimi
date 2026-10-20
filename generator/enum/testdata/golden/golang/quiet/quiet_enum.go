package quiet

import (
	"errors"
	"fmt"
)

// String returns Tone's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Tone) String() string {
	switch v {
	case Low:
		return "Low"
	case High:
		return "High"
	}
	return fmt.Sprintf("Tone(%d)", int(v))
}

// ErrUnknownTone reports text that names no declared variant.
var ErrUnknownTone = errors.New("quiet: unknown Tone value")

// ParseTone maps text back to a Tone.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseTone(s string) (Tone, error) {
	switch s {
	case "Low":
		return Low, nil
	case "High":
		return High, nil
	}
	var zero Tone
	return zero, ErrUnknownTone
}

// ToneValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func ToneValues() []Tone {
	return []Tone{
		Low,
		High,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Tone) IsValid() bool {
	switch v {
	case Low:
		return true
	case High:
		return true
	}
	return false
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Tone) UnmarshalText(text []byte) error {
	parsed, err := ParseTone(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns Switch's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Switch) String() string {
	switch v {
	case Up:
		return "Up"
	case Down:
		return "Down"
	}
	return fmt.Sprintf("Switch(%t)", bool(v))
}

// ErrUnknownSwitch reports text that names no declared variant.
var ErrUnknownSwitch = errors.New("quiet: unknown Switch value")

// ParseSwitch maps text back to a Switch.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseSwitch(s string) (Switch, error) {
	switch s {
	case "Up":
		return Up, nil
	case "Down":
		return Down, nil
	}
	var zero Switch
	return zero, ErrUnknownSwitch
}

// SwitchValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func SwitchValues() []Switch {
	return []Switch{
		Up,
		Down,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Switch) IsValid() bool {
	switch v {
	case Up:
		return true
	case Down:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Switch legal as a map key —
// which a JSON marshaller alone does not.
func (v Switch) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Switch) UnmarshalText(text []byte) error {
	parsed, err := ParseSwitch(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Phase) IsValid() bool {
	switch v {
	case Warm:
		return true
	case Cool:
		return true
	}
	return false
}

// String returns Speed's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Speed) String() string {
	switch v {
	case Slow:
		return "Slow"
	case Fast:
		return "Fast"
	}
	return fmt.Sprintf("Speed(%d)", int(v))
}

// ErrUnknownSpeed reports text that names no declared variant.
var ErrUnknownSpeed = errors.New("quiet: unknown Speed value")

// ParseSpeed maps text back to a Speed.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseSpeed(s string) (Speed, error) {
	switch s {
	case "Slow":
		return Slow, nil
	case "Fast":
		return Fast, nil
	}
	var zero Speed
	return zero, ErrUnknownSpeed
}

// SpeedValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func SpeedValues() []Speed {
	return []Speed{
		Slow,
		Fast,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Speed) IsValid() bool {
	switch v {
	case Slow:
		return true
	case Fast:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Speed legal as a map key —
// which a JSON marshaller alone does not.
func (v Speed) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// String returns Level's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Level) String() string {
	switch v {
	case Dim:
		return "Dim"
	case Bright:
		return "Bright"
	}
	return fmt.Sprintf("Level(%g)", float64(v))
}

// ErrUnknownLevel reports text that names no declared variant.
var ErrUnknownLevel = errors.New("quiet: unknown Level value")

// ParseLevel maps text back to a Level.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "Dim":
		return Dim, nil
	case "Bright":
		return Bright, nil
	}
	var zero Level
	return zero, ErrUnknownLevel
}

// LevelValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func LevelValues() []Level {
	return []Level{
		Dim,
		Bright,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Level) IsValid() bool {
	switch v {
	case Dim:
		return true
	case Bright:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Level legal as a map key —
// which a JSON marshaller alone does not.
func (v Level) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Level) UnmarshalText(text []byte) error {
	parsed, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns Flag's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Flag) String() string {
	switch v {
	case Lit:
		return "Lit"
	}
	return fmt.Sprintf("Flag(%t)", bool(v))
}

// ErrUnknownFlag reports text that names no declared variant.
var ErrUnknownFlag = errors.New("quiet: unknown Flag value")

// ParseFlag maps text back to a Flag.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseFlag(s string) (Flag, error) {
	switch s {
	case "Lit":
		return Lit, nil
	}
	var zero Flag
	return zero, ErrUnknownFlag
}

// FlagValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func FlagValues() []Flag {
	return []Flag{
		Lit,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Flag) IsValid() bool {
	switch v {
	case Lit:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Flag legal as a map key —
// which a JSON marshaller alone does not.
func (v Flag) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Flag) UnmarshalText(text []byte) error {
	parsed, err := ParseFlag(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
