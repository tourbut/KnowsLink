"""Independent read-only HTTP checks for U2/U4/U7 names UX against the isolated names fixture (synthetic A/B/C only).
Stage 1 of 2: nothing here changes fixture state. mutate.py runs afterwards (irreversible POSTs)."""
import json, os, re, sys, traceback, html as H
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__))); from client import *

M = json.load(open(os.environ.get("NAMES_MANIFEST", "C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json")))
AB, CA = M["friend_ab"], M["request_ca"]
A_NAME, B_AGENT, A_AGENT = "가" * 40, "😀" * 40, "긴한글" * 13
REFUSE = "이 요청의 상태를 확인할 수 없습니다. 내 에이전트와 요청 ID를 확인하세요."
RES = []; EVID = {}


def check(name, covers):
    def deco(fn):
        try: fn(); RES.append({"check": name, "covers": covers, "result": "pass"})
        except Exception as e: RES.append({"check": name, "covers": covers, "result": "fail", "detail": f"{type(e).__name__}: {e}", "trace": traceback.format_exc(limit=3)})
        return fn
    return deco


PAGES = ["/home", "/home/friends", "/home/account", "/home/agents", "/home/legacy", "/home/receipts", "/home/new"]
S = {w: as_(w) for w in "abc"}
P = {w: {p: S[w].get(p) for p in PAGES} for w in "abc"}
for w in "abc": P[w]["/home/friends/" + AB] = S[w].get("/home/friends/" + AB)
tt = lambda b: H.unescape(text(main(b)))
title = lambda b: H.unescape(re.search(r"<title>(.*?)</title>", b, re.S).group(1))
CODES = {w: re.search(r'<textarea id="my-code"[^>]*>([^<]*)</textarea>', P[w]["/home/account"][2]).group(1) for w in "abc"}
parts = lambda s: [x.strip() for x in H.unescape(s).split(" / ")]


@check("landmarks_title_h1_skip_main_current", ["U7 title/h1/main/skip", "U7 aria-current"])
def _():
    titles = {}
    for w in "abc":
        for p in PAGES:
            s, h, b = P[w][p]; assert s == 200, (w, p, s)
            assert re.search(r'<html lang="ko"', b) and 'name="viewport"' in b, (w, p)
            t = title(b); assert t.strip() and t.endswith("KnowsLink"), (p, t); titles.setdefault(p, set()).add(t)
            assert len(re.findall(r"<h1[ >]", b)) == 1, (w, p)
            assert re.search(r'<a class="skip" href="#main">', b) and re.search(r'<main id="main"', b), (w, p)
            assert 'aria-current="page"' in b or p == "/home/legacy", (w, p)  # legacy is one scrolling page with no nav
    assert all(len(v) == 1 for v in titles.values()), titles       # same page, same title per account
    assert len({next(iter(v)) for v in titles.values()}) == len(titles), titles  # distinct per page
    EVID["titles"] = {p: next(iter(v)) for p, v in titles.items()}


@check("legacy_title_differs_from_agents_and_current_nav_is_not_faked", ["U7 unique title", "U7 no fake current"])
def _():
    t = EVID["titles"]; assert t["/home/legacy"] != t["/home/agents"] and "기존 연결" in t["/home/legacy"]
    b = P["a"]["/home/legacy"][2]
    assert 'aria-current="page"' not in b, "legacy single page must not mark a nav item current by scroll guess"


