"""Independent A/B/C QA against the live final-SHA relay.

Does not overwrite the designer session file, print cookies, OTPs, friend codes,
or message bodies. A failure report records statuses and check names only.
"""
import email
import html
import http.client
import json
import os
import pathlib
import re
import secrets
import subprocess
import sys
import time
import urllib.parse
from html.parser import HTMLParser

BASE = pathlib.Path(sys.argv[1]).resolve()
SHA = sys.argv[2]
assert re.fullmatch(r"/tmp/knowslink-humanpeer001-[0-9a-f]{12}-[0-9a-f]{8}", str(BASE))
assert re.fullmatch(r"[0-9a-f]{40}", SHA)
assert subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=BASE / "repo", text=True).strip() == SHA
assert not subprocess.check_output(["git", "status", "--porcelain"], cwd=BASE / "repo", text=True).strip()
record = max(BASE.glob("qa-*/runtime.json"), key=lambda path: path.stat().st_mtime_ns)
metadata = json.loads(record.read_text())
assert metadata["sha"] == SHA
QA = record.parent
REPO = pathlib.Path(metadata.get("repo", str(BASE / "repo")))
origin = urllib.parse.urlparse(metadata["origin"])
assert origin.hostname == "127.0.0.1" and origin.scheme == "http"
os.umask(0o077)
OUT = QA / ("tester-independent-" + secrets.token_hex(4))
OUT.mkdir(mode=0o700)
checks = []
requests = []
TEXT_A = "검수독립발송일번"
TEXT_B = "검수독립에이전트답장"
TEXT_REPLAY = "검수독립재전송확인"
slept = False


class Forms(HTMLParser):
    def __init__(self, source):
        super().__init__()
        self.forms, self.current = [], None
        self.feed(source)

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == "form":
            self.current = {"action": attrs.get("action", ""), "values": {}}
            self.forms.append(self.current)
        if tag == "input" and self.current is not None and "name" in attrs:
            self.current["values"][attrs["name"]] = attrs.get("value", "")

    def handle_endtag(self, tag):
        if tag == "form":
            self.current = None


def mark(name, ok):
    checks.append({"name": name, "ok": bool(ok)})
    if not ok:
        raise AssertionError(name)


def call(cookie, method, path, values=None, evil=False, timeout=20):
    global slept
    route = urllib.parse.urlparse(path).path
    for attempt in (1, 2):
        connection = http.client.HTTPConnection(origin.hostname, origin.port, timeout=timeout)
        headers = {"Origin": "https://invalid.example" if evil else metadata["origin"]}
        if cookie:
            headers["Cookie"] = cookie
        body = urllib.parse.urlencode(values) if values is not None else None
        if body is not None:
            headers["Content-Type"] = "application/x-www-form-urlencoded"
        connection.request(method, path, body, headers)
        response = connection.getresponse()
        status = response.status
        location = response.getheader("Location", "")
        text = response.read().decode("utf-8", "replace")
        cookies = {}
        for key, value in response.getheaders():
            if key.lower() == "set-cookie":
                part = value.split(";", 1)[0]
                name, _, raw = part.partition("=")
                cookies[name] = raw
        connection.close()
        if status != 429 or attempt == 2:
            recorded = "/home/friends/{id}" if route.startswith("/home/friends/") and route.count("/") == 3 else route
            requests.append({"method": method, "path": recorded, "status": status, "evil": evil})
            return status, text, location, cookies
        if not slept:
            time.sleep(65)
            slept = True
        else:
            requests.append({"method": method, "path": route, "status": status, "evil": evil})
            return status, text, location, cookies
    raise AssertionError("unreachable")


