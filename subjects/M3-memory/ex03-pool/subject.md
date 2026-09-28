# M3-ex03 · fixed slot pool

## Goal

Implement `pool.c` backing the provided `pool.h` — a **fixed slot pool**: 8 slots × 32
bytes, where any call to `get` either hands out a whole unused slot or returns `NULL`. Object
pools are how real systems avoid per-object `malloc` churn (web servers pool request
structs, kernels pool descriptors, jemalloc pool-allocates).

```c
void  *pool_get(void);        /* next free slot, or NULL when all 8 are borrowed */
void   pool_put(void *slot);  /* return a borrowed slot to the pool */
size_t pool_borrowed(void);   /* count of slots currently handed out */
```

`forge` compiles `main.c` (harness) + `pool.c`, runs `./test`, diffs stdout.

## Constraints

- The pool is your `pool.c` static 8×32 memory; `_Alignas(8)` the backing array; every slot
  pointer must be 8-aligned (so the harness's writes are always legal).
- `pool_put` must **validate** its argument: `NULL` and any pointer that is not one of the 8
  slots is rejected *silently* (no output, no crash, no state change). The harness passes a
  stack variable to prove it.
- A slot may be borrowed only once (get never returns a slot that is out).
- No `malloc`, no `mmap`; the pool is the only memory.
- No `CFLAGS` redefines.

## Acceptance criteria

- [ ] 8 gets in a row succeed; the 9th returns `NULL`; `pool_borrowed() == 8`
- [ ] writing distinct bytes into every slot and re-reading proves no slot overlaps another
- [ ] returning 3 slots makes `pool_borrowed() == 5`, and a get works again
- [ ] `pool_put(NULL)` and `pool_put(&stackvar)` are ignored without crash or count change
- [ ] repeated get-all/put-all churn never shrinks capacity
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What does an empty pool return from get?: <answer>
How is a slot pointer proven to belong to the pool?: <answer>
```

## Readings

- **Reading ladder** — start with the object-pool rationale and `_Alignas`
  (references below), then TLPI §7.1 for why malloc is not free.
- "The Linux Programming Interface", §7.1 "Allocating Memory on the Heap" for why malloc
  (in `glibc` via its own pool structure) is not free — the motivating background of pools.
- Object pooling rationale:
  https://en.wikipedia.org/wiki/Object_pool_pattern (skim; the "Benefits" section).
- Hash-free, allocation-free bookkeeping: think bitset, not linked list — that is the lesson.
- `_Alignas`/`_Alignof` again: cppreference https://en.cppreference.com/w/c/language/_Alignas

## How you are graded

- `build`: strict compile + harness stdout diff, exit 0, empty stderr. The harness tests the
  8-slot bound, NULL-on-empty, isolation by tagging, borrow accounting, invalid-put
  rejection, and 100-round churn. A pool that lets two gets hand out the same slot fails the
  isolation check; one that lacks bounds-checking on put crashes the churn round.