#!python -u
import os
import sys
import glob
import time
import signal
import queue
from subprocess import Popen, getoutput, PIPE
from multiprocessing import Process, Queue
from watchdog.observers import Observer

def continuous_read(stream, q, k):
    while True:
        l = stream.readline()
        q.put((k, l))
        if len(l) == 0:
            # process is dead and we read everything there is
            if k=='std':
                print(R+'* server is down!!!')
            break
    

class FileEventManager:
    def __init__(self):
        self.event_queue = []

    def dispatch(self, event):
        self.event_queue.append(event)


class Server:
    def __init__(self):
        self.server = None
        self.ps = None
        self.pe = None
        self.q = Queue()

    def run(self):
        self.server = Popen('go run ./back', stdout=PIPE, stderr=PIPE, encoding='utf-8', universal_newlines=True, bufsize=1, shell=True, preexec_fn=os.setsid)
        self.ps = Process(target=continuous_read, args=(self.server.stdout, self.q, 'std'))
        self.pe = Process(target=continuous_read, args=(self.server.stderr, self.q, 'err'))
        self.ps.start()
        self.pe.start()

    def stop(self):
        os.killpg(os.getpgid(self.server.pid), signal.SIGTERM)
        self.ps.kill()
        self.pe.kill()


# console colors
R  = '\033[1;33m' # bold!
G  = '\033[1;32m' # bold!
W  = '\033[0m'

if __name__=='__main__':

    path = os.environ.get('TD4_ROOT', '.') + '/'
    print(G+'* watching', path)
    observer = Observer()
    file_event_manager = FileEventManager()
    observer.schedule(file_event_manager, path, recursive=True)
    observer.start()

    print(G+'* initializing server')
    server = Server()
    server.run()

    last_time = time.time() # now!
    should_run_sqlc     = False
    should_rerun_server = False
    while True:
        try:
            try:
                # read server output
                l = server.q.get(timeout=0.3)
                if l[0]=='std':
                    print(W+l[1].strip())
                else:
                    print(W+l[1].strip(), file=sys.stderr)
            except queue.Empty:
                # timeout, time to check file events
                pass

            # check file events, at most 0.3 secs apart

            while len(file_event_manager.event_queue) > 0:
                e = file_event_manager.event_queue.pop()
                if e.src_path.endswith('.sql'):
                    should_run_sqlc = True
                    last_time = time.time()
                    continue
                if e.src_path.endswith('.go'):
                    should_rerun_server = True
                    last_time = time.time()
                    continue

            if time.time() - last_time < 0.35:
                # events are too clustered, wait a bit before running sqlc
                continue

            if should_run_sqlc:
                print(G+'* running sqlc')
                o = getoutput('sqlc generate').strip()
                if len(o) > 0:
                    print(W+o)
                should_run_sqlc = False

            if time.time() - last_time < 0.35:
                # events too clustered for killing and rerunning the server...
                continue

            if should_rerun_server:
                print(G+'* Killing server')
                server.stop()
                print(G+'* starting server')
                server.run()
                should_rerun_server = False
        except KeyboardInterrupt:
            print(G+'* stopping server')
            server.stop()
            print(G+'* exit!')
            sys.exit()

    observer.join()
