import socket
import threading

# Correct persistent KV gateway: keeps shared in-memory state across all
# connections, so the store is persistent for the process lifetime.
data = {}
lock = threading.Lock()
PORT = 18240


def on(c):
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == 'PUT':
            if len(p) >= 3:
                with lock:
                    data[p[1]] = p[2]
            f.write(b"ok\n")
        elif cmd == 'GET':
            with lock:
                v = data.get(p[1]) if len(p) >= 2 else None
            if v is None:
                f.write(b"nil\n")
            else:
                f.write((v + "\n").encode())
        elif cmd == 'DEL':
            with lock:
                if len(p) >= 2:
                    data.pop(p[1], None)
            f.write(b"ok\n")
        else:
            f.write(b"err\n")
        f.flush()
    c.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', PORT))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=on, args=(c,)).start()