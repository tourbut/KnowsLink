"""Verify synthetic MVP in a private Compose project without touching existing services or Tunnel."""

import os
from pathlib import Path
import socket
import secrets
import subprocess
import tempfile
import uuid

ROOT = Path(__file__).resolve().parents[1]


def run(command, environment):
    result = subprocess.run(command, cwd=ROOT, env=environment, text=True)
    print(f"command: {command}; exit: {result.returncode}", flush=True)
    if result.returncode:
        raise RuntimeError("synthetic verification command failed")


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


def main():
    environment = dict(os.environ)
    test_password = secrets.token_urlsafe(24)
    relay_port, database_port = free_port(), free_port()
    environment.update(POSTGRES_PASSWORD=test_password, DATABASE_URL=f"postgres://knowslink:{test_password}@postgres:5432/knowslink?sslmode=disable", RELAY_PORT=str(relay_port), TUNNEL_TOKEN="", COMPOSE_PROFILES="", KNOWSLINK_TEST_AGENTS="trial_codex,trial_grok")
    project = "knowslink-mvp-" + uuid.uuid4().hex[:10]
    with tempfile.TemporaryDirectory(prefix="knowslink-mvp-") as temporary:
        override = Path(temporary) / "test.yaml"
        override.write_text(f'services:\n  postgres:\n    ports: ["127.0.0.1:{database_port}:5432"]\n    networks: [database, ingress]\n')
        compose = ["docker", "compose", "--project-name", project, "--env-file", ".env.example", "-f", "compose.yaml", "-f", str(override)]
        try:
            run(compose + ["up", "--build", "--wait", "relay"], environment)
            tests = dict(environment, TEST_SYNTHETIC_DATABASE="1", TEST_DATABASE_URL=f"postgres://knowslink:{test_password}@127.0.0.1:{database_port}/knowslink?sslmode=disable")
            run(["go", "test", "-tags=integration", "-race", "-count=1", "-v", "./internal/relay"], tests)
            run(["node", "adapters/dist/synthetic.js", f"http://127.0.0.1:{relay_port}"], environment)
            run(["node", "adapters/dist/synthetic.js", f"http://127.0.0.1:{relay_port}", "--seed"], environment)
            run(["node", "adapters/dist/trial-check.js", f"http://127.0.0.1:{relay_port}"], environment)
        finally:
            run(compose + ["down", "--volumes"], environment)
    print("PASS: isolated Compose migration, Postgres races, TS adapter and Go owner UI; Tunnel unused")


if __name__ == "__main__":
    main()
