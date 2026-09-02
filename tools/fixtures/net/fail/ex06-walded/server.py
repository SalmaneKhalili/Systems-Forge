import socket
import threading

# BUG: writes are applied to memory but NEVER appended to the write-ahead log.
# State "survives" in-session, yet the durable record is empty -- a crash would
# lose everything. The `wal` command must show committed ops; here it lies.
PORT = 17670


def serve(conn):
    f = conn.makefile("rwb")
    data = {}
    txn = None
    for line in f:
        t = line.decode().strip().split()
        if not t:
            continue
        cmd, p = t[0], t[1:]
        if cmd == "set":
            if txn is not None:
                txn[p[0]] = p[1]
            else:
                data[p[0]] = p[1]
                # BUG: nothing appended to any WAL.
            f.write(b"ok\n")
        elif cmd == "get":
            f.write((data.get(p[0], "nil") + "\n").encode())
        elif cmd == "begin":
            txn = {}
            f.write(b"ok\n")
        elif cmd in ("commit", "rollback"):
            if cmd == "commit" and txn is not None:
                data.update(txn)
            if cmd == "commit":
                pass  # BUG: committed ops never reach the durable record.
            txn = None
            f.write(b"ok\n")
        elif cmd == "wal":
            # BUG: no WAL exists at all.
            f.write(b"\n")
        elif cmd == "list":
            f.write(("\n".join(sorted(data)) + "\n").encode())
        else:
            f.write(b"err\n")
        f.flush()
    conn.close()


s = socket.socket()
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(("", PORT))
s.listen(5)
while True:
    c, _ = s.accept()
    threading.Thread(target=serve, args=(c,)).start()