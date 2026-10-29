// Package dual keeps the items of the shop in stores that are the subjects of
// two suites.
package dual

// Reader reads the items of the shop.
//
//dokimi:suite
type Reader interface {
	// Read returns the item of the key.
	Read(key string) (string, error)
}

// Writer writes the items of the shop.
//
//dokimi:suite
type Writer interface {
	// Write stores the item under the key.
	//
	//dokimi:budget allocs=0
	Write(key, item string) error
	// Flush writes the stored items through.
	Flush() error
}

// Both reads and writes the items of the shop. It states a budget of
// Reader.Read, and excuses the ceiling of Writer.Write.
//
//dokimi:subject Reader new=NewBoth
//dokimi:subject Writer new=NewBoth
//dokimi:budget Reader.Read allocs=1
//dokimi:excuse Writer.Write/allocation-budget reason="the map allocates on each write"
type Both struct {
	items map[string]string
}

// NewBoth returns an empty store.
func NewBoth() *Both {
	return &Both{items: map[string]string{}}
}

// Read returns the item of the key.
func (b *Both) Read(key string) (string, error) {
	return b.items[key], nil
}

// Write stores the item under the key.
func (b *Both) Write(key, item string) error {
	b.items[key] = item
	return nil
}

// Flush writes nothing through.
func (*Both) Flush() error { return nil }

// Twice is the subject of Reader twice.
//
//dokimi:subject Reader new=NewTwice
//dokimi:subject Reader new=NewTwice
type Twice struct{}

// NewTwice returns a new Twice.
func NewTwice() *Twice { return &Twice{} }

// Read returns no item.
func (*Twice) Read(string) (string, error) { return "", nil }

// Unruly is the subject of Reader and Writer, and excuses a check of a
// declaration that is neither.
//
//dokimi:subject Reader new=NewUnruly
//dokimi:subject Writer new=NewUnruly
//dokimi:excuse Lister.List/keeps-the-order reason="the store has no list"
type Unruly struct{}

// NewUnruly returns a new Unruly.
func NewUnruly() *Unruly { return &Unruly{} }

// Read returns no item.
func (*Unruly) Read(string) (string, error) { return "", nil }

// Write drops the item.
func (*Unruly) Write(string, string) error { return nil }

// Flush writes nothing through.
func (*Unruly) Flush() error { return nil }

// Mixed is the subject of Reader and Writer. It excuses a check that the
// suite of Reader does not have, and then a check of Writer.
//
//dokimi:subject Reader new=NewMixed
//dokimi:subject Writer new=NewMixed
//dokimi:excuse Reader.Read/keeps-the-key reason="the store forgets each key"
//dokimi:excuse Writer.Write/allocation-budget reason="the map allocates on each write"
type Mixed struct{}

// NewMixed returns a new Mixed.
func NewMixed() *Mixed { return &Mixed{} }

// Read returns no item.
func (*Mixed) Read(string) (string, error) { return "", nil }

// Write drops the item.
func (*Mixed) Write(string, string) error { return nil }

// Flush writes nothing through.
func (*Mixed) Flush() error { return nil }

// Spent is the subject of Reader and Writer. It states a budget of
// Reader.Read that no check can apply, and then a budget of Writer.Flush.
//
//dokimi:subject Reader new=NewSpent
//dokimi:subject Writer new=NewSpent
//dokimi:budget Reader.Read allocs=-1
//dokimi:budget Writer.Flush allocs=2
type Spent struct{}

// NewSpent returns a new Spent.
func NewSpent() *Spent { return &Spent{} }

// Read returns no item.
func (*Spent) Read(string) (string, error) { return "", nil }

// Write drops the item.
func (*Spent) Write(string, string) error { return nil }

// Flush writes nothing through.
func (*Spent) Flush() error { return nil }

// Overspent is the subject of Reader and Writer, and states a budget of a
// declaration that is neither.
//
//dokimi:subject Reader new=NewOverspent
//dokimi:subject Writer new=NewOverspent
//dokimi:budget Lister.List allocs=1
type Overspent struct{}

// NewOverspent returns a new Overspent.
func NewOverspent() *Overspent { return &Overspent{} }

// Read returns no item.
func (*Overspent) Read(string) (string, error) { return "", nil }

// Write drops the item.
func (*Overspent) Write(string, string) error { return nil }

// Flush writes nothing through.
func (*Overspent) Flush() error { return nil }

// Unflushed is the subject of Writer, and excuses a check of Flush that the
// suite of Writer does not have.
//
//dokimi:subject Writer new=NewUnflushed
//dokimi:excuse Writer.Flush/writes-through reason="the store keeps no items"
type Unflushed struct{}

// NewUnflushed returns a new Unflushed.
func NewUnflushed() *Unflushed { return &Unflushed{} }

// Write drops the item.
func (*Unflushed) Write(string, string) error { return nil }

// Flush writes nothing through.
func (*Unflushed) Flush() error { return nil }

// Stray is the subject of Reader, and excuses a check of Writer, which is
// not one of its declarations.
//
//dokimi:subject Reader new=NewStray
//dokimi:excuse Writer.Write/allocation-budget reason="the store is not a writer"
type Stray struct{}

// NewStray returns a new Stray.
func NewStray() *Stray { return &Stray{} }

// Read returns no item.
func (*Stray) Read(string) (string, error) { return "", nil }
