# M0-ex04 · Address Sanitizer

The checker establishes a precise report; AddressSanitizer now makes a memory bug
visible in that same build-and-verify loop. This exercise gives you a planted
use-after-free and asks you to use the report as evidence, not to suppress it. You
deliver a fixed `micro.c` that builds and runs cleanly.

## Shape

`micro.c` in this directory has a subtle bug caught instantly by AddressSanitizer (ASan).
Find it with the sanitizer, fix it, and make the program build and run cleanly. The fixed
program must still print exactly:

```text
micro ok
```

Read `micro.c` carefully: there is a **use-after-free**. Run the exercise as-is with
`forge check`; the grader flags the build with a `heap-use-after-free` report, and that
report is your map. Fix the bug by moving the write _before_ the `free` (or removing the
line), not by tricking the sanitizer. The final Makefile must not contain
`-fno-sanitize=address`, and there must be no leaks.

The grader builds with `-fsanitize=address,undefined` and injects CFLAGS; do not undo
them. Your Makefile needs the targets `all`, `fclean`, and `re`. The program must return
0; printing `micro ok` and then returning nonzero is a failure. No global state and no
suppression files are allowed.

## Acceptance

`forge check M0-ex04` must pass on the fixed version. `make all` under the injected
sanitizers must succeed; an ASan abort makes the build step fail. The resulting `./micro`
must then print exactly the line above and return 0.

- The fix removes the use-after-free instead of hiding the report.
- `Makefile` is present with `all`/`fclean`/`re`, and the sanitizer flags remain active.
- You can quote the exact line ASan flagged in `quiz.txt`.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` requires the sanitizer-enabled `make all`, exact
`micro ok` output, and exit 0; `quiz` reads the answers from `quiz.txt`.

## Readings

- AddressSanitizer docs: https://clang.llvm.org/docs/AddressSanitizer.html
  (read "Usage" and the "Additional Checks" list — the use-after-free example) and
  the common bugs list.
- How ASan works — what shadows? https://github.com/google/sanitizers/wiki/AddressSanitizerAlgorithm
  (this is the "how it works" page: red zones, quarantine, and shadow memory).
- TLPI §7 (Memory Allocation — malloc/free semantics, why "free" means "unsafe to keep using").
- `man 2 munmap`, `man 3 malloc`.

## Quiz

1. What exact report did ASan emit for this program?: <the report line, e.g. "heap-use-after-free...">
2. Which syscall does the libc free() ultimately rely on?: <munmap or sbrk — pick one and justify>
