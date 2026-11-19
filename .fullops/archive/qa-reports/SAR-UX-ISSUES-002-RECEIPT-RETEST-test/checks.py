"""Independent HTTP checks for U6/U7 receipt lookup against the isolated fixture (synthetic A/B/C only)."""
import base64, http.client, json, re, sys, traceback, urllib.parse
sys.path.insert(0, "."); from client import *

AB = "0199a3f2-4c10-7a11-8b22-334455667788"; BC = "0199a3f2-4c10-7a11-8b22-3344556677aa"; NONE = "0199a3f2-4c10-7a11-8b22-334455667799"
REFUSE = "이 요청의 상태를 확인할 수 없습니다. 내 에이전트와 요청 ID를 확인하세요."
RES = []; EVID = {}


def check(name, covers):
    def deco(fn):
        try:
            fn(); RES.append({"check": name, "covers": covers, "result": "pass"})
        except Exception as e:
            RES.append({"check": name, "covers": covers, "result": "fail", "detail": f"{type(e).__name__}: {e}", "trace": traceback.format_exc(limit=3)})
        return fn
    return deco


def main(b): return re.search(r"<main.*</main>", strip_style(b), re.S).group(0)
def val_id(b): return re.search(r'<input id="receipt-id" name="id" value="([^"]*)"', b).group(1)
def sel_agent(b):
    m = re.search(r'<option value="([^"]+)" selected>', b); return m.group(1) if m else ""
def options(b): return re.findall(r'<option value="([^"]*)"', b)
def cookie_obj(br): return json.loads(base64.urlsafe_b64decode(br.cookies["__Host-kl_receipt"] + "=="))
def setck(br, o): br.cookies["__Host-kl_receipt"] = base64.urlsafe_b64encode(json.dumps(o, separators=(",", ":")).encode()).decode().rstrip("=")
def norm(b, i): return main(b).replace(i, "<ID>")
def new(who):
    br = Browser(); assert br.post("/fixture/recover")[0] == 204  # clear the fixture's 40/member lookup window so volume of this suite is not the variable
    assert br.post(f"/fixture/as/{who}")[0] == 204; return br


# ---- A: first, success, refusals, invalid forms
A = new("a")
LIM = 38  # keep A's own session under the 40-lookup window between recovers
s, h, first = A.get("/home/receipts")


@check("first_view_form_and_selector", ["U6 first", "U2 name"])
def _():
    b = first
    assert s == 200 and "<title>요청 상태 조회" in b
    m = main(b)
    for t in ("내 에이전트</label>", "요청 ID 또는 답장 ID</label>", "상태 확인</button>", "현재 온라인 상태를 확인하는 기능은 아닙니다"): assert t in m, t
    assert options(b) == ["", "agent_a", "agent_a2"], options(b)
    assert "agent_b" not in b and "agent_c" not in b
    opt = re.findall(r"<option value=\"agent_a\"[^>]*>([^<]*)</option>", b)[0]
    assert opt.count(" / ") == 2 and opt.endswith("agent_a") and "Node 22" in opt, opt
    assert m.count("<h1>") == 1 and 'class="skip" href="#main"' in b and '<main id="main"' in b
    assert m.count('aria-current="page"') == 1
    assert re.search(r'<form id="receipt-form" method="get"', m) and 'method="post"' not in m
    assert len(re.findall(r"\bautofocus\b", b)) == 1


s, h, ok = A.get(q(agent="agent_a", id=AB))


@check("success_splits_delivery_processing_reply_and_keeps_input", ["U6 success"])
def _():
    assert s == 200, s
    m = main(ok)
    for t in ("<dt>전달</dt>", "<dt>처리</dt>", "<dt>관련 답장</dt>", "<summary>기술 상세</summary>"): assert t in m, t
    assert "새 key로 송신" not in ok
    assert sel_agent(ok) == "agent_a" and val_id(ok) == AB
    assert 'role="status"' in m and "<code>%s</code>" % AB in m
    assert 'href="#receipt-id"' in m
    assert 'action="/home/receipts"' in m and 'method="post"' not in m  # lookup never POSTs


