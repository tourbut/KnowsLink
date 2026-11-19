#!/usr/bin/env python3
"""Observe the shared http:new budget on the fixed candidate.

One anonymous source sends requests that the per-IP cap rejects.
The script prints status counts and bucket lengths only.
"""

import json
import os
import signal
import socket
import subprocess
import time
import urllib.parse
from collections import Counter
from pathlib import Path
import http.client
import base64
import hashlib
import quopri
import re
import secrets

EXPECT = "59b66ada8b36802484cc6d7e22523257b50572cc"
QA = Path("/tmp/knowslink-public-identity-qa-59b66ad")
WORK = Path("/tmp/knowslink-pi-f1")
CONTAINER = "knowslink-pi-f1-pg"
IMAGE = "postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24"
RESULT = Path(
    "/home/shin/orca/workspaces/KnowsLink/fullops-tester/.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER-test/f1-result.json"
)
IP_A = "198.51.100.10"
IP_B = "203.0.113.10"
IP_MEMBER = "192.0.2.10"


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


def hash_token(value):
    return base64.urlsafe_b64encode(hashlib.sha256(value.encode()).digest()).decode().rstrip("=")


class Client:
    def __init__(self, port, source_ip):
        self.port = port
        self.source_ip = source_ip
        self.cookies = {}

    def call(self, method, path, form=None):
        connection = http.client.HTTPConnection("127.0.0.1", self.port, timeout=20)
        headers = {"Host": "localhost", "CF-Connecting-IP": self.source_ip}
        if self.cookies:
            headers["Cookie"] = "; ".join(f"{name}={value}" for name, value in self.cookies.items())
        body = None
        if form is not None:
            body = urllib.parse.urlencode(form)
            headers["Content-Type"] = "application/x-www-form-urlencoded"
        connection.request(method, path, body=body, headers=headers)
        response = connection.getresponse()
        response.read()
        for name, value in response.getheaders():
            if name.lower() != "set-cookie":
                continue
            pair, _, _attrs = value.partition(";")
            cookie_name, _, cookie_value = pair.partition("=")
            self.cookies[cookie_name] = cookie_value
        status = response.status
        location = response.getheader("Location") or ""
        connection.close()
        return status, location


def psql(sql):
    result = subprocess.run(
        ["docker", "exec", "-i", CONTAINER, "psql", "-U", "knowslink", "-d", "knowslink",
         "-v", "ON_ERROR_STOP=1", "-X", "-q", "-t", "-A", "-c", sql],
        check=False, capture_output=True)
    if result.returncode != 0:
        raise RuntimeError("psql failed")
    return result.stdout.decode().strip()


def bucket_len(key):
    quoted = key.replace("'", "''")
    raw = psql(
        "SELECT COALESCE(jsonb_array_length(data->'Rates'->'" + quoted + "'), 0) FROM relay_state;"
    )
    return int(raw or "0")


def post_invalid(client, count):
    counts = Counter()
    for _ in range(count):
        status, _location = client.call("POST", "/auth/start", {"email": "not-an-address"})
        counts[status] += 1
    return dict(counts)


def take_code(directory, email):
    deadline = time.time() + 5
    while time.time() < deadline:
        for path in sorted(directory.glob("*.eml")):
            raw = path.read_bytes()
            path.unlink()
            head, _, body = raw.partition(b"\r\n\r\n")
            text = quopri.decodestring(body or raw).decode("utf-8", "replace")
            header = head.decode("utf-8", "replace")
            if email not in header and email not in text:
                continue
            match = re.search(r"(\d{6})", text)
            if not match:
                raise RuntimeError("mail had no code")
            return match.group(1)
        time.sleep(0.05)
    raise RuntimeError("mail missing")


