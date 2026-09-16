#include "arena.h"

static _Alignas(8) unsigned char memory[512];
static size_t used;

void *arena_alloc(size_t n)
{
	size_t aligned = (n + 7) & ~(size_t)7;

	if (n == 0 || aligned > 512 - used)
		return NULL;
	void *p = memory + used;
	used += aligned;
	return p;
}

size_t arena_left(void)
{
	return 512 - used;
}

void arena_reset(void)
{
	used = 0;
}