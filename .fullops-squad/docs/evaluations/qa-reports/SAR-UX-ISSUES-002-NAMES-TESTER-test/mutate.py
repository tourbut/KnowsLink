"""Stage 2 of 2: state-changing checks on the isolated names fixture (synthetic A/B/C only). IRREVERSIBLE; run once, after checks.py and the read-only browser stages.
Order: request validation -> request indistinguishability -> agent rename -> code rotation -> cancel/old code/new code -> deny -> stale decision -> accept (0-agent C)."""
import html as H, json, os, re, sys, traceback
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__))); from client import *

M = json.load(open(os.environ.get("NAMES_MANIFEST", "C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json")))
AB, CA = M["friend_ab"], M["request_ca"]
A_NAME, B_AGENT, A_AGENT = "가" * 40, "😀" * 40, "긴한글" * 13
RES = []; EVID = {}


def check(name, covers):
    def deco(fn):
        try: fn(); RES.append({"check": name, "covers": covers, "result": "pass"})
        except Exception as e: RES.append({"check": name, "covers": covers, "result": "fail", "detail": f"{type(e).__name__}: {e}", "trace": traceback.format_exc(limit=3)})
        return fn
    return deco


A, B, C = as_("a"), as_("b"), as_("c")
tt = lambda b: H.unescape(text(main(b)))
codeof = lambda br: re.search(r'<textarea id="my-code"[^>]*>([^<]*)</textarea>', br.get("/home/account")[2]).group(1)
statuses = lambda b: [H.unescape(text(x)) for x in re.findall(r'role="status"[^>]*>(.*?)</(?:p|div|section)>', main(b), re.S)]
alerts = lambda b: [H.unescape(text(x)) for x in re.findall(r'role="alert"[^>]*>(.*?)</(?:p|div)>', main(b), re.S)]
def received(br):  # [(name, friend, generation, deadline)]
    b = main(br.get("/home/friends")[2])
    return [(H.unescape(n), f, g) for n, f, g in re.findall(r'<li><h3>(.*?)의 친구 요청</h3>.*?name="friend" value="([^"]+)"><input type="hidden" name="generation" value="(\d+)"', b, re.S)]
def sent(br):
    b = main(br.get("/home/friends")[2])
    return re.findall(r'<li><h3>보낸 친구 요청</h3>.*?name="friend" value="([^"]+)"><input type="hidden" name="generation" value="(\d+)">', b, re.S)
def deadlines(br):
    return re.findall(r"요청 기한: (\d{4}-\d\d-\d\d \d\d:\d\d:\d\d KST)", main(br.get("/home/friends")[2]))
def follow(r, br):
    s, h, b = r
    if s in (301, 302, 303) and h.get("location"): return br.get(h["location"])
    return r
def request(br, code, consent=True):
    f = {"code": code}
    if consent: f["consent"] = "yes"
    return follow(br.post("/home/friends/request", f), br)
def decide(br, friend, gen, d):
    return follow(br.post("/home/friends/decision", {"friend": friend, "generation": gen, "decision": d, "consent": "yes"}), br)
def rename(br, agent, name): return follow(br.post("/home/agent-name", {"agent": agent, "name": name}), br)
def agents_of(br):
    return re.findall(r'<h2 class="name-preview">([^<]*)</h2><p class="muted">[^<]*/ (agent_\w+) /', br.get("/home/agents")[2])


CODE = {w: codeof(br) for w, br in (("a", A), ("b", B), ("c", C))}; EVID["code_lengths"] = {w: len(c) for w, c in CODE.items()}
C0 = {"sent": sent(C), "dl": deadlines(C), "a_recv": received(A), "a_dl": deadlines(A)}
EVID["before"] = {k: v for k, v in C0.items() if k != "a_recv"}; EVID["before"]["a_recv"] = [r[0] for r in C0["a_recv"]]


# ---------- M1: validation failures keep input, create nothing
@check("request_validation_422_keeps_input_and_creates_nothing", ["U4 empty/consent 422", "U4 input preserved", "U7 alert+describedby"])
def _():
    n0 = len(sent(C))
    for lab, code, consent in (("empty", "", True), ("spaces", "   ", True)):
        s, h, b = request(C, code, consent); assert s == 422, (lab, s)
        assert any("친구 코드를 입력하세요" in a for a in alerts(b)), (lab, alerts(b)); assert 'id="friend-code"' in b and 'aria-invalid="true"' in b and 'value=""' in b.split('id="friend-code"')[1][:300], lab
        assert "요청을 보냈습니다" not in tt(b)
    s, h, b = request(C, "abc-" + "x" * 30, consent=False); assert s == 422, s
    fc = re.search(r'<input id="friend-code"[^>]*>', b).group(0); assert 'value="abc-' + "x" * 30 + '"' in fc, fc  # typed code survives the consent error
    assert alerts(b) and "동의" in alerts(b)[0], alerts(b)
    cons = re.search(r'<input id="friend-consent"[^>]*>', b).group(0); assert "checked" not in cons and ("aria-describedby=\"friend-scope" in cons)
    assert re.search(r'aria-describedby="[^"]*(problem|error)[^"]*"', cons), "consent checkbox should point at the error: " + cons  # DEV low fix 3
    for lab, code in (("long500", "q" * 500), ("nul", "ab\x00cd"), ("bidi", "ab\u202ecd"), ("newline", "ab\ncd")):
        s, h, b = request(C, code); assert s in (400, 422), (lab, s)
        assert "\x00" not in b and "\u202e" not in b and ("q" * 500) not in b, lab
    assert len(sent(C)) == n0, "validation failures must not create requests"


