import socket, threading

# BUG: lists every member including dropped/dead ones; beat never revives suspect.
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
                names = sorted(states)      # BUG: includes dead
            f.write(("|".join(names) + "\n").encode()); f.flush()
        elif cmd == 'status':
            with lock:
                s = states.get(p[1], 'dead')
            f.write((s + "\n").encode()); f.flush()
        elif cmd == 'beat':
            with lock:
                states[p[1]] = 'dead' if states.get(p[1]) == 'dead' else states.get(p[1], 'alive')
            f.write(b"ok\n"); f.flush()
        elif cmd == 'suspect':
            with lock:
                states[p[1]] = 'suspect'
            f.write(b"ok\n"); f.flush()
        elif cmd == 'drop':
            with lock:
                states[p[1]] = 'dead'
            f.write(b"ok\n"); f.flush()
    c.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('', 17017))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=on, args=(c,)).start()