@check("repeat_and_refresh_are_idempotent_read_only", ["U6 refresh", "no send"])
def _():
    s2, _, ok2 = A.get(q(agent="agent_a", id=AB)); s3, _, ok3 = A.get(q(agent="agent_a", id=AB))
    d = lambda x: re.sub(r"\d{4}-\d\d-\d\d \d\d:\d\d:\d\d KST", "<T>", main(x))
    assert s2 == s3 == 200 and d(ok2) == d(ok3) == d(ok)


REF = {}
for lab, agent, i in [("none", "agent_a", NONE), ("bc_other", "agent_a", BC), ("ab_via_nonparticipant_own_agent", "agent_a2", AB), ("foreign_agent_name", "agent_b", AB), ("foreign_agent_c", "agent_c", BC)]:
    REF[lab] = A.get(q(agent=agent, id=i))


@check("refusals_are_indistinguishable_and_keep_owned_selection", ["U6 none/other/nonparticipant"])
def _():
    for lab in ("none", "bc_other", "ab_via_nonparticipant_own_agent"):
        st, hd, bd = REF[lab]
        assert st == 403 and REFUSE in bd and 'role="alert"' in bd, (lab, st)
        assert "agent_c" not in bd and "agent_b" not in bd, lab
    assert sel_agent(REF["none"][2]) == "agent_a" and sel_agent(REF["bc_other"][2]) == "agent_a"
    assert sel_agent(REF["ab_via_nonparticipant_own_agent"][2]) == "agent_a2"
    assert norm(REF["none"][2], NONE) == norm(REF["bc_other"][2], BC)
    pick = lambda r: {k: v for k, v in r[1].items() if k in ("cache-control", "content-security-policy")}
    assert pick(REF["none"]) == pick(REF["bc_other"])


@check("foreign_agent_in_url_is_not_reflected", ["U6 other account", "account isolation"])
def _():
    for lab in ("foreign_agent_name", "foreign_agent_c"):
        st, hd, bd = REF[lab]
        assert st == 403 and REFUSE in bd
        assert sel_agent(bd) == "" and val_id(bd) == "", lab
        assert "agent_b" not in bd and "agent_c" not in bd and AB not in bd and BC not in bd


bad_x = A.get(q(agent="agent_a", id='nope"><script>alert(1)</script>'))


@check("invalid_format_is_escaped_and_editable_in_place", ["U6 invalid format"])
def _():
    st, hd, bd = bad_x
    assert st == 422 and "요청 ID 형식이 올바르지 않습니다" in bd
    assert "<script>alert" not in bd and "&lt;script&gt;" in bd
    assert sel_agent(bd) == "agent_a" and 'aria-invalid="true"' in bd and 'role="alert"' in bd
    assert 'href="#receipt-id"' in bd and "receipt-problem" in re.search(r'<input id="receipt-id"[^>]*>', bd).group(0)


@check("oversized_or_control_input_is_dropped_selection_kept", ["U6 safe length", "escape"])
def _():
    cases = {"ascii200": "a" * 200, "kr43_129bytes": "가" * 43, "nul": "ab\x00cd", "newline": "ab\ncd", "bidi_rlo": "ab\u202ecd", "bidi_isolate": "ab\u2066cd"}
    for lab, v in cases.items():
        st, hd, bd = A.get(q(agent="agent_a", id=v))
        assert st in (400, 422), (lab, st)
        assert val_id(bd) == "" and sel_agent(bd) == "agent_a", (lab, val_id(bd)[:20], sel_agent(bd))
        assert v not in bd and "\u202e" not in bd and "\x00" not in bd, lab


@check("empty_id_keeps_selection", ["U6 empty"])
def _():
    st, hd, bd = A.get(q(agent="agent_a", id=""))
    assert st in (400, 422) and sel_agent(bd) == "agent_a" and 'role="alert"' in bd, st


