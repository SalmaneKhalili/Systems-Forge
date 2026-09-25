# M3-ex03 · Fixed Slot Pool

The bump arena hands out variable sizes; a **fixed slot pool** trades that for
speed: 8 slots × 32 bytes, where any call to `get` either hands out a whole
unused slot or returns `NULL`. Object pools are how real systems avoid
per-object `malloc` churn — web servers pool request structs, kernels pool
descriptors, jemalloc pool-allocates. You implement `pool.c` behind the
provided header.

## Shape

The exercise provides `pool.h` and a harness `main.c`; you write **`pool.c`**:

```c
void  *pool_get(void);        /* next free slot, or NULL when all 8 are borrowed */
void   pool_put(void *slot);  /* return a borrowed slot to the pool */
size_t pool_borrowed(void);   /* count of slots currently handed out */
```

`forge` compiles `main.c` (harness) + `pool.c`, runs `./test`, and diffs
stdout. The pool is your `pool.c`'s static 8×32 memory. The contract:

- `_Alignas(8)` the backing array; every slot pointer must be 8-aligned so the
  harness's writes are always legal.
- `pool_put` must **validate** its argument: `NULL` and any pointer that is
  not one of the 8 slots is rejected *silently* (no output, no crash, no state
  change). The harness passes a stack variable to prove it.
- A slot may be borrowed only once (`get` never returns a slot that is out).
- No `malloc`, no `mmap`; the pool is the only memory.
- No `CFLAGS` redefines.

## Acceptance

Graded `build`: strict compile + harness stdout diff, exit 0, empty stderr.
The harness tests the 8-slot bound, NULL-on-empty, isolation by tagging,
borrow accounting, invalid-put rejection, and 100-round churn:

- 8 gets in a row succeed; the 9th returns `NULL`; `pool_borrowed() == 8`.
- Writing distinct bytes into every slot and re-reading proves no slot overlaps
  another; a pool that lets two gets hand out the same slot fails the
  isolation check.
- Returning 3 slots makes `pool_borrowed() == 5`, and a get works again.
- `pool_put(NULL)` and `pool_put(&stackvar)` are ignored without crash or
  count change.
- Repeated get-all/put-all churn never shrinks capacity; a pool that lacks
  bounds-checking on put crashes the churn round.

`quiz.txt` is complete (see Quiz).

## Readings

- **Reading ladder** — start with the object-pool rationale and `_Alignas`
  (references below), then TLPI §7.1 for why malloc is not free.
- The Linux Programming Interface, §7.1 "Allocating Memory on the Heap" for why malloc
  (in `glibc` via its own pool structure) is not free — the motivating background of pools.
- Object pooling rationale:
  https://en.wikipedia.org/wiki/Object_pool_pattern (skim; the "Benefits" section).
- Hash-free, allocation-free bookkeeping: think bitset, not linked list — that is the lesson.
- `_Alignas`/`_Alignof` again: cppreference https://en.cppreference.com/w/c/language/_Alignas

## Quiz

1. What does an empty pool return from get?
2. How is a slot pointer proven to belong to the pool?