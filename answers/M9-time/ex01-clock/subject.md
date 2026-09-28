# M9-ex01 · Lamport logical clock

## Goal

Implement a Lamport logical clock in `clock.go`. The provided `main.go`
plays a fixed, deterministic script of local events, sends, and receives on
two clocks and prints the transcript. Your clock must produce exactly the
reference transcript.

Complete a type `Lamport` with three methods (add any private fields you
need):

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

The receive rule is the heart of Lamport's clock: a message tells the
receiver "the sender was at least at `other`", so the receiver's counter
must jump past it before counting the receive itself.

## Constraints

- Go, standard library only; file is `clock.go`.
- `make all` must build a `test` binary and `./test` must print the exact
  reference transcript (`expected.txt` is authoritative; whitespace is
  normalized).
- Never read the wall clock. No time, no sleeps.
- quiz.txt answered.

## Acceptance

Run `python3 - <<'EOF'`-style check of anyone with a broken clock — transcripts
will differ. The reference prints:

```
e0 t1
p0 send m0 t2
p1 recv m0 t3
e1 t4
e2 t3
p1 send m1 t5
p0 recv m1 t6
e3 t7
```

`e2 t3` is the tell: p0's local events proceed independently of p1, and when
p0 finally receives m1 (stamp 5) it must jump to 6, never stay at 3.

## Readings

- **Reading ladder** — the paper is short and original; read the idea first, then
  Kleppmann's modern treatment, then the Go you already know.
- Lamport, "Time, Clocks, and the Ordering of Events in a Distributed System" (1978):
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
- Martin Kleppmann, *Distributed Systems* lecture notes — logical clocks:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Go `sync` (mutex around the counter) and `time` — only as needed to parallel the trace.
