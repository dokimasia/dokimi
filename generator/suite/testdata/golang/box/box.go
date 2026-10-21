// Package box holds generic declarations that have no suite.
package box

// Box holds one value.
//
//dokimi:suite
type Box[T any] struct {
	value T
}

// Get returns the value of the box.
func (b *Box[T]) Get() T {
	return b.value
}

// Mapped returns the result of f for each item.
//
//dokimi:suite
func Mapped[T any](items []T, f func(T) T) []T {
	out := make([]T, 0, len(items))
	for _, item := range items {
		out = append(out, f(item))
	}
	return out
}
