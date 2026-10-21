package gap

import (
	"errors"
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

// String returns Gap's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
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

// ErrUnknownGap reports text that names no declared variant.
var ErrUnknownGap = errors.New("gap: unknown Gap value")

// ParseGap maps text back to a Gap.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
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
	return zero, ErrUnknownGap
}

// GapValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func GapValues() []Gap {
	return []Gap{
		Two,
		Three,
		Five,
		Six,
		Eleven,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Gap) IsValid() bool {
	switch v {
	case Two:
		return true
	case Three:
		return true
	case Five:
		return true
	case Six:
		return true
	case Eleven:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Gap legal as a map key —
// which a JSON marshaller alone does not.
func (v Gap) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
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

var (
	_Unum_index_0 = [...]uint8{0, 6, 9}
)

// String returns Unum's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
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

// ErrUnknownUnum reports text that names no declared variant.
var ErrUnknownUnum = errors.New("gap: unknown Unum value")

// ParseUnum maps text back to a Unum.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
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
	return zero, ErrUnknownUnum
}

// UnumValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func UnumValues() []Unum {
	return []Unum{
		Naught,
		One,
		Far,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Unum) IsValid() bool {
	switch v {
	case Naught:
		return true
	case One:
		return true
	case Far:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Unum legal as a map key —
// which a JSON marshaller alone does not.
func (v Unum) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
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

// String returns Prime's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Prime) String() string {
	if str, ok := _Prime_map[v]; ok {
		return str
	}
	return "Prime(" + strconv.FormatInt(int64(v), 10) + ")"
}

// ErrUnknownPrime reports text that names no declared variant.
var ErrUnknownPrime = errors.New("gap: unknown Prime value")

// ParsePrime maps text back to a Prime.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
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
	return zero, ErrUnknownPrime
}

// PrimeValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
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

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Prime) IsValid() bool {
	switch v {
	case P2:
		return true
	case P3:
		return true
	case P5:
		return true
	case P7:
		return true
	case P11:
		return true
	case P13:
		return true
	case P17:
		return true
	case P19:
		return true
	case P23:
		return true
	case P29:
		return true
	case P31:
		return true
	case P37:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Prime legal as a map key —
// which a JSON marshaller alone does not.
func (v Prime) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
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

// String returns Rank's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Rank) String() string {
	idx := int(v) - 1
	if v < 1 || idx >= len(_Rank_index)-1 {
		return "Rank(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Rank_name[_Rank_index[idx]:_Rank_index[idx+1]]
}

// ErrUnknownRank reports text that names no declared variant.
var ErrUnknownRank = errors.New("gap: unknown Rank value")

// ParseRank maps text back to a Rank.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
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
	return zero, ErrUnknownRank
}

// RankValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func RankValues() []Rank {
	return []Rank{
		Gold,
		Silver,
		Bronze,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Rank) IsValid() bool {
	switch v {
	case Gold:
		return true
	case Silver:
		return true
	case Bronze:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Rank legal as a map key —
// which a JSON marshaller alone does not.
func (v Rank) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
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

// String returns Odd's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
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

// ErrUnknownOdd reports text that names no declared variant.
var ErrUnknownOdd = errors.New("gap: unknown Odd value")

// ParseOdd maps text back to a Odd.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
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
	return zero, ErrUnknownOdd
}

// OddValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
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

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Odd) IsValid() bool {
	switch v {
	case O1:
		return true
	case O3:
		return true
	case O5:
		return true
	case O7:
		return true
	case O9:
		return true
	case O11:
		return true
	case O13:
		return true
	case O15:
		return true
	case O17:
		return true
	case O19:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Odd legal as a map key —
// which a JSON marshaller alone does not.
func (v Odd) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
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

// String returns Num's textual form.
//
// A value outside the declared set renders as itself, so a corrupt value
// is visible rather than indistinguishable from a declared one.
func (v Num) String() string {
	idx := int(v) - -1
	if v < -1 || idx >= len(_Num_index)-1 {
		return "Num(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Num_name[_Num_index[idx]:_Num_index[idx+1]]
}

// ErrUnknownNum reports text that names no declared variant.
var ErrUnknownNum = errors.New("gap: unknown Num value")

// ParseNum maps text back to a Num.
//
// The inverse of String for every declared variant, which is what lets a
// value survive a round trip through anything that carries it as text — a
// log line, a JSON field, a database column.
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
	return zero, ErrUnknownNum
}

// NumValues returns every declared variant in declaration order.
//
// A fresh slice per call: a package-level one would let any caller reorder
// or overwrite the set every other caller reads.
func NumValues() []Num {
	return []Num{
		Minus,
		Zero,
		Plus,
	}
}

// IsValid reports whether v is a declared variant.
//
// A conversion admits any value of the underlying type, so this is the
// only thing standing between an out-of-range value and the rest of the
// program.
func (v Num) IsValid() bool {
	switch v {
	case Minus:
		return true
	case Zero:
		return true
	case Plus:
		return true
	}
	return false
}

// MarshalText implements [encoding.TextMarshaler].
//
// Text rather than JSON: encoding/json reaches for TextMarshaler on its
// own and so does YAML, and it is what makes Num legal as a map key —
// which a JSON marshaller alone does not.
func (v Num) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], refusing text that
// names no declared variant rather than admitting it as an unnamed value.
func (v *Num) UnmarshalText(text []byte) error {
	parsed, err := ParseNum(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
