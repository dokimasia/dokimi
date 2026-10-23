package gap

import (
	"errors"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Two-2]
	_ = x[Three-3]
	_ = x[Five-5]
	_ = x[Six-6]
	_ = x[Eleven-11]
}

const (
	_Gap_name_0 = "TwoThree"
	_Gap_name_1 = "FiveSix"
	_Gap_name_2 = "Eleven"
)

var (
	_Gap_index_0 = [...]uint8{0, 3, 8}
	_Gap_index_1 = [...]uint8{0, 4, 7}
)

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Gap) String() string {
	switch {
	case 2 <= v && v <= 3:
		v -= 2
		return _Gap_name_0[_Gap_index_0[v]:_Gap_index_0[v+1]]
	case 5 <= v && v <= 6:
		v -= 5
		return _Gap_name_1[_Gap_index_1[v]:_Gap_index_1[v+1]]
	case v == 11:
		return _Gap_name_2
	default:
		return "Gap(" + strconv.FormatInt(int64(v), 10) + ")"
	}
}

// ErrUnknownGap is the error of a text that does not match a declared
// variant. The error of [ParseGap] wraps it with the text.
var ErrUnknownGap = errors.New("gap: unknown Gap value")

// ParseGap returns the variant of Gap whose text is s. It is the inverse
// of [Gap.String] on the declared variants. It returns an error that wraps
// [ErrUnknownGap] for a text that does not match a declared variant.
func ParseGap(s string) (Gap, error) {
	switch s {
	case "Two":
		return Two, nil
	case "Three":
		return Three, nil
	case "Five":
		return Five, nil
	case "Six":
		return Six, nil
	case "Eleven":
		return Eleven, nil
	}
	var zero Gap
	return zero, fmt.Errorf("%w %q", ErrUnknownGap, s)
}

// GapValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func GapValues() []Gap {
	return []Gap{
		Two,
		Three,
		Five,
		Six,
		Eleven,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Gap) IsValid() bool {
	switch v {
	case Two, Three, Five, Six, Eleven:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownGap] for a value that [Gap.IsValid] refuses, whose
// text [ParseGap] does not parse.
func (v Gap) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownGap, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseGap] for a text that does not match a declared variant.
func (v *Gap) UnmarshalText(text []byte) error {
	parsed, err := ParseGap(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Naught-0]
	_ = x[One-1]
	_ = x[Far-253]
}

const (
	_Unum_name_0 = "NaughtOne"
	_Unum_name_1 = "Far"
)

var _Unum_index_0 = [...]uint8{0, 6, 9}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Unum) String() string {
	switch {
	case v <= 1:
		return _Unum_name_0[_Unum_index_0[v]:_Unum_index_0[v+1]]
	case v == 253:
		return _Unum_name_1
	default:
		return "Unum(" + strconv.FormatInt(int64(v), 10) + ")"
	}
}

// ErrUnknownUnum is the error of a text that does not match a declared
// variant. The error of [ParseUnum] wraps it with the text.
var ErrUnknownUnum = errors.New("gap: unknown Unum value")

// ParseUnum returns the variant of Unum whose text is s. It is the inverse
// of [Unum.String] on the declared variants. It returns an error that
// wraps [ErrUnknownUnum] for a text that does not match a declared
// variant.
func ParseUnum(s string) (Unum, error) {
	switch s {
	case "Naught":
		return Naught, nil
	case "One":
		return One, nil
	case "Far":
		return Far, nil
	}
	var zero Unum
	return zero, fmt.Errorf("%w %q", ErrUnknownUnum, s)
}

// UnumValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func UnumValues() []Unum {
	return []Unum{
		Naught,
		One,
		Far,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Unum) IsValid() bool {
	switch v {
	case Naught, One, Far:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownUnum] for a value that [Unum.IsValid] refuses,
// whose text [ParseUnum] does not parse.
func (v Unum) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownUnum, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseUnum] for a text that does not match a declared variant.
func (v *Unum) UnmarshalText(text []byte) error {
	parsed, err := ParseUnum(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[P2-2]
	_ = x[P3-3]
	_ = x[P5-5]
	_ = x[P7-7]
	_ = x[P11-11]
	_ = x[P13-13]
	_ = x[P17-17]
	_ = x[P19-19]
	_ = x[P23-23]
	_ = x[P29-29]
	_ = x[P31-31]
	_ = x[P37-37]
}

const _Prime_name = "P2P3P5P7P11P13P17P19P23P29P31P37"

var _Prime_map = map[Prime]string{
	2:  _Prime_name[0:2],
	3:  _Prime_name[2:4],
	5:  _Prime_name[4:6],
	7:  _Prime_name[6:8],
	11: _Prime_name[8:11],
	13: _Prime_name[11:14],
	17: _Prime_name[14:17],
	19: _Prime_name[17:20],
	23: _Prime_name[20:23],
	29: _Prime_name[23:26],
	31: _Prime_name[26:29],
	37: _Prime_name[29:32],
}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Prime) String() string {
	if str, ok := _Prime_map[v]; ok {
		return str
	}
	return "Prime(" + strconv.FormatInt(int64(v), 10) + ")"
}

