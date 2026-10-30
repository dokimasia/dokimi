// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus_test

import (
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"

	corpus "go.dokimi.dev/dokimi/corpus/go"
)

func TestFamilies(t *testing.T) {
	t.Parallel()

	t.Run("Families", func(t *testing.T) {
		t.Parallel()

		t.Run("contain each subject of the corpus once and have a control each", func(t *testing.T) {
			t.Parallel()

			subjects, err := corpus.Subjects(os.DirFS("."))
			assert.NoError(t, err, "the corpus reads")
			assert.NotEmpty(t, subjects, "the corpus declares subjects")
			assert.Empty(t, corpus.Problems(corpus.Families, subjects), "the families fit the corpus")
		})
	})

	t.Run("Subjects", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each type after the directive subject", func(t *testing.T) {
			t.Parallel()

			fsys := fstest.MapFS{
				"a/atest/inmemory.go": {Data: []byte("package atest\n\n// InMemory is a subject.\n//\n" +
					"//dokimi:subject a.Store new=New\ntype InMemory struct{}\n\ntype helper struct{}\n")},
				"b/control.go": {
					Data: []byte("package b\n\n//dokimi:subject a.Store new=New\ntype Map map[string]int\n"),
				},
				"b/types.go/notes": {Data: []byte("a directory whose name ends in .go\n")},
				"b/notes.txt":      {Data: []byte("//dokimi:subject a.Store new=New\ntype Text struct{}\n")},
				"c/generic.go": {
					Data: []byte("package c\n\n//dokimi:subject a.Store new=New\ntype Store[K comparable] struct{}\n"),
				},
			}

			got, err := corpus.Subjects(fsys)

			assert.NoError(t, err, "the files read")
			assert.Equal(t, got, []string{"a/atest.InMemory", "b.Map", "c.Store"},
				"the subjects are the types after the directive")
		})

		t.Run("returns the error of a directory that does not exist", func(t *testing.T) {
			t.Parallel()

			_, err := corpus.Subjects(os.DirFS("absent"))

			assert.HasError(t, err, "the walk fails")
		})

		t.Run("returns the error of a file that does not read", func(t *testing.T) {
			t.Parallel()

			_, err := corpus.Subjects(fstest.MapFS{"a/link.go": {Mode: fs.ModeSymlink, Data: []byte("absent.go")}})

			assert.HasError(t, err, "the read fails")
		})
	})

	t.Run("Problems", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name         string
			giveFamilies []corpus.Family
			giveSubjects []string
			want         []string
		}{
			{
				name: "returns nothing for families that fit the subjects",
				giveFamilies: []corpus.Family{
					{Name: "stores", Subjects: []string{"a.Store"}, Controls: []string{"b.Store"}},
				},
				giveSubjects: []string{"a.Store", "b.Store"},
			},
			{
				name:         "returns a family without a control",
				giveFamilies: []corpus.Family{{Name: "stores", Subjects: []string{"a.Store"}}},
				giveSubjects: []string{"a.Store"},
				want:         []string{"the family stores has no control"},
			},
			{
				name: "returns a subject of a family that the corpus does not declare",
				giveFamilies: []corpus.Family{
					{Name: "stores", Subjects: []string{"a.Store", "c.Store"}, Controls: []string{"b.Store"}},
				},
				giveSubjects: []string{"a.Store", "b.Store"},
				want:         []string{"the family stores lists c.Store, which the corpus does not declare"},
			},
			{
				name: "returns a subject in no family",
				giveFamilies: []corpus.Family{
					{Name: "stores", Subjects: []string{"a.Store"}, Controls: []string{"b.Store"}},
				},
				giveSubjects: []string{"a.Store", "b.Store", "c.Log"},
				want:         []string{"c.Log is in no family"},
			},
			{
				name: "returns a subject in two families",
				giveFamilies: []corpus.Family{
					{Name: "stores", Subjects: []string{"a.Store"}, Controls: []string{"b.Store"}},
					{Name: "logs", Subjects: []string{"a.Store"}, Controls: []string{"c.Log"}},
				},
				giveSubjects: []string{"a.Store", "b.Store", "c.Log"},
				want:         []string{"a.Store is in the families stores and logs"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, corpus.Problems(tt.giveFamilies, tt.giveSubjects), tt.want,
					"the problems are the misfits of the families")
			})
		}
	})
}
