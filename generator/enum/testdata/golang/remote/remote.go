// Package remote declares an enumeration whose values read the constants of
// another package.
package remote

import "time"

// Season is a season of a year. Each value reads a month of the package time,
// so the frontend stamps no exact value of a variant.
//
//dokimi:enum
type Season int

// The seasons of a year, by their first month.
const (
	Spring Season = Season(time.March)
	Autumn Season = Season(time.September)
)
