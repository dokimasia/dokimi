package commented

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestGrade checks the functions and the methods of the enumeration Grade.
func TestGrade(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Grade]bool{}
		for _, v := range []Grade{
			Poor,
			Fair,
			Good,
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
				want := []Grade{
					Poor,
					Fair,
					Good,
				}
				assert.Equal(t, GradeValues(), want, "GradeValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := GradeValues(), GradeValues()
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
			give Grade
			want string
		}{
			{name: "returns the declared text of Poor", give: Poor, want: "poor"},
			{name: "returns the declared text of Fair", give: Fair, want: "fair enough"},
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
			for _, v := range []Grade{
				Poor,
				Fair,
				Good,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Grade(3)
				for _, v := range []Grade{
					Poor,
					Fair,
					Good,
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
			give Grade
		}{
			{name: "accepts Poor", give: Poor},
			{name: "accepts Fair", give: Fair},
			{name: "accepts Good", give: Good},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Grade(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Poor", func(t *testing.T) {
			t.Parallel()
			var zero Grade
			assert.Equal(t, zero, Poor, "the zero value is Poor")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseGrade", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Grade
		}{
			{name: "returns Poor for its text", give: Poor},
			{name: "returns Fair for its text", give: Fair},
			{name: "returns Good for its text", give: Good},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseGrade(tt.give.String())
				assert.NoError(t, err, "ParseGrade accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseGrade returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseGrade("????????????")
			assert.ErrorIs(
				t,
				err,
				ErrUnknownGrade,
				"ParseGrade refuses the text with ErrUnknownGrade",
			)
			assert.Contains(t, err.Error(), "????????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Grade
		}{
			{name: "encodes Poor as its text", give: Poor},
			{name: "encodes Fair as its text", give: Fair},
			{name: "encodes Good as its text", give: Good},
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
			out := Grade(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownGrade,
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
			give Grade
		}{
			{name: "decodes Poor from its encoding", give: Poor},
			{name: "decodes Fair from its encoding", give: Fair},
			{name: "decodes Good from its encoding", give: Good},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Grade
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Grade
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
		seen := map[Mood]bool{}
		for _, v := range []Mood{
			Calm,
			Angry,
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
				want := []Mood{
					Calm,
					Angry,
				}
				assert.Equal(t, MoodValues(), want, "MoodValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := MoodValues(), MoodValues()
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
			for _, v := range []Mood{
				Calm,
				Angry,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Mood(2)
				for _, v := range []Mood{
					Calm,
					Angry,
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
			give Mood
		}{
			{name: "accepts Calm", give: Calm},
			{name: "accepts Angry", give: Angry},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Mood(2).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Calm", func(t *testing.T) {
			t.Parallel()
			var zero Mood
			assert.Equal(t, zero, Calm, "the zero value is Calm")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseMood", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Mood
		}{
			{name: "returns Calm for its text", give: Calm},
			{name: "returns Angry for its text", give: Angry},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseMood(tt.give.String())
				assert.NoError(t, err, "ParseMood accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseMood returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseMood("??????")
			assert.ErrorIs(t, err, ErrUnknownMood, "ParseMood refuses the text with ErrUnknownMood")
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Mood
		}{
			{name: "encodes Calm as its text", give: Calm},
			{name: "encodes Angry as its text", give: Angry},
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
			out := Mood(2)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownMood,
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
			give Mood
		}{
			{name: "decodes Calm from its encoding", give: Calm},
			{name: "decodes Angry from its encoding", give: Angry},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Mood
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Mood
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
