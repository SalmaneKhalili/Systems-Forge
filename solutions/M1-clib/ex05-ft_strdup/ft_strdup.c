#include <stdlib.h>

#include "ft_strdup.h"

char *ft_strdup(const char *s)
{
	size_t n = 0;
	char *out;
	size_t i;

	while (s[n])
		n++;
	out = malloc(n + 1);
	if (out == NULL)
		return NULL;
	for (i = 0; i < n; i++)
		out[i] = s[i];
	out[n] = '\0';
	return out;
}
