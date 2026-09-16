#include <stdio.h>
#include <signal.h>

static volatile sig_atomic_t got;

static void on_sigusr1(int sig)
{
	(void)sig;
	got = 1;
}

int main(void)
{
	if (signal(SIGUSR1, on_sigusr1) == SIG_ERR) {
		printf("install failed\n");
		return 1;
	}
	printf("installing\n");
	if (raise(SIGUSR1) != 0) {
		printf("raise failed\n");
		return 1;
	}
	if (got)
		printf("caught SIGUSR1\n");
	else
		printf("handler did not run\n");
	return 0;
}
