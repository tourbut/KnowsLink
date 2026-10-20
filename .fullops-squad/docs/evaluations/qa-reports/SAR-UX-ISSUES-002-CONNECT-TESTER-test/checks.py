"""Read-only HTTP checks for the connect-guide / resume fixture (U8, U7, O-NM-03, low 3). Synthetic A/B/C only; GET plus /fixture/as.
Run: python -I -X utf8 checks.py   -> http-result.json, exit 1 on any FAIL.
"""
import json, re, sys
sys.path.insert(0, ".")
from client import as_, main, text

M = json.load(open("C:/Users/shin/AppData/Local/KnowsLinkDevTools/connect-ux002-dev/fixture-manifest.json"))
C, AG = M["connections"], M["agents"]
results = []


def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": detail if not ok else ""})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:400])


def card(m, agent):  # one <li> of ul.agents for an agent id
    for li in re.findall(r"<li[^>]*>.*?</li>", m, re.S):
        if f"식별자: {agent}" in text(li) or f"#agent-{agent}" in li or f'id="agent-{agent}"' in li:
            return li
    return ""


def cards(m):
    ul = re.search(r'<ul class="agents">(.*?)</ul>\s*(?:<p|<section|<details|<form|<h)', m, re.S)
    return re.findall(r'<li id="agent-[^"]+"', m)


USES = ["web", "grok", "dots", "claude", "other"]
CHOICES = [("웹에서 바로 사용", "설치 필요 없음"), ("Grok Bot (내 컴퓨터)", "연결 안내 있음 · 기존 계정 사용 확인"),
           ("Dots (원격)", "연결 안내 있음 · 기존 계정 사용 확인"),
           ("Claude Code 작업 세션 (내 컴퓨터)", "연결 안내 있음 · 시험 환경 확인 (Claude Code 2.1.302, Windows)"),
           ("지원 확인이 필요한 호스트", "Codex, Hermes Agent, 그 밖의 호스트")]
OTHER_TEXT = "아직 연결 방법을 확인하는 중입니다. 지금은 웹에서 대화를 이용하세요."

# ---------------- A: agents page, selector, guides
a = as_("a")
s, h, base = a.get("/home/agents"); mb = main(base)
chk("agents_a_200_title_h1", s == 200 and re.search(r"<title>[^<]+</title>", base) and len(re.findall(r"<h1[ >]", base)) == 1, (s,))
form = re.search(r'<form method="get" action="/home/agents#use-guide">.*?</form>', mb, re.S)
form = form.group(0) if form else ""
labels = re.findall(r'<label class="choice"><input type="radio" name="use" value="(\w+)"[^>]*><span>(.*?)<br><span class="muted">(.*?)</span></span></label>', form)
chk("choices_five_in_design_order_and_text", [(t, sub) for _, t, sub in labels] == CHOICES and [v for v, _, _ in labels] == USES, labels)
chk("choices_start_first_continue", re.findall(r'name="start" value="(\w+)"', form) == ["first", "continue"] and "처음 연결" in form and "기존 연결 계속" in form, form[-400:])
chk("no_choice_preselected_without_query", " checked" not in re.search(r'<legend>사용할 곳</legend>.*?</fieldset>', form, re.S).group(0), "")
chk("no_guide_panel_without_choice", 'id="use-guide"' not in mb, "")
chk("selector_is_get_form_only", 'method="get"' in form and "<form method=\"post\"" not in form, "")
n_base = len(cards(mb))
guides = {}
for use, st in [(u, "first") for u in USES] + [("claude", "continue"), ("web", "continue")]:
    s, h, d = a.get(f"/home/agents?use={use}&start={st}"); m = main(d)
    g = re.search(r'<div id="use-guide".*?</div>(?=\s*</section>|\s*<section|$)', m, re.S)
    guides[(use, st)] = (s, g.group(0) if g else "", m, d)
