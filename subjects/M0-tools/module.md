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
| ex05 | **Gate: mini-forge**  | `mini_forge.sh` — a real mini grader |

---

## Prerequisites

None — this module IS the onboarding. You need a terminal, a text editor, and the
willingness to type commands. Everything here teaches the tools you'll use for every
subsequent module.

---

## So what? (interview / portfolio)

Everything you build from here ships through the same machinery you install in this
module: strict compiles under Werror, sanitizer runs, and a grader that is an exact
byte-compare. That tooling discipline is the cheapest thing to mention in an interview
("I'll gate every build with `-Werror` and a sanitizer") and the easiest to screen-share
as proof of rigor (a `check.py` that audits a tree is a tiny stand-in for CI).

**Interview questions this module arms you for:**
- What does a `Makefile` target like `all`/`fclean`/`re` actually do, and why enforce them?
- What do AddressSanitizer and UBSan catch that the compiler does not?
- Why does a script run under `set -euo pipefail` (and how does it differ from dash's `sh`)?
- How would you automate a build-and-verify loop for a fresh checkout?

**Portfolio artifact:** M0-ex05 `mini_forge.sh` — a tiny self-contained grader that audits
and verifies a tree the same way the real `forge` does.
