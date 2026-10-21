// Package gap declares enumerations of integers whose values are not one run
// of consecutive values from one.
package gap

// Gap has three runs of consecutive values.
//
//dokimi:enum
type Gap int

// The values of Gap: two runs of two values and one run of one.
const (
	Two    Gap = 2
	Three  Gap = 3
	Five   Gap = 5
	Six    Gap = 6
	Eleven Gap = 11
)

// Unum is unsigned, and its first run starts at zero.
//
//dokimi:enum
type Unum uint

// The values of Unum: a run from zero and one value far after it.
const (
	Naught Unum = iota
	One
	Far Unum = 253
)

// Prime has more runs of consecutive values than a switch takes.
//
//dokimi:enum
type Prime int

// The primes up to 37: a run of 2 and 3, and ten runs of one value.
const (
	P2  Prime = 2
	P3  Prime = 3
	P5  Prime = 5
	P7  Prime = 7
	P11 Prime = 11
	P13 Prime = 13
	P17 Prime = 17
	P19 Prime = 19
	P23 Prime = 23
	P29 Prime = 29
	P31 Prime = 31
	P37 Prime = 37
)

// Rank declares its values in decreasing order, which are one run.
//
//dokimi:enum
type Rank int

// The ranks, from the highest.
const (
	Gold   Rank = 3
	Silver Rank = 2
	Bronze Rank = 1
)

// Odd has ten runs of one value, as many runs as a switch takes.
//
//dokimi:enum
type Odd int

// The odd numbers below 20.
const (
	O1  Odd = 1
	O3  Odd = 3
	O5  Odd = 5
	O7  Odd = 7
	O9  Odd = 9
	O11 Odd = 11
	O13 Odd = 13
	O15 Odd = 15
	O17 Odd = 17
	O19 Odd = 19
)

// Num is signed, and its one run starts below zero.
//
//dokimi:enum
type Num int

// The values of Num, from minus one.
const (
	Minus Num = iota - 1
	Zero
	Plus
)
