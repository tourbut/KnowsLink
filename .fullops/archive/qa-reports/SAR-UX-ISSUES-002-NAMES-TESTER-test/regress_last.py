"""Last, smallest state-changing regression for 5fa3333 on the shared fixture: friend-request non-disclosure (U4).
C submits (a) a missing code (adds one placeholder 'sent' row), (b) A's real code (C already has a pending request to A: no new request), (c) A's code with changed case, (d) A's code with padding.
All four must give the same status and status text. No friend decision, no code rotation, no rename is saved."""
import html as H, json, os, re, sys, traceback
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__))); from client import *
A, C = as_("a"), as_("c")
RES = []; EVID = {}
code_a = re.search(r'<textarea id="my-code"[^>]*>([^<]*)</textarea>', A.get("/home/account")[2]).group(1)
sent = lambda: len(re.findall(r"<h3>보낸 친구 요청</h3>", main(C.get("/home/friends")[2])))
rcv = lambda: re.findall(r"<h3>(.*?)의 친구 요청</h3>", main(A.get("/home/friends")[2]))
sig = lambda r: (r[0], [H.unescape(text(x)) for x in re.findall(r'role="status"[^>]*>(.*?)</(?:p|div|section)>', main(r[2]), re.S)])
def follow(r):
    s, h, b = r
    return C.get(h["location"]) if s in (301, 302, 303) and h.get("location") else r
try:
    n0, r0 = sent(), rcv()
    out = {}
    for lab, v in (("missing", "Zq" + "7" * 41), ("a_real_duplicate", code_a), ("a_case_changed", code_a.swapcase()), ("a_padded", "  " + code_a + "  ")):
        out[lab] = sig(follow(C.post("/home/friends/request", {"code": v, "consent": "yes"})))
    EVID["signatures"] = out; EVID["sent_rows_before_after"] = [n0, sent()]; EVID["a_received_before_after"] = [r0, rcv()]
    ok = len({json.dumps(v, ensure_ascii=False) for v in out.values()}) == 1 and out["missing"][0] == 200 and out["missing"][1]
    ok = ok and rcv() == r0 and sent() <= n0 + 2  # first run used <= n0 + 1 (wrong: a case-swapped code is a second unknown code, so two placeholders are expected); the first run's result file is kept as is
    RES.append({"check": "friend_request_answer_is_identical_for_missing_real_duplicate_case_changed_and_padded_codes_and_creates_no_extra_request", "result": "pass" if ok else "fail", "detail": "" if ok else str(EVID)[:600]})
except Exception as e:
    RES.append({"check": "friend_request_non_disclosure", "result": "fail", "detail": f"{type(e).__name__}: {e}", "trace": traceback.format_exc(limit=3)})
EVID["member_window_429_retries"] = Browser.waited429
json.dump({"origin": ORIGIN, "fixed_sha": "5fa3333b5948fd36792662b202915bc6243cac14", "evidence": EVID, "results": RES}, open(os.path.join(os.path.dirname(os.path.abspath(__file__)), "regress-last-5fa3333.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
for r in RES: print(r["result"].upper(), r["check"], r.get("detail", ""))
sys.exit(0 if all(r["result"] == "pass" for r in RES) else 1)
