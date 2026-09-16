#include <arpa/inet.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#define REQ "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\npong\n"
#define NOTFOUND "HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"

static int find_term(const char *buf, size_t used)
{
	size_t i;
	for (i = 0; i + 3 < used; i++) {
		if (buf[i] == '\r' && buf[i + 1] == '\n' &&
		    buf[i + 2] == '\r' && buf[i + 3] == '\n')
			return (int)i;
	}
	return -1;
}

static void serve(int client)
{
	char buf[4096];
	char path[129];
	size_t used = 0;
	ssize_t n;
	int end;

	for (;;) {
		n = read(client, buf + used, sizeof(buf) - used);
		if (n <= 0)
			break;
		used += (size_t)n;
		for (;;) {
			end = find_term(buf, used);
			if (end < 0) {
				if (used >= sizeof(buf)) {
					(void)write(client, NOTFOUND, sizeof(NOTFOUND) - 1);
					used = 0;
					break;
				}
				break;
			}
			buf[end] = '\0';
			if (sscanf(buf, "GET %128s HTTP/1.1", path) == 1 &&
			    strcmp(path, "/ping") == 0)
				(void)write(client, REQ, sizeof(REQ) - 1);
			else
				(void)write(client, NOTFOUND, sizeof(NOTFOUND) - 1);
			used -= (size_t)end + 4;
			memmove(buf, buf + end + 4, used);
		}
	}
	close(client);
}

int main(void)
{
	struct sockaddr_in sa;
	socklen_t alen = sizeof(sa);
	const char *portEnv;
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
		serve(client);
	}
	return 0;
}