# CS-AC-06 (member D): '이 시도 정리' deletes at once; consented-unsaved attempt completes normally before cleanup; unsaved attempt is purged after exp+24h; keyed agent untouched.
from lib import *
D = Browser().login("d")
def unsaved(folder, name):
    r = hnode("start", URL, folder); cid = conn_id(r["url"])
    b = Browser().login("d"); google_login(b, cid, "d"); consent(b, cid, name); return cid
def has(cid): return any(a["id"] == cid for a in attempts("d"))
# (a) immediate cleanup button
cid = unsaved(W + "/cl-now", "정리 대상")
s, raw, _ = D.get("/home/agents")
form = re.search(r'<form[^>]*action="/home/attempt-remove"[\s\S]*?</form>', raw)
sec = D.text(raw)
rec("CS-AC-06", "unsaved consented attempt shows with '이 시도 정리' button on 내 에이전트 (M-02)", bool(form) and cid in raw and "이 시도 정리" in sec, (form.group(0)[:200] if form else "no form") )
before = has(cid)
s, raw2, h = D.post("/home/attempt-remove", {"connection": cid})
rec("CS-AC-06", "'이 시도 정리' removes the attempt immediately (server state)", before and not has(cid) and s in (200, 303), f"before={before} after={has(cid)} status={s}")
p = hnode("poll", W + "/cl-now")
rec("CS-AC-06", "after cleanup the app's poll is refused (no phantom agent)", "error" in p and not agents_active("d") , json.dumps(p, ensure_ascii=False))
# (b) save before cleanup completes normally (exp already passed, still within exp+24h)
fx("/fixture/drop-complete?on=1")
r = hnode("start", URL, W + "/cl-save"); cid = conn_id(r["url"])
b = Browser().login("d"); google_login(b, cid, "d"); consent(b, cid, "정리 전 저장")
p = hnode("poll", W + "/cl-save")                # approved seen, key save lost on the way
fx("/fixture/drop-complete?on=0")
fx(f"/fixture/exp?connection={cid}&delta=-2h", "POST")
st_before = [a for a in attempts("d") if a["id"] == cid]
code, ev, err = cli(URL, W + "/cl-save")
last = ev[-1] if ev else {}
rec("CS-AC-06", "first save attempt dropped: poll shows an error and no agent exists yet (A-01)", "error" in p or p.get("state") != "connected", json.dumps(p, ensure_ascii=False)[:300] + f" attempts={[a['state'] for a in st_before]}")
rec("CS-AC-06", "after exp passed (<24h) the saved-late completion still connects normally", code == 0 and last.get("state") == "connected" and last.get("name") == "정리 전 저장", json.dumps(last, ensure_ascii=False))
keyed = [a for a in agents_active("d") if a["keys"] > 0]
rec("CS-AC-06", "the completed agent has a stored key and its chosen name", len(keyed) == 1 and keyed[0]["name"] == "정리 전 저장", json.dumps(keyed, ensure_ascii=False))
# (c) purge after exp + 24h; keyed agent untouched
cid = unsaved(W + "/cl-old", "오래된 미저장")
fx(f"/fixture/exp?connection={cid}&delta=-23h", "POST"); D.get("/home/agents")
still = has(cid)
fx(f"/fixture/exp?connection={cid}&delta=-25h", "POST"); D.get("/home/agents"); D.get("/home/agents")
gone = not has(cid)
rec("CS-AC-06", "unsaved attempt 23h past exp is kept, 25h past exp is purged by the relay", still and gone, f"kept@23h={still} purged@25h={gone}")
p = hnode("poll", W + "/cl-old")
rec("CS-AC-06", "purged attempt cannot be completed later", "error" in p, json.dumps(p, ensure_ascii=False))
keyed2 = [a for a in agents_active("d") if a["keys"] > 0]
rec("CS-AC-06", "keyed agent is untouched by cleanup (A-03)", len(keyed2) == 1 and keyed2[0]["id"] == keyed[0]["id"], json.dumps(keyed2, ensure_ascii=False))
