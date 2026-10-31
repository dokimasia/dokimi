package quiet_test

import (
	"testing"

	"example.com/shop/quiet"
	"go.dokimi.dev/assert"
)

// TestMode checks the functions and the methods of the enumeration Mode.
func TestMode(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Mode]bool{}
		for _, v := range []quiet.Mode{
			quiet.Off,
			quiet.On,
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
			for _, v := range []quiet.Mode{
				quiet.Off,
				quiet.On,
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
			give quiet.Mode
		}{
			{name: "accepts Off", give: quiet.Off},
			{name: "accepts On", give: quiet.On},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, quiet.Mode(2).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Off", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Mode
			assert.Equal(t, zero, quiet.Off, "the zero value is Off")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Mode
		}{
			{name: "encodes Off as its text", give: quiet.Off},
			{name: "encodes On as its text", give: quiet.On},
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
}

// TestPlain checks the functions and the methods of the enumeration Plain.
func TestPlain(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Plain]bool{}
		for _, v := range []quiet.Plain{
			quiet.Loose,
			quiet.Tight,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})
}

// TestGauge checks the functions and the methods of the enumeration Gauge.
func TestGauge(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Gauge]bool{}
		for _, v := range []quiet.Gauge{
			quiet.Narrow,
			quiet.Wide,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})
}

