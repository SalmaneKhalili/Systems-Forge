#include <pthread.h>
#include <stdlib.h>

#include "queue.h"

struct Queue {
	pthread_mutex_t lock;
	pthread_cond_t not_full, not_empty;
	size_t cap, count, head, tail;
	int *buf;
};

Queue *queue_new(size_t cap)
{
	Queue *q;

	if (cap == 0)
		return NULL;
	q = calloc(1, sizeof(*q));
	if (q == NULL)
		return NULL;
	q->buf = calloc(cap, sizeof(int));
	if (q->buf == NULL) {
		free(q);
		return NULL;
	}
	q->cap = cap;
	pthread_mutex_init(&q->lock, NULL);
	pthread_cond_init(&q->not_full, NULL);
	pthread_cond_init(&q->not_empty, NULL);
	return q;
}

void queue_free(Queue *q)
{
	if (q == NULL)
		return;
	free(q->buf);
	pthread_mutex_destroy(&q->lock);
	pthread_cond_destroy(&q->not_full);
	pthread_cond_destroy(&q->not_empty);
	free(q);
}

int queue_push(Queue *q, int v)
{
	pthread_mutex_lock(&q->lock);
	while (q->count == q->cap)
		pthread_cond_wait(&q->not_full, &q->lock);
	q->buf[q->head] = v;
	q->head = (q->head + 1) % q->cap;
	q->count++;
	pthread_cond_signal(&q->not_empty);
	pthread_mutex_unlock(&q->lock);
	return 0;
}

int queue_pop(Queue *q, int *out)
{
	pthread_mutex_lock(&q->lock);
	while (q->count == 0)
		pthread_cond_wait(&q->not_empty, &q->lock);
	*out = q->buf[q->tail];
	q->tail = (q->tail + 1) % q->cap;
	q->count--;
	pthread_cond_signal(&q->not_full);
	pthread_mutex_unlock(&q->lock);
	return 0;
}