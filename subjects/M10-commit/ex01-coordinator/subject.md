# M10-ex01 · The coordinator

## Goal

Two-phase commit (2PC) is how a **coordinator** makes a set of participants
agree on one transaction: either it **commits** everywhere or it **aborts**
everywhere. Phase one is *prepare*: the coordinator asks every participant
to vote. Phase two is *commit/abort*: the coordinator tells every participant
the final decision.

Implement `coordinator.go` as a small deterministic state machine:

```go
type Voter func(id int) bool // returns the participant's vote (true = commit)

// Coordinator drives 2PC over n participants. It returns the final
// decision: "commit" only if EVERY participant voted commit, else "abort".
type Coordinator struct{ n int }

func (c *Coordinator) Run(vote Voter) string
```

The provided `main.go` runs the coordinator over a fixed participant set and
prints the transcript: round one asks each participant (`prepare id → vote`),
then round two emits the decision.

## Constraints

- Go, standard library only; file is `coordinator.go`.
- `make all` must build `test`; `./test` must print the reference transcript
  exactly (`expected.txt`, whitespace normalized).
- Never read the wall clock. No sleeps.
- A participant that votes no makes the whole outcome `abort` — even if
  others voted yes.

## Acceptance

Reference transcript:

```
prepare 0 -> yes
prepare 1 -> yes
prepare 2 -> no
DECISION abort
commit 0
abort 1
abort 2
```

The coordinator always emits a decision line, then one per-participant action
(`commit` for every participant on a commit; `abort` for every participant on
an abort). The tell: even though participants 0 and 1 said yes, participant 2
said no, so **everyone** is aborted — partial commit is not allowed.

## Readings

- **Reading ladder** — one chapter covers both atomic commit and 2PC; Kleppmann's notes make
  the coordinator's state machine concrete.
- *Designing Data-Intensive Applications*, Chapter 9 "Consistency and Consensus" — "Atomic
  commit and two-phase commit (2PC)".
- Martin Kleppmann, *Distributed Systems* notes — atomic commit / 2PC:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
