// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	sdkrules "go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/rules"
)

// Derived is one check that a claim of the conformance spec derives for a
// callable, which the template writes as a check that calls its probe.
//
// Field and Type are the field and the type of the callable, as its
// [Callable] contains them, and Method and Segment are the method and the
// segment of its identity. Claim and Probe are the names of the constants of
// the package spec of the Go runtime, Machineries those of the machineries of
// the check, the deepest first, and Entails the identities of the shallower
// checks that the law of the check implies. The runtime takes the reason code
// and the law of the check from its claim. Args are the
// arguments of the probe after the handle. Defect is the name of the function
// that returns the planted defect, which every check with the same break of
// the same callable shares. Bubble reports a probe that runs inside a
// testing/synctest bubble: a probe of the clock, of deliveries that arrive on
// other goroutines, or of calls that wait on one computation.
//
// Generated reports a check whose inputs the property engine generates. The
// template writes it as a property check that makes the draws Draws for each
// case before it calls its probe. Base is the method of the identity of the
// sample binding of a generated binding, as ZeroOnError, and empty for any
// other check. The suite declares no method of the identity of a generated
// binding, and names it through the method Generated of the identity of its
// sample binding.
//
// Sequence reports a check that runs generated sequences. The template writes
// it as a property check that passes the case to its probe, with the types
// TypeArgs after the subject type as the type arguments of the probe. Each
// type argument is one type, or the two types of a probe.Pair. Whole is the
// identity of a check of the whole declaration, which the method of its
// identity returns, and empty for a check of a callable.
//
// Cased reports a check whose probe takes the case of a property check
// without generated sequences, such as a probe whose clients the scheduler of
// the case runs. The template writes it as a property check that passes the
// case to its probe.
//
// Door is the name of the function that opens the planted defect, and empty
// for a defect without one. A check of a world probe passes it as the door
// Host of its defect, which opens the defect over the environment of an
// incarnation. Any other check passes it as the door Open, which opens the
// defect over a medium.
//
// World reports a check of a world probe, whose one input is the world of the
// declaration. The template writes it with the environment conform.Bubble and
// with the options that probe.WorldOptions returns for its property run.
type Derived struct {
	Field       string
	Type        string
	Method      string
	Base        string
	Segment     string
	Statement   string
	Claim       string
	Probe       string
	Defect      string
	Door        string
	Whole       string
	Machineries []string
	Entails     []Identity
	Args        []Arg
	Draws       []*Draw
	TypeArgs    [][]*emit.TypeRef
	Shallow     bool
	Bubble      bool
	Generated   bool
	Sequence    bool
	Cased       bool
	World       bool
}

// Arg is one argument of a probe: a closure that calls the subject, an
// expression, a projection of a field, a duration, a method of a handle, the
// calls of a durable writer, or a part of the generated sequences, such as
// the key of a call. Ref is the name of a function of the suite, such as the
// actions of the declaration, which the template writes as a name that is not
// exported. Budget is the most steps of the check of a history, which the
// template writes as the option Budget of the package history of the assert
// library.
type Arg struct {
	Closure    *Closure
	Expr       *Expr
	Projection *Projection
	Apply      *Apply
	Written    *Written
	Countable  *Countable
	Method     *HandleMethod
	Key        *Key
	Durable    *Durable
	Ref        string
	Duration   time.Duration
	Budget     int
}

// Durable is the probe.Durable of a callable that writes and its paired read.
// For each sample, Writes and Reads are the write and the read as closures on
// the subject that they take, Wants is the value that the sample writes, and
// Keys are the arguments of the read after its context, which are the key of
// the sample. Value is the type of the value of the read.
type Durable struct {
	Value  *emit.TypeRef
	Writes [2]*Closure
	Reads  [2]*Closure
	Wants  [2]Expr
	Keys   [2][]Expr
}

// HandleMethod is the method Name of the handle of the type Type, as a method
// expression, which takes the handle as its first argument.
type HandleMethod struct {
	Type *emit.TypeRef
	Name string
}

// Key is a function that returns the key of a call: the list of the arguments
// of the call at the positions Args.
type Key struct {
	Args []int
}

// Projection is a function that returns the field Field of a value of the
// type Type, which has the type Result. A projection without a field returns
// the whole value, and its Result is Type.
type Projection struct {
	Type   *emit.TypeRef
	Result *emit.TypeRef
	Field  string
}

// Expr is one expression of a derived check: the name of a variable, a
// value, the value of a draw, the field Select of another expression, a pair
// of two expressions, an empty stream, or a word of Go such as nil or true.
// Where Hook is set, the expression is a function of the function type Hook
// that calls the function Name and ignores its arguments. Where Compute is
// set, the expression is a computation with the signature Compute. It records
// a call of the callable Name on a new dependency of the injector of the
// handle and returns its zero results. Where Convert is set, the expression is
// the conversion of the variable Name to the type Convert. Stream are the
// types of the elements of an empty stream: the value of an iter.Seq, or the
// two values of an iter.Seq2. Built is a value that the builder of a callable
// returns.
type Expr struct {
	Value   *emit.Value
	Drawn   *Drawn
	Of      *Expr
	Pair    *PairExpr
	Built   *Built
	Hook    *emit.TypeRef
	Compute *emit.Method
	Convert *emit.TypeRef
	Name    string
	Select  string
	Word    string
	Stream  []*emit.TypeRef
}

// PairExpr is a probe.Pair of two expressions of the types Types.
type PairExpr struct {
	Types  [2]*emit.TypeRef
	Values [2]Expr
}

// Closure is a function that a derived check passes to its probe. It takes
// the context of the call and Param, calls each builder of Builds, declares
// Lets, and calls the callable Method on Receiver, or the subject itself
// where Method is empty, with Args.
// It assigns the results of the call to Assign and returns Results. Without
// Assign it returns the results of the call, and with Discard it makes the
// call a statement of its own.
//
// Returns are the types of its results before the error, which it returns
// where Error is set, and before the flag, which it returns where Flag is
// set. Pair reports that its two types are the two values of a probe.Pair.
// Later are the assignments that it makes on each call after its first.
// Stamps is the field of each element of the result that it returns in place
// of the elements. Drain is the drain of a stream that the callable returns,
// which the closure returns as a list of the types Returns: drainElements,
// drainUntilError or drainPairs. A closure that drains a stream returns the
// stream Results[0] after the error Results[1] of the call. Pull reports a
// closure that returns a pull iteration of the stream in place of its drain,
// which yields the first value of each element of the stream. Channel reports
// a closure that returns the channel Results[0] of the call as a channel of
// the values Returns that the closure receives from. Absent reports a closure
// of a read with a presence flag, which returns fs.ErrNotExist beside its
// value where the flag reports no value.
type Closure struct {
	Param    *ClosureParam
	Receiver *Expr
	Method   string
	Stamps   string
	Drain    string
	Returns  []*emit.TypeRef
	Builds   []Build
	Lets     []Let
	Later    []Let
	Args     []Expr
	Assign   []string
	Results  []Expr
	Discard  bool
	Error    bool
	Flag     bool
	Pair     bool
	Pull     bool
	Channel  bool
	Absent   bool
}

