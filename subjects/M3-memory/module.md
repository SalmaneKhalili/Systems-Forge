# M3 · Memory

Processes run in address space; memory is the material the address space is made of. This
module leaves `malloc` behind on purpose: you will build three allocators with your own hands
(a bump arena, a fixed-slot pool, an `mmap` page allocator), run a copy-on-write experiment,
and finish by wiring graceful out-of-memory handling into a growable buffer. After M3, an
allocation failure is never a crash — it is a branch you already wrote.

**Rules of the module**

- Two shapes: some exercises give you a `main.c` harness + a header — you implement the
  allocator in your own `.c` (never edit the harness or header). The virtual-memory exercise
  is a full program of yours.
- You write every `Makefile` (same contract as M1/M2: `all`→`./test`, `fclean`, `re`; never
  redefine `CFLAGS`/`LDFLAGS`).
- `forge` compiles with `-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined`. The
  sanitizers make a mess of a hand-rolled allocator the same way they make a mess of a
  careless one — garbage pointers, overflows of globals, and leaks all surface.
- Determinism rule from M2 still rules: no raw addresses in output, no sleeps, every failure
  path returns gracefully instead of a segfault.
- "OOM" here is *bounded memory exhaustion* (an arena/pool that is genuinely full) — real
  and deterministic, nothing that depends on the machine's RAM.

| ex | topic | artifact you deliver |
|----|-------|----------------------|
| ex01 | bump arena | `arena.c`: 512-byte arena, 8-aligned allocs, reset |
| ex02 | copy-on-write | whole `main.c`: two `mmap`s, fork, MAP_PRIVATE vs MAP_SHARED |
| ex03 | fixed slot pool | `pool.c`: 8×32-byte slots, borrow/return, NULL on empty |
| ex04 | mmap allocator | `ma_alloc.c`: page-rounded anonymous `mmap`/`munmap` |
| ex05 | **Gate: OOM-safe buffer** | `buf.c`: growable buffer on a 64-byte arena that refuses growth gracefully |

---

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