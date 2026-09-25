# M8 · Messaging & the Switch

This module builds the message layer every later module's cluster speaks: a
length-prefixed frame codec, an ordered delivery queue, and — the centrepiece —
the **fault-injection switch**, a transparent TCP relay that can duplicate,
drop, or hold frames on command. You build the relay in ex03–ex05; every
grader from M9 onward runs its distributed exercises through your switch.

## The build

- **ex01 · Message framing** — a `PutFrame`/`GetFrame` codec (4-byte
  big-endian length, then payload) hardened against split streams and
  one-byte-at-a-time readers; the wire format the switch relays and M9+
  replicas speak.
- **ex02 · Ordered delivery** — a (partition, offset) queue that turns a
  tangled arrival order into the one legal delivery order; the contract that
  "delivery order is a guarantee" is first proven here.
- **ex03 · The switch: transparent relay** — binds `TARGETPORT`, brings up a
  provided counter backend, and relays newline-framed lines unchanged, one
  fresh backend connection per client.
- **ex04 · The switch: fault injection** — a fault file (`dup:N`, `drop:N`,
  `hold:N`) damages the client→backend stream in precise, per-connection ways.
- **ex05 · Gate: the switch in anger** — drop, dup and hold in one stream
  across two connections, with an abrupt client disconnect mid-bookkeeping;
  the switch that survives this is the injector every later grader reuses.

## Rules

- Go throughout. ex01/ex02 are harness-shape: the exercise ships `main.go` +
  `go.mod`, you write the codec/queue. ex03–ex05 are whole-program: you write
  `switch.go` and, with a provided `backend/main.go` you never modify, a
  `Makefile` that builds `./switch` and `./backend`.
- The switch prints nothing to stdout; diagnostics may go to stderr.
- Nothing printed may depend on wall-clock time.
- Fault effects are proven at the **data level** through the counter backend's
  reply numbering and echoed content: a dup answers twice, a drop shifts the
  sequence, a hold swaps reply order — no two scenarios can produce the same
  transcript.
- Fault directives come from a fault file whose path is injected via
  `TARGETFAULTS` (a `StartEnv`, no control channel, no stdin pipe); the
  default is `./faults.txt`.

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