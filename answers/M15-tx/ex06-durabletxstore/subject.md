# M15-ex06 · durable txstore

## Goal

The **Micro-App** that fuses M15 with M16. ex01–ex05 gave you a transactional
store and, separately, a persistence engine. Here they meet in one artifact: a
**durable transactional KV gateway**. Committed writes are atomically applied
*and* recorded in a write-ahead log (WAL) — the very record a crash replay
would restore. This is the storage engine you'd actually hand to an
application.

`make all` must produce `./durable`. `forge` starts it once and drives it with
transactions and writes.

## Commands (line protocol, one per line)

| Command | Reply | Meaning |
|---|---|---|
| `set k v` | `ok` | commit `k=v` durably (append to WAL) |
| `get k` | value / `nil` | read committed store (a staged value wins inside a txn) |
| `status k` | `alive` / `dead` | key existence |
| `list` | keys, one per line | committed keys, sorted |
| `begin` | `ok` | open a transaction (staged writes only) |
| `set k v` (in txn) | `ok` | stage a write — **not yet applied** |
| `commit` | `ok` | apply all staged writes atomically, tracking each in the WAL |
| `rollback` | `ok` | discard staged writes; **no WAL side effect** |
| `wal` | WAL, one op per line | the durable commit record, in order |

## Acceptance

Conn 1 (persisted across the session):

```
set a 1   ->  ok
set b 2   ->  ok
begin     ->  ok
set a 3   ->  ok      (staged — not committed yet)
commit    ->  ok
get a     ->  3
list      ->  a
              b
wal       ->  set a=1
              set b=2
              set a=3
```

Conn 2 (a fresh dial, same server — committed state must persist):

```
get a     ->  3
begin     ->  ok
set c 4   ->  ok
rollback  ->  ok
get c     ->  nil
wal       ->  set a=1      (rollback left nothing behind)
              set b=2
              set a=3
```

`wal` is the whole point: committed writes appear there **once, in commit
order**, and rolled-back writes never touch it.

## Constraints

- Go, standard library only; file is `durable.go`.
- The WAL must record every committed `set` as a line `set k=v` in commit order.
- No wall clock — all steps are client-driven.
- `make all` must build `./durable`.

## Readings

- M15-ex01 (staged writes / atomic commit), M15-ex02 (WAL replay).
- M16-ex03 (persistence): a WAL is the durability guarantee; replay restores
  state. Here you can observe the WAL directly via the `wal` command.