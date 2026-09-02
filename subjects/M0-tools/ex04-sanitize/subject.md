# M0-ex04 · Address Sanitizer

## Goal

`micro.c` in this directory has a subtle bug caught instantly by AddressSanitizer
(ASan). Find it with the sanitizer, fix it, and make the program build and run cleanly.

The fixed program must still print exactly:

```
micro ok
```

## The planted bug

- Read `micro.c` carefully. There is a **use-after-free**.
- Run the exercise as-is with `forge check` — the grader flags the build with a
  heap-use-after-free report. That report is your map.
- Fix by moving the write *before* the `free` (or removing the line), not by tricking
  the sanitizer (no `-fno-sanitize=address` in your final Makefile, no leaks).

## Constraints

- The grader builds with `-fsanitize=address,undefined` (it injects CFLAGS; don't undo them).
- Your Makefile targets: `all`, `fclean`, `re`.
- The program must return 0; printing the line then returning nonzero is a fail.
- No global state, no suppression files.

## Acceptance criteria

- [ ] `forge check M0-ex04` passes on your fixed version
- [ ] you can quote the exact line ASan flagged (write it in `quiz.txt`)
- [ ] `Makefile` present with `all`/`fclean`/`re`
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
What exact report did ASan emit for this program?: <the report line, e.g. "heap-use-after-free...">
Which syscall does the libc free() ultimately rely on?: <munmap or sbrk — pick one and justify>
```

## Readings

- AddressSanitizer docs: https://clang.llvm.org/docs/AddressSanitizer.html
  (read the "how it works" section — what shadows?) and the common bugs list.
- TLPI §7 (Memory Allocation — malloc/free semantics, why "free" means "unsafe to keep using").
- `man 2 munmap`, `man 3 malloc`.

## How you are graded

- `build`: `make all` under sanitizers must succeed (ASan abort ⇒ build step fails),
  then `./micro` must print exactly `micro ok` and return 0.
- `quiz`: answers from `quiz.txt`.