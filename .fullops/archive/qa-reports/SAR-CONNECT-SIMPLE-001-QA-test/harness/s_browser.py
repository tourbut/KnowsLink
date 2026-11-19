# CS-AC-12 (link revisit by starter / other member / logged out) and CS-AC-13 (Google failure causes). Part "signup" fills the member cap: run it last.
from lib import *
part = sys.argv[1] if len(sys.argv) > 1 else "main"
GEN = "유효하지 않습니다"
def callback_for(b, cid, acct):
    return google_login(b, cid, acct)
if part == "main":
    # W-01: completed link
    def ok(url): approve("a", url)
    code, ev, err = cli(URL, W + "/w01", "--host", "claude", on_waiting=ok)
    link = URL + "/connect/" + conn_id(ev[0]["url"]); cid = conn_id(ev[0]["url"]); name = ev[-1].get("name")
    A = Browser().login("a"); s, raw, _ = A.get(f"/connect/{cid}"); t = A.text(raw)
    rec("CS-AC-12", "starter A reopens a completed link: completion screen with the agent name", s == 200 and "이 연결은 완료됐습니다" in t and name in t, f"status={s} name={name} text={t[:300]}")
    rec("CS-AC-12", "completed screen does not claim expiry/cancel/invalid", not any(w in t for w in ("만료", "취소됐", GEN)), t[:200])
    Bx = Browser().login("b"); s, raw, _ = Bx.get(f"/connect/{cid}"); t = Bx.text(raw)
    rec("CS-AC-12", "other member B opens the same link: current generic refusal, no agent name/state/email", s == 401 and GEN in t and name not in t and "이 연결은 완료됐습니다" not in t and "csqa-" not in t, f"status={s} text={t[-300:]}")
    out = Browser(); s, raw, _ = out.get(f"/connect/{cid}"); t = out.text(raw)
    rec("CS-AC-12", "logged-out browser opens the link: generic refusal, no name/state", s == 401 and GEN in t and name not in t and "이 연결은 완료됐습니다" not in t, f"status={s} text={t[-300:]}")
    # approved but key not yet saved
    fx("/fixture/drop-complete?on=1")
    r = hnode("start", URL, W + "/w01-pending"); pid = conn_id(r["url"]); b = Browser().login("a"); google_login(b, pid, "a"); consent(b, pid, "저장 대기")
    hnode("poll", W + "/w01-pending"); fx("/fixture/drop-complete?on=0")
    s, raw, _ = b.get(f"/connect/{pid}"); t = b.text(raw)
    rec("CS-AC-12", "approved but key not saved: starter sees the D04 text (title 연결을 승인했습니다; 'Approval alone doesn't connect it yet') and not 'connected/완료' (W-01, W-04)", s == 200 and "연결을 승인했습니다" in t and "승인만으로는 아직 연결되지 않았습니다" in t and "이 연결은 완료됐습니다" not in t, f"status={s} text={t[:400]}")
    # cancelled / expired revisit by starter
    r = hnode("start", URL, W + "/w01-cancel"); xid = conn_id(r["url"]); b = Browser().login("a"); google_login(b, xid, "a"); b.post("/home/cancel", {"connection": xid})
    s, raw, _ = b.get(f"/connect/{xid}"); t = b.text(raw)
    rec("CS-AC-12", "cancelled link, starter: 'cancelled' + connect-again-from-app guidance", s == 200 and "이 연결은 취소됐습니다" in t and "시작한 앱에서 다시 연결" in t, f"status={s} text={t[:400]}")
    r = hnode("start", URL, W + "/w01-expire"); yid = conn_id(r["url"]); b = Browser().login("a"); google_login(b, yid, "a"); fx(f"/fixture/exp?connection={yid}&delta=-1m", "POST")
    s, raw, _ = b.get(f"/connect/{yid}"); t = b.text(raw)
    rec("CS-AC-12", "expired link, starter: reason + connect-again-from-app guidance (W-01)", s == 200 and "이 연결은 만료됐습니다" in t and "시작한 앱에서 다시 연결" in t, f"status={s} text={t[:400]}")
    # W-02: a different account / expired login
    r = hnode("start", URL, W + "/w02-other"); oid = conn_id(r["url"])
    a = Browser().login("a"); google_login(a, oid, "a")                        # A claims the attempt
    bb = a                                                                     # same browser (A's session) but Google account B is chosen
    s, raw, _ = google_login(bb, oid, "b"); t = bb.text(raw)
    other = t
    rec("CS-AC-13", "other account started this connection: specific message, retry form back to Google account choice", s == 403 and "다른 Google 계정으로 시작됐습니다" in t and 'action="/auth/google"' in raw and "다른 계정" in t, f"status={s} text={t[-300:]}")
    rec("CS-AC-13", "other-account response has no email/name/id of A or B", not any(x in raw for x in ("csqa-a@", "csqa-b@", "합성 회원 A", "합성 회원 B", "google-qa-")), "searched emails, member names, subjects")
    r = hnode("start", URL, W + "/w02-exp"); eid = conn_id(r["url"]); bx = Browser()
    s0, raw0, h0 = bx.post("/auth/google", {"connection": eid}); q = urllib.parse.parse_qs(urllib.parse.urlparse(h0["Location"]).query)
    Browser().req("GET", urllib.parse.urlparse(h0["Location"]).path + "?" + urllib.parse.urlparse(h0["Location"]).query)
    fx(f"/fixture/exp?connection={eid}&delta=-1m", "POST")
    s, raw, _ = bx.get("/auth/google/callback?" + urllib.parse.urlencode({"code": "a." + q["nonce"][0], "state": q["state"][0]})); t = bx.text(raw)
    rec("CS-AC-13", "login deadline passed (attempt expired during Google): specific 'login time passed' message, no retry-as-other-account", s == 403 and "로그인 시간이 지났습니다" in t and "시작한 앱에서" in t and "csqa-" not in raw, f"status={s} text={t[-300:]}")
    rec("CS-AC-13", "the two messages differ", other[-300:] != t[-300:] and "다른 Google 계정으로 시작됐습니다" not in t, "")
else:
    n = fx("/fixture/fill-members", "POST")[1]
    r = hnode("start", URL, W + "/w02-signup"); sid = conn_id(r["url"]); b = Browser()
    s, raw, _ = google_login(b, sid, "newbie"); t = b.text(raw)
    rec("CS-AC-13", f"signups closed (cap filled with {n} synthetic members): specific message for a new Google account, retry to account choice, no other-account info", s == 403 and "새 가입을 받지 않습니다" in t and 'action="/auth/google"' in raw and "csqa-" not in raw, f"status={s} text={t[-300:]}")
