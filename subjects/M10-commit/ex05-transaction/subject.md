# M10-ex05 · The transaction

M10 ends by composing the coordinator, participant promise and quorum work on
the wire. The result is a real 2PC gateway whose commit decision is strict and
all-or-nothing across every client connection.

## Shape

You write **`gate.go`**, a TCP server on `TARGETPORT`. A provided `votes.txt`
supplies the participants' votes: a line `3` means transaction 3 has prepared
votes from every participant, while a line `1` means transaction 1 does not.

For each client line `tx <id>`:

- if the id is in `votes.txt` (all participants prepared), the coordinator
  replies `COMMIT\n`;
- otherwise (someone did not prepare / could not commit), it replies `ABORT\n`.

Use Go and the standard library only. `make all` must build `gate`; the grader
starts `./gate` and dials multiple connections, which an accept loop must
serve concurrently. Replies must match byte-for-byte. The coordinator's
decision is strict all-or-nothing: a missing entry is an abort. The port comes
from `TARGETPORT`; never read the wall clock, and put no timestamps or sleeps
in replies.

## Acceptance

The scaffold's `votes.txt` lists transactions whose participants all prepared:

```text
3
2
7
```

Therefore `tx 3`, `tx 2` and `tx 7` must receive `COMMIT`; every other id,
including 0 and 1, must receive `ABORT`. The transcript is byte-deterministic.
A gateway that replies `COMMIT` for an absent id, or ties the decision to
anything other than the whole-file-committed check, fails.

## Readings

- **Reading ladder** — the wire gate composes ex01–ex04; the protocol is 2PC, the framing is M8.
- *Designing Data-Intensive Applications*, Chapter 9 — 2PC end to end.
- Reread M10-ex01 (coordinator) and M10-ex02 (the promise).
- Go `net` / `bufio` docs: https://pkg.go.dev/net, https://pkg.go.dev/bufio
