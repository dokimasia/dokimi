// Package edge declares an interface with the names that the methods of a
// double must not reuse.
package edge

// Edge has methods whose parameters and results take the names that a double
// gives a parameter, its receiver, its locals and the locals of its results.
//
//dokimi:stub
type Edge interface {
	// Skip has a parameter named _ and a parameter named e.
	Skip(_ string, e int)
	// Read has a result named r and a result named _.
	Read() (r int, _ error)
	// Keep has a parameter named r.
	Keep(r int) int
	// Use has a parameter named s, which is the name of the receiver.
	Use(s string) error
	// Record has a parameter named call, which is the name of the local of
	// the call.
	Record(call string) (n int, err error)
	// Clash has a result named like the receiver of a double, so the
	// methods of the double declare no names of its results.
	Clash() (s string, err error)
	// Shift has a parameter named _ and a parameter named a0, so the
	// parameter _ takes another name.
	Shift(_ string, a0 int)
	// Fetch has a parameter named _ and a result named a0, so the parameter _
	// takes another name.
	Fetch(_ string) (a0 int, err error)
	// Wide has a parameter named _ at the positions 1 and 10, and a parameter
	// named a1. The parameter _ at the position 1 takes the name a10, so the
	// parameter _ at the position 10 takes another name.
	Wide(a, _, a1, b, c, d, e, f, g, h, _ int)
	// Pair has two results with fields of their own, whose parameters of the
	// setter Returns would have one name.
	Pair() (a_b int, ab int)
}
