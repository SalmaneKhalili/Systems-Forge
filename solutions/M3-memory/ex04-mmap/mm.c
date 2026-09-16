#include <stddef.h>
#include <stdint.h>
#include <sys/mman.h>

#include "mm.h"

#define PAGE 4096UL
#define MAXLIVE 64

static struct region {
	void *addr;
	size_t size;
} live[MAXLIVE];
static size_t nlive;
static size_t outstanding;

static size_t roundp(size_t n)
{
	return ((n + PAGE - 1) / PAGE) * PAGE;
}

void *mm_calloc(size_t n)
{
	size_t sz;
	void *m;

	if (n == 0 || nlive >= MAXLIVE)
		return NULL;
	sz = roundp(n);
	m = mmap(NULL, sz, PROT_READ | PROT_WRITE,
		 MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
	if (m == MAP_FAILED)
		return NULL;
	live[nlive].addr = m;
	live[nlive].size = sz;
	nlive++;
	outstanding += sz;
	return m;
}

void mm_delloc(void *p)
{
	size_t i;

	if (p == NULL)
		return;
	for (i = 0; i < nlive; i++) {
		if (live[i].addr == p) {
			munmap(p, live[i].size);
			outstanding -= live[i].size;
			live[i] = live[nlive - 1];
			nlive--;
			return;
		}
	}
}

size_t mm_allocated(void)
{
	return outstanding;
}