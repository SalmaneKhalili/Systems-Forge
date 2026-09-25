# M15-ex05 · txstore gateway

The transaction, replay, conflict, and drill exercises now need a network
boundary. This gateway keeps staged writes private to each connection while
sharing committed, durable state across sockets, then proves the isolation
boundary and rollback path over its line protocol. You deliver `txstore.go`, the
TCP form of the store built in ex01–ex04.

## Shape

Whole-program. You write **`txstore.go`** as a single `main` package using only
the Go standard library. One command arrives per line without CRLF, and each
reply occupies one line:

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

Bind to `os.Getenv("TARGETPORT")` with fallback `17440`, using
`net.Listen("tcp", ":"+port)`. Run one goroutine per connection and protect
committed state with a mutex. A connection's staged writes must never be visible
to another connection before `commit`, and `rollback` must discard them without
touching the committed store. The protocol is fully deterministic; the only
wall-clock use permitted is the grader's socket timeout. `make all` must produce
`./txstore`, and `forge` starts the server and checks the protocol transcript.

## Acceptance

`make all` must produce `./txstore`; `forge` starts the server, opens sockets on
`TARGETPORT`, and requires this exact two-connection transcript:

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

- **Reading ladder** — start with how a WAL makes commits durable (replay),
  then revisit the staging exercises below.
- Revisit ex01 (transaction staging), ex02 (replay/durability), ex03 (conflict).
- The Linux Programming Interface: sockets — §56 (accept loop), framing as in M6/M16.