def main():
    WORK.mkdir(mode=0o700, parents=True, exist_ok=True)
    os.chmod(WORK, 0o700)
    password = secrets.token_hex(16)
    db_port = free_port()
    env_file = WORK / "pg.env"
    env_file.write_text(
        "POSTGRES_USER=knowslink\nPOSTGRES_DB=knowslink\nPOSTGRES_PASSWORD=" + password + "\n",
        encoding="utf-8",
    )
    os.chmod(env_file, 0o600)
    url = f"postgres://knowslink:{password}@127.0.0.1:{db_port}/knowslink?sslmode=disable"
    (WORK / "db.url").write_text(url, encoding="utf-8")
    os.chmod(WORK / "db.url", 0o600)
    sink = None
    relay = None
    created = False
    try:
        sha = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=QA, text=True).strip()
        porcelain = subprocess.check_output(["git", "status", "--porcelain"], cwd=QA, text=True)
        if sha != EXPECT or porcelain != "":
            raise RuntimeError("qa checkout is not the clean candidate")
        subprocess.run(
            ["docker", "rm", "-f", CONTAINER],
            check=False, capture_output=True)
        run = subprocess.run(
            ["docker", "run", "-d", "--name", CONTAINER, "--env-file", str(env_file),
             "-p", f"127.0.0.1:{db_port}:5432", IMAGE],
            check=False, capture_output=True, text=True)
        if run.returncode != 0:
            raise RuntimeError("postgres start failed")
        created = True
        deadline = time.time() + 40
        while time.time() < deadline:
            ready = subprocess.run(
                ["docker", "exec", CONTAINER, "pg_isready", "-U", "knowslink", "-d", "knowslink"],
                check=False, capture_output=True)
            try:
                with socket.create_connection(("127.0.0.1", db_port), 0.2):
                    host_open = True
            except OSError:
                host_open = False
            if ready.returncode == 0 and host_open:
                break
            time.sleep(0.3)
        else:
            raise RuntimeError("database timeout")
        environment = os.environ.copy()
        environment["DATABASE_URL"] = url
        environment["MIGRATIONS_DIR"] = str(QA / "db/migrations")
        migrated = None
        for _ in range(20):
            migrated = subprocess.run(
                [str(QA / "build/migrate")], cwd=QA, env=environment, check=False, capture_output=True)
            if migrated.returncode == 0:
                break
            time.sleep(0.5)
        if migrated.returncode != 0:
            note = (migrated.stdout + migrated.stderr).decode("utf-8", "replace").replace(password, "redacted")
            (Path("/tmp") / "knowslink-pi-f1-migrate.txt").write_text(note[-800:], encoding="utf-8")
            raise RuntimeError("migrate failed")
        sink_port = free_port()
        mail_dir = WORK / "mail"
        mail_dir.mkdir(mode=0o700)
        sink_log = WORK / "sink.log"
        sink = subprocess.Popen(
            ["python3", str(QA / "scripts/mail_sink.py"), str(mail_dir), "--listen", "127.0.0.1", "--port", str(sink_port)],
            cwd=QA, stdout=sink_log.open("ab"), stderr=subprocess.STDOUT, start_new_session=True)
        relay_port = free_port()
        relay_log = WORK / "relay.log"
        environment["RELAY_ADDR"] = f"127.0.0.1:{relay_port}"
        environment["KNOWSLINK_SMTP_URL"] = f"smtp://127.0.0.1:{sink_port}"
        environment["KNOWSLINK_MAIL_FROM"] = "noreply@example.test"
        environment["KNOWSLINK_CLIENT_IP_HEADER"] = "CF-Connecting-IP"
        environment.pop("KNOWSLINK_SYNTHETIC_SIGNUP", None)
        relay = subprocess.Popen(
            [str(QA / "build/relay")], cwd=QA, env=environment,
            stdout=relay_log.open("ab"), stderr=subprocess.STDOUT, start_new_session=True)
        deadline = time.time() + 10
        while time.time() < deadline:
            if relay.poll() is not None:
                raise RuntimeError("relay exited")
            try:
                status, _location = Client(relay_port, IP_B).call("GET", "/healthz")
                if status == 200:
                    break
            except OSError:
                time.sleep(0.05)
        else:
            raise RuntimeError("relay health timeout")
        psql("UPDATE relay_state SET data = '{}'::jsonb, clock = clock_timestamp();")
        first = post_invalid(Client(relay_port, IP_A), 30)
        rest = post_invalid(Client(relay_port, IP_A), 170)
        after_a = {
            "http_new": bucket_len("http:new"),
            "ip_a": bucket_len("http:ip:" + IP_A),
            "ip_b": bucket_len("http:ip:" + IP_B),
        }
        other = Client(relay_port, IP_B).call("POST", "/auth/start", {"email": "not-an-address"})[0]
        phase1 = {
            "ip_a_first_30": first,
            "ip_a_next_170": rest,
            "buckets_before_other": after_a,
            "ip_b_first_status": other,
            "http_new_after_other": bucket_len("http:new"),
            "ip_b_after_other": bucket_len("http:ip:" + IP_B),
        }
        psql("UPDATE relay_state SET data = '{}'::jsonb, clock = clock_timestamp();")
        member = Client(relay_port, IP_MEMBER)
        email = "f1-member@example.test"
        start_status, start_location = member.call("POST", "/auth/start", {"email": email})
        code = take_code(mail_dir, email) if start_status == 303 else ""
        verify_status = 0
        if code:
            verify_status, _location = member.call("POST", "/auth/verify", {"code": code})
        after_signup = bucket_len("http:new")
        need = 200 - after_signup
        allowed = min(30, need)
        rejected = need - allowed
        flood_ok = post_invalid(Client(relay_port, IP_A), allowed)
        flood_rejected = post_invalid(Client(relay_port, IP_A), rejected)
        before_home = bucket_len("http:new")
        home_status, _location = member.call("GET", "/home")
        fresh_status, _location = Client(relay_port, IP_B).call("POST", "/auth/start", {"email": "not-an-address"})
        logout_status, logout_location = member.call("POST", "/auth/logout")
        phase2 = {
            "signup_start": start_status,
            "signup_start_location": start_location,
            "signup_verify": verify_status,
            "http_new_after_signup": after_signup,
            "flood_allowed": flood_ok,
            "flood_rejected": flood_rejected,
            "http_new_before_home": before_home,
            "home_status": home_status,
            "fresh_ip_status": fresh_status,
            "logout_status": logout_status,
            "logout_location": logout_location,
            "cleanup_bucket": bucket_len("cleanup"),
            "member_bucket": bucket_len("http:member:" + "unused"),
        }
        state = psql("SELECT data::text FROM relay_state;")
        member_keys = [key for key in json.loads(state).get("Rates", {}) if key.startswith("http:member:")]
        phase2["member_rate_keys"] = len(member_keys)
        if member_keys:
            phase2["member_bucket"] = bucket_len(member_keys[0])
        log_text = relay_log.read_text(encoding="utf-8", errors="replace") if relay_log.exists() else ""
        leaked = bool(code and code in log_text)
        payload = {
            "candidate": EXPECT,
            "phase1": phase1,
            "phase2": phase2,
            "relay_log_contains_code": leaked,
            "qa_porcelain_empty": subprocess.check_output(["git", "status", "--porcelain"], cwd=QA, text=True) == "",
        }
        text = json.dumps(payload, ensure_ascii=True, indent=2)
        if "@" in text or "postgres://" in text or (code and code in text) or email in text:
            raise RuntimeError("result scrub")
        RESULT.write_text(text + "\n", encoding="utf-8")
        print("F1_OBSERVE_DONE", flush=True)
        return 0
    finally:
        if relay and relay.poll() is None:
            os.killpg(relay.pid, signal.SIGTERM)
            relay.wait(5)
        if sink and sink.poll() is None:
            os.killpg(sink.pid, signal.SIGTERM)
            sink.wait(5)
        if created:
            subprocess.run(["docker", "rm", "-f", CONTAINER], check=False, capture_output=True)
        if WORK.exists():
            subprocess.run(["rm", "-rf", str(WORK)], check=False)


if __name__ == "__main__":
    raise SystemExit(main())
