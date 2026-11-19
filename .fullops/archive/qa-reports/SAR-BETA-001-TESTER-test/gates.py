#!/usr/bin/env python3
"""Fail-closed checks for Access proof and deploy gates. Does not move the live checkout."""

import importlib.util
import json
import os
import re
import shutil
import subprocess
import tempfile
from pathlib import Path

LIVE = Path("/home/shin/deploy/knowslink")
SCRIPT = LIVE / "deploy/knowslink"
FIXED = "f824015314c66bcab42940cfe3db2edabb22e1dd"
AUD = "a" * 64
OTHER = "b" * 64
CHECKS = []


def redact(text):
    return re.sub(r"[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}", "REDACTED_EMAIL", text)


def check(name, ok, detail):
    CHECKS.append({"name": name, "pass": bool(ok), "detail": detail})
    print(f"{'PASS' if ok else 'FAIL'} {name}: {detail}", flush=True)


def run(args, env=None):
    completed = subprocess.run(args, text=True, capture_output=True, env=env)
    return completed.returncode, redact(completed.stdout), redact(completed.stderr)


def load_access(state):
    os.environ["KNOWSLINK_STATE_DIR"] = str(state)
    os.environ.pop("CF_API_TOKEN_FILE", None)
    spec = importlib.util.spec_from_file_location("beta_access_apply", SCRIPT / "access_apply.py")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def expose(state):
    env = os.environ.copy()
    env["KNOWSLINK_STATE_DIR"] = str(state)
    env["KNOWSLINK_DEPLOY"] = str(LIVE)
    env.pop("CF_API_TOKEN_FILE", None)
    return run([str(SCRIPT / "beta.sh"), "expose"], env)


def write_proof(state, *, aud=AUD, live_aud=AUD, team="scshin88", age_seconds=0, idps=None):
    tunnel = state / "tunnel"
    tunnel.mkdir(parents=True, exist_ok=True)
    (state / "access.json").write_text(json.dumps({"app": "app", "policy": "policy", "policy_created": False}))
    (state / "access.aud").write_text(aud)
    app = {
        "type": "self_hosted", "domain": "link.knowslog.com", "aud": live_aud,
        "destinations": [{"type": "public", "uri": "link.knowslog.com"}],
        "policies": [{"id": "policy"}], "allowed_idps": idps or ["idp"],
        "options_preflight_bypass": False,
    }
    policy = {"id": "policy", "decision": "allow", "include": [{"email": {"email": "owner@example.test"}}], "exclude": [], "require": []}
    (state / "access.live.json").write_text(json.dumps({"app": app, "policy": policy}))
    if age_seconds:
        stamp = (state / "access.live.json").stat().st_mtime - age_seconds
        os.utime(state / "access.live.json", (stamp, stamp))
    (tunnel / "config.yml").write_text(f"access:\n  required: true\n  teamName: {team}\n  audTag:\n    - {live_aud}\n")


