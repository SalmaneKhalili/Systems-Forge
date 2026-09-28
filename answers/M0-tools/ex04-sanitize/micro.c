#include <stdio.h>
#include <stdlib.h>

int main(void)
{
    int* p = malloc(sizeof *p);
    *p = 1;
    *p = 2; /* use-after-free */
    free(p);
    printf("micro ok\n");
    return 0;
}
