#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

#define RCHUNK 7
#define FCHUNK 11

int main(int argc, char **argv)
{
	char buf[FCHUNK];
	unsigned long rtotal = 0, rsum = 0, ftotal = 0, fsum = 0;
	FILE *f;
	size_t got, i;
	ssize_t n;
	int fd;

	if (argc != 2) {
		fprintf(stderr, "usage: %s FILE\n", argv[0]);
		return 1;
	}

	fd = open(argv[1], O_RDONLY);
	if (fd < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	while ((n = read(fd, buf, RCHUNK)) > 0) {
		rtotal += (unsigned long)n;
		for (i = 0; i < (size_t)n; i++)
			rsum += (unsigned char)buf[i];
	}
	if (n < 0) {
		fprintf(stderr, "read failed\n");
		return 1;
	}
	close(fd);

	f = fopen(argv[1], "r");
	if (f == NULL) {
		fprintf(stderr, "fopen failed\n");
		return 1;
	}
	while ((got = fread(buf, 1, FCHUNK, f)) > 0) {
		ftotal += (unsigned long)got;
		for (i = 0; i < got; i++)
			fsum += (unsigned char)buf[i];
	}
	fclose(f);

	printf("%s read %lu checksum %lu\n", argv[1], rtotal, rsum);
	printf("%s fread %lu checksum %lu\n", argv[1], ftotal, fsum);
	if (rtotal != ftotal || rsum != fsum) {
		printf("mismatch\n");
		return 1;
	}
	printf("match 1\n");
	printf("ALL PASS\n");
	return 0;
}