// Package generic declares generic structs under the directive builder.
package generic

// Page has a type parameter with a witness that the rules derive.
//
//dokimi:builder
type Page[T any] struct {
	// Items refers to the type parameter.
	Items []T
	// Next has a default.
	//
	//dokimi:default 1
	Next int
}
