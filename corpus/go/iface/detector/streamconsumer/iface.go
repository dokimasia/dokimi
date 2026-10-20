// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package streamconsumer is the detector-axis fixture for the streamconsumer shape:
// a stream in rather than out.
//
// It is not quite the mirror of a stream reader, and the asymmetry is worth
// knowing: eidos recognises a produced stream as an iter.Seq2 return, but a
// consumed stream as an *interface* parameter. A method taking iter.Seq is a
// method taking a func, and no amount of naming makes the detector see a
// stream in it.
//
// No directive appears here: the classification comes from the signature
// alone, so a detector that misfires shows up as wrong generated output
// rather than as a directive that was misread.
package streamconsumer

import "context"

// Value is the element type the fixture ingests.
type Value struct{ Key, Body string }

// Source is the stream Ingest consumes. It has to be an interface: the
// detector tests whether the parameter is one, and an iter.Seq is a func type
// however stream-like it reads.
//
// Both fixture values are declared rather than derived, because an interface
// has no literal form and neither half could be worked out. Without the
// sample directive every check needing a Source is declined, which is what
// the named pair below is here to close.
//
//dokimi:stub
//dokimi:out path=streamconsumertest/
type Source interface {
	Next(ctx context.Context) (Value, bool, error)
}

// StreamConsumer is the fixture interface.
//
// The model tier is unarmed and the reason is on the one method: a stream
// is consumed by reading it, and the fixture holds one value where driving
// two sides needs a fresh stream each. The suite tier's checks over Ingest
// are the ones a stored value can answer, and they exist because Source
// declares its pair — which is what this fixture is here to show.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=streamconsumertest/
type StreamConsumer interface {
	// Ingest takes one non-context parameter and returns one value plus an
	// error. Both counts are load-bearing: a second parameter or a second
	// value takes the method out of the shape.
	Ingest(ctx context.Context, src Source) (int, error)
}

// SourceDefaults is the seed companion the harness reads for a Source, found
// by name the way every Defaults companion is.
func SourceDefaults() Source { return &listSource{} }

// SourceAlternate is the second value, named by the fixture directive because
// nothing can derive one that differs from the first.
func SourceAlternate() Source { return &listSource{items: []Value{{Key: "k", Body: "b"}}} }

// listSource walks a fixed list and then reports exhaustion.
type listSource struct{ items []Value }

func (s *listSource) Next(context.Context) (Value, bool, error) {
	if len(s.items) == 0 {
		return Value{}, false, nil
	}
	v := s.items[0]
	s.items = s.items[1:]
	return v, true, nil
}
