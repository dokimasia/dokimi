// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim

import (
	"errors"
	"slices"
	"strconv"
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/position"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/internal/lang"
	"go.dokimi.dev/dokimi/spec"
)

// ID is the name of the plugin.
const ID plugin.ID = "sim"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.1.0"

// Word is the word of the family of the simulations and of the family of
// their tests. TestTag is the tag of the family of the tests.
const (
	Word    = "sim"
	TestTag = sdk.TagTest
)

// The tree of the target contains the file template of the simulations and
// the file template of their tests.
const (
	simTemplate  = "sim.tmpl"
	testTemplate = "sim.test.tmpl"
)

// SimDirective puts a machine interface under sim. Its params StartParam,
// StopParam and CancelParam name the events of the lifecycle, TolerateParam
// the faults in force at once, RequireParam the share of the cases in which
// each fault strikes, StrategyParam the strategy of the scheduling, and
// IsolationParam the isolation level of the transactions of the clients.
// TokenDirective marks the type of a token, and with its positional param
// KindParam set to reply the type of a reply token. ReplyDirective marks the
// effect of a reply. Its params ValueParam, ErrParam and ToParam name its
// fields, and UnknownParam names the errors of a reply that leave the
// outcome of its call unknown. PerformDirective names in its param FuncParam
// the function that performs an effect, or in its param CallParam the method
// of a dependency that performs it and in its param ResultParam the event of
// the result. Its param TimeoutParam names the field of the effect that
// bounds its call. The other directives mark the
// functions of the consumer: SystemDirective the system, ActionDirective an
// action with the positional params CallParam and FinalParam, SpecDirective
// the spec of the calls, InvariantDirective an invariant with the param
// DefectParam, ConvergesDirective a convergence with the param WithinParam,
// FingerprintDirective the fingerprint and ScenarioDirective a scenario.
const (
	SimDirective         directive.Name     = "sim"
	TokenDirective       directive.Name     = "token"
	ReplyDirective       directive.Name     = "reply"
	PerformDirective     directive.Name     = "perform"
	SystemDirective      directive.Name     = "system"
	ActionDirective      directive.Name     = "action"
	SpecDirective        directive.Name     = "spec"
	InvariantDirective   directive.Name     = "invariant"
	ConvergesDirective   directive.Name     = "converges"
	FingerprintDirective directive.Name     = "fingerprint"
	ScenarioDirective    directive.Name     = "scenario"
	StartParam           directive.ParamKey = "start"
	StopParam            directive.ParamKey = "stop"
	CancelParam          directive.ParamKey = "cancel"
	TolerateParam        directive.ParamKey = "tolerate"
	RequireParam         directive.ParamKey = "require"
	StrategyParam        directive.ParamKey = "strategy"
	IsolationParam       directive.ParamKey = "isolation"
	KindParam            directive.ParamKey = "kind"
	ValueParam           directive.ParamKey = "value"
	ErrParam             directive.ParamKey = "err"
	ToParam              directive.ParamKey = "to"
	UnknownParam         directive.ParamKey = "unknown"
	FuncParam            directive.ParamKey = "func"
	ResultParam          directive.ParamKey = "result"
	TimeoutParam         directive.ParamKey = "timeout"
	CallParam            directive.ParamKey = "call"
	FinalParam           directive.ParamKey = "final"
	DefectParam          directive.ParamKey = "defect"
	WithinParam          directive.ParamKey = "within"
)

// KeyDefect and KeyWithin are the keys of the facts that the annotate half
// stamps on an invariant and on a convergence, and contain the identities of
// the functions of the params defect and within. KeyCall and KeyFinal report
// an action call and a final action. KeyToken contains token for a type under
// the directive token, and reply for a reply token. KeyReply contains the
// fields of the value, the error and the token of a reply, and KeyUnknown the
// identities of the errors of the reply that leave the outcome of its call
// unknown. KeyPerform contains the name of the function that performs an
// effect, and KeyTimeout the field of the effect that bounds its call.
// KeyCallee contains the identity of the method of a dependency that performs
// an effect, and KeyResult the identity of the event of its result.
const (
	KeyDefect  meta.KeyName = "sim.defect"
	KeyWithin  meta.KeyName = "sim.within"
	KeyCall    meta.KeyName = "sim.call"
	KeyFinal   meta.KeyName = "sim.final"
	KeyToken   meta.KeyName = "sim.token"
	KeyReply   meta.KeyName = "sim.reply"
	KeyUnknown meta.KeyName = "sim.unknown"
	KeyPerform meta.KeyName = "sim.perform"
	KeyTimeout meta.KeyName = "sim.timeout"
	KeyCallee  meta.KeyName = "sim.callee"
	KeyResult  meta.KeyName = "sim.result"
)

// The spellings of the positional params of the directives token and action,
// and of the strategies and the isolation levels of the directive sim.
const (
	replyKind         = "reply"
	tokenKind         = "token"
	callWord          = "call"
	finalWord         = "final"
	uniformStrategy   = "uniform"
	pctStrategy       = "pct"
	serializable      = "serializable"
	snapshotIsolation = "snapshot-isolation"
)

// The words of the names of the declarations of a simulation, and the name
// of the step method and the error method of a machine.
const (
	anatomyWord    = "anatomy"
	performWord    = "perform"
	simulationWord = "simulation"
	stepMethod     = "Step"
	errMethod      = "Err"
)

// faultErrors is the handle of the errors that the directive fault of a
// method declares, which the annotator fault stamps on the method.
var faultErrors = meta.Named[[]symbol.Identity](fault.KeyErrors)

// faults maps the spelling of each fault that the params tolerate and
// require name to the name of its constant in the module go.dokimi.dev/sim.
var faults = map[string]string{
	"crash": "FaultCrash", "restart": "FaultRestart", "pause": "FaultPause", "partition": "FaultPartition",
	"drop": "FaultDrop", "delay": "FaultDelay", "duplicate": "FaultDuplicate", "lose-before": "FaultLoseBefore",
	"lose-after": "FaultLoseAfter", "outage": "FaultOutage", "slow": "FaultSlow", "clock-jump": "FaultClockJump",
	"clock-strobe": "FaultClockStrobe",
}

// Model is the data of the record of the simulation of one machine
// interface, and the template of the family of the simulations reads it.
//
// Machine is the name of the interface, and Names the names of the
// declarations of its simulation. Event and Effect are its event and effect
// interfaces. Started, Stopped and Cancelled are the events of its lifecycle,
// and Cancelled is nil without the param cancel. Calls are the inbound calls
// with their reply tokens, Tokens the effects with their tokens, Replies the
// types of its replies, with an instantiation of a generic reply for each
// type of its union, and Performs the effects with their performers, each in
// the order of their identities. Unknown are the errors of the replies that
// leave the outcome of their call unknown, in the order of the replies and
// of their directives.
//
// System is the function of the system. Spec is the spec of the calls of the
// clients, and nil without one. State is the state of the spec. Fingerprint
// is the function of the fingerprint, and empty without one. Rules are the
// checks of the claims of the spec whose generator is sim. Actions,
// Invariants, Converges and Scenarios are the functions of the consumer, in
// the order of their identities.
//
// Tolerate and Require are the faults of the params tolerate and require, in
// their order. Strategy is uniform, pct or empty for the default, and Depth
// the depth of pct. Isolation is the isolation level, and empty for none.
// Timeouts reports whether an effect has a field that bounds its call.
type Model struct {
	Event       *emit.TypeRef
	Effect      *emit.TypeRef
	State       *emit.TypeRef
	Stopped     *emit.TypeRef
	Cancelled   *Lifecycle
	Spec        *Spec
	Started     Lifecycle
	Names       Names
	Machine     string
	System      string
	Fingerprint string
	Strategy    string
	Isolation   string
	Calls       []Field
	Tokens      []Field
	Replies     []Reply
	Unknown     []*emit.TypeRef
	Performs    []Perform
	Rules       []Rule
	Actions     []Action
	Invariants  []Invariant
	Converges   []Convergence
	Scenarios   []Scenario
	Tolerate    []Bound
	Require     []Share
	Depth       int
	Timeouts    bool
}

