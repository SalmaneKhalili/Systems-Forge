#include <stdio.h>
#include <string.h>

#include "ft_memset.h"

static int failures;

static void one(const char *name, int c, size_t off, size_t n)
{
	unsigned char buf[16];
	unsigned char *p;
	size_t i;

	for (i = 0; i < sizeof(buf); i++)
		buf[i] = 0xaa;

	p = ft_memset(buf + off, c, n);

	if (p != buf + off) {
		printf("%s: FAIL return pointer\n", name);
		failures++;
		return;
	}
	for (i = off; i < off + n; i++) {
		if (buf[i] != (unsigned char)c) {
			printf("%s: FAIL byte %zu is %02x want %02x\n",
			       name, i, buf[i], (unsigned char)c);
			failures++;
			return;
		}
	}
	for (i = 0; i < off; i++) {
		if (buf[i] != 0xaa) {
			printf("%s: FAIL prefix byte %zu touched\n", name, i);
			failures++;
			return;
		}
	}
	for (i = off + n; i < sizeof(buf); i++) {
		if (buf[i] != 0xaa) {
			printf("%s: FAIL suffix byte %zu touched\n", name, i);
			failures++;
			return;
		}
	}
	printf("%s: PASS\n", name);
}

int main(void)
{
	one("fill-ff", 0xff, 0, 6);
	one("fill-00", 0, 2, 4);
	one("middle", 0x01, 8, 3);
	one("edge-end", 0x7f, 13, 3);
	one("zero-n", 0x42, 3, 0);

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}