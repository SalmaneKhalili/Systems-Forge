# M9-ex05 · The sequencer

M9 ends by putting the logical clock on the network. A **sequencer** turns
concurrent clients' unordered streams into one strictly ordered stream by
assigning every frame an increasing sequence number; that single-writer order
is the foundation replicated state machines build on.

## Shape

You write **`srv.go`**: a TCP server that listens on `TARGETPORT` and applies
the Lamport **receive rule** to every frame. Frames are lines:

```text
<content> <stamp>
```

`<stamp>` is the client's claimed Lamport value. For each frame the server
assigns:

```text
seq = max(seq, stamp) + 1
```

It then replies `R<seq> <content> <stamp>`. The sequence counter lives for the
life of the process and is never reset per connection. The grader dials at
least two connections back-to-back, so an accept loop must serve multiple
concurrent connections and replies must match the reference transcript
byte-for-byte.

Use Go and the standard library only. `make all` must build `srv`; the grader
starts `./srv`. The port comes from `TARGETPORT`. Never read the wall clock:
replies contain no timestamps and the server performs no sleeps.

## Acceptance

The reference transcript (`a` = conn 1, `b` = conn 2) is:

```text
R1 a 0
R10 b 9
R11 c 9
R12 d 11
R13 e 0
R21 f 20
```

- Conn 2's first reply is `R12`, proving the counter carries across connections.
- `c 9` after `b 9` proves the merge rule: the client's *claimed* stamp does
  not reset the sequencer; wall-clock intuition fails here, while Lamport's
  max-then-add rule governs the reply.

## Readings

- **Reading ladder** — the sequencer is the Lamport receive rule wrapped in a TCP server; the
  wire framing is the part you already trained in M8.
- Lamport's receive rule, "Time, Clocks…" §2.2:
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
- Kleppmann's total-order-broadcast treatment:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Go `net` docs — the TCP accept loop: https://pkg.go.dev/net