chk("guide_each_status_200_and_selected_checked", all(v[0] == 200 and re.search(r'value="%s"[^>]*checked' % u, v[2]) and re.search(r'name="start" value="%s" checked' % st, v[2]) for (u, st), v in guides.items()), {k: v[0] for k, v in guides.items()})
chk("selection_creates_no_agent_cards_unchanged", all(len(cards(v[2])) == n_base for v in guides.values()), (n_base, {k: len(cards(v[2])) for k, v in guides.items()}))
G = {u: text(guides[(u, "first")][1]) for u in USES}
chk("guide_titles_match_choice", all(G[u].startswith(CHOICES[i][0]) for i, u in enumerate(USES)), {u: G[u][:40] for u in USES})
chk("guide_web_links_friends_and_new", 'href="/home/friends"' in guides[("web", "first")][1] and 'href="/home/new"' in guides[("web", "first")][1] and "설치 필요 없음" in G["web"], G["web"])
chk("guide_grok_readme_link_and_prereqs", "README.md#grok-신규-설치" in guides[("grok", "first")][1] and "준비물" in G["grok"] and "권한" in G["grok"], G["grok"])
chk("guide_dots_readme_link_and_start_in_host", "README.md#dots-신규-연결" in guides[("dots", "first")][1] and "Dots에서 연결을 시작" in G["dots"], G["dots"])
gc = guides[("claude", "first")][1]
chk("guide_claude_adapters_link_and_session_on", "adapters/README.md#claude-code-작업-세션-연결--sar-session-peer-001-dev" in gc and "<code>/knowslink-session on</code>" in gc and "시험 환경은 Claude Code 2.1.302, Windows" in G["claude"], G["claude"])
chk("guide_claude_session_unobserved_says_app_check", "앱에서 확인 필요" in G["claude"] and "/knowslink-session status" in G["claude"] and "knowslink_text_receive" in G["claude"] and "연결되었습니다" not in G["claude"], G["claude"])
chk("guide_claude_continue_has_no_new_agent_instruction", "새 에이전트나 새 설치 폴더를 만들지 않습니다" in text(guides[("claude", "continue")][1]) and 'href="#resume"' in guides[("claude", "continue")][1], text(guides[("claude", "continue")][1]))
go = guides[("other", "first")][1]
chk("guide_other_exact_text_and_web_link", OTHER_TEXT in G["other"] and 'href="/home/agents?use=web#use-guide"' in go, G["other"])
chk("guide_other_no_install_button_command_badge_tool_name", not re.search(r"<(button|code|pre|form|input)\b", go) and not re.search(r"npm|node |curl|install|설치|class=\"(badge|pill|ok|tag)", go) and not re.search(r"Codex|Hermes|Claude|Grok|Dots", G["other"]) and "완료" not in G["other"] and "지원됨" not in G["other"], go)
allm = {k: v[2] for k, v in guides.items()}; allm["base"] = mb
chk("hidden_hosts_openclaw_meta_muse_absent_all_variants", not any(re.search(r"OpenClaw|Open Claw|Meta Muse|메타 뮤즈|Muse", t, re.I) for t in allm.values()), "")
# invalid/forged query values are ignored, not echoed, no 500
s, h, d = a.get("/home/agents?use=openclaw%3Cb%3E&start=zzz"); m = main(d)
chk("invalid_use_start_ignored_not_reflected", s == 200 and "openclaw" not in m.lower() and "zzz" not in m and 'id="use-guide"' not in m and len(cards(m)) == n_base, (s, text(m)[:80]))
s, h, d = a.get("/home/agents?use=other&start=first&use=claude"); m = main(d)
chk("duplicate_use_param_no_crash", s == 200 and 'id="use-guide"' in m, s)

# ---------------- 온라인 badge / claims
def online_claims(m):
    out = []
    for sent in re.split(r"(?<=[.。])\s|<[^>]+>", m):
        t = text(sent)
        if "온라인" in t and not re.search(r"확인하지 않|관측하지 않|않습니다", t): out.append(t)
    return out
chk("no_online_badge_or_claim_agents_a", not online_claims(mb) and not re.search(r'class="[^"]*(online|presence|status-dot)', mb), online_claims(mb))
chk("online_word_only_in_negated_sentence", all("확인하지 않" in x for x in re.findall(r"[^.]*온라인[^.]*\.", text(mb))), re.findall(r"[^.]*온라인[^.]*\.", text(mb)))

