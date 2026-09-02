import os
import socket


def main():
    port = int(os.environ.get("TARGETPORT", "17005"))
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind(("127.0.0.1", port))
    srv.listen(4)
    while True:
        conn, _ = srv.accept()
        try:
            line = conn.recv(4096)
            if line.strip() == b"PING":
                conn.sendall(b"PONG\n")
        finally:
            conn.close()


if __name__ == "__main__":
    main()