// Names are the names of the declarations of the simulation of a machine
// interface such as Machine: the variable of the anatomy machineAnatomy, the
// function that performs the effects performMachine, and the function that
// returns the simulation MachineSimulation.
type Names struct {
	Anatomy    string
	Perform    string
	Simulation string
}

// Lifecycle is an event that the world builds from a reply token: its type,
// the field of the reply token and the type of that field.
type Lifecycle struct {
	Type  *emit.TypeRef
	Token *emit.TypeRef
	Field string
}

// Field is an event or an effect with the field of its token: an inbound call
// with the field of its reply token, or an effect that starts a call or a
// wait with the field of its token.
type Field struct {
	Type  *emit.TypeRef
	Field string
}

// Reply is one type of a reply effect, with the fields of its value, its
// error and the reply token of its call.
type Reply struct {
	Type  *emit.TypeRef
	Value string
	Err   string
	To    string
}

// Perform is one effect with the function or the call that performs it, and
// the field of the effect that bounds its call, empty for none. Func is the
// name of the function, and empty for an effect that Call performs.
type Perform struct {
	Type    *emit.TypeRef
	Call    *Call
	Func    string
	Timeout string
}

// Call is the call of a method of a dependency that performs an effect.
// Dependency is the type of the dependency, and Name its name in the faults
// of a run, such as remote.Writer. Args are the arguments of the call, in the
// order of the parameters of Method. Event is the event of the result. It
// takes the token of the effect from its field EffectToken in its field
// Token, the error of the call in its field Err, and each value of the call
// in its field of Values. Fails reports a method with an error result.
// Errors are the errors that the directive fault of the method declares.
type Call struct {
	Dependency  *emit.TypeRef
	Event       *emit.TypeRef
	Name        string
	Method      string
	Token       string
	EffectToken string
	Err         string
	Values      []string
	Args        []Arg
	Errors      []*emit.TypeRef
	Fails       bool
}

// Arg is one argument of a call that performs an effect: the context of the
// world, or the field of the effect named Field.
type Arg struct {
	Field   string
	Context bool
}

// Rule is the check of one claim of the spec whose generator is sim, on a
// machine interface: its identity, such as Machine/replies-once, its
// statement, and the identifiers of its claim, its probe, its reason and the
// form of its defect.
type Rule struct {
	ID        string
	Statement string
	Claim     string
	Probe     string
	Reason    string
	Form      string
}

// Action is one action of the consumer: the name of its function, whether it
// is a call of a client, and whether the run takes it after the faults end.
type Action struct {
	Func  string
	Call  bool
	Final bool
}

// Invariant is one invariant of the consumer: its statement, the name of its
// function and the name of the function of its defect system.
type Invariant struct {
	Statement string
	Func      string
	Defect    string
}

// Convergence is one convergence of the consumer: its statement, the name of
// its function and the name of the function of its bound, empty for none.
type Convergence struct {
	Statement string
	Func      string
	Within    string
}

// Scenario is one scenario of the consumer: the name of its function, which is
// the name of the scenario.
type Scenario struct {
	Func string
}

// Spec is the spec of the calls of the clients: the name of its function,
// and whether the function returns the subject that history.SpecFrom wraps.
type Spec struct {
	Func string
	From bool
}

// Bound is the most faults of one kind in force at once: the constant of the
// kind in the module go.dokimi.dev/sim, such as FaultCrash, and the count.
type Bound struct {
	Fault string
	N     int
}

// Share is the share of the cases in which one kind of fault strikes: the
// constant of the kind and the share as written, such as 0.2.
type Share struct {
	Fault string
	Share string
}

// Checks is the data of the record of the tests of the simulation of one
// machine interface, and the template of the family of the tests reads it:
// the name of the interface and the name of the function of its simulation.
type Checks struct {
	Machine    string
	Simulation string
}

// part is the part of sim for one language.
type part interface {
	lang.Part
}

// parts are the parts of sim. RFC-0025 designs the Go runtime alone.
var parts = []part{golangPart{}}

// simSchema is the directive of a machine interface.
var simSchema = directive.Schema{
	Name: SimDirective,
	Params: []directive.ParamSpec{
		{
			Key: StartParam, Type: directive.TypeReference, Resolution: directive.ResolveTypeInScope, Required: true,
			Doc: "the event that starts a process, an inbound call with one field of the reply token",
		},
		{
			Key: StopParam, Type: directive.TypeReference, Resolution: directive.ResolveTypeInScope, Required: true,
			Doc: "the event that stops a process",
		},
		{
			Key: CancelParam, Type: directive.TypeReference, Resolution: directive.ResolveTypeInScope,
			Doc: "the event that ends the wait of the caller of an inbound call, with one field of the reply token",
		},
		{
			Key: TolerateParam, Type: directive.TypeList, ListOf: directive.TypeString,
			Doc: "the most faults of each kind in force at once, such as [crash:1, partition:1]",
		},
		{
			Key: RequireParam, Type: directive.TypeList, ListOf: directive.TypeString,
			Doc: "the share of the cases in which each kind of fault strikes, such as [crash:0.2]",
		},
		{
			Key: StrategyParam, Type: directive.TypeString,
			Doc: "the strategy that picks the next process: uniform, or pct:<depth>, pct:3 by default",
		},
		{
			Key: IsolationParam, Type: directive.TypeString, Choices: []string{serializable, snapshotIsolation},
			Doc: "the isolation level that the list-append transactions of the clients must have",
		},
	},
	Doc: "generates the deterministic simulation of a system of the sans-IO machines of the interface, and its tests",
}

// tokenSchema is the directive of a token type.
var tokenSchema = directive.Schema{
	Name: TokenDirective,
	Positional: []directive.ParamSpec{{
		Key: KindParam, Type: directive.TypeString, Choices: []string{replyKind},
		Doc: "reply for the type of the field that links an inbound call to the effect that replies to it",
	}},
	Doc: "marks the type of the field that links an effect to the event of its result",
}

// replySchema is the directive of the effect of a reply.
var replySchema = directive.Schema{
	Name: ReplyDirective,
	Params: []directive.ParamSpec{
		{Key: ValueParam, Type: directive.TypeString, Default: "Value", Doc: "the field of the value of the reply"},
		{Key: ErrParam, Type: directive.TypeString, Default: "Err", Doc: "the field of the error of the reply"},
		{Key: ToParam, Type: directive.TypeString, Default: "To", Doc: "the field of the reply token of the call"},
		{
			Key: UnknownParam, Type: directive.TypeList, ListOf: directive.TypeReference,
			Resolution: directive.ResolvePackageVar,
			Doc:        "the errors of a reply that leave the outcome of its call unknown, such as [ErrStopped]",
		},
	},
	Doc: "marks the effect that replies to an inbound call",
}

