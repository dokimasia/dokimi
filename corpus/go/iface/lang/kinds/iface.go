// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package kinds is the fixture of inputs of each kind of reflect at each
// position: a parameter, a field, an element of a list and a result of a
// method of an input. It also has a clock whose timers return a receive
// channel, channels of each direction, an interface with an unexported
// method, an interface as the key of a map, and values of the standard
// library that the property engine does not read. The suite gives each type a
// generator, and the subject in memory drains each channel and calls each
// function and method that it takes.
package kinds

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Fields has a field of each kind of reflect, named after the kind, in the
// order of their alignment.
type Fields struct {
	Struct        struct{}
	Interface     fmt.Stringer
	Chan          chan int
	Func          func() int
	Map           map[string]int
	Ptr           *int
	UnsafePointer unsafe.Pointer
	String        string
	Slice         []int
	Complex128    complex128
	Array         [2]int
	Int           int
	Int64         int64
	Uint          uint
	Uint64        uint64
	Uintptr       uintptr
	Float64       float64
	Complex64     complex64
	Int32         int32
	Uint32        uint32
	Float32       float32
	Int16         int16
	Uint16        uint16
	Bool          bool
	Int8          int8
	Uint8         uint8
}

// Elements has a list of each kind of reflect, named after the kind.
type Elements struct {
	Bool          []bool
	Int           []int
	Int8          []int8
	Int16         []int16
	Int32         []int32
	Int64         []int64
	Uint          []uint
	Uint8         []uint8
	Uint16        []uint16
	Uint32        []uint32
	Uint64        []uint64
	Uintptr       []uintptr
	Float32       []float32
	Float64       []float64
	Complex64     []complex64
	Complex128    []complex128
	Array         [][2]int
	Chan          []chan int
	Func          []func() int
	Interface     []fmt.Stringer
	Map           []map[string]int
	Ptr           []*int
	Slice         [][]int
	String        []string
	Struct        []struct{}
	UnsafePointer []unsafe.Pointer
}

// Results has a method that returns a value of each kind of reflect, named
// after the kind.
//
//nolint:interfacebloat // the fixture has a method of each kind of reflect
type Results interface {
	// Bool returns a bool.
	Bool() bool
	// Int returns an int.
	Int() int
	// Int8 returns an int8.
	Int8() int8
	// Int16 returns an int16.
	Int16() int16
	// Int32 returns an int32.
	Int32() int32
	// Int64 returns an int64.
	Int64() int64
	// Uint returns a uint.
	Uint() uint
	// Uint8 returns a uint8.
	Uint8() uint8
	// Uint16 returns a uint16.
	Uint16() uint16
	// Uint32 returns a uint32.
	Uint32() uint32
	// Uint64 returns a uint64.
	Uint64() uint64
	// Uintptr returns a uintptr.
	Uintptr() uintptr
	// Float32 returns a float32.
	Float32() float32
	// Float64 returns a float64.
	Float64() float64
	// Complex64 returns a complex64.
	Complex64() complex64
	// Complex128 returns a complex128.
	Complex128() complex128
	// Array returns an array.
	Array() [2]int
	// Chan returns a channel of both directions.
	Chan() chan int
	// Func returns a function.
	Func() func() int
	// Interface returns a value of an interface of the standard library.
	Interface() fmt.Stringer
	// Map returns a map.
	Map() map[string]int
	// Ptr returns a pointer.
	Ptr() *int
	// Slice returns a slice.
	Slice() []int
	// String returns a string.
	String() string
	// Struct returns a struct without fields.
	Struct() struct{}
	// UnsafePointer returns an unsafe pointer.
	UnsafePointer() unsafe.Pointer
}

// Timer sends the time on its channel once.
type Timer interface {
	// C returns the channel that receives the time of the timer.
	C() <-chan time.Time
	// Stop stops the timer, and reports whether the timer was active.
	Stop() bool
}

