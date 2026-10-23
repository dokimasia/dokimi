package specced_test

import "example.com/shop/specced"

// newStore returns the store of the conformance test of Store.
func newStore() *specced.Store {
	return &specced.Store{}
}

// newVault returns the vault of the conformance test of Vault.
func newVault() *specced.Vault {
	return &specced.Vault{}
}
