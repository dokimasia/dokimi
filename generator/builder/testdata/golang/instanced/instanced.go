// Package instanced declares an alias of an instance of a generic struct
// under the directive builder.
package instanced

import "example.com/shop/generic"

// Pages is an alias of an instance of the generic struct Page.
//
//dokimi:builder
type Pages = generic.Page[int]