// ErrUnknownPrime is the error of a text that does not match a declared
// variant. The error of [ParsePrime] wraps it with the text.
var ErrUnknownPrime = errors.New("gap: unknown Prime value")

// ParsePrime returns the variant of Prime whose text is s. It is the
// inverse of [Prime.String] on the declared variants. It returns an error
// that wraps [ErrUnknownPrime] for a text that does not match a declared
// variant.
func ParsePrime(s string) (Prime, error) {
	switch s {
	case "P2":
		return P2, nil
	case "P3":
		return P3, nil
	case "P5":
		return P5, nil
	case "P7":
		return P7, nil
	case "P11":
		return P11, nil
	case "P13":
		return P13, nil
	case "P17":
		return P17, nil
	case "P19":
		return P19, nil
	case "P23":
		return P23, nil
	case "P29":
		return P29, nil
	case "P31":
		return P31, nil
	case "P37":
		return P37, nil
	}
	var zero Prime
	return zero, fmt.Errorf("%w %q", ErrUnknownPrime, s)
}

// PrimeValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func PrimeValues() []Prime {
	return []Prime{
		P2,
		P3,
		P5,
		P7,
		P11,
		P13,
		P17,
		P19,
		P23,
		P29,
		P31,
		P37,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Prime) IsValid() bool {
	switch v {
	case P2, P3, P5, P7, P11, P13, P17, P19, P23, P29, P31, P37:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownPrime] for a value that [Prime.IsValid] refuses,
// whose text [ParsePrime] does not parse.
func (v Prime) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownPrime, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParsePrime] for a text that does not match a declared variant.
func (v *Prime) UnmarshalText(text []byte) error {
	parsed, err := ParsePrime(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Gold-3]
	_ = x[Silver-2]
	_ = x[Bronze-1]
}

const _Rank_name = "BronzeSilverGold"

var _Rank_index = [...]uint8{0, 6, 12, 16}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Rank) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Rank_index)-1 {
		return "Rank(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Rank_name[_Rank_index[idx]:_Rank_index[idx+1]]
}

// ErrUnknownRank is the error of a text that does not match a declared
// variant. The error of [ParseRank] wraps it with the text.
var ErrUnknownRank = errors.New("gap: unknown Rank value")

// ParseRank returns the variant of Rank whose text is s. It is the inverse
// of [Rank.String] on the declared variants. It returns an error that
// wraps [ErrUnknownRank] for a text that does not match a declared
// variant.
func ParseRank(s string) (Rank, error) {
	switch s {
	case "Gold":
		return Gold, nil
	case "Silver":
		return Silver, nil
	case "Bronze":
		return Bronze, nil
	}
	var zero Rank
	return zero, fmt.Errorf("%w %q", ErrUnknownRank, s)
}

// RankValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func RankValues() []Rank {
	return []Rank{
		Gold,
		Silver,
		Bronze,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Rank) IsValid() bool {
	switch v {
	case Gold, Silver, Bronze:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownRank] for a value that [Rank.IsValid] refuses,
// whose text [ParseRank] does not parse.
func (v Rank) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownRank, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseRank] for a text that does not match a declared variant.
func (v *Rank) UnmarshalText(text []byte) error {
	parsed, err := ParseRank(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[O1-1]
	_ = x[O3-3]
	_ = x[O5-5]
	_ = x[O7-7]
	_ = x[O9-9]
	_ = x[O11-11]
	_ = x[O13-13]
	_ = x[O15-15]
	_ = x[O17-17]
	_ = x[O19-19]
}

const (
	_Odd_name_0 = "O1"
	_Odd_name_1 = "O3"
	_Odd_name_2 = "O5"
	_Odd_name_3 = "O7"
	_Odd_name_4 = "O9"
	_Odd_name_5 = "O11"
	_Odd_name_6 = "O13"
	_Odd_name_7 = "O15"
	_Odd_name_8 = "O17"
	_Odd_name_9 = "O19"
)

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Odd) String() string {
	switch {
	case v == 1:
		return _Odd_name_0
	case v == 3:
		return _Odd_name_1
	case v == 5:
		return _Odd_name_2
	case v == 7:
		return _Odd_name_3
	case v == 9:
		return _Odd_name_4
	case v == 11:
		return _Odd_name_5
	case v == 13:
		return _Odd_name_6
	case v == 15:
		return _Odd_name_7
	case v == 17:
		return _Odd_name_8
	case v == 19:
		return _Odd_name_9
	default:
		return "Odd(" + strconv.FormatInt(int64(v), 10) + ")"
	}
}

