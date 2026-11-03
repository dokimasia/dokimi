package validated

// Pitch is a pitch that the kanon directive in the doc of an alias names.
//
//dokimi:enum
type Pitch uint8

// The pitches.
const (
	Low Pitch = iota + 1
	High
)

// Tempo is a tempo that the kanon directive in the doc of a constant names.
//
//dokimi:enum
type Tempo uint8

// The tempos.
const (
	Slow Tempo = iota + 1
	Fast
)

// Meter is a meter that the kanon directive in the doc of a variable names.
//
//dokimi:enum
type Meter uint8

// The meters.
const (
	Duple Meter = iota + 1
	Triple
)

// Chord is a chord that the kanon directive in the doc of an interface names.
//
//dokimi:enum
type Chord uint8

// The chords.
const (
	Major Chord = iota + 1
	Minor
)

// Scale is a scale that the kanon directive in the doc of a function names.
//
//dokimi:enum
type Scale uint8

// The scales.
const (
	Diatonic Scale = iota + 1
	Chromatic
)

// Voice is a voice that the kanon directive in the doc of a struct names.
//
//dokimi:enum
type Voice uint8

// The voices.
const (
	Soprano Voice = iota + 1
	Alto
)

// Clef is a clef that the kanon directive in the doc of a method names.
//
//dokimi:enum
type Clef uint8

// The clefs.
const (
	Treble Clef = iota + 1
	Bass
)

// Shell is a shell that a directive of the executable of kanon on Windows
// names.
//
//go:generate C:\tools\kanon.exe --type=Shell -canonical
//dokimi:enum
type Shell uint8

// The shells.
const (
	Bourne Shell = iota + 1
	Korn
)

// Shade is a shade that a directive of a command in a directory names.
//
//go:generate go run ./tools/ -type=Shade
//dokimi:enum
type Shade uint8

// The shades.
const (
	Light Shade = iota + 1
	Dark
)

// Note is a note that a positional argument of a kanon directive names.
//
//go:generate go tool kanon type=Note
//dokimi:enum
type Note uint8

// The notes.
const (
	Whole Note = iota + 1
	Half
)
