import socket
import threading

# BUG: staged writes are applied to the shared store immediately (at 'set'),
# so an uncommitted write from one connection leaks into the live set before
# 'commit'. A correct gateway keeps staged writes invisible until commit.
data = {}
mu = threading.Lock()
PORT = 17641


def on(c):
    f = c.makefile('rwb')
    staged = {}
    in_txn = False
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == 'begin':
            staged = {}
            in_txn = True
            f.write(b"ok\n")
        elif cmd == 'commit':
            in_txn = False
            f.write(b"ok\n")
        elif cmd == 'rollback':
            staged = {}
            in_txn = False
            f.write(b"ok\n")
        elif cmd == 'set' and len(p) >= 3:
            staged[p[1]] = p[2]
            if in_txn:
                # BUG: eagerly publish staged write to the shared live set.
                with mu:
                    data[p[1]] = p[2]
            else:
                with mu:
                    data[p[1]] = p[2]
            f.write(b"ok\n")
        elif cmd == 'get' and len(p) >= 2:
            if in_txn and p[1] in staged:
                f.write((staged[p[1]] + "\n").encode())
            else:
                with mu:
                    v = data.get(p[1])
                if v is None:
                    f.write(b"nil\n")
                else:
                    f.write((v + "\n").encode())
        elif cmd == 'drop' and len(p) >= 2:
            with mu:
                data.pop(p[1], None)
            f.write(b"ok\n")
        elif cmd == 'list':
            with mu:
                keys = sorted(data)
            f.write(("\n".join(keys) + "\n").encode())
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