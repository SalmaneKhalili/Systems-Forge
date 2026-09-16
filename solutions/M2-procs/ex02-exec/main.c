#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main(int argc, char **argv)
{
	pid_t pid;
	int st = -1;

	if (argc < 2)
		return 2;
	pid = fork();
	if (pid < 0) {
		printf("fork failed\n");
		return 1;
	}
	if (pid == 0) {
		if (execvp(argv[1], &argv[1]) == -1) {
			printf("exec failed for %s\n", argv[1]);
			return 127;
		}
	}
	if (waitpid(pid, &st, 0) == pid && WIFEXITED(st)) {
		int code = WEXITSTATUS(st);
		if (code == 127)
			printf("parent: exec failed and was reported (127)\n");
		else if (code == 0)
			printf("parent: exec ok\n");
		else
			printf("parent: child exited with status %d\n", code);
		return code;
	}
	printf("parent: reap failed\n");
	return 1;
}
