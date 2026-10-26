// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deep declares the values of the parts of a tool, which only the
// packages of the tree lib import.
package deep

import "errors"

// ErrWorn reports a worn part.
var ErrWorn = errors.New("deep: the part is worn")

// High is the highest level of a part.
const High = 9
