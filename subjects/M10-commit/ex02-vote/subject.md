# M10-ex02 · Vote

The coordinator from ex01 can only drive a final outcome because each
participant's vote is a **promise**. This exercise makes that promise explicit:
a participant votes `commit` only when it can and will commit, and votes
`abort` when it cannot make that guarantee.

## Shape

When a participant votes `commit`, it may never later switch to abort because
the coordinator may already have told other participants to commit based on
that vote. If a resource is unavailable or a local write is doomed, the
participant votes `abort`. You write **`vote.go`** and implement:

```go
// VoteFor returns the vote for a mobile transaction given its environment.
// The participant votes commit iff it can guarantee the write, else abort.
func VoteFor(canCommit bool) string
```

The exercise ships `main.go`, which feeds a fixed sequence of environments and
prints one vote for each. `make all` must build `test`, and `./test` must print
the reference transcript exactly. Use Go and the standard library only. The
reference transcript is `expected.txt`, with whitespace normalized; never read
the wall clock.

## Acceptance

The reference transcript is:

```text
t0: commit
t1: abort
t2: commit
t3: abort
```

The exercise proves the **finality** of a commit vote. This transcript lists
only votes, while the quiz captures why an `abort` after a `commit` vote is
forbidden.

## Readings

- **Reading ladder** — a vote is a *promise*; read why the coordinator must treat it as
  binding, in both the commit and abort directions.
- *Designing Data-Intensive Applications*, Chapter 9 — the participant promise inside 2PC.
- Kleppmann notes, the 2PC participant protocol:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
