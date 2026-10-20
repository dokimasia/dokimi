// Package aliased declares aliases under the directive builder.
package aliased

import "example.com/shop/catalog"

// Mirror is an alias of a struct of another package.
//
//dokimi:builder
type Mirror = catalog.Role

// Seconds is an alias of a type that is no struct.
//
//dokimi:builder
type Seconds = int
