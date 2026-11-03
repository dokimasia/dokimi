---
rfc: 0033
title: The isolation levels of transactions
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-02
updated: 2026-11-03
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0033: The isolation levels of transactions

## Summary

A declaration with the contract `tx` and the mixin `serializable` or `snapshot-isolation` gets the claim of that isolation level. The generated check runs whole transactions on three clients at once: a begin, a few writes and reads in the scope of the transaction, and a commit or a rollback. It records each transaction as one call of the operation `txn` in the history of assert-go, and checks the history with `history.Serializable` or `history.HasSnapshotIsolation`.

## Motivation

A declaration with the mixin `serializable` or `snapshot-isolation` states the isolation level of its concurrent transactions. Both mixins have the stance `needs` today, so the suite generator does not check either level.

The checkers of assert-go read each transaction of a history as one call of the operation `txn`, in the model of Elle. Each argument of the call is a micro-operation, such as the append of a value to a list, the write of a register, or a read of either. Each step of the generated sequences calls one callable, and the contract `tx` states only begin, commit and rollback. A history of a suite therefore has no call of `txn`.

## Detailed design

### The claims

| Claim | Licence | Machineries | Probe | Reason | Defect form |
|---|---|---|---|---|---|
| `serializable` | The mixin `serializable` with its param `read` | `concurrent-callers`, `generated-sequences` | `serializable-history` | `not-serializable` | `reads-uncommitted` |
| `snapshot-isolation` | The mixin `snapshot-isolation` with its param `read` | `concurrent-callers`, `generated-sequences` | `snapshot-history` | `not-snapshot-isolated` | `writes-without-a-snapshot` |

Both claims have the scope `declaration`. Each takes its licence from the callable with its mixin and binds once for the declaration of the role `begin`, as `linearizable` binds once for its declaration. Both list the gate `transactions`.

`history.HasSnapshotIsolation` fails a history for each anomaly that snapshot isolation excludes, such as a dirty read or a lost update. The reason `not-snapshot-isolated` covers each of those anomalies, as `not-serializable` covers each anomaly of a history that is not serializable.

### The scope of a transaction

The callables of one transaction are the roles of an instance of the contract `tx` and the two callables of the mixin:

| Part | Callable |
|---|---|
| Begin | The role `begin`. It returns the handle of the transaction |
| Write | The callable with the mixin. It takes the handle, a key and a value |
| Read | The callable that the param `read` of the mixin names. It takes the handle and the key, and returns the value of the key |
| Commit and rollback | The roles `commit` and `rollback`. Each takes the handle |

A read that returns a list of the type of the written value makes a history of list appends, and the write appends its value to the list of the key. Any other read makes a history of registers, and the write sets the value of the key.

A begin can also return a handle whose own methods write, read, commit and roll back, as `database/sql` returns a `*sql.Tx`. eidos joins the callables of one package that have a role in a contract into one instance, also when they are methods of two types. So the roles `commit` and `rollback` and the mixin can be on the methods of the handle type, and the param `read` names a method of the same type first. The claim then binds for the declaration of the role `begin`, and each call of the transaction calls the method of the handle:

```go
type Store interface {
	//dokimi:contract tx role=begin closed=ErrClosed
	Begin(ctx context.Context) (Tx, error)
}

type Tx interface {
	//dokimi:contract tx role=commit
	Commit(ctx context.Context) error

	//dokimi:contract tx role=rollback
	Rollback(ctx context.Context) error

	//dokimi:mixin serializable read=Get
	Append(ctx context.Context, key string, v int) error

	Get(ctx context.Context, key string) ([]int, error)
}
```

A handle type of another package, such as `*sql.Tx`, cannot carry a directive. The gate `transactions` refuses it, and its remedy is an interface of the package that the handle implements.

```go
type Store interface {
	//dokimi:contract tx role=begin closed=ErrClosed
	Begin(ctx context.Context) (Tx, error)

	//dokimi:contract tx role=commit
	Commit(ctx context.Context, tx Tx) error

	//dokimi:contract tx role=rollback
	Rollback(ctx context.Context, tx Tx) error

	//dokimi:mixin serializable read=Get
	Append(ctx context.Context, tx Tx, key string, v int) error

	Get(ctx context.Context, tx Tx, key string) ([]int, error)
}
```

