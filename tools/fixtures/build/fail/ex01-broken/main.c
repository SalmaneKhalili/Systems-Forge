#include <stdio.h>

int main(void)
{
	printf("build ok\n") /* no semicolon: must not compile under -Werror */
	return 0;
}