// ClosureParam is the parameter of a closure after the context. Hook and
// Body report the hook func() of a registrar and the body probe.Act of a
// scope, and Type is the type of any other parameter. Other is the name of a
// second parameter of the type Type, for a closure that combines two values,
// and empty for a closure of one parameter.
type ClosureParam struct {
	Type  *emit.TypeRef
	Name  string
	Other string
	Hook  bool
	Body  bool
}

// Let is a variable that a closure declares before its call, or the
// assignment of its field Select where Select is set.
type Let struct {
	Name   string
	Select string
	Value  Expr
}

// form is the Go form of one input of a probe.
type form int

// The forms of the inputs of the probes, by RFC-0027 of the Go probes.
const (
	// formAct calls the callable for its error.
	formAct form = iota
	// formCall returns the value, or the pair of values, beside the error.
	formCall
	// formLookup returns the value, or the pair of values, beside the flag.
	formLookup
	// formWritten passes a sample to the written input and returns it.
	formWritten
	// formRewritten is formWritten with the first sample on its first call
	// and the second sample on every later call.
	formRewritten
	// formActs calls the callable for its error with the first sample on its
	// first call and the second sample on every later call, so two writes
	// add two elements.
	formActs
	// formTaking takes a value for the first input and returns the error.
	formTaking
	// formInverse takes a value for the first input and returns the value
	// beside the error.
	formInverse
	// formRoundTrip returns the sample of the first input and the result.
	formRoundTrip
	// formHook takes the hook for the first input.
	formHook
	// formScope takes the body for the first input.
	formScope
	// formHandle takes the handle of a call for the first input, or calls
	// the method of the handle that a param names.
	formHandle
	// formDelay passes a sample to the first input and returns it.
	formDelay
	// formStamps returns the stamp of each entry of the callable.
	formStamps
	// formZero is the zero value, or the pair of zero values, of the
	// results of another input.
	formZero
	// formValue is the value of a param.
	formValue
	// formName is the name of the declaration that a param refers to, as a
	// string, such as the callable that the param open of the mixin leak-free
	// names.
	formName
	// formInt is the int of a param.
	formInt
	// formDuration is the duration of a param.
	formDuration
	// formProjection returns the field that a param names of the result of
	// another input.
	formProjection
	// formMatched calls a writer that a predicate guards for its error, with
	// the first sample on its first call and the field of the param field of
	// the second sample on every later call.
	formMatched
	// formMatch calls the predicate of a guarded writer with the value of the
	// later calls of formMatched, and returns its verdict beside the error.
	formMatch
	// formCombine takes values for the first two inputs and returns the
	// value beside the error.
	formCombine
	// formInput is a value that a generated binding draws for an input of
	// the callable, by the position of the input and the index of the value.
	formInput
	// formBatch returns the values of a batch read beside the error, with
	// the key of each of the two writes of its pair.
	formBatch
	// formPosition takes a position for each integer input and returns the
	// error.
	formPosition
	// formActions is the function of the actions of the declaration.
	formActions
	// formOperation is the name of the callable of the input, as the history
	// records it.
	formOperation
	// formApply calls the callable of the input with the arguments of a call
	// of another input.
	formApply
	// formReference is the function that applies an operation to a reference
	// by the name of its callable.
	formReference
	// formWrittenArg returns the written value of a call of another input.
	formWrittenArg
	// formStamp returns the field that a param names of an element of the
	// result of another input.
	formStamp
	// formCountable returns the use that the callable of the input returns.
	formCountable
	// formCount is the int of a param, or the default of its binding.
	formCount
	// formTotal returns the total that the callable of the input reads,
	// beside the error.
	formTotal
	// formPull starts a pull iteration of the stream that the callable of the
	// input returns.
	formPull
	// formChannel returns the channel that the callable of the input returns,
	// beside the error.
	formChannel
	// formNext is the method that a param names on the handle that the
	// callable of another input returns. The method takes a wait and returns
	// an update beside a presence flag.
	formNext
	// formStop is the method that a param names on the handle that the
	// callable of another input returns. The method takes and returns
	// nothing.
	formStop
	// formElementArg returns the argument of a call of the callable of the
	// input that has the type of the elements of another input: the first
	// values of a stream, the values of a channel or the updates of a handle.
	formElementArg
	// formKey returns the key of a call of the callable of the input: the
	// arguments of the inputs that it shares by name and type with the
	// callable of another input.
	formKey
	// formClients is the number of the clients of a concurrent section after
	// the first, as stateful.Clients counts them.
	formClients
	// formBudget is the most steps of the check of a history.
	formBudget
	// formComputing is formCall with a computation for each input of a
	// function type, which records a call on a dependency of the injector.
	formComputing
	// formDurable is the write of the callable and the read of its partner
	// for each sample, with the value and the key that each sample writes.
	formDurable
	// formRequestField is the field that a param names of the request of
	// another input, the input that has such a field: its first sample in a
	// sample binding, and the field of the drawn request in a generated
	// binding.
	formRequestField
	// formFielded is formCall whose sample binding sets the field that the
	// param of another input names in the request of the call to its first
	// sample, so the field of the request is not its zero value.
	formFielded
	// formWorld is the world of the subject of the declaration.
	formWorld
)

// concurrentClients is the number of the clients of a concurrent section
// after the first, so three clients call the subject at once. historyBudget
// is the most steps that the checker of a history takes.
const (
	concurrentClients = 2
	historyBudget     = 200_000
)

// slot is one input of a probe: its name in the entry of the spec, its form,
// and the input whose results a zero or a projection reads. A fresh slot
// calls its callable on a new subject of the factory. gate is the gate that
// the suite generator judges on the input, which refuses the check where the
// callable of the input does not fit the form of the input.
type slot struct {
	name  string
	of    string
	gate  spec.GateID
	form  form
	fresh bool
}

