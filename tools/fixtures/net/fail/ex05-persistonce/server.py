import socket
import threading

# BUG: each connection uses a fresh local map, so writes from an earlier
# connection are never visible (not persistent/shared).
lock = threading.Lock()
PORT = 18241


def on(c):
    # Fresh per-connection store: the bug.
    data = {}
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == 'PUT':
            if len(p) >= 3:
                data[p[1]] = p[2]
            f.write(b"ok\n")
        elif cmd == 'GET':
            v = data.get(p[1]) if len(p) >= 2 else None
            if v is None:
                f.write(b"nil\n")
            else:
                f.write((v + "\n").encode())
        elif cmd == 'DEL':
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