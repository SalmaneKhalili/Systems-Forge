import socket, threading

shards = [{}, {}, {}]
lock = threading.Lock()


def shard_of(key):
    if key <= 'm':
        return 0
    if key <= 't':
        return 1
    return 2


def on(c):
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        if p[0] == 'set':
            n = shard_of(p[1])
            with lock:
                shards[n][p[1]] = p[2]
            f.write(f"set {n}\n".encode()); f.flush()
        elif p[0] == 'get':
            n = shard_of(p[1])
            with lock:
                val = shards[n].get(p[1], '?')
            f.write(f"{n}:{val}\n".encode()); f.flush()
    c.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', 17014))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=on, args=(c,)).start()