// The names of the inputs of the probes that more than one probe takes.
const (
	inCall      = "call"
	inZero      = "zero"
	inSentinel  = "sentinel"
	inRead      = "read"
	inWrite     = "write"
	inObserve   = "observe"
	inPut       = "put"
	inDrain     = "drain"
	inAttempts  = "attempts"
	inClose     = "close"
	inReplay    = "replay"
	inAppend    = "append"
	inStream    = "stream"
	inWritten   = "written"
	inSubscribe = "subscribe"
	inWatch     = "watch"
	inTrigger   = "trigger"
)

// The inputs that many probes share.
var (
	durable  = slot{name: "durable", form: formDurable}
	act      = slot{name: inCall, form: formAct}
	call     = slot{name: inCall, form: formCall}
	zero     = slot{name: inZero, of: inCall, form: formZero}
	sentinel = slot{name: inSentinel, form: formValue}
	read     = slot{name: inRead, form: formCall}
	drain    = slot{name: inDrain, form: formCall}
	writes   = slot{name: inWrite, form: formActs}
	put      = slot{name: inPut, form: formAct}
	observe  = slot{name: inObserve, form: formCall}
	attempts = slot{name: inAttempts, form: formInt}
	clients  = slot{name: "clients", form: formClients}
	world    = slot{name: "world", form: formWorld}
)

