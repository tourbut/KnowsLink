"""Check local Compose readiness and safe failure without activating a Tunnel."""

import json
import os
from pathlib import Path
import socket
import shlex
import subprocess
import tempfile
import uuid

ROOT = Path(__file__).resolve().parents[1]


def main():
    environment = dict(os.environ)
    environment.update(
        POSTGRES_PASSWORD="example-local-only",
        DATABASE_URL="postgres://knowslink:example-local-only@postgres:5432/knowslink?sslmode=disable",
        TUNNEL_TOKEN="",
        COMPOSE_PROFILES="",
    )
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        environment["RELAY_PORT"] = str(listener.getsockname()[1])
    project = "knowslink-check-" + uuid.uuid4().hex[:10]
    compose = ["docker", "compose", "--env-file", ".env.example", "-p", project]

    def run(command, expected=0, diagnostic=None):
        result = subprocess.run(command, cwd=ROOT, env=environment, text=True, capture_output=True)
        print(f"command: {shlex.join(command)}; exit: {result.returncode}", flush=True)
        print(result.stdout + result.stderr, end="", flush=True)
        assert (result.returncode == 0) if expected == 0 else (result.returncode != 0)
        if diagnostic:
            assert diagnostic in result.stdout + result.stderr, diagnostic
        return result.stdout

    # Temporary storage avoids deleting or retaining a persistent Postgres volume.
    # This override changes only validation storage; product Compose stays unchanged.
    configuration = subprocess.run(compose + ["--profile", "tunnel", "config", "--format", "json"],
                                   cwd=ROOT, env=environment, text=True, capture_output=True, check=True)
    cloudflared = json.loads(configuration.stdout)["services"]["cloudflared"]["image"]
    with tempfile.TemporaryDirectory(prefix="knowslink-runtime-") as temporary:
        override = Path(temporary) / "compose.json"
        override.write_text(json.dumps({"services": {"postgres": {
            "volumes": [{"type": "tmpfs", "target": "/var/lib/postgresql/data"}],
        }}}))
        compose += ["-f", str(ROOT / "compose.yaml"), "-f", str(override)]
        try:
            run(compose + ["up", "--build", "--wait", "--wait-timeout", "120", "relay"])
            run(compose + ["logs", "migrate"], diagnostic="SQL migrations completed")
            tables = run(compose + ["exec", "-T", "postgres", "psql", "-U", "knowslink", "-d", "knowslink", "-Atc",
                                   "SELECT count(*) FROM information_schema.tables WHERE table_schema='public'"])
            assert tables.strip() == "2", "relay_state and goose version table must exist"
            run(compose + ["exec", "-T", "relay", "wget", "-q", "-O", "-", "http://127.0.0.1:8080/healthz"], diagnostic="ok")
            for binary in ("relay", "migrate"):
                run(compose + ["exec", "-T", "relay", "env", "DATABASE_URL=", f"/app/{binary}"],
                    expected=1, diagnostic="DATABASE_URL is required")
                run(compose + ["exec", "-T", "relay", "env", "DATABASE_URL=postgres://localhost:1/knowslink", f"/app/{binary}"],
                    expected=1, diagnostic="database ping failed")
            run(compose + ["stop", "postgres"])
            run(compose + ["exec", "-T", "relay", "wget", "-q", "-O", "-", "http://127.0.0.1:8080/healthz"],
                expected=1, diagnostic="503")
        finally:
            run(compose + ["down", "--volumes"])
    run(["docker", "run", "--rm", "--network", "none", "-e", "TUNNEL_TOKEN=",
         cloudflared, "tunnel", "--no-autoupdate", "run"],
        expected=1, diagnostic="requires the ID or name of the tunnel")
    run(["npm", "run", "start", "--prefix", "adapters"], diagnostic='"state":"unconfigured"')
    print("PASS: local startup, SQL migration, readiness, missing settings, DB failure, disabled Tunnel, unconfigured adapter")


if __name__ == "__main__":
    main()
