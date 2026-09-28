#include "arena.h"
#include <unistd.h>
#include <stddef.h>


// bump allocator
// -----------
// sbrk(512); -> address dyal 512
//

#define ARENA_SIZE 512

static unsigned char *start;
static unsigned char *current;
static size_t left = ARENA_SIZE;

void *arena_alloc(size_t n)
{
  if (start == NULL)
  {
    start = sbrk(ARENA_SIZE);
    if (start == (void *)-1)
      return NULL;

    current = start;
  }

  size_t remainder = n % 8;

  if (remainder != 0)
    n += 8 - remainder;

  if (n > left)
    return NULL;

  void *ptr = current;
  current += n;
  left -= n;

  return ptr;
}

size_t arena_left(void){
  return left;
}
void arena_reset(void){
  left = ARENA_SIZE;
  current = start;
}