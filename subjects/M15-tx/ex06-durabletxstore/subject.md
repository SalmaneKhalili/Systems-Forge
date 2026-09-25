# M15-ex06 · durable txstore

The module's Micro-App fuses M15 with M16: transactional staging and atomic
commit meet a durable write-ahead log in one gateway. Every committed write is
applied atomically and recorded in the WAL that crash replay would consume, while
rollback leaves no durable trace. You deliver `durable.go`, a storage engine an
application can use across client connections and restarts.

## Shape

Whole-program. You write **`durable.go`** using only the Go standard library.
The line protocol accepts one command per line:

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

The WAL records every committed `set` as one `set k=v` line in commit order.
Rollback writes nothing to it, and all steps are client-driven with no wall
clock. `make all` must produce `./durable`; `forge` starts it once and drives it
with transactions and writes.

## Acceptance

`make all` must produce `./durable`; `forge` starts it once and its first
connection must produce exactly:

```text
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

A fresh second connection to the same server must observe the persisted state and
produce exactly:

```text
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

## Readings

- **Reading ladder** — start from the WAL durability core (M16-ex03, ex02
  replay), then the revisit list below.
- M15-ex01 (staged writes / atomic commit), M15-ex02 (WAL replay).
- M16-ex03 (persistence): a WAL is the durability guarantee; replay restores
  state. Here you can observe the WAL directly via the `wal` command.
