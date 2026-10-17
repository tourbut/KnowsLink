"""SAR-SESSION-PEER-001-TESTER fixture: synthetic members on an isolated relay.

A,B: friends (accepted). C: no friendship. D: friend request from A, never accepted.
Each member gets one node-local agent connected into a private folder. Secrets stay in those folders.
usage: python3 fixture.py <relay> <mail-dir> <private-dir-prefix> <connect.js> <out.json>
"""
import json
import re
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

relay, mail_dir, prefix, connect_js, out_json = sys.argv[1:6]
NODE = r"C:\Users\shin\AppData\Local\KnowsLinkDevTools\node-v22.22.2-win-x64\node.EXE"
DOMAIN = "sessionpeerqa4.invalid"


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *a, **k):
        return None


class Browser:
    def __init__(self):
        self.cookies = {}

    def do(self, method, path, form=None):
        data = urllib.parse.urlencode(form).encode() if form is not None else None
        req = urllib.request.Request(relay + path, data=data, method=method)
        if data is not None:
            req.add_header("Content-Type", "application/x-www-form-urlencoded")
        if self.cookies:
            req.add_header("Cookie", "; ".join(f"{k}={v}" for k, v in self.cookies.items()))
        try:
            resp = urllib.request.build_opener(NoRedirect).open(req)
        except urllib.error.HTTPError as e:
            resp = e
        for header in resp.headers.get_all("Set-Cookie") or []:
            name, _, rest = header.partition("=")
            if "max-age=0" in header.lower():
                self.cookies.pop(name, None)
            else:
                self.cookies[name] = rest.split(";", 1)[0]
        return resp.status, resp.read().decode("utf-8", "replace"), resp.headers.get("Location", "")


def expect(result, status, *contains):
    code, body, location = result
    assert code == status, f"expected {status}, got {code}: {body[:300]}"
    for c in contains:
        assert c in body or c in location, f"missing {c!r}: {body[:300]} {location}"
    return body


def mail_code(email):
    deadline = time.time() + 30
    while time.time() < deadline:
        for p in sorted(Path(mail_dir).glob("*.eml")):
            text = p.read_text("utf-8", "replace")
            if email in text.lower():
                m = re.search(r"\b\d{6}\b", text)
                if m:
                    p.rename(p.with_suffix(".used"))
                    return m.group(0)
        time.sleep(0.5)
    raise RuntimeError("no mail code")


def login(email):
    b = Browser()
    for _ in range(8):  # relay allows one code mail per address per minute
        started = b.do("POST", "/auth/start", {"email": email})
        if started[0] != 429:
            break
        time.sleep(20)
    expect(started, 303, "/auth/verify")
    expect(b.do("POST", "/auth/verify", {"code": mail_code(email)}), 303, "/home")
    return b


def member(label, name, agent_name):
    email = f"{label.lower()}@{DOMAIN}"
    folder = f"{prefix}-{label.lower()}-private"
    b = login(email)
    expect(b.do("POST", "/home/profile", {"name": name}), 303)
    expect(b.do("POST", "/home/agents", {"name": agent_name}), 303)
    page = expect(b.do("GET", "/home/agents"), 200)
    agent = re.findall(r'식별자: <code class="select">(agent_[a-f0-9]+)</code>', page)[-1]
    page = expect(b.do("POST", "/home/connect", {"agent": agent, "client": "node-local", "mode": "register"}), 200, "연결 대기")
    token = re.search(r'aria-labelledby="token-l">([A-Za-z0-9_-]{43})</code>', page).group(1)
    out = subprocess.run([NODE, connect_js, "prepare", relay, folder], input=token, capture_output=True, text=True, check=True).stdout
    connection = re.search(r"자기 브라우저에서 확인: \S*?([A-Za-z0-9_-]{43})\s*$", out, re.M).group(1)
    expect(b.do("POST", "/home/confirm", {"connection": connection}), 303)
    subprocess.run([NODE, connect_js, "complete", folder], capture_output=True, text=True, check=True)
    return b, agent, folder


def friend_code(b):
    page = expect(b.do("GET", "/home/account"), 200)
    return re.search(r'aria-label="내 친구 코드">([A-Za-z0-9_-]+)</code>', page).group(1)


def request_friend(src, dst_code):
    expect(src.do("POST", "/home/friends/request", {"code": dst_code, "consent": "yes"}), 303, "n=requested")


def accept(dst, who_name):
    page = expect(dst.do("GET", "/home/friends"), 200, who_name)
    friend = re.search(r'name="friend" value="([^"]+)"><input type="hidden" name="generation" value="(\d+)"><input type="hidden" name="consent"', page)
    expect(dst.do("POST", "/home/friends/decision", {"friend": friend.group(1), "generation": friend.group(2), "consent": "yes", "decision": "accept"}), 303)


RESUME = len(sys.argv) > 6 and sys.argv[6] == "--resume"  # members exist (agents connected); redo login + friendships only
members = {}
for label, name, agent_name in (("a", "합성 소유자 A", "QA 세션 A"), ("b", "합성 소유자 B", "QA 세션 B"), ("c", "합성 소유자 C", "QA 세션 C"), ("d", "합성 소유자 D", "QA 세션 D")):
    if RESUME:
        folder = f"{prefix}-{label}-private"
        members[label] = (login(f"{label}@{DOMAIN}"), json.loads((Path(folder) / "agent.json").read_text())["agent"], folder)
        members[label] = (members[label][0], members[label][1], folder)
    else:
        members[label] = member(label, name, agent_name)
request_friend(members["a"][0], friend_code(members["b"][0]))
accept(members["b"][0], "합성 소유자 A")
request_friend(members["a"][0], friend_code(members["d"][0]))  # pending only: D never accepts
result = {k: {"agent": v[1], "folder": v[2], "email": f"{k}@{DOMAIN}"} for k, v in members.items()}
result["relay"] = relay
Path(out_json).write_text(json.dumps(result, ensure_ascii=False, indent=1), encoding="utf-8")
print(json.dumps(result, ensure_ascii=False))
