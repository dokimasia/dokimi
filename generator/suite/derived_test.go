// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
)

func TestDerived(t *testing.T) {
	t.Parallel()

	got := generate(t, goTree, derivedPackage).suites["Store"]

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("records a derived check for each check of the derivation", func(t *testing.T) {
			t.Parallel()

			assert.Length(t, got.Derived, len(got.Derivation.Checks), "each derived check has its record")
		})

		t.Run("names the constants of the spec of each derived check", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &got, "Get", "read-after-write")
			expect.Equal(t, d.Method, "ReadAfterWrite", "the method of the identity is the segment in Pascal case")
			expect.Equal(t, d.Claim, "ReadAfterWrite", "the claim is the constant of the claim")
			expect.Equal(t, d.Probe, "WriteThenRead", "the probe is the constant of the probe")
			expect.Equal(t, d.Machineries, []string{"NothingBeyondACaller"}, "the machineries are their constants")
			expect.False(t, d.Bubble, "the probe runs outside a bubble")
		})

		t.Run("names the sample binding of a generated binding", func(t *testing.T) {
			t.Parallel()

			expect.Equal(t, derivedNamed(t, &got, "Get", "smoke#generated").Base, "Smoke",
				"the generated binding names the method of the identity of its sample binding")
			expect.Empty(t, derivedNamed(t, &got, "Get", "smoke").Base, "a sample binding names no other binding")
		})

		t.Run("passes the read and the write that returns the written value to a probe of a pair", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &got, "Get", "read-after-write")
			assert.Length(t, d.Args, 2, "the probe takes the read and the write")
			read, write := d.Args[0].Closure, d.Args[1].Closure
			assert.NotNil(t, read, "the read is a closure")
			assert.NotNil(t, write, "the write is a closure")
			expect.Equal(t, read.Method, "Get", "the read calls Get")
			expect.Equal(t, write.Method, "Put", "the write calls its partner Put")
			assert.NotNil(t, read.Args[1].Of, "the read selects a field of a value")
			expect.Equal(t, *read.Args[1].Of, write.Lets[0].Value, "the read selects the session of the write")
			expect.Equal(t, read.Args[1].Select, "ID", "the read takes the ID of the written session")
			expect.Equal(t, write.Args[1].Name, "w", "the write passes the written session")
			expect.Equal(t, write.Results[0].Name, "w", "the write returns the written session")
		})

		t.Run("shares one planted defect between the checks with one break of one callable", func(t *testing.T) {
			t.Parallel()

			names := make([]string, 0, len(got.Defects))
			for _, d := range got.Defects {
				names = append(names, d.Name)
			}
			assert.NoDuplicates(t, func() ([]string, error) { return names, nil }, "each break has one defect")
			smoke := derivedNamed(t, &got, "Get", "smoke")
			expect.Equal(t, smoke.Defect, "StoreGetPanics", "the defect of a panic of Get is StoreGetPanics")
			raw := derivedNamed(t, &got, "Get", "read-after-write")
			expect.Equal(
				t,
				raw.Defect,
				"StoreGetReadsAnotherValue",
				"the defect of read-after-write breaks the read Get",
			)
		})

		vaults := generate(t, goTree, vaultPackage).suites

		t.Run("passes the write and the read of each sample to a durable probe", func(t *testing.T) {
			t.Parallel()

			v := vaults["Vault"]
			d := derivedNamed(t, &v, "Put", "synced-writes-survive")
			assert.Length(t, d.Args, 1, "the probe takes the calls of the durable writer")
			durable := d.Args[0].Durable
			assert.NotNil(t, durable, "the argument is the calls of a durable writer")
			for i := range 2 {
				write, read := durable.Writes[i], durable.Reads[i]
				expect.Equal(t, []any{write.Receiver.Name, write.Method, read.Receiver.Name, read.Method},
					[]any{"s", "Put", "s", "Get"}, "the write and the read call the subject that they take")
				expect.Equal(t, durable.Wants[i], write.Args[2], "the value of a sample is the value of its write")
				expect.Equal(t, durable.Keys[i], []suite.Expr{write.Args[1]},
					"the key of a sample is the argument of its read after the context")
				expect.Equal(t, read.Args[1], write.Args[1], "the read of a sample reads the key of its write")
			}
			expect.NotEqual(t, durable.Wants[0], durable.Wants[1], "the two samples write two values")
			expect.Equal(t, d.Door, "VaultPutAcknowledgesBeforeSyncOpen",
				"the check passes the door of its planted defect")
		})

		t.Run("passes the declared error of a corrupted record to the probe of corruption", func(t *testing.T) {
			t.Parallel()

			v := vaults["Vault"]
			d := derivedNamed(t, &v, "Put", "corruption-is-reported")
			assert.Length(t, d.Args, 2, "the probe takes the calls of the durable writer and the error")
			assert.NotNil(t, d.Args[1].Expr, "the error is an expression")
			assert.NotNil(t, d.Args[1].Expr.Value, "the error is a value")
			expect.Equal(t, d.Args[1].Expr.Value.Callee.Name, "ErrCorrupt", "the error is the name that corrupt sets")
		})

		t.Run("returns the error of a write that also returns a value", func(t *testing.T) {
			t.Parallel()

			j := vaults["Journal"]
			write := derivedNamed(t, &j, "Append", "synced-writes-survive").Args[0].Durable.Writes[0]
			expect.Equal(t, write.Assign, []string{"_", "err"}, "the write keeps the error of Append alone")
			expect.Equal(t, write.Results, []suite.Expr{{Name: "err"}}, "the write returns the error of Append")
		})

		t.Run("passes the world of the declaration to the probe of a world", func(t *testing.T) {
			t.Parallel()

			r := vaults["Registry"]
			d := derivedNamed(t, &r, "Value", "linearizable-under-faults")
			assert.NotNil(t, r.World, "the suite of Registry has a world")
			expect.Equal(t, []any{d.Whole, d.World, d.Cased, d.Sequence}, []any{
				"Registry/linearizable-under-faults", true, true, false,
			}, "the check of the declaration runs its world in a case of a property check")
			assert.Length(t, d.Args, 1, "the probe takes the world alone")
			expect.Equal(t, d.Args[0].Ref, r.World.Name, "the argument is the world of Registry")
			expect.Equal(
				t,
				[]string{d.Defect, d.Door},
				[]string{"RegistryAmnesiacRecovery", "RegistryAmnesiacRecoveryHost"},
				"the check passes the door Host of its planted defect",
			)
		})

		t.Run("binds the recovery of a world on the first callable of the declaration", func(t *testing.T) {
			t.Parallel()

			l := vaults["Ledger"]
			d := derivedNamed(t, &l, "Record", "completes-after-recovery")
			expect.Equal(t, []any{d.Whole, d.World, d.Defect}, []any{
				"Ledger/completes-after-recovery", true, "LedgerBlocksAfterRecovery",
			}, "the check of the declaration plants the defect that blocks after a restart")
		})

		t.Run("reads a key without a value as an error through a read with a presence flag", func(t *testing.T) {
			t.Parallel()

			l := vaults["Ledger"]
			d := derivedNamed(t, &l, "Record", "no-torn-record")
			read := d.Args[0].Durable.Reads[0]
			expect.True(t, read.Absent, "the read returns an error for a key without a value")
			expect.Equal(t, read.Assign, []string{"r0", "ok"}, "the read keeps the value and the presence flag")
			corruption := func(d suite.Derived) bool { return d.Claim == "CorruptionIsReported" }
			expect.False(t, slices.ContainsFunc(l.Derived, corruption),
				"a writer without the param corrupt has no check of corruption")
		})

		tallies := generate(t, goTree, drawnPackage).suites["Tallies"]
		lacks := func(s *suite.Suite, field, segment string) bool {
			return !slices.ContainsFunc(s.Derived, func(d suite.Derived) bool {
				return d.Field == field && d.Segment == segment
			})
		}

		t.Run("records the generated binding of a check whose calls take samples", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &tallies, "Limit", "smoke#generated")
			expect.True(t, d.Generated, "the check generates the inputs of its calls")
			expect.Equal(t, d.Method, "SmokeGenerated", "the method of the identity adds the word Generated")
			expect.Equal(t, d.Machineries, []string{"GeneratedInputs", "NothingBeyondACaller"},
				"the machineries add generated-inputs before nothing-beyond-a-caller")
			expect.False(t, derivedNamed(t, &tallies, "Limit", "smoke").Generated,
				"the sample binding passes the samples of the inputs")
		})

		t.Run("records no generated binding of a check whose calls take no input", func(t *testing.T) {
			t.Parallel()

			derivedNamed(t, &tallies, "Count", "smoke")
			expect.True(t, lacks(&tallies, "Count", "smoke#generated"), "the check of Count draws no input")
		})

		t.Run("passes the drawn inputs of a combining callable to the probe of commutative", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &tallies, "Merge", "commutative")
			expect.True(t, d.Generated, "the claim of generated-inputs generates the inputs of its calls")
			assert.Length(t, d.Args, 3, "the probe takes the call and the two inputs")
			expect.NotNil(t, d.Args[0].Closure, "the call is a closure")
			expect.Equal(t, drawnOf(t, d.Args[1].Expr).Draw.Name, "drawnA", "the first input is the drawn a")
			expect.Equal(t, drawnOf(t, d.Args[2].Expr).Draw.Name, "drawnB", "the second input is the drawn b")
		})

		t.Run("passes a second value of the second input to the probe of associative", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &tallies, "Merge", "associative")
			assert.Length(t, d.Args, 4, "the probe takes the call and the three inputs")
			third := drawnOf(t, d.Args[3].Expr)
			expect.Equal(t, third.Draw.Name, "drawnB", "the third input is a value of b")
			expect.True(t, third.Second, "the third input is the second value of b")
			expect.False(t, drawnOf(t, d.Args[2].Expr).Second, "the second input is the first value of b")
			expect.True(t, third.Draw.Pair, "the check draws two distinct values of b")
		})

		faults := generate(t, goTree, faultsPackage).suites

		t.Run("passes the names of the callables of open and close to the probe of a release", func(t *testing.T) {
			t.Parallel()

			copier := faults["Copier"]
			d := derivedNamed(t, &copier, "Copy", "releases-on-every-path")
			assert.Length(t, d.Args, 3, "the probe takes the call, the acquire and the release")
			assert.NotNil(t, d.Args[0].Closure, "the call is a closure")
			expect.Equal(t, d.Args[0].Closure.Method, "Copy", "the call calls Copy")
			expect.Equal(t, d.Args[1].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Open"))},
				"the acquire is the name of the callable that open names")
			expect.Equal(t, d.Args[2].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Close"))},
				"the release is the name of the callable that close names")
		})

		t.Run("passes the read of the mixin atomic to the probe of a failed call", func(t *testing.T) {
			t.Parallel()

			writer := faults["Writer"]
			d := derivedNamed(t, &writer, "Put", "no-partial-write")
			assert.Length(t, d.Args, 2, "the probe takes the call and the read")
			call, read := d.Args[0].Closure, d.Args[1].Closure
			assert.NotNil(t, call, "the call is a closure")
			assert.NotNil(t, read, "the read is a closure")
			expect.Equal(t, call.Method, "Put", "the call calls Put")
			expect.Equal(t, read.Method, "Get", "the read calls the reader that read names")
		})

		families := generate(t, goTree, familiesPackage).suites
		family := func(name string) *suite.Suite {
			s := families[name]
			return &s
		}
		t.Run("records no generated binding of a check whose probe fills each input of its calls", func(t *testing.T) {
			t.Parallel()

			screened := family("Screened")
			derivedNamed(t, screened, "Store", "validates")
			expect.True(t, lacks(screened, "Store", "validates#generated"),
				"the probe of validates passes each value that Store takes")
		})

		t.Run("passes the actions, the reference, the clients and the budget to the probe of linearizable",
			func(t *testing.T) {
				t.Parallel()

				board := generate(t, goTree, locksPackage).suites["Board"]
				d := derivedNamed(t, &board, "Read", "linearizable")
				expect.Equal(
					t,
					[]any{d.Probe, d.Whole, d.Sequence},
					[]any{"LinearizableHistory", "Board/linearizable", true},
					"the check of the whole declaration runs generated sequences",
				)
				assert.NotNil(t, board.Sequences, "the declaration has the generated sequences of its reference")
				expect.Equal(t, d.Args, []suite.Arg{
					{Ref: board.Sequences.Actions},
					{Ref: board.Sequences.Reference},
					{Expr: &suite.Expr{Word: "2"}},
					{Budget: 200_000},
				}, "three clients call the subject, and the checker takes at most 200,000 steps")
			})

		t.Run("passes the actions, the read, the read again and the clients to the probe of concurrent reads",
			func(t *testing.T) {
				t.Parallel()

				shelf := generate(t, goTree, locksPackage).suites["Shelf"]
				d := derivedNamed(t, &shelf, "Count", "reads-agree-under-concurrency")
				expect.Equal(t, []any{d.Probe, d.Sequence}, []any{"ConcurrentReadsAgree", true},
					"the check runs generated sequences")
				assert.NotNil(t, shelf.Sequences, "the declaration has generated sequences")
				assert.Length(t, d.Args, 4, "the probe takes the actions, the read, the read again and the clients")
				expect.Equal(t, d.Args[0], suite.Arg{Ref: shelf.Sequences.Actions}, "the actions are those of Shelf")
				expect.Equal(t, d.Args[1].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Count"))},
					"the read is the operation Count")
				assert.NotNil(t, d.Args[2].Apply, "the read again applies Count to the arguments of a read")
				expect.Equal(t, d.Args[3], suite.Arg{Expr: &suite.Expr{Word: "2"}}, "three clients read at once")
			})

		t.Run("passes the actions, the read, the write and the clients to the probe of a pinned read",
			func(t *testing.T) {
				t.Parallel()

				snapshots := family("Snapshots")
				d := derivedNamed(t, snapshots, "Get", "pins-the-snapshot-under-writes")
				expect.Equal(t, []any{d.Probe, d.Sequence}, []any{"PinnedUnderConcurrentWrites", true},
					"the check runs generated sequences")
				assert.NotNil(t, snapshots.Sequences, "the declaration has generated sequences")
				expect.Equal(t, d.Args, []suite.Arg{
					{Ref: snapshots.Sequences.Actions},
					{Expr: &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Get"))}},
					{Expr: &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Store"))}},
					{Expr: &suite.Expr{Word: "2"}},
				}, "the probe reads with Get beside the writes of Store on three clients")
				expect.Equal(t, d.Entails,
					[]suite.Identity{{Field: "Get", Method: "PointInTime", ID: "Snapshots.Get/point-in-time"}},
					"the check entails the check of a single caller of its law")
			})

		t.Run("passes the actions, the read, the write, the version and the clients to the probe of a causal order",
			func(t *testing.T) {
				t.Parallel()

				stamped := family("Stamped")
				d := derivedNamed(t, stamped, "Get", "causal-history")
				expect.Equal(t, []any{d.Probe, d.Sequence}, []any{"VersionsFollowTheCausalOrder", true},
					"the check runs generated sequences")
				assert.NotNil(t, stamped.Sequences, "the declaration has generated sequences")
				assert.Length(t, d.Args, 5,
					"the probe takes the actions, the read, the write, the version and the clients")
				expect.Equal(t, d.Args[:3], []suite.Arg{
					{Ref: stamped.Sequences.Actions},
					{Expr: &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Get"))}},
					{Expr: &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Store"))}},
				}, "the probe reads with Get beside the writes of Store")
				assert.NotNil(t, d.Args[3].Projection, "the version is a projection")
				expect.Equal(t, d.Args[3].Projection.Field, "Rev", "the version reads the field that version names")
				expect.Equal(t, d.Args[4], suite.Arg{Expr: &suite.Expr{Word: "2"}}, "three clients call at once")
				expect.Equal(t, d.Entails, []suite.Identity{{Field: "Get", Method: "Causal", ID: "Stamped.Get/causal"}},
					"the check entails the check of a single caller of its law")
			})

		flight := generate(t, goTree, locksPackage).suites["Flight"]

		t.Run("passes a call with a computation for each input of a function type to the probe of one computation",
			func(t *testing.T) {
				t.Parallel()

				d := derivedNamed(t, &flight, "Do", "computes-once")
				expect.Equal(t, []any{d.Probe, d.Bubble}, []any{"SharesOneComputation", true},
					"the probe waits inside a bubble until each call blocks")
				assert.Length(t, d.Args, 2, "the probe takes the call and the clients")
				expect.Equal(t, d.Args[1], suite.Arg{Expr: &suite.Expr{Word: "2"}}, "three clients call at once")
				call := d.Args[0].Closure
				assert.NotNil(t, call, "the call is a closure")
				expect.Equal(t, []any{call.Method, call.Error, len(call.Returns)}, []any{"Do", true, 1},
					"the call returns the value of Do beside its error")
				assert.Length(t, call.Args, 3, "the call passes the context, the key and the load")
				expect.Nil(t, call.Args[1].Compute, "the key keeps its sample")
				load := call.Args[2]
				assert.NotNil(t, load.Compute, "the load is a computation")
				expect.Equal(t, load.Name, "Do", "the computation records a call of Do")
				signature := make([]string, 0, load.Compute.Params.Len()+load.Compute.Returns.Len())
				for _, p := range load.Compute.Params.Items() {
					signature = append(signature, spelled(p.Type))
				}
				for _, r := range load.Compute.Returns.Items() {
					signature = append(signature, r.Name+" "+spelled(r.Type))
				}
				expect.Equal(t, signature, []string{"context.Context", "string", "r0 string", "r1 error"},
					"the computation has the parameters and the named results of Loader")
			})

		locker := generate(t, goTree, locksPackage).suites["Locker"]

		t.Run("passes the case, the call, the sentinel and the clients to the probe of one holder", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &locker, "Lock", "one-holder-at-a-time")
			expect.Equal(t, []any{d.Probe, d.Cased, d.Sequence}, []any{"ExclusiveUnderConcurrency", true, false},
				"the probe takes the case, whose scheduler runs the clients of a planted defect")
			assert.Length(t, d.Args, 3, "the probe takes the call, the sentinel and the clients")
			assert.NotNil(t, d.Args[0].Closure, "the call is a closure")
			expect.Equal(t, d.Args[0].Closure.Method, "Lock", "the call calls Lock")
			assert.NotNil(t, d.Args[1].Expr, "the sentinel is an expression")
			expect.Equal(t, d.Args[1].Expr.Value.Callee.Name, "ErrLocked", "the sentinel is the value that taken names")
			expect.Equal(t, d.Args[2], suite.Arg{Expr: &suite.Expr{Word: "2"}}, "three clients call at once")
			expect.Equal(t, d.Entails, []suite.Identity{{Field: "Lock", Method: "Lease", ID: "Locker.Lock/lease"}},
				"the check entails the check of a single caller of its law")
		})

		t.Run("passes a computation beside the drawn key in the generated binding of one computation",
			func(t *testing.T) {
				t.Parallel()

				d := derivedNamed(t, &flight, "Do", "computes-once#generated")
				assert.Length(t, d.Draws, 1, "the check draws the key alone")
				expect.Equal(t, d.Draws[0].Name, "drawnKey", "the draw is the key")
				call := d.Args[0].Closure
				assert.NotNil(t, call, "the call is a closure")
				expect.Equal(t, drawnOf(t, &call.Args[1]).Draw, d.Draws[0], "the call passes the drawn key")
				expect.NotNil(t, call.Args[2].Compute, "the load is a computation")
			})

		records := generate(t, goTree, bufferedPackage).suites["Records"]

		t.Run("passes the call, the field of the result and the field of the request to the probe of a field beside "+
			"an error", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &records, "Read", "returns-on-error")
			expect.Equal(t, d.Probe, "FieldBesideError", "the probe compares the field beside an error")
			assert.Length(t, d.Args, 3,
				"the probe takes the call, the field of the result and the field of the request")
			call := d.Args[0].Closure
			assert.NotNil(t, call, "the call is a closure")
			expect.Equal(t, call.Method, "Read", "the call calls Read")
			projection := d.Args[1].Projection
			assert.NotNil(t, projection, "the field of the result is a projection")
			expect.Equal(t, []string{projection.Type.Spelling, projection.Result.Spelling, projection.Field},
				[]string{"Result", "[]byte", "Record"}, "the projection reads the field that result names")
			assert.Length(t, call.Lets, 2, "the call declares the request and sets its field")
			expect.Equal(t, []string{call.Lets[1].Name, call.Lets[1].Select}, []string{call.Lets[0].Name, "Dst"},
				"the call sets the field of the request that field names")
			expect.Equal(t, call.Args[1], suite.Expr{Name: call.Lets[0].Name}, "the call passes the request")
			want := d.Args[2].Expr
			assert.NotNil(t, want, "the field of the request is an expression")
			expect.Equal(t, *want, call.Lets[1].Value, "the probe expects the value of the field that the call sets")
			assert.NotNil(t, want.Value, "the value of the field is a sample")
			expect.Length(t, want.Value.Fields, 1, "the sample of the field has one byte, so it is not the zero value")
		})

		t.Run("projects the whole result of a callable without the param result", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &records, "Fill", "returns-on-error")
			projection := d.Args[1].Projection
			assert.NotNil(t, projection, "the field of the result is a projection")
			expect.Equal(t, []string{projection.Type.Spelling, projection.Result.Spelling, projection.Field},
				[]string{"[]byte", "[]byte", ""}, "the projection returns the whole result")
		})

		t.Run("passes the field of the drawn request to the generated binding of a field beside an error",
			func(t *testing.T) {
				t.Parallel()

				d := derivedNamed(t, &records, "Read", "returns-on-error#generated")
				call := d.Args[0].Closure
				assert.NotNil(t, call, "the call is a closure")
				expect.Empty(t, call.Lets, "the call passes the drawn request, whose field the case draws")
				want := d.Args[2].Expr
				assert.NotNil(t, want, "the field of the request is an expression")
				expect.Equal(t, drawnOf(t, want.Of).Draw, drawnOf(t, &call.Args[1]).Draw,
					"the field is a field of the drawn request of the call")
			})

		reach := generate(t, goTree, wrapPackage, toolPackage, deepPackage).suites["Reach"]
		refusals := []struct {
			name     string
			give     *suite.Suite
			giveFrom string
			want     string
		}{
			{
				name: "refuses a check of a cursor whose handle is not an interface", give: family("Sessions"),
				giveFrom: "Open", want: "cursor-open",
			},
			{
				name: "refuses a check whose call passes a sample that the suite cannot import", give: &reach,
				giveFrom: "Store", want: "smoke",
			},
			{
				name: "refuses a check whose write declares a sample that the suite cannot import", give: &reach,
				giveFrom: "Current", want: "read-after-write",
			},
			{
				name: "refuses a check whose planted defect returns a sample that the suite cannot import",
				give: &reach, giveFrom: "Gear", want: "zero-on-error",
			},
			{
				name: "refuses a check of a field beside an error whose field no input of the call has", give: &records,
				giveFrom: "Peek", want: "returns-on-error",
			},
		}
		for _, tt := range refusals {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				refused := slices.ContainsFunc(tt.give.Derivation.Refusals, func(r derive.Refusal) bool {
					return r.Callable == tt.giveFrom && string(r.Claim) == tt.want
				})
				expect.True(t, refused, "the suite refuses the claim "+tt.want+" of "+tt.giveFrom)
				derived := slices.ContainsFunc(tt.give.Derived, func(d suite.Derived) bool {
					return d.Field == tt.giveFrom && d.Segment == tt.want
				})
				expect.False(t, derived, "the suite derives no check of the refused claim")
			})
		}
	})
}
