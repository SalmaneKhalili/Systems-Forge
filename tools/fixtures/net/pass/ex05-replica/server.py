import socket, threading
entries=[]; commit=-1; lock=threading.Lock()
def on(c):
    global commit
    f=c.makefile('rwb')
    for line in f:
        parts=line.decode().strip().split()
        if not parts: continue
        if parts[0]=='append':
            with lock: entries.append(parts[1]); idx=len(entries)-1
            f.write(f"ok {idx}\n".encode()); f.flush()
        elif parts[0]=='commit':
            with lock: commit=commit if commit>int(parts[1]) else int(parts[1]); idx=commit
            f.write(f"commit {idx}\n".encode()); f.flush()
        elif parts[0]=='read':
            n=commit+1
            with lock: out="|".join(entries[:n])
            f.write((out+"\n").encode()); f.flush()
    c.close()
s=socket.socket(); s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1); s.bind(('',17010)); s.listen(5)
while True:
    c,_=s.accept(); threading.Thread(target=on,args=(c,)).start()
