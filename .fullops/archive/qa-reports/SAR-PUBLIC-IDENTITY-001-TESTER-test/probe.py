#!/usr/bin/env python3
"""Independent loopback QA for the fixed public-identity candidate.

Drives the real relay process, a temporary Postgres, and the local SMTP sink.
Prints check ids only. Does not print addresses, codes, cookies, or tokens.
"""

import base64
import hashlib
import hmac
import http.client
import json
import os
import quopri
import re
import signal
import socket
import subprocess
import threading
import time
import urllib.parse
from pathlib import Path

EXPECT = "59b66ada8b36802484cc6d7e22523257b50572cc"
QA = Path(os.environ.get("QA_ROOT", "/tmp/knowslink-public-identity-qa-59b66ad"))
CONTAINER = os.environ.get("PG_CONTAINER", "knowslink-pi-qa-pg")
WORK = Path("/tmp/knowslink-pi-qa-run")
RESULT = Path(os.environ.get(
    "QA_RESULT",
    "/home/shin/orca/workspaces/KnowsLink/fullops-tester/.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER-test/result.json",
))
DOMAIN = "example.test"


def address(label):
    return f"{label}@{DOMAIN}"


def hash_token(value):
    return base64.urlsafe_b64encode(hashlib.sha256(value.encode()).digest()).decode().rstrip("=")


def csrf(token, gate_id):
    mac = hmac.new(token.encode(), b"gate:" + gate_id.encode(), hashlib.sha256).digest()
    return base64.urlsafe_b64encode(mac).decode().rstrip("=")


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


class Report:
    def __init__(self):
        self.rows = []
        self.secrets = []

    def secret(self, value):
        if value:
            self.secrets.append(str(value))

    def check(self, check_id, ok, note="", http=None):
        self.rows.append({"id": check_id, "ok": bool(ok), "note": note, "http": http})
        print(f"{'PASS' if ok else 'FAIL'} {check_id} {note}", flush=True)
        return bool(ok)

    def skip(self, check_id, note):
        self.rows.append({"id": check_id, "ok": None, "skipped": True, "note": note})
        print(f"SKIP {check_id} {note}", flush=True)


REPORT = Report()


def db_url():
    return (WORK / "db.url").read_text().strip()


def psql(sql):
    result = subprocess.run(
        ["docker", "exec", "-i", CONTAINER, "psql", "-U", "knowslink", "-d", "knowslink",
         "-v", "ON_ERROR_STOP=1", "-X", "-q", "-t", "-A", "-c", sql],
        check=False, capture_output=True)
    if result.returncode != 0:
        (WORK / "psql.err").write_bytes(result.stderr)
        raise RuntimeError("psql failed")
    return result.stdout.decode()


def pg_time(interval="interval '0 seconds'"):
    return psql(
        "SELECT to_char(clock_timestamp() AT TIME ZONE 'UTC' - "
        + interval + ", 'YYYY-MM-DD\"T\"HH24:MI:SS.US\"Z\"')"
    ).strip()


def load():
    raw = psql("SELECT data::text FROM relay_state;").rstrip("\n")
    value = json.loads(raw or "{}")
    return value if isinstance(value, dict) else {}


def save(state):
    payload = json.dumps(state, separators=(",", ":"), ensure_ascii=True)
    if "$qa$" in payload:
        raise RuntimeError("state payload rejected")
    psql("UPDATE relay_state SET data = $qa$" + payload + "$qa$::jsonb, clock = clock_timestamp();")


def reset_state():
    psql("UPDATE relay_state SET data = '{}'::jsonb, clock = clock_timestamp();")


def member_id(email):
    return (load().get("Identities") or {}).get("knowslink-email-otp|" + email)


def age_space(email):
    state = load()
    key = "send:space:" + hash_token(email)
    rates = state.setdefault("Rates", {})
    if not rates.get(key):
        return
    rates[key] = [pg_time("interval '61 seconds'")]
    save(state)


def set_member_times(email, **fields):
    state = load()
    target = member_id(email)
    for session in (state.get("Sessions") or {}).values():
        if session.get("Member") == target:
            session.update(fields)
    save(state)


def count_challenges(email):
    return sum(1 for item in (load().get("Challenges") or {}).values() if item.get("Email") == email)


class Client:
    def __init__(self, port, source_ip=None):
        self.port = port
        self.source_ip = source_ip
        self.cookies = {}
        self.cookie_flags = {}

    def call(self, method, path, form=None, extra=None):
        connection = http.client.HTTPConnection("127.0.0.1", self.port, timeout=20)
        headers = {"Host": "localhost", "Cache-Control": "no-store"}
        if self.cookies:
            headers["Cookie"] = "; ".join(f"{name}={value}" for name, value in self.cookies.items())
        if self.source_ip:
            headers["CF-Connecting-IP"] = self.source_ip
        body = None
        if form is not None:
            body = urllib.parse.urlencode(form)
            headers["Content-Type"] = "application/x-www-form-urlencoded"
        if extra:
            headers.update(extra)
        connection.request(method, path, body=body, headers=headers)
        response = connection.getresponse()
        payload = response.read()
        flags = {}
        for name, value in response.getheaders():
            if name.lower() != "set-cookie":
                continue
            pair, _, attrs = value.partition(";")
            cookie_name, _, cookie_value = pair.partition("=")
            self.cookies[cookie_name] = cookie_value
            flags[cookie_name] = attrs
            REPORT.secret(cookie_value)
        self.cookie_flags.update(flags)
        location = response.getheader("Location") or ""
        text = payload.decode("utf-8", "replace")
        connection.close()
        wanted = ("cache-control", "content-security-policy", "referrer-policy", "x-content-type-options")
        found = {name.lower(): value for name, value in response.getheaders()}
        return response.status, location, text, {name: found.get(name, "") for name in wanted}