# ---------------- resume rows (A)
ra = {k: text(card(mb, AG[k])) for k in ["waiting", "prepared", "approved", "registered"]}
chk("resume_local_waiting", "연결 상태: 키 등록 전 · 내 컴퓨터의 준비 대기" in ra["waiting"] and "연결 수단은 다시 표시하지 않습니다" in ra["waiting"] and "연결 상태 확인" in ra["waiting"] and "연결 필요" in ra["waiting"], ra["waiting"][:300])
chk("resume_local_prepared", "연결 상태: 키 등록 전 · 지문 확인 필요" in ra["prepared"] and "다르면 취소하세요" in ra["prepared"] and "지문 확인" in ra["prepared"], ra["prepared"][:300])
chk("resume_local_approved", "연결 상태: 브라우저 동의 완료 · 클라이언트 저장 확인 필요" in ra["approved"] and "승인만으로는 연결되지 않습니다" in ra["approved"], ra["approved"][:300])
chk("resume_registered_key_only_claim", "연결 상태: 키 등록됨" in ra["registered"] and "앱에서 연결 상태를 확인하세요" in ra["registered"] and "기존 키를 덮어쓰지 말고" in ra["registered"] and "✓ 등록됨" in ra["registered"], ra["registered"][:300])
chk("resume_registered_never_claims_connected_online_working", not re.search(r"연결되었습니다|연결 완료|온라인입니다|수신 중|정상 작동", ra["registered"] + ra["approved"]), "")
chk("resume_each_state_has_one_next_action_link_set", all("연결·키 관리" in ra["registered"] or "연결 상태 확인" in ra[k] or "지문 확인" in ra[k] for k in ra), "")
_fa = as_("a").get("/home/friends")[2]
_pend_a = len([r for r in re.findall(r"<(?:li|article|section)[^>]*>(?:(?!</(?:li|article|section)>).)*?보낸 친구 요청(?:(?!</(?:li|article|section)>).)*?</(?:li|article|section)>", main(_fa), re.S) if "요청 대기" in text(r) and "합성 회원 C의 친구 요청" not in r])
chk("resume_summary_counts_pending_sent_requests", f"보낸 친구 요청 {_pend_a}건이 상대 확인을 기다립니다" in text(mb) and 'href="/home/friends#requests-h"' in mb, (_pend_a, text(mb)[:400]))
chk("resume_summary_explains_unobserved_items", "앱 설치·도구 검색·세션 연결·온라인 여부는 이 화면에서 확인하지 않습니다" in text(mb), "")
chk("resume_cards_do_not_expose_connection_token_commands_or_keys", not re.search(r"connect\.js|complete <|private[_ ]?key|BEGIN [A-Z ]*KEY|code_verifier|access_token", mb), "")
# new agent slot form still present (one slot) and idempotency attempt field
chk("create_form_has_attempt_field", bool(re.search(r'<form method="post" action="/home/agents">.*?name="attempt"', mb, re.S)), "")

# ---------------- connection pages
def conn(b, key):
    s, h, d = b.get("/home/connections/" + C[key]); return s, h, d, main(d)
step = {}
for k, want in [("waiting", 1), ("prepared", 2), ("approved", 3), ("consumed", 3)]:
    s, h, d, m = conn(a, k)
    cur = re.findall(r'<li aria-current="step"><span class="n">(\d)</span>([^<]*)</li>', m)
    step[k] = (s, cur, text(m), d, m)
    chk(f"connection_{k}_200_exactly_one_current_step_{want}", s == 200 and len(cur) == 1 and int(cur[0][0]) == want and m.count("aria-current") == 1, (s, cur))
    chk(f"connection_{k}_landmarks_title_h1_skip", re.search(r"<title>에이전트 연결 · KnowsLink</title>", d) and len(re.findall(r"<h1[ >]", d)) == 1 and "<main" in d and 'class="skip" href="#main"' in d and 'id="main"' in d, "")
    chk(f"connection_{k}_return_link_exact", f'href="/home/agents#agent-{AG["registered" if k == "consumed" else k]}"' in m, re.findall(r'href="([^"]+)"', m))
chk("connection_consumed_label_key_registered_not_complete", "3 키 등록됨" in step["consumed"][2] and "키 등록됨. 앱에서 연결 상태를 확인하세요" in step["consumed"][2] and "이 화면은 온라인 여부를 확인하지 않습니다" in step["consumed"][2] and "연결 취소" not in step["consumed"][2], step["consumed"][2][:300])
chk("connection_waiting_not_connected_and_refresh", "아직 연결되지 않았습니다" in step["waiting"][2] and "상태 새로고침" in step["waiting"][2] and "연결 취소" in step["waiting"][2], step["waiting"][2][:300])
chk("connection_approved_not_connected_yet_instruction", "승인했습니다. 아직 새 키가 연결되지 않았습니다" in step["approved"][2] and "complete" in step["approved"][2], "")
chk("connection_prepared_fingerprint_compare_and_cancel", "SHA256:" in step["prepared"][2] and "다르면 취소하세요" in step["prepared"][2] and "연결 취소" in step["prepared"][2], "")
# cross-account: b cannot open a's connection, anonymous cannot
bb = as_("b")
s, h, d = bb.get("/home/connections/" + C["waiting"])
chk("connection_of_other_member_not_visible", s in (303, 403, 404) or ("agent_fx_waiting" not in d and "합성 대기 도우미" not in d), (s, text(d)[:100]))

