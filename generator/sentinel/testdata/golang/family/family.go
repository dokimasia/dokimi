// Package family reads the policy of a procedure. Its error values form a
// family: each rule wraps ErrPolicy, and ErrAnonymous wraps the rule
// ErrAccess.
//
//dokimi:sentinel
package family

import (
	"errors"
	"fmt"
)

// ErrPolicy is returned for a method whose options break a rule.
var ErrPolicy = errors.New("family: invalid policy")

// The errors of the rules. Each of them wraps ErrPolicy.
var (
	// ErrExposure is returned for a method without an exposure.
	ErrExposure = fmt.Errorf("%w: no exposure", ErrPolicy)

	// ErrAccess is returned for a method without an access rule.
	ErrAccess = fmt.Errorf("%w: no access rule", ErrPolicy)
)

// ErrAnonymous is returned for an access rule without a principal.
var ErrAnonymous = fmt.Errorf("%w: no principal", ErrAccess)

// errQuota is the rule of the quotas, and no error value of the package.
var errQuota = fmt.Errorf("%w: quota", ErrPolicy)

// ErrOverQuota is returned for a method over its quota. It wraps ErrPolicy
// through errQuota.
var ErrOverQuota = fmt.Errorf("%w: over the limit", errQuota)

// ErrClosed is returned for a closed procedure, and wraps no error value.
var ErrClosed = errors.New("family: closed")
