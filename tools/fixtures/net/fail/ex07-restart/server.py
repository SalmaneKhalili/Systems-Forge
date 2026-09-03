import os
import signal
import socket
import sys


def main():
    port = int(os.environ.get("TARGETPORT", "17021"))
    host = os.environ.get("TARGETHOST", "127.0.0.1")

    # The negative control: a server that swallows SIGTERM and refuses to die,
    # so a supervisor cannot restart it. The grader's restart step sends TERM,
    # waits for exit, times out, and flags "server did not exit before restart".
    def ignore(signum, frame):
        pass

    signal.signal(signal.SIGTERM, ignore)

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
