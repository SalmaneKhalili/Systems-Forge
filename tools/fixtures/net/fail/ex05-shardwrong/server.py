import socket, threading

# BUG: everything routes to shard 0.
shards = [{}, {}, {}]
lock = threading.Lock()


def on(c):
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        if p[0] == 'set':
            with lock:
                shards[0][p[1]] = p[2]
            f.write(b"set 0\n"); f.flush()
        elif p[0] == 'get':
            with lock:
                val = shards[0].get(p[1], '?')
            f.write(f"0:{val}\n".encode()); f.flush()
    c.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', 17015))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=on, args=(c,)).start()
