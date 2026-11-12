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
            storage = run(compose + ["exec", "-T", "postgres", "psql", "-U", "knowslink", "-d", "knowslink", "-Atc",
                                     "SELECT atttypid::regtype, attcompression FROM pg_attribute "
                                     "WHERE attrelid='relay_state'::regclass AND attname='state'"])
            assert storage.strip() == "json|l", "relay_state.state must be json with lz4 compression"
            run(compose + ["exec", "-T", "relay", "wget", "-q", "-O", "-", "http://127.0.0.1:8080/healthz"], diagnostic="ok")
            # Messages are volatile: Redis keeps no RDB/AOF and refuses writes at its cap instead of evicting.
            for setting, value in (("save", ""), ("appendonly", "no"), ("maxmemory-policy", "noeviction")):
                got = run(compose + ["exec", "-T", "redis", "redis-cli", "--raw", "CONFIG", "GET", setting]).splitlines()
                assert got[1:2] == [value] or (value == "" and got[1:] in ([], [""])), f"Redis {setting} must be {value!r}"
            run(compose + ["exec", "-T", "relay", "env", "REDIS_URL=", "/app/relay"], expected=1, diagnostic="REDIS_URL is required")
            run(compose + ["exec", "-T", "relay", "env", "REDIS_URL=redis://localhost:1/0", "/app/relay"],
                expected=1, diagnostic="redis ping failed")
            # A restart keeps nothing: no RDB/AOF means the probe key is gone afterwards.
            run(compose + ["exec", "-T", "redis", "redis-cli", "SET", "knowslink:restart-probe", "1"], diagnostic="OK")
            run(compose + ["stop", "redis"])
            run(compose + ["exec", "-T", "relay", "wget", "-q", "-O", "-", "http://127.0.0.1:8080/healthz"],
                expected=1, diagnostic="503")
            run(compose + ["start", "redis"])
            run(compose + ["exec", "-T", "redis", "sh", "-c",
                           "for i in $(seq 20); do redis-cli ping >/dev/null 2>&1 && break; sleep 1; done; redis-cli EXISTS knowslink:restart-probe"],
                diagnostic="0")
            run(compose + ["exec", "-T", "relay", "sh", "-c",
                           "for i in $(seq 20); do wget -q -O - http://127.0.0.1:8080/healthz && exit 0; sleep 1; done; exit 1"],
                diagnostic="ok")
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
    print("PASS: local startup, SQL migration, readiness, missing settings, volatile Redis, Redis and DB failure, disabled Tunnel, unconfigured adapter")


if __name__ == "__main__":
    main()
