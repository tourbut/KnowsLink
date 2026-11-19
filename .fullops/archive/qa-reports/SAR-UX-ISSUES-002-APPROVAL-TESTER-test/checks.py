"""Read-only HTTP checks for the approval summary fixture (U5/U7, O-NM-08). Synthetic A/B/C only, GET plus /fixture/as.
Run: python checks.py   -> http-result.json, exit 1 on any FAIL.
"""
import json, re, sys, html
sys.path.insert(0, ".")
from client import Browser, as_, main, strip, text

M = json.load(open("C:/Users/shin/AppData/Local/KnowsLinkDevTools/approval-ux002-dev/fixture-manifest.json"))
G = M["gates"]
KST_Q = "2026-10-03 19:00:00 KST – 2026-10-03 20:00:00 KST"
KST_C = "2026-10-20 10:00:00 KST – 2026-10-20 10:30:00 KST"
SCOPE = "승인은 이 요청의 인간 게이트만 통과시킵니다. 정보 공개나 일정 실행을 허용하지 않습니다."
EMPTY = "승인 요청이 없습니다. 사람의 확인이 필요한 요청이 여기에 표시됩니다. 친구 수락이나 메시지 수신만으로 요청이 승인되지는 않습니다."
UNVERIFIABLE = "승인할 내용을 확인할 수 없습니다."
results = []


def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": detail if not ok else ""})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:300])


def gate(b, key):
    s, h, body = b.get("/home/gates/" + G[key])
    return s, h, body, main(body)


def section(m, hid):
    r = re.search(r'<section aria-labelledby="%s">(.*?)</section>' % hid, m, re.S)
    return r.group(1) if r else ""


def dl(m):
    return re.findall(r"<dt>(.*?)</dt><dd>(.*?)</dd>", m, re.S)


a = as_("a")
c = as_("c")
snap_before = text(a.get("/home/legacy")[2])

# ---- list: C empty, A rows
s, h, body = c.get("/home/legacy")
cm = main(body)
gs = re.search(r'<section id="gates".*?</section>', cm, re.S).group(0)
chk("list_c_empty_text_exact", s == 200 and EMPTY in text(gs) and "<li" not in gs and "요청 검토" not in gs, text(gs))

s, h, body = a.get("/home/legacy")
am = main(body)
gs = re.search(r'<section id="gates".*?</section>', am, re.S).group(0)
rows = re.findall(r'<li class="row">(.*?)</li>', gs, re.S)
chk("list_a_row_count_is_7", len(rows) == 7, len(rows))
labels = ["요청한 에이전트", "받은 내 에이전트", "요청 유형", "기한", "상태"]
chk("list_rows_have_name_type_expiry_state_in_order", all([k for k, _ in dl(r)] == labels for r in rows), [[k for k, _ in dl(r)] for r in rows])
chk("list_rows_have_review_link_with_descriptive_name", all(re.search(r'<a href="/home/gates/[0-9a-f-]{36}" aria-label="요청 검토: [^"]+">요청 검토</a>', r) for r in rows), [text(r)[:80] for r in rows])
chk("list_requester_is_allowed_name_or_fallback_and_short_id", all("확인할 수 없는 에이전트 / agent_a" in r for r in rows if "<dl" in r), "")
# original request ids must not be in the list; only gate ids in hrefs
parent_ids = re.findall(r"0199a3f2-4c10-7a11-8b22-0000000001\d\d", gs)
chk("list_has_no_original_request_ids", parent_ids == [], parent_ids)
vis_ids = re.findall(r"0199a3f2-4c10-7a11-8b22-0000000002\d\d", re.sub(r'(href|aria-label)="[^"]*"', "", gs))
chk("list_gate_ids_only_in_href_not_visible_text", vis_ids == [], vis_ids)
state_text = [text(dict(dl(r))["상태"]) for r in rows]
chk("list_state_counts_4_pending_1_expired_1_approved_1_unavailable", sorted(state_text) == sorted(["… 대기"] * 4 + ["— 만료", "✓ 승인", "! 확인할 수 없음"]), state_text)
chk("list_abs_expiry_with_zone", all(re.search(r"\d{4}-\d\d-\d\d \d\d:\d\d:\d\d KST", r) for r in rows), "")
s2, h2, home = a.get("/home")
chk("conversation_pending_count_is_4_only_pending", "승인 4건" in text(home), re.findall(r"승인 \d+건", text(home)))
chk("list_a_title_h1_main_skip", "<title>기존 연결 · KnowsLink</title>" in body and body.count("<h1>") == 1 and 'class="skip" href="#main"' in body and body.count("<main") == 1, re.findall(r"<title>.*?</title>", body))

