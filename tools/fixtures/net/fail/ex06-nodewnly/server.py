import socket
import threading

# BUG: proposals are never replicated to followers. The 'leader' keeps its own
# log but sends no AppendEntries, so peer logs stay empty and diverge (the
# replication invariant M12-ex06 exists to prove is broken).
PORT = 17650


class Peer(object):
    def __init__(self):
        self.term = 0
        self.lead = False
        self.log = []


peers = {"a": Peer(), "b": Peer(), "c": Peer()}
lead = [None]


def serve_client(conn):
    f = conn.makefile("rwb")
    for line in f:
        p = line.decode().strip().split()
        if not p:
            continue
        cmd = p[0]
        if cmd == "elect" and len(p) >= 2:
            node = p[1]
            peer = peers[node]
            peer.term += 1
            peer.lead = True
            lead[0] = node
            # BUG: no RequestVote round; pretend a majority voted YES.
            f.write(("ELECTED %s term %d\n" % (node, peer.term)).encode())
        elif cmd == "propose" and len(p) >= 2:
            if lead[0] is None:
                f.write(b"not leader\n")
            else:
                leader = peers[lead[0]]
                idx = len(leader.log)
                leader.log.append(" ".join(p[1:]))
                # BUG: never send APP to followers -> their logs stay empty.
                f.write(("ok %d\n" % idx).encode())
        elif cmd == "read":
            leader = peers[lead[0]] if lead[0] else None
            f.write(("|".join(leader.log) + "\n").encode() if leader else b"\n")
        elif cmd in ("log", "term") and len(p) >= 2:
            peer = peers[p[1]]
            if cmd == "log":
                f.write(("|".join(peer.log) + "\n").encode())
            else:
                f.write(("term %d\n" % peer.term).encode())
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
    threading.Thread(target=serve_client, args=(c,)).start()