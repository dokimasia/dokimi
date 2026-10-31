package status_test

import (
	"testing"

	"example.com/shop/status"
	"go.dokimi.dev/assert"
)

// TestStatus checks the functions and the methods of the enumeration
// Status.
func TestStatus(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[status.Status]bool{}
		for _, v := range []status.Status{
			status.Draft,
			status.Published,
			status.Archived,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("StatusValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []status.Status{
					status.Draft,
					status.Published,
					status.Archived,
				}
				assert.Equal(
					t,
					status.StatusValues(),
					want,
					"StatusValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := status.StatusValues(), status.StatusValues()
			assert.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []status.Status{
				status.Draft,
				status.Published,
				status.Archived,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := status.Status(4)
				for _, v := range []status.Status{
					status.Draft,
					status.Published,
					status.Archived,
				} {
					assert.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give status.Status
		}{
			{name: "accepts Draft", give: status.Draft},
			{name: "accepts Published", give: status.Published},
			{name: "accepts Archived", give: status.Archived},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, status.Status(4).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero status.Status
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseStatus", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give status.Status
		}{
			{name: "returns Draft for its text", give: status.Draft},
			{name: "returns Published for its text", give: status.Published},
			{name: "returns Archived for its text", give: status.Archived},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := status.ParseStatus(tt.give.String())
				assert.NoError(t, err, "ParseStatus accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseStatus returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := status.ParseStatus("??????????")
			assert.ErrorIs(
				t,
				err,
				status.ErrUnknownStatus,
				"ParseStatus refuses the text with ErrUnknownStatus",
			)
			assert.Contains(t, err.Error(), "??????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give status.Status
		}{
			{name: "encodes Draft as its text", give: status.Draft},
			{name: "encodes Published as its text", give: status.Published},
			{name: "encodes Archived as its text", give: status.Archived},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				assert.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			out := status.Status(4)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				status.ErrUnknownStatus,
				"MarshalText refuses the value with ErrUnknownStatus",
			)
			assert.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give status.Status
		}{
			{name: "decodes Draft from its encoding", give: status.Draft},
			{name: "decodes Published from its encoding", give: status.Published},
			{name: "decodes Archived from its encoding", give: status.Archived},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got status.Status
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got status.Status
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
