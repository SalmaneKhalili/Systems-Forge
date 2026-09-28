# M3-ex01 · bump arena

## Goal

Implement `arena.c` backing the provided `arena.h` — a **bump allocator**: a fixed 512-byte
region of memory where every allocation is served by advancing a single offset pointer. It is
the simplest allocator that exists, and the foundation of real systems like jemalloc arenas.

```c
void  *arena_alloc(size_t n);   /* 8-aligned, NULL if n==0 or not enough room */
size_t arena_left(void);        /* bytes still available */
void   arena_reset(void);       /* all memory available again */
```

`forge` compiles `main.c` (provided harness) + `arena.c`, runs `./test`, diffs stdout.

## Constraints

- The arena lives in your `arena.c`: a `static` 512-byte buffer, nothing else global.
- Every returned pointer must be **8-byte aligned** — even after calls like `arena_alloc(3)`.
  You may pad the offset up to an 8-boundary; prefer `_Alignas(8)` on the buffer so the base
  never wastes a byte.
- Never hand out memory past the end of the 512 bytes; when `n` doesn't fit, return `NULL`.
  `arena_alloc(0)` returns `NULL` too.
- No `malloc`, no `mmap` — the arena is the only source of memory.
- `arena_reset` makes the whole 512 bytes available again (no clearing needed).
- No `CFLAGS` redefines.

## Acceptance criteria

- [ ] `arena_alloc(3)+(5)+(7)+(1)` all succeed, all aligned
- [ ] 60 more `arena_alloc(8)` calls succeed; the *next* call returns `NULL` (512 bytes exact)
- [ ] `arena_left()` reports exactly 0 at exhaustion
- [ ] no block overwrites its neighbor (the harness tags every block and re-reads them)
- [ ] after `arena_reset`: `arena_left() == 512`, `arena_alloc(512)` succeeds,
      `arena_alloc(1)` fails cleanly
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
An allocator that assigns memory only by advancing a pointer is called what?: <answer>
Which C11 keyword forces a static buffer to a required alignment?: <answer>
```

## Readings

- **Reading ladder** — start with `man 3 malloc`'s contract, then the jemalloc
  arena overview below, then TLPI §7.1 for the heap model.
- The Linux Programming Interface, §7.1 "Allocating Memory on the Heap" is the background;
  for the *arena* pattern itself, read the jemalloc overview below.
- jemalloc's arena overview (skim): https://jemalloc.net/jemalloc.3.html search "arena".
- `_Alignas` / `_Alignof`: cppreference https://en.cppreference.com/w/c/language/_Alignas
- Basic pointer-arithmetic refresher if needed: K&R §5.4 "Address Arithmetic".

## How you are graded

- `build`: strict compile + harness stdout diff (whitespace-insensitive), exit 0, empty
  stderr. The harness checks alignment, exact 512-byte exhaustion, tag-based non-overlap,
  reset semantics, and the `NULL`-on-when-full contract. A bumper that overwrites its
  boundary fails the tag check (and likely trips ASan's global redzone).