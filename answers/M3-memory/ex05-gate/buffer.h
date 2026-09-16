#ifndef BUFFER_H
#define BUFFER_H

#include <stddef.h>

#define BUF_MAX 1024

typedef struct Buffer {
	size_t len;
	size_t cap;
} Buffer;

void buf_init(Buffer *b);
int buf_put(Buffer *b, const char *s, size_t n);
size_t buf_len(const Buffer *b);
const char *buf_data(const Buffer *b);

#endif