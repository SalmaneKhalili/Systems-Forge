//
// Created by salmane on 9/5/26.
//

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>

int main()
{
    int arr[2];
    char buff[64];
    pid_t pid;
    ssize_t n;
    ssize_t total = 0;
    int flag = 0;


    if (pipe(arr) == -1)
    {
       printf("pipe failed\n");
        return 1;
    }

    switch (pid = fork())
    {
    case -1:
        printf("fork failed\n");
        return 1;
    case 0 :
        // child process
        close(arr[0]);
        write(arr[1], "hello over the pipe", 19); //fd of stdout
        close(arr[1]);
        return 0;
    default:
        close(arr[1]);
        while ((n = read(arr[0], buff, sizeof(buff))) > 0)
        {
            if (n < 0)
            {
                printf("read failed\n");
                return 1;
            }
            total += n;
        }
        close(arr[0]);
        buff[total] = '\0';
        waitpid(pid, NULL, 0);
        if (total == 19 && memcmp(buff, "hello over the pipe", 19) == 0)
            flag = 1;
        if (flag) {
            printf("read %zd bytes: %s\n", total, buff);
            printf("pipe ok\n");
        } else
        {
            printf("mismatch (read %zd bytes)\n", total);
        }
    }
    return !flag;
}
