#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main(void)
{
	int fds[2];
	pid_t pid;
	char buf[64];
	ssize_t n = 0;
	int got = 0;

	if (pipe(fds) == -1) {
		printf("pipe failed\n");
		return 1;
	}
	pid = fork();
	if (pid < 0) {
		printf("fork failed\n");
		return 1;
	}
	if (pid == 0) {
		close(fds[0]);
		write(fds[1], "hello over the pipe", 19);
		close(fds[1]);
		return 0;
	}
	close(fds[1]);
	while (1) {
		ssize_t r = read(fds[0], buf + n, (int)(sizeof(buf) - 1 - (size_t)n));
		if (r < 0) {
			printf("read failed\n");
			return 1;
		}
		if (r == 0)
			break;
		n += r;
	}
	close(fds[0]);
	buf[n] = '\0';
	waitpid(pid, NULL, 0);
	if (n == 19 && memcmp(buf, "hello over the pipe", 19) == 0)
		got = 1;
	if (got) {
		printf("read %zd bytes: %s\n", n, buf);
		printf("pipe ok\n");
	} else
		printf("mismatch (read %zd bytes)\n", n);
	//return got ? 0 : 1;
	return !got;
}
