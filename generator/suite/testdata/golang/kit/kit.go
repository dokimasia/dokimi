// Package kit declares the containers of the shop, whose suites the directive
// out routes into the kit kittest.
package kit

// Box is a box of items under the directive out, and the subject of its own
// suite. Its factory is a function of the external test package of the kit.
//
//dokimi:suite
//dokimi:out path=kittest/
//dokimi:subject Box new=newBox
type Box struct {
	items int
}

// Len returns the number of items in the box.
func (b *Box) Len() int {
	return b.items
}

// Crate is a crate of items under the directive out, and a subject of the
// suite of Box. Its factory is a function of the external test package of the
// package kit, which the test in the kit cannot import.
//
//dokimi:out path=kittest/
//dokimi:subject Box new=newCrate
type Crate struct {
	items int
}

// Len returns the number of items in the crate.
func (c *Crate) Len() int {
	return c.items
}
