import socket, threading

states = {}
lock = threading.Lock()


def on(c):
    f = c.makefile('rwb')
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == 'list':
            with lock:
                names = sorted(n for n, s in states.items() if s in ('alive', 'suspect'))
            f.write(("|".join(names) + "\n").encode()); f.flush()
        elif cmd == 'status':
            with lock:
                s = states.get(p[1], 'dead')
            f.write((s + "\n").encode()); f.flush()
        elif cmd == 'beat':
            with lock:
                states[p[1]] = 'alive'
            f.write(b"ok\n"); f.flush()
        elif cmd == 'suspect':
            with lock:
                if states.get(p[1]) == 'alive':
                    states[p[1]] = 'suspect'
            f.write(b"ok\n"); f.flush()
        elif cmd == 'drop':
            with lock:
                states[p[1]] = 'dead'
            f.write(b"ok\n"); f.flush()
    c.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', 17016))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=on, args=(c,)).start()
