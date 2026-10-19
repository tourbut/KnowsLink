"""Recheck of the 5fa3333 fix (4 items) at HTTP level. Read-only GETs plus POSTs that the product must refuse (nothing is saved).
Item 1: new-conversation radios (labels, defaults, explicit values, forged values, other accounts). Item 2: the two links on my-agents.
Item 3: paths other than /home/agent-name are outside the fixture allow-list, so they are only checked to be blocked (not_run for the product text).
Item 4: refusal notice markup. The browser stage (browser.mjs recheck) adds keyboard, accessibility tree and layout."""
import html as H, json, os, re, sys, traceback
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__))); from client import *

M = json.load(open(os.environ.get("NAMES_MANIFEST", "C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json")))
AB, CA = M["friend_ab"], M["request_ca"]
A_AGENT, B_AGENT = "긴한글" * 13, "😀" * 40
KIND = "Node 22 로컬 클라이언트"
RES = []; EVID = {}


def check(name, covers):
    def deco(fn):
        try: fn(); RES.append({"check": name, "covers": covers, "result": "pass"})
        except Exception as e: RES.append({"check": name, "covers": covers, "result": "fail", "detail": f"{type(e).__name__}: {e}", "trace": traceback.format_exc(limit=3)})
        return fn
    return deco


A, B, C = as_("a"), as_("b"), as_("c")


def groups(b):
    """{'from': [(value, label, checked)], 'to': [...], 'legends': {...}} from the new-conversation form."""
    out = {"legends": {}}
    for g in ("from", "to"):
        m = re.search(rf'<fieldset id="{g}" role="radiogroup"><legend>([^<]*)</legend>(.*?)</fieldset>', b, re.S)
        if not m: out[g] = None; continue
        out["legends"][g] = H.unescape(m.group(1))
        out[g] = [(v, H.unescape(lab), "checked" in attrs) for v, attrs, lab in re.findall(r'<input type="radio" name="%s" value="([^"]+)"([^>]*)><span>([^<]*)</span>' % g, m.group(2))]
    return out


parts = lambda s: [x.strip() for x in s.split(" / ")]
checked = lambda g, side: [v for v, _, c in g[side] if c]


@check("new_conversation_self_is_two_radio_groups_with_name_kind_discriminator_labels_and_old_defaults", ["fix 1 radio", "U2 label order", "defaults unchanged"])
def _():
    s, h, b = A.get("/home/new?peer=self"); assert s == 200
    assert "<select" not in b and "<option" not in b
    g = groups(b); EVID["self_groups"] = g
    assert g["legends"] == {"from": "내 발신 에이전트", "to": "내 수신 에이전트"}, g["legends"]
    for side in ("from", "to"):
        assert [v for v, _, _ in g[side]] == ["agent_a1", "agent_a2"], g[side]
        labs = [l for _, l, _ in g[side]]; assert len(set(labs)) == 2, labs
        for v, l, _ in g[side]: assert parts(l) == [A_AGENT, KIND, v], l
        assert len(re.findall(rf'<input type="radio" name="{side}" value="[^"]+" required', b)) == 2, "every radio required"
    assert checked(g, "from") == ["agent_a1"] and checked(g, "to") == ["agent_a2"], (checked(g, "from"), checked(g, "to"))
    assert re.search(r'<form method="get" action="/home/conversation">', b) and 'class="primary">대화 열기</button>' in b


@check("new_conversation_friend_has_peer_legend_single_receiver_radio_and_first_sender_default", ["fix 1 radio", "peer legend", "emoji label"])
def _():
    s, h, b = A.get("/home/new?peer=" + AB); assert s == 200 and "<select" not in b
    g = groups(b); EVID["peer_groups"] = g
    assert g["legends"]["to"] == "합성 회원 b의 수신 에이전트" and g["legends"]["from"] == "내 발신 에이전트"
    assert [(v, parts(l)) for v, l, _ in g["to"]] == [("agent_b1", [B_AGENT, KIND, "agent_b1"])], g["to"]
    assert checked(g, "from") == ["agent_a1"] and checked(g, "to") == ["agent_b1"]


@check("new_conversation_other_accounts_see_only_their_own_and_friends_agents", ["authority unchanged", "U2 no leak"])
def _():
    gb = groups(B.get("/home/new?peer=" + AB)[2]); EVID["b_groups"] = gb
    assert [v for v, _, _ in gb["from"]] == ["agent_b1"] and [v for v, _, _ in gb["to"]] == ["agent_a1", "agent_a2"], gb
    assert len({l for _, l, _ in gb["to"]}) == 2 and checked(gb, "to") == ["agent_a1"]
    s, h, b = C.get("/home/new?peer=" + AB)  # C is not B's friend and has no agent
    EVID["c_new_with_foreign_peer"] = [s, text(main(b))[:200]]
    assert "agent_a" not in b and "agent_b1" not in b and A_AGENT not in b and B_AGENT not in b, "foreign peer must not leak agent names"
    s, h, b = C.get("/home/new?peer=self"); assert "<input type=\"radio\"" not in b and "<select" not in b and "첫 에이전트를 추가하세요" in b


