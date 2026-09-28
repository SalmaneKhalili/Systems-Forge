# M11-ex05 · The replica

## Goal

A **replica** serves one consistent committed log to many clients. It holds a
log with a commit index (like ex02) and answers three line commands over TCP:

```
append <cmd>   -> append at the end; reply "ok <index>"
commit <idx>   -> raise the commit index; reply "commit <idx>"
read           -> reply the committed entries joined by "|", or empty line
```

Entries beyond the commit index are **never** returned by `read` — a client
would otherwise see a value that a later leader change could roll back.

Write `replica.go`: a TCP server on `TARGETPORT` that keeps **one log for the
process lifetime** (the state is shared across all connections). `make all`
must build `replica`; the grader starts `./replica`.

## Constraints

- Go, standard library only; file is `replica.go`.
- `make all` must build `replica` (the grader starts `./replica`).
- Multiple concurrent connections must be served; state persists across them.
- `read` must never reveal an entry beyond the commit index.
- Never read the wall clock; no sleeps, no timestamps in replies.
- Port number comes from `TARGETPORT`.

## Acceptance

The grader appends then commits on conn 1, then opens conn 2:

```
conn1: append x -> ok 0
conn1: append y -> ok 1
conn1: commit 1 -> commit 1
conn1: read     -> x|y
conn2: append z -> ok 2
conn2: read     -> x|y
```

`z` sits at index 2, beyond the commit index (still 1), so conn 2's `read`
returns only `x|y`. A replica that leaks `z` (returns it uncommitted) is the
bug.

## Readings

- **Reading ladder** — the gate is a TCP replica serving one committed log; the wire part is M8.
- Raft paper, §5.3 (log replication and commit index): https://raft.github.io/raft.pdf
- Go `net` / `bufio` docs: https://pkg.go.dev/net