// performSchema is the directive of the performer of an effect.
var performSchema = directive.Schema{
	Name: PerformDirective,
	Params: []directive.ParamSpec{
		{
			Key: FuncParam, Type: directive.TypeString,
			Doc: "the function of the package of the simulation that performs the effect",
		},
		{
			Key: CallParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
			Doc: "the method of the dependency that performs the effect, such as store.Writer.Commit",
		},
		{
			Key: ResultParam, Type: directive.TypeReference, Resolution: directive.ResolveTypeInScope,
			Doc: "the event of the result of the call",
		},
		{Key: TimeoutParam, Type: directive.TypeString, Doc: "the field of the effect that bounds its call"},
	},
	Doc: "names the function or the call of a dependency that performs the effect",
}

// functional is a directive of the functions of the consumer: its schema,
// and the method that binds a function under it to the model of a
// simulation.
type functional struct {
	bind   func(in *invocation, model *Model, fn *node.Function)
	schema directive.Schema
}

// functionals are the directives of the functions of the consumer, in the
// order in which a simulation binds their functions.
var functionals = []functional{
	{
		bind: (*invocation).bindSystem,
		schema: directive.Schema{
			Name: SystemDirective, Doc: "marks the function that builds the processes of the system",
		},
	},
	{
		bind: (*invocation).bindAction,
		schema: directive.Schema{
			Name: ActionDirective,
			Positional: []directive.ParamSpec{
				{
					Key: CallParam, Type: directive.TypeString, Choices: []string{callWord},
					Doc: "call for a call of a client",
				},
				{
					Key: FinalParam, Type: directive.TypeString, Choices: []string{finalWord},
					Doc: "final for a call that the run makes once after the faults end",
				},
			},
			Doc: "marks the function that takes the inputs of an action from the case and returns its run",
		},
	},
	{
		bind: (*invocation).bindSpec,
		schema: directive.Schema{
			Name: SpecDirective, Doc: "marks the function that returns the sequential spec of the calls of the clients",
		},
	},
	{
		bind: (*invocation).bindInvariant,
		schema: directive.Schema{
			Name: InvariantDirective,
			Params: []directive.ParamSpec{{
				Key: DefectParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
				Required: true, Doc: "the function that builds the system with a planted defect of the invariant",
			}},
			Doc: "marks a safety property that the run checks after every step",
		},
	},
	{
		bind: (*invocation).bindConvergence,
		schema: directive.Schema{
			Name: ConvergesDirective,
			Params: []directive.ParamSpec{{
				Key: WithinParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
				Doc: "the function that returns the bound of virtual time within which the property becomes true",
			}},
			Doc: "marks a liveness property that the run checks once the faults end",
		},
	},
	{
		bind: (*invocation).bindFingerprint,
		schema: directive.Schema{
			Name: FingerprintDirective, Doc: "marks the function that returns a hash of the state of the system",
		},
	},
	{
		bind:   (*invocation).bindScenario,
		schema: directive.Schema{Name: ScenarioDirective, Doc: "marks a fixed run whose events the crash points cover"},
	},
}

// generator is one instance of sim: the handles of the keys of its facts.
type generator struct {
	defect  meta.Key[symbol.Identity]
	within  meta.Key[symbol.Identity]
	call    meta.Key[bool]
	final   meta.Key[bool]
	token   meta.Key[string]
	reply   meta.Key[[]string]
	unknown meta.Key[[]symbol.Identity]
	perform meta.Key[string]
	timeout meta.Key[string]
	callee  meta.Key[symbol.Identity]
	result  meta.Key[symbol.Identity]
}

// New returns sim, a plugin with an annotate half and a generate half.
//
// The annotate half stamps the facts of the directives token, reply and
// perform on their types, and the facts of the directives of the consumer on
// their functions. The generate half records a [Model] into the family [Word]
// for each interface under [SimDirective], and a [Checks] into the family
// [Word] under the tag [TestTag]. The template of the family writes the
// anatomy of the machine, the function that performs its effects and the
// function that returns its simulation. The template of the tests writes
// TestSimulation, TestSimulationFalsify and FuzzSimulation. sim serves Go,
// the language of its part, and leaves a declaration of another language
// alone. It reports the findings of the package documentation and records
// nothing for a machine interface that it reports. Each call returns a new
// plugin.
func New() sdk.Dual {
	g := &generator{}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Keys(g.register).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: simTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, External: true, Word: Word, Template: testTemplate})
	lang.Register(b, parts)
	handlers := make([]sdk.Rule, 0, 4+len(functionals))
	handlers = append(handlers,
		sdk.Directive(simSchema, sdk.OnInterface(lang.Served(parts, g.simulationOf))),
		sdk.Directive(tokenSchema, sdk.OnAlias(lang.Served(parts, g.concludeToken))),
		sdk.Directive(replySchema, sdk.OnStruct(lang.Served(parts, g.concludeReply))),
		sdk.Directive(performSchema, sdk.OnStruct(lang.Served(parts, g.concludePerform))),
	)
	for i := range functionals {
		handlers = append(handlers,
			sdk.Directive(functionals[i].schema, sdk.OnFunction(lang.Served(parts, g.concludeFunction))))
	}
	return b.Handle(handlers...).Dual()
}

// register claims the namespace of sim and registers the keys of its facts.
//
// # Errors
//
// register returns the error of the claim, such as a namespace that another
// registrant claimed, joined with the errors of the registrations, such as a
// key that the registrant registered with another type.
func (g *generator) register(r *meta.Registry) error {
	on := func(kind symbol.Kind, name meta.KeyName, doc string) meta.KeySpec {
		return meta.KeySpec{Name: name, Kinds: []symbol.Kind{kind}, Doc: doc}
	}
	var errs [12]error
	errs[0] = r.ClaimNamespace(KeyDefect.Namespace())
	g.defect, errs[1] = meta.Register[symbol.Identity](r, on(symbol.KindFunction, KeyDefect,
		"contains the identity of the function that builds the system with the planted defect of the invariant"))
	g.within, errs[2] = meta.Register[symbol.Identity](r, on(symbol.KindFunction, KeyWithin,
		"contains the identity of the function that returns the bound of the convergence"))
	g.call, errs[3] = meta.Register[bool](r, on(symbol.KindFunction, KeyCall,
		"reports that the action is a call of a client"))
	g.final, errs[4] = meta.Register[bool](r, on(symbol.KindFunction, KeyFinal,
		"reports that the run takes the action once after the faults end"))
	g.token, errs[5] = meta.Register[string](r, on(symbol.KindAlias, KeyToken,
		"contains token for the type of a token, and reply for the type of a reply token"))
	g.reply, errs[6] = meta.Register[[]string](r, on(symbol.KindStruct, KeyReply,
		"contains the fields of the value, the error and the reply token of the reply"))
	g.perform, errs[7] = meta.Register[string](r, on(symbol.KindStruct, KeyPerform,
		"contains the name of the function that performs the effect"))
	g.timeout, errs[8] = meta.Register[string](r, on(symbol.KindStruct, KeyTimeout,
		"contains the field of the effect that bounds its call"))
	g.unknown, errs[9] = meta.Register[[]symbol.Identity](r, on(symbol.KindStruct, KeyUnknown,
		"contains the identities of the errors of the reply that leave the outcome of its call unknown"))
	g.callee, errs[10] = meta.Register[symbol.Identity](r, on(symbol.KindStruct, KeyCallee,
		"contains the identity of the method of the dependency that performs the effect"))
	g.result, errs[11] = meta.Register[symbol.Identity](r, on(symbol.KindStruct, KeyResult,
		"contains the identity of the event of the result of the call that performs the effect"))
	return errors.Join(errs[:]...)
}

