"""Re-request equality (V-CN-01/F1/O-CT-01): A re-sends the same code on three paths (blocked by D, receive limit at E, unknown cancelled).
Run once on a fresh fixture: python -I -X utf8 rerequest.py -> rerequest-result.json. Synthetic accounts only; POST /home/friends/request and /home/friends/decision.
Compared: A's sent-row count/order/state/tail/deadline format per path; D, E (and control B, C) pages byte-equal before and after every step."""
import json, re, sys
sys.path.insert(0, ".")
from client import as_, main, text

M = json.load(open("C:/Users/shin/AppData/Local/KnowsLinkDevTools/connect-ux002-dev/fixture-manifest.json"))
CODES = M["rerequest_codes"]  # key -> full code the requester types
TAIL = {k: v[-4:] for k, v in CODES.items()}
results, obs = [], {}

def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": detail if not ok else ""})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:500])

ROW = re.compile(r"<(?:li|article|section)[^>]*>(?:(?!</(?:li|article|section)>).)*?(?:보낸 친구 요청|친구 요청 기록)(?:(?!</(?:li|article|section)>).)*?</(?:li|article|section)>", re.S)
def sent_rows(html):
    return [r for r in ROW.findall(main(html)) if "합성 회원 C의 친구 요청" not in r]
def norm(row):  # drop only the tail and the clock digits; state words, labels and layout stay
    t = text(row)
    for k, v in TAIL.items(): t = t.replace("…" + v, "…TAIL")
    return re.sub(r"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}", "DEADLINE", t)
def row_for(html, key):
    hit = [r for r in sent_rows(html) if "…" + TAIL[key] in text(r)]
    return hit[0] if len(hit) == 1 else None
def snap(who):
    s, h, d = as_(who).get("/home/friends"); assert s == 200, (who, s); return main(d)

a = as_("a")
page = lambda: a.get("/home/friends")[2]
others0 = {w: snap(w) for w in "bcde"}
rows0 = sent_rows(page())
chk("start_initial_state_a_has_7_sent_rows_3_ended", len(rows0) == 7 and sum("보낸 친구 요청" not in text(r) for r in rows0) == 3, [text(r)[:80] for r in rows0])
chk("start_each_code_row_is_ended_not_pending", all(row_for(page(), k) and "보낸 친구 요청" not in text(row_for(page(), k)) for k in CODES), [(k, text(row_for(page(), k) or "")[:80]) for k in CODES])

def send(code):
    s, h, d = a.post("/home/friends/request", {"code": code, "consent": "yes"}, patient=True)
    return {"status": s, "location": h.get("location"), "text": norm(re.search(r"<main.*</main>", d, re.S).group(0)) if s == 200 else ""}
def stable(step):
    for w, before in others0.items():
        now = snap(w)
        chk(f"{step}_{w}_page_byte_equal_to_start" + ("_(control)" if w in "bc" else "_(inference check)"), now == before, text(now)[:200])

# ---- P1: re-send each code once
resp1 = {k: send(c) for k, c in CODES.items()}
obs["post_responses_p1"] = resp1
chk("p1_post_response_identical_on_three_paths", len({json.dumps(v, sort_keys=True) for v in resp1.values()}) == 1, resp1)
pg = page(); rows1 = sent_rows(pg)
chk("p1_a_row_count_unchanged_7", len(rows1) == 7, len(rows1))
shapes1 = {k: norm(row_for(pg, k)) if row_for(pg, k) else None for k in CODES}
obs["p1_shapes"] = shapes1
chk("p1_three_rows_each_pending_with_tail_same_shape", None not in shapes1.values() and len(set(shapes1.values())) == 1 and "요청 대기" in shapes1["blocked"] and "입력한 코드 …TAIL" in shapes1["blocked"], shapes1)
pos = {k: [i for i, r in enumerate(rows1) if "…" + TAIL[k] in text(r)] for k in CODES}
obs["p1_positions"] = pos
chk("p1_each_code_exactly_one_row", all(len(v) == 1 for v in pos.values()), pos)
btn = lambda html, k: re.findall(r'aria-label="(입력한 코드 …%s 친구 요청 취소, 기한 [^"]+)"' % TAIL[k], html)
chk("p1_cancel_button_name_has_tail_and_deadline_same_format", all(len(btn(pg, k)) == 1 for k in CODES) and len({re.sub(r"…\w{4}", "…T", re.sub(r"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}", "D", btn(pg, k)[0])) for k in CODES}) == 1, {k: btn(pg, k) for k in CODES})
dead = {k: re.findall(r"요청 기한: (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} KST)", text(row_for(pg, k))) for k in CODES}
chk("p1_deadline_24h_format_same_three_paths", all(len(v) == 1 for v in dead.values()) and len({x[0][:10] for x in dead.values()}) <= 2, dead)
stable("p1")

