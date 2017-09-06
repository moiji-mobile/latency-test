#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include <sys/socket.h>
#include <sys/types.h>
#include <sys/uio.h>

#include <netinet/in.h>


int main(int argc, char **argv)
{
	if (argc < 2) {
		fprintf(stderr, "%s port\n", argv[0]);
		return EXIT_FAILURE;
	}

	int server = socket(PF_INET, SOCK_STREAM, 0);
	struct sockaddr_in addr = { 0, };
	addr.sin_family = PF_INET;
	addr.sin_port = htons(atoi(argv[1]));
	addr.sin_addr.s_addr = INADDR_ANY;

	int rc = bind(server, (struct sockaddr *) &addr, sizeof(addr));
	if (rc != 0) {
		fprintf(stderr, "Failed to bind port(%d) %s\n", atoi(argv[1]), strerror(errno));
		return EXIT_FAILURE;
	}

	rc = listen(server, 10);
	if (rc != 0) {
		fprintf(stderr, "Failed to listen %s\n", strerror(errno));
		return EXIT_FAILURE;
	}

	do {
		int fd = accept(server, NULL, NULL);
		do {
			char buf[4096];
			rc = read(fd, &buf[0], 2);
			if (rc != 2) {
				fprintf(stderr, "Short read...rc(%d) errno(%d)\n", rc, errno);
				break;
			}
			const unsigned short len = ntohs(*((unsigned short *)&buf[0]));
			if (len > sizeof(buf) - 2) {
				fprintf(stderr, "Len too long: len(%d) errno(%d)\n", len, errno);
				break;
			}
			rc = read(fd, &buf[2], len);
			if (rc != len) {
				fprintf(stderr, "Short read...rc(%d) len(%d) errno(%d)\n", rc, len, errno);
				break;
			}
			rc = write(fd, buf, len + 2);
			if (rc != len + 2) {
				fprintf(stderr, "Short write...rc(%d) errno(%d)\n", rc, errno);
				break;
			}
		} while (1);
		printf("Closed connection.\n");
		close(fd);
	} while (1);
	

	return EXIT_SUCCESS;
}
