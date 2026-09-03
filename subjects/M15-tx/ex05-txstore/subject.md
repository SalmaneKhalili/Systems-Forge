# M15-ex05 · txstore gateway

## Goal

Expose the transactional store you built in ex01–ex04 over TCP. Write `txstore.go` — a
line-based gateway: each connection speaks a small command language. Committed writes are
shared and durable across connections; a transaction's staged writes are visible only to the
connection that staged them, until `commit`.

`make all` must produce `./txstore`. `forge` starts the server, opens sockets on
`TARGETPORT`, and checks the protocol transcript.

## Protocol

One command per line, CRLF-free; each reply is one line:

| Command | Reply | Meaning |
|---|---|---|
| `set <k> <v>` | `ok` | store `k=v` (staged if in a txn, else committed) |
| `get <k>` | `<v>` or `nil` | committed value (own staged write wins inside a txn) |
| `begin` | `ok` | start a transaction |
| `commit` | `ok` | apply all staged writes atomically |
| `rollback` | `ok` | discard all staged writes; abort the txn |
| `status <k>` | `alive` / `dead` | is `k` in the committed live set |
| `drop <k>` | `ok` | remove `k` from the live set |
| `list` | keys, one per line | the committed live set, sorted; blank line if empty |

## Constraints

- Go, standard library only; file is `txstore.go` (single `main` package).
- Bind to `os.Getenv("TARGETPORT")` (fallback `17440`); `net.Listen("tcp", ":"+port)`.
- One goroutine per connection; committed state shared behind a mutex.
- A connection's staged writes must never be visible to another connection before `commit`.
- `rollback` must discard staged writes without touching the committed store.
- No wall clock beyond the grader's socket timeout — the protocol is fully deterministic.

## Acceptance

The grader's two connections must observe this exact transcript:

- **conn 1:** `status a` → `dead`; `set a 1`; `get a` → `1`; `begin`; `set a 2`;
  `get a` → `2` (staged, own txn); `status a` → `alive`; `commit`; `get a` → `2`;
  `list` → `a`
- **conn 2** (same server, fresh socket): `list` → `a`; `drop a`; `status a` → `dead`;
  `list` → (blank); `begin`; `set b 3`; `rollback`; `get b` → `nil`

The grader then sends **SIGTERM** to the process and requires it to exit **0**.
This re-deploys the M7-ex04 graceful-shutdown skill: a server should stop
accepting and exit cleanly, not be killed with a non-zero code.

Staged writes leaking to another connection before `commit`, or a `rollback` that applies
its writes anyway, is the bug this gate exists to catch. A gateway that dies on SIGTERM
(non-zero exit) fails the graceful-shutdown step.

## Readings

- Revisit ex01 (transaction staging), ex02 (replay/durability), ex03 (conflict).
- The Linux Programming Interface: sockets — §56 (accept loop), framing as in M6/M16.