// probes contains the inputs of each probe that a suite calls, in the order
// of the parameters of its Go function after the handle. The probes of the
// enumerations and the errors belong to the generators enum and sentinel.
var probes = map[spec.ProbeID][]slot{
	spec.Survives:        {act},
	spec.UnderCancelled:  {act},
	spec.ZeroBesideError: {call, zero},
	spec.FieldBesideError: {
		{name: inCall, of: "want", form: formFielded},
		{name: "field", of: inCall, form: formProjection},
		{name: "want", of: inCall, form: formRequestField, gate: spec.ClassificationParams},
	},
	spec.ZeroPairBesideError:   {call, zero},
	spec.ZeroBesideFalse:       {{name: inCall, form: formLookup}, zero},
	spec.ZeroPairBesideFalse:   {{name: inCall, form: formLookup}, zero},
	spec.NilBesideAMiss:        {call},
	spec.ZeroBesideAMiss:       {call, zero},
	spec.RepeatAgrees:          {call},
	spec.SurvivesRepeat:        {act},
	spec.ReportsSentinel:       {call, sentinel},
	spec.RepeatIsQuiet:         {act},
	spec.NeverDecreases:        {call},
	spec.RepeatAscends:         {call},
	spec.EscapesTheUnsafe:      {{name: inCall, of: "unsafe", form: formCall}, {name: "unsafe", form: formValue}},
	spec.AcceptsTheEdge:        {{name: inCall, form: formTaking}, {name: "edge", form: formValue}},
	spec.RemainsPoisoned:       {act, {name: "induce", form: formAct}},
	spec.SucceedsWithin:        {act, attempts},
	spec.RefusesTheSecondWrite: {act, sentinel},
	spec.RefusesPastTheBurst:   {act, {name: "burst", form: formInt}, sentinel},
	spec.RefusesEarly:          {act, sentinel},
	spec.WriteThenRead:         {read, {name: inWrite, form: formWritten}},
	spec.PinnedAcrossWrite:     {read, {name: inWrite, form: formAct}},
	spec.SecondWriteReplaces:   {{name: inCall, form: formRewritten}, read},
	spec.KeepsPartitionsApart:  {{name: inWrite, form: formRewritten}, read},
	spec.RoundTrips:            {{name: inCall, form: formRoundTrip}, {name: "inverse", form: formInverse}},
	spec.WorksAfter:            {act, {name: "ready", form: formAct}},
	spec.ReportsSealed:         {act, {name: inClose, form: formAct}, sentinel},
	spec.SealsAfterClose: {
		call, {name: inClose, of: inCall, form: formHandle}, {name: "after", of: inCall, form: formHandle}, sentinel,
	},
	spec.MovesObservation:    {act, observe},
	spec.EachCallCompounds:   {act, observe},
	spec.SettlesThenObserves: {act, {name: "settle", form: formAct}, observe},
	spec.DetectsTamper:       {act, {name: "verify", form: formAct}, {name: "tamper", form: formAct}},
	spec.ConvergesAfter:      {act, {name: "begin", form: formAct}, attempts},
	spec.RunsTheHook:         {act, {name: "register", form: formHook}},
	spec.RefusesTheInvalid: {
		{name: inCall, form: formTaking}, {name: "validate", form: formTaking}, {name: "invalid", form: formValue},
	},
	spec.RefusesWhatThePredicateRejects: {{name: inCall, form: formMatched}, {name: "match", form: formMatch}},
	spec.ResignsLeadership:              {act, {name: "resign", form: formAct}, {name: "leading", form: formCall}},
	spec.RollsBackOnFailure: {
		{name: inCall, form: formScope}, {name: inWrite, form: formAct}, {name: inRead, form: formAct}, sentinel,
	},
	spec.AgreesWithTheBacking: {call, {name: "backing", form: formCall}, put},
	spec.TakesBackWhatItLent:  {call, {name: inPut, form: formTaking}},
	spec.ContinuesFromTheCursor: {
		{name: inCall, form: formInverse},
		{name: "cursor", of: inCall, form: formProjection},
	},
	spec.CompensatesTheStep: {act, {name: "compensate", form: formAct}},
	spec.RefusedLeavesNothing: {
		{name: inCall, form: formTaking},
		{name: "reader", form: formCall},
		{name: "refused", form: formValue},
		{name: inZero, of: "reader", form: formZero},
	},
	spec.RemovesOnDelete:          {act, put, {name: inRead, form: formAct}, sentinel},
	spec.NoElementTwice:           {drain, writes},
	spec.DrainsInOneOrder:         {drain, writes},
	spec.GainsWhatWentIn:          {drain, writes},
	spec.KeepsWhatWentIn:          {drain, writes},
	spec.ReadsNeverRegress:        {read, {name: inWrite, form: formAct}, version},
	spec.WritesAdvance:            {read, {name: inWrite, form: formCall}, version},
	spec.ReplaysInStampOrder:      {{name: inCall, form: formActs}, {name: inReplay, form: formStamps}},
	spec.ReplayGrows:              {act, {name: inReplay, form: formCall}},
	spec.ReturnsWithin:            {act, {name: "budget", form: formDuration}},
	spec.CountsWithinWindow:       {act, {name: "count", form: formCall}, {name: "window", form: formDuration}},
	spec.RunsOnSchedule:           {{name: inCall, form: formDelay}, {name: "count", form: formCall}},
	spec.MovesWithTheClock:        {call, put},
	spec.SteadyWhileTheClockStops: {call, put},
	spec.ExpiresAfter: {
		read, {name: inWrite, form: formAct}, {name: "lifetime", form: formDuration}, sentinel,
	},
	spec.ReportsAFailedDependency: {act},
	spec.FailsWithoutEffect:       {act, read},
	spec.SucceedsAfterAFailure:    {act, attempts},
	spec.RollsBackUnderFaults: {
		{name: inCall, form: formScope}, {name: inWrite, form: formAct}, {name: inRead, form: formAct}, sentinel,
	},
	spec.RefusesWithoutTheDependency: {act},
	spec.UsesOneInstance:             {act},
	spec.ReleasesWhatItAcquires:      {act, {name: "open", form: formName}, {name: inClose, form: formName}},
	spec.EitherOrderAgrees:           {combined, firstInput, secondInput},
	spec.EitherGroupingAgrees:        {combined, firstInput, secondInput, {name: "third", form: formInput}},
	spec.ReadsWhatWasStored:          {{name: inCall, form: formBatch}, {name: inWrite, form: formRewritten}},
	spec.ReportsNoError:              {act},
	spec.BoundedByTheSize:            {{name: inCall, form: formPosition}, {name: "size", form: formCall}},
	spec.AgreesWithTheSibling:        {call, {name: "sibling", form: formCall, fresh: true}},
	spec.VerifiesAfterAppends:        {act, {name: inAppend, form: formActs}},
	spec.GrantsAfterRelease:          {act, {name: "release", form: formAct}},
	spec.MovesWithALoan:              {call, {name: "get", form: formAct}},
	spec.ResultsNeverFall:            {actions, readOp},
	spec.ReadVersionsNeverFall:       {actions, readOp, version},
	spec.DrainsAgreeBetweenWrites:    {actions, drainOp},
	spec.WrittenVersionsRise:         {actions, writeOp, {name: inRead, of: inWrite, form: formApply}, version},
	spec.ReadsSeeTheLastWrite:        {actions, readOp, writeOp, readback, version},
	spec.WritesRiseAboveReads:        {actions, writeOp, readOp, readback, version},
	spec.PositionsRise:               {actions, {name: inAppend, form: formOperation}},
	spec.ReplaysFollowTheStamps: {
		actions, {name: inReplay, form: formOperation}, {name: "stamp", of: inReplay, form: formStamp},
	},
	spec.DrainsWithoutRepeats:   {actions, drainOp},
	spec.DrainsContainTheWrites: {actions, drainOp, writeOp, written},
	spec.DrainsGainTheWrites:    {actions, drainOp, writeOp, written},
	spec.EachCallMovesTheObservation: {
		actions, {name: inCall, form: formOperation}, {name: inObserve, of: inCall, form: formApply},
	},
	spec.ReadsReturnTheLastWrite:        {actions, writeOp, readback, written},
	spec.ReadsAgreeWithTheBacking:       {actions, readOp, {name: "backing", of: inRead, form: formApply}},
	spec.SequencesAgreeWithTheReference: {actions, {name: "reference", form: formReference}},
	spec.LinearizableHistory: {
		actions,
		{name: "reference", form: formReference},
		clients,
		{name: "budget", form: formBudget},
	},
	spec.SerializableHistory:          {actions, clients},
	spec.SnapshotHistory:              {actions, clients},
	spec.SurvivesEveryCrashPoint:      {durable},
	spec.WholeAtEveryCrashPoint:       {durable},
	spec.ReportsEachFlippedBit:        {durable, {name: "corrupt", form: formValue}},
	spec.RefusesOnAFullMedium:         {durable},
	spec.FaultyWorldHistory:           {world},
	spec.WorldRecovery:                {world},
	spec.ConcurrentReadsAgree:         {actions, readOp, {name: "again", of: inRead, form: formApply}, clients},
	spec.SharesOneComputation:         {{name: inCall, form: formComputing}, clients},
	spec.ExclusiveUnderConcurrency:    {act, sentinel, clients},
	spec.PinnedUnderConcurrentWrites:  {actions, readOp, writeOp, clients},
	spec.VersionsFollowTheCausalOrder: {actions, readOp, writeOp, version, clients},
	spec.UseIsWithinTheBounds: {
		actions,
		{name: inObserve, form: formCountable, gate: spec.CountableUse},
		{name: "min", form: formCount},
		{name: "limit", form: formCount},
	},
	spec.MergesCommute: {
		actions, {name: "apply", form: formReference}, {name: "merge", form: formOperation}, writeOp, readOp,
	},
	spec.StatesMoveAlongTheGraph: {
		actions,
		{name: "fn", form: formOperation},
		{name: inObserve, of: "fn", form: formApply},
		{name: "transitions", form: formValue},
	},
	spec.MovesKeepTheTotal: {
		actions, {name: inCall, form: formOperation}, {name: inObserve, form: formTotal, gate: spec.NumericTotal},
	},
	spec.FailedCallsLeaveNoTrace: {
		actions, {name: inCall, form: formOperation}, {name: "readback", of: inCall, form: formApply},
	},
	spec.IterationsSeeTheMutations: {
		actions,
		{name: inStream, form: formPull, gate: spec.StreamMutations},
		{name: "mutate", form: formOperation},
		{name: "delete", form: formOperation},
		{name: inWritten, of: inStream, form: formElementArg, gate: spec.StreamMutations},
		{name: "removed", of: inStream, form: formElementArg, gate: spec.StreamMutations},
	},
	spec.AppendsArrive:          {actions, {name: inAppend, form: formOperation}, subscribe, delivered},
	spec.DeliveriesMatchTheMode: {actions, publishOp, subscribe, delivered, {name: "mode", form: formName}},
	spec.RedeliveriesArrive: {
		actions,
		publishOp,
		{name: "redeliver", form: formOperation},
		subscribe,
		delivered,
		{name: "redelivered", of: inSubscribe, form: formElementArg, gate: spec.DeliversOnAChannel},
	},
	spec.WatchesSeeTheTriggers: {
		actions,
		{name: inWatch, form: formOperation},
		{name: inTrigger, form: formOperation},
		{name: "next", of: inWatch, form: formNext, gate: spec.SubscriptionMethods},
		{name: "stop", of: inWatch, form: formStop, gate: spec.SubscriptionMethods},
		{name: "retrigger", of: inTrigger, form: formApply},
		{name: "watch-key", of: inTrigger, form: formKey},
		{name: "trigger-key", of: inWatch, form: formKey},
		{name: "triggered", of: "next", form: formElementArg, gate: spec.SubscriptionMethods},
	},
}