@check("name_first_kind_discriminator_on_every_surface", ["U2 name first", "U2 kind", "U2 discriminator"])
def _():
    rb = P["a"]["/home/receipts"][2]
    radios = re.findall(r'<input type="radio" name="agent" value="([^"]+)"[^>]*><span>([^<]*)</span>', rb)
    assert [v for v, _ in radios] == ["agent_a1", "agent_a2"], radios
    for v, lab in radios:
        assert parts(lab) == [A_AGENT, "Node 22 로컬 클라이언트", v], lab
    assert len({lab for _, lab in radios}) == 2
    ab = P["a"]["/home/agents"][2]
    heads = re.findall(r'<h2 class="name-preview">([^<]*)</h2><p class="muted">([^<]*) / <span', ab)
    assert [(H.unescape(n), H.unescape(k)) for n, k in heads] == [(A_AGENT, "Node 22 로컬 클라이언트 / agent_a1"), (A_AGENT, "Node 22 로컬 클라이언트 / agent_a2")], heads
    lb = P["a"]["/home/legacy"][2]
    sums = re.findall(r"<summary><h3>([^<]*)</h3>", lb)
    assert [parts(x) for x in sums] == [[A_AGENT, "Node 22 로컬 클라이언트", "agent_a1"], [A_AGENT, "Node 22 로컬 클라이언트", "agent_a2"]], sums
    nb = S["a"].get("/home/new?peer=self")[2]
    for sel in ("from", "to"):
        opts = re.findall(r'<option value="([^"]+)"[^>]*>([^<]*)</option>', re.search(rf'<select id="{sel}".*?</select>', nb, re.S).group(0))
        assert len(opts) == 2 and len({l for _, l in opts}) == 2, opts
        for v, lab in opts: assert parts(lab) == [A_AGENT, "Node 22 로컬 클라이언트", v], lab
    fb = P["a"]["/home/friends/" + AB][2]
    assert re.search(r"<h2>%s</h2><p class=\"muted\">Node 22 로컬 클라이언트 / agent_b1 / " % B_AGENT, fb), "friend detail agent line"
    assert re.search(r'<option value="agent_b1">%s / Node 22 로컬 클라이언트 / agent_b1</option>' % B_AGENT, S["a"].get("/home/new?peer=" + AB)[2])
    # B sees A's two same-name agents on the friend detail, told apart by the discriminator
    bb = P["b"]["/home/friends/" + AB][2]
    assert len(re.findall(r"<h2>%s</h2><p class=\"muted\">Node 22 로컬 클라이언트 / agent_a[12] / " % A_AGENT, bb)) == 2 and "agent_a1" in bb and "agent_a2" in bb


@check("colliding_names_are_distinguishable_on_each_surface", ["U2 collision"])
def _():
    for p, pat in (("/home/receipts", r"<span>([^<]*)</span></label>"), ("/home/legacy", r"<summary><h3>([^<]*)</h3>"), ("/home/agents", r'aria-label="에이전트 상세: ([^"]*)"')):
        labs = [H.unescape(x) for x in re.findall(pat, P["a"][p][2])]
        assert len(labs) >= 2 and len(set(labs)) == len(labs), (p, labs)


@check("non_friend_agent_shows_unverifiable_label_and_only_allowed_id", ["U2 non-friend", "U2 no extra lookup"])
def _():
    lb = main(P["a"]["/home/legacy"][2]); tt_ = H.unescape(lb)
    assert "— 확인할 수 없는 에이전트" in tt_, "label"
    row = re.search(r'<li class="row"><p>(.*?)</p>', lb, re.S).group(1)
    assert "agent_legacy_x" not in H.unescape(row), "ID belongs in the details, not in the label"
    assert re.search(r"<summary>식별자</summary><p><code>agent_a1</code> / <code>agent_legacy_x</code></p>", lb)
    for w in "bc":  # a stranger / non-owner never sees A's agent names
        for p in PAGES:
            allowed = w == "b" and p in ("/home/new", "/home/friends/" + AB)
            if p == "/home/friends/" + AB and w == "c": continue
            bd = P[w][p][2]
            assert allowed or (A_AGENT not in bd and "agent_a1" not in bd), (w, p)
    cf = P["c"]["/home/friends/" + AB]; assert cf[0] in (403, 404), cf[0]
    assert A_NAME not in cf[2] and B_AGENT not in cf[2] and A_AGENT not in cf[2] and "agent_a1" not in cf[2] and 'role="alert"' in cf[2]