// concludeToken stamps the kind of a token type. It reports [BadToken] at a
// type that is not an unsigned integer type, and stamps nothing for it.
func (g *generator) concludeToken(m *sdk.AliasMatch, _ part, st *sdk.Stamper) error {
	a := m.Alias
	b := m.Rules()
	if !unsigned(&b, a.Target) {
		m.Errorf(BadToken, "%s is not an unsigned integer type, so it is no token", a.Name)
		return nil
	}
	kind := tokenKind
	if args := m.Directive().Args; len(args) > 0 {
		kind = args[0].Str
	}
	sdk.Stamp(st, g.token, kind)
	return nil
}

// concludeReply stamps the fields of a reply and the errors of an unknown
// outcome.
func (g *generator) concludeReply(m *sdk.StructMatch, _ part, st *sdk.Stamper) error {
	d := m.Directive()
	sdk.Stamp(st, g.reply, []string{d.Params[ValueParam].Str, d.Params[ErrParam].Str, d.Params[ToParam].Str})
	unknown := d.Params[UnknownParam].List
	errs := make([]symbol.Identity, 0, len(unknown))
	for i := range unknown {
		errs = append(errs, unknown[i].Target)
	}
	sdk.Stamp(st, g.unknown, errs)
	return nil
}

// concludePerform stamps the function or the call that performs an effect,
// the event of the result of the call, and the field that bounds the call. It
// reports [BadPerform] at a directive that names both a function and a call
// or neither, and at a call without a result, and stamps nothing for it.
func (g *generator) concludePerform(m *sdk.StructMatch, _ part, st *sdk.Stamper) error {
	d := m.Directive()
	fn, function := d.Param(FuncParam)
	call, called := d.Param(CallParam)
	result, resulted := d.Param(ResultParam)
	name := m.Struct.Name
	switch {
	case function && called:
		m.Errorf(BadPerform, "the directive perform of %s names both a function and a call. Name one of them", name)
		return nil
	case !function && !called:
		m.Errorf(BadPerform, "the directive perform of %s names neither a function nor a call. Write "+
			"func=<function>, or call=<callable> result=<event>", name)
		return nil
	case called && !resulted:
		m.Errorf(BadPerform, "the directive perform of %s names a call without a result. Write result=<event>", name)
		return nil
	case called:
		sdk.Stamp(st, g.callee, call.Target)
		sdk.Stamp(st, g.result, result.Target)
	default:
		sdk.Stamp(st, g.perform, fn.Str)
	}
	if timeout, stated := d.Param(TimeoutParam); stated {
		sdk.Stamp(st, g.timeout, timeout.Str)
	}
	return nil
}

// concludeFunction stamps the facts of the params of the directive of a
// function of the consumer.
func (g *generator) concludeFunction(m *sdk.FunctionMatch, _ part, st *sdk.Stamper) error {
	d := m.Directive()
	if defect, stated := d.Param(DefectParam); stated {
		sdk.Stamp(st, g.defect, defect.Target)
	}
	if within, stated := d.Param(WithinParam); stated {
		sdk.Stamp(st, g.within, within.Target)
	}
	for i := range d.Args {
		switch d.Args[i].Str {
		case callWord:
			sdk.Stamp(st, g.call, true)
		case finalWord:
			sdk.Stamp(st, g.final, true)
		}
	}
	return nil
}

// unsigned reports whether the type ref is an unsigned integer type.
func unsigned(b *rules.Bound, ref *node.TypeRef) bool {
	switch rules.Unaliased(ref, b.View()).Spelling {
	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return true
	}
	return false
}

// invocation is the state of the simulation of one machine interface: the
// match of the interface, the generator, its event and effect interfaces,
// the rules of its language, the package of the simulation, and whether a
// finding refused the interface.
type invocation struct {
	m       *sdk.InterfaceMatch
	g       *generator
	event   *node.TypeRef
	effect  *node.TypeRef
	b       rules.Bound
	pkg     symbol.Identity
	refused bool
}

// refuse reports the finding of code at pos, and refuses the interface, so
// sim records nothing for it.
func (in *invocation) refuse(code diag.Code, pos position.Pos, format string, args ...any) {
	in.m.ErrorfAt(code, pos, format, args...)
	in.refused = true
}

// simulationOf records the simulation of a machine interface under the
// directive, and the record of its tests.
func (g *generator) simulationOf(m *sdk.InterfaceMatch, _ part, e *sdk.Emitter) error {
	pkg, _ := e.File().Package()
	in := &invocation{m: m, b: m.Rules(), g: g, pkg: pkg}
	model, ok := in.model()
	if !ok {
		return nil
	}
	i := m.Interface
	record := &emit.Function{Origin: i.ID, Name: model.Names.Simulation, Visibility: symbol.VisibilityPublic}
	e.File().Record(model, record)
	test := &emit.Struct{Origin: i.ID, Name: model.Names.Simulation, Visibility: symbol.VisibilityPublic}
	test.Fields.Append(&emit.Field{
		Name: simulationWord,
		Type: &emit.TypeRef{Spelling: model.Names.Simulation, Target: i.ID, Emitted: true},
	})
	e.File(TestTag).Record(Checks{Machine: i.Name, Simulation: model.Names.Simulation}, test)
	return nil
}

// model returns the model of the simulation of the interface, and false
// after a finding. It checks the params, the lifecycle, the anatomy and the
// functions of a machine in turn, so one run reports the findings of each.
func (in *invocation) model() (Model, bool) {
	i := in.m.Interface
	if in.machine(); in.refused {
		return Model{}, false
	}
	model := Model{
		Machine: i.Name,
		Event:   rules.EmitRef(in.event),
		Effect:  rules.EmitRef(in.effect),
		Names: Names{
			Anatomy:    naming.Camel(i.Name + "Anatomy"),
			Perform:    naming.Camel(performWord + i.Name),
			Simulation: naming.Pascal(i.Name + "Simulation"),
		},
	}
	in.params(&model)
	events := in.members(in.event)
	in.lifecycle(&model, events)
	in.anatomy(&model, events, in.members(in.effect))
	in.functions(&model)
	if len(model.Replies) == 0 && slices.ContainsFunc(model.Actions, func(a Action) bool { return a.Call }) {
		in.refuse(NoReply, i.Pos, "%s has an action call and no effect under the directive reply, so the run "+
			"cannot complete a call in the history", i.Name)
	}
	if in.refused {
		return Model{}, false
	}
	model.Rules = rulesOf(i.Name)
	return model, true
}

// machine checks that the interface has the step method and the error method
// of a sans-IO machine, and sets the event and the effect interfaces. It
// reports [BadMachine] for each method that the interface lacks.
func (in *invocation) machine() {
	i := in.m.Interface
	set, _ := in.b.MembersOf(i)
	var step, failure *node.Method
	for i := range set.Members {
		if method, isMethod := set.Members[i].Symbol.(*node.Method); isMethod {
			switch method.Name {
			case stepMethod:
				step = method
			case errMethod:
				failure = method
			}
		}
	}
	if step == nil || !in.stepForm(step) {
		in.refuse(BadMachine, i.Pos, "%s has no method Step(ev E, out []F) []F with interface types E and F, "+
			"so it is no sans-IO machine", i.Name)
	}
	if failure == nil || !fits(paramTypes(failure.Params)) || !fits(returnTypes(failure.Returns), in.isError) {
		in.refuse(BadMachine, i.Pos, "%s has no method Err() error, so it is no sans-IO machine", i.Name)
	}
}

