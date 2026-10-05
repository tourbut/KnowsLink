"""Local QA SMTP sink: stores each received message as a 0600 file. Never use as a real mail service."""

import argparse
import os
from pathlib import Path
import socketserver
import time


class Sink(socketserver.StreamRequestHandler):
    def reply(self, line):
        self.wfile.write(line.encode() + b"\r\n")

    def handle(self):
        self.reply("220 knowslink-qa-sink")
        while line := self.rfile.readline():
            command = line.strip().upper()
            if command.startswith(b"EHLO") or command.startswith(b"HELO"):
                self.reply("250 knowslink-qa-sink")
            elif command == b"DATA":
                self.reply("354 end with .")
                data = b""
                while (part := self.rfile.readline()) not in (b".\r\n", b""):
                    data += part
                path = self.server.directory / f"{time.time_ns()}.eml"
                descriptor = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
                with os.fdopen(descriptor, "wb") as output:
                    output.write(data)
                print(f"stored message {path.name}", flush=True)
                self.reply("250 stored")
            elif command == b"QUIT":
                self.reply("221 bye")
                return
            else:
                self.reply("250 ok")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("directory", type=Path, help="private output folder, created 0700")
    parser.add_argument("--listen", default="127.0.0.1", help="use the Docker bridge address only for local Compose QA")
    parser.add_argument("--port", type=int, default=2525)
    args = parser.parse_args()
    args.directory.mkdir(mode=0o700, parents=True, exist_ok=True)
    os.chmod(args.directory, 0o700)
    with socketserver.ThreadingTCPServer((args.listen, args.port), Sink) as server:
        server.directory = args.directory
        print(f"QA sink on {args.listen}:{args.port}; messages in {args.directory}", flush=True)
        server.serve_forever()


if __name__ == "__main__":
    main()
