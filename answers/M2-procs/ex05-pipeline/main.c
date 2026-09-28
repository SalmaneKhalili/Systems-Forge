//
// Created by Salmane on 9/15/26.
//

#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/wait.h>

int main(const int argc, char** argv)
{
    if (argc == 2)
    {
        int fds[2];
        pid_t pid_child_a;
        pid_t pid_child_b;
        int st_a;
        int st_b;


        if (pipe(fds) == -1)
        {
            printf("Pipe failed\n");
            return 1;
        }

        switch (pid_child_a = fork())
        {
        case -1:
            printf("fork failed\n");
            return 1;
        case 0:
            close(fds[0]);
            dup2(fds[1], STDOUT_FILENO);
            close(fds[1]);

            if (execlp("echo", "echo", argv[1], NULL) == 0)
            {
                printf("echo failed\n");
                exit(127);
            }
            break;
        default:
            switch (pid_child_b = fork())
            {
            case -1:
                printf("fork failed\n");
                return 1;
            case 0:
                close(fds[1]);
                dup2(fds[0], STDIN_FILENO);
                close(fds[0]);
                if (execlp("tr", "tr", "a-z", "A-Z", NULL) == 0)
                {
                    printf("tr failed\n");
                    exit(127);
                }
                break;
            default:
                close(fds[0]);
                close(fds[1]);
                if ((waitpid(pid_child_a, &st_a, 0) == pid_child_a && WIFEXITED(st_a) && WEXITSTATUS(st_a) == 0) &&
         (waitpid(pid_child_b, &st_b, 0) == pid_child_b && WIFEXITED(st_b) && WEXITSTATUS(st_b) == 0))
                {
                    printf("pipeline done (2 children)\n");
                }
            }
            return 0;
        }
    }
}
