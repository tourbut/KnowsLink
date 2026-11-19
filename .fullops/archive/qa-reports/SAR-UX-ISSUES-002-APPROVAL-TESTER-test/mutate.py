"""State-changing checks for the approval fixture. Only pending_decide_1 and pending_decide_2 receive POSTs, and only as the last step.
Run: python mutate.py 1   (token/foreign rejections, approve decide_1, duplicate clicks, stash decide_2 stale form)
     [browser.mjs flow denies decide_2 with real keys]
     python mutate.py 2   (stale approve/deny replay on decide_2, final state)
Result: mutate-result-<n>.json, exit 1 on any FAIL.
"""
import json, re, sys, html
sys.path.insert(0, ".")
from client import Browser, as_, main, text

M = json.load(open("C:/Users/shin/AppData/Local/KnowsLinkDevTools/approval-ux002-dev/fixture-manifest.json"))
G = M["gates"]
D1, D2 = G["pending_decide_1"], G["pending_decide_2"]
STAGE = sys.argv[1]
STASH = "stash.json"
results = []


def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": detail if not ok else ""})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:300])


def csrf(body):
    return re.search(r'name="csrf" value="([^"]+)"', body).group(1)


def state(b, gid):
    s, h, body = b.get("/home/gates/" + gid)
    m = main(body)
    st = re.search(r'<h2 id="state-h">.*?<dd>(.*?)</dd>', m, re.S)
    return {"http": s, "badge": text(st.group(1)) if st else "", "has_form": "<form" in m, "has_approve": "이 요청 승인" in m}


def pend_count(b):
    return re.search(r"승인 (\d+)건", text(b.get("/home")[2])).group(1)


def decide(b, gid, decision, token, extra=None):
    path = f"/home/gates/{gid}" + ("/deny" if decision == "deny" else "")
    form = {"csrf": token, "decision": decision}
    form.update(extra or {})
    return b.post(path, form)


