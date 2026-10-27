package durable

import "context"

// Journal is the Keeper over a medium in memory, whose function of the
// options supplies its doors.
//
//dokimi:subject Keeper new=NewJournal options=JournalOptions
type Journal struct {
	values map[string]string
}

// NewJournal returns an empty journal.
func NewJournal() *Journal { return &Journal{values: map[string]string{}} }

// Put stores the value under the key.
func (j *Journal) Put(_ context.Context, key, value string) error {
	j.values[key] = value
	return nil
}

// Get returns the value under the key.
func (j *Journal) Get(_ context.Context, key string) (string, error) {
	return j.values[key], nil
}

// JournalOptions returns the doors of a journal.
func JournalOptions() []KeeperOption { return nil }

// Tuned is the journal under a budget of its own for Get, with 20 cases for
// each property run and an excusal, beside the options of the journal.
//
//dokimi:subject Keeper new=NewTuned cases=20 options=JournalOptions
//dokimi:excuse Keeper.Get/keeps-each-put-across-a-restart reason="the tuned journal keeps no medium"
//dokimi:budget Keeper.Get allocs=2
type Tuned struct {
	*Journal
}

// NewTuned returns an empty tuned journal.
func NewTuned() *Tuned { return &Tuned{Journal: NewJournal()} }
