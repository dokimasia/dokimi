// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

import "go.dokimi.dev/eidos/plugin/shape"

// ClaimID identifies a claim. It is the segment of the identity of each
// check of the claim, such as zero-on-error.
type ClaimID string

// Claim is a statement that a classification licenses, with the law that its
// check witnesses. A claim binds to a callable that has one of its licences
// and passes all its gates, and then becomes one check.
type Claim struct {
	ID ClaimID `json:"id"`

	// Generator is the generator that binds the claim: suite, enum,
	// sentinel or sim.
	Generator string `json:"generator"`

	// Licences are the classifications that license the claim. A callable
	// needs one of them. A claim of enum, sentinel or sim has none, and
	// binds to each subject of its generator.
	Licences []Licence `json:"licences,omitempty"`

	// Unless withdraws the claim from a subject that has it.
	Unless *Unless `json:"unless,omitempty"`

	// Scope is callable, the default, or declaration for a claim on the
	// whole declaration, whose checks have no callable in their identity.
	Scope string `json:"scope,omitempty"`

	// Gates are the conditions that the facts of the subject must meet.
	Gates []GateID `json:"gates,omitempty"`

	// Partners are the other callables that the probe calls.
	Partners []Partner `json:"partners,omitempty"`

	// Law is the law that the check witnesses.
	Law LawID `json:"law"`

	// Machineries are what the run supplies before the check exercises its
	// claim. The first is the deepest.
	Machineries []MachineryID `json:"machineries"`

	// Statement is the sentence of the check. The generator replaces
	// {declaration}, {callable} and the name of each partner in braces.
	Statement string `json:"statement"`

	// Probe is the probe that the check calls.
	Probe ProbeID `json:"probe"`

	// Inputs are the bindings of the inputs of the probe.
	Inputs []Binding `json:"inputs"`

	// Defect is the planted defect of the check.
	Defect Defect `json:"defect"`

	Since string `json:"since"`

	// Entails are the claims on the same callable whose laws the law of
	// this claim implies.
	Entails []ClaimID `json:"entails,omitempty"`
}

func (c *Claim) key() ClaimID { return c.ID }

// Classification is a shape, a mixin, or a role of a contract of the shape
// catalog. Exactly one of Shape, Mixin and Contract is set.
type Classification struct {
	Shape    shape.Shape    `json:"shape,omitempty"`
	Mixin    shape.Mixin    `json:"mixin,omitempty"`
	Contract shape.Contract `json:"contract,omitempty"`

	// Role is the role of the callable in Contract.
	Role shape.Role `json:"role,omitempty"`
}

// Licence is one classification that licenses a claim. Exactly one of
// Callable, Context and the classification is set.
type Licence struct {
	Classification

	// With is a second classification that the callable must also have,
	// for a claim of two classifications together.
	With *Classification `json:"with,omitempty"`

	// NeedsParam is a param of the classification that the directive must
	// set, such as attempts of the mixin retry-succeeds.
	NeedsParam string `json:"needs-param,omitempty"`

	// Conflicts are the classifications whose presence on the callable
	// withdraws the licence.
	Conflicts []Classification `json:"conflicts,omitempty"`

	// Callable licenses the claim on every callable.
	Callable bool `json:"callable,omitempty"`

	// Context licenses the claim on every callable with a parameter in the
	// role context.
	Context bool `json:"context,omitempty"`
}

// Unless withdraws a claim from a subject. A claim of suite names a mixin.
// A claim of enum or sentinel names a param of the directive and the value
// that withdraws it. A YAML file writes the mixin alone as a scalar, and the
// rendered file writes it under the key mixin.
type Unless struct {
	Mixin shape.Mixin `json:"mixin,omitempty"`
	Param string      `json:"param,omitempty"`
	Value string      `json:"value,omitempty"`
}

// PartnerMode is how the derivation finds the partner of a claim.
type PartnerMode string

