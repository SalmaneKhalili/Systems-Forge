import socket
import threading

# Correct telemetry gateway: shared counter registry; conn_total increments
# on every accepted connection; SNAP lists metrics sorted by name.
counters = {}
lock = threading.Lock()
PORT = 18540


def on(c):
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == 'INC':
            if len(p) >= 3:
                with lock:
                    counters[p[1]] = counters.get(p[1], 0) + int(p[2])
            f.write(b"ok\n")
        elif cmd == 'GET':
            with lock:
                v = counters.get(p[1], 0) if len(p) >= 2 else 0
            f.write((str(v) + "\n").encode())
        elif cmd == 'SNAP':
            with lock:
                lines = "\n".join("%s=%d" % (k, counters[k]) for k in sorted(counters))
            f.write((lines + "\n").encode())
        else:
            f.write(b"err\n")
        f.flush()
    c.close()


def bump_conn():
    with lock:
        counters['conn_total'] = counters.get('conn_total', 0) + 1


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', PORT))
s.listen(5)
while True:
    c, _ = s.accept()
    bump_conn()
    threading.Thread(target=on, args=(c,)).start()