# ---------------- B: expired / cancelled / device_lost / remote
s, h, bagents = bb.get("/home/agents"); bm = main(bagents)
rb = {k: text(card(bm, AG[k])) for k in ["expired", "cancelled", "device_lost", "remote_returning", "remote_lost", "remote_exchanged"]}
chk("resume_expired", "이전 시도 만료 · 키 없음" in rb["expired"] and "이전 연결 수단은 다시 쓸 수 없습니다" in rb["expired"], rb["expired"][:300])
chk("resume_cancelled", "이전 시도 취소 · 키 없음" in rb["cancelled"] and "이전 연결 수단은 다시 쓸 수 없습니다" in rb["cancelled"], rb["cancelled"][:300])
chk("resume_device_lost_unknown_not_blind_retry", "확인 필요 · 승인 뒤 클라이언트 저장을 확인하지 못함" in rb["device_lost"] and "자동 철회는 없습니다" in rb["device_lost"] and "새로" not in rb["device_lost"], rb["device_lost"][:300])
chk("resume_remote_returning", "동의 완료 · 앱 복귀 확인 필요" in rb["remote_returning"] and "앱으로 이동한 것만으로 연결되지 않습니다" in rb["remote_returning"], rb["remote_returning"][:300])
chk("resume_remote_lost_unknown", "확인 필요 · 앱의 토큰 교환을 확인하지 못함" in rb["remote_lost"] and "자동으로 철회되지 않습니다" in rb["remote_lost"] and "앱에서 연결을 다시 시작" in rb["remote_lost"], rb["remote_lost"][:300])
chk("resume_remote_exchanged_key_registered_only", "키 등록됨 · 원격 앱의 토큰 교환 확인" in rb["remote_exchanged"] and "이 화면은 온라인 여부를 확인하지 않습니다" in rb["remote_exchanged"], rb["remote_exchanged"][:300])
chk("resume_b_summary_no_sent_pending_line", "보낸 친구 요청" not in text(bm) and not online_claims(bm), "")
expc = {}
for k in ["expired", "cancelled", "device_lost"]:
    s, h, d, m = conn(bb, k); expc[k] = (s, text(m), m, d)
    chk(f"connection_{k}_terminal_no_step_no_action_return_link", s == 200 and "aria-current" not in m and not re.search(r"<button[^>]*>\s*(연결 취소|연결 승인)", m) and f'href="/home/agents#agent-{AG[k]}"' in m, (s, text(m)[:160]))
chk("connection_expired_cancelled_text", "만료되었습니다" in expc["expired"][1] and "취소되었습니다" in expc["cancelled"][1] and all("이전 연결 수단은 다시 쓸 수 없습니다" in expc[k][1] for k in ("expired", "cancelled")), "")
# card vs connection page agreement for the approved-then-lost device attempt (server state: approved, never consumed)
dl_card, dl_page = rb["device_lost"], expc["device_lost"][1]
chk("device_lost_connection_page_agrees_with_card_unknown_state", "확인 필요" in dl_page and "새로 시작" not in dl_page, {"card": dl_card[:200], "page": dl_page[:260]})
chk("device_lost_page_same_unknown_wording_as_card_and_no_expired_or_not_connected_claim", "확인 필요 · 승인 뒤 클라이언트 저장을 확인하지 못함" in dl_page and "시작한 클라이언트의 저장 결과를 먼저 확인하세요" in dl_page and "자동 철회는 없습니다" in dl_page and not re.search(r"만료되었습니다|취소되었습니다|새 키가 연결되지 않았습니다|새로 시작", dl_page) and "aria-current" not in expc["device_lost"][2], dl_page[:300])
chk("device_lost_page_no_blind_new_start_button", not re.search(r"<button[^>]*>\s*(새로 시작|새 키|연결 시작|연결 승인)", expc["device_lost"][2]), "")
chk("device_lost_expired_and_cancelled_pages_keep_their_own_terminal_wording", "새 키가 연결되지 않았습니다" in expc["expired"][1] and "새 키가 연결되지 않았습니다" in expc["cancelled"][1], "")
for who, nm in [(bb, "unknown_id"), (a, "other_members_id")]:  # V-CN-02: rejected connection id returns to 내 에이전트
    s_, h_, d_ = who.get("/home/connections/" + ("x" * 43 if nm == "unknown_id" else C["device_lost"])); m_ = main(d_)
    chk(f"rejected_connection_{nm}_403_returns_to_my_agents", s_ == 403 and "내 에이전트로 돌아가기" in text(m_) and 'href="/home/agents"' in m_ and "홈으로" not in text(m_) and "agent_fx" not in d_, (s_, text(m_)[:200], re.findall(r'href="([^"]+)"', m_)))

