# M8-ex03 · The switch: transparent relay

This is the core of the fault-injection **switch**, the centrepiece of the
module: a transparent TCP relay that sits between a client and a service. It
binds `TARGETPORT`, must bring up the service itself, and forwards every byte
unchanged in both directions. If a client vanishes mid-relay, the switch must
tear down that relay and keep serving the next connection — the discipline
that ex04 turns into fault injection and ex05 tests under fire.

## Shape

Whole-program. The exercise ships the **backend** (`backend/main.go` — a
provided line server you must not modify, built as `./svc`), `go.mod`, and a
`Makefile` that builds both `./switch` and `./backend`. You write
**`switch.go`** (the whole program, `package main`).

The wire format is **newline-terminated lines**. The backend answers the k-th
frame of each connection with `R<k> <content>` (one reply per frame, per
connection). The switch must not invent, merge, split or reorder frames. The
contract:

1. Launch `./svc` yourself on a `127.0.0.1` port of your choosing, and pass it
   `TARGETPORT` in its environment (your own `TARGETPORT` env tells you where
   the world expects *you*).
2. Accept client connections on `TARGETPORT`.
3. For each client connection, dial a **fresh** backend connection and relay:
   client lines → backend, backend replies → client, both ways at once.
4. Close the whole relay the moment either side disconnects; a closed client
   connection must never corrupt or stall the next one.

Nothing may be printed to stdout; diagnostics may go to stderr.

## Acceptance

`make all` must build `switch` and `backend`. The grader then runs `./switch`
and plays this script over the network:

- connection 1: `ping` → expect `R1 ping`, then `alpha` → `R2 alpha`, then
  `beta` → `R3 beta`.
- connection 2 (fresh): `one` → expect `R1 one`, then `two` → `R2 two`.

The backend numbers per connection: if connection 2 got an `R3`, the relay
reused a stale backend connection — wrong.

Graded `net`.

## Readings

- **Reading ladder** — the failure half of networking: a transparent relay and
  where the teardown rules live.
- Go `net` docs — TCP listener and accept loop: https://pkg.go.dev/net
- Go `bufio` docs — line framing: https://pkg.go.dev/bufio
- Tanenbaum & Wetherall, *Computer Networks* — the proxy / relay chapter.