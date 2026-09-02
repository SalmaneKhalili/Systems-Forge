#ifndef QUEUE_H
#define QUEUE_H

#include <stddef.h>

typedef struct Queue Queue;

Queue *queue_new(size_t cap);
void queue_free(Queue *q);
int queue_push(Queue *q, int v);
int queue_pop(Queue *q, int *out);

#endif