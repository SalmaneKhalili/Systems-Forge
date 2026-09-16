#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

int main(void)
{
	int a, b, t;

	printf("stdin %d\n", fileno(stdin));
	printf("stdout %d\n", fileno(stdout));
	printf("stderr %d\n", fileno(stderr));

	a = open("/dev/null", O_RDONLY);
	if (a < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	printf("first open %d\n", a);
	close(a);

	a = open("/dev/null", O_RDONLY);
	if (a < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	printf("reuse after close %d\n", a);

	close(a);
	a = open("/dev/null", O_RDONLY);
	b = open("/dev/null", O_RDONLY);
	if (a < 0 || b < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	if (a > b) {
		t = a;
		a = b;
		b = t;
	}
	printf("two opens %d and %d\n", a, b);
	close(a);
	close(b);

	a = open("/dev/null", O_RDONLY);
	if (a < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	printf("lowest available %d\n", a);
	close(a);

	printf("ALL PASS\n");
	return 0;
}