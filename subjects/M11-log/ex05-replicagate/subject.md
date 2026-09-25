# M11-ex05 · The replica

M11 ends with the log mechanics behind a shared TCP boundary. The **replica**
serves one process-wide committed log to many clients, so a connection can
append or commit without learning entries a later leader change may roll back.

## Shape

A **replica** serves one consistent committed log to many clients. It holds a
log with a commit index like ex02 and answers three line commands over TCP:

```text
append <cmd>   -> append at the end; reply "ok <index>"
commit <idx>   -> raise the commit index; reply "commit <idx>"
read           -> reply the committed entries joined by "|", or empty line
```

Entries beyond the commit index are **never** returned by `read`; a client
would otherwise see a value a later leader change could roll back.

You write **`replica.go`**, a TCP server on `TARGETPORT` that keeps **one log
for the process lifetime** and shares that state across all connections. Use
Go and the standard library only. `make all` must build `replica`; the grader
starts `./replica`. An accept loop must serve multiple concurrent connections,
and state persists across them. `read` must never reveal an entry beyond the
commit index. Never read the wall clock; replies contain no sleeps or
timestamps. The port comes from `TARGETPORT`.

## Acceptance

The grader appends and commits on conn 1, then opens conn 2:

```text
conn1: append x -> ok 0
conn1: append y -> ok 1
conn1: commit 1 -> commit 1
conn1: read     -> x|y
conn2: append z -> ok 2
conn2: read     -> x|y
```

`z` occupies index 2, beyond the still-1 commit index, so conn 2's `read`
returns only `x|y`. Returning `z` before it commits leaks uncommitted state
and fails.

## Readings

- **Reading ladder** — the gate is a TCP replica serving one committed log; the wire part is M8.
- Raft paper, §5.3 (log replication and commit index): https://raft.github.io/raft.pdf
- Go `net` / `bufio` docs: https://pkg.go.dev/net
