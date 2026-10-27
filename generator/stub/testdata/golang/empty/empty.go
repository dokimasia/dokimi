// Package empty declares an interface without methods and a struct without
// exported methods under the directive stub.
package empty

// Nothing declares no method, so a double of it implements nothing.
//
//dokimi:stub
type Nothing interface{}

// Hollow declares no exported method, so a double of it implements nothing.
//
//dokimi:stub
type Hollow struct {
	n int
}

// count returns the count of the struct.
func (h *Hollow) count() int {
	return h.n
}
