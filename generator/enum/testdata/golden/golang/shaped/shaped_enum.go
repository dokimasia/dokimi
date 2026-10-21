package shaped

import (
	"errors"
	"fmt"
)

// String returns Ratio's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Ratio) String() string {
	switch v {
	case Single:
		return "Single"
	case Double:
		return "Double"
	}
	return fmt.Sprintf("Ratio(%g)", float64(v))
}

// ErrUnknownRatio reports text that names no declared variant.
var ErrUnknownRatio = errors.New("shaped: unknown Ratio value")

// ParseRatio maps text back to a Ratio.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseRatio(s string) (Ratio, error) {
	switch s {
	case "Single":
		return Single, nil
	case "Double":
		return Double, nil
	}
	var zero Ratio
	return zero, ErrUnknownRatio
}

// RatioValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func RatioValues() []Ratio {
	return []Ratio{
		Single,
		Double,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Ratio) IsValid() bool {
	switch v {
	case Single:
		return true
	case Double:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Ratio legal as a map key —
// which a JSON marshaller alone does not.
func (v Ratio) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Ratio) UnmarshalText(text []byte) error {
	parsed, err := ParseRatio(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns Tier's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Tier) String() string {
	switch v {
	case Bronze:
		return "Bronze"
	case Silver:
		return "Silver"
	}
	return fmt.Sprintf("Tier(%d)", Base(v))
}

// ErrUnknownTier reports text that names no declared variant.
var ErrUnknownTier = errors.New("shaped: unknown Tier value")

// ParseTier maps text back to a Tier.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
func ParseTier(s string) (Tier, error) {
	switch s {
	case "Bronze":
		return Bronze, nil
	case "Silver":
		return Silver, nil
	}
	var zero Tier
	return zero, ErrUnknownTier
}

// TierValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func TierValues() []Tier {
	return []Tier{
		Bronze,
		Silver,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Tier) IsValid() bool {
	switch v {
	case Bronze:
		return true
	case Silver:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Tier legal as a map key —
// which a JSON marshaller alone does not.
func (v Tier) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Tier) UnmarshalText(text []byte) error {
	parsed, err := ParseTier(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
