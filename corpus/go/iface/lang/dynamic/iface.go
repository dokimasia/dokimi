// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package dynamic is the language-axis fixture for parameters of any type and
// for the pointers that the property engine reads as values of their own
// shape.
package dynamic

import (
	"math/big"
	"time"
)

// Diary records messages with values of any type, in a zone and at a scale.
//
//dokimi:suite
//dokimi:out path=dynamictest/
type Diary interface {
	// Print records the message with its values.
	Print(msg string, values ...any)
	// Field records the value of the field.
	Field(name string, value any)
	// In moves the diary into the zone loc.
	In(loc *time.Location)
	// Scale sets the scale of the diary to n.
	Scale(n *big.Int)
	// Count returns the number of records of the diary.
	Count() int
}
