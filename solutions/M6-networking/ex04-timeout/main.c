#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <unistd.h>

#define RTIMEOUT_MS 400

static void serve(int client)
{
	char buf[4096];
	size_t used = 0;
	ssize_t n;
	char *nl;

	for (;;) {
		n = read(client, buf + used, sizeof(buf) - used);
		if (n > 0) {
			used += (size_t)n;
			while ((nl = memchr(buf, '\n', used)) != NULL) {
				size_t linelen = (size_t)(nl - buf);
				(void)write(client, buf, linelen);
				(void)write(client, "\n", 1);
				used -= linelen + 1;
				memmove(buf, nl + 1, used);
			}
			continue;
		}
		if (n < 0 && (errno == EAGAIN || errno == EWOULDBLOCK)) {
			(void)write(client, "TIMEOUT\n", 8);
			break;
		}
		break; /* EOF or hard error */
	}
	close(client);
}

int main(void)
{
	struct sockaddr_in sa;
	struct timeval tv;
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
	tv.tv_sec = 0;
	tv.tv_usec = RTIMEOUT_MS * 1000;
	setsockopt(server_fd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));
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