@check("forty_char_names_are_fully_present_in_markup", ["U2 40 chars", "U2 emoji"])
def _():
    assert len(A_NAME) == 40 and len(B_AGENT) == 40
    assert re.search(r'<strong class="wrap">%s</strong>' % A_NAME, P["a"]["/home/friends"][2])
    assert A_NAME in H.unescape(P["b"]["/home/friends"][2]) and A_NAME in P["a"]["/home/account"][2]
    assert B_AGENT in P["a"]["/home/friends/" + AB][2] and A_NAME in P["b"]["/home/friends/" + AB][2]
    assert "<h1>%s</h1>" % A_NAME in P["b"]["/home/friends/" + AB][2]
    for p in ("/home/friends", "/home/friends/" + AB):
        b = main(P["a"][p][2]); assert "…" not in b and "..." not in b, p


@check("friend_code_is_labeled_readonly_unsubmitted_and_not_leaked", ["U2 own friend code", "U7 label"])
def _():
    for w in "abc":
        for p in ("/home/account", "/home/friends"):
            b = P[w][p][2]
            m = re.search(r'<label for="my-code">내 친구 코드</label><textarea id="my-code"([^>]*)>([^<]*)</textarea>', b)
            assert m, (w, p); at, v = m.groups()
            assert " readonly" in at and 'rows="2"' in at and 'aria-describedby="my-code-help"' in at and 'name=' not in at, at
            assert re.fullmatch(r"[A-Za-z0-9_-]{43}", v), v
            assert 'id="my-code-help"' in b and "전체 선택" in b and "로그인이나 연결 권한이 아닙니다" in b
            before = b[:b.index('id="my-code"')]; assert before.count("<form") == before.count("</form>"), (w, p)  # not inside a form
    assert len(set(CODES.values())) == 3
    for w in "abc":
        for o in "abc":
            if o == w: continue
            for p in PAGES + ["/home/friends/" + AB]:
                assert CODES[o] not in P[w][p][2], (w, o, p)
    for p in PAGES:  # codes appear only on the two code pages
        if p not in ("/home/account", "/home/friends"): assert CODES["a"] not in P["a"][p][2], p
    h = P["a"]["/home/account"][1]; assert "no-store" in h.get("cache-control", ""), h


@check("account_page_has_code_rotation_with_consequence_text", ["U2 code rotation UI"])
def _():
    b = main(P["a"]["/home/account"][2])
    assert '<details><summary>친구 코드 교체</summary><p>새 코드를 받으면 이전 코드로 새 요청을 받을 수 없습니다. 기존 친구와 이미 받은 요청은 유지됩니다.</p><form method="post" action="/home/friend-code"><button>새 친구 코드 받기</button></form></details>' in b


@check("friends_page_request_form_has_name_code_scope_consent_and_no_selector", ["U4 form", "U4 zero agent no selector"])
def _():
    for w, nm in (("a", A_NAME), ("b", "합성 회원 b"), ("c", "합성 회원 c")):
        b = main(P[w]["/home/friends"][2])
        f = re.search(r'<form method="post" action="/home/friends/request">.*?</form>', b, re.S).group(0)
        assert f"보내는 사람: <strong class=\"wrap\">{nm}</strong>" in f and 'href="/home/account">내 이름 바꾸기' in f
        assert re.search(r'<label for="friend-code">친구가 준 친구 코드</label><input id="friend-code" name="code" required autocomplete="off" spellcheck="false" value=""', f)
        assert 'id="friend-scope"' in f and "서로의 현재·앞으로 추가할 활성 에이전트가 연결됩니다" in f and "내 에이전트가 아직 없어도 요청할 수 있습니다" in f
        assert re.search(r'<input id="friend-consent" type="checkbox" name="consent" value="yes" required aria-describedby="friend-scope"><span>내 이름 ‘%s’을 상대에게 공개하고' % nm, f)
        assert "<select" not in f and "<option" not in f and f.rstrip().endswith('<button class="primary">친구 요청 보내기</button></form>')
        assert "미리 보이지 않습니다" in f
    assert "<select" not in P["c"]["/home/friends"][2]


