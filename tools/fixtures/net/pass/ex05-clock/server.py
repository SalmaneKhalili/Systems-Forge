import os
import socket
import threading


def main():
    port = int(os.environ.get("TARGETPORT", "17007"))
    seq = [0]
    lock = threading.Lock()
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind(("127.0.0.1", port))
    srv.listen(4)
    while True:
        conn, _ = srv.accept()
        threading.Thread(
            target=serve, args=(conn, seq, lock), daemon=True
        ).start()


def serve(conn, seq, lock):
    with conn:
        buf = b""
        while True:
            data = conn.recv(4096)
            if not data:
                return
            buf += data
            while b"\n" in buf:
                line, buf = buf.split(b"\n", 1)
                text = line.decode()
                try:
                    content, stamp = text.rsplit(" ", 1)
                    claimed = int(stamp)
                except ValueError:
                    return
                with lock:
                    if claimed > seq[0]:
                        seq[0] = claimed
                    seq[0] += 1
                    n = seq[0]
                conn.sendall(b"R%d %s %s\n" % (n, content.encode(), stamp.encode()))


if __name__ == "__main__":
    main()