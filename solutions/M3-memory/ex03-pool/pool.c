#include "pool.h"

#define NSLOTS 8
#define SLOT 32

static _Alignas(8) unsigned char cells[NSLOTS][SLOT];
static unsigned char inuse;

void *pool_get(void)
{
	size_t i;
	for (i = 0; i < NSLOTS; i++)
		if (!(inuse & (unsigned char)(1u << i))) {
			inuse |= (unsigned char)(1u << i);
			return cells[i];
		}
	return NULL;
}

void pool_put(void *slot)
{
	unsigned char *p = slot;
	size_t i;

	if (slot == NULL)
		return;
	for (i = 0; i < NSLOTS; i++) {
		if (p == cells[i]) {
			inuse &= (unsigned char)~(1u << i);
			return;
		}
	}
}

size_t pool_borrowed(void)
{
	size_t n = 0, i;
	for (i = 0; i < NSLOTS; i++)
		if (inuse & (unsigned char)(1u << i))
			n++;
	return n;
}