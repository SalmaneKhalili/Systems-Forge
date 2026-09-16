#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main(int argc, char **argv)
{
	int fds[2];
	pid_t p1, p2;
	int s1 = -1, s2 = -1;
	const char *word = "supply chain";

	if (argc > 1)
		word = argv[1];
	if (pipe(fds) == -1) {
		printf("pipe failed\n");
		return 1;
	}
	p1 = fork();
	if (p1 < 0) {
		printf("fork failed\n");
		return 1;
	}
	if (p1 == 0) {
		close(fds[0]);
		if (dup2(fds[1], STDOUT_FILENO) == -1)
			return 126;
		close(fds[1]);
		execlp("echo", "echo", word, (char *)NULL);
		printf("echo failed\n");
		return 127;
	}
	p2 = fork();
	if (p2 < 0) {
		printf("fork failed\n");
		return 1;
	}
	if (p2 == 0) {
		close(fds[1]);
		if (dup2(fds[0], STDIN_FILENO) == -1)
			return 126;
		close(fds[0]);
		execlp("tr", "tr", "a-z", "A-Z", (char *)NULL);
		printf("tr failed\n");
		return 127;
	}
	close(fds[0]);
	close(fds[1]);
	waitpid(p1, &s1, 0);
	waitpid(p2, &s2, 0);
	if (WIFEXITED(s1) && WIFEXITED(s2) &&
	    WEXITSTATUS(s1) == 0 && WEXITSTATUS(s2) == 0)
		printf("pipeline done (2 children)\n");
	else
		printf("pipeline FAILED\n");
	return 0;
}
