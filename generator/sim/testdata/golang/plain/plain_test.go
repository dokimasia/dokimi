package plain_test

import "example.com/shop/plain"

// foreign is an effect of a type that a test of the package declares. It
// embeds a Beta, so it implements plain.Effect.
type foreign struct{ plain.Beta }
