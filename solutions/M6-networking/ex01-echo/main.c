#include <arpa/inet.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

int main(void)
{
	struct sockaddr_in sa;
	socklen_t alen = sizeof(sa);
	const char *portEnv;
	char buf[4096];
	ssize_t n, w;
	int one = 1;
	int port, server_fd, client;

	signal(SIGPIPE, SIG_IGN);
	portEnv = getenv("TARGETPORT");
	if (portEnv == NULL) {
		fprintf(stderr, "TARGETPORT unset\n");
		return 1;
	}
	port = atoi(portEnv);
	if (port <= 0) {
		fprintf(stderr, "bad TARGETPORT %s\n", portEnv);
		return 1;
	}

	server_fd = socket(AF_INET, SOCK_STREAM, 0);
	if (server_fd < 0) {
		fprintf(stderr, "socket failed\n");
		return 1;
	}
	setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &one, sizeof(one));
	memset(&sa, 0, sizeof(sa));
	sa.sin_family = AF_INET;
	sa.sin_port = htons((uint16_t)port);
	sa.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
	if (bind(server_fd, (struct sockaddr *)&sa, sizeof(sa)) < 0 ||
	    listen(server_fd, 4) < 0) {
		fprintf(stderr, "bind/listen failed\n");
		return 1;
	}

	for (;;) {
		client = accept(server_fd, (struct sockaddr *)&sa, &alen);
		if (client < 0)
			continue;
		while ((n = read(client, buf, sizeof(buf))) > 0) {
			w = write(client, buf, (size_t)n);
			if (w != n)
				break;
		}
		close(client);
	}
	return 0;
}