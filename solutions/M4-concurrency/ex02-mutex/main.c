#include <pthread.h>
#include <stdio.h>

#define PER_THREAD 200000

static volatile long tally;
static pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;

static void *adder(void *arg)
{
	long i;
	(void)arg;
	for (i = 0; i < PER_THREAD; i++) {
		pthread_mutex_lock(&lock);
		tally++;
		pthread_mutex_unlock(&lock);
	}
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
	if (tally != (long)2 * PER_THREAD) {
		printf("total %ld (want %d)\n", tally, 2 * PER_THREAD);
		return 1;
	}
	printf("total 400000\n");
	printf("locks ok\n");
	printf("ALL PASS\n");
	return 0;
}