@check("return_and_open_redirect_params_are_ignored", ["U7 safe return"])
def _():
    for k in ("return", "return_to", "next", "redirect", "continue", "url"):
        for ev in ("https://evil.example/x", "//evil.example", "javascript:alert(1)"):
            st, hd, bd = A.get(f"/home/receipts?{k}={urllib.parse.quote(ev)}&agent=agent_a&id={AB}")
            assert st == 200 and "location" not in hd and "evil.example" not in bd and "javascript:" not in bd, (k, ev)
    hrefs = set(re.findall(r'href="([^"]+)"', strip_style(first)))
    assert all(x.startswith("/") or x.startswith("#") for x in hrefs), hrefs


# ---- cookie
@check("receipt_cookie_is_minimal_and_hardened", ["recovery context"])
def _():
    A2 = new("a")
    conn = http.client.HTTPConnection(HOST, PORT)
    conn.request("GET", q(agent="agent_a", id=AB), headers={"Cookie": "; ".join(f"{k}={v}" for k, v in A2.cookies.items())})
    r = conn.getresponse(); r.read()
    sc = [v for k, v in r.getheaders() if k.lower() == "set-cookie" and v.startswith("__Host-kl_receipt=")][0]
    for t in ("Secure", "HttpOnly", "SameSite=Strict", "Path=/"): assert t in sc, sc
    assert "Domain" not in sc and int(re.search(r"Max-Age=(\d+)", sc).group(1)) <= 4200
    A2.get(q(agent="agent_a", id=AB))
    o = cookie_obj(A2); assert set(o) == {"Member", "Agent", "ID", "Exp"}, o
    raw = base64.urlsafe_b64decode(A2.cookies["__Host-kl_receipt"] + "==").decode()
    assert A2.cookies["__Host-kl_session"] not in raw and "@" not in raw and "전달" not in raw and "긴한글" not in raw
    assert o["Agent"] == "agent_a" and o["ID"] == AB


# ---- restore
@check("bare_lookup_restores_last_own_input_same_session", ["U6 back/restore"])
def _():
    A3 = new("a"); A3.get(q(agent="agent_a", id=AB))
    st, hd, bd = A3.get("/home/receipts"); assert st == 200 and sel_agent(bd) == "agent_a" and val_id(bd) == AB, (sel_agent(bd), val_id(bd))
    A3.get(q(agent="agent_a2", id=NONE)); st, hd, bd = A3.get("/home/receipts")
    assert sel_agent(bd) == "agent_a2" and val_id(bd) == NONE


@check("account_B_and_anonymous_never_see_A_content", ["isolation"])
def _():
    A4 = new("a"); A4.get(q(agent="agent_a", id=AB)); ck = dict(A4.cookies)
    A4.post("/fixture/as/b"); st, hd, bd = A4.get("/home/receipts")
    assert "__Host-kl_receipt" in A4.cookies
    assert st == 200 and options(bd) == ["", "agent_b"], options(bd)
    assert val_id(bd) == "" and sel_agent(bd) == "" and AB not in bd and "agent_a" not in bd
    st, hd, bd = A4.get(q(agent="agent_a", id=AB))
    assert st == 403 and "agent_a" not in bd and sel_agent(bd) == "" and val_id(bd) == ""
    anon = Browser(); anon.cookies["__Host-kl_receipt"] = ck["__Host-kl_receipt"]
    for p in ("/home/receipts", q(agent="agent_a", id=AB)):
        st, hd, bd = anon.get(p)
        assert AB not in bd and "agent_a" not in bd and "Node 22" not in bd and "긴한글" not in bd, p
        assert st in (200, 302, 303, 401), st


