package commented_test

import (
	"testing"

	"example.com/shop/commented"
	"go.dokimi.dev/assert"
)

// TestGrade checks the functions and the methods of the enumeration Grade.
func TestGrade(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[commented.Grade]bool{}
		for _, v := range []commented.Grade{
			commented.Poor,
			commented.Fair,
			commented.Good,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("GradeValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []commented.Grade{
					commented.Poor,
					commented.Fair,
					commented.Good,
				}
				assert.Equal(
					t,
					commented.GradeValues(),
					want,
					"GradeValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := commented.GradeValues(), commented.GradeValues()
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

		tests := []struct {
			name string
			give commented.Grade
			want string
		}{
			{name: "returns the declared text of Poor", give: commented.Poor, want: "poor"},
			{name: "returns the declared text of Fair", give: commented.Fair, want: "fair enough"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(
					t,
					tt.give.String(),
					tt.want,
					"String returns the declared text of the variant",
				)
			})
		}

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []commented.Grade{
				commented.Poor,
				commented.Fair,
				commented.Good,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := commented.Grade(3)
				for _, v := range []commented.Grade{
					commented.Poor,
					commented.Fair,
					commented.Good,
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
			give commented.Grade
		}{
			{name: "accepts Poor", give: commented.Poor},
			{name: "accepts Fair", give: commented.Fair},
			{name: "accepts Good", give: commented.Good},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, commented.Grade(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Poor", func(t *testing.T) {
			t.Parallel()
			var zero commented.Grade
			assert.Equal(t, zero, commented.Poor, "the zero value is Poor")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseGrade", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give commented.Grade
		}{
			{name: "returns Poor for its text", give: commented.Poor},
			{name: "returns Fair for its text", give: commented.Fair},
			{name: "returns Good for its text", give: commented.Good},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := commented.ParseGrade(tt.give.String())
				assert.NoError(t, err, "ParseGrade accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseGrade returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := commented.ParseGrade("????????????")
			assert.ErrorIs(
				t,
				err,
				commented.ErrUnknownGrade,
				"ParseGrade refuses the text with ErrUnknownGrade",
			)
			assert.Contains(t, err.Error(), "????????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give commented.Grade
		}{
			{name: "encodes Poor as its text", give: commented.Poor},
			{name: "encodes Fair as its text", give: commented.Fair},
			{name: "encodes Good as its text", give: commented.Good},
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
			out := commented.Grade(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				commented.ErrUnknownGrade,
				"MarshalText refuses the value with ErrUnknownGrade",
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
			give commented.Grade
		}{
			{name: "decodes Poor from its encoding", give: commented.Poor},
			{name: "decodes Fair from its encoding", give: commented.Fair},
			{name: "decodes Good from its encoding", give: commented.Good},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got commented.Grade
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got commented.Grade
			assert.HasError(
				t,
				got.UnmarshalText([]byte("????????????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestMood checks the functions and the methods of the enumeration Mood.
func TestMood(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[commented.Mood]bool{}
		for _, v := range []commented.Mood{
			commented.Calm,
			commented.Angry,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("MoodValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []commented.Mood{
					commented.Calm,
					commented.Angry,
				}
				assert.Equal(
					t,
					commented.MoodValues(),
					want,
					"MoodValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := commented.MoodValues(), commented.MoodValues()
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
			for _, v := range []commented.Mood{
				commented.Calm,
				commented.Angry,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := commented.Mood(2)
				for _, v := range []commented.Mood{
					commented.Calm,
					commented.Angry,
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
			give commented.Mood
		}{
			{name: "accepts Calm", give: commented.Calm},
			{name: "accepts Angry", give: commented.Angry},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, commented.Mood(2).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Calm", func(t *testing.T) {
			t.Parallel()
			var zero commented.Mood
			assert.Equal(t, zero, commented.Calm, "the zero value is Calm")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseMood", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give commented.Mood
		}{
			{name: "returns Calm for its text", give: commented.Calm},
			{name: "returns Angry for its text", give: commented.Angry},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := commented.ParseMood(tt.give.String())
				assert.NoError(t, err, "ParseMood accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseMood returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := commented.ParseMood("??????")
			assert.ErrorIs(
				t,
				err,
				commented.ErrUnknownMood,
				"ParseMood refuses the text with ErrUnknownMood",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give commented.Mood
		}{
			{name: "encodes Calm as its text", give: commented.Calm},
			{name: "encodes Angry as its text", give: commented.Angry},
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
			out := commented.Mood(2)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				commented.ErrUnknownMood,
				"MarshalText refuses the value with ErrUnknownMood",
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
			give commented.Mood
		}{
			{name: "decodes Calm from its encoding", give: commented.Calm},
			{name: "decodes Angry from its encoding", give: commented.Angry},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got commented.Mood
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got commented.Mood
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
