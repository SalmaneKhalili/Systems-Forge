# M8 "Messaging + the Switch" — module plan

Focus: a small message layer (framing, ordering) and the module's centerpiece — the
**fault-injection switch** (`switch/`) the learner builds once and every later module's
grader reuses.

## Prerequisites

Before starting M8, you should be comfortable with everything from M0–M7, plus:

- Use `io.Reader` / `io.Writer` interfaces — explain that `net.Conn` implements both.
- Call `io.ReadFull(reader, buf)` and handle `io.ErrUnexpectedEOF` (incomplete frame).
- Explain `io.EOF`: returned when the remote side closes the connection.
- Use `encoding/binary.BigEndian.PutUint32`/`Uint64` to encode integers into byte slices.
  (You do NOT need to know endianness in depth — just that network byte order is big-endian
  and the library handles the conversion.)
- Spawn a subprocess with `exec.Command("path/to/binary")` / `.Start()` / `.Wait()`.
  Explain that `.Wait()` blocks until the process exits.
- Use `bufio.Scanner` with a custom `Split` function (for length-prefixed framing).
- Write a goroutine that reads from a `net.Conn` and forwards to another `net.Conn` (relay).
- Understand the concept of "fault injection": deliberately dropping, delaying, or corrupting
  messages to test resilience.

You do NOT need to know: gRPC, protobuf, or any serialization library. The switch uses raw
byte framing (length-prefix) — you build the codec here.

---

## Shape

Go throughout. ex01/ex02 are the messaging core (harness-shape: provided `main.go` +
`go.mod`, student writes the codec/queue). ex03–ex05 are the switch: the student writes
`switch.go`; the exercise scaffold also ships a provided **counter backend** (`backend`,
a tiny line server that numbers the frames it receives on each connection:
`R<SEQ> <content>`). The switch spawns the backend itself (exec on an ephemeral
localhost port) and proxies `TARGETPORT` → backend. The `net` runner only spawns
`./switch`, so all remaining exercise grading is single-process and deterministic.

## Determinism rules (this module)

- Nothing printed may depend on wall-clock time (same as M7). The switch prints nothing.
- All fault effects are proven at the **data level** via the counter backend's reply
  numbering + echoed content:
  - `dup` → one client frame reaches the backend twice → two replies (exact-bytes expect).
  - `drop` → one client frame never reaches the backend → its reply is absent and the
    backend sequence shifts (a later frame answers with the earlier `R<s>`), byte-distinct.
  - `hold` → a frame is withheld until the next client frame arrives, then released in
    order → reply *order* swaps, and because replies echo content (`R<SEQ> <content>`),
    no two scenarios can ever produce the same transcript.
- Fault directives are read from a **fault file** whose path is injected via
  `TARGETFAULTS` env (a `StartEnv`). No control channel, no stdin pipe, no new grader
  behaviour: the engine already forwards `StartEnv` to the spawned process.

## Exercises

1. `ex01-framing` — length-prefixed message envelope (`PutFrame`/`GetFrame` on a stream),
   resisted to partial writes/reads. Harness `main.go` packs a known set then replays a
   poisoned stream (split mid-frame) and expects identical contents back. build+stdout+quiz.
2. `ex02-orders` — ordered delivery: a single writer must never reorder; the harness
   interleaves two producers over a shared outbound channel and asserts FIFO per producer
   and stable merge order; a frame dropped by a faulty producer must not deadlock the queue.
   build+stdout+quiz.
3. `ex03-switch` — transparent line relay: bind `TARGETPORT`, spawn+prox` the provided
   backend, relay client↔backend frames (newline-framed) unchanged, survive many
   connections (one backend per switch, seq resets per connection). Graded `net`:
   round-trips `ping`/`pong`-style frames through the switch; multi-connection sessions.
4. `ex04-switchfaults` — fault injection from `TARGETFAULTS` (`dup:2`, `drop:1`, `hold:1`):
   only the listed frames are affected, everything else relays cleanly. Graded `net` with
   the counter-backend replies above proving each fault both ways.
5. `ex05-gate` — **the switch in anger**: one client sends an ordered batch through a
   multi-fault scenario (drop, dup, hold across one connection and a second reconnecting
   connection); the learner must implement the switch so the *client-visible transcript*
   matches exactly. build+net+quiz.

Reuse contract (later modules): later exams run the learner's switch as their fault
injector by starting `./switch` with `TARGETFAULTS`; the message framing from ex01 is the
wire format the M9+ replicas speak under the switch.

---

## So what? (interview / portfolio)

The switch is a miniature of the thing that makes distributed testing tractable in
industry: a fault-injection proxy between a client and a real backend. Its `dup`/`drop`/
`hold` coverage is exactly the "what could go wrong in the middle?" reasoning interviewers
love, and the learner's switch doubles as a reusable test tool for every later module —
a piece of infrastructure you built and that everything after it depends on.

**Interview questions this module arms you for:**
- How do you prove a duplicate frame was actually duplicated at the *data* level?
- Why does framing matter on a stream, and how do you survive a split mid-frame?
- How would you inject faults into a system that was never designed to be tested?
- Single-writer ordering: how do you guarantee FIFO and avoid deadlock with a faulty producer?

**Portfolio artifact:** M8-ex05 `gate` — a fault-injection switch that replays exact
client-visible transcripts across drop/dup/hold, and the injector every later module reuses.