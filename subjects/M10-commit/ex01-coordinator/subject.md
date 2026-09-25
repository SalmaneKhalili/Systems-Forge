# M10-ex01 · The coordinator

M10 begins with the process that turns independent participant votes into one
atomic decision. You build the deterministic 2PC state machine that asks every
participant to prepare, chooses commit or abort once, and makes every
participant follow that same final outcome.

## Shape

Two-phase commit (2PC) begins with **prepare**: the coordinator asks every
participant to vote. The coordinator then broadcasts the final commit/abort
decision in phase two. You write **`coordinator.go`** as a small deterministic
state machine:

```go
type Voter func(id int) bool // returns the participant's vote (true = commit)

// Coordinator drives 2PC over n participants. It returns the final
// decision: "commit" only if EVERY participant voted commit, else "abort".
type Coordinator struct{ n int }

func (c *Coordinator) Run(vote Voter) string
```

The exercise ships `main.go`, which runs the coordinator over a fixed
participant set and prints the transcript. Round one asks each participant
(`prepare id → vote`); round two emits the decision.

Use Go and the standard library only. `make all` must build `test`, and
`./test` must print the reference transcript exactly; `expected.txt` is
authoritative and whitespace is normalized. Never read the wall clock or sleep.
A participant that votes no makes the whole outcome `abort`, even if every other
participant voted yes.

## Acceptance

The reference transcript is:

```text
prepare 0 -> yes
prepare 1 -> yes
prepare 2 -> no
DECISION abort
commit 0
abort 1
abort 2
```

The coordinator always emits a decision line, then one per-participant action:
`commit` for every participant after a commit, or `abort` for every participant
after an abort. Here participants 0 and 1 voted yes, but participant 2 voted
no, so **everyone** is aborted. Partial commit is not allowed.

## Readings

- **Reading ladder** — one chapter covers both atomic commit and 2PC; Kleppmann's notes make
  the coordinator's state machine concrete.
- *Designing Data-Intensive Applications*, Chapter 9 "Consistency and Consensus" — "Atomic
  commit and two-phase commit (2PC)".
- Martin Kleppmann, *Distributed Systems* notes — atomic commit / 2PC:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