def signup(label):
    address = label.lower() + "-" + secrets.token_hex(4) + "@humanpeer.invalid"
    status, _, _, cookies = call("", "POST", "/auth/start", {"email": address})
    mark(label + " otp start", status == 303 and "__Host-kl_pending" in cookies)
    pending = "__Host-kl_pending=" + cookies["__Host-kl_pending"]
    code = None
    started = time.time()
    recent = set()
    deadline = started + 20
    while time.time() < deadline and not code:
        for message_file in (QA / "mail").glob("*.eml"):
            if message_file.stat().st_mtime < started - 2:
                continue
            recent.add(message_file.name)
            raw = message_file.read_bytes()
            if address.encode() not in raw:
                continue
            message = email.message_from_bytes(raw)
            payload = message.get_payload(decode=True)
            text = payload.decode("utf-8", "replace") if isinstance(payload, bytes) else str(payload)
            match = re.search(r"코드: ([0-9]{6})", text)
            if match:
                code = match.group(1)
                break
        if not code:
            time.sleep(0.25)
    mark(label + " otp delivered to local sink recent=" + str(len(recent)), bool(code))
    status, _, _, cookies = call(pending, "POST", "/auth/verify", {"code": code})
    mark(label + " otp verify", status == 303 and "__Host-kl_session" in cookies)
    return "__Host-kl_session=" + cookies["__Host-kl_session"]


def node(action, values):
    prefix = "let raw='';for await(const chunk of process.stdin)raw+=chunk;const p=JSON.parse(raw);"
    code = {
        "prepare": "import {prepare} from './adapters/dist/connect.js';" + prefix + "process.stdout.write(JSON.stringify(await prepare(p.origin,p.folder,p.token)));",
        "complete": "import {complete} from './adapters/dist/connect.js';" + prefix + "await complete(p.folder);process.stdout.write('null');",
        "receive": "import {memberTransport} from './adapters/dist/text.js';import {writeFile} from 'node:fs/promises';import {randomUUID} from 'node:crypto';" + prefix + "const t=await memberTransport(p.folder);const r=await t.receive(m=>writeFile(p.folder+'/copy-'+randomUUID()+'.json',JSON.stringify(m),{mode:0o600,flag:'wx'}));process.stdout.write(JSON.stringify(r?{id:r.id,authorship:r.authorship}:null));",
        "send": "import {memberTransport} from './adapters/dist/text.js';" + prefix + "const t=await memberTransport(p.folder);process.stdout.write(JSON.stringify(await t.sendText(p.peer,p.text,p.key,true)));",
    }[action]
    result = subprocess.run(["node", "--input-type=module", "-e", code], input=json.dumps(values), cwd=REPO, text=True, capture_output=True)
    mark("adapter " + action, result.returncode == 0)
    return json.loads(result.stdout)


def create_agent(cookie, label, name):
    status, before, _, _ = call(cookie, "GET", "/home/agents")
    mark(label + " agents page", status == 200)
    before_ids = set(re.findall(r'name="agent" value="([^"]+)"', before))
    status, _, _, _ = call(cookie, "POST", "/home/agents", {"name": name})
    mark(label + " agent create", status == 303)
    status, after, _, _ = call(cookie, "GET", "/home/agents")
    added = set(re.findall(r'name="agent" value="([^"]+)"', after)) - before_ids
    mark(label + " agent visible", status == 200 and len(added) == 1)
    agent = added.pop()
    status, page, _, _ = call(cookie, "POST", "/home/connect", {"agent": agent, "client": "node-local", "mode": "register"})
    match = re.search(r'<code[^>]*aria-labelledby="token-l"[^>]*>([^<]+)</code>', page)
    mark(label + " connect grant", status == 200 and bool(match))
    token = html.unescape(match.group(1)).strip()
    folder = str(OUT / label)
    prepared = node("prepare", {"origin": metadata["origin"], "folder": folder, "token": token})
    connection_id = urllib.parse.urlparse(prepared["connection"]).path.rsplit("/", 1)[1]
    status, _, _, _ = call(cookie, "POST", "/home/confirm", {"connection": connection_id})
    mark(label + " confirm", status == 303)
    node("complete", {"folder": folder})
    return agent, folder


def friend_code(cookie):
    status, page, _, _ = call(cookie, "GET", "/home/account")
    match = re.search(r'<code[^>]*aria-label="내 친구 코드"[^>]*>([^<]+)</code>', page)
    mark("friend code issued", status == 200 and bool(match))
    return html.unescape(match.group(1)).strip()


