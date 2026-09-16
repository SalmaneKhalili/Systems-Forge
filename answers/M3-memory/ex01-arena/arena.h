#ifndef ARENA_H
#define ARENA_H

#include <stddef.h>

void *arena_alloc(size_t n);
size_t arena_left(void);
void arena_reset(void);

#endif