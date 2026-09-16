#include <stdio.h>

#include "ft_strcmp.h"

static int failures;

static void one(const char *name, const char *a, const char *b, int sig)
{
	int ok;
	int r = ft_strcmp(a, b);

	if (sig == 0)
		ok = (r == 0);
	else if (sig < 0)
		ok = (r < 0);
	else
		ok = (r > 0);

	if (ok)
		printf("%s: PASS\n", name);
	else {
		printf("%s: FAIL got %d\n", name, r);
		failures++;
	}
}

int main(void)
{
	one("equal", "hello", "hello", 0);
	one("less", "abc", "abd", -1);
	one("more", "abd", "abc", 1);
	one("prefix", "abc", "abcd", -1);
	one("case", "A", "a", -1);
	one("empty-a", "", "a", -1);
	one("empty-b", "a", "", 1);
	one("high-bit", "\xff", "\x00", 1);
	one("tail-equal", "same", "same", 0);

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}