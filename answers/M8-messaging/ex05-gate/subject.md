# M8-ex05 — Gate: the switch in anger

## Goal

The gate of M8: your switch, complete with broadcast-plane and teardown
discipline, must survive a **full fault scenario** across two connections —
including an abrupt client disconnect in the middle of bookkeeping — and hand
the next connection a clean, per-connection-faulted stream. This switch is the
injector every later module's grader will throw in front of a cluster.

## Shape

Whole-program, evolving from M8-ex04. The switch must now prove:

- **Mixed faults in one stream**, in any order (dup, drop, hold can all fire
  in the same connection).
- **Mid-scenario teardown.** A client may drop its connection while the switch
  is holding a frame or mid-relay; the switch must close the relay's backend
  side (unreleased held frames are simply gone — they were never forwarded)
  and stay ready for the next connection.
- **Re-applied faults, reset numbering.** Connect 2 restarts frame numbers and
  the backend repeats its `R1…` sequence.

`TARGETFAULTS` (default `./faults.txt`) carries the directives; per-connection
numbering and the `R<k> <content>` reply contract are exactly as in M8-ex04.

## Acceptance criteria (all graded)

With the shipped `faults.txt` (`dup:1`, `drop:2`, `hold:4`), `make all` and the
grader's two connections must produce **exactly**:

connection 1:
- `M0` → `R1 M0`, `R2 M0` (duplicated);
- `M1`, `M2` → `R3 M2` (M1 dropped, M2 clean);
- `M3`, `M4` → `R4 M4`, `R5 M3` (M3 held until M4 went through).

connection 2 (fresh connection, faults re-applied, backend reset):
- `target` → `R1 target`, `R2 target`;
- `power`, `up` → `R3 power`, `R4 up`.

A switch that releases a held frame before the following frame (or not at
all), or that leaks fault state across connections, prints a different reply
order and fails.

## Readings

- The M8-ex04 fault contract — per-connection numbering, release-after-next.
- Go `net`/`bufio` — half-close and teardown semantics.

## Quiz

1. After a client disconnects mid-scenario, what must a correct switch do?
2. Does the dup directive also apply to a freshly connected client?
3. Which ordering of two replies proves a hold was released after the following frame?