package shaped

import (
	"errors"
	"fmt"
)

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Ratio) String() string {
	switch v {
	case Single:
		return "Single"
	case Double:
		return "Double"
	}
	return fmt.Sprintf("Ratio(%g)", float64(v))
}

// ErrUnknownRatio is the error of text that names no declared variant. The
// error of [ParseRatio] wraps it with the text.
var ErrUnknownRatio = errors.New("shaped: unknown Ratio value")

// ParseRatio returns the variant of Ratio whose text is s. It is the
// inverse of [Ratio.String] on the declared variants. It returns an error
// that wraps [ErrUnknownRatio] for text that names no declared variant.
func ParseRatio(s string) (Ratio, error) {
	switch s {
	case "Single":
		return Single, nil
	case "Double":
		return Double, nil
	}
	var zero Ratio
	return zero, fmt.Errorf("%w %q", ErrUnknownRatio, s)
}

// RatioValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func RatioValues() []Ratio {
	return []Ratio{
		Single,
		Double,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Ratio) IsValid() bool {
	switch v {
	case Single, Double:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Ratio) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseRatio] for text that names no declared variant.
func (v *Ratio) UnmarshalText(text []byte) error {
	parsed, err := ParseRatio(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Tier) String() string {
	switch v {
	case Bronze:
		return "Bronze"
	case Silver:
		return "Silver"
	}
	return fmt.Sprintf("Tier(%d)", Base(v))
}

// ErrUnknownTier is the error of text that names no declared variant. The
// error of [ParseTier] wraps it with the text.
var ErrUnknownTier = errors.New("shaped: unknown Tier value")

// ParseTier returns the variant of Tier whose text is s. It is the inverse
// of [Tier.String] on the declared variants. It returns an error that
// wraps [ErrUnknownTier] for text that names no declared variant.
func ParseTier(s string) (Tier, error) {
	switch s {
	case "Bronze":
		return Bronze, nil
	case "Silver":
		return Silver, nil
	}
	var zero Tier
	return zero, fmt.Errorf("%w %q", ErrUnknownTier, s)
}

// TierValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func TierValues() []Tier {
	return []Tier{
		Bronze,
		Silver,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Tier) IsValid() bool {
	switch v {
	case Bronze, Silver:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Tier) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseTier] for text that names no declared variant.
func (v *Tier) UnmarshalText(text []byte) error {
	parsed, err := ParseTier(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