// The partner modes. The derivation implements each mode once.
const (
	// ModeNamed is the callable that a param of the licence names. The
	// licensed callable reads, and the partner writes.
	ModeNamed PartnerMode = "named"

	// ModeReader is the callable that a param of the licence names, with
	// the halves swapped: the licensed callable writes, and the partner
	// reads.
	ModeReader PartnerMode = "reader"

	// ModeDerived is the one other callable that takes the type of the
	// result field that a param of the licence names.
	ModeDerived PartnerMode = "derived"

	// ModeRole is the callable with the named role of the same contract
	// instance.
	ModeRole PartnerMode = "role"

	// ModeInterposed is a named partner that the probe calls twice, with
	// reads before and after the second call, and a field of the written
	// value that the second call changes.
	ModeInterposed PartnerMode = "interposed"

	// ModeRewriter is a reader partner that the probe calls twice, with the
	// same field as ModeInterposed.
	ModeRewriter PartnerMode = "rewriter"

	// ModeEstablishing is the one other callable of the declaration with
	// the classification of the partner that takes every input of the
	// licensed callable.
	ModeEstablishing PartnerMode = "establishing"
)

// Partner is another callable that the probe of a claim calls.
type Partner struct {
	// Name is the name of the partner in the inputs and the statement of
	// the claim.
	Name string `json:"name"`

	Mode PartnerMode `json:"mode"`

	// Param is the param of the licence that the modes named, reader,
	// derived, interposed and rewriter read.
	Param string `json:"param,omitempty"`

	// Classification is the role of the same contract that the mode role
	// reads, or the shape or the mixin that the mode establishing reads.
	Classification
}

// Binding is one input of a probe and the source of its value.
type Binding struct {
	// Input is the name of the input in the probe.
	Input string `json:"name"`

	Source
}

// SourceKind is the kind of fact from which a check takes the value of an
// input.
type SourceKind string

// The kinds of source.
const (
	// FromCallable is a call of the callable itself with sample inputs.
	FromCallable SourceKind = "callable"

	// FromPartner is a call of a partner of the claim.
	FromPartner SourceKind = "partner"

	// FromSample is a sample of the type of a result.
	FromSample SourceKind = "sample"

	// FromZero is the zero value of the type of a result.
	FromZero SourceKind = "zero"

	// FromParam is the value of a param of the classification, or of the
	// directive for a claim of enum or sentinel.
	FromParam SourceKind = "param"

	// FromClock is the controlled clock of the run.
	FromClock SourceKind = "clock"

	// FromMembers are the members of an enumeration with their names.
	FromMembers SourceKind = "members"

	// FromName is the function that returns the name of a member of an
	// enumeration.
	FromName SourceKind = "name"

	// FromParse is the function that parses a name into a member.
	FromParse SourceKind = "parse"

	// FromList is the function that lists the members of an enumeration.
	FromList SourceKind = "list"

	// FromValid is the function that reports whether a value is a member of
	// an enumeration.
	FromValid SourceKind = "valid"

	// FromOutside is a value of the type of an enumeration that no member
	// has.
	FromOutside SourceKind = "outside"

	// FromErrorValues are the error values of a package with their names.
	FromErrorValues SourceKind = "error-values"

	// FromErrorTypes are the values of the error types of a package with
	// their names.
	FromErrorTypes SourceKind = "error-types"

	// FromCauses are the functions of the error types of a package with a
	// cause. Each makes an error of its type with a given cause.
	FromCauses SourceKind = "causes"
)

// Source is the fact of the subject from which a check takes the value of
// an input. The subject is a callable, an enumeration, the errors of a
// package or a machine interface.
type Source struct {
	Kind SourceKind `json:"from"`

	// Param is the param of a FromParam source.
	Param string `json:"param,omitempty"`

	// Partner is the name of the partner of a FromPartner source.
	Partner string `json:"partner,omitempty"`

	// Sample is the index of the sample of a FromSample source, 0 or 1.
	Sample int `json:"sample,omitempty"`
}

// Defect is the planted defect of a claim.
type Defect struct {
	// Form is the defect form that each template tree writes.
	Form DefectFormID `json:"form"`

	// Clause describes the defect. The generator replaces the names in
	// braces, as in the statement of the claim.
	Clause string `json:"clause"`

	// Reason is the reason code with which the check must fail on the
	// defect.
	Reason ReasonID `json:"reason"`
}

var claims = section[ClaimID, Claim]("claims.json", "claims")

// Claims returns every claim in the order of claims.json, which is the order
// in which the derivation visits them. The slice is shared with every other
// call, so a caller does not modify it.
func Claims() []Claim { return claims().entries }

// ClaimOf returns the claim with the identifier id, and false when the spec
// has no such claim.
func ClaimOf(id ClaimID) (Claim, bool) { return claims().of(id) }
