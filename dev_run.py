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

# TODO: put `sqlc` in a `SyncService`

class FileEventManager:
    def __init__(self):
        self.event_queue = []

    def dispatch(self, event):
        self.event_queue.append(event)


class AsyncProcess:
    def __init__(self, display_name, command, cwd, q, prefix_path_filter, postfix_path_filter):
        self.display_name = display_name
        self.command = command
        self.cwd = cwd
        self.server = None
        self.ps = None
        self.pe = None
        self.q = q
        self.prefix_path_filter = prefix_path_filter
        self.postfix_path_filter = postfix_path_filter
        self.should_rerun = False
        self.last_time = time.time()

    def path_event_passes(self, p):
        # TODO: better filtering
        return any(p.endswith(pf) for pf in self.postfix_path_filter) and \
            any(p.startswith(pp) for pp in self.prefix_path_filter)

    def run(self):
        self.server = Popen(self.command, stdout=PIPE, stderr=PIPE, encoding='utf-8', universal_newlines=True, bufsize=1, shell=True, preexec_fn=os.setsid, cwd=self.cwd)
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

    print(G+'* initializing services')
    q = Queue()
    services = [
        AsyncProcess(
            display_name='API',
            command='go run .',
            cwd='./back',
            q=q,
            prefix_path_filter=[path+'back', path+'sql/db'],
            postfix_path_filter=['.go']
        ),
        AsyncProcess(
            display_name='WRK',
            command='go run .',
            cwd='./test_worker',
            q=q,
            prefix_path_filter=[path+'test_worker', path+'sql/db'],
            postfix_path_filter=['.go']
        )
    ]
    for s in services:
        print(G+'* starting '+s.display_name)
        s.run()

    sqlc_last_time = time.time() # now!
    should_run_sqlc     = False
    while True:
        try:
            try:
                # read services output
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
                    sqlc_last_time = time.time()
                    continue
                for s in services:
                    if not s.path_event_passes(e.src_path):
                        continue
                    s.should_rerun = True
                    s.last_time = time.time()


            if should_run_sqlc and time.time() - sqlc_last_time > 0.35:
                print(G+'* running sqlc')
                o = getoutput('sqlc generate').strip()
                if len(o) > 0:
                    print(W+o)
                should_run_sqlc = False
                print(G+'* done with sqlc')

            for s in services:
                if s.should_rerun and time.time() - s.last_time > 0.35:
                    print(G+f'* Killing {s.display_name}')
                    s.stop()
                    print(G+f'* starting {s.display_name}')
                    s.run()
                    s.should_rerun = False

        except KeyboardInterrupt:
            print()
            for s in  services:
                print(G+f'* stopping {s.display_name}')
                s.stop()
            print(G+'* exit!')
            sys.exit()

    observer.join()
