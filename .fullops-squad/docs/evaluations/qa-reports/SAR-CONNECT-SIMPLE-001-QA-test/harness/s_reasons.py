# CS-AC-11 live injections for the client output: unreachable, agent_rate, unconfirmed, rate_limited (KO/EN). expired/cancelled are in s_expiry_real.py / s_attempts.py / s_folders.py.
from lib import *
import http.client, threading
for lang, key in (("ko", "인터넷 연결을 확인"), ("en", "Check your internet")):
    code, ev, err = cli("http://127.0.0.1:1", W + f"/r-unreach-{lang}", "--host", "claude", "--lang", lang)
    last = ev[-1] if ev else {}
    rec("CS-AC-11", f"unreachable ({lang}): code, text, next action; exit 2", code == 2 and last.get("reason") == "unreachable" and key in last.get("message", ""), json.dumps(last, ensure_ascii=False))
# agent_rate: B made 11 agents in s_limits; reach 20 within 24h, the 21st must be refused (web), and a device connection too
B = Browser().login("b")
n = len(state("b")["agents"] or [])
for i in range(20 - n):
    s, raw, _ = B.get("/home/agents"); s, raw, h = B.post("/home/agents", {"attempt": field(raw, "attempt"), "name": f"B 속도 {i}"})
    assert s == 303, (s, raw[:200])
s, raw, _ = B.get("/home/agents"); s, raw, h = B.post("/home/agents", {"attempt": field(raw, "attempt"), "name": "B 21번째"})
t = B.text(raw)
rec("CS-AC-04", "web: 21st new agent in 24h is refused with agent_rate wording and a time", s == 429 and "오늘 새 에이전트를 많이 만들었습니다" in t, f"status={s} text={t[-400:]}")
def go(url):
    approve("b", url)
for lang in ("ko", "en"):
    code, ev, err = cli(URL, W + f"/r-agent-rate-{lang}", "--host", "claude", "--lang", lang, on_waiting=go)
    last = ev[-1] if ev else {}
    rec("CS-AC-11", f"agent_rate ({lang}): code, text with time, next action; retry_at present; exit 2", code == 2 and last.get("reason") == "agent_rate" and bool(last.get("retry_at")) and ("기존 에이전트" in last.get("message", "") or "existing agents" in last.get("message", "")), json.dumps(last, ensure_ascii=False))
st = state("b")
rec("CS-AC-04", "after agent_rate refusals: B still has the same 20 agents (10 revoked kept, 10 active), none newly made", len(st["agents"]) == 20 and len([a for a in st["agents"] if not a["revoked"]]) == 10, f"total={len(st['agents'])}")
# unconfirmed: consent given, key save lost -> client keeps agent and name and says unconfirmed (exit 3) once its wait ends; here we use h.mjs for a single poll
fx("/fixture/drop-complete?on=1")
r = hnode("start", URL, W + "/r-unconfirmed", env={"KNOWSLINK_HOST": "claude"}); cid = conn_id(r["url"])
b = Browser().login("e"); google_login(b, cid, "e"); consent(b, cid, "미확인 이름")
p = hnode("poll", W + "/r-unconfirmed")
fx("/fixture/drop-complete?on=0")
rec("CS-AC-11", "dropped key-save answer: single poll errors, nothing claims 'connected', no agent exists yet", "error" in p and not [a for a in agents_active("e") if a["keys"]], json.dumps(p, ensure_ascii=False))
code, ev, err = cli(URL, W + "/r-unconfirmed", "--host", "claude")
rec("CS-AC-11", "retry on the same folder checks that attempt first and connects with the chosen name", code == 0 and ev[-1].get("state") == "connected" and ev[-1].get("name") == "미확인 이름", json.dumps(ev[-1] if ev else {}, ensure_ascii=False))
# rate_limited: exhaust the anonymous per-IP budget, then one connection start
for i in range(45):
    c = http.client.HTTPConnection(BASE); c.request("POST", "/v1/connect/poll", "{}", {"Content-Type": "application/json"}); c.getresponse().read()
for lang in ("ko", "en"):
    code, ev, err = cli(URL, W + f"/r-rate-{lang}", "--host", "claude", "--lang", lang)
    last = ev[-1] if ev else {}
    rec("CS-AC-11", f"rate_limited ({lang}): code, text with time, next action", code == 2 and last.get("reason") == "rate_limited", json.dumps(last, ensure_ascii=False))
