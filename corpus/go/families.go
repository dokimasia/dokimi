// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package corpus contains the shape families of the Go corpus. A family is a
// group of subjects whose state has one representation, such as a map keyed
// by the input. Each family has at least one control: a correct
// implementation of the declaration of one of its subjects, with another
// representation, which the corpus applies the suite to. A check that fails
// on a control forbids behaviour that the declaration permits.
package corpus

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
)

// Family is one shape family of the corpus. Each subject and each control is
// the directory of its file relative to the corpus, a dot and the name of its
// type, such as iface/detector/reader/readertest.InMemory.
type Family struct {
	// Name is the name of the family, such as stores-and-readers.
	Name string

	// Representation is the representation of the state of its subjects.
	Representation string

	// Subjects are the subjects of the family.
	Subjects []string

	// Controls are the controls of the family.
	Controls []string
}

// subjectDirective opens the comment line that declares a subject.
const subjectDirective = "//dokimi:subject "

// Subjects returns every type of the Go files of fsys whose doc comment has
// the directive subject, in the form of [Family], in the lexical order of
// their files.
func Subjects(fsys fs.FS) ([]string, error) {
	var found []string
	err := fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || path.Ext(name) != ".go" {
			return nil
		}
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return fmt.Errorf("corpus: %w", err)
		}
		declared := false
		for line := range strings.Lines(string(data)) {
			if strings.HasPrefix(line, subjectDirective) {
				declared = true
				continue
			}
			if typ, isType := strings.CutPrefix(line, "type "); isType && declared {
				found = append(found, path.Dir(name)+"."+strings.FieldsFunc(typ, isNameEnd)[0])
				declared = false
			}
		}
		return nil
	})
	return found, err
}

// isNameEnd reports whether r ends the name of a type in its declaration: a
// space, or the bracket of its type parameters.
func isNameEnd(r rune) bool { return r == ' ' || r == '[' }

// Problems returns the problems of families against the subjects that the
// corpus declares: a family without a control, a subject or a control of a
// family that the corpus does not declare, and a declared subject that no
// family contains or that two families contain.
func Problems(families []Family, subjects []string) []string {
	var found []string
	listed := map[string]string{}
	for _, f := range families {
		if len(f.Controls) == 0 {
			found = append(found, fmt.Sprintf("the family %s has no control", f.Name))
		}
		for _, s := range slices.Concat(f.Subjects, f.Controls) {
			if other, twice := listed[s]; twice {
				found = append(found, fmt.Sprintf("%s is in the families %s and %s", s, other, f.Name))
			}
			listed[s] = f.Name
			if !slices.Contains(subjects, s) {
				found = append(found,
					fmt.Sprintf("the family %s lists %s, which the corpus does not declare", f.Name, s))
			}
		}
	}
	for _, s := range subjects {
		if _, ok := listed[s]; !ok {
			found = append(found, s+" is in no family")
		}
	}
	return found
}

