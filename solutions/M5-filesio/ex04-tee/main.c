#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

#define CHUNK 17

int main(void)
{
	char buf[CHUNK];
	unsigned long total = 0, wrote = 0;
	ssize_t n, w;
	int fd;

	fd = open("out.txt", O_CREAT | O_TRUNC | O_WRONLY, 0644);
	if (fd < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	while ((n = read(STDIN_FILENO, buf, CHUNK)) > 0) {
		w = write(STDOUT_FILENO, buf, (size_t)n);
		if (w != n) {
			fprintf(stderr, "stdout write short\n");
			return 1;
		}
		w = write(fd, buf, (size_t)n);
		if (w != n) {
			fprintf(stderr, "file write short\n");
			return 1;
		}
		total += (unsigned long)n;
	}
	if (n < 0) {
		fprintf(stderr, "read failed\n");
		return 1;
	}
	close(fd);

	fd = open("out.txt", O_RDONLY);
	if (fd < 0) {
		fprintf(stderr, "reopen failed\n");
		return 1;
	}
	while ((n = read(fd, buf, CHUNK)) > 0)
		wrote += (unsigned long)n;
	close(fd);

	printf("wrote %lu\n", total);
	printf("readback %lu\n", wrote);
	if (total != wrote) {
		printf("mismatch\n");
		return 1;
	}
	printf("ALL PASS\n");
	return 0;
}