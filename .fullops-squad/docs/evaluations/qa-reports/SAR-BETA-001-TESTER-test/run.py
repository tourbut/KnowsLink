#!/usr/bin/env python3
"""Narrow QA of the live 437f143 beta runtime. Prints status codes and modes only."""

import hashlib
import http.client
import json
import re
import subprocess
from pathlib import Path

DEPLOY = Path("/home/shin/deploy/knowslink")
STATE = Path("/home/shin/deploy/knowslink-state")
FIXED = "437f1432a158670a485413c1aba159debc3759e5"
PRODUCT = "78b1d92c8aa626245d3349ffaf7367d28f1dd3ef"
EVIDENCE = Path(__file__).resolve().parent
SCRIPT = DEPLOY / "deploy/knowslink"
ALLOWED = {"invalid_auth", "invalid_csrf", "sender_not_allowed"}
CHECKS = []


def redact(text):
    text = re.sub(r"postgres://\S+", "postgres://REDACTED", text)
    text = re.sub(r"(POSTGRES_PASSWORD=)\S+", r"\1REDACTED", text)
    text = re.sub(r"(?i)(authorization:\s*)\S+", r"\1REDACTED", text)
    text = re.sub(r"(?i)(bearer\s+)\S+", r"\1REDACTED", text)
    text = re.sub(r"(?i)(basic\s+)\S+", r"\1REDACTED", text)
    return re.sub(r"[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}", "REDACTED_EMAIL", text)


def check(name, ok, detail):
    CHECKS.append({"name": name, "pass": bool(ok), "detail": detail})
    print(f"{'PASS' if ok else 'FAIL'} {name}: {detail}", flush=True)


def capture(args):
    completed = subprocess.run(args, text=True, capture_output=True)
    out, err = redact(completed.stdout), redact(completed.stderr)
    return completed.returncode, out, err


def save(name, text):
    (EVIDENCE / name).write_text(text, encoding="utf-8")


def request(method, path, headers=None, body=None):
    connection = http.client.HTTPConnection("127.0.0.1", 8080, timeout=15)
    try:
        connection.request(method, path, body=body, headers=headers or {})
        response = connection.getresponse()
        return response.status, response.read(100000)
    finally:
        connection.close()


def safe_body(payload):
    text = payload.decode(errors="replace").strip()
    if text in ALLOWED:
        return text
    try:
        code = json.loads(text).get("error")
    except json.JSONDecodeError:
        return f"len={len(payload)}"
    return code if code in ALLOWED else f"len={len(payload)}"


def mode_of(path):
    return oct(path.stat().st_mode & 0o777)[2:]


def shared_state():
    code, out, err = capture(["python3", str(SCRIPT / "verify.py"), "regression"])
    return code, out, err


def limits():
    selected = {}
    for name in ("knowslink-postgres-1", "knowslink-relay-1", "knowslink-migrate-1"):
        raw = subprocess.check_output(["docker", "inspect", name])
        host = json.loads(raw)[0]["HostConfig"]
        selected[name] = {
            "RestartPolicy": host.get("RestartPolicy", {}).get("Name"),
            "Memory": host.get("Memory"),
            "NanoCpus": host.get("NanoCpus"),
            "PidsLimit": host.get("PidsLimit"),
            "CapDrop": host.get("CapDrop"),
            "ReadonlyRootfs": host.get("ReadonlyRootfs"),
            "SecurityOpt": host.get("SecurityOpt"),
            "PortBindings": host.get("PortBindings"),
            "LogConfig": host.get("LogConfig"),
        }
    return selected


def credential_modes():
    files = [path for path in (STATE / "tunnel").iterdir() if path.suffix == ".json"]
    return sorted(mode_of(path) for path in files), len(files)


