package family_test

import (
	"strings"
	"testing"

	"example.com/shop/family"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestErrPolicy checks the sentinel error ErrPolicy. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrPolicy(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrExposure", family.ErrExposure},
		{"ErrAccess", family.ErrAccess},
		{"ErrAnonymous", family.ErrAnonymous},
		{"ErrOverQuota", family.ErrOverQuota},
		{"ErrClosed", family.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, family.ErrPolicy, "ErrPolicy is assigned")
			assert.NotEqual(t, family.ErrPolicy.Error(), "", "ErrPolicy has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrPolicy.Error(),
				"family: ",
				"the message of ErrPolicy opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					family.ErrPolicy.Error(),
					other.err.Error(),
					"ErrPolicy has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(family.ErrPolicy.Error(), other.err.Error()),
					"the message of ErrPolicy opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					family.ErrPolicy,
					other.err,
					"ErrPolicy does not match "+other.name,
				)
			}
		})
	})
}

// TestErrExposure checks the sentinel error ErrExposure. A caller matches
// it with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrExposure(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrAccess", family.ErrAccess},
		{"ErrAnonymous", family.ErrAnonymous},
		{"ErrOverQuota", family.ErrOverQuota},
		{"ErrClosed", family.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, family.ErrExposure, "ErrExposure is assigned")
			assert.NotEqual(t, family.ErrExposure.Error(), "", "ErrExposure has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrExposure.Error(),
				"family: ",
				"the message of ErrExposure opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			expect.NotEqual(
				t,
				family.ErrExposure.Error(),
				family.ErrPolicy.Error(),
				"ErrExposure has the message of ErrPolicy",
			)
			for _, other := range others {
				expect.NotEqual(
					t,
					family.ErrExposure.Error(),
					other.err.Error(),
					"ErrExposure has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of its parent", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrExposure.Error(),
				family.ErrPolicy.Error(),
				"the message of ErrExposure opens with the message of ErrPolicy",
			)
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(family.ErrExposure.Error(), other.err.Error()),
					"the message of ErrExposure opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches its parent", func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, family.ErrExposure, family.ErrPolicy, "ErrExposure matches ErrPolicy")
		})

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					family.ErrExposure,
					other.err,
					"ErrExposure does not match "+other.name,
				)
			}
		})
	})
}

// TestErrAccess checks the sentinel error ErrAccess. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrAccess(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrExposure", family.ErrExposure},
		{"ErrAnonymous", family.ErrAnonymous},
		{"ErrOverQuota", family.ErrOverQuota},
		{"ErrClosed", family.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, family.ErrAccess, "ErrAccess is assigned")
			assert.NotEqual(t, family.ErrAccess.Error(), "", "ErrAccess has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrAccess.Error(),
				"family: ",
				"the message of ErrAccess opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			expect.NotEqual(
				t,
				family.ErrAccess.Error(),
				family.ErrPolicy.Error(),
				"ErrAccess has the message of ErrPolicy",
			)
			for _, other := range others {
				expect.NotEqual(
					t,
					family.ErrAccess.Error(),
					other.err.Error(),
					"ErrAccess has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of its parent", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrAccess.Error(),
				family.ErrPolicy.Error(),
				"the message of ErrAccess opens with the message of ErrPolicy",
			)
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(family.ErrAccess.Error(), other.err.Error()),
					"the message of ErrAccess opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches its parent", func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, family.ErrAccess, family.ErrPolicy, "ErrAccess matches ErrPolicy")
		})

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					family.ErrAccess,
					other.err,
					"ErrAccess does not match "+other.name,
				)
			}
		})
	})
}

// TestErrAnonymous checks the sentinel error ErrAnonymous. A caller
// matches it with errors.Is and reads its message, so it has a message and
// an identity of its own.
func TestErrAnonymous(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrExposure", family.ErrExposure},
		{"ErrOverQuota", family.ErrOverQuota},
		{"ErrClosed", family.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, family.ErrAnonymous, "ErrAnonymous is assigned")
			assert.NotEqual(t, family.ErrAnonymous.Error(), "", "ErrAnonymous has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrAnonymous.Error(),
				"family: ",
				"the message of ErrAnonymous opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			expect.NotEqual(
				t,
				family.ErrAnonymous.Error(),
				family.ErrAccess.Error(),
				"ErrAnonymous has the message of ErrAccess",
			)
			for _, other := range others {
				expect.NotEqual(
					t,
					family.ErrAnonymous.Error(),
					other.err.Error(),
					"ErrAnonymous has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of its parent", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrAnonymous.Error(),
				family.ErrAccess.Error(),
				"the message of ErrAnonymous opens with the message of ErrAccess",
			)
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(family.ErrAnonymous.Error(), other.err.Error()),
					"the message of ErrAnonymous opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches its parent", func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(
				t,
				family.ErrAnonymous,
				family.ErrAccess,
				"ErrAnonymous matches ErrAccess",
			)
		})

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					family.ErrAnonymous,
					other.err,
					"ErrAnonymous does not match "+other.name,
				)
			}
		})
	})
}

// TestErrOverQuota checks the sentinel error ErrOverQuota. A caller
// matches it with errors.Is and reads its message, so it has a message and
// an identity of its own.
func TestErrOverQuota(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrExposure", family.ErrExposure},
		{"ErrAccess", family.ErrAccess},
		{"ErrAnonymous", family.ErrAnonymous},
		{"ErrClosed", family.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, family.ErrOverQuota, "ErrOverQuota is assigned")
			assert.NotEqual(t, family.ErrOverQuota.Error(), "", "ErrOverQuota has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrOverQuota.Error(),
				"family: ",
				"the message of ErrOverQuota opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			expect.NotEqual(
				t,
				family.ErrOverQuota.Error(),
				family.ErrPolicy.Error(),
				"ErrOverQuota has the message of ErrPolicy",
			)
			for _, other := range others {
				expect.NotEqual(
					t,
					family.ErrOverQuota.Error(),
					other.err.Error(),
					"ErrOverQuota has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of its parent", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrOverQuota.Error(),
				family.ErrPolicy.Error(),
				"the message of ErrOverQuota opens with the message of ErrPolicy",
			)
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(family.ErrOverQuota.Error(), other.err.Error()),
					"the message of ErrOverQuota opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches its parent", func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(
				t,
				family.ErrOverQuota,
				family.ErrPolicy,
				"ErrOverQuota matches ErrPolicy",
			)
		})

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					family.ErrOverQuota,
					other.err,
					"ErrOverQuota does not match "+other.name,
				)
			}
		})
	})
}

// TestErrClosed checks the sentinel error ErrClosed. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrClosed(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrPolicy", family.ErrPolicy},
		{"ErrExposure", family.ErrExposure},
		{"ErrAccess", family.ErrAccess},
		{"ErrAnonymous", family.ErrAnonymous},
		{"ErrOverQuota", family.ErrOverQuota},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, family.ErrClosed, "ErrClosed is assigned")
			assert.NotEqual(t, family.ErrClosed.Error(), "", "ErrClosed has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				family.ErrClosed.Error(),
				"family: ",
				"the message of ErrClosed opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					family.ErrClosed.Error(),
					other.err.Error(),
					"ErrClosed has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(family.ErrClosed.Error(), other.err.Error()),
					"the message of ErrClosed opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					family.ErrClosed,
					other.err,
					"ErrClosed does not match "+other.name,
				)
			}
		})
	})
}