class Sink:
    def __init__(self):
        self.port = free_port()
        self.directory = WORK / "mail"
        self.proc = None
        self.log = WORK / "sink.log"

    def start(self):
        self.directory.mkdir(parents=True, exist_ok=True)
        self.proc = subprocess.Popen(
            ["python3", str(QA / "scripts/mail_sink.py"), str(self.directory), "--listen", "127.0.0.1", "--port", str(self.port)],
            cwd=QA, stdout=self.log.open("ab"), stderr=subprocess.STDOUT, start_new_session=True)
        deadline = time.time() + 5
        while time.time() < deadline:
            if self.proc.poll() is not None:
                raise RuntimeError("sink exited")
            try:
                with socket.create_connection(("127.0.0.1", self.port), 0.2):
                    return
            except OSError:
                time.sleep(0.05)
        raise RuntimeError("sink timeout")

    def stop(self):
        if self.proc and self.proc.poll() is None:
            os.killpg(self.proc.pid, signal.SIGTERM)
            self.proc.wait(5)

    def take_code(self, email):
        deadline = time.time() + 5
        while time.time() < deadline:
            for path in sorted(self.directory.glob("*.eml")):
                mode = path.stat().st_mode & 0o777
                raw = path.read_bytes()
                path.unlink()
                if mode != 0o600:
                    raise RuntimeError("mail mode")
                head, _, body = raw.partition(b"\r\n\r\n")
                text = quopri.decodestring(body or raw).decode("utf-8", "replace")
                header = head.decode("utf-8", "replace")
                if email not in header and email not in text:
                    continue
                match = re.search(r"(\d{6})", text)
                if not match:
                    raise RuntimeError("mail had no code")
                REPORT.secret(match.group(1))
                REPORT.secret(email)
                return match.group(1)
            time.sleep(0.05)
        raise RuntimeError("mail missing")


class Relay:
    def __init__(self, extra, log_name):
        self.port = free_port()
        self.extra = extra
        self.log_path = WORK / log_name
        self.proc = None

    def start(self):
        environment = os.environ.copy()
        environment["DATABASE_URL"] = db_url()
        environment["RELAY_ADDR"] = f"127.0.0.1:{self.port}"
        for name in ("KNOWSLINK_SYNTHETIC_SIGNUP", "KNOWSLINK_CLIENT_IP_HEADER", "KNOWSLINK_TEST_AGENTS", "KNOWSLINK_SMTP_URL", "KNOWSLINK_MAIL_FROM"):
            environment.pop(name, None)
        environment.update(self.extra)
        self.proc = subprocess.Popen(
            [str(QA / "build/relay")], cwd=QA, env=environment,
            stdout=self.log_path.open("ab"), stderr=subprocess.STDOUT, start_new_session=True)
        deadline = time.time() + 10
        while time.time() < deadline:
            if self.proc.poll() is not None:
                raise RuntimeError("relay exited before health")
            try:
                status, _, text, _ = Client(self.port).call("GET", "/healthz")
                if status == 200 and text.strip() == "ok":
                    return
            except OSError:
                pass
            time.sleep(0.05)
        raise RuntimeError("relay health timeout")

    def stop(self):
        if self.proc and self.proc.poll() is None:
            os.killpg(self.proc.pid, signal.SIGTERM)
            try:
                self.proc.wait(5)
            except subprocess.TimeoutExpired:
                os.killpg(self.proc.pid, signal.SIGKILL)
                self.proc.wait(3)


def wait_database():
    deadline = time.time() + 30
    while time.time() < deadline:
        result = subprocess.run(
            ["docker", "exec", CONTAINER, "pg_isready", "-U", "knowslink", "-d", "knowslink"],
            check=False, capture_output=True)
        if result.returncode == 0:
            return
        time.sleep(0.3)
    raise RuntimeError("database timeout")


def migrate():
    environment = os.environ.copy()
    environment["DATABASE_URL"] = db_url()
    environment["MIGRATIONS_DIR"] = str(QA / "db/migrations")
    result = subprocess.run([str(QA / "build/migrate")], cwd=QA, env=environment, check=False, capture_output=True)
    (WORK / "migrate.log").write_bytes(result.stdout + result.stderr)
    if result.returncode != 0:
        raise RuntimeError("migrate failed")


def phrases(check_id, text, required):
    missing = [item for item in required if item not in text]
    return REPORT.check(check_id, not missing, note="" if not missing else "missing " + " | ".join(missing))


def login(client, sink, email, check_id):
    status, location, _, _ = client.call("POST", "/auth/start", {"email": email})
    if not REPORT.check(check_id + "-start", status == 303 and location == "/auth/verify", http=status):
        return ""
    code = sink.take_code(email.lower())
    pending = client.cookies.get("__Host-kl_pending", "")
    state = load()
    bound = (state.get("Challenges") or {}).get(hash_token(pending), {})
    REPORT.check(
        check_id + "-binding",
        bound.get("Code") == hash_token(pending + ":" + code) and bound.get("Email") == email.lower(),
        note="stored challenge matches the sink message",
    )
    status, location, _, _ = client.call("POST", "/auth/verify", {"code": code})
    if not REPORT.check(check_id + "-verify", status == 303 and location == "/home", http=status):
        return ""
    status, _, body, _ = client.call("GET", "/home")
    found = re.search(r"mem_[A-Za-z0-9_-]+", body)
    ok = status == 200 and found and member_id(email.lower()) == found.group(0)
    REPORT.check(check_id + "-home", ok, http=status)
    return found.group(0) if found else ""


