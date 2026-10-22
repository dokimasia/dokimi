package status

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Draft-1]
	_ = x[Published-2]
	_ = x[Archived-3]
}

const _Status_name = "DraftPublishedArchived"

var _Status_index = [...]uint8{0, 5, 14, 22}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Status) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Status_index)-1 {
		return "Status(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Status_name[_Status_index[idx]:_Status_index[idx+1]]
}

// ErrUnknownStatus is the error of text that names no declared variant.
// The error of [ParseStatus] wraps it with the text.
var ErrUnknownStatus = errors.New("status: unknown Status value")

// ParseStatus returns the variant of Status whose text is s. It is the
// inverse of [Status.String] on the declared variants. It returns an error
// that wraps [ErrUnknownStatus] for text that names no declared variant.
func ParseStatus(s string) (Status, error) {
	switch s {
	case "Draft":
		return Draft, nil
	case "Published":
		return Published, nil
	case "Archived":
		return Archived, nil
	}
	var zero Status
	return zero, fmt.Errorf("%w %q", ErrUnknownStatus, s)
}

// StatusValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func StatusValues() []Status {
	return []Status{
		Draft,
		Published,
		Archived,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Status) IsValid() bool {
	switch v {
	case Draft, Published, Archived:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key.
func (v Status) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseStatus] for text that names no declared variant.
func (v *Status) UnmarshalText(text []byte) error {
	parsed, err := ParseStatus(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
