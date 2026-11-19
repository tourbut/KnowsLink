#!/usr/bin/env python3
"""Run the RF-01 narrow probe against a private Compose project."""

import os
from pathlib import Path
import secrets
import socket
import subprocess
import tempfile
import uuid

ROOT = Path(__file__).resolve().parents[5]
PROBE = ROOT / ".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL-test/probe.mjs"


def run(command, environment):
    completed = subprocess.run(command, cwd=ROOT, env=environment)
    print(f"command: {' '.join(command)}; exit: {completed.returncode}", flush=True)
    return completed.returncode


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


def main():
    environment = dict(os.environ)
    password = secrets.token_urlsafe(24)
    relay_port, database_port = free_port(), free_port()
    environment.update(
        POSTGRES_PASSWORD=password,
        DATABASE_URL=f"postgres://knowslink:{password}@postgres:5432/knowslink?sslmode=disable",
        RELAY_PORT=str(relay_port),
        TUNNEL_TOKEN="",
        COMPOSE_PROFILES="",
    )
    project = "knowslink-final-" + uuid.uuid4().hex[:10]
    with tempfile.TemporaryDirectory(prefix="knowslink-final-") as temporary:
        override = Path(temporary) / "test.yaml"
        override.write_text(
            "services:\n  postgres:\n    ports: [\"127.0.0.1:%s:5432\"]\n    networks: [database, ingress]\n" % database_port
        )
        compose = ["docker", "compose", "--project-name", project, "--env-file", ".env.example", "-f", "compose.yaml", "-f", str(override)]
        probe_env = dict(environment, RELAY_URL=f"http://127.0.0.1:{relay_port}", COMPOSE_PROJECT=project, COMPOSE_OVERRIDE=str(override))
        status = 1
        try:
            status = run(compose + ["up", "--build", "--wait", "relay"], environment)
            if status == 0:
                status = run(["node", str(PROBE)], probe_env)
        finally:
            down = run(compose + ["down", "--volumes"], environment)
            if status == 0:
                status = down
        raise SystemExit(status)


if __name__ == "__main__":
    main()