# ---- P2: repeat while pending is a no-op, same on every path
pg_before = page()
resp2 = {k: send(c) for k, c in CODES.items()}
chk("p2_repeat_response_identical_three_paths", len({json.dumps(v, sort_keys=True) for v in resp2.values()}) == 1, resp2)
chk("p2_repeat_while_pending_changes_nothing_in_a", main(page()) == main(pg_before), "")
stable("p2")

# ---- P3: cancel each; ended record same on three paths
def cancel(k):
    pg = page()
    for f in re.findall(r'<form method="post" action="/home/friends/decision">.*?</form>', pg, re.S):
        if "…" + TAIL[k] in f and 'value="cancel"' in f:
            fid = re.search(r'name="friend" value="([^"]+)"', f).group(1); gen = re.search(r'name="generation" value="([^"]+)"', f).group(1)
            s, h, d = a.post("/home/friends/decision", {"friend": fid, "generation": gen, "decision": "cancel"}, patient=True)
            return {"status": s, "location": h.get("location"), "generation": gen}
    return None
c3 = {k: cancel(k) for k in CODES}
obs["cancel_responses"] = c3
chk("p3_cancel_found_and_response_same_three_paths", None not in c3.values() and len({(v["status"], v["location"]) for v in c3.values()}) == 1, c3)
chk("p3_cancel_generation_hidden_value_same_three_paths", len({v["generation"] for v in c3.values() if v}) == 1, {k: v["generation"] for k, v in c3.items() if v})
pg = page(); rows3 = sent_rows(pg)
shapes3 = {k: norm(row_for(pg, k)) if row_for(pg, k) else None for k in CODES}
obs["p3_shapes"] = shapes3
chk("p3_row_count_still_7_and_ended_shape_same", len(rows3) == 7 and None not in shapes3.values() and len(set(shapes3.values())) == 1 and "취소" in shapes3["blocked"] and "보낸 친구 요청" not in shapes3["blocked"], shapes3)
chk("p3_row_order_same_as_p1", {k: [i for i, r in enumerate(rows3) if "…" + TAIL[k] in text(r)] for k in CODES} == pos, "")
stable("p3")

# ---- P4: re-request after cancel
resp4 = {k: send(c) for k, c in CODES.items()}
chk("p4_post_response_identical_three_paths", len({json.dumps(v, sort_keys=True) for v in resp4.values()}) == 1, resp4)
pg = page(); rows4 = sent_rows(pg)
shapes4 = {k: norm(row_for(pg, k)) if row_for(pg, k) else None for k in CODES}
chk("p4_after_cancel_rerequest_7_rows_pending_same_shape", len(rows4) == 7 and None not in shapes4.values() and len(set(shapes4.values())) == 1 and "요청 대기" in shapes4["blocked"], shapes4)
chk("p4_row_order_same_as_p1", {k: [i for i, r in enumerate(rows4) if "…" + TAIL[k] in text(r)] for k in CODES} == pos, "")
chk("p4_p1_shape_equals_p4_shape", shapes1 == shapes4, (shapes1, shapes4))
stable("p4")

# ---- D/E inference: what D and E can still do/see (read-only)
d_now, e_now = snap("d"), snap("e")
chk("d_still_shows_block_record_no_request_card", "— 차단" in text(d_now) and re.findall(r'name="decision" value="(\w+)"', d_now) == ["unblock"] and "합성 회원 A" not in text(d_now), text(d_now)[:300])
chk("e_still_ten_pending_plus_old_ended_no_new_from_A", text(e_now).count("친구 요청 서로의") == 10 and "합성 회원 A" not in text(e_now), text(e_now).count("친구 요청 서로의"))

json.dump({"origin": M["origin"] if False else "http://127.0.0.1:18932", "results": results, "obs": obs}, open("rerequest-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
fails = [r for r in results if r["result"] == "FAIL"]
print(f"\n{len(results) - len(fails)} PASS / {len(fails)} FAIL"); sys.exit(1 if fails else 0)
