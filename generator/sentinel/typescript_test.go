// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	"go.dokimi.dev/eidos/lang/typescript"
	tsbackend "go.dokimi.dev/eidos/lang/typescript/backend"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The modules of the TypeScript tree that one spec reads.
const (
	parseModule         = "parse"
	labelModule         = "label"
	unconstructedModule = "unconstructed"
	badModule           = "bad"
	relayModule         = "relay"
)

// The part of TypeScript declares the target typescript with the template of
// the family. It constructs each error class through a public constructor,
// opens no message with a prefix by default and leaves the fields of Error
// out of the members.
func TestTypeScript(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of the family for the target typescript", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, sentinel.New(), typescript.Target, testTemplate)
	})

	t.Run("writes the checks of the golden tree typescript", func(t *testing.T) {
		t.Parallel()

		p := sentinel.New()
		f := fixture(t, typescript.Target, tsTree, p, storeModule, parseModule, labelModule, unconstructedModule,
			relayModule)
		res := f.Generate(t, p)
		assert.NoError(t, res.Err, "the generator runs whole")
		golden.MatchTree(t, "typescript", gentest.Rendered(t, f, res, p, tsbackend.New()), golden.ShouldUpdate())
	})

	t.Run("records each error class under the param remote as remote", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, relayModule).checks[relayModule]
		remote := map[string]bool{}
		for _, k := range got.Kinds {
			remote[k.Name] = k.Remote
		}
		expect.Equal(t, remote, map[string]bool{"CallError": true, "TimeoutError": false},
			"CallError carries the message of the remote store")
	})

	t.Run("records each public error class with its members", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, storeModule).checks[storeModule]
		expect.Equal(t, describe(got.Kinds), []string{
			"ClosedError", "MissingError", "NotFoundError members=key:text", "WrappedError cause=cause members=op:text",
		}, "each class has the facts of its contract, and the fields of Error are no members")
		expect.Length(t, got.Sentinels, 0, "the rules of TypeScript accept no name for an error value")
	})

	t.Run("records no cause of a class whose parameter cause a required parameter follows", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, labelModule).checks[labelModule]
		expect.Equal(t, describe(got.Kinds), []string{"LabelError members=op:text"},
			"a check cannot give the cause without the parameters after it")
	})

	prefixes := []struct {
		name   string
		module string
		option string
		want   string
	}{
		{name: "opens each message with the prefix of the directive", module: storeModule, want: "kv: "},
		{name: "opens no message with a prefix by default", module: parseModule},
		{name: "opens no message with the option prefix", module: parseModule, option: "acme"},
	}
	for _, tt := range prefixes {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := sentinel.New()
			p.(plugin.OptionsProvider).Options().(*sentinel.Options).Prefix = tt.option
			got := run(t, p, typescript.Target, tsTree, tt.module).checks[tt.module]
			expect.Equal(t, got.Prefix, tt.want, "the prefix opens each message")
		})
	}

	methods := []struct {
		name   string
		module string
		give   func(sentinel.Kind) int
		want   map[string]string
	}{
		{
			name: "constructs each error class with a sample of each parameter", module: storeModule,
			give: func(k sentinel.Kind) int { return k.Method },
			want: map[string]string{
				"ClosedError":   "got = new ClosedError()",
				"MissingError":  `got = new MissingError("key-a")`,
				"NotFoundError": `got = new NotFoundError("key-a")`,
				"WrappedError":  `got = new WrappedError("op-a", "cause-a")`,
			},
		},
		{
			name: "gives the cause to the parameter of the name of the field cause", module: storeModule,
			give: func(k sentinel.Kind) int { return k.CauseMethod },
			want: map[string]string{
				"ClosedError": "", "MissingError": "", "NotFoundError": "",
				"WrappedError": `got = new WrappedError("op-a", cause)`,
			},
		},
		{
			name: "lists a construction of each other error class", module: storeModule,
			give: func(k sentinel.Kind) int { return k.PeersMethod },
			want: map[string]string{
				"ClosedError": `others = [new MissingError("key-a"), new NotFoundError("key-a"), ` +
					`new WrappedError("op-a", "cause-a")]`,
				"MissingError": `others = [new ClosedError(), new NotFoundError("key-a"), ` +
					`new WrappedError("op-a", "cause-a")]`,
				"NotFoundError": `others = [new ClosedError(), new MissingError("key-a"), ` +
					`new WrappedError("op-a", "cause-a")]`,
				"WrappedError": `others = [new ClosedError(), new MissingError("key-a"), ` +
					`new NotFoundError("key-a")]`,
			},
		},
		{
			name: "lists no other error class in a module of one", module: labelModule,
			give: func(k sentinel.Kind) int { return k.PeersMethod },
			want: map[string]string{"LabelError": ""},
		},
		{
			name: "lists the other error class in a module of two", module: parseModule,
			give: func(k sentinel.Kind) int { return k.PeersMethod },
			want: map[string]string{
				"LexError":   `others = [new ParseError("text-a", 1.5)]`,
				"ParseError": `others = [new LexError("text-a")]`,
			},
		},
		{
			name: "passes no argument to a rest parameter", module: parseModule,
			give: func(k sentinel.Kind) int { return k.Method },
			want: map[string]string{
				"LexError": `got = new LexError("text-a")`, "ParseError": `got = new ParseError("text-a", 1.5)`,
			},
		},
		{
			name: "gives the cause to a first parameter of the name of the field cause", module: unconstructedModule,
			give: func(k sentinel.Kind) int { return k.CauseMethod },
			want: map[string]string{"ChainError": "got = new ChainError(cause)", "KeptError": "", "LostError": ""},
		},
		{
			name:   "constructs an error class without a constructor of its own through the constructor of Error",
			module: unconstructedModule, give: func(k sentinel.Kind) int { return k.Method },
			want: map[string]string{
				"ChainError": `got = new ChainError("cause-a")`, "KeptError": "got = new KeptError()",
				"LostError": `got = new LostError("message-a")`,
			},
		},
		{
			name:   "passes a sample to a parameter of the name cause that a required parameter follows",
			module: labelModule, give: func(k sentinel.Kind) int { return k.Method },
			want: map[string]string{"LabelError": `got = new LabelError("cause-a", "op-a")`},
		},
		{
			name:   "gives no cause to a parameter of the name cause that a required parameter follows",
			module: labelModule, give: func(k sentinel.Kind) int { return k.CauseMethod },
			want: map[string]string{"LabelError": ""},
		},
	}
	for _, tt := range methods {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := sentinel.New()
			res := fixture(t, typescript.Target, tsTree, p, tt.module).Generate(t, p)
			assert.NoError(t, res.Err, "the generator runs whole")
			record := gentest.Records(t, res, testTemplate, sentinel.TestTag,
				func(c sentinel.Checks) string { return c.Subject })[tt.module]
			assert.NotNil(t, record.Decl, "the record of the module is a struct")
			got := map[string]string{}
			for _, k := range record.Data.Kinds {
				got[k.Name] = ""
				if at := tt.give(k); at >= 0 {
					got[k.Name] = statements(record.Decl.Methods.Items()[at].Body.Stmts)
				}
			}
			expect.Equal(t, got, tt.want, "each method of the record declares its local")
		})
	}

	findings := []struct {
		name     string
		module   string
		give     string
		wantCode diag.Code
		wantMsg  string
	}{
		{
			name: "reports Unconstructed at an error class with a private constructor", module: unconstructedModule,
			give: "SealedError", wantCode: sentinel.Unconstructed,
			wantMsg: "the error type SealedError has no public constructor whose parameters have samples, " +
				"so it has no checks",
		},
		{
			name:   "reports Unconstructed at an error class with a required parameter without a sample",
			module: unconstructedModule, give: "HandlerError", wantCode: sentinel.Unconstructed,
			wantMsg: "the error type HandlerError has no public constructor whose parameters have samples, " +
				"so it has no checks",
		},
		{
			name:     "reports Unconstructed at an error class with a default without a sample before a required parameter",
			module:   unconstructedModule,
			give:     "DefaultedError",
			wantCode: sentinel.Unconstructed,
			wantMsg: "the error type DefaultedError has no public constructor whose parameters have samples, " +
				"so it has no checks",
		},
	}
	for _, tt := range findings {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := sentinel.New()
			f := fixture(t, typescript.Target, tsTree, p, tt.module)
			res := f.Generate(t, p)
			assert.NoError(t, res.Err, "the generator runs whole")
			got := gentest.Finding(t, res, f.Find(t, symbol.KindStruct, "", tt.give))
			expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
			expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the class and the reason")
			kinds := gentest.Records(t, res, testTemplate, sentinel.TestTag,
				func(c sentinel.Checks) string { return c.Subject })[tt.module].Data.Kinds
			expect.Equal(t, describe(kinds), []string{"ChainError cause=cause", "KeptError", "LostError"},
				"the checks leave out each class that they report")
		})
	}

	t.Run("reports NoErrors at a module without an error class", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, badModule)
		assert.Length(t, got.findings, 1, "sentinel reports the module once")
		expect.Equal(t, got.findings[0].Code, sentinel.NoErrors, "the module has no error to check")
		expect.Length(t, got.checks, 0, "sentinel records no checks of the module")
	})
}
