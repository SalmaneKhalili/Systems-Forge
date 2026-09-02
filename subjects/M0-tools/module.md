# M0 · Tools of the Trade

The piscine begins with the instrument kit you will use every single day from here on:
`make`, a strict compile pipeline, `bash`, and a preferred editor. Nothing here is busywork —
the grader for _this_ module is the same shape as the grader you will build later.

**Rules of the module**

- Every C project must build with `-std=gnu11 -Wall -Wextra -Werror` and, when the subject
  asks for it, `-fsanitize=address,undefined`. `forge` enforces this itself.
- You hand-write `Makefiles` for every exercise — exact targets `all`, `fclean`, `re`.
- Scripts you own must start with a shebang and run under `set -euo pipefail`.
- Read the assigned sections before you touch the keyboard; the quiz checks it.

| ex   | topic                 | artifact you deliver                                 |
| ---- | --------------------- | ---------------------------------------------------- |
| ex01 | Makefile discipline   | a `Makefile` whose targets pass the strict build     |
| ex02 | Shell hygiene         | a `solve.sh` that survives `-u`, `-e`, `-o pipefail` |
| ex03 | A five-minute checker | a stdlib-only `check.py` that audits a tree          |
| ex04 | Address Sanitizer     | a `micro.c` freed of its use-after-free              |
| ex05 | **Gate: mini-forge**  | `mini_forge.sh` — a real mini grader                 |