# ---------- M2: indistinguishable submit
@check("request_submit_same_status_and_text_for_missing_valid_padded_cased_and_existing", ["U4 same copy", "U4 no existence oracle"])
def _():
    a_recv0, a_dl0, c_dl0 = [r[1:] for r in received(A)], deadlines(A), deadlines(C)
    rnd = "Zq" + "9" * 41
    cases = [("missing", C, rnd), ("valid_dup_pending", C, CODE["a"]), ("valid_padded", C, "  " + CODE["a"] + "\t"), ("valid_lowercase_variant", C, CODE["a"].swapcase()), ("valid_new_B", C, CODE["b"]),
             ("friend_existing_B_to_A", B, CODE["a"]), ("missing_from_B", B, rnd[::-1])]
    out = {}
    for lab, br, code in cases:
        s, h, b = request(br, code); out[lab] = (s, statuses(b), alerts(b), h.get("cache-control"))
    EVID["submit"] = out
    ref = out["missing"]; assert ref[0] == 200 and ref[1] and "요청을 보냈습니다. 가능한 요청은 상대가 확인할 수 있습니다." in ref[1][0] and not ref[2], ref
    for lab, v in out.items(): assert v == ref, (lab, v, ref)
    # effects: a duplicate to A must not renew the deadline or add a second request at A; existing friend adds nothing
    assert [r[1:] for r in received(A)] == a_recv0 and deadlines(A) == a_dl0, "A's received request changed"
    assert all(d in deadlines(C) for d in c_dl0), "repeat request moved C's earlier deadline"
    assert sum(1 for r in received(A) if r[0].startswith("합성 회원 c")) == 1
    assert any(r[0].startswith("합성 회원 c") for r in received(B)), "valid code to B should create B's received request"
    assert len(received(A)) == 1, "no extra request may appear at A"


@check("sent_pending_view_reveals_nothing_about_the_receiver_even_for_valid_code", ["U4 sent pending"])
def _():
    b = main(C.get("/home/friends")[2]); sec = b.split("친구 목록")[0]
    assert "보낸 친구 요청" in sec and "요청 취소" in sec
    for leak in (A_NAME, "합성 회원 b", A_AGENT, "agent_a", B_AGENT): assert leak not in sec, leak
    items = re.findall(r'<li><h3>보낸 친구 요청</h3>(.*?)</li>', sec, re.S); assert len(items) >= 3, len(items)
    shape = {re.sub(r"[0-9A-Za-z_-]{22}|\d{4}-\d\d-\d\d \d\d:\d\d:\d\d", "<X>", re.sub(r"\s+", " ", i)) for i in items}
    assert len(shape) == 1, "valid and missing codes must render the same sent row shape: %r" % shape


# ---------- M3: rename
@check("agent_rename_rejects_unsafe_names_without_change", ["U2 control chars", "U2 40 limit"])
def _():
    before = agents_of(A); assert [a for _, a in before] == ["agent_a1", "agent_a2"], before
    bad = {"empty": "", "spaces": "   ", "bel": "ab\x07cd", "nul": "ab\x00cd", "bidi": "ab\u202ecd", "isolate": "ab\u2066cd", "newline": "ab\ncd", "len41": "가" * 41, "emoji41": "😀" * 41}
    for lab, nm in bad.items():
        s, h, b = rename(A, "agent_a2", nm); assert s in (400, 422), (lab, s)
        assert "\x00" not in b and "\u202e" not in b and "\x07" not in b, lab
        assert agents_of(A) == before, ("changed by", lab)
    s, h, b = rename(B, "agent_a2", "해킹 시도"); assert s in (400, 403, 404, 422), s
    assert agents_of(A) == before and "해킹" not in A.get("/home/agents")[2] and "해킹" not in B.get("/home/friends/" + AB)[2]
    s, h, b = rename(A, "agent_b1", "남의 에이전트"); assert s in (400, 403, 404, 422), s
    assert "남의" not in B.get("/home/friends/" + AB)[2] and B_AGENT in A.get("/home/friends/" + AB)[2]


