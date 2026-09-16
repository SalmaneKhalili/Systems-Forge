//
// Created by salmane on 9/5/26.
//
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main(void)
{
    pid_t pid;
    int st;

    switch (pid = fork())
    {
    case -1:
        printf("fork failed\n");
        return 1;
    case 0:
        printf("child: hello from the child\n");
        exit(7);
    default:
        if (waitpid(pid, &st, 0) == pid &&
     WIFEXITED(st) && WEXITSTATUS(st) == 7)
            printf("parent: reaped the child, exit status 7\n");
    }
    return 0;
}
