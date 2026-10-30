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
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Params:
//
// Params.Bool/smoke: Bool returns without a panic for derived inputs.
//
// Params.Int/smoke: Int returns without a panic for derived inputs.
//
// Params.Int8/smoke: Int8 returns without a panic for derived inputs.
//
// Params.Int16/smoke: Int16 returns without a panic for derived inputs.
//
// Params.Int32/smoke: Int32 returns without a panic for derived inputs.
//
// Params.Int64/smoke: Int64 returns without a panic for derived inputs.
//
// Params.Uint/smoke: Uint returns without a panic for derived inputs.
//
// Params.Uint8/smoke: Uint8 returns without a panic for derived inputs.
//
// Params.Uint16/smoke: Uint16 returns without a panic for derived inputs.
//
// Params.Uint32/smoke: Uint32 returns without a panic for derived inputs.
//
// Params.Uint64/smoke: Uint64 returns without a panic for derived inputs.
//
// Params.Uintptr/smoke: Uintptr returns without a panic for derived
// inputs.
//
// Params.Float32/smoke: Float32 returns without a panic for derived
// inputs.
//
// Params.Float64/smoke: Float64 returns without a panic for derived
// inputs.
//
// Params.Complex64/smoke: Complex64 returns without a panic for derived
// inputs.
//
// Params.Complex128/smoke: Complex128 returns without a panic for derived
// inputs.
//
// Params.Array/smoke: Array returns without a panic for derived inputs.
//
// Params.Func/smoke: Func returns without a panic for derived inputs.
//
// Params.Map/smoke: Map returns without a panic for derived inputs.
//
// Params.Slice/smoke: Slice returns without a panic for derived inputs.
//
// Params.String/smoke: String returns without a panic for derived inputs.

// The derived checks of Params that witness their laws at the shallow end:
//
// Params.Bool/smoke
//
// Params.Int/smoke
//
// Params.Int8/smoke
//
// Params.Int16/smoke
//
// Params.Int32/smoke
//
// Params.Int64/smoke
//
// Params.Uint/smoke
//
// Params.Uint8/smoke
//
// Params.Uint16/smoke
//
// Params.Uint32/smoke
//
// Params.Uint64/smoke
//
// Params.Uintptr/smoke
//
// Params.Float32/smoke
//
// Params.Float64/smoke
//
// Params.Complex64/smoke
//
// Params.Complex128/smoke
//
// Params.Array/smoke
//
// Params.Func/smoke
//
// Params.Map/smoke
//
// Params.Slice/smoke
//
// Params.String/smoke

// The unclaimed classifications of Params:
//
// The shape:mutator of Params.Bool has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Int has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Int8 has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Int16 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Int32 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Int64 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Uint has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Uint8 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Uint16 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Uint32 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Uint64 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Uintptr has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Float32 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Float64 has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Complex64 has the stance silent. A mutation
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Complex128 has the stance silent. A mutation
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Array has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Chan has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Func has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Interface has the stance silent. A mutation
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Map has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Ptr has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Slice has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.String has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.Struct has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Params.UnsafePointer has the stance silent. A
// mutation is an input of the other checks, and its result is an error
// alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Params writes.
type (
	ParamsHandle = conform.Handle[Params]
	ParamsCheck  = conform.Check[Params]
	ParamsOption = conform.Option[Params]
)

// ParamsChecks contains the name of each callable of Params, with the
// identities of its checks.
var ParamsChecks = struct {
	Bool          paramsBool
	Int           paramsInt
	Int8          paramsInt8
	Int16         paramsInt16
	Int32         paramsInt32
	Int64         paramsInt64
	Uint          paramsUint
	Uint8         paramsUint8
	Uint16        paramsUint16
	Uint32        paramsUint32
	Uint64        paramsUint64
	Uintptr       paramsUintptr
	Float32       paramsFloat32
	Float64       paramsFloat64
	Complex64     paramsComplex64
	Complex128    paramsComplex128
	Array         paramsArray
	Chan          paramsChan
	Func          paramsFunc
	Interface     paramsInterface
	Map           paramsMap
	Ptr           paramsPtr
	Slice         paramsSlice
	String        paramsString
	Struct        paramsStruct
	UnsafePointer paramsUnsafePointer
}{
	Bool:          paramsBool{"Params.Bool"},
	Int:           paramsInt{"Params.Int"},
	Int8:          paramsInt8{"Params.Int8"},
	Int16:         paramsInt16{"Params.Int16"},
	Int32:         paramsInt32{"Params.Int32"},
	Int64:         paramsInt64{"Params.Int64"},
	Uint:          paramsUint{"Params.Uint"},
	Uint8:         paramsUint8{"Params.Uint8"},
	Uint16:        paramsUint16{"Params.Uint16"},
	Uint32:        paramsUint32{"Params.Uint32"},
	Uint64:        paramsUint64{"Params.Uint64"},
	Uintptr:       paramsUintptr{"Params.Uintptr"},
	Float32:       paramsFloat32{"Params.Float32"},
	Float64:       paramsFloat64{"Params.Float64"},
	Complex64:     paramsComplex64{"Params.Complex64"},
	Complex128:    paramsComplex128{"Params.Complex128"},
	Array:         paramsArray{"Params.Array"},
	Chan:          paramsChan{"Params.Chan"},
	Func:          paramsFunc{"Params.Func"},
	Interface:     paramsInterface{"Params.Interface"},
	Map:           paramsMap{"Params.Map"},
	Ptr:           paramsPtr{"Params.Ptr"},
	Slice:         paramsSlice{"Params.Slice"},
	String:        paramsString{"Params.String"},
	Struct:        paramsStruct{"Params.Struct"},
	UnsafePointer: paramsUnsafePointer{"Params.UnsafePointer"},
}

// paramsBool is the name of Params.Bool, with the identities of its
// checks.
type paramsBool struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Bool returns without a
// panic for derived inputs.
func (c paramsBool) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsInt is the name of Params.Int, with the identities of its checks.
type paramsInt struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Int returns without a panic
// for derived inputs.
func (c paramsInt) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsInt8 is the name of Params.Int8, with the identities of its
// checks.
type paramsInt8 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Int8 returns without a
// panic for derived inputs.
func (c paramsInt8) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsInt16 is the name of Params.Int16, with the identities of its
// checks.
type paramsInt16 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Int16 returns without a
// panic for derived inputs.
func (c paramsInt16) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsInt32 is the name of Params.Int32, with the identities of its
// checks.
type paramsInt32 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Int32 returns without a
// panic for derived inputs.
func (c paramsInt32) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsInt64 is the name of Params.Int64, with the identities of its
// checks.
type paramsInt64 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Int64 returns without a
// panic for derived inputs.
func (c paramsInt64) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsUint is the name of Params.Uint, with the identities of its
// checks.
type paramsUint struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Uint returns without a
// panic for derived inputs.
func (c paramsUint) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsUint8 is the name of Params.Uint8, with the identities of its
// checks.
type paramsUint8 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Uint8 returns without a
// panic for derived inputs.
func (c paramsUint8) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsUint16 is the name of Params.Uint16, with the identities of its
// checks.
type paramsUint16 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Uint16 returns without a
// panic for derived inputs.
func (c paramsUint16) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsUint32 is the name of Params.Uint32, with the identities of its
// checks.
type paramsUint32 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Uint32 returns without a
// panic for derived inputs.
func (c paramsUint32) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsUint64 is the name of Params.Uint64, with the identities of its
// checks.
type paramsUint64 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Uint64 returns without a
// panic for derived inputs.
func (c paramsUint64) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsUintptr is the name of Params.Uintptr, with the identities of its
// checks.
type paramsUintptr struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Uintptr returns without a
// panic for derived inputs.
func (c paramsUintptr) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsFloat32 is the name of Params.Float32, with the identities of its
// checks.
type paramsFloat32 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Float32 returns without a
// panic for derived inputs.
func (c paramsFloat32) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsFloat64 is the name of Params.Float64, with the identities of its
// checks.
type paramsFloat64 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Float64 returns without a
// panic for derived inputs.
func (c paramsFloat64) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsComplex64 is the name of Params.Complex64, with the identities of
// its checks.
type paramsComplex64 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Complex64 returns without a
// panic for derived inputs.
func (c paramsComplex64) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsComplex128 is the name of Params.Complex128, with the identities
// of its checks.
type paramsComplex128 struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Complex128 returns without
// a panic for derived inputs.
func (c paramsComplex128) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsArray is the name of Params.Array, with the identities of its
// checks.
type paramsArray struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Array returns without a
// panic for derived inputs.
func (c paramsArray) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsChan is the name of Params.Chan, with the identities of its
// checks.
type paramsChan struct {
	conform.Callable[Params]
}

// paramsFunc is the name of Params.Func, with the identities of its
// checks.
type paramsFunc struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Func returns without a
// panic for derived inputs.
func (c paramsFunc) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsInterface is the name of Params.Interface, with the identities of
// its checks.
type paramsInterface struct {
	conform.Callable[Params]
}

// paramsMap is the name of Params.Map, with the identities of its checks.
type paramsMap struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Map returns without a panic
// for derived inputs.
func (c paramsMap) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsPtr is the name of Params.Ptr, with the identities of its checks.
type paramsPtr struct {
	conform.Callable[Params]
}

// paramsSlice is the name of Params.Slice, with the identities of its
// checks.
type paramsSlice struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that Slice returns without a
// panic for derived inputs.
func (c paramsSlice) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsString is the name of Params.String, with the identities of its
// checks.
type paramsString struct {
	conform.Callable[Params]
}

// Smoke returns the identity of the check that String returns without a
// panic for derived inputs.
func (c paramsString) Smoke() conform.ID[Params] {
	return c.Check("smoke")
}

// paramsStruct is the name of Params.Struct, with the identities of its
// checks.
type paramsStruct struct {
	conform.Callable[Params]
}

// paramsUnsafePointer is the name of Params.UnsafePointer, with the
// identities of its checks.
type paramsUnsafePointer struct {
	conform.Callable[Params]
}

// ParamsSuite returns the checks of Params, and the calls and the budgets
// of its callables.
func ParamsSuite() conform.Suite[Params] {
	return conform.Suite[Params]{
		Declaration: "Params",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Params]{
			{
				ID:          ParamsChecks.Bool.Smoke(),
				Statement:   "Bool returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Bool(bool(true))
						return nil
					})
				},
				Defect: paramsBoolPanics,
			},
			{
				ID:          ParamsChecks.Int.Smoke(),
				Statement:   "Int returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Int(int(42))
						return nil
					})
				},
				Defect: paramsIntPanics,
			},
			{
				ID:          ParamsChecks.Int8.Smoke(),
				Statement:   "Int8 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Int8(int8(42))
						return nil
					})
				},
				Defect: paramsInt8Panics,
			},
			{
				ID:          ParamsChecks.Int16.Smoke(),
				Statement:   "Int16 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Int16(int16(42))
						return nil
					})
				},
				Defect: paramsInt16Panics,
			},
			{
				ID:          ParamsChecks.Int32.Smoke(),
				Statement:   "Int32 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Int32(int32(42))
						return nil
					})
				},
				Defect: paramsInt32Panics,
			},
			{
				ID:          ParamsChecks.Int64.Smoke(),
				Statement:   "Int64 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Int64(int64(42))
						return nil
					})
				},
				Defect: paramsInt64Panics,
			},
			{
				ID:          ParamsChecks.Uint.Smoke(),
				Statement:   "Uint returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Uint(uint(42))
						return nil
					})
				},
				Defect: paramsUintPanics,
			},
			{
				ID:          ParamsChecks.Uint8.Smoke(),
				Statement:   "Uint8 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Uint8(uint8(1))
						return nil
					})
				},
				Defect: paramsUint8Panics,
			},
			{
				ID:          ParamsChecks.Uint16.Smoke(),
				Statement:   "Uint16 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Uint16(uint16(42))
						return nil
					})
				},
				Defect: paramsUint16Panics,
			},
			{
				ID:          ParamsChecks.Uint32.Smoke(),
				Statement:   "Uint32 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Uint32(uint32(42))
						return nil
					})
				},
				Defect: paramsUint32Panics,
			},
			{
				ID:          ParamsChecks.Uint64.Smoke(),
				Statement:   "Uint64 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Uint64(uint64(42))
						return nil
					})
				},
				Defect: paramsUint64Panics,
			},
			{
				ID:          ParamsChecks.Uintptr.Smoke(),
				Statement:   "Uintptr returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Uintptr(uintptr(42))
						return nil
					})
				},
				Defect: paramsUintptrPanics,
			},
			{
				ID:          ParamsChecks.Float32.Smoke(),
				Statement:   "Float32 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Float32(float32(1.5))
						return nil
					})
				},
				Defect: paramsFloat32Panics,
			},
			{
				ID:          ParamsChecks.Float64.Smoke(),
				Statement:   "Float64 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Float64(float64(1.5))
						return nil
					})
				},
				Defect: paramsFloat64Panics,
			},
			{
				ID:          ParamsChecks.Complex64.Smoke(),
				Statement:   "Complex64 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Complex64(complex64(1.5))
						return nil
					})
				},
				Defect: paramsComplex64Panics,
			},
			{
				ID:          ParamsChecks.Complex128.Smoke(),
				Statement:   "Complex128 returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Complex128(complex128(1.5))
						return nil
					})
				},
				Defect: paramsComplex128Panics,
			},
			{
				ID:          ParamsChecks.Array.Smoke(),
				Statement:   "Array returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Array([2]int{42})
						return nil
					})
				},
				Defect: paramsArrayPanics,
			},
			{
				ID:          ParamsChecks.Func.Smoke(),
				Statement:   "Func returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Func(func() int { return 0 })
						return nil
					})
				},
				Defect: paramsFuncPanics,
			},
			{
				ID:          ParamsChecks.Map.Smoke(),
				Statement:   "Map returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Map(map[string]int{"v-a": 42})
						return nil
					})
				},
				Defect: paramsMapPanics,
			},
			{
				ID:          ParamsChecks.Slice.Smoke(),
				Statement:   "Slice returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Slice([]int{42})
						return nil
					})
				},
				Defect: paramsSlicePanics,
			},
			{
				ID:          ParamsChecks.String.Smoke(),
				Statement:   "String returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ParamsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().String(string("v-a"))
						return nil
					})
				},
				Defect: paramsStringPanics,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:mutator",
				Callable:       "Params.Bool",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Int",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Int8",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Int16",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Int32",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Int64",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Uint",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Uint8",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Uint16",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Uint32",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Uint64",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Uintptr",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Float32",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Float64",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Complex64",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Complex128",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Array",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Chan",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Func",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Interface",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Map",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Ptr",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Slice",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.String",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.Struct",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Params.UnsafePointer",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Params]{
			{
				Calls:    paramsBoolCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Bool.Callable,
			},
			{
				Calls:    paramsIntCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Int.Callable,
			},
			{
				Calls:    paramsInt8Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Int8.Callable,
			},
			{
				Calls:    paramsInt16Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Int16.Callable,
			},
			{
				Calls:    paramsInt32Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Int32.Callable,
			},
			{
				Calls:    paramsInt64Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Int64.Callable,
			},
			{
				Calls:    paramsUintCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Uint.Callable,
			},
			{
				Calls:    paramsUint8Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Uint8.Callable,
			},
			{
				Calls:    paramsUint16Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Uint16.Callable,
			},
			{
				Calls:    paramsUint32Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Uint32.Callable,
			},
			{
				Calls:    paramsUint64Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Uint64.Callable,
			},
			{
				Calls:    paramsUintptrCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Uintptr.Callable,
			},
			{
				Calls:    paramsFloat32Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Float32.Callable,
			},
			{
				Calls:    paramsFloat64Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Float64.Callable,
			},
			{
				Calls:    paramsComplex64Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Complex64.Callable,
			},
			{
				Calls:    paramsComplex128Calls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Complex128.Callable,
			},
			{
				Calls:    paramsArrayCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Array.Callable,
			},
			{
				Calls:    paramsChanCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Chan.Callable,
			},
			{
				Calls:    paramsFuncCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Func.Callable,
			},
			{
				Calls:    paramsInterfaceCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Interface.Callable,
			},
			{
				Calls:    paramsMapCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Map.Callable,
			},
			{
				Calls:    paramsPtrCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Ptr.Callable,
			},
			{
				Calls:    paramsSliceCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Slice.Callable,
			},
			{
				Calls:    paramsStringCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.String.Callable,
			},
			{
				Calls:    paramsStructCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.Struct.Callable,
			},
			{
				Calls:    paramsUnsafePointerCalls,
				Over:     newParamsOver,
				Callable: ParamsChecks.UnsafePointer.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(paramsStructValues),
			prop.Using(conform.Stream[chan int](prop.Of[int]())),
			prop.Using(prop.Composite(func(*prop.Case) unsafe.Pointer { return nil })),
		},
	}
}