// The inputs that the probes of deliveries share: the subscription, the
// operation of a publish, and the value that a call of a callable passes to
// the subscription.
var (
	subscribe = slot{name: inSubscribe, form: formChannel, gate: spec.DeliversOnAChannel}
	publishOp = slot{name: "publish", form: formOperation}
	delivered = slot{name: inWritten, of: inSubscribe, form: formElementArg, gate: spec.DeliversOnAChannel}
)

// The inputs that the probes of generated sequences share: the actions, the
// operations of a read, a write and a drain, the version of a read, the read
// after a write, and the written value of a write.
var (
	actions  = slot{name: "actions", form: formActions}
	readOp   = slot{name: inRead, form: formOperation}
	writeOp  = slot{name: inWrite, form: formOperation}
	drainOp  = slot{name: inDrain, form: formOperation}
	version  = slot{name: "version", of: inRead, form: formProjection}
	readback = slot{name: "readback", of: inWrite, form: formApply}
	written  = slot{name: inWritten, of: inWrite, form: formWrittenArg}
)

// The inputs that the probes of a callable that combines two values share.
var (
	combined    = slot{name: inCall, form: formCombine}
	firstInput  = slot{name: "first", form: formInput}
	secondInput = slot{name: "second", form: formInput}
)

// bubbled are the probes that run inside a testing/synctest bubble: the
// probes that advance the clock of the handle, and the probes of deliveries
// and of one computation, which wait until each goroutine of the bubble is
// blocked.
var bubbled = []spec.ProbeID{
	spec.ReturnsWithin, spec.CountsWithinWindow, spec.RunsOnSchedule, spec.MovesWithTheClock,
	spec.SteadyWhileTheClockStops, spec.ExpiresAfter, spec.AppendsArrive, spec.DeliveriesMatchTheMode,
	spec.RedeliveriesArrive, spec.WatchesSeeTheTriggers, spec.SharesOneComputation,
}

// cased are the probes that take the case of a property check without
// generated sequences: the probes whose clients the scheduler of the case runs.
var cased = []spec.ProbeID{spec.ExclusiveUnderConcurrency}

// goName returns the name of the constant of id in the package spec of the
// Go runtime, as the naming table of the spec spells it.
func goName[T ~string](id T) string { return naming.Pascal(string(id)) }

// identityOf returns the identity of the derived check c among the callables
// of s: its field, its method and its segment. A claim that binds once for
// each instance of a contract has the id of the instance after its segment.
func identityOf(s *Suite, c *derive.Check) (Callable, string, string) {
	segment := c.Identity[strings.Index(c.Identity, "/")+1:]
	i := slices.IndexFunc(s.Facts, func(f derive.Facts) bool { return f.Callable == c.Facts.Callable })
	return s.Callables[i], methodOf(c.Identity), segment
}

// methodOf returns the name of the method that returns the identity of a
// derived check: the segment of the identity in Pascal case, with the id of an
// instance as a word of its own, as Expires0 for Store.Get/expires#0.
func methodOf(identity string) string {
	segment := identity[strings.Index(identity, "/")+1:]
	return naming.Pascal(strings.ReplaceAll(segment, "#", "-"))
}

// bases returns the method of the identity of the sample binding of each
// generated binding among checks, by the identity of the generated binding.
// The map has no entry for the identity of any other check.
func bases(checks []derive.Check) map[string]string {
	out := make(map[string]string, len(checks))
	for i := range checks {
		out[derive.GeneratedOf(checks[i].Identity)] = methodOf(checks[i].Identity)
	}
	return out
}

// derived records into s the derived check of each check of the derivation of
// s, with its planted defect. A generated binding draws the inputs of its
// calls by the rules of the call generators of s.
func (k *caller) derived(s *Suite) {
	base := bases(s.Derivation.Checks)
	for i := range s.Derivation.Checks {
		c := &s.Derivation.Checks[i]
		callable, method, segment := identityOf(s, c)
		d := Derived{
			Field: callable.Field, Type: callable.Type, Method: method, Base: base[c.Identity], Segment: segment,
			Statement: c.Statement, Claim: goName(c.Claim.ID), Probe: goName(c.Claim.Probe), Shallow: c.Shallow,
			Bubble: slices.Contains(bubbled, c.Claim.Probe), Generated: c.Generated,
			World: slices.Contains(worldProbes, c.Claim.Probe),
		}
		d.Cased = d.World || slices.Contains(cased, c.Claim.Probe)
		for _, m := range c.Machineries() {
			d.Machineries = append(d.Machineries, goName(m))
		}
		for _, e := range s.Derivation.Edges {
			if e.Deeper != c.Identity {
				continue
			}
			j := slices.IndexFunc(s.Derivation.Checks, func(o derive.Check) bool { return o.Identity == e.Shallower })
			shallow, shallowMethod, _ := identityOf(s, &s.Derivation.Checks[j])
			d.Entails = append(d.Entails, Identity{
				Field: shallow.Field, Method: shallowMethod, Base: base[e.Shallower], ID: e.Shallower,
			})
		}
		b := binder{k: k, check: c, s: s}
		if c.Claim.Probe == spec.RedeliveriesArrive && !b.redelivers() {
			b.refused = spec.AtLeastOnceMode
		}
		if c.Generated {
			b.drawn = k.newDraws(s.Facts, s.Calls)
		}
		d.Args = b.args()
		if c.Generated {
			d.Draws = used(d.Args)
		}
		if sequenced(c) && !d.World {
			d.Sequence, d.TypeArgs = true, b.typeArgs()
		}
		if c.Claim.Scope == scopeDeclaration {
			d.Whole = c.Identity
		}
		switch {
		case c.Claim.Probe == spec.MergesCommute && !k.replicable(s, c.Facts):
			b.refused = spec.TakesAReplica
		case c.Claim.Probe == spec.StatesMoveAlongTheGraph && !k.graphed(b.param(transitionsParam)):
			b.refused = spec.GraphOfStates
		default:
		}
		if b.refused != "" {
			s.Derivation.Refusals = append(s.Derivation.Refusals, derive.Refuse(c, b.refused))
			continue
		}
		defect, known := k.defect(s, c, &b)
		if defect.Body == bodyLiveHandle && defect.Handle == nil {
			s.Derivation.Refusals = append(s.Derivation.Refusals, derive.Refuse(c, spec.InterfaceHandle))
			continue
		}
		if !k.referableArgs(d.Args) || !k.referableDefect(defect) {
			s.Derivation.Refusals = append(s.Derivation.Refusals, derive.Refuse(c, spec.ImportableSamples))
			continue
		}
		if !known {
			s.Defects = append(s.Defects, defect)
		}
		d.Defect, d.Door = defect.Name, defect.Door
		s.Derived = append(s.Derived, d)
	}
}

