#ifndef FIX_H
#define FIX_H

#include <stddef.h>

#define TOKENS_MAX 64
#define TOKEN_LEN 32

typedef struct {
	char keys[TOKENS_MAX][TOKEN_LEN];
	size_t n;
} Tokens;

int parse_log(int fd, Tokens *t);

#endif