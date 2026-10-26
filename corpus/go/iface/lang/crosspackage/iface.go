// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package crosspackage is the fixture of a suite whose inputs are types of
// another tree. The builder of a part is in another package than the suite,
// and the default of a field and an error of the directive fault are values
// of an internal package that the suite cannot import.
package crosspackage

import (
	"context"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/crosspackage/lib/part"
)

// Fitter fits the parts of a tool.
//
//dokimi:suite
//dokimi:out path=crosspackagetest/
type Fitter interface {
	// Fit fits the part.
	Fit(ctx context.Context, p part.Part) error
	// Limit limits the fit to the limits.
	Limit(ctx context.Context, l part.Limits) error
	// Sense senses the part through the sensor.
	Sense(ctx context.Context, s part.Sensor) error
}
