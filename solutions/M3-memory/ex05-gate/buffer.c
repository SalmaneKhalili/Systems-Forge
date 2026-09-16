#include <string.h>

#include "buffer.h"

static _Alignas(8) char store[BUF_MAX];

void buf_init(Buffer *b)
{
	b->len = 0;
	b->cap = 64;
}

size_t buf_len(const Buffer *b)
{
	return b->len;
}

const char *buf_data(const Buffer *b)
{
	(void)b;
	return store;
}

static int grow(Buffer *b, size_t need)
{
	while (need > b->cap) {
		if (b->cap == BUF_MAX)
			return -1;
		size_t nc = b->cap * 2;
		if (nc > BUF_MAX)
			nc = BUF_MAX;
		b->cap = nc;
	}
	return 0;
}

int buf_put(Buffer *b, const char *s, size_t n)
{
	if (n == 0)
		return 0;
	if (s == NULL)
		return -1;
	if (n > BUF_MAX - b->len)
		return -1;
	if (grow(b, b->len + n) != 0)
		return -1;
	memmove(store + b->len, s, n);
	b->len += n;
	return 0;
}