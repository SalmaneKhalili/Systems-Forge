#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "ft_strdup.h"

static int failures;

static void one(const char *name, const char *src)
{
	char *cp;

	cp = ft_strdup(src);
	if (cp == NULL) {
		printf("%s: FAIL NULL\n", name);
		failures++;
		return;
	}
	if (cp == src) {
		printf("%s: FAIL returned source address\n", name);
		free(cp);
		failures++;
		return;
	}
	if (strcmp(cp, src) != 0) {
		printf("%s: FAIL content mismatch\n", name);
		free(cp);
		failures++;
		return;
	}
	cp[0] = cp[0] == 'z' ? 'y' : 'z';
	if (strcmp(cp, src) == 0) {
		printf("%s: FAIL copy is not independent\n", name);
		free(cp);
		failures++;
		return;
	}
	free(cp);
	printf("%s: PASS\n", name);
}

int main(void)
{
	one("empty", "");
	one("hello", "hello");
	one("spaces", "  two spaces  ");
	one("long", "0123456789abcdef0123456789abcdef");
	one("high-bit", "caf\xc3\xa9 \xe2\x98\x83");

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}