def phase_a(sink):
    reset_state()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{sink.port}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
    }, "relay-a.log")
    relay.start()
    try:
        client = Client(relay.port)
        status, _, body, headers = client.call("GET", "/")
        REPORT.check("p01-start", status == 200, http=status)
        phrases("p07-ux01-start", body, [
            "이메일로 가입·로그인",
            "관리자 아이디, 서버 접속, 별도 초대는 필요하지 않습니다",
            "코드를 확인하기 전에는 로그인되지 않습니다",
            "확인 코드 받기",
        ])
        REPORT.check(
            "p01-security-headers",
            headers["cache-control"] == "no-store" and "default-src 'none'" in headers["content-security-policy"]
            and headers["referrer-policy"] == "no-referrer" and headers["x-content-type-options"] == "nosniff",
        )
        status, location, _, _ = client.call("GET", "/home")
        REPORT.check("p03-anonymous-home", status == 303 and location.startswith("/?n=expired"), http=status)
        status, _, body, _ = client.call("POST", "/auth/start", {"email": "Name <x@y.co>"})
        REPORT.check("p01-bad-format", status == 422 and "형식" in body, http=status)
        status, _, body, _ = client.call("POST", "/v1/owners", {})
        REPORT.check("p04-synthetic-closed", status == 403 and "sender_not_allowed" in body and "credential" not in body, http=status)
        status, _, body, _ = client.call("POST", "/v1/agents", {})
        REPORT.check("p04-sessionless-agent-api", status == 401 and "invalid_auth" in body, http=status)

        email = address("fixture-a")
        status, location, _, _ = client.call("POST", "/auth/start", {"email": email})
        before = len(load().get("Members") or {})
        code = sink.take_code(email) if status == 303 else ""
        REPORT.check("p01-no-member-before-verify", status == 303 and before == 0 and count_challenges(email) == 1 and code != "", http=status)
        status, _, body, _ = client.call("GET", "/auth/verify")
        phrases("p07-ux02-verify", body, [
            "확인 대기",
            "메일 발송은 신원 확인이 아닙니다",
            "아직 로그인되지 않았습니다",
            "5번 틀리면 새 코드를 받아야 합니다",
            "다시 받기는 이전 요청 60초 뒤부터 가능합니다",
            "새 코드를 받으면 이전 코드는 쓸 수 없습니다",
            "다른 주소로 시작",
        ])
        REPORT.check("p01-verify-mask", (email[0] + "***@") in body)
        wrong = "000000" if code != "000000" else "111111"
        status, _, body, _ = client.call("POST", "/auth/verify", {"code": wrong})
        REPORT.check("p01-wrong-code", status == 401 and "남은 시도: 4회" in body and member_id(email) is None, http=status)
        status, location, home, _ = client.call("POST", "/auth/verify", {"code": code})
        status, _, home, _ = client.call("GET", "/home")
        found = re.search(r"mem_[A-Za-z0-9_-]+", home)
        first = found.group(0) if found else ""
        REPORT.check("p01-verified-home", status == 200 and first.startswith("mem_") and member_id(email) == first, http=status)
        phrases("p07-ux03-home", home, [
            "확인된 이메일",
            "회원 식별자",
            "60분 동안 활동이 없으면 먼저 끝납니다",
            "아직 연결한 agent가 없습니다",
            "상태: 준비 중",
            "지금은 연결할 수 없습니다",
            "로그아웃과 전체 로그아웃은 브라우저 세션만 끝냅니다",
            "별도로 연결한 agent의 키와 자격은 철회하지 않습니다",
            "계정 비활성화",
            "아직 요청할 수 없습니다",
            'action="/auth/logout">',
            'action="/auth/logout-all"',
        ])
        REPORT.check("p01-home-hides-address", email not in home and "fixture-a" not in home and "/v1/owners" not in home and "Basic" not in home)
        flags = client.cookie_flags.get("__Host-kl_session", "")
        lowered = flags.lower()
        REPORT.check(
            "p03-cookie-flags",
            "secure" in lowered and "httponly" in lowered and "samesite=strict" in lowered and "path=/" in lowered and "domain=" not in lowered,
        )
        status, _, body, _ = client.call("POST", "/auth/verify", {"code": code})
        REPORT.check("p01-code-reuse", status == 401 and "만료" in body, http=status)
        status, _, body, _ = client.call("GET", "/owner")
        REPORT.check("p04-owner-screen-without-session", status == 401, http=status)
        owner_client = Client(relay.port)
        owner_client.cookies = dict(client.cookies)
        status, _, body, _ = owner_client.call("GET", "/owner")
        REPORT.check("p04-member-cookie-is-not-admin", status == 401 and "KnowsLink Owner 작업 화면" not in body, http=status)
        status, _, body, _ = owner_client.call("GET", "/v1/contacts", extra={"Authorization": "Bearer " + client.cookies.get("__Host-kl_session", "")})
        REPORT.check("p04-session-cookie-is-not-agent", status == 401, http=status)

        state = load()
        state.setdefault("Identities", {})["other-issuer|" + email] = "mem_otherissuerfixture01"
        save(state)
        age_space(email)
        other = Client(relay.port)
        again = login(other, sink, "Fixture-A@Example.Test", "p02-relogin")
        prefixes = sorted({key.split("|", 1)[0] for key in (load().get("Identities") or {}) if key.endswith("|" + email)})
        REPORT.check("p02-same-member", again == first and again != "")
        REPORT.check("p02-other-issuer-not-merged", prefixes == ["knowslink-email-otp", "other-issuer"] and member_id(email) == first)

        held = dict(other.cookies)
        relay.stop()
        relay.start()
        restarted = Client(relay.port)
        restarted.cookies = dict(held)
        status, _, body, _ = restarted.call("GET", "/home")
        REPORT.check("p03-restart-keeps-session", status == 200 and first in body, http=status)

        owner = load()["Members"][first]["Owner"]
        agent_secret = "fixture-agent-secret-value"
        REPORT.secret(agent_secret)
        state = load()
        state.setdefault("Agents", {})["agent_fixture"] = {
            "Owner": owner, "Credential": hash_token(agent_secret),
            "Keys": {"k1": {"Public": base64.b64encode(b"\x11" * 32).decode(), "Revoked": False, "Changed": pg_time()}},
        }
        save(state)
        status, _, body, _ = Client(relay.port).call("GET", "/v1/contacts", extra={"Authorization": "Bearer " + agent_secret})
        REPORT.check("p03-agent-before-logout", status == 200 and body.strip() == "[]", http=status)
        status, location, _, _ = restarted.call("POST", "/auth/logout", {})
        REPORT.check("p03-logout", status == 303 and location == "/?n=logout", http=status)
        status, _, body, _ = Client(relay.port).call("GET", "/?n=logout")
        REPORT.check("p03-logout-notice", status == 200 and "연결한 agent의 키와 자격은 그대로입니다." in body, http=status)
        replay = Client(relay.port)
        replay.cookies = dict(held)
        status, location, _, _ = replay.call("GET", "/home")
        REPORT.check("p03-old-session-rejected", status == 303 and location.startswith("/?n=expired"), http=status)
        status, _, body, _ = Client(relay.port).call("GET", "/v1/contacts", extra={"Authorization": "Bearer " + agent_secret})
        REPORT.check("p03-agent-survives-logout", status == 200 and body.strip() == "[]", http=status)

        age_space(email)
        left = Client(relay.port)
        right = Client(relay.port)
        left_id = login(left, sink, email, "p03-second-browser")
        REPORT.check("p03-second-browser-same-member", left_id == first)
        status, location, _, _ = left.call("POST", "/auth/logout-all", {})
        REPORT.check("p03-logout-all-fresh", status == 303 and location == "/?n=all", http=status)
        status, location, _, _ = restarted.call("GET", "/home")
        # restarted was already logged out; the live second browser is `right` only after its own login.
        status, location, _, _ = right.call("GET", "/home")
        REPORT.check("p03-right-still-anonymous", status == 303, http=status)
        age_space(email)
        right_id = login(right, sink, email, "p03-live-peer")
        REPORT.check("p03-peer-same-member", right_id == first)
        left = Client(relay.port)
        age_space(email)
        login(left, sink, email, "p03-live-current")
        status, location, _, _ = left.call("POST", "/auth/logout-all", {})
        status_peer, location_peer, _, _ = right.call("GET", "/home")
        REPORT.check("p03-logout-all-revokes-peer", status == 303 and status_peer == 303 and location_peer.startswith("/?n=expired"), http=status)

        email_b = address("fixture-b")
        stale = Client(relay.port)
        login(stale, sink, email_b, "p03-stale-login")
        set_member_times(email_b, Verified=pg_time("interval '6 minutes'"))
        status, _, body, _ = stale.call("GET", "/home")
        REPORT.check(
            "p03-reauth-required",
            status == 200 and "이메일 다시 확인" in body and 'action="/auth/logout-all"' not in body and 'action="/auth/reauth"' in body,
            http=status,
        )
        status, _, body, _ = stale.call("POST", "/auth/logout-all", {})
        REPORT.check("p03-logout-all-stale", status == 403 and "5분 안에" in body, http=status)
        age_space(email_b)
        status, location, _, _ = stale.call("POST", "/auth/reauth", {})
        reauth_ok = REPORT.check("p03-reauth-start", status == 303 and location == "/auth/verify", http=status)
        if reauth_ok:
            code = sink.take_code(email_b)
            status, location, _, _ = stale.call("POST", "/auth/verify", {"code": code})
            status, _, body, _ = stale.call("GET", "/home")
            REPORT.check("p03-reauth-restores-logout-all", status == 200 and 'action="/auth/logout-all"' in body, http=status)
        else:
            REPORT.check("p03-reauth-restores-logout-all", False, note="reauth did not send")
        peer = Client(relay.port)
        age_space(email_b)
        login(peer, sink, email_b, "p03-reauth-peer")
        status, _, _, _ = stale.call("POST", "/auth/logout-all", {})
        status_peer, location_peer, _, _ = peer.call("GET", "/home")
        REPORT.check("p03-logout-all-after-reauth", status == 303 and status_peer == 303 and location_peer.startswith("/?n=expired"))

        email_c = address("fixture-c")
        idle = Client(relay.port)
        login(idle, sink, email_c, "p03-idle-over-login")
        set_member_times(email_c, Seen=pg_time("interval '61 minutes'"))
        status, location, _, _ = idle.call("GET", "/home")
        REPORT.check("p03-idle-over", status == 303 and location.startswith("/?n=expired"), http=status)
        status, _, _, _ = Client(relay.port).call("GET", "/v1/contacts", extra={"Authorization": "Bearer " + agent_secret})
        REPORT.check("p03-agent-survives-idle", status == 200, http=status)

    finally:
        relay.stop()


