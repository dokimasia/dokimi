// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package part declares the parts of a tool. Their values refer to the
// internal package of the tree lib, and the builder of a part is in the
// package parttest.
//
//dokimi:out path=parttest/
package part

import (
	"errors"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/crosspackage/lib/internal/deep"
)

// ErrJammed reports a jammed part.
var ErrJammed = errors.New("part: the part is jammed")

// MaxLevel is the highest level of a part.
const MaxLevel = deep.High

// Part is a part of a tool, with a builder.
//
//dokimi:builder
type Part struct {
	Name  string
	Count int
}

// Limits are the limits of a part. The default of Level is a constant of the
// internal package of the tree lib.
type Limits struct {
	Level int `default:"deep.High"`
	Count int `default:"2"`
}

// Sensor senses a part. Its errors are error values of the internal package of
// the tree lib and of this package.
type Sensor interface {
	// Sense senses the part.
	//
	//dokimi:fault errors=[deep.ErrWorn, ErrJammed]
	Sense() error
}
