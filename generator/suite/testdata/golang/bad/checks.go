package bad

import "example.com/shop/store"

// Empty is a declaration under the directive suite without callables.
//
//dokimi:suite
type Empty interface{}

// Normalize is a function under the directive suite, which a check names
// without a method.
//
//dokimi:suite
func Normalize(text string) string { return text }

// TrimsTheText checks that Trim of Normalize trims the text.
//
//dokimi:check Normalize.Trim defect=broken
func TrimsTheText(h any) {}

// NamesNoCallable checks that Empty has a callable.
//
//dokimi:check Empty defect=broken
func NamesNoCallable(h any) {}

// GetsNothing checks that Get of a declaration without the directive suite
// finds nothing.
//
//dokimi:check Lister.List defect=broken
func GetsNothing(h any) {}

// FetchesTheItem checks that Fetch returns the item.
//
//dokimi:check store.Store.Fetch defect=broken
func FetchesTheItem(h any) {}

// KeepsAnItem checks that Get returns the item of one of two stores.
//
//dokimi:check Store.Get defect=broken
func KeepsAnItem(h any) {}

// KeepsTheName checks that Get of the twin store returns the name.
//
//dokimi:check twin.Store.Get defect=broken
func KeepsTheName(h any) {}

// PutsTheName checks that Put of the twin store keeps the name.
//
//dokimi:check twin.Store.Put defect=broken
func PutsTheName(h any) {}

// FindsNothing checks that Find of a declaration without callables finds
// nothing.
//
//dokimi:check Empty.Find defect=broken
func FindsNothing(h any) {}

// Undocumented returns the item of the store.
//
//dokimi:check store.Store.Get defect=broken
func Undocumented(h any) {}

// TakesTwo checks that Get takes two handles.
//
//dokimi:check store.Store.Get defect=broken
func TakesTwo(h, other any) {}

// ReturnsOne checks that Get returns a result.
//
//dokimi:check store.Store.Get defect=broken
func ReturnsOne(h any) error { return nil }

// BuildsNothing checks that Get returns the item.
//
//dokimi:check store.Store.Get defect=shapeless
func BuildsNothing(h any) {}

// BuildsTwo checks that Get returns the item.
//
//dokimi:check store.Store.Get defect=paired
func BuildsTwo(h any) {}

// broken returns a store that breaks every check.
func broken(any) store.Store { return nil }

// shapeless returns a store without the handle of a check.
func shapeless() store.Store { return nil }

// paired returns a store and an error.
func paired(any) (store.Store, error) { return nil, nil }
