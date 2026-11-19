"""Follow-up after mutate.py's two first-run failures (kept in mutate-result.json, not edited).
1) 'spaces' expected value="" — wrong expectation: the 422 keeps the typed text, so spaces stay. Re-asserted here with the contract (input preserved).
2) 'bel' asserted that rejected control/bidi characters are not echoed — mutate.py stopped there, so the later rejected-name cases were never reached. Re-run here without stopping, and the echo is kept as its own check.
Idempotent: every case here is refused by the product, so it can run any number of times."""
import html as H, json, os, re, sys, traceback
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__))); from client import *

M = json.load(open(os.environ.get("NAMES_MANIFEST", "C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json")))
AB = M["friend_ab"]; A_AGENT, B_AGENT = "긴한글" * 13, "😀" * 40
RES = []; EVID = {}


ONLY = set(filter(None, os.environ.get("ONLY_CHECKS", "").split(",")))  # recheck 5fa3333: run just these checks (state-changing chain is not replayed on the shared fixture)


def check(name, covers):
    def deco(fn):
        if ONLY and name not in ONLY: return fn
        try: fn(); RES.append({"check": name, "covers": covers, "result": "pass"})
        except Exception as e: RES.append({"check": name, "covers": covers, "result": "fail", "detail": f"{type(e).__name__}: {e}", "trace": traceback.format_exc(limit=3)})
        return fn
    return deco


def follow(r, br):
    s, h, b = r
    return br.get(h["location"]) if s in (301, 302, 303) and h.get("location") else r


A, B, C = as_("a"), as_("b"), as_("c")
alerts = lambda b: [H.unescape(text(x)) for x in re.findall(r'role="alert"[^>]*>(.*?)</(?:p|div)>', main(b), re.S)]
n_sent = lambda: len(re.findall(r"<h3>보낸 친구 요청</h3>", main(C.get("/home/friends")[2])))
agents = lambda: re.findall(r'<h2 class="name-preview">([^<]*)</h2><p class="muted">[^<]*/ (agent_\w+) /', A.get("/home/agents")[2])
BAD_NAMES = {"empty": "", "spaces": "   ", "bel": "ab\x07cd", "nul": "ab\x00cd", "bidi_rlo": "ab\u202ecd", "bidi_isolate": "ab\u2066cd", "newline": "ab\ncd", "len41": "가" * 41, "emoji41": "😀" * 41}
INVIS = lambda s: bool(re.search(r"[\x00-\x08\x0b-\x1f\x7f\u202a-\u202e\u2066-\u2069]", s))


@check("friend_request_invalid_inputs_422_preserve_typed_text_and_create_nothing", ["U4 empty/consent 422", "U4 input preserved", "U7 alert + describedby"])
def _():
    n0 = n_sent(); f = lambda b: re.search(r'<input id="friend-code"[^>]*>', b).group(0)
    s, h, b = follow(C.post("/home/friends/request", {"code": "", "consent": "yes"}), C); assert s == 422 and alerts(b) and 'value=""' in f(b) and 'aria-invalid="true"' in f(b)
    s, h, b = C.post("/home/friends/request", {"code": "   ", "consent": "yes"})[0:3]; assert s == 422 and 'value="   "' in f(b), f(b)
    s, h, b = C.post("/home/friends/request", {"code": "abc-xyz"}); assert s == 422 and 'value="abc-xyz"' in f(b) and "동의" in alerts(b)[0] and "checked" not in re.search(r'<input id="friend-consent"[^>]*>', b).group(0)
    assert re.search(r'<input id="friend-consent"[^>]*aria-describedby="friend-scope problem"[^>]*aria-invalid="true"', b), "consent checkbox must point at the error"
    assert n_sent() == n0, "invalid submissions created requests"


