// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package kindstest

import (
	"fmt"
	"io"
	"time"
	"unsafe"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/kinds"
)

// InMemory implements [kinds.Every]. It receives from each channel until the
// channel closes, sends to each send channel, and calls each function and
// method that it takes, so each value that the suite draws has to end and has
// to take each send. The corpus applies the suite of Every to it as a
// control.
//
//dokimi:subject kinds.Every new=NewInMemory
type InMemory struct{}

// NewInMemory returns an InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Params implements [kinds.Every].
func (*InMemory) Params(
	_ bool, _ int, _ int8, _ int16, _ int32, _ int64, _ uint, _ uint8, _ uint16, _ uint32,
	_ uint64, _ uintptr, _ float32, _ float64, _ complex64, _ complex128, _ [2]int,
	ch chan int, fn func() int, str fmt.Stringer, _ map[string]int, _ *int, _ []int, _ string,
	_ struct{}, _ unsafe.Pointer,
) int {
	return drained(ch) + fn() + len(str.String()) + 23
}

// Fields implements [kinds.Every].
func (*InMemory) Fields(f *kinds.Fields) int {
	return drained(f.Chan) + f.Func() + len(f.Interface.String())
}

// Elements implements [kinds.Every].
func (*InMemory) Elements(e *kinds.Elements) int {
	n := 0
	for _, ch := range e.Chan {
		n += drained(ch)
	}
	for _, fn := range e.Func {
		n += fn()
	}
	for _, s := range e.Interface {
		n += len(s.String())
	}
	return n
}

// Results implements [kinds.Every].
func (*InMemory) Results(r kinds.Results) int {
	n := drained(r.Chan())
	if fn := r.Func(); fn != nil {
		n += fn()
	}
	if s := r.Interface(); s != nil {
		n += len(s.String())
	}
	_, _, _, _, _ = r.Bool(), r.Int(), r.String(), r.Struct(), r.UnsafePointer()
	return n
}

// Place implements [kinds.Every].
func (*InMemory) Place(clock kinds.Clock, key string) int {
	if timer := clock.NewTimer(time.Second); timer != nil {
		<-timer.C()
		timer.Stop()
	}
	return (clock.Now().Second() + len(key)) % 8
}

// Route implements [kinds.Every].
func (*InMemory) Route(in <-chan int, out chan<- int, both chan int) int {
	n := 0
	for v := range in {
		out <- v
		n++
	}
	return n + drained(both)
}

// Drain implements [kinds.Every].
func (*InMemory) Drain(f kinds.Feed) int {
	n := 0
	for _, in := range []<-chan int{f.In, f.Both, f.Ticks} {
		for v := range in {
			f.Out <- v
			n++
		}
	}
	return n
}

// Watch implements [kinds.Every].
func (*InMemory) Watch(s kinds.Source) int {
	n := 0
	for v := range s.Updates() {
		s.Results() <- v
		n++
	}
	return n + drained(s.Pipe())
}

// Measure implements [kinds.Every].
func (*InMemory) Measure(s kinds.Shape) float64 { return s.Area() }

// Index implements [kinds.Every].
func (*InMemory) Index(names map[fmt.Stringer]int) int {
	n := 0
	for name := range names {
		_ = name.String()
		n++
	}
	return n
}

// Count implements [kinds.Every].
func (*InMemory) Count(shapes map[kinds.Shape]int) int { return len(shapes) }

// Use implements [kinds.Every].
func (*InMemory) Use(r *kinds.Resources) int {
	n := 0
	if r.Err != nil {
		n++
	}
	if scope := r.Scope.Context(); scope != nil && scope.Err() == nil {
		n++
	}
	if read, _ := io.ReadAll(r.Reader); len(read) > 0 {
		n++
	}
	if read, _ := io.ReadAll(r.Buffered); len(read) > 0 {
		n++
	}
	if _, err := r.Writer.WriteString("use"); err == nil && r.Writer.Flush() == nil {
		n++
	}
	r.Slog.Info("use")
	r.Log.Print("use")
	if r.Mutex != nil {
		r.Mutex.Lock()
		defer r.Mutex.Unlock()
		n++
	}
	if r.Timer != nil {
		for range r.Timer.C {
			n++
		}
	}
	return n + r.Builder.Len() + r.Box.Value()
}

// drained returns the number of the values of ch, which it receives until ch
// closes.
func drained(ch <-chan int) int {
	n := 0
	for range ch {
		n++
	}
	return n
}

// The compiler checks that InMemory implements [kinds.Every].
var _ kinds.Every = (*InMemory)(nil)
