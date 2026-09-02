# M9-ex05 · The sequencer

## Goal

A **sequencer** is a single process that turns concurrent clients' unordered
streams into one strictly ordered stream: every frame it receives is assigned a
sequence number, and the numbers only ever increase. This is the "single-writer
total order" that replicated state machines build on.

Write `srv.go`: a TCP server that listens on `TARGETPORT` and applies the
Lamport **receive rule** to every frame. Frames are lines

```
<content> <stamp>
```

where `<stamp>` is the client's claimed Lamport value. For each frame the
server assigns

```
seq = max(seq, stamp) + 1
```

and replies `R<seq> <content> <stamp>`. The sequence counter lives for the
life of the process — never reset per connection.

The grader dials at least two connections back-to-back; replies must match the
reference transcript byte-for-byte.

## Constraints

- Go, standard library only; file is `srv.go`.
- `make all` must build `srv` (the grader starts `./srv`).
- Multiple concurrent connections must be served (an accept loop).
- Never read the wall clock; no timestamps, no sleeps in replies.
- Port number comes from `TARGETPORT`.

## Acceptance

Reference transcript (`a` = conn 1, `b` = conn 2):

```
R1 a 0
R10 b 9
R11 c 9
R12 d 11
R13 e 0
R21 f 20
```

Conn 2's first reply is `R12` — the counter carries across connections. And
`c 9` after `b 9` shows the merge rule: the client's *claimed* stamp does not
reset the sequencer; wall-clock intuition breaks here, Lamport's max-then-add
is the law.