// TestTone checks the functions and the methods of the enumeration Tone.
func TestTone(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Tone]bool{}
		for _, v := range []quiet.Tone{
			quiet.Low,
			quiet.High,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("ToneValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []quiet.Tone{
					quiet.Low,
					quiet.High,
				}
				assert.Equal(
					t,
					quiet.ToneValues(),
					want,
					"ToneValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := quiet.ToneValues(), quiet.ToneValues()
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
			for _, v := range []quiet.Tone{
				quiet.Low,
				quiet.High,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := quiet.Tone(3)
				for _, v := range []quiet.Tone{
					quiet.Low,
					quiet.High,
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
			give quiet.Tone
		}{
			{name: "accepts Low", give: quiet.Low},
			{name: "accepts High", give: quiet.High},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, quiet.Tone(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Tone
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseTone", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Tone
		}{
			{name: "returns Low for its text", give: quiet.Low},
			{name: "returns High for its text", give: quiet.High},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := quiet.ParseTone(tt.give.String())
				assert.NoError(t, err, "ParseTone accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseTone returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := quiet.ParseTone("?????")
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownTone,
				"ParseTone refuses the text with ErrUnknownTone",
			)
			assert.Contains(t, err.Error(), "?????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Tone
		}{
			{name: "encodes Low as its text", give: quiet.Low},
			{name: "encodes High as its text", give: quiet.High},
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
}

// TestSwitch checks the functions and the methods of the enumeration
// Switch.
func TestSwitch(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Switch]bool{}
		for _, v := range []quiet.Switch{
			quiet.Up,
			quiet.Down,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("SwitchValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []quiet.Switch{
					quiet.Up,
					quiet.Down,
				}
				assert.Equal(
					t,
					quiet.SwitchValues(),
					want,
					"SwitchValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := quiet.SwitchValues(), quiet.SwitchValues()
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
			for _, v := range []quiet.Switch{
				quiet.Up,
				quiet.Down,
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
			give quiet.Switch
		}{
			{name: "accepts Up", give: quiet.Up},
			{name: "accepts Down", give: quiet.Down},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("accepts the zero value Down", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Switch
			assert.Equal(t, zero, quiet.Down, "the zero value is Down")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseSwitch", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Switch
		}{
			{name: "returns Up for its text", give: quiet.Up},
			{name: "returns Down for its text", give: quiet.Down},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := quiet.ParseSwitch(tt.give.String())
				assert.NoError(t, err, "ParseSwitch accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseSwitch returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := quiet.ParseSwitch("?????")
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownSwitch,
				"ParseSwitch refuses the text with ErrUnknownSwitch",
			)
			assert.Contains(t, err.Error(), "?????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Switch
		}{
			{name: "encodes Up as its text", give: quiet.Up},
			{name: "encodes Down as its text", give: quiet.Down},
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
			give quiet.Switch
		}{
			{name: "decodes Up from its encoding", give: quiet.Up},
			{name: "decodes Down from its encoding", give: quiet.Down},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got quiet.Switch
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got quiet.Switch
			assert.HasError(t, got.UnmarshalText([]byte("?????")), "UnmarshalText refuses the text")
		})
	})
}

// TestPhase checks the functions and the methods of the enumeration Phase.
func TestPhase(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Phase]bool{}
		for _, v := range []quiet.Phase{
			quiet.Warm,
			quiet.Cool,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("PhaseValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []quiet.Phase{
					quiet.Warm,
					quiet.Cool,
				}
				assert.Equal(
					t,
					quiet.PhaseValues(),
					want,
					"PhaseValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := quiet.PhaseValues(), quiet.PhaseValues()
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
			for _, v := range []quiet.Phase{
				quiet.Warm,
				quiet.Cool,
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
			give quiet.Phase
		}{
			{name: "accepts Warm", give: quiet.Warm},
			{name: "accepts Cool", give: quiet.Cool},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, quiet.Phase(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Phase
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})
}

// TestSpeed checks the functions and the methods of the enumeration Speed.
func TestSpeed(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Speed]bool{}
		for _, v := range []quiet.Speed{
			quiet.Slow,
			quiet.Fast,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("SpeedValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []quiet.Speed{
					quiet.Slow,
					quiet.Fast,
				}
				assert.Equal(
					t,
					quiet.SpeedValues(),
					want,
					"SpeedValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := quiet.SpeedValues(), quiet.SpeedValues()
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
			for _, v := range []quiet.Speed{
				quiet.Slow,
				quiet.Fast,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := quiet.Speed(3)
				for _, v := range []quiet.Speed{
					quiet.Slow,
					quiet.Fast,
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
			give quiet.Speed
		}{
			{name: "accepts Slow", give: quiet.Slow},
			{name: "accepts Fast", give: quiet.Fast},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, quiet.Speed(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Speed
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseSpeed", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Speed
		}{
			{name: "returns Slow for its text", give: quiet.Slow},
			{name: "returns Fast for its text", give: quiet.Fast},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := quiet.ParseSpeed(tt.give.String())
				assert.NoError(t, err, "ParseSpeed accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseSpeed returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := quiet.ParseSpeed("?????")
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownSpeed,
				"ParseSpeed refuses the text with ErrUnknownSpeed",
			)
			assert.Contains(t, err.Error(), "?????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Speed
		}{
			{name: "encodes Slow as its text", give: quiet.Slow},
			{name: "encodes Fast as its text", give: quiet.Fast},
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
			out := quiet.Speed(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownSpeed,
				"MarshalText refuses the value with ErrUnknownSpeed",
			)
			assert.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})
}

// TestLevel checks the functions and the methods of the enumeration Level.
func TestLevel(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Level]bool{}
		for _, v := range []quiet.Level{
			quiet.Dim,
			quiet.Bright,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("LevelValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []quiet.Level{
					quiet.Dim,
					quiet.Bright,
				}
				assert.Equal(
					t,
					quiet.LevelValues(),
					want,
					"LevelValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := quiet.LevelValues(), quiet.LevelValues()
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
			for _, v := range []quiet.Level{
				quiet.Dim,
				quiet.Bright,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := quiet.Level(2.5)
				for _, v := range []quiet.Level{
					quiet.Dim,
					quiet.Bright,
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
			give quiet.Level
		}{
			{name: "accepts Dim", give: quiet.Dim},
			{name: "accepts Bright", give: quiet.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, quiet.Level(2.5).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Level
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseLevel", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Level
		}{
			{name: "returns Dim for its text", give: quiet.Dim},
			{name: "returns Bright for its text", give: quiet.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := quiet.ParseLevel(tt.give.String())
				assert.NoError(t, err, "ParseLevel accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseLevel returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := quiet.ParseLevel("???????")
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownLevel,
				"ParseLevel refuses the text with ErrUnknownLevel",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Level
		}{
			{name: "encodes Dim as its text", give: quiet.Dim},
			{name: "encodes Bright as its text", give: quiet.Bright},
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
			out := quiet.Level(2.5)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownLevel,
				"MarshalText refuses the value with ErrUnknownLevel",
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
			give quiet.Level
		}{
			{name: "decodes Dim from its encoding", give: quiet.Dim},
			{name: "decodes Bright from its encoding", give: quiet.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got quiet.Level
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got quiet.Level
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestFlag checks the functions and the methods of the enumeration Flag.
func TestFlag(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[quiet.Flag]bool{}
		for _, v := range []quiet.Flag{
			quiet.Lit,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("FlagValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []quiet.Flag{
					quiet.Lit,
				}
				assert.Equal(
					t,
					quiet.FlagValues(),
					want,
					"FlagValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := quiet.FlagValues(), quiet.FlagValues()
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
			for _, v := range []quiet.Flag{
				quiet.Lit,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := quiet.Flag(false)
				for _, v := range []quiet.Flag{
					quiet.Lit,
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
			give quiet.Flag
		}{
			{name: "accepts Lit", give: quiet.Lit},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, quiet.Flag(false).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero quiet.Flag
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseFlag", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Flag
		}{
			{name: "returns Lit for its text", give: quiet.Lit},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := quiet.ParseFlag(tt.give.String())
				assert.NoError(t, err, "ParseFlag accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseFlag returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := quiet.ParseFlag("????")
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownFlag,
				"ParseFlag refuses the text with ErrUnknownFlag",
			)
			assert.Contains(t, err.Error(), "????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give quiet.Flag
		}{
			{name: "encodes Lit as its text", give: quiet.Lit},
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
			out := quiet.Flag(false)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				quiet.ErrUnknownFlag,
				"MarshalText refuses the value with ErrUnknownFlag",
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
			give quiet.Flag
		}{
			{name: "decodes Lit from its encoding", give: quiet.Lit},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got quiet.Flag
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got quiet.Flag
			assert.HasError(t, got.UnmarshalText([]byte("????")), "UnmarshalText refuses the text")
		})
	})
}
