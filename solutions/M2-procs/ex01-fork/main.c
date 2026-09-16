#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main(void)
{
	pid_t pid;
	int st = -1;

	pid = fork();
	if (pid < 0) {
		printf("fork failed\n");
		return 1;
	}
	if (pid == 0) {
		printf("child: hello from the child\n");
		return 7;
	}
	if (waitpid(pid, &st, 0) == pid &&
	    WIFEXITED(st) && WEXITSTATUS(st) == 7)
		printf("parent: reaped the child, exit status 7\n");
	else
		printf("parent: reap failed\n");
	return 0;
}
