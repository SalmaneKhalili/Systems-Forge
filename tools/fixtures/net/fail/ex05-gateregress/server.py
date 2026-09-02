import socket, threading

term = 0
entries = []
lock = threading.Lock()


def on(c):
    global term
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        if p[0] == 'see':
            t = int(p[1])
            with lock:
                term = t          # BUG: allows a stale term to regress
                cur = term
            f.write(f"term {cur}\n".encode()); f.flush()
        elif p[0] == 'append':
            with lock:
                if term < 1:
                    f.write(b"not leader\n"); f.flush(); continue
                entries.append(p[1]); idx = len(entries) - 1
            f.write(f"ok {idx}\n".encode()); f.flush()
        elif p[0] == 'read':
            with lock:
                out = "|".join(entries)
            f.write((out + "\n").encode()); f.flush()
    c.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', 17013))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=on, args=(c,)).start()
