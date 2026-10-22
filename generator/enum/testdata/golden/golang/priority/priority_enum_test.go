package priority

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestPriority checks the functions and the methods of the enumeration
// Priority.
func TestPriority(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Priority]bool{}
		for _, v := range []Priority{
			Low,
			Medium,
			High,
			Critical,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("PriorityValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []Priority{
					Low,
					Medium,
					High,
					Critical,
				}
				assert.Equal(
					t,
					PriorityValues(),
					want,
					"PriorityValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := PriorityValues(), PriorityValues()
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
			for _, v := range []Priority{
				Low,
				Medium,
				High,
				Critical,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Priority(4)
				for _, v := range []Priority{
					Low,
					Medium,
					High,
					Critical,
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
			give Priority
		}{
			{name: "accepts Low", give: Low},
			{name: "accepts Medium", give: Medium},
			{name: "accepts High", give: High},
			{name: "accepts Critical", give: Critical},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Priority(4).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Low", func(t *testing.T) {
			t.Parallel()
			var zero Priority
			assert.Equal(t, zero, Low, "the zero value is Low")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParsePriority", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Priority
		}{
			{name: "returns Low for its text", give: Low},
			{name: "returns Medium for its text", give: Medium},
			{name: "returns High for its text", give: High},
			{name: "returns Critical for its text", give: Critical},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParsePriority(tt.give.String())
				assert.NoError(t, err, "ParsePriority accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParsePriority returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParsePriority("?????????")
			assert.ErrorIs(
				t,
				err,
				ErrUnknownPriority,
				"ParsePriority refuses the text with ErrUnknownPriority",
			)
			assert.Contains(t, err.Error(), "?????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Priority
		}{
			{name: "encodes Low as its text", give: Low},
			{name: "encodes Medium as its text", give: Medium},
			{name: "encodes High as its text", give: High},
			{name: "encodes Critical as its text", give: Critical},
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
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Priority
		}{
			{name: "decodes Low from its encoding", give: Low},
			{name: "decodes Medium from its encoding", give: Medium},
			{name: "decodes High from its encoding", give: High},
			{name: "decodes Critical from its encoding", give: Critical},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Priority
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Priority
			assert.HasError(
				t,
				got.UnmarshalText([]byte("?????????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