// Clock reads the time, and makes the timers that send the time after a
// delay.
type Clock interface {
	// Now returns the time.
	Now() time.Time
	// NewTimer returns a timer that sends the time after d.
	NewTimer(d time.Duration) Timer
}

// Ticks is a channel of ticks that a subject receives from.
type Ticks <-chan int

// Feed has a channel of each direction and a defined channel type.
type Feed struct {
	In    <-chan int
	Out   chan<- int
	Both  chan int
	Ticks Ticks
}

// Source returns the channels of a feed.
type Source interface {
	// Updates returns the channel that the source sends its updates to.
	Updates() <-chan int
	// Results returns the channel that takes the results of the updates.
	Results() chan<- int
	// Pipe returns a channel of both directions.
	Pipe() chan int
}

// Shape is a shape of the shop. Its method area is not exported, so the
// package kinds alone implements it.
type Shape interface {
	// Area returns the area of the shape.
	Area() float64
	area()
}

// Circle is a shape whose value implements Shape.
type Circle struct {
	Radius float64
}

// Area returns the area of the circle.
func (c Circle) Area() float64 { return 3 * c.Radius * c.Radius }

func (Circle) area() {}

// Square is a shape whose pointer implements Shape.
type Square struct {
	Side float64
}

// Area returns the area of the square.
func (s *Square) Area() float64 { return s.Side * s.Side }

func (*Square) area() {}

// Box is a generic box whose value the package does not export.
type Box[T any] struct {
	v T
}

// Value returns the value of the box.
func (b Box[T]) Value() T { return b.v }

// Scope returns the context of a scope.
type Scope interface {
	// Context returns the context of the scope.
	Context() context.Context
}

// Resources has a field of each type of the standard library that the
// property engine does not read, an error outside a callable, and a scope
// whose context is outside a callable.
type Resources struct {
	Err      error
	Scope    Scope
	Reader   io.Reader
	Buffered *bufio.Reader
	Writer   *bufio.Writer
	Slog     *slog.Logger
	Log      *log.Logger
	Mutex    *sync.Mutex
	Timer    *time.Timer
	Builder  strings.Builder
	Box      Box[int]
}

// Every takes a value of each kind of reflect at each position, and the
// values of the standard library that the property engine does not read.
//
//dokimi:suite
//dokimi:out path=kindstest/
//nolint:interfacebloat // the fixture has a callable of each position of a value
type Every interface {
	// Params counts the parameters that it takes, one of each kind.
	Params(
		b bool, i int, i8 int8, i16 int16, i32 int32, i64 int64, u uint, u8 uint8, u16 uint16, u32 uint32,
		u64 uint64, uptr uintptr, f32 float32, f64 float64, c64 complex64, c128 complex128, arr [2]int,
		ch chan int, fn func() int, str fmt.Stringer, m map[string]int, p *int, s []int, text string,
		st struct{}, up unsafe.Pointer,
	) int
	// Fields counts the fields of f, and drains its channel.
	Fields(f *Fields) int
	// Elements counts the elements of e, and drains its channels.
	Elements(e *Elements) int
	// Results calls each method of r once, and counts the values that they
	// return.
	Results(r Results) int
	// Place waits for a timer of the clock, and returns the shard of the key
	// by the time of the clock.
	Place(clock Clock, key string) int
	// Route sends each value of in to out, drains both, and returns the
	// number of values that it sent.
	Route(in <-chan int, out chan<- int, both chan int) int
	// Drain drains the channels of f into its channel Out, and returns the
	// number of values that it sent.
	Drain(f Feed) int
	// Watch sends each update of s to its results, drains its pipe, and
	// returns the number of updates.
	Watch(s Source) int
	// Measure returns the area of the shape.
	Measure(s Shape) float64
	// Index returns the number of the keys of names.
	Index(names map[fmt.Stringer]int) int
	// Count returns the number of the keys of shapes.
	Count(shapes map[Shape]int) int
	// Use reads or calls each value of r, and returns the number of the values
	// that it used.
	Use(r *Resources) int
}