// RunParamsSuite applies the checks of Params, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunParamsSuite[T Params](t *testing.T, newSubject func() T, opts ...ParamsOption) {
	t.Helper()
	RunParamsSuiteFromHandle(t, func(*ParamsHandle) T { return newSubject() }, opts...)
}

// RunParamsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunParamsSuiteFromHandle[T Params](
	t *testing.T,
	newSubject func(h *ParamsHandle) T,
	opts ...ParamsOption,
) {
	t.Helper()
	conform.Run(t, ParamsSuite(), paramsSubjects(newSubject), opts...)
}

// FuzzParamsSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzParamsSuite[T Params](
	f *testing.F,
	callable conform.Callable[Params],
	newSubject func(h *ParamsHandle) T,
	opts ...ParamsOption,
) {
	f.Helper()
	conform.Fuzz(f, ParamsSuite(), callable, paramsSubjects(newSubject), opts...)
}

// RunParamsBudgets applies the allocation ceilings of Params and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunParamsBudgets[T Params](
	t *testing.T,
	newSubject func(h *ParamsHandle) T,
	opts ...ParamsOption,
) {
	t.Helper()
	conform.RunBudgets(t, ParamsSuite(), paramsSubjects(newSubject), opts...)
}

// BenchParamsSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchParamsSuite[T Params](
	b conform.B,
	callable conform.Callable[Params],
	newSubject func(h *ParamsHandle) T,
	opts ...ParamsOption,
) {
	b.Helper()
	conform.Bench(b, ParamsSuite(), callable, paramsSubjects(newSubject), opts...)
}

// paramsSubjects returns the factory of the subjects of Params that calls
// newSubject with the handle of each check.
func paramsSubjects[T Params](newSubject func(h *ParamsHandle) T) func(h *ParamsHandle) Params {
	return func(h *ParamsHandle) Params { return newSubject(h) }
}

// paramsBoolCalls generates the calls of Params.Bool with the inputs that
// each case draws.
var paramsBoolCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[bool](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Bool(v)
		return nil
	}
})

// paramsIntCalls generates the calls of Params.Int with the inputs that
// each case draws.
var paramsIntCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[int](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Int(v)
		return nil
	}
})

// paramsInt8Calls generates the calls of Params.Int8 with the inputs that
// each case draws.
var paramsInt8Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[int8](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Int8(v)
		return nil
	}
})

// paramsInt16Calls generates the calls of Params.Int16 with the inputs
// that each case draws.
var paramsInt16Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[int16](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Int16(v)
		return nil
	}
})

// paramsInt32Calls generates the calls of Params.Int32 with the inputs
// that each case draws.
var paramsInt32Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[int32](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Int32(v)
		return nil
	}
})

// paramsInt64Calls generates the calls of Params.Int64 with the inputs
// that each case draws.
var paramsInt64Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[int64](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Int64(v)
		return nil
	}
})

// paramsUintCalls generates the calls of Params.Uint with the inputs that
// each case draws.
var paramsUintCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[uint](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Uint(v)
		return nil
	}
})

// paramsUint8Calls generates the calls of Params.Uint8 with the inputs
// that each case draws.
var paramsUint8Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[uint8](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Uint8(v)
		return nil
	}
})

// paramsUint16Calls generates the calls of Params.Uint16 with the inputs
// that each case draws.
var paramsUint16Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[uint16](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Uint16(v)
		return nil
	}
})

// paramsUint32Calls generates the calls of Params.Uint32 with the inputs
// that each case draws.
var paramsUint32Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[uint32](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Uint32(v)
		return nil
	}
})

// paramsUint64Calls generates the calls of Params.Uint64 with the inputs
// that each case draws.
var paramsUint64Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[uint64](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Uint64(v)
		return nil
	}
})

// paramsUintptrCalls generates the calls of Params.Uintptr with the inputs
// that each case draws.
var paramsUintptrCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[uintptr](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Uintptr(v)
		return nil
	}
})

// paramsFloat32Calls generates the calls of Params.Float32 with the inputs
// that each case draws.
var paramsFloat32Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[float32](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Float32(v)
		return nil
	}
})

// paramsFloat64Calls generates the calls of Params.Float64 with the inputs
// that each case draws.
var paramsFloat64Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[float64](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Float64(v)
		return nil
	}
})

// paramsComplex64Calls generates the calls of Params.Complex64 with the
// inputs that each case draws.
var paramsComplex64Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[complex64](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Complex64(v)
		return nil
	}
})

// paramsComplex128Calls generates the calls of Params.Complex128 with the
// inputs that each case draws.
var paramsComplex128Calls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[complex128](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Complex128(v)
		return nil
	}
})

// paramsArrayCalls generates the calls of Params.Array with the inputs
// that each case draws.
var paramsArrayCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[[2]int](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Array(v)
		return nil
	}
})

// paramsChanCalls generates the calls of Params.Chan with the inputs that
// each case draws.
var paramsChanCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[chan int](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Chan(v)
		return nil
	}
})

// paramsFuncCalls generates the calls of Params.Func with the inputs that
// each case draws.
var paramsFuncCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	vResults := c.Draw(prop.List(prop.Of[int]()), "v")
	var vAt atomic.Int64
	v := func() (r0 int) {
		i := int(vAt.Add(1) - 1)
		if i >= len(vResults) {
			return r0
		}
		return vResults[i]
	}
	return func(_ context.Context, subject Params) error {
		subject.Func(v)
		return nil
	}
})

// paramsInterfaceCalls generates the calls of Params.Interface with the
// inputs that each case draws.
var paramsInterfaceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	vStringResults := c.Draw(prop.List(prop.Of[string]()), "v-string")
	var vStringAt atomic.Int64
	v := &paramsStringerInput{
		onString: func() (r0 string) {
			i := int(vStringAt.Add(1) - 1)
			if i >= len(vStringResults) {
				return r0
			}
			return vStringResults[i]
		},
	}
	return func(_ context.Context, subject Params) error {
		subject.Interface(v)
		return nil
	}
})

// paramsMapCalls generates the calls of Params.Map with the inputs that
// each case draws.
var paramsMapCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[map[string]int](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Map(v)
		return nil
	}
})

// paramsPtrCalls generates the calls of Params.Ptr with the inputs that
// each case draws.
var paramsPtrCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[int](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Ptr(&v)
		return nil
	}
})

// paramsSliceCalls generates the calls of Params.Slice with the inputs
// that each case draws.
var paramsSliceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[[]int](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Slice(v)
		return nil
	}
})

// paramsStringCalls generates the calls of Params.String with the inputs
// that each case draws.
var paramsStringCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[string](), "v")
	return func(_ context.Context, subject Params) error {
		subject.String(v)
		return nil
	}
})

// paramsStructCalls generates the calls of Params.Struct with the inputs
// that each case draws.
var paramsStructCalls = prop.Composite(func(c *prop.Case) func(context.Context, Params) error {
	v := c.Draw(prop.Of[struct{}](), "v")
	return func(_ context.Context, subject Params) error {
		subject.Struct(v)
		return nil
	}
})

// paramsUnsafePointerCalls generates the calls of Params.UnsafePointer
// with the inputs that each case draws.
var paramsUnsafePointerCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Params) error {
		v := c.Draw(prop.Of[unsafe.Pointer](), "v")
		return func(_ context.Context, subject Params) error {
			subject.UnsafePointer(v)
			return nil
		}
	},
)

// paramsStructValues generates the starting value of struct{}, which has
// no field to change.
var paramsStructValues = prop.Composite(func(*prop.Case) struct{} {
	return struct{}{}
})

// paramsStringerInputEmbedded is fmt.Stringer under the name of the field
// that paramsStringerInput embeds.
type paramsStringerInputEmbedded = fmt.Stringer

// paramsStringerInput implements fmt.Stringer. Its methods return the
// results that the generators of the suite of Params draw for them, in
// turn.
type paramsStringerInput struct {
	paramsStringerInputEmbedded
	onString func() string
}

// String calls the function of String.
func (in paramsStringerInput) String() string {
	return in.onString()
}

// paramsOverEmbedded is Params under the name of the field that paramsOver
// embeds.
type paramsOverEmbedded = Params

// paramsOver is a subject of Params whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type paramsOver struct {
	paramsOverEmbedded
	spend func()
}

// Bool calls spend and returns its zero results.
func (o paramsOver) Bool(bool) {
	o.spend()
}

// Int calls spend and returns its zero results.
func (o paramsOver) Int(int) {
	o.spend()
}

// Int8 calls spend and returns its zero results.
func (o paramsOver) Int8(int8) {
	o.spend()
}

// Int16 calls spend and returns its zero results.
func (o paramsOver) Int16(int16) {
	o.spend()
}

// Int32 calls spend and returns its zero results.
func (o paramsOver) Int32(int32) {
	o.spend()
}

// Int64 calls spend and returns its zero results.
func (o paramsOver) Int64(int64) {
	o.spend()
}

// Uint calls spend and returns its zero results.
func (o paramsOver) Uint(uint) {
	o.spend()
}

// Uint8 calls spend and returns its zero results.
func (o paramsOver) Uint8(uint8) {
	o.spend()
}

// Uint16 calls spend and returns its zero results.
func (o paramsOver) Uint16(uint16) {
	o.spend()
}

// Uint32 calls spend and returns its zero results.
func (o paramsOver) Uint32(uint32) {
	o.spend()
}

// Uint64 calls spend and returns its zero results.
func (o paramsOver) Uint64(uint64) {
	o.spend()
}

// Uintptr calls spend and returns its zero results.
func (o paramsOver) Uintptr(uintptr) {
	o.spend()
}

// Float32 calls spend and returns its zero results.
func (o paramsOver) Float32(float32) {
	o.spend()
}

// Float64 calls spend and returns its zero results.
func (o paramsOver) Float64(float64) {
	o.spend()
}

// Complex64 calls spend and returns its zero results.
func (o paramsOver) Complex64(complex64) {
	o.spend()
}

// Complex128 calls spend and returns its zero results.
func (o paramsOver) Complex128(complex128) {
	o.spend()
}

// Array calls spend and returns its zero results.
func (o paramsOver) Array([2]int) {
	o.spend()
}

// Chan calls spend and returns its zero results.
func (o paramsOver) Chan(chan int) {
	o.spend()
}

// Func calls spend and returns its zero results.
func (o paramsOver) Func(func() int) {
	o.spend()
}

// Interface calls spend and returns its zero results.
func (o paramsOver) Interface(fmt.Stringer) {
	o.spend()
}

// Map calls spend and returns its zero results.
func (o paramsOver) Map(map[string]int) {
	o.spend()
}

// Ptr calls spend and returns its zero results.
func (o paramsOver) Ptr(*int) {
	o.spend()
}

// Slice calls spend and returns its zero results.
func (o paramsOver) Slice([]int) {
	o.spend()
}

// String calls spend and returns its zero results.
func (o paramsOver) String(string) {
	o.spend()
}

// Struct calls spend and returns its zero results.
func (o paramsOver) Struct(struct{}) {
	o.spend()
}

// UnsafePointer calls spend and returns its zero results.
func (o paramsOver) UnsafePointer(unsafe.Pointer) {
	o.spend()
}

// newParamsOver returns the planted defect of a ceiling of Params, whose
// callables call spend.
func newParamsOver(_ *ParamsHandle, spend func()) Params {
	return paramsOver{spend: spend}
}

// paramsBoolPanicsDefect is the planted defect of a Params whose Bool
// panics.
type paramsBoolPanicsDefect struct {
	paramsOverEmbedded
}

// Bool is the override of the planted defect.
func (*paramsBoolPanicsDefect) Bool(p0 bool) {
	panic("the planted defect of a Params whose Bool panics")
}

