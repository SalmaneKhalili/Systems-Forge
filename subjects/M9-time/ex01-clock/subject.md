# M9-ex01 · Lamport logical clock

M9 starts with the scalar clock that gives every later ordering mechanism a
causal baseline. The provided `main.go` runs a fixed script of local events,
sends and receives across two clocks; you implement the counter so its receive
rule preserves the send's position in happens-before.

## Shape

The exercise ships `main.go` with a fixed, deterministic local-event/send/
receive script. You write **`clock.go`** in the same package and complete a type
`Lamport` with three methods, adding any private fields you need:

```go
// Tick counts a local event (including a send) and returns the new value.
//   c = c + 1
func (l *Lamport) Tick() int

// Now returns the current value without advancing.
func (l *Lamport) Now() int

// Add reflects a received event that carries the stamp other:
//   c = max(c, other) + 1
func (l *Lamport) Add(other int) int
```

The receive rule is the core operation: a message tells the receiver that the
sender was at least at `other`, so the receiver jumps past that stamp before
counting the receive itself.

Use Go and the standard library only. `make all` must build a `test` binary,
and `./test` must print the exact reference transcript; `expected.txt` is
authoritative and whitespace is normalized. Never read the wall clock: no time,
no sleeps. `quiz.txt` must be answered.

## Acceptance

The `python3 - <<'EOF'`-style check exposes a broken clock through a differing
transcript. The reference prints exactly:

```text
e0 t1
p0 send m0 t2
p1 recv m0 t3
e1 t4
e2 t3
p1 send m1 t5
p0 recv m1 t6
e3 t7
```

- `e2 t3` proves p0's local events advance independently of p1; the final
  receive of m1 must jump from 3 to 6, never leave the counter at 3.
- A counter that does not fold the received stamp into its own produces a
  different line and fails the transcript.

## Readings

- **Reading ladder** — the paper is short and original; read the idea first, then
  Kleppmann's modern treatment, then the Go you already know.
- Lamport, "Time, Clocks, and the Ordering of Events in a Distributed System" (1978):
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
- Martin Kleppmann, *Distributed Systems* lecture notes — logical clocks:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Go `sync` (mutex around the counter) and `time` — only as needed to parallel the trace.
