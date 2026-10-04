"""Start an isolated Postgres relay and record the independent trial observer."""
import json
import os
from pathlib import Path
import secrets
import shutil
import socket
import subprocess
import tempfile
import uuid

ROOT = Path(__file__).resolve().parent
CLONE = Path("/tmp/sar-mvp-003-qa")
COUNTS = {
    "owners": "SELECT count(*) FROM relay_state, LATERAL jsonb_object_keys(COALESCE(data->'Owners','{}'::jsonb)) AS k(name)",
    "gates": "SELECT count(*) FROM relay_state, LATERAL jsonb_object_keys(COALESCE(data->'Gates','{}'::jsonb)) AS k(name)",
    "trial": "SELECT count(*) FROM relay_state, LATERAL jsonb_each(COALESCE(data->'Messages','{}'::jsonb)) AS t(id, msg) WHERE msg->'Receipt'->>'intent' = 'relay.test.message'",
    "business": "SELECT count(*) FROM relay_state, LATERAL jsonb_each(COALESCE(data->'Messages','{}'::jsonb)) AS t(id, msg) WHERE COALESCE(msg->'Receipt'->>'intent','') <> '' AND msg->'Receipt'->>'intent' <> 'relay.test.message'",
    "payload": "SELECT count(*) FROM relay_state, LATERAL jsonb_each(COALESCE(data->'Messages','{}'::jsonb)) AS t(id, msg) WHERE octet_length(COALESCE(msg->>'Envelope','')) > 4 OR octet_length(COALESCE(msg->>'Inbox','')) > 4",
}


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


def run(command, env, cwd, out):
    result = subprocess.run(command, cwd=cwd, env=env, text=True, capture_output=True)
    secret = env.get("POSTGRES_PASSWORD", "")
    stdout = result.stdout.replace(secret, "[redacted]") if secret else result.stdout
    stderr = result.stderr.replace(secret, "[redacted]") if secret else result.stderr
    out.write_text(f"exit={result.returncode}\n{stdout}\n{stderr}")
    return result.returncode


def counts(compose, env):
    found = {}
    for name, query in COUNTS.items():
        result = subprocess.run(compose + ["exec", "-T", "postgres", "psql", "-U", "knowslink", "-d", "knowslink", "-Atc", query], cwd=CLONE, env=env, text=True, capture_output=True)
        if result.returncode != 0:
            found[name] = f"psql:{result.returncode}"
        else:
            found[name] = result.stdout.strip()
    return found


def node_phase(name, env, out_dir):
    phase_env = dict(env)
    phase_env["RESULT_PATH"] = str(out_dir / f"{name}.json")
    code = run(["node", str(ROOT / "observe.mjs"), name], phase_env, "/tmp", out_dir / f"{name}.log")
    result = json.loads(Path(phase_env["RESULT_PATH"]).read_text())
    return code, result


def main():
    out_dir = ROOT / "logs"
    out_dir.mkdir(exist_ok=True)
    password = secrets.token_urlsafe(24)
    port = free_port()
    env = dict(os.environ)
    env.update(POSTGRES_PASSWORD=password, DATABASE_URL=f"postgres://knowslink:{password}@postgres:5432/knowslink?sslmode=disable", RELAY_PORT=str(port), TUNNEL_TOKEN="", COMPOSE_PROFILES="", KNOWSLINK_TEST_AGENTS="trial_codex,trial_grok", CLONE=str(CLONE), RELAY_URL=f"http://127.0.0.1:{port}")
    project = "knowslink-qa003-" + uuid.uuid4().hex[:10]
    compose = ["docker", "compose", "--project-name", project, "--env-file", ".env.example", "-f", "compose.yaml"]
    summary = {"project": project, "relay_port": port, "phases": {}}
    state = Path(tempfile.mkdtemp(prefix="knowslink-qa-state-")) / "state.json"
    env["STATE_PATH"] = str(state)
    error = ""
    try:
        summary["up"] = run(compose + ["up", "--build", "--wait", "relay"], env, CLONE, out_dir / "compose-up.log")
        summary["before"] = counts(compose, env)
        code, result = node_phase("held", env, out_dir)
        summary["phases"]["held"] = {"exit": code, "ok": result.get("ok"), "failed": [item["name"] for item in result.get("checks", []) if not item.get("pass")]}
        summary["after_held"] = counts(compose, env)
        code, result = node_phase("live", env, out_dir)
        summary["phases"]["live"] = {"exit": code, "ok": result.get("ok"), "ids": result.get("ids"), "failed": [item["name"] for item in result.get("checks", []) if not item.get("pass")]}
        summary["after_live"] = counts(compose, env)
        closed_env = dict(env)
        closed_env["KNOWSLINK_TEST_AGENTS"] = ""
        summary["recreate"] = run(compose + ["up", "-d", "--no-deps", "--force-recreate", "--wait", "relay"], closed_env, CLONE, out_dir / "compose-recreate.log")
        shown = subprocess.run(compose + ["exec", "-T", "relay", "printenv", "KNOWSLINK_TEST_AGENTS"], cwd=CLONE, env=closed_env, text=True, capture_output=True)
        summary["allowlist_empty"] = shown.stdout == ""
        if state.exists():
            code, result = node_phase("closed", closed_env, out_dir)
            summary["phases"]["closed"] = {"exit": code, "ok": result.get("ok"), "failed": [item["name"] for item in result.get("checks", []) if not item.get("pass")]}
        summary["after_closed"] = counts(compose, closed_env)
    except Exception as exc:
        error = type(exc).__name__
    finally:
        summary["down"] = run(compose + ["down", "--volumes"], env, CLONE, out_dir / "compose-down.log")
        shutil.rmtree(state.parent, ignore_errors=True)
    summary["error"] = error
    (out_dir / "summary.json").write_text(json.dumps(summary, indent=2) + "\n")
    expect_empty = {"owners": "0", "gates": "0", "trial": "0", "business": "0", "payload": "0"}
    expect_live = {"owners": "2", "gates": "0", "trial": "4", "business": "0", "payload": "0"}
    ok = summary.get("up") == 0 and summary.get("down") == 0 and summary.get("before") == expect_empty and summary.get("after_held") == expect_empty and summary.get("after_live") == expect_live and summary.get("after_closed") == expect_live and summary.get("allowlist_empty") is True and all(phase.get("ok") for phase in summary["phases"].values()) and set(summary["phases"]) == {"held", "live", "closed"}
    print(json.dumps({"ok": ok, "phases": summary.get("phases"), "counts": {key: summary.get(key) for key in ("before", "after_held", "after_live", "after_closed")}}, indent=2))
    raise SystemExit(0 if ok else 1)


if __name__ == "__main__":
    main()