def phase_b(sink):
    reset_state()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{sink.port}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
        "KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP",
    }, "relay-b.log")
    relay.start()
    try:
        email = address("fixture-race")
        left, right = Client(relay.port, "192.0.2.10"), Client(relay.port, "192.0.2.11")
        results = []

        def start(client):
            results.append(client.call("POST", "/auth/start", {"email": email})[:2])

        threads = [threading.Thread(target=start, args=(item,)) for item in (left, right)]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()
        wins = [item for item in results if item == (303, "/auth/verify")]
        blocked = [item for item in results if item[0] == 429]
        REPORT.check("p02-concurrent-start", len(wins) == 1 and len(blocked) == 1 and len(results) == 2, http=wins[0][0] if wins else None)
        winner = left if left.cookies.get("__Host-kl_pending") and (left.call("GET", "/auth/verify")[0] == 200) else right
        # The probe above consumed the verify page. Read the code and finish one login.
        code = sink.take_code(email)
        verified, verified_location, _, _ = winner.call("POST", "/auth/verify", {"code": code})
        status, _, body, _ = winner.call("GET", "/home")
        REPORT.check(
            "p02-concurrent-start-one-member",
            verified == 303 and verified_location == "/home" and status == 200 and len(load().get("Members") or {}) == 1 and "mem_" in body,
            http=verified,
        )

        reset_state()
        email = address("fixture-once")
        state = load()
        challenges = {}
        clients = []
        for index in range(2):
            pending = f"pending-race-{index}-token-value"
            code = f"{100000 + index}"
            REPORT.secret(pending)
            REPORT.secret(code)
            challenges[hash_token(pending)] = {
                "Email": email, "Code": hash_token(pending + ":" + code), "Exp": pg_time("interval '-10 minutes'"), "Attempts": 0,
            }
            client = Client(relay.port, f"192.0.2.{20 + index}")
            client.cookies["__Host-kl_pending"] = pending
            clients.append((client, code))
        state["Challenges"] = challenges
        save(state)
        statuses = []

        def verify(client, code):
            statuses.append(client.call("POST", "/auth/verify", {"code": code})[0])

        threads = [threading.Thread(target=verify, args=item) for item in clients]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()
        state = load()
        REPORT.check(
            "p02-concurrent-same-email",
            statuses.count(303) == 2 and len(state.get("Members") or {}) == 1 and len(state.get("Identities") or {}) == 1 and len(state.get("Owners") or {}) == 1,
            note=f"statuses {sorted(statuses)}",
        )

        reset_state()
        email = address("fixture-single")
        pending = "pending-single-use-token-value"
        code = "654321"
        REPORT.secret(pending)
        REPORT.secret(code)
        state = load()
        state["Challenges"] = {hash_token(pending): {
            "Email": email, "Code": hash_token(pending + ":" + code), "Exp": pg_time("interval '-10 minutes'"), "Attempts": 0,
        }}
        save(state)
        statuses = []
        workers = []
        for index in range(6):
            client = Client(relay.port, f"198.51.100.{index + 1}")
            client.cookies["__Host-kl_pending"] = pending
            workers.append(client)

        def verify_same(client):
            statuses.append(client.call("POST", "/auth/verify", {"code": code})[0])

        threads = [threading.Thread(target=verify_same, args=(item,)) for item in workers]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()
        REPORT.check("p02-one-use-parallel", statuses.count(303) == 1 and len(load().get("Members") or {}) == 1, note=f"wins {statuses.count(303)}")

        reset_state()
        alice_email, bob_email = address("fixture-alice"), address("fixture-bob")
        alice, bob = Client(relay.port, "203.0.113.10"), Client(relay.port, "203.0.113.11")
        alice_id = login(alice, sink, alice_email, "p04-alice")
        bob_id = login(bob, sink, bob_email, "p04-bob")
        gate_id = "gatefixtureb0001"
        state = load()
        state.setdefault("Gates", {})[gate_id] = {
            "ID": gate_id, "Owner": state["Members"][bob_id]["Owner"], "Parent": "parentfixture",
            "Digest": "digest", "From": "agent_a", "To": "agent_b", "Policy": "fixture",
            "State": "pending", "Generation": 1, "Exp": pg_time("interval '-10 minutes'"), "Consumed": False,
        }
        save(state)
        status, _, body, _ = alice.call("GET", "/home")
        REPORT.check("p04-home-hides-other-member", status == 200 and bob_id not in body and bob_email not in body, http=status)
        status, _, body, _ = alice.call("GET", f"/home/gates/{gate_id}")
        REPORT.check("p04-other-gate-read", status == 403 and "sender_not_allowed" in body, http=status)
        status, _, body, _ = alice.call("POST", f"/home/gates/{gate_id}", {
            "csrf": csrf(alice.cookies["__Host-kl_session"], gate_id), "decision": "approve",
        })
        REPORT.check("p04-other-gate-decide", status == 409 and "invalid_gate" in body, http=status)
        status, _, body, _ = bob.call("GET", "/home")
        REPORT.check("p04-own-gate-link", status == 200 and f"/home/gates/{gate_id}" in body and alice_id not in body, http=status)
        status, _, body, _ = bob.call("GET", f"/home/gates/{gate_id}")
        REPORT.check("p04-own-gate-without-body", status == 200 and "원문 부재: 승인 불가" in body and "승인·거절 버튼 비활성" in body, http=status)
        status, _, body, _ = bob.call("POST", f"/home/gates/{gate_id}", {
            "csrf": csrf(bob.cookies["__Host-kl_session"], gate_id), "decision": "approve",
        }, extra={"Sec-Fetch-Site": "cross-site"})
        REPORT.check("p04-cross-site-decision", status == 403 and "cross-origin" in body, http=status)
        REPORT.check("p04-cross-site-did-not-approve", load()["Gates"][gate_id]["State"] != "approved")
        status, _, body, _ = bob.call("GET", "/owner/gates/" + gate_id, extra={"Authorization": "Bearer not-an-owner-token"})
        REPORT.check("p04-agent-bearer-not-owner-page", status == 401, http=status)

        base, dotted, plus = address("fixturea"), address("fixture.a"), address("fixture-a+tag")
        base_id = login(Client(relay.port, "192.0.2.30"), sink, base, "p02-dot-base")
        dot_id = login(Client(relay.port, "192.0.2.31"), sink, dotted, "p02-dot-alias")
        plus_id = login(Client(relay.port, "192.0.2.32"), sink, plus, "p02-plus-alias")
        REPORT.check("p02-dot-and-plus-distinct", len({base_id, dot_id, plus_id}) == 3 and "" not in {base_id, dot_id, plus_id})

        wrong_email = address("fixture-wrong")
        wrong_client = Client(relay.port, "192.0.2.40")
        status, _, _, _ = wrong_client.call("POST", "/auth/start", {"email": wrong_email})
        real = sink.take_code(wrong_email)
        alternate = "000000" if real != "000000" else "111111"
        seen = []
        for remaining in (4, 3, 2, 1):
            status, _, body, _ = wrong_client.call("POST", "/auth/verify", {"code": alternate})
            seen.append(status == 401 and f"남은 시도: {remaining}회" in body)
        status, _, body, _ = wrong_client.call("POST", "/auth/verify", {"code": alternate})
        fifth = status == 401 and "만료" in body and "남은 시도" not in body
        status, _, body, _ = wrong_client.call("POST", "/auth/verify", {"code": real})
        REPORT.check("p01-five-wrong-then-expired", all(seen) and fifth and status == 401 and member_id(wrong_email) is None, http=status)

        expired_email = address("fixture-expired")
        expired = Client(relay.port, "192.0.2.41")
        expired.call("POST", "/auth/start", {"email": expired_email})
        code = sink.take_code(expired_email)
        state = load()
        for item in state.get("Challenges", {}).values():
            if item.get("Email") == expired_email:
                item["Exp"] = pg_time("interval '1 minute'")
        save(state)
        status, _, body, _ = expired.call("POST", "/auth/verify", {"code": code})
        REPORT.check("p01-expired-code", status == 401 and "만료" in body and member_id(expired_email) is None, http=status)

        before_files = list(sink.directory.glob("*.eml"))
        status, _, body, _ = Client(relay.port, "192.0.2.42").call(
            "POST", "/auth/start", {"email": address("fixture-cross")}, extra={"Sec-Fetch-Site": "cross-site"})
        status_origin, _, body_origin, _ = Client(relay.port, "192.0.2.43").call(
            "POST", "/auth/start", {"email": address("fixture-origin")}, extra={"Origin": "https://evil.example"})
        after_files = list(sink.directory.glob("*.eml"))
        REPORT.check("p01-cross-site-start", status == 403 and "cross-origin" in body and len(after_files) == len(before_files), http=status)
        REPORT.check("p01-mismatched-origin", status_origin == 403 and "cross-origin" in body_origin, http=status_origin)

        bounds = (
            ("fixture-idle-in", "192.0.2.90", {"Seen": pg_time("interval '59 minutes'")}, 200, "p03-idle-inside"),
            ("fixture-abs-over", "192.0.2.91", {"Created": pg_time("interval '12 hours 1 minute'"), "Seen": pg_time()}, 303, "p03-absolute-over"),
            ("fixture-abs-edge", "192.0.2.92", {"Created": pg_time("interval '12 hours'"), "Seen": pg_time()}, 303, "p03-absolute-boundary"),
            ("fixture-abs-in", "192.0.2.93", {"Created": pg_time("interval '11 hours 58 minutes'"), "Seen": pg_time()}, 200, "p03-absolute-inside"),
        )
        for label, source, fields, expect, check_id in bounds:
            email = address(label)
            client = Client(relay.port, source)
            found = login(client, sink, email, check_id + "-login")
            set_member_times(email, **fields)
            status, location, body, _ = client.call("GET", "/home")
            if expect == 200:
                REPORT.check(check_id, status == 200 and found != "" and found in body, http=status)
            else:
                REPORT.check(check_id, status == 303 and location.startswith("/?n=expired"), http=status)
    finally:
        relay.stop()


