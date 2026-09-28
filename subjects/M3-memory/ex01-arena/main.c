#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "arena.h"

static int failures;

static const size_t small_w[4] = {3, 5, 7, 1};

static void expect_aligned(const char* tag, void* p)
{
    if (((uintptr_t)p % 8) != 0)
    {
        printf("FAIL %s: unaligned %p\n", tag, p);
        failures++;
    }
}

int main(void)
{
    void* a[64];
    size_t widths[64];
    size_t ok8 = 0;
    size_t blocks;
    size_t i, j;
    int small_ok = 1;

    if (arena_left() != 512)
    {
        printf("FAIL: arena does not start full (got %zu)\n", arena_left());
        failures++;
    }

    for (i = 0; i < 4; i++)
    {
        a[i] = arena_alloc(small_w[i]);
        widths[i] = small_w[i];
        if (a[i] == NULL)
        {
            printf("FAIL: arena_alloc(%zu) NULL\n", small_w[i]);
            failures++;
            small_ok = 0;
            continue;
        }
        expect_aligned("small", a[i]);
    }
    if (small_ok)
        printf("align ok\n");
    else
        goto overlap_done;

    for (; ok8 < 60; ok8++)
    {
        a[4 + ok8] = arena_alloc(8);
        if (a[4 + ok8] == NULL)
            break;
        widths[4 + ok8] = 8;
    }
    blocks = 4 + ok8;
    if (blocks != 64)
    {
        printf("FAIL: arena exhausted after %zu blocks (want 64)\n", blocks);
        failures++;
    }
    if (arena_left() != 0)
    {
        printf("FAIL: arena_left()=%zu want 0\n", arena_left());
        failures++;
    }
    if (arena_alloc(1) != NULL)
    {
        printf("FAIL: expected NULL after exhaustion\n");
        failures++;
    }

    for (i = 0; i < blocks; i++)
    {
        unsigned char* w = a[i];
        for (j = 0; j < widths[i]; j++)
            w[j] = (unsigned char)(i + 1);
    }
    for (i = 0; i < blocks; i++)
    {
        unsigned char* w = a[i];
        for (j = 0; j < widths[i]; j++)
        {
            if (w[j] != (unsigned char)(i + 1))
            {
                printf("FAIL: block %zu corrupted at %zu\n", i, j);
                failures++;
                goto overlap_done;
            }
        }
    }
    printf("no overlap\n");
overlap_done:

    arena_reset();
    if (arena_left() != 512)
    {
        printf("FAIL: arena_left() after reset = %zu\n", arena_left());
        failures++;
    }
    else
    {
        printf("reset ok\n");
    }
    if (arena_alloc(512) == NULL)
    {
        printf("FAIL: exact-fit 512 allocation\n");
        failures++;
    }
    else
    {
        printf("exact fit ok\n");
    }
    if (arena_alloc(1) != NULL)
    {
        printf("FAIL: 513th byte handed out\n");
        failures++;
    }
    arena_reset();
    if (arena_alloc(513) != NULL)
    {
        printf("FAIL: arena_alloc(513) should be NULL\n");
        failures++;
    }

    if (failures == 0)
        printf("ALL PASS\n");
    else
        printf("SOME FAIL\n");
    return failures == 0 ? 0 : 1;
}