// Families are the shape families of the Go corpus.
var Families = []Family{
	{
		Name:           "stores-and-readers",
		Representation: "a map keyed by the input, with a counter or a revision beside it, or a single value",
		Subjects: []string{
			"iface/composite/batched-mixins/batchedmixinstest.InMemory",
			"iface/contract/batch-writer/batchwritertest.InMemory",
			"iface/contract/cas/castest.InMemory",
			"iface/contract/circuit-breaker/circuitbreakertest.InMemory",
			"iface/contract/if-absent/ifabsenttest.InMemory",
			"iface/contract/if-match/ifmatchtest.InMemory",
			"iface/contract/pagination/paginationtest.InMemory",
			"iface/contract/persister/persistertest.InMemory",
			"iface/contract/updater/updatertest.InMemory",
			"iface/contract/upserter/upsertertest.InMemory",
			"iface/detector/answeringwriter/answeringwritertest.InMemory",
			"iface/detector/batchreader/batchreadertest.InMemory",
			"iface/detector/compositewriter/compositewritertest.InMemory",
			"iface/detector/deleter/deletertest.InMemory",
			"iface/detector/lookup/lookuptest.InMemory",
			"iface/detector/multiargwriter/multiargwritertest.InMemory",
			"iface/detector/multireader/multireadertest.InMemory",
			"iface/detector/mutator/mutatortest.InMemory",
			"iface/detector/pointerreader/pointerreadertest.InMemory",
			"iface/detector/reader/readertest.InMemory",
			"iface/detector/readernoerror/readernoerrortest.InMemory",
			"iface/detector/readerwithbool/readerwithbooltest.InMemory",
			"iface/detector/writer/writertest.InMemory",
			"iface/door/durable/durabletest.Journal",
			"iface/lang/embedded/embeddedtest.InMemory",
			"iface/lang/embeddedforeign/embeddedforeigntest.InMemory",
			"iface/lang/generic/memory.Store",
			"iface/lang/genericbound/genericboundtest.InMemory",
			"iface/lang/multireturn/multireturntest.InMemory",
			"iface/lang/namedreturns/namedreturnstest.InMemory",
			"iface/lang/partnernaming/partnernamingtest.InMemory",
			"iface/lang/roledtypes/roledtypestest.InMemory",
			"iface/lang/sampled/sampledtest.InMemory",
			"iface/lang/seededreader/seededreadertest.InMemory",
			"iface/lang/suitedhandle/suitedhandletest.InMemoryFS",
			"iface/lang/variadic/variadictest.InMemory",
			"iface/mixin/accumulates/accumulatestest.InMemory",
			"iface/mixin/atomic/atomictest.InMemory",
			"iface/mixin/causal/causaltest.InMemory",
			"iface/mixin/concurrent/concurrenttest.InMemory",
			"iface/mixin/concurrentreaders/concurrentreaderstest.InMemory",
			"iface/mixin/conservative/conservativetest.InMemory",
			"iface/mixin/crdtmerge/crdtmergetest.InMemory",
			"iface/mixin/defaultonerror/defaultonerrortest.InMemory",
			"iface/mixin/deleteremoves/deleteremovestest.InMemory",
			"iface/mixin/deprecated/deprecatedtest.InMemory",
			"iface/mixin/errors/errorstest.InMemory",
			"iface/mixin/idempotent/idempotenttest.InMemory",
			"iface/mixin/injectionsafe/injectionsafetest.InMemory",
			"iface/mixin/monotonicreads/monotonicreadstest.InMemory",
			"iface/mixin/monotonicwrites/monotonicwritestest.InMemory",
			"iface/mixin/nilsafe/nilsafetest.InMemory",
			"iface/mixin/noduplicates/noduplicatestest.InMemory",
			"iface/mixin/partition/partitiontest.InMemory",
			"iface/mixin/readafterwrite/readafterwritetest.InMemory",
			"iface/mixin/readyourwrites/readyourwritestest.InMemory",
			"iface/mixin/restrictedpool/restrictedpooltest.InMemory",
			"iface/mixin/retrysucceeds/retrysucceedstest.InMemory",
			"iface/mixin/scope/scopetest.InMemory",
			"iface/mixin/sideeffect/sideeffecttest.InMemory",
			"iface/mixin/tamperevident/tamperevidenttest.InMemory",
			"iface/mixin/validates/validatestest.InMemory",
			"iface/mixin/writesfollowreads/writesfollowreadstest.InMemory",
		},
		Controls: []string{"iface/mixin/readafterwrite/sorted.Store"},
	},
	{
		Name:           "appenders-and-logs",
		Representation: "a slice that grows at its end, with a position beside it",
		Subjects: []string{
			"iface/composite/causal-chain/archive.Log",
			"iface/composite/causal-chain/replica.Log",
			"iface/composite/paginated-reader/paginatedreadertest.InMemory",
			"iface/contract/appender/appendertest.InMemory",
			"iface/contract/chain/chaintest.InMemory",
			"iface/contract/cursor/cursortest.InMemory",
			"iface/contract/cursor-open/cursoropentest.InMemory",
			"iface/contract/saga/sagatest.InMemory",
			"iface/detector/aggregator/aggregatortest.InMemory",
			"iface/detector/multiaggregator/multiaggregatortest.InMemory",
			"iface/detector/predicate/predicatetest.InMemory",
			"iface/detector/streamconsumer/streamconsumertest.InMemory",
			"iface/detector/streamreader/streamreadertest.InMemory",
			"iface/lang/dynamic/dynamictest.InMemory",
			"iface/mixin/bounded/boundedtest.InMemory",
			"iface/mixin/indexed/indexedtest.InMemory",
			"iface/mixin/overmatch/overmatchtest.InMemory",
			"iface/mixin/permutation/permutationtest.InMemory",
			"iface/mixin/stableorder/stableordertest.InMemory",
			"iface/mixin/streamreflectsmutations/streamreflectsmutationstest.InMemory",
		},
		Controls: []string{"iface/mixin/stableorder/chunked.List"},
	},
	{
		Name:           "caches",
		Representation: "a map in front of a backing store",
		Subjects: []string{
			"iface/contract/cache/cachetest.InMemory",
			"iface/contract/singleflight/singleflighttest.InMemory",
			"iface/mixin/boundedcache/boundedcachetest.InMemory",
			"iface/mixin/cacheable/cacheabletest.InMemory",
			"iface/mixin/sticky/stickytest.InMemory",
		},
		Controls: []string{"iface/contract/cache/oneentry.Cache"},
	},
	{
		Name:           "queues",
		Representation: "a channel or a slice that the oldest element leaves first",
		Subjects: []string{
			"iface/contract/pool/pooltest.InMemory",
			"iface/door/queue/queuetest.InMemory",
			"iface/mixin/eventually/eventuallytest.InMemory",
			"iface/mixin/leakfree/leakfreetest.InMemory",
		},
		Controls: []string{"iface/door/queue/ring.Buffer"},
	},
	{
		Name:           "publishers",
		Representation: "a slice of the channels or the functions of the subscribers",
		Subjects: []string{
			"iface/contract/outbox/outboxtest.InMemory",
			"iface/contract/publisher/publishertest.InMemory",
			"iface/contract/publisher-atleastonce/publisheratleastoncetest.InMemory",
			"iface/contract/publisher-atmostonce/publisheratmostoncetest.InMemory",
			"iface/contract/publisher-exactlyonce/publisherexactlyoncetest.InMemory",
			"iface/contract/publisher-redeliver/publisherredelivertest.InMemory",
			"iface/contract/watcher/watchertest.InMemory",
			"iface/mixin/hooks/hookstest.InMemory",
		},
		Controls: []string{"iface/contract/publisher/mailbox.Publisher"},
	},
	{
		Name:           "flags-and-lifecycles",
		Representation: "a flag or a small state machine",
		Subjects: []string{
			"iface/composite/tx-with-retry/txwithretrytest.InMemory",
			"iface/contract/workflow/workflowtest.InMemory",
			"iface/detector/closer/closertest.InMemory",
			"iface/detector/lifecycle/lifecycletest.InMemory",
			"iface/detector/poisonaccessor/poisonaccessortest.InMemory",
			"iface/detector/voidlifecycle/voidlifecycletest.InMemory",
			"iface/mixin/idempotentclose/idempotentclosetest.InMemory",
			"iface/mixin/integrationonly/integrationonlytest.InMemory",
			"iface/mixin/lifecycleafterclose/lifecycleafterclosetest.InMemory",
			"iface/mixin/orderafter/orderaftertest.InMemory",
			"iface/mixin/poisonable/poisonabletest.InMemory",
			"iface/mixin/wrappedvia/wrappedviatest.InMemory",
		},
		Controls: []string{"iface/detector/lifecycle/atomicstate.Lifecycle"},
	},
	{
		Name:           "counters",
		Representation: "one integer",
		Subjects: []string{
			"iface/lang/crosspackage/crosspackagetest.InMemory",
			"iface/lang/nested/nestedtest.InMemory",
			"iface/lang/nocontext/nocontexttest.InMemory",
			"iface/lang/opaque/opaquetest.InMemory",
			"iface/mixin/associative/associativetest.InMemory",
			"iface/mixin/commutative/commutativetest.InMemory",
			"iface/mixin/monotonic/monotonictest.InMemory",
			"iface/mixin/sample/sampletest.InMemory",
		},
		Controls: []string{"iface/mixin/commutative/perkey.Counter"},
	},
	{
		Name:           "stateless",
		Representation: "no state, or a value that the constructor sets",
		Subjects: []string{
			"iface/contract/codec/codectest.InMemory",
			"iface/contract/codec-lossy/codeclossytest.InMemory",
			"iface/detector/pure/puretest.InMemory",
			"iface/lang/kinds/kindstest.InMemory",
			"iface/lang/recursive/recursivetest.InMemory",
			"iface/lang/sampledparam/sampledparamtest.SHA256",
			"iface/lang/selfnamed/selfnamedtest.InMemory",
			"iface/lang/suitedhandle/suitedhandletest.InMemoryFile",
			"iface/mixin/pure/puretest.InMemory",
			"iface/mixin/total/totaltest.InMemory",
			"iface/mixin/xsssafe/xsssafetest.InMemory",
		},
		Controls: []string{"iface/mixin/xsssafe/htmlescape.Renderer"},
	},
	{
		Name:           "clocked",
		Representation: "instants of a clock beside the values",
		Subjects: []string{
			"iface/contract/rate-limit/ratelimittest.InMemory",
			"iface/mixin/scheduled/scheduledtest.InMemory",
			"iface/mixin/timeaware/timeawaretest.InMemory",
			"iface/mixin/timeout/timeouttest.InMemory",
			"iface/mixin/ttl/ttltest.InMemory",
			"iface/mixin/ttlperwrite/ttlperwritetest.InMemory",
			"iface/mixin/windowed/windowedtest.InMemory",
		},
		Controls: []string{"iface/mixin/ttl/deadline.Store"},
	},
	{
		Name:           "transactional",
		Representation: "staged writes beside committed ones, a snapshot, or a history",
		Subjects: []string{
			"iface/contract/transaction/transactiontest.InMemory",
			"iface/contract/tx/txtest.InMemory",
			"iface/mixin/pointintime/pointintimetest.InMemory",
			"iface/mixin/serializable/serializabletest.InMemory",
			"iface/mixin/snapshotisolation/snapshotisolationtest.InMemory",
		},
		Controls: []string{"iface/contract/tx/writelog.Store"},
	},
	{
		Name:           "leases",
		Representation: "the holders of keys, which a release or the end of a context frees",
		Subjects: []string{
			"iface/composite/leased-idempotent-writer/leasedidempotentwritertest.InMemory",
			"iface/contract/leader-election/leaderelectiontest.InMemory",
			"iface/contract/lease/leasetest.InMemory",
			"iface/mixin/waits/waitstest.InMemory",
		},
		Controls: []string{"iface/contract/lease/sortedkeys.Table"},
	},
}