// referableArgs reports whether the file of the suite can refer to each
// sample that a closure of args passes or declares. The other values of a
// check are referable wherever its declaration is: a zero value is a literal
// of a type of the callable, and a param names a declaration in the scope of
// the declaration. A read selects a field of the value that its write passes
// or declares, so the sample of the write covers the field.
func (k *caller) referableArgs(args []Arg) bool {
	var exprs []Expr
	for i := range args {
		if c := args[i].Closure; c != nil {
			exprs = append(exprs, c.Args...)
			lets := slices.Concat(c.Lets, c.Later)
			for j := range lets {
				exprs = append(exprs, lets[j].Value)
			}
		}
		if a := args[i].Apply; a != nil {
			exprs = append(exprs, a.Args...)
		}
		if c := args[i].Countable; c != nil {
			exprs = append(exprs, c.Args...)
		}
	}
	return !slices.ContainsFunc(exprs, func(e Expr) bool { return !k.referableExpr(&e) })
}

// referableExpr reports whether the file of the suite can refer to the value
// of the expression e, and reports true for an expression without one.
func (k *caller) referableExpr(e *Expr) bool { return e.Value == nil || k.referable(e.Value) }

// referableDefect reports whether the file of the suite can refer to the
// sample of each value that the planted defect d and its overrides return.
// The zero value of a type is a literal of the type, which the file can refer
// to with the type.
func (k *caller) referableDefect(d *Defect) bool {
	return !slices.ContainsFunc(append([]*Defect{d}, d.Overrides...), func(o *Defect) bool {
		return slices.ContainsFunc(o.Results, func(r Outcome) bool {
			return r.Sample != nil && !k.referableExpr(r.Sample)
		})
	})
}

// binder binds the inputs of a derived check of the suite s to the arguments
// of its probe. drawn are the draws of a generated binding, and nil for a
// sample binding. refused is the gate that the suite generator judges on the
// code of the check where the gate fails, and empty otherwise.
type binder struct {
	k       *caller
	s       *Suite
	check   *derive.Check
	drawn   *draws
	refused spec.GateID
}

// scopeDeclaration is the scope of a claim on a whole declaration, and
// transitionsParam the param of the contract workflow that names its graph.
const (
	scopeDeclaration = "declaration"
	transitionsParam = "transitions"
)

// typeArgs returns the type arguments after the subject type of the probe of a
// check of generated sequences that the arguments of the probe do not infer:
// the type of the value of a read or an append, and the type of the elements
// of a drain.
func (b *binder) typeArgs() [][]*emit.TypeRef {
	switch b.check.Claim.Probe {
	case spec.ResultsNeverFall:
		f, _ := b.facts(inRead)
		return [][]*emit.TypeRef{{b.k.typeOf(b.k.resultRef(f, 0))}}
	case spec.PositionsRise:
		f, _ := b.facts(inAppend)
		return [][]*emit.TypeRef{{b.k.typeOf(b.k.resultRef(f, 0))}}
	case spec.DrainsWithoutRepeats:
		f, _ := b.facts(inDrain)
		return [][]*emit.TypeRef{b.k.elementOf(f)}
	default:
		return nil
	}
}

// pairOf returns the pair of the check whose read is target, and nil where no
// pair reads with target. The write of that pair is the callable whose
// arguments a probe applies target to, because each claim binds an input of
// the kind apply to the partner that it pairs with the input of the call.
func (b *binder) pairOf(target *derive.Facts) *derive.Pair {
	for _, p := range b.check.Claim.Partners {
		if pair := b.check.Pairs[p.Name]; same(pair.Read, target) {
			return pair
		}
	}
	return nil
}

// args returns the argument of the probe of the check for each input of the
// probe, in the order of the parameters of its Go function after the handle.
// A gate that refuses an input refuses the check, so args binds no input
// after it.
func (b *binder) args() []Arg {
	slots := probes[b.check.Claim.Probe]
	out := make([]Arg, 0, len(slots))
	for i := range slots {
		out = append(out, b.arg(&slots[i]))
		if b.refused != "" {
			break
		}
	}
	return out
}

// refuse sets the gate of the input sl as refused, and returns no argument.
func (b *binder) refuse(sl *slot) Arg {
	b.refused = sl.gate
	return Arg{}
}

// binding returns the binding of the input name of the claim of the check.
func (b *binder) binding(name string) *rules.Binding {
	i := slices.IndexFunc(b.check.Claim.Inputs, func(in rules.Binding) bool { return in.Input == name })
	return &b.check.Claim.Inputs[i]
}

// facts returns the facts of the callable of the input name: the callable of
// the check, or the partner that the binding names.
func (b *binder) facts(name string) (*derive.Facts, *derive.Pair) {
	in := b.binding(name)
	if in.Kind == rules.FromPartner {
		pair := b.check.Pairs[in.Partner]
		return pair.Partner, pair
	}
	return b.check.Facts, b.ownPair()
}

// ownPair returns the pair of the check in which the callable of the check
// reads or writes, and nil where it does neither. The partners of a claim are
// visited in the order of the claim, so one pair is chosen for every run.
func (b *binder) ownPair() *derive.Pair {
	for _, p := range b.check.Claim.Partners {
		pair := b.check.Pairs[p.Name]
		if same(pair.Read, b.check.Facts) || same(pair.Write, b.check.Facts) {
			return pair
		}
	}
	return nil
}

// same reports whether a is the facts of the callable of b. a is nil for the
// read and the write of a pair where neither writes.
func same(a, b *derive.Facts) bool { return a != nil && a.ID == b.ID }

// param returns the value of the param key of the classifications of the
// callable of the check, and nil where none sets it.
func (b *binder) param(key string) any { return paramOf(b.check.Facts, key) }

// paramOf returns the value of the param key of the classifications of the
// callable of f, and nil where none sets it.
func paramOf(f *derive.Facts, key string) any {
	for _, params := range f.Class.Params {
		if v, ok := params[directive.ParamKey(key)]; ok {
			return v
		}
	}
	return nil
}

// redelivers reports whether the publish of the check states the delivery
// mode at-least-once, the one mode under which a redelivery delivers an event
// once more. A second delivery breaks the other modes, so their check of the
// mode judges the redeliveries.
func (b *binder) redelivers() bool {
	publish, _ := b.facts(publishOp.name)
	return nameOf(paramOf(publish, modeParam)) == atLeastOnce
}

