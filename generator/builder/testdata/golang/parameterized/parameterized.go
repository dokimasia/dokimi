// Package parameterized declares a generic alias of a struct without type
// parameters under the directive builder.
package parameterized

import "example.com/shop/catalog"

// Roled is a generic alias of the struct Role.
//
//dokimi:builder
type Roled[T any] = catalog.Role
