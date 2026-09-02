#include <stdint.h>
#include <stdio.h>

#include "pool.h"

static int failures;

static int aligned8(const void *p)
{
	return ((uintptr_t)p % 8) == 0;
}

static void tag_blocks(void *slots[8])
{
	size_t i, j;
	for (i = 0; i < 8; i++)
		for (j = 0; j < 32; j++)
			((unsigned char *)slots[i])[j] = (unsigned char)(i * 32 + j);
}

static int tag_ok(void *slots[8])
{
	size_t i, j;
	for (i = 0; i < 8; i++)
		for (j = 0; j < 32; j++)
			if (((unsigned char *)slots[i])[j] != (unsigned char)(i * 32 + j))
				return 0;
	return 1;
}

static int in_set(const void *p, void *set[8], size_t n)
{
	size_t i;
	for (i = 0; i < n; i++)
		if (p == set[i])
			return 1;
	return 0;
}

int main(void)
{
	void *slots[8];
	void *u[3];
	void *x;
	size_t i, r;
	unsigned char dummy;

	if (pool_borrowed() != 0) {
		printf("FAIL: pool not empty at start\n");
		failures++;
	} else {
		printf("pool empty ok\n");
	}

	for (i = 0; i < 8; i++) {
		slots[i] = pool_get();
		if (slots[i] == NULL) {
			printf("FAIL: get %zu returned NULL\n", i);
			failures++;
		}
	}
	x = pool_get();
	if (x != NULL) {
		printf("FAIL: 9th get handed out memory\n");
		failures++;
	} else {
		printf("9th get NULL ok\n");
	}
	if (pool_borrowed() != 8) {
		printf("FAIL: borrowed %zu want 8\n", pool_borrowed());
		failures++;
	} else {
		printf("borrowed 8 ok\n");
	}

	for (i = 0; i < 8; i++) {
		if (!aligned8(slots[i])) {
			printf("FAIL: slot %zu unaligned\n", i);
			failures++;
		}
	}
	if (failures == 0)
		printf("slots aligned ok\n");

	for (i = 0; i < 8; i++)
		for (r = i + 1; r < 8; r++)
			if (slots[i] == slots[r]) {
				printf("FAIL: slots %zu and %zu are the same block\n", i, r);
				failures++;
			}

	tag_blocks(slots);
	if (tag_ok(slots)) {
		printf("isolation ok\n");
	} else {
		printf("FAIL: slots overlap\n");
		failures++;
	}

	pool_put(slots[0]);
	pool_put(slots[1]);
	pool_put(slots[2]);
	if (pool_borrowed() != 5) {
		printf("FAIL: borrowed after 3 puts = %zu want 5\n", pool_borrowed());
		failures++;
	}
	x = pool_get();
	if (x == NULL || !aligned8(x)) {
		printf("FAIL: refill get\n");
		failures++;
	}
	pool_put(x);
	if (pool_borrowed() != 5) {
		printf("FAIL: borrowed after refill = %zu\n", pool_borrowed());
		failures++;
	}

	pool_put(NULL);
	pool_put(&dummy);
	if (pool_borrowed() != 5) {
		printf("FAIL: invalid put changed the count\n");
		failures++;
	} else {
		printf("invalid puts ignored ok\n");
	}

	r = 0;
	while (r < 100) {
		int round_ok = 1;
		pool_put(slots[3]);
		pool_put(slots[4]);
		pool_put(slots[5]);
		for (i = 0; i < 3; i++) {
			u[i] = pool_get();
			if (u[i] == NULL || !aligned8(u[i]) ||
			    !in_set(u[i], slots, 6)) {
				round_ok = 0;
				break;
			}
		}
		if (u[0] == u[1] || u[1] == u[2] || u[0] == u[2])
			round_ok = 0;
		for (i = 0; i < 3; i++)
			pool_put(u[i]);
		if (!round_ok) {
			printf("FAIL: churn round %zu\n", r);
			failures++;
			break;
		}
		r++;
	}
	if (r == 100)
		printf("churn 100 ok\n");

	pool_put(slots[6]);
	pool_put(slots[7]);
	if (pool_borrowed() != 0) {
		printf("FAIL: borrowed after release = %zu\n", pool_borrowed());
		failures++;
	} else {
		printf("release all ok\n");
	}

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}