// isError reports whether the type ref is the type error.
func (in *invocation) isError(ref *node.TypeRef) bool {
	return in.b.TypeOf(ref).Ref == rules.WellKnownError
}

// stepForm reports whether the method has the form Step(ev E, out []F) []F
// with interface types E and F, and sets the event and the effect interfaces.
// The element of a list is its one element.
func (in *invocation) stepForm(step *node.Method) bool {
	if len(step.Params) != 2 || len(step.Returns) != 1 {
		return false
	}
	ev, out, ret := step.Params[0].Type, step.Params[1].Type, step.Returns[0].Type
	if out.Form != symbol.FormList || ret.Form != symbol.FormList || out.Elems[0].Target != ret.Elems[0].Target ||
		!in.isInterface(ev) || !in.isInterface(out.Elems[0]) {

		return false
	}
	in.event, in.effect = ev, out.Elems[0]
	return true
}

// isInterface reports whether the type ref names an interface of the view.
func (in *invocation) isInterface(ref *node.TypeRef) bool {
	decl, _ := in.b.View().Lookup(ref.Target)
	_, isInterface := decl.(*node.Interface)
	return isInterface
}

// members returns the structs of the view that implement the interface ref,
// in identity order.
func (in *invocation) members(ref *node.TypeRef) []*node.Struct {
	var out []*node.Struct
	for _, id := range in.b.Implementers(ref) {
		decl, _ := in.b.View().Lookup(id)
		if s, isStruct := decl.(*node.Struct); isStruct {
			out = append(out, s)
		}
	}
	return out
}

// params reads the params tolerate, require, strategy and isolation of the
// directive into the model. It reports [BadParam] at each value that names
// no fault or no strategy. A share that does not parse is 0, infinite or NaN,
// so its range refuses it.
func (in *invocation) params(model *Model) {
	d := in.m.Directive()
	i := in.m.Interface
	if tolerate, stated := d.Param(TolerateParam); stated {
		for at := range tolerate.List {
			entry := tolerate.List[at].Str
			kind, count, _ := strings.Cut(entry, ":")
			if n, err := strconv.Atoi(count); faults[kind] != "" && err == nil && n > 0 {
				model.Tolerate = append(model.Tolerate, Bound{Fault: faults[kind], N: n})
			} else {
				in.refuse(BadParam, i.Pos, "the param tolerate of %s has the entry %q, which is no fault with a "+
					"count of at least 1, such as crash:1", i.Name, entry)
			}
		}
	}
	if require, stated := d.Param(RequireParam); stated {
		for at := range require.List {
			entry := require.List[at].Str
			kind, share, _ := strings.Cut(entry, ":")
			if f, _ := strconv.ParseFloat(share, 64); faults[kind] != "" && f > 0 && f <= 1 {
				model.Require = append(model.Require, Share{Fault: faults[kind], Share: share})
			} else {
				in.refuse(BadParam, i.Pos, "the param require of %s has the entry %q, which is no fault with a "+
					"share above 0 and at most 1, such as crash:0.2", i.Name, entry)
			}
		}
	}
	if strategy, stated := d.Param(StrategyParam); stated {
		name, depth, _ := strings.Cut(strategy.Str, ":")
		n, err := strconv.Atoi(depth)
		switch {
		case strategy.Str == uniformStrategy:
			model.Strategy = uniformStrategy
		case name == pctStrategy && err == nil && n > 0:
			model.Strategy, model.Depth = pctStrategy, n
		default:
			in.refuse(BadParam, i.Pos, "the param strategy of %s is %q, which is neither uniform nor pct with a "+
				"depth of at least 1, such as pct:3", i.Name, strategy.Str)
		}
	}
	if isolation, stated := d.Param(IsolationParam); stated {
		model.Isolation = isolation.Str
	}
}

// lifecycle sets the events of the lifecycle from the params start, stop
// and cancel. It reports [BadLifecycle] at each param that names no event of
// its form.
func (in *invocation) lifecycle(model *Model, events []*node.Struct) {
	d := in.m.Directive()
	i := in.m.Interface
	find := func(key directive.ParamKey) *node.Struct {
		at := slices.IndexFunc(events, func(s *node.Struct) bool { return s.ID == d.Params[key].Target })
		if at < 0 {
			return nil
		}
		return events[at]
	}
	if started := find(StartParam); in.replied(started) {
		model.Started = in.lifecycleOf(started)
	} else {
		in.refuse(BadLifecycle, i.Pos, "the param start of %s names no event of %s with one field of a reply "+
			"token", i.Name, in.event.Spelling)
	}
	if stopped := find(StopParam); stopped != nil {
		model.Stopped = refOf(stopped)
	} else {
		in.refuse(BadLifecycle, i.Pos, "the param stop of %s names no event of %s", i.Name, in.event.Spelling)
	}
	if _, stated := d.Param(CancelParam); !stated {
		return
	}
	if cancelled := find(CancelParam); in.replied(cancelled) {
		cancel := in.lifecycleOf(cancelled)
		model.Cancelled = &cancel
	} else {
		in.refuse(BadLifecycle, i.Pos, "the param cancel of %s names no event of %s with one field of a reply "+
			"token", i.Name, in.event.Spelling)
	}
}

// replied reports whether the struct s is an event with one field of a reply
// token, and false where s is nil.
func (in *invocation) replied(s *node.Struct) bool {
	return s != nil && len(in.tokenFields(s, replyKind)) == 1
}

// lifecycleOf returns the lifecycle event s with the type and the name of its
// field of a reply token.
func (in *invocation) lifecycleOf(s *node.Struct) Lifecycle {
	f := in.tokenFields(s, replyKind)[0]
	return Lifecycle{Type: refOf(s), Token: rules.EmitRef(f.Type), Field: f.Name}
}

// tokenFields returns the fields of the struct s whose type is a token of
// the kind, token or reply.
func (in *invocation) tokenFields(s *node.Struct, kind string) []*node.Field {
	var out []*node.Field
	for _, f := range s.Fields {
		if k, _ := sdk.FactOf(in.m, f.Type.Target, in.g.token); k == kind {
			out = append(out, f)
		}
	}
	return out
}

