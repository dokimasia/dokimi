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

// Node is a node of the route of a parcel. Its children are nodes.
type Node interface {
	// Name returns the name of the node.
	Name() string
	// Children returns the children of the node.
	Children() ([]Node, error)
}

// Branch is a branch of a route, whose next hop leads to the next branch.
type Branch interface {
	// Next returns the next hop of the branch.
	Next() Hop
}

// Hop is a hop of a route, through the branch Branch.
type Hop struct {
	Via    string
	Branch Branch
}

// Step is a step of a route, which returns the next step.
type Step func(at string) Step

// Ping is the first of two interfaces whose methods return each other.
type Ping interface {
	// Pong returns the reply to the ping.
	Pong() Pong
}

// Pong is the second of two interfaces whose methods return each other.
type Pong interface {
	// Ping returns the next ping.
	Ping() Ping
}

// Leaf is a leaf of a choice, which returns the name of the leaf.
type Leaf func() string

// Pair is a pair of leaves. Its results contain one function type twice, and
// not the pair itself.
type Pair interface {
	// Left returns the left leaf.
	Left() Leaf
	// Right returns the right leaf.
	Right() Leaf
}

// Choice is a choice between the leaves of a pair.
type Choice struct {
	Pair Pair
}

// Rally serves the first ping of a rally. Its results lead to Ping and Pong,
// which return each other, and not to the rally itself.
type Rally interface {
	// Serve returns the first ping.
	Serve() Ping
}

// Match is a match of one rally.
type Match struct {
	Rally Rally
}

// Chain is a chain of links, whose next link is a chain.
type Chain struct {
	Next *Chain
	Name string
}

// Linker returns the chain of a route. Its results lead to the chain, whose
// next link is a chain, and not to the linker itself.
type Linker interface {
	// Link returns the chain.
	Link() Chain
}

// Links are the links of a route, made by the linker Linker.
type Links struct {
	Linker Linker
}

// Tracker tracks the routes of the parcels of the shop.
//
//dokimi:suite
type Tracker interface {
	// Route counts the nodes of the route.
	Route(root Node) int
	// Follow counts the hops of the branch.
	Follow(b Branch) int
	// Walk counts the steps of the route.
	Walk(s Step) int
	// Volley counts the pings of the volley.
	Volley(p Ping) int
	// Track counts the labels of the labeler.
	Track(l Labeler) int
	// Pick counts the leaves of the choice.
	Pick(c Choice) int
	// Play counts the pings of the match.
	Play(m Match) int
	// Link counts the links of the route.
	Link(l Links) int
}

// Hire hires the labelers that hire returns, and returns how many it hired.
//
//dokimi:suite
func Hire(hire func() Labeler) int {
	return 0
}