@check("tampered_recovery_cookie_does_not_cross_accounts_or_agents", ["tamper"])
def _():
    A5 = new("a"); A5.get(q(agent="agent_a", id=AB)); good = cookie_obj(A5)
    B5 = new("b"); B5.get(q(agent="agent_b", id=AB)); bo = cookie_obj(B5)
    variants = {
        "member_swapped_to_B": dict(good, Member=bo["Member"]),
        "agent_swapped_to_B_agent": dict(good, Agent="agent_b"),
        "agent_C": dict(good, Agent="agent_c"),
        "id_overlong": dict(good, ID="a" * 500),
        "id_html": dict(good, ID='"><img src=x onerror=1>'),
        "id_bidi": dict(good, ID="ab\u202ecd"),
        "expired": dict(good, Exp="2020-01-01T00:00:00Z"),
        "extra_fields": dict(good, Result="delivered", Email="x@y.z"),
    }
    for lab, o in variants.items():
        c = Browser(); c.cookies.update(A5.cookies); setck(c, o)
        st, hd, bd = c.get("/home/receipts")
        # an unowned Agent in the cookie is answered like a foreign agent in the URL (403 refusal, once), then the cookie is reset
        # an unowned Agent in the cookie is answered like a foreign agent in the URL (403 refusal, once), then the cookie is reset
        # a syntactically unsafe-but-short ID is a plain format error (escaped, selection kept); unsafe length/control/expired are dropped silently
        want = {"agent_swapped_to_B_agent": 403, "agent_C": 403, "id_html": 422}.get(lab, 200)
        assert st == want, (lab, st, want)
        if st == 403: assert REFUSE in bd and c.get("/home/receipts")[0] == 200, lab
        assert "onerror=1>" not in bd and "<img src=x" not in bd and "‮" not in bd, lab
        # AB is A's own participant receipt, so agent_b legitimately appears only in the restored success detail
        assert "x@y.z" not in bd and "Result" not in bd and ("agent_b" not in bd and "agent_c" not in bd or lab == "extra_fields"), lab
        if lab == "id_html": assert "&#34;&gt;&lt;img" in bd and sel_agent(bd) == "agent_a"
        elif lab in ("agent_swapped_to_B_agent", "agent_C"): assert val_id(bd) == "" and sel_agent(bd) == ""
        elif lab in ("member_swapped_to_B", "expired", "id_overlong", "id_bidi"): assert val_id(bd) == "" and sel_agent(bd) == "", (lab, sel_agent(bd), val_id(bd))
        else: assert sel_agent(bd) == "agent_a" and val_id(bd) == AB  # unknown fields ignored, not echoed
    for lab, raw in {"garbage": "!!!!", "not_json": base64.urlsafe_b64encode(b"hello").decode(), "huge": "A" * 6000}.items():
        c = Browser(); c.cookies.update(A5.cookies); c.cookies["__Host-kl_receipt"] = raw
        st, hd, bd = c.get("/home/receipts"); assert st == 200 and val_id(bd) == "", (lab, st)


# ---- 429 / 503 / expire (state-changing; recover each)
@check("429_keeps_input_shows_retry_and_manual_only", ["U6 429"])
def _():
    L = new("a"); L.get(q(agent="agent_a", id=AB)); assert L.post("/fixture/limit")[0] == 204
    st, hd, bd = L.get(q(agent="agent_a", id=AB))
    EVID["429"] = {"status": st, "retry-after": hd.get("retry-after"), "alert": re.findall(r'role="alert"[^>]*>([^<]*)', main(bd))}
    try:
        assert st == 429 and sel_agent(bd) == "agent_a" and val_id(bd) == AB, (st, sel_agent(bd), val_id(bd))
        assert hd.get("retry-after") or re.search(r"\d{4}-\d\d-\d\d|\d+\s*(초|분)", main(bd)), "no retry time"
        assert "refresh" not in hd and "http-equiv" not in bd and "<script" not in bd
        st, hd, bd = L.get(q(agent="agent_a", id=NONE)); assert st == 429 and REFUSE not in bd  # no existence oracle while limited
    finally:
        assert L.post("/fixture/recover")[0] == 204
    st, hd, bd = L.get(q(agent="agent_a", id=AB)); assert st == 200 and "<dt>전달</dt>" in bd


