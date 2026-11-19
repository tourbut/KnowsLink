# CS-AC-05 / CS-AC-06 / CS-AC-07 (device path): five failed-attempt kinds on member C leave agent counts unchanged; next connect works.
from lib import *
ba, bb, _ = friends("c", "e")          # C asks, E accepts
accept(bb)
st0 = state("c")
rec("CS-AC-05", "setup: C and E are friends, C has 0 agents and 0 attempts", st0["friends"] == 1 and not (st0["agents"] or []) and not (st0["attempts"] or []), f"friends={st0['friends']} agents={st0['agents']} attempts={st0['attempts']}")
C = Browser().login("c")
def counts():
    s, raw, _ = C.get("/home/agents"); t = C.text(raw)
    n_cards = len(re.findall(r'id="agent-agent_', raw))
    E = Browser().login("e"); s, rawe, _ = E.get("/home/friends")
    link = re.search(r'href="(/home/friends/[^"]+)"', rawe)
    seen = 0
    if link:
        s, rawf, _ = E.get(link.group(1)); seen = len(re.findall(r"/ [0-9a-f]{8} /", E.text(rawf)))
    return len(agents_active("c")), n_cards, seen, t, raw, rawe
base_counts = counts()
ids = {}
# 1 cancelled (CLI, KO), cancel button pressed on consent page
def cancel_flow(url):
    b, cid, _ = approve_stop(url); ids["cancel"] = cid; b.post("/home/cancel", {"connection": cid})
def approve_stop(url):
    b = Browser().login("c"); cid = conn_id(url); google_login(b, cid, "c"); return b, cid, None
code, ev, err = cli(URL, W + "/at-cancel", "--host", "claude", on_waiting=cancel_flow)
last = ev[-1]
rec("CS-AC-11", "cancelled (KO): reason code, text, next action; exit 2", code == 2 and last.get("reason") == "cancelled" and "다시 연결" in last["message"] and last.get("state") == "failed", json.dumps(last, ensure_ascii=False))
code, ev, err = cli(URL, W + "/at-cancel-en", "--host", "claude", "--lang", "en", on_waiting=cancel_flow)
last = ev[-1]
rec("CS-AC-11", "cancelled (EN): reason code, text, next action", code == 2 and last.get("reason") == "cancelled" and "Connect again" in last["message"], json.dumps(last, ensure_ascii=False))
# 2 server-side expiry after the browser login (relay clock moved), CLI side
def expire_flow(url):
    b, cid, _ = approve_stop(url); ids["expire"] = cid; fx(f"/fixture/exp?connection={cid}&delta=-1m", "POST")
code, ev, err = cli(URL, W + "/at-expire", "--host", "claude", on_waiting=expire_flow)
last = ev[-1]
rec("CS-AC-11", "relay-expired attempt: client ends with a reason (not bare failed)", code == 2 and last.get("reason") in ("expired", "cancelled") and last.get("message"), json.dumps(last, ensure_ascii=False))
# 3 app closed right after start (nobody opened the link)
r = hnode("start", URL, W + "/at-idle"); ids["idle"] = conn_id(r["url"])
# 4 Google login done, no consent (prepared)
r = hnode("start", URL, W + "/at-prepared"); ids["prepared"] = conn_id(r["url"]); b4 = Browser().login("c"); google_login(b4, ids["prepared"], "c")
# 5 consent given, key never saved (app died)
r = hnode("start", URL, W + "/at-unsaved"); ids["unsaved"] = conn_id(r["url"]); b5 = Browser().login("c"); google_login(b5, ids["unsaved"], "c")
dflt5, s5, raw5, _ = consent(b5, ids["unsaved"], "미저장 시도")
rec("CS-AC-05", "attempt 5 consent accepted (303) but no key saved", s5 == 303, f"status={s5} default_name={dflt5}")
after = counts()
rec("CS-AC-05", "after 5 attempts (cancel, expire, idle, prepared, consented-unsaved): C active agents and 내 에이전트 cards unchanged (0)", after[0] == 0 and after[1] == 0, f"state_agents={after[0]} cards={after[1]}")
rec("CS-AC-05", "friend E's view of C (friend detail page) lists no agent of C", after[2] == base_counts[2], f"agents on E's view of C before={base_counts[2]} after={after[2]}")
t = after[3]
sec = re.search(r"진행 중인 연결[\s\S]{0,900}", t)
rec("CS-AC-05", "'진행 중인 연결' section lists the open attempts (M-02) on C's 내 에이전트", bool(sec), (sec.group(0)[:500] if sec else "no section found"))
open(os.path.join(OUT, "evidence-c-agents.txt"), "w", encoding="utf-8").write(t)
# next normal connection succeeds (CLI, host claude → default name "Claude Code")
def normal(url):
    b, cid, c = approve("c", url); ids["normal_default"] = c[0]
code, ev, err = cli(URL, W + "/at-normal", "--host", "claude", on_waiting=normal)
last = ev[-1]
rec("CS-AC-05", "next normal connection succeeds after the failed attempts", code == 0 and last.get("state") == "connected" and last.get("agent"), json.dumps(last, ensure_ascii=False))
rec("CS-AC-07", "device path (host=claude): consent default name is 'Claude Code' and agent saved with that name", ids.get("normal_default") == "Claude Code" and last.get("name") == "Claude Code", f"default={ids.get('normal_default')} agent.name={last.get('name')}")
fin = counts()
rec("CS-AC-05", "after success: C has exactly 1 active agent, visible to E; failed attempts still not agents", fin[0] == 1 and fin[1] == 1 and fin[2] == base_counts[2] + 1, f"active={fin[0]} cards={fin[1]} E-sees={fin[2]} (base {base_counts[2]})")
# second claude → numbered
def second(url):
    b, cid, c = approve("c", url); ids["second_default"] = c[0]
code, ev, err = cli(URL, W + "/at-second", "--host", "claude", on_waiting=second)
rec("CS-AC-07", "second same-host agent gets the next number ('Claude Code 2')", ids.get("second_default") == "Claude Code 2" and ev[-1].get("name") == "Claude Code 2", f"default={ids.get('second_default')} saved={ev[-1].get('name')}")
def third(url):
    b, cid, c = approve("c", url); ids["codex_default"] = c[0]
code, ev, err = cli(URL, W + "/at-codex", "--host", "codex", on_waiting=third)
rec("CS-AC-07", "host=codex default name is the codex display name", bool(ids.get("codex_default")) and ids["codex_default"].startswith("Codex"), f"default={ids.get('codex_default')}")
# no host → 내 에이전트 N
def fourth(url):
    b, cid, c = approve("c", url); ids["nohost_default"] = c[0]
code, ev, err = cli(URL, W + "/at-nohost", on_waiting=fourth)
rec("CS-AC-07", "no host: default 내 에이전트 <번호>", (ids.get("nohost_default") or "").startswith("내 에이전트 "), f"default={ids.get('nohost_default')}")
json.dump(ids, open(os.path.join(OUT, "ids-attempts.json"), "w"))
