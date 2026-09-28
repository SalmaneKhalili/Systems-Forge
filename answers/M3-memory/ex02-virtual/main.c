#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

int	main(void)
{
	pid_t			pid;
	int				st;
	unsigned char	*f_arr;
	unsigned char	*s_arr;

	f_arr = (unsigned char *)mmap(NULL, 2, PROT_READ | PROT_WRITE,
			MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
	s_arr = (unsigned char *)mmap(NULL, 2, PROT_READ | PROT_WRITE,
			MAP_SHARED | MAP_ANONYMOUS, -1, 0);
	f_arr[1] = '\0';
	s_arr[1] = '\0';
	f_arr[0] = 'A';
	s_arr[0] = 'A';
	switch (pid = fork())
	{
	case -1:
		perror("Forking failed");
		return (1);
	case 0:
		f_arr[0] = 'B';
		s_arr[0] = 'B';
		_exit(0);
	default:
		if (waitpid(pid, &st, 0) == pid && WIFEXITED(st)
			&& WEXITSTATUS(st) == 0)
		{
			if (f_arr[0] == 'A')
				printf("private page kept A after child write (copy-on-write)\n");
			if (s_arr[0] == 'B')
				printf("shared page shows B (MAP_SHARED)");
		}
	}
	munmap(f_arr, 2);
	munmap(s_arr, 2);
	return (0);
}