def phase_rates(sink):
    reset_state()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{sink.port}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
        "KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP",
    }, "relay-rate.log")
    relay.start()
    try:
        email = address("fixture-space")
        client = Client(relay.port, "192.0.2.50")
        status, _, _, _ = client.call("POST", "/auth/start", {"email": email})
        sink.take_code(email)
        status_now, _, body, _ = client.call("POST", "/auth/start", {"email": email})
        REPORT.check("p05-resend-immediate", status == 303 and status_now == 429 and "이후 다시 시도" in body, http=status_now)
        print("WAIT resend window 61s", flush=True)
        time.sleep(61)
        status_later, location, _, _ = client.call("POST", "/auth/start", {"email": email})
        REPORT.check("p05-resend-after-61s", status_later == 303 and location == "/auth/verify", http=status_later)
        sink.take_code(email)

        limited = address("fixture-hour")
        state = load()
        key = "send:email:" + hash_token(limited)
        state.setdefault("Rates", {})[key] = [pg_time(f"interval '{10 + index} minutes'") for index in range(5)]
        save(state)
        status, _, body, _ = Client(relay.port, "192.0.2.51").call("POST", "/auth/start", {"email": limited})
        remained = len((load().get("Rates") or {}).get(key, []))
        REPORT.check("p05-email-hour-cap", status == 429 and remained == 5 and "이후 다시 시도" in body, http=status)
        relay.stop()
        relay.start()
        status, _, _, _ = Client(relay.port, "192.0.2.51").call("POST", "/auth/start", {"email": limited})
        REPORT.check("p05-email-cap-survives-restart", status == 429, http=status)

        ip_key = "send:ip:192.0.2.60"
        state = load()
        state.setdefault("Rates", {})[ip_key] = [pg_time("interval '10 minutes'") for _ in range(20)]
        save(state)
        denied, _, _, _ = Client(relay.port, "192.0.2.60").call("POST", "/auth/start", {"email": address("fixture-ip60")})
        allowed, location, _, _ = Client(relay.port, "192.0.2.61").call("POST", "/auth/start", {"email": address("fixture-ip61")})
        REPORT.check("p05-ip-hour-cap", denied == 429 and allowed == 303 and location == "/auth/verify", http=denied)
        sink.take_code(address("fixture-ip61"))

        flood = Client(relay.port, "198.51.100.80")
        codes = [flood.call("POST", "/auth/verify", {"code": "000000"})[0] for _ in range(30)]
        blocked, _, body, _ = flood.call("POST", "/auth/verify", {"code": "000000"})
        other, location, _, _ = Client(relay.port, "198.51.100.81").call("POST", "/auth/start", {"email": address("fixture-otherip")})
        REPORT.check("p05-anonymous-30", codes.count(401) == 30 and blocked == 429 and "이후 다시 시도" in body, http=blocked)
        REPORT.check("p05-anonymous-isolates-source", other == 303 and location == "/auth/verify", http=other)
        sink.take_code(address("fixture-otherip"))

        person = Client(relay.port, "203.0.113.80")
        login(person, sink, address("fixture-principal"), "p05-principal")
        homes = [person.call("GET", "/home")[0] for _ in range(39)]
        blocked, _, body, _ = person.call("GET", "/home")
        logout, location, _, _ = person.call("POST", "/auth/logout", {})
        REPORT.check("p05-principal-40", homes.count(200) == 39 and blocked == 429 and "이후 다시 시도" in body, http=blocked)
        REPORT.check("p05-cleanup-separate", logout == 303 and location == "/?n=logout", http=logout)

        starter = Client(relay.port, "203.0.113.81")
        login(starter, sink, address("fixture-cleanup"), "p05-global-login")
        state = load()
        state.setdefault("Rates", {})["http:new"] = [pg_time("interval '10 seconds'") for _ in range(200)]
        save(state)
        denied, _, _, _ = Client(relay.port, "203.0.113.82").call("POST", "/auth/start", {"email": address("fixture-global-http")})
        logout, location, _, _ = starter.call("POST", "/auth/logout", {})
        REPORT.check("p05-global-http-200", denied == 429 and logout == 303 and location == "/?n=logout", http=denied)

        state = load()
        state.setdefault("Rates", {})["send"] = [pg_time("interval '5 minutes'") for _ in range(100)]
        save(state)
        status, _, _, _ = Client(relay.port, "203.0.113.83").call("POST", "/auth/start", {"email": address("fixture-global-send")})
        REPORT.check("p05-global-send-100", status == 429, http=status)
    finally:
        relay.stop()


