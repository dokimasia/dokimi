// Package fielded declares structs under the directive stub with the param
// field.
package fielded

// Cell lists the shards of each node. The field OnNode that configures Node
// under the default template would have the name of the method OnNode.
//
//dokimi:stub field={{.Method}}Stub
type Cell struct {
	shards map[string][]string
}

// Node reports whether the cell has the node n.
func (c *Cell) Node(n string) bool {
	_, ok := c.shards[n]
	return ok
}

// OnNode returns the shards on the node n.
func (c *Cell) OnNode(n string) []string {
	return c.shards[n]
}

// Fixed has a param field without {{.Method}}, so each method would have the
// same field.
//
//dokimi:stub field=Stub
type Fixed struct{}

// Get returns nothing.
func (f *Fixed) Get() {}

// Named has a param field with an action other than {{.Method}}.
//
//dokimi:stub field=On{{.Name}}
type Named struct{}

// Get returns nothing.
func (n *Named) Get() {}
