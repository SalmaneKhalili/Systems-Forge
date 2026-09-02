# M13-ex05 · The shard gateway

## Goal

The capstone of the module: a **shard gateway** over TCP that routes each key
to the shard that owns it and keeps each shard's data isolated. Routing is by
key range with two sorted boundaries — `m` and `t` — giving three shards:
`0`=`(-∞, m]`, `1`=`(m, t]`, `2`=`(t, +∞)`.

The gateway answers two line commands:

```
set <key> <val>   -> route to the owning shard, store it; reply "set <n>"
get <key>         -> route to the owning shard, reply "<n>:<val>" (or "<n>:?" if absent)
```

Each shard is an independent key→value store; a key is always routed to the
same shard, so a value set on one connection is visible from another.

Write `gate.go`: a TCP server on `TARGETPORT` that keeps **one shared set of
shards for the process lifetime**. `make all` must build `gate`; the grader
starts `./gate`.

## Constraints

- Go, standard library only; file is `gate.go`.
- `make all` must build `gate` (the grader starts `./gate`).
- Multiple concurrent connections must be served; shard state persists across them.
- Routing is by key range: shard 0 if `key <= "m"`, shard 1 if `key <= "t"`,
  else shard 2.
- Never read the wall clock; no sleeps, no timestamps in replies.
- Port number comes from `TARGETPORT`.

## Acceptance

The grader drives one connection, then a second shares the same shards:

```
conn1: set apple red     -> set 0        (apple <= m)
conn1: get apple         -> 0:red
conn1: set mango yellow  -> set 1        (m < mango <= t)
conn1: get mango         -> 1:yellow
conn2: get apple         -> 0:red        (persists)
conn2: set zebra striped -> set 2        (zebra > t)
conn2: get zebra         -> 2:striped
conn2: get melon         -> 1:?          (melon routes to shard 1; absent)
```

`melon` routes to shard 1 (m < melon < t) but was never set, so it reads
`1:?`. A gateway that routes a key to the wrong shard, or that returns another
shard's value, is the bug.