// arg returns the argument of the probe for the input sl.
func (b *binder) arg(sl *slot) Arg {
	switch sl.form {
	case formZero:
		f, _ := b.facts(sl.of)
		return Arg{Expr: b.k.zeroOf(f)}
	case formValue:
		return Arg{Expr: b.k.valueOfParam(b.param(b.binding(sl.name).Param))}
	case formName:
		name := nameOf(b.param(b.binding(sl.name).Param))
		return Arg{Expr: &Expr{Value: new(emit.Literal(emit.LiteralString, name))}}
	case formInt:
		return Arg{Expr: &Expr{Word: intOf(b.param(b.binding(sl.name).Param))}}
	case formClients:
		return Arg{Expr: &Expr{Word: strconv.Itoa(concurrentClients)}}
	case formBudget:
		return Arg{Budget: historyBudget}
	case formDuration:
		d, _ := time.ParseDuration(nameOf(b.param(b.binding(sl.name).Param)))
		return Arg{Duration: d}
	case formProjection:
		f, _ := b.facts(sl.of)
		field := nameOf(b.param(b.binding(sl.name).Param))
		ref := b.k.resultRef(f, 0)
		p := &Projection{Type: b.k.typeOf(ref), Result: b.k.typeOf(ref), Field: field}
		if field != "" {
			p.Result = b.k.fieldRef(ref, field)
		}
		return Arg{Projection: p}
	case formRequestField:
		f, _ := b.facts(sl.of)
		field := nameOf(b.param(b.binding(sl.name).Param))
		at := requestAt(f, field)
		if at < 0 {
			return b.refuse(sl)
		}
		if b.drawn == nil {
			first, _, _ := b.k.fieldSamples(f, at, field)
			return Arg{Expr: &first}
		}
		request := b.k.argOf(b.drawn, f, at, 0)
		return Arg{Expr: &Expr{Of: &request, Select: field}}
	case formInput:
		bound := b.binding(sl.name)
		v := b.drawn.value(b.check.Facts, bound.Position, bound.Sample)
		return Arg{Expr: &v}
	case formActions, formOperation, formApply, formReference, formWrittenArg, formStamp, formCountable, formCount,
		formNext, formStop, formElementArg, formKey:
		return b.sequenceArg(sl)
	case formDurable:
		_, pair := b.facts(sl.name)
		return Arg{Durable: b.k.durableOf(b.drawn, pair)}
	case formWorld:
		return Arg{Ref: b.s.World.Name}
	case formTotal:
		// Only a number has a scalar class.
		if f, _ := b.facts(sl.name); len(f.Results) == 0 || f.Results[0].Shape.Class == 0 {
			return b.refuse(sl)
		}
	case formPull:
		if f, _ := b.facts(sl.name); streamAt(f) < 0 {
			return b.refuse(sl)
		}
	case formChannel:
		if f, _ := b.facts(sl.name); !b.k.subscribes(f) {
			return b.refuse(sl)
		}
	default:
	}
	in := &input{sl: sl, drawn: b.drawn}
	switch sl.form {
	case formHandle:
		call, _ := b.facts(sl.of)
		in.handle = b.k.resultRef(call, 0)
		if bound := b.binding(sl.name); bound.Kind == rules.FromParam {
			in.method = nameOf(b.param(bound.Param))
			return Arg{Closure: b.k.closure(in)}
		}
	case formStamps:
		// The stamp of an entry is the field that the param version of the
		// mixin causal of the callable names.
		in.stamps = nameOf(b.param(stampParam))
	case formCall:
		if sl.of != "" {
			in.with = b.k.valueOfParam(b.param(b.binding(sl.of).Param))
		}
	case formMatched, formMatch:
		in.changed = nameOf(b.param(fieldParam))
	case formFielded:
		in.field = nameOf(b.param(b.binding(sl.of).Param))
	default:
	}
	in.f, in.pair = b.facts(sl.name)
	return Arg{Closure: b.k.closure(in)}
}

// requestAt returns the position of the request of the callable of f: the
// first input whose type has the field name. It returns -1 where no input has
// the field.
func requestAt(f *derive.Facts, name string) int {
	return slices.IndexFunc(f.Inputs, func(v derive.Value) bool {
		return slices.ContainsFunc(v.Fields, func(fd derive.Field) bool { return fd.Name == name })
	})
}

// sequenceArg returns the argument of a probe of generated sequences for the
// input sl. It sets the gate of sl as refused for an input that does not fit
// its form, such as an observer whose value is no use.
func (b *binder) sequenceArg(sl *slot) Arg {
	q := b.s.Sequences
	switch sl.form {
	case formActions:
		return Arg{Ref: q.Actions}
	case formReference:
		return Arg{Ref: q.Reference}
	case formOperation:
		f, _ := b.facts(sl.name)
		return Arg{Expr: &Expr{Value: new(emit.Literal(emit.LiteralString, f.Callable))}}
	case formApply:
		target, _ := b.facts(sl.name)
		source, _ := b.facts(sl.of)
		return Arg{Apply: b.k.applyOf(b.s, q, target, source, b.pairOf(target))}
	case formWrittenArg:
		// A written value is an input of the partner of a pair. The first
		// input stands in for a pair without one.
		write, pair := b.facts(sl.of)
		index := max(pair.Written, 0)
		return Arg{Written: &Written{Type: b.k.typeOf(b.k.inputRef(write, index)), Index: index}}
	case formStamp:
		f, _ := b.facts(sl.of)
		field := nameOf(b.param(b.binding(sl.name).Param))
		ref := b.k.elementRefs(f)[0]
		return Arg{Projection: &Projection{Type: b.k.typeOf(ref), Result: b.k.fieldRef(ref, field), Field: field}}
	case formCountable:
		f, _ := b.facts(sl.name)
		c, ok := b.k.countableOf(b.s, q, f)
		if !ok {
			return b.refuse(sl)
		}
		return Arg{Countable: c}
	case formNext, formStop:
		m, _, ok := b.handleMethod(sl)
		if !ok {
			return b.refuse(sl)
		}
		return Arg{Method: m}
	case formElementArg:
		f, _ := b.facts(sl.name)
		elem := b.k.b.TypeOf(b.k.instance(b.element(sl.of)))
		at := slices.IndexFunc(f.Inputs, func(v derive.Value) bool { return reflect.DeepEqual(v.Shape, elem) })
		if at < 0 {
			return b.refuse(sl)
		}
		return Arg{Written: &Written{Type: b.k.typeOf(b.k.inputRef(f, at)), Index: at}}
	case formKey:
		f, _ := b.facts(sl.name)
		other, _ := b.facts(sl.of)
		return Arg{Key: &Key{Args: shared(f, other)}}
	default:
		bound := b.binding(sl.name)
		count := intOf(b.param(bound.Param))
		if count == "" {
			count = bound.Default
		}
		return Arg{Expr: &Expr{Word: count}}
	}
}

