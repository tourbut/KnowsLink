"""Expose only the four task-owned synthetic loopback fixtures through verified SSH.

The manifest binds forwarding to the prepared fixture SHA; no operational service is forwarded.
"""
import json
import select
import socketserver
import threading
import urllib.parse

from remote import HERE, connect


def main():
    location = json.loads((HERE / 'fixture-location.json').read_text(encoding='utf-8'))
    base, sha = location['base'], location['sha']
    if not base.startswith(f'/tmp/knowslink-landing001-{sha[:12]}-') or not all(c.isalnum() or c in '/-' for c in base):
        raise ValueError('unexpected task fixture')
    client, _ = connect()
    with client.open_sftp() as sftp:
        with sftp.open(base + '/fixture.json') as stream:
            targets = json.load(stream)
    if set(targets) != {'google', 'both', 'email', 'none'}:
        raise ValueError('unexpected fixture modes')
    manifest = {'sha': sha, 'base': base, 'container': location['container'], 'origins': {}}
    servers = []

    def handler_for(port):
        class Forward(socketserver.BaseRequestHandler):
            def handle(self):
                channel = client.get_transport().open_channel('direct-tcpip', ('127.0.0.1', port), self.request.getpeername())
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
        return Forward

    try:
        for offset, mode in enumerate(['google', 'both', 'email', 'none']):
            target = urllib.parse.urlparse(targets[mode])
            if target.scheme != 'http' or target.hostname != '127.0.0.1' or not target.port:
                raise ValueError('non-loopback fixture target')
            server = socketserver.ThreadingTCPServer(('127.0.0.1', 18781 + offset), handler_for(target.port))
            server.daemon_threads = True
            servers.append(server)
            threading.Thread(target=server.serve_forever, daemon=True).start()
            manifest['origins'][mode] = f'http://127.0.0.1:{18781 + offset}'
        (HERE / 'browser-fixture.json').write_text(json.dumps(manifest, indent=2), encoding='utf-8')
        print('Task loopback forwarding ready on ports 18781-18784', flush=True)
        threading.Event().wait(24 * 60 * 60)
    finally:
        for server in servers:
            server.shutdown()
            server.server_close()
        client.close()


if __name__ == '__main__':
    main()