# ---------------- C: zero agents
cc = as_("c")
s, h, cagents = cc.get("/home/agents"); cm = main(cagents)
chk("agents_c_zero_agents_first_agent_prompt_and_selector", "첫 에이전트를 추가하세요" in text(cm) and len(cards(cm)) == 0 and len(re.findall('name="use"', cm)) == 5, text(cm)[:200])
chk("agents_c_summary_sent_request_line", "보낸 친구 요청 1건이 상대 확인을 기다립니다" in text(cm), "")

# ---------------- O-NM-03 (friends pages)
s, h, fa = a.get("/home/friends"); fam = main(fa)
RC = M["rerequest_codes"]
b_code = re.search(r"fixture-code-b-\w+", bb.get("/home/friends")[2]).group(0); b_tail = b_code[-4:]
a_code = re.search(r"fixture-code-a-[\w-]+", fa).group(0); a_tail = a_code[-4:]
sent_rows = re.findall(r"<(?:li|article|section)[^>]*>(?:(?!</(?:li|article|section)>).)*?(?:보낸 친구 요청|친구 요청 기록)(?:(?!</(?:li|article|section)>).)*?</(?:li|article|section)>", fam, re.S)
RE_TAIL = re.compile(r"입력한 코드 (…[A-Za-z0-9_-]{4}|형식 다름)|이전 요청")
trows = [text(r) for r in sent_rows]
chk("onm03_sent_rows_found_at_least_7_original", len(trows) >= 7, trows)  # cap.py cancel records may add ended rows
chk("onm03_each_sent_row_has_tail_or_format_or_old_marker", all(RE_TAIL.search(r) for r in trows), trows)
chk("onm03_expected_tails_present", all(any(x in r for r in trows) for x in ["입력한 코드 …" + b_tail, "입력한 코드 …zz99", "입력한 코드 형식 다름", "이전 요청"] + ["입력한 코드 …" + c[-4:] for c in RC.values()]), trows)
shape = lambda r: re.sub(r"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}", "DEADLINE", re.sub(r"입력한 코드 …[A-Za-z0-9_-]{4}", "TAIL", r))
pend = [shape(r).removeprefix("친구 요청과 기록 ") for r in trows if "요청 대기" in r and "TAIL" in shape(r)]  # first row carries its section heading
chk("onm03_real_and_nonexistent_code_rows_identical_shape", len(pend) >= 2 and len(set(pend)) == 1, pend)
chk("onm03_case_kept_exact_tail_not_lowercased", "…" + b_tail in fam and "…" + b_tail.lower() not in fam, b_tail)
chk("onm03_short_input_ab_shows_format_unknown_not_raw", "입력한 코드 형식 다름" in fam and not re.search(r"입력한 코드[^<]*\bab\b", fam), "")
chk("onm03_cancel_button_accessible_name_has_tail", all(re.search(r'aria-label="입력한 코드 …%s 친구 요청 취소[^"]*"' % t, fam) for t in [b_tail, "zz99"]), re.findall(r'aria-label="[^"]*취소[^"]*"', fam))
_ended = [r for r in trows if "보낸 친구 요청" not in r]
chk("onm03_ended_row_keeps_tail_for_requester", all("입력한 코드 …" in r and "새 친구 요청과 수락이 필요합니다" in r for r in _ended) and (not _ended or any("취소" in r or "만료" in r for r in _ended)), (len(_ended), _ended))
chk("onm03_full_typed_code_not_in_requester_page_or_links", b_code not in fa and b_code[:-4] not in fa, "")
recv = text(re.search(r"합성 회원 C의 친구 요청.*?거절", fam, re.S).group(0)) if "합성 회원 C의 친구 요청" in fam else ""
chk("onm03_received_request_card_has_no_tail", recv != "" and "입력한 코드" not in recv and not re.search(r"…\w{4}", recv), recv)
s, h, fb = bb.get("/home/friends"); fbm = main(fb)
recvb = text(re.search(r"합성 회원 A의 친구 요청.*?거절", fbm, re.S).group(0)) if "합성 회원 A의 친구 요청" in fbm else ""
chk("onm03_recipient_b_sees_no_sender_input_tail", recvb != "" and "입력한 코드" not in fbm and "zz99" not in fb and "형식 다름" not in fb and "이전 요청" not in fb and "…" not in recvb, recvb)
chk("onm03_recipient_b_no_existence_signal_beyond_own_request", b_tail in fb and fb.count(b_tail) == fb.count(b_code), fb.count(b_tail))
s, h, fc = cc.get("/home/friends"); fcm = main(fc)
chk("onm03_c_sent_row_tail_of_as_code_and_no_full_code", "입력한 코드 …" + a_tail in text(fcm) and a_code not in fc, text(fcm)[:200])
chk("onm03_no_tail_in_urls_forms_hidden_fields", not re.search(r'(href|action|value)="[^"]*(%s|zz99|ca11)' % b_tail, fa), "")

