#include <arpa/inet.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

static int conn_no;

static void serve(int client)
{
	char buf[4096];
	char hdr[32];
	size_t used = 0;
	ssize_t n;
	char *nl;
	int hl;

	conn_no++;
	hl = snprintf(hdr, sizeof(hdr), "HELLO %d\n", conn_no);
	if (write(client, hdr, (size_t)hl) != hl)
		goto done;

	for (;;) {
		n = read(client, buf + used, sizeof(buf) - used);
		if (n <= 0)
			break;
		used += (size_t)n;
		while ((nl = memchr(buf, '\n', used)) != NULL) {
			size_t linelen = (size_t)(nl - buf);
			if (linelen == 4 && memcmp(buf, "PING", 4) == 0) {
				(void)write(client, "PONG\n", 5);
			} else if (linelen > 5 && memcmp(buf, "ECHO ", 5) == 0) {
				(void)write(client, "ECHO ", 5);
				(void)write(client, buf + 5, linelen - 5);
				(void)write(client, "\n", 1);
			} else {
				(void)write(client, "BAD\n", 4);
			}
			used -= linelen + 1;
			memmove(buf, nl + 1, used);
		}
	}
done:
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