def phase_capacity(sink):
    reset_state()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{sink.port}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
        "KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP",
    }, "relay-capacity.log")
    relay.start()
    try:
        state = {"Owners": {}, "Members": {}, "Identities": {}, "Agents": {}, "Pairs": {}, "Messages": {}, "Idempotency": {}, "Gates": {}, "Sessions": {}, "Challenges": {}, "Rates": {}}
        existing = address("fixture-existing")
        existing_id = ""
        created = pg_time()
        for index in range(100):
            member = f"mem_seed{index:03d}xxxxxxxxxxxxxx"
            owner = f"owner_seed{index:03d}"
            email = address(f"seed{index:03d}")
            if index == 0:
                email = existing
                existing_id = member
            state["Owners"][owner] = {"Credential": "", "Active": True}
            state["Members"][member] = {"Owner": owner, "Email": email, "Issuer": "knowslink-email-otp", "Active": True, "Created": created}
            state["Identities"]["knowslink-email-otp|" + email] = member
        save(state)
        client = Client(relay.port, "192.0.2.70")
        found = login(client, sink, existing, "p05-existing-at-capacity")
        fresh = Client(relay.port, "192.0.2.71")
        status, _, _, _ = fresh.call("POST", "/auth/start", {"email": address("fixture-newcomer")})
        code = sink.take_code(address("fixture-newcomer"))
        status, _, body, _ = fresh.call("POST", "/auth/verify", {"code": code})
        REPORT.check(
            "p05-member-cap-100",
            found == existing_id and status == 503 and "신규 가입을 받을 수 없습니다" in body and len(load().get("Members") or {}) == 100,
            http=status,
        )
    finally:
        relay.stop()


