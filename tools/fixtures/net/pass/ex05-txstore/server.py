import socket
import threading

# Correct txstore gateway: committed state is shared across connections,
# staged writes are visible only to their own connection until 'commit'.
data = {}
mu = threading.Lock()
PORT = 17640


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
            if in_txn:
                with mu:
                    for k, v in staged.items():
                        data[k] = v
                staged = {}
                in_txn = False
            f.write(b"ok\n")
        elif cmd == 'rollback':
            staged = {}
            in_txn = False
            f.write(b"ok\n")
        elif cmd == 'set' and len(p) >= 3:
            if in_txn:
                staged[p[1]] = p[2]
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