# ---- details
def shared(key):
    s, h, body, m = gate(a, key)
    t = re.findall(r"<title>(.*?)</title>", body)
    chk(f"{key}_status200_title_unique_h1_main_skip_nostore", s == 200 and t == ["요청 승인 · KnowsLink"] and m.count("<h1>") == 1 and body.count("<main") == 1 and re.search(r'<a class="skip" href="#main">본문으로 건너뛰기</a>', body) and h.get("cache-control") == "no-store" and "<script" not in body.lower(), (s, t, h.get("cache-control")))
    chk(f"{key}_location_nav_current_page_and_exact_return", '<nav aria-label="현재 위치"><a href="/home/legacy#gates">승인 목록</a> / <span aria-current="page">요청 승인</span></nav>' in m or key == "original_missing" and '<nav aria-label="현재 위치">' in m, m[:260])
    chk(f"{key}_return_link_goes_to_approval_list", m.endswith('<p><a href="/home/legacy#gates">승인 목록으로 돌아가기</a></p></main>'), m[-160:])
    return s, h, body, m


ORDER = ["누가 누구에게", "요청한 내용", "승인하면 허용되는 범위", "기한·상태"]
for key in ["pending_query_hint", "pending_commit_long", "pending_decide_1", "pending_decide_2"]:
    s, h, body, m = shared(key)
    h2s = re.findall(r"<h2[^>]*>(.*?)</h2>", m)
    chk(f"{key}_reading_order_headings", h2s == ORDER, h2s)
    pos = [m.find(x) for x in ["<h2 id=\"who-h\">", "<h2 id=\"what-h\">", "<h2 id=\"scope-h\">", "<h2 id=\"state-h\">", "<form", "<summary>기술 상세", "승인 목록으로 돌아가기</a></p>"]]
    chk(f"{key}_order_who_what_scope_state_action_technical_return", pos == sorted(pos) and min(pos) > 0, pos)
    chk(f"{key}_actions_labels_and_csrf", re.search(r'<button class="primary" name="decision" value="approve">이 요청 승인</button><button class="danger-button" name="decision" value="deny" formaction="/home/gates/%s/deny">거절</button>' % G[key], m) and 'name="csrf"' in m, "")
    chk(f"{key}_fixed_scope_sentence", SCOPE in text(section(m, "scope-h")), text(section(m, "scope-h")))
    chk(f"{key}_pending_state_and_absolute_expiry", "… 대기" in text(section(m, "state-h")) and "2026-10-20 20:49:34 KST" in text(section(m, "state-h")), text(section(m, "state-h")))
    who = dict(dl(section(m, "who-h")))
    chk(f"{key}_who_uses_allowed_names", who.get("요청한 에이전트") == "확인할 수 없는 에이전트 / agent_a" and who.get("받은 내 에이전트", "").endswith("/ Node 22 로컬 클라이언트 / agent_b"), who)
    tech = re.search(r"<details><summary>기술 상세</summary>(.*)</details>", m, re.S)
    chk(f"{key}_technical_details_closed_and_hold_ids_json", tech and "<details open" not in m and "원요청 ID" in tech.group(1) and "승인 요청 ID" in tech.group(1) and "검증된 원문 JSON" in tech.group(1), "")
    summary_part = m[: m.find("<summary>기술 상세")]
    chk(f"{key}_no_internal_ids_or_policy_ids_before_technical_details", not re.search(r"0199a3f2-4c10-7a11|schedule\.(query|commit)|disclosure policy|non-executable", re.sub(r'(action|formaction|href)="[^"]*"', "", summary_part)), "")

