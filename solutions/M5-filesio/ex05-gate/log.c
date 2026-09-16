#include <string.h>
#include <unistd.h>

#include "fix.h"

#define BUF 4096

int parse_log(int fd, Tokens *t)
{
	char buf[BUF];
	char *line, *end, *eq;
	size_t used = 0;
	ssize_t n;

	t->n = 0;
	while ((n = read(fd, buf + used, sizeof(buf) - used)) > 0)
		used += (size_t)n;
	if (n < 0)
		return -1;
	if (used == 0)
		return 0;

	line = buf;
	do {
		end = memchr(line, '\n', (size_t)((buf + used) - line));
		if (end != NULL)
			*end = '\0';
		if (line != end) {
			eq = strchr(line, '=');
			if (eq == NULL || eq == line)
				return -1;
			if ((size_t)(eq - line) >= TOKEN_LEN)
				return -1;
			memcpy(t->keys[t->n], line, (size_t)(eq - line));
			t->keys[t->n][(size_t)(eq - line)] = '\0';
			t->n++;
			if (t->n >= TOKENS_MAX)
				return -1;
		}
		if (end == NULL)
			break;
		line = end + 1;
	} while (line < (buf + used));
	return 0;
}