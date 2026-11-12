// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/generator/suite"
)

// bodiesOf returns the method and the body of each override of a planted
// defect, in order.
func bodiesOf(overrides []*suite.Defect) [][2]string {
	out := make([][2]string, 0, len(overrides))
	for _, o := range overrides {
		out = append(out, [2]string{o.Method.Name, o.Body})
	}
	return out
}

func TestTxn(t *testing.T) {
	t.Parallel()

	got := generate(t, goTree, isolatedPackage)
	lists, registers := got.suites["Lists"], got.suites["Registers"]

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the actions and the clients to the probe of an isolation level", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &lists, "Append", "serializable")
			expect.Equal(t, []any{d.Probe, d.Whole, d.Sequence},
				[]any{string(spec.SerializableHistory), "Lists/serializable", true},
				"the check of the declaration runs generated sequences")
			assert.NotNil(t, lists.Sequences, "the declaration has generated sequences")
			expect.Equal(t, d.Args, []suite.Arg{{Ref: lists.Sequences.Actions}, {Expr: &suite.Expr{Word: "2"}}},
				"three clients run transactions")
		})

		t.Run("writes the action of a transaction whose calls take the handle", func(t *testing.T) {
			t.Parallel()

			txn := lists.Sequences.Txn
			assert.NotNil(t, txn, "the actions have the action of a transaction")
			expect.Equal(t, []any{txn.Key.Spelling, txn.Value.Spelling, txn.Function, txn.Appends, txn.Text},
				[]any{"string", "int", "append", true, false}, "the transaction appends numbers to the lists of keys")
			steps := lists.Sequences.Steps
			expect.Equal(t, []string{txn.Begin.Call, txn.Commit.Call, txn.Rollback.Call, txn.Write.Call, txn.Read.Call},
				[]string{steps[0].Call, steps[1].Call, steps[2].Call, steps[3].Call, steps[4].Call},
				"each call of the transaction calls the step of its callable")
			assert.Length(t, txn.Begin.Args, 1, "the begin takes its label")
			label := txn.Begin.Args[0].Value
			assert.NotNil(t, label, "the label is a sample")
			assert.NotNil(t, label.Inner, "the sample converts a literal to the type of the label")
			expect.Equal(t, label.Inner.Text, "label-a", "the begin takes the first sample of its label")
			expect.Equal(t, txn.Write.Args, []suite.Expr{{Name: "tx"}, {Name: "key"}, {Name: "v"}},
				"the write takes the handle, the key and the value")
			expect.Equal(t, txn.Read.Args, []suite.Expr{{Name: "tx"}, {Name: "key"}},
				"the read takes the handle and the key")
			expect.Empty(t, txn.Steps, "the transaction calls no method of a handle")
			expect.Equal(t, txn.Pool, "keysString", "the keys come from the pool of their type")
		})

		t.Run("writes the action of a transaction whose calls are methods of the handle", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &registers, "Begin", "snapshot-isolation")
			expect.Equal(
				t,
				[]any{d.Probe, d.Whole},
				[]any{string(spec.SnapshotHistory), "Registers/snapshot-isolation"},
				"the check of the declaration binds on the begin of the handle",
			)
			txn := registers.Sequences.Txn
			assert.NotNil(t, txn, "the actions have the action of a transaction")
			expect.Equal(t, []any{txn.Value.Spelling, txn.Function, txn.Appends, txn.Text},
				[]any{"string", "write", false, true}, "the transaction writes the texts of numbers to registers")
			handled := []bool{
				txn.Begin.Handled, txn.Write.Handled, txn.Read.Handled, txn.Commit.Handled, txn.Rollback.Handled,
			}
			expect.Equal(t, handled, []bool{false, true, true, true, true},
				"each call after the begin is a method of the handle")
			names := make([]string, 0, len(txn.Steps))
			for _, st := range txn.Steps {
				names = append(names, st.Name+" "+st.Receiver.Spelling)
			}
			expect.Equal(t, names, []string{"Set Txn", "Load Txn", "Commit Txn", "Rollback Txn"},
				"the transaction adds a step of each method of the handle")
			expect.Equal(t, txn.Write.Args, []suite.Expr{{Name: "key"}, {Name: "v"}},
				"the method of the handle takes the key and the value")
		})

		t.Run("plants the defect of an isolation level on the callables of a transaction", func(t *testing.T) {
			t.Parallel()

			defect := defectNamed(t, &lists, "ListsReadsUncommitted")
			expect.Equal(t, bodiesOf(defect.Overrides), [][2]string{
				{"Append", "keep-open"}, {"Get", "read-open"}, {"Commit", "drop-open"}, {"Rollback", "drop-open"},
			}, "the write keeps its value open until its transaction settles")
			write := defect.Overrides[0].Txn
			expect.Equal(t, []string{write.Handle, write.Key, write.Value}, []string{"p1", "p2", "p3"},
				"the write keeps the value of its key for the handle")
			expect.Nil(t, defect.Wrapper, "the defect wraps no handle")
		})

		t.Run("wraps the handle of a transaction whose calls are methods of the handle", func(t *testing.T) {
			t.Parallel()

			defect := defectNamed(t, &registers, "RegistersWritesWithoutASnapshot")
			expect.Equal(t, bodiesOf(defect.Overrides), [][2]string{{"Begin", "wrap-handle"}},
				"the begin wraps the handle that it opens")
			assert.NotNil(t, defect.Wrapper, "the defect wraps the handle")
			expect.Equal(t, bodiesOf(defect.Wrapper.Overrides), [][2]string{
				{"Set", "stage"}, {"Load", "read-staged"}, {"Commit", "apply-late"}, {"Rollback", "drop-staged"},
			}, "the write stages its value, and the commit applies it in a late transaction")
			expect.Equal(t, defect.Wrapper.EmbeddedField(), "Txn", "the wrapper embeds the handle")
			expect.Equal(t, defect.Wrapper.Overrides[0].Txn.Handle, "t", "the wrapper is the handle of its calls")
		})
	})
}
