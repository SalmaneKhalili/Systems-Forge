# M3-ex04 · The Mmap Allocator

The pool was the heap's bookkeeping; this allocator removes the heap
altogether. Its *only* memory source is the virtual-memory layer itself —
`mmap` on the way in, `munmap` on the way out. Every region you hand out is an
anonymous private page (or several). This is the pattern behind `malloc`'s
first large request: glibc serves big allocations straight from `mmap`, so
once you have written this, you have written `calloc`'s engine room.

## Shape

The exercise provides `mm.h` and a harness `main.c`; you write **`mm.c`**:

```c
void  *mm_calloc(size_t n);  /* >= n zeroed bytes, page-aligned; NULL if n==0 or on failure */
void   mm_delloc(void *p);   /* release exactly the region earlier handed out; others ignored */
size_t mm_allocated(void);   /* total bytes currently backed by the kernel */
```

`forge` compiles `main.c` (harness) + `mm.c`, runs `./test`, and diffs stdout.
The contract:

- `mmap` with `MAP_PRIVATE|MAP_ANONYMOUS`, `PROT_READ|PROT_WRITE`; failure
  returns `NULL`.
- Round every request **up to a whole page** (4096) — the kernel cannot give
  you less. `mm_calloc(0)` returns `NULL`.
- The returned region must be byte-zeroed (anonymous mappings are),
  4096-aligned, and never overlap a live neighbor.
- `mm_delloc` accepts only pointers handed out and still live; `NULL` and
  foreign pointers (like a stack address) are ignored *silently*.
- `mm_allocated()` accounts exactly for the rounded sizes of live regions.
- Bookkeeping is a static table (this is an exercise, not jemalloc); you may
  cap it at 64 live regions and fail cleanly past that.
- No `malloc`, no `free`, no global heap. `munmap` everything a region gave you.
- No `CFLAGS` redefines.

## Acceptance

Graded `build`: strict compile + harness stdout diff, exit 0, empty stderr.
The harness checks alignment, zero-fill, disjointness, exact rounded
accounting, silent invalid-free, and the churn/clean-exit cycle:

- `mm_calloc(1)`, `mm_calloc(1500)`, `mm_calloc(8200)` all non-NULL,
  4096-aligned, mutually disjoint, and byte-zeroed.
- `mm_allocated()` reports exactly `4096 + 4096 + 12288` after those three.
- `mm_delloc` of the middle region drops `mm_allocated()` by 4096.
- `mm_delloc(NULL)` and `mm_delloc(&stackvar)` change nothing.
- 40-region churn stays accounting-clean; releasing everything returns to 0.
- No leak: every region is freed by the end of the harness.

A free that un-maps a foreign region (or skips a live one) fails the
accounting check; a region that overlaps another fails disjointness (likely
with an ASan report too). `quiz.txt` is complete (see Quiz).

## Readings

- **Reading ladder** — start with `man 2 mmap` (MAP_ANONYMOUS), then TLPI
  §49.1 "Overview" and §49.7 below.
- The Linux Programming Interface, §49.1 "Overview" and §49.7 "Anonymous Mappings".
- `man 2 mmap` — the `MAP_ANONYMOUS` paragraph, and why the kernel allocates in pages
  (the `length` fiddling notes near the end).
- TLPI §49.2 "Creating a Mapping" including the `map()` helper for page-rounding, and
  the classic note: "the kernel rounds length up".

## Quiz

1. Which syscall hands back an anonymous region to the kernel?
2. What granularity does mmap allocate in?