// anatomy sets the inbound calls, the tokens, the replies and the performers
// of the model from the events and the effects of the machine. It reports
// [NoPerform], [BadReply] and [BadFunction] at the effects at fault.
func (in *invocation) anatomy(model *Model, events, effects []*node.Struct) {
	for _, ev := range events {
		if model.Cancelled != nil && ev.ID == model.Cancelled.Type.Target {
			continue
		}
		if fields := in.tokenFields(ev, replyKind); len(fields) == 1 {
			model.Calls = append(model.Calls, Field{Type: refOf(ev), Field: fields[0].Name})
		}
	}
	for _, f := range effects {
		if names, marked := sdk.FactOf(in.m, f.ID, in.g.reply); marked {
			in.reply(model, f, names)
			continue
		}
		if fields := in.tokenFields(f, tokenKind); len(fields) == 1 {
			model.Tokens = append(model.Tokens, Field{Type: refOf(f), Field: fields[0].Name})
		}
		fn, performed := sdk.FactOf(in.m, f.ID, in.g.perform)
		callee, called := sdk.FactOf(in.m, f.ID, in.g.callee)
		timeout, _ := sdk.FactOf(in.m, f.ID, in.g.timeout)
		switch {
		case called:
			if call, bound := in.called(f, &callee); bound {
				model.Performs = append(model.Performs, Perform{Type: refOf(f), Call: call, Timeout: timeout})
			}
		case !performed:
			in.refuse(NoPerform, f.Pos, "the effect %s of %s has no directive perform. Write //dokimi:perform "+
				"func=<function> on %s", f.Name, in.m.Interface.Name, f.Name)
		case !in.performer(fn, f):
			in.refuse(BadFunction, f.Pos, "the package of the simulation of %s has no function %s with the "+
				"signature %s, with T the effect %s", in.m.Interface.Name, fn, performForm, f.Name)
		default:
			model.Performs = append(model.Performs, Perform{Type: refOf(f), Func: fn, Timeout: timeout})
		}
	}
	model.Timeouts = slices.ContainsFunc(model.Performs, func(p Perform) bool { return p.Timeout != "" })
}

// called returns the call of the method callee of an interface that performs
// the effect f. It binds each parameter of the method to the field of f with
// its name in PascalCase, and the token of f, the error and each value of the
// method to the fields of the event of the result by their types. It reports
// [BadPerform] at f where callee is no method of an interface, where f has no
// field of a token, and where a parameter or a result has no field to bind.
func (in *invocation) called(f *node.Struct, callee *symbol.Identity) (*Call, bool) {
	view := in.b.View()
	decl, _ := view.Lookup(*callee)
	method, isMethod := decl.(*node.Method)
	var host *node.Interface
	if isMethod {
		owner, _ := view.Lookup(method.Host)
		host, _ = owner.(*node.Interface)
	}
	tokens := in.tokenFields(f, tokenKind)
	switch {
	case host == nil:
		in.refuse(BadPerform, f.Pos, "the call %s of %s is no method of an interface", callee.Name, f.Name)
		return nil, false
	case len(tokens) != 1:
		in.refuse(BadPerform, f.Pos, "the effect %s has no field of a token, which the event of its result takes",
			f.Name)
		return nil, false
	}
	result, _ := sdk.FactOf(in.m, f.ID, in.g.result)
	var fields []*node.Field
	d, _ := view.Lookup(result)
	if ev, isStruct := d.(*node.Struct); isStruct {
		fields = ev.Fields
	}
	field := func(fits func(t *node.TypeRef) bool) string {
		if at := slices.IndexFunc(fields, func(x *node.Field) bool { return fits(x.Type) }); at >= 0 {
			return fields[at].Name
		}
		return ""
	}
	pkg, _ := in.m.Reader().PackageOf(host.ID.PackageIdentity())
	label := host.Name + "." + method.Name
	call := &Call{
		Dependency: &emit.TypeRef{Spelling: host.Name, Target: host.ID},
		Event:      &emit.TypeRef{Spelling: result.Name, Target: result},
		Name:       pkg.Name + "." + host.Name, Method: method.Name, EffectToken: tokens[0].Name,
		Token: field(func(t *node.TypeRef) bool { return t.Target == tokens[0].Type.Target }),
		Err:   field(in.isError),
	}
	unbound := func(what, role string) (*Call, bool) {
		in.refuse(BadPerform, f.Pos, "the event %s of the result of %s has no field of the type %s for the %s",
			result.Name, f.Name, what, role)
		return nil, false
	}
	switch {
	case call.Token == "":
		return unbound(tokens[0].Type.Spelling, "token of "+f.Name)
	case call.Err == "":
		return unbound("error", "error of "+label)
	}
	sig, _ := in.b.CallableOf(method)
	for _, p := range sig.Params {
		if p.Role == rules.ParamContext {
			call.Args = append(call.Args, Arg{Context: true})
			continue
		}
		name := naming.Pascal(p.Name)
		if !slices.ContainsFunc(f.Fields, func(x *node.Field) bool { return x.Name == name }) {
			in.refuse(BadPerform, f.Pos, "the effect %s has no field %s for the parameter %s of %s",
				f.Name, name, p.Name, label)
			return nil, false
		}
		call.Args = append(call.Args, Arg{Field: name})
	}
	for _, r := range sig.Returns {
		if r.Role == rules.ReturnError {
			call.Fails = true
			continue
		}
		value := field(func(t *node.TypeRef) bool { return sameType(t, r.Ref) })
		if value == "" {
			return unbound(r.Ref.Spelling, "value of "+label)
		}
		call.Values = append(call.Values, value)
	}
	errs, _ := sdk.FactOf(in.m, *callee, faultErrors)
	for _, id := range errs {
		call.Errors = append(call.Errors, &emit.TypeRef{Spelling: id.Name, Target: id})
	}
	return call, true
}

// sameType reports whether the type refs a and b name one type: one
// declaration, or one builtin type by its spelling.
func sameType(a, b *node.TypeRef) bool {
	return a.Target == b.Target && (a.Target != symbol.Identity{} || a.Spelling == b.Spelling)
}

// reply adds the instantiations of the reply effect f, whose directive names
// the fields of its value, its error and its reply token, to the model. It
// reports [BadReply] at f where it lacks one of the fields, and where it is
// generic over a constraint whose instantiations no list contains.
func (in *invocation) reply(model *Model, f *node.Struct, names []string) {
	missing := slices.IndexFunc(names, func(name string) bool {
		return !slices.ContainsFunc(f.Fields, func(field *node.Field) bool { return field.Name == name })
	})
	instances, listed := in.instances(f)
	switch {
	case missing >= 0:
		in.refuse(BadReply, f.Pos, "the reply %s has no field %s that its directive names", f.Name, names[missing])
	case !listed:
		in.refuse(BadReply, f.Pos, "the reply %s is generic over another constraint than a union of types without "+
			"an approximation, so the anatomy cannot list its instantiations", f.Name)
	default:
		for _, t := range instances {
			model.Replies = append(model.Replies, Reply{Type: t, Value: names[0], Err: names[1], To: names[2]})
		}
		unknown, _ := sdk.FactOf(in.m, f.ID, in.g.unknown)
		for _, id := range unknown {
			model.Unknown = append(model.Unknown, &emit.TypeRef{Spelling: id.Name, Target: id})
		}
	}
}

// instances returns the instantiations of the reply f: f itself without a
// type parameter, and f at each type of the union that constrains its one
// type parameter otherwise, in the order of the union. The union is the
// constraint of the type parameter, or the one constraint element of the
// interface that the type parameter names. instances reports false for any
// other generic reply. A union with an approximation, such as ~int, has
// instantiations that no list contains, and the frontend keeps the mark of
// the approximation in the spelling of the union alone.
func (in *invocation) instances(f *node.Struct) ([]*emit.TypeRef, bool) {
	if len(f.TypeParams) == 0 {
		return []*emit.TypeRef{refOf(f)}, true
	}
	if len(f.TypeParams) > 1 {
		return nil, false
	}
	union := f.TypeParams[0].Bounds[0]
	decl, _ := in.b.View().Lookup(union.Target)
	if constraint, isInterface := decl.(*node.Interface); isInterface && len(constraint.Constraints) == 1 {
		union = constraint.Constraints[0]
	}
	if union.Form != symbol.FormUnion || strings.Contains(union.Spelling, "~") {
		return nil, false
	}
	out := make([]*emit.TypeRef, 0, len(union.Elems))
	for _, term := range union.Elems {
		out = append(out, &emit.TypeRef{Spelling: f.Name, Target: f.ID, Args: []*emit.TypeRef{rules.EmitRef(term)}})
	}
	return out, true
}

