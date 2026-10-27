// Package counter declares a struct and an interface under the directive stub
// with methods that are not exported.
package counter

// step is the change of one call of add.
type step struct {
	delta int
}

// Counter is a struct under the directive stub whose double is in the package
// countertest. Its method add takes the unexported type step, which
// countertest cannot import.
//
//dokimi:stub
//dokimi:out path=countertest/
type Counter struct {
	n int
}

// Inc adds one to the count.
func (c *Counter) Inc() {
	c.add(step{delta: 1})
}

// Value returns the count.
func (c *Counter) Value() int {
	return c.n
}

// add adds the delta of s to the count.
func (c *Counter) add(s step) {
	c.n += s.delta
}

// Gauge is a struct under the directive stub whose double is in its own
// package. Its unexported method adjust takes an int.
//
//dokimi:stub
type Gauge struct {
	level int
}

// Level returns the level of the gauge.
func (g *Gauge) Level() int {
	return g.level
}

// adjust sets the level of the gauge to n.
func (g *Gauge) adjust(n int) {
	g.level = n
}

// Tally is an interface under the directive stub with a method that is not
// exported, so only a type of its package implements it. Its double is in its
// own package.
//
//dokimi:stub
type Tally interface {
	// Total returns the total of the tally.
	Total() int
	// seal limits the implementations of Tally to the package. It can fail.
	seal() error
}
