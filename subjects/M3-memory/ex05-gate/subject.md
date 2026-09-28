# M3-ex05 gate · OOM-safe buffer

## Goal

Implement `buffer.c` backing the provided `buffer.h`: a **growable byte buffer with a hard
capacity**. Real servers die when memory is exhausted; production code survives because it
*refuses*. This gate's hard cap is a machine-checkable stand-in for malloc failure: when the
buffer cannot grow any further, `buf_put` returns `-1` — cleanly, without touching the data
already stored.

```c
typedef struct Buffer {
	size_t len;
	size_t cap;           /* current claimed capacity, a power of two in [64, BUF_MAX] */
} Buffer;

void          buf_init(Buffer *b);   /* len 0, cap 64 */
int           buf_put(Buffer *b, const char *s, size_t n);  /* 0 ok | -1 refused */
size_t        buf_len(const Buffer *b);
const char   *buf_data(const Buffer *b);  /* the stored bytes, growing/copy-safe */
```

`forge` compiles `main.c` (harness) + `buffer.c`, runs `./test`, diffs stdout.

## Constraints

- The buffer's memory lives in `buffer.c` as a **static 1024-byte arena** (`_Alignas(8)`),
  the only storage. `cap` starts at 64 and, on growth, doubles up to `BUF_MAX` (1024).
- Growth must **preserve all earlier bytes**: `buf_data` must always return the exact
  sequence previously appended, in order. Copy when you grow — nothing is handed to you for
  free.
- `buf_put` must refuse without side effects: when a push would need more than `BUF_MAX`
  bytes, return `-1` **leaving `len` and the stored bytes exactly as they were**.
- `buf_put(b, NULL, 0)` succeeds trivially (returns 0, unchanged); a non-NULL `n == 0` too.
  A NULL source with `n > 0` is refused.
- No `malloc`, no `mmap`. No `CFLAGS` redefines.
- The queue after refusal must still work: `buf_put` smaller data succeeds again.

## Acceptance criteria

- [ ] starts at `len 0, cap 64`
- [ ] five appends of 200 bytes raise `len` to 1000 across several growths, and every step
      the stored bytes match the source exactly
- [ ] after a final 24-byte append, `len == 1024` and `cap == 1024`
- [ ] a 1025th byte is refused with `-1`; `len` stays 1024 and all 1024 bytes are intact
- [ ] after `buf_init`, the buffer accepts appends again (a full buffer stays valid until
      reset)
- [ ] `buf_put(b, NULL, 0)` and `buf_put(b, NULL, 5)` leave it unchanged (0 / -1)
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
How must a bounded allocator signal that memory is exhausted?: <answer>
What is the buffer's hard capacity in bytes?: <answer>
```

## Readings

- **Reading ladder** — start with `man 3 realloc`'s contract, then TLPI §7.1
  below for why allocators fail.
- The Linux Programming Interface, §7.1 "Allocating Memory on the Heap" — why every real
  allocator can fail, and what robust code must do. An allocation failure is a `NULL`; this
  exercise makes that failure deterministic.
- `realloc(3)`'s contract is the model: "The contents of the object shall be unchanged in
  the range from the start to the common part of the old and new sizes." `man 3 realloc`.
- How production crawlers/servers structure bounded buffers (skim; the pattern is yours):
  https://github.com/valyala/bytebufferpool (Go) — "get habit" of the pattern, not the code.

## How you are graded

- `build`: strict compile + harness stdout diff, exit 0, empty stderr. The harness forces
  growth past every power-of-two and then *past the cap*, checking data integrity at each
  step and the no-side-effect refusal at the end. A student that ignores refusal overruns
  the 1024-byte arena and trips ASan's redzone; a student that drops bytes on growth fails
  the content checks.
- `quiz`: `quiz.txt` answers must match.