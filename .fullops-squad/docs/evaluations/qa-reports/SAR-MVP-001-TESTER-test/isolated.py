#!/usr/bin/env python3
"""Run the SAR-MVP-001 probe in a private Compose project and capture owner pages."""

import os
from pathlib import Path
import secrets
import socket
import subprocess
import tempfile
import uuid

ROOT = Path(__file__).resolve().parents[5]
UI = ROOT / ".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui"
PROBE = ROOT / ".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/probe.mjs"


def run(command, environment, cwd=ROOT):
    completed = subprocess.run(command, cwd=cwd, env=environment)
    print(f"command: {' '.join(command)}; exit: {completed.returncode}", flush=True)
    return completed.returncode


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


def capture(environment):
    for page in sorted(UI.glob("*.html")):
        image = page.with_suffix(".png")
        code = run([
            "google-chrome", "--headless=new", "--disable-gpu", "--no-sandbox", "--hide-scrollbars",
            "--window-size=1280,900", f"--user-data-dir={tempfile.mkdtemp(prefix='knowslink-chrome-')}",
            f"--screenshot={image}", page.as_uri(),
        ], environment)
        print(f"screenshot: {page.name}; exit: {code}", flush=True)


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
    project = "knowslink-qa-" + uuid.uuid4().hex[:10]
    with tempfile.TemporaryDirectory(prefix="knowslink-qa-") as temporary:
        override = Path(temporary) / "test.yaml"
        override.write_text(
            f"services:\n  postgres:\n    ports: [\"127.0.0.1:{database_port}:5432\"]\n    networks: [database, ingress]\n"
        )
        compose = ["docker", "compose", "--project-name", project, "--env-file", ".env.example", "-f", "compose.yaml", "-f", str(override)]
        probe_env = dict(environment, RELAY_URL=f"http://127.0.0.1:{relay_port}", COMPOSE_PROJECT=project, COMPOSE_OVERRIDE=str(override))
        status = 1
        try:
            status = run(compose + ["up", "--build", "--wait", "relay"], environment)
            if status == 0:
                status = run(["node", str(PROBE)], probe_env)
                capture(environment)
        finally:
            run(compose + ["down", "--volumes"], environment)
        raise SystemExit(status)


if __name__ == "__main__":
    main()
