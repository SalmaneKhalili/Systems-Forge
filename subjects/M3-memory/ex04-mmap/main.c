#include <stdint.h>
#include <stdio.h>

#include "mm.h"

#define P ((size_t)4096)

static int failures;

static int aligned(const void *p)
{
	return ((uintptr_t)p % P) == 0;
}

static int ranges_disjoint(uintptr_t a, size_t as, uintptr_t b, size_t bs)
{
	return b >= a + as || a >= b + bs;
}

int main(void)
{
	unsigned char *a, *b, *c;
	unsigned char *churn[40];
	size_t expect, i;
	unsigned char dummy;
	int ok;

	a = mm_calloc(1);
	if (a == NULL || !aligned(a)) {
		printf("FAIL: mm_calloc(1)\n");
		failures++;
	}
	ok = 1;
	for (i = 0; i < 16; i++)
		if (a[i] != 0)
			ok = 0;
	if (!ok) {
		printf("FAIL: region A not zeroed\n");
		failures++;
	} else {
		printf("calloc zeroed ok\n");
	}

	b = mm_calloc(1500);
	if (b == NULL || !aligned(b) ||
	    !ranges_disjoint((uintptr_t)a, P, (uintptr_t)b, P)) {
		printf("FAIL: mm_calloc(1500) bad\n");
		failures++;
	}

	c = mm_calloc(8200);
	if (c == NULL || !aligned(c) ||
	    !ranges_disjoint((uintptr_t)a, P, (uintptr_t)c, 3 * P) ||
	    !ranges_disjoint((uintptr_t)b, P, (uintptr_t)c, 3 * P)) {
		printf("FAIL: mm_calloc(8200) bad\n");
		failures++;
	}
	if (failures == 0)
		printf("regions aligned and disjoint ok\n");

	if (a != NULL && b != NULL && c != NULL) {
	for (i = 0; i < 64; i++)
		a[i] = (unsigned char)i;
	for (i = 0; i < 1500; i++)
		b[i] = (unsigned char)(i * 7);
	for (i = 0; i < 8200; i++)
		c[i] = (unsigned char)(i * 13);
	ok = 1;
	for (i = 0; i < 64 && ok; i++)
		if (a[i] != (unsigned char)i)
			ok = 0;
	for (i = 0; i < 1500 && ok; i++)
		if (b[i] != (unsigned char)(i * 7))
			ok = 0;
	for (i = 0; i < 8200 && ok; i++)
		if (c[i] != (unsigned char)(i * 13))
			ok = 0;
	if (!ok) {
		printf("FAIL: region contents corrupted\n");
		failures++;
	} else {
		printf("region contents ok\n");
	}
	} else {
		printf("FAIL: cannot test contents\n");
		failures++;
	}

	if (mm_allocated() != 5 * P) {
		printf("FAIL: allocated %zu want %zu\n", mm_allocated(), 5 * P);
		failures++;
	} else {
		printf("accounting rounds to pages ok\n");
	}

	mm_delloc(b);
	if (mm_allocated() != 4 * P) {
		printf("FAIL: after free allocated %zu want %zu\n", mm_allocated(), 4 * P);
		failures++;
	}
	mm_delloc(NULL);
	mm_delloc(&dummy);
	if (mm_allocated() != 4 * P) {
		printf("FAIL: invalid free changed accounting\n");
		failures++;
	} else {
		printf("invalid free ignored ok\n");
	}

	expect = 4 * P;
	for (i = 0; i < 40; i++) {
		churn[i] = mm_calloc(129);
		if (churn[i] == NULL || !aligned(churn[i])) {
			printf("FAIL: churn alloc %zu\n", i);
			failures++;
			break;
		}
		expect += P;
		if (mm_allocated() != expect) {
			printf("FAIL: churn accounting %zu\n", i);
			failures++;
			break;
		}
	}
	for (i = 0; i < 40; i++)
		if (churn[i])
			mm_delloc(churn[i]);
	if (mm_allocated() != 4 * P) {
		printf("FAIL: churn frees left %zu\n", mm_allocated());
		failures++;
	}
	if (failures == 0)
		printf("churn 40 ok\n");

	mm_delloc(c);
	if (mm_allocated() != P) {
		printf("FAIL: after C free allocated %zu\n", mm_allocated());
		failures++;
	}
	mm_delloc(a);
	if (mm_allocated() != 0) {
		printf("FAIL: not all released (%zu)\n", mm_allocated());
		failures++;
	} else {
		printf("clean exit ok\n");
	}

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}