@check("pending_received_and_sent_states", ["U4 received", "U4 sent", "U4 no peer preview"])
def _():
    ab = main(P["a"]["/home/friends"][2]); cb = main(P["c"]["/home/friends"][2])
    assert "합성 회원 c의 친구 요청" in ab and "서로의 현재·앞으로 추가할 에이전트가 연결됩니다" in ab and "수락 전에는 서로 메시지를 보낼 수 없습니다" in ab
    exp = re.search(r"요청 기한: (\d{4}-\d\d-\d\d \d\d:\d\d:\d\d KST)", ab).group(1); EVID["expires"] = exp
    assert f"요청 기한: {exp}" in cb
    assert f'<button class="primary" name="decision" value="accept" aria-label="친구 요청 수락: 합성 회원 c / {CA}">수락</button><button name="decision" value="deny" aria-label="친구 요청 거절: 합성 회원 c / {CA}">거절</button>' in ab
    assert "보낸 친구 요청" in cb and "요청 대기" in cb and f'name="decision" value="cancel" aria-label="보낸 친구 요청 취소: {CA}">요청 취소' in cb
    for bad in (A_NAME, A_AGENT, "agent_a", 'value="accept"', 'value="deny"'): assert bad not in cb, bad  # sender sees nothing about the receiver and cannot decide
    assert "아직 친구가 없습니다" in cb and "합성 회원 b" not in cb


@check("zero_agent_account_sees_add_agent_guidance_and_no_empty_selector", ["U4 zero agents"])
def _():
    for p in ("/home/agents", "/home/new"):
        b = main(P["c"][p][2]); assert "첫 에이전트를 추가하세요" in b and "<select" not in b, p
    assert '<a class="button primary" href="/home/agents">에이전트 추가</a>' in P["c"]["/home/new"][2]


@check("get_requests_never_create_or_accept_relationships", ["U4 GET no side effect"])
def _():
    snap = lambda w: re.sub(r"<textarea.*?</textarea>", "", tt(S[w].get("/home/friends")[2]))
    before = {w: snap(w) for w in "abc"}
    A = S["a"]; seen = []
    for q in (f"/home/friends?friend={CA}&decision=accept&generation=1", f"/home/friends?decision=accept&friend={CA}&consent=yes", f"/home/friends/{CA}?decision=accept", f"/home/friends/request?code={CODES['c']}&consent=yes", f"/home/friends/decision?friend={CA}&generation=1&decision=accept&consent=yes", "/home/friend-code", f"/home/new?peer={CA}", "/home/friends/request", "/home/agent-name?agent=agent_a1&name=x"):
        s, h, b = A.get(q); seen.append((q.split("?")[0], s)); assert s in (200, 400, 403, 404, 405, 422), (q, s)
    C = S["c"]
    for q in (f"/home/friends?code={CODES['a']}&consent=yes", f"/home/friends/request?code={CODES['a']}&consent=yes", "/home/friend-code"):
        C.get(q)
    after = {w: snap(w) for w in "abc"}
    assert before == after, {w: (before[w][:200], after[w][:200]) for w in "abc" if before[w] != after[w]}
    assert "합성 회원 c의 친구 요청" in before["a"] and "친구 목록 합성 회원 b" in before["a"]
    EVID["get_probe_statuses"] = seen


