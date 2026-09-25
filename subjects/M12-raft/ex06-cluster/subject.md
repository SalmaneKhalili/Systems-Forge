# M12-ex06 · Raft cluster

M12 ends with the pure logic from ex01–ex05 running as a real **3-node raft
cluster**. A single binary boots three distinct nodes with separate TCP
endpoints; they exchange RequestVote and AppendEntries over real sockets until
the leader's log is replicated across the cluster.

## Shape

The binary contains **three distinct nodes (a, b, c)**, each with its own TCP
endpoint. `make all` must produce `./cluster`. `forge` starts the binary and
drives it with explicit commands; there is no wall clock, and the transcript
is byte-deterministic.

The client talks to the whole cluster on `TARGETPORT`. The three nodes listen
on `TARGETPORT+1`, `+2` and `+3`, and message one another over those sockets:

| Command | Reply | Meaning |
|---|---|---|
| `elect <node>` | `ELECTED …` / `LOST …` | `<node>` runs an election at term clusterMax+1; wins on a majority |
| `propose <cmd>` | `ok <idx>` (or `not leader`) | the leader appends `<cmd>` and replicates it to the followers |
| `read` | committed log, `|`-joined | the current leader's log |
| `log <node>` | `<node>`'s full log | **the replication proof** — every node's log must match |
| `term <node>` | `term <t>` | `<node>`'s current election term |

Use Go and the standard library only in `cluster.go`. There must be exactly 3
peers, and `make all` must build `./cluster`. An election needs a strict
majority: 2 of 3 votes, from the candidate and one peer. Every message between
nodes must cross a real TCP connection. Every step is driven by the grading
client, never a wall clock. A node that observes a higher term steps down to
follower immediately.

## Acceptance

This transcript must hold on conn 1 and then a fresh conn 2:

```text
read        ->  (empty: no leader yet)
elect a     ->  ELECTED a term 1
propose x   ->  ok 0
propose y   ->  ok 1
read        ->  x|y
log a       ->  x|y
log b       ->  x|y
log c       ->  x|y
-- conn 2 --
read        ->  x|y
term a      ->  term 1
log c       ->  x|y
```

`log b` and `log c` must match the leader's `read`. That equality is the
replication proof: entries proposed to the leader really reached the followers
over TCP.

## Readings

- **Reading ladder** — start from the Raft paper's log-replication core, then
  revisit the term / vote / quorum exercises below.
- Revisit ex01 (term), ex02 (up-to-date vote rule), ex04 (majority quorum).
- The module intro: Raft keeps one consistent log across a cluster, safely.
