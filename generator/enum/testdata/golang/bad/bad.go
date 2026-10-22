// Package bad declares the enumerations that enum reports.
package bad

import "context"

// Shade has two variants of one value.
//
//dokimi:enum
type Shade int

// The shades, two of one value.
const (
	Dark  Shade = 1
	Black Shade = 1
)

// Code has two variants of one text, because the text of CodeA is A.
//
//dokimi:enum
type Code int

// The codes, two of one text.
const (
	CodeA Code = 1
	A     Code = 2
)

// Size has two variants whose trailing comments state one text.
//
//dokimi:enum text=comment
type Size int

// The sizes, two of one text.
const (
	Small Size = iota // s
	Large             // s
)

// Tint has a parse function of bytes, whose form conflicts with the parse
// function of the surface.
//
//dokimi:enum
type Tint int

// The tints of a paint.
const (
	Pale Tint = iota + 1
	Deep
)

// ParseTint returns the tint of the bytes b.
func ParseTint(b []byte) (Tint, error) {
	if string(b) == "Deep" {
		return Deep, nil
	}
	return Pale, nil
}

// Hue has a constant under the name of the error of its parse function.
//
//dokimi:enum
type Hue int

// The hues of a light.
const (
	Warm Hue = iota + 1
	Cold
)

// ErrUnknownHue is a text, so it is no error of a parse function.
const ErrUnknownHue = "bad: unknown hue"

// Glow has a validity method that returns a text.
//
//dokimi:enum
type Glow int

// The glows of a lamp.
const (
	Faint Glow = iota + 1
	Bright
)

// IsValid returns the name of the glow.
func (g Glow) IsValid() string {
	if g == Bright {
		return "bright"
	}
	return "faint"
}

// Beam has a render with a parameter, whose form conflicts with the render of
// the surface.
//
//dokimi:enum
type Beam int

// The beams of a headlight.
const (
	Dipped Beam = iota + 1
	Main
)

// String returns the name of the beam, in capitals where loud is true.
func (b Beam) String(loud bool) string {
	if loud {
		return "BEAM"
	}
	return "beam"
}

// Knot has declarations of other types under the names of the members of its
// surface.
//
//dokimi:enum
type Knot int

// The knots of a rope.
const (
	Reef Knot = iota + 1
	Hitch
)

// ParseKnot returns a shade, so it is no parse function of Knot.
func ParseKnot(s string) (Shade, error) { return Dark, nil }

// ErrUnknownKnot is a type, so it is no error of a parse function.
type ErrUnknownKnot struct{}

// KnotValues returns texts, so it is no list of the knots.
func KnotValues() []string { return []string{"Reef", "Hitch"} }

// MarshalText returns a text, so it is no encode of Knot.
func (k Knot) MarshalText() (string, error) { return "knot", nil }

// Loop has declarations of other forms under the names of the members of its
// surface.
//
//dokimi:enum
type Loop int

// The loops of a program.
const (
	For Loop = iota + 1
	While
)

// ParseLoop returns no error, so it is no parse function of Loop.
func ParseLoop(s string) Loop { return For }

// ErrUnknownLoop is a function, so it is no error of a parse function.
func ErrUnknownLoop() {}

// LoopValues returns an array, so it is no list of the loops.
func LoopValues() [2]Loop { return [2]Loop{For, While} }

// Mesh has a parse function with a context and a variable under the name of
// the function that lists its variants.
//
//dokimi:enum
type Mesh int

// The meshes of a sieve.
const (
	Fine Mesh = iota + 1
	Coarse
)

// ParseMesh takes a context, so a check cannot call it with a text alone.
func ParseMesh(ctx context.Context, s string) (Mesh, error) { return Fine, ctx.Err() }

// MeshValues is a variable, so it is no list of the meshes.
var MeshValues = []Mesh{Fine, Coarse}
