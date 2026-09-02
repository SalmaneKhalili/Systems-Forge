# M5 · Files & I/O

Every byte in and out of your program travels through a small table of small integers:
**file descriptors**. `open` hands you the next free one, `read`/`write` move bytes through
them, and `dup2` lets you swap what a handle means — that is all a redirection really is.
Once you own these, `printf` and `fgets` stop being magic: they are thin wrappers over the
same table.

This module is the "mechanical keyboard" of systems programming — you build the tools
(nothing but `open`, `read`, `write`, and `dup2`) that the rest of the curriculum leans on.

| exercise | kind       | what you build |
|----------|-----------|----------------|
| ex01     | descriptors | the fd table: numbers, reuse, lowest-available |
| ex02     | readwrite   | `read()`/`write()` vs stdio, byte for byte |
| ex03     | redirect    | repoint stdout with `dup2` across a fork |
| ex04     | tee         | a real `tee`: stdin → stdout + file |
| ex05     | **gate**    | a log parser over an fd (harness shape) |

All five are graded `build` + `quiz` under `-std=gnu11 -Wall -Wextra -Werror` with
ASan/UBSan. ex04's copied file is additionally asserted by an `artifact` check.

Determinism is the whole point: verdicts are exact byte counts and exact strings, printed
only after every fd is closed or every child is reaped — no PIDs, no addresses, no races.