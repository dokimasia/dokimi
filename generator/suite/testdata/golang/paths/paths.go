// Package paths routes the parcels of the shop. Each type that a call takes is
// on one path of the inputs alone.
package paths

import (
	"context"
	"errors"

	"example.com/shop/inputs"
)

// ErrLost reports a parcel that the depot lost.
var ErrLost = errors.New("paths: the parcel is lost")

// ErrGone reports a parcel that the depot does not have.
var ErrGone = errors.New("paths: the parcel is gone")

// ErrClosed reports a read of the rows after their close.
var ErrClosed = errors.New("paths: the rows are closed")

// Weight is the weight of a parcel.
type Weight int

// The weights of a parcel.
const (
	Light Weight = iota + 1
	Heavy
)

// Speed is the speed of a parcel.
type Speed int

// The speeds of a parcel.
const (
	Slow Speed = iota + 1
	Fast
)

// Zone is the zone of a parcel.
type Zone int

// The zones of a parcel.
const (
	North Zone = iota + 1
	South
)

// Size is the size of a parcel.
type Size int

// The sizes of a parcel.
const (
	Small Size = iota + 1
	Large
)

// Grade is the grade of a crate.
type Grade int

// The grades of a crate.
const (
	Low Grade = iota + 1
	High
)

// Crate is a crate of parcels, whose grade no caller sets.
type Crate struct {
	Label string
	grade Grade
}

// Item is one item of a parcel, with a declared default count. The prop tag of
// Name keeps the field, and the prop tag of Memo leaves it out.
type Item struct {
	Count int    `default:"1"`
	Name  string `prop:"name"`
	Memo  string `prop:"-"`
	owner string
}

// Slot is a generic slot, with a declared default size.
type Slot[T any] struct {
	Size  int `default:"2"`
	Value T
}

// Score is a generic score, with a declared default.
//
//dokimi:default 3
type Score[T any] int

// Label is a transparent alias of a string, with a declared default.
//
//dokimi:default "new"
type Label = string

// Depot keeps the items of the parcels.
type Depot interface {
	// Fetch returns the item with the identity id.
	//
	//dokimi:fault errors=[ErrLost]
	//dokimi:mixin not-found sentinel=ErrGone
	Fetch(ctx context.Context, id string) (Item, error)
	// Drop removes the item with the identity id.
	//
	//dokimi:mixin delete-removes read=Fetch sentinel=ErrGone
	Drop(ctx context.Context, id string) error
}

// Ledger counts the parcels.
type Ledger interface {
	// Total returns the number of parcels.
	Total() int
}

// Rows are the rows of a manifest, under the contract cursor, whose sentinel
// is no declared error of Next.
type Rows interface {
	// Next returns the next row, and false where the rows have none.
	//
	//dokimi:contract cursor role=next sentinel=ErrClosed
	Next(ctx context.Context) (string, bool, error)
	// Close closes the rows.
	//
	//dokimi:contract cursor role=close
	Close(ctx context.Context) error
}

// Router routes the parcels of the shop.
//
//dokimi:suite
type Router interface {
	// Weigh weighs the parcels.
	Weigh(weights ...Weight)
	// Pace reads the speeds of the channel.
	Pace(ctx context.Context, speeds <-chan Speed) error
	// Ship sends the parcels to the channel.
	Ship(out chan<- string)
	// Place places the parcel in the zone, which is never nil.
	Place(zone *Zone) error
	// Measure measures the parcel with fn.
	Measure(fn func() Size)
	// Pack packs the crate.
	Pack(c Crate)
	// Stock stocks an item of each package.
	Stock(item Item, other inputs.Item)
	// Hold holds the slot with its score and its label.
	Hold(s Slot[int], score Score[int], label Label)
	// Serve serves the parcels of the depot and the ledger.
	Serve(depot Depot, ledger Ledger) error
	// Read reads the rows of a manifest.
	Read(rows Rows) error
	// Scan counts the parcels with several readers at once.
	//
	//dokimi:mixin concurrent-readers
	Scan(ctx context.Context) (int, error)
}

// Idle routes no parcel.
//
//dokimi:suite
type Idle interface{}

// Route routes the parcels of the depot.
//
//dokimi:suite
func Route(depot Depot) int {
	return 0
}

// Memory routes the parcels in memory, under a budget of its own for Scan.
//
//dokimi:subject Router new=NewMemory
//dokimi:budget Router.Scan p99=1ms
type Memory struct {
	Router
}

// NewMemory returns a router in memory.
func NewMemory() *Memory {
	return &Memory{}
}
