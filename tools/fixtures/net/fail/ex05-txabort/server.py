import os
import socket
import threading


def main():
    port = int(os.environ.get("TARGETPORT", "17010"))
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind(("127.0.0.1", port))
    srv.listen(4)
    while True:
        conn, _ = srv.accept()
        threading.Thread(target=serve, args=(conn,), daemon=True).start()


def serve(conn):
    with conn:
        buf = b""
        while True:
            data = conn.recv(4096)
            if not data:
                return
            buf += data
            while b"\n" in buf:
                line, buf = buf.split(b"\n", 1)
                parts = line.decode().split()
                if len(parts) == 2 and parts[0] == "tx":
                    conn.sendall(b"COMMIT\n")


if __name__ == "__main__":
    main()