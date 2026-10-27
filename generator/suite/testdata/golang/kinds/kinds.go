// Package kinds declares the inputs of the suites of the shop with a value of
// each kind of reflect at each position: a parameter, a field, an element of
// a list and a result of a method of an input.
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
	"testing"
	"time"
	"unsafe"
)

// Params takes a value of each kind of reflect as a parameter. Each callable
// has the name of its kind, as reflect.Kind.String returns it.
//
//dokimi:suite
type Params interface {
	// Bool takes a bool.
	Bool(v bool)
	// Int takes an int.
	Int(v int)
	// Int8 takes an int8.
	Int8(v int8)
	// Int16 takes an int16.
	Int16(v int16)
	// Int32 takes an int32.
	Int32(v int32)
	// Int64 takes an int64.
	Int64(v int64)
	// Uint takes a uint.
	Uint(v uint)
	// Uint8 takes a uint8.
	Uint8(v uint8)
	// Uint16 takes a uint16.
	Uint16(v uint16)
	// Uint32 takes a uint32.
	Uint32(v uint32)
	// Uint64 takes a uint64.
	Uint64(v uint64)
	// Uintptr takes a uintptr.
	Uintptr(v uintptr)
	// Float32 takes a float32.
	Float32(v float32)
	// Float64 takes a float64.
	Float64(v float64)
	// Complex64 takes a complex64.
	Complex64(v complex64)
	// Complex128 takes a complex128.
	Complex128(v complex128)
	// Array takes an array.
	Array(v [2]int)
	// Chan takes a channel of both directions.
	Chan(v chan int)
	// Func takes a function.
	Func(v func() int)
	// Interface takes an interface of the standard library.
	Interface(v fmt.Stringer)
	// Map takes a map.
	Map(v map[string]int)
	// Ptr takes a pointer.
	Ptr(v *int)
	// Slice takes a slice.
	Slice(v []int)
	// String takes a string.
	String(v string)
	// Struct takes a struct without fields.
	Struct(v struct{})
	// UnsafePointer takes an unsafe pointer.
	UnsafePointer(v unsafe.Pointer)
}

