# M11-ex01 · Append

## Goal

A replicated log is a **append-only** sequence of commands, indexed from 0.
The two operations a replica offers are `append` (put a new command at the
end) and `read` (return the current log). Nothing is ever overwritten or
removed here — snapshotting is a later exercise.

Implement `log.go`:

```go
type Log struct {
	entries []string
}

// Append adds cmd at the end and returns its index.
func (l *Log) Append(cmd string) int

// Len returns how many entries are in the log.
func (l *Log) Len() int

// All returns a copy of the entries.
func (l *Log) All() []string
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `log.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `All` must not let the caller mutate the log by editing the returned slice.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
len=0 all=
append set -> 0
append add -> 1
append del -> 2
len=3 all=set|add|del
at 1 = add
```

Indexes are assigned strictly at the end: `set` is 0, `add` is 1, `del` is 2.

## Readings

- **Reading ladder** — the Log essay frames *why* append-only; the Raft paper routes the log
  mechanics; Go slices do the data structure.
- Jay Kreps, *The Log* — https://www.confluent.io/blog/log-what-every-software-engineer-should-know-about-real-time-datas-unifying/
- Raft paper, §5.3 "Log replication": https://raft.github.io/raft.pdf
