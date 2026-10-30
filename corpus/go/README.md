<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# The Go corpus

The corpus contains a fixture for each classification that the suite generator handles. `dokimi run` writes the suites of the fixtures, and `go test ./...` applies each suite to its subjects.

## Subjects and controls

A subject is a type whose doc comment has the directive `subject`. The generator writes a test that applies the suite of its declaration to it with the option `isolate`. Each fixture has one subject, its `InMemory` in most fixtures.

A control implements the declaration of another subject with another representation of its state. When a check fails on a control, the check rejects behaviour that the declaration permits. Each application also fails a check that never engages its claim, as vacuous, unless the subject excuses the check with a reason.

## Shape families

A shape family is a group of subjects whose state has one representation. `families.go` declares the family of each subject, and its spec fails for a subject outside every family and for a family without a control.

| Family | Representation of its subjects | Subjects | Control | Representation of the control |
|---|---|---|---|---|
| stores-and-readers | A map keyed by the input, with a counter or a revision beside it, or a single value | 63 | `iface/mixin/readafterwrite/sorted` | A sorted slice of keys, searched by binary search |
| appenders-and-logs | A slice that grows at its end, with a position beside it | 20 | `iface/mixin/stableorder/chunked` | A list of blocks of four values |
| caches | A map in front of a backing store | 5 | `iface/contract/cache/oneentry` | A cache of one entry that each miss replaces |
| queues | A channel, or a slice that the oldest element leaves first | 4 | `iface/door/queue/ring` | A ring buffer of two items under a mutex |
| publishers | A slice of the channels or the functions of the subscribers | 8 | `iface/contract/publisher/mailbox` | An unbounded mailbox and a delivering goroutine per subscriber |
| flags-and-lifecycles | A flag or a small state machine | 12 | `iface/detector/lifecycle/atomicstate` | An atomic integer without a mutex |
| counters | One integer | 8 | `iface/mixin/commutative/perkey` | A total for each key, summed on each read |
| stateless | No state, or a value that the constructor sets | 11 | `iface/mixin/xsssafe/htmlescape` | `html.EscapeString` of the standard library |
| clocked | Instants of a clock beside the values | 7 | `iface/mixin/ttl/deadline` | The deadline of each value in place of the instant of its write |
| transactional | Staged writes beside committed ones, a snapshot, or a history | 5 | `iface/contract/tx/writelog` | One log of every write, read through the order of the commits |
| leases | The holders of keys, which a release or the end of a context frees | 4 | `iface/contract/lease/sortedkeys` | A sorted slice of the held keys |

## Excusals

| Subject | Check | Reason |
|---|---|---|
| `iface/composite/causal-chain/replica` | `Log.Append/chain`, `Log.Append/chain-causal`, `Log.Replay/keeps-the-body`, `Log.Replay/keeps-every-body` | The replica refuses each append, so its replay never grows. |

## Seeded subjects

A subject whose checks read a value before they write one has a value at construction, so each check engages its claim:

- `iface/composite/paginated-reader` starts with three values.
- `iface/mixin/monotonicreads` and `iface/mixin/pointintime` start with a value under the key `Key-a`, which their checks read first.
