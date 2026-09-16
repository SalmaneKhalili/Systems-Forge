#include "ft_strlen.h"

size_t ft_strlen(const char *s)
{
	const char *p = s;

	while (*p)
		p++;
	return (size_t)(p - s);
}