// ErrUnknownOdd is the error of a text that does not match a declared
// variant. The error of [ParseOdd] wraps it with the text.
var ErrUnknownOdd = errors.New("gap: unknown Odd value")

// ParseOdd returns the variant of Odd whose text is s. It is the inverse
// of [Odd.String] on the declared variants. It returns an error that wraps
// [ErrUnknownOdd] for a text that does not match a declared variant.
func ParseOdd(s string) (Odd, error) {
	switch s {
	case "O1":
		return O1, nil
	case "O3":
		return O3, nil
	case "O5":
		return O5, nil
	case "O7":
		return O7, nil
	case "O9":
		return O9, nil
	case "O11":
		return O11, nil
	case "O13":
		return O13, nil
	case "O15":
		return O15, nil
	case "O17":
		return O17, nil
	case "O19":
		return O19, nil
	}
	var zero Odd
	return zero, fmt.Errorf("%w %q", ErrUnknownOdd, s)
}

// OddValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func OddValues() []Odd {
	return []Odd{
		O1,
		O3,
		O5,
		O7,
		O9,
		O11,
		O13,
		O15,
		O17,
		O19,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Odd) IsValid() bool {
	switch v {
	case O1, O3, O5, O7, O9, O11, O13, O15, O17, O19:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownOdd] for a value that [Odd.IsValid] refuses, whose
// text [ParseOdd] does not parse.
func (v Odd) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownOdd, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseOdd] for a text that does not match a declared variant.
func (v *Odd) UnmarshalText(text []byte) error {
	parsed, err := ParseOdd(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

func _() {
	// An "invalid array index" compiler error means that the value of a
	// variant changed. Run dokimi again to write the table of the new values.
	var x [1]struct{}
	_ = x[Minus - -1]
	_ = x[Zero-0]
	_ = x[Plus-1]
}

const _Num_name = "MinusZeroPlus"

var _Num_index = [...]uint8{0, 5, 9, 13}

// String returns the text of v. A value outside the declared set renders
// as its underlying value.
func (v Num) String() string {
	idx := int(v) - -1
	if v < -1 || idx >= len(_Num_index)-1 {
		return "Num(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Num_name[_Num_index[idx]:_Num_index[idx+1]]
}

// ErrUnknownNum is the error of a text that does not match a declared
// variant. The error of [ParseNum] wraps it with the text.
var ErrUnknownNum = errors.New("gap: unknown Num value")

// ParseNum returns the variant of Num whose text is s. It is the inverse
// of [Num.String] on the declared variants. It returns an error that wraps
// [ErrUnknownNum] for a text that does not match a declared variant.
func ParseNum(s string) (Num, error) {
	switch s {
	case "Minus":
		return Minus, nil
	case "Zero":
		return Zero, nil
	case "Plus":
		return Plus, nil
	}
	var zero Num
	return zero, fmt.Errorf("%w %q", ErrUnknownNum, s)
}

// NumValues returns the declared variants in the order of their
// declaration. Each call returns a new slice.
func NumValues() []Num {
	return []Num{
		Minus,
		Zero,
		Plus,
	}
}

// IsValid reports whether v is a declared variant. A conversion from the
// underlying type can make any other value.
func (v Num) IsValid() bool {
	switch v {
	case Minus, Zero, Plus:
		return true
	default:
		return false
	}
}

// MarshalText implements [encoding.TextMarshaler] with the text of v.
// encoding/json uses it for a value and for a map key. It returns an error
// that wraps [ErrUnknownNum] for a value that [Num.IsValid] refuses, whose
// text [ParseNum] does not parse.
func (v Num) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrUnknownNum, v.String())
	}
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It returns the
// error of [ParseNum] for a text that does not match a declared variant.
func (v *Num) UnmarshalText(text []byte) error {
	parsed, err := ParseNum(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
