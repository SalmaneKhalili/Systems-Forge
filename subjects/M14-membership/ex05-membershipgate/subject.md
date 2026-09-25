# M14-ex05 · The membership gateway

The heartbeat countdown, suspect lifecycle, gossip merge, and staleness bound now
meet at the module's TCP gate. The gateway must expose one process-wide view in
which commands register, suspect, revive, evict, and report members without any
wall-clock dependency. You deliver `watch.go`, including concurrent service and
clean shutdown under `SIGTERM`.

## Shape

Whole-program. You write **`watch.go`**, a TCP server on `TARGETPORT` that keeps
**one shared membership view for the process lifetime**. Each member's lifecycle
state is `alive`, `suspect`, or `dead`, and the gateway answers these line
commands:

```text
list           -> the names of live + suspected members joined by "|" (or empty line)
status <name>  -> the member's state ("alive" | "suspect" | "dead")
beat <name>    -> mark the member alive (register or refresh); reply "ok"
suspect <name> -> move an alive member to "suspect"; reply "ok"
drop <name>    -> evict a member to "dead"; reply "ok"
```

Only `alive` and `suspect` members appear in `list`, sorted by name and joined
by `|`. A `dead` member is evicted; a suspect member remains listed, and a fresh
`beat` returns it to `alive`. Use only the Go standard library. `make all` must
build `watch`; the grader starts `./watch`. Serve multiple concurrent
connections, preserve membership across them, and never read the wall clock or
add sleeps or timestamps to replies. The port number comes from `TARGETPORT`.

## Acceptance

`make all` must build `watch`; the grader starts `./watch`, drives one
connection, then uses a second connection against the same membership:

```text
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

Listing a `dead` member fails eviction, and failing to revive a `suspect` member
on `beat` breaks the lifecycle. A gateway that dies on `SIGTERM` with a non-zero
exit also fails the graceful-shutdown step.

## Readings

- **Reading ladder** — the gate is a TCP view of the lifecycle you trained; wiring is M8.
- SWIM paper, protocol states: https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Go `net` / `bufio` docs: https://pkg.go.dev/net