def negative():
    fixture = json.loads((DEPLOY / "build/qa-fixture.json").read_text(encoding="utf-8"))
    owner = fixture["b"]["owner"]["credential"]
    agent = fixture["b"]["credential"]
    del fixture
    owner_headers = {"Authorization": "Bearer " + owner}
    agent_headers = {"Authorization": "Bearer " + agent}

    status, body = request("GET", "/owner")
    check("NEG-no-auth", status == 401 and safe_body(body) == "invalid_auth", f"{status} {safe_body(body)}")
    status, body = request("GET", "/owner", agent_headers)
    check("NEG-agent-owner", status == 401 and safe_body(body) == "invalid_auth", f"{status} {safe_body(body)}")
    status, page = request("GET", "/owner", owner_headers)
    check("NEG-owner-still-active-before", status == 200, str(status))
    status, body = request("POST", "/v1/owner-revoke", {**agent_headers, "Content-Type": "application/json"}, b"{}")
    check("NEG-agent-revoke", status == 401 and safe_body(body) == "invalid_auth", f"{status} {safe_body(body)}")
    status, _page = request("GET", "/owner", owner_headers)
    check("NEG-owner-still-active-after", status == 200, str(status))
    status, before = request("GET", "/owner", owner_headers)
    before_gates = re.findall(r'상태: [^<]+', before.decode(errors="replace")) if status == 200 else None
    form = {"Authorization": "Bearer " + owner, "Content-Type": "application/x-www-form-urlencoded"}
    status, body = request("POST", "/owner/gates/qa-absent", form, b"csrf=bad&decision=approve")
    check("NEG-bad-csrf", status == 403 and safe_body(body) == "invalid_csrf", f"{status} {safe_body(body)}")
    status, after = request("GET", "/owner", owner_headers)
    after_gates = re.findall(r'상태: [^<]+', after.decode(errors="replace")) if status == 200 else None
    same = before_gates is not None and before_gates == after_gates
    check("NEG-csrf-no-decision", same, f"states={len(before_gates or [])} unchanged={same}")


