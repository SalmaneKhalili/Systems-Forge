#ifndef MM_H
#define MM_H

#include <stddef.h>

void *mm_calloc(size_t n);
void mm_delloc(void *p);
size_t mm_allocated(void);

#endif