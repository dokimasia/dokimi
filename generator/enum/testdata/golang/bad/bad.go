// Package bad declares the enumerations that enum reports.
package bad

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
