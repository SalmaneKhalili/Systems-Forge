#include <stdio.h>
#include <string.h>

#include "buffer.h"

static int failures;

static void expect(const char *name, int cond)
{
	if (cond) {
		printf("%s ok\n", name);
	} else {
		printf("FAIL %s\n", name);
		failures++;
	}
}

static int pow2(size_t v)
{
	return v != 0 && (v & (v - 1)) == 0;
}

int main(void)
{
	static char pattern[1024];
	static char expected[1024];
	Buffer b;
	size_t i;
	int r;

	for (i = 0; i < 1024; i++)
		pattern[i] = (char)("0123456789abcdefghijklmnopqrstuvwxyz"[i % 36]);
	for (i = 0; i < 1000; i++)
		expected[i] = pattern[i];
	for (i = 0; i < 24; i++)
		expected[1000 + i] = pattern[i];

	buf_init(&b);
	expect("init", b.len == 0 && b.cap == 64 && buf_len(&b) == 0);

	for (i = 1; i <= 5; i++) {
		r = buf_put(&b, pattern + 200 * (i - 1), 200);
		if (r != 0) {
			printf("FAIL chunk %zu\n", i);
			failures++;
			goto content_done;
		}
		if (buf_len(&b) != 200 * i || b.cap < 200 * i || !pow2(b.cap) ||
		    memcmp(buf_data(&b), pattern, 200 * i) != 0) {
			printf("FAIL grow after chunk %zu (len %zu cap %zu)\n",
			       i, buf_len(&b), b.cap);
			failures++;
			goto content_done;
		}
	}
	printf("grow integrity ok\n");
content_done:

	if (buf_len(&b) == 1000) {
		r = buf_put(&b, pattern, 24);
		if (r != 0 || buf_len(&b) != 1024 || b.cap != 1024) {
			printf("FAIL final append (r=%d len=%zu cap=%zu)\n",
			       r, buf_len(&b), b.cap);
			failures++;
		} else {
			printf("full 1024 cap ok\n");
		}
	}

	if (buf_len(&b) == 1024) {
		if (memcmp(buf_data(&b), expected, 1024) != 0) {
			printf("FAIL contents before refusal\n");
			failures++;
		}
		r = buf_put(&b, pattern, 1);
		if (r != -1 || buf_len(&b) != 1024) {
			printf("FAIL refusal (r=%d len=%zu)\n", r, buf_len(&b));
			failures++;
		} else if (memcmp(buf_data(&b), expected, 1024) == 0) {
			printf("refusal no side effect ok\n");
		} else {
			printf("FAIL refusal corrupted data\n");
			failures++;
		}
		r = buf_put(&b, pattern, 200);
		if (r != -1) {
			printf("FAIL big append after cap\n");
			failures++;
		}
	}

	r = buf_put(&b, NULL, 0);
	expect("NULL0", r == 0 && buf_len(&b) == 1024);
	r = buf_put(&b, NULL, 5);
	expect("NULL5 refused", r == -1 && buf_len(&b) == 1024);

	buf_init(&b);
	r = buf_put(&b, pattern, 8);
	if (r != 0 || buf_len(&b) != 8 ||
	    memcmp(buf_data(&b), pattern, 8) != 0) {
		printf("FAIL reuse after refusal\n");
		failures++;
	} else {
		printf("reuse after refusal ok\n");
	}

	if (failures == 0)
		printf("ALL PASS\n");
	else
		printf("SOME FAIL\n");
	return failures == 0 ? 0 : 1;
}