package priority

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Low-0]
	_ = x[Medium-1]
	_ = x[High-2]
	_ = x[Critical-3]
}

const _Priority_name = "LowMediumHighCritical"

var _Priority_index = [...]uint8{0, 3, 9, 13, 21}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Priority) String() string {
	idx := int(v) - 0
	if v < 0 || idx >= len(_Priority_index)-1 {
		return "Priority(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Priority_name[_Priority_index[idx]:_Priority_index[idx+1]]
}

// ErrUnknownPriority is the error of a text that does not match a declared
// variant. The error of [ParsePriority] wraps it with the text.
var ErrUnknownPriority = errors.New("priority: unknown Priority value")

// ParsePriority returns the variant of Priority whose text is s. It is the
// inverse of [Priority.String] on the declared variants. It returns an
// error that wraps [ErrUnknownPriority] for a text that does not match a
// declared variant.
func ParsePriority(s string) (Priority, error) {
	switch s {
	case "Low":
		return Low, nil
	case "Medium":
		return Medium, nil
	case "High":
		return High, nil
	case "Critical":
		return Critical, nil
	}
	var zero Priority
	return zero, fmt.Errorf("%w %q", ErrUnknownPriority, s)
}

// PriorityValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func PriorityValues() []Priority {
	return []Priority{
		Low,
		Medium,
		High,
		Critical,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Priority) IsValid() bool {
	switch v {
	case Low, Medium, High, Critical:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownPriority] for a value that [Priority.IsValid]
// refuses, whose text [ParsePriority] does not parse.
func (v Priority) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownPriority, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParsePriority] for a text that does not match a declared
// variant.
func (v *Priority) UnmarshalText(text []byte) error {
	parsed, err := ParsePriority(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
