# CS-AC-01 (7 active agents, no limit wording) and CS-AC-02 (10 revoked records kept, new agent still created).
from lib import *
BAD = ["회원당", "5개", "먼저 철회", "철회하세요", "남은 ", "한도", "최대 5", "per member", "5 agents", "revoke unused", "limit"]
def scan(b, pages, label):
    hits = {}
    for p in pages:
        s, raw, _ = b.get(p)
        t = b.text(raw)
        for w in BAD:
            if w in t: hits.setdefault(p, []).append(w)
    return hits
A = Browser().login("a")
made = []
for i in range(7):
    s, raw, h = A.get("/home/agents"); att = field(raw, "attempt")
    s, raw, h = A.post("/home/agents", {"attempt": att, "name": f"합성 에이전트 {i+1}"})
    made.append(s)
act = agents_active("a")
rec("CS-AC-01", "member A creates 7 active agents on the web; none refused", made == [303] * 7 and len(act) == 7, f"statuses={made} active={len(act)}")
s, raw, h = A.get("/home/agents"); att = field(raw, "attempt")
s, raw, h = A.post("/home/agents", {"attempt": att, "name": "합성 에이전트 8"})
rec("CS-AC-01", "8th agent also accepted (no count wall at 5 or 7)", s == 303 and len(agents_active("a")) == 8, f"status={s}")
pages = ["/home", "/home/agents", "/home/new", "/home/friends", "/home/account", "/home/legacy", "/home/receipts"]
ko = scan(A, pages, "ko")
rec("CS-AC-01", "KO screens contain no limit wording (회원당/5개/먼저 철회/남은/한도) with 8 agents", not ko, json.dumps(ko, ensure_ascii=False))
A.req("GET", "/lang/en"); en = scan(A, pages, "en")
rec("CS-AC-01", "EN screens contain no limit wording with 8 agents", not en, json.dumps(en, ensure_ascii=False))
# a refused web create (bad name) must not use limit wording either
s, raw, _ = A.post("/home/agents", {"attempt": "x" * 20, "name": ""}); t = A.text(raw)
rec("CS-AC-01", "invalid-name refusal has no limit wording", not any(w in t for w in BAD), f"status={s}")
# CS-AC-02: B keeps 10 revoked agents, then makes a new one (web) and a device one is covered in s_attempts
B = Browser().login("b")
ids = []
for i in range(10):
    s, raw, h = B.get("/home/agents"); att = field(raw, "attempt")
    s, raw, h = B.post("/home/agents", {"attempt": att, "name": f"B 합성 {i+1}"})
    assert s == 303, (s, raw[:200])
    ids.append(h["Location"].split("#agent-")[1])
rev = []
for i in ids:
    s, raw, h = B.post("/home/agent-revoke", {"agent": i}); rev.append(s)
st = state("b"); revoked = [a for a in st["agents"] if a["revoked"]]
rec("CS-AC-02", "B revokes 10 agents; all 10 records kept (revoked=true)", len(revoked) == 10 and all(s in (200, 303) for s in rev), f"revoke statuses={set(rev)} revoked={len(revoked)}")
s, raw, h = B.get("/home/agents"); att = field(raw, "attempt")
s, raw, h = B.post("/home/agents", {"attempt": att, "name": "B 새 에이전트"})
rec("CS-AC-02", "with 10 revoked records kept, B still creates a new agent (web)", s == 303 and len(agents_active("b")) == 1, f"status={s} loc={h.get('Location')}")
ko2 = scan(B, pages, "ko-b")
rec("CS-AC-02", "no limit wording on B's screens", not ko2, json.dumps(ko2, ensure_ascii=False))
