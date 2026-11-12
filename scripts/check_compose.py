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
environment.pop("KNOWSLINK_SYNTHETIC_SIGNUP", None)
command = [
    "docker", "compose", "--env-file", ".env.example", "--profile", "tunnel",
    "config", "--format", "json",
]
result = subprocess.run(command, env=environment, text=True, capture_output=True)
if result.returncode:
    raise SystemExit("Compose example configuration failed")
services = json.loads(result.stdout)["services"]
assert set(services) == {"postgres", "redis", "migrate", "relay", "cloudflared"}
assert not services["postgres"].get("ports"), "Postgres must not publish ports"
redis = services["redis"]
assert not redis.get("ports") and list(redis["networks"]) == ["database"], "Redis must stay on the internal network"
assert redis["command"][redis["command"].index("--save") + 1] == "" and redis["command"][redis["command"].index("--appendonly") + 1] == "no", "Redis must not persist to disk"
assert redis["command"][redis["command"].index("--maxmemory-policy") + 1] == "noeviction", "Redis must refuse writes, not evict messages"
assert services["relay"]["depends_on"]["redis"]["condition"] == "service_healthy"
assert services["migrate"]["command"] == ["/app/migrate"]
assert services["migrate"]["restart"] == "no"
assert services["migrate"]["depends_on"]["postgres"]["condition"] == "service_healthy"
assert services["relay"]["depends_on"]["migrate"]["condition"] == "service_completed_successfully"
assert services["relay"]["ports"][0]["host_ip"] == "127.0.0.1"
assert services["cloudflared"]["profiles"] == ["tunnel"]
assert services["cloudflared"]["environment"]["TUNNEL_TOKEN"] == ""
assert services["relay"]["environment"]["KNOWSLINK_SYNTHETIC_SIGNUP"] == "", "Synthetic signup must be closed by default"
for variable in ("DATABASE_URL", "POSTGRES_PASSWORD"):
    missing = dict(environment)
    missing[variable] = ""
    failure = subprocess.run(command, env=missing, text=True, capture_output=True)
    assert failure.returncode != 0 and f"{variable} is required" in failure.stderr
ops = dict(environment, KNOWSLINK_STATE_DIR="/srv/knowslink-state", KNOWSLINK_UID="1000", KNOWSLINK_GID="1000")
result = subprocess.run(command[:2] + ["-f", "compose.yaml", "-f", "deploy/knowslink/compose.ops.yaml"] + command[2:], env=ops, text=True, capture_output=True)
if result.returncode:
    raise SystemExit("Compose OPS override configuration failed:\n" + result.stderr)
merged = json.loads(result.stdout)["services"]
assert set(merged) == set(services) and not merged["postgres"].get("ports"), "OPS override must keep the services and private Postgres"
assert all(merged[name]["read_only"] for name in ("redis", "migrate", "relay", "cloudflared")), "OPS containers must be read-only"
assert not merged["redis"].get("ports") and merged["redis"]["mem_limit"], "OPS Redis must stay private and bounded"
assert any(volume["target"] == "/install" and volume.get("read_only") for volume in merged["relay"]["volumes"]), "relay /install must be read-only"
assert not (merged["cloudflared"].get("environment") or {}).get("TUNNEL_TOKEN"), "OPS cloudflared must use the config file, not a token"
print("Compose: five services, one-shot migrate, private Postgres, volatile private Redis, loopback relay, optional Tunnel, synthetic signup closed, OPS override merges")
