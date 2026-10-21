package store

// KeepsTheItem checks that Get returns the item that Put stored.
//
//dokimi:check Store.Get defect=forgetful
func KeepsTheItem(h *StoreHandle) {}

// RefusesACancelledPut checks that Put reports a cancelled context. Its
// directive has the prefix of the plugin, and its callable has the name of
// the package.
//
//dokimi:suite:check store.Store.Put defect=careless
func RefusesACancelledPut(h *StoreHandle) {}

// forgetful returns a Store that drops every item.
func forgetful(*StoreHandle) Store { return nil }

// careless returns a Store whose Put ignores its context.
func careless(*StoreHandle) Store { return nil }