def main():
    code, head, _err = run(["git", "-C", str(LIVE), "rev-parse", "HEAD"])
    before = head.strip()
    _code, relay_before, _err = run(["docker", "inspect", "-f", "{{.Id}}", "knowslink-relay-1"])
    _code, dig_before, _err = run(["dig", "+short", "link.knowslog.com", "A"])
    check("live-head-before", code == 0 and before == FIXED, before)

    code, out, err = run(["python3", str(SCRIPT / "access_apply.py"), "selftest"])
    check("access-selftest", code == 0 and "selftest ok" in out, f"exit {code}")

    with tempfile.TemporaryDirectory(prefix="knowslink-qa-gates-") as temporary:
        root = Path(temporary)
        missing = root / "missing"
        missing.mkdir()
        module = load_access(missing)
        try:
            module.check()
            check("missing-proof", False, "check returned")
        except FileNotFoundError as exc:
            check("missing-proof", "access.json" in redact(str(exc)), "access.json absent")

        stale = root / "stale"
        write_proof(stale, age_seconds=601)
        module = load_access(stale)
        try:
            module.check()
            check("stale-proof", False, "check returned")
        except SystemExit as exc:
            check("stale-proof", "older than 10 minutes" in str(exc), "snapshot older than 10 minutes")

        mismatch = root / "mismatch"
        write_proof(mismatch, live_aud=OTHER)
        module = load_access(mismatch)
        saved = json.loads((mismatch / "access.live.json").read_text())
        try:
            module.verify_live(saved["app"], saved["policy"], "owner@example.test", AUD, "scshin88", (mismatch / "tunnel" / "config.yml").read_text())
            check("mismatched-aud", False, "verify_live returned")
        except AssertionError as exc:
            check("mismatched-aud", "aud" in str(exc).lower(), redact(str(exc)))

        fresh = root / "fresh"
        write_proof(fresh)
        module = load_access(fresh)
        try:
            module.check()
            check("foreign-email-rejected", False, "check returned")
        except AssertionError as exc:
            check("foreign-email-rejected", "owner email" in str(exc), redact(str(exc)))
        app = json.loads((fresh / "access.live.json").read_text())["app"]
        policy = json.loads((fresh / "access.live.json").read_text())["policy"]
        config = (fresh / "tunnel" / "config.yml").read_text()
        try:
            module.verify_live(app | {"allowed_idps": ["one", "two"]}, policy, "owner@example.test", AUD, "scshin88", config)
            check("extra-idp", False, "verify_live returned")
        except AssertionError as exc:
            check("extra-idp", "identity provider" in str(exc), redact(str(exc)))

        code, out, err = expose(missing)
        check("expose-missing-aud", code == 1 and "Access app not recorded" in err, f"exit {code}")
        jwt = root / "jwt"
        write_proof(jwt)
        (jwt / "tunnel" / "config.yml").write_text("ingress: []\n")
        code, out, err = expose(jwt)
        check("expose-missing-jwt", code == 1 and "origin JWT check missing" in err, f"exit {code}")
        code, out, err = expose(stale)
        blocked = code == 1 and "live Access verification failed" in err and "tunnel route" not in err
        check("expose-stale-before-dns", blocked, f"exit {code}")

    snap = Path(tempfile.mkdtemp(prefix="knowslink-qa-snap-"))
    state = Path(tempfile.mkdtemp(prefix="knowslink-qa-state-"))
    try:
        code, out, err = run(["git", "clone", "--quiet", str(LIVE), str(snap)])
        check("snapshot-clone", code == 0, f"exit {code}")
        env = os.environ.copy()
        env["KNOWSLINK_DEPLOY"] = str(snap)
        env["KNOWSLINK_STATE_DIR"] = str(state)
        env.pop("CF_API_TOKEN_FILE", None)
        code, out, err = run([str(SCRIPT / "beta.sh"), "deploy", "0" * 40], env)
        check("deploy-unknown-sha", code == 1 and "unknown commit" in err, f"exit {code}")
        code, base, _err = run(["git", "-C", str(snap), "rev-parse", "HEAD"])
        migration = snap / "db/migrations/9999_qa_gate.sql"
        migration.write_text("-- temporary gate fixture\n")
        run(["git", "-C", str(snap), "add", "db/migrations/9999_qa_gate.sql"])
        code, out, err = run(["git", "-C", str(snap), "commit", "-q", "-m", "temporary migration gate"])
        check("snapshot-commit", code == 0, f"exit {code}")
        code, sha, _err = run(["git", "-C", str(snap), "rev-parse", "HEAD"])
        code, out, err = run(["git", "-C", str(snap), "checkout", "-q", "--detach", base.strip()])
        check("snapshot-reset", code == 0, "temporary checkout only")
        code, out, err = run([str(SCRIPT / "beta.sh"), "deploy", sha.strip()], env)
        check("deploy-migration-blocked", code == 1 and "db/migrations differ" in err, f"exit {code}")
        check("snapshot-no-dump", not any(state.rglob("*.dump")), "no dump")
    finally:
        shutil.rmtree(snap, ignore_errors=True)
        shutil.rmtree(state, ignore_errors=True)

    code, head, _err = run(["git", "-C", str(LIVE), "rev-parse", "HEAD"])
    _code, relay_after, _err = run(["docker", "inspect", "-f", "{{.Id}}", "knowslink-relay-1"])
    _code, dig_after, _err = run(["dig", "+short", "link.knowslog.com", "A"])
    _code, status, _err = run(["git", "-C", str(LIVE), "status", "--porcelain"])
    check("live-unchanged", head.strip() == before == FIXED and relay_before.strip() == relay_after.strip() and dig_before == dig_after and status.strip() == "", "head relay dns clean")
    evidence = Path(__file__).resolve().parent
    (evidence / "gates.json").write_text(json.dumps(CHECKS, indent=2) + "\n")
    failed = [item["name"] for item in CHECKS if not item["pass"]]
    print(f"failed={failed}", flush=True)
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
