#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

int main(void)
{
	char buf[64];
	ssize_t n;
	int fd, status;
	pid_t pid;

	fd = open("out.txt", O_CREAT | O_TRUNC | O_WRONLY, 0644);
	if (fd < 0) {
		fprintf(stderr, "open failed\n");
		return 1;
	}
	pid = fork();
	if (pid < 0) {
		fprintf(stderr, "fork failed\n");
		return 1;
	}
	if (pid == 0) {
		if (dup2(fd, STDOUT_FILENO) < 0)
			_exit(1);
		close(fd);
		printf("redirected\n");
		fflush(stdout); /* _exit skips the stdio flush */
		_exit(0);
	}
	close(fd);
	waitpid(pid, &status, 0);
	if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
		fprintf(stderr, "child died\n");
		return 1;
	}
	fd = open("out.txt", O_RDONLY);
	if (fd < 0) {
		fprintf(stderr, "reopen failed\n");
		return 1;
	}
	n = read(fd, buf, sizeof(buf) - 1);
	close(fd);
	if (n < 0) {
		fprintf(stderr, "read back failed\n");
		return 1;
	}
	buf[n] = '\0';
	if (strcmp(buf, "redirected\n") != 0) {
		printf("captured: %s\n", buf);
		return 1;
	}
	printf("captured: redirected\n");
	printf("ALL PASS\n");
	return 0;
}