@check("odd_friend_codes_get_the_same_generic_success_as_a_missing_code_and_are_never_echoed", ["U4 no existence oracle", "U4 safe input"])
def _():
    ref = follow(C.post("/home/friends/request", {"code": "Zq" + "6" * 41, "consent": "yes"}), C); sig = lambda r: (r[0], [H.unescape(text(x)) for x in re.findall(r'role="status"[^>]*>(.*?)</(?:p|div|section)>', main(r[2]), re.S)])
    assert sig(ref)[0] == 200 and sig(ref)[1], sig(ref); out = {}
    for lab, v in (("nul", "ab\x00cd"), ("bel", "ab\x07cd"), ("bidi_rlo", "ab\u202ecd"), ("newline", "ab\ncd"), ("long500", "q" * 500), ("long2000", "q" * 2000), ("emoji", "😀" * 60)):
        r = follow(C.post("/home/friends/request", {"code": v, "consent": "yes"}), C); out[lab] = [r[0], sig(r)[1]]
        assert sig(r) == sig(ref), (lab, sig(r)); assert not INVIS(main(r[2])) and v not in r[2], lab
    EVID["odd_code_results"] = out


@check("sender_with_many_pending_requests_still_gets_the_same_answer_for_new_missing_codes", ["U4 sender cap not an oracle"])
def _():
    seq = []
    for i in range(14):
        r = follow(C.post("/home/friends/request", {"code": "Cap%03d" % i + "x" * 36, "consent": "yes"}), C)
        seq.append([r[0], [H.unescape(text(x)) for x in re.findall(r'role="status"[^>]*>(.*?)</(?:p|div|section)>', main(r[2]), re.S)][:1], len(alerts(r[2]))])
    EVID["sender_cap_sequence_distinct"] = sorted({json.dumps(x, ensure_ascii=False) for x in seq}); EVID["sent_rows_now"] = n_sent()
    assert len({json.dumps(x, ensure_ascii=False) for x in seq}) == 1, EVID["sender_cap_sequence_distinct"]


@check("rename_with_unsafe_name_is_refused_without_change_or_leak", ["U2 control chars", "U2 40 limit", "U2 foreign agent"])
def _():
    before = agents(); assert [a for _, a in before] == ["agent_a1", "agent_a2"]
    msg = {}
    for lab, nm in BAD_NAMES.items():
        s, h, b = A.post("/home/agent-name", {"agent": "agent_a2", "name": nm}); assert s in (400, 422), (lab, s)
        msg[lab] = alerts(b)[:1]; assert agents() == before, ("changed by", lab)
    EVID["rename_error_messages"] = msg
    s, h, b = B.post("/home/agent-name", {"agent": "agent_a2", "name": "해킹 시도"}); EVID["foreign_rename_by_B"] = s; assert s in (400, 403, 404, 422)
    s2, h2, b2 = A.post("/home/agent-name", {"agent": "agent_b1", "name": "남의 에이전트"}); EVID["foreign_rename_by_A"] = s2; assert s2 in (400, 403, 404, 422)
    assert agents() == before and "해킹" not in A.get("/home/agents")[2] and "남의" not in B.get("/home/friends/" + AB)[2]


@check("rejected_name_error_explains_control_and_direction_characters_and_is_not_echoed_invisibly", ["U2 control chars", "U7 error text", "U6-style safe echo"])
def _():
    fails = []
    for lab in ("bel", "nul", "bidi_rlo", "bidi_isolate", "newline"):
        s, h, b = A.post("/home/agent-name", {"agent": "agent_a2", "name": BAD_NAMES[lab]})
        inp = re.search(r'<input id="agent-name-agent_a2"[^>]*>', b).group(0); msg = " ".join(alerts(b))
        if INVIS(inp): fails.append(f"{lab}: invisible character echoed into the edit field")
        if not re.search(r"제어|방향|사용할 수 없|허용", msg): fails.append(f"{lab}: error says only {msg!r}")
    EVID["echo_failures"] = fails; assert not fails, fails


if __name__ == "__main__":
    EVID["member_window_429_retries"] = Browser.waited429
    out = {"origin": ORIGIN, "fixed_sha": "5fa3333b5948fd36792662b202915bc6243cac14", "evidence": EVID, "results": RES}
    json.dump(out, open(os.path.join(os.path.dirname(os.path.abspath(__file__)), "followup-result-5fa3333.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
    for r in RES: print(r["result"].upper(), r["check"], r.get("detail", "")[:500])
    sys.exit(0 if all(r["result"] == "pass" for r in RES) else 1)
