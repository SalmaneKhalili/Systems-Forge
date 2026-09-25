# M3-ex01 · Bump Arena

M3 leaves `malloc` behind. Your first allocator is the simplest one that
exists: a **bump arena** — a fixed 512-byte region where every allocation is
served by advancing a single offset pointer. It is the foundation of real
allocators (jemalloc's arenas) and the base the gate's buffer (ex05) grows out
of. You implement `arena.c` behind the provided header.

## Shape

The exercise provides `arena.h` and a harness `main.c`; you write **`arena.c`**:

```c
void  *arena_alloc(size_t n);   /* 8-aligned, NULL if n==0 or not enough room */
size_t arena_left(void);        /* bytes still available */
void   arena_reset(void);       /* all memory available again */
```

`forge` compiles `main.c` (provided harness) + `arena.c`, runs `./test`, and
diffs stdout. The arena lives in your `arena.c` as a `static` 512-byte buffer —
nothing else global. The contract:

- Every returned pointer must be **8-byte aligned**, even after calls like
  `arena_alloc(3)`. Pad the offset up to an 8-boundary, and prefer
  `_Alignas(8)` on the buffer so the base never wastes a byte.
- Never hand out memory past the end of the 512 bytes: when `n` doesn't fit,
  return `NULL`, and `arena_alloc(0)` returns `NULL` too.
- No `malloc`, no `mmap` — the arena is the only source of memory.
- `arena_reset` makes the whole 512 bytes available again (no clearing needed).
- No `CFLAGS` redefines.

## Acceptance

Graded `build`: strict compile + harness stdout diff (whitespace-insensitive),
exit 0, empty stderr. The harness verifies alignment, exact 512-byte
exhaustion, tag-based non-overlap, reset semantics, and the NULL-on-full
contract:

- `arena_alloc(3)` + `(5)` + `(7)` + `(1)` all succeed, all aligned.
- 60 more `arena_alloc(8)` calls succeed; the *next* call returns `NULL`
  (512 bytes exact).
- `arena_left()` reports exactly 0 at exhaustion.
- No block overwrites its neighbor (the harness tags every block and re-reads
  them); a bumper that overwrites its boundary fails that check and likely
  trips ASan's global redzone.
- After `arena_reset`: `arena_left() == 512`, `arena_alloc(512)` succeeds, and
  `arena_alloc(1)` fails cleanly.

`quiz.txt` is complete (see Quiz).

## Readings

- **Reading ladder** — start with `man 3 malloc`'s contract, then the jemalloc
  arena overview below, then TLPI §7.1 for the heap model.
- The Linux Programming Interface, §7.1 "Allocating Memory on the Heap" is the background;
  for the *arena* pattern itself, read the jemalloc overview below.
- jemalloc's arena overview (skim): https://jemalloc.net/jemalloc.3.html search "arena".
- `_Alignas` / `_Alignof`: cppreference https://en.cppreference.com/w/c/language/_Alignas
- Basic pointer-arithmetic refresher if needed: K&R §5.4 "Address Arithmetic".

## Quiz

1. An allocator that assigns memory only by advancing a pointer is called what?
2. Which C11 keyword forces a static buffer to a required alignment?