// Fields has a field of each kind of reflect, named after the kind.
type Fields struct {
	Bool          bool
	Int           int
	Int8          int8
	Int16         int16
	Int32         int32
	Int64         int64
	Uint          uint
	Uint8         uint8
	Uint16        uint16
	Uint32        uint32
	Uint64        uint64
	Uintptr       uintptr
	Float32       float32
	Float64       float64
	Complex64     complex64
	Complex128    complex128
	Array         [2]int
	Chan          chan int
	Func          func() int
	Interface     fmt.Stringer
	Map           map[string]int
	Ptr           *int
	Slice         []int
	String        string
	Struct        struct{}
	UnsafePointer unsafe.Pointer
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
	// Interface returns an interface of the standard library.
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

// Every takes a value of each kind of reflect in a field, in an element of a
// list and in a result of a method of an input.
//
//dokimi:suite
type Every interface {
	// Fields takes a value of each kind in a field.
	Fields(f Fields)
	// Elements takes a value of each kind in an element of a list.
	Elements(e Elements)
	// Results takes a value of each kind in a result of a method of an input.
	Results(r Results)
}

// Timer fires once at a time.
type Timer interface {
	// C returns the channel that receives the time of the timer.
	C() <-chan time.Time
	// Stop stops the timer, and reports whether the timer had not fired.
	Stop() bool
}

// Clock reads the time, and makes the timers that fire after a delay.
type Clock interface {
	// Now returns the time.
	Now() time.Time
	// NewTimer returns a timer that fires after d.
	NewTimer(d time.Duration) Timer
}

// Placer places the shards of the shop by the time of a clock.
//
//dokimi:suite
type Placer interface {
	// Place returns the shard of the key, by the time of the clock.
	Place(clock Clock, key string) (int, error)
}

// Ticks is a channel of ticks that a subject receives from.
type Ticks <-chan int

// Feed has a channel of each direction and a defined channel type.
type Feed struct {
	In    <-chan string
	Out   chan<- string
	Both  chan string
	Ticks Ticks
}

// Source returns the channels of a feed.
type Source interface {
	// Updates returns the channel that the source sends its updates to.
	Updates() <-chan string
	// Results returns the channel that takes the results of the updates.
	Results() chan<- string
	// Pipe returns a channel of both directions.
	Pipe() chan string
}

// Router routes the messages of the shop over channels.
//
//dokimi:suite
type Router interface {
	// Route routes the messages of the feed.
	Route(ctx context.Context, f Feed) error
	// Read reads the messages of in until in closes.
	Read(in <-chan string) int
	// Write writes the messages to out.
	Write(out chan<- string)
	// Pass receives the messages of both and sends them back.
	Pass(both chan string)
	// Watch routes the updates of the source.
	Watch(s Source) error
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

// Corner is a shape of one corner, whose value implements Shape.
type Corner int

// The corners of a shape.
const (
	Left Corner = iota + 1
	Right
)

// Area returns the area of the corner.
func (Corner) Area() float64 { return 0 }

func (Corner) area() {}

// Blob is a generic shape.
type Blob[T any] struct {
	Value T
}

// Area returns the area of the blob.
func (Blob[T]) Area() float64 { return 1 }

func (Blob[T]) area() {}

// dot is a shape that the package does not export.
type dot struct{}

// Area returns the area of the dot.
func (dot) Area() float64 { return 0 }

func (dot) area() {}

// Polygon is a shape with corners.
type Polygon interface {
	Shape
	// Corners returns the number of corners.
	Corners() int
}

// Sealed has a method that the package does not export, and the package
// exports no implementation of it.
type Sealed interface {
	seal()
}

// Picker picks a sealed value.
type Picker interface {
	// Pick returns the sealed value.
	Pick() Sealed
}

// Drawer draws the shapes of the shop.
//
//dokimi:suite
type Drawer interface {
	// Draw draws the shape, and returns its area.
	Draw(s Shape) float64
	// Count counts the shapes of the map.
	Count(shapes map[Shape]int) int
	// Name names the values of the map.
	Name(names map[fmt.Stringer]string) int
	// Pick picks a value with the picker.
	Pick(p Picker) int
	// Seal seals the value.
	Seal(s Sealed)
	// Test reports to the test.
	Test(tb testing.TB)
}

// Label is a defined interface type over fmt.Stringer.
type Label fmt.Stringer

// Anything is a defined type over any.
type Anything any

// Namer names a value.
type Namer interface {
	// Name returns the name of the value.
	Name() string
}

// Ring1 and Ring2 are aliases of each other, which no compiler accepts.
type (
	Ring1 = Ring2
	Ring2 = Ring1
)

// Odd takes interfaces without a declaration of their own.
//
//dokimi:suite
type Odd interface {
	// Tag tags the label.
	Tag(l Label)
	// Keep keeps the value.
	Keep(a Anything)
	// Embed sizes the value with a name.
	Embed(v interface {
		Namer
		Size() int
		Reset()
	}) int
	// Loop takes a type that no compiler accepts.
	Loop(r Ring1)
}

// Box is a generic box whose value the package does not export.
type Box[T any] struct {
	v T
}

// Twins has two fields of one JSON name.
type Twins struct {
	A int `json:"x"`
	B int `json:"x"`
}

// Clash has a field with the JSON name of the type that it embeds.
type Clash struct {
	Circle
	Name string `json:"Circle"`
}

// Holder has a field of each type that the property engine does not read,
// and two inline structs that it reads.
type Holder struct {
	Err     error
	Ctx     context.Context
	Ptr     unsafe.Pointer
	Box     Box[int]
	Mutex   *sync.Mutex
	Timer   *time.Timer
	Builder strings.Builder
	Named   interface{ Name() string }
	Private struct{ n int }
	Point   struct{ X, Y int }
	Inside  struct{ Circle }
}

// Std takes values of the standard library and errors.
//
//dokimi:suite
type Std interface {
	// Copy copies src to dst.
	Copy(dst io.Writer, src io.Reader) error
	// Close closes c.
	Close(c io.ReadCloser) error
	// From reads into r.
	From(r io.ReaderFrom) error
	// At writes into w.
	At(w io.WriterAt) error
	// Log logs to both loggers.
	Log(l *slog.Logger, o *log.Logger)
	// Buffer reads r into w.
	Buffer(r *bufio.Reader, w *bufio.Writer) error
	// Lock locks mu.
	Lock(mu *sync.Mutex)
	// Fail fails with err.
	Fail(err error) error
	// Within runs within ctx and parent.
	Within(ctx context.Context, parent context.Context) error
	// Hold holds the values of h.
	Hold(h Holder)
	// Pair pairs the twins.
	Pair(t []Twins)
	// Meet meets the clash.
	Meet(c Clash)
}
