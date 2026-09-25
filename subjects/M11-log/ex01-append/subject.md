# M11-ex01 · Append

M11 starts with the storage contract every replicated state machine depends
on: commands enter an append-only log and receive stable indexes. You implement
that first operation now, leaving overwrite and compaction to the later
snapshot exercise.

## Shape

A replicated log is a **append-only** sequence of commands indexed from 0. Its
operations are `append`, which puts a new command at the end, and `read`, which
returns the current log. Nothing is overwritten or removed here. You write
**`log.go`** and implement:

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

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. `All` must return a copy, so editing the returned slice
cannot mutate the log. Never read the wall clock.

## Acceptance

The reference transcript is:

```text
len=0 all=
append set -> 0
append add -> 1
append del -> 2
len=3 all=set|add|del
at 1 = add
```

Indexes are assigned strictly at the end: `set` is 0, `add` is 1 and `del` is
2. Any gap or reuse fails the transcript.

## Readings

- **Reading ladder** — the Log essay frames *why* append-only; the Raft paper routes the log
  mechanics; Go slices do the data structure.
- Jay Kreps, *The Log* — https://www.confluent.io/blog/log-what-every-software-engineer-should-know-about-real-time-datas-unifying/
- Raft paper, §5.3 "Log replication": https://raft.github.io/raft.pdf
