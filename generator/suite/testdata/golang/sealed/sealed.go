// Package sealed seals and opens messages, counts ticks and keeps items,
// through callables whose builders return inputs that each call accepts.
package sealed

import (
	"context"
	"crypto/cipher"
	"time"
)

// AEAD is an authenticated cipher with a nonce of NonceSize bytes.
//
//dokimi:suite
type AEAD interface {
	cipher.AEAD

	// Seal encrypts and authenticates plaintext, as [cipher.AEAD] documents.
	//
	//dokimi:mixin sample builder=SealInputs
	Seal(dst, nonce, plaintext, additionalData []byte) []byte

	// Open authenticates and decrypts ciphertext, as [cipher.AEAD] documents.
	//
	//dokimi:mixin sample builder=OpenInputs
	Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error)
}

// SealInputs returns the inputs of a call of Seal on a, with a nonce of
// NonceSize bytes.
func SealInputs(a AEAD) (dst, nonce, plaintext, additionalData []byte) {
	return nil, make([]byte, a.NonceSize()), []byte("plaintext"), nil
}

// OpenInputs returns the inputs of a call of Open on a: the ciphertext that
// Seal returns for the inputs of SealInputs.
func OpenInputs(a AEAD) (dst, nonce, ciphertext, additionalData []byte) {
	_, nonce, plaintext, additionalData := SealInputs(a)
	return nil, nonce, a.Seal(nil, nonce, plaintext, additionalData), additionalData
}

// Metronome counts ticks, with a builder that takes no subject.
//
//dokimi:suite
type Metronome interface {
	// Tick advances the ticker by n ticks, and returns the count of ticks.
	//
	//dokimi:mixin sample builder=Ticks
	Tick(ctx context.Context, n int) (int, error)

	// Pace advances the ticker by the ticks of its parameter, and returns the
	// count of ticks.
	//
	//dokimi:mixin sample builder=Ticks
	Pace(context.Context, int) (int, error)
}

// Ticks returns the input of a call of Tick.
func Ticks() int { return 3 }

// Short has a builder that returns fewer values than its inputs.
//
//dokimi:suite
type Short interface {
	// Put stores n under key.
	//
	//dokimi:mixin sample builder=Key
	Put(ctx context.Context, key string, n int) error
}

// Key returns a key.
func Key() string { return "a" }

// Locker keeps values on a medium, with a builder of the inputs of its
// durable write.
//
//dokimi:suite machineries=[medium]
type Locker interface {
	// Put stores the value under the key.
	//
	//dokimi:mixin durable
	//dokimi:mixin sample builder=Entry
	Put(ctx context.Context, key, value string) error

	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
}

// Entry returns the key and the value of a call of Put.
func Entry() (key, value string) { return "k", "v" }

// Shelf keeps numbers by key, with a builder of the inputs of its write.
//
//dokimi:suite
type Shelf interface {
	// Put stores n under key.
	//
	//dokimi:mixin sample builder=Shelved
	Put(ctx context.Context, key string, n int) error

	// Get returns the number under key.
	//
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, key string) (int, error)
}

// Shelved returns the key and the number of a call of Put.
func Shelved() (key string, n int) { return "a", 1 }

// Wrong has a builder that returns values of other types than its inputs.
//
//dokimi:suite
type Wrong interface {
	// Put stores n under key.
	//
	//dokimi:mixin sample builder=Strings
	Put(ctx context.Context, key string, n int) error
}

// Strings returns two strings.
func Strings() (key, n string) { return "a", "b" }

// Taker has a builder that takes another parameter than the subject.
//
//dokimi:suite
type Taker interface {
	// Get returns the value of the key.
	//
	//dokimi:mixin sample builder=Takes
	Get(ctx context.Context, key string) (string, error)
}

// Takes returns a key for the number n.
func Takes(n int) string { return string(rune('a' + n)) }

// Concealed has a builder that the file of its suite cannot import.
//
//dokimi:suite
//dokimi:out path=sealedtest/
type Concealed interface {
	// Count returns the count of n.
	//
	//dokimi:mixin sample builder=hiddenInputs
	Count(ctx context.Context, n int) (int, error)
}

// hiddenInputs returns the input of a call of Count.
func hiddenInputs() int { return 1 }

// Track is a track of a mix.
type Track struct{ Name string }

// Mixer mixes tracks, with a builder of the inputs of its mix.
//
//dokimi:suite
type Mixer interface {
	// Mix mixes the track for the duration time at the levels subject and
	// subject0.
	//
	//dokimi:mixin sample builder=Mixed
	Mix(ctx context.Context, _ int, time time.Duration, subject, subject0 string, track *Track) error
}

// Mixed returns the inputs of a call of Mix.
func Mixed() (int, time.Duration, string, string, *Track) {
	return 1, time.Second, "low", "high", &Track{Name: "a"}
}

// Item is an item of a depot.
type Item struct {
	SKU  string
	Name string
}

// Depot keeps items by their SKU, with a builder of the inputs of each
// callable.
//
//dokimi:suite
type Depot interface {
	// Stock keeps the item.
	//
	//dokimi:mixin sample builder=Stocked
	Stock(ctx context.Context, item Item) error

	// Find returns the item with the SKU, where match accepts it.
	//
	//dokimi:mixin read-after-write write=Stock
	//dokimi:mixin sample builder=Found
	Find(ctx context.Context, sku string, match func(Item) bool) (Item, error)
}

// Stocked returns the item of a call of Stock.
func Stocked() Item { return Item{SKU: "a", Name: "b"} }

// Found returns the SKU and the match of a call of Find.
func Found() (string, func(Item) bool) { return "a", func(Item) bool { return true } }

// Stamp returns the stamp of the text, with a builder that takes another
// parameter than the subject.
//
//dokimi:suite
//dokimi:mixin sample builder=Takes
func Stamp(text string) string { return text }

// Stall counts items, with a builder that takes the struct of its suite.
//
//dokimi:suite
type Stall struct{ n int }

// Add adds n items to the stall.
//
//dokimi:mixin sample builder=StallInputs
func (s *Stall) Add(n int) { s.n += n }

// StallInputs returns the input of a call of Add on s.
func StallInputs(s Stall) int { return s.n + 1 }
