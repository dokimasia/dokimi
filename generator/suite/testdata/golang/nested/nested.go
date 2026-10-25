// Package nested ships the parcels of the shop with requests whose fields have
// function types and interface types.
package nested

import (
	"context"
	"errors"
	"io"
)

// ErrTorn reports a label that the labeler tore.
var ErrTorn = errors.New("nested: the label is torn")

// Receipt is the receipt of a label, which Void voids.
type Receipt struct {
	Text string
	Void func()
}

// Labeler prints the labels of the parcels.
type Labeler interface {
	// Label returns the receipt of the label of the parcel with the identity
	// id.
	//
	//dokimi:fault errors=[ErrTorn]
	Label(id string) (Receipt, error)
}

// Scale weighs values of one type.
type Scale[T any] interface {
	// Weigh returns the weight of the next value.
	Weigh() T
}

// Ack acknowledges the parcel with the identity id.
type Ack func(id string) error

// ShipRequest is the request of a shipment. Its fields have an interface type,
// a function type, a defined function type, a generic interface and an
// interface of the standard library.
type ShipRequest struct {
	Labeler Labeler
	Done    func() error
	Ack     Ack
	Scale   Scale[int]
	Body    io.Reader
	ID      string
}

// Batch is a batch of shipments, with a spare labeler and the ack of the
// batch.
type Batch struct {
	Requests []ShipRequest
	Spare    Labeler
	Ack      Ack
}

// Courier ships the parcels of the shop.
//
//dokimi:suite
type Courier interface {
	// Ship ships the parcel of the request.
	Ship(ctx context.Context, req *ShipRequest) error
	// Send ships each request of the batch.
	Send(b Batch) error
}

// Hire hires the labelers that hire returns, and returns how many it hired.
//
//dokimi:suite
func Hire(hire func() Labeler) int {
	return 0
}