@check("new_conversation_get_ignores_from_to_query_and_forged_values_never_become_choices", ["defaults not steerable by URL", "forged agent value ignored"])
def _():
    # /home/new takes From/To only when re-rendered after a failed open (POST/GET /home/conversation is outside the fixture), so a URL cannot preselect or inject a choice
    forged = {}
    for q in ("peer=self&from=agent_a2&to=agent_a1", f"peer={AB}&from=agent_a2", "peer=self&from=agent_b1", "peer=self&to=agent_b1", "peer=self&from=agent_b1&to=agent_b1", f"peer={AB}&to=agent_a1", f"peer={AB}&from=agent_b1", "peer=self&from=%3Cscript%3E&to=x%22%3E", "peer=self&from=agent_a1&to=agent_a1"):
        s, h, b = A.get("/home/new?" + q); g = groups(b)
        assert s == 200 and g["from"] and g["to"], (q, s)
        assert [v for v, _, _ in g["from"]] == ["agent_a1", "agent_a2"], (q, g["from"])
        assert [v for v, _, _ in g["to"]] == (["agent_b1"] if AB in q else ["agent_a1", "agent_a2"]), (q, g["to"])
        assert len(checked(g, "from")) == 1 and len(checked(g, "to")) == 1, (q, g)
        assert "<script" not in main(b) and "x\"&gt;" not in b and "x\">" not in b
        forged[q] = [checked(g, "from"), checked(g, "to")]
    EVID["forged_new_page_selection"] = forged
    assert all(v == [["agent_a1"], ["agent_b1"]] if AB in k else v == [["agent_a1"], ["agent_a2"]] for k, v in forged.items()), forged


@check("conversation_open_endpoint_is_outside_the_fixture_so_open_result_is_browser_request_only", ["fixture limit"])
def _():
    for q in ("from=agent_a2&to=agent_a1&peer=self", "from=agent_b1&to=agent_a1", "from=agent_a1&to=agent_b1"):
        s, h, b = A.get("/home/conversation?" + q); EVID.setdefault("conversation_get", {})[q] = [s, text(b)[:60]]
        assert s == 403 and "outside names fixture scope" in b


@check("my_agents_two_links_are_separate_paragraphs_with_exact_targets", ["fix 2 link split", "targets"])
def _():
    b = A.get("/home/agents")[2]
    for ag in ("agent_a1", "agent_a2"):
        m = re.search(rf'<p><a href="/home/receipts\?agent={ag}#receipt-form" aria-label="요청 상태 조회: [^"]+">요청 상태 조회</a></p><p><a href="/home/legacy#agents">연결·키 관리</a></p>', b)
        assert m, ag
    assert "요청 상태 조회</a><a" not in b and "요청 상태 조회</a> <a" not in b
    r = A.get("/home/receipts?agent=agent_a2"); assert r[0] == 200 and re.search(r'value="agent_a2" required checked', r[2]), "lookup preselects the agent"
    l = A.get("/home/legacy"); assert l[0] == 200 and 'id="agents"' in l[2], "manage target anchor exists"
    EVID["manage_title"] = H.unescape(re.search(r"<title>(.*?)</title>", l[2]).group(1))


@check("received_request_has_refusal_notice_between_accept_and_deny_described_by_id", ["fix 4 notice", "aria-describedby"])
def _():
    b = main(A.get("/home/friends")[2])
    m = re.search(r'<form method="post" action="/home/friends/decision">(.*?)</form>', b, re.S); assert m
    f = m.group(1)
    kids = re.findall(r'<(button|p)[^>]*?(?: id="([^"]*)")?[^>]*>', f)
    EVID["decision_form_children"] = kids
    pid = f"deny-help-{CA}"
    assert f.index('value="accept"') < f.index(f'id="{pid}"') < f.index('value="deny"'), "notice sits between accept and deny"
    assert H.unescape(text(re.search(rf'<p class="muted" id="{pid}">(.*?)</p>', f).group(1))) == "거절하면 요청한 사람의 기록에 거절로 표시됩니다. 알리지 않으려면 응답하지 않아도 기한이 지나면 요청이 끝납니다."
    deny = re.search(r'<button name="decision" value="deny"[^>]*>', f).group(0)
    assert f'aria-describedby="{pid}"' in deny and "aria-describedby" not in re.search(r'<button class="primary" name="decision" value="accept"[^>]*>', f).group(0)
    assert len(re.findall(r'id="deny-help-', b)) == len(re.findall(r'value="deny"', b)), "one notice per deny button, unique ids"
    assert len(set(re.findall(r'id="(deny-help-[^"]+)"', b))) == len(re.findall(r'id="deny-help-', b))
    assert "거절하면 요청한 사람" not in main(C.get("/home/friends")[2]), "requester's own view does not show the receiver's notice"


@check("other_name_paths_are_outside_the_fixture_so_only_blocked_not_verified", ["fixture limit: profile / add-agent"])
def _():
    for path, form in (("/home/profile", {"name": "ab\x07cd"}), ("/home/agents", {"name": "ab\x07cd"})):
        s, h, b = A.post(path, form); EVID.setdefault("blocked_name_posts", {})[path] = s
        assert s == 403 and "outside names fixture scope" in b, (path, s)


if __name__ == "__main__":
    EVID["member_window_429_retries"] = Browser.waited429
    out = {"origin": ORIGIN, "fixed_sha": "5fa3333b5948fd36792662b202915bc6243cac14", "evidence": EVID, "results": RES}
    json.dump(out, open(os.path.join(os.path.dirname(os.path.abspath(__file__)), "recheck-result-5fa3333.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
    for r in RES: print(r["result"].upper(), r["check"], r.get("detail", "")[:500])
    sys.exit(0 if all(r["result"] == "pass" for r in RES) else 1)
