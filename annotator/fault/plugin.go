// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fault

import (
	"errors"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// ID is the name of the plugin.
const ID plugin.ID = "fault"

// Version is the version of the plugin. The fingerprint of a composition
// folds it in.
const Version = "0.1.0"

// Capability is the capability that fault provides. The stub generator
// requires it, so it runs after fault.
const Capability plugin.Capability = "fault"

// FaultDirective configures the faults of the double of a method. ErrorsParam
// lists the error values of the package that a helper of the double arms one
// at a time. RetryParam is the attempt that a retry schedule lets succeed, and
// PartitionParam is the parameter whose argument selects the calls that a
// partition fault fails.
const (
	FaultDirective directive.Name     = "fault"
	ErrorsParam    directive.ParamKey = "errors"
	RetryParam     directive.ParamKey = "retry"
	PartitionParam directive.ParamKey = "partition"
)

// KeyErrors is the key of the fact that contains the identities of the error
// values of errors, in the order of the param. KeyRetry is the key of the
// fact that contains the attempt of retry. KeyPartition is the key of the
// fact that contains the identity of the parameter of partition.
const (
	KeyErrors    meta.KeyName = "fault.errors"
	KeyRetry     meta.KeyName = "fault.retry"
	KeyPartition meta.KeyName = "fault.partition"
)

// leastRetry is the least attempt of a retry schedule. The schedule fails
// each call before the attempt, so an earlier attempt fails no call.
const leastRetry = 2

// schema is the directive of fault. eidos resolves each name of the param
// errors to a variable of the package, and the param partition to a
// parameter of the method.
var schema = directive.Schema{
	Name: FaultDirective,
	Params: []directive.ParamSpec{
		{
			Key: ErrorsParam, Type: directive.TypeList, ListOf: directive.TypeReference,
			Resolution: directive.ResolvePackageVar,
			Doc:        "the error values of the package, each of which a helper of the double arms for the next call",
		},
		{
			Key: RetryParam, Type: directive.TypeInt,
			Doc: "the attempt that a retry schedule lets succeed, 2 or more, which each call before it fails",
		},
		{
			Key: PartitionParam, Type: directive.TypeReference, Resolution: directive.ResolveHostParam,
			Doc: "the parameter whose argument selects the calls that a partition fault fails",
		},
	},
	Doc: "generates the helpers of the stub of the method that arm an error value, a retry schedule and the faults " +
		"of a partition",
}

// conclusion is the state of one instance of fault: the handles of its keys,
// which the registration sets and the rule stamps through.
//
// # Concurrency
//
// The composition registers the keys before the first phase call, and the
// handler reads them afterwards, so the workers of a phase call read them
// concurrently.
type conclusion struct {
	values    meta.Key[[]symbol.Identity]
	retry     meta.Key[int64]
	partition meta.Key[symbol.Identity]
}

// New returns fault, an annotator. For each instance of [FaultDirective] on
// a method, it stamps the identities of the error values of errors under
// [KeyErrors], the attempt of retry under [KeyRetry] and the identity of the
// parameter of partition under [KeyPartition]. It reports [Infallible] at an
// instance on a method that cannot fail, and stamps nothing of that instance.
// It reports [BadRetry] at a retry below 2, and stamps no retry of that
// instance. Each call returns a new plugin, because a plugin instance belongs
// to one workspace.
func New() plugin.Annotator {
	c := &conclusion{}
	return sdk.NewPlugin(ID).
		Version(Version).
		Provides(Capability).
		Keys(c.register).
		Handle(sdk.Directive(schema, sdk.OnMethod(c.declared))).
		Annotator()
}

// register claims the namespace of fault and registers [KeyErrors],
// [KeyRetry] and [KeyPartition].
//
// # Errors
//
// register returns the error of the claim, such as a namespace that another
// registrant claimed, joined with the errors of the three registrations, such
// as a key that the registrant registered with another type.
func (c *conclusion) register(r *meta.Registry) error {
	errClaim := r.ClaimNamespace(KeyRetry.Namespace())
	var errValues, errRetry, errPartition error
	c.values, errValues = meta.Register[[]symbol.Identity](r, meta.KeySpec{
		Name:  KeyErrors,
		Kinds: []symbol.Kind{symbol.KindMethod},
		Doc:   "contains the error values that the helpers of the double of a method arm, one helper each",
	})
	c.retry, errRetry = meta.Register[int64](r, meta.KeySpec{
		Name:  KeyRetry,
		Kinds: []symbol.Kind{symbol.KindMethod},
		Doc:   "contains the attempt that the retry schedule of the double of a method lets succeed",
	})
	c.partition, errPartition = meta.Register[symbol.Identity](r, meta.KeySpec{
		Name:  KeyPartition,
		Kinds: []symbol.Kind{symbol.KindMethod},
		Doc:   "contains the parameter whose argument selects the calls that a partition fault of a double fails",
	})
	return errors.Join(errClaim, errValues, errRetry, errPartition)
}

// declared stamps the params of one instance of the directive on its method.
// A method fails where the error model of its language gives it a failure
// channel, such as a Go method whose last result is an error.
func (c *conclusion) declared(m *sdk.MethodMatch, st *sdk.Stamper) error {
	d := m.Directive()
	if view, _ := m.Rules().CallableOf(m.Method); view.Errors == rules.ErrorsNone {
		m.ErrorfAt(Infallible, d.Pos, "%s cannot fail, so a fault of its double has no error to return",
			m.Method.Name)
		return nil
	}
	if listed, stated := d.Param(ErrorsParam); stated {
		values := make([]symbol.Identity, 0, len(listed.List))
		for i := range listed.List {
			values = append(values, listed.List[i].Target)
		}
		sdk.Stamp(st, c.values, values)
	}
	if retry, stated := d.Param(RetryParam); stated {
		if retry.Int < leastRetry {
			m.ErrorfAt(BadRetry, d.Pos, "retry=%d is below %d, the first attempt that a retry schedule can let succeed",
				retry.Int, leastRetry)
		} else {
			sdk.Stamp(st, c.retry, retry.Int)
		}
	}
	if partition, stated := d.Param(PartitionParam); stated {
		sdk.Stamp(st, c.partition, partition.Target)
	}
	return nil
}
