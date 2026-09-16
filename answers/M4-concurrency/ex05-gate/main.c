#include <pthread.h>
#include <stdio.h>

#include "queue.h"

#define ITEMS 1000

static Queue *q;
static int pop_err;
static int fifo_ok = 1;
static long sum;

static void *producer(void *arg)
{
	int i;
	(void)arg;
	for (i = 0; i < ITEMS; i++)
		if (queue_push(q, i) != 0)
			return (void *)1;
	return NULL;
}

static void *consumer(void *arg)
{
	int v, i;
	(void)arg;
	for (i = 0; i < ITEMS; i++) {
		if (queue_pop(q, &v) != 0) {
			pop_err++;
			return (void *)1;
		}
		if (v != i)
			fifo_ok = 0;
		sum += (long)v;
	}
	return NULL;
}

int main(void)
{
	pthread_t p, c;
	void *r1, *r2;

	q = queue_new(16);
	if (q == NULL) {
		fprintf(stderr, "queue_new failed\n");
		return 1;
	}
	if (pthread_create(&p, NULL, producer, NULL) != 0 ||
	    pthread_create(&c, NULL, consumer, NULL) != 0) {
		fprintf(stderr, "create failed\n");
		return 1;
	}
	pthread_join(p, &r1);
	pthread_join(c, &r2);
	queue_free(q);

	if (r1 != NULL || r2 != NULL) {
		printf("worker error\n");
		return 1;
	}
	if (pop_err != 0) {
		printf("pop errors %d\n", pop_err);
		return 1;
	}
	printf("produced 1000\n");
	printf("consumed 1000\n");
	if (!fifo_ok) {
		printf("fifo violated\n");
		return 1;
	}
	printf("fifo ok\n");
	if (sum != (long)ITEMS * (ITEMS - 1) / 2) {
		printf("checksum %ld (want 499500)\n", sum);
		return 1;
	}
	printf("checksum 499500\n");
	printf("ALL PASS\n");
	return 0;
}