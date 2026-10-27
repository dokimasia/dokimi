// Package tagged hashes the leaves and the nodes of a tree under tagged roles.
package tagged

import "context"

// Role is the role of a hash. A leaf takes a role below 0x80, and a node takes
// a role of 0x80 or above.
type Role uint8

// leafRole is the role of a leaf, which code outside the package cannot refer
// to.
const leafRole Role = 0x01

// Size is the size of a truncated digest, with an authored sample.
//
//dokimi:sample value=16
type Size int

// Hasher hashes the leaves and the nodes of a tree.
//
//dokimi:suite
type Hasher interface {
	// HashTagged hashes the data of a leaf under the role r, which is below
	// 0x80.
	HashTagged(
		//dokimi:sample value=0x01
		r Role,
		data []byte,
	) [32]byte
	// CombineTagged hashes two digests under the role r, which is 0x80 or
	// above.
	CombineTagged(
		//dokimi:sample value=0x81
		r Role,
		left, right [32]byte,
	) [32]byte
	// HashRoles hashes the roles, each of which is below 0x80.
	HashRoles(
		//dokimi:sample value=0x02
		roles ...Role,
	) [32]byte
	// Truncate truncates the digest to the size n.
	Truncate(d [32]byte, n Size) []byte
	// Seal seals the digest.
	Seal(
		//dokimi:sample value=nil
		ctx context.Context,
		d [32]byte,
	) error
}

// Digest hashes the data under the role r, which is below 0x80.
//
//dokimi:suite
func Digest(
	//dokimi:sample value=0x03
	r Role,
	data []byte,
) [32]byte {
	return [32]byte{byte(r), byte(len(data))}
}

// Leaves hashes the leaves of a tree in the kit taggedtest.
//
//dokimi:suite
//dokimi:out path=taggedtest/
type Leaves interface {
	// Hash hashes the data of a leaf under the role r.
	Hash(
		//dokimi:sample value=leafRole
		r Role,
		data []byte,
	) [32]byte
}

// Leaf hashes the data of a leaf under the role r in the kit taggedtest.
//
//dokimi:suite
//dokimi:out path=taggedtest/
func Leaf(
	//dokimi:sample value=leafRole
	r Role,
	data []byte,
) [32]byte {
	return [32]byte{byte(r), byte(len(data))}
}
