# M10-ex02 · Vote

## Goal

A participant's vote is a **promise**. When it votes `commit`, it guarantees
it can and will commit once told — it may never later switch to abort,
because the coordinator may already have instructed other participants to
commit based on that vote. When a participant cannot make that guarantee (a
resource is unavailable, a local write is doomed), it must vote `abort`.

Implement `vote.go`:

```go
// VoteFor returns the vote for a mobile transaction given its environment.
// The participant votes commit iff it can guarantee the write, else abort.
func VoteFor(canCommit bool) string
```

The provided `main.go` feeds a fixed sequence of environments and prints one
vote each. `make all` must build `test`; `./test` must print the reference
transcript exactly.

## Constraints

- Go, standard library only; file is `vote.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock.

## Acceptance

Reference transcript:

```
t0: commit
t1: abort
t2: commit
t3: abort
```

The point of the exercise is the **finality** of a commit vote: the transcript
lists only votes here, but the quiz captures why an `abort` after a `commit`
vote is forbidden.