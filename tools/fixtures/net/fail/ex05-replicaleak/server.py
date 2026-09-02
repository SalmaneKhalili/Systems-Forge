import socket, threading
entries=[]; commit=-1; lock=threading.Lock()
def on(c):
    global commit
    f=c.makefile('rwb')
    for line in f:
        parts=line.decode().strip().split()
        if not parts: continue
        if parts[0]=='append':
            entries.append(parts[1]); f.write(f"ok {len(entries)-1}\n".encode()); f.flush()
        elif parts[0]=='commit':
            commit=int(parts[1]); f.write(f"commit {commit}\n".encode()); f.flush()
        elif parts[0]=='read':
            f.write(("|".join(entries)+"\n").encode()); f.flush()   # BUG: leaks uncommitted
    c.close()
s=socket.socket(); s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1); s.bind(('',17010)); s.listen(5)
while True:
    c,_=s.accept(); threading.Thread(target=on,args=(c,)).start()