def phase_header_and_smtp(sink):
    reset_state()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{sink.port}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
    }, "relay-default.log")
    relay.start()
    try:
        flood = Client(relay.port)
        codes = [flood.call("POST", "/auth/verify", {"code": "000000"})[0] for _ in range(30)]
        blocked, _, _, _ = flood.call("POST", "/auth/verify", {"code": "000000"})
        spoofed = Client(relay.port, "198.51.100.90")
        # Default mode ignores the header helper unless the process is configured. This client always sends the header.
        status, _, _, _ = spoofed.call("POST", "/auth/start", {"email": address("fixture-spoof")})
        REPORT.check("p05-default-ip-30", codes.count(401) == 30 and blocked == 429 and status == 429, http=status)
    finally:
        relay.stop()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{sink.port}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
        "KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP",
    }, "relay-header.log")
    relay.start()
    try:
        still, _, _, _ = Client(relay.port).call("POST", "/auth/verify", {"code": "000000"})
        allowed, location, _, _ = Client(relay.port, "198.51.100.91").call("POST", "/auth/start", {"email": address("fixture-header")})
        REPORT.check("p05-header-honored-only-when-configured", still == 429 and allowed == 303 and location == "/auth/verify", http=allowed)
        sink.take_code(address("fixture-header"))
    finally:
        relay.stop()

    closed = free_port()
    relay = Relay({
        "KNOWSLINK_SMTP_URL": f"smtp://127.0.0.1:{closed}",
        "KNOWSLINK_MAIL_FROM": "noreply@example.test",
        "KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP",
    }, "relay-closed.log")
    relay.start()
    try:
        email = address("fixture-fail")
        client = Client(relay.port, "192.0.2.80")
        status, _, body, _ = client.call("POST", "/auth/start", {"email": email})
        retry, _, retry_body, _ = client.call("POST", "/auth/start", {"email": email})
        REPORT.check(
            "p01-send-failure",
            status == 503 and "보내지 못했습니다" in body and "로그인되지 않았습니다" in body and count_challenges(email) == 0 and member_id(email) is None and retry == 429 and "이후 다시 시도" in retry_body,
            http=status,
        )
    finally:
        relay.stop()

    relay = Relay({"KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP"}, "relay-unconfigured.log")
    relay.start()
    try:
        status, _, body, _ = Client(relay.port, "192.0.2.81").call("POST", "/auth/start", {"email": address("fixture-unconfigured")})
        REPORT.check("p01-mail-unconfigured", status == 503 and "보내지 못했습니다" in body and member_id(address("fixture-unconfigured")) is None, http=status)
    finally:
        relay.stop()

    received = bytearray()

    def fake_server(port):
        server = socket.socket()
        server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        server.bind(("192.168.0.74", port))
        server.listen(1)
        server.settimeout(8)
        try:
            connection, _ = server.accept()
        except OSError:
            server.close()
            return
        with connection:
            connection.settimeout(5)
            connection.sendall(b"220 qa\r\n")
            while True:
                try:
                    chunk = connection.recv(1024)
                except OSError:
                    break
                if not chunk:
                    break
                received.extend(chunk)
                command = chunk.strip().upper()
                if command.startswith(b"EHLO") or command.startswith(b"HELO"):
                    connection.sendall(b"250-qa\r\n250 HELP\r\n")
                elif command == b"QUIT":
                    connection.sendall(b"221 bye\r\n")
                    break
                else:
                    connection.sendall(b"250 ok\r\n")
        server.close()

    port = free_port()
    try:
        thread = threading.Thread(target=fake_server, args=(port,))
        thread.start()
        relay = Relay({
            "KNOWSLINK_SMTP_URL": f"smtp://192.168.0.74:{port}",
            "KNOWSLINK_MAIL_FROM": "noreply@example.test",
            "KNOWSLINK_CLIENT_IP_HEADER": "CF-Connecting-IP",
        }, "relay-cleartext.log")
        relay.start()
        try:
            email = address("fixture-clear")
            status, _, body, _ = Client(relay.port, "192.0.2.82").call("POST", "/auth/start", {"email": email})
            digits = re.findall(rb"\d{6}", bytes(received))
            REPORT.check(
                "p01-cleartext-non-loopback-refused",
                status == 503 and count_challenges(email) == 0 and not digits,
                http=status,
            )
        finally:
            relay.stop()
        thread.join(5)
    except OSError:
        REPORT.skip("p01-cleartext-non-loopback-refused", "non-loopback fixture socket was unavailable")

    environment = os.environ.copy()
    environment["DATABASE_URL"] = db_url()
    environment["RELAY_ADDR"] = f"127.0.0.1:{free_port()}"
    environment["KNOWSLINK_CLIENT_IP_HEADER"] = "X-Forwarded-For"
    result = subprocess.run([str(QA / "build/relay")], cwd=QA, env=environment, check=False, capture_output=True, timeout=8)
    text = (result.stdout + result.stderr).decode("utf-8", "replace")
    REPORT.secret(db_url())
    REPORT.check("p05-bad-ip-header-rejected", result.returncode != 0 and "CF-Connecting-IP" in text and "@" not in text, http=result.returncode)


