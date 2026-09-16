# M14-ex05 · The membership gateway

## Goal

The capstone of the module: a **membership gateway** over TCP that reports and
mutates the cluster's live set. Each member has a lifecycle state — `alive`,
`suspect`, or `dead` — and the gateway answers these line commands:

```
list           -> the names of live + suspected members joined by "|" (or empty line)
status <name>  -> the member's state ("alive" | "suspect" | "dead")
beat <name>    -> mark the member alive (register or refresh); reply "ok"
suspect <name> -> move an alive member to "suspect"; reply "ok"
drop <name>    -> evict a member to "dead"; reply "ok"
```

Only `alive`/`suspect` members appear in `list`; a `dead` member is evicted.
A `suspect` member is still listed but flagged; a fresh `beat` returns it to
`alive`.

Write `watch.go`: a TCP server on `TARGETPORT` that keeps **one shared
membership view for the process lifetime**. `make all` must build `watch`; the
grader starts `./watch`.

## Constraints

- Go, standard library only; file is `watch.go`.
- `make all` must build `watch` (the grader starts `./watch`).
- Multiple concurrent connections must be served; membership persists across them.
- `list` returns member names sorted, joined by `|`.
- Never read the wall clock; no sleeps, no timestamps in replies.
- Port number comes from `TARGETPORT`.

## Acceptance

The grader drives one connection, then a second shares the same membership:

```
conn1: status a -> dead          (never seen)
conn1: beat a / beat b / beat c  -> ok ok ok
conn1: status a -> alive
conn1: list     -> a|b|c
conn1: suspect b -> ok
conn1: status b -> suspect
conn1: status a -> alive
conn1: list     -> a|b|c         (b still listed)
conn2: list     -> a|b|c         (persists)
conn2: drop c   -> ok
conn2: status c -> dead
conn2: list     -> a|b           (c evicted)
conn2: beat b   -> ok
conn2: status b -> alive         (revived from suspect)

The grader then sends **SIGTERM** to the process and requires it to exit **0**.
This re-deploys the M7-ex04 graceful-shutdown skill: a server should stop
accepting and exit cleanly, not be killed with a non-zero code.
```

A gateway that lists a `dead` member, or that fails to revive a `suspect`
member on `beat`, is the bug. A gateway that dies on SIGTERM (non-zero
exit) fails the graceful-shutdown step.

## Readings

- **Reading ladder** — the gate is a TCP view of the lifecycle you trained; wiring is M8.
- SWIM paper, protocol states: https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Go `net` / `bufio` docs: https://pkg.go.dev/net