@check("repeated_buttons_have_target_and_discriminator_in_accessible_name", ["U7 repeated action names"])
def _():
    def names(b, visible):
        out = []
        for m in re.finditer(r"<(button|summary|a)\b([^>]*)>([^<]*)</\1>", b):
            tag, at, vis = m.groups(); lab = re.search(r'aria-label="([^"]*)"', at)
            if vis.strip() in visible: out.append((vis.strip(), H.unescape(lab.group(1)) if lab else vis.strip()))
        return out
    for p, vis in (("/home/agents", {"이름 저장", "이 에이전트 철회", "이름·연결 관리", "요청 상태 조회"}), ("/home/legacy", {"연결 수단 발급", "이 에이전트 철회 (모든 키·관계 종료)", "기존 연결 철회"})):
        ns = names(P["a"][p][2], vis); assert ns, p
        for v, n in ns:
            assert n != v and ("agent_a1" in n or "agent_a2" in n), (p, v, n)
        by = {}
        for v, n in ns: by.setdefault(v, []).append(n)
        for v, l in by.items():
            if len(l) > 1: assert len(set(l)) == len(l), (p, v, l)
    fa = P["a"]["/home/friends/" + AB][2]
    labs = re.findall(r'aria-label="(친구 (?:해제|차단)[^"]*)"', fa); assert len(labs) == 2 and all("합성 회원 b" in x for x in labs), labs
    labs = re.findall(r'aria-label="(친구 (?:해제|차단)[^"]*)"', P["b"]["/home/friends/" + AB][2]); assert len(labs) == 2 and all(A_NAME in x for x in labs), labs


@check("return_links_are_exact_internal_and_ignore_return_params", ["U7 exact return", "U7 no open redirect"])
def _():
    rb = main(P["a"]["/home/receipts"][2])
    assert '<a href="/home/agents">내 에이전트로 돌아가기</a>' in rb and '<a href="/home/legacy#receipts">관리 화면으로 돌아가기</a>' in rb
    assert '<nav aria-label="현재 위치"><a href="/home/agents">내 에이전트</a> / <a href="/home/receipts#receipt-form" aria-current="page">요청 상태 조회</a></nav>' in rb
    assert '<a href="/home/friends">친구 목록</a>' in P["a"]["/home/friends/" + AB][2] and '<a href="/home">대화 목록</a>' in P["a"]["/home/new"][2]
    for w in "abc":
        for p in P[w]:
            for href in re.findall(r'href="([^"]*)"', strip(P[w][p][2])):
                assert (href.startswith("/") and not href.startswith("//")) or href.startswith("#"), (w, p, href)
    A = S["a"]
    for k in ("return", "return_to", "next", "redirect", "continue", "url"):
        for ev in ("https://evil.example/x", "//evil.example", "javascript:alert(1)"):
            for p in ("/home/friends", "/home/receipts", "/home/account"):
                s, h, b = A.get(f"{p}?{k}={ev}"); assert "location" not in h and "evil.example" not in b and "javascript:" not in b, (p, k, ev)


@check("receipt_radio_selector_keeps_selection_and_input_u6", ["1단위 radio", "U6 preserve"])
def _():
    A = S["a"]
    s, h, b = A.get("/home/receipts")
    assert '<fieldset id="receipt-agent" role="radiogroup" tabindex="-1" autofocus><legend>내 에이전트</legend>' in b and ' checked' not in b.split("</fieldset>")[0]
    assert len(re.findall(r'<input type="radio" name="agent"[^>]* required', b)) == 2 and "<select" not in b
    s, h, b = A.get("/home/receipts?agent=agent_a2"); assert s == 200 and 'value="agent_a2" required checked' in b and 'value="agent_a1" required checked' not in b
    s, h, b = A.get('/home/receipts?agent=agent_a2&id=bad%22%3E%3Cb%3E'); assert s == 422 and "bad&#34;&gt;&lt;b&gt;" in b and "<b>" not in b and 'value="agent_a2" required checked' in b and 'aria-invalid="true"' in b and 'href="#receipt-id"' in b and 'role="alert"' in b
    NONE = "0199a3f2-4c10-7a11-8b22-334455667799"
    r1 = A.get(f"/home/receipts?agent=agent_a1&id={NONE}"); r2 = A.get(f"/home/receipts?agent=agent_a2&id={NONE[:-1]}8")
    for r in (r1, r2):
        assert r[0] == 403 and REFUSE in r[2] and 'role="alert"' in r[2]
    n = lambda r: re.sub(r'0199a3f2-[0-9a-f-]+', "<ID>", re.sub(r' checked', "", re.sub(r"\s+", " ", main(r[2]))))
    assert n(r1) == n(r2)
    f = A.get(f"/home/receipts?agent=agent_b1&id={NONE}"); assert f[0] == 403 and "agent_b1" not in f[2] and B_AGENT not in f[2] and ' checked' not in f[2]
    for bad in ("a" * 200, "ab\x00cd", "ab\u202ecd"):
        s, h, b = A.get("/home/receipts?" + urllib.parse.urlencode({"agent": "agent_a1", "id": bad})); assert s in (400, 422) and bad not in b and 'value="agent_a1" required checked' in b, repr(bad[:5])