def git_ok():
    sha = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=QA, text=True).strip()
    porcelain = subprocess.check_output(["git", "status", "--porcelain"], cwd=QA, text=True)
    return sha == EXPECT and porcelain == ""


def main():
    WORK.mkdir(parents=True, exist_ok=True)
    os.chmod(WORK, 0o700)
    sink = None
    try:
        if not git_ok():
            raise RuntimeError("qa checkout is not the clean candidate")
        wait_database()
        migrate()
        sink = Sink()
        sink.start()
        directory_mode = sink.directory.stat().st_mode & 0o777
        REPORT.check("env-mail-dir-mode", directory_mode == 0o700, note=oct(directory_mode))
        phase_a(sink)
        phase_b(sink)
        phase_rates(sink)
        phase_capacity(sink)
        phase_header_and_smtp(sink)
        REPORT.check("env-qa-tree-clean", git_ok())
    except Exception as exc:
        import traceback
        frame = traceback.extract_tb(exc.__traceback__)[-1]
        (WORK / "probe.err").write_text(f"{type(exc).__name__} {frame.name}:{frame.lineno}\n", encoding="utf-8")
        REPORT.check("probe-exception", False, note=f"{type(exc).__name__} {frame.name}:{frame.lineno}")
    finally:
        if sink:
            sink.stop()
    logs = "\n".join(path.read_text(encoding="utf-8", errors="replace") for path in WORK.glob("relay*.log"))
    leaked = [item for item in REPORT.secrets if item and item in logs]
    REPORT.check("env-relay-log-redacted", not leaked, note=f"matches {len(leaked)}")
    REPORT.skip("p06-public-boundary", "운영 공개 후보와 현재 hostname의 owner-only 경계는 이 Run에 없다. 로컬 synthetic 403은 공개 PASS가 아니다.")
    REPORT.skip("p07-human-email-and-designer-ui", "실제 사용자 이메일과 designer의 직접 화면 검수는 미실행이다. HTML 문구 검사는 fixture 자동 근거다.")
    failed = [row for row in REPORT.rows if row["ok"] is False]
    payload = {
        "candidate": EXPECT,
        "qa_checkout": str(QA),
        "passed": sum(1 for row in REPORT.rows if row["ok"] is True),
        "failed": len(failed),
        "skipped": sum(1 for row in REPORT.rows if row.get("skipped")),
        "checks": REPORT.rows,
    }
    text = json.dumps(payload, ensure_ascii=True, indent=2)
    unsafe = [item for item in REPORT.secrets if item and item in text]
    if "@" in text or "postgres://" in text or unsafe:
        RESULT.write_text('{"failed": 1, "note": "result scrub blocked the detailed file"}\n', encoding="utf-8")
        print("FAIL result-scrub", flush=True)
        return 1
    RESULT.parent.mkdir(parents=True, exist_ok=True)
    RESULT.write_text(text + "\n", encoding="utf-8")
    os.chmod(RESULT, 0o644)
    print(f"SUMMARY passed {payload['passed']} failed {payload['failed']} skipped {payload['skipped']}", flush=True)
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
