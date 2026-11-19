#!/usr/bin/env python3
"""Narrow fail-closed checks for the 28bd1bb Access gate. Temp state only."""

import hashlib
import json
import os
import re
import shutil
import subprocess
import tempfile
import time
from pathlib import Path

SNAP = Path("/tmp/knowslink-beta-review-28bd1bb")
LIVE = Path("/home/shin/deploy/knowslink")
STATE = Path("/home/shin/deploy/knowslink-state")
SCRIPT = SNAP / "deploy/knowslink/access_apply.py"
FIXED = "28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2"
AUD = "a" * 64
OTHER = "b" * 64
APP_ID = "app-fixture"
POLICY_ID = "policy-fixture"
EVIDENCE = Path(__file__).resolve().parent
CHECKS = []


def redact(text):
    text = re.sub(r"[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}", "REDACTED_EMAIL", text)
    return re.sub(r"(?i)(bearer|token|secret|password|credential)\s*[:=]\s*\S+", r"\1=REDACTED", text)


def meta_digest(root):
    if not root.exists():
        return "absent"
    lines = []
    for path in sorted(root.rglob("*")):
        info = path.lstat()
        lines.append(f"{path.relative_to(root)} {stat_mode(info.st_mode)} {info.st_size} {info.st_mtime_ns}")
    return hashlib.sha256("\n".join(lines).encode()).hexdigest()


def stat_mode(mode):
    return oct(mode & 0o777)


def run(args, env=None):
    completed = subprocess.run(args, text=True, capture_output=True, env=env)
    return completed.returncode, redact(completed.stdout), redact(completed.stderr)


def record(name, ok, exit_code, detail):
    CHECKS.append({"name": name, "pass": bool(ok), "exit": exit_code, "detail": detail})
    print(f"{'PASS' if ok else 'FAIL'} {name} exit={exit_code} {detail}", flush=True)


def write_proof(root, *, live_aud=AUD, app_id=APP_ID, policy_id=POLICY_ID,
                saved_app=APP_ID, saved_policy=POLICY_ID, age_seconds=0, future_seconds=0,
                live=True, access_json=True):
    root.mkdir(mode=0o700)
    os.chmod(root, 0o700)
    email = root / "owner-email"
    email.write_text("owner@example.test\n")
    os.chmod(email, 0o600)
    if access_json:
        (root / "access.json").write_text(json.dumps({
            "app": saved_app, "policy": saved_policy, "policy_created": False}))
        os.chmod(root / "access.json", 0o600)
    (root / "access.aud").write_text(AUD + "\n")
    os.chmod(root / "access.aud", 0o600)
    tunnel = root / "tunnel"
    tunnel.mkdir(mode=0o700)
    (tunnel / "config.yml").write_text(
        "tunnel: fixture\naccess:\n  required: true\n  teamName: fixture-team\n  audTag:\n"
        f"    - {live_aud}\n")
    if not live:
        return
    app = {
        "id": app_id, "type": "self_hosted", "domain": "link.knowslog.com", "aud": live_aud,
        "destinations": [{"type": "public", "uri": "link.knowslog.com"}],
        "policies": [{"id": policy_id}], "allowed_idps": ["idp-fixture"],
        "options_preflight_bypass": False,
    }
    policy = {
        "id": policy_id, "decision": "allow",
        "include": [{"email": {"email": "owner@example.test"}}],
        "exclude": [], "require": [],
    }
    live_path = root / "access.live.json"
    live_path.write_text(json.dumps({"app": app, "policy": policy}))
    os.chmod(live_path, 0o600)
    if age_seconds or future_seconds:
        stamp = time.time() + future_seconds - age_seconds
        os.utime(live_path, (stamp, stamp))


def invoke(root, command, *, flags=(), optimize=None):
    env = os.environ.copy()
    env["KNOWSLINK_STATE_DIR"] = str(root)
    env["KNOWSLINK_OWNER_EMAIL_FILE"] = str(root / "owner-email")
    env["KNOWSLINK_TEAM"] = "fixture-team"
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    env.pop("CF_API_TOKEN_FILE", None)
    if optimize is None:
        env.pop("PYTHONOPTIMIZE", None)
    else:
        env["PYTHONOPTIMIZE"] = optimize
    return run(["python3", "-B", *flags, str(SCRIPT), command], env)


def expect(name, code, out, err, exit_code, needle):
    blob = out + err
    ok = code == exit_code and needle in blob and "REDACTED" not in blob
    if exit_code != 0:
        ok = ok and "agree" not in out
    record(name, ok, code, f"expect {exit_code} needle={needle!r}")


