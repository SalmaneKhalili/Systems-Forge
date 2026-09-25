# M12-ex05 · The gate node

The gate composes ex01–ex04 into one networked Raft node. Its monotonic term
and committed log share process lifetime across clients, giving the three-node
cluster in ex06 a concrete node contract.

## Shape

A **raft gate node** over TCP holds a monotonic term and a committed log, and
answers three line commands:

```text
see <term>     -> observe a term (never moves your term down); reply "term <t>"
append <cmd>   -> as leader, append and commit (single-node quorum); reply "ok <idx>"
read           -> reply the committed entries joined by "|", or empty line
```

A fresh node is a follower at term 0 with an empty log, so an append before any
term has been seen is refused with `not leader`. Once it has seen a term, it
is the single-node leader: every append commits immediately, and `read` never
leaks uncommitted entries because none exist.

You write **`node.go`**, a TCP server on `TARGETPORT` that keeps **one node
state for the process lifetime**, with the term and committed log shared
across all connections. Use Go and the standard library only. `make all` must
build `node`; the grader starts `./node`. An accept loop must serve multiple
concurrent connections, with term and log persisting across them. The term is
strictly monotonic, so `see 2` at term 3 remains `term 3`. Appending before
any term is refused; appends after commit land at the end. Never read the wall
clock, and put no sleeps or timestamps in replies. The port comes from
`TARGETPORT`.

## Acceptance

The grader drives one connection, then a second uses the same node state:

```text
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

A node whose term drops, that lets stale `see` regress it, or whose `read`
returns anything but the committed prefix fails the gate.

## Readings

- **Reading ladder** — the gate ties term, vote, log-match and commit into one node;
  §5.1–§5.4 in a single read.
- Raft paper, §5.1–§5.4: https://raft.github.io/raft.pdf
- Go `net` / `bufio` docs: https://pkg.go.dev/net