// functions binds the functions of the consumer in the package of the
// simulation to the model, in the order of the directives and of their
// identities. A function under a directive in another package belongs to the
// simulation of another machine interface. functions reports [BadFunction] at
// a function that does not have the signature of its directive, at a second
// function of a directive that a simulation has once, and at the interface
// where the package has no function under the directive system.
func (in *invocation) functions(model *Model) {
	reader := in.b.View().Reader()
	for i := range functionals {
		d := &functionals[i]
		// The reader keys a directive by its spelling as written, bare or
		// canonical, and returns the functions of a spelling in identity order.
		var carriers []*node.Function
		canonical := directive.Schema{Plugin: string(ID), Name: d.schema.Name}.Canonical()
		for _, name := range []directive.Name{d.schema.Name, canonical} {
			for sym := range reader.ByDirective(name) {
				// The validation of the directives admits a function alone.
				if fn, _ := sym.(*node.Function); fn.ID.Package == in.pkg.Package {
					carriers = append(carriers, fn)
				}
			}
		}
		if d.schema.Name == SystemDirective && len(carriers) == 0 {
			in.refuse(BadFunction, in.m.Interface.Pos, "the package of the simulation of %s has no function under "+
				"the directive system, so the simulation has no processes", in.m.Interface.Name)
		}
		slices.SortFunc(carriers, func(a, b *node.Function) int { return a.ID.Compare(b.ID) })
		for _, fn := range carriers {
			d.bind(in, model, fn)
		}
	}
}

// The import paths of the packages whose types the signatures of the
// functions of the consumer name.
const (
	simPackage     = "go.dokimi.dev/sim"
	propPackage    = "go.dokimi.dev/assert/prop"
	historyPackage = "go.dokimi.dev/assert/history"
	timePackage    = "time"
)

// The forms of the signatures of the functions of the consumer, in the words
// of a finding. E and F are the event and the effect interfaces of the
// machine, S is the state of a spec, and T is the type of an effect.
const (
	systemForm      = "func(w *sim.World[E, F])"
	checkForm       = "func(h *sim.Check, w *sim.World[E, F])"
	actionForm      = "func(c *prop.Case) func(w *sim.World[E, F])"
	callForm        = "func(c *prop.Case) func(w *sim.World[E, F], client int)"
	specForm        = "func() history.Spec[S] or func() history.Subject"
	withinForm      = "func() time.Duration"
	fingerprintForm = "func(w *sim.World[E, F]) uint64"
	performForm     = "func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E]"
)

// refuseForm reports [BadFunction] at the function fn under the directive
// name, which needs the signature form.
func (in *invocation) refuseForm(fn *node.Function, name directive.Name, form string) {
	in.refuse(BadFunction, fn.Pos, "%s is under the directive %s, so it needs the signature %s, with E and F the "+
		"event and the effect interfaces of %s", fn.Name, name, form, in.m.Interface.Name)
}

// twice reports [BadFunction] at the function fn, the second function under
// the directive name after the function first.
func (in *invocation) twice(fn *node.Function, name directive.Name, first string) {
	in.refuse(BadFunction, fn.Pos, "%s is a second function under the directive %s of %s, after %s", fn.Name,
		name, in.m.Interface.Name, first)
}

// bindSystem sets the system of the model to the function fn.
func (in *invocation) bindSystem(model *Model, fn *node.Function) {
	switch {
	case !in.system(fn):
		in.refuseForm(fn, SystemDirective, systemForm)
	case model.System != "":
		in.twice(fn, SystemDirective, model.System)
	default:
		model.System = fn.Name
	}
}

// bindScenario adds the scenario fn to the model.
func (in *invocation) bindScenario(model *Model, fn *node.Function) {
	if in.system(fn) {
		model.Scenarios = append(model.Scenarios, Scenario{Func: fn.Name})
	} else {
		in.refuseForm(fn, ScenarioDirective, systemForm)
	}
}

// bindFingerprint sets the fingerprint of the model to the function fn.
func (in *invocation) bindFingerprint(model *Model, fn *node.Function) {
	switch {
	case !fits(paramTypes(fn.Params), in.world) || !fits(returnTypes(fn.Returns), isUint64):
		in.refuseForm(fn, FingerprintDirective, fingerprintForm)
	case model.Fingerprint != "":
		in.twice(fn, FingerprintDirective, model.Fingerprint)
	default:
		model.Fingerprint = fn.Name
	}
}

// bindSpec sets the spec of the model to the function fn, and its state to
// the type argument of the spec that fn returns.
func (in *invocation) bindSpec(model *Model, fn *node.Function) {
	spec, ok := specOf(fn)
	switch {
	case !ok:
		in.refuseForm(fn, SpecDirective, specForm)
	case model.Spec != nil:
		in.twice(fn, SpecDirective, model.Spec.Func)
	case spec.From:
		model.Spec = spec
	default:
		model.Spec, model.State = spec, rules.EmitRef(fn.Returns[0].Type.Args[0])
	}
}

// bindAction adds the action fn to the model, a call of a client where its
// directive states call.
func (in *invocation) bindAction(model *Model, fn *node.Function) {
	call, _ := sdk.FactOf(in.m, fn.ID, in.g.call)
	final, _ := sdk.FactOf(in.m, fn.ID, in.g.final)
	switch {
	case in.action(fn, call):
		model.Actions = append(model.Actions, Action{Func: fn.Name, Call: call, Final: final})
	case call:
		in.refuseForm(fn, ActionDirective, callForm)
	default:
		in.refuseForm(fn, ActionDirective, actionForm)
	}
}

// bindInvariant adds the invariant fn with the defect of its directive to the
// model. It reports [BadFunction] at fn where the defect is no system.
func (in *invocation) bindInvariant(model *Model, fn *node.Function) {
	defect, _ := sdk.FactOf(in.m, fn.ID, in.g.defect)
	d, _ := in.b.View().Lookup(defect)
	// The validation of the directive resolves the defect to a function alone.
	planted, _ := d.(*node.Function)
	switch {
	case !in.check(fn):
		in.refuseForm(fn, InvariantDirective, checkForm)
	case !in.system(planted):
		in.refuse(BadFunction, fn.Pos, "the defect %s of the invariant %s needs the signature %s", defect.Name,
			fn.Name, systemForm)
	default:
		model.Invariants = append(model.Invariants, Invariant{
			Statement: statementOf(fn), Func: fn.Name, Defect: defect.Name,
		})
	}
}

// bindConvergence adds the convergence fn with the bound of its directive to
// the model. It reports [BadFunction] at fn where the bound is no duration.
func (in *invocation) bindConvergence(model *Model, fn *node.Function) {
	within, bounded := sdk.FactOf(in.m, fn.ID, in.g.within)
	w, _ := in.b.View().Lookup(within)
	// The validation of the directive resolves the bound to a function alone.
	limit, _ := w.(*node.Function)
	switch {
	case !in.check(fn):
		in.refuseForm(fn, ConvergesDirective, checkForm)
	case bounded && !bound(limit):
		in.refuse(BadFunction, fn.Pos, "the bound %s of the convergence %s needs the signature %s", within.Name,
			fn.Name, withinForm)
	default:
		model.Converges = append(model.Converges, Convergence{
			Statement: statementOf(fn), Func: fn.Name, Within: within.Name,
		})
	}
}

