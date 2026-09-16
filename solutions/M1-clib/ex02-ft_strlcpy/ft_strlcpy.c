#include "ft_strlcpy.h"

size_t ft_strlcpy(char *dst, const char *src, size_t size)
{
	size_t n = 0;
	size_t i;

	while (src[n])
		n++;
	if (size > 0) {
		size_t copy = n < size - 1 ? n : size - 1;
		for (i = 0; i < copy; i++)
			dst[i] = src[i];
		dst[copy] = '\0';
	}
	return n;
}
