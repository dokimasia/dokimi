package validated

// Ring is a ring that the kanon directive in its doc names.
//
//go:generate go run go.thesmos.sh/kanon/cmd/kanon@v0.5.0 -type Ring
//dokimi:enum
type Ring uint8

// The rings.
const (
	Inner Ring = iota + 1
	Outer
)
