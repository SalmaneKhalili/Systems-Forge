#include <stdio.h>
#include <sys/mman.h>
#include <sys/wait.h>
#include <unistd.h>

#define PAGE 4096

int main(void)
{
	char *priv, *shar;
	pid_t pid;
	int st;
	char pv, sv;

	priv = mmap(NULL, PAGE, PROT_READ | PROT_WRITE,
		    MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
	shar = mmap(NULL, PAGE, PROT_READ | PROT_WRITE,
		    MAP_SHARED | MAP_ANONYMOUS, -1, 0);
	if (priv == MAP_FAILED || shar == MAP_FAILED) {
		printf("FAIL: mmap\n");
		return 1;
	}
	priv[0] = 'A';
	shar[0] = 'A';

	pid = fork();
	if (pid < 0) {
		printf("FAIL: fork\n");
		return 1;
	}
	if (pid == 0) {
		priv[0] = 'B';
		shar[0] = 'B';
		_exit(0);
	}
	if (waitpid(pid, &st, 0) != pid) {
		printf("FAIL: waitpid\n");
		return 1;
	}

	pv = priv[0];
	sv = shar[0];
	if (pv != 'A') {
		printf("FAIL: private page shows %c (want A)\n", pv);
		return 1;
	}
	if (sv != 'B') {
		printf("FAIL: shared page shows %c (want B)\n", sv);
		return 1;
	}
	printf("private page kept A after child write (copy-on-write)\n");
	printf("shared page shows B (MAP_SHARED)\n");

	munmap(priv, PAGE);
	munmap(shar, PAGE);
	return 0;
}