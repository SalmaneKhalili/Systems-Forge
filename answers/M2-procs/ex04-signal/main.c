//
// Created by salmane on 9/15/26.
//

#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

static volatile sig_atomic_t got = 0;
static void on_sigusr1(int signum)
{
    if (signum == SIGUSR1)
        got = 1;
}

int main(void)
{
    // installing handler
    struct sigaction sa;
    sa.sa_handler = on_sigusr1;
    sa.sa_flags = 0;
    sigemptyset(&sa.sa_mask);

    printf("installing\n");
    if (sigaction(SIGUSR1, &sa, NULL) == -1)
    {
        perror("Failure to install handler.");
        exit(EXIT_FAILURE);
    }

    if (raise(SIGUSR1) != 0)
    {
        printf("failed to raise signal");
        exit(EXIT_FAILURE);
    }

    if (got)
        printf("caught SIGUSR1\n");
    return EXIT_SUCCESS;
}
