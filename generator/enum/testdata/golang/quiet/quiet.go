// Package quiet declares enumerations whose surface the directive leaves
// out.
package quiet

// Mode is a mode of a lamp with its own render.
//
//dokimi:enum methods=off
type Mode int

// The modes of a lamp.
const (
	Off Mode = iota
	On
)

// String returns the name of the mode.
func (m Mode) String() string {
	if m == On {
		return "on"
	}
	return "off"
}

// IsValid reports whether m is a mode.
func (m Mode) IsValid() bool { return m == Off || m == On }

// MarshalText writes the name of the mode.
func (m Mode) MarshalText() ([]byte, error) { return []byte(m.String()), nil }

// Plain is a mode of a switch without a member of the surface, whose
// surface the directive leaves out.
//
//dokimi:enum methods=off
type Plain int

// The modes of a switch.
const (
	Loose Plain = iota + 1
	Tight
)

// Tone is a tone of a bell with its own encode.
//
//dokimi:enum
type Tone int

// The tones of a bell.
const (
	Low Tone = iota + 1
	High
)

// MarshalText writes the text of the tone.
func (t Tone) MarshalText() ([]byte, error) { return []byte(t.String()), nil }

// Switch is a switch whose variants take both truth values, so no value
// lies outside them.
//
//dokimi:enum
type Switch bool

// The positions of a switch.
const (
	Up   Switch = true
	Down Switch = false
)

// Phase is a phase of a lamp with its own render and no parse function.
//
//dokimi:enum
type Phase int

// The phases of a lamp.
const (
	Warm Phase = iota + 1
	Cool
)

// String returns the name of the phase.
func (p Phase) String() string {
	if p == Cool {
		return "cool"
	}
	return "warm"
}

// Speed is a speed of a fan with its own decode.
//
//dokimi:enum
type Speed int

// The speeds of a fan.
const (
	Slow Speed = iota + 1
	Fast
)

// UnmarshalText reads a speed from its digit.
func (s *Speed) UnmarshalText(text []byte) error {
	if string(text) == "2" {
		*s = Fast
		return nil
	}
	*s = Slow
	return nil
}

// Level is a level of a lamp, a float.
//
//dokimi:enum
type Level float64

// The levels of a lamp.
const (
	Dim    Level = 0.5
	Bright Level = 1.5
)

// Flag is a flag of a lamp, a truth value.
//
//dokimi:enum
type Flag bool

// The flags of a lamp.
const (
	Lit Flag = true
)
