import os
import socket
import struct
import subprocess
import sys


def pick_free_port():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("127.0.0.1", 0))
    port = s.getsockname()[1]
    s.close()
    return port


def wait_up(port, tries=50):
    for _ in range(tries):
        try:
            c = socket.create_connection(("127.0.0.1", port), timeout=0.2)
            c.close()
            return
        except OSError:
            pass


def forward(src, dst):
    try:
        while True:
            data = src.recv(4096)
            if not data:
                return
            dst.sendall(data)
    finally:
        try:
            dst.shutdown(socket.SHUT_WR)
        except OSError:
            pass


def main():
    port = int(os.environ.get("TARGETPORT", "17005"))
    worker = pick_free_port()
    subprocess.Popen(
        [sys.executable, "worker.py"],
        env={**os.environ, "TARGETPORT": str(worker)},
    )
    wait_up(worker)

    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind(("127.0.0.1", port))
    srv.listen(4)
    while True:
        conn, _ = srv.accept()
        backend = socket.create_connection(("127.0.0.1", worker))
        import threading

        threading.Thread(target=forward, args=(conn, backend), daemon=True).start()
        forward(backend, conn)


if __name__ == "__main__":
    main()