NEW = "변경후-" + "나" * 29 + "🙂" * 7


@check("agent_rename_valid_40_chars_keeps_id_links_relationship_and_creates_no_agent", ["U2 rename keeps ID", "U2 no duplicate"])
def _():
    assert len(NEW) == 40
    for _ in range(2):  # a retry of the same request must not create or duplicate anything
        s, h, b = rename(A, "agent_a2", NEW); assert s == 200 or s in (303,), s
    ag = agents_of(A); assert [(H.unescape(n), i) for n, i in ag] == [(A_AGENT, "agent_a1"), (NEW, "agent_a2")], ag
    ab = A.get("/home/agents")[2]
    assert 'href="/home/receipts?agent=agent_a2#receipt-form"' in ab and 'name="agent" value="agent_a2"' in ab and f"aria-label=\"에이전트 이름 저장: {NEW} / agent_a2\"" in H.unescape(ab)
    rb = A.get("/home/receipts")[2]; assert f"<span>{NEW} / Node 22 로컬 클라이언트 / agent_a2</span>" in rb and rb.count('type="radio"') == 2
    lb = A.get("/home/legacy")[2]; assert f"<h3>{NEW} / Node 22 로컬 클라이언트 / agent_a2</h3>" in lb and f"<h3>{A_AGENT} / Node 22 로컬 클라이언트 / agent_a1</h3>" in lb
    nb = A.get("/home/new?peer=self")[2]; assert f">{NEW} / Node 22 로컬 클라이언트 / agent_a2<" in nb
    bf = B.get("/home/friends/" + AB); assert bf[0] == 200 and f"<h2>{NEW}</h2><p class=\"muted\">Node 22 로컬 클라이언트 / agent_a2 /" in bf[2] and "agent_a1" in bf[2]
    assert "합성 회원 b" in A.get("/home/friends")[2] and A.get("/home/friends/" + AB)[0] == 200, "friend relationship must be unchanged"
    EVID["rename"] = {"new_len": len(NEW), "agents_after": ag}


@check("agent_rename_back_restores_original_labels", ["fixture hygiene"])
def _():
    rename(A, "agent_a2", A_AGENT); assert [(H.unescape(n), i) for n, i in agents_of(A)] == [(A_AGENT, "agent_a1"), (A_AGENT, "agent_a2")]


# ---------- M4: rotate / cancel / old code / new code / deny / stale / accept
OLD = CODE["a"]; STATE = {}


@check("code_rotation_replaces_code_keeps_friends_and_pending_received_requests", ["U2 rotation", "U4 retained received"])
def _():
    recv0 = [r[1:] for r in received(A)]; assert recv0, "C's pending request should still exist at A"
    s, h, b = follow(A.post("/home/friend-code"), A); assert s == 200, s
    new = codeof(A); STATE["new"] = new
    assert new != OLD and re.fullmatch(r"[A-Za-z0-9_-]{43}", new), (OLD[:6], new[:6])
    for p in ("/home/friends", "/home/account", "/home/agents", "/home/legacy"): assert OLD not in A.get(p)[2], p
    assert [r[1:] for r in received(A)] == recv0, "received request must survive rotation"
    assert "합성 회원 b" in A.get("/home/friends")[2] and A.get("/home/friends/" + AB)[0] == 200
    assert codeof(B) == CODE["b"] and codeof(C) == CODE["c"], "other members' codes must not change"
    EVID["rotation"] = {"old_prefix": OLD[:4], "new_prefix": new[:4]}


@check("old_code_after_cancel_creates_no_request_but_looks_identical_and_new_code_does", ["U2 old code unusable", "U4 same copy", "U4 cancel"])
def _():
    (fid, gen), = [x for x in sent(C) if x[0] == CA] or [(None, None)]
    assert fid, "C's pending to A must still exist before cancel"
    s, h, b = decide(C, fid, gen, "cancel"); assert s == 200, s
    assert CA not in main(C.get("/home/friends")[2]) and not received(A), "cancel must remove A's received request"
    ref = request(C, "Zq" + "8" * 41)
    s, h, b = request(C, OLD); assert (s, statuses(b), alerts(b)) == (ref[0], statuses(ref[2]), alerts(ref[2])), "old code must look like any missing code"
    assert not received(A), "old code must not create a request"
    s2, h2, b2 = request(C, STATE["new"]); assert (s2, statuses(b2)) == (ref[0], statuses(ref[2]))
    r = received(A); assert len(r) == 1 and r[0][0].startswith("합성 회원 c"), r
    STATE["req"] = r[0][1:]; EVID["new_request"] = {"friend_changed": r[0][1] != CA, "generation": r[0][2]}


