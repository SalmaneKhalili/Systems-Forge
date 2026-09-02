#include <stdio.h>
#include <stdlib.h>

int main(void)
{
	int *p = malloc(sizeof *p);
	*p = 1;
	free(p);
	*p = 2; /* use-after-free */
	printf("micro ok\n");
	return 0;
}
