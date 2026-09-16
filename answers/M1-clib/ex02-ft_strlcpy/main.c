#include <stdio.h>
#include <string.h>

#include "ft_strlcpy.h"

static int failures;

static void one(const char *name, const char *src, size_t size)
{
	char dst[16];
	size_t i;
	size_t n;
	size_t got;
	size_t copy;

	for (i = 0; i < sizeof(dst); i++)
		dst[i] = 'x';

	n = strlen(src);
	got = ft_strlcpy(dst, src, size);

	if (got != n) {
		printf("%s: FAIL return %zu want %zu\n", name, got, n);
		failures++;
		return;
	}
	if (size > 0) {
		copy = n < size - 1 ? n : size - 1;
		if (memcmp(dst, src, copy) != 0) {
			printf("%s: FAIL copied bytes\n", name);
			failures++;
			return;
		}
		if (dst[copy] != '\0') {
			printf("%s: FAIL not NUL-terminated at index %zu\n", name, copy);
			failures++;
			return;
		}
		for (i = copy + 1; i < sizeof(dst); i++) {
			if (dst[i] != 'x') {
				printf("%s: FAIL tail byte %zu touched\n", name, i);
				failures++;
				return;
			}
		}
	} else if (memcmp(dst, "xxxxxxxxxxxxxxxx", sizeof(dst)) != 0) {
		printf("%s: FAIL wrote anything at size 0\n", name);
		failures++;
		return;
	}
	printf("%s: PASS\n", name);
}

int main(void)
{
	one("full", "hello", 16);
	one("tight", "hello", 3);
	one("zero", "hello", 0);
	one("empty-dst", "hello", 1);
	one("exact", "hello world", 12);
	one("truncated", "hello world", 6);

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}