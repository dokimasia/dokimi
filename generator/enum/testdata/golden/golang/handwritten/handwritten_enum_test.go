package handwritten

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestWeekday checks the functions and the methods of the enumeration
// Weekday.
func TestWeekday(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Weekday]bool{}
		for _, v := range []Weekday{
			Monday,
			Tuesday,
			Wednesday,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []Weekday{
				Monday,
				Tuesday,
				Wednesday,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Weekday
		}{
			{name: "accepts Monday", give: Monday},
			{name: "accepts Tuesday", give: Tuesday},
			{name: "accepts Wednesday", give: Wednesday},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Weekday(4).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Weekday
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Weekday
		}{
			{name: "encodes Monday as its text", give: Monday},
			{name: "encodes Tuesday as its text", give: Tuesday},
			{name: "encodes Wednesday as its text", give: Wednesday},
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
			give Weekday
		}{
			{name: "decodes Monday from its encoding", give: Monday},
			{name: "decodes Tuesday from its encoding", give: Tuesday},
			{name: "decodes Wednesday from its encoding", give: Wednesday},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Weekday
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Weekday
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
