#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

#include "fix.h"

int main(int argc, char **argv)
{
	Tokens t;
	int fd, i;

	if (argc != 2) {
		fprintf(stderr, "usage: %s FILE\n", argv[0]);
		return 1;
	}
	fd = open(argv[1], O_RDONLY);
	if (fd < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	if (parse_log(fd, &t) != 0) {
		fprintf(stderr, "parse failed\n");
		return 1;
	}
	close(fd);
	for (i = 0; i < (int)t.n; i++)
		printf("%s\n", t.keys[i]);
	printf("tokens %zu\n", t.n);
	printf("ALL PASS\n");
	return 0;
}