// system reports whether the function fn has the signature of a system: one
// world of the machine and no result.
func (in *invocation) system(fn *node.Function) bool {
	return fits(paramTypes(fn.Params), in.world) && fits(returnTypes(fn.Returns))
}

// check reports whether the function fn has the signature of an invariant
// or a convergence: the handle of a check, a world of the machine and no
// result.
func (in *invocation) check(fn *node.Function) bool {
	return fits(paramTypes(fn.Params), isCheck, in.world) && fits(returnTypes(fn.Returns))
}

// action reports whether the function fn has the signature of an action: it
// takes a case and returns the function that runs the action on a world of
// the machine, and on the client of the call for an action call. The run
// returns nothing. Only a function type splits its elements into parameters
// and results, so a type whose elements are all parameters is a function
// type.
func (in *invocation) action(fn *node.Function, call bool) bool {
	run := []func(*node.TypeRef) bool{in.world}
	if call {
		run = append(run, isInt)
	}
	returns := returnTypes(fn.Returns)
	return fits(paramTypes(fn.Params), isCase) && len(returns) == 1 && returns[0].Split == len(returns[0].Elems) &&
		fits(returns[0].Elems, run...)
}

// performer reports whether the package of the simulation declares the
// function named name with the signature of the performer of the effect f.
func (in *invocation) performer(name string, f *node.Struct) bool {
	id := symbol.Identity{Lang: f.ID.Lang, Package: in.pkg.Package, Name: name, Kind: symbol.KindFunction}
	d, _ := in.b.View().Lookup(id)
	fn, isFunction := d.(*node.Function)
	effect := func(ref *node.TypeRef) bool { return ref.Target == f.ID }
	return isFunction && fits(paramTypes(fn.Params), in.world, in.process, effect) &&
		fits(returnTypes(fn.Returns), in.outcome)
}

// world reports whether the type ref is a pointer to the world of the
// machine, *sim.World[E, F].
func (in *invocation) world(ref *node.TypeRef) bool {
	return in.machined(ref, "World")
}

// process reports whether the type ref is a pointer to a process of the
// machine, *sim.Process[E, F].
func (in *invocation) process(ref *node.TypeRef) bool {
	return in.machined(ref, "Process")
}

// outcome reports whether the type ref is the outcome of an effect of the
// machine, sim.Outcome[E].
func (in *invocation) outcome(ref *node.TypeRef) bool {
	return named(ref, simPackage, "Outcome") && slices.EqualFunc(ref.Args, []*node.TypeRef{in.event}, sameTarget)
}

// machined reports whether the type ref is a pointer to the generic type
// name of the module sim at the event and the effect interfaces of the
// machine, such as *sim.World[E, F].
func (in *invocation) machined(ref *node.TypeRef, name string) bool {
	t := rules.Unaliased(pointed(ref), in.b.View())
	return named(t, simPackage, name) && slices.EqualFunc(t.Args, []*node.TypeRef{in.event, in.effect}, sameTarget)
}

// specOf returns the spec of the function fn: a function without parameters
// that returns a history.Spec, or a history.Subject that history.SpecFrom
// wraps. It reports false for any other signature.
func specOf(fn *node.Function) (*Spec, bool) {
	returns := returnTypes(fn.Returns)
	if !fits(paramTypes(fn.Params)) || len(returns) != 1 {
		return nil, false
	}
	switch {
	case named(returns[0], historyPackage, "Spec"):
		return &Spec{Func: fn.Name}, true
	case named(returns[0], historyPackage, "Subject"):
		return &Spec{Func: fn.Name, From: true}, true
	}
	return nil, false
}

// bound reports whether the function fn has the signature of the bound of a
// convergence: no parameter and one time.Duration.
func bound(fn *node.Function) bool {
	return fits(paramTypes(fn.Params)) && fits(returnTypes(fn.Returns), isDuration)
}

// fits reports whether refs has one type ref for each check, and whether each
// type ref passes the check at its position.
func fits(refs []*node.TypeRef, checks ...func(*node.TypeRef) bool) bool {
	if len(refs) != len(checks) {
		return false
	}
	for i, check := range checks {
		if !check(refs[i]) {
			return false
		}
	}
	return true
}

// isCheck reports whether the type ref is *sim.Check.
func isCheck(ref *node.TypeRef) bool {
	return named(pointed(ref), simPackage, "Check")
}

// isCase reports whether the type ref is *prop.Case.
func isCase(ref *node.TypeRef) bool {
	return named(pointed(ref), propPackage, "Case")
}

// isDuration reports whether the type ref is time.Duration.
func isDuration(ref *node.TypeRef) bool {
	return named(ref, timePackage, "Duration")
}

// isInt reports whether the type ref is int.
func isInt(ref *node.TypeRef) bool {
	return ref.Spelling == "int"
}

// isUint64 reports whether the type ref is uint64.
func isUint64(ref *node.TypeRef) bool {
	return ref.Spelling == "uint64"
}

// sameTarget reports whether the type refs a and b refer to the same
// declaration.
func sameTarget(a, b *node.TypeRef) bool {
	return a.Target == b.Target
}

// named reports whether the type ref names the declaration name of the
// package pkg.
func named(ref *node.TypeRef, pkg, name string) bool {
	return ref != nil && ref.Target.Package == pkg && ref.Target.Name == name
}

// pointed returns the element of the pointer type ref, and nil for any other
// type.
func pointed(ref *node.TypeRef) *node.TypeRef {
	if ref.Form != symbol.FormOptional {
		return nil
	}
	return ref.Elems[0]
}

// paramTypes returns the types of the parameters params.
func paramTypes(params []*node.Param) []*node.TypeRef {
	out := make([]*node.TypeRef, 0, len(params))
	for _, p := range params {
		out = append(out, p.Type)
	}
	return out
}

// returnTypes returns the types of the results returns.
func returnTypes(returns []*node.Return) []*node.TypeRef {
	out := make([]*node.TypeRef, 0, len(returns))
	for _, r := range returns {
		out = append(out, r.Type)
	}
	return out
}

// statementOf returns the first sentence of the doc comment of the function.
func statementOf(fn *node.Function) string {
	text, _, _ := strings.Cut(strings.Join(fn.Doc, " "), ". ")
	return strings.TrimSuffix(strings.TrimSpace(text), ".")
}

// rulesOf returns the check of each claim of the spec whose generator is
// sim, on the machine interface named machine.
func rulesOf(machine string) []Rule {
	var out []Rule
	claims := spec.Claims()
	for i := range claims {
		c := &claims[i]
		if c.Generator != string(Word) {
			continue
		}
		out = append(out, Rule{
			ID:        machine + "/" + string(c.ID),
			Statement: strings.ReplaceAll(c.Statement, "{declaration}", machine),
			Claim:     string(c.ID),
			Probe:     string(c.Probe),
			Reason:    string(c.Defect.Reason),
			Form:      string(c.Defect.Form),
		})
	}
	return out
}

// refOf returns the reference to the struct s.
func refOf(s *node.Struct) *emit.TypeRef {
	return &emit.TypeRef{Spelling: s.Name, Target: s.ID}
}