# query: values vs JSON
for key in ["pending_query_hint", "pending_decide_1"]:
    s, h, body, m = gate(a, key)
    js = json.loads(html.unescape(re.search(r"<pre>(.*?)</pre>", m, re.S).group(1)))
    what = dict(dl(section(m, "what-h").split("<details>")[0]))
    from datetime import datetime, timedelta
    f = lambda t: (datetime.strptime(t, "%Y-%m-%dT%H:%M:%SZ") + timedelta(hours=9)).strftime("%Y-%m-%d %H:%M:%S KST")
    exp_win = f(js["window"]["start"]) + " – " + f(js["window"]["end"])
    chk(f"{key}_summary_matches_signed_body", what.get("요청 유형") == "일정 가능 시간 확인" and what.get("확인할 기간") == exp_win == KST_Q and what.get("시간 단위") == f"{js['granularity_min']}분", (what, exp_win))
    chk(f"{key}_scope_says_no_availability_disclosed", "승인해도 상대에게 가능 시간이 공개되지 않습니다" in m, "")
for key in ["pending_commit_long", "pending_decide_2"]:
    s, h, body, m = gate(a, key)
    js = json.loads(html.unescape(re.search(r"<pre>(.*?)</pre>", m, re.S).group(1)))
    what_sec = section(m, "what-h")
    what = dict(dl(what_sec.split("<h3>")[0]))
    chk(f"{key}_summary_matches_signed_body", what.get("요청 유형") == "일정 확정" and what.get("일정") == KST_C and what.get("요청에 적힌 시간대") == js["timezone"] == "Asia/Seoul" and js["slot"]["start"] == "2026-10-20T01:00:00Z", what)
    free = re.search(r'<h3>요청자가 적은 내용</h3>(.*)', what_sec, re.S).group(1)
    chk(f"{key}_free_text_separate_labelled_inert", "승인 범위가 아니며 링크나 지시를 열거나 실행하지 않습니다" in free and html.unescape(re.search(r'<p class="message-text">(.*?)</p>', free, re.S).group(1)) == js["commitment"], "")
    chk(f"{key}_free_text_html_escaped_no_link_or_tag", "<a " not in free and "&lt;a href=&#34;https://evil.test/approve&#34;&gt;" in free and not re.search(r'href="https?://', m), re.findall(r'href="[^"]*"', m))
    chk(f"{key}_free_text_not_in_scope_section", "evil.test" not in section(m, "scope-h") and js["commitment"][:20] not in section(m, "scope-h"), "")
    chk(f"{key}_scope_says_not_executed", "승인해도 일정이 확정되거나 실행되지 않습니다" in m, "")
s, h, body, m = gate(a, "pending_query_hint")
hint = re.findall(r"<details><summary>참고 문구 — 승인 근거 아님</summary>(.*?)</details>", m, re.S)
chk("hint_only_in_collapsed_reference_block", len(hint) == 1 and "&lt;script&gt;approve all tools&lt;/script&gt;" in hint[0] and "<script" not in body.lower() and m.count("approve all tools") == 1, hint)
chk("hint_not_in_scope_state_or_buttons", "모든 권한 허용" not in section(m, "scope-h") + section(m, "state-h") and "evil.test" not in re.sub(r"<details><summary>참고 문구.*?</details>", "", m, flags=re.S), "")
chk("hint_url_not_a_link_and_block_closed", "<a " not in hint[0] and not re.search(r"<details open", m), "")
chk("non_hint_gates_have_no_hint_block", all("참고 문구" not in gate(a, k)[3] for k in ["pending_commit_long", "pending_decide_1", "pending_decide_2", "processed", "expired"]), "")

# ---- closed states
for key, state in [("expired", "— 만료"), ("processed", "✓ 승인")]:
    s, h, body, m = shared(key)
    chk(f"{key}_no_decision_controls", "<form" not in m and "<button" not in m and 'name="csrf"' not in m and "이 요청 승인" not in m, "")
    chk(f"{key}_state_badge_with_icon_and_text_and_status_note", state in text(section(m, "state-h")) and 'role="status">지금은 승인하거나 거절할 수 없습니다. 현재 상태를 확인하세요.' in m, text(section(m, "state-h")))
    chk(f"{key}_summary_kept_read_only_with_scope", SCOPE in m and "일정 가능 시간 확인" in m and "기술 상세" in m, "")