// The package and the name of the type of the wait of the method next of a
// handle.
const (
	timePackage  = "time"
	durationType = "Duration"
)

// handleMethod returns the method that the param of the input sl names on
// the handle that the callable of the input sl.of returns, as a method
// expression, with the facts of the method. It returns false where the param
// names no method of the handle, such as a field, and where the method does
// not fit the form of sl: next takes a wait and returns an update beside a
// presence flag, and stop takes and returns nothing.
func (b *binder) handleMethod(sl *slot) (*HandleMethod, *derive.Facts, bool) {
	watch, _ := b.facts(sl.of)
	handle := b.k.resultRef(watch, 0)
	ref := handle
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	decl, _ := b.k.m.Reader().Lookup(ref.Target)
	bound := b.k.b
	methods := methodsOf(&bound, decl)
	name := nameOf(b.param(b.binding(sl.name).Param))
	at := slices.IndexFunc(methods, func(m *node.Method) bool { return m.Name == name })
	if at < 0 {
		return nil, nil, false
	}
	f := b.k.factsOf(methods[at])
	s := &f.Signature
	fits := len(s.Params) == 0 && len(s.Returns) == 0
	if sl.form == formNext {
		fits = len(s.Params) == 1 && is(&s.Params[0].Ref.Target, timePackage, durationType) &&
			slices.ContainsFunc(s.Returns, func(r sdkrules.ReturnView) bool { return r.Role == sdkrules.ReturnOkBool })
	}
	return &HandleMethod{Type: b.k.typeOf(handle), Name: name}, &f, fits
}

// element returns the type of the elements of the input name of the probe of
// the check, which passed its gate: the first value of each element of a
// stream, the values of a channel, or the update that the method next of a
// handle returns.
func (b *binder) element(name string) *node.TypeRef {
	slots := probes[b.check.Claim.Probe]
	sl := &slots[slices.IndexFunc(slots, func(s slot) bool { return s.name == name })]
	f, _ := b.facts(name)
	switch sl.form {
	case formPull:
		return b.k.elementRefs(f)[0]
	case formChannel:
		return b.k.channelElem(f)
	default:
		_, next, _ := b.handleMethod(sl)
		return b.k.resultRef(next, 0)
	}
}

// subscribes reports whether the callable of f returns first a channel that a
// caller receives from.
func (k *caller) subscribes(f *derive.Facts) bool {
	if len(f.Results) == 0 {
		return false
	}
	under := k.underlying(k.resultRef(f, 0))
	return under.Form == symbol.FormStream && !sends(under.Spelling)
}

// channelElem returns the type of the values of the channel that the callable
// of f returns first.
func (k *caller) channelElem(f *derive.Facts) *node.TypeRef {
	return k.underlying(k.resultRef(f, 0)).Elems[0]
}

// shared returns the positions of the inputs of the callable of f that the
// callable of other takes under the same name, in the order of their names,
// so the key of a watch and the key of its trigger list the same inputs in one
// order. The derivation pairs a watch with its trigger only where the trigger
// takes each input of the watch with its name and its type.
func shared(f, other *derive.Facts) []int {
	var at []int
	for i := range f.Inputs {
		v := &f.Inputs[i]
		if slices.ContainsFunc(other.Inputs, func(o derive.Value) bool { return o.Name == v.Name }) {
			at = append(at, i)
		}
	}
	slices.SortFunc(at, func(x, y int) int { return strings.Compare(f.Inputs[x].Name, f.Inputs[y].Name) })
	return at
}

// stampParam is the param of the mixin causal that names the field of the
// stamp of an entry, and fieldParam the param of the contract if-match that
// names the field that the predicate judges.
const (
	stampParam = "version"
	fieldParam = "field"
)

// nameOf returns the name of the declaration that a param refers to, such as
// the field of version=Rev, or the text of a param that is a string. It
// returns the empty string for a param that is not set.
func nameOf(v any) string {
	switch p := v.(type) {
	case symbol.Identity:
		return p.Name
	case string:
		return p
	}
	return ""
}

// fieldRef returns the type of the field name of the struct that ref refers
// to, at the witnesses of the declaration.
func (k *caller) fieldRef(ref *node.TypeRef, name string) *emit.TypeRef {
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	sym, _ := k.m.Reader().Lookup(ref.Target)
	st, _ := sym.(*node.Struct)
	// The derivation binds a field param only to a field of the struct.
	at := slices.IndexFunc(st.Fields, func(f *node.Field) bool { return f.Name == name })
	return k.typeOf(st.Fields[at].Type)
}

// resultRef returns the type of the i-th result in the role value of the
// callable of f. Each value of a Go callable precedes its flag and its
// error, so the i-th value is its i-th result.
func (k *caller) resultRef(f *derive.Facts, i int) *node.TypeRef {
	return k.instance(f.Signature.Returns[i].Ref)
}

// zeroOf returns the zero value of the results of the callable of f: the
// zero of its one value, or the pair of the zeros of its two values.
func (k *caller) zeroOf(f *derive.Facts) *Expr {
	if len(f.Results) == 1 {
		return &Expr{Value: typed(k.typeOf(k.resultRef(f, 0)), &f.Results[0].Zero)}
	}
	p := &PairExpr{}
	for i := range 2 {
		p.Types[i] = k.typeOf(k.resultRef(f, i))
		p.Values[i] = Expr{Value: typed(p.Types[i], &f.Results[i].Zero)}
	}
	return &Expr{Pair: p}
}

// typed returns v as a value of the type ref: the conversion of a literal or
// a name, whose constant can be untyped, and any other value as it is, which
// already has its type.
func typed(ref *emit.TypeRef, v *emit.Value) *emit.Value {
	if v.Kind == emit.ValueLiteral || v.Kind == emit.ValueName {
		return new(emit.Conversion(ref, *v))
	}
	return v
}

// valueOfParam returns the expression of the value of a param, which names a
// declaration of the package, such as sentinel=ErrMissing.
func (*caller) valueOfParam(v any) *Expr {
	id, _ := v.(symbol.Identity)
	return &Expr{Value: new(emit.Name(id))}
}

// intOf returns the decimal text of an int param, which the classification
// gives as an int64 or as its text.
func intOf(v any) string {
	if n, ok := v.(int64); ok {
		return strconv.FormatInt(n, 10)
	}
	return nameOf(v)
}
