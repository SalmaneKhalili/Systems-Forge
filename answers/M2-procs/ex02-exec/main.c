//
// Created by salmane on 9/5/26.
//


#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

int main(int argc, char **argv)
{
    pid_t pid;
    int st;
    char *cmd;

    if (argc < 0)
        return 0;
    switch (pid = fork())
    {
    case -1:

    case 0:
        cmd = argv[1];
        argv++;
        execvp(cmd, argv); // never returns on success
        printf("exec failed for %s\n", cmd);
        exit(127);
    default:
        if (waitpid(pid, &st, 0) == pid &&
            WEXITSTATUS(st) == 0 && WIFEXITED(st))
        {
            printf("parent: exec ok\n");
            return 0;
        } else
        {
            printf("parent: exec failed and was reported (127)\n");
            return 127;
        }
    }
}
