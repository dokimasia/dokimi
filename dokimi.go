// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package dokimi

import (
	"slices"

	"go.dokimi.dev/eidos/core/authored"
	"go.dokimi.dev/eidos/core/workspace"
	golang "go.dokimi.dev/eidos/lang/go"
	goannotate "go.dokimi.dev/eidos/lang/go/annotate"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gofrontend "go.dokimi.dev/eidos/lang/go/frontend"
	gorules "go.dokimi.dev/eidos/lang/go/rules"
	"go.dokimi.dev/eidos/lang/typescript"
	tsbackend "go.dokimi.dev/eidos/lang/typescript/backend"
	tsfrontend "go.dokimi.dev/eidos/lang/typescript/frontend"
	tsrules "go.dokimi.dev/eidos/lang/typescript/rules"
	"go.dokimi.dev/eidos/plugin/shape/catalog"
	"go.dokimi.dev/eidos/sdk/output"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/generator/suite"
)

// Brand is the brand of dokimi. A carrier sets a directive as
// //dokimi:builder, the config file is .dokimi.yaml, and the state directory
// is .dokimi/.
const Brand output.Brand = "dokimi"

// GoPlan and TypeScriptPlan are the names of the two plans.
const (
	GoPlan         = "go"
	TypeScriptPlan = "typescript"
)

// Base returns the composition of dokimi without its plans. It registers the
// brand dokimi, the frontends and the rules of Go and TypeScript, the targets
// golang and typescript, and the annotators. The annotators are the kernel
// annotators of the authored samples and witnesses, the annotator of the Go
// satellite, defaults, fault, and the annotators of the shape catalog. The
// stub generator reads the faults that fault stamps and the mixins that the
// shape catalog stamps. Each call returns new plugin instances.
func Base() *workspace.Builder {
	return workspace.New().
		Brand(Brand).
		Frontends(gofrontend.New(nil), tsfrontend.New()).
		Rules(gorules.New(), tsrules.New()).
		Targets(golang.Target, typescript.Target).
		Annotators(authored.Sample(), authored.Witness(), goannotate.New(), defaults.New(), fault.New()).
		Annotators(catalog.Annotators()...)
}

// Plans returns the two plans. The plan go reads the Go packages of the
// workspace and writes Go through the Go backend. The plan typescript reads
// the TypeScript sources of the workspace and writes TypeScript through the
// TypeScript backend. Each plan reads the sources of its own language alone,
// so a generator writes each declaration in its own language. Its language
// alone scopes each plan, so a plan also reads the packages of its language
// that a store provides, such as the standard library of Go. A double of an
// interface that embeds io.Closer then has the method Close. Both plans share
// one instance of each generator: builder, stub, enum and sentinel. The plan
// go also runs suite, which writes the suites of Go. builder also has an
// annotate half, which the workspace runs with the annotators. It stamps the
// names of each builder, which the stub generator reads for the default
// results of a double. Each call returns new plugin instances.
func Plans() []workspace.Plan {
	generators := []plugin.Generator{builder.New(), stub.New(), enum.New(), sentinel.New()}
	return []workspace.Plan{
		{
			Name:       GoPlan,
			Sources:    workspace.Sources{Lang: golang.Lang},
			Generators: append(slices.Clone(generators), suite.New()),
			Backend:    gobackend.New(),
		},
		{
			Name:       TypeScriptPlan,
			Sources:    workspace.Sources{Lang: typescript.Lang},
			Generators: generators,
			Backend:    tsbackend.New(),
		},
	}
}

// Compose returns Base with Plans, the composition of the binary.
func Compose() *workspace.Builder {
	return Base().Plans(Plans()...)
}
