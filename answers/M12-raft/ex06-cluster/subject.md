# M12-ex06 · Raft cluster

## Goal

Put it all together: a **real 3-node raft cluster**. Ex01–ex05 taught the
mechanics as pure logic (term, vote rule, log-match, quorum). Here a single
binary boots **three distinct nodes (a, b, c)**, each with its own TCP
endpoint, and they replicate a log the real distributed way — nodes exchange
RequestVote and AppendEntries **over real sockets**, the leader commits on a
majority, and every node ends with the same log.

`make all` must produce `./cluster`. `forge` starts it and drives the cluster
with explicit commands (no wall clock — the transcript is byte-deterministic).

## Cluster layout

- The client talks to the whole cluster on `TARGETPORT`.
- The three nodes listen on `TARGETPORT+1`, `+2`, `+3` and message each other
  over those sockets.

## Commands (to the client socket)

| Command | Reply | Meaning |
|---|---|---|
| `elect <node>` | `ELECTED …` / `LOST …` | `<node>` runs an election at term clusterMax+1; wins on a majority |
| `propose <cmd>` | `ok <idx>` (or `not leader`) | the leader appends `<cmd>` and replicates it to the followers |
| `read` | committed log, `|`-joined | the current leader's log |
| `log <node>` | `<node>`'s full log | **the replication proof** — every node's log must match |
| `term <node>` | `term <t>` | `<node>`'s current election term |

## Constraints

- Go, standard library only; file is `cluster.go`.
- Exactly 3 peers; `make all` must build `./cluster`.
- Election needs a strict majority: 2 of 3 votes (self + one peer).
- All message passing between nodes must go over a real TCP connection.
- No wall clock: every step is driven by the grading client.
- A node that observes a higher term steps down to follower immediately.

## Acceptance

This transcript must hold (conn 1 then a fresh conn 2):

```
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

`log b` and `log c` matching the leader's `read` is the whole point: entries
proposed to the leader are actually replicated to the followers over TCP.

## Readings

- Revisit ex01 (term), ex02 (up-to-date vote rule), ex04 (majority quorum).
- The module intro: Raft keeps one consistent log across a cluster, safely.