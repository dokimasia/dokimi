// Package commented declares enumerations whose variants have trailing
// comments.
package commented

// Grade is the grade of a review. Its directive takes the text of a variant
// from the trailing comment of the variant.
//
//dokimi:enum text=comment
type Grade int

// The grades of a review. The comment of Fair has two spaces after its
// marker, and Good has no comment, so its text is its name.
const (
	Poor Grade = iota // poor
	Fair              //  fair enough
	Good
)

// Mood is the mood of a reader. Its trailing comments explain each variant,
// and its directive leaves the text of each variant at its name.
//
//dokimi:enum
type Mood int

// The moods of a reader.
const (
	Calm  Mood = iota // the reader reads on
	Angry             // the reader writes a review
)
