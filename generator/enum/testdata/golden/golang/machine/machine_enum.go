package machine

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[idle-0]
	_ = x[running-1]
}

const _state_name = "idlerunning"

var _state_index = [...]uint8{0, 4, 11}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v state) String() string {
	idx := int(v) - 0
	if v < 0 || idx >= len(_state_index)-1 {
		return "state(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _state_name[_state_index[idx]:_state_index[idx+1]]
}

// errUnknownState is the error of a text that does not match a declared
// variant. The error of [parseState] wraps it with the text.
var errUnknownState = errors.New("machine: unknown state value")

// parseState returns the variant of state whose text is s. It is the
// inverse of [state.String] on the declared variants. It returns an error
// that wraps [errUnknownState] for a text that does not match a declared
// variant.
func parseState(s string) (state, error) {
	switch s {
	case "idle":
		return idle, nil
	case "running":
		return running, nil
	}
	var zero state
	return zero, fmt.Errorf("%w %q", errUnknownState, s)
}

// stateValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func stateValues() []state {
	return []state{
		idle,
		running,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v state) IsValid() bool {
	switch v {
	case idle, running:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [errUnknownState] for a value that [state.IsValid] refuses,
// whose text [parseState] does not parse.
func (v state) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", errUnknownState, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [parseState] for a text that does not match a declared variant.
func (v *state) UnmarshalText(text []byte) error {
	parsed, err := parseState(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
