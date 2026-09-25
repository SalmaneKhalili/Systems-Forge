# M15-ex03 · Conflict

Staged transactions prevent partial commits, but concurrent transactions can
still target the same key and turn one write into a silent lost update. This
exercise makes overlap explicit through write-set comparison and resolves the
merged state with the local transaction's value protected. You deliver
`conflict.go` with detection, deterministic resolution, and stable key order.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`conflict.go`** using only the Go standard library and implement:

```go
func Conflict(writesA, writesB map[string]string) bool // same key in both -> true
func Resolve(mine, theirs map[string]string) map[string]string // merge; your key wins
func SortedKeys(m map[string]string) []string          // map keys sorted (for tests)
```

`Conflict` is true when both write sets contain the same key. `Resolve` returns
a new map, merges keys found only in `theirs`, and keeps `mine`'s value for a
key in both. `SortedKeys` provides sorted map keys for the tests. Never read the
wall clock. The reference transcript is `expected.txt` with whitespace
normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
conflicts -> true
x -> 1
y -> c
```

`Conflict({x:1},{x:2})` is true because both transactions write `x`.
`Resolve({x:1},{x:b,y:c})` keeps the local `x=1` instead of accepting the
conflicting `x=b`, and merges `y=c`. Replacing the local value with the other
transaction's value is the lost update this exercise rules out.

## Readings

- **Reading ladder** — a conflict is a lost update in waiting; comparing write sets is the
  detector.
- *Designing Data-Intensive Applications*, Chapter 7 — write-write conflicts and lost updates.
- Wikipedia, "Isolation (database systems)": https://en.wikipedia.org/wiki/Isolation_(database_systems)