// paramsBoolPanics returns the planted defect of a Params whose Bool
// panics over a fresh subject.
func paramsBoolPanics(h *ParamsHandle) Params {
	return &paramsBoolPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsIntPanicsDefect is the planted defect of a Params whose Int
// panics.
type paramsIntPanicsDefect struct {
	paramsOverEmbedded
}

// Int is the override of the planted defect.
func (*paramsIntPanicsDefect) Int(p0 int) {
	panic("the planted defect of a Params whose Int panics")
}

// paramsIntPanics returns the planted defect of a Params whose Int panics
// over a fresh subject.
func paramsIntPanics(h *ParamsHandle) Params {
	return &paramsIntPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsInt8PanicsDefect is the planted defect of a Params whose Int8
// panics.
type paramsInt8PanicsDefect struct {
	paramsOverEmbedded
}

// Int8 is the override of the planted defect.
func (*paramsInt8PanicsDefect) Int8(p0 int8) {
	panic("the planted defect of a Params whose Int8 panics")
}

// paramsInt8Panics returns the planted defect of a Params whose Int8
// panics over a fresh subject.
func paramsInt8Panics(h *ParamsHandle) Params {
	return &paramsInt8PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsInt16PanicsDefect is the planted defect of a Params whose Int16
// panics.
type paramsInt16PanicsDefect struct {
	paramsOverEmbedded
}

// Int16 is the override of the planted defect.
func (*paramsInt16PanicsDefect) Int16(p0 int16) {
	panic("the planted defect of a Params whose Int16 panics")
}

// paramsInt16Panics returns the planted defect of a Params whose Int16
// panics over a fresh subject.
func paramsInt16Panics(h *ParamsHandle) Params {
	return &paramsInt16PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsInt32PanicsDefect is the planted defect of a Params whose Int32
// panics.
type paramsInt32PanicsDefect struct {
	paramsOverEmbedded
}

// Int32 is the override of the planted defect.
func (*paramsInt32PanicsDefect) Int32(p0 int32) {
	panic("the planted defect of a Params whose Int32 panics")
}

// paramsInt32Panics returns the planted defect of a Params whose Int32
// panics over a fresh subject.
func paramsInt32Panics(h *ParamsHandle) Params {
	return &paramsInt32PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsInt64PanicsDefect is the planted defect of a Params whose Int64
// panics.
type paramsInt64PanicsDefect struct {
	paramsOverEmbedded
}

// Int64 is the override of the planted defect.
func (*paramsInt64PanicsDefect) Int64(p0 int64) {
	panic("the planted defect of a Params whose Int64 panics")
}

// paramsInt64Panics returns the planted defect of a Params whose Int64
// panics over a fresh subject.
func paramsInt64Panics(h *ParamsHandle) Params {
	return &paramsInt64PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsUintPanicsDefect is the planted defect of a Params whose Uint
// panics.
type paramsUintPanicsDefect struct {
	paramsOverEmbedded
}

// Uint is the override of the planted defect.
func (*paramsUintPanicsDefect) Uint(p0 uint) {
	panic("the planted defect of a Params whose Uint panics")
}

// paramsUintPanics returns the planted defect of a Params whose Uint
// panics over a fresh subject.
func paramsUintPanics(h *ParamsHandle) Params {
	return &paramsUintPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsUint8PanicsDefect is the planted defect of a Params whose Uint8
// panics.
type paramsUint8PanicsDefect struct {
	paramsOverEmbedded
}

// Uint8 is the override of the planted defect.
func (*paramsUint8PanicsDefect) Uint8(p0 uint8) {
	panic("the planted defect of a Params whose Uint8 panics")
}

// paramsUint8Panics returns the planted defect of a Params whose Uint8
// panics over a fresh subject.
func paramsUint8Panics(h *ParamsHandle) Params {
	return &paramsUint8PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsUint16PanicsDefect is the planted defect of a Params whose Uint16
// panics.
type paramsUint16PanicsDefect struct {
	paramsOverEmbedded
}

// Uint16 is the override of the planted defect.
func (*paramsUint16PanicsDefect) Uint16(p0 uint16) {
	panic("the planted defect of a Params whose Uint16 panics")
}

// paramsUint16Panics returns the planted defect of a Params whose Uint16
// panics over a fresh subject.
func paramsUint16Panics(h *ParamsHandle) Params {
	return &paramsUint16PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsUint32PanicsDefect is the planted defect of a Params whose Uint32
// panics.
type paramsUint32PanicsDefect struct {
	paramsOverEmbedded
}

// Uint32 is the override of the planted defect.
func (*paramsUint32PanicsDefect) Uint32(p0 uint32) {
	panic("the planted defect of a Params whose Uint32 panics")
}

// paramsUint32Panics returns the planted defect of a Params whose Uint32
// panics over a fresh subject.
func paramsUint32Panics(h *ParamsHandle) Params {
	return &paramsUint32PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsUint64PanicsDefect is the planted defect of a Params whose Uint64
// panics.
type paramsUint64PanicsDefect struct {
	paramsOverEmbedded
}

// Uint64 is the override of the planted defect.
func (*paramsUint64PanicsDefect) Uint64(p0 uint64) {
	panic("the planted defect of a Params whose Uint64 panics")
}

// paramsUint64Panics returns the planted defect of a Params whose Uint64
// panics over a fresh subject.
func paramsUint64Panics(h *ParamsHandle) Params {
	return &paramsUint64PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsUintptrPanicsDefect is the planted defect of a Params whose
// Uintptr panics.
type paramsUintptrPanicsDefect struct {
	paramsOverEmbedded
}

// Uintptr is the override of the planted defect.
func (*paramsUintptrPanicsDefect) Uintptr(p0 uintptr) {
	panic("the planted defect of a Params whose Uintptr panics")
}

// paramsUintptrPanics returns the planted defect of a Params whose Uintptr
// panics over a fresh subject.
func paramsUintptrPanics(h *ParamsHandle) Params {
	return &paramsUintptrPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsFloat32PanicsDefect is the planted defect of a Params whose
// Float32 panics.
type paramsFloat32PanicsDefect struct {
	paramsOverEmbedded
}

// Float32 is the override of the planted defect.
func (*paramsFloat32PanicsDefect) Float32(p0 float32) {
	panic("the planted defect of a Params whose Float32 panics")
}

// paramsFloat32Panics returns the planted defect of a Params whose Float32
// panics over a fresh subject.
func paramsFloat32Panics(h *ParamsHandle) Params {
	return &paramsFloat32PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsFloat64PanicsDefect is the planted defect of a Params whose
// Float64 panics.
type paramsFloat64PanicsDefect struct {
	paramsOverEmbedded
}

// Float64 is the override of the planted defect.
func (*paramsFloat64PanicsDefect) Float64(p0 float64) {
	panic("the planted defect of a Params whose Float64 panics")
}

// paramsFloat64Panics returns the planted defect of a Params whose Float64
// panics over a fresh subject.
func paramsFloat64Panics(h *ParamsHandle) Params {
	return &paramsFloat64PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsComplex64PanicsDefect is the planted defect of a Params whose
// Complex64 panics.
type paramsComplex64PanicsDefect struct {
	paramsOverEmbedded
}

// Complex64 is the override of the planted defect.
func (*paramsComplex64PanicsDefect) Complex64(p0 complex64) {
	panic("the planted defect of a Params whose Complex64 panics")
}

// paramsComplex64Panics returns the planted defect of a Params whose
// Complex64 panics over a fresh subject.
func paramsComplex64Panics(h *ParamsHandle) Params {
	return &paramsComplex64PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsComplex128PanicsDefect is the planted defect of a Params whose
// Complex128 panics.
type paramsComplex128PanicsDefect struct {
	paramsOverEmbedded
}

// Complex128 is the override of the planted defect.
func (*paramsComplex128PanicsDefect) Complex128(p0 complex128) {
	panic("the planted defect of a Params whose Complex128 panics")
}

// paramsComplex128Panics returns the planted defect of a Params whose
// Complex128 panics over a fresh subject.
func paramsComplex128Panics(h *ParamsHandle) Params {
	return &paramsComplex128PanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsArrayPanicsDefect is the planted defect of a Params whose Array
// panics.
type paramsArrayPanicsDefect struct {
	paramsOverEmbedded
}

// Array is the override of the planted defect.
func (*paramsArrayPanicsDefect) Array(p0 [2]int) {
	panic("the planted defect of a Params whose Array panics")
}

// paramsArrayPanics returns the planted defect of a Params whose Array
// panics over a fresh subject.
func paramsArrayPanics(h *ParamsHandle) Params {
	return &paramsArrayPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsFuncPanicsDefect is the planted defect of a Params whose Func
// panics.
type paramsFuncPanicsDefect struct {
	paramsOverEmbedded
}

// Func is the override of the planted defect.
func (*paramsFuncPanicsDefect) Func(p0 func() int) {
	panic("the planted defect of a Params whose Func panics")
}

// paramsFuncPanics returns the planted defect of a Params whose Func
// panics over a fresh subject.
func paramsFuncPanics(h *ParamsHandle) Params {
	return &paramsFuncPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsMapPanicsDefect is the planted defect of a Params whose Map
// panics.
type paramsMapPanicsDefect struct {
	paramsOverEmbedded
}

// Map is the override of the planted defect.
func (*paramsMapPanicsDefect) Map(p0 map[string]int) {
	panic("the planted defect of a Params whose Map panics")
}

// paramsMapPanics returns the planted defect of a Params whose Map panics
// over a fresh subject.
func paramsMapPanics(h *ParamsHandle) Params {
	return &paramsMapPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsSlicePanicsDefect is the planted defect of a Params whose Slice
// panics.
type paramsSlicePanicsDefect struct {
	paramsOverEmbedded
}

// Slice is the override of the planted defect.
func (*paramsSlicePanicsDefect) Slice(p0 []int) {
	panic("the planted defect of a Params whose Slice panics")
}

// paramsSlicePanics returns the planted defect of a Params whose Slice
// panics over a fresh subject.
func paramsSlicePanics(h *ParamsHandle) Params {
	return &paramsSlicePanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// paramsStringPanicsDefect is the planted defect of a Params whose String
// panics.
type paramsStringPanicsDefect struct {
	paramsOverEmbedded
}

// String is the override of the planted defect.
func (*paramsStringPanicsDefect) String(p0 string) {
	panic("the planted defect of a Params whose String panics")
}

// paramsStringPanics returns the planted defect of a Params whose String
// panics over a fresh subject.
func paramsStringPanics(h *ParamsHandle) Params {
	return &paramsStringPanicsDefect{paramsOverEmbedded: h.Fresh()}
}

// The derived checks of Every:
//
// Every.Fields/smoke: Fields returns without a panic for derived inputs.
//
// Every.Elements/smoke: Elements returns without a panic for derived
// inputs.

// The derived checks of Every that witness their laws at the shallow end:
//
// Every.Fields/smoke
//
// Every.Elements/smoke

// The unclaimed classifications of Every:
//
// The shape:mutator of Every.Fields has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Every.Elements has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Every.Results has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Every writes.
type (
	EveryHandle = conform.Handle[Every]
	EveryCheck  = conform.Check[Every]
	EveryOption = conform.Option[Every]
)

// EveryChecks contains the name of each callable of Every, with the
// identities of its checks.
var EveryChecks = struct {
	Fields   everyFields
	Elements everyElements
	Results  everyResults
}{
	Fields:   everyFields{"Every.Fields"},
	Elements: everyElements{"Every.Elements"},
	Results:  everyResults{"Every.Results"},
}

// everyFields is the name of Every.Fields, with the identities of its
// checks.
type everyFields struct {
	conform.Callable[Every]
}

// Smoke returns the identity of the check that Fields returns without a
// panic for derived inputs.
func (c everyFields) Smoke() conform.ID[Every] {
	return c.Check("smoke")
}

// everyElements is the name of Every.Elements, with the identities of its
// checks.
type everyElements struct {
	conform.Callable[Every]
}

// Smoke returns the identity of the check that Elements returns without a
// panic for derived inputs.
func (c everyElements) Smoke() conform.ID[Every] {
	return c.Check("smoke")
}

// everyResults is the name of Every.Results, with the identities of its
// checks.
type everyResults struct {
	conform.Callable[Every]
}

// EverySuite returns the checks of Every, and the calls and the budgets of
// its callables.
func EverySuite() conform.Suite[Every] {
	return conform.Suite[Every]{
		Declaration: "Every",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Every]{
			{
				ID:          EveryChecks.Fields.Smoke(),
				Statement:   "Fields returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *EveryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Fields(Fields{Bool: true})
						return nil
					})
				},
				Defect: everyFieldsPanics,
			},
			{
				ID:          EveryChecks.Elements.Smoke(),
				Statement:   "Elements returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *EveryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Elements(Elements{Bool: []bool{true}})
						return nil
					})
				},
				Defect: everyElementsPanics,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:mutator",
				Callable:       "Every.Fields",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Every.Elements",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Every.Results",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Every]{
			{
				Calls:    everyFieldsCalls,
				Over:     newEveryOver,
				Callable: EveryChecks.Fields.Callable,
			},
			{
				Calls:    everyElementsCalls,
				Over:     newEveryOver,
				Callable: EveryChecks.Elements.Callable,
			},
			{
				Calls:    everyResultsCalls,
				Over:     newEveryOver,
				Callable: EveryChecks.Results.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(everyStructValues),
			prop.Using(prop.Composite(func(c *prop.Case) func() int {
				funcResults := c.Draw(prop.List(prop.Of[int]()), "func")
				var funcAt atomic.Int64
				return func() (r0 int) {
					i := int(funcAt.Add(1) - 1)
					if i >= len(funcResults) {
						return r0
					}
					return funcResults[i]
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) fmt.Stringer {
				stringerStringResults := c.Draw(prop.List(prop.Of[string]()), "stringer-string")
				var stringerStringAt atomic.Int64
				return &everyStringerInput{
					onString: func() (r0 string) {
						i := int(stringerStringAt.Add(1) - 1)
						if i >= len(stringerStringResults) {
							return r0
						}
						return stringerStringResults[i]
					},
				}
			})),
			prop.Using(conform.Stream[chan int](prop.Of[int]())),
			prop.Using(prop.Composite(func(*prop.Case) unsafe.Pointer { return nil })),
		},
	}
}

// RunEverySuite applies the checks of Every, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunEverySuite[T Every](t *testing.T, newSubject func() T, opts ...EveryOption) {
	t.Helper()
	RunEverySuiteFromHandle(t, func(*EveryHandle) T { return newSubject() }, opts...)
}

// RunEverySuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunEverySuiteFromHandle[T Every](
	t *testing.T,
	newSubject func(h *EveryHandle) T,
	opts ...EveryOption,
) {
	t.Helper()
	conform.Run(t, EverySuite(), everySubjects(newSubject), opts...)
}

// FuzzEverySuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzEverySuite[T Every](
	f *testing.F,
	callable conform.Callable[Every],
	newSubject func(h *EveryHandle) T,
	opts ...EveryOption,
) {
	f.Helper()
	conform.Fuzz(f, EverySuite(), callable, everySubjects(newSubject), opts...)
}

// RunEveryBudgets applies the allocation ceilings of Every and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunEveryBudgets[T Every](
	t *testing.T,
	newSubject func(h *EveryHandle) T,
	opts ...EveryOption,
) {
	t.Helper()
	conform.RunBudgets(t, EverySuite(), everySubjects(newSubject), opts...)
}

// BenchEverySuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchEverySuite[T Every](
	b conform.B,
	callable conform.Callable[Every],
	newSubject func(h *EveryHandle) T,
	opts ...EveryOption,
) {
	b.Helper()
	conform.Bench(b, EverySuite(), callable, everySubjects(newSubject), opts...)
}

// everySubjects returns the factory of the subjects of Every that calls
// newSubject with the handle of each check.
func everySubjects[T Every](newSubject func(h *EveryHandle) T) func(h *EveryHandle) Every {
	return func(h *EveryHandle) Every { return newSubject(h) }
}

// everyFieldsCalls generates the calls of Every.Fields with the inputs
// that each case draws.
var everyFieldsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Every) error {
	f := c.Draw(prop.Of[Fields](), "f")
	return func(_ context.Context, subject Every) error {
		subject.Fields(f)
		return nil
	}
})

// everyElementsCalls generates the calls of Every.Elements with the inputs
// that each case draws.
var everyElementsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Every) error {
	e := c.Draw(prop.Of[Elements](), "e")
	return func(_ context.Context, subject Every) error {
		subject.Elements(e)
		return nil
	}
})

// everyResultsCalls generates the calls of Every.Results with the inputs
// that each case draws.
var everyResultsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Every) error {
	rBoolResults := c.Draw(prop.List(prop.Of[bool]()), "r-bool")
	var rBoolAt atomic.Int64
	rIntResults := c.Draw(prop.List(prop.Of[int]()), "r-int")
	var rIntAt atomic.Int64
	rInt8Results := c.Draw(prop.List(prop.Of[int8]()), "r-int8")
	var rInt8At atomic.Int64
	rInt16Results := c.Draw(prop.List(prop.Of[int16]()), "r-int16")
	var rInt16At atomic.Int64
	rInt32Results := c.Draw(prop.List(prop.Of[int32]()), "r-int32")
	var rInt32At atomic.Int64
	rInt64Results := c.Draw(prop.List(prop.Of[int64]()), "r-int64")
	var rInt64At atomic.Int64
	rUintResults := c.Draw(prop.List(prop.Of[uint]()), "r-uint")
	var rUintAt atomic.Int64
	rUint8Results := c.Draw(prop.List(prop.Of[uint8]()), "r-uint8")
	var rUint8At atomic.Int64
	rUint16Results := c.Draw(prop.List(prop.Of[uint16]()), "r-uint16")
	var rUint16At atomic.Int64
	rUint32Results := c.Draw(prop.List(prop.Of[uint32]()), "r-uint32")
	var rUint32At atomic.Int64
	rUint64Results := c.Draw(prop.List(prop.Of[uint64]()), "r-uint64")
	var rUint64At atomic.Int64
	rUintptrResults := c.Draw(prop.List(prop.Of[uintptr]()), "r-uintptr")
	var rUintptrAt atomic.Int64
	rFloat32Results := c.Draw(prop.List(prop.Of[float32]()), "r-float32")
	var rFloat32At atomic.Int64
	rFloat64Results := c.Draw(prop.List(prop.Of[float64]()), "r-float64")
	var rFloat64At atomic.Int64
	rComplex64Results := c.Draw(prop.List(prop.Of[complex64]()), "r-complex64")
	var rComplex64At atomic.Int64
	rComplex128Results := c.Draw(prop.List(prop.Of[complex128]()), "r-complex128")
	var rComplex128At atomic.Int64
	rArrayResults := c.Draw(prop.List(prop.Of[[2]int]()), "r-array")
	var rArrayAt atomic.Int64
	rChanResults := c.Draw(prop.List(prop.Of[chan int]()), "r-chan")
	rChanEnd := conform.Closed[chan int]()
	var rChanAt atomic.Int64
	rFuncResults := c.Draw(prop.List(prop.Of[func() int]()), "r-func")
	var rFuncAt atomic.Int64
	rInterfaceResults := c.Draw(prop.List(prop.Of[fmt.Stringer]()), "r-interface")
	var rInterfaceAt atomic.Int64
	rMapResults := c.Draw(prop.List(prop.Of[map[string]int]()), "r-map")
	var rMapAt atomic.Int64
	rPtrResults := c.Draw(prop.List(prop.Of[*int]()), "r-ptr")
	var rPtrAt atomic.Int64
	rSliceResults := c.Draw(prop.List(prop.Of[[]int]()), "r-slice")
	var rSliceAt atomic.Int64
	rStringResults := c.Draw(prop.List(prop.Of[string]()), "r-string")
	var rStringAt atomic.Int64
	rStructResults := c.Draw(prop.List(prop.Of[struct{}]()), "r-struct")
	var rStructAt atomic.Int64
	rUnsafePointerResults := c.Draw(prop.List(prop.Of[unsafe.Pointer]()), "r-unsafe-pointer")
	var rUnsafePointerAt atomic.Int64
	r := &everyResultsInput{
		onBool: func() (r0 bool) {
			i := int(rBoolAt.Add(1) - 1)
			if i >= len(rBoolResults) {
				return r0
			}
			return rBoolResults[i]
		},
		onInt: func() (r0 int) {
			i := int(rIntAt.Add(1) - 1)
			if i >= len(rIntResults) {
				return r0
			}
			return rIntResults[i]
		},
		onInt8: func() (r0 int8) {
			i := int(rInt8At.Add(1) - 1)
			if i >= len(rInt8Results) {
				return r0
			}
			return rInt8Results[i]
		},
		onInt16: func() (r0 int16) {
			i := int(rInt16At.Add(1) - 1)
			if i >= len(rInt16Results) {
				return r0
			}
			return rInt16Results[i]
		},
		onInt32: func() (r0 int32) {
			i := int(rInt32At.Add(1) - 1)
			if i >= len(rInt32Results) {
				return r0
			}
			return rInt32Results[i]
		},
		onInt64: func() (r0 int64) {
			i := int(rInt64At.Add(1) - 1)
			if i >= len(rInt64Results) {
				return r0
			}
			return rInt64Results[i]
		},
		onUint: func() (r0 uint) {
			i := int(rUintAt.Add(1) - 1)
			if i >= len(rUintResults) {
				return r0
			}
			return rUintResults[i]
		},
		onUint8: func() (r0 uint8) {
			i := int(rUint8At.Add(1) - 1)
			if i >= len(rUint8Results) {
				return r0
			}
			return rUint8Results[i]
		},
		onUint16: func() (r0 uint16) {
			i := int(rUint16At.Add(1) - 1)
			if i >= len(rUint16Results) {
				return r0
			}
			return rUint16Results[i]
		},
		onUint32: func() (r0 uint32) {
			i := int(rUint32At.Add(1) - 1)
			if i >= len(rUint32Results) {
				return r0
			}
			return rUint32Results[i]
		},
		onUint64: func() (r0 uint64) {
			i := int(rUint64At.Add(1) - 1)
			if i >= len(rUint64Results) {
				return r0
			}
			return rUint64Results[i]
		},
		onUintptr: func() (r0 uintptr) {
			i := int(rUintptrAt.Add(1) - 1)
			if i >= len(rUintptrResults) {
				return r0
			}
			return rUintptrResults[i]
		},
		onFloat32: func() (r0 float32) {
			i := int(rFloat32At.Add(1) - 1)
			if i >= len(rFloat32Results) {
				return r0
			}
			return rFloat32Results[i]
		},
		onFloat64: func() (r0 float64) {
			i := int(rFloat64At.Add(1) - 1)
			if i >= len(rFloat64Results) {
				return r0
			}
			return rFloat64Results[i]
		},
		onComplex64: func() (r0 complex64) {
			i := int(rComplex64At.Add(1) - 1)
			if i >= len(rComplex64Results) {
				return r0
			}
			return rComplex64Results[i]
		},
		onComplex128: func() (r0 complex128) {
			i := int(rComplex128At.Add(1) - 1)
			if i >= len(rComplex128Results) {
				return r0
			}
			return rComplex128Results[i]
		},
		onArray: func() (r0 [2]int) {
			i := int(rArrayAt.Add(1) - 1)
			if i >= len(rArrayResults) {
				return r0
			}
			return rArrayResults[i]
		},
		onChan: func() (r0 chan int) {
			i := int(rChanAt.Add(1) - 1)
			if i >= len(rChanResults) {
				r0 = rChanEnd
				return r0
			}
			return rChanResults[i]
		},
		onFunc: func() (r0 func() int) {
			i := int(rFuncAt.Add(1) - 1)
			if i >= len(rFuncResults) {
				return r0
			}
			return rFuncResults[i]
		},
		onInterface: func() (r0 fmt.Stringer) {
			i := int(rInterfaceAt.Add(1) - 1)
			if i >= len(rInterfaceResults) {
				return r0
			}
			return rInterfaceResults[i]
		},
		onMap: func() (r0 map[string]int) {
			i := int(rMapAt.Add(1) - 1)
			if i >= len(rMapResults) {
				return r0
			}
			return rMapResults[i]
		},
		onPtr: func() (r0 *int) {
			i := int(rPtrAt.Add(1) - 1)
			if i >= len(rPtrResults) {
				return r0
			}
			return rPtrResults[i]
		},
		onSlice: func() (r0 []int) {
			i := int(rSliceAt.Add(1) - 1)
			if i >= len(rSliceResults) {
				return r0
			}
			return rSliceResults[i]
		},
		onString: func() (r0 string) {
			i := int(rStringAt.Add(1) - 1)
			if i >= len(rStringResults) {
				return r0
			}
			return rStringResults[i]
		},
		onStruct: func() (r0 struct{}) {
			i := int(rStructAt.Add(1) - 1)
			if i >= len(rStructResults) {
				return r0
			}
			return rStructResults[i]
		},
		onUnsafePointer: func() (r0 unsafe.Pointer) {
			i := int(rUnsafePointerAt.Add(1) - 1)
			if i >= len(rUnsafePointerResults) {
				return r0
			}
			return rUnsafePointerResults[i]
		},
	}
	return func(_ context.Context, subject Every) error {
		subject.Results(r)
		return nil
	}
})

// everyStructValues generates the starting value of struct{}, which has no
// field to change.
var everyStructValues = prop.Composite(func(*prop.Case) struct{} {
	return struct{}{}
})

// everyResultsInputEmbedded is Results under the name of the field that
// everyResultsInput embeds.
type everyResultsInputEmbedded = Results

// everyResultsInput implements Results. Its methods return the results
// that the generators of the suite of Every draw for them, in turn.
type everyResultsInput struct {
	everyResultsInputEmbedded
	onBool          func() bool
	onInt           func() int
	onInt8          func() int8
	onInt16         func() int16
	onInt32         func() int32
	onInt64         func() int64
	onUint          func() uint
	onUint8         func() uint8
	onUint16        func() uint16
	onUint32        func() uint32
	onUint64        func() uint64
	onUintptr       func() uintptr
	onFloat32       func() float32
	onFloat64       func() float64
	onComplex64     func() complex64
	onComplex128    func() complex128
	onArray         func() [2]int
	onChan          func() chan int
	onFunc          func() func() int
	onInterface     func() fmt.Stringer
	onMap           func() map[string]int
	onPtr           func() *int
	onSlice         func() []int
	onString        func() string
	onStruct        func() struct{}
	onUnsafePointer func() unsafe.Pointer
}

// Bool calls the function of Bool.
func (in everyResultsInput) Bool() bool {
	return in.onBool()
}

// Int calls the function of Int.
func (in everyResultsInput) Int() int {
	return in.onInt()
}

// Int8 calls the function of Int8.
func (in everyResultsInput) Int8() int8 {
	return in.onInt8()
}

// Int16 calls the function of Int16.
func (in everyResultsInput) Int16() int16 {
	return in.onInt16()
}

// Int32 calls the function of Int32.
func (in everyResultsInput) Int32() int32 {
	return in.onInt32()
}

// Int64 calls the function of Int64.
func (in everyResultsInput) Int64() int64 {
	return in.onInt64()
}

// Uint calls the function of Uint.
func (in everyResultsInput) Uint() uint {
	return in.onUint()
}

// Uint8 calls the function of Uint8.
func (in everyResultsInput) Uint8() uint8 {
	return in.onUint8()
}

// Uint16 calls the function of Uint16.
func (in everyResultsInput) Uint16() uint16 {
	return in.onUint16()
}

// Uint32 calls the function of Uint32.
func (in everyResultsInput) Uint32() uint32 {
	return in.onUint32()
}

// Uint64 calls the function of Uint64.
func (in everyResultsInput) Uint64() uint64 {
	return in.onUint64()
}

// Uintptr calls the function of Uintptr.
func (in everyResultsInput) Uintptr() uintptr {
	return in.onUintptr()
}

// Float32 calls the function of Float32.
func (in everyResultsInput) Float32() float32 {
	return in.onFloat32()
}

// Float64 calls the function of Float64.
func (in everyResultsInput) Float64() float64 {
	return in.onFloat64()
}

// Complex64 calls the function of Complex64.
func (in everyResultsInput) Complex64() complex64 {
	return in.onComplex64()
}

// Complex128 calls the function of Complex128.
func (in everyResultsInput) Complex128() complex128 {
	return in.onComplex128()
}

// Array calls the function of Array.
func (in everyResultsInput) Array() [2]int {
	return in.onArray()
}

// Chan calls the function of Chan.
func (in everyResultsInput) Chan() chan int {
	return in.onChan()
}

// Func calls the function of Func.
func (in everyResultsInput) Func() func() int {
	return in.onFunc()
}

// Interface calls the function of Interface.
func (in everyResultsInput) Interface() fmt.Stringer {
	return in.onInterface()
}

// Map calls the function of Map.
func (in everyResultsInput) Map() map[string]int {
	return in.onMap()
}

// Ptr calls the function of Ptr.
func (in everyResultsInput) Ptr() *int {
	return in.onPtr()
}

// Slice calls the function of Slice.
func (in everyResultsInput) Slice() []int {
	return in.onSlice()
}

// String calls the function of String.
func (in everyResultsInput) String() string {
	return in.onString()
}

// Struct calls the function of Struct.
func (in everyResultsInput) Struct() struct{} {
	return in.onStruct()
}

// UnsafePointer calls the function of UnsafePointer.
func (in everyResultsInput) UnsafePointer() unsafe.Pointer {
	return in.onUnsafePointer()
}

// everyStringerInputEmbedded is fmt.Stringer under the name of the field
// that everyStringerInput embeds.
type everyStringerInputEmbedded = fmt.Stringer

// everyStringerInput implements fmt.Stringer. Its methods return the
// results that the generators of the suite of Every draw for them, in
// turn.
type everyStringerInput struct {
	everyStringerInputEmbedded
	onString func() string
}

// String calls the function of String.
func (in everyStringerInput) String() string {
	return in.onString()
}

// everyOverEmbedded is Every under the name of the field that everyOver
// embeds.
type everyOverEmbedded = Every

// everyOver is a subject of Every whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type everyOver struct {
	everyOverEmbedded
	spend func()
}

// Fields calls spend and returns its zero results.
func (o everyOver) Fields(Fields) {
	o.spend()
}

// Elements calls spend and returns its zero results.
func (o everyOver) Elements(Elements) {
	o.spend()
}

// Results calls spend and returns its zero results.
func (o everyOver) Results(Results) {
	o.spend()
}

// newEveryOver returns the planted defect of a ceiling of Every, whose
// callables call spend.
func newEveryOver(_ *EveryHandle, spend func()) Every {
	return everyOver{spend: spend}
}

// everyFieldsPanicsDefect is the planted defect of a Every whose Fields
// panics.
type everyFieldsPanicsDefect struct {
	everyOverEmbedded
}

// Fields is the override of the planted defect.
func (*everyFieldsPanicsDefect) Fields(p0 Fields) {
	panic("the planted defect of a Every whose Fields panics")
}

// everyFieldsPanics returns the planted defect of a Every whose Fields
// panics over a fresh subject.
func everyFieldsPanics(h *EveryHandle) Every {
	return &everyFieldsPanicsDefect{everyOverEmbedded: h.Fresh()}
}

// everyElementsPanicsDefect is the planted defect of a Every whose
// Elements panics.
type everyElementsPanicsDefect struct {
	everyOverEmbedded
}

// Elements is the override of the planted defect.
func (*everyElementsPanicsDefect) Elements(p0 Elements) {
	panic("the planted defect of a Every whose Elements panics")
}

// everyElementsPanics returns the planted defect of a Every whose Elements
// panics over a fresh subject.
func everyElementsPanics(h *EveryHandle) Every {
	return &everyElementsPanicsDefect{everyOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Placer writes.
type (
	PlacerHandle = conform.Handle[Placer]
	PlacerCheck  = conform.Check[Placer]
	PlacerOption = conform.Option[Placer]
)

// PlacerChecks contains the name of each callable of Placer, with the
// identities of its checks.
var PlacerChecks = struct {
	Place placerPlace
}{
	Place: placerPlace{"Placer.Place"},
}

// placerPlace is the name of Placer.Place, with the identities of its
// checks.
type placerPlace struct {
	conform.Callable[Placer]
}

// PlacerSuite returns the checks of Placer, and the calls and the budgets
// of its callables.
func PlacerSuite() conform.Suite[Placer] {
	return conform.Suite[Placer]{
		Declaration: "Placer",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Calls: []conform.Call[Placer]{
			{
				Calls:    placerPlaceCalls,
				Over:     newPlacerOver,
				Callable: PlacerChecks.Place.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) Timer {
				timerCResults := c.Draw(prop.List(prop.Of[<-chan time.Time]()), "timer-c")
				timerCEnd := conform.Closed[<-chan time.Time]()
				var timerCAt atomic.Int64
				timerStopResults := c.Draw(prop.List(prop.Of[bool]()), "timer-stop")
				var timerStopAt atomic.Int64
				return &placerTimerInput{
					onC: func() (r0 <-chan time.Time) {
						i := int(timerCAt.Add(1) - 1)
						if i >= len(timerCResults) {
							r0 = timerCEnd
							return r0
						}
						return timerCResults[i]
					},
					onStop: func() (r0 bool) {
						i := int(timerStopAt.Add(1) - 1)
						if i >= len(timerStopResults) {
							return r0
						}
						return timerStopResults[i]
					},
				}
			})),
			prop.Using(conform.Stream[<-chan time.Time](prop.Of[time.Time]())),
		},
	}
}

// RunPlacerSuite applies the checks of Placer, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunPlacerSuite[T Placer](t *testing.T, newSubject func() T, opts ...PlacerOption) {
	t.Helper()
	RunPlacerSuiteFromHandle(t, func(*PlacerHandle) T { return newSubject() }, opts...)
}

// RunPlacerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPlacerSuiteFromHandle[T Placer](
	t *testing.T,
	newSubject func(h *PlacerHandle) T,
	opts ...PlacerOption,
) {
	t.Helper()
	conform.Run(t, PlacerSuite(), placerSubjects(newSubject), opts...)
}

// FuzzPlacerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPlacerSuite[T Placer](
	f *testing.F,
	callable conform.Callable[Placer],
	newSubject func(h *PlacerHandle) T,
	opts ...PlacerOption,
) {
	f.Helper()
	conform.Fuzz(f, PlacerSuite(), callable, placerSubjects(newSubject), opts...)
}

// RunPlacerBudgets applies the allocation ceilings of Placer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunPlacerBudgets[T Placer](
	t *testing.T,
	newSubject func(h *PlacerHandle) T,
	opts ...PlacerOption,
) {
	t.Helper()
	conform.RunBudgets(t, PlacerSuite(), placerSubjects(newSubject), opts...)
}

// BenchPlacerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPlacerSuite[T Placer](
	b conform.B,
	callable conform.Callable[Placer],
	newSubject func(h *PlacerHandle) T,
	opts ...PlacerOption,
) {
	b.Helper()
	conform.Bench(b, PlacerSuite(), callable, placerSubjects(newSubject), opts...)
}

// placerSubjects returns the factory of the subjects of Placer that calls
// newSubject with the handle of each check.
func placerSubjects[T Placer](newSubject func(h *PlacerHandle) T) func(h *PlacerHandle) Placer {
	return func(h *PlacerHandle) Placer { return newSubject(h) }
}

// placerPlaceCalls generates the calls of Placer.Place with the inputs
// that each case draws.
var placerPlaceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Placer) error {
	clockNowResults := c.Draw(prop.List(prop.Of[time.Time]()), "clock-now")
	var clockNowAt atomic.Int64
	clockNewTimerResults := c.Draw(prop.List(prop.Of[Timer]()), "clock-new-timer")
	var clockNewTimerAt atomic.Int64
	clock := &placerClockInput{
		onNow: func() (r0 time.Time) {
			i := int(clockNowAt.Add(1) - 1)
			if i >= len(clockNowResults) {
				return r0
			}
			return clockNowResults[i]
		},
		onNewTimer: func(time.Duration) (r0 Timer) {
			i := int(clockNewTimerAt.Add(1) - 1)
			if i >= len(clockNewTimerResults) {
				return r0
			}
			return clockNewTimerResults[i]
		},
	}
	key := c.Draw(prop.Of[string](), "key")
	return func(_ context.Context, subject Placer) error {
		_, err := subject.Place(clock, key)
		return err
	}
})

