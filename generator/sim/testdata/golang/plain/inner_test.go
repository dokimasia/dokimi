package plain

// local is an effect of a type that a test file of the package declares. It
// embeds a Beta, so it implements Effect.
type local struct{ Beta }
