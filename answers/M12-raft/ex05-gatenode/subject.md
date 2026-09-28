# M12-ex05 · The gate node

## Goal

The capstone of the module: a **raft gate node** over TCP that ties the
safety rules together. It holds a monotonic term and a committed log, and it
answers three line commands:

```
see <term>     -> observe a term (never moves your term down); reply "term <t>"
append <cmd>   -> as leader, append and commit (single-node quorum); reply "ok <idx>"
read           -> reply the committed entries joined by "|", or empty line
```

A fresh node is a follower at term 0 with an empty log: appending before any
term has been seen is refused (`not leader`). Once it has seen a term it is
the single-node leader, so every append is immediately committed and a `read`
never leaks uncommitted entries (there are none).

Write `node.go`: a TCP server on `TARGETPORT` that keeps **one node state for
the process lifetime** (term + committed log shared across all connections).
`make all` must build `node`; the grader starts `./node`.

## Constraints

- Go, standard library only; file is `node.go`.
- `make all` must build `node` (the grader starts `./node`).
- Multiple concurrent connections must be served; term and log persist across them.
- The term must be strictly monotonic: `see 2` at term 3 stays `term 3`.
- Appending before any term refuses; appends after commit append at the end.
- Never read the wall clock; no sleeps, no timestamps in replies.
- Port number comes from `TARGETPORT`.

## Acceptance

The grader drives one connection, then a second shares the same node state:

```
conn1: append x -> not leader      (still term 0, no leader)
conn1: see 3    -> term 3          (follower steps up to term 3)
conn1: append a -> ok 0
conn1: append b -> ok 1
conn1: read     -> a|b
conn1: see 2    -> term 3          (stale term ignored)
conn2: read     -> a|b             (committed log persists)
conn2: see 9    -> term 9          (higher term steps the node up)
conn2: append c -> ok 2
conn2: read     -> a|b|c
```

A node whose term drops (or that lets a stale `see` regress it), or that
answers a `read` with anything but the committed prefix, is the bug.

## Readings

- **Reading ladder** — the gate ties term, vote, log-match and commit into one node;
  §5.1–§5.4 in a single read.
- Raft paper, §5.1–§5.4: https://raft.github.io/raft.pdf
- Go `net` / `bufio` docs: https://pkg.go.dev/net