// placerClockInputEmbedded is Clock under the name of the field that
// placerClockInput embeds.
type placerClockInputEmbedded = Clock

// placerClockInput implements Clock. Its methods return the results that
// the generators of the suite of Placer draw for them, in turn.
type placerClockInput struct {
	placerClockInputEmbedded
	onNow      func() time.Time
	onNewTimer func(a0 time.Duration) Timer
}

// Now calls the function of Now.
func (in placerClockInput) Now() time.Time {
	return in.onNow()
}

// NewTimer calls the function of NewTimer.
func (in placerClockInput) NewTimer(a0 time.Duration) Timer {
	return in.onNewTimer(a0)
}

// placerTimerInputEmbedded is Timer under the name of the field that
// placerTimerInput embeds.
type placerTimerInputEmbedded = Timer

// placerTimerInput implements Timer. Its methods return the results that
// the generators of the suite of Placer draw for them, in turn.
type placerTimerInput struct {
	placerTimerInputEmbedded
	onC    func() <-chan time.Time
	onStop func() bool
}

// C calls the function of C.
func (in placerTimerInput) C() <-chan time.Time {
	return in.onC()
}

// Stop calls the function of Stop.
func (in placerTimerInput) Stop() bool {
	return in.onStop()
}

// placerOverEmbedded is Placer under the name of the field that placerOver
// embeds.
type placerOverEmbedded = Placer

// placerOver is a subject of Placer whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type placerOver struct {
	placerOverEmbedded
	spend func()
}

// Place calls spend and returns its zero results.
func (o placerOver) Place(Clock, string) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newPlacerOver returns the planted defect of a ceiling of Placer, whose
// callables call spend.
func newPlacerOver(_ *PlacerHandle, spend func()) Placer {
	return placerOver{spend: spend}
}

// The refused claims of Router:
//
// Router.Route/cancel: eidos derives no sample of the parameter f of
// Route. The remedy is to author a sample of the parameter f with the
// directive sample.
//
// Router.Read/deterministic: eidos derives no sample of the parameter in
// of Read. The remedy is to author a sample of the parameter in with the
// directive sample.

// The unclaimed classifications of Router:
//
// The shape:writer of Router.Route has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Router.Write has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Router.Pass has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Router.Watch has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Router writes.
type (
	RouterHandle = conform.Handle[Router]
	RouterCheck  = conform.Check[Router]
	RouterOption = conform.Option[Router]
)

// RouterChecks contains the name of each callable of Router, with the
// identities of its checks.
var RouterChecks = struct {
	Route routerRoute
	Read  routerRead
	Write routerWrite
	Pass  routerPass
	Watch routerWatch
}{
	Route: routerRoute{"Router.Route"},
	Read:  routerRead{"Router.Read"},
	Write: routerWrite{"Router.Write"},
	Pass:  routerPass{"Router.Pass"},
	Watch: routerWatch{"Router.Watch"},
}

// routerRoute is the name of Router.Route, with the identities of its
// checks.
type routerRoute struct {
	conform.Callable[Router]
}

// routerRead is the name of Router.Read, with the identities of its
// checks.
type routerRead struct {
	conform.Callable[Router]
}

// routerWrite is the name of Router.Write, with the identities of its
// checks.
type routerWrite struct {
	conform.Callable[Router]
}

// routerPass is the name of Router.Pass, with the identities of its
// checks.
type routerPass struct {
	conform.Callable[Router]
}

// routerWatch is the name of Router.Watch, with the identities of its
// checks.
type routerWatch struct {
	conform.Callable[Router]
}

// RouterSuite returns the checks of Router, and the calls and the budgets
// of its callables.
func RouterSuite() conform.Suite[Router] {
	return conform.Suite[Router]{
		Declaration: "Router",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Router.Route",
				Why:      "eidos derives no sample of the parameter f of Route",
				Remedy:   "author a sample of the parameter f with the directive sample",
			},
			{
				Claim:    spec.Deterministic,
				Callable: "Router.Read",
				Why:      "eidos derives no sample of the parameter in of Read",
				Remedy:   "author a sample of the parameter in with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Router.Route",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Router.Write",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Router.Pass",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Router.Watch",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Router]{
			{
				Calls:    routerRouteCalls,
				Over:     newRouterOver,
				Callable: RouterChecks.Route.Callable,
			},
			{
				Calls:    routerReadCalls,
				Over:     newRouterOver,
				Callable: RouterChecks.Read.Callable,
			},
			{
				Calls:    routerWriteCalls,
				Over:     newRouterOver,
				Callable: RouterChecks.Write.Callable,
			},
			{
				Calls:    routerPassCalls,
				Over:     newRouterOver,
				Callable: RouterChecks.Pass.Callable,
			},
			{
				Calls:    routerWatchCalls,
				Over:     newRouterOver,
				Callable: RouterChecks.Watch.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(conform.Stream[<-chan string](prop.Of[string]())),
			prop.Using(conform.Stream[chan string](prop.Of[string]())),
			prop.Using(conform.Stream[Ticks](prop.Of[int]())),
		},
		Sinks: []conform.Sink{
			conform.SinkOf[chan<- string](),
		},
	}
}

// RunRouterSuite applies the checks of Router, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunRouterSuite[T Router](t *testing.T, newSubject func() T, opts ...RouterOption) {
	t.Helper()
	RunRouterSuiteFromHandle(t, func(*RouterHandle) T { return newSubject() }, opts...)
}

// RunRouterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRouterSuiteFromHandle[T Router](
	t *testing.T,
	newSubject func(h *RouterHandle) T,
	opts ...RouterOption,
) {
	t.Helper()
	conform.Run(t, RouterSuite(), routerSubjects(newSubject), opts...)
}

// FuzzRouterSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRouterSuite[T Router](
	f *testing.F,
	callable conform.Callable[Router],
	newSubject func(h *RouterHandle) T,
	opts ...RouterOption,
) {
	f.Helper()
	conform.Fuzz(f, RouterSuite(), callable, routerSubjects(newSubject), opts...)
}

// RunRouterBudgets applies the allocation ceilings of Router and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRouterBudgets[T Router](
	t *testing.T,
	newSubject func(h *RouterHandle) T,
	opts ...RouterOption,
) {
	t.Helper()
	conform.RunBudgets(t, RouterSuite(), routerSubjects(newSubject), opts...)
}

// BenchRouterSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchRouterSuite[T Router](
	b conform.B,
	callable conform.Callable[Router],
	newSubject func(h *RouterHandle) T,
	opts ...RouterOption,
) {
	b.Helper()
	conform.Bench(b, RouterSuite(), callable, routerSubjects(newSubject), opts...)
}

// routerSubjects returns the factory of the subjects of Router that calls
// newSubject with the handle of each check.
func routerSubjects[T Router](newSubject func(h *RouterHandle) T) func(h *RouterHandle) Router {
	return func(h *RouterHandle) Router { return newSubject(h) }
}

// routerRouteCalls generates the calls of Router.Route with the inputs
// that each case draws.
var routerRouteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Router) error {
	f := c.Draw(prop.Of[Feed](), "f")
	return func(ctx context.Context, subject Router) error {
		return subject.Route(ctx, f)
	}
})

// routerReadCalls generates the calls of Router.Read with the inputs that
// each case draws.
var routerReadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Router) error {
	inValues := c.Draw(prop.List(prop.Of[string]()), "in")
	return func(_ context.Context, subject Router) error {
		_ = subject.Read(func() <-chan string {
			ch := make(chan string, len(inValues))
			for _, v := range inValues {
				ch <- v
			}
			close(ch)
			return ch
		}())
		return nil
	}
})

// routerWriteCalls generates the calls of Router.Write with the inputs
// that each case draws.
var routerWriteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Router) error {
	out := c.Draw(prop.Of[chan<- string](), "out")
	return func(_ context.Context, subject Router) error {
		subject.Write(out)
		return nil
	}
})

// routerPassCalls generates the calls of Router.Pass with the inputs that
// each case draws.
var routerPassCalls = prop.Composite(func(c *prop.Case) func(context.Context, Router) error {
	both := c.Draw(prop.Of[chan string](), "both")
	return func(_ context.Context, subject Router) error {
		subject.Pass(both)
		return nil
	}
})

// routerWatchCalls generates the calls of Router.Watch with the inputs
// that each case draws.
var routerWatchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Router) error {
	sUpdatesResults := c.Draw(prop.List(prop.Of[<-chan string]()), "s-updates")
	sUpdatesEnd := conform.Closed[<-chan string]()
	var sUpdatesAt atomic.Int64
	sResultsResults := c.Draw(prop.List(prop.Of[chan<- string]()), "s-results")
	sResultsEnd := c.Draw(prop.Of[chan<- string](), "s-results-end")
	var sResultsAt atomic.Int64
	sPipeResults := c.Draw(prop.List(prop.Of[chan string]()), "s-pipe")
	sPipeEnd := conform.Closed[chan string]()
	var sPipeAt atomic.Int64
	s := &routerSourceInput{
		onUpdates: func() (r0 <-chan string) {
			i := int(sUpdatesAt.Add(1) - 1)
			if i >= len(sUpdatesResults) {
				r0 = sUpdatesEnd
				return r0
			}
			return sUpdatesResults[i]
		},
		onResults: func() (r0 chan<- string) {
			i := int(sResultsAt.Add(1) - 1)
			if i >= len(sResultsResults) {
				r0 = sResultsEnd
				return r0
			}
			return sResultsResults[i]
		},
		onPipe: func() (r0 chan string) {
			i := int(sPipeAt.Add(1) - 1)
			if i >= len(sPipeResults) {
				r0 = sPipeEnd
				return r0
			}
			return sPipeResults[i]
		},
	}
	return func(_ context.Context, subject Router) error {
		return subject.Watch(s)
	}
})

// routerSourceInputEmbedded is Source under the name of the field that
// routerSourceInput embeds.
type routerSourceInputEmbedded = Source

// routerSourceInput implements Source. Its methods return the results that
// the generators of the suite of Router draw for them, in turn.
type routerSourceInput struct {
	routerSourceInputEmbedded
	onUpdates func() <-chan string
	onResults func() chan<- string
	onPipe    func() chan string
}

// Updates calls the function of Updates.
func (in routerSourceInput) Updates() <-chan string {
	return in.onUpdates()
}

