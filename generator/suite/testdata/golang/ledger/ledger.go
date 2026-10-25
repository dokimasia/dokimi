// Package ledger keeps the balances of the accounts of the shop.
package ledger

// Ledger keeps a balance for each account.
//
//dokimi:suite
type Ledger interface {
	// Credit adds the amount to the balance of the account.
	Credit(account string, amount int)
	// Balance returns the balance of the account.
	Balance(account string) int
}

// Book is a Ledger in memory.
//
//dokimi:subject Ledger new=NewBook
type Book struct {
	balances map[string]int
}

// NewBook returns an empty book.
func NewBook() *Book {
	return &Book{balances: map[string]int{}}
}

// Credit adds the amount to the balance of the account.
func (b *Book) Credit(account string, amount int) {
	b.balances[account] += amount
}

// Balance returns the balance of the account.
func (b *Book) Balance(account string) int {
	return b.balances[account]
}