@check("deny_shows_requester_record_and_new_request_is_required", ["U4 deny", "DEV judgement item 1"])
def _():
    fid, gen = STATE["req"]; s, h, b = decide(A, fid, gen, "deny"); assert s == 200, s
    assert not received(A)
    cb = main(C.get("/home/friends")[2]); sec = cb.split("친구 목록")[0]
    EVID["requester_view_after_deny"] = tt(C.get("/home/friends")[2]).split("친구 목록")[0][-600:]
    EVID["requester_sees_declined_word"] = "거절" in sec
    assert "친구 목록" in cb and "합성 회원 b" not in sec
    s2, h2, b2 = decide(A, fid, gen, "accept"); EVID["stale_accept_status"] = s2
    assert s2 in (400, 403, 404, 409, 422), "a stale-generation decision must be refused, got %s" % s2
    assert "합성 회원 c" not in main(A.get("/home/friends")[2]).split("친구 추가")[0].split("친구 목록")[1], "stale accept must not create a friendship"


@check("re_request_after_deny_creates_new_generation_and_duplicate_accept_is_idempotent", ["U4 re-request", "U4 duplicate accept"])
def _():
    s, h, b = request(C, STATE["new"]); assert s == 200 and "요청을 보냈습니다" in "".join(statuses(b)), (s, statuses(b), alerts(b))
    r = received(A); assert len(r) == 1, r; fid, gen = r[0][1:]; EVID["second_request_generation"] = gen
    assert gen != STATE["req"][1] or fid != STATE["req"][0], "a re-request after deny must be a new generation"
    s1, _, b1 = decide(A, fid, gen, "accept"); assert s1 == 200, s1
    s2, _, b2 = decide(A, fid, gen, "accept"); EVID["second_accept_status"] = s2
    af = main(A.get("/home/friends")[2]); lst = af.split("친구 목록")[1].split("친구 추가")[0]
    assert lst.count("합성 회원 c") == 1 and lst.count("합성 회원 b") == 1, lst
    assert not received(A)


@check("zero_agent_c_is_friend_without_empty_selector_and_with_add_agent_guidance", ["U4 zero agents after accept", "U2 friend detail"])
def _():
    cf = C.get("/home/friends"); lst = main(cf[2]).split("친구 목록")[1].split("친구 추가")[0]
    m = re.search(r'href="/home/friends/([^"]+)">%s</a>' % A_NAME, lst); assert m, lst
    fid = m.group(1); s, h, b = C.get("/home/friends/" + fid); assert s == 200, s
    mb = main(b); EVID["c_friend_detail_text"] = tt(b)[:500]
    assert f"<h1>{A_NAME}</h1>" in mb and "<select" not in mb
    assert f"<h2>{A_AGENT}</h2>" in mb and mb.count("agent_a1") >= 1 and mb.count("agent_a2") >= 1
    nb = C.get("/home/new?peer=" + fid)[2]; assert "<select" not in nb and "에이전트 추가" in nb and "/home/agents" in nb, tt(nb)
    # A sees C as a friend that has no agent: no conversation action and no empty list of agents
    s, h, b = A.get("/home/friends"); m2 = re.search(r'href="/home/friends/([^"]+)">합성 회원 c</a>', b); assert m2
    ab = main(A.get("/home/friends/" + m2.group(1))[2]); EVID["a_view_of_c_text"] = text(ab)[:400]
    assert "<h1>합성 회원 c</h1>" in ab and "<select" not in ab and "href=\"/home/new?peer=" not in ab, text(ab)
    EVID["friend_ids"] = {"ac": fid, "ca": m2.group(1)}


@check("final_relationship_set_is_exactly_ab_and_ac_and_b_has_pending_from_c", ["U4 GET/side-effect final state"])
def _():
    names = lambda br: re.findall(r'class="person-name name-preview" href="/home/friends/[^"]+">([^<]*)</a>', br.get("/home/friends")[2])
    assert sorted(names(A)) == sorted([H.unescape("합성 회원 b"), "합성 회원 c"]), names(A)
    assert names(B) == [A_NAME], names(B)
    assert any(r[0].startswith("합성 회원 c") for r in received(B)), "B's pending request from C should remain"
    assert names(C) == [A_NAME]


if __name__ == "__main__":
    EVID["member_window_429_retries"] = Browser.waited429
    out = {"origin": ORIGIN, "fixed_sha": "b2e683ae677c78cd7bdc47c75820bcb191b94838", "evidence": EVID, "results": RES}
    json.dump(out, open(os.path.join(os.path.dirname(os.path.abspath(__file__)), "mutate-result.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
    for r in RES: print(r["result"].upper(), r["check"], r.get("detail", ""))
    sys.exit(0 if all(r["result"] == "pass" for r in RES) else 1)
