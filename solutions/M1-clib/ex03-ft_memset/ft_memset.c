#include "ft_memset.h"

void *ft_memset(void *s, int c, size_t n)
{
	unsigned char *p = s;
	size_t i;

	for (i = 0; i < n; i++)
		p[i] = (unsigned char)c;
	return s;
}
