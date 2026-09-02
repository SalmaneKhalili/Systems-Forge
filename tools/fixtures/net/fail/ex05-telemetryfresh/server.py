import socket
import threading

# BUG: each connection starts from a fresh counter registry and never
# increments conn_total per connection, so counters reset and conn_total is 0.
lock = threading.Lock()
PORT = 18541


def on(c):
    counters = {}  # fresh per-connection registry: the bug
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == 'INC':
            if len(p) >= 3:
                counters[p[1]] = counters.get(p[1], 0) + int(p[2])
            f.write(b"ok\n")
        elif cmd == 'GET':
            v = counters.get(p[1], 0) if len(p) >= 2 else 0
            f.write((str(v) + "\n").encode())
        elif cmd == 'SNAP':
            lines = "\n".join("%s=%d" % (k, counters[k]) for k in sorted(counters))
            f.write((lines + "\n").encode())
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