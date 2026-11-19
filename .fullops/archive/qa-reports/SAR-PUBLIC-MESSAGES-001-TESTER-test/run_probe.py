#!/usr/bin/env python3
"""Run the independent messages QA probe against a private Postgres. Do not print secrets."""

import os
import re
import secrets
import socket
import subprocess
import sys
from pathlib import Path

CLONE = Path("/tmp/knowslink-messages-qa-09c523d")
EVIDENCE = Path(__file__).resolve().parent
PROBE_SRC = EVIDENCE / "probe_test.go.src"
PROBE_DST = CLONE / "internal" / "relay" / "qa_messages_probe_test.go"
LOG_DIR = Path("/tmp/knowslink-messages-qa-logs")
SECRET = re.compile(r"postgres(?:ql)?://\S+|-----BEGIN [\s\S]*?-----END [^-]+-----")
CANDIDATE = "09c523da8a3407288d9f5d711e1834af12bc7808"


def redact(text):
    return SECRET.sub("[redacted]", text)


def run(command, cwd, env, log_name):
    result = subprocess.run(command, cwd=cwd, env=env, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    (LOG_DIR / log_name).write_text(redact(result.stdout), encoding="utf-8")
    print(f"{log_name} exit {result.returncode}", flush=True)
    return result.returncode


def main():
    LOG_DIR.mkdir(parents=True, exist_ok=True)
    before = subprocess.run(["docker", "ps", "-a", "--format", "{{.ID}} {{.Names}} {{.Status}}"], text=True, stdout=subprocess.PIPE, check=True).stdout
    (EVIDENCE / "docker-before.txt").write_text(before, encoding="utf-8")
    head = subprocess.run(["git", "rev-parse", "HEAD"], cwd=CLONE, text=True, stdout=subprocess.PIPE, check=True).stdout.strip()
    status = subprocess.run(["git", "status", "--porcelain"], cwd=CLONE, text=True, stdout=subprocess.PIPE, check=True).stdout
    (EVIDENCE / "clone-before.txt").write_text(f"{head}\n{status}", encoding="utf-8")
    if head != CANDIDATE or status.strip():
        print("clone is not the clean fixed candidate", flush=True)
        return 2
    password = secrets.token_urlsafe(24)
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        port = listener.getsockname()[1]
    project = "knowslink-messagesqa-" + secrets.token_hex(4)
    override = LOG_DIR / "postgres.yaml"
    override.write_text(
        "services:\n  postgres:\n    ports: [\"127.0.0.1:%d:5432\"]\n    networks: [database, ingress]\n" % port,
        encoding="utf-8",
    )
    env = dict(os.environ)
    env["POSTGRES_PASSWORD"] = password
    compose = [
        "docker", "compose", "--project-name", project, "--env-file", ".env.example",
        "-f", "compose.yaml", "-f", str(override),
    ]
    code = 1
    try:
        up = run(compose + ["up", "-d", "--wait", "postgres"], CLONE, env, "probe-postgres-up.log")
        if up != 0:
            code = up
        else:
            built = run(["go", "build", "-o", str(LOG_DIR / "migrate"), "./cmd/migrate"], CLONE, env, "probe-migrate-build.log")
            if built != 0:
                code = built
            else:
                migrate_env = dict(env)
                migrate_env["DATABASE_URL"] = f"postgres://knowslink:{password}@127.0.0.1:{port}/knowslink?sslmode=disable"
                migrate_env["MIGRATIONS_DIR"] = "db/migrations"
                migrated = run([str(LOG_DIR / "migrate")], CLONE, migrate_env, "probe-migrate.log")
                if migrated != 0:
                    code = migrated
                else:
                    text = PROBE_SRC.read_text(encoding="utf-8")
                    formatted = subprocess.run(["gofmt"], input=text, text=True, stdout=subprocess.PIPE)
                    if formatted.returncode != 0:
                        (LOG_DIR / "probe-gofmt.log").write_text(formatted.stderr, encoding="utf-8")
                        print("probe-gofmt.log exit 1", flush=True)
                        code = formatted.returncode
                    else:
                        PROBE_SRC.write_text(formatted.stdout, encoding="utf-8")
                        PROBE_DST.write_text(formatted.stdout, encoding="utf-8")
                        test_env = dict(migrate_env)
                        test_env["TEST_SYNTHETIC_DATABASE"] = "1"
                        test_env["TEST_DATABASE_URL"] = migrate_env["DATABASE_URL"]
                        test_env["QA_CLONE"] = str(CLONE)
                        test_env["QA_ROUNDTRIP"] = str(EVIDENCE / "qa-roundtrip.mjs")
                        test_env["QA_DEFECTS"] = str(EVIDENCE / "admission-defects.json")
                        code = run(
                            ["go", "test", "-tags=integration", "-race", "-count=1", "-timeout", "15m", "-v", "-run", "TestQAMessagesIndependent", "./internal/relay"],
                            CLONE, test_env, "probe-go.log",
                        )
                        (EVIDENCE / "probe-go.log").write_text((LOG_DIR / "probe-go.log").read_text(encoding="utf-8"), encoding="utf-8")
    finally:
        PROBE_DST.unlink(missing_ok=True)
        down_env = dict(os.environ)
        down_env["POSTGRES_PASSWORD"] = password
        down = subprocess.run(compose + ["down", "--volumes"], cwd=CLONE, env=down_env, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        (LOG_DIR / "probe-postgres-down.log").write_text(redact(down.stdout), encoding="utf-8")
        print(f"probe-postgres-down.log exit {down.returncode}", flush=True)
        after = subprocess.run(["docker", "ps", "-a", "--format", "{{.ID}} {{.Names}} {{.Status}}"], text=True, stdout=subprocess.PIPE, check=True).stdout
        (EVIDENCE / "docker-after.txt").write_text(after, encoding="utf-8")
        end_head = subprocess.run(["git", "rev-parse", "HEAD"], cwd=CLONE, text=True, stdout=subprocess.PIPE, check=True).stdout.strip()
        end_status = subprocess.run(["git", "status", "--porcelain"], cwd=CLONE, text=True, stdout=subprocess.PIPE, check=True).stdout
        (EVIDENCE / "clone-after.txt").write_text(f"{end_head}\n{end_status}", encoding="utf-8")
        if code == 0 and down.returncode != 0:
            code = down.returncode
    return code


if __name__ == "__main__":
    raise SystemExit(main())
