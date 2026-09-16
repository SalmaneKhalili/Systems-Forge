#include <stdatomic.h>
#include <pthread.h>
#include <stdio.h>

#define N 200000

static atomic_ulong tally;

static void *adder(void *arg)
{
	long i;
	(void)arg;
	for (i = 0; i < N; i++)
		atomic_fetch_add(&tally, 1UL);
	return NULL;
}

int main(void)
{
	pthread_t a, b;

	if (pthread_create(&a, NULL, adder, NULL) != 0 ||
	    pthread_create(&b, NULL, adder, NULL) != 0) {
		fprintf(stderr, "create failed\n");
		return 1;
	}
	pthread_join(a, NULL);
	pthread_join(b, NULL);
	if (atomic_load(&tally) != (unsigned long)2 * N) {
		printf("total %lu (want %d)\n", atomic_load(&tally), 2 * N);
		return 1;
	}
	printf("total 400000\n");
	printf("atomic ok\n");
	printf("ALL PASS\n");
	return 0;
}