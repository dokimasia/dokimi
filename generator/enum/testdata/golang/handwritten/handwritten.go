// Package handwritten declares an enumeration with its own render and parse
// function.
package handwritten

import "errors"

// ErrUnknownWeekday reports a text of no weekday.
var ErrUnknownWeekday = errors.New("handwritten: unknown weekday")

// Weekday is a day of a work week.
//
//dokimi:enum
type Weekday int

// The days of a work week.
const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
)

// String returns the short name of the day.
func (w Weekday) String() string {
	switch w {
	case Monday:
		return "mon"
	case Tuesday:
		return "tue"
	case Wednesday:
		return "wed"
	}
	return "unknown"
}

// ParseWeekday returns the day of a short name.
func ParseWeekday(s string) (Weekday, error) {
	for _, w := range []Weekday{Monday, Tuesday, Wednesday} {
		if w.String() == s {
			return w, nil
		}
	}
	return 0, ErrUnknownWeekday
}