REFUSED = "이번 제출은 적용되지 않았습니다"
if STAGE == "1":
    a = as_("a")
    s1 = state(a, D1); s2 = state(a, D2)
    chk("pre_both_decide_gates_pending_with_buttons", s1["badge"].endswith("대기") and s2["badge"].endswith("대기") and s1["has_approve"] and s2["has_approve"], (s1, s2))
    c0 = pend_count(a)
    t1 = csrf(a.get("/home/gates/" + D1)[2]); t2 = csrf(a.get("/home/gates/" + D2)[2])
    chk("csrf_tokens_differ_per_gate_page", t1 != t2, "")
    json.dump({"t1": t1, "t2": t2, "pending_before": c0}, open(STASH, "w"))

    # cross-gate token: decide_1's token on decide_2 must not decide anything
    for dec in ("approve", "deny"):
        s, h, body = decide(a, D2, dec, t1)
        chk(f"cross_gate_token_{dec}_rejected_403", s == 403, (s, text(body)[:160]))
    chk("cross_gate_token_left_both_gates_pending", state(a, D1)["badge"].endswith("대기") and state(a, D2)["badge"].endswith("대기"), "")
    # wrong / missing / empty token
    for label, tok in (("wrong", "AAAA" + t1[4:]), ("empty", ""), ("truncated", t1[:10])):
        s, h, body = decide(a, D1, "approve", tok)
        chk(f"csrf_{label}_rejected_403_gate_unchanged", s == 403 and state(a, D1)["badge"].endswith("대기"), (s, text(body)[:120]))
    s, h, body = a.post(f"/home/gates/{D1}", {"decision": "approve"})
    chk("csrf_absent_rejected_403_gate_unchanged", s == 403 and state(a, D1)["badge"].endswith("대기"), s)
    # unknown decision value must not approve
    s, h, body = decide(a, D1, "maybe", t1)
    chk("unknown_decision_value_does_not_decide", s >= 400 and state(a, D1)["badge"].endswith("대기"), (s, text(body)[:120]))
    # foreign members with A's token
    for who in "bc":
        x = as_(who)
        for dec in ("approve", "deny"):
            s, h, body = decide(x, D1, dec, t1)
            chk(f"foreign_member_{who}_{dec}_with_stolen_token_403_no_leak", s == 403 and "긴한글" not in body and "agent_a" not in body, (s, text(body)[:120]))
    # GET with query tampering cannot decide
    s, h, body = a.get(f"/home/gates/{D1}?decision=approve&csrf={t1}")
    chk("get_with_decision_query_does_not_decide", state(a, D1)["badge"].endswith("대기"), s)
    # approve decide_1 (A), real first decision
    s, h, body = decide(a, D1, "approve", t1)
    chk("approve_decide_1_redirects_303_to_detail", s == 303, (s, h.get("location")))
    st = state(a, D1)
    chk("approved_detail_shows_approved_without_controls", "승인" in st["badge"] and not st["has_form"] and not st["has_approve"], st)
    sl, hl, bl = a.get("/home/legacy")
    row = [r for r in re.findall(r'<li class="row">(.*?)</li>', main(bl), re.S) if D1 in r][0]
    chk("list_row_decide_1_shows_approved", "✓ 승인" in text(row), text(row))
    chk("pending_count_dropped_by_one", int(pend_count(a)) == int(c0) - 1, (c0, pend_count(a)))
    # duplicate clicks / stale form replays on the decided gate
    for dec in ("approve", "deny"):
        s, h, body = decide(a, D1, dec, t1)
        mm = main(body)
        chk(f"duplicate_{dec}_after_approval_409_with_safe_message_and_return", s == 409 and REFUSED in mm and "/home/legacy#gates" in mm and "<form" not in mm, (s, text(mm)[:200]))
        chk(f"duplicate_{dec}_page_has_title_h1_alert", "<title>요청 승인 · KnowsLink</title>" in body and body.count("<h1>") == 1 and 'role="alert"' in body, "")
    open("409-page.html", "w", encoding="utf-8").write(body)
    chk("first_decision_kept_after_duplicates", "승인" in state(a, D1)["badge"] and "거절" not in state(a, D1)["badge"], state(a, D1))
    # fresh CSRF from a new page load cannot reopen it either
    s, h, body = a.get("/home/gates/" + D1)
    chk("approved_page_offers_no_new_csrf_token", 'name="csrf"' not in body, "")
    # foreign member still locked out after the decision
    xb = as_("b"); s, h, body = xb.get("/home/gates/" + D1)
    chk("foreign_member_still_403_after_decision", s == 403 and "긴한글" not in body, s)
    chk("decide_2_untouched_still_pending", state(a, D2)["badge"].endswith("대기") and state(a, D2)["has_approve"], state(a, D2))

if STAGE == "2":
    a = as_("a")
    t = json.load(open(STASH))
    st = state(a, D2)
    chk("decide_2_denied_by_real_browser_flow", "거절" in st["badge"] and not st["has_form"], st)
    # the denial came from the browser's own session, so the stage-1 token belongs to another session: refused (403) with the same safe message, nothing decided
    for dec in ("approve", "deny"):
        s, h, body = decide(a, D2, dec, t["t2"])
        mm = main(body)
        chk(f"other_session_token_{dec}_on_denied_gate_refused_403_safe_message_and_return", s == 403 and REFUSED in mm and "/home/legacy#gates" in mm and "<form" not in mm, (s, text(mm)[:200]))
    chk("denial_kept_after_stale_resubmits", "거절" in state(a, D2)["badge"] and "승인" not in state(a, D2)["badge"].replace("거절", ""), state(a, D2))
    chk("decide_1_still_approved", "승인" in state(a, D1)["badge"], state(a, D1))
    chk("pending_count_final_is_2", pend_count(a) == "2", pend_count(a))
    sl, hl, bl = a.get("/home/legacy")
    row = [r for r in re.findall(r'<li class="row">(.*?)</li>', main(bl), re.S) if D2 in r][0]
    chk("list_row_decide_2_shows_denied", "거절" in text(row), text(row))
    others = {k: state(a, G[k])["badge"] for k in ("pending_query_hint", "pending_commit_long")}
    chk("other_pending_gates_unaffected", all(v.endswith("대기") for v in others.values()), others)

json.dump({"stage": STAGE, "checks": results}, open(f"mutate-result-{STAGE}.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
f = [r for r in results if r["result"] == "FAIL"]
print(f"{len(results)-len(f)} PASS / {len(f)} FAIL")
sys.exit(1 if f else 0)