def decision_form(cookie):
    status, page, _, _ = call(cookie, "GET", "/home/friends")
    mark("friends page", status == 200)
    forms = [item for item in Forms(page).forms if item["action"] == "/home/friends/decision" and "friend" in item["values"]]
    mark("one friend decision form", len(forms) == 1)
    return forms[0]["values"], page


def composer(cookie, sender, target):
    path = "/home/conversation?" + urllib.parse.urlencode({"from": sender, "to": target})
    status, page, _, _ = call(cookie, "GET", path)
    mark("composer page " + sender[:8], status == 200 and 'action="/home/conversation/send"' in page)
    form = next(item for item in Forms(page).forms if item["action"] == "/home/conversation/send")
    return form["values"], page


def main():
    mark("healthz", call("", "GET", "/healthz")[0] == 200)
    mark("anonymous home", call("", "GET", "/home")[0] == 303)
    mark("anonymous legacy invite", call("", "POST", "/v1/invites", {})[0] in (401, 403, 422))
    a = signup("A")
    b = signup("B")
    c = signup("C")
    status, page, _, _ = call(a, "POST", "/home/profile", {"name": "가" * 41})
    mark("41 rune name rejected", status == 422 and "1–40" in page)
    status, _, _, _ = call(a, "POST", "/home/profile", {"name": "가" * 40})
    mark("40 rune name accepted", status == 303)
    status, _, _, _ = call(a, "POST", "/home/profile", {"name": "검수회원가"})
    mark("A named", status == 303)
    status, _, _, _ = call(b, "POST", "/home/profile", {"name": "검수회원가"})
    mark("same display name allowed", status == 303)
    status, _, _, _ = call(b, "POST", "/home/profile", {"name": "검수회원나"})
    mark("B named", status == 303)
    status, _, _, _ = call(c, "POST", "/home/profile", {"name": "검수회원다"})
    mark("C named", status == 303)
    status, page, _, _ = call(a, "POST", "/home/agents", {"name": "가" * 41})
    mark("41 rune agent name rejected", status == 422 and "1–40" in page)
    agent_a, folder_a = create_agent(a, "a-main", "검수에이전트가")
    agent_a2, _ = create_agent(a, "a-extra", "검수에이전트가둘")
    agent_b, folder_b = create_agent(b, "b-main", "검수에이전트나")
    agent_c, _ = create_agent(c, "c-main", "검수에이전트다")
    values, _ = composer(a, agent_a, agent_a2)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text="검수자기연결"))
    mark("self connect send without friend", status == 303)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", {"from": agent_a, "to": agent_b, "text": TEXT_A, "attempt": "tester-before-friend-01"})
    mark("send before consent denied", status == 403)
    code_b = friend_code(b)
    status, _, location, _ = call(c, "POST", "/home/friends/request", {"code": "not-a-real-friend-code", "consent": "yes"})
    mark("unknown code still redirects", status == 303 and "n=requested" in location)
    status, page, _, _ = call(c, "GET", "/home/friends")
    mark("unknown code does not reveal B", status == 200 and "검수회원나" not in page and "요청 대기" in page)
    _, cancel_page = decision_form(c)
    mark("outgoing hides target name", "검수회원나" not in cancel_page)
    status, _, _, _ = call(c, "POST", "/home/friends/decision", dict(decision_form(c)[0], decision="cancel"))
    mark("unknown request cancelled", status == 303)
    status, _, _, _ = call(a, "POST", "/home/friends/request", {"code": code_b, "consent": "yes"})
    mark("A requests B", status == 303)
    status, page, _, _ = call(a, "GET", "/home/friends")
    mark("outgoing hides B name", status == 200 and "요청 대기" in page and "검수회원나" not in page)
    status, page, _, _ = call(b, "GET", "/home/friends")
    mark("incoming shows A name", status == 200 and "검수회원가" in page)
    form, _ = decision_form(b)
    generation = int(form["generation"])
    status, _, _, _ = call(c, "POST", "/home/friends/decision", dict(form, decision="accept", consent="yes"))
    mark("C cannot accept A-B", status == 403)
    status, _, _, _ = call(a, "POST", "/home/friends/decision", dict(form, decision="accept", consent="yes"))
    mark("inviter cannot accept", status == 403)
    status, _, _, _ = call(b, "POST", "/home/friends/decision", dict(form, decision="accept", consent="yes"), evil=True)
    mark("cross origin accept denied", status == 403)
    status, _, _, _ = call(b, "POST", "/home/friends/decision", dict(form, decision="accept", consent="yes"))
    mark("B accepts", status == 303)
    status, _, _, _ = call(b, "POST", "/home/friends/decision", dict(form, decision="accept", consent="yes"))
    mark("duplicate accept is idempotent", status == 303)
    status, page, _, _ = call(a, "GET", "/home/friends/" + form["friend"])
    mark("accepted friend shows current agent", status == 200 and "검수에이전트나" in page)
    values, page = composer(a, agent_a, agent_b)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text=TEXT_A), evil=True)
    mark("cross origin send denied", status == 403)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text=TEXT_A))
    mark("owner send accepted", status == 303)
    received = node("receive", {"folder": folder_b})
    mark("owner authorship preserved", received.get("authorship") == "owner" and bool(received.get("id")))
    sent = node("send", {"folder": folder_b, "peer": agent_a, "text": TEXT_B, "key": "tester-qa-" + secrets.token_hex(12)})
    received = node("receive", {"folder": folder_a})
    mark("agent reply authorship preserved", received.get("authorship") == "agent" and received.get("id") == sent.get("id"))
    status, page, _, _ = call(a, "GET", "/home/conversation?" + urllib.parse.urlencode({"from": agent_a, "to": agent_b}))
    mark("both directions visible to A", status == 200 and TEXT_A in page and TEXT_B in page)
    status, page, _, _ = call(c, "GET", "/home/conversation?" + urllib.parse.urlencode({"from": agent_a, "to": agent_b}))
    mark("C cannot read A-B conversation", status == 403 and TEXT_A not in page and TEXT_B not in page)
    status, _, _, _ = call(c, "POST", "/home/conversation/send", {"from": agent_c, "to": agent_b, "text": TEXT_A, "attempt": "tester-third-party-01"})
    mark("C cannot send into A-B", status == 403)
    status, _, _, _ = call("", "POST", "/home/conversation/send", dict(values, text=TEXT_A))
    mark("logged out send rejected", status in (303, 401, 403))
    status, page, _, _ = call("", "GET", "/home")
    mark("logged out home rejected", status == 303 and TEXT_A not in page)
    values, _ = composer(a, agent_a, agent_b)
    payload = dict(values, text=TEXT_REPLAY)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", payload)
    mark("replay original accepted", status == 303)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", payload)
    mark("same attempt and body is idempotent", status == 303)
    status, page, _, _ = call(a, "POST", "/home/conversation/send", dict(payload, text=TEXT_REPLAY + "변경"))
    mark("same attempt different body conflicts", status == 409)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text="   ", attempt=values["attempt"] + "x"))
    mark("blank text rejected", status == 422)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text="a" * 4097, attempt=values["attempt"] + "y"))
    mark("over 4096 bytes rejected", status == 422)
    values, _ = composer(a, agent_a2, agent_b)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text="검수추가에이전트발송"))
    mark("second agent sends without new friend accept", status == 303)
    code_a = friend_code(a)
    code_c = friend_code(c)
    status, _, _, _ = call(a, "POST", "/home/friends/request", {"code": code_c, "consent": "yes"})
    mark("A requests C", status == 303)
    status, _, _, _ = call(c, "POST", "/home/friends/request", {"code": code_a, "consent": "yes"})
    mark("opposite request does not auto-accept", status == 303)
    status, page, _, _ = call(a, "GET", "/home/friends")
    mark("opposite request stays pending", status == 200 and "검수회원다" not in page and "요청 대기" in page)
    form, _ = decision_form(a)
    status, _, _, _ = call(a, "POST", "/home/friends/decision", dict(form, decision="cancel"))
    mark("A cancels outgoing", status == 303)
    status, _, _, _ = call(c, "POST", "/home/friends/request", {"code": code_a, "consent": "yes"})
    mark("C requests A after cancel", status == 303)
    form, page = decision_form(a)
    mark("deny shows inviter", "검수회원다" in page)
    status, _, _, _ = call(a, "POST", "/home/friends/decision", dict(form, decision="deny"))
    mark("A denies C", status == 303)
    status, page, _, _ = call(a, "GET", "/home/friends")
    mark("denied pair is not an active friend", status == 200 and "검수회원다" not in page)
    status, page, _, _ = call(a, "GET", "/home/friends")
    friend_ids = re.findall(r'href="/home/friends/([^"]+)"', page)
    mark("active friend page exists", len(friend_ids) == 1)
    status, page, _, _ = call(b, "GET", "/home/friends/" + friend_ids[0])
    block = next(item for item in Forms(page).forms if item["action"] == "/home/friends/decision")
    status, _, _, _ = call(b, "POST", "/home/friends/decision", dict(block["values"], decision="block"))
    mark("B blocks A", status == 303)
    status, page, _, _ = call(a, "GET", "/home/conversation?" + urllib.parse.urlencode({"from": agent_a, "to": agent_b}))
    mark("blocked sender keeps read-only history", status == 200 and TEXT_A in page and 'action="/home/conversation/send"' not in page)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", {"from": agent_a, "to": agent_b, "text": TEXT_REPLAY, "attempt": payload["attempt"]})
    mark("old attempt rejected after block", status == 403)
    status, page, _, _ = call(b, "GET", "/home/friends")
    unblock = next(item for item in Forms(page).forms if item["values"].get("friend") == friend_ids[0])
    status, _, _, _ = call(b, "POST", "/home/friends/decision", dict(unblock["values"], decision="unblock"))
    mark("B unblocks without restoring friendship", status == 303)
    status, page, _, _ = call(a, "GET", "/home/friends")
    mark("unblock is not active friendship", status == 200 and 'href="/home/friends/' not in page)
    status, _, _, _ = call(a, "POST", "/home/friends/request", {"code": code_b, "consent": "yes"})
    mark("re-request after unblock", status == 303)
    form, page = decision_form(b)
    mark("re-accept is a new generation", int(form["generation"]) > generation and "검수회원가" in page)
    status, _, _, _ = call(b, "POST", "/home/friends/decision", dict(form, decision="accept", consent="yes"))
    mark("B re-accepts", status == 303)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", {"from": agent_a, "to": agent_b, "text": TEXT_REPLAY, "attempt": payload["attempt"]})
    mark("pre-block attempt stays rejected", status == 403)
    values, _ = composer(a, agent_a, agent_b)
    status, _, _, _ = call(a, "POST", "/home/conversation/send", dict(values, text="검수재연결새발송"))
    mark("new generation can send", status == 303)
    status, _, _, _ = call(a, "GET", "/v1/contacts")
    mark("browser session is not agent contacts credential", status in (401, 403))
    status, page, _, _ = call(a, "GET", "/home/legacy")
    mark("legacy page does not bypass current friendship", status in (200, 303, 403) and "검수독립발송일번" not in page)


report = {"sha": SHA, "origin_port": origin.port, "qa": str(QA), "evidence": str(OUT), "scope": "tester independent live A/B/C; designer alice/bob session file not replaced", "secrets_in_output": False}
try:
    main()
    report["result"] = "PASS"
    code = 0
except Exception as error:
    report["result"] = "FAIL"
    report["error"] = type(error).__name__
    report["failed_check"] = str(error)[:180]
    code = 1
report["checks"] = checks
report["requests"] = requests
(OUT / "independent-qa.json").write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
print(json.dumps(report, ensure_ascii=False, indent=2))
sys.exit(code)
