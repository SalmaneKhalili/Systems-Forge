# M3 · Memory

Processes run in address space; memory is the material the address space is
made of. This module leaves `malloc` behind on purpose: you build three
allocators with your own hands — a bump arena, a fixed-slot pool, and an
`mmap` page allocator — run a copy-on-write experiment, and finish by wiring
graceful out-of-memory handling into a growable buffer. After M3, an
allocation failure is never a crash; it is a branch you already wrote.

## The build

- **ex01 · Bump arena** — `arena.c`: a 512-byte arena with 8-aligned
  allocations and reset — the simplest allocator that exists, and the base
  the gate grows out of.
- **ex02 · Copy-on-write** — a whole `main.c`: two `mmap`s, a `fork`, and the
  `MAP_PRIVATE` vs `MAP_SHARED` difference proven byte-for-byte.
- **ex03 · Fixed slot pool** — `pool.c`: 8×32-byte slots with borrow/return
  accounting and `NULL` on empty; the pool pattern behind request structs and
  kernel descriptor tables.
- **ex04 · The mmap allocator** — `mm.c`: page-rounded anonymous
  `mmap`/`munmap`, the engine room of `calloc`'s big-request path.
- **ex05 · Gate: OOM-safe buffer** — `buffer.c`: a growable buffer on a
  64→1024-byte arena that **refuses** growth gracefully, leaving its data
  untouched.

## Rules

- Two shapes: ex01/ex03/ex04/ex05 give you a `main.c` harness plus a header —
  you implement the allocator in your own `.c` (never edit the harness or
  header). ex02 is a full program of yours (`main.c` + `Makefile`).
- You write every `Makefile` (same contract as M1/M2: `all`→`./test`,
  `fclean`, `re`; never redefine `CFLAGS`/`LDFLAGS`).
- `forge` compiles with `-std=gnu11 -Wall -Wextra -Werror
  -fsanitize=address,undefined`. The sanitizers make a mess of a hand-rolled
  allocator the same way they make a mess of a careless one — garbage
  pointers, overflows of globals, and leaks all surface.
- The M2 determinism rule still rules: no raw addresses in output, no sleeps,
  and every failure path returns gracefully instead of a segfault.
- "OOM" here is *bounded memory exhaustion* — an arena or pool that is
  genuinely full: real and deterministic, nothing that depends on the
  machine's RAM.

## Prerequisites

Before starting M3, you should be comfortable with everything from M0–M2, plus:

- Explain what a pointer is: `int *p = &x;` means `p` holds the address of `x`.
- Dereference a pointer: `*p = 42;` writes 42 to the address `p` holds.
- Do pointer arithmetic: `p + 1` advances by `sizeof(int)` bytes, not 1 byte.
- Explain `sizeof(int)` vs `sizeof(char)` and why `(char *)p + 1` advances by 1 byte.
- Cast between pointer types: `(char *)ptr` to treat raw bytes as a byte stream.
- Call `malloc(N)` and check if it returned `NULL`. Free what you allocated.
- Explain what `mmap` conceptually does: gives you a block of memory (or a file mapping).
- Know what a page is: the OS manages memory in fixed-size chunks (usually 4096 bytes).
- Explain what `fork` does at a high level: duplicates the process (you learned this in M2).

You do NOT need to know: threads, mutexes, or concurrency. Those come in M4.

## So what? (interview / portfolio)

Hand-building allocators is the clearest evidence you understand what `malloc` actually is
and how memory really works — a top signal for systems/OS and embedded interviews. After
this module "the heap is a data structure" is something you can *show*, and the OOM gate
teaches the safety habit that most production code gets wrong: failing gracefully instead
of one giant segfault.

**Interview questions this module arms you for:**
- How does a bump arena allocate, and where is it faster than `malloc`?
- Copy-on-write vs. plain copy: what does `fork` + a shared `mmap` actually do?
- When would a fixed-slot pool beat the general allocator, and at what cost?
- What do `mmap`/`munmap` give you that `malloc`/`free` do not?

**Portfolio artifact:** M3-ex05 `buf.c` — a growable buffer that handles out-of-memory as
a handled branch, not a crash.