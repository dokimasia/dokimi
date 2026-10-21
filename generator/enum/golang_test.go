// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"fmt"
	"strings"
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/files"
	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The packages of the Go tree that this spec alone reads.
const (
	gapPackage    = "example.com/shop/gap"
	remotePackage = "example.com/shop/remote"
	shapedPackage = "example.com/shop/shaped"
)

// The part of Go declares the target golang with the template of each
// family, and names the members of the surface of an enumeration. It prints a
// value outside the variants through the type under the enumeration, and
// quotes the text of each variant. It returns the table of the texts of an
// enumeration of integers, which the render slices. The checks refer to the
// functions and the variable of the surface through its package.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of each family for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, enum.New(), golang.Target, enumTemplate, testTemplate)
	})

	t.Run("writes the surfaces and the checks of the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t), golden.ShouldUpdate())
	})

	t.Run("writes Go that vets and whose tests pass", func(t *testing.T) {
		t.Parallel()

		a := gotesting.New()
		if !toolchain.Require(t, a) {
			return
		}
		project := gentest.Project(rendered(t), goTree, "example.com/shop")
		gotesting.AssertVets(t.Context(), t, project)
		toolchain.AssertTestsPass(t.Context(), t, a, project)
	})

	t.Run("names the methods, the functions and the error of the surface", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage).models[statusName].Names
		assert.Equal(t, got, enum.Names{
			Render: "String", Parse: "ParseStatus", Values: "StatusValues", Valid: "IsValid",
			Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
		}, "the surface of Status has the names of Go")
	})

	fallbacks := []struct {
		name       string
		pkg        string
		give       string
		wantType   string
		wantFormat string
	}{
		{
			name: "prints a number through its type with the verb d", pkg: statusPackage, give: statusName,
			wantType: "int", wantFormat: "%d",
		},
		{
			name: "prints a float through its type with the verb g", pkg: quietPackage, give: levelName,
			wantType: "float64", wantFormat: "%g",
		},
		{
			name: "prints a truth value through its type with the verb t", pkg: quietPackage, give: "Flag",
			wantType: "bool", wantFormat: "%t",
		},
	}
	for _, tt := range fallbacks {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, tt.pkg).models[tt.give]
			assert.NotNil(t, got.Fallback, tt.give+" converts a value outside the variants")
			expect.Equal(t, got.Fallback.Spelling, tt.wantType, "the value converts to the type under "+tt.give)
			expect.Equal(t, got.Format, tt.wantFormat, "the verb prints a value of that type")
		})
	}

	{
		tests := []struct {
			name    string
			givePkg string
			give    string
			want    *enum.Table
		}{
			{
				name: "records the table of an enumeration of integers in one run", givePkg: statusPackage,
				give: statusName,
				want: &enum.Table{
					Values: []enum.Entry{
						{Name: "Draft", Value: "1"},
						{Name: "Published", Value: "2"},
						{Name: "Archived", Value: "3"},
					},
					Runs: []enum.Run{{
						Name: `"DraftPublishedArchived"`, Low: "1", High: "3", Index: []int{0, 5, 14, 22}, Bits: 8,
					}},
					Indexed: true,
				},
			},
			{
				name: "records a run from the zero of a signed type as closed", givePkg: priorityPackage,
				give: priorityName,
				want: &enum.Table{
					Values: []enum.Entry{
						{Name: "Low", Value: "0"},
						{Name: "Medium", Value: "1"},
						{Name: "High", Value: "2"},
						{Name: "Critical", Value: "3"},
					},
					Runs: []enum.Run{{
						Name: `"LowMediumHighCritical"`, Low: "0", High: "3", Index: []int{0, 3, 9, 13, 21}, Bits: 8,
					}},
					Indexed: true,
				},
			},
			{
				name: "records the run of values that the declaration lists in decreasing order", givePkg: gapPackage,
				give: "Rank",
				want: &enum.Table{
					Values: []enum.Entry{
						{Name: "Gold", Value: "3"},
						{Name: "Silver", Value: "2"},
						{Name: "Bronze", Value: "1"},
					},
					Runs: []enum.Run{{
						Name: `"BronzeSilverGold"`, Low: "1", High: "3", Index: []int{0, 6, 12, 16}, Bits: 8,
					}},
					Indexed: true,
				},
			},
			{
				name: "records a run that starts below zero", givePkg: gapPackage, give: "Num",
				want: &enum.Table{
					Values: []enum.Entry{
						{Name: "Minus", Value: "-1"},
						{Name: "Zero", Value: "0"},
						{Name: "Plus", Value: "1"},
					},
					Runs: []enum.Run{{
						Name: `"MinusZeroPlus"`, Low: "-1", High: "1", Index: []int{0, 5, 9, 13}, Bits: 8,
					}},
					Indexed: true,
				},
			},
			{
				name: "records each run of consecutive values in increasing order", givePkg: gapPackage, give: "Gap",
				want: &enum.Table{
					Values: []enum.Entry{
						{Name: "Two", Value: "2"},
						{Name: "Three", Value: "3"},
						{Name: "Five", Value: "5"},
						{Name: "Six", Value: "6"},
						{Name: "Eleven", Value: "11"},
					},
					Runs: []enum.Run{
						{Name: `"TwoThree"`, Low: "2", High: "3", Index: []int{0, 3, 8}, Bits: 8},
						{Name: `"FiveSix"`, Low: "5", High: "6", Index: []int{0, 4, 7}, Bits: 8},
						{Name: `"Eleven"`, Low: "11", High: "11", Index: []int{0, 6}, Bits: 8, Single: true},
					},
					Indexed: true,
				},
			},
			{
				name: "records a run from the zero of an unsigned type as open", givePkg: gapPackage, give: "Unum",
				want: &enum.Table{
					Values: []enum.Entry{
						{Name: "Naught", Value: "0"},
						{Name: "One", Value: "1"},
						{Name: "Far", Value: "253"},
					},
					Runs: []enum.Run{
						{Name: `"NaughtOne"`, Low: "0", High: "1", Index: []int{0, 6, 9}, Bits: 8, Open: true},
						{Name: `"Far"`, Low: "253", High: "253", Index: []int{0, 3}, Bits: 8, Single: true},
					},
					Indexed: true,
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, golang.Target, goTree, tt.givePkg).models[tt.give]
				assert.Equal(t, got.Table, tt.want, "the render of "+tt.give+" slices its texts out of the table")
			})
		}
	}

	t.Run("records no index and no map of ten runs of one value", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, gapPackage).models["Odd"].Table
		assert.NotNil(t, got, "Odd has a table")
		expect.Length(t, got.Runs, 10, "each odd number is a run of its own")
		expect.False(t, got.Map, "a switch takes ten runs")
		expect.False(t, got.Indexed, "a run of one value has no index")
	})

	t.Run("records a map of the texts of an enumeration with more runs than a switch takes", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, gapPackage).models["Prime"].Table
		assert.NotNil(t, got, "Prime has a table")
		expect.True(t, got.Map, "a map looks up the text of a value of eleven runs")
		expect.Length(t, got.Runs, 11, "2 and 3 make one run, and each other prime a run of its own")
		expect.Equal(t, got.Name, `"P2P3P5P7P11P13P17P19P23P29P31P37"`, "the map slices one string of the texts")
		expect.Equal(t, got.Entries, []enum.Entry{
			{Name: "P2", Value: "2", Start: 0, End: 2},
			{Name: "P3", Value: "3", Start: 2, End: 4},
			{Name: "P5", Value: "5", Start: 4, End: 6},
			{Name: "P7", Value: "7", Start: 6, End: 8},
			{Name: "P11", Value: "11", Start: 8, End: 11},
			{Name: "P13", Value: "13", Start: 11, End: 14},
			{Name: "P17", Value: "17", Start: 14, End: 17},
			{Name: "P19", Value: "19", Start: 17, End: 20},
			{Name: "P23", Value: "23", Start: 20, End: 23},
			{Name: "P29", Value: "29", Start: 23, End: 26},
			{Name: "P31", Value: "31", Start: 26, End: 29},
			{Name: "P37", Value: "37", Start: 29, End: 32},
		}, "each value maps to the bounds of its text")
	})

	{
		tests := []struct {
			name    string
			givePkg string
			give    string
		}{
			{name: "records no table of an enumeration of text", givePkg: regionPackage, give: regionName},
			{name: "records no table of an enumeration of floats", givePkg: quietPackage, give: levelName},
			{name: "records no table of an enumeration of truth values", givePkg: quietPackage, give: "Switch"},
			{
				name:    "records no table of an enumeration whose values read another package",
				givePkg: remotePackage, give: "Season",
			},
			{name: "records no table of a float whose values are integers", givePkg: shapedPackage, give: "Ratio"},
			{
				name:    "records no table of an enumeration over another defined type",
				givePkg: shapedPackage, give: "Tier",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, recorded := generate(t, golang.Target, goTree, tt.givePkg).models[tt.give]
				assert.True(t, recorded, "enum records the surface of "+tt.give)
				expect.Nil(t, got.Table, "the render of "+tt.give+" switches over its variants")
			})
		}
	}

	{
		tests := []struct {
			name     string
			giveSize int
			want     int
		}{
			{name: "records an index of 8 bits for texts below 256 bytes", giveSize: 255, want: 8},
			{name: "records an index of 16 bits for texts of 256 bytes or more", giveSize: 256, want: 16},
			{name: "records an index of 32 bits for texts of 65536 bytes or more", giveSize: 1 << 16, want: 32},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				source := fmt.Sprintf("package wide\n\n// Wide has a long text.\n//\n//dokimi:enum text=comment\n"+
					"type Wide int\n\n// The variants of Wide.\nconst (\n\tShort Wide = iota // a\n\tLong // %s\n)\n",
					strings.Repeat("x", tt.giveSize-1))
				dir := files.Workspace(t, files.Tree{
					"go.mod":       files.Text("module example.com/wide\n\ngo 1.27\n"),
					"wide/wide.go": files.Text(source),
				})
				got := generate(t, golang.Target, dir, "example.com/wide/wide").models["Wide"].Table
				assert.NotNil(t, got, "Wide has a table")
				assert.Length(t, got.Runs, 1, "the values of Wide are one run")
				expect.Equal(t, got.Runs[0].Bits, tt.want, "the index of texts of the size takes the bits of the case")
			})
		}
	}

	t.Run("prints a value of text outside the variants without a conversion", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, regionPackage).models[regionName]
		expect.Nil(t, got.Fallback, "Region renders a value outside the variants as its text")
	})

	t.Run("writes a variant as the constant of its package", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage).models[statusName]
		assert.NotEmpty(t, got.Variants, "Status has variants")
		expect.Equal(t, got.Variants[0].Member, "Draft", "a template writes Draft alone")
	})

	t.Run("declares the package members of the surface beside the enumeration", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage)
		assert.Equal(t, got.beside[statusName], []string{"ParseStatus", "StatusValues", "ErrUnknownStatus"},
			"the record of the surface declares the parse function, the list of the variants and the error")
	})

	t.Run("refers to each package member of the surface from the record of the checks", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage)
		assert.Equal(t, gentest.Referents(got.tests[statusName]),
			[]string{"ParseStatus", "StatusValues", "ErrUnknownStatus"},
			"the fields of the record of the checks refer to the declarations of the surface in their order")
	})

	t.Run("declares nothing beside an enumeration with its own parse function", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, handwrittenPackage)
		expect.Empty(t, got.beside[weekdayName], "the surface of Weekday has no parse function of its own")
	})

	t.Run("refers to nothing from the checks of an enumeration with its own parse function", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, handwrittenPackage)
		assert.NotNil(t, got.tests[weekdayName], "Weekday has checks")
		expect.Empty(t, gentest.Referents(got.tests[weekdayName]), "the checks of Weekday call no parse function")
	})
}

// rendered runs enum over the packages of the Go tree whose enumerations
// have surfaces, and returns the files that the Go backend renders from
// them.
func rendered(t *testing.T) fstest.MapFS {
	t.Helper()

	p := enum.New()
	res := fixture(t, golang.Target, goTree, p, statusPackage, "example.com/shop/priority", regionPackage,
		"example.com/shop/split", "example.com/shop/handwritten", quietPackage, commentedPackage, gapPackage,
		remotePackage, shapedPackage).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	return gentest.Rendered(t, res, p, gobackend.New())
}
