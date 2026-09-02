import os
import socket


class OnceServer:
    """Serves exactly one connection, then closes the listening socket."""

    def __init__(self):
        self.served = 0

    def serve_loop(self, srv):
        while self.served < 1:
            conn, _ = srv.accept()
            try:
                while True:
                    data = conn.recv(4096)
                    if not data:
                        break
                    conn.sendall(b"PONG\n")
            finally:
                conn.close()
            self.served += 1


def main():
    port = int(os.environ.get("TARGETPORT", "17004"))
    host = os.environ.get("TARGETHOST", "127.0.0.1")
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind((host, port))
    srv.listen(4)
    OnceServer().serve_loop(srv)


if __name__ == "__main__":
    main()