// Results calls the function of Results.
func (in routerSourceInput) Results() chan<- string {
	return in.onResults()
}

// Pipe calls the function of Pipe.
func (in routerSourceInput) Pipe() chan string {
	return in.onPipe()
}

// routerOverEmbedded is Router under the name of the field that routerOver
// embeds.
type routerOverEmbedded = Router

// routerOver is a subject of Router whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type routerOver struct {
	routerOverEmbedded
	spend func()
}

// Route calls spend and returns its zero results.
func (o routerOver) Route(context.Context, Feed) (r0 error) {
	o.spend()
	return r0
}

// Read calls spend and returns its zero results.
func (o routerOver) Read(<-chan string) (r0 int) {
	o.spend()
	return r0
}

// Write calls spend and returns its zero results.
func (o routerOver) Write(chan<- string) {
	o.spend()
}

// Pass calls spend and returns its zero results.
func (o routerOver) Pass(chan string) {
	o.spend()
}

// Watch calls spend and returns its zero results.
func (o routerOver) Watch(Source) (r0 error) {
	o.spend()
	return r0
}

// newRouterOver returns the planted defect of a ceiling of Router, whose
// callables call spend.
func newRouterOver(_ *RouterHandle, spend func()) Router {
	return routerOver{spend: spend}
}

// The refused claims of Drawer:
//
// Drawer.Draw/deterministic: eidos derives no sample of the parameter s of
// Draw. The remedy is to author a sample of the parameter s with the
// directive sample.
//
// Drawer.Count/deterministic: eidos derives no sample of the parameter
// shapes of Count. The remedy is to author a sample of the parameter
// shapes with the directive sample.
//
// Drawer.Name/deterministic: eidos derives no sample of the parameter
// names of Name. The remedy is to author a sample of the parameter names
// with the directive sample.
//
// Drawer.Pick/deterministic: eidos derives no sample of the parameter p of
// Pick. The remedy is to author a sample of the parameter p with the
// directive sample.

// The unclaimed classifications of Drawer:
//
// The shape:mutator of Drawer.Seal has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Drawer.Test has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Drawer writes.
type (
	DrawerHandle = conform.Handle[Drawer]
	DrawerCheck  = conform.Check[Drawer]
	DrawerOption = conform.Option[Drawer]
)

// DrawerChecks contains the name of each callable of Drawer, with the
// identities of its checks.
var DrawerChecks = struct {
	Draw  drawerDraw
	Count drawerCount
	Name  drawerName
	Pick  drawerPick
	Seal  drawerSeal
	Test  drawerTest
}{
	Draw:  drawerDraw{"Drawer.Draw"},
	Count: drawerCount{"Drawer.Count"},
	Name:  drawerName{"Drawer.Name"},
	Pick:  drawerPick{"Drawer.Pick"},
	Seal:  drawerSeal{"Drawer.Seal"},
	Test:  drawerTest{"Drawer.Test"},
}

// drawerDraw is the name of Drawer.Draw, with the identities of its
// checks.
type drawerDraw struct {
	conform.Callable[Drawer]
}

// drawerCount is the name of Drawer.Count, with the identities of its
// checks.
type drawerCount struct {
	conform.Callable[Drawer]
}

// drawerName is the name of Drawer.Name, with the identities of its
// checks.
type drawerName struct {
	conform.Callable[Drawer]
}

// drawerPick is the name of Drawer.Pick, with the identities of its
// checks.
type drawerPick struct {
	conform.Callable[Drawer]
}

// drawerSeal is the name of Drawer.Seal, with the identities of its
// checks.
type drawerSeal struct {
	conform.Callable[Drawer]
}

// drawerTest is the name of Drawer.Test, with the identities of its
// checks.
type drawerTest struct {
	conform.Callable[Drawer]
}

// DrawerSuite returns the checks of Drawer, and the calls and the budgets
// of its callables.
func DrawerSuite() conform.Suite[Drawer] {
	return conform.Suite[Drawer]{
		Declaration: "Drawer",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Deterministic,
				Callable: "Drawer.Draw",
				Why:      "eidos derives no sample of the parameter s of Draw",
				Remedy:   "author a sample of the parameter s with the directive sample",
			},
			{
				Claim:    spec.Deterministic,
				Callable: "Drawer.Count",
				Why:      "eidos derives no sample of the parameter shapes of Count",
				Remedy:   "author a sample of the parameter shapes with the directive sample",
			},
			{
				Claim:    spec.Deterministic,
				Callable: "Drawer.Name",
				Why:      "eidos derives no sample of the parameter names of Name",
				Remedy:   "author a sample of the parameter names with the directive sample",
			},
			{
				Claim:    spec.Deterministic,
				Callable: "Drawer.Pick",
				Why:      "eidos derives no sample of the parameter p of Pick",
				Remedy:   "author a sample of the parameter p with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:mutator",
				Callable:       "Drawer.Seal",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Drawer.Test",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Drawer]{
			{
				Calls:    drawerDrawCalls,
				Over:     newDrawerOver,
				Callable: DrawerChecks.Draw.Callable,
			},
			{
				Calls:    drawerCountCalls,
				Over:     newDrawerOver,
				Callable: DrawerChecks.Count.Callable,
			},
			{
				Calls:    drawerNameCalls,
				Over:     newDrawerOver,
				Callable: DrawerChecks.Name.Callable,
			},
			{
				Calls:    drawerPickCalls,
				Over:     newDrawerOver,
				Callable: DrawerChecks.Pick.Callable,
			},
			{
				Calls:    drawerSealCalls,
				Over:     newDrawerOver,
				Callable: DrawerChecks.Seal.Callable,
			},
			{
				Calls:    drawerTestCalls,
				Over:     newDrawerOver,
				Callable: DrawerChecks.Test.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(Left, Right)),
			prop.Using(prop.Composite(func(c *prop.Case) fmt.Stringer {
				stringerStringResults := c.Draw(prop.List(prop.Of[string]()), "stringer-string")
				var stringerStringAt atomic.Int64
				return &drawerStringerInput{
					onString: func() (r0 string) {
						i := int(stringerStringAt.Add(1) - 1)
						if i >= len(stringerStringResults) {
							return r0
						}
						return stringerStringResults[i]
					},
				}
			})),
			prop.Using(prop.OneOf(
				prop.Composite(func(c *prop.Case) Shape {
					return c.Draw(prop.Of[Circle](), "circle")
				}),
				prop.Composite(func(c *prop.Case) Shape {
					return c.Draw(prop.Of[Corner](), "corner")
				}),
				prop.Composite(func(c *prop.Case) Shape {
					v := c.Draw(prop.Of[Square](), "square")
					return &v
				}),
			)),
		},
	}
}

// RunDrawerSuite applies the checks of Drawer, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunDrawerSuite[T Drawer](t *testing.T, newSubject func() T, opts ...DrawerOption) {
	t.Helper()
	RunDrawerSuiteFromHandle(t, func(*DrawerHandle) T { return newSubject() }, opts...)
}

// RunDrawerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDrawerSuiteFromHandle[T Drawer](
	t *testing.T,
	newSubject func(h *DrawerHandle) T,
	opts ...DrawerOption,
) {
	t.Helper()
	conform.Run(t, DrawerSuite(), drawerSubjects(newSubject), opts...)
}

// FuzzDrawerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzDrawerSuite[T Drawer](
	f *testing.F,
	callable conform.Callable[Drawer],
	newSubject func(h *DrawerHandle) T,
	opts ...DrawerOption,
) {
	f.Helper()
	conform.Fuzz(f, DrawerSuite(), callable, drawerSubjects(newSubject), opts...)
}

// RunDrawerBudgets applies the allocation ceilings of Drawer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunDrawerBudgets[T Drawer](
	t *testing.T,
	newSubject func(h *DrawerHandle) T,
	opts ...DrawerOption,
) {
	t.Helper()
	conform.RunBudgets(t, DrawerSuite(), drawerSubjects(newSubject), opts...)
}

// BenchDrawerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchDrawerSuite[T Drawer](
	b conform.B,
	callable conform.Callable[Drawer],
	newSubject func(h *DrawerHandle) T,
	opts ...DrawerOption,
) {
	b.Helper()
	conform.Bench(b, DrawerSuite(), callable, drawerSubjects(newSubject), opts...)
}

// drawerSubjects returns the factory of the subjects of Drawer that calls
// newSubject with the handle of each check.
func drawerSubjects[T Drawer](newSubject func(h *DrawerHandle) T) func(h *DrawerHandle) Drawer {
	return func(h *DrawerHandle) Drawer { return newSubject(h) }
}

// drawerDrawCalls generates the calls of Drawer.Draw with the inputs that
// each case draws.
var drawerDrawCalls = prop.Composite(func(c *prop.Case) func(context.Context, Drawer) error {
	s := c.Draw(prop.Of[Shape](), "s")
	return func(_ context.Context, subject Drawer) error {
		_ = subject.Draw(s)
		return nil
	}
})

// drawerCountCalls generates the calls of Drawer.Count with the inputs
// that each case draws.
var drawerCountCalls = prop.Composite(func(c *prop.Case) func(context.Context, Drawer) error {
	shapes := c.Draw(prop.Of[map[Shape]int](), "shapes")
	return func(_ context.Context, subject Drawer) error {
		_ = subject.Count(shapes)
		return nil
	}
})

// drawerNameCalls generates the calls of Drawer.Name with the inputs that
// each case draws.
var drawerNameCalls = prop.Composite(func(c *prop.Case) func(context.Context, Drawer) error {
	names := c.Draw(prop.Of[map[fmt.Stringer]string](), "names")
	return func(_ context.Context, subject Drawer) error {
		_ = subject.Name(names)
		return nil
	}
})

// drawerPickCalls generates the calls of Drawer.Pick with the inputs that
// each case draws.
var drawerPickCalls = prop.Composite(func(c *prop.Case) func(context.Context, Drawer) error {
	pPickResults := c.Draw(prop.List(prop.Of[Sealed]()), "p-pick")
	var pPickAt atomic.Int64
	p := &drawerPickerInput{
		onPick: func() (r0 Sealed) {
			i := int(pPickAt.Add(1) - 1)
			if i >= len(pPickResults) {
				return r0
			}
			return pPickResults[i]
		},
	}
	return func(_ context.Context, subject Drawer) error {
		_ = subject.Pick(p)
		return nil
	}
})

// drawerSealCalls generates the calls of Drawer.Seal with the inputs that
// each case draws.
var drawerSealCalls = prop.Composite(func(c *prop.Case) func(context.Context, Drawer) error {
	s := c.Draw(prop.Of[Sealed](), "s")
	return func(_ context.Context, subject Drawer) error {
		subject.Seal(s)
		return nil
	}
})

// drawerTestCalls generates the calls of Drawer.Test with the inputs that
// each case draws.
var drawerTestCalls = prop.Composite(func(c *prop.Case) func(context.Context, Drawer) error {
	tb := c.Draw(prop.Of[testing.TB](), "tb")
	return func(_ context.Context, subject Drawer) error {
		subject.Test(tb)
		return nil
	}
})

// drawerPickerInputEmbedded is Picker under the name of the field that
// drawerPickerInput embeds.
type drawerPickerInputEmbedded = Picker

// drawerPickerInput implements Picker. Its methods return the results that
// the generators of the suite of Drawer draw for them, in turn.
type drawerPickerInput struct {
	drawerPickerInputEmbedded
	onPick func() Sealed
}

// Pick calls the function of Pick.
func (in drawerPickerInput) Pick() Sealed {
	return in.onPick()
}

// drawerStringerInputEmbedded is fmt.Stringer under the name of the field
// that drawerStringerInput embeds.
type drawerStringerInputEmbedded = fmt.Stringer

// drawerStringerInput implements fmt.Stringer. Its methods return the
// results that the generators of the suite of Drawer draw for them, in
// turn.
type drawerStringerInput struct {
	drawerStringerInputEmbedded
	onString func() string
}

// String calls the function of String.
func (in drawerStringerInput) String() string {
	return in.onString()
}

// drawerOverEmbedded is Drawer under the name of the field that drawerOver
// embeds.
type drawerOverEmbedded = Drawer

// drawerOver is a subject of Drawer whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type drawerOver struct {
	drawerOverEmbedded
	spend func()
}

// Draw calls spend and returns its zero results.
func (o drawerOver) Draw(Shape) (r0 float64) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o drawerOver) Count(map[Shape]int) (r0 int) {
	o.spend()
	return r0
}

// Name calls spend and returns its zero results.
func (o drawerOver) Name(map[fmt.Stringer]string) (r0 int) {
	o.spend()
	return r0
}

// Pick calls spend and returns its zero results.
func (o drawerOver) Pick(Picker) (r0 int) {
	o.spend()
	return r0
}

// Seal calls spend and returns its zero results.
func (o drawerOver) Seal(Sealed) {
	o.spend()
}

// Test calls spend and returns its zero results.
func (o drawerOver) Test(testing.TB) {
	o.spend()
}

// newDrawerOver returns the planted defect of a ceiling of Drawer, whose
// callables call spend.
func newDrawerOver(_ *DrawerHandle, spend func()) Drawer {
	return drawerOver{spend: spend}
}

// The derived checks of Odd:
//
// Odd.Keep/smoke: Keep returns without a panic for derived inputs.

// The derived checks of Odd that witness their laws at the shallow end:
//
// Odd.Keep/smoke

// The refused claims of Odd:
//
// Odd.Embed/deterministic: eidos derives no sample of the parameter v of
// Embed. The remedy is to author a sample of the parameter v with the
// directive sample.

// The unclaimed classifications of Odd:
//
// The shape:mutator of Odd.Tag has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Odd.Keep has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Odd.Loop has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Odd writes.
type (
	OddHandle = conform.Handle[Odd]
	OddCheck  = conform.Check[Odd]
	OddOption = conform.Option[Odd]
)

// OddChecks contains the name of each callable of Odd, with the identities
// of its checks.
var OddChecks = struct {
	Tag   oddTag
	Keep  oddKeep
	Embed oddEmbed
	Loop  oddLoop
}{
	Tag:   oddTag{"Odd.Tag"},
	Keep:  oddKeep{"Odd.Keep"},
	Embed: oddEmbed{"Odd.Embed"},
	Loop:  oddLoop{"Odd.Loop"},
}

// oddTag is the name of Odd.Tag, with the identities of its checks.
type oddTag struct {
	conform.Callable[Odd]
}

// oddKeep is the name of Odd.Keep, with the identities of its checks.
type oddKeep struct {
	conform.Callable[Odd]
}

// Smoke returns the identity of the check that Keep returns without a
// panic for derived inputs.
func (c oddKeep) Smoke() conform.ID[Odd] {
	return c.Check("smoke")
}

// oddEmbed is the name of Odd.Embed, with the identities of its checks.
type oddEmbed struct {
	conform.Callable[Odd]
}

// oddLoop is the name of Odd.Loop, with the identities of its checks.
type oddLoop struct {
	conform.Callable[Odd]
}

