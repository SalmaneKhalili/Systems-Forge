# M10-ex05 · The transaction

## Goal

Bring two-phase commit to the wire. Write `gate.go`: a TCP server on
`TARGETPORT` whose participants' votes come from a provided `votes.txt`
(the scaffold supplies it — a line `3` means "transaction 3 has prepared
votes from every participant", a line `1` means "transaction 1 does not").

For each client line `tx <id>`:

- if the id is in `votes.txt` (all participants prepared), the coordinator
  reply `COMMIT\n`;
- otherwise (someone did not prepare / could not commit), reply `ABORT\n`.

`make all` must build `gate`; the grader starts `./gate` and dials multiple
connections; replies must match byte-for-byte. The coordinator's decision is
a strict all-or-nothing: a missing entry is an abort.

## Constraints

- Go, standard library only; file is `gate.go`.
- `make all` must build `gate` (the grader starts `./gate`).
- Multiple concurrent connections must be served (an accept loop).
- Never read the wall clock; no timestamps, no sleeps in replies.
- Port number comes from `TARGETPORT`.

## Acceptance

`votes.txt` (scaffold) lists transactions whose participants all prepared:

```
3
2
7
```

So `tx 3`, `tx 2`, `tx 7` must `COMMIT`, and every other id (0, 1, …) must
`ABORT`. The transcript is byte-deterministic. A gateway that replies
`COMMIT` for an id absent from the file — or that ties the decision to
anything other than the whole-file-committed check — fails.