def main():
    EVIDENCE.mkdir(parents=True, exist_ok=True)
    code, out, _err = capture(["git", "-C", str(DEPLOY), "rev-parse", "HEAD"])
    head = out.strip()
    code_short, short_out, _err = capture(["git", "-C", str(DEPLOY), "rev-parse", "--short", "HEAD"])
    short = short_out.strip()
    code_status, status, _err = capture(["git", "-C", str(DEPLOY), "status", "--porcelain"])
    ancestor, _, _ = capture(["git", "-C", str(DEPLOY), "merge-base", "--is-ancestor", FIXED, head or "HEAD"])
    code_log, reflog, reflog_err = capture(["git", "-C", str(DEPLOY), "reflog", "-8", "--date=iso"])
    save("reflog.log", f"[exit {code_log}]\n{reflog}{reflog_err}")
    clean = code == 0 and code_short == 0 and code_status == 0 and status == "" and ancestor == 0 and bool(head)
    check("deploy-lineage", clean, f"head={head} short={short} fixed_ancestor={ancestor == 0} equal_fixed={head == FIXED}")
    code, diff, _err = capture(["git", "-C", str(DEPLOY), "diff", "--stat", PRODUCT, "HEAD", "--", ":!.fullops-squad", ":!deploy"])
    check("product-unchanged", code == 0 and diff.strip() == "", "empty" if diff.strip() == "" else "diff")

    bound = limits()
    save("limits.json", json.dumps(bound, indent=2) + "\n")
    relay = bound["knowslink-relay-1"]
    postgres = bound["knowslink-postgres-1"]
    migrate = bound["knowslink-migrate-1"]
    loopback = relay["PortBindings"] == {"8080/tcp": [{"HostIp": "127.0.0.1", "HostPort": "8080"}]}
    check("loopback-8080", loopback and postgres["PortBindings"] == {}, "relay 127.0.0.1:8080 postgres unbound")
    code, listeners, _err = capture(["ss", "-ltn"])
    wide = re.search(r"(0\.0\.0\.0|\[::\]|\*):8080\b", listeners)
    local_bind = re.search(r"127\.0\.0\.1:8080\b", listeners)
    check("ss-loopback", code == 0 and local_bind and not wide, "127.0.0.1:8080 only" if not wide else "wide bind")
    resources_ok = (
        postgres["RestartPolicy"] == "unless-stopped" and postgres["Memory"] == 536870912 and postgres["NanoCpus"] == 1000000000
        and relay["RestartPolicy"] == "unless-stopped" and relay["Memory"] == 268435456 and relay["CapDrop"] == ["ALL"]
        and relay["ReadonlyRootfs"] is True and relay["SecurityOpt"] == ["no-new-privileges:true"]
        and migrate["RestartPolicy"] == "no" and migrate["Memory"] == 134217728 and migrate["CapDrop"] == ["ALL"]
    )
    check("restart-limits", resources_ok, "unless-stopped memory caps drop")
    log_ok = all(item["LogConfig"]["Config"].get("max-size") == "10m" and item["LogConfig"]["Config"].get("max-file") == "3" for item in bound.values())
    check("log-limit", log_ok, "10m x 3")
    modes, count = credential_modes()
    check("secret-mode", mode_of(STATE / ".env") == "600" and modes == ["600"] and count == 1 and mode_of(STATE) == "700", f"env={mode_of(STATE / '.env')} cred={modes} state={mode_of(STATE)}")
    env_size = (STATE / ".env").stat().st_size
    aud_before = (STATE / "access.aud").exists()
    baseline_hash = hashlib.sha256((STATE / "shared-baseline.json").read_bytes()).hexdigest()
    dumps_before = sorted(path.name for path in (STATE / "backups").iterdir())

    code, out, err = capture(["python3", str(SCRIPT / "verify.py"), "local"])
    save("verify-local.log", f"[exit {code}]\n{out}{err}")
    check("verify-local", code == 0, f"exit {code}")
    code, out, err = shared_state()
    save("verify-regression-before.log", f"[exit {code}]\n{out}{err}")
    check("verify-regression-before", code == 0, f"exit {code}")
    code, out, err = capture(["python3", str(SCRIPT / "access_apply.py"), "selftest"])
    save("access-selftest.log", f"[exit {code}]\n{out}{err}")
    check("access-selftest", code == 0 and aud_before is False, f"exit {code} aud_absent={not aud_before}")

    try:
        negative()
    except Exception as exc:
        check("negative-probe", False, type(exc).__name__)

    code, out, err = capture([str(SCRIPT / "beta.sh"), "backup"])
    save("backup.log", f"[exit {code}]\n{out}{err}")
    dumps_after = sorted(path.name for path in (STATE / "backups").iterdir())
    created = [name for name in dumps_after if name not in dumps_before]
    backup_ok = code == 0 and len(created) == 1 and created[0].startswith(short + "-") and mode_of(STATE / "backups" / created[0]) == "600"
    check("backup", backup_ok, created[0] if created else f"exit {code}")
    kept = all(name in dumps_after for name in dumps_before)
    check("prior-backup-kept", kept, f"before={len(dumps_before)} after={len(dumps_after)}")

    if backup_ok:
        code, out, err = capture([str(SCRIPT / "beta.sh"), "restore-verify", str(STATE / "backups" / created[0])])
        save("restore.log", f"[exit {code}]\n{out}{err}")
        _code, names, _err = capture(["docker", "ps", "-aq", "--filter", "name=knowslink-restore-check"])
        health, _body = request("GET", "/healthz")
        volume_code, volume, _err = capture(["docker", "volume", "ls", "-q", "--filter", "name=knowslink_postgres-data"])
        restored = code == 0 and "tables=" in out and names.strip() == "" and health == 200 and volume.strip() == "knowslink_postgres-data"
        check("restore-isolated", restored, redact(out.strip()) or f"exit {code}")
        if names.strip():
            capture(["docker", "rm", "-f", "knowslink-restore-check"])
            check("restore-container-removed", False, "leftover removed by this QA")
    else:
        check("restore-isolated", False, "backup did not create one dump")

    if aud_before:
        check("expose-blocked", False, "access.aud present; expose not run")
    else:
        code, out, err = capture([str(SCRIPT / "beta.sh"), "expose"])
        save("expose.log", f"[exit {code}]\n{out}{err}")
        _code, dig, _err = capture(["dig", "+short", "link.knowslog.com"])
        _code, containers, _err = capture(["docker", "ps", "-a", "--filter", "name=knowslink-cloudflared", "--format", "{{.Names}}"])
        blocked = code == 1 and "Access app not recorded" in err and dig.strip() == "" and containers.strip() == "" and not (STATE / "access.aud").exists()
        check("expose-blocked", blocked, f"exit {code} dns_empty={dig.strip() == ''}")

    code, out, err = shared_state()
    save("verify-regression-after.log", f"[exit {code}]\n{out}{err}")
    check("verify-regression-after", code == 0, f"exit {code}")
    env_same = (STATE / ".env").stat().st_size == env_size and mode_of(STATE / ".env") == "600"
    baseline_same = hashlib.sha256((STATE / "shared-baseline.json").read_bytes()).hexdigest() == baseline_hash
    code, status, _err = capture(["git", "-C", str(DEPLOY), "status", "--porcelain"])
    code_end, end_out, _err = capture(["git", "-C", str(DEPLOY), "rev-parse", "HEAD"])
    check("deploy-stable", code_end == 0 and end_out.strip() == head, end_out.strip())
    check("deploy-untouched", code == 0 and status == "" and env_same and baseline_same, "checkout clean env unchanged baseline unchanged")
    save("results.json", json.dumps(CHECKS, indent=2) + "\n")
    failed = [item["name"] for item in CHECKS if not item["pass"]]
    print(f"failed={failed}", flush=True)
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
