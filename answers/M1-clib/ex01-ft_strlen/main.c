#include <stdio.h>
#include <string.h>

#include "ft_strlen.h"

static int failures;

static void one(const char *name, const char *s)
{
	size_t got = ft_strlen(s);
	size_t want = strlen(s);

	if (got == want)
		printf("%s: PASS (%zu)\n", name, got);
	else {
		printf("%s: FAIL got %zu want %zu\n", name, got, want);
		failures++;
	}
}

int main(void)
{
	one("empty", "");
	one("hello", "hello");
	one("with-spaces", "with spaces");
	one("utf8", "\xc3\xa9\xe2\x98\x83");
	one("embedded-nul", "a\x00" "b");
	one("long", "0123456789abcdef0123456789abcdef0123456789abcdef");

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}