@check("503_outage_hides_previous_input_keeps_cookie_and_recovers", ["U6 503", "fail-closed"])
def _():
    O = new("a"); O.get(q(agent="agent_a", id=AB)); before = O.cookies["__Host-kl_receipt"]
    assert O.post("/fixture/outage")[0] == 204
    try:
        st, hd, bd = O.get("/home/receipts"); st2, hd2, bd2 = O.get(q(agent="agent_a", id=AB))
        EVID["503"] = {"status": [st, st2], "cache": hd.get("cache-control")}
        assert st == 503 and st2 == 503, (st, st2)
        for x in (bd, bd2): assert AB not in x and "agent_a" not in x and "Node 22" not in x and "긴한글" not in x
        assert "no-store" in hd.get("cache-control", "")
        assert O.cookies.get("__Host-kl_receipt") == before and "__Host-kl_session" in O.cookies
    finally:
        assert O.post("/fixture/recover")[0] == 204
    st, hd, bd = O.get("/home/receipts"); assert st == 200 and sel_agent(bd) == "agent_a" and val_id(bd) == AB


@check("session_expiry_then_same_member_restores_other_member_does_not", ["U6 session", "same-account re-login"])
def _():
    E = new("a"); E.get(q(agent="agent_a", id=AB)); E.post("/fixture/expire")
    st, hd, bd = E.get("/home/receipts"); EVID["expire"] = {"status": st, "location": hd.get("location")}
    assert AB not in bd and "agent_a" not in bd and "Node 22" not in bd and "긴한글" not in bd
    assert st in (200, 302, 303, 401) and ("location" in hd or "로그인" in bd or "login" in bd.lower()), (st, bd[:200])
    E.post("/fixture/as/a"); st, hd, bd = E.get("/home/receipts")
    assert st == 200 and sel_agent(bd) == "agent_a" and val_id(bd) == AB, (st, sel_agent(bd), val_id(bd))
    E.post("/fixture/expire"); E.post("/fixture/as/b"); st, hd, bd = E.get("/home/receipts")
    assert st == 200 and val_id(bd) == "" and sel_agent(bd) == "" and "agent_a" not in bd


# ---- revoke (last; irreversible for this fixture)
import os
RV = new("a"); RV.get(q(agent="agent_a", id=AB))


def revoke_check():
    assert RV.post("/fixture/revoke")[0] == 204
    st, hd, bd = RV.get(q(agent="agent_a", id=AB))
    EVID.setdefault("revoke", []).append({"step": "query revoked agent", "status": st, "options": options(bd), "id": val_id(bd), "alert": re.findall(r'role="alert"[^>]*>([^<]*)', bd)})
    # revoked agent is shown, disabled, labelled; never silently replaced by agent_a2
    assert options(bd) == ["", "agent_a2", "agent_a"], options(bd)
    assert re.search(r'<option value="agent_a" selected disabled>[^<]*/ 철회</option>', bd) and '<option value="agent_a2" selected' not in bd
    assert "선택한 내 에이전트가 철회되었습니다. 사용할 내 에이전트를 다시 선택하세요." in bd and 'aria-invalid="true"' in bd
    assert val_id(bd) == AB and st in (403, 422)
    st, hd, bd = RV.get("/home/receipts")  # recovery context still names the revoked agent: same notice, still no silent switch
    EVID["revoke"].append({"step": "bare after", "status": st, "options": options(bd), "selected": sel_agent(bd), "id": val_id(bd)})
    assert st in (200, 403) and '<option value="agent_a2" selected' not in bd and "철회되었습니다" in bd
    st, hd, bd = RV.get(q(agent="agent_a2", id=AB)); assert st == 403 and REFUSE in bd and sel_agent(bd) == "agent_a2"  # non-participant agent still refused


if os.environ.get("RUN_REVOKE") == "1":  # irreversible for this fixture: run last
    check("revoked_agent_is_not_silently_replaced", ["U6 revoke"])(revoke_check)

if __name__ == "__main__":
    out = {"origin": ORIGIN, "evidence": EVID, "results": RES}
    json.dump(out, open("http-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
    for r in RES: print(r["result"].upper(), r["check"], r.get("detail", ""))
    sys.exit(0 if all(r["result"] == "pass" for r in RES) else 1)
