#ifndef POOL_H
#define POOL_H

#include <stddef.h>

void *pool_get(void);
void pool_put(void *slot);
size_t pool_borrowed(void);

#endif