// OddSuite returns the checks of Odd, and the calls and the budgets of its
// callables.
func OddSuite() conform.Suite[Odd] {
	return conform.Suite[Odd]{
		Declaration: "Odd",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Odd]{
			{
				ID:          OddChecks.Keep.Smoke(),
				Statement:   "Keep returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OddHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Keep(Anything("a-a"))
						return nil
					})
				},
				Defect: oddKeepPanics,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Deterministic,
				Callable: "Odd.Embed",
				Why:      "eidos derives no sample of the parameter v of Embed",
				Remedy:   "author a sample of the parameter v with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:mutator",
				Callable:       "Odd.Tag",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Odd.Keep",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Odd.Loop",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Odd]{
			{
				Calls:    oddTagCalls,
				Over:     newOddOver,
				Callable: OddChecks.Tag.Callable,
			},
			{
				Calls:    oddKeepCalls,
				Over:     newOddOver,
				Callable: OddChecks.Keep.Callable,
			},
			{
				Calls:    oddEmbedCalls,
				Over:     newOddOver,
				Callable: OddChecks.Embed.Callable,
			},
			{
				Calls:    oddLoopCalls,
				Over:     newOddOver,
				Callable: OddChecks.Loop.Callable,
			},
		},
	}
}

// RunOddSuite applies the checks of Odd, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunOddSuite[T Odd](t *testing.T, newSubject func() T, opts ...OddOption) {
	t.Helper()
	RunOddSuiteFromHandle(t, func(*OddHandle) T { return newSubject() }, opts...)
}

// RunOddSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunOddSuiteFromHandle[T Odd](
	t *testing.T,
	newSubject func(h *OddHandle) T,
	opts ...OddOption,
) {
	t.Helper()
	conform.Run(t, OddSuite(), oddSubjects(newSubject), opts...)
}

// FuzzOddSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzOddSuite[T Odd](
	f *testing.F,
	callable conform.Callable[Odd],
	newSubject func(h *OddHandle) T,
	opts ...OddOption,
) {
	f.Helper()
	conform.Fuzz(f, OddSuite(), callable, oddSubjects(newSubject), opts...)
}

// RunOddBudgets applies the allocation ceilings of Odd and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunOddBudgets[T Odd](t *testing.T, newSubject func(h *OddHandle) T, opts ...OddOption) {
	t.Helper()
	conform.RunBudgets(t, OddSuite(), oddSubjects(newSubject), opts...)
}

// BenchOddSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchOddSuite[T Odd](
	b conform.B,
	callable conform.Callable[Odd],
	newSubject func(h *OddHandle) T,
	opts ...OddOption,
) {
	b.Helper()
	conform.Bench(b, OddSuite(), callable, oddSubjects(newSubject), opts...)
}

// oddSubjects returns the factory of the subjects of Odd that calls
// newSubject with the handle of each check.
func oddSubjects[T Odd](newSubject func(h *OddHandle) T) func(h *OddHandle) Odd {
	return func(h *OddHandle) Odd { return newSubject(h) }
}

// oddTagCalls generates the calls of Odd.Tag with the inputs that each
// case draws.
var oddTagCalls = prop.Composite(func(c *prop.Case) func(context.Context, Odd) error {
	lStringResults := c.Draw(prop.List(prop.Of[string]()), "l-string")
	var lStringAt atomic.Int64
	l := &oddStringerInput{
		onString: func() (r0 string) {
			i := int(lStringAt.Add(1) - 1)
			if i >= len(lStringResults) {
				return r0
			}
			return lStringResults[i]
		},
	}
	return func(_ context.Context, subject Odd) error {
		subject.Tag(l)
		return nil
	}
})

// oddKeepCalls generates the calls of Odd.Keep with the inputs that each
// case draws.
var oddKeepCalls = prop.Composite(func(*prop.Case) func(context.Context, Odd) error {
	a := &oddAnythingInput{}
	return func(_ context.Context, subject Odd) error {
		subject.Keep(a)
		return nil
	}
})

// oddEmbedCalls generates the calls of Odd.Embed with the inputs that each
// case draws.
var oddEmbedCalls = prop.Composite(func(c *prop.Case) func(context.Context, Odd) error {
	vSizeResults := c.Draw(prop.List(prop.Of[int]()), "v-size")
	var vSizeAt atomic.Int64
	vNameResults := c.Draw(prop.List(prop.Of[string]()), "v-name")
	var vNameAt atomic.Int64
	v := &oddInterfaceInput{
		onSize: func() (r0 int) {
			i := int(vSizeAt.Add(1) - 1)
			if i >= len(vSizeResults) {
				return r0
			}
			return vSizeResults[i]
		},
		onReset: func() {},
		onName: func() (r0 string) {
			i := int(vNameAt.Add(1) - 1)
			if i >= len(vNameResults) {
				return r0
			}
			return vNameResults[i]
		},
	}
	return func(_ context.Context, subject Odd) error {
		_ = subject.Embed(v)
		return nil
	}
})

// oddLoopCalls generates the calls of Odd.Loop with the inputs that each
// case draws.
var oddLoopCalls = prop.Composite(func(c *prop.Case) func(context.Context, Odd) error {
	r := c.Draw(prop.Of[Ring1](), "r")
	return func(_ context.Context, subject Odd) error {
		subject.Loop(r)
		return nil
	}
})

// oddStringerInputEmbedded is fmt.Stringer under the name of the field
// that oddStringerInput embeds.
type oddStringerInputEmbedded = fmt.Stringer

// oddStringerInput implements fmt.Stringer. Its methods return the results
// that the generators of the suite of Odd draw for them, in turn.
type oddStringerInput struct {
	oddStringerInputEmbedded
	onString func() string
}

// String calls the function of String.
func (in oddStringerInput) String() string {
	return in.onString()
}

// oddAnythingInputEmbedded is any under the name of the field that
// oddAnythingInput embeds.
type oddAnythingInputEmbedded = any

// oddAnythingInput implements any. Its methods return the results that the
// generators of the suite of Odd draw for them, in turn.
type oddAnythingInput struct {
	oddAnythingInputEmbedded
}

// oddInterfaceInputEmbedded is interface { Namer Size() int Reset() }
// under the name of the field that oddInterfaceInput embeds.
type oddInterfaceInputEmbedded = interface {
	Namer
	Size() int
	Reset()
}

// oddInterfaceInput implements interface { Namer Size() int Reset() }. Its
// methods return the results that the generators of the suite of Odd draw
// for them, in turn.
type oddInterfaceInput struct {
	oddInterfaceInputEmbedded
	onSize  func() int
	onReset func()
	onName  func() string
}

// Size calls the function of Size.
func (in oddInterfaceInput) Size() int {
	return in.onSize()
}

// Reset calls the function of Reset.
func (in oddInterfaceInput) Reset() {
	in.onReset()
}

// Name calls the function of Name.
func (in oddInterfaceInput) Name() string {
	return in.onName()
}

// oddOverEmbedded is Odd under the name of the field that oddOver embeds.
type oddOverEmbedded = Odd

// oddOver is a subject of Odd whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type oddOver struct {
	oddOverEmbedded
	spend func()
}

// Tag calls spend and returns its zero results.
func (o oddOver) Tag(Label) {
	o.spend()
}

// Keep calls spend and returns its zero results.
func (o oddOver) Keep(Anything) {
	o.spend()
}

// Embed calls spend and returns its zero results.
func (o oddOver) Embed(interface {
	Namer
	Size() int
	Reset()
},
) (r0 int) {
	o.spend()
	return r0
}

// Loop calls spend and returns its zero results.
func (o oddOver) Loop(Ring1) {
	o.spend()
}

// newOddOver returns the planted defect of a ceiling of Odd, whose
// callables call spend.
func newOddOver(_ *OddHandle, spend func()) Odd {
	return oddOver{spend: spend}
}

// oddKeepPanicsDefect is the planted defect of a Odd whose Keep panics.
type oddKeepPanicsDefect struct {
	oddOverEmbedded
}

// Keep is the override of the planted defect.
func (*oddKeepPanicsDefect) Keep(p0 Anything) {
	panic("the planted defect of a Odd whose Keep panics")
}

// oddKeepPanics returns the planted defect of a Odd whose Keep panics over
// a fresh subject.
func oddKeepPanics(h *OddHandle) Odd {
	return &oddKeepPanicsDefect{oddOverEmbedded: h.Fresh()}
}

// The derived checks of Std:
//
// Std.Within/smoke: Within returns without a panic for derived inputs.
//
// Std.Within/cancel: Within returns an error for a cancelled context.
//
// Std.Pair/smoke: Pair returns without a panic for derived inputs.
//
// Std.Meet/smoke: Meet returns without a panic for derived inputs.

// The derived checks of Std that witness their laws at the shallow end:
//
// Std.Within/smoke
//
// Std.Pair/smoke
//
// Std.Meet/smoke

// The unclaimed classifications of Std:
//
// The shape:composite-writer of Std.Copy has the stance silent. A write is
// an input of the other checks, and its result is an error alone.
//
// The shape:writer of Std.Close has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Std.From has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Std.At has the stance silent. A write is an input of
// the other checks, and its result is an error alone.
//
// The shape:composite-writer of Std.Buffer has the stance silent. A write
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Std.Lock has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Std.Fail has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Std.Within has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:mutator of Std.Hold has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Std.Pair has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Std.Meet has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Std writes.
type (
	StdHandle = conform.Handle[Std]
	StdCheck  = conform.Check[Std]
	StdOption = conform.Option[Std]
)

// StdChecks contains the name of each callable of Std, with the identities
// of its checks.
var StdChecks = struct {
	Copy   stdCopy
	Close  stdClose
	From   stdFrom
	At     stdAt
	Log    stdLog
	Buffer stdBuffer
	Lock   stdLock
	Fail   stdFail
	Within stdWithin
	Hold   stdHold
	Pair   stdPair
	Meet   stdMeet
}{
	Copy:   stdCopy{"Std.Copy"},
	Close:  stdClose{"Std.Close"},
	From:   stdFrom{"Std.From"},
	At:     stdAt{"Std.At"},
	Log:    stdLog{"Std.Log"},
	Buffer: stdBuffer{"Std.Buffer"},
	Lock:   stdLock{"Std.Lock"},
	Fail:   stdFail{"Std.Fail"},
	Within: stdWithin{"Std.Within"},
	Hold:   stdHold{"Std.Hold"},
	Pair:   stdPair{"Std.Pair"},
	Meet:   stdMeet{"Std.Meet"},
}

// stdCopy is the name of Std.Copy, with the identities of its checks.
type stdCopy struct {
	conform.Callable[Std]
}

// stdClose is the name of Std.Close, with the identities of its checks.
type stdClose struct {
	conform.Callable[Std]
}

// stdFrom is the name of Std.From, with the identities of its checks.
type stdFrom struct {
	conform.Callable[Std]
}

// stdAt is the name of Std.At, with the identities of its checks.
type stdAt struct {
	conform.Callable[Std]
}

// stdLog is the name of Std.Log, with the identities of its checks.
type stdLog struct {
	conform.Callable[Std]
}

// stdBuffer is the name of Std.Buffer, with the identities of its checks.
type stdBuffer struct {
	conform.Callable[Std]
}

// stdLock is the name of Std.Lock, with the identities of its checks.
type stdLock struct {
	conform.Callable[Std]
}

// stdFail is the name of Std.Fail, with the identities of its checks.
type stdFail struct {
	conform.Callable[Std]
}

// stdWithin is the name of Std.Within, with the identities of its checks.
type stdWithin struct {
	conform.Callable[Std]
}

// Smoke returns the identity of the check that Within returns without a
// panic for derived inputs.
func (c stdWithin) Smoke() conform.ID[Std] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Within returns an error
// for a cancelled context.
func (c stdWithin) Cancel() conform.ID[Std] {
	return c.Check("cancel")
}

// stdHold is the name of Std.Hold, with the identities of its checks.
type stdHold struct {
	conform.Callable[Std]
}

// stdPair is the name of Std.Pair, with the identities of its checks.
type stdPair struct {
	conform.Callable[Std]
}

// Smoke returns the identity of the check that Pair returns without a
// panic for derived inputs.
func (c stdPair) Smoke() conform.ID[Std] {
	return c.Check("smoke")
}

// stdMeet is the name of Std.Meet, with the identities of its checks.
type stdMeet struct {
	conform.Callable[Std]
}

// Smoke returns the identity of the check that Meet returns without a
// panic for derived inputs.
func (c stdMeet) Smoke() conform.ID[Std] {
	return c.Check("smoke")
}

// StdSuite returns the checks of Std, and the calls and the budgets of its
// callables.
func StdSuite() conform.Suite[Std] {
	return conform.Suite[Std]{
		Declaration: "Std",
		Package:     "example.com/shop/kinds",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Std]{
			{
				ID:          StdChecks.Within.Smoke(),
				Statement:   "Within returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StdHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Within(ctx, ctx)
					})
				},
				Defect: stdWithinPanics,
			},
			{
				ID:          StdChecks.Within.Cancel(),
				Statement:   "Within returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StdHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Within(ctx, ctx)
					})
				},
				Defect: stdWithinIgnoresCancel,
			},
			{
				ID:          StdChecks.Pair.Smoke(),
				Statement:   "Pair returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StdHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Pair([]Twins{{A: 42}})
						return nil
					})
				},
				Defect: stdPairPanics,
			},
			{
				ID:          StdChecks.Meet.Smoke(),
				Statement:   "Meet returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StdHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Meet(Clash{Name: "Name-a"})
						return nil
					})
				},
				Defect: stdMeetPanics,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Std.Copy",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Std.Close",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Std.From",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Std.At",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Std.Buffer",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Std.Lock",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Std.Fail",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Std.Within",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Std.Hold",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Std.Pair",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Std.Meet",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Std]{
			{
				Calls:    stdCopyCalls,
				Over:     newStdOver,
				Callable: StdChecks.Copy.Callable,
			},
			{
				Calls:    stdCloseCalls,
				Over:     newStdOver,
				Callable: StdChecks.Close.Callable,
			},
			{
				Calls:    stdFromCalls,
				Over:     newStdOver,
				Callable: StdChecks.From.Callable,
			},
			{
				Calls:    stdAtCalls,
				Over:     newStdOver,
				Callable: StdChecks.At.Callable,
			},
			{
				Calls:    stdLogCalls,
				Over:     newStdOver,
				Callable: StdChecks.Log.Callable,
			},
			{
				Calls:    stdBufferCalls,
				Over:     newStdOver,
				Callable: StdChecks.Buffer.Callable,
			},
			{
				Calls:    stdLockCalls,
				Over:     newStdOver,
				Callable: StdChecks.Lock.Callable,
			},
			{
				Calls:    stdFailCalls,
				Over:     newStdOver,
				Callable: StdChecks.Fail.Callable,
			},
			{
				Calls:    stdWithinCalls,
				Over:     newStdOver,
				Callable: StdChecks.Within.Callable,
			},
			{
				Calls:    stdHoldCalls,
				Over:     newStdOver,
				Callable: StdChecks.Hold.Callable,
			},
			{
				Calls:    stdPairCalls,
				Over:     newStdOver,
				Callable: StdChecks.Pair.Callable,
			},
			{
				Calls:    stdMeetCalls,
				Over:     newStdOver,
				Callable: StdChecks.Meet.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(stdMutexValues),
			prop.Using(stdBoxValues),
			prop.Using(stdBuilderValues),
			prop.Using(stdStructValues),
			prop.Using(prop.Composite(func(c *prop.Case) interface{ Name() string } {
				interfaceNameResults := c.Draw(prop.List(prop.Of[string]()), "interface-name")
				var interfaceNameAt atomic.Int64
				return &stdInterfaceInput{
					onName: func() (r0 string) {
						i := int(interfaceNameAt.Add(1) - 1)
						if i >= len(interfaceNameResults) {
							return r0
						}
						return interfaceNameResults[i]
					},
				}
			})),
			prop.Using(conform.Stream[<-chan time.Time](prop.Of[time.Time]())),
			prop.Using(prop.Composite(func(c *prop.Case) io.Writer {
				return conform.NewBytes(c.Draw(prop.Of[[]byte](), "writer"))
			})),
			prop.Using(prop.Composite(func(c *prop.Case) io.Reader {
				return conform.NewBytes(c.Draw(prop.Of[[]byte](), "reader"))
			})),
			prop.Using(prop.Composite(func(c *prop.Case) io.ReadCloser {
				return conform.NewBytes(c.Draw(prop.Of[[]byte](), "read-closer"))
			})),
			prop.Using(prop.Composite(func(c *prop.Case) *bufio.Reader {
				return bufio.NewReader(conform.NewBytes(c.Draw(prop.Of[[]byte](), "reader")))
			})),
			prop.Using(
				prop.Composite(
					func(*prop.Case) *slog.Logger { return slog.New(slog.DiscardHandler) },
				),
			),
			prop.Using(
				prop.Composite(func(*prop.Case) *log.Logger { return log.New(io.Discard, "", 0) }),
			),
			prop.Using(
				prop.Composite(
					func(*prop.Case) *bufio.Writer { return bufio.NewWriter(io.Discard) },
				),
			),
			prop.Using(
				prop.Composite(func(*prop.Case) context.Context { return context.Background() }),
			),
			prop.Using(prop.Composite(func(*prop.Case) unsafe.Pointer { return nil })),
			prop.Using(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		},
	}
}

// RunStdSuite applies the checks of Std, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunStdSuite[T Std](t *testing.T, newSubject func() T, opts ...StdOption) {
	t.Helper()
	RunStdSuiteFromHandle(t, func(*StdHandle) T { return newSubject() }, opts...)
}

// RunStdSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStdSuiteFromHandle[T Std](
	t *testing.T,
	newSubject func(h *StdHandle) T,
	opts ...StdOption,
) {
	t.Helper()
	conform.Run(t, StdSuite(), stdSubjects(newSubject), opts...)
}

// FuzzStdSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStdSuite[T Std](
	f *testing.F,
	callable conform.Callable[Std],
	newSubject func(h *StdHandle) T,
	opts ...StdOption,
) {
	f.Helper()
	conform.Fuzz(f, StdSuite(), callable, stdSubjects(newSubject), opts...)
}

// RunStdBudgets applies the allocation ceilings of Std and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunStdBudgets[T Std](t *testing.T, newSubject func(h *StdHandle) T, opts ...StdOption) {
	t.Helper()
	conform.RunBudgets(t, StdSuite(), stdSubjects(newSubject), opts...)
}

// BenchStdSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchStdSuite[T Std](
	b conform.B,
	callable conform.Callable[Std],
	newSubject func(h *StdHandle) T,
	opts ...StdOption,
) {
	b.Helper()
	conform.Bench(b, StdSuite(), callable, stdSubjects(newSubject), opts...)
}

// stdSubjects returns the factory of the subjects of Std that calls
// newSubject with the handle of each check.
func stdSubjects[T Std](newSubject func(h *StdHandle) T) func(h *StdHandle) Std {
	return func(h *StdHandle) Std { return newSubject(h) }
}

// stdCopyCalls generates the calls of Std.Copy with the inputs that each
// case draws.
var stdCopyCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	dst := c.Draw(prop.Of[io.Writer](), "dst")
	src := c.Draw(prop.Of[io.Reader](), "src")
	return func(_ context.Context, subject Std) error {
		return subject.Copy(dst, src)
	}
})

// stdCloseCalls generates the calls of Std.Close with the inputs that each
// case draws.
var stdCloseCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	c0 := c.Draw(prop.Of[io.ReadCloser](), "c")
	return func(_ context.Context, subject Std) error {
		return subject.Close(c0)
	}
})

// stdFromCalls generates the calls of Std.From with the inputs that each
// case draws.
var stdFromCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	rReadFromResults0 := c.Draw(prop.List(prop.Of[int64]()), "r-read-from-0")
	rReadFromResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(rReadFromResults0)),
			prop.MaxSize(len(rReadFromResults0)),
		),
		"r-read-from-1",
	)
	var rReadFromAt atomic.Int64
	r := &stdReaderFromInput{
		onReadFrom: func(io.Reader) (r0 int64, r1 error) {
			i := int(rReadFromAt.Add(1) - 1)
			if i >= len(rReadFromResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return rReadFromResults0[i],
				rReadFromResults1[i]
		},
	}
	return func(_ context.Context, subject Std) error {
		return subject.From(r)
	}
})

// stdAtCalls generates the calls of Std.At with the inputs that each case
// draws.
var stdAtCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	wWriteAtResults0 := c.Draw(prop.List(prop.Of[int]()), "w-write-at-0")
	wWriteAtResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(wWriteAtResults0)),
			prop.MaxSize(len(wWriteAtResults0)),
		),
		"w-write-at-1",
	)
	var wWriteAtAt atomic.Int64
	w := &stdWriterAtInput{
		onWriteAt: func([]byte, int64) (r0 int, r1 error) {
			i := int(wWriteAtAt.Add(1) - 1)
			if i >= len(wWriteAtResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return wWriteAtResults0[i],
				wWriteAtResults1[i]
		},
	}
	return func(_ context.Context, subject Std) error {
		return subject.At(w)
	}
})

// stdLogCalls generates the calls of Std.Log with the inputs that each
// case draws.
var stdLogCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	l := c.Draw(prop.Of[*slog.Logger](), "l")
	o := c.Draw(prop.Of[*log.Logger](), "o")
	return func(_ context.Context, subject Std) error {
		subject.Log(l, o)
		return nil
	}
})

// stdBufferCalls generates the calls of Std.Buffer with the inputs that
// each case draws.
var stdBufferCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	r := c.Draw(prop.Of[*bufio.Reader](), "r")
	w := c.Draw(prop.Of[*bufio.Writer](), "w")
	return func(_ context.Context, subject Std) error {
		return subject.Buffer(r, w)
	}
})

// stdLockCalls generates the calls of Std.Lock with the inputs that each
// case draws.
var stdLockCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	mu := c.Draw(prop.Of[sync.Mutex](), "mu")
	return func(_ context.Context, subject Std) error {
		subject.Lock(&mu)
		return nil
	}
})

// stdFailCalls generates the calls of Std.Fail with the inputs that each
// case draws.
var stdFailCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	err0 := c.Draw(prop.Of[error](), "err")
	return func(_ context.Context, subject Std) error {
		return subject.Fail(err0)
	}
})

// stdWithinCalls generates the calls of Std.Within with the inputs that
// each case draws.
var stdWithinCalls = prop.Composite(func(*prop.Case) func(context.Context, Std) error {
	return func(ctx context.Context, subject Std) error {
		return subject.Within(ctx, ctx)
	}
})

// stdHoldCalls generates the calls of Std.Hold with the inputs that each
// case draws.
var stdHoldCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	h := c.Draw(prop.Of[Holder](), "h")
	return func(_ context.Context, subject Std) error {
		subject.Hold(h)
		return nil
	}
})

// stdPairCalls generates the calls of Std.Pair with the inputs that each
// case draws.
var stdPairCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	t := c.Draw(prop.Of[[]Twins](), "t")
	return func(_ context.Context, subject Std) error {
		subject.Pair(t)
		return nil
	}
})

// stdMeetCalls generates the calls of Std.Meet with the inputs that each
// case draws.
var stdMeetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Std) error {
	c0 := c.Draw(prop.Of[Clash](), "c")
	return func(_ context.Context, subject Std) error {
		subject.Meet(c0)
		return nil
	}
})

// stdMutexValues generates the starting value of sync.Mutex, which has no
// field to change.
var stdMutexValues = prop.Composite(func(*prop.Case) sync.Mutex {
	return sync.Mutex{}
})

// stdBoxValues generates the starting value of Box, which has no field to
// change.
var stdBoxValues = prop.Composite(func(*prop.Case) Box[int] {
	return Box[int]{}
})

// stdBuilderValues generates the starting value of strings.Builder, which
// has no field to change.
var stdBuilderValues = prop.Composite(func(*prop.Case) strings.Builder {
	return strings.Builder{}
})

// stdStructValues generates the starting value of struct{ n int }, which
// has no field to change.
var stdStructValues = prop.Composite(func(*prop.Case) struct{ n int } {
	return struct{ n int }{}
})

// stdReaderFromInputEmbedded is io.ReaderFrom under the name of the field
// that stdReaderFromInput embeds.
type stdReaderFromInputEmbedded = io.ReaderFrom

// stdReaderFromInput implements io.ReaderFrom. Its methods return the
// results that the generators of the suite of Std draw for them, in turn.
type stdReaderFromInput struct {
	stdReaderFromInputEmbedded
	onReadFrom func(a0 io.Reader) (int64, error)
}

// ReadFrom calls the function of ReadFrom.
func (in stdReaderFromInput) ReadFrom(a0 io.Reader) (int64, error) {
	return in.onReadFrom(a0)
}

// stdWriterAtInputEmbedded is io.WriterAt under the name of the field that
// stdWriterAtInput embeds.
type stdWriterAtInputEmbedded = io.WriterAt

// stdWriterAtInput implements io.WriterAt. Its methods return the results
// that the generators of the suite of Std draw for them, in turn.
type stdWriterAtInput struct {
	stdWriterAtInputEmbedded
	onWriteAt func(a0 []byte, a1 int64) (int, error)
}

// WriteAt calls the function of WriteAt.
func (in stdWriterAtInput) WriteAt(a0 []byte, a1 int64) (int, error) {
	return in.onWriteAt(a0, a1)
}

// stdInterfaceInputEmbedded is interface{ Name() string } under the name
// of the field that stdInterfaceInput embeds.
type stdInterfaceInputEmbedded = interface{ Name() string }

// stdInterfaceInput implements interface{ Name() string }. Its methods
// return the results that the generators of the suite of Std draw for
// them, in turn.
type stdInterfaceInput struct {
	stdInterfaceInputEmbedded
	onName func() string
}

// Name calls the function of Name.
func (in stdInterfaceInput) Name() string {
	return in.onName()
}

// stdOverEmbedded is Std under the name of the field that stdOver embeds.
type stdOverEmbedded = Std

// stdOver is a subject of Std whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type stdOver struct {
	stdOverEmbedded
	spend func()
}

// Copy calls spend and returns its zero results.
func (o stdOver) Copy(io.Writer, io.Reader) (r0 error) {
	o.spend()
	return r0
}

// Close calls spend and returns its zero results.
func (o stdOver) Close(io.ReadCloser) (r0 error) {
	o.spend()
	return r0
}

// From calls spend and returns its zero results.
func (o stdOver) From(io.ReaderFrom) (r0 error) {
	o.spend()
	return r0
}

// At calls spend and returns its zero results.
func (o stdOver) At(io.WriterAt) (r0 error) {
	o.spend()
	return r0
}

// Log calls spend and returns its zero results.
func (o stdOver) Log(*slog.Logger, *log.Logger) {
	o.spend()
}

// Buffer calls spend and returns its zero results.
func (o stdOver) Buffer(*bufio.Reader, *bufio.Writer) (r0 error) {
	o.spend()
	return r0
}

// Lock calls spend and returns its zero results.
func (o stdOver) Lock(*sync.Mutex) {
	o.spend()
}

// Fail calls spend and returns its zero results.
func (o stdOver) Fail(error) (r0 error) {
	o.spend()
	return r0
}

// Within calls spend and returns its zero results.
func (o stdOver) Within(context.Context, context.Context) (r0 error) {
	o.spend()
	return r0
}

// Hold calls spend and returns its zero results.
func (o stdOver) Hold(Holder) {
	o.spend()
}

// Pair calls spend and returns its zero results.
func (o stdOver) Pair([]Twins) {
	o.spend()
}

// Meet calls spend and returns its zero results.
func (o stdOver) Meet(Clash) {
	o.spend()
}

// newStdOver returns the planted defect of a ceiling of Std, whose
// callables call spend.
func newStdOver(_ *StdHandle, spend func()) Std {
	return stdOver{spend: spend}
}

// stdWithinPanicsDefect is the planted defect of a Std whose Within
// panics.
type stdWithinPanicsDefect struct {
	stdOverEmbedded
}

// Within is the override of the planted defect.
func (*stdWithinPanicsDefect) Within(p0, p1 context.Context) (r0 error) {
	panic("the planted defect of a Std whose Within panics")
}

// stdWithinPanics returns the planted defect of a Std whose Within panics
// over a fresh subject.
func stdWithinPanics(h *StdHandle) Std {
	return &stdWithinPanicsDefect{stdOverEmbedded: h.Fresh()}
}

// stdWithinIgnoresCancelDefect is the planted defect of a Std whose Within
// ignores the context.
type stdWithinIgnoresCancelDefect struct {
	stdOverEmbedded
}

// Within is the override of the planted defect.
func (d *stdWithinIgnoresCancelDefect) Within(p0, p1 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.stdOverEmbedded.Within(p0, p1)
}

// stdWithinIgnoresCancel returns the planted defect of a Std whose Within
// ignores the context over a fresh subject.
func stdWithinIgnoresCancel(h *StdHandle) Std {
	return &stdWithinIgnoresCancelDefect{stdOverEmbedded: h.Fresh()}
}

// stdPairPanicsDefect is the planted defect of a Std whose Pair panics.
type stdPairPanicsDefect struct {
	stdOverEmbedded
}

// Pair is the override of the planted defect.
func (*stdPairPanicsDefect) Pair(p0 []Twins) {
	panic("the planted defect of a Std whose Pair panics")
}

// stdPairPanics returns the planted defect of a Std whose Pair panics over
// a fresh subject.
func stdPairPanics(h *StdHandle) Std {
	return &stdPairPanicsDefect{stdOverEmbedded: h.Fresh()}
}

// stdMeetPanicsDefect is the planted defect of a Std whose Meet panics.
type stdMeetPanicsDefect struct {
	stdOverEmbedded
}

// Meet is the override of the planted defect.
func (*stdMeetPanicsDefect) Meet(p0 Clash) {
	panic("the planted defect of a Std whose Meet panics")
}

// stdMeetPanics returns the planted defect of a Std whose Meet panics over
// a fresh subject.
func stdMeetPanics(h *StdHandle) Std {
	return &stdMeetPanicsDefect{stdOverEmbedded: h.Fresh()}
}
