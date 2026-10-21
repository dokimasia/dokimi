// Package billing declares a value whose builder is in another package than
// the store of the package built.
package billing

// Invoice is a value with a builder.
//
//dokimi:builder
type Invoice struct {
	Number string
}
