<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# The Go corpus

The corpus contains a fixture for each classification that the suite generator handles. `dokimi run` writes the suites of the fixtures, and `go test ./...` applies each suite to its subjects.

## Expected outcomes

`corpus/expect` contains one file of expected outcomes for each fixture, and `overlay.yaml` declares the differences of Go from those outcomes. Go is the reference corpus, so its overlay declares none.

`dokimi conform --lang go` checks that the generated files are current and runs the tests without the cache. It fails on each difference of the lock and of the suite records from the expected outcomes. After a change of the generator, `dokimi conform --lang go --update` rewrites the declarations of each file from the run. Review the diff of each file as a change of the contract of its fixture.

## Subjects and controls

A subject is a type whose doc comment has the directive `subject`. The generator writes a test that applies the suite of its declaration to it with the option `isolate`. Most fixtures have one subject, their `InMemory`. `iface/door/reference` has three, `iface/composite/causal-chain` and `iface/lang/suitedhandle` have two, and 23 fixtures have none, such as those of enumerations, errors and structs.

A control implements the declaration of another subject with another representation of its state. When a check fails on a control, the check rejects behaviour that the declaration permits. Each application also fails a check that never engages its claim, as vacuous, unless the subject excuses the check with a reason.

## Shape families

A shape family is a group of subjects whose state has one representation. `families.go` declares the family of each subject, and its spec fails for a subject outside every family and for a family without a control.

| Family | Representation of its subjects | Subjects | Control | Representation of the control |
|---|---|---|---|---|
| stores-and-readers | A map keyed by the input, with a counter or a revision beside it, or a single value | 69 | `iface/mixin/readafterwrite/sorted` | A sorted slice of keys, searched by binary search |
| appenders-and-logs | A slice that grows at its end, with a position beside it | 21 | `iface/mixin/stableorder/chunked` | A list of blocks of four values |
| caches | A map in front of a backing store | 7 | `iface/contract/cache/oneentry` | A cache of one entry that each miss replaces |
| queues | A channel, or a slice that the oldest element leaves first | 5 | `iface/door/queue/ring` | A ring buffer of two items under a mutex |
| publishers | A slice of the channels or the functions of the subscribers | 8 | `iface/contract/publisher/mailbox` | An unbounded mailbox and a delivering goroutine per subscriber |
| flags-and-lifecycles | A flag or a small state machine | 14 | `iface/detector/lifecycle/atomicstate` | An atomic integer without a mutex |
| counters | One integer | 8 | `iface/mixin/commutative/perkey` | A total for each key, summed on each read |
| stateless | No state, or a value that the constructor sets | 15 | `iface/mixin/xsssafe/htmlescape` | `html.EscapeString` of the standard library |
| clocked | Instants of a clock beside the values | 9 | `iface/mixin/ttl/deadline` | The deadline of each value in place of the instant of its write |
| transactional | Staged writes beside committed ones, or a snapshot | 6 | `iface/contract/tx/writelog` | One log of every write, read through the order of the commits |
| leases | The holders of keys, which a release or the end of a context frees | 5 | `iface/contract/lease/sortedkeys` | A sorted slice of the held keys |

## Excusals

| Subject | Check | Reason |
|---|---|---|
| `iface/composite/causal-chain/replica` | `Log.Append/chain`, `Log.Append/chain-causal`, `Log.Append/replays-in-order-over-sequences`, `Log.Replay/keeps-the-body`, `Log.Replay/keeps-every-body` | The replica refuses each append, so its replay never grows. |

## Seeded subjects

A subject whose checks read a value before they write one has a value at construction, so each check engages its claim:

- `iface/composite/paginated-reader` starts with three values.
- `iface/mixin/monotonicreads` and `iface/mixin/pointintime` start with a value under the key `Key-a`, which their checks read first.
