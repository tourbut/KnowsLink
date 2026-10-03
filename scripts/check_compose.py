"""Validate isolated example Compose settings and locked service boundaries."""

import json
import os
import subprocess

environment = dict(os.environ)
environment.update(
    POSTGRES_PASSWORD="example-local-only",
    DATABASE_URL="postgres://knowslink:example-local-only@postgres:5432/knowslink?sslmode=disable",
    TUNNEL_TOKEN="",
    RELAY_PORT="8080",
)
command = [
    "docker", "compose", "--env-file", ".env.example", "--profile", "tunnel",
    "config", "--format", "json",
]
result = subprocess.run(command, env=environment, text=True, capture_output=True)
if result.returncode:
    raise SystemExit("Compose example configuration failed")
services = json.loads(result.stdout)["services"]
assert set(services) == {"postgres", "migrate", "relay", "cloudflared"}
assert not services["postgres"].get("ports"), "Postgres must not publish ports"
assert services["migrate"]["command"] == ["/app/migrate"]
assert services["migrate"]["restart"] == "no"
assert services["migrate"]["depends_on"]["postgres"]["condition"] == "service_healthy"
assert services["relay"]["depends_on"]["migrate"]["condition"] == "service_completed_successfully"
assert services["relay"]["ports"][0]["host_ip"] == "127.0.0.1"
assert services["cloudflared"]["profiles"] == ["tunnel"]
assert services["cloudflared"]["environment"]["TUNNEL_TOKEN"] == ""
for variable in ("DATABASE_URL", "POSTGRES_PASSWORD"):
    missing = dict(environment)
    missing[variable] = ""
    failure = subprocess.run(command, env=missing, text=True, capture_output=True)
    assert failure.returncode != 0 and f"{variable} is required" in failure.stderr
print("Compose: four services, one-shot migrate, private Postgres, loopback relay, optional Tunnel")
