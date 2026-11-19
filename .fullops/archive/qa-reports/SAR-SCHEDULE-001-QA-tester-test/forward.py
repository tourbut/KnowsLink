"""Forward only this task's synthetic loopback fixture to local port 18961.

The remote origin must be http://127.0.0.1. No operational service is forwarded.
"""
import json
import select
import socketserver
import threading
import urllib.parse

from remote import HERE, SHA, connect

LOCAL_PORT = 18961


def main():
    location = json.loads((HERE / 'fixture-location.json').read_text(encoding='utf-8'))
    base, sha = location['base'], location['sha']
    if sha != SHA or not base.startswith(f'/tmp/knowslink-schedqa-{SHA[:12]}-'):
        raise ValueError('unexpected task fixture')
    manifest = json.loads((HERE / 'fixture.json').read_text(encoding='utf-8'))
    target = urllib.parse.urlparse(manifest['origin'])
    if target.scheme != 'http' or target.hostname != '127.0.0.1' or not target.port:
        raise ValueError('non-loopback fixture target')
    client, _ = connect()

    class Forward(socketserver.BaseRequestHandler):
        def handle(self):
            channel = client.get_transport().open_channel('direct-tcpip', ('127.0.0.1', target.port), self.request.getpeername())
            if channel is None:
                return
            try:
                while True:
                    ready, _, _ = select.select([self.request, channel], [], [], 60)
                    if self.request in ready:
                        data = self.request.recv(65536)
                        if not data:
                            return
                        channel.sendall(data)
                    if channel in ready:
                        data = channel.recv(65536)
                        if not data:
                            return
                        self.request.sendall(data)
            finally:
                channel.close()

    server = socketserver.ThreadingTCPServer(('127.0.0.1', LOCAL_PORT), Forward)
    server.daemon_threads = True
    public = {'sha': SHA, 'local_origin': f'http://127.0.0.1:{LOCAL_PORT}', 'gates': manifest['gates'], 'accounts': manifest['accounts'], 'expires_at': manifest['expires_at']}
    (HERE / 'browser-fixture.json').write_text(json.dumps(public, indent=2), encoding='utf-8')
    print(f'Task loopback forwarding ready on port {LOCAL_PORT}', flush=True)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    threading.Event().wait(12 * 60 * 60)


if __name__ == '__main__':
    main()
