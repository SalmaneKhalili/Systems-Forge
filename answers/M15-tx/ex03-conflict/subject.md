# M15-ex03 · Conflict

## Goal

When two transactions edit overlapping data a **conflict** (potential lost
update) arises. Optimistic schemes detect the conflict and then resolve by
merging — but a write you made must never be silently overwritten by the other
transaction's write to the same key (that would be a lost update).

Implement `conflict.go`:

```go
func Conflict(writesA, writesB map[string]string) bool // same key in both -> true
func Resolve(mine, theirs map[string]string) map[string]string // merge; your key wins
func SortedKeys(m map[string]string) []string          // map keys sorted (for tests)
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `conflict.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `Resolve` returns a new map; keys only in `theirs` are merged in; a key in
  BOTH keeps `mine`'s value.
- No wall clock.

## Acceptance

Reference transcript:

```
conflicts -> true
x -> 1
y -> c
```

`Conflict({x:1},{x:2})` is true (both write `x`). `Resolve({x:1},{x:b,y:c})`
must keep `x=1` (yours wins over the conflicting `x=b`) and merge `y=c`.
Overwriting your own `x` with the other's `b` is the bug.
