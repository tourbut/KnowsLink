"""Loopback SMTP sink for the synthetic ART member: stores each message body in <out>/last.txt (no real mail is sent)."""
import socketserver
import sys
from pathlib import Path

OUT = Path(sys.argv[1])


class H(socketserver.StreamRequestHandler):
    def handle(self):
        def w(s):
            self.wfile.write((s + "\r\n").encode())
            self.wfile.flush()

        w("220 sink ESMTP")
        data, body = False, []
        for raw in self.rfile:
            line = raw.decode("utf-8", "replace").rstrip("\r\n")
            if data:
                if line == ".":
                    (OUT / "last.txt").write_text("\n".join(body), encoding="utf-8")
                    data, body = False, []
                    w("250 ok")
                else:
                    body.append(line)
                continue
            cmd = line.upper()
            if cmd.startswith(("EHLO", "HELO")):
                w("250 sink")
            elif cmd.startswith("DATA"):
                data = True
                w("354 go")
            elif cmd.startswith("QUIT"):
                w("221 bye")
                return
            else:
                w("250 ok")


socketserver.ThreadingTCPServer.allow_reuse_address = True
socketserver.ThreadingTCPServer(("127.0.0.1", int(sys.argv[2])), H).serve_forever()
