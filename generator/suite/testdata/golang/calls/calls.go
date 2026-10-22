// Package calls places the points of a game.
package calls

import (
	"context"
	"time"

	shop "example.com/shop/store"
)

// Shape is the shape of a point.
type Shape int

// The shapes of a point.
const (
	Round Shape = iota + 1
	Square
)

// Point is a point of a board.
type Point struct {
	X, Y  int
	Shape Shape
}

// Tint is the tint of a pair.
type Tint int

// The tints of a pair.
const (
	Light Tint = iota + 1
	Dark
)

// Pair is a pair of values of one type.
type Pair[T any] struct {
	First, Second T
}

// Ink is the ink of a stamp.
type Ink int

// The inks of a stamp.
const (
	Red Ink = iota + 1
	Blue
)

// Stamp is the stamp of a cell.
type Stamp struct {
	Ink Ink
}

// Cell is a cell of a board with a stamp.
type Cell struct {
	Stamp
	Note string
}

// Hue is the hue of a palette.
type Hue int

// The hues of a palette.
const (
	Warm Hue = iota + 1
	Cool
)

// Palette is a list of hues.
type Palette []Hue

// Board places the points of a game.
//
//dokimi:suite
type Board interface {
	// Place places the point, which is never nil.
	Place(p *Point) error
	// Clash takes parameters with the names that a generated call takes for
	// itself.
	Clash(ctx context.Context, c, subject, err string, prop int) error
	// Blank takes parameters without names.
	Blank(string, int)
	// Index counts the points of each shape.
	Index(points map[Shape][]Point) int
	// Label takes a parameter with the name of its type.
	Label(string string)
	// Wait takes a parameter with the name of the package of a list of
	// delays.
	Wait(time int, delays []time.Duration)
	// Hold takes a parameter with the name of the package of a type argument.
	Hold(time int, pair Pair[time.Month])
	// Swap swaps the tints of a pair.
	Swap(pair Pair[Tint]) Pair[Tint]
	// Mark marks the cell.
	Mark(cell Cell) error
	// Paint paints the board in the palette.
	Paint(palette Palette)
	// Repeat takes a parameter with the name of the package of the board.
	Repeat(calls int)
	// Keep takes a parameter with the name of the package of an item, which
	// the file imports under another name.
	Keep(store int, item shop.Item)
	// Stock takes a parameter with the name under which the file imports the
	// package of an item.
	Stock(shop int, item shop.Item)
	// Discard takes a parameter with the blank name.
	Discard(_ string)
	// Twice takes a parameter with the name that another parameter takes in
	// its place.
	Twice(c, c0 string)
	// Drain takes a parameter with the name of the list of a channel.
	Drain(ints <-chan int, intsValues int)
	// Gather takes parameters with the names of the list and of the counter of
	// a generated function.
	Gather(fn func() int, fnResults, fnAt int)
}