@check("error_pages_use_alert_and_do_not_leak", ["U7 alert", "U2 no existence leak"])
def _():
    C = S["c"]
    for u in ("/home/friends/zzzzzzzzzzzzzzzzzzzzzz", "/home/friends/" + AB, "/home/friends/%00", "/home/friends/" + "a" * 300):
        s, h, b = C.get(u); assert s in (400, 403, 404), (u, s)
        assert 'role="alert"' in b and "처리할 수 없습니다" in b and A_NAME not in b and "합성 회원" not in b, u
    s1 = C.get("/home/friends/zzzzzzzzzzzzzzzzzzzzzz"); s2 = C.get("/home/friends/" + AB)
    n = lambda r: (r[0], re.sub(r"\s+", " ", main(r[2])))
    assert n(s1) == n(s2), "nonexistent and other-people friend ids must look the same"


@check("blocked_endpoints_stay_blocked_by_fixture", ["fixture allow-list"])
def _():
    A = S["a"]
    for p in ("/home/gate", "/gate", "/home/conversation?from=agent_a1&to=agent_a2", "/auth/logout", "/fixture/recover"):
        assert A.get(p)[0] == 403, p
    assert A.req("POST", "/home/friends/request", {"code": "x", "consent": "yes"}, origin="http://evil.example", patient=True)[0] in (400, 403)  # cross-site POST refused


@check("cache_and_security_headers_on_code_pages", ["security"])
def _():
    for w in "abc":
        for p in ("/home/account", "/home/friends"):
            h = P[w][p][1]; assert "no-store" in h["cache-control"] and h["x-content-type-options"] == "nosniff" and h["referrer-policy"] == "no-referrer", h
            csp = h["content-security-policy"]; assert "script-src" not in csp and "default-src 'none'" in csp and "unsafe-eval" not in csp
        assert not re.search(r"<script|onclick=|onload=|javascript:", P[w]["/home/account"][2]), "reading/selecting the code needs no JS"


@check("anonymous_429_screen_has_no_misleading_empty_agent_text", ["1단위 anon 429"])
def _():
    an = Browser(); codes = []; last = None
    for i in range(40):
        s, h, b = an.req("GET", "/home/receipts", patient=False); codes.append(s)
        if s == 429: last = (s, h, b); break
    EVID["anon429_codes"] = codes
    assert last, codes
    m = tt(last[2]); EVID["anon429_text"] = m
    assert "조회할 내 에이전트가 없습니다" not in m and "내 에이전트가 없" not in m, m
    assert 'role="alert"' in last[2] and ("retry-after" in last[1] or re.search(r"\d+\s*(초|분)|\d{4}-\d\d-\d\d", m)), (last[1], m)


if __name__ == "__main__":
    EVID["member_window_429_retries"] = Browser.waited429
    out = {"origin": ORIGIN, "fixed_sha": "b2e683ae677c78cd7bdc47c75820bcb191b94838", "evidence": EVID, "results": RES}
    json.dump(out, open(os.path.join(os.path.dirname(os.path.abspath(__file__)), "http-result.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
    for r in RES: print(r["result"].upper(), r["check"], r.get("detail", ""))
    sys.exit(0 if all(r["result"] == "pass" for r in RES) else 1)
