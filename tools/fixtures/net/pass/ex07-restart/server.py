import os
import signal
import socket
import sys


def main():
    port = int(os.environ.get("TARGETPORT", "17020"))
    host = os.environ.get("TARGETHOST", "127.0.0.1")

    def handle(signum, frame):
        os._exit(0)

    signal.signal(signal.SIGTERM, handle)

    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind((host, port))
    srv.listen(4)
    while True:
        conn, _ = srv.accept()
        try:
            while True:
                data = conn.recv(4096)
                if not data:
                    break
                if data.strip() == b"PING":
                    conn.sendall(b"PONG\n")
        finally:
            conn.close()


if __name__ == "__main__":
    main()
