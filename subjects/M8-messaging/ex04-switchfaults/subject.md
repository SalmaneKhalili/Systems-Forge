# M8-ex04 · The switch: fault injection

ex03 gave you a transparent relay; this exercise adds the **fault plane** that
makes it a test instrument. A small text file leaves the switch transparent by
default, but lets it damage the client→backend stream in three precise,
per-connection ways:

- `dup:N` — the N-th client frame is forwarded to the backend **twice**.
- `drop:N` — the N-th client frame is **not forwarded at all**; every other
  frame relays untouched.
- `hold:N` — the N-th client frame is withheld and forwarded only **after the
  next** client frame has been forwarded (a reordering fault).

Frames are newline-terminated lines. Frame numbers count the client frames of
the **current connection**, starting at 1. The directive parser reads one
`op:index` per line from the **fault file**; a missing index is ignored and a
missing file means "no faults". `TARGETFAULTS`, when set, overrides the path
(the default is `./faults.txt`).

## Shape

Whole-program, evolving from M8-ex03. The scaffold is the previous switch:
`backend/main.go` (provided, unmodified), `go.mod`, `Makefile`, and a relay
`switch.go` you extend. Fault directives must be parsed **before** the accept
loop, and applied per connection (a fresh connection restarts the numbering
and re-applies every directive).

Why the backend's numbered replies prove faults: the backend still answers
`R<k> <content>` per received frame. A dropped frame is visible because the
next relayed frame answers with the *earlier* number. A duplicated frame is
visible because one send gets two replies. A held frame is visible because the
reply ordering swaps (`R4 M4` before `R5 M3`, never the other way).

Releases must be deterministic: the held frame goes out **after** the
following forwarded frame, in that exact order.

Nothing may be printed to stdout; diagnostics may go to stderr.

## Acceptance

`make all`, then the grader runs `./switch` with the shipped `faults.txt`
(`dup:1`, `drop:2`, `hold:3`) and plays, over two connections:

- connection 1: `dupme` → expect `R1 dupme` and `R2 dupme`; then
  `dropme`, `held`, `w` → expect `R3 w` and `R4 held`; then `z` → `R5 z`.
  (The dropped `dropme` never gets a reply; the held `held` comes back only
  after `w`.)
- connection 2 (fresh, same faults): `a` → `R1 a`, `R2 a`; then `b`, `c` →
  `R3 c`.

Graded `net`.

## Readings

- The concept of a **fault injection switch** — the M8 module notes and the
  later capstones: every later grader drives its chaos through a switch like
  this one; the `TARGETFAULTS` env is the channel it uses.
- Go `os/exec`, `strconv`, `strings` — spawning, parsing.

## Quiz

1. Where does the switch read its fault directives from?
2. Which directive forwards one client frame to the backend twice?
3. In what situation is the frame a drop directive targets placed?