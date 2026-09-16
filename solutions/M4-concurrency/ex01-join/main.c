#include <inttypes.h>
#include <pthread.h>
#include <stdint.h>
#include <stdio.h>

#define NTHREADS 5

static void *worker(void *arg)
{
	intptr_t i = (intptr_t)arg;
	intptr_t val = i * 2;

	if (i == 2)
		pthread_exit((void *)val);
	return (void *)val;
}

int main(void)
{
	pthread_t th[NTHREADS];
	intptr_t v;
	intptr_t i;
	int jerr = 0;
	long sum = 0;

	for (i = 1; i <= NTHREADS; i++) {
		if (pthread_create(&th[i - 1], NULL, worker, (void *)i) != 0) {
			fprintf(stderr, "create failed\n");
			return 1;
		}
	}
	for (i = 1; i <= NTHREADS; i++) {
		if (pthread_join(th[i - 1], (void **)&v) != 0) {
			jerr++;
			continue;
		}
		printf("thread %" PRIdPTR " returned %" PRIdPTR "\n", i, v);
		sum += (long)v;
	}
	if (jerr != 0) {
		printf("join errors %d\n", jerr);
		return 1;
	}
	printf("join errors 0\n");
	if (sum != (long)NTHREADS * (NTHREADS + 1)) {
		printf("sum %ld (want %ld)\n", sum, (long)NTHREADS * (NTHREADS + 1));
		return 1;
	}
	printf("sum 30\n");
	printf("ALL PASS\n");
	return 0;
}