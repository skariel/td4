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

# TODO: separate `.go` file events for worker and server
# TODO: add worker (run, stop, etc.)
    

class FileEventManager:
    def __init__(self):
        self.event_queue = []

    def dispatch(self, event):
        self.event_queue.append(event)


# TODO: change to general process, so worker can fit here too
class ManagedProcess:
    def __init__(self, display_name, command, q):
        self.display_name = display_name
        self.command = command
        self.server = None
        self.ps = None
        self.pe = None
        self.q = q

    def run(self):
        self.server = Popen(self.command, stdout=PIPE, stderr=PIPE, encoding='utf-8', universal_newlines=True, bufsize=1, shell=True, preexec_fn=os.setsid)
        self.ps = Process(target=self.continuous_read, args=(self.server.stdout, 'std'))
        self.pe = Process(target=self.continuous_read, args=(self.server.stderr, 'err'))
        self.ps.start()
        self.pe.start()

    def stop(self):
        os.killpg(os.getpgid(self.server.pid), signal.SIGTERM)
        self.ps.kill()
        self.pe.kill()

    def continuous_read(self, stream, k):
        while True:
            l = stream.readline()
            if len(l) == 0:
                # process is dead and we read everything there is
                if k=='std':
                    print(R+f'* {self.display_name} is down!!!')
                break
            # process is till alive!
            self.q.put((self.display_name, k, l))


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
    q = Queue()
    server = ManagedProcess(display_name='API', command='go run ./back', q=q)
    server.run()

    last_time = time.time() # now!
    should_run_sqlc     = False
    should_rerun_server = False
    while True:
        try:
            try:
                # read server output
                l = q.get(timeout=0.3)
                if l[1]=='std':
                    print(W+f' {l[0]} '+l[2].strip())
                else:
                    print(W+f' {l[0]} '+l[2].strip(), file=sys.stderr)
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
                print(G+f'* Killing {server.display_name}')
                server.stop()
                print(G+f'* starting {server.display_name}')
                server.run()
                should_rerun_server = False
        except KeyboardInterrupt:
            print()
            print(G+f'* stopping {server.display_name}')
            server.stop()
            print(G+'* exit!')
            sys.exit()

    observer.join()