# ---------------- U7 secondary pages
s, h, lg = a.get("/home/legacy")
chk("legacy_find_on_this_screen_visible_heading_no_false_current", s == 200 and re.search(r'<nav class="toc" aria-labelledby="(\w[\w-]*)"><span id="\1"[^>]*>이 화면에서 찾기</span>', lg) is not None and "aria-current" not in main(lg), (s, re.findall(r'<nav class="toc"[^>]*>.{0,80}', lg, re.S), re.findall(r'aria-current="[^"]+"', lg)))
chk("legacy_toc_heading_is_visible_text_not_hidden", "이 화면에서 찾기" in text(lg) and not re.search(r'class="(sr-only|visually-hidden)"[^>]*>이 화면에서 찾기', lg), "")
chk("legacy_title_h1_main_skip", re.search(r"<title>[^<]+</title>", lg) and len(re.findall(r"<h1[ >]", lg)) == 1 and "<main" in lg and 'class="skip"' in lg, "")
for name, body in [("agents", base), ("friends", fa), ("connection_waiting", step["waiting"][3])]:
    chk(f"page_{name}_unique_title_h1_main_skip_lang", re.search(r"<title>[^<]+</title>", body) and len(re.findall(r"<h1[ >]", body)) == 1 and "<main" in body and 'class="skip"' in body and 'lang="ko"' in body, "")
titles = [re.search(r"<title>([^<]+)</title>", x).group(1) for x in (base, fa, lg, step["waiting"][3])]
chk("page_titles_distinct", len(set(titles)) == len(titles), titles)
cur_page = lambda body: re.findall(r'aria-current="page"[^>]*>([^<]*)', body)
chk("agents_nav_marks_only_my_agents_current_page", cur_page(base) == ["내 에이전트"], cur_page(base))
chk("friends_nav_marks_only_friends_current_page", cur_page(fa) == ["친구"], cur_page(fa))
chk("connection_page_marks_no_page_current", not cur_page(step["waiting"][3]), cur_page(step["waiting"][3]))

# ---------------- anonymous
from client import Browser
an = Browser(); s, h, d = an.get("/home/agents")
chk("anonymous_agents_redirects_to_login_no_content", s in (302, 303) and "합성" not in d, (s, h.get("location")))
s, h, d = an.get("/home/connections/" + C["waiting"])
chk("anonymous_connection_redirects_no_content", s in (302, 303) and "agent_fx" not in d, (s, h.get("location")))

# ---------------- state unchanged by all reads
s, h, after = a.get("/home/agents")
chk("reads_left_a_agent_count_and_rows_unchanged", len(cards(main(after))) == n_base and text(main(after)) == text(mb), (n_base, len(cards(main(after)))))
json.dump({"origin": "http://127.0.0.1:18932", "results": results, "waited429": __import__("client").Browser.waited429}, open("http-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
fails = [r for r in results if r["result"] == "FAIL"]
print(f"\n{len(results) - len(fails)} PASS / {len(fails)} FAIL")
sys.exit(1 if fails else 0)