### The action of a transaction

The actions of the declaration gain one action, `txn`. Its input takes one to four micro-operations over the pool of keys from the case, and whether the transaction commits. Most transactions commit. Each write takes a value that no other write of the case takes: the next value of a counter of the case, converted to the type of the value.

The action begins the transaction, calls the write or the read of each micro-operation with the handle, and then commits or rolls back. It records one call of `txn` whose arguments are the micro-operations, such as `[]any{"append", "a", 3}` and `[]any{"read", "a", nil}`:

| End of the transaction | Completion of the call |
|---|---|
| The commit returns no error | `OK`, with each read and its list or value |
| The rollback runs, or a call of the transaction fails | `Fail` |
| The context of the case ends | `Unknown` |

### The gate transactions

```yaml
transactions:
  passes: The roles begin, commit and rollback of one instance of the contract tx are in the package of the declaration, the write takes the handle that begin returns or is a method of it, the write takes a key and a value of an integer or a string type, and the read takes the handle or is a method of it and takes the key.
  why: "{callable} and the read of its mixin do not run in the scope of a transaction of {declaration}"
  remedy: "give {callable} and its read the handle that begin returns, or declare them on an interface of the package that the handle implements, with a key and a value of an integer or a string type"
  since: 1.2.0
```

### The probes

```go
// SerializableHistory runs generated sequences of the actions on the
// subject of h, with a concurrent section on the clients 0 to clients that
// the schedule of h runs. After the last step it checks the transactions of
// the history of c with history.Serializable. It engages h when a
// transaction commits.
func SerializableHistory[S any](h *conform.Handle[S], c *prop.Case, actions Actions[S], clients int)

// SnapshotHistory runs the sequences of SerializableHistory, and checks the
// transactions of the history of c with history.HasSnapshotIsolation.
func SnapshotHistory[S any](h *conform.Handle[S], c *prop.Case, actions Actions[S], clients int)
```

Both checkers decide every history, so neither probe has a budget. Each probe is vacuous where no transaction of a run commits.

### The planted defects

| Form | Changes | What the defect does |
|---|---|---|
| `reads-uncommitted` | declaration | The write keeps its value in a field of the defect until its transaction commits or rolls back. A read returns the value that the latest open transaction kept for the key, also when another transaction kept it |
| `writes-without-a-snapshot` | declaration | The write keeps its value in a field of the defect, and a read returns a value that its own transaction kept. The commit applies the kept writes in a transaction of its own, which begins at the commit, so two transactions that write one key both commit |

Called one transaction at a time, each defect returns what the subject returns. Each races through `Handle.Yield` after each call of a transaction, so the scheduler of a planted defect interleaves the transactions.

## Alternatives considered

### A history that the subject reports

The serializable and snapshot fixtures of the corpus give the declaration a method that returns the transactions that the subject ran, and a check would read them.

**Why not:** the subject then reports its own history. A subject that records a serial history and runs another passes, and every subject needs a method that only a test calls.

### List appends alone

Elle reads list appends, because each read of a list reveals the order of the appends before it.

**Why not:** a store of single values has no list to read. A history of registers shows fewer anomalies, and it still finds a dirty read and a lost update.

## Drawbacks

- **Value types.** A unique value of each write needs an integer or a string type. A declaration whose value has another type refuses both claims with the gate `transactions`.
- **A read outside the scope.** A read that takes no handle, such as the read of the committed state of the `tx` fixture of the corpus, refuses both claims.
- **A handle of another package.** A begin that returns `*sql.Tx` refuses both claims until the consumer declares an interface of the handle.
- **Larger histories.** Each step runs up to six calls of the subject, so a case of three clients records more calls than a case of single calls.

## Open questions

None.

## References

| What | Where |
|---|---|
| Elle, inferring isolation anomalies from histories | <https://arxiv.org/abs/2003.10554> |
| assert-go, the package history | <https://github.com/dokimasia/assert-go/tree/main/history> |
| The mixins `serializable` and `snapshot-isolation` of the shape catalog | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec/mixins> |