def main():
    before = {
        "snap": run(["git", "-C", str(SNAP), "rev-parse", "HEAD"])[1].strip(),
        "snap_porcelain": run(["git", "-C", str(SNAP), "status", "--porcelain"])[1],
        "live": run(["git", "-C", str(LIVE), "rev-parse", "HEAD"])[1].strip(),
        "live_porcelain": run(["git", "-C", str(LIVE), "status", "--porcelain"])[1],
        "script": hashlib.sha256(SCRIPT.read_bytes()).hexdigest(),
        "state": meta_digest(STATE),
    }
    record("snap-head", before["snap"] == FIXED and before["snap_porcelain"] == "", None, before["snap"])
    record("live-head", before["live"] == FIXED and before["live_porcelain"] == "", None, before["live"])

    with tempfile.TemporaryDirectory(prefix="knowslink-public-gate-") as temporary:
        root = Path(temporary)
        code, out, err = invoke(root / "unused", "selftest")
        expect("selftest", code, out, err, 0, "access_apply selftest ok")
        code, out, err = invoke(root / "unused-o", "selftest", flags=("-O",))
        expect("selftest-O", code, out, err, 0, "access_apply selftest ok")

        cases = [
            ("fresh", {}, 0, "live Access app, owner-only policy and tunnel config agree", ()),
            ("fresh-O", {}, 0, "live Access app, owner-only policy and tunnel config agree", ("-O",)),
            ("aud-mismatch", {"live_aud": OTHER}, 1, "app aud differs from the recorded aud", ()),
            ("aud-mismatch-O", {"live_aud": OTHER}, 1, "app aud differs from the recorded aud", ("-O",)),
            ("aud-mismatch-OO", {"live_aud": OTHER}, 1, "app aud differs from the recorded aud", ("-OO",)),
            ("aud-mismatch-env", {"live_aud": OTHER}, 1, "app aud differs from the recorded aud", (), "1"),
            ("proof-app-id", {"saved_app": "other-app"}, 1, "snapshot app/policy ids differ", ()),
            ("proof-app-id-env", {"saved_app": "other-app"}, 1, "snapshot app/policy ids differ", (), "1"),
            ("proof-policy-id", {"saved_policy": "other-policy"}, 1, "snapshot app/policy ids differ", ()),
            ("missing-live", {"live": False}, 1, "access.live.json", ()),
            ("missing-access-json", {"access_json": False}, 1, "access.json", ()),
            ("stale-proof", {"age_seconds": 660}, 1, "last 10 minutes", ()),
            ("stale-proof-env", {"age_seconds": 660}, 1, "last 10 minutes", (), "1"),
            ("future-proof", {"future_seconds": 300}, 1, "last 10 minutes", ()),
            ("future-proof-env", {"future_seconds": 300}, 1, "last 10 minutes", (), "1"),
        ]
        for item in cases:
            name, kwargs, exit_code, needle, flags = item[:5]
            optimize = item[5] if len(item) == 6 else None
            target = root / name
            write_proof(target, **kwargs)
            code, out, err = invoke(target, "check", flags=flags, optimize=optimize)
            expect(name, code, out, err, exit_code, needle)

    after = {
        "snap": run(["git", "-C", str(SNAP), "rev-parse", "HEAD"])[1].strip(),
        "snap_porcelain": run(["git", "-C", str(SNAP), "status", "--porcelain"])[1],
        "live": run(["git", "-C", str(LIVE), "rev-parse", "HEAD"])[1].strip(),
        "live_porcelain": run(["git", "-C", str(LIVE), "status", "--porcelain"])[1],
        "state": meta_digest(STATE),
        "pyc": list(SNAP.rglob("__pycache__")),
    }
    record("snap-unchanged", after["snap"] == FIXED and after["snap_porcelain"] == "" and not after["pyc"], None, after["snap"])
    record("live-unchanged", after["live"] == FIXED and after["live_porcelain"] == "" and after["state"] == before["state"], None, "head porcelain state-meta")
    (EVIDENCE / "gates.json").write_text(json.dumps({
        "sha": FIXED,
        "script_sha256": before["script"],
        "state_meta_unchanged": before["state"] == after["state"],
        "checks": CHECKS,
    }, indent=2) + "\n")
    failed = [item["name"] for item in CHECKS if not item["pass"]]
    print(f"failed={failed}", flush=True)
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