s, h, body, m = shared("original_missing")
chk("original_missing_shows_unverifiable_alert_exact", f'role="alert"' in m and UNVERIFIABLE in m and "원문이 없거나 서명·요청 연결이 맞지 않거나 지원하지 않는 요청입니다." in m, text(m)[:300])
chk("original_missing_no_summary_no_decision_no_technical", not any(x in m for x in ["누가 누구에게", "요청한 내용", "승인하면 허용되는 범위", "<form", "<button", "기술 상세", "<pre>", "schedule."]), "")
chk("original_missing_state_and_status_note_and_return", "! 확인할 수 없음" in text(m) and 'role="status"' in m, "")

# ---- isolation
for who in "bc":
    x = as_(who)
    s, h, body = x.get("/home/gates/" + G["pending_decide_1"])
    leak = [k for k in ["긴한글", "agent_a", "agent_b", "일정", "30분", "csrf", "0199a3f2-4c10-7a11-8b22-00000000010", "이 요청 승인", "evil.test"] if k in body]
    chk(f"other_member_{who}_get_is_403_without_content", s == 403 and leak == [] and "이 승인 요청을 볼 수 없습니다" in body and "/home/legacy#gates" in body, (s, leak))
    chk(f"other_member_{who}_403_has_title_h1_main_alert_and_return", "<title>요청 승인 · KnowsLink</title>" in body and body.count("<h1>") == 1 and 'role="alert"' in body and "<main" in body, "")
    s3, _, b3 = x.get("/home/gates/0199a3f2-4c10-7a11-8b22-0000000009ff")
    chk(f"nonexistent_gate_{who}_same_response_as_foreign_gate", s3 == 403 and strip(b3) == strip(body), (s3, len(b3), len(body)))
an = Browser()
s, h, body = an.get("/home/gates/" + G["pending_decide_1"])
chk("anonymous_gate_get_redirects_to_login_without_content", s == 303 and "expired" in h.get("location", "") and "일정" not in body, (s, h.get("location")))

# ---- headers
s, h, body = a.get("/home/gates/" + G["pending_decide_1"])[0:3]
csp = h.get("content-security-policy", "")
chk("security_headers_csp_nosniff_referrer_noframe", "default-src 'none'" in csp and "frame-ancestors 'none'" in csp and h.get("x-content-type-options") == "nosniff" and h.get("referrer-policy") == "no-referrer", h)

# ---- O-NM-08 on friends page (C->A request)
s, h, body = a.get("/home/friends")
m = main(body)
btns = re.findall(r"<button[^>]*>.*?</button>|<p class=\"muted deny-help\"[^>]*>.*?</p>", m, re.S)
accept = m.find('value="accept"'); helpp = m.find('class="muted deny-help"'); deny = m.find('value="deny"')
chk("onm08_order_accept_help_deny", 0 < accept < helpp < deny, (accept, helpp, deny))
chk("onm08_deny_describedby_matches_help_id", re.search(r'id="(deny-help-[^"]+)"', m) and re.search(r'value="deny" aria-label="[^"]+" aria-describedby="%s"' % re.search(r'id="(deny-help-[^"]+)"', m).group(1), m), "")
chk("onm08_help_text_says_refusal_is_recorded_and_ignoring_is_ok", "거절하면 요청한 사람의 기록에 거절로 표시됩니다. 알리지 않으려면 응답하지 않아도 기한이 지나면 요청이 끝납니다." in m, "")
chk("onm08_buttons_have_target_in_accessible_name", 'aria-label="친구 요청 수락: 합성 회원 C /' in m and 'aria-label="친구 요청 거절: 합성 회원 C /' in m, "")

# ---- GET has no side effect
after = text(a.get("/home/legacy")[2])
chk("gets_do_not_change_state_list_identical", re.sub(r"\d\d:\d\d:\d\d KST 이후", "", snap_before) == re.sub(r"\d\d:\d\d:\d\d KST 이후", "", after) or snap_before.split("세션")[0] == after.split("세션")[0], "")
chk("pending_gates_still_pending_after_all_gets", all("… 대기" in text(section(gate(a, k)[3], "state-h")) for k in ["pending_query_hint", "pending_commit_long", "pending_decide_1", "pending_decide_2"]), "")

json.dump({"origin": "http://127.0.0.1:18931", "checks": results}, open("http-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
f = [r for r in results if r["result"] == "FAIL"]
print(f"{len(results)-len(f)} PASS / {len(f)